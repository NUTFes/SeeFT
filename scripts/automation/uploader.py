"""ホスティング先への HTML アップローダ (スタブ版)。

このスタブ実装はローカル `scripts/automation/out/` にコピーして `file://` URL を返すだけ。
ホスティング先 (GitHub Pages / Cloudflare R2 / 等) が決まったら、
`upload()` 関数の中身を本実装に差し替えるだけで他のパイプライン部品は無変更で動く。

差し替え時に保つべき I/F:
  def upload(local_html_path: str, key: str) -> str
    - local_html_path: アップロード対象の HTML ファイル
    - key: 配信先での識別子 (例: "01_44th_配線マニュアル")
    - 戻り値: 部門長・確認担当者がブラウザでアクセスできる URL
"""

from __future__ import annotations

import os
import re
import shutil
from typing import Optional


SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
PROJECT_ROOT = os.path.normpath(os.path.join(SCRIPT_DIR, "..", ".."))

# スタブ版の出力先 (リポジトリ外、git に commit したくない)
DEFAULT_STAGING_DIR = os.path.join(PROJECT_ROOT, "scripts", "automation", "out")


def _safe_key(key: str) -> str:
    """key (Sheets の値、外部入力) をファイル名として安全な形に正規化する。"""
    name = os.path.basename(key).strip()
    name = re.sub(r"[^\w.-]+", "_", name, flags=re.UNICODE)
    if not name or name in {".", ".."}:
        raise ValueError(f"Invalid key: {key!r}")
    return name


def upload(local_html_path: str, key: str, staging_dir: Optional[str] = None) -> str:
    """生成済 HTML を「ホスティング先」に配置して URL を返す (スタブ実装)。

    本実装ではないため、ローカルにコピーして `file://` URL を返すだけ。
    確認担当者がブラウザで開ける真の URL を返すのが本来の責務。

    引数:
      local_html_path: アップロード対象の HTML
      key: 配信先での識別子 (例: マニュアル名)
      staging_dir: 配置先ディレクトリ (None なら DEFAULT_STAGING_DIR)
    """
    if not os.path.isfile(local_html_path):
        raise FileNotFoundError(f"Local HTML not found: {local_html_path}")

    target_dir = staging_dir or DEFAULT_STAGING_DIR
    os.makedirs(target_dir, exist_ok=True)

    # key はスプシ由来の外部入力なので、ファイル名として安全な形に正規化してから使う
    target_filename = f"{_safe_key(key)}.html"
    target_path = os.path.join(target_dir, target_filename)
    shutil.copy2(local_html_path, target_path)

    # file:// URL を返す (ローカル開発・dry-run 用)
    abs_path = os.path.abspath(target_path)
    return f"file://{abs_path}"


def upload_text(local_text_path: str, key: str, staging_dir: Optional[str] = None) -> str:
    """検証レポート等のテキストファイルを配置して URL を返す。

    現在は使われていないが、将来「列 22 を再導入したい」となった場合の拡張余地として残す。
    """
    if not os.path.isfile(local_text_path):
        raise FileNotFoundError(f"Local text not found: {local_text_path}")

    target_dir = staging_dir or DEFAULT_STAGING_DIR
    os.makedirs(target_dir, exist_ok=True)

    target_path = os.path.join(target_dir, f"{_safe_key(key)}.txt")
    shutil.copy2(local_text_path, target_path)

    abs_path = os.path.abspath(target_path)
    return f"file://{abs_path}"
