# 解説マニュアル自動生成パイプライン

ステータス: v1 稼働中、ブラッシュアップ進行中
最終更新: 2026-05-11
現 PM: 上林（かんばやし）
引き継ぎ予定: 2026-08月
執行部説明予定: 2026-05 今週中

## 何を解決するか

技大祭のマニュアル類は Google Document で作成されているが、当日スマホで見るには
情報密度が高すぎる・ナビゲーションが弱い・PDF を1ページずつスクロールする運用は
読みづらい、といった問題があった。

このパイプラインは Google Doc → スマホ最適化された解説 HTML への自動変換を担う。

入力:
- Google Doc を「ダウンロード → HTML」エクスポートしたファイル
- ドキュメント内の画像が同じフォルダにある状態

出力:
- 自己完結した HTML 1ファイル（CSS・JavaScript インライン、画像 base64 埋め込み）
- 任意のスマホ・PC ブラウザで開けて、オフラインでも動く
- カード形式 UI、フローティング目次ボタン、章ごとの折りたたみ等を備える

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
| Anthropic API | `scripts/generate_manual_slide.py` | `ANTHROPIC_API_KEY` | base64 inline（vision 有） | 従量課金 | 高 |
| Sakura AI Engine | `scripts/sakura-slide/generate_slide.py` | `SAKURA_API_KEY` | ファイル名のみ | 従量（安） | 中 |
| Claude Agent SDK | `scripts/claude-slide/generate_slide.py` | `claude login`（サブスク） | ファイル名のみ | サブスク枠（追加0円） | 高 |

メイン運用想定は Claude Agent SDK 版。理由は以下:
- Anthropic API: 質は高いが従量課金で予算予測しにくい
- Sakura: 安いが gpt-oss-120b はカード形式の指示への追従が弱い
- Claude Agent SDK: サブスク $20-100/月で月額固定、Opus 4.7 が使えて品質は API 版と同等

## プロンプト

共有プロンプト: `.claude/manual-prompt-card.md`

3バックエンドが同じプロンプトを参照する設計。プロンプト改善が全バックエンドに自動で波及する。

主要な要件（プロンプトに記述済）:
- SeeFT デザインシステムのカラー（teal #009688、ゴールド枠線 #C9A227 等）
- 目次セクション + 章ごとの折りたたみ + ライトボックス
- カンバン形式の役割分担カード（1枚で完結する情報密度）
- TOC ボタンカード、フローティング目次ボタン、章末ナビ、サブ目次
- スマホ縦持ち最適化、本文 15px 以上
- 元 Markdown にない情報を絶対に追加しない（メタ説明文、略称、UI ヒント等を禁止）
- 元 Markdown のコンテンツを絶対に省略しない

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

各マニュアル 3-6 分、合計 30-40 分程度。

### 生成された HTML の使い方

- 各マニュアルディレクトリの `slide_claude.card.html` を開けば閲覧可能
- 自己完結なので、Slack に添付・GitHub Pages 公開・LINE 共有等で配布できる
- スマホでは縦持ち推奨

## コスト

サブスク認証時:
- 月額 $20 (Pro) または $100 (Max 5x) の固定
- 各生成は API 換算 $0.5〜1 だが、サブスクから追加課金は発生しない
- レート制限: Pro で 5h あたり ~190 本生成可能、Max 5x なら ~950 本

API 認証時（参考）:
- 1本あたり $0.5〜1
- 8本生成で $5-8

技大祭の規模なら Pro で十分。

## 既知の挙動

### max_turns について
Claude Agent SDK 版では `max_turns=20` を設定済。

理由: 大入力（55KB+ markdown）で Claude が稀に内部的にツール試行する挙動があり、
`max_turns=10` では `disallowed_tools` 拒否で turn が消費されて flaky に失敗する事象を確認したため、
余裕大きめに `max_turns=20` で安全マージンを確保している。実際の num_turns は 1 で済むことが多い。

### 画像配置精度
LLM はファイル名と Markdown 文脈（figcaption の順序等）から推測して画像を配置する。
Opus 4.7 はこの推論が強く、配線マニュアル 25 枚で全て正しく配置できた。
ファイル名が完全に非記述的（`image1.png` 等）でも、文脈が十分なら正確に配置される。

vision を入れる検討は `manual-slide-vision-todo.md` を参照。

## 残課題（執行部発表前に解決したい）

優先度高:
- 全8マニュアルの HTML 品質検証（実際にスマホで開いて読みやすさ確認）
- 説明資料 HTML の磨き込み（執行部向けプレゼン用）

優先度中:
- 残バグの発見と修正
- 引き継ぎチェックリストの確定

優先度低（v2 候補、引き継ぎ後）:
- vision 化（画像も LLM に渡す版）
- compare_manual_versions.sh の3バックエンド対応
- 自己レビューループ（生成 → プロンプト準拠チェック → 修正）

## スケジュール（45th 技大祭向け）

| 期間 | 内容 |
| --- | --- |
| 2026-05 今週 | 執行部発表（ハード締切） |
| 2026-05 後半 | 全マニュアル試作品の品質検証、フィードバック収集 |
| 2026-06 〜 07月 | 本番運用準備、必要な改善 |
| 2026-08月 | 新 PM への引き継ぎ開始 |
| 2026-09月前半 | 45th 技大祭で実運用 |
| 2026-09月から | 現 PM 海外実務訓練、新 PM が継続運用 |

## 引き継ぎ

引き継ぎ用の詳細資料は別ファイルに分離してある:

- `docs/proposals/manual-slide-handover.html` — **新 PM 向け引き継ぎ資料**（環境セットアップ、3バックエンド詳細、運用コマンド、既知挙動、トラブルシューティング、v2 候補、チェックリスト）

本 MD は技術リファレンスとして残し、引き継ぎ実務は HTML 側を参照する。

## 関連ドキュメント

- `.claude/manual-prompt-card.md` — 現行のカード形式プロンプト（メイン）
- `.claude/manual-prompt.md` — 初代スライド形式プロンプト（参考）
- `.claude/manual-pipeline.md` — 設計検討の経緯
- `docs/proposals/manual-slide-pipeline.html` — 執行部発表用スライド（非技術的）
- `docs/proposals/manual-slide-handover.html` — 新 PM 向け引き継ぎ資料（技術的）
- `docs/proposals/manual-slide-vision-todo.md` — vision 化検討（保留）
- `scripts/compare_manual_versions.sh` — Sakura 版で新旧プロンプトを並列生成（補助ツール）
