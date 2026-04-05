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
import subprocess
import sys

# ─── 設定 ───
MODEL = "claude-sonnet-4-20250514"
MAX_TOKENS = 16000

# ─── プロンプト（.claude/manual-prompt.md から読み込み） ───
PROMPT_MD_PATH = os.path.join(os.path.dirname(__file__), "..", ".claude", "manual-prompt.md")


def _load_prompt_from_md(path: str) -> tuple[str, str]:
    """manual-prompt.md から SYSTEM_PROMPT と USER_PROMPT_TEMPLATE を抽出"""
    with open(path, "r", encoding="utf-8") as f:
        content = f.read()
    # ```...``` で囲まれたコードブロックを順に抽出
    blocks = re.findall(r"```\n(.*?)```", content, re.DOTALL)
    if len(blocks) < 2:
        raise ValueError(f"Expected 2 code blocks in {path}, found {len(blocks)}")
    return blocks[0].strip(), blocks[1].strip()


SYSTEM_PROMPT, USER_PROMPT_TEMPLATE = _load_prompt_from_md(PROMPT_MD_PATH)


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

    # HTML → MD変換 + Google Docsノイズ除去
    result = subprocess.run(
        ["pandoc", "-f", "html", "-t", "markdown", html_file],
        capture_output=True, 
        text=True
    )
    md_content = result.stdout
    md_content = re.sub(r'\{[^}]*style="[^"]*"[^}]*\}', '', md_content)
    md_content = re.sub(r'\{\.c\d+\}', '', md_content)
    html_content = md_content

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
