"""スプシをポーリングして process_one.py を起動する常駐スクリプト。

トリガー (現状の MVP):
  1. 列 18 (HTML生成ステータス) = "実行中" かつ 列 20 (最終生成日時) 空
     → 初回生成
  2. 列 22 (比較確認状況) = "再生成依頼"
     → 再生成 (列 23 の修正提案がプロンプトに注入される)

トリガーしないこと (重要):
  - Doc URL 単独入力では発火しない (PM が明示的にステータスを変えるまで待つ)
  - 任意セル編集では発火しない (列 18, 22 だけを見る)
  - 完了済 (列 18=完了 で 最終生成日時 あり) は二重処理しない

使い方:
  # ローカル CSV モード (dry-run / 開発)
  uv run --project scripts/automation python scripts/automation/watcher.py \\
      --csv-path docs/spread_sheets/45th_マニュアル割り当て.csv \\
      --interval 60

  # Sheets API モード (本番)
  uv run --project scripts/automation python scripts/automation/watcher.py \\
      --spreadsheet-id <ID> \\
      --interval 60

  # 1 周だけスキャンして終了 (cron 等から呼ぶ用途)
  uv run --project scripts/automation python scripts/automation/watcher.py \\
      --csv-path ... --once
"""

from __future__ import annotations

import argparse
import os
import subprocess
import sys
import time


SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
PROJECT_ROOT = os.path.normpath(os.path.join(SCRIPT_DIR, "..", ".."))
PROCESS_ONE_SCRIPT = os.path.join(SCRIPT_DIR, "process_one.py")


def find_pending_rows(client) -> list[tuple[str, str]]:
    """処理対象の (マニュアル名, mode) を返す。

    mode は "first_gen" または "regen"。
    """
    pending: list[tuple[str, str]] = []

    # トリガー 1: 実行中 かつ 最終生成日時 が空
    for row in client.find_rows_by_status("HTML生成ステータス", "実行中"):
        if not row.get("最終生成日時"):
            pending.append((row.get("マニュアル名"), "first_gen"))

    # トリガー 2: 再生成依頼
    for row in client.find_rows_by_status("比較確認状況", "再生成依頼"):
        name = row.get("マニュアル名")
        if not name:
            continue
        # first_gen と重複しないように name 単位で dedupe
        if any(n == name for n, _ in pending):
            continue
        pending.append((name, "regen"))

    return pending


def run_process_one(manual_name: str, mode: str, args) -> int:
    """process_one.py を subprocess で呼ぶ。"""
    cmd = [
        sys.executable, PROCESS_ONE_SCRIPT,
        manual_name,
    ]
    if args.csv_path:
        cmd.extend(["--csv-path", args.csv_path])
    elif args.spreadsheet_id:
        cmd.extend(["--spreadsheet-id", args.spreadsheet_id])

    print(f"[{mode}] {manual_name} を処理開始...")
    proc = subprocess.run(cmd)
    if proc.returncode == 0:
        print(f"[{mode}] {manual_name} 完了 ✓")
    else:
        print(f"[{mode}] {manual_name} 失敗 (exit={proc.returncode})")
    return proc.returncode


def scan_once(client, args) -> int:
    """1 周スキャンして、検出された全マニュアルを処理。処理件数を返す。"""
    pending = find_pending_rows(client)
    if not pending:
        return 0
    print(f"処理対象: {len(pending)} 件")
    for name, mode in pending:
        try:
            run_process_one(name, mode, args)
        except Exception as e:
            print(f"ERROR processing {name}: {e}", file=sys.stderr)
    return len(pending)


def main() -> int:
    parser = argparse.ArgumentParser(
        description="スプシのステータスをポーリングして process_one.py を起動する",
    )
    group = parser.add_mutually_exclusive_group(required=True)
    group.add_argument("--csv-path", help="ローカル CSV モード")
    group.add_argument("--spreadsheet-id", help="Sheets API モード")

    parser.add_argument(
        "--interval", type=int, default=60,
        help="ポーリング間隔 (秒)。default: 60",
    )
    parser.add_argument(
        "--once", action="store_true",
        help="1 周スキャンして終了 (cron 等から定期実行する場合)",
    )
    args = parser.parse_args()

    sys.path.insert(0, SCRIPT_DIR)
    from sheets_client import SheetsClient

    client = SheetsClient(
        spreadsheet_id=args.spreadsheet_id,
        csv_path=args.csv_path,
    )
    print(f"=== watcher 起動 ({client.mode} モード, interval={args.interval}s) ===")

    if args.once:
        n = scan_once(client, args)
        print(f"=== 1 周完了 ({n} 件処理) ===")
        return 0

    # 常駐モード
    try:
        while True:
            n = scan_once(client, args)
            if n == 0:
                print(f"  処理対象なし。{args.interval}s 待機...")
            time.sleep(args.interval)
    except KeyboardInterrupt:
        print("\n=== watcher 停止 ===")
        return 0


if __name__ == "__main__":
    sys.exit(main())
