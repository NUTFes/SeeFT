# AGENTS.md

このファイルは AI コーディングアシスタント（Claude Code, CodeRabbit 等）と新規参加メンバーが SeeFT のコードベースを理解するためのガイドです。

## 1. 概要

SeeFT は技大祭（NUTFes）の運営マニュアル管理システムです。当日スタッフが自分の担当タスクのマニュアル・手順を確認するためのアプリケーションで、シフト管理システムではありません。

### 構成

3 つのコンポーネントで構成されています。

- `api/`: Go + Echo + GORM + PostgreSQL のバックエンド API
- `mobile/lib/`: Flutter（モバイル / Web 配信）のクライアントアプリ
- `gas/`: Google Apps Script。スプレッドシートと API を連携させ、運営側のデータ管理を行う

`admin/`(Next.js) と `raw/`(メモ類) のディレクトリは現在使用していません。

## 2. 開発コマンド

主要操作はリポジトリルートの `Makefile` に集約されています。

### バックエンド / DB

```
make up        # docker compose up（API + DB + admin）
make up-db     # DB のみ起動
make up-api    # DB 起動 → API 起動
make down      # 全サービス停止
make build     # コンテナビルド
make exec      # api コンテナにシェルでログイン
make seed      # DB シード投入
make tidy      # go mod tidy（コンテナ内）
```

Mac 環境向けには `mac-up` / `mac-build` / `mac-seed`、本番向けには `prod-up` / `prod-build` / `prod-seed` を使用します。

### モバイル

```
make mobile-up   # fvm flutter run（Web ターゲット、ポート 45029）
```

Flutter のバージョン管理は `fvm`、環境変数は `mobile/env/.env` で渡します。

### compose ファイル

- `docker-compose.yml`: 開発デフォルト
- `docker-compose.mac.yml`: Mac 向け（プラットフォーム差分対応）
- `docker-compose.prod.yml`: 本番

## 3. アーキテクチャ

### `api/`（Go バックエンド）

Clean Architecture を採用しています。

```
api/
├── main.go                # エントリポイント
└── lib/
    ├── di/                # 依存性注入（di.InitializeServer）
    ├── entity/            # ビジネスエンティティ
    ├── usecase/           # ビジネスロジック層
    ├── internals/
    │   ├── controller/    # HTTP ハンドラ（Echo）
    │   └── repository/    # DB アクセス層
    └── externals/
        ├── db/            # DB 接続設定
        ├── server/        # サーバ初期化
        └── slack/         # Slack 通知
```

各層の責務は次のとおりです。

- **Controller**: HTTP リクエスト / レスポンスの処理。バリデーション、JSON シリアライズ。DB アクセスは行わない。
- **UseCase**: ビジネスロジック。Repository から受け取った `*sql.Rows` を Scan して entity に変換する責務も持つ。
- **Repository**: SQL クエリの組み立てと実行。entity への変換は行わず `*sql.Row` / `*sql.Rows` を返す。

依存性注入は `lib/di/di.go` の `InitializeServer()` で集約。各層の Factory 関数（`NewXxxController` / `NewXxxUseCase` / `NewXxxRepository`）でコンストラクタ注入を行います。

### `mobile/lib/`（Flutter クライアント）

```
mobile/lib/
├── main.dart
├── pages/      # 画面（StatefulWidget）
├── widgets/    # 再利用 Widget
├── models/     # データモデル
├── utils/      # api.dart, logger.dart, permanent_store.dart
├── configs/    # 定数、importer.dart（共通 import 集約）
└── theme/      # 色・フォント・トークン
```

状態管理は **StatefulWidget + `setState()` のみ**。Riverpod / Provider / Bloc / GetX などの状態管理ライブラリは導入していません（意図的選択）。永続化は `Hive`（複数行データ）と `SharedPreferences`（単純な設定値）を使い分けます。

API 通信は `utils/api.dart` の `Api` クラスシングルトン経由で行います。`package:http` を直接 import してよいのは `api.dart` 内のみ。

### `gas/`（Google Apps Script）

ドメインごとにディレクトリ分割しています。

```
gas/
├── shift/    # シフト管理
├── task/     # タスク管理
├── user/     # ユーザー管理
└── rescue/   # 救援要請
```

各ディレクトリに `コード.js`（メイン処理）、`onChange.js`（イベントハンドラ）、`条件付き書式.js`（書式設定）等が配置されます。スプレッドシートに直接バインドして動作します。

API URL は `PropertiesService.getScriptProperties()` から取得します（ハードコード禁止）。

## 4. コーディング規約（新規コード向け）

このセクションは **新規に書くコードの規約** です。既存コードは過渡期で混在している箇所があります（第 6 章参照）。

### Go (`api/`)

#### Repository

- SQL は **必ずプレースホルダ（`$1, $2, ...`）を使用**。文字列連結（`+ id` 等）は禁止。
- 入力値は引数として渡す: `db.QueryContext(ctx, query, id)`
- メソッド名は CRUD 動詞: `All` / `Find` / `FindByXxx` / `Create` / `Update` / `Destroy`

```go
// NG（文字列連結。SQL インジェクションの危険）
query := "SELECT * FROM bureaus WHERE id =" + id

// OK（プレースホルダ）
query := "SELECT * FROM bureaus WHERE id = $1"
rows, err := db.QueryContext(ctx, query, id)
```

#### Controller

- HTTP I/O のみを担当。DB アクセス・SQL 実行・`*sql.Row` の Scan は行わない。
- エラーレスポンスは JSON で返す: `return c.JSON(status, map[string]string{"error": err.Error()})`
- HTTP ステータスは 200 / 201 / 204 / 400 / 404 / 500 を基本に使用。

#### UseCase

- ビジネスロジックの集約点。Repository から受け取った `*sql.Rows` を Scan して entity に詰める。
- オプショナルフィールドは `sql.NullString` / `sql.NullInt64` で受け取り、`.Valid` チェック後に entity の `string` / `int` に変換。

#### 命名規則

- Interface 名: `XxxxController` / `XxxxUseCase` / `XxxxRepository`（PascalCase）
- 実装 struct 名: `xxxController` / `xxxUseCase` / `xxxRepository`（lowerCamelCase）
- Factory 関数: `NewXxxController(deps...)`（公開）
- ファイル名: `snake_case`（例: `user_repository.go`）

#### 空リスト

- レスポンスで空リストを返す場合は `[]Type{}` または `make([]Type, 0)` を使用。`nil` は避ける（JSON で `null` ではなく `[]` になるよう保証）。

#### コメント

- 日本語で記述。GoDoc 形式（関数名で始まる英語コメント）は使わない。

### Flutter (`mobile/lib/`)

- 状態管理は **`StatefulWidget` + `setState()`**。Riverpod 等の状態管理ライブラリは導入しない（意図的）。
- API 呼び出しは `utils/api.dart` のシングルトン経由のみ。`package:http` は `api.dart` 以外で直接 import 禁止。
- ログは `utils/logger.dart` の `logger`（`logger.i / e / w`）を使用。新規コードで `print()` は使わない。
- ファイル名は `snake_case`、Widget クラスは `PascalCase`、State クラスは `_WidgetNameState`（必ず private）。
- private メンバ・関数は `_` プレフィックスを付ける。
- 非同期処理後の `setState()` / `Navigator` 呼び出し前に `if (!mounted) return;` でチェック。
- 色・フォントは `theme/tokens.dart` の `AppColors` / `AppFontSizes` 等のトークン経由で参照する（生の `Color(0xFF...)` は避ける）。
- Widget には可能な限り `const` コンストラクタを使用。
- 共通 import は `configs/importer.dart` 経由で集約する。

### GAS (`gas/`)

- 破壊的操作（DB 更新・削除）の前に `ui.alert(..., ui.ButtonSet.OK_CANCEL)` で確認ダイアログを表示する。CANCEL の場合は `Logger.log("...キャンセルされました"); return;` の慣例に従う。
- API URL は `PropertiesService.getScriptProperties().getProperty("API_BASE_URL")` から取得。**ハードコード禁止**。
- バッチ操作は `LockService.getScriptLock()` を `waitLock(30000)` で取得し、try-catch-finally で実行。`finally` で必ず `releaseLock()` を呼ぶ。
- `doPost` / `doGet` の戻り値は必ず `ContentService.createTextOutput(...).setMimeType(ContentService.MimeType.TEXT)`。
- 新規コードは `var` を使わず `const`（必要に応じて `let`）で統一。
- ログは `Logger.log` を基本に使用（GAS のサーバサイドログとして残る）。デバッグ用途の `console.log` 併用は許容。

### 共通

- コメントは日本語で記述。
- 設定値（API URL、シークレット）はハードコードせず、環境変数または `PropertiesService` から取得。
- コミットメッセージは日本語、`feat:` / `fix:` 等のプレフィックスを使用（既存履歴に倣う）。
- ブランチ名は `feat/{username}/{issue-number}/{description}` または `fix/...` の形式。
- ドキュメント変更（AGENTS.md・README 等）であっても issue → branch → PR の正規フローを通す。

## 5. CI / レビュー

導入計画中の項目を含みます。実装が進み次第このセクションを拡充します。

- **AI コードレビュー**: CodeRabbit を導入予定。日本語化設定は `.coderabbit.yaml` で管理する。
- **Go 静的解析**: `golangci-lint` を GitHub Actions 上で実行する計画。`gosec` の有効化により SQL インジェクション検出を行う。
- **Flutter 静的解析**: `dart analyze` + `flutter_lints` を `analysis_options.yaml` で設定する計画。
- **エディタ統一**: `.editorconfig` を将来的に追加予定。

CI 設定ファイルが追加された段階で、本セクションに具体的なコマンドと運用ルールを追記します。

## 6. 既知の過渡期事項

新旧の規約が混在している箇所があります。新規コードは第 4 章のルールに従い、既存コードは関連 issue に沿って段階的に整理します。

### Repository 層の SQL 文字列連結（→ #266）

`api/lib/internals/repository/` 配下の 8 ファイル（`bureau_repository.go`, `department_repository.go`, `grade_repository.go`, `place_repository.go`, `review_repository.go`, `shift_repository.go`, `task_repository.go`, `time_repository.go` 等）で、SQL を文字列連結で組み立てている箇所が約 41 件残存しています。`user_repository.go` / `session_repository.go` は対応済み。詳細と進め方は #266 を参照。

### JSON キー命名の混在

- 古い entity（`User`, `Task`, `Shift` 等）は **camelCase**（例: `createdAt`, `taskID`）。
- 新しい entity（`ActionLog` 等）は **snake_case**（例: `user_id`, `created_at`）。

クライアント（mobile / gas）が現フォーマットに依存しているため、既存の camelCase は維持します。新規 entity は snake_case を推奨。

### エラーレスポンスのスタイル混在

- 古い controller: `return err`（Echo のデフォルトハンドラに委ねる）
- 新しい controller: `return c.JSON(status, map[string]string{"error": err.Error()})`

新規コードは map 形式で統一します。既存はリファクタの機会に置き換え。

### 空リスト返却の混在

- 古い UseCase: `nil` を返す箇所が一部あり。
- 新しい UseCase: `[]Type{}` または `make([]Type, 0)` で空配列を明示。

### UseCase 層の責務肥厚

Repository が `*sql.Rows` を返して UseCase 側で Scan する現状の設計は当面維持します。「Repository が entity を返す」設計へのリファクタは大規模になるため、現時点では方針外。

---

参考: 同種ドキュメントとして `FinanSu/AGENTS.md`, `group-manager-2/AGENTS.md` を参照しました。本ドキュメントは [Issue #269](https://github.com/NUTFes/SeeFT/issues/269) で初版を作成しています。
