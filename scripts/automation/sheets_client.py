"""Google Sheets / ローカル CSV 両対応の薄いクライアント。

設計方針:
- 本番運用: Sheets API (OAuth) で 45th_マニュアル割り当て.csv 相当のスプシを読み書き
- 開発・dry-run: ローカル CSV を直接読み書き (OAuth セットアップ不要)
- 列名 ↔ A1 表記の変換は本ファイルに集約 (列順依存を呼び出し側から隠す)

OAuth セットアップ (本番モード): `docs/proposals/manual-proposal-v4-slides/automation-design.md`
Section 5 を参照。`~/.config/seeft-pipeline/credentials.json` を配置 →
初回実行でブラウザ認証 → `token.json` 自動生成。
"""

from __future__ import annotations

import csv
import os
from dataclasses import dataclass
from typing import Optional


# 25 列スキーマ。順番は CSV と完全に一致させる。
# (docs/spread_sheets/45th_マニュアル割り当て.csv の 1 行目)
COLUMNS: list[str] = [
    "区分",                  # 1 / A
    "マニュアル名",           # 2 / B
    "進捗",                  # 3 / C
    "副委員長候補確認",        # 4 / D
    "執行部確認",             # 5 / E
    "担当局",                # 6 / F
    "担当部門",              # 7 / G
    "",                     # 8 / H ← 44th 慣例の空ヘッダー
    "担当者名",              # 9 / I
    "備考",                  # 10 / J
    "39thマニュアル",         # 11 / K
    "41stマニュアル",         # 12 / L
    "42ndマニュアル",         # 13 / M
    "43rdマニュアル",         # 14 / N
    "44thマニュアル",         # 15 / O
    "メモ",                  # 16 / P
    "Google Doc URL",        # 17 / Q
    "HTML生成ステータス",     # 18 / R
    "生成HTML URL",          # 19 / S
    "最終生成日時",           # 20 / T
    "AI比較サマリー",         # 21 / U
    "比較確認状況",           # 22 / V
    "修正提案",              # 23 / W
    "配信用HTML URL",        # 24 / X
    "完成日時",              # 25 / Y
]

# Sheets API スコープ
SCOPES = [
    "https://www.googleapis.com/auth/spreadsheets",
]

# 認証ファイル配置場所 (リポジトリ外、git 事故防止)
CONFIG_DIR = os.path.expanduser("~/.config/seeft-pipeline")
CREDENTIALS_PATH = os.path.join(CONFIG_DIR, "credentials.json")
TOKEN_PATH = os.path.join(CONFIG_DIR, "token.json")

DEFAULT_SHEET_NAME = "Sheet1"  # スプシ側のシート名 (運用開始時に確認・調整)


def col_idx_to_letter(idx: int) -> str:
    """0-based 列インデックスを A1 表記に変換。0 → 'A', 25 → 'Z', 26 → 'AA'。"""
    letter = ""
    n = idx
    while True:
        letter = chr(ord("A") + (n % 26)) + letter
        n = n // 26 - 1
        if n < 0:
            break
    return letter


def column_name_to_idx(name: str) -> int:
    """列名 → 0-based インデックス。見つからなければ ValueError。"""
    try:
        return COLUMNS.index(name)
    except ValueError:
        raise ValueError(f"Unknown column: {name!r}. Known columns: {COLUMNS}")


def column_name_to_letter(name: str) -> str:
    """列名 → A1 表記の列文字。"""
    return col_idx_to_letter(column_name_to_idx(name))


@dataclass
class Row:
    """1 マニュアルの行データ。

    `data` は列名 → 値の dict。空ヘッダーの列 8 はキー `""` で格納される。
    `row_index` は 1-based のスプシ行番号 (ヘッダーが行 1、最初のデータ行が行 2)。
    """
    row_index: int
    data: dict[str, str]

    def get(self, column_name: str, default: str = "") -> str:
        return self.data.get(column_name, default)


# ---------------------------------------------------------------------------
# CSV バックエンド (dry-run / オフライン用)
# ---------------------------------------------------------------------------

class CsvBackend:
    """ローカル CSV を直接読み書きするバックエンド。

    本番運用ではなく、OAuth セットアップ前の開発・テスト用。
    """

    def __init__(self, csv_path: str):
        self.csv_path = os.path.abspath(csv_path)
        if not os.path.isfile(self.csv_path):
            raise FileNotFoundError(f"CSV not found: {self.csv_path}")

    def _load_all(self) -> list[list[str]]:
        with open(self.csv_path, newline="", encoding="utf-8") as f:
            return list(csv.reader(f))

    def _save_all(self, rows: list[list[str]]):
        with open(self.csv_path, "w", newline="", encoding="utf-8") as f:
            writer = csv.writer(f)
            for row in rows:
                writer.writerow(row)

    def read_row(self, manual_name: str) -> Optional[Row]:
        rows = self._load_all()
        if not rows:
            return None
        header = rows[0]
        name_idx = header.index("マニュアル名") if "マニュアル名" in header else 1
        for i, row in enumerate(rows[1:], start=2):
            # 行が短ければ空で埋める
            row = row + [""] * (len(header) - len(row))
            if row[name_idx] == manual_name:
                return Row(row_index=i, data=dict(zip(header, row)))
        return None

    def write_cells(self, manual_name: str, updates: dict[str, str]):
        rows = self._load_all()
        if not rows:
            raise RuntimeError("CSV is empty")
        header = rows[0]
        name_idx = header.index("マニュアル名") if "マニュアル名" in header else 1
        for i, row in enumerate(rows[1:], start=1):
            # 列幅を header に合わせる
            row = row + [""] * (len(header) - len(row))
            rows[i + 0] = row  # 念のため再代入
            if row[name_idx] == manual_name:
                for col_name, value in updates.items():
                    if col_name not in header:
                        raise ValueError(f"Unknown column: {col_name}")
                    col_idx = header.index(col_name)
                    row[col_idx] = value
                rows[i + 0] = row
                self._save_all(rows)
                return
        raise ValueError(f"Manual not found in CSV: {manual_name}")


# ---------------------------------------------------------------------------
# Sheets API バックエンド (本番運用)
# ---------------------------------------------------------------------------

class SheetsBackend:
    """Google Sheets API バックエンド。

    OAuth 認証は初回のみブラウザ起動。以降は token.json を使い回す。
    認証エラーは呼び出し側で握る (例: token 期限切れで再認証が要る場合)。
    """

    def __init__(self, spreadsheet_id: str, sheet_name: str = DEFAULT_SHEET_NAME):
        self.spreadsheet_id = spreadsheet_id
        self.sheet_name = sheet_name
        self._service = None  # lazy init: 認証を実行時まで遅延

    def _ensure_service(self):
        if self._service is not None:
            return self._service
        # google-auth-oauthlib に依存。インポートを遅延させて CSV モード単独実行を妨げない
        from google.auth.transport.requests import Request
        from google.oauth2.credentials import Credentials
        from google_auth_oauthlib.flow import InstalledAppFlow
        from googleapiclient.discovery import build

        creds = None
        if os.path.exists(TOKEN_PATH):
            creds = Credentials.from_authorized_user_file(TOKEN_PATH, SCOPES)
        if not creds or not creds.valid:
            if creds and creds.expired and creds.refresh_token:
                creds.refresh(Request())
            else:
                if not os.path.exists(CREDENTIALS_PATH):
                    raise FileNotFoundError(
                        f"OAuth credentials not found at {CREDENTIALS_PATH}. "
                        f"Google Cloud Console で credentials.json を取得し、"
                        f"{CONFIG_DIR}/ に配置してください。詳細は "
                        f"docs/proposals/manual-proposal-v4-slides/automation-design.md Section 5"
                    )
                flow = InstalledAppFlow.from_client_secrets_file(CREDENTIALS_PATH, SCOPES)
                creds = flow.run_local_server(port=0)
            os.makedirs(CONFIG_DIR, exist_ok=True)
            with open(TOKEN_PATH, "w") as token:
                token.write(creds.to_json())

        self._service = build("sheets", "v4", credentials=creds)
        return self._service

    def _read_all(self) -> list[list[str]]:
        service = self._ensure_service()
        # 列は Y まで固定し、行数は可変で取得する
        range_name = f"{self.sheet_name}!A:Y"
        result = service.spreadsheets().values().get(
            spreadsheetId=self.spreadsheet_id,
            range=range_name,
            valueRenderOption="UNFORMATTED_VALUE",
        ).execute()
        values = result.get("values", [])
        return values

    def read_row(self, manual_name: str) -> Optional[Row]:
        values = self._read_all()
        if not values:
            return None
        header = values[0]
        name_idx = header.index("マニュアル名") if "マニュアル名" in header else 1
        for i, row in enumerate(values[1:], start=2):
            row = row + [""] * (len(header) - len(row))
            if str(row[name_idx]) == manual_name:
                # 値はすべて str に揃える (Sheets API は数値型を返すことがある)
                return Row(row_index=i, data={h: str(v) for h, v in zip(header, row)})
        return None

    def write_cells(self, manual_name: str, updates: dict[str, str]):
        service = self._ensure_service()
        # 対象行を特定するために read を 1 回挟む
        row = self.read_row(manual_name)
        if row is None:
            raise ValueError(f"Manual not found in spreadsheet: {manual_name}")

        # batchUpdate で 1 リクエストにまとめる
        data = []
        for col_name, value in updates.items():
            letter = column_name_to_letter(col_name)
            range_str = f"{self.sheet_name}!{letter}{row.row_index}"
            data.append({"range": range_str, "values": [[value]]})
        body = {"valueInputOption": "USER_ENTERED", "data": data}
        service.spreadsheets().values().batchUpdate(
            spreadsheetId=self.spreadsheet_id, body=body
        ).execute()

    def find_rows_by_status(self, status_column: str, status_value: str) -> list[Row]:
        """指定列の値が特定の値である行を全て返す (watcher 用)。"""
        values = self._read_all()
        if not values:
            return []
        header = values[0]
        col_idx = header.index(status_column) if status_column in header else -1
        if col_idx < 0:
            return []
        out: list[Row] = []
        for i, row in enumerate(values[1:], start=2):
            row = row + [""] * (len(header) - len(row))
            if str(row[col_idx]) == status_value:
                out.append(Row(row_index=i, data={h: str(v) for h, v in zip(header, row)}))
        return out


# ---------------------------------------------------------------------------
# 統一インターフェース
# ---------------------------------------------------------------------------

class SheetsClient:
    """CSV / Sheets のどちらでも同じ I/F で使えるラッパー。

    使い方:
      # CSV モード (dry-run、OAuth 不要)
      client = SheetsClient(csv_path="docs/spread_sheets/45th_マニュアル割り当て.csv")

      # Sheets API モード (本番)
      client = SheetsClient(spreadsheet_id="1jz_870-...")

      row = client.read_row("配線マニュアル")
      client.write_cells("配線マニュアル", {"HTML生成ステータス": "完了"})
    """

    def __init__(
        self,
        spreadsheet_id: Optional[str] = None,
        csv_path: Optional[str] = None,
        sheet_name: str = DEFAULT_SHEET_NAME,
    ):
        if csv_path and spreadsheet_id:
            raise ValueError("csv_path と spreadsheet_id は片方だけ指定してください")
        if csv_path:
            self.backend = CsvBackend(csv_path)
            self.mode = "csv"
        elif spreadsheet_id:
            self.backend = SheetsBackend(spreadsheet_id, sheet_name=sheet_name)
            self.mode = "sheets"
        else:
            raise ValueError("csv_path か spreadsheet_id のどちらかを指定してください")

    def read_row(self, manual_name: str) -> Optional[Row]:
        return self.backend.read_row(manual_name)

    def write_cells(self, manual_name: str, updates: dict[str, str]):
        return self.backend.write_cells(manual_name, updates)

    def find_rows_by_status(self, status_column: str, status_value: str) -> list[Row]:
        if hasattr(self.backend, "find_rows_by_status"):
            return self.backend.find_rows_by_status(status_column, status_value)
        # CsvBackend にはまだないので、自前で実装
        rows: list[Row] = []
        with open(self.backend.csv_path, newline="", encoding="utf-8") as f:  # type: ignore[attr-defined]
            reader = list(csv.reader(f))
        header = reader[0]
        if status_column not in header:
            return rows
        col_idx = header.index(status_column)
        for i, row in enumerate(reader[1:], start=2):
            row = row + [""] * (len(header) - len(row))
            if row[col_idx] == status_value:
                rows.append(Row(row_index=i, data=dict(zip(header, row))))
        return rows
