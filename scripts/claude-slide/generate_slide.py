"""Claude Agent SDK 版 解説スライド生成スクリプト

サブスク経由（Claude Code OAuth）で Opus 4.7 を使う。Anthropic API キー不要、
Sakura API キーも不要。`claude login` 済みのサブスク契約で動く。

scripts/generate_manual_slide.py (Anthropic API版) と
scripts/sakura-slide/generate_slide.py (Sakura版) の Claude Agent SDK 版。

使い方:
  uv run --project scripts/claude-slide python scripts/claude-slide/generate_slide.py docs/manuals/01_44th_配線マニュアル
  uv run --project scripts/claude-slide python scripts/claude-slide/generate_slide.py --prompt card docs/manuals/01_44th_配線マニュアル

出力:
  default: 引数ディレクトリ配下の slide_claude.html
  card   : 引数ディレクトリ配下の slide_claude.card.html
"""

import anyio
import argparse
import base64
import mimetypes
import os
import re
import subprocess
import sys

from claude_agent_sdk import (
    AssistantMessage,
    ClaudeAgentOptions,
    ResultMessage,
    TextBlock,
    query,
)


SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
PROJECT_ROOT = os.path.normpath(os.path.join(SCRIPT_DIR, "..", ".."))

PROMPT_VARIANTS = {
    "default": "manual-prompt.md",
    "card": "manual-prompt-card.md",
    # 文章不変ポリシー版。機械検証 (verify_slide_mechanical.py) と組み合わせて使う前提
    "card-strict": "manual-prompt-card-strict.md",
}

# 単発生成: ツール使用は完全に禁止し、テキスト出力のみさせる。
# 画像は LLM には渡さない（プロンプトにファイル名のみ）。
# Markdown 文脈と figcaption から配置を推測させる戦略。
DISALLOWED_TOOLS = [
    "Read", "Write", "Edit", "Bash",
    "Task", "WebFetch", "WebSearch",
    "Grep", "Glob", "TodoWrite",
    "NotebookEdit",
]


def _load_prompt_from_md(path: str) -> tuple[str, str]:
    with open(path, "r", encoding="utf-8") as f:
        content = f.read()
    # Prefer 4-backtick outer fences (allows nested 3-backtick blocks inside).
    # Falls back to 3-backtick for prompts without inner code blocks.
    blocks = re.findall(r"````\s*\n(.*?)\n````", content, re.DOTALL)
    if len(blocks) < 2:
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
    candidates = sorted(
        f for f in os.listdir(manual_dir)
        if f.endswith(".html")
        and not f.startswith("slide")
        and not f.startswith("verify")
    )
    if not candidates:
        raise FileNotFoundError(f"No source HTML found in {manual_dir}")
    if "source.html" in candidates:
        html_file = os.path.join(manual_dir, "source.html")
    elif len(candidates) == 1:
        html_file = os.path.join(manual_dir, candidates[0])
    else:
        raise RuntimeError(
            f"Ambiguous source HTML in {manual_dir}: {candidates}. "
            "Keep exactly one source HTML (or name it source.html)."
        )

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


def _resolve_image_src(src: str, images: dict[str, str]) -> str | None:
    """`<img src="...">` の src 値がローカル画像を指すなら data URI を返す。違えば None。

    card-strict では LLM が `{{fname}}` を守らず、素のファイル名 (`image1.jpg`) や
    相対パス (`images/image1.png`, `./images/image1.png`) で src を書くことがある。
    その揺れを決定的に吸収して必ず base64 埋め込みする保険。

    引数:
      src    : img タグの生の src 値（例: "image1.jpg" / "images/image1.png" /
               "./images/image1.png" / "data:image/png;base64,..."）
      images : {ファイル名 -> data URI} の辞書（load_images_base64 の返り値）

    注意: 既に `data:` 化済みの src は再処理しても None を返すこと（冪等性）。
    """
    if src.startswith("data:"):
        return None  # 既に埋め込み済み。再処理しても壊さない（冪等）
    return images.get(os.path.basename(src))


def replace_placeholders(html: str, images: dict[str, str]) -> str:
    # 1) {{fname}} プレースホルダー（プロンプトが指示する正規ルート）
    for fname, data_uri in images.items():
        html = html.replace(f"{{{{{fname}}}}}", data_uri)

    # 2) 素のファイル名・相対パスで書かれた <img src> も決定的に base64 化する保険。
    #    LLM が {{}} を守らなかった場合（card-strict で多発）でも画像が必ず埋まる。
    def _embed(m: "re.Match[str]") -> str:
        quote = m.group("q")
        data_uri = _resolve_image_src(m.group("src"), images)
        if data_uri is None:
            return m.group(0)
        return f"src={quote}{data_uri}{quote}"

    html = re.sub(r'src=(?P<q>["\'])(?P<src>[^"\']*)(?P=q)', _embed, html)
    return html


def extract_html(response: str) -> str:
    # 非貪欲 (.*?) だと、生成 HTML の中に ``` が現れた時点で切れる（CSS コメントや
    # コードサンプル、モデルが入れ子でフェンスを書いた場合）。プロンプト自体が html と
    # css の 2 つのフェンスを例示しているため、モデルが真似る余地は十分にある。
    # 「最初の ```html から最後の ```」を貪欲に取り、その中の最初の <!DOCTYPE〜最後の
    # </html> をさらに切り出すことで、前後の説明文や入れ子フェンスの影響を受けなくする。
    match = re.search(r"```html\s*\n(.*)```", response, re.DOTALL)
    body = match.group(1) if match else response
    match = re.search(r"(<!DOCTYPE html>.*</html>)", body, re.DOTALL | re.IGNORECASE)
    if match:
        return match.group(1).strip()
    return body.strip()


def is_complete_html(html: str) -> bool:
    # 応答が途中で切れていても SDK は is_error=False を返す。構造的に閉じているかで判定する
    return html.rstrip().lower().endswith("</html>")


async def call_claude_sdk(
    system_prompt: str,
    user_prompt_template: str,
    md_content: str,
    image_files: list[str],
    model: str | None,
    extra_instructions: str = "",
) -> tuple[str, dict]:
    image_list = "\n".join(f"- {f}" for f in image_files)
    user_text = user_prompt_template.format(
        image_list=image_list,
        html_content=md_content,
    )
    # 部門長レビューの修正案など、この回だけの追加指示を末尾に注入する。
    # 共有プロンプトを書き換えずにマニュアル単位で再生成を回すための仕組み。
    if extra_instructions.strip():
        user_text += (
            "\n\n# この回の追加・修正指示（前回レビューの反映。最優先で従うこと）\n"
            + extra_instructions.strip()
        )

    # max_turns を 20 に: 大きい入力（55KB+ markdown）で Claude が tool 試行する場合に
    # disallowed_tools 拒否で turn が消費されるため、余裕大きめ。
    # お化け屋敷で max_turns=10 では flaky に失敗する事象を確認したため。
    options_kwargs: dict = {
        "system_prompt": system_prompt,
        "max_turns": 20,
        "disallowed_tools": DISALLOWED_TOOLS,
    }
    if model:
        options_kwargs["model"] = model

    options = ClaudeAgentOptions(**options_kwargs)

    print("  Calling Claude via Agent SDK (subscription)...")
    print(f"  Model: {model or '(Claude Code default)'}")
    print(f"  Input: ~{len(user_text)//1000}K chars text, {len(image_files)} images (filename only)")

    result_text = ""
    usage: dict = {}

    async for message in query(prompt=user_text, options=options):
        if isinstance(message, AssistantMessage):
            for block in message.content:
                if isinstance(block, TextBlock):
                    result_text += block.text
        elif isinstance(message, ResultMessage):
            usage = {
                "duration_ms": getattr(message, "duration_ms", None),
                "num_turns": getattr(message, "num_turns", None),
                "total_cost_usd": getattr(message, "total_cost_usd", None),
                "is_error": getattr(message, "is_error", None),
            }

    print(f"  Response: {len(result_text)//1000}K chars")
    if usage:
        print(f"  Usage: {usage}")

    return result_text, usage


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Claude Agent SDK 版 解説スライド生成 (subscription auth, no API key)",
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
    parser.add_argument(
        "--model",
        default=None,
        help="使用するモデル（例: claude-opus-4-7, claude-sonnet-4-6）。未指定なら Claude Code のデフォルト",
    )
    parser.add_argument(
        "--embed-only",
        action="store_true",
        help="LLM を呼ばず、既存の出力 HTML に画像を base64 再埋め込みするだけ（決定的・再生成不要）",
    )
    args = parser.parse_args()

    manual_dir = resolve_manual_dir(args.manual_dir)
    variant = args.prompt
    output_filename = "slide_claude.html" if variant == "default" else f"slide_claude.{variant}.html"
    output_path = os.path.join(manual_dir, output_filename)

    # 既存 HTML への画像再埋め込みのみ（壊れた card-strict の決定的な復旧用）
    if args.embed_only:
        if not os.path.isfile(output_path):
            print(f"  ERROR: 出力 HTML が見つかりません: {output_path}", file=sys.stderr)
            return 1
        with open(output_path, "r", encoding="utf-8") as f:
            slide_html = f.read()
        before = len(slide_html)
        slide_html = replace_placeholders(slide_html, load_images_base64(manual_dir))
        with open(output_path, "w", encoding="utf-8") as f:
            f.write(slide_html)
        print("=== 画像再埋め込み（--embed-only）===")
        print(f"  {output_path}: {before//1024}KB → {os.path.getsize(output_path)//1024}KB")
        return 0

    prompt_path = _resolve_prompt_path(variant)
    system_prompt, user_prompt_template = _load_prompt_from_md(prompt_path)

    print("=== Claude Agent SDK 版 解説スライド生成 ===")
    print(f"  Source: {manual_dir}")
    print(f"  Prompt: {variant} ({PROMPT_VARIANTS[variant]})")
    print(f"  Output: {output_path}")

    md_content, image_files = load_source(manual_dir)
    print(f"  Markdown: {len(md_content)//1024}KB, Images: {len(image_files)} files")

    # マニュアルごとの追加・修正指示（部門長レビューの反映用）。あれば末尾に注入する。
    extra_instructions = ""
    instr_path = os.path.join(manual_dir, "instructions.md")
    if os.path.isfile(instr_path):
        with open(instr_path, "r", encoding="utf-8") as f:
            extra_instructions = f.read()
        print(f"  追加指示: {instr_path} ({len(extra_instructions)} chars)")

    response_text, usage = anyio.run(
        call_claude_sdk,
        system_prompt,
        user_prompt_template,
        md_content,
        image_files,
        args.model,
        extra_instructions,
    )
    if usage.get("is_error"):
        print(f"  ERROR: Claude Agent SDK returned is_error=True: {usage}", file=sys.stderr)
        return 1
    if not response_text.strip():
        print("  ERROR: Empty response from Claude Agent SDK", file=sys.stderr)
        return 1

    # 抽出の前に生の応答を残す。抽出や埋め込みに失敗しても、10 分かけて再生成せず
    # ここから復旧できる（LLM 呼び出しは非決定的で、再生成すると別の HTML になる）
    raw_path = output_path + ".raw.md"
    with open(raw_path, "w", encoding="utf-8") as f:
        f.write(response_text)
    print(f"  Raw response: {raw_path} ({os.path.getsize(raw_path)//1024}KB)")

    slide_html = extract_html(response_text)
    if not is_complete_html(slide_html):
        # SDK は途中で切れた応答も is_error=False で返すため、ここで止める。
        # 生の応答は上で保存済みなので、原因を見て extract_html を直せば再生成なしで済む
        print(
            f"  ERROR: 抽出した HTML が </html> で閉じていません（{len(slide_html)//1024}KB）。"
            f" 生の応答 {raw_path} を確認してください",
            file=sys.stderr,
        )
        return 1

    images = load_images_base64(manual_dir)
    slide_html = replace_placeholders(slide_html, images)

    with open(output_path, "w", encoding="utf-8") as f:
        f.write(slide_html)

    print(f"  Output: {output_path} ({os.path.getsize(output_path)//1024}KB)")
    print("=== 完了 ===")
    return 0


if __name__ == "__main__":
    sys.exit(main())
