"""Google Doc のエクスポート zip を、生成スクリプトが読める形に展開する。

Google ドキュメントの「ウェブページ (.html, zip)」でダウンロードした zip を渡すと、
`docs/manuals/{マニュアル名}/` に `source.html` と `images/` を配置する。

この工程を機械化するのは、手作業だと次の3点を毎回間違えるため。

  - zip 内の HTML は名前が切り詰められている（実例: `45th_企画マニュアル_ホールインワン.zip`
    の中身は `45th_.html`）。何のマニュアルか分からなくなる
  - 生成スクリプトは「`slide` と `verify` で始まらない .html がちょうど1つ」を要求する。
    リネームを忘れると曖昧だとして停止する
  - `images/` の位置を間違えると画像が埋め込まれない

依存は標準ライブラリのみ。uv も claude_agent_sdk も不要で、python3 だけで動く。

使い方:
  python3 scripts/automation/prepare_manual.py docs/manuals/_zips/45th_企画マニュアル_縁日.zip
  python3 scripts/automation/prepare_manual.py --name 45th_企画マニュアル_縁日 path/to/export.zip
"""

import argparse
import os
import shutil
import sys
import tempfile
import zipfile


SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
PROJECT_ROOT = os.path.normpath(os.path.join(SCRIPT_DIR, "..", ".."))
MANUALS_DIR = os.path.join(PROJECT_ROOT, "docs", "manuals")

IMAGE_EXTS = {".png", ".jpg", ".jpeg", ".gif", ".webp", ".bmp", ".svg"}


def derive_name(zip_path: str) -> str:
    """zip のファイル名からマニュアル名を導出する。

    Google の書き出しは zip 名にドキュメントのタイトルを保つ（中の HTML は切り詰められる）。
    そのためタイトルの取得元としては zip 名のほうが信頼できる。
    """
    return os.path.splitext(os.path.basename(zip_path))[0]


def _is_safe_member(name: str) -> bool:
    """zip 内のパスが展開先の外を指していないかを判定する（zip slip 対策）。"""
    if name.startswith("/") or name.startswith("\\"):
        return False
    parts = name.replace("\\", "/").split("/")
    return ".." not in parts


def classify_members(names: list[str]) -> tuple[list[str], list[str]]:
    """zip のエントリを (HTML候補, 画像) に振り分ける。

    Google の書き出しは「HTML 1つ + images/ サブフォルダ」だが、階層の付き方は
    書き出し時期やブラウザで揺れる。深さを問わず拡張子で拾う。
    """
    html_files = []
    image_files = []
    for n in names:
        if n.endswith("/") or not _is_safe_member(n):
            continue
        ext = os.path.splitext(n)[1].lower()
        if ext == ".html":
            html_files.append(n)
        elif ext in IMAGE_EXTS:
            image_files.append(n)
    return sorted(html_files), sorted(image_files)


def prepare(zip_path: str, manual_dir: str, force: bool) -> tuple[int, str]:
    """zip を manual_dir へ展開する。戻り値は (画像枚数, 元のHTMLファイル名)。"""
    with zipfile.ZipFile(zip_path) as zf:
        names = zf.namelist()
        html_files, image_files = classify_members(names)

        if not html_files:
            raise RuntimeError(
                "zip の中に .html がありません。"
                "ダウンロード形式が「ウェブページ (.html, zip 形式)」か確認してください"
            )
        if len(html_files) > 1:
            raise RuntimeError(f"zip の中に .html が複数あります: {html_files}")

        source_dest = os.path.join(manual_dir, "source.html")
        if os.path.exists(source_dest) and not force:
            raise RuntimeError(
                f"既に配置済みです: {source_dest}\n"
                "  上書きするなら --force を付けてください"
            )

        os.makedirs(manual_dir, exist_ok=True)

        # 一時ディレクトリへ展開してから移す。展開先を直接 docs/manuals 配下にすると、
        # zip 内の階層がそのまま残って images/ の位置がずれる
        with tempfile.TemporaryDirectory() as tmp:
            members = [n for n in names if _is_safe_member(n)]
            zf.extractall(tmp, members=members)

            shutil.copyfile(os.path.join(tmp, html_files[0]), source_dest)

            img_dest = os.path.join(manual_dir, "images")
            os.makedirs(img_dest, exist_ok=True)
            for rel in image_files:
                shutil.copyfile(
                    os.path.join(tmp, rel),
                    os.path.join(img_dest, os.path.basename(rel)),
                )

    return len(image_files), html_files[0]


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Google Doc のエクスポート zip を docs/manuals/ 配下に展開する",
    )
    parser.add_argument("zip_path", help="ダウンロードした zip のパス")
    parser.add_argument(
        "--name",
        default=None,
        help="マニュアル名（省略時は zip のファイル名から導出）",
    )
    parser.add_argument(
        "--force",
        action="store_true",
        help="既に source.html がある場合も上書きする",
    )
    args = parser.parse_args()

    zip_path = os.path.expanduser(args.zip_path)
    if not os.path.isfile(zip_path):
        print(f"  ERROR: zip が見つかりません: {zip_path}", file=sys.stderr)
        return 1
    if not zipfile.is_zipfile(zip_path):
        print(f"  ERROR: zip ファイルではありません: {zip_path}", file=sys.stderr)
        return 1

    name = args.name or derive_name(zip_path)
    manual_dir = os.path.join(MANUALS_DIR, name)

    print("=== マニュアルの展開 ===")
    print(f"  zip : {zip_path}")
    print(f"  名前: {name}")

    try:
        image_count, original_html = prepare(zip_path, manual_dir, args.force)
    except (RuntimeError, zipfile.BadZipFile) as e:
        sys.stdout.flush()
        print(f"  ERROR: {e}", file=sys.stderr)
        return 1

    print(f"  HTML: {original_html} → source.html")
    print(f"  画像: {image_count}枚 → images/")
    print(f"  配置: {manual_dir}")
    print()
    # マニュアル名は対応表A列のキーになるが、一致すべき相手はタスク一覧M列であって
    # ドキュメントのタイトルではない。両者がずれていると VLOOKUP がエラーも出さずに
    # 外れるため、ここで必ず目視させる
    print("  次にやること:")
    print("    1. 上の「名前」がタスク一覧M列の値と完全に一致するか確認する")
    print("       （末尾スペース・全角アンダースコアなどで静かに外れる）")
    print(f"    2. uv run --project scripts/claude-slide python scripts/claude-slide/generate_slide.py \\")
    print(f"         --prompt card-strict --model claude-opus-4-7 docs/manuals/{name}")
    print("=== 完了 ===")
    return 0


if __name__ == "__main__":
    sys.exit(main())
