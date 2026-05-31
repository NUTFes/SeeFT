# 解説マニュアル生成パイプライン — 自動化実装設計

ステータス: Phase 1.5 完了、Phase 2 実装中
最終更新: 2026-05-13
担当: 上林（PM）
位置付け: **実装者視点の単一正本**。新 PM への引き継ぎはまずこの MD を読む

このドキュメントは「解説マニュアル生成パイプライン」の自動化部分について、**コードと同期して育てる**実装設計メモ。執行部視点 / インフラ視点 / 機能視点の他 MD（後述）に対して、本 MD は「**実装する人**」が手を動かす際の正本として機能する。

## 1. 全体ビジョン: 「スプシを軸」にする

PM が触るのは Google Sheets だけ。コード（generate / verify / アップロード）は全てスプシのデータを読み書きする側に回る。

```
[PM]
  │ Google Sheets でステータス列を「再生成」に変える
  ▼
[watcher.py（ローカル常駐 or cron）]
  │ Sheets API でステータス変更を検知
  ▼
[process_one.py（1 本のマニュアル処理）]
  │ ① スプシから行を読む（Doc URL / 個別生成指示 / 解説HTML生成 等）
  │ ② Drive API で Doc を HTML として取得
  │ ③ pandoc → generate_slide.py で解説 HTML 生成
  │ ④ verify_slide.py で AI 検証
  │ ⑤ ホスティング先にアップロード（インフラ部門と決定後）
  │ ⑥ Sheets API で結果を書き戻し
  ▼
[Google Sheets]
  ステータス「検証済(OK/NG)」、生成HTML URL、件数、最終生成日時 が反映される
```

## 2. Phase の進捗

| Phase | 内容 | 状態 | 関連物 |
| --- | --- | --- | --- |
| **Phase 1** | スプシ可視化（手動運用、CSV で状態管理） | 完了 | `docs/spread_sheets/45th_マニュアル生成ステータス.csv` |
| **Phase 1.5** | 個別生成指示の手動同期（PM が CSV 経由でプロンプトに追加指示を流せる） | 完了 | `generate_slide.py` の CSV ルックアップ機能、19 列に拡張済み |
| **Phase 2** | スプシ軸の完全自動化（Sheets API + Drive API + watcher） | **実装中** | 本 MD の Step 1-6 を参照 |
| **Phase 3（任意）** | GAS onEdit + webhook で polling を webhook に置換 | 未着手 | qa.md TODO 参照 |

## 3. 構成部品

### 既存（Phase 1.5 まで）

| 部品 | 役割 | 場所 |
| --- | --- | --- |
| `generate_slide.py` | Claude Agent SDK + サブスク認証で解説 HTML 生成。CSV から個別生成指示も注入 | `scripts/claude-slide/generate_slide.py` |
| `verify_slide.py` | 元 Doc と生成 HTML の AI 検証、VERDICT: OK/NG + 件数を返す | `scripts/claude-slide/verify_slide.py` |
| `.claude/manual-prompt-card.md` | カード形式生成用の共通プロンプト | `.claude/manual-prompt-card.md` |
| `.claude/manual-verify-prompt.md` | AI 検証用プロンプト | `.claude/manual-verify-prompt.md` |
| ステータス CSV | 19 列の状態管理 | `docs/spread_sheets/45th_マニュアル生成ステータス.csv` |
| `build_status_csv.py` | xlsx + verify レポートから CSV を再構築するユーティリティ | `/tmp/build_status_csv.py`（将来 `scripts/` に昇格予定） |

### 新規（Phase 2 で作る）

| 部品 | 役割 | 工数目安 | 依存 |
| --- | --- | --- | --- |
| `scripts/automation/sheets_client.py` | Google Sheets API ラッパー、行の読み書き、enum 値の検証 | 2-3h | OAuth 認証必須 |
| `scripts/automation/drive_client.py` | Google Drive API で Doc を HTML エクスポート | 1-2h | OAuth 認証必須 |
| `scripts/automation/process_one.py` | マニュアル 1 本の生成 → 検証 → スプシ書き戻し（メインエンジン） | 1-2h | sheets_client + 既存 generate/verify |
| `scripts/automation/watcher.py` | スプシ polling でステータス変更検知、`process_one.py` を起動 | 1-2h | sheets_client |
| `scripts/automation/uploader.py` | 生成 HTML を配信先にアップロード（インターフェースだけ先に切る、実装はインフラ部門の決定後） | 30 分（スタブ）+ TBD（本実装） | ホスティング先依存 |

### Phase 2 のディレクトリ構成（提案）

```
scripts/
├── claude-slide/                 # 既存 (Phase 1.5 まで)
│   ├── generate_slide.py
│   ├── verify_slide.py
│   ├── pyproject.toml            # claude_agent_sdk dep
│   └── ...
└── automation/                    # 新規 (Phase 2)
    ├── sheets_client.py
    ├── drive_client.py
    ├── process_one.py
    ├── watcher.py
    ├── uploader.py
    └── pyproject.toml             # google-api-python-client 系 dep
```

`claude-slide` と `automation` を分けるのは依存が違うため（LLM 系 vs Google API 系）。お互いを subprocess で呼び合う / 共有モジュールを参照する形にする。

## 4. データフロー詳細

### 4-1. 起動からスプシ書き戻しまで

```
[watcher.py] スプシを 30-60 秒ごとに polling
   ↓ 「再生成」または「未生成」のステータスを検出
[process_one.py manual_name]
   ↓
sheets_client.read_row(manual_name)
   → 行データを取得 (マニュアル名 / Doc URL / 解説HTML生成 / 個別生成指示 等)
   ↓
[分岐 1] 解説HTML生成 == "生成しない":
   ├ Doc URL をそのまま manual_url としてマーク (このサイクル終了)
   └ パイプライン状態 = "Doc 直配信" にしてスプシに書き戻し

[分岐 2] 解説HTML生成 == "生成する":
   ├ drive_client.export_doc_as_html(doc_url, dest=docs/manuals/{name}/)
   │    → docs/manuals/{name}/source.html を作成
   ├ pandoc で source.html → markdown 化（既存ロジック）
   ├ generate_slide.py を呼ぶ
   │    → docs/manuals/{name}/slide_claude.card.html を生成
   ├ verify_slide.py を呼ぶ
   │    → docs/manuals/{name}/verify_claude.card.txt を保存
   │    → VERDICT (OK/NG) と件数 (追加/欠落/改変) を取得
   ├ uploader.upload(slide_html, name) → URL を取得 (Phase 2 後半)
   └ sheets_client.write_row(manual_name, {
        "生成HTML URL": url,
        "パイプライン状態": "検証済(OK)" or "検証済(NG)",
        "最終生成日時": now,
        "検証VERDICT": "OK" or "NG",
        "追加件数": N, "欠落件数": M, "改変件数": K,
        "検証レポート": "docs/manuals/{name}/verify_claude.card.txt",
        "確認結果": "未確認" if NG else "",
     })
```

### 4-2. SeeFT モバイル側との接続

`mobile/lib/widgets/manual_viewer.dart:8-27` の既存ロジックがそのまま動く:

- `manual_url` が `docs.google.com` → Doc プレビュー iframe（PR #258）
- `manual_url` が `.pdf` → pdf.js viewer
- `manual_url` がその他 → 普通の iframe（生成 HTML はこの分岐）

→ **モバイル側は無改修**。詳細は `selective-html-generation.md` 参照。

## 5. 認証 (OAuth セットアップ)

### 5-1. 必要なもの

- Google Cloud Console プロジェクト（PM の Google アカウントで作る）
- 有効化する API: **Google Sheets API**、**Google Drive API**
- OAuth 同意画面: User Type = 外部、テストユーザーに PM の Gmail を追加（未公開のままで OK）
- OAuth クライアント ID: アプリケーションの種類 = **デスクトップアプリ**
- ダウンロードした `credentials.json`

### 5-2. ファイル配置

```
~/.config/seeft-pipeline/
├── credentials.json    # GCP からダウンロード、リポジトリには絶対に置かない
└── token.json          # 初回スクリプト実行時に自動生成、以降は再利用
```

リポジトリ外（`~/.config/`）に置く理由:

- git の事故防止（誤コミットで認証情報が公開される事態を避ける）
- 新 PM 引き継ぎ時に「このディレクトリをコピーすれば動く」状態にできる
- メモリ `project/ignore_convention.md` の「ブランチ依存生成物は `.git/info/exclude` に寄せる」原則を、認証情報には拡張適用

### 5-3. 初回起動フロー

```
1. PM が任意の automation/* スクリプトを初回実行
2. credentials.json を読む
3. ブラウザが自動で開く → Google アカウントにログイン → スコープ同意
4. token.json が自動生成・保存される
5. 以降は token.json を使い回し（有効期限切れ時は自動リフレッシュ）
```

### 5-4. スコープ

- `https://www.googleapis.com/auth/spreadsheets` — Sheets 読み書き
- `https://www.googleapis.com/auth/drive.readonly` — Doc 取得（読み取りのみで十分）

書き込みスコープを Drive に与えないことで、PM の Drive 内ファイルを誤って改変するリスクをゼロにする。

## 6. ステータス CSV スキーマ（19 列）

| # | カラム名 | 型 | 記入主体 | enum 値 / 例 |
| --- | --- | --- | --- | --- |
| 1 | マニュアル名 | str | PM | "配線マニュアル" |
| 2 | 担当局 | str | PM | "総務" "渉外" "財務" "企画" |
| 3 | 担当部門 | str | PM | "会場" "副局長" "広報" 等 |
| 4 | 担当者名 | str | PM | "赤嶺" "黒木康士朗" |
| 5 | Google Doc URL | URL | PM | "https://docs.google.com/document/d/.../edit" |
| 6 | **解説HTML生成** | enum | PM | "生成する" / "生成しない" |
| 7 | 生成HTML URL | URL | automation | "https://manuals.../wiring.html" |
| 8 | パイプライン状態 | enum | automation | "未生成" "生成中" "検証中" "検証済(OK)" "検証済(NG)" "Doc 直配信" "エラー" |
| 9 | 最終生成日時 | datetime | automation | "2026-05-13 13:30" |
| 10 | 検証VERDICT | enum | automation | "OK" / "NG" / "" |
| 11 | 追加件数 | int | automation | 0, 1, 2, ... |
| 12 | 欠落件数 | int | automation | 0, 1, 2, ... |
| 13 | 改変件数 | int | automation | 0, 1, 2, ... |
| 14 | 検証レポート | path | automation | "docs/manuals/.../verify_claude.card.txt" |
| 15 | 確認担当者 | str | PM | "赤嶺" |
| 16 | 確認結果 | enum | 確認担当者 | "未確認" / "訂正OK" / "要修正" / "再生成依頼" |
| 17 | 確認備考 | str | 確認担当者 | フリーテキスト |
| 18 | 個別生成指示 | markdown | PM | "- 配線番号 ①〜㉖ はバッジで強調" |
| 19 | 備考 | str | 全員 | フリーテキスト |

「PM」記入の列は人が触る、「automation」記入の列は自動化が書き戻す、と境界が明確。

### enum 値のプルダウン推奨設定（Google Sheets 側）

| カラム | プルダウン値（リストを直接指定） |
| --- | --- |
| 解説HTML生成 | `生成する,生成しない` |
| パイプライン状態 | `未生成,生成中,検証中,検証済(OK),検証済(NG),Doc 直配信,エラー` |
| 検証VERDICT | `OK,NG` |
| 確認結果 | `未確認,訂正OK,要修正,再生成依頼` |

## 7. 実装ステップ

| Step | 内容 | 状態 | 主作業者 |
| --- | --- | --- | --- |
| 1 | 「解説HTML生成」列を CSV に追加（19 列に拡張） | 完了 | Claude |
| 2 | Google Cloud Console で OAuth セットアップ、credentials.json を `~/.config/seeft-pipeline/` に配置 | **進行中** | PM |
| 3 | `scripts/automation/sheets_client.py` 実装、認証フローと行読み書きをテスト | 未着手 | Claude |
| 4 | `scripts/automation/process_one.py` 実装、1 マニュアルで end-to-end 動作確認 | 未着手 | Claude |
| 5 | `scripts/automation/drive_client.py` 実装、Doc URL → HTML エクスポート | 未着手 | Claude |
| 6 | `scripts/automation/watcher.py` 実装、polling でステータス変更検知 | 未着手 | Claude |
| 並行 | `uploader.py` インターフェースだけ先に切る | 未着手 | Claude |
| 並行 | ホスティング先決定（→ 本実装） | 未着手 | PM × インフラ部門 |

## 8. 依存パッケージ・ファイル配置

### Python 依存（`scripts/automation/pyproject.toml`）

```toml
[project]
name = "seeft-automation"
version = "0.1.0"
requires-python = ">=3.11"
dependencies = [
    "google-api-python-client>=2.0",
    "google-auth-httplib2>=0.2",
    "google-auth-oauthlib>=1.2",
]
```

### システム依存

- pandoc (`brew install pandoc`)
- uv (`brew install uv`)
- Claude Code CLI (`claude login` 済み)

### 認証ファイル

- `~/.config/seeft-pipeline/credentials.json` — GCP からダウンロード
- `~/.config/seeft-pipeline/token.json` — 初回起動で自動生成

### スプシ ID

- 環境変数 or 設定ファイル: `SEEFT_STATUS_SHEET_ID=1jz_870-Id89UYS-00F9ozZUNWL92IPntRqtNYYCRF0c`
- ハードコードせず外部から差し替え可能に

## 9. 新 PM 引き継ぎ用クイックスタート

新 PM が初日に動かすための最短手順:

```bash
# 1. リポジトリを clone
git clone https://github.com/NUTFes/SeeFT.git
cd SeeFT

# 2. 環境を作る
brew install pandoc uv
uv sync --project scripts/claude-slide
uv sync --project scripts/automation

# 3. Claude にログイン
claude login

# 4. 上林さんから OAuth 認証情報を受け取り、所定の場所に置く
mkdir -p ~/.config/seeft-pipeline
# credentials.json と token.json を ~/.config/seeft-pipeline/ にコピー

# 5. 試運転
uv run --project scripts/automation python scripts/automation/process_one.py 配線マニュアル

# 6. 結果がスプシに反映されれば OK
```

## 10. 関連ドキュメント

| MD | 視点 | 主な内容 |
| --- | --- | --- |
| `index.html` | 執行部 | 13 スライドの提案資料 |
| `infra-hosting-discussion.md` | インフラ部門 | ホスティング先決定のための相談材料 |
| `selective-html-generation.md` | 機能設計 | 「解説HTML生成」列と Doc 直表示の両立 |
| `45th_マニュアル可視化方針（仮）.md` | 可視化方針 | 全体方針（参照） |
| `../manual-proposal-v4.md` | 提案本文 | 執行部 MT 向け Q1-Q6 形式 |
| `../manual-slide-pipeline.md` | 技術リファレンス | 既存パイプラインの解説 |
| `../manual-slide-pipeline-qa.md` | Q&A | 各種疑問への回答集、TODO 高レベル方針 |
| `../agent-sdk-usage.md` | 実装者 | `generate_slide.py` / `verify_slide.py` が裏で使う Claude Agent SDK の使い方リファレンス |

## 11. コード上の連携先

- `scripts/claude-slide/generate_slide.py` — Phase 1.5 で CSV 読み込み機能を追加済み、`load_per_manual_instructions()` 関数
- `scripts/claude-slide/verify_slide.py` — AI 検証本体
- `mobile/lib/widgets/manual_viewer.dart:8-27` — URL パターン分岐、改修不要
- `mobile/lib/widgets/shift_card.dart:165-171` — マニュアルボタン、改修不要

## 12. メモ・運用ルール

- **「生成しない」のマニュアル**は Doc URL を manual_url にそのまま入れる。生成パイプラインは走らない（process_one.py が早期 return する）
- **エラー時の動作**: Doc 取得失敗 / generate 失敗 / verify 失敗、いずれもパイプライン状態 = "エラー"、備考に詳細を書く。次回再生成で復活可能
- **同時実行制御**: process_one.py は単一プロセスで走らせる前提（同じマニュアルを並行処理しない）。watcher は処理中のマニュアルを再起動しないようロック
- **権限分離**: PM はスプシのオーナー権限、確認担当者は編集権限のみ、関係者外はビュー権限。スコープを絞ることで誤編集を防ぐ
- **PM 不在時のフォールバック**: 自動化が動かなくなった場合、PM の Mac で `generate_slide.py` を直接叩く運用に戻す（Phase 1 と同じフロー）。コードを残しておくこと

---

このドキュメントは Phase 2 実装の進捗に合わせて更新する。実装した内容は同日中に本 MD に反映するルールにすると、コードと文書がズレない。
