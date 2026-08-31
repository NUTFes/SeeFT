#!/usr/bin/env bash
# auto.md の①〜④（Googleドキュメントのzip → 解説HTML → 検証 → アップロード）を
# 通しで実行するラッパー。⑤（シフトスプシへの紐付け）はスプシ・GAS側の作業なので
# 対象外。最後に④の出力（対応表に貼る1行）をそのまま表示する。
#
# 使い方:
#   scripts/automation/run_pipeline.sh <zipのパス> --id <公開ID> --doc-url <ドキュメントURL> [--prompt card-strict] [--model claude-opus-4-7]
#
# 例:
#   scripts/automation/run_pipeline.sh ~/Downloads/45th_企画マニュアル_縁日.zip \
#     --id en-nichi \
#     --doc-url "https://docs.google.com/document/d/xxxx/edit"
#
# ③（中身の確認）だけは目視が要るため、確認用にブラウザで開いたあと
# 続行してよいか一度だけ確認を挟む。 --yes を付けるとその確認をスキップする。

set -euo pipefail

PROMPT_VARIANT="card-strict"
MODEL="claude-opus-4-7"
MANUAL_ID=""
DOC_URL=""
ZIP_PATH=""
ASSUME_YES=0

usage() {
  echo "使い方: $0 <zipのパス> --id <公開ID> --doc-url <ドキュメントURL> [--prompt <variant>] [--model <モデル>] [--yes]" >&2
  exit 1
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --id) MANUAL_ID="$2"; shift 2 ;;
    --doc-url) DOC_URL="$2"; shift 2 ;;
    --prompt) PROMPT_VARIANT="$2"; shift 2 ;;
    --model) MODEL="$2"; shift 2 ;;
    --yes) ASSUME_YES=1; shift ;;
    -h|--help) usage ;;
    *)
      if [[ -z "$ZIP_PATH" ]]; then
        ZIP_PATH="$1"; shift
      else
        echo "不明な引数: $1" >&2
        usage
      fi
      ;;
  esac
done

[[ -z "$ZIP_PATH" ]] && usage
[[ -z "$MANUAL_ID" ]] && { echo "ERROR: --id は必須です" >&2; usage; }

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

for cmd in python3 uv; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "ERROR: $cmd が見つかりません。'brew install pandoc uv' などで準備してください（auto.md参照）" >&2
    exit 1
  fi
done

echo "########################################"
echo "# ① Googleドキュメントを展開"
echo "########################################"
PREPARE_OUTPUT="$(python3 "$SCRIPT_DIR/prepare_manual.py" "$ZIP_PATH")"
echo "$PREPARE_OUTPUT"

MANUAL_DIR="$(printf '%s\n' "$PREPARE_OUTPUT" | sed -n 's/^  配置: //p')"
if [[ -z "$MANUAL_DIR" ]]; then
  echo "ERROR: prepare_manual.py の出力からマニュアルディレクトリを特定できませんでした" >&2
  exit 1
fi
MANUAL_NAME="$(basename "$MANUAL_DIR")"

echo
echo "  → マニュアル名: $MANUAL_NAME"
echo "  この名前がシフトスプシ「タスク一覧」M列の値と一致しているか、後で必ず確認してください。"

if [[ "$ASSUME_YES" -ne 1 ]]; then
  read -r -p "続行しますか？ [y/N] " reply
  [[ "$reply" =~ ^[Yy]$ ]] || { echo "中断しました"; exit 1; }
fi

echo
echo "########################################"
echo "# ② 解説HTMLに変換（${PROMPT_VARIANT} / ${MODEL}）"
echo "########################################"
uv run --project "$PROJECT_ROOT/scripts/claude-slide" python "$PROJECT_ROOT/scripts/claude-slide/generate_slide.py" \
  --prompt "$PROMPT_VARIANT" --model "$MODEL" "$MANUAL_DIR"

echo
echo "########################################"
echo "# ③ 中身を確認"
echo "########################################"
uv run --project "$PROJECT_ROOT/scripts/claude-slide" python "$PROJECT_ROOT/scripts/claude-slide/verify_slide_mechanical.py" \
  "$MANUAL_DIR"

SLIDE_HTML="$MANUAL_DIR/slide_claude.${PROMPT_VARIANT}.html"
if [[ "$PROMPT_VARIANT" == "default" ]]; then
  SLIDE_HTML="$MANUAL_DIR/slide_claude.html"
fi

if command -v open >/dev/null 2>&1; then
  open "$SLIDE_HTML"
fi

echo
echo "  verify_report.card-strict.md（文章差分レポート）と、開いたブラウザの見た目"
echo "  （スマホ幅・375px程度）を確認してください。"

if [[ "$ASSUME_YES" -ne 1 ]]; then
  read -r -p "内容に問題なく、アップロードへ進んでよいですか？ [y/N] " reply
  [[ "$reply" =~ ^[Yy]$ ]] || { echo "中断しました。アップロードは行っていません。"; exit 1; }
fi

echo
echo "########################################"
echo "# ④ サーバーにアップロード"
echo "########################################"
python3 "$SCRIPT_DIR/upload_manual.py" \
  --id "$MANUAL_ID" \
  --doc-url "$DOC_URL" \
  "$MANUAL_DIR"

echo
echo "########################################"
echo "# ⑤ タスクへの紐づけ（手作業）"
echo "########################################"
echo "  上に表示された「マニュアルURL」シートに貼る行をコピーして、"
echo "  シフトスプシ 45th_シフト_ver0 の「マニュアルURL」シートに貼ってください。"
echo "  以降は auto.md の「⑤ タスクに紐づける」を参照してください。"
