# 解説マニュアル 実運用手順（Slack スレッド中心）

ステータス: 生成・検証・部門長レビューの手順として有効。配信と紐付けは `../development/manual-html-operations.md` が正
最終更新: 2026-06-18（配信経路の記述は 2026-08-20 の自前APIゲート導入で置き換わった）
担当: 上林（PM）

> 注（2026-08-31）: 本書の手順4「GitHub Pages へ配置して URL 発行」以降は、issue #444 / #448 で
> 実装した自前の配信・アップロードAPI（`PUT /manuals/:id`）と、シフトスプシ経由の
> `tasks.manual_url` 紐付けに置き換わった。**アップロードから紐付けまでは
> `docs/development/manual-html-operations.md` を参照すること。**
> 本書の手順1〜3（生成・画像確認・機械検証）と手順6（部門長レビューの回し方）は現在も有効。

このドキュメントは「解説マニュアル生成パイプライン」を **実際に回すときの手順書** である。
状態管理を Google Sheets で行う設計（`manual-proposal-v4-slides/automation-design.md`）に対し、
本番は **Slack スレッドを軸にした手作業フロー** で運用する。スプシ常駐監視（watcher / sheets_client /
drive_client）は当面使わない（[スプシ自動化との関係](#スプシ自動化との関係) を参照）。

技術的な背景は `manual-slide-pipeline.md`、生成エンジンの内部は `agent-sdk-usage.md` を参照。

## 登場人物と責務

- SeeFT 担当（PM ほか）: 生成・検証・公開・Slack への投稿・再生成までの「手を動かす」役。状態管理は SeeFT の仕事。
- 部門長: 生成された解説マニュアルを確認し、修正案を Slack で返す。スプシは開かない。
- 執行部: 部門長 OK 後に最終チェックして公開を確定する。

## 全体フロー

```
[SeeFT] develop で card-strict 生成
   │  scripts/claude-slide/generate_slide.py --prompt card-strict
   ▼
[SeeFT] 画像が埋まっているか確認（壊れていれば --embed-only で復旧）
   ▼
[SeeFT] 文章の機械検証（AI なし・決定的）
   │  scripts/claude-slide/verify_slide_mechanical.py
   │  → 部門長向けレポート verify_report.card-strict.md を生成
   ▼
[SeeFT] GitHub Pages リポジトリへ配置して URL 発行
   │  scripts/automation/publish.py
   ▼
[SeeFT] Slack スレッドに「URL」+「部門長向けレポート本文」を投稿
   ▼
[部門長] 解説 HTML を開いて確認 → 修正案をスレッドに返信
   ▼
[SeeFT] 修正案を instructions.md に書いて再生成 → 再検証 → 再公開（OK が出るまで繰り返す）
   ▼
[執行部] 最終チェック → 公開確定・完了
```

## 前提セットアップ（初回のみ）

```bash
brew install pandoc uv
claude login
uv sync --project scripts/claude-slide
```

GitHub Pages 用リポジトリ（解説 HTML の公開先。private で可）をローカルに clone し、環境変数で指す。
公開先パスと URL はリポジトリにハードコードしない。

```bash
export SEEFT_PAGES_REPO=~/work/seeft-manuals-pages
export SEEFT_PAGES_BASE_URL=https://nutfes.github.io/seeft-manuals
```

## ステップ詳細

作業は develop（または当該タスクのブランチ）で行う。Issue → ブランチ → PR の正規フローを省略しない。

### 1. 生成（変換 = AI / プロンプトは card-strict）

文章を一字一句変えない card-strict で生成する。元 Doc の文章は執行部と部門長の間で調整済みのため、
変換時に文言を動かしたくない、という理由で strict を採る。

```bash
uv run --project scripts/claude-slide python scripts/claude-slide/generate_slide.py --prompt card-strict --model claude-opus-4-7 docs/manuals/01_44th_のぼり広告設置マニュアル
```

出力は同ディレクトリの `slide_claude.card-strict.html`。

### 2. 画像が埋まっているか確認（壊れていれば復旧）

card-strict では稀に画像が base64 埋め込みされず、ファイルが極端に小さくなることがある（数十 KB）。
画像のあるマニュアルなら通常は数 MB になる。小さすぎる場合は再生成せず、決定的に再埋め込みできる。

```bash
uv run --project scripts/claude-slide python scripts/claude-slide/generate_slide.py --prompt card-strict --embed-only docs/manuals/01_44th_のぼり広告設置マニュアル
```

`--embed-only` は LLM を呼ばず、既存 HTML の `<img>` をローカル画像で埋め直すだけ。何度実行しても結果は同じ。

### 3. 文章の機械検証（比較 = AI なしコード）

元 Doc と生成 HTML の **文章** を決定的に比較する。デザインは検証しない。

```bash
uv run --project scripts/claude-slide python scripts/claude-slide/verify_slide_mechanical.py docs/manuals/01_44th_のぼり広告設置マニュアル
```

2 つのファイルが出る。

- `verify_report.card-strict.md` — 部門長向け。Slack にそのまま貼れる日本語サマリ。本文の「消えた・書き換わった・増えた」だけを示し、見出し番号や図番号の整理などレイアウト差は「確認不要」に畳んである。
- `verify_mechanical.card-strict.txt` — 開発デバッグ用の詳細（文字レベル差分・件数）。Slack には貼らない。

### 4. GitHub Pages へ配置して URL 発行

マニュアルは番号で識別する（`01_…` → `01`）。先頭が数字でないマニュアルは `--number` で明示する。

```bash
python3 scripts/automation/publish.py docs/manuals/01_44th_のぼり広告設置マニュアル --push
```

`<pages-repo>/manuals/01/index.html` に配置・commit・push し、`https://…/manuals/01/` を出力する。
`--push` を付けなければ commit で止まり、push コマンドを表示する（確認してから公開したいとき用）。

### 5. Slack スレッドに投稿

当該マニュアルのスレッドに、SeeFT 担当が次の 2 つを続けて投稿する。

- 公開 URL（手順 4 の出力）
- 部門長向けレポート `verify_report.card-strict.md` の本文

投稿時に部門長へメンションして確認を依頼する。状態管理はスプシではなくこのスレッドで完結させる。

### 6. 部門長レビューと修正ループ

部門長は解説 HTML を開いて確認し、修正案（読みにくい箇所・文言の直し）をスレッドに返信する。
SeeFT 担当はその修正案を、共有プロンプトではなく **マニュアルごとの追加指示ファイル** に書く。

```bash
# docs/manuals/01_44th_のぼり広告設置マニュアル/instructions.md に修正案を箇条書きで書く
```

`instructions.md` があると、生成時にその内容が「この回の追加・修正指示（最優先）」として末尾に注入される。
あとは手順 1（再生成）→ 2（画像確認）→ 3（再検証）→ 4（再公開）→ 5（再投稿）を OK が出るまで繰り返す。

instructions.md の例:

```markdown
- 「準備日タイムスケジュール」の表は、時間列を左端に固定して読みやすく
- 緊急連絡先カードは本文より目立つゴールド枠にする
- 「のぼり内容一覧」の見出しは元の文言のまま（「企業広告のぼりイメージ」に変えない）
```

### 7. 執行部の最終チェック

部門長 OK の後、執行部が最終チェックする。問題なければその URL を公開確定とし、当該マニュアルは完了。

## 1 本を通しで回すコマンド例

```bash
uv run --project scripts/claude-slide python scripts/claude-slide/generate_slide.py --prompt card-strict --model claude-opus-4-7 docs/manuals/01_44th_のぼり広告設置マニュアル
uv run --project scripts/claude-slide python scripts/claude-slide/verify_slide_mechanical.py docs/manuals/01_44th_のぼり広告設置マニュアル
python3 scripts/automation/publish.py docs/manuals/01_44th_のぼり広告設置マニュアル --push
```

このあと、出力された URL と `verify_report.card-strict.md` を Slack スレッドへ貼る。

## 部門長向けレポートの読み方

「確認してほしい本文差」だけ見れば良い。ここに出るのは、元 Doc にあって生成 HTML に見当たらない文・
書き換わった文・元になかった文。大きな表やドキュメントヘッダ（実行委員会名・担当者・作成日）が
「消えたかも」に出ることがあるので、その表や情報が必要かは部門長が判断する。

「確認不要の差」は、見出しの採番・キャプションの括弧・折りたたみボタンの文言・図番号の整理など、
レイアウト上の差で本文内容は変わっていない。図番号ラベルが整理されても画像そのものは含まれている。

## 既知の限界

- 元 Doc 側がもともと雑然としているマニュアル（巨大な表が複数、画像が大量）は、表の中身が
  「書き換わったかも／消えたかも」に多めに出る。これは表の行が再構成された結果で、本文の意味は
  保たれていることが多い。件数ではなく中身を一目見て判断する。
- 機械検証はあくまで文章の差分検出。デザインの良し悪し・読みやすさは部門長の目視で確認する。

## スプシ自動化との関係

`scripts/automation/` の watcher・sheets_client・drive_client、および
`manual-proposal-v4-slides/automation-design.md` の 19 列スキーマは「スプシを軸」とする将来案。
本番は本手順書の Slack 中心フローで回すため、これらは当面 critical path 外（保留）。コードは
残しておき、状態管理を自動化したくなったときに再評価する。

## 関連ドキュメント

- `manual-slide-pipeline.md` — パイプラインの技術リファレンス
- `agent-sdk-usage.md` — 生成エンジンが使う Claude Agent SDK の使い方
- `manual-proposal-v4-slides/automation-design.md` — スプシ軸の自動化設計（保留中の将来案）
