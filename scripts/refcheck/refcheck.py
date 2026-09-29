#!/usr/bin/env python3
# Markdown に書いたファイルパスとシンボルが、リポジトリにまだ存在するかを確かめる。
#
# 点検するのはインラインコード（`...`）だけ。書き方の決まりは docs/decisions/README.md。
#   `api/lib/usecase/shift_usecase.go`                       パス
#   `docs/operations/`                                        ディレクトリ
#   `api/lib/usecase/shift_usecase.go#UpdateShiftsFromGAS`    シンボル（ファイル中に単語として現れるか）
#   `api/lib/usecase/shift_usecase.go#shiftUseCase.UpdateShiftsFromGAS`  型.メソッド
#
# 使い方:
#   python3 scripts/refcheck/refcheck.py docs/decisions/*.md             # git 管理下のファイルだけを「存在する」とみなす
#   python3 scripts/refcheck/refcheck.py --worktree path/to/notes.md     # 手元のディスクにあれば「存在する」とみなす
#   python3 scripts/refcheck/refcheck.py --base HEAD^1 docs/decisions/*.md  # 変更前の版にあった名前も「パスらしい」とみなす（CI 用）
#
# 終了コード: 0 = 問題なし、1 = 見つからないものがある、2 = 使い方の誤り

import argparse
import os
import re
import subprocess
import sys
from dataclasses import dataclass
from pathlib import Path

INLINE_CODE = re.compile(r"(`+)(.+?)\1")
FENCE = re.compile(r"^\s*(`{3,}|~{3,})")
CLOSING_FENCE = re.compile(r"^\s*(`{3,}|~{3,})\s*$")
LINE_NUMBER = re.compile(r"^(.+?):\d+(?:-\d+)?$")

IGNORE_LINE = "<!-- refcheck:ignore -->"
OFF = "<!-- refcheck:off -->"
ON = "<!-- refcheck:on -->"

# これを含むものはプレースホルダとみなして点検しない
PLACEHOLDER_CHARS = set("*?{}<>")


@dataclass(frozen=True)
class Finding:
    line: int
    kind: str  # "missing path" / "missing symbol" / "line-number"
    text: str


class Repo:
    # tracked=True なら git ls-files の結果だけを、False なら手元のディスクを見る
    def __init__(self, root: Path, tracked: bool, base: str | None = None):
        self.root = root
        self.tracked = tracked
        self._files: set[str] = set()
        self._dirs: set[str] = set()
        if tracked:
            out = subprocess.run(
                ["git", "-C", str(root), "ls-files", "-z"],
                check=True,
                capture_output=True,
            ).stdout.decode("utf-8")
            for f in filter(None, out.split("\0")):
                self._files.add(f)
                parts = f.split("/")
                for i in range(1, len(parts)):
                    self._dirs.add("/".join(parts[:i]))
            self._top = {f.split("/")[0] for f in self._files}
            self._root_files = {f for f in self._files if "/" not in f}
        else:
            entries = os.listdir(root)
            self._top = set(entries)
            self._root_files = {e for e in entries if (root / e).is_file()}
        # PR で直下のディレクトリやファイルを丸ごと消すと、今の版だけではパスだと分からなくなる。
        # 変更前の版の直下の名前も足しておき、消えた参照を「missing path」として拾う
        if base is not None:
            out = subprocess.run(
                ["git", "-C", str(root), "ls-tree", "-z", "--name-only", base],
                check=True,
                capture_output=True,
            ).stdout.decode("utf-8")
            for name in filter(None, out.split("\0")):
                self._top.add(name)
                self._root_files.add(name)

    def is_file(self, path: str) -> bool:
        if self.tracked:
            return path in self._files
        return (self.root / path).is_file()

    def is_dir(self, path: str) -> bool:
        if self.tracked:
            return path in self._dirs
        return (self.root / path).is_dir()

    # リポジトリ内のパスを書いたものとみなすか。
    # 1段目がリポジトリ直下に実在する名前のときだけ対象にする（`application/json` などを拾わないため）
    def looks_like_path(self, path: str) -> bool:
        if "/" in path:
            return path.split("/")[0] in self._top
        return path in self._root_files

    # .gitignore で管理から外したもの（.env や生成物）は CI では確かめようがないので点検しない
    def is_ignored(self, path: str) -> bool:
        if not self.tracked:
            return False
        r = subprocess.run(["git", "-C", str(self.root), "check-ignore", "-q", path])
        return r.returncode == 0

    def read(self, path: str) -> str:
        return (self.root / path).read_text(encoding="utf-8", errors="replace")


def _has_word(text: str, word: str) -> bool:
    return re.search(r"(?<!\w)" + re.escape(word) + r"(?!\w)", text) is not None


def check_span(span: str, repo: Repo) -> Finding | None:
    s = span.strip()
    if not s or any(c.isspace() for c in s):
        return None
    if "://" in s or s.startswith(("~", "/")):
        return None
    if PLACEHOLDER_CHARS & set(s):
        return None

    m = LINE_NUMBER.match(s)
    if m and repo.looks_like_path(m.group(1).split("#")[0]):
        return Finding(0, "line-number", s)

    path, _, symbol = s.partition("#")
    if not repo.looks_like_path(path):
        return None
    if repo.is_ignored(path):
        return None

    if not symbol:
        if repo.is_file(path.rstrip("/")) or repo.is_dir(path.rstrip("/")):
            return None
        return Finding(0, "missing path", s)

    if not repo.is_file(path):
        return Finding(0, "missing path", s)
    parts = symbol.split(".")
    text = repo.read(path)
    if all(parts) and all(_has_word(text, p) for p in parts):
        return None
    return Finding(0, "missing symbol", s)


def check_text(text: str, repo: Repo) -> list[Finding]:
    findings: list[Finding] = []
    fence: str | None = None
    off = False
    for lineno, line in enumerate(text.splitlines(), start=1):
        m = FENCE.match(line)
        if fence is None and m:
            fence = m.group(1)
            continue
        if fence is not None:
            # 閉じるのは、開きと同じ記号が同じ数以上並び、後ろに空白しかない行だけ
            cm = CLOSING_FENCE.match(line)
            if cm and cm.group(1)[0] == fence[0] and len(cm.group(1)) >= len(fence):
                fence = None
            continue
        if OFF in line:
            off = True
        if ON in line:
            off = False
            continue
        if off or IGNORE_LINE in line:
            continue
        for cm in INLINE_CODE.finditer(line):
            f = check_span(cm.group(2), repo)
            if f is not None:
                findings.append(Finding(lineno, f.kind, f.text))
    return findings


def _git_root(start: Path) -> Path:
    r = subprocess.run(
        ["git", "-C", str(start), "rev-parse", "--show-toplevel"],
        capture_output=True,
        text=True,
    )
    return Path(r.stdout.strip()) if r.returncode == 0 else start


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description="Markdown に書いたパスとシンボルが存在するかを確かめる")
    parser.add_argument("files", nargs="+", help="点検する Markdown ファイル")
    mode = parser.add_mutually_exclusive_group()
    mode.add_argument("--tracked", action="store_true", help="git 管理下のファイルだけを見る（既定）")
    mode.add_argument("--worktree", action="store_true", help="手元のディスクにあるファイルを見る")
    parser.add_argument("--root", type=Path, help="リポジトリのルート（既定は今いる場所の git のルート）")
    parser.add_argument("--base", help="変更前の版（例: HEAD^1）。その版の直下にあった名前もパスとして点検する")
    args = parser.parse_args(argv)

    root = (args.root or _git_root(Path.cwd())).resolve()
    try:
        repo = Repo(root, tracked=not args.worktree, base=args.base)
    except subprocess.CalledProcessError as e:
        print(f"git の実行に失敗しました: {' '.join(e.cmd)}", file=sys.stderr)
        return 2

    status = 0
    for name in args.files:
        p = Path(name)
        if not p.is_file():
            print(f"{name}: ファイルがありません", file=sys.stderr)
            return 2
        for f in check_text(p.read_text(encoding="utf-8"), repo):
            print(f"{name}:{f.line}: {f.kind} {f.text}")
            status = 1
    return status


if __name__ == "__main__":
    sys.exit(main())
