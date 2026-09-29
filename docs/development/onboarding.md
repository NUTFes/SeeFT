# オンボーディング: 技術スタック

SeeFT に新しく入った人向けの資料です。Go も Flutter も Apps Script も初めて、という人を想定しています。

この文書が扱うのは**機能ではなく技術**です。何を使っているか、なぜそれがあるのか、それぞれについて何を学ぶべきか、学んだことを確かめるために自分で何を作るべきかを説明します。SeeFT が技大祭で何をしているかにはほとんど触れません。それは issue とレビューを通じて身につきますし、運用については [docs/operations/](../operations/README.md) にまとまっています。

上から順に読んでください。順番には意味があり、各節は前の節を前提にしています。コードの書き方の約束事は [AGENTS.md](../../AGENTS.md) にあります。この文書を読み終えてから AGENTS.md を読むと、なぜその約束事があるのかが分かるはずです。

---

## 0. システムの全体像

1つのリポジトリに、動くものが3種類入っています。

```text
  スプレッドシート ──► Apps Script (GAS) ──┐
  （シフト・名簿・タスク・レスキュー）       │ HTTP / JSON
                                            ▼
  スマホのブラウザ ──► mobile (Flutter Web) ──► api (Go + Echo) ──► PostgreSQL
                                                   │
                                                   └──► Slack（本人への DM）
```

データの入口は2つあります。

- **GAS** がスプレッドシートの内容（名簿・タスク・シフト）を API に送り込みます。シフトを組むのは人間で、その作業場所はスプレッドシートです。SeeFT はそれを受け取って配る側です。
- **mobile** は参加者がスマホのブラウザで開くアプリで、自分のシフトを見たり、レスキュー（人手不足やトラブルの連絡）を送ったりします。

api はどちらからも同じ HTTP API で呼ばれます。シフトが変わったときやレスキューの対応状況が変わったときは、api が Slack で本人に DM を送ります。

言語は場所ごとに違います。api は Go、mobile は Dart、GAS は JavaScript です。フロントエンドとバックエンドで型を共有する仕組みはありません。**api と mobile / GAS の間の約束は JSON のキー名だけ**で、それが崩れてもコンパイラは何も言いません。これが最初に腹落ちさせるべき最大のポイントです。AGENTS.md が「既存 entity の JSON キー命名変更」を事前相談にしているのはこのためです。

リポジトリの構成:

```text
api/                 Go の API サーバー（port 1234）
  main.go            起動するだけ。中身は lib/di
  lib/               本体（後述の3層構成）
  cmd/migrate        DB のテーブル作成・変更を流すコマンド
  cmd/seed           初期データを入れるコマンド
  cmd/send-notifications  未送信の Slack 通知を手動で流すコマンド
mobile/              Flutter Web のアプリ（port 45029）
  lib/               本体
  test/              widget テスト
  python/server.py   本番でビルド済みのファイルを配る小さなサーバー
gas/                 Apps Script の「写し」（6節で詳しく）
postgresql/db/       テーブル定義（schema/）、変更（migrations/）、初期データ（seed.sql）
docs/                この文書を含む資料
admin/, raw/         使っていない（読まなくてよい）
```

---

## 1. Go

**何か。** Google が作った、静的型付けでコンパイルする言語です。文法が小さく、書き方の選択肢が少ないので、他人のコードが読みやすいのが特徴です。バージョンは Go 1.26（`api/go.mod` と `api/Dockerfile` で固定）です。

**私たちの使い方。** api のすべてです。ローカルでは Docker コンテナの中で動かすので、手元に Go を入れなくても api は起動できます。ただしエディタの補完（gopls）のために、手元にも同じバージョンの Go を入れておくことを勧めます。

Python や JavaScript から来た人がつまずくのは、次の3つです。

- **例外が無い。** 失敗する関数は最後の戻り値に `error` を返し、呼んだ側が毎回 `if err != nil { return ..., err }` と書きます。冗長に見えますが、どこで失敗しうるかがコードに全部書いてあるということです。
- **クラスが無い。** `struct`（データ）にメソッドを生やし、`interface`（メソッドの一覧）で抽象化します。ある型が interface を満たすかどうかは、メソッドが揃っているかだけで決まります。`implements` のような宣言はありません。
- **大文字で始まる名前だけが外から見える。** `NewBureauController` はパッケージの外から呼べて、`bureauController` は呼べません。SeeFT では interface を大文字、実装の struct を小文字にしてこれを使い分けています。

**学ぶこと:**

- 変数、関数、複数の戻り値、`struct`、メソッド、ポインタ（`*T` と `&x`）
- `error` の扱い方、`errors.Wrap` でエラーに文脈を足すこと
- `interface` と、それを満たすとはどういうことか
- スライス（`[]T`）とマップ、`for ... range`
- パッケージと import、大文字・小文字による公開範囲
- `context.Context` が何のために関数の第1引数で引き回されているか
- goroutine と `select`（api 本体では Slack 送信とスケジューラで使う程度。深入りは後回しでよい）

**やってみる:** [A Tour of Go](https://go.dev/tour/) を最後まで。2〜3日あれば終わります。

---

## 2. Echo と3層構成（api）

**何か。** Echo は Go の Web フレームワークで、URL とメソッドの対応づけ（ルーティング）、JSON の読み書き、ミドルウェアを提供します。フレームワークとしては薄く、構造を決めているのは Echo ではなく私たちの規約です。

**3層構成。** `api/lib/` の下はこう分かれています。

| 層 | 場所 | 役割 |
| --- | --- | --- |
| **router** | `lib/router/router.go` | URL とメソッドを controller の関数に対応づける。全エンドポイントの一覧はここ |
| **controller** | `lib/internals/controller/` | HTTP の入出力だけ。リクエストを読み、usecase を呼び、JSON で返す。DB には触らない |
| **usecase** | `lib/usecase/` | 業務ロジック。repository から受け取った行を entity に詰める |
| **repository** | `lib/internals/repository/` | SQL を実行するだけ。結果は `*sql.Rows` のまま返す |
| **entity** | `lib/entity/` | データの形。JSON のキー名はここの struct タグで決まる |

リクエストは router → controller → usecase → repository → DB の順に下り、結果は逆順に上がってきます。

**依存性注入（DI）は手書きです。** `lib/di/di.go` を開くと、repository → usecase → controller → router の順に `NewXxx(...)` を呼んで、作ったものを次の層に渡しているのが分かります。

```go
// api/lib/di/di.go（抜粋）
bureauRepository := repository.NewBureauRepository(client, crud)
bureauUseCase := usecase.NewBureauUseCase(bureauRepository)
bureauController := controller.NewBureauController(bureauUseCase)
```

usecase は repository の**具体的な型ではなく interface** を受け取ります。だからテストでは本物の DB の代わりに偽物を渡せます（8節）。新しいエンドポイントを足すときは、4つの層にファイルを足し、di.go で配線し、router.go にルートを書く、という流れになります。di.go は配線だけに留め、業務の判断は usecase に書きます。

**repository が `*sql.Rows` を返し、usecase が `Scan` する**のは SeeFT の少し変わったところです。

```go
// api/lib/usecase/bureau_usecase.go（抜粋）
rows, err := a.rep.All(c)          // repository は SQL を実行するだけ
...
for rows.Next() {
    err := rows.Scan(&bureau.ID, &bureau.Bureau, &bureau.Color, ...)  // 詰めるのは usecase
```

repository の多くは `SELECT *` で、`Scan` の引数はテーブルの列の順番と一致していなければなりません。列を1つ足すと、そのテーブルを `Scan` している箇所をすべて直す必要があります。直し漏れはコンパイルでは見つからず、実行時に `Scan` がエラーを返します。

**横断的に使っている仕組み:**

- **ミドルウェア**（`lib/externals/server/server.go`）: パニックからの復帰、アクセスログ、CORS。CORS は許可するオリジン（ドメインとポート）を列挙しているので、mobile を 45029 以外のポートで動かすとブラウザが API の応答を拒みます。
- **Slack 通知**（`lib/externals/slack/`、`lib/externals/scheduler/`）: シフト変更の通知は、変更を `action_logs` に記録しておき、5分ごとに動くループがまとめて DM します。レスキューの通知はその場で送ります。詳しくは AGENTS.md の「通知の定期実行」。
- **ホットリロード**: ローカルのコンテナは [air](https://github.com/air-verse/air) で起動しているので、`.go` ファイルを保存すると自動で再ビルドされます。`di.go` を変えてうまく動かないときは、コンテナを一度 down して up し直してください。

**新旧の書き方が混在しています。** 古いファイルには、今は禁止している書き方が残っています。手本にするファイルを間違えないでください。

- SQL を文字列連結で組み立てている（`bureau_repository.go` など）。新しいコードはプレースホルダ（`$1`）で書きます。`review_repository.go` がプレースホルダで書き直された例です
- エラーを `return err` で Echo に任せている（`bureau_controller.go` など）。新しいコードは `c.JSON(status, map[string]string{"error": ...})` で返します。`review_controller.go` が新しい書き方です
- 空の一覧で `nil` を返している。新しいコードは `[]entity.Xxx{}` を返します（`nil` だと JSON が `null` になり、mobile が壊れる）

**学ぶこと:**

- HTTP そのもの: メソッド（GET / POST / PUT / DELETE）、ステータスコード、ヘッダー、JSON ボディ
- Echo: `e.GET(...)` によるルーティング、`c.Param`、`c.Bind`、`c.JSON`
- interface を受け取る設計と、それがテストのしやすさにつながる理由
- CORS とは何か、なぜブラウザだけが気にするのか

**どこを見るか。** `review` の一式（`review_controller.go` → `review_usecase.go` → `review_repository.go` → `entity/review.go`）が、新しい書き方で CRUD が揃った例です。repository にはテスト（`review_repository_sqlmock_test.go`）もあります。

エンドポイントの一覧は `router.go` を読んでください。`/swagger/index.html` にも API ドキュメントがありますが、2023年から更新されておらず今のエンドポイントと合っていません。

---

## 3. PostgreSQL と migration

**PostgreSQL** がデータベースです。ローカルでは Docker の `postgres:18` で動きます。api からは Go 標準の `database/sql` で SQL を直接書いてアクセスします。ORM（GORM）も入っていますが、使っているのは `shift_card_repository.go` などごく一部です。**SeeFT で DB を触るとは、ほぼ SQL を書くことです。**

**テーブル定義は2か所にあります。**

1. **`postgresql/db/schema/`**: `create1Bureaus.sql` のような `CREATE TABLE` のファイル。ファイル名の数字が実行順で、外部キーで参照される側ほど番号が小さくなっています（`create1Bureaus` → `create2Users` → `create4Shifts`）。
2. **`postgresql/db/migrations/`**: 既存のテーブルへの変更（列の追加など）。`2026082201_add_manual_url_to_tasks.up.sql` と `.down.sql` の組で書き、[golang-migrate](https://github.com/golang-migrate/migrate) が流します。

`make migrate` は schema を番号順に流し、続けて migrations を流します。どちらも「どのファイルを適用済みか」を DB に記録しているので、何度実行しても同じファイルは二度流れません。

**適用済みのファイルは編集しないこと。** schema はファイルの中身のハッシュも記録していて、適用済みのファイルが書き換わっていると `make migrate` がエラーで止まります。変更は migrations に新しいファイルを足してください（列の追加の置き場所が schema と migrations に分かれている件は #461 で整理中です）。DB スキーマの変更は AGENTS.md で「事前相談」の対象です。

**学ぶこと:**

- SQL の基礎: `SELECT`, `JOIN`, `WHERE`, `GROUP BY`, `ORDER BY`、`INSERT` / `UPDATE` / `DELETE`
- 主キーと外部キー、1対多、なぜ中間テーブルが要るのか
- プレースホルダ（`$1`）と、文字列連結がなぜ SQL インジェクションになるのか。`review_repository_sqlmock_test.go` の冒頭にある「アポストロフィを含むコメント」のテストが、連結だと何が壊れるかの具体例です
- `database/sql`: `QueryContext`（複数行）、`QueryRowContext`（1行）、`ExecContext`（更新）、`Rows.Scan`
- インデックスとは何か。`shifts` はインデックスが無く、件数が増えると遅くなる課題があります（#500）
- `psql` でテーブルを眺めること。`make exec` ではなく、`docker compose exec db psql -U seeft -d seeft_db` で DB のコンテナに入れます（抜けるときは `\q`）

**どこを見るか。** まず `postgresql/er.png` で全体の関係を見て、`schema/` を番号順に読んでください。`users`、`tasks`、`shifts` の3つが中心で、残りの多くは「局」「学年」「天気」「時間帯」のような選択肢の表です。

---

## 4. Dart と Flutter Web（mobile）

**何か。** Flutter は Google の UI フレームワークで、言語は Dart です。1つのコードから Android / iOS / Web を作れますが、**SeeFT は Web としてだけ動かしています。** 参加者はアプリをインストールせず、スマホのブラウザで URL を開きます。`mobile/android/` や `mobile/ios/` は残っていますが使っていません。

**Flutter のバージョンは [fvm](https://fvm.app/) で固定しています**（`mobile/.fvmrc` で 3.27.3）。`flutter` コマンドは必ず `fvm flutter ...` の形で打ってください。素の `flutter` を打つと手元に入っている別のバージョンが動き、依存の解決やビルド結果が変わります。Flutter 自体を上げる作業は、他の変更と混ぜず専用のブランチで行います（#514）。

**画面は Widget の木です。** 画面のすべて（文字、ボタン、余白、並べ方）が Widget で、Widget の中に Widget を入れて画面を組み立てます。

- `StatelessWidget`: 引数だけで見た目が決まるもの
- `StatefulWidget`: 自分で状態を持ち、`setState` で状態を変えると描き直されるもの

```dart
// 状態を変えるときは setState の中で。Flutter が build を呼び直す
setState(() => _isLoading = false);
```

**設定値はビルド時に焼き込まれます。** API の URL などは `mobile/env/.env` に書き、`--dart-define-from-file=env/.env` で渡しています。コードからは `String.fromEnvironment('API_BASE_URL')`（`lib/configs/constant.dart`）で読みますが、これは**ビルドした瞬間に値が埋め込まれる**仕組みです。値を変えたら、ビルドし直すまで反映されません。

**学ぶこと:**

- Dart: 型、`final` と `const`、null 安全（`String?`、`!`、`??`）、`async` / `await` と `Future`
- Widget の木、`build` メソッド、`StatelessWidget` と `StatefulWidget` の違い
- レイアウト: `Column`, `Row`, `Padding`, `Expanded`, `ListView`
- `setState` と `initState` / `dispose` の役割
- `Navigator` による画面遷移
- ブラウザの開発者ツールで、Flutter Web が出す通信（Network タブ）を見ること

**やってみる:** 公式の [Write your first Flutter app](https://docs.flutter.dev/get-started/codelab) を、Web を対象にしてやってください。

---

## 5. mobile の構成と状態の持ち方

`mobile/lib/` の下はこう分かれています。

```text
pages/     画面。1ファイル1画面（my_shift_page.dart がマイシフト画面）
widgets/   複数の画面で使う部品（shift_card.dart がシフトの1コマ）
models/    データの形
utils/     api.dart（API 呼び出し）、logger.dart、permanent_store.dart（保存）
configs/   constant.dart（設定値）、importer.dart（よく使う import をまとめたもの）
theme/     tokens.dart（色・文字サイズ）
```

**API はシングルトン経由で呼びます。** `lib/utils/api.dart` の `api` という1つのオブジェクトに、エンドポイントごとの関数（`api.getShiftCardsByUserAndDateAndWeather(...)`、`api.postReview(...)` など）が並んでいます。画面側は `import 'package:seeft_mobile/configs/importer.dart';` だけで `api` を使えます。HTTP の細部を1か所に閉じ込めるため、`package:http` を `api.dart` 以外から import するのは禁止です。

**状態管理ライブラリは使いません。** Riverpod や Provider のような状態管理ライブラリは入れず、各画面の `State` クラスの中に状態を持ち、`setState` で描き直します。導入は AGENTS.md で禁止しています。画面数が少なく、画面をまたいで共有する状態がほとんど無いので、それで足りています。

**非同期の後は `mounted` を確かめます。** API の応答を待っている間に利用者が別の画面に移ると、その画面の `State` は捨てられています。捨てられた `State` で `setState` を呼ぶと例外になります。

```dart
final data = await api.fetchData();
if (!mounted) return;   // 待っている間に画面が閉じられていたら何もしない
setState(() => _data = data);
```

`manual_list_page.dart` の `_loadData` も、同じ考え方で `mounted` を確かめてから `setState` しています。

**端末への保存は2種類あります**（`lib/utils/permanent_store.dart`）。

- **SharedPreferences**: ログイン中のユーザー ID のような小さな値
- **Hive**: シフトカードやレスキューの一覧のような、まとまったデータ。マイシフト画面は、前回取得したシフトカードを日付・天気ごとに Hive から読み出しています（`my_shift_page.dart`）

Web なので、どちらも実体はブラウザのストレージです。検証中に表示が古いままのときは、ブラウザのサイトデータを消すと直ることがあります。

**見た目の決まり。** 色と文字サイズは `AppColors.main` や `AppFontSizes.md`（`lib/theme/tokens.dart`）を使い、`Color(0xFF...)` を直接書きません。ログは `print` ではなく `logger.i(...)` / `logger.e(...)` で出します。

**学ぶこと:**

- `StatefulWidget` の一生: `initState` → `build` → `dispose`
- `FutureBuilder` と、`initState` で Future を作っておく理由（`widgets/first_jump_selector.dart` のコメントに書いてあります）
- `mounted` がなぜ必要か
- SharedPreferences と Hive の使い分け

**どこを見るか。** 起動の流れは `main.dart` → `widgets/first_jump_selector.dart`（ログイン済みかを見て、ログイン画面かマイシフト画面に振り分ける）です。画面を1つ読むなら `pages/manual_list_page.dart` が、API 呼び出し・読み込み中表示・エラー表示・検索が1画面に収まっていて手頃です。

---

## 6. Google Apps Script（GAS）と clasp

**何か。** Google のスプレッドシートに紐づけて動かせる JavaScript です。SeeFT では、スプレッドシートに作ったメニューから名簿・タスク・シフトを API に送ったり、レスキューを受け付けたりしています。

**`gas/` は実物ではなく写しです。** GAS の本体は各スプレッドシートに紐づいたクラウド上のプロジェクトで、リポジトリにあるのは [clasp](https://github.com/google/clasp)（GAS をコマンドラインから取得・反映するツール）で取ってきた、ある時点のコピーです。誰かがスプレッドシートのエディタで直接編集すると、リポジトリは変わらないまま本物だけが進みます。**GAS を触る作業は、リポジトリを読むことではなく、`clasp` で本物を取ってくることから始めます。** 手順は [gas/README.md](../../gas/README.md) にあります。

**Go や Dart と違って、GAS には CI が無く、型もありません。** 壊れたコードを push しても誰も止めてくれません。しかも全ファイルが同じスコープで読み込まれるので、別ファイルで同じ名前の `const` を定義しただけでスクリプト全体が読み込めなくなり、スプレッドシートのメニューごと消えます。push の前に `node --check` で構文を確かめるのはこのためです。

GAS 特有の道具:

- `UrlFetchApp.fetch(url, options)`: API を呼ぶ
- `PropertiesService.getScriptProperties()`: API の URL やトークンの置き場。コードに直接書かない
- `LockService`: 同じスクリプトが同時に2回動かないようにする鍵。取ったら `finally` で必ず放す
- `doPost(e)`: GAS をウェブアプリとして公開し、外から POST を受ける入口（レスキューで使用）
- `onChange(e)`: スプレッドシートが変更されたときに動くトリガー

**学ぶこと:**

- JavaScript の基礎: `const` / `let`、関数、配列の `map` / `filter`、オブジェクト、`JSON.stringify`
- スプレッドシートの操作: `SpreadsheetApp`、`getRange`、`getValues` / `setValues`（1セルずつ読み書きすると遅いので、範囲でまとめて読む）
- GAS の実行時間の上限（6分）と、それを超えたときに何が起きるか。シフトの送信は件数が多く、この上限に当たります（[docs/operations/](../operations/README.md) に記録があります）
- clasp: `clone`、`pull`、`push`、`deploy`
- ウェブアプリの「デプロイ」と「バージョン」。`clasp push` しただけでは公開中のウェブアプリは変わりません

**どこを見るか。** `gas/rescue/コード.js` の `doPost` が、外からの POST を受けてスプレッドシートに書き、API にも送る一連の流れです。

---

## 7. データの流れを1本追う

リポジトリを読むときは、1つの操作を端から端まで追ってください。例として、mobile でタスクのレビュー（人手の足り具合とマニュアルの評価）を送る操作を追います。

1. `mobile/lib/widgets/review_bottom_sheet.dart`: 送信ボタンで `api.postReview(...)` を呼ぶ
2. `mobile/lib/utils/api.dart`: `POST /reviews` に JSON を送る
3. `api/lib/router/router.go`: `e.POST("/reviews", r.reviewController.CreateReview)`
4. `api/lib/internals/controller/review_controller.go`: `CreateReview` が JSON を読んで usecase を呼ぶ
5. `api/lib/usecase/review_usecase.go`: `CreateReview` がタスク名から ID を引き、repository に保存させる
6. `api/lib/internals/repository/review_repository.go`: `INSERT INTO reviews ...` をプレースホルダで実行
7. `postgresql/db/schema/create4Review.sql`: `reviews` テーブルの定義

JSON のキー名（`user_id`、`task_name` など）が、Dart 側で組み立てるときと Go 側の struct タグで一致していることを確かめてください。これが0節で書いた「約束は JSON のキー名だけ」の実物です。

これを一度、丁寧にやることは、どんなチュートリアルよりも多くを教えてくれます。慣れたら、GAS からシフトを送る流れ（`gas/shift/コード.js` → `POST /api/update_shifts` → `shift_controller.go` の `UpdateShiftsFromGAS` → `shift_usecase.go`）も追ってみてください。こちらはずっと大きく、SeeFT の中心です。

---

## 8. テスト

**api（Go）。** `go test` で動く標準のテストに、[testify](https://github.com/stretchr/testify)（`assert` / `require`）と [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) を足して書いています。テストファイルは対象と同じディレクトリに `xxx_test.go` として置きます。

repository も usecase も、本物の DB ではなく go-sqlmock の偽の DB を相手にテストします。「この SQL がこの引数で呼ばれるはず」と先に宣言しておき、呼ばれた中身を確かめる形です。2節の DI が、ここで偽物を差し込む口になっています。

```go
// api/lib/internals/repository/review_repository_sqlmock_test.go（抜粋）
mock.ExpectExec(`INSERT INTO reviews \(user_id, task_id, staffing_rating, manual_rating, comment\)`).
    WithArgs("12", "34", "5", "4", testReviewComment).
    WillReturnResult(sqlmock.NewResult(1, 1))
```

新しいテストは `sqlmock_helper_test.go` の `newDBMock` を使うと短く書けます。どの層からどう書き進めるかの方針は [test-roadmap.md](test-roadmap.md) にあり、テストが書かれていない関数の issue（`[test]` で始まるもの）が練習に向いています。

**mobile（Flutter）。** `flutter_test` で widget テストを書いています。画面の部品を仮想的に描画して、表示や押したときの動きを確かめます。

```bash
cd mobile && fvm flutter test --platform chrome
```

`--platform chrome` が必須です（`-p` と略すことはできません）。付けないとテストが Web 向けのコードを読み込めずに落ちます。`test/widget_test.dart` は Flutter が最初に作ったまま壊れているので、落ちても気にしないでください（#493 で除去予定）。

**学ぶこと:**

- Go: `func TestXxx(t *testing.T)`、テーブル駆動テスト、`t.Run`、`assert` と `require` の違い
- go-sqlmock: `ExpectQuery` / `ExpectExec`、`WithArgs`、`WillReturnRows`、`ExpectationsWereMet`
- Flutter: `testWidgets`、`pumpWidget`、`find.text`、`tester.tap`
- 何をテストする価値があるか。境界の値や、過去に事故が起きた入力（アポストロフィを含む文字列、空の値）

**GAS には自動テストがありません。** 検証環境のスプレッドシートで実際に動かして確かめます（[staging-rehearsal.md](../operations/staging-rehearsal.md)）。

---

## 9. ツール、Docker、CI

- **Docker / Docker Compose**: ローカルでは DB と api がコンテナで動きます（`docker-compose.yml`）。mobile はコンテナではなく手元の fvm で動かします。本番も Docker Compose で、api・mobile・外部公開用のトンネルがコンテナで動いています（`docker-compose.prod.yml`、手順は [deploy.md](../operations/deploy.md)）。
- **Makefile**: よく使う docker compose のコマンドに短い名前を付けたものです。一覧は AGENTS.md の Commands にあります。Mac 用の compose ファイル（`docker-compose.mac.yml`）を使う `mac-` 付きのコマンドもあります。
- **GitHub Actions**: PR を出すと、変更した場所に応じて次が走ります。
  - `api/` を変えたとき: golangci-lint（`go-lint.yml`、PR で増えた指摘だけを見る）と `go test`（`go-test.yml`）
  - `mobile/` を変えたとき: `flutter analyze --fatal-infos`（`flutter-lint.yml`、info レベルの指摘でも落ちる）
  - `gas/` を変えたとき: 何も走らない
- **CodeRabbit**: PR に AI がレビューコメントを付けます。指摘は参考です。すべてに従う必要はなく、スコープ外のものは理由を書いて見送って構いません。
- **Git**: issue を立て、`feat/{ユーザー名}/{issue番号}/{内容}` のブランチで作業し、PR を出します。コミットメッセージは日本語で `feat:` / `fix:` / `docs:` を付けます。詳しくは AGENTS.md の Git Workflow。

**学ぶこと:**

- Docker の基礎: image と container の違い、`docker compose up / down / logs / exec`、volume、port
- lint のエラーを読み、黙らせるのではなく原因を直すこと
- workflow ファイルを読める程度の YAML
- Git: ブランチ、merge、コンフリクトの解消、1つの PR に1つの目的

---

## 10. 手元で動かす

最初に一度、全部を手元で動かしてください。Docker Desktop、fvm、Git が要ります。以下はリポジトリのルートで実行します。

DB と api を用意する:

```bash
make build
make migrate
make seed
make up-api
```

`make up-api` はログを出し続けるので、そのターミナルは開いたままにします。ブラウザで http://localhost:1234/bureaus を開いて、JSON が返ってくれば api は動いています。

別のターミナルで mobile を起動する（初回だけ、設定ファイルの雛形をコピーし、Flutter を入れます）:

```bash
cp mobile/env/.env.example mobile/env/.env
cd mobile && fvm install && cd ..
make mobile-up
```

http://localhost:45029 を開くとログイン画面が出ます。ポートは 45029 から変えないでください（2節の CORS の理由で、API の応答をブラウザが拒みます）。

片付けるときは `make down`。DB の中身ごと消したいときは `docker compose down -v` です。

---

## 11. 現実的な学習順序

全部を一度に学ぼうとしないでください。担当する場所から始めれば十分です。api と mobile の両方を触る人でも、だいたい2〜3週間です。

| 段階 | 重点 |
| --- | --- |
| 1 | 10節で手元で動かす。Git と GitHub の流れ（issue → ブランチ → PR）。 |
| 2 | 担当が api なら Go の基礎（A Tour of Go）、mobile なら Dart と Flutter の基礎。 |
| 3 | SQL の基礎と `schema/` の読解。`psql` でテーブルを眺める。 |
| 4 | api: 3層構成と DI。`review` の一式を読む。mobile: 画面1つ（`manual_list_page.dart`）を読む。 |
| 5 | 7節のデータの流れを1本追う。 |
| 6 | テスト（go-sqlmock または widget テスト）を1本書く。Docker と CI の基礎。 |
| 7 | GAS と clasp（GAS に関わる issue を持ったときでよい）。 |

---

## 12. 何か作る

読むだけでは足りません。同じ技術を使った**自分の小さなアプリを、新しいリポジトリで**作ってください。SeeFT とは関係ないもの（本棚、習慣の記録、サークルの備品の貸し出し表など）で構いません。題材は重要ではなく、重要なのは各部分のつなぎ込みです。

### 最低限

**api: Go + Echo**

- Docker Compose 上の PostgreSQL
- 関連する2つ以上のテーブル（1対多）を `CREATE TABLE` の SQL で定義し、リポジトリにコミットする
- controller / usecase / repository / entity の4つに分け、手書きの DI（SeeFT の `di.go` と同じ形）で配線する
- 1つのテーブルについて一覧・1件取得・作成・更新・削除のエンドポイント
- SQL はすべてプレースホルダで書く
- 見つからないときは 404、入力が不正なときは 400 を、`{"error": "..."}` の JSON で返す
- 空の一覧は `null` ではなく `[]` を返す
- go-sqlmock を使った repository のテストを1本以上

**アプリ: Flutter Web**

- 一覧画面と、一覧から開く詳細画面
- API 呼び出しを1つのファイル（SeeFT の `api.dart` のような）にまとめる
- 読み込み中の表示と、失敗したときの表示
- 作成フォームを1つ。送信後に一覧が新しい内容で表示される
- `await` の後の `setState` の前に `mounted` を確かめる
- API の URL は `--dart-define-from-file` で渡し、コードに直接書かない
- widget テストを1本以上（`fvm flutter test --platform chrome` で通ること）

### 発展（時間があれば）

- スプレッドシートを1つ作り、GAS からあなたの API にデータを送るメニューを付ける。URL は `PropertiesService` に置き、`LockService` で二重実行を防ぐ
- 既存のテーブルに列を足す migration を golang-migrate で書き、up と down の両方を試す
- 一覧のデータを Hive にキャッシュし、API が落ちているときは前回の内容を表示する
- 一定間隔で動く処理（SeeFT の `scheduler` のような ticker ループ）を api に足す
- PR ごとに lint とテストを走らせる GitHub Actions の workflow

### 「完了」の意味

次を、調べずに説明できること。

- SeeFT の usecase が、repository の具体的な型ではなく interface を受け取るのはなぜか
- テーブルに列を1つ足したとき、Go のコードで直す必要があるのはどこか。直し忘れるといつ気づくか
- SQL を文字列連結で組み立てると、どんな入力で何が起きるか
- `await` の後で `mounted` を確かめないと何が起きるか
- `mobile/env/.env` の値を変えたのに画面が変わらないのはなぜか
- `gas/` のコードを読んで「本番はこう動いている」と判断してはいけないのはなぜか

これらに答えられるようになったら、本物の issue に取り組む準備ができています。

---

## 13. ブックマーク

- A Tour of Go — https://go.dev/tour/
- Effective Go — https://go.dev/doc/effective_go
- Echo — https://echo.labstack.com/docs
- `database/sql` のチュートリアル — https://go.dev/doc/tutorial/database-access
- go-sqlmock — https://github.com/DATA-DOG/go-sqlmock
- golang-migrate — https://github.com/golang-migrate/migrate
- PostgreSQL チュートリアル — https://www.postgresql.org/docs/current/tutorial.html
- Dart の言語ツアー — https://dart.dev/language
- Flutter — https://docs.flutter.dev/
- fvm — https://fvm.app/
- Hive — https://pub.dev/packages/hive_flutter
- Apps Script — https://developers.google.com/apps-script
- clasp — https://github.com/google/clasp
- Docker Compose — https://docs.docker.com/compose/

リポジトリ内:

- [AGENTS.md](../../AGENTS.md)：コードの書き方の約束事
- [docs/operations/](../operations/README.md)：技大祭での運用と、過去の事故の記録
- [test-roadmap.md](test-roadmap.md)：テストの進め方
- [gas/README.md](../../gas/README.md)：GAS の取得と反映

---

## 14. 助けの求め方

しばらく自分で試してから、聞いてください。丸一日詰まったままにしないこと、そして2分で聞かないこと。聞くときは次を含めてください。

- 何をしようとしていたか
- 何が起きると思っていたか
- 実際には何が起きたか
- 正確なエラー文（スクリーンショットより、コピーした文字のほうが検索できて助かります）
- すでに試したこと

この形なら早く答えが返ってきますし、issue を書くときの書き方でもあります。詰まっている issue があるなら、Slack で PM に連絡してください。
