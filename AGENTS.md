# AGENTS.md

SeeFT は技大祭（NUTFes）のシフト管理システムです。

## Tech Stack

- `api/`: Go 1.26 + Echo v4 + `database/sql`（一部 GORM v1.25）+ PostgreSQL
- `mobile/lib/`: Flutter 3.27.3（Dart >= 3.6.0）、`fvm` 管理、Web のみで運用、`Hive` + `SharedPreferences` で永続化
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
├── router/            # URL と controller の対応づけ（全エンドポイントの一覧）
├── entity/            # ビジネスエンティティ
├── usecase/           # ビジネスロジック + *sql.Rows の Scan
├── internals/
│   ├── controller/    # HTTP I/O（Echo）。DB アクセス禁止
│   └── repository/    # SQL 実行のみ。*sql.Rows を返す
└── externals/{db,server,slack,scheduler}   # scheduler: 通知の定期実行 ticker ループ

mobile/lib/
├── pages/      # 画面（StatefulWidget）
├── widgets/    # 再利用 Widget
├── models/     # データモデル
├── utils/      # api.dart, logger.dart, permanent_store.dart
├── configs/    # importer.dart（共通 import 集約）
└── theme/      # tokens.dart（AppColors, AppFontSizes）

gas/{shift,task,user,rescue,manual-assignment}/   # ドメイン別。コード.js / onChange.js 等
```

`di/di.go` は部品の組み立て（Repository → UseCase → Controller → Router）だけを書く。定期実行やバックグラウンドの処理は `externals/` の下に専用パッケージを作り、`di.go` では組み立てて `Start()` を呼ぶだけにする（理由は `docs/decisions/0002-di-wiring-only.md`）。

## 通知の定期実行

未送信の通知ログ（`action_logs` の `is_sent = false`）は、`externals/scheduler` の ticker ループが API プロセス内で5分間隔に `NotificationUseCase.ProcessUnsentNotifications` を呼び、Slack DM へ flush します（`di.go` で組み立てて起動。`cmd/send-notifications` は手動 flush 用として併存）。

**API は単一インスタンス前提**です。複数レプリカで動かすと各プロセスの ticker が同じ未送信ログを拾って二重送信します。本番（`docker-compose.prod.yml`）は API を1レプリカで運用しているため現状は問題ありません。複数レプリカ化する場合はリーダー選出や排他制御が必要です。

レスキューの対応状況（未対応→対応中→対応済み）や本部からの返答が変わると、`*RescueUseCase.Update*Rescue` が更新前後の値を `RescueNotifier` に渡し、知らせるべき変化ならその場で送信者本人に Slack DM します（キューや scheduler は使いません。送信は goroutine で裏に回し、失敗しても送り直しません）。戻る変更（対応済み→未対応など）や返答を消しただけの変更は通知しません。GAS が押し直しの重複を「対応済み」にしてまとめる書き込みは `"notify": false` を付けて送り、通知しません。`SLACK_BOT_TOKEN` が無いか `RESCUE_NOTIFICATION_DISABLED=true` のときは通知しません。

## Code Style

### Go (`api/`)

**SQL は必ずプレースホルダで書く**

```go
query := "SELECT * FROM bureaus WHERE id = $1"
rows, err := db.QueryContext(ctx, query, id)
```

文字列連結（`"... " + id`）は SQL インジェクション脆弱性のため禁止。golangci-lint の gosec は SeeFT の連結 SQL（`abstract.Crud` 経由や、`DB()` からのメソッドチェーン）を検出できない。lint が通っても連結が無い証明にはならないので、レビューで確かめる。

**INSERT した行の id をあとで使うときは、`RETURNING id` で受け取る**

作成のあとに最新の行（`ORDER BY id DESC LIMIT 1` など）を読み直すと、同時に作られた別の行を掴む（#536）。id を使わない INSERT（`actionLogRepository.Create` など）は、これまでどおり実行するだけでよい。`user_repository.go` の `Create` などが今の書き方。

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

dispose 後に `setState` を呼ぶと例外になるため必須。失敗の分岐（`else`・`catch`）で `await` のあとに `setState` や `ScaffoldMessenger.of(context)` を呼ぶときも同じ。`flutter analyze` の `use_build_context_synchronously` が0件でも、失敗の分岐が抜けていることがある（PR #383）。

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

仕事の進め方の全体は `docs/development/workflow.md`。コードを書くときに要るものは次のとおり。

- ブランチ名: `feat/{username}/{issue-number}/{description}`、`fix/...`、`docs/...`。develop から切る（main は使っていない）
- コミットメッセージは日本語、`feat:` / `fix:` / `docs:` プレフィックス
- ドキュメント変更（AGENTS.md・README 等）も issue → branch → PR の正規フローを通す
- PoC は試作用のブランチで試し、完成形が見えたら develop から切った新しいブランチに要るファイルだけを持ち込んで PR にする
- ほかの人の PR を引き継ぐときは、元の PR のブランチにコミットを積む。develop は `git fetch origin` の後に `git merge origin/develop` で取り込み、force push しない
- PR は `.github/pull_request_template.md` のフォーマットに従う
- PR を出す前の点検は、差分全体を対象にする。誤りを直したら、同じ誤りがほかにもないか `git grep` で探して直す
- PR 本文で `resolve #XXX` と書くと issue が自動 close される。番号ごとに `gh issue view <N> --json title` で中身を確かめてから書く
- issue・PR でコードを引用するときは、言語を指定したコードブロックに入れる

## 判断の記録（ADR）

機能を足す・見送る、設計や運用の方針を選ぶ、Ask First の項目を決めた、といった判断をしたら、`docs/decisions/` に ADR を書く。書き方は `docs/decisions/README.md`、雛形は `docs/decisions/template.md`。

- 決めたことだけでなく、理由と、選ばなかった候補を書く。見送り・先送りも ADR にする
- ADR にするのは、システムの作りや本番の運用に関わり、後から戻すのに手間がかかる判断だけ。チームの約束事（PR の書き方、点検のしかた、タスクの割り振り方など）は ADR にせず、`docs/development/workflow.md` に理由と一緒に書く
- 候補を比べて決めるときは、決める前に状態を「提案」にして PR に出し、決まったら「採用」か「見送り」にする。Ask First の項目は、実装の前にこの流れで進める
- 決定の信頼度（高・中・低。後から起こした ADR で当時の信頼度が分からなければ「記録なし」）と、なぜその信頼度かの根拠を書く。低いときは、何が分かれば見直すかを「前提」に書く
- 理由が分からなければ「理由の記録なし」と書く。推測で埋めない
- 本文は書き換えない。判断が変わったら新しい番号で書き、古い方の状態を `置き換え（→ NNNN）` にする。ただし、コードの名前やパスが変わっただけで判断が変わらないときは、「前提」欄の参照先を直し、追記に日付と理由を1行書く
- 前提にしたコードは `path#Symbol` の形で書き、行番号は書かない。PR ごとに `scripts/refcheck/refcheck.py` が存在を確かめる
- 公開リポジトリなので、個人名（役職で書く）・スプレッドシートや Drive の ID・サーバーの IP・未公開のセキュリティ問題は書かない

## Boundaries

### Always Do
- 新規 SQL はプレースホルダで書く
- 設定値（API URL・シークレット）は環境変数 / `PropertiesService`（GASのみ）から取得
- INSERT した行の id をあとで使うときは、`RETURNING id` で受け取る
- Flutter で非同期またぎ後の `setState()` 前に `mounted` チェック（失敗の分岐も）
- GAS で `LockService` 取得後は `finally` で `releaseLock()`
- 空リストは `[]Type{}` を返す
- 機能の採否や設計・運用の方針を決めたら `docs/decisions/` に ADR を書く

### Ask First
- 新規ライブラリの導入（特に Flutter の状態管理系）
- 既存 entity の JSON キー命名変更（mobile / gas に影響）
- API レスポンス形式の変更（キーを変えるなら、クライアントは新旧どちらの形でも受けられるようにする。マージしても本番の反映までは古い API が動き続けるため。#479）
- DB スキーマ変更（マイグレーション）

### Never Do
- SQL を文字列連結で組み立てる
- 適用済みの migration のファイルを消す・書き換える（消すと、`schema_migrations` に記録された版のファイルが見つからず migrate が止まる。書き換えても、適用済みの DB には反映されない）
- API URL やシークレットをハードコードする
- `package:http` を `mobile/lib/utils/api.dart` 以外で import する
- `print()` を新規コードで使う（`mobile/lib/`）
- `mobile/lib/` に Riverpod / Provider / Bloc 等の状態管理ライブラリを導入する
- `.env` や認証情報をコミットする

## Known Transitional Issues

新旧の規約が混在する箇所があります。新規コードは上記ルールに従い、既存は段階的に整理します。

- **Repository の SQL 文字列連結**（残 7 ファイル / 約 33 件） → #266
- **JSON キー命名**: 古い entity は camelCase、新しい entity は snake_case。クライアント影響のため既存維持
- **エラーレスポンス**: 古い controller は `return err`、新しいものは map 形式
- **空リスト返却**: 一部 UseCase が `nil` を返す箇所あり（順次 `[]Type{}` へ）
- **作成後の最新行の読み直し**: place・task・review の UseCase が、作成後に `FindNewRecord` で最新の行を読み直している（新規コードは `RETURNING id`）
