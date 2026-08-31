"""生成済みマニュアルHTMLを配信APIへアップロードし、対応表へ貼る行を出力する。

`PUT /manuals/:id` を実行して公開URL（`manual_url`）を受け取り、シフトスプシの
「マニュアルURL」シートへそのまま貼れるタブ区切り1行を組み立てる。

手作業の curl を置き換えるのは、2026-08-31 の実運用で事故が紐付け工程に集中したため。
生成やアップロード自体では失敗しておらず、失敗したのは「どの値をどの列に書くか」だった。

  - マニュアル名を手で打ち、末尾スペースで VLOOKUP が静かに外れた
  - ドキュメント版(tasks.url)とスライド版(tasks.manual_url)の列を取り違えた

そこで、貼る内容を機械が組み立てて列順を固定する。

トークンは環境変数 MANUAL_UPLOAD_TOKEN から読む。未設定なら入力を促す（画面にもシェルの
履歴にも残らない）。`export MANUAL_UPLOAD_TOKEN=...` と直接書くと履歴に平文で残るため。

依存は標準ライブラリのみ。python3 だけで動く。

使い方:
  python3 scripts/automation/upload_manual.py --id en-nichi docs/manuals/45th_企画マニュアル_縁日
  python3 scripts/automation/upload_manual.py --id en-nichi --doc-url https://docs.google.com/... docs/manuals/45th_企画マニュアル_縁日
"""

import argparse
import getpass
import json
import os
import re
import sys
import urllib.error
import urllib.request


SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
PROJECT_ROOT = os.path.normpath(os.path.join(SCRIPT_DIR, "..", ".."))

# HTMLの探索と完全性検査は embed_images.py と同じ規則を使う。重複させると
# 片方だけ直したときに挙動がずれるため、実装を1つに保つ
sys.path.insert(0, os.path.join(PROJECT_ROOT, "scripts", "claude-slide"))
from embed_images import find_slide_html, is_complete_html  # noqa: E402

DEFAULT_BASE_URL = "https://seeft-api.nutfes.net"

# api/lib/usecase/manual_usecase.go の manualIDRe と揃える。
# サーバ側でも弾かれるが、20MBを送ってから400を受け取るのは無駄なので手前で検査する
MANUAL_ID_RE = re.compile(r"^[a-z0-9_-]{1,64}$")

# api/lib/internals/controller/manual_controller.go の manualUploadMaxBytes と揃える
MAX_UPLOAD_BYTES = 20 << 20


def resolve_manual_dir(arg: str) -> str:
    arg = arg.rstrip("/")
    if os.path.isabs(arg):
        return arg
    if os.path.isdir(arg):
        return os.path.abspath(arg)
    return os.path.join(PROJECT_ROOT, arg)


def read_token() -> str:
    """トークンを取得する。環境変数になければ入力を促す（履歴に残さないため）。"""
    token = os.environ.get("MANUAL_UPLOAD_TOKEN", "").strip()
    if token:
        return token
    return getpass.getpass("アップロードトークンを貼り付けてEnter（表示されません）: ").strip()


def upload(base_url: str, manual_id: str, html_path: str, token: str) -> dict:
    """PUT /manuals/:id を実行し、レスポンスのJSONを返す。"""
    with open(html_path, "rb") as f:
        body = f.read()

    url = f"{base_url.rstrip('/')}/manuals/{manual_id}"
    req = urllib.request.Request(url, data=body, method="PUT")
    req.add_header("Authorization", f"Bearer {token}")
    req.add_header("Content-Type", "text/html")

    with urllib.request.urlopen(req, timeout=180) as resp:
        return json.loads(resp.read().decode("utf-8"))


def describe_http_error(e: urllib.error.HTTPError) -> str:
    """サーバが返すエラーを、運用者が次の手を判断できる文言に変える。"""
    hints = {
        401: "トークンが違うか、サーバ側で MANUAL_UPLOAD_TOKEN が未設定です"
             "（設定漏れで素通りするのを防ぐため、未設定でも401になります）",
        400: "マニュアルIDが不正です。a-z 0-9 _ - の1〜64文字のみ使えます",
        413: "ファイルが20MBを超えています。画像の枚数やサイズを確認してください",
        500: "サーバ側の保存に失敗しました。APIのログを確認してください",
    }
    hint = hints.get(e.code, "")
    return f"HTTP {e.code}" + (f"\n         {hint}" if hint else "")


def main() -> int:
    parser = argparse.ArgumentParser(
        description="マニュアルHTMLを配信APIへアップロードし、対応表へ貼る行を出力する",
    )
    parser.add_argument("manual_dir", help="マニュアルディレクトリ")
    parser.add_argument(
        "--id",
        required=True,
        help="公開URLに使う識別子（例: en-nichi）。a-z 0-9 _ - の1〜64文字",
    )
    parser.add_argument(
        "--doc-url",
        default=None,
        help="Googleドキュメントの共有URL。渡すと対応表のB列も埋めた行を出力する",
    )
    parser.add_argument(
        "--file",
        default=None,
        help="アップロードするHTMLを明示する（省略時は manual_dir 内の slide*.html を自動検出）",
    )
    parser.add_argument(
        "--base-url",
        default=os.environ.get("SEEFT_API_BASE_URL", DEFAULT_BASE_URL),
        help=f"APIのベースURL（既定: {DEFAULT_BASE_URL}）",
    )
    args = parser.parse_args()

    if not MANUAL_ID_RE.match(args.id):
        print(
            f"  ERROR: --id が命名規則に反しています: {args.id!r}\n"
            "         a-z 0-9 _ - の1〜64文字のみ（日本語・大文字・スラッシュは不可）",
            file=sys.stderr,
        )
        return 1

    manual_dir = resolve_manual_dir(args.manual_dir)
    if not os.path.isdir(manual_dir):
        print(f"  ERROR: ディレクトリがありません: {manual_dir}", file=sys.stderr)
        return 1

    try:
        html_path = args.file or find_slide_html(manual_dir)
    except (FileNotFoundError, RuntimeError) as e:
        print(f"  ERROR: {e}", file=sys.stderr)
        return 1

    size = os.path.getsize(html_path)
    print("=== マニュアルのアップロード ===")
    print(f"  対象: {html_path} ({size//1024}KB)")
    print(f"  ID  : {args.id}")

    with open(html_path, "r", encoding="utf-8") as f:
        html = f.read()

    # 途中で切れたHTMLを配信すると、閲覧者には白紙や崩れたページが見える。
    # アップロードは即時反映されるため、送る前に止める
    if not is_complete_html(html):
        sys.stdout.flush()
        print(
            "  ERROR: HTMLが </html> で閉じていません。応答が途中で切れています。\n"
            "         生成をやり直してください",
            file=sys.stderr,
        )
        return 1

    if size > MAX_UPLOAD_BYTES:
        sys.stdout.flush()
        print(
            f"  ERROR: {size//1024//1024}MB は上限の20MBを超えています",
            file=sys.stderr,
        )
        return 1

    token = read_token()
    if not token:
        print("  ERROR: トークンが空です", file=sys.stderr)
        return 1

    try:
        result = upload(args.base_url, args.id, html_path, token)
    except urllib.error.HTTPError as e:
        sys.stdout.flush()
        print(f"  ERROR: アップロードに失敗しました: {describe_http_error(e)}", file=sys.stderr)
        return 1
    except urllib.error.URLError as e:
        sys.stdout.flush()
        print(f"  ERROR: APIに接続できません: {e.reason}", file=sys.stderr)
        return 1

    manual_url = result.get("manual_url", "")
    if not manual_url:
        print(f"  ERROR: レスポンスに manual_url がありません: {result}", file=sys.stderr)
        return 1

    print(f"  成功: {manual_url}")
    print()

    manual_name = os.path.basename(manual_dir)
    doc_url = args.doc_url or ""
    print("  「マニュアルURL」シートに貼る行（タブ区切り）:")
    print()
    print(f"{manual_name}\t{doc_url}\t{manual_url}")
    print()
    # A列はタスク一覧M列と完全一致していなければ VLOOKUP が外れる。ここで組み立てた
    # 名前はディレクトリ名の写しであって、M列の値である保証はない
    print("  貼る前に確認すること:")
    print("    1行目の値がタスク一覧M列と完全に一致しているか。ずれていると")
    print("    エラーを出さずにS列が空になり、シフトカードにボタンが出ない")
    if not doc_url:
        print("    B列（ドキュメントURL）が空。--doc-url を渡すか、手で埋める")
    print("=== 完了 ===")
    return 0


if __name__ == "__main__":
    sys.exit(main())
