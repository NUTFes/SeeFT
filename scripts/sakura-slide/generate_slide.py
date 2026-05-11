"""Sakura AI Engine 版 解説スライド生成スクリプト

既存の scripts/generate_manual_slide.py の Sakura 版。
重要な差分: gpt-oss-120b は vision 非対応のため、画像を LLM に渡さない。
画像はファイル名のみを user prompt のリストで渡し、配置は LLM の推測に委ねる。

使い方:
  export SAKURA_API_KEY=...
  uv run --project scripts/sakura-slide python scripts/sakura-slide/generate_slide.py docs/manuals/01_44th_配線マニュアル
  uv run --project scripts/sakura-slide python scripts/sakura-slide/generate_slide.py --prompt card docs/manuals/01_44th_配線マニュアル

出力:
  default: 引数ディレクトリ配下の slide_sakura.html
  card   : 引数ディレクトリ配下の slide_sakura.card.html
"""

import argparse
import base64
import mimetypes
import os
import re
import subprocess
import sys

from openai import OpenAI


SAKURA_BASE_URL = "https://api.ai.sakura.ad.jp/v1"
MODEL = "gpt-oss-120b"
MAX_TOKENS = 16000

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
PROJECT_ROOT = os.path.normpath(os.path.join(SCRIPT_DIR, "..", ".."))

PROMPT_VARIANTS = {
    "default": "manual-prompt.md",
    "card": "manual-prompt-card.md",
}


def _load_prompt_from_md(path: str) -> tuple[str, str]:
    with open(path, "r", encoding="utf-8") as f:
        content = f.read()
    blocks = re.findall(r"```\n(.*?)```", content, re.DOTALL)
    if len(blocks) < 2:
        raise ValueError(f"Expected 2 code blocks in {path}, found {len(blocks)}")
    return blocks[0].strip(), blocks[1].strip()


def _resolve_prompt_path(variant: str) -> str:
    return os.path.join(PROJECT_ROOT, ".claude", PROMPT_VARIANTS[variant])


def resolve_manual_dir(arg: str) -> str:
    arg = arg.rstrip("/")
    if os.path.isabs(arg):
        return arg
    if os.path.isdir(arg):
        return os.path.abspath(arg)
    return os.path.join(PROJECT_ROOT, arg)


def load_source(manual_dir: str) -> tuple[str, list[str]]:
    html_file = None
    for f in os.listdir(manual_dir):
        if f.endswith(".html") and not f.startswith("slide"):
            html_file = os.path.join(manual_dir, f)
            break
    if not html_file:
        raise FileNotFoundError(f"No source HTML found in {manual_dir}")

    result = subprocess.run(
        ["pandoc", "-f", "html", "-t", "markdown", html_file],
        capture_output=True,
        text=True,
        check=True,
    )
    md_content = result.stdout
    md_content = re.sub(r'\{[^}]*style="[^"]*"[^}]*\}', "", md_content)
    md_content = re.sub(r"\{\.c\d+\}", "", md_content)

    img_dir = os.path.join(manual_dir, "images")
    image_files: list[str] = []
    if os.path.isdir(img_dir):
        image_files = sorted(os.listdir(img_dir))

    return md_content, image_files


def load_images_base64(manual_dir: str) -> dict[str, str]:
    images: dict[str, str] = {}
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


def call_sakura_api(
    system_prompt: str,
    user_prompt_template: str,
    html_content: str,
    image_files: list[str],
) -> tuple[str, dict]:
    api_key = os.environ.get("SAKURA_API_KEY")
    if not api_key:
        raise RuntimeError("SAKURA_API_KEY が未設定です")

    client = OpenAI(base_url=SAKURA_BASE_URL, api_key=api_key)

    image_list = "\n".join(f"- {f}" for f in image_files)
    user_text = user_prompt_template.format(
        image_list=image_list,
        html_content=html_content,
    )

    print(f"  Calling Sakura API ({MODEL})...")
    print(f"  Input: ~{len(user_text)//1000}K chars text, {len(image_files)} images (filename only)")

    response = client.chat.completions.create(
        model=MODEL,
        max_tokens=MAX_TOKENS,
        messages=[
            {"role": "system", "content": system_prompt},
            {"role": "user", "content": user_text},
        ],
    )

    message = response.choices[0].message
    content = message.content or ""

    usage_info = {
        "prompt_tokens": response.usage.prompt_tokens,
        "completion_tokens": response.usage.completion_tokens,
        "total_tokens": response.usage.total_tokens,
        "finish_reason": response.choices[0].finish_reason,
        "reasoning_chars": len(getattr(message, "reasoning_content", "") or ""),
    }

    print(
        f"  Response: {len(content)//1000}K chars, "
        f"finish_reason={usage_info['finish_reason']}"
    )
    print(
        f"  Usage: prompt={usage_info['prompt_tokens']}, "
        f"completion={usage_info['completion_tokens']}, "
        f"reasoning={usage_info['reasoning_chars']} chars"
    )

    return content, usage_info


def replace_placeholders(html: str, images: dict[str, str]) -> str:
    for fname, data_uri in images.items():
        html = html.replace(f"{{{{{fname}}}}}", data_uri)
    return html


def extract_html(response: str) -> str:
    match = re.search(r"```html\s*\n(.*?)```", response, re.DOTALL)
    if match:
        return match.group(1).strip()
    match = re.search(r"(<!DOCTYPE html>.*?</html>)", response, re.DOTALL | re.IGNORECASE)
    if match:
        return match.group(1).strip()
    return response.strip()


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Sakura AI Engine 版 解説スライド生成",
    )
    parser.add_argument(
        "manual_dir",
        help="マニュアルディレクトリ（絶対パス or プロジェクトルートからの相対パス）",
    )
    parser.add_argument(
        "--prompt",
        choices=list(PROMPT_VARIANTS.keys()),
        default="default",
        help="使用するプロンプトのバリアント（default=従来のスライド形式 / card=カード形式）",
    )
    args = parser.parse_args()

    manual_dir = resolve_manual_dir(args.manual_dir)
    variant = args.prompt
    output_filename = "slide_sakura.html" if variant == "default" else f"slide_sakura.{variant}.html"
    output_path = os.path.join(manual_dir, output_filename)

    prompt_path = _resolve_prompt_path(variant)
    system_prompt, user_prompt_template = _load_prompt_from_md(prompt_path)

    print("=== Sakura AI Engine 版 解説スライド生成 ===")
    print(f"  Source: {manual_dir}")
    print(f"  Prompt: {variant} ({PROMPT_VARIANTS[variant]})")
    print(f"  Output: {output_path}")

    md_content, image_files = load_source(manual_dir)
    print(f"  Markdown: {len(md_content)//1024}KB, Images: {len(image_files)} files")

    response_text, _usage = call_sakura_api(
        system_prompt,
        user_prompt_template,
        md_content,
        image_files,
    )

    slide_html = extract_html(response_text)

    images = load_images_base64(manual_dir)
    slide_html = replace_placeholders(slide_html, images)

    with open(output_path, "w", encoding="utf-8") as f:
        f.write(slide_html)

    print(f"  Output: {output_path} ({os.path.getsize(output_path)//1024}KB)")
    print("=== 完了 ===")
    return 0


if __name__ == "__main__":
    sys.exit(main())
