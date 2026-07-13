"""Google Drive API ラッパー: Google Doc URL → ローカル展開済 HTML + images/。

設計方針:
- Drive Export API で「Web ページ (zip)」形式を取得 = HTML 1 個 + images/ サブフォルダ
- これを generate_slide.py の入力規約 (manual_dir に .html と images/) と互換にする
- 既存の `docs/manuals/<name>/` ディレクトリ構造をそのまま生成可能

OAuth は sheets_client.py と同じトークン (CONFIG_DIR/token.json) を共有。
スコープに drive.readonly を追加する必要があるので、初回認証は両スコープで通す。
"""

from __future__ import annotations

import io
import os
import re
import zipfile
from typing import Optional
from urllib.parse import urlparse


# Drive (読み取り) + Sheets スコープを同じ token で扱う
DRIVE_SCOPES = [
    "https://www.googleapis.com/auth/spreadsheets",
    "https://www.googleapis.com/auth/drive.readonly",
]

CONFIG_DIR = os.path.expanduser("~/.config/seeft-pipeline")
CREDENTIALS_PATH = os.path.join(CONFIG_DIR, "credentials.json")
TOKEN_PATH = os.path.join(CONFIG_DIR, "token.json")


DOC_ID_PATTERNS = [
    re.compile(r"/document/d/([a-zA-Z0-9_-]+)"),   # 通常の Doc URL
    re.compile(r"/spreadsheets/d/([a-zA-Z0-9_-]+)"),  # 万一スプシ URL が渡された場合
    re.compile(r"/file/d/([a-zA-Z0-9_-]+)"),       # Drive 共有 URL
    re.compile(r"id=([a-zA-Z0-9_-]+)"),             # クエリパラメータ形式
]


def extract_doc_id(url: str) -> str:
    """Google Doc URL から ID を抽出。

    対応する URL 形式:
      https://docs.google.com/document/d/<ID>/edit
      https://docs.google.com/document/d/<ID>/preview
      https://docs.google.com/document/d/<ID>/view?usp=sharing
      https://drive.google.com/file/d/<ID>/view
    """
    if not url:
        raise ValueError("Empty URL")
    for pattern in DOC_ID_PATTERNS:
        m = pattern.search(url)
        if m:
            return m.group(1)
    raise ValueError(f"Could not extract Doc ID from URL: {url}")


def _ensure_drive_service():
    """Drive API service を作成 (sheets_client と同じ token を共有)。"""
    from google.auth.transport.requests import Request
    from google.oauth2.credentials import Credentials
    from google_auth_oauthlib.flow import InstalledAppFlow
    from googleapiclient.discovery import build

    creds = None
    if os.path.exists(TOKEN_PATH):
        creds = Credentials.from_authorized_user_file(TOKEN_PATH, DRIVE_SCOPES)
        if creds and not creds.has_scopes(DRIVE_SCOPES):
            creds = None
    if not creds or not creds.valid:
        if creds and creds.expired and creds.refresh_token:
            creds.refresh(Request())
        else:
            if not os.path.exists(CREDENTIALS_PATH):
                raise FileNotFoundError(
                    f"OAuth credentials not found at {CREDENTIALS_PATH}. "
                    f"Google Cloud Console で credentials.json を取得し、"
                    f"{CONFIG_DIR}/ に配置してください。"
                )
            flow = InstalledAppFlow.from_client_secrets_file(CREDENTIALS_PATH, DRIVE_SCOPES)
            creds = flow.run_local_server(port=0)
        os.makedirs(CONFIG_DIR, exist_ok=True)
        with open(TOKEN_PATH, "w") as token:
            token.write(creds.to_json())
        os.chmod(TOKEN_PATH, 0o600)

    return build("drive", "v3", credentials=creds)


def export_doc_as_html_zip(doc_id: str) -> bytes:
    """Google Doc を 「Web ページ (zip)」 形式でエクスポートして bytes を返す。"""
    service = _ensure_drive_service()
    request = service.files().export_media(
        fileId=doc_id,
        mimeType="application/zip",
    )
    fh = io.BytesIO()
    from googleapiclient.http import MediaIoBaseDownload
    downloader = MediaIoBaseDownload(fh, request)
    done = False
    while not done:
        _status, done = downloader.next_chunk()
    fh.seek(0)
    return fh.read()


def unpack_html_zip_to_manual_dir(zip_bytes: bytes, manual_dir: str) -> str:
    """ZIP を展開し、manual_dir の中に .html と images/ を配置する。

    展開ルール:
      ZIP 内の最上位 .html を `<manual_dir>/source.html` にリネーム
      ZIP 内の `images/` 配下を `<manual_dir>/images/` に配置
      その他のファイル (style.css 等) は破棄

    戻り値: 展開された HTML の絶対パス
    """
    os.makedirs(manual_dir, exist_ok=True)
    images_dir = os.path.join(manual_dir, "images")
    os.makedirs(images_dir, exist_ok=True)

    html_path: Optional[str] = None
    with zipfile.ZipFile(io.BytesIO(zip_bytes)) as zf:
        for name in zf.namelist():
            # ディレクトリエントリはスキップ
            if name.endswith("/"):
                continue
            data = zf.read(name)
            basename = os.path.basename(name)
            if name.endswith(".html") and html_path is None:
                # 最上位 (or 最初に見つかった) HTML を source.html として保存
                html_path = os.path.join(manual_dir, "source.html")
                with open(html_path, "wb") as f:
                    f.write(data)
            elif "/images/" in name or name.startswith("images/"):
                with open(os.path.join(images_dir, basename), "wb") as f:
                    f.write(data)
            # それ以外 (style.css 等) は無視

    if html_path is None:
        raise RuntimeError("ZIP に .html が含まれていませんでした")

    return html_path


def download_doc_to_manual_dir(doc_url: str, manual_dir: str) -> str:
    """Doc URL を受け取り、manual_dir 配下に HTML + images/ を展開する。

    戻り値: 展開された HTML の絶対パス。
    """
    doc_id = extract_doc_id(doc_url)
    zip_bytes = export_doc_as_html_zip(doc_id)
    return unpack_html_zip_to_manual_dir(zip_bytes, manual_dir)
