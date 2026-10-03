# 解説マニュアル自動生成パイプライン

ステータス: v1稼働中、ブラッシュアップ進行中
最終更新: 2026-05-11
現PM: 上林（かんばやし）
引き継ぎ予定: 2026-08月
執行部説明予定: 2026-05今週中

## 何を解決するか

技大祭のマニュアル類はGoogle Documentで作成されているが、当日スマホで見るには
情報密度が高すぎる・ナビゲーションが弱い・PDFを1ページずつスクロールする運用は
読みづらい、といった問題があった。

このパイプラインは、Google Docを、スマホ向けに最適化した解説HTMLへ自動で変換する。

入力は、Google Docを「ダウンロード → HTML」でエクスポートしたファイルと、同じフォルダに置かれたドキュメント内の画像である。

出力は、自己完結したHTML 1ファイル（CSS・JavaScriptはインライン、画像はbase64で埋め込み）である。
任意のスマホ・PCブラウザで開けて、オフラインでも動く。カード形式のUI、フローティング目次ボタン、章ごとの折りたたみなどを備える。

## 全体構成

```
[ Google Doc ]
       │ Drive エクスポート（手動）
       ▼
[ ソース HTML + images/ ]
       │
       │ pandoc で HTML → Markdown 変換
       │ regex で Google Docs のノイズ除去
       │
       ▼
[ Markdown text + 画像ファイル名リスト ]
       │
       │ LLM 呼び出し（3バックエンドから選択）
       │ system_prompt = .claude/manual-prompt-card.md
       │
       ▼
[ HTML with {{filename}} プレースホルダー ]
       │
       │ replace_placeholders() で画像を base64 化して埋め込み
       │
       ▼
[ slide_xxx.card.html ] （自己完結）
```

## 3バックエンド比較

| バックエンド | スクリプト | 認証 | 画像 | コスト | 品質 |
| --- | --- | --- | --- | --- | --- |
| Anthropic API | `scripts/generate_manual_slide.py` | `ANTHROPIC_API_KEY` | base64 inline（vision有） | 従量課金 | 高 |
| Sakura AI Engine | `scripts/sakura-slide/generate_slide.py` | `SAKURA_API_KEY` | ファイル名のみ | 従量（安） | 中 |
| Claude Agent SDK | `scripts/claude-slide/generate_slide.py` | `claude login`（サブスク） | ファイル名のみ | サブスク枠（追加0円） | 高 |

メインの運用にはClaude Agent SDK版を想定している。Anthropic APIは、質は高いが従量課金なので予算を見積もりにくい。Sakuraは安いが、gpt-oss-120bはカード形式の指示に従う力が弱い。Claude Agent SDKは、サブスク $20-100/月の月額固定で、Opus 4.7が使えて品質はAPI版と同等である。

## プロンプト

共有プロンプトは`.claude/manual-prompt-card.md`である。

3バックエンドが同じプロンプトを参照する設計なので、プロンプトを改善すると全バックエンドに自動で反映される。

主な要件は次のとおりで、どれもプロンプトに記述済み。
- SeeFTデザインシステムのカラー（teal #009688、ゴールド枠線#C9A227等）
- 目次セクション + 章ごとの折りたたみ + ライトボックス
- カンバン形式の役割分担カード（1枚で完結する情報密度）
- TOCボタンカード、フローティング目次ボタン、章末ナビ、サブ目次
- スマホ縦持ち最適化、本文15px以上
- 元Markdownにない情報を絶対に追加しない（メタ説明文、略称、UIヒント等を禁止）
- 元Markdownのコンテンツを絶対に省略しない

## 運用方法

### 前提セットアップ（初回のみ）

```bash
brew install pandoc
brew install uv
claude login  # Claude Agent SDK 版のサブスク認証
```

### 単一マニュアルを生成

```bash
uv run --project scripts/claude-slide python scripts/claude-slide/generate_slide.py --prompt card --model claude-opus-4-7 docs/manuals/01_44th_配線マニュアル
```

### 全マニュアルを一括生成

```bash
for d in docs/manuals/*/; do name=$(basename "$d"); echo "===== $name ====="; uv run --project scripts/claude-slide python scripts/claude-slide/generate_slide.py --prompt card --model claude-opus-4-7 "$d" || echo "FAIL: $name"; done
```

各マニュアル3-6分、合計30-40分程度。

### 生成されたHTMLの使い方

各マニュアルディレクトリの`slide_claude.card.html`を開けば閲覧できる。自己完結しているので、Slackへの添付・GitHub Pagesでの公開・LINEでの共有などで配布できる。スマホでは縦持ちを推奨する。

## コスト

サブスク認証のときは、月額 $20 (Pro)または $100 (Max 5x)の固定である。各生成はAPI換算で $0.5〜1だが、サブスクからの追加の課金はない。レート制限は、Proで5hあたり ~190本、Max 5xなら ~950本まで生成できる。

API認証のとき（参考）は、1本あたり $0.5〜1で、8本生成すると $5-8になる。

技大祭の規模ならProで十分。

## 既知の挙動

### max_turnsについて
Claude Agent SDK版では`max_turns=20`を設定済み。

大きい入力（55KB+ markdown）では、Claudeがまれに内部でツールを試す挙動がある。
`max_turns=10`では、`disallowed_tools`の拒否でturnが消費されてflakyに失敗する事象を確認した。そのため、
余裕を大きめに取って`max_turns=20`にし、安全マージンを確保している。実際のnum_turnsは1で済むことが多い。

### 画像配置精度
LLMはファイル名とMarkdown文脈（figcaptionの順序等）から推測して画像を配置する。
Opus 4.7はこの推測が得意で、配線マニュアルでは25枚すべてを正しく配置できた。
ファイル名が完全に非記述的（`image1.png`等）でも、文脈が十分なら正確に配置される。

visionを入れる検討は`manual-slide-vision-todo.md`を参照。

## 残課題（執行部発表前に解決したい）

優先度が高いものは次のとおり。

- 全8マニュアルのHTML品質検証（実際にスマホで開いて読みやすさ確認）
- 説明資料HTMLの仕上げ（執行部向けプレゼン用）

優先度が中くらいのものは次のとおり。

- 残バグの発見と修正
- 引き継ぎチェックリストの確定

優先度が低いもの（v2候補、引き継ぎ後）は次のとおり。

- vision化（画像もLLMに渡す版）
- compare_manual_versions.shの3バックエンド対応
- 自己レビューループ（生成 → プロンプト準拠チェック → 修正）

## スケジュール（45th技大祭向け）

| 期間 | 内容 |
| --- | --- |
| 2026-05今週 | 執行部発表（ハード締切） |
| 2026-05後半 | 全マニュアル試作品の品質検証、フィードバック収集 |
| 2026-06 〜 07月 | 本番運用準備、必要な改善 |
| 2026-08月 | 新PMへの引き継ぎ開始 |
| 2026-09月前半 | 45th技大祭で実運用 |
| 2026-09月から | 現PM海外実務訓練、新PMが継続運用 |

## 引き継ぎ

引き継ぎ用の詳細資料は、別ファイルの`docs/proposals/manual-slide-handover.html`に分けてある。新PM向けの引き継ぎ資料で、環境セットアップ、3バックエンドの詳細、運用コマンド、既知の挙動、トラブルシューティング、v2候補、チェックリストを載せている。

本MDは技術リファレンスとして残し、引き継ぎ実務はHTML側を参照する。

## 関連ドキュメント

- `.claude/manual-prompt-card.md`は現行のカード形式プロンプト（メイン）
- `.claude/manual-prompt.md`は初代のスライド形式プロンプト（参考）
- `.claude/manual-pipeline.md`は設計検討の経緯
- `docs/proposals/manual-slide-pipeline.html`は執行部発表用のスライド（非技術的）
- `docs/proposals/manual-slide-handover.html`は新PM向けの引き継ぎ資料（技術的）
- `docs/proposals/manual-slide-vision-todo.md`はvision化の検討（保留）
- `scripts/compare_manual_versions.sh`は、Sakura版で新旧プロンプトを並列に生成する補助ツール
