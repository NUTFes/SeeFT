# Claude Agent SDK の使い方（このプロジェクトでの運用パターン）

ステータス: 運用中
最終更新: 2026-05-21
担当: 上林（PM）
対象読者: 新 PM 引き継ぎ、`scripts/claude-slide/` 系の修正・拡張をする人
位置付け: `scripts/claude-slide/generate_slide.py` と `scripts/claude-slide/verify_slide.py` が裏で使っている Claude Agent SDK の使い方リファレンス。新規実装ではなく、既存コードを読む・直す・拡張するための知識

## 1. Anthropic SDK との違い

```
Anthropic SDK (anthropic パッケージ)
  └─ Claude モデルを API キーで叩く。普通の SDK。Anthropic にトークン課金

Claude Agent SDK (claude_agent_sdk パッケージ)
  └─ Claude Code CLI のサブスク認証を流用する SDK
     ・OAuth トークンで認証 (API キー不要)
     ・「Claude Code が裏で持っている対話セッション」を
       スクリプトから呼べるイメージ
     ・Agent としての構造化 (system_prompt, tools, max_turns) を持つ
```

このプロジェクトで Agent SDK を選んでいる理由は、`project/manual_slide_pipeline.md` で確定した「Max プラン共用」運用と整合するため。API キー方式だと API 課金が発生するが、Agent SDK は `claude login` 済の Max プラン容量を消費する。

代償として、**Claude Code CLI がローカルでバックエンドとして必要**。Mac の引き継ぎで `claude login` をやり直す手間や、CI からは動かしにくい等の運用制約が付く。

## 2. インポートする 5 つのもの

Agent SDK は API 表面積が小さく作られていて、ほぼこれだけ覚えれば良い:

```python
from claude_agent_sdk import (
    AssistantMessage,    # アシスタント (Claude) からのメッセージ
    ClaudeAgentOptions,  # 呼び出しオプション
    ResultMessage,       # 終了時の集計メッセージ (usage, duration 等)
    TextBlock,           # メッセージ内のテキストブロック
    query,               # 唯一のエントリポイント関数 (async generator)
)
```

## 3. 基本パターン (generate_slide.py から抜粋)

### 3-1. オプション組み立て

```python
DISALLOWED_TOOLS = [
    "Read", "Write", "Edit", "Bash",
    "Task", "WebFetch", "WebSearch",
    "Grep", "Glob", "TodoWrite",
    "NotebookEdit",
]

options = ClaudeAgentOptions(
    system_prompt=system_prompt,    # .md ファイルから読んだプロンプト
    max_turns=20,                   # 内部ループの上限
    disallowed_tools=DISALLOWED_TOOLS,
    model="claude-opus-4-7",        # 省略可。省略すると Claude Code デフォルト
)
```

#### `disallowed_tools` で「テキストだけ吐かせる」のがコツ

Claude Code は本来「Read/Write/Bash 等のツールを使ってコードを書く」エージェント。だが本スクリプトでは「マニュアル HTML を 1 個返してくれればいい」のでツール不要。全ツール禁止することで「ツール試行で容量を浪費する」のを防いでいる。

#### `max_turns` の落とし穴

Claude Code は内部で「ツールを使って試行錯誤」する設計。テキスト返答 1 回でも、内部で複数 turn (思考 → ツール試行 → 結果見て次の手) を経ることがある。disallowed_tools で全部ブロックしてても、Claude が「Read を使おうとして拒否される」を何回も繰り返すと turn を消費する。

generate_slide.py には経験則として以下のコメントがある:

```python
# max_turns を 20 に: 大きい入力（55KB+ markdown）で Claude が tool 試行する場合に
# disallowed_tools 拒否で turn が消費されるため、余裕大きめ。
# お化け屋敷で max_turns=10 では flaky に失敗する事象を確認したため。
```

つまり `max_turns` は **「思考の上限」ではなく「ツール試行も含む内部ループの上限」**。大きい入力で flaky に失敗するようなら、まず `max_turns` を上げる。

### 3-2. query() を async で叩いて結果を吸う

```python
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

return result_text, usage
```

ポイント:

- `query()` は **async generator**。`await` で 1 個ずつ取得ではなく、`async for` でストリーミング取得
- 1 回のリクエストで複数 message が流れてくる:
  - `AssistantMessage` → Claude のテキスト出力 (1 回または複数回)
  - `ResultMessage` → 最後に必ず 1 回、集計情報 (duration, num_turns, cost)
- `message.content` は **ブロックのリスト**。`TextBlock` 以外 (将来的に思考ブロックなど) が混入する可能性があるので、明示的に `isinstance(block, TextBlock)` で絞る
- `result_text += block.text` と**累積加算**しているのは、1 回の応答が長文だと SDK が部分的に小分けして送ってくることがあるため

呼び出し側 (main) では:

```python
import anyio

response_text, usage = anyio.run(
    call_claude_sdk,
    system_prompt, user_prompt_template, md_content, image_files, args.model,
)
```

`anyio.run()` で async 関数を同期的に起動。`asyncio.run` でも同じことができるが、`anyio` を使ってるのは Agent SDK が anyio ベースだから合わせている。

### 3-3. ResultMessage の使いどころ

`total_cost_usd` は Anthropic API 換算の推定値。Max プラン経由でも、参考値として API キー使用時のコスト相当が出る。**「Max プラン共用で API 課金相当をどれくらい節約できているか」を測れる**。

`is_error` フラグは現状コードで参照されていない。本番運用に乗せる前に下記を足したい:

```python
if usage.get("is_error"):
    raise RuntimeError(f"Agent SDK error: {usage}")
```

## 4. プロンプトを `.md` ファイルから読む

プロンプトは Python 文字列リテラルではなく `.claude/manual-prompt-card.md` のような Markdown ファイルから読む。`_load_prompt_from_md()` がその実装:

```python
def _load_prompt_from_md(path: str) -> tuple[str, str]:
    with open(path, "r", encoding="utf-8") as f:
        content = f.read()
    # Prefer 4-backtick outer fences (allows nested 3-backtick blocks inside).
    blocks = re.findall(r"````\s*\n(.*?)\n````", content, re.DOTALL)
    if len(blocks) < 2:
        blocks = re.findall(r"```\n(.*?)```", content, re.DOTALL)
    if len(blocks) < 2:
        raise ValueError(f"Expected 2 code blocks in {path}, found {len(blocks)}")
    return blocks[0].strip(), blocks[1].strip()
```

- 1 つ目の 4-backtick ブロック → system_prompt
- 2 つ目の 4-backtick ブロック → user_prompt_template
- 4-backtick で外側を囲うのは、中に普通のコードブロック (3-backtick) を書きたいから

### なぜこのパターンが良いか

- プロンプトを git で diff 可能に管理できる
- Markdown のシンタックスハイライトが効くのでエディタの編集体験が良い
- プロンプトには `## オーバーライド` `## 厳守事項` 等の**解説セクション**を書け、コードブロック内だけが実際に AI に渡る
- 「ドキュメントとしてのプロンプト」と「実行されるプロンプト」を同じファイルで管理できる

### プロンプトファイルの実例

- `.claude/manual-prompt.md` — デフォルト (スライド形式)
- `.claude/manual-prompt-card.md` — カード形式
- `.claude/manual-prompt-card-strict.md` — 文章不変ポリシー版
- `.claude/manual-verify-prompt.md` — 検証用

`generate_slide.py` の `PROMPT_VARIANTS` 辞書で対応関係を管理している。新しいプロンプトを追加するには:

1. `.claude/<name>.md` を作る (4-backtick ブロック 2 つを含む)
2. `PROMPT_VARIANTS` に `"<key>": "<name>.md"` を追加
3. `--prompt <key>` で呼び出せるようになる

## 5. レスポンスから HTML を取り出す

Claude が `` ```html ... ``` `` で囲んで返してくる前提のパーサ:

```python
def extract_html(response: str) -> str:
    match = re.search(r"```html\s*\n(.*?)```", response, re.DOTALL)
    if match:
        return match.group(1).strip()
    match = re.search(r"(<!DOCTYPE html>.*?</html>)", response, re.DOTALL | re.IGNORECASE)
    if match:
        return match.group(1).strip()
    return response.strip()
```

fallback として `<!DOCTYPE html>...</html>` 直接抽出も。AI 応答のフォーマットを 100% 信用しない健全な防衛コード。

## 6. 認証 (コードに出てこないが必須)

Agent SDK が動くための前提条件:

```bash
claude login
# ブラウザが開く → Anthropic アカウントでログイン
# Mac のキーチェーン or ~/.claude/ 配下に OAuth トークン保存
# 以降、claude_agent_sdk.query() は自動でそのトークンを読む
```

運用ルール (`project/manual_slide_pipeline.md` から):

- 1 Mac = 1 アカウント運用、同時複数 Mac 使用を避ける
- ベトナム期間中は他デバイスでログアウト
- 引き継ぎ時に新 PM の Mac で 1 回だけ `claude login`
- 月 1 で Anthropic 利用状況を確認

コード側は何もしない。シェル環境 (= claude CLI の状態) を信頼するモデル。

## 7. このパターンの強みと弱み

### 強み

- **API キー課金ゼロ** (Max プラン容量で吸収)
- **API 表面積が小さく覚えやすい** (`query` と 4 つの型だけ)
- **プロンプトを `.md` で管理** → git diff・PR レビューが効く
- **disallowed_tools パターン** → Claude Code の Agent 特性を逆手に取って「テキスト返答に専念させる」
- **生成と検証で同じ構造** → 学習コストが運用全体で 1 つで済む

### 弱み

- **Mac が必要**: Claude Code CLI が動くマシン依存。CI や VPS では動かしにくい
- **`max_turns` チューニング**: ツール拒否で turn を消費する特性が直感に反する。たまに `flaky` 失敗が出る
- **`is_error` 未活用**: 現コードはエラー時の挙動が緩い。本番化前に強化したい
- **Anthropic TOS グレー領域**: 個人サブスクをチームで共用する運用は厳密には「単一ユーザー前提」と擦れる (詳細: `manual_slide_pipeline.md`)

## 8. 別バックエンド (API キー方式) への移行コスト

万一サブスク共用が運用できなくなった場合、`scripts/generate_manual_slide.py` (Anthropic API キー版) に切り替える。移行で書き換える部分:

| 部分 | Agent SDK 版 | API キー版 (移行後) |
| --- | --- | --- |
| インポート | `from claude_agent_sdk import ...` | `import anthropic` |
| クライアント生成 | (暗黙、`query()` が自動認証) | `client = anthropic.Anthropic(api_key=os.environ["ANTHROPIC_API_KEY"])` |
| リクエスト送信 | `async for m in query(...)` | `client.messages.create(...)` |
| メッセージ解析 | `AssistantMessage` の `TextBlock` を集める | `response.content` の `TextBlock` を集める |
| ツール禁止 | `disallowed_tools=[...]` | tools 引数を渡さない (デフォルトでツールなし) |
| `max_turns` | あり | 概念なし。`max_tokens` で出力長を制御 |
| 認証 | `claude login` | 環境変数 `ANTHROPIC_API_KEY` |

`call_claude_sdk()` 関数を `call_anthropic_api()` に書き換える形で、他のロジック (prompt 読み込み・画像 base64 化・HTML 抽出) は再利用可能。

## 9. 参考: コード上の出現箇所

- `scripts/claude-slide/generate_slide.py:1-322` — 生成本体、`call_claude_sdk()` がメイン
- `scripts/claude-slide/verify_slide.py:1-253` — 検証本体、ほぼ同じパターン
- `scripts/claude-slide/pyproject.toml` — 依存宣言 (`claude-agent-sdk>=0.1.80`, `anyio>=4.0`)
- `.claude/manual-prompt-card.md` — プロンプトの実例 (4-backtick パターン)
- `.claude/manual-verify-prompt.md` — 検証用プロンプトの実例

## 10. 関連ドキュメント

- `docs/proposals/manual-proposal-v4-slides/automation-design.md` — 解説マニュアル生成パイプライン全体の実装設計 (本ドキュメントの親)
- `docs/proposals/manual-slide-pipeline.md` — 既存パイプラインの技術リファレンス
- `docs/proposals/manual-slide-pipeline-qa.md` — Q&A 形式の運用解説
- `project/manual_slide_pipeline.md` (memory) — バックエンド方針 (Max 共用 / API キー / Sakura 不採用)
