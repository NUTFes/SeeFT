"""生成済みスライドHTMLに画像を base64 で埋め込む（LLM・AIアカウント不要）。

解説マニュアルの生成は、非決定的な工程と決定的な工程に分かれている。

    pandoc で Markdown 化      決定的
    LLM で HTML 生成           非決定的  ← ここだけAIアカウントが要る
    画像の base64 埋め込み      決定的    ← 本スクリプト
    完全性の検査               決定的    ← 本スクリプト

LLM工程は任意のAIセッション（Claude Code / Codex / ブラウザのチャット等）で代替できる。
プロンプトは .claude/manual-prompt-card-strict.md にあり、そこで得た応答を
manual_dir に保存してから本スクリプトを通せば、どのサブスクを使っても同じ成果物になる。

画像の埋め込みをLLMに任せられないのは、生成HTMLが数MBになりその大半が base64 文字列で、
LLMの出力長に収まらないため。ここは決定的なコードで処理する必要がある。

依存は標準ライブラリのみ。uv も claude_agent_sdk も不要で、python3 だけで動く。

使い方:
  python3 scripts/claude-slide/embed_images.py docs/manuals/45th_企画マニュアル_縁日
  python3 scripts/claude-slide/embed_images.py --file out.html docs/manuals/45th_企画マニュアル_縁日
"""

import argparse
import base64
import mimetypes
import os
import re
import sys


SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
PROJECT_ROOT = os.path.normpath(os.path.join(SCRIPT_DIR, "..", ".."))


def resolve_manual_dir(arg: str) -> str:
    arg = arg.rstrip("/")
    if os.path.isabs(arg):
        return arg
    if os.path.isdir(arg):
        return os.path.abspath(arg)
    return os.path.join(PROJECT_ROOT, arg)


def load_images_base64(manual_dir: str) -> dict[str, str]:
    """manual_dir/images/ の全画像を {ファイル名 -> data URI} で返す。"""
    images: dict[str, str] = {}
    img_dir = os.path.join(manual_dir, "images")
    if not os.path.isdir(img_dir):
        return images
    for f in sorted(os.listdir(img_dir)):
        path = os.path.join(img_dir, f)
        mime = mimetypes.guess_type(path)[0] or "image/png"
        with open(path, "rb") as fh:
            b64 = base64.b64encode(fh.read()).decode()
        images[f] = f"data:{mime};base64,{b64}"
    return images


def _resolve_image_src(src: str, images: dict[str, str]) -> str | None:
    """`<img src="...">` の src 値がローカル画像を指すなら data URI を返す。違えば None。

    プロンプトは `{{fname}}` プレースホルダーを指示しているが、LLMがそれを守らず
    素のファイル名 (`image1.jpg`) や相対パス (`images/image1.png`, `./images/image1.png`)
    で書くことがある。チャット経由の生成では特に揺れやすいため、その揺れを決定的に
    吸収して必ず base64 埋め込みする。

    注意: 既に `data:` 化済みの src は再処理しても None を返すこと（冪等性）。
    """
    if src.startswith("data:"):
        return None  # 既に埋め込み済み。再処理しても壊さない（冪等）
    return images.get(os.path.basename(src))


def replace_placeholders(html: str, images: dict[str, str]) -> str:
    """`{{fname}}` と素の `src` の両方を data URI に置き換える。"""
    # 1) {{fname}} プレースホルダー（プロンプトが指示する正規ルート）
    for fname, data_uri in images.items():
        html = html.replace(f"{{{{{fname}}}}}", data_uri)

    # 2) 素のファイル名・相対パスで書かれた <img src> も決定的に base64 化する保険
    def _embed(m: "re.Match[str]") -> str:
        quote = m.group("q")
        data_uri = _resolve_image_src(m.group("src"), images)
        if data_uri is None:
            return m.group(0)
        return f"src={quote}{data_uri}{quote}"

    return re.sub(r'src=(?P<q>["\'])(?P<src>[^"\']*)(?P=q)', _embed, html)


def referenced_images(html: str, images: dict[str, str]) -> set[str]:
    """HTMLが参照している画像のファイル名の集合を返す（埋め込み前に数えるため）。

    images/ にあるのに参照されていない画像は、本文から抜け落ちた可能性がある。
    チャット経由の生成ではLLMが画像を落とすことがあるため、件数で検出できるようにする。
    """
    found: set[str] = set()
    for fname in images:
        if f"{{{{{fname}}}}}" in html:
            found.add(fname)
    for m in re.finditer(r'src=(["\'])(?P<src>[^"\']*)\1', html):
        src = m.group("src")
        if src.startswith("data:"):
            continue
        base = os.path.basename(src)
        if base in images:
            found.add(base)
    return found


def already_embedded_images(html: str, images: dict[str, str]) -> set[str]:
    """既に data URI として埋め込まれている画像のファイル名を返す。

    埋め込み後のHTMLは `src` からファイル名が消えるため、これを数えずに
    referenced_images だけで判定すると、再実行時に「本文から抜けている」と
    誤警告してしまう。実際には埋め込み済みなので、区別する必要がある。
    """
    return {fname for fname, uri in images.items() if uri in html}


def is_complete_html(html: str) -> bool:
    """応答が途中で切れていないかを構造で判定する。

    LLMの応答は途中で切れてもエラーにならないことがある。チャット経由の生成では
    人間が検査役になるため見落としやすく、壊れたHTMLをアップロードする事故につながる。
    """
    return html.rstrip().lower().endswith("</html>")


def find_slide_html(manual_dir: str) -> str:
    """manual_dir 内のスライドHTMLを1つに決める。曖昧なら候補を示して失敗する。"""
    candidates = sorted(
        f for f in os.listdir(manual_dir)
        if f.startswith("slide") and f.endswith(".html")
    )
    if not candidates:
        raise FileNotFoundError(
            f"スライドHTMLが見つかりません: {manual_dir}\n"
            "  AIの応答を slide_claude.card-strict.html として保存してください"
        )
    if len(candidates) > 1:
        raise RuntimeError(
            f"スライドHTMLが複数あります: {candidates}\n"
            "  --file でどれを処理するか指定してください"
        )
    return os.path.join(manual_dir, candidates[0])


def main() -> int:
    parser = argparse.ArgumentParser(
        description="生成済みスライドHTMLに画像を埋め込む（LLM不要・標準ライブラリのみ）",
    )
    parser.add_argument(
        "manual_dir",
        help="マニュアルディレクトリ（images/ を含む。絶対パス or プロジェクトルートからの相対パス）",
    )
    parser.add_argument(
        "--file",
        default=None,
        help="処理するHTMLを明示する（省略時は manual_dir 内の slide*.html を自動検出）",
    )
    args = parser.parse_args()

    manual_dir = resolve_manual_dir(args.manual_dir)
    if not os.path.isdir(manual_dir):
        print(f"  ERROR: ディレクトリがありません: {manual_dir}", file=sys.stderr)
        return 1

    try:
        html_path = args.file or find_slide_html(manual_dir)
    except (FileNotFoundError, RuntimeError) as e:
        print(f"  ERROR: {e}", file=sys.stderr)
        return 1

    if not os.path.isfile(html_path):
        print(f"  ERROR: ファイルがありません: {html_path}", file=sys.stderr)
        return 1

    with open(html_path, "r", encoding="utf-8") as f:
        html = f.read()

    print("=== 画像の埋め込み ===")
    print(f"  対象: {html_path} ({os.path.getsize(html_path)//1024}KB)")

    # 埋め込む前に完全性を見る。切れているHTMLに画像を入れても意味がないうえ、
    # 壊れたまま配信してしまうため、ここで止める
    if not is_complete_html(html):
        # 直前までの進捗表示(stdout)を吐き出してからエラー(stderr)を出す。
        # リダイレクト時に順序が入れ替わって読みにくくなるのを防ぐ
        sys.stdout.flush()
        print(
            "  ERROR: HTMLが </html> で閉じていません。応答が途中で切れています。\n"
            "         AIに続きを出力させるか、生成をやり直してください",
            file=sys.stderr,
        )
        return 1

    images = load_images_base64(manual_dir)
    if not images:
        print("  images/ が無いか空です。埋め込む画像はありません")
        return 0

    used = referenced_images(html, images)
    already = already_embedded_images(html, images)
    # 参照もされておらず埋め込み済みでもない画像だけが「本文から抜けている」候補
    missing = sorted(set(images) - used - already)

    html = replace_placeholders(html, images)

    with open(html_path, "w", encoding="utf-8") as f:
        f.write(html)

    print(f"  画像: {len(images)}枚中 {len(used)}枚を埋め込み")
    if already:
        print(f"  {len(already)}枚は既に埋め込み済み（再実行のため変更なし）")
    if missing:
        print(f"  警告: HTMLから参照されていない画像が {len(missing)}枚あります")
        for f_ in missing:
            print(f"    - {f_}")
        print("         本文から抜け落ちていないか確認してください")
    print(f"  出力: {html_path} ({os.path.getsize(html_path)//1024}KB)")
    print("=== 完了 ===")
    return 0


if __name__ == "__main__":
    sys.exit(main())
