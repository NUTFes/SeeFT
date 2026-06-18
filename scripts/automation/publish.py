"""生成した解説マニュアル HTML を GitHub Pages 用リポジトリに配置する補助ツール。

新フロー（Slack スレッド中心・手作業運用）の「URL 発行」ステップを担う。
PM/SeeFT 担当が手で叩く想定。Slack への貼り付けは人が行う（このツールは URL を出すだけ）。

配置構造:
    <pages-repo>/
    └── manuals/
        ├── 01/index.html
        ├── 02/index.html
        └── ...
    マニュアルは番号で識別する。番号はディレクトリ名の先頭（`01_44th_…` → `01`）から
    自動推定するが、先頭が数字でないマニュアルは --number で明示する。

依存: 標準ライブラリのみ（`python3 scripts/automation/publish.py …` で直接動く）。
公開先リポジトリのパスと公開 URL はリポジトリにハードコードせず環境変数で外部化する。

使い方:
    SEEFT_PAGES_REPO=~/work/seeft-manuals-pages \
    SEEFT_PAGES_BASE_URL=https://nutfes.github.io/seeft-manuals \
    python3 scripts/automation/publish.py docs/manuals/01_44th_のぼり広告設置マニュアル

    # 番号を明示し、commit までで止めず push まで実行
    python3 scripts/automation/publish.py docs/manuals/44th_企画マニュアル_お化け屋敷 \
        --number 08 --push
"""

from __future__ import annotations

import argparse
import os
import re
import shutil
import subprocess
import sys

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
PROJECT_ROOT = os.path.normpath(os.path.join(SCRIPT_DIR, "..", ".."))

DEFAULT_HTML_NAME = "slide_claude.card-strict.html"
_LEADING_NUM_RE = re.compile(r"^(\d+)")


def resolve_manual_dir(arg: str) -> str:
    arg = arg.rstrip("/")
    if os.path.isabs(arg):
        return arg
    if os.path.isdir(arg):
        return os.path.abspath(arg)
    return os.path.join(PROJECT_ROOT, arg)


def derive_number(manual_dir: str, explicit: str | None) -> str:
    """マニュアル番号を決める。--number 優先、無ければディレクトリ名先頭の数字。"""
    if explicit:
        return explicit.zfill(2)
    name = os.path.basename(manual_dir.rstrip("/"))
    m = _LEADING_NUM_RE.match(name)
    if not m:
        raise SystemExit(
            f"ERROR: '{name}' から番号を推定できません。--number NN で明示してください。"
        )
    return m.group(1).zfill(2)


def git(repo: str, *args: str) -> None:
    subprocess.run(["git", "-C", repo, *args], check=True)


def main() -> int:
    parser = argparse.ArgumentParser(
        description="解説マニュアル HTML を GitHub Pages リポジトリへ配置（URL 発行補助）",
    )
    parser.add_argument("manual_dir", help="マニュアルディレクトリ")
    parser.add_argument("--number", default=None, help="マニュアル番号（例: 01）。未指定ならディレクトリ名から推定")
    parser.add_argument("--html-name", default=DEFAULT_HTML_NAME, help=f"配置する HTML 名 (default: {DEFAULT_HTML_NAME})")
    parser.add_argument(
        "--pages-repo",
        default=os.environ.get("SEEFT_PAGES_REPO"),
        help="GitHub Pages 用リポジトリのローカルパス（env SEEFT_PAGES_REPO でも可）",
    )
    parser.add_argument(
        "--base-url",
        default=os.environ.get("SEEFT_PAGES_BASE_URL", ""),
        help="公開 URL のベース（env SEEFT_PAGES_BASE_URL でも可。例: https://nutfes.github.io/seeft-manuals）",
    )
    parser.add_argument("--push", action="store_true", help="commit 後に push まで実行（既定は commit で停止し push コマンドを表示）")
    parser.add_argument("--message", default=None, help="commit メッセージ（未指定なら自動生成）")
    args = parser.parse_args()

    manual_dir = resolve_manual_dir(args.manual_dir)
    src_html = os.path.join(manual_dir, args.html_name)
    if not os.path.isfile(src_html):
        print(f"ERROR: HTML が見つかりません: {src_html}", file=sys.stderr)
        return 2

    if not args.pages_repo:
        print(
            "ERROR: 公開先リポジトリが未指定です。SEEFT_PAGES_REPO か --pages-repo で指定してください。",
            file=sys.stderr,
        )
        return 2
    pages_repo = os.path.abspath(os.path.expanduser(args.pages_repo))
    if not os.path.isdir(os.path.join(pages_repo, ".git")):
        print(f"ERROR: git リポジトリではありません: {pages_repo}", file=sys.stderr)
        return 2

    number = derive_number(manual_dir, args.number)
    rel_path = os.path.join("manuals", number, "index.html")
    dest = os.path.join(pages_repo, rel_path)
    os.makedirs(os.path.dirname(dest), exist_ok=True)
    shutil.copy2(src_html, dest)

    size_kb = os.path.getsize(dest) // 1024
    print(f"=== 配置 ===")
    print(f"  {os.path.basename(manual_dir)} → manuals/{number}/index.html ({size_kb}KB)")

    message = args.message or f"manuals/{number}: {os.path.basename(manual_dir)} を更新"
    git(pages_repo, "add", rel_path)
    git(pages_repo, "commit", "-m", message)

    if args.push:
        git(pages_repo, "push")
        print("  push 完了")
    else:
        print(f"  commit まで完了。公開するには: git -C {pages_repo} push")

    url = f"{args.base_url.rstrip('/')}/manuals/{number}/" if args.base_url else f"(SEEFT_PAGES_BASE_URL 未設定)/manuals/{number}/"
    print(f"  公開 URL: {url}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
