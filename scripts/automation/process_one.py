"""1 マニュアル分の end-to-end オーケストレーション。

責務:
  1. スプシ (or CSV) から行を読む
  2. ステージ 3 (副委員長 + 執行部) 完了済か検証
  3. ステータスを「処理中」にセット (watcher の lease/claim と同じ値で揃える)
  4. Drive から Google Doc を HTML (zip) でダウンロード → manual_dir 配下に展開
  5. 「修正提案」列の内容があれば instructions.md へ書き込む (CSV/Sheets 両対応)
  6. generate_slide.py を呼んで Claude で HTML 生成 (instructions.md を追加指示として注入)
  7. uploader.py で配信先にアップロード
  8. スプシに「生成HTML URL」「最終生成日時」「HTML生成ステータス=完了」を書き戻す
     (--mode regen のときは「比較確認状況=再確認待ち」も書き戻す)

例外時はステータスを「エラー」にする。--mode regen の場合は「比較確認状況」を
「再生成依頼」へ戻し再試行できるようにする。詳細は備考に書く。

使い方:
  # CSV モード (OAuth 不要、Drive ダウンロードはスキップ前提)
  uv run --project scripts/automation python scripts/automation/process_one.py \\
      --csv-path docs/spread_sheets/45th_マニュアル割り当て.csv \\
      --skip-drive \\
      --manual-dir docs/manuals/01_44th_配線マニュアル \\
      配線マニュアル

  # Sheets API モード (本番)
  uv run --project scripts/automation python scripts/automation/process_one.py \\
      --spreadsheet-id <ID> \\
      配線マニュアル
"""

from __future__ import annotations

import argparse
import datetime as dt
import os
import subprocess
import sys
import traceback


SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
PROJECT_ROOT = os.path.normpath(os.path.join(SCRIPT_DIR, "..", ".."))

# 既存の解説スライド生成スクリプト
GENERATE_SCRIPT = os.path.join(PROJECT_ROOT, "scripts", "claude-slide", "generate_slide.py")
GENERATE_PROJECT = os.path.join(PROJECT_ROOT, "scripts", "claude-slide")
DEFAULT_PROMPT_VARIANT = "card"

# ステータス列の値 (sheets_client.COLUMNS の列名と整合)
STATUS_COL = "HTML生成ステータス"
URL_COL = "生成HTML URL"
DATETIME_COL = "最終生成日時"
NOTE_COL = "備考"
COMPARISON_COL = "比較確認状況"
SUGGESTION_COL = "修正提案"

# watcher.py の CLAIMED_VALUE と同じ値 (lease 表示用。二重起動防止の実体は watcher 側で完結)
CLAIMED_VALUE = "処理中"
REGEN_DONE_VALUE = "再確認待ち"
REGEN_REQUEST_VALUE = "再生成依頼"

# ステージ 3 完了とみなす値 (44th xlsx 慣例: 完成 / 確認不要)
STAGE_DONE_VALUES = {"完成", "確認不要"}


def _now_str() -> str:
    return dt.datetime.now().strftime("%Y-%m-%d %H:%M")


def _slug_for_dir(manual_name: str) -> str:
    """マニュアル名 → manual_dir ディレクトリ名のサジェスト。

    44th 既存ディレクトリと衝突しないよう、45th_ プレフィックスを付ける。
    既に PM が手動で manual_dir を切ってる場合は --manual-dir で上書き可能。
    """
    separators = tuple(s for s in (os.sep, os.altsep) if s)
    if not manual_name or any(separator in manual_name for separator in separators):
        raise ValueError("マニュアル名にパス区切りは使用できません")
    return f"45th_{manual_name}"


def validate_stage_complete(row, force: bool) -> tuple[bool, str]:
    """ステージ 3 完了済か (副委員長 + 執行部) 検証。force=True なら警告のみ。"""
    fukuiincho = row.get("副委員長候補確認")
    shikkoubu = row.get("執行部確認")
    missing = []
    if fukuiincho not in STAGE_DONE_VALUES:
        missing.append(f"副委員長候補確認={fukuiincho!r}")
    if shikkoubu not in STAGE_DONE_VALUES:
        missing.append(f"執行部確認={shikkoubu!r}")
    if missing:
        msg = f"ステージ 3 未完了: {', '.join(missing)}"
        if force:
            print(f"WARNING: {msg} (--force で強行)")
            return True, msg
        return False, msg
    return True, ""


def write_instructions(manual_dir: str, suggestion: str) -> None:
    """「修正提案」列の内容を manual_dir/instructions.md へ書く。

    generate_slide.py はこのファイルを直接読み、「この回の追加・修正指示 (最優先)」
    として生成プロンプトの末尾に注入する。空文字なら何もしない (PM が
    instructions.md を直接手編集している既存運用を壊さないため)。
    """
    if not suggestion.strip():
        return
    instr_path = os.path.join(manual_dir, "instructions.md")
    with open(instr_path, "w", encoding="utf-8") as f:
        f.write(suggestion.strip() + "\n")
    print(f"  修正提案を書き込み: {instr_path}")


def run_generate(manual_dir: str, prompt_variant: str) -> str:
    """既存 generate_slide.py を subprocess 起動。

    generate_slide.py は --status-csv を受け付けない (存在しない CLI フラグ)。
    「修正提案」列の内容は呼び出し元 (process_one) が事前に write_instructions() で
    manual_dir/instructions.md へ書いており、generate_slide.py がそれを直接読む。

    戻り値: 生成された HTML ファイルの絶対パス
    """
    cmd = [
        "uv", "run", "--project", GENERATE_PROJECT, "python", GENERATE_SCRIPT,
        "--prompt", prompt_variant,
        manual_dir,
    ]

    print(f"  Running: {' '.join(cmd)}")
    subprocess.run(cmd, check=True)

    # 出力パスは generate_slide.py のルール (slide_claude.{variant}.html)
    output_filename = (
        "slide_claude.html" if prompt_variant == "default"
        else f"slide_claude.{prompt_variant}.html"
    )
    output_path = os.path.join(manual_dir, output_filename)
    if not os.path.isfile(output_path):
        raise RuntimeError(f"Generated HTML not found at expected path: {output_path}")
    return output_path


def process_one(args) -> int:
    # sheets_client / uploader / drive_client を遅延 import (CSV モードで Drive 不要なため)
    from sheets_client import SheetsClient
    from uploader import upload

    client = SheetsClient(
        spreadsheet_id=args.spreadsheet_id,
        csv_path=args.csv_path,
    )
    print(f"=== process_one: {args.manual_name} ({client.mode} モード) ===")

    # 1. 行を読む
    row = client.read_row(args.manual_name)
    if row is None:
        print(f"ERROR: Manual not found: {args.manual_name}", file=sys.stderr)
        return 2
    print(f"  行 {row.row_index} を読み込み")

    # 2. ステージ 3 完了済か検証
    ok, msg = validate_stage_complete(row, force=args.force)
    if not ok:
        print(f"ABORT: {msg}", file=sys.stderr)
        return 3

    # 3. Doc URL 存在チェック
    raw_doc_url = row.get("Google Doc URL")
    doc_url = raw_doc_url.strip() if isinstance(raw_doc_url, str) else ""
    if not doc_url and not args.skip_drive:
        print("ABORT: Google Doc URL が空です", file=sys.stderr)
        return 4

    # 4. ステータスを「処理中」に (watcher の claim と同じ値。二重起動防止の可視化)
    if not args.dry_run:
        client.write_cells(args.manual_name, {STATUS_COL: CLAIMED_VALUE})
        print(f"  ステータス → {CLAIMED_VALUE}")

    try:
        # 5. manual_dir を確定
        manual_dir = args.manual_dir
        if not manual_dir:
            manual_dir = os.path.join(PROJECT_ROOT, "docs", "manuals", _slug_for_dir(args.manual_name))
        manual_dir = os.path.abspath(manual_dir)
        os.makedirs(manual_dir, exist_ok=True)
        print(f"  Manual dir: {manual_dir}")

        # 6. Drive ダウンロード (skip 可)
        if not args.skip_drive:
            from drive_client import download_doc_to_manual_dir
            print(f"  Doc ダウンロード: {doc_url}")
            html_path = download_doc_to_manual_dir(doc_url, manual_dir)
            print(f"  Doc HTML: {html_path}")

        # 7. 修正提案 → instructions.md、generate_slide.py 起動
        if not args.skip_generation:
            write_instructions(manual_dir, row.get(SUGGESTION_COL))
            generated_path = run_generate(manual_dir, args.prompt)
            print(f"  生成完了: {generated_path}")
        else:
            # スキップ時は既存ファイルを使う
            variant = args.prompt
            fname = "slide_claude.html" if variant == "default" else f"slide_claude.{variant}.html"
            generated_path = os.path.join(manual_dir, fname)
            if not os.path.isfile(generated_path):
                raise RuntimeError(f"--skip-generation 指定だが既存ファイルなし: {generated_path}")
            print(f"  生成スキップ、既存ファイルを使用: {generated_path}")

        # 8. アップロード
        url = upload(generated_path, key=args.manual_name)
        print(f"  Upload URL: {url}")

        # 9. スプシ書き戻し
        updates = {
            STATUS_COL: "完了",
            URL_COL: url,
            DATETIME_COL: _now_str(),
        }
        if args.mode == "regen":
            # 再確認待ちへ遷移させ、次回スキャンで同じ行を再生成し続けないようにする
            updates[COMPARISON_COL] = REGEN_DONE_VALUE
        if args.dry_run:
            print(f"  [dry-run] スプシ書き戻しはスキップ。更新内容: {updates}")
        else:
            client.write_cells(args.manual_name, updates)
            print("  スプシ更新完了")

        print("=== 完了 ===")
        return 0

    except Exception as e:
        print(f"ERROR: {e}", file=sys.stderr)
        traceback.print_exc()
        if not args.dry_run:
            error_updates = {
                STATUS_COL: "エラー",
                NOTE_COL: f"[{_now_str()}] {type(e).__name__}: {e}",
            }
            if args.mode == "regen":
                # 再試行できるよう claim 前の状態へ戻す
                error_updates[COMPARISON_COL] = REGEN_REQUEST_VALUE
            client.write_cells(args.manual_name, error_updates)
        return 1


def main() -> int:
    parser = argparse.ArgumentParser(
        description="1 マニュアル分の生成パイプラインを end-to-end で走らせる",
    )
    parser.add_argument("manual_name", help="スプシの「マニュアル名」列の値")
    # バックエンド選択 (片方だけ)
    group = parser.add_mutually_exclusive_group(required=True)
    group.add_argument("--csv-path", help="ローカル CSV モード (dry-run/開発用)")
    group.add_argument("--spreadsheet-id", help="Sheets API モード (本番)")
    # オプション
    parser.add_argument(
        "--manual-dir",
        help="マニュアルディレクトリ (省略時は docs/manuals/45th_<name>/ を生成)",
    )
    parser.add_argument(
        "--prompt", default=DEFAULT_PROMPT_VARIANT,
        help=f"使うプロンプト variant (default: {DEFAULT_PROMPT_VARIANT})",
    )
    parser.add_argument(
        "--skip-drive", action="store_true",
        help="Drive ダウンロードをスキップ (manual_dir に既存ファイルがある前提)",
    )
    parser.add_argument(
        "--skip-generation", action="store_true",
        help="generate_slide.py 実行をスキップ (デバッグ用、Upload とスプシ書き戻しのみ)",
    )
    parser.add_argument(
        "--dry-run", action="store_true",
        help="スプシへの書き戻しを行わない (ただし他は走る)",
    )
    parser.add_argument(
        "--force", action="store_true",
        help="ステージ 3 未完了でも強行実行 (警告のみ)",
    )
    parser.add_argument(
        "--mode", choices=["first_gen", "regen"], default=None,
        help=(
            "watcher.py から渡される呼び出し種別。regen 指定時は成功で"
            "「比較確認状況=再確認待ち」、失敗で「再生成依頼」に戻す。"
            "PM が直接 CLI 実行する場合は省略可"
        ),
    )
    args = parser.parse_args()

    # sheets_client を import するために path を通す
    sys.path.insert(0, SCRIPT_DIR)

    return process_one(args)


if __name__ == "__main__":
    sys.exit(main())
