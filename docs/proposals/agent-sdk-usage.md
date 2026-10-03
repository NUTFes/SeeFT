# Claude Agent SDKの使い方（このプロジェクトでの運用パターン）

ステータス: 運用中
最終更新: 2026-05-21
担当: 上林（PM）
対象読者: 新PM引き継ぎ、`scripts/claude-slide/`系の修正・拡張をする人
位置付け: `scripts/claude-slide/generate_slide.py`と`scripts/claude-slide/verify_slide.py`が裏で使っているClaude Agent SDKの使い方リファレンス。新しく実装するための知識ではなく、既存コードを読む・直す・拡張するための知識

## 1. Anthropic SDKとの違い

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

このプロジェクトでAgent SDKを選んでいるのは、`project/manual_slide_pipeline.md`で決めた「Maxプラン共用」の運用に合うためだ。APIキー方式だとAPIの料金がかかるが、Agent SDKは`claude login`済みのMaxプランの容量を使う。

その代わり、**Claude Code CLIがローカルでバックエンドとして必要**になる。Macを引き継ぐときに`claude login`をやり直す手間がかかり、CIからは動かしにくいといった運用上の制約もある。

## 2. インポートする5つのもの

Agent SDKはAPIの表面積が小さく作られていて、覚えるのはほぼ次の5つだけでよい。

```python
from claude_agent_sdk import (
    AssistantMessage,    # アシスタント (Claude) からのメッセージ
    ClaudeAgentOptions,  # 呼び出しオプション
    ResultMessage,       # 終了時の集計メッセージ (usage, duration 等)
    TextBlock,           # メッセージ内のテキストブロック
    query,               # 唯一のエントリポイント関数 (async generator)
)
```

## 3. 基本パターン(generate_slide.pyから抜粋)

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

#### `disallowed_tools`でテキストだけを返させる

Claude Codeは本来、Read・Write・Bashなどのツールを使ってコードを書くエージェントだ。本スクリプトではマニュアルHTMLを1個返してくれればよいので、ツールは要らない。ツールを全部禁止して、ツールを試すことで容量を無駄に使うのを防いでいる。

#### `max_turns`で気をつけたい点

Claude Codeは、内部でツールを使って試行錯誤するように作られている。テキストを1回返すだけでも、内部で複数のturn (思考 → ツールの試行 → 結果を見て次の手)を経ることがある。disallowed_toolsで全部ブロックしていても、ClaudeがReadを使おうとして拒否されることを何回も繰り返すと、turnを消費する。

generate_slide.pyには、経験から書いた次のコメントがある。

```python
# max_turns を 20 に: 大きい入力（55KB+ markdown）で Claude が tool 試行する場合に
# disallowed_tools 拒否で turn が消費されるため、余裕大きめ。
# お化け屋敷で max_turns=10 では flaky に失敗する事象を確認したため。
```

つまり`max_turns`は、思考の上限ではなく、**ツールの試行も含む内部ループの上限**である。大きい入力でflakyに失敗するようなら、まず`max_turns`を上げる。

### 3-2. query()をasyncで呼んで結果を受け取る

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

`query()`はasync generatorなので、`await`で1個ずつ取得するのではなく、`async for`でストリーミングして取得する。1回のリクエストで、複数のmessageが届く。

- `AssistantMessage` → Claudeのテキスト出力(1回または複数回)
- `ResultMessage` → 最後に必ず1回、集計情報(duration, num_turns, cost)

`message.content`はブロックのリストである。`TextBlock`以外のブロック(将来は思考ブロックなど)が入ってくる可能性があるので、`isinstance(block, TextBlock)`で明示的に絞っている。`result_text += block.text`と足し合わせているのは、1回の応答が長いとSDKが小分けにして送ってくることがあるためだ。

呼び出し側(main)のコードは次のとおり。

```python
import anyio

response_text, usage = anyio.run(
    call_claude_sdk,
    system_prompt, user_prompt_template, md_content, image_files, args.model,
)
```

`anyio.run()`でasync関数を同期的に起動する。`asyncio.run`でも同じことができるが、Agent SDKがanyioベースなので、それに合わせて`anyio`を使っている。

### 3-3. ResultMessageの使いどころ

`total_cost_usd`はAnthropic APIに換算した推定値である。Maxプラン経由でも、APIキーを使ったときのコスト相当が参考値として出る。この値で、Maxプラン共用によってAPI課金相当をどれくらい節約できているかを測れる。

`is_error`フラグは、今のコードでは参照されていない。本番運用に乗せる前に、次の処理を足したい。

```python
if usage.get("is_error"):
    raise RuntimeError(f"Agent SDK error: {usage}")
```

## 4. プロンプトを`.md`ファイルから読む

プロンプトはPython文字列リテラルではなく`.claude/manual-prompt-card.md`のようなMarkdownファイルから読む。その実装が`_load_prompt_from_md()`である。

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

- 1つ目の4-backtickブロック → system_prompt
- 2つ目の4-backtickブロック → user_prompt_template

外側を4-backtickで囲うのは、中に普通のコードブロック(3-backtick)を書くためだ。

### なぜこのパターンが良いか

プロンプトをgitでdiffを取れる形で管理できる。Markdownのシンタックスハイライトが付くので、エディタで編集しやすい。プロンプトのファイルには`## オーバーライド` `## 厳守事項`などの解説の節を書けて、実際にAIに渡るのはコードブロックの中だけになる。そのため、ドキュメントとしてのプロンプトと実行されるプロンプトを、同じファイルで管理できる。

### プロンプトファイルの実例

- `.claude/manual-prompt.md`はデフォルト(スライド形式)
- `.claude/manual-prompt-card.md`はカード形式
- `.claude/manual-prompt-card-strict.md`は文章不変ポリシー版
- `.claude/manual-verify-prompt.md`は検証用

対応関係は`generate_slide.py`の`PROMPT_VARIANTS`辞書で管理している。新しいプロンプトを追加する手順は次のとおり。

1. `.claude/<name>.md`を作る(4-backtickブロック2つを含む)
2. `PROMPT_VARIANTS`に`"<key>": "<name>.md"`を追加
3. `--prompt <key>`で呼び出せるようになる

## 5. レスポンスからHTMLを取り出す

Claudeが`` ```html ... ``` ``で囲んで返してくる前提で、次のパーサで取り出す。

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

fallbackとして、`<!DOCTYPE html>...</html>`を直接取り出す処理もある。AIの応答の形式を100%は信用せずに備えておく、妥当な作りのコードである。

## 6. 認証(コードに出てこないが必須)

Agent SDKを動かすには、次の準備が前提になる。

```bash
claude login
# ブラウザが開く → Anthropic アカウントでログイン
# Mac のキーチェーン or ~/.claude/ 配下に OAuth トークン保存
# 以降、claude_agent_sdk.query() は自動でそのトークンを読む
```

運用ルールは次のとおり(`project/manual_slide_pipeline.md`から)。

- 1 Mac = 1アカウント運用、同時複数Mac使用を避ける
- ベトナム期間中は他デバイスでログアウト
- 引き継ぎ時に新PMのMacで1回だけ`claude login`
- 月1でAnthropic利用状況を確認

コード側は何もしない。シェル環境(= claude CLIの状態)を信頼する作りである。

## 7. このパターンの強みと弱み

### 強み

APIキーの課金はゼロで、Maxプランの容量の中に収まる。APIの表面積が小さく、`query`と4つの型だけ覚えればよい。プロンプトを`.md`で管理するので、git diffとPRレビューで変更を確かめられる。disallowed_toolsのパターンは、Claude CodeのAgentとしての性質を逆手に取って、テキストを返すことだけに専念させる。生成と検証が同じ構造なので、運用全体で覚えることが1つで済む。

### 弱み

Macが必要になる。Claude Code CLIが動くマシンに依存するので、CIやVPSでは動かしにくい。`max_turns`の調整も要る。ツールの拒否でturnを消費する性質は直感に反していて、たまに`flaky`な失敗が出る。`is_error`はまだ使っていないので、今のコードはエラー時の扱いが緩い。本番で使う前に強化したい。AnthropicのTOSの上ではグレーな領域にある。個人のサブスクをチームで共用する運用は、厳密には「単一ユーザー前提」と食い違う(詳細は`manual_slide_pipeline.md`)。

## 8. 別バックエンド(APIキー方式)への移行コスト

万一サブスク共用で運用できなくなった場合は、`scripts/generate_manual_slide.py` (Anthropic APIキー版)に切り替える。移行で書き換える部分は次の表のとおり。

| 部分 | Agent SDK版 | APIキー版(移行後) |
| --- | --- | --- |
| インポート | `from claude_agent_sdk import ...` | `import anthropic` |
| クライアント生成 | (暗黙、`query()`が自動認証) | `client = anthropic.Anthropic(api_key=os.environ["ANTHROPIC_API_KEY"])` |
| リクエスト送信 | `async for m in query(...)` | `client.messages.create(...)` |
| メッセージ解析 | `AssistantMessage`の`TextBlock`を集める | `response.content`の`TextBlock`を集める |
| ツール禁止 | `disallowed_tools=[...]` | tools引数を渡さない(デフォルトでツールなし) |
| `max_turns` | あり | 概念なし。`max_tokens`で出力長を制御 |
| 認証 | `claude login` | 環境変数`ANTHROPIC_API_KEY` |

`call_claude_sdk()`関数を`call_anthropic_api()`に書き換えれば、他のロジック(promptの読み込み・画像のbase64化・HTMLの抽出)はそのまま使える。

## 9. 参考: コード上の出現箇所

- `scripts/claude-slide/generate_slide.py:1-322`は生成の本体。中心は`call_claude_sdk()`
- `scripts/claude-slide/verify_slide.py:1-253`は検証の本体。ほぼ同じパターン
- `scripts/claude-slide/pyproject.toml`は依存の宣言(`claude-agent-sdk>=0.1.80`, `anyio>=4.0`)
- `.claude/manual-prompt-card.md`はプロンプトの実例(4-backtickパターン)
- `.claude/manual-verify-prompt.md`は検証用プロンプトの実例

## 10. 関連ドキュメント

- `docs/proposals/manual-proposal-v4-slides/automation-design.md`は解説マニュアル生成パイプライン全体の実装設計(本ドキュメントの親)
- `docs/proposals/manual-slide-pipeline.md`は既存パイプラインの技術リファレンス
- `docs/proposals/manual-slide-pipeline-qa.md`はQ&A形式の運用解説
- `project/manual_slide_pipeline.md` (memory)はバックエンド方針(Max共用/APIキー/Sakura不採用)
