#!/usr/bin/env python3
"""
Google Doc HTML → 解説スライドHTML 生成スクリプト

使い方:
  python3 scripts/generate_manual_slide.py docs/manuals/01_44th_配線マニュアル

処理:
  1. マニュアルディレクトリの元HTML + 画像を読み込み
  2. Claude APIに送信して解説スライドHTMLを生成
  3. 画像プレースホルダーをBase64に置換して slide_v2.html を出力
"""

import anthropic
import base64
import mimetypes
import os
import re
import sys

# ─── 設定 ───
MODEL = "claude-sonnet-4-20250514"
MAX_TOKENS = 16000

# ─── プロンプト ───
SYSTEM_PROMPT = """\
あなたはマニュアルデザイナーです。大学の学園祭の運営マニュアルを、当日スマホで見る実行委員向けの解説スライドHTMLに再構成します。

## 出力形式（厳守）

- 自己完結したHTML（<!DOCTYPE html> から </html> まで）のみを出力。説明文やコードブロック記法は不要。
- 画像は `{{ファイル名}}` のプレースホルダーで参照する（例: `<img src="{{image1.png}}">`）。後処理でBase64 data URIに置換される。
- CSSはHTML内の <style> に記述する。外部ファイル参照禁止。

## スタイル仕様（厳守）

- カラースキーム:
  - ゴールド系: #C89932（メインアクセント）, #D8B660（テーブルヘッダ等）, #FFD966（ハイライト）
  - 背景: #E0DACF（ウォームベージュ）
  - 警告・強調: #96514D（ブラウニッシュレッド）
  - テキスト: #000000（メイン）, #595959（サブ）
  - カード・ボックス背景: #FFFFFF
- フォント: font-family: system-ui, sans-serif
- 各セクション: min-height: 100svh; scroll-snap-align: start
- html: scroll-snap-type: y proximity
- スマホ最適化: font-size は clamp() で最小15px以上を保証。padding は vw 単位。

## 除外ルール（厳守）

- 目次ページ（ページ番号の羅列）は除外する
- 注釈（[a], [b], [c] 等）は除外する

## 必須UI機能（厳守）

### 1. 目次（アンカーリンク付き）
- `<body>` の直後、最初のセクションの前に目次セクションを配置する
- 各章に `id="section-1"` 等のIDを付与し、目次から `<a href="#section-1">` でジャンプできるようにする
- 目次セクション自体にも `id="toc"` を付与する

### 2. 章ごとの折りたたみトグル
- 各章のコンテンツ部分を `<details>/<summary>` で囲み、開閉できるようにする
- `<details open>` にして最初から開いた状態にする
- `summary` には章タイトルを入れ、スタイルを `.section-title` に合わせる

### 3. 画像の拡大表示（ライトボックス）
- 画像タップで全画面拡大できるよう、純粋なCSS+最小限のJavaScriptで実装する
- 実装: 画像クリックで固定オーバーレイ（半透明黒背景）に拡大表示、クリックで閉じる
- `<script>` はHTML末尾に1箇所まとめて記述する

### 4. カード型レイアウト
- 全ての情報ブロック（リスト、テーブル、説明文のまとまり）をカード（白背景、角丸、影、ゴールドの枠線）で囲む
- カードは `padding: clamp(1rem, 3vw, 2rem)` で余白を確保する

## 必須コンテンツ（厳守）

### 表紙セクション（最初のセクション）
- マニュアル名、サブタイトル（担当シフト名）を大きく表示する
- 元ドキュメントに担当者名・作成者名の記載があれば、表紙に明記する

### 緊急連絡先セクション（最後のセクション）
- 「代表連絡先」「緊急連絡先」のセクションがあれば、マニュアルの**最後**に配置する
- 電話番号は `<a href="tel:090-xxxx-xxxx">090-xxxx-xxxx</a>` 形式でタップ発信できるようにする
- 責任者名と役職（部門長等）を明記する

### 設営・運営・片付けの切り分け
- マニュアルに設営（準備）・運営（当日）・片付け（撤収）のフェーズが含まれる場合、それぞれを**独立したセクション群**として切り分ける
- 各フェーズの先頭に「設営」「運営」「片付け」のフェーズ見出しを置く
- 大規模企画（役割が複数ある場合）は役割ごとにさらに分割してよい

### カンバン形式（役割分担があるとき）
- 役割・担当者の情報がある場合は、カンバンカード形式（横並びカード）で表示する
- 各カードに「役割名」「人数」「場所」「やること」を明示する
- CSS: `display: grid; grid-template-columns: repeat(auto-fit, minmax(160px, 1fr)); gap: 1em;`

## コンテンツ設計指針（判断はあなたに任せる）

- このマニュアルを読む人は、当日初めてその作業をする人。「何を」「どこで」「いつ」「どの順で」が一目でわかるように。
- 情報の優先順位: 基本情報（場所・時間・人数）→ タスク・手順 → スケジュール → 役割分担 → 必要物品 の順を基本とする。
- 元の文書の構造に縛られず、読み手視点で再構成してよい。
- 1画面（セクション）の情報密度が薄すぎず、はみ出しもしないように。
- 同一マニュアル内で同じ内容のセクションを重複させない。
- 画像は関連する説明の直後に配置し、キャプションをつける。
- 大見出しと小見出しのスタイルに明確な差をつける（色帯 vs 下線 等）。
- 表が大きすぎる場合はカテゴリに分割する。
- 元ドキュメントにURLがあれば必ずリンクとして反映する。Google Docs経由のリダイレクトURL（`https://www.google.com/url?q=...`）は、qパラメータの実URLを抽出して使用する。リンクは `href="#"` にせず、必ず実際のURLを設定すること。"""

USER_PROMPT_TEMPLATE = """\
以下はGoogle Documentからエクスポートされたマニュアルの HTML です。
この内容を読み取り、解説スライドHTMLを生成してください。

画像ファイルは以下が利用可能です:
{image_list}

元HTML:
{html_content}"""


def load_source(manual_dir: str):
    """元HTMLと画像一覧を読み込む"""
    # HTML ファイルを探す
    html_file = None
    for f in os.listdir(manual_dir):
        if f.endswith(".html") and not f.startswith("slide"):
            html_file = os.path.join(manual_dir, f)
            break
    if not html_file:
        raise FileNotFoundError(f"No source HTML found in {manual_dir}")

    with open(html_file, "r", encoding="utf-8") as fh:
        html_content = fh.read()

    # 画像一覧
    img_dir = os.path.join(manual_dir, "images")
    image_files = []
    if os.path.isdir(img_dir):
        image_files = sorted(os.listdir(img_dir))

    return html_content, image_files


def load_images_base64(manual_dir: str) -> dict:
    """画像をBase64 data URIとして読み込む"""
    images = {}
    img_dir = os.path.join(manual_dir, "images")
    if not os.path.isdir(img_dir):
        return images
    for f in sorted(os.listdir(img_dir)):
        path = os.path.join(img_dir, f)
        mime = mimetypes.guess_type(path)[0] or "image/png"
        with open(path, "rb") as fh:
            b64 = base64.b64encode(fh.read()).decode()
        images[f] = f"data:{mime};base64,{b64}"
    return images


def build_image_content_blocks(manual_dir: str, image_files: list) -> list:
    """Claude APIに送る画像コンテンツブロックを構築（最大10枚、大きすぎるものはスキップ）"""
    blocks = []
    img_dir = os.path.join(manual_dir, "images")
    for f in image_files:
        path = os.path.join(img_dir, f)
        size = os.path.getsize(path)
        # 500KB以上の画像はスキップ（API制限対策）
        if size > 500_000:
            continue
        mime = mimetypes.guess_type(path)[0] or "image/png"
        with open(path, "rb") as fh:
            b64 = base64.b64encode(fh.read()).decode()
        blocks.append({
            "type": "text",
            "text": f"[画像: {f}]"
        })
        blocks.append({
            "type": "image",
            "source": {
                "type": "base64",
                "media_type": mime,
                "data": b64,
            }
        })
        # API制限: 最大20画像ブロック
        if len([b for b in blocks if b["type"] == "image"]) >= 20:
            break
    return blocks


def call_claude_api(html_content: str, image_files: list, image_blocks: list) -> str:
    """Claude APIを呼び出して解説スライドHTMLを生成"""
    client = anthropic.Anthropic()

    image_list = "\n".join(f"- {f}" for f in image_files)
    user_text = USER_PROMPT_TEMPLATE.format(
        image_list=image_list,
        html_content=html_content,
    )

    # コンテンツブロック: テキスト + 画像
    content = [{"type": "text", "text": user_text}]
    content.extend(image_blocks)

    print(f"  Calling Claude API ({MODEL})...")
    print(f"  Input: ~{len(user_text)//1000}K chars text + {len([b for b in image_blocks if b['type'] == 'image'])} images")

    message = client.messages.create(
        model=MODEL,
        max_tokens=MAX_TOKENS,
        system=SYSTEM_PROMPT,
        messages=[{"role": "user", "content": content}],
    )

    # レスポンスからテキストを抽出
    result = ""
    for block in message.content:
        if block.type == "text":
            result += block.text

    print(f"  Response: {len(result)//1000}K chars, stop_reason={message.stop_reason}")
    print(f"  Usage: input={message.usage.input_tokens}, output={message.usage.output_tokens}")

    return result


def replace_placeholders(html: str, images: dict) -> str:
    """{{filename}} プレースホルダーをBase64 data URIに置換"""
    for fname, data_uri in images.items():
        html = html.replace(f"{{{{{fname}}}}}", data_uri)
    return html


def extract_html(response: str) -> str:
    """レスポンスからHTML部分を抽出（コードブロック記法があれば除去）"""
    # ```html ... ``` で囲まれている場合
    match = re.search(r"```html\s*\n(.*?)```", response, re.DOTALL)
    if match:
        return match.group(1).strip()
    # <!DOCTYPE から始まる場合
    match = re.search(r"(<!DOCTYPE html>.*?</html>)", response, re.DOTALL | re.IGNORECASE)
    if match:
        return match.group(1).strip()
    return response.strip()


def main():
    if len(sys.argv) < 2:
        print("Usage: python3 scripts/generate_manual_slide.py <manual_dir>")
        print("Example: python3 scripts/generate_manual_slide.py docs/manuals/01_44th_配線マニュアル")
        sys.exit(1)

    manual_dir = sys.argv[1].rstrip("/")
    output_path = os.path.join(manual_dir, "slide_api.html")

    print(f"=== 解説スライド生成 ===")
    print(f"  Source: {manual_dir}")

    # 1. ソース読み込み
    html_content, image_files = load_source(manual_dir)
    print(f"  HTML: {len(html_content)//1024}KB, Images: {len(image_files)} files")

    # 2. 画像ブロック構築（API送信用）
    image_blocks = build_image_content_blocks(manual_dir, image_files)

    # 3. Claude API 呼び出し
    response = call_claude_api(html_content, image_files, image_blocks)

    # 4. HTML抽出
    slide_html = extract_html(response)

    # 5. プレースホルダーをBase64に置換
    images = load_images_base64(manual_dir)
    slide_html = replace_placeholders(slide_html, images)

    # 6. 出力
    with open(output_path, "w", encoding="utf-8") as f:
        f.write(slide_html)

    print(f"  Output: {output_path} ({os.path.getsize(output_path)//1024}KB)")
    print(f"=== 完了 ===")


if __name__ == "__main__":
    main()
