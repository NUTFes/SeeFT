#!/usr/bin/env bash
# 同じマニュアルディレクトリから新旧両方のプロンプトでスライドHTMLを生成する。
# バックエンドは Sakura AI Engine 版（gpt-oss-120b）を使用する。
# 料金重視のため Anthropic API は呼ばない。Anthropic 版は scripts/generate_manual_slide.py
# を直接呼び出すこと。
#
# 前提:
#   export SAKURA_API_KEY=...
#   uv が利用可能 (scripts/sakura-slide/.venv は uv sync 済み)
#
# 使い方:
#   scripts/compare_manual_versions.sh <manual_dir>
# 例:
#   scripts/compare_manual_versions.sh docs/manuals/01_44th_配線マニュアル
#
# 出力:
#   <manual_dir>/slide_sakura.html       （manual-prompt.md      / 従来のスライド形式）
#   <manual_dir>/slide_sakura.card.html  （manual-prompt-card.md / カード形式）

set -euo pipefail

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <manual_dir>" >&2
  echo "Example: $0 docs/manuals/01_44th_配線マニュアル" >&2
  exit 1
fi

MANUAL_DIR="$1"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SAKURA_PROJECT="$SCRIPT_DIR/sakura-slide"
GENERATOR="$SAKURA_PROJECT/generate_slide.py"

if [ ! -d "$MANUAL_DIR" ]; then
  echo "Error: manual_dir not found: $MANUAL_DIR" >&2
  exit 1
fi

if [ -z "${SAKURA_API_KEY:-}" ]; then
  echo "Error: SAKURA_API_KEY is not set" >&2
  echo "  export SAKURA_API_KEY=... before running this script" >&2
  exit 1
fi

echo "=== [1/2] default プロンプト (manual-prompt.md) ==="
uv run --project "$SAKURA_PROJECT" python "$GENERATOR" "$MANUAL_DIR"

echo ""
echo "=== [2/2] card プロンプト (manual-prompt-card.md) ==="
uv run --project "$SAKURA_PROJECT" python "$GENERATOR" --prompt card "$MANUAL_DIR"

echo ""
echo "=== 完了 ==="
echo "  既存版: $MANUAL_DIR/slide_sakura.html"
echo "  新版  : $MANUAL_DIR/slide_sakura.card.html"
