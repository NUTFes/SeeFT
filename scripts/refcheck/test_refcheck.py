# refcheck.py のテスト。実行: python3 -m unittest discover -s scripts/refcheck
import os
import subprocess
import tempfile
import unittest
from pathlib import Path

from refcheck import Repo, check_text

GO_SRC = """package usecase

type shiftUseCase struct{}

func (s *shiftUseCase) UpdateShiftsFromGAS() error { return nil }

const MaxShift = 10
"""


class RefcheckTest(unittest.TestCase):
    def setUp(self):
        self._tmp = tempfile.TemporaryDirectory()
        self.root = Path(self._tmp.name)
        files = {
            "api/lib/usecase/shift_usecase.go": GO_SRC,
            "gas/shift/コード.js": "const YEAR_ID = 45;\nfunction 送信() {}\n",
            "docs/operations/README.md": "# 運用\n",
            "AGENTS.md": "# AGENTS\n",
        }
        for path, body in files.items():
            p = self.root / path
            p.parent.mkdir(parents=True, exist_ok=True)
            p.write_text(body, encoding="utf-8")
        (self.root / ".gitignore").write_text("mobile/env/.env\n", encoding="utf-8")
        subprocess.run(["git", "init", "-q", str(self.root)], check=True)
        subprocess.run(["git", "-C", str(self.root), "add", "."], check=True)
        # git 管理外のファイル（raw/ のような置き場所）
        (self.root / "raw").mkdir()
        (self.root / "raw" / "memo.md").write_text("メモ\n", encoding="utf-8")
        self.repo = Repo(self.root, tracked=True)

    def tearDown(self):
        self._tmp.cleanup()

    def kinds(self, text, repo=None):
        return [(f.line, f.kind, f.text) for f in check_text(text, repo or self.repo)]

    def test_existing_paths_pass(self):
        text = "\n".join(
            [
                "`api/lib/usecase/shift_usecase.go`",
                "`docs/operations/`",
                "`docs/operations`",
                "`AGENTS.md`",
                "`gas/shift/コード.js`",
            ]
        )
        self.assertEqual(self.kinds(text), [])

    def test_missing_path(self):
        self.assertEqual(
            self.kinds("前提: `api/lib/usecase/task_usecase.go`"),
            [(1, "missing path", "api/lib/usecase/task_usecase.go")],
        )

    def test_missing_directory(self):
        self.assertEqual(self.kinds("`docs/decisions/`"), [(1, "missing path", "docs/decisions/")])

    def test_symbols(self):
        ok = "\n".join(
            [
                "`api/lib/usecase/shift_usecase.go#UpdateShiftsFromGAS`",
                "`api/lib/usecase/shift_usecase.go#shiftUseCase.UpdateShiftsFromGAS`",
                "`api/lib/usecase/shift_usecase.go#MaxShift`",
                "`gas/shift/コード.js#YEAR_ID`",
                "`gas/shift/コード.js#送信`",
            ]
        )
        self.assertEqual(self.kinds(ok), [])

    def test_missing_symbol(self):
        text = "\n".join(
            [
                "`api/lib/usecase/shift_usecase.go#RegisterTask`",
                # 部分一致では通さない
                "`api/lib/usecase/shift_usecase.go#UpdateShifts`",
                "`api/lib/usecase/shift_usecase.go#taskUseCase.UpdateShiftsFromGAS`",
                "`api/lib/usecase/nothing.go#Foo`",
            ]
        )
        self.assertEqual(
            self.kinds(text),
            [
                (1, "missing symbol", "api/lib/usecase/shift_usecase.go#RegisterTask"),
                (2, "missing symbol", "api/lib/usecase/shift_usecase.go#UpdateShifts"),
                (3, "missing symbol", "api/lib/usecase/shift_usecase.go#taskUseCase.UpdateShiftsFromGAS"),
                (4, "missing path", "api/lib/usecase/nothing.go#Foo"),
            ],
        )

    def test_line_numbers_are_rejected(self):
        text = "`api/lib/usecase/shift_usecase.go:5`\n`api/lib/usecase/shift_usecase.go:5-7`"
        self.assertEqual(
            self.kinds(text),
            [
                (1, "line-number", "api/lib/usecase/shift_usecase.go:5"),
                (2, "line-number", "api/lib/usecase/shift_usecase.go:5-7"),
            ],
        )

    def test_not_paths_are_ignored(self):
        text = "\n".join(
            [
                "`application/json`",
                "`localhost:45029`",
                "`RegisterTask`",
                "`feat/{username}/{issue}`",
                "`docs/decisions/<番号>-<題>.md`",
                "`manuals/*.html`",
                "`~/.claude/memory/x.md`",
                "`/tmp/x`",
                "`https://github.com/NUTFes/SeeFT/api/x`",
                "`make up-api`",
            ]
        )
        self.assertEqual(self.kinds(text), [])

    def test_fenced_block_is_ignored(self):
        text = "```text\n`api/nothing.go`\napi/nothing.go:10\n```\n`api/nothing2.go`"
        self.assertEqual(self.kinds(text), [(5, "missing path", "api/nothing2.go")])

    def test_fence_closes_only_on_bare_fence(self):
        # 後ろに文字がある行や、開きより短い行ではブロックは閉じない
        text = "\n".join(
            [
                "```text",
                "```not-a-close",
                "`api/nothing.go`",
                "```",
                "`api/nothing2.go`",
                "````",
                "```",
                "`api/nothing3.go`",
                "````",
            ]
        )
        self.assertEqual(self.kinds(text), [(5, "missing path", "api/nothing2.go")])

    def test_deleted_top_level_is_detected_with_base(self):
        # 直下のディレクトリを丸ごと消しても、変更前の版を渡せば消えた参照を拾う
        env = {"GIT_AUTHOR_NAME": "t", "GIT_AUTHOR_EMAIL": "t@example.com",
               "GIT_COMMITTER_NAME": "t", "GIT_COMMITTER_EMAIL": "t@example.com"}
        git = ["git", "-C", str(self.root)]
        subprocess.run(git + ["commit", "-q", "-m", "base"], check=True, env={**os.environ, **env})
        subprocess.run(git + ["rm", "-q", "-r", "gas", "AGENTS.md"], check=True)
        text = "`gas/shift/コード.js`\n`AGENTS.md`\n`application/json`"
        self.assertEqual(self.kinds(text, Repo(self.root, tracked=True)), [])
        self.assertEqual(
            self.kinds(text, Repo(self.root, tracked=True, base="HEAD")),
            [(1, "missing path", "gas/shift/コード.js"), (2, "missing path", "AGENTS.md")],
        )

    def test_ignore_markers(self):
        text = "\n".join(
            [
                "`api/nothing.go` <!-- refcheck:ignore -->",
                "<!-- refcheck:off -->",
                "`api/nothing2.go`",
                "<!-- refcheck:on -->",
                "`api/nothing3.go`",
            ]
        )
        self.assertEqual(self.kinds(text), [(5, "missing path", "api/nothing3.go")])

    def test_gitignored_paths_are_ignored(self):
        # .env のように、わざと管理から外したものは CI では確かめられない
        self.assertEqual(self.kinds("`mobile/env/.env`"), [])

    def test_tracked_and_worktree(self):
        # git 管理外の raw/ は、tracked では対象外（1段目が管理下に無い）、worktree では存在する
        self.assertEqual(self.kinds("`raw/memo.md`"), [])
        worktree = Repo(self.root, tracked=False)
        self.assertEqual(self.kinds("`raw/memo.md`", worktree), [])
        self.assertEqual(
            self.kinds("`raw/nothing.md`", worktree),
            [(1, "missing path", "raw/nothing.md")],
        )


if __name__ == "__main__":
    unittest.main()
