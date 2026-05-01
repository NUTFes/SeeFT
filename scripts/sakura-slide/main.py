"""Sakura AI Engine 疎通確認スクリプト

実行:
  export SAKURA_API_KEY=...
  uv run --directory scripts/sakura-slide python main.py
"""

import os
import sys

from openai import OpenAI


SAKURA_BASE_URL = "https://api.ai.sakura.ad.jp/v1"
MODEL = "gpt-oss-120b"


def main() -> int:
    api_key = os.environ.get("SAKURA_API_KEY")
    if not api_key:
        print("ERROR: SAKURA_API_KEY が未設定です", file=sys.stderr)
        return 1

    client = OpenAI(base_url=SAKURA_BASE_URL, api_key=api_key)

    response = client.chat.completions.create(
        model=MODEL,
        messages=[{"role": "user", "content": "こんにちは"}],
    )

    message = response.choices[0].message
    usage = response.usage

    print("=== content ===")
    print(message.content)
    print()
    print("=== reasoning_content ===")
    print(getattr(message, "reasoning_content", None))
    print()
    print(
        f"=== usage === prompt={usage.prompt_tokens} "
        f"completion={usage.completion_tokens} total={usage.total_tokens}"
    )
    print(f"=== model === {response.model}")
    print(f"=== finish_reason === {response.choices[0].finish_reason}")

    return 0


if __name__ == "__main__":
    sys.exit(main())
