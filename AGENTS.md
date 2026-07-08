# AGENTS.md

SeeFT は技大祭（NUTFes）のシフト管理システムです。

## Tech Stack

- `api/`: Go 1.16 + Echo v4 + GORM v1.25 + PostgreSQL
- `mobile/lib/`: Flutter（Dart >= 3.6.0）、`fvm` 管理、`Hive` + `SharedPreferences` で永続化
- `gas/`: Google Apps Script（スプレッドシートにバインド）

`admin/` と `raw/` は使用していません。

## Commands

```bash
make up           # docker compose up（API + DB + admin）
make up-db        # DB のみ起動
make up-api       # DB → API の順で起動
make down         # 全サービス停止
make build        # コンテナビルド
make exec         # api コンテナにシェルログイン
make seed         # DB シード投入
make tidy         # go mod tidy（コンテナ内）
make test         # go test ./...（コンテナ内）
make mobile-up    # fvm flutter run -d web-server --web-port 45029 --dart-define-from-file=env/.env
```

Mac 環境は `mac-up` / `mac-build` / `mac-seed`、本番は `prod-up` / `prod-build` / `prod-seed`。

## Architecture

```text
api/lib/
├── di/                # 依存性注入（di.InitializeServer）
├── entity/            # ビジネスエンティティ
├── usecase/           # ビジネスロジック + *sql.Rows の Scan
├── internals/
│   ├── controller/    # HTTP I/O（Echo）。DB アクセス禁止
│   └── repository/    # SQL 実行のみ。*sql.Rows を返す
└── externals/{db,server,slack}

mobile/lib/
├── pages/      # 画面（StatefulWidget）
├── widgets/    # 再利用 Widget
├── models/     # データモデル
├── utils/      # api.dart, logger.dart, permanent_store.dart
├── configs/    # importer.dart（共通 import 集約）
└── theme/      # tokens.dart（AppColors, AppFontSizes）

gas/{shift,task,user,rescue}/   # ドメイン別。コード.js / onChange.js 等
```

## Code Style

### Go (`api/`)

**SQL は必ずプレースホルダで書く**

```go
query := "SELECT * FROM bureaus WHERE id = $1"
rows, err := db.QueryContext(ctx, query, id)
```

文字列連結（`"... " + id`）は SQL インジェクション脆弱性のため禁止。

**エラーレスポンスは JSON 形式で返す**

```go
return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
```

`return err` で Echo の自動処理に委ねるのは旧スタイル。

**空リストは `[]Type{}` を返す**

```go
return []entity.Task{}, nil
```

`nil` を返すと JSON が `null` になりクライアントが壊れる。

**命名規則**: Interface は `XxxxController` / `XxxxUseCase` / `XxxxRepository`（PascalCase）、実装 struct は `xxxController`（lowerCamelCase）、Factory は `NewXxxController(deps...)`、ファイル名は `snake_case`。Repository メソッドは `All` / `Find` / `FindByXxx` / `Create` / `Update` / `Destroy`。

**その他**: HTTP ステータスは 200 / 201 / 204 / 400 / 404 / 500 を基本。Optional フィールドは `sql.NullString` / `sql.NullInt64` で受け取り `.Valid` チェック後に entity の `string` / `int` に詰める。コメントは日本語、GoDoc 形式は使わない。

### Flutter (`mobile/lib/`)

**API はシングルトン経由のみ**

```dart
import 'package:seeft_mobile/configs/importer.dart';
final users = await api.getUsers();
```

`package:http` を `mobile/lib/utils/api.dart` 以外で直接 import するのは禁止。

**非同期またぎは `mounted` チェック**

```dart
final data = await api.fetchData();
if (!mounted) return;
setState(() => _data = data);
```

dispose 後に `setState` を呼ぶと例外になるため必須。

**ログは `logger`、`print` 禁止**

```dart
logger.i('loaded: ${items.length}');
logger.e('failed', error: e, stackTrace: st);
```

`print()` は新規コードでは使わない（既存コードの残骸は順次置換）。

**その他**: ファイル名は `snake_case`、Widget は `PascalCase`、State は `_WidgetNameState`。private メンバ・関数は `_` プレフィックス。Widget には可能な限り `const` コンストラクタ。色・フォントは `AppColors.main` / `AppFontSizes.md` 経由（生 hex 禁止）。

### GAS (`gas/`)

**API URL は `PropertiesService` から取得**

```javascript
const baseUrl = PropertiesService.getScriptProperties().getProperty("API_BASE_URL");
```

URL のハードコードは禁止。本番 / 開発の切り替えに `PropertiesService` を使う。

**破壊的操作の前に確認ダイアログ**

```javascript
const confirm = ui.alert("削除しますか？", ui.ButtonSet.OK_CANCEL);
if (confirm === ui.Button.CANCEL) {
  Logger.log("キャンセルされました");
  return;
}
```

**ロックは try-catch-finally で必ず解放**

```javascript
const lock = LockService.getScriptLock();
lock.waitLock(30000);
try {
  // バッチ操作
} finally {
  lock.releaseLock();
}
```

**その他**: `doPost` / `doGet` は `ContentService.createTextOutput(...).setMimeType(ContentService.MimeType.TEXT)` を返す。新規コードは `const` 使用（`var` 禁止）。ログは `Logger.log` を基本（`console.log` はデバッグ用途で併用可）。

## Git Workflow

- ブランチ名: `feat/{username}/{issue-number}/{description}` または `fix/...`
- コミットメッセージは日本語、`feat:` / `fix:` プレフィックス
- ドキュメント変更（AGENTS.md・README 等）も issue → branch → PR の正規フローを通す
- PR は `.github/pull_request_template.md` のフォーマットに従う
- PR 本文で `resolve #XXX` と書くと issue が自動 close される

## Boundaries

### Always Do
- 新規 SQL はプレースホルダで書く
- 設定値（API URL・シークレット）は環境変数 / `PropertiesService` から取得
- Flutter で非同期またぎ後の `setState()` 前に `mounted` チェック
- GAS で `LockService` 取得後は `finally` で `releaseLock()`
- 空リストは `[]Type{}` を返す

### Ask First
- 新規ライブラリの導入（特に Flutter の状態管理系）
- 既存 entity の JSON キー命名変更（mobile / gas に影響）
- API レスポンス形式の変更
- DB スキーマ変更（マイグレーション）

### Never Do
- SQL を文字列連結で組み立てる
- API URL やシークレットをハードコードする
- `package:http` を `mobile/lib/utils/api.dart` 以外で import する
- `print()` を新規コードで使う（`mobile/lib/`）
- `mobile/lib/` に Riverpod / Provider / Bloc 等の状態管理ライブラリを導入する
- `.env` や認証情報をコミットする

## Known Transitional Issues

新旧の規約が混在する箇所があります。新規コードは上記ルールに従い、既存は段階的に整理します。

- **Repository の SQL 文字列連結**（残 8 ファイル / 約 41 件） → #266
- **JSON キー命名**: 古い entity は camelCase、新しい entity は snake_case。クライアント影響のため既存維持
- **エラーレスポンス**: 古い controller は `return err`、新しいものは map 形式
- **空リスト返却**: 一部 UseCase が `nil` を返す箇所あり（順次 `[]Type{}` へ）
