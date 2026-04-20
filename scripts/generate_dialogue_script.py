#!/usr/bin/env python3
"""
Google Doc HTML → 対話台本JSON 生成スクリプト (Claude Code CLI サブスク経由)

使い方:
  python3 scripts/generate_dialogue_script.py docs/manuals/01_44th_配線マニュアル

処理:
  1. 元HTML を Pandoc で MD に変換
  2. claude -p 経由でプロンプトを送信 (Claude Code のサブスクリプション認証を使用)
  3. JSON Schema バリデーション付きで対話台本JSONを取得
  4. dialogue.json を出力
"""

import json
import os
import re
import subprocess
import sys

PROMPT_MD_PATH = os.path.join(os.path.dirname(__file__), "..", ".claude", "dialogue-prompt.md")
MODEL = "opus"  # Claude Code alias: 最新 Opus
TIMEOUT_SEC = 600


def _load_prompt_from_md(path: str) -> tuple[str, str]:
    """dialogue-prompt.md から SYSTEM_PROMPT と USER_PROMPT_TEMPLATE を抽出"""
    with open(path, "r", encoding="utf-8") as f:
        content = f.read()
    blocks = re.findall(r"```\n(.*?)```", content, re.DOTALL)
    if len(blocks) < 2:
        raise ValueError(f"Expected 2 code blocks in {path}, found {len(blocks)}")
    return blocks[0].strip(), blocks[1].strip()


SYSTEM_PROMPT, USER_PROMPT_TEMPLATE = _load_prompt_from_md(PROMPT_MD_PATH)


DIALOGUE_SCHEMA = {
    "type": "object",
    "properties": {
        "manual_name": {"type": "string"},
        "sections": {
            "type": "array",
            "items": {
                "type": "object",
                "properties": {
                    "section_title": {"type": "string"},
                    "estimated_duration_sec": {"type": "number"},
                    "turns": {
                        "type": "array",
                        "items": {
                            "type": "object",
                            "properties": {
                                "speaker": {"type": "string", "enum": ["Expert", "Novice"]},
                                "text": {"type": "string"},
                                "image": {"type": ["string", "null"]},
                                "tags": {
                                    "type": "array",
                                    "items": {
                                        "type": "string",
                                        "enum": [
                                            "short pause", "medium pause", "long pause",
                                            "uhm", "laughing", "whispering", "shouting", "extremely fast",
                                        ],
                                    },
                                },
                            },
                            "required": ["speaker", "text", "image", "tags"],
                        },
                    },
                },
                "required": ["section_title", "estimated_duration_sec", "turns"],
            },
        },
    },
    "required": ["manual_name", "sections"],
}


def load_source(manual_dir: str):
    """元HTML を MD に変換し、画像一覧とあわせて返す"""
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
    )
    md_content = result.stdout
    md_content = re.sub(r'\{[^}]*style="[^"]*"[^}]*\}', '', md_content)
    md_content = re.sub(r'\{\.c\d+\}', '', md_content)

    img_dir = os.path.join(manual_dir, "images")
    image_files = []
    if os.path.isdir(img_dir):
        image_files = sorted(os.listdir(img_dir))

    return md_content, image_files, img_dir


def call_claude_code(user_prompt: str, system_prompt: str, manual_dir: str) -> str:
    """claude -p でプロンプトを実行し、結果JSON文字列を返す"""
    cmd = [
        "claude", "-p",
        "--system-prompt", system_prompt,
        "--allowedTools", "Read,Glob",
        "--add-dir", os.path.abspath(manual_dir),
        "--output-format", "json",
        "--json-schema", json.dumps(DIALOGUE_SCHEMA, ensure_ascii=False),
        "--model", MODEL,
        "--permission-mode", "bypassPermissions",
        user_prompt,
    ]

    print(f"  Calling claude -p (subscription, model={MODEL})...")
    print(f"  Prompt: ~{len(user_prompt)//1000}K chars, tools=Read,Glob")

    result = subprocess.run(
        cmd,
        capture_output=True,
        text=True,
        timeout=TIMEOUT_SEC,
    )

    if result.returncode != 0:
        print(f"  stderr: {result.stderr[:2000]}")
        raise RuntimeError(f"claude -p failed with code {result.returncode}")

    outer = json.loads(result.stdout)

    debug_path = os.path.join(os.path.dirname(__file__), "_debug_last_envelope.json")
    with open(debug_path, "w", encoding="utf-8") as f:
        json.dump(outer, f, ensure_ascii=False, indent=2)

    print(f"  Duration: {outer.get('duration_ms', 0)}ms, num_turns: {outer.get('num_turns')}, stop_reason: {outer.get('stop_reason')}")
    print(f"  Usage: in={outer.get('usage', {}).get('input_tokens', '?')} out={outer.get('usage', {}).get('output_tokens', '?')}")

    if outer.get("is_error") or outer.get("subtype") != "success":
        raise RuntimeError(f"claude -p returned error: subtype={outer.get('subtype')}, is_error={outer.get('is_error')}")

    # --json-schema 指定時は structured_output に JSON 本体が入る
    structured = outer.get("structured_output")
    if structured is not None:
        return json.dumps(structured, ensure_ascii=False)

    # フォールバック: result テキストから JSON を抽出
    result_text = outer.get("result", "")
    if not result_text.strip():
        raise RuntimeError(f"Both structured_output and result are empty. See {debug_path}")
    return result_text


def extract_json(response: str) -> dict:
    """--json-schema 指定時は基本的に生JSONが返るが、念のためフォールバックを残す"""
    match = re.search(r"```json\s*\n(.*?)```", response, re.DOTALL)
    if match:
        return json.loads(match.group(1).strip())
    match = re.search(r"(\{.*\})", response, re.DOTALL)
    if match:
        return json.loads(match.group(1).strip())
    return json.loads(response.strip())


def validate_dialogue(dialogue: dict, image_files: list):
    """生成されたJSONをドメインルール観点で検証し、警告を表示"""
    allowed_tags = {
        "short pause", "medium pause", "long pause",
        "uhm", "laughing", "whispering", "shouting", "extremely fast",
    }
    allowed_speakers = {"Expert", "Novice"}
    image_set = set(image_files)

    warnings = []
    for i, section in enumerate(dialogue.get("sections", [])):
        dur = section.get("estimated_duration_sec", 0)
        if not (120 <= dur <= 300):
            warnings.append(f"section[{i}] estimated_duration_sec={dur} out of [120, 300]")
        for j, turn in enumerate(section.get("turns", [])):
            if turn.get("speaker") not in allowed_speakers:
                warnings.append(f"section[{i}].turns[{j}] invalid speaker: {turn.get('speaker')}")
            img = turn.get("image")
            if img is not None and img not in image_set:
                warnings.append(f"section[{i}].turns[{j}] image not in list: {img}")
            for tag in turn.get("tags", []):
                if tag not in allowed_tags:
                    warnings.append(f"section[{i}].turns[{j}] invalid tag: {tag}")

    if warnings:
        print(f"  ⚠ Validation warnings ({len(warnings)}):")
        for w in warnings[:10]:
            print(f"    - {w}")
        if len(warnings) > 10:
            print(f"    ... and {len(warnings) - 10} more")
    else:
        print(f"  ✓ Validation passed")


def main():
    if len(sys.argv) < 2:
        print("Usage: python3 scripts/generate_dialogue_script.py <manual_dir>")
        print("Example: python3 scripts/generate_dialogue_script.py docs/manuals/01_44th_配線マニュアル")
        sys.exit(1)

    manual_dir = sys.argv[1].rstrip("/")
    manual_name = os.path.basename(manual_dir)
    output_path = os.path.join(manual_dir, "dialogue.json")

    print(f"=== 対話台本生成 ===")
    print(f"  Source: {manual_dir}")

    md_content, image_files, img_dir = load_source(manual_dir)
    print(f"  MD: {len(md_content)//1024}KB, Images: {len(image_files)} files")

    image_list = "\n".join(f"- {f}" for f in image_files)
    user_text = USER_PROMPT_TEMPLATE.format(
        manual_name=manual_name,
        image_list=image_list,
        md_content=md_content,
    )

    abs_img_dir = os.path.abspath(img_dir)
    user_text += (
        f"\n\n---\n"
        f"作業写真の絶対ディレクトリ: {abs_img_dir}\n"
        f"各ターンで使う image を選ぶ前に、必要に応じて Read ツールで画像を確認してよい。\n"
    )

    response = call_claude_code(user_text, SYSTEM_PROMPT, manual_dir)

    dialogue = extract_json(response)
    validate_dialogue(dialogue, image_files)

    with open(output_path, "w", encoding="utf-8") as f:
        json.dump(dialogue, f, ensure_ascii=False, indent=2)

    print(f"  Output: {output_path} ({os.path.getsize(output_path)//1024}KB)")
    print(f"  Sections: {len(dialogue.get('sections', []))}")
    total_turns = sum(len(s.get("turns", [])) for s in dialogue.get("sections", []))
    print(f"  Total turns: {total_turns}")
    print(f"=== 完了 ===")


if __name__ == "__main__":
    main()
