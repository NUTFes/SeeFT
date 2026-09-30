# オンボーディング: 技術スタック

SeeFT に新しく入った人向けの資料です。Go も Flutter も Apps Script も初めて、という人を想定しています。

この文書が扱うのは**機能ではなく技術**です。何を使っているか、なぜそれがあるのか、それぞれについて何を学ぶべきか、学んだことを確かめるために自分で何を作るべきかを説明します。SeeFT が技大祭で何をしているかにはほとんど触れません。機能については、issue に取り組み、レビューを受けるうちに分かってきます。運用の手順は [docs/operations/](../operations/README.md) にまとまっています。

## この資料の読み方

初回は上から通して読んでください。前の節の言葉を後の節で使うので、順番に意味があります。分からない言葉は、まず1節の早見表を見てください。

通読したら、あとは自分の担当に合わせて深く学びます。担当（api・mobile・GAS のどれを触るか）は、PM が issue を割り振るときに決まります。各節の見出しの後ろに、その節がとくに必要な担当を書きました。「全員」の節は担当にかかわらず読んでください。

コードの書き方の規約は [AGENTS.md](../../AGENTS.md) にあります。AGENTS.md はもともと AI のコーディング支援ツール向けに書いたものですが、人が読んでもそのまま SeeFT の規約です。この資料を読み終えてから読むと、なぜその規約があるのかが分かるはずです。

---

## 0. システムの全体像（全員）

1つのリポジトリに、動くものが3種類入っています。

```text
  スプレッドシート（シフト・名簿・タスク）
        │ メニューから送信
        ▼
  GAS ─────────────────┐
        ▲              │ HTTP / JSON
        │ レスキュー    ▼
        └──────────── api (Go + Echo) ◄──── mobile (Flutter Web) ◄── 参加者のスマホのブラウザ
                        │   ▲
                        │   └─ レスキューの対応状況（スプレッドシートの変更を GAS が送る）
                        ├──► PostgreSQL
                        └──► Slack（本人への DM）
```

データの流れは3本あります。

- **シフト・名簿・タスク**：シフトを組むのは人間で、その作業場所はスプレッドシートです。GAS（スプレッドシートに付いたスクリプト）が、その内容を api に送ります。SeeFT は、組まれたシフトを受け取り、参加者に表示・通知する側です。
- **シフトの閲覧**：参加者はスマホのブラウザで mobile を開き、mobile が api から自分のシフトを取ってきて表示します。
- **レスキュー**（当日の人手不足やトラブルの連絡）：mobile から api に送られ、api がそれを GAS に送り、GAS がスプレッドシートに書き込みます。本部の人がスプレッドシート上で対応状況を変えると、GAS がそれを api に送り返します。

シフトが変わったときやレスキューの対応状況が変わったときは、api が Slack で本人に DM を送ります。

3つは動く場所も違います。

| | 書く言語 | 動く場所 | リポジトリにあるもの |
| --- | --- | --- | --- |
| api | Go | サーバー上の Docker コンテナ | 本体のコード |
| mobile | Dart | 参加者のブラウザ（サーバーはビルド済みのファイルを配るだけ） | 本体のコード |
| GAS | JavaScript | Google のクラウド | ある時点の写し（8節） |

言語が場所ごとに違うので、api と mobile / GAS の間で型を共有する仕組みはありません。両者の間の約束は、**どの URL にどのメソッドで、どんな形の JSON を送り、何が返ってくるか**だけです。JSON の形とは、キー名、値の型（文字列か数値か）、値が `null` になりうるか、配列かどうかです。キー名が1文字食い違っても、数値を文字列で送っても、コンパイルエラーにはならず、動かして初めて壊れていると分かります。最初に理解してほしいのはこの点です。AGENTS.md が「既存 entity の JSON キー命名変更」や「API レスポンス形式の変更」を事前相談にしているのはこのためです。

環境は3つあります。この資料で「ローカル」と書いたら自分の PC、「本番」と書いたら技大祭で参加者が使うサーバーです。その間に、本番に入れる前に通しで試す検証環境があります（[staging-rehearsal.md](../operations/staging-rehearsal.md)）。ローカルは何度壊してもやり直せますが、本番の DB やスプレッドシートは壊すと取り返しがつきません。

リポジトリの構成:

```text
api/                 Go の API サーバー（ローカルでは port 1234）
  main.go            起動するだけ。中身は lib/di
  lib/               本体（5節）
  cmd/migrate        DB のテーブルを作る・変えるコマンド（4節）
  cmd/seed           初期データを入れるコマンド
  cmd/send-notifications  未送信の Slack 通知を手動で送るコマンド
mobile/              Flutter Web のアプリ（ローカルでは port 45029）
  lib/               本体（7節）
  test/              画面の部品のテスト
  python/server.py   本番でビルド済みのファイルを配る小さなサーバー
gas/                 Apps Script の写し（8節）
postgresql/db/       テーブル定義（schema/）、既存テーブルへの変更（migrations/）、初期データ（seed.sql）
docs/                この資料を含む文書
admin/               使っていない（読まなくてよい）
```

---

## 1. 用語の早見表（全員）

この資料で説明なしに使う言葉です。最初は全部覚えなくてよいので、出てきたときに戻ってきてください。

| 言葉 | 意味 |
| --- | --- |
| issue | GitHub 上のタスク票。やることと理由が書いてある |
| PR（プルリクエスト） | 自分の変更を取り込んでもらうための申請。ここでレビュー（他の人による確認）を受ける |
| CI | PR を出すたびに GitHub が自動で走らせる検査（11節） |
| lint | コードの書き方の問題を機械的に見つけるツール |
| Docker、コンテナ | アプリを、必要なものごと箱に詰めて動かす仕組み。その箱がコンテナ |
| make | よく使う長いコマンドに短い名前を付けて呼ぶ道具。定義はリポジトリ直下の `Makefile` |
| ポート | 1台の PC の中で、どのプログラムに通信を届けるかを区別する番号 |
| HTTP、API | HTTP はブラウザとサーバーがやり取りする決まり。API はプログラムから呼べる窓口で、SeeFT では `GET /reviews` のような URL とメソッドの組（エンドポイント）の集まり |
| JSON | `{"user_id": 1, "comment": "..."}` のような、データをやり取りするための書き方 |
| ライブラリ、フレームワーク | 他人が書いた部品。フレームワークはプログラムの骨組みまで決めるもの |
| ビルド | ソースコードを、動かせる形（実行ファイルや配信用ファイル）に変換すること |
| デプロイ | ビルドしたものを本番のサーバーに置いて動かすこと |
| migration | DB のテーブルの作成や変更を、ファイルに書いて順に適用するやり方（4節） |
| seed | 開発用に DB に入れる初期データ |
| モック、偽物 | テストで本物の代わりに差し込む作り物（10節） |
| 天気（weather） | 技大祭のシフトは晴れ用と雨用の2通りを組むので、シフトは日付と天気で引く |

表記の約束として、小文字の **api** は `api/` ディレクトリとそのサーバーを、大文字の **API** はサーバーが外に見せている窓口（エンドポイントの集まり）を指します。**mobile** は `mobile/` とそのアプリのことです。

---

## 2. 手元で動かす（全員）

読み進める前に、一度 SeeFT を手元で動かしてください。動くものを触りながら読むほうが、ずっと頭に入ります。

用意するもの:

- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [fvm](https://fvm.app/)（Flutter のバージョン管理ツール。6節）
- Git
- Google Chrome（mobile のテストに使う）
- Node.js（GAS を触るときだけ）

Mac なら make は最初から入っています。

リポジトリを取ってくる:

```bash
git clone https://github.com/NUTFes/SeeFT.git
cd SeeFT
```

以下はリポジトリのルートで実行します。まず DB と api を用意します。

```bash
make build
make migrate
make seed
make up-api
```

`make build` は使っていない admin もビルドするので、初回は時間がかかります。`make up-api` はログを出し続けるので、そのターミナルは開いたままにします。ブラウザで http://localhost:1234/bureaus を開いて、JSON が返ってくれば api は動いています。

次に、別のターミナルで mobile を起動します。初回だけ、設定ファイルの雛形をコピーし、Flutter を入れます。

```bash
cp mobile/env/.env.example mobile/env/.env
cd mobile && fvm install && cd ..
make mobile-up
```

ブラウザで **http://localhost:45029** を開くとログイン画面が出ます。`127.0.0.1` で開いたり、api の許可リストに無いポートに変えたりすると、api の応答をブラウザが画面に渡さず、何も表示されません（5節の CORS）。ログインに使うテスト用アカウントは PM に聞いてください。seed にはシフトが入っていないので、ログインできてもマイシフトは空です。

ローカルでは、次の機能は動かないのが正常です。

- Slack の通知：`SLACK_BOT_TOKEN` が無いので、起動時にログが1行出て無効になります
- マニュアルの閲覧（`/manuals`）：Google ログインの設定が無いので、ルートごと無効になります
- レスキューのスプレッドシートへの送信：送り先の GAS の URL が無いので失敗します

api のコンテナは起動のたびに `go mod tidy`（依存ライブラリの整理）を実行するので、起動にはネット接続が要ります。

片付けるときは `make down`。DB の中身ごと消したいときは `docker compose down -v` です。

---

## 3. Go（api 担当）

**何か。** Google が作った、静的型付け（変数の型をコンパイル時に検査する）の言語です。文法が小さく、書き方の選択肢が少ないので、他人のコードが読みやすいのが特徴です。バージョンは Go 1.26 で、`api/go.mod` と `api/Dockerfile` で指定しています。

**私たちの使い方。** api はすべて Go で書かれています。ローカルでは Docker コンテナの中で動かすので、手元に Go を入れなくても api は起動できます。ただしエディタの補完（VS Code の Go 拡張が使う gopls）のために、手元にも同じバージョンの Go を入れておくことを勧めます。[公式のダウンロードページ](https://go.dev/dl/) から入れられます。

Python から来た人がつまずくのは、次の3つです。

- **失敗は例外ではなく戻り値で知らせる。** 失敗しうる関数は最後の戻り値に `error` を返し、呼んだ側が毎回 `if err != nil { return ..., err }` と書きます。冗長に見えますが、どこで失敗しうるかがコードに全部書いてあるということです。`panic` という仕組みもありますが、これは続行できない異常のためのもので、普段のエラー処理には使いません。
- **クラスが無い。** `struct`（データの入れ物）にメソッドを定義し、`interface`（メソッドの一覧）で抽象化します。ある型が interface を満たすかどうかは、メソッドが揃っているかだけで決まります。Java の `implements` のような宣言はありません。
- **大文字で始まる名前だけがパッケージの外から使える。** `NewBureauController` は外から呼べて、`bureauController` は呼べません。SeeFT では interface を大文字、それを実装する struct を小文字にして、外からは interface 越しにしか触れないようにしています。

**学ぶこと：**

- 変数、関数、複数の戻り値、`struct`、メソッド、ポインタ（`*T` と `&x`）
- `error` の扱い方。エラーに「どこで失敗したか」を足すには、標準では `fmt.Errorf("...: %w", err)` を使います。SeeFT の古いコードでは外部ライブラリの `errors.Wrap`（`github.com/pkg/errors`）も使っています
- `interface` と、それを満たすとはどういうことか
- スライス（`[]T`）とマップ、`for ... range`
- パッケージと import、大文字・小文字による公開範囲
- `context.Context` が何のために関数の第1引数で引き回されているか（リクエストが打ち切られたことを下の処理に伝えるため）
- goroutine（軽量な並行処理）。自分で書く機会は少ないですが、Go の HTTP サーバー（標準ライブラリの `net/http`。Echo はその上で handler として動く）は、接続ごとに goroutine を立てて、複数のリクエストを並行に処理します。複数のリクエストから同時に触られる変数を作らない、という前提だけは覚えておいてください

**やってみる:** [A Tour of Go](https://go.dev/tour/) を最後まで。2〜3日あれば終わります。

---

## 4. PostgreSQL と SQL（api 担当）

**PostgreSQL** がデータベースです。ローカルでは Docker の `postgres:18` で動きます。api からは Go 標準の `database/sql` で SQL を直接書いてアクセスします。ORM（SQL を書かずにオブジェクトの操作で DB を扱うライブラリ）の GORM も入っていますが、使っているのは `shift_card_repository.go` の1ファイルだけです。**SeeFT で DB を触るとは、ほぼ SQL を書くことです。**

**テーブルを作るファイルと、変えるファイルがあります。**

1. **`postgresql/db/schema/`**：`create1Bureaus.sql` のような、テーブルを作る `CREATE TABLE` のファイル。ファイル名の数字が実行順で、外部キーで参照される側ほど番号が小さくなっています（`create1Bureaus` → `create2Users` → `create4Shifts`）。
2. **`postgresql/db/migrations/`**：既存のテーブルへの変更（列の追加など）。`2026082201_add_manual_url_to_tasks.up.sql`（適用するときの SQL）と `.down.sql`（取り消すときの SQL）の組で書き、[golang-migrate](https://github.com/golang-migrate/migrate) が適用します。

`make migrate` は schema を番号順に適用し、続けて migrations を適用します。何度実行しても、同じ変更が二重に適用されることはありません。ただし、覚えている中身が2つで違います。

- schema は、ファイルごとに適用済みかどうかと中身のハッシュ（中身から計算した指紋）を DB に記録します。適用済みのファイルを書き換えると、`make migrate` がエラーで止まります。
- golang-migrate は、適用済みの一番新しい番号だけを記録します。そのため、自分より新しい番号の migration が先に適用された環境では、自分の migration は番号が古いという理由で適用されずに飛ばされます。マージする前に、自分のファイルの番号（先頭の日時）が develop にある migration のどれよりも新しいことを確かめてください。

**適用済みのファイルは編集しないでください。** 変更は migrations に新しいファイルを足します。過去に schema 側へ列の追加を書いたファイル（`create6AddSlackUserIdToUsers.sql`）が1つ残っていますが、これは migration の仕組みを入れる前の名残で、整理中です（#461）。DB スキーマの変更は、AGENTS.md で事前相談の対象です。

**学ぶこと：**

- SQL の基礎：`SELECT`, `JOIN`, `WHERE`, `GROUP BY`, `ORDER BY`、`INSERT` / `UPDATE` / `DELETE`
- 主キーと外部キー、1対多、なぜ中間テーブルが要るのか
- `SELECT *` と、列を名前で並べる `SELECT id, name, ...` の違い（次の節で効いてきます）
- プレースホルダ（`WHERE id = $1` と書き、値は別に渡す）と、値を文字列連結で SQL に埋め込むとなぜ SQL インジェクションになるのか
- インデックスとは何か。`shifts` は主キー以外の検索に使う列（`user_id` など）にインデックスが無く、件数が増えると遅くなる課題があります（#500）
- `psql` でテーブルを眺めること。`docker compose exec db psql -U seeft -d seeft_db` で DB のコンテナに入れます（抜けるときは `\q`）

**どこを見るか。** まず `postgresql/er.png`（テーブルの関係を描いた ER 図）で全体を見て、`schema/` を番号順に読んでください。`users`、`tasks`、`shifts` の3つが中心で、残りの多くは局・学年・天気・時間帯のような選択肢の表です。

---

## 5. api の構成：Echo と層（api 担当）

**Echo とは。** Go の Web フレームワークで、URL とメソッドの対応づけ（ルーティング）、JSON の読み書き、ミドルウェア（すべてのリクエストの前後に挟む共通処理）を提供します。Echo が提供する機能は少なく、ディレクトリの分け方や層の構成は SeeFT の規約で決めています。

**層の構成。** リクエストは入口の router から、controller → usecase → repository の3つの層を順に下って DB に届き、結果は逆順に上がってきます。entity は層ではなく、どの層からも使うデータの型です。

| 名前 | 場所 | 役割 |
| --- | --- | --- |
| router（入口） | `lib/router/router.go` | URL とメソッドを controller の関数に対応づける。全エンドポイントの一覧はここ |
| controller | `lib/internals/controller/` | HTTP の入出力だけ。リクエストを読み、usecase を呼び、JSON で返す。DB には触らない |
| usecase | `lib/usecase/` | 業務の処理。repository から受け取った行を entity に詰める |
| repository | `lib/internals/repository/` | SQL を実行するだけ。結果は行（`*sql.Rows`）のまま返す |
| entity（共通の型） | `lib/entity/` | データの形。応答の JSON のキー名は主にここの struct タグ（`json:"user_id"`）で決まる |

リクエストの JSON のキー名は、entity ではなく controller の中で定義した struct で決めていることもあります（`review_controller.go` の `CreateReview` など）。

**部品の組み立ては `di.go` に手で書いてあります。** 部品が使う別の部品を、部品の中で作らずに外から渡すことを、依存性注入（DI）と呼びます。NestJS や Spring のように DI を自動でやってくれるフレームワークもありますが、SeeFT は `lib/di/di.go` で、repository → usecase → controller → router の順に `NewXxx(...)` を呼び、作ったものを次の層に渡しています。

```go
// api/lib/di/di.go（抜粋）
bureauRepository := repository.NewBureauRepository(client, crud)
bureauUseCase := usecase.NewBureauUseCase(bureauRepository)
bureauController := controller.NewBureauController(bureauUseCase)
```

外から渡すので、テストでは本物の代わりに偽物を渡せます。SeeFT のテストでは、一番下で DB につないでいる部品（`db.Client`）を偽の DB に差し替えています（10節）。新しいエンドポイントを足すときは、controller・usecase・repository（と必要なら entity）にファイルを足し、`di.go` で組み立て、`router.go` にルートを書く、という流れになります。`di.go` には組み立てだけを書き、業務の判断は usecase に書きます。

**repository は行を返し、usecase がそれを struct に詰めます。** ここは SeeFT の少し変わったところです。

```go
// api/lib/usecase/bureau_usecase.go（抜粋）
rows, err := a.rep.All(c)          // repository は SQL を実行して行を返すだけ
...
for rows.Next() {
    err := rows.Scan(&bureau.ID, &bureau.Bureau, &bureau.Color, ...)  // 行の各列を struct に詰めるのは usecase
```

repository の多くは `SELECT *` で全列を取っていて、`Scan` に渡す変数はテーブルの列の順番と一致していなければなりません。テーブルに列を1つ足すと、そのテーブルを `SELECT *` して `Scan` している箇所をすべて直す必要があります。直し漏れは、コンパイルでは見つかりません。実行時に、列の数が合わなければ `Scan` がエラーを返します。列の数が合っていて順番だけずれた場合は、型が変換できる限りエラーにならず、違う列の値が黙って入ります。

**横断的に使っている仕組み：**

- **ミドルウェア**（`lib/externals/server/server.go`）：処理中の予期しない異常（panic）からの復帰、アクセスログ、CORS。
- **CORS**：ブラウザは、ページを配ったのとは別のオリジン（`http://localhost:45029` のような、スキーム・ホスト・ポートの組）への通信について、相手のサーバーが許可していなければ、応答をページに渡しません。api は許可するオリジンを列挙していて、ローカルでは `http://localhost:45029` などが入っています。`127.0.0.1` で開いたページは許可リストと一致しないので、画面に何も出ません。これはブラウザだけがする検査で、GAS や `curl` からの通信には関係しません。
- **Slack 通知**（`lib/externals/slack/`、`lib/externals/scheduler/`）：シフトの変更は、api が `action_logs` テーブルに記録しておき、5分ごとに動く処理がまとめて本人に DM します。レスキューの通知はその場で送ります。詳しくは AGENTS.md の「通知の定期実行」。
- **自動再ビルド**：ローカルのコンテナは [air](https://github.com/air-verse/air) で起動しているので、`.go` ファイルを保存すると自動で再ビルドしてサーバーを起動し直します。`di.go` を変えてうまく動かないときは、コンテナを一度 `make down` してから起動し直してください。

**手本にするファイルを選んでください。** 古いファイルには、今は禁止している書き方が残っています。全部の規約を満たした完全な手本は、まだありません。規約ごとに手本を選んでください。

| 規約 | 古い書き方の例 | 新しい書き方の例 |
| --- | --- | --- |
| SQL はプレースホルダで書く | `bureau_repository.go`（文字列連結） | `review_repository.go` |
| エラーは `{"error": "..."}` の JSON で返す | `bureau_controller.go`（`return err` で Echo に任せる） | `review_controller.go` |
| 空の一覧は `null` ではなく `[]` を返す | `bureau_usecase.go`、`review_usecase.go` の `GetReviewsGAS`（0件で `nil`） | `var res = []entity.Xxx{}` で初期化する |
| 見つからないときは 404 を返す | `review_controller.go` の `ShowReview`（500 を返す） | `sql.ErrNoRows` を見て 404 にする |

`nil` の一覧を JSON にすると `null` になり、配列を期待している mobile が壊れます。

**学ぶこと：**

- HTTP そのもの：メソッド（GET / POST / PUT / DELETE）、ステータスコード（200 / 201 / 400 / 404 / 500）、ヘッダー、JSON ボディ
- Echo：`e.GET(...)` によるルーティング、`c.Param`、`c.Bind`、`c.JSON`
- interface を受け取る設計と、それがテストのしやすさにつながる理由
- CORS と、なぜブラウザだけが検査するのか

**どこを見るか。** `review` の一式（`review_controller.go` → `review_usecase.go` → `review_repository.go` → `entity/review.go`、`entity/reviewGAS.go`）が、一覧・取得・作成・更新・削除の揃った例です。上の表のとおり規約に合っていない箇所もあるので、それを見つけながら読んでください。repository にはテスト（`review_repository_sqlmock_test.go`）もあります。

エンドポイントの一覧は `router.go` を読んでください。`/swagger/index.html` にも API の説明ページがありますが、2025年7月から更新されておらず、今のエンドポイントとは一致しません。

---

## 6. Dart と Flutter Web（mobile 担当）

**何か。** Flutter は Google の UI フレームワークで、言語は Dart です。1つのコードから Android / iOS / Web などを作れますが、**SeeFT は Web としてだけ動かしています。** 参加者はアプリをインストールせず、スマホのブラウザで URL を開きます。`mobile/` の下に `android/`・`ios/`・`macos/`・`windows/` も残っていますが、使っていません。

**Flutter のバージョンは [fvm](https://fvm.app/) で固定しています**（`mobile/.fvmrc` で 3.27.3）。`flutter` コマンドは必ず `fvm flutter ...` の形で打ってください。素の `flutter` を打つと、手元に入っている別のバージョンが動き、依存ライブラリの解決やビルド結果が変わります。Flutter 自体を上げる作業は、他の変更と混ぜず、専用のブランチで行います（#514）。

**画面は Widget の木です。** 画面のすべて（文字、ボタン、余白、並べ方）が Widget で、Widget の中に Widget を入れて画面を組み立てます。

- `StatelessWidget`：引数だけで見た目が決まるもの
- `StatefulWidget`：状態を持つもの。状態は、組になる `State` クラスのオブジェクトが持ちます。Widget 自体は描き直すたびに作り直されますが、`State` は画面にある間ずっと残ります

```dart
// 状態を変えるときは setState の中で。Flutter が build を呼び直して描き直す
setState(() => _isLoading = false);
```

**設定値はビルド時に埋め込まれます。** api の URL などは `mobile/env/.env` に書き、`--dart-define-from-file=env/.env` で渡しています。コードからは `String.fromEnvironment('API_BASE_URL')`（`lib/configs/constant.dart`）で読みますが、これはビルドした時点で値がコードに埋め込まれる仕組みです。値を変えたら、ビルドし直すまで反映されません。また、埋め込まれた値はブラウザに配られる JavaScript にそのまま入り、誰でも読めます。`.env` には公開してよい値だけを書き、トークンやパスワードは絶対に書かないでください。

**学ぶこと：**

- Dart：型、`final` と `const`、null 安全（`String?`、`!`、`??`）、`async` / `await` と `Future`
- Widget の木、`build` メソッド、`StatelessWidget` と `StatefulWidget` の違い
- レイアウト：`Column`, `Row`, `Padding`, `Expanded`, `ListView`
- `setState`、`initState`、`dispose` の役割
- `Navigator` による画面遷移
- Chrome の開発者ツールで、mobile が出す通信（Network タブ）とログ（Console タブ）を見ること

**やってみる:** 公式の [Write your first Flutter app](https://docs.flutter.dev/get-started/codelab) を、Web を対象にしてやってください。

---

## 7. mobile の構成と状態の持ち方（mobile 担当）

`mobile/lib/` の下はこう分かれています。

```text
pages/     画面。1ファイル1画面（my_shift_page.dart がマイシフト画面）。rescue/ はレスキューの画面群
widgets/   複数の画面で使う部品（shift_card.dart がシフトの1コマ）
models/    データの形
utils/     api.dart（API の呼び出し）、logger.dart、permanent_store.dart（端末への保存）
configs/   constant.dart（設定値）、importer.dart（よく使う import をまとめたもの）
theme/     tokens.dart（色・文字サイズ）、theme.dart（アプリ全体の見た目）
```

**API は1つのオブジェクト経由で呼びます。** `lib/utils/api.dart` にある `api` というオブジェクト（アプリ全体で1つだけ作るので、シングルトンと呼びます）に、エンドポイントごとの関数（`api.getAllManual()`、`api.postReview(...)` など）が並んでいます。画面側は `import 'package:seeft_mobile/configs/importer.dart';` だけで `api` を使えます。HTTP 通信の処理を `api.dart` だけにまとめるため、`package:http` を `api.dart` 以外から import するのは禁止です。

**状態管理ライブラリは使いません。** Riverpod や Provider のような、画面をまたいで状態を共有するためのライブラリは入れず、各画面の `State` の中に状態を持ち、`setState` で描き直します。AGENTS.md で導入を禁止しています。画面数が少なく、画面をまたいで共有する状態がほとんど無いので、それで足りています。

**非同期の処理を待った後は `mounted` を確かめます。** api の応答を待っている間に利用者がその画面を閉じると、その画面の `State` は破棄されます。破棄された `State` で `setState` を呼ぶと例外になります。

```dart
final manuals = await api.getAllManual();
if (!mounted) return;   // 待っている間に画面が閉じられていたら何もしない
setState(() => _allManuals = manuals);
```

`manual_list_page.dart` の `_loadData` も、同じ考え方で `mounted` を確かめてから `setState` しています。

**端末への保存は2種類あります**（`lib/utils/permanent_store.dart`）。

- **SharedPreferences**：ログイン中のユーザー ID のような小さな値
- **Hive**：シフトカードやレスキューの一覧のような、まとまったデータ。マイシフト画面は、前回取得したシフトカードを日付・天気ごとに Hive から読み出しています（`my_shift_page.dart`）

Web なので、どちらも実体はブラウザのストレージです。ブラウザがデータを消すこともあるので、消えても api から取り直せるものだけを置きます。検証中に表示が古いままのときは、ブラウザのサイトデータを消すと直ることがあります。

**見た目とログの規約。** 色と文字サイズは `AppColors.main` や `AppFontSizes.md`（`lib/theme/tokens.dart`）を使い、`Color(0xFF...)` を直接書きません。ログは `print` ではなく `logger.i(...)` / `logger.e(...)` で出します。

**学ぶこと：**

- `StatefulWidget` のライフサイクル：`initState` → `build` → `dispose`
- `FutureBuilder` と、`initState` で Future を作っておく理由（`widgets/first_jump_selector.dart` のコメントに書いてあります）
- `mounted` がなぜ必要か
- SharedPreferences と Hive の使い分け

**どこを見るか。** 起動の流れは `main.dart` → `widgets/first_jump_selector.dart`（ログイン済みかを見て、ログイン画面かマイシフト画面に振り分ける）です。画面を1つ読むなら `pages/manual_list_page.dart` が、API の呼び出し・読み込み中の表示・エラーの表示・検索が1画面に収まっていて手頃です。

---

## 8. Google Apps Script（GAS）と clasp（GAS を触る人）

**何か。** Google のスプレッドシートに紐づけて動かせる JavaScript です。SeeFT では、スプレッドシートに作ったメニューから名簿・タスク・シフトを api に送ったり、api から届いたレスキューをスプレッドシートに書き込んだりしています。

**`gas/` はクラウド上のコードの写しです。** GAS のコードは、各スプレッドシートに紐づいた Google のクラウド上のプロジェクトにあり、動いているのはそちらです。リポジトリの `gas/` にあるのは、[clasp](https://github.com/google/clasp)（GAS のコードをコマンドラインで取得・反映するツール）で取ってきた、ある時点のコピーです。誰かがスプレッドシートのエディタで直接編集すると、クラウド上のプロジェクトだけが更新され、`gas/` は古いままになります。**GAS を触る作業は、`gas/` を読むことではなく、`clasp` でクラウド上の最新のコードを取ってくることから始めます。** 手順と、clasp の入れ方（`npm --prefix gas install`）は [gas/README.md](../../gas/README.md) にあります。

**GAS には、壊れたコードを止める仕組みがありません。** api と mobile には CI がありますが、`gas/` を変えても何も走りません。JavaScript には静的な型の検査も無く、`clasp push` で壊れたコードを反映しても、それを検出する仕組みがありません。しかも GAS は全ファイルを1つのスコープで読み込むので、別ファイルで同じ名前の `const` を定義しただけでスクリプト全体が読み込めなくなり、スプレッドシートのメニューごと消えます。gas/README.md の手順が、全ファイルを1つにつなげてから `node --check` で構文を確かめるのはこのためです（1ファイルずつ確かめても、ファイルをまたぐ重複は見つかりません）。

GAS 特有の道具:

- `UrlFetchApp.fetch(url, options)`：api を呼ぶ
- `PropertiesService.getScriptProperties()`：api の URL やトークンの置き場。コードに直接書かない
- `LockService`：同じスクリプトが同時に動かないよう、順番待ちにする鍵。取ったら `finally` で必ず放す。順番待ちにするだけなので、メニューを2回押せば、2回目は1回目の後に普通に動きます。二重の送信を防ぐには、送る側・受ける側で重複を見分ける処理が別に要ります
- `doPost(e)`：GAS を Web アプリとして公開し、外から POST を受ける入口。api から届くレスキューをここで受ける
- `onChange(e)`：スプレッドシートが変更されたときに動かす関数。スプレッドシートを開いたときなどに自動で動く特別な名前の関数とは違い、トリガーとして登録しないと動きません

**学ぶこと：**

- JavaScript の基礎：`const` / `let`、関数、配列の `map` / `filter`、オブジェクト、`JSON.stringify`
- スプレッドシートの操作：`SpreadsheetApp`、`getRange`、`getValues` / `setValues`（1セルずつ読み書きすると遅いので、範囲でまとめて読む）
- GAS の実行時間の上限（6分）。シフトの送信は件数が多く、1回では送り切れないので、途中で区切り、次の実行で続きから送る作りになっています（`gas/shift/コード.js`）
- clasp：`clone`、`pull`、`push`、`deploy`
- Web アプリの「デプロイ」と「バージョン」。`clasp push` しただけでは、公開中の Web アプリは変わりません

**どこを見るか。** `gas/rescue/コード.js` の `doPost` が、api から届いたレスキューをスプレッドシートに書き込む処理です。`gas/rescue/onChange.js` が、スプレッドシートで変えた対応状況を api に送り返す処理です。

---

## 9. データの流れを1本追う（全員）

リポジトリを読むときは、1つの操作を端から端まで追ってください。例として、mobile でタスクのレビュー（人手の足り具合とマニュアルの評価）を送る操作を追います。

1. `mobile/lib/widgets/review_bottom_sheet.dart`：送信ボタンで `api.postReview(...)` を呼ぶ
2. `mobile/lib/utils/api.dart`：`POST /reviews` に JSON を送る
3. `api/lib/router/router.go`：`e.POST("/reviews", r.reviewController.CreateReview)`
4. `api/lib/internals/controller/review_controller.go`：`CreateReview` が JSON を struct に読み込み、usecase を呼ぶ
5. `api/lib/usecase/review_usecase.go`：`CreateReview` がタスク名からタスクの ID を引き、repository に保存させる
6. `api/lib/internals/repository/review_repository.go`：`INSERT INTO reviews ...` をプレースホルダで実行する
7. `postgresql/db/schema/create4Review.sql`：`reviews` テーブルの定義

手順2で Dart が組み立てる JSON のキー名（`user_id`、`task_name` など）と、手順4の controller の struct タグが一致していることを確かめてください。0節で書いた「約束は JSON の形だけ」の具体例がこれです。どちらかを変えても、コンパイルは通ります。

この追跡を一度丁寧にやると、チュートリアルを読むより多くのことが分かります。慣れたら、次の2本も追ってみてください。

- **シフトの送信**：`gas/shift/コード.js` → `POST /api/update_shifts` → `shift_controller.go` の `UpdateShiftsFromGAS` → `shift_usecase.go`。関わるファイルと処理が多く、SeeFT で最も重要な処理です
- **レスキュー**：`mobile/lib/pages/rescue/` → `POST /rescues` → `rescue_unified_controller.go` → `rescue_unified_usecase.go` の `SendRescueToGAS` → `gas/rescue/コード.js` の `doPost`。api から GAS へ向かう、逆向きの通信が入っています

---

## 10. テスト（api 担当・mobile 担当）

**api（Go）。** 標準の `go test` で動くテストに、[testify](https://github.com/stretchr/testify)（`assert` / `require` で結果を確かめるライブラリ）と [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)（偽の DB を作るライブラリ）を足して書いています。テストファイルは、対象と同じディレクトリに `xxx_test.go` として置きます。

repository も usecase も、実際の DB ではなく go-sqlmock の偽の DB を相手にテストします。「この SQL がこの引数で実行されるはず」と先に宣言しておき、実行された中身を確かめる形です。5節の DI のおかげで、`db.Client` を偽の DB に差し替えるだけで、本物の repository をそのまま使ってテストできます。usecase のテストでも、偽の repository は作らず、偽の DB を載せた本物の repository を渡しています。repository が行（`*sql.Rows`）を返す作りなので、偽の repository を作るのが難しいためです。

```go
// api/lib/internals/repository/review_repository_sqlmock_test.go（抜粋）
mock.ExpectExec(`INSERT INTO reviews \(user_id, task_id, staffing_rating, manual_rating, comment\)`).
    WithArgs("12", "34", "5", "4", testReviewComment).
    WillReturnResult(sqlmock.NewResult(1, 1))
```

go-sqlmock が確かめられるのは「送った SQL の文字列と引数」だけです。SQL が実際のテーブルで正しく動くか、`Scan` の列の順番が合っているかは確かめられません。

repository のテストは、`sqlmock_helper_test.go` の `newDBMock` を使うと短く書けます（repository のパッケージの中でだけ使えます）。どの層からどう書き進めるかの方針は [test-roadmap.md](test-roadmap.md) にあり、テストが書かれていない関数の issue（題名が `[test]` で始まるもの）が練習に向いています。

`review_repository_sqlmock_test.go` の冒頭には、アポストロフィ（`'`）を含むコメントのテストがあります。SQL を文字列連結で組み立てると、この `'` が SQL の文字列の終わりと解釈されて SQL が壊れます。プレースホルダがなぜ要るのかの具体例として読んでください。

**mobile（Flutter）。** `flutter_test` で widget テストを書いています。画面の部品を仮想的に描画して、表示や、押したときの動きを確かめます。

```bash
cd mobile && fvm flutter test --platform chrome
```

`--platform chrome` が必須です（`-p` と略すことはできません）。付けないと、Web 向けのコードを読み込めずにテストが落ちます。`test/widget_test.dart` は Flutter の雛形のカウンタを確かめる中身が残っていて、必ず落ちます。気にしないでください（#493 で除去予定）。

**学ぶこと：**

- Go：`func TestXxx(t *testing.T)`、テーブル駆動テスト（入力と期待値の組を表にして並べる書き方）、`t.Run`、`assert` と `require` の違い
- go-sqlmock：`ExpectQuery` / `ExpectExec`、`WithArgs`、`WillReturnRows`、`ExpectationsWereMet`
- Flutter：`testWidgets`、`pumpWidget`、`find.text`、`tester.tap`
- 何をテストする価値があるか。境界の値や、過去に事故が起きた入力（アポストロフィを含む文字列、空の値）

**GAS には自動テストがありません。** 検証環境のスプレッドシートで実際に動かして確かめます（[staging-rehearsal.md](../operations/staging-rehearsal.md)）。

---

## 11. ツール、Docker、CI（全員）

- **Docker / Docker Compose**：ローカルでは DB と api がコンテナで動きます（`docker-compose.yml`）。mobile はコンテナではなく手元の fvm で動かします。本番も Docker Compose で、api・mobile・admin・外部公開用のトンネル（Cloudflare Tunnel）がコンテナで動いています。本番の DB は compose の外にあります（`docker-compose.prod.yml`、手順は [deploy.md](../operations/deploy.md)）。
- **Makefile**：よく使う docker compose のコマンドに短い名前を付けたものです。一覧はリポジトリ直下の `Makefile` を開くのが確実です。Mac 用の compose ファイル（`docker-compose.mac.yml`）を使う `mac-` 付きのコマンドもあります。
- **GitHub Actions（CI）**：PR を出すと、変更した場所に応じて次が走ります。
  - `api/` を変えたとき：golangci-lint（`go-lint.yml`。PR で新しく増えた指摘だけを見る）と `go test`（`go-test.yml`）
  - `mobile/` を変えたとき：`flutter analyze --fatal-infos`（`flutter-lint.yml`。最も軽い info レベルの指摘でも落ちる）
  - `gas/` を変えたとき：何も走らない
- **CodeRabbit**：PR に AI がレビューコメントを付けます。指摘は参考です。すべてに従う必要はなく、スコープ外のものは理由を書いて見送って構いません。
- **Git**：issue を立て、`feat/{ユーザー名}/{issue番号}/{内容}` のブランチで作業し、PR を出します。コミットメッセージは日本語で、`feat:` / `fix:` / `docs:` を付けます。詳しくは [開発の進め方](workflow.md)。

**学ぶこと：**

- Docker の基礎：image（コンテナの元になる型）と container の違い、`docker compose up / down / logs / exec`、volume、port
- lint の指摘を読み、無効化のコメントで抑えるのではなく原因を直すこと
- workflow ファイルを読める程度の YAML
- Git：ブランチ、merge、コンフリクトの解消、1つの PR に1つの目的

---

## 12. 現実的な学習順序

全部を一度に学ぼうとしないでください。通読した後は、担当する場所から深く学べば十分です。api と mobile の両方を触る人でも、だいたい2〜3週間です。

| 段階 | 重点 |
| --- | --- |
| 1 | 2節で手元で動かす。Git と GitHub の流れ（issue → ブランチ → PR） |
| 2 | 担当が api なら Go の基礎（A Tour of Go）、mobile なら Dart と Flutter の基礎 |
| 3 | api：SQL の基礎と `schema/` の読解。`psql` でテーブルを眺める |
| 4 | api：層の構成と DI。`review` の一式を読む。mobile：画面1つ（`manual_list_page.dart`）を読む |
| 5 | 9節のデータの流れを1本追う |
| 6 | テスト（go-sqlmock または widget テスト）を1本書く。Docker と CI の基礎 |
| 7 | GAS と clasp（GAS に関わる issue を持ったときでよい） |

---

## 13. 何か作る

読むだけでは足りません。同じ技術を使った**自分の小さなシステムを、新しいリポジトリで**作ってください。SeeFT とは関係ないもの（本棚、習慣の記録、サークルの備品の貸し出し表など）で構いません。題材は何でも構いません。大事なのは、画面・api・DB をつなぐ部分を自分で書くことです。

担当が片方だけなら、その側だけを作り、もう片側は SeeFT のローカル環境で代用して構いません。締め切りや報告の仕方は、PM と相談して決めてください。

### 最低限

**api：Go + Echo**

- Docker Compose 上の PostgreSQL
- 関連する2つ以上のテーブル（1対多）を `CREATE TABLE` の SQL で定義し、リポジトリにコミットする
- controller / usecase / repository の3層と entity に分け、SeeFT の `di.go` と同じ形で手で組み立てる
- 1つのテーブルについて、一覧・1件取得・作成・更新・削除のエンドポイント
- SQL はすべてプレースホルダで書く
- 見つからないときは 404、入力が不正なときは 400 を、`{"error": "..."}` の JSON で返す
- 空の一覧は `null` ではなく `[]` を返す
- go-sqlmock を使った repository のテストを1本以上

**画面：Flutter Web**

- 一覧画面と、一覧から開く詳細画面
- API の呼び出しを1つのファイル（SeeFT の `api.dart` のような）にまとめる
- 読み込み中の表示と、失敗したときの表示
- 作成フォームを1つ。送信後に一覧が新しい内容で表示される
- `await` の後、`setState` の前に `mounted` を確かめる
- api の URL は `--dart-define-from-file` で渡し、コードに直接書かない
- widget テストを1本以上（`fvm flutter test --platform chrome` で通ること）

### 発展（時間があれば）

- スプレッドシートを1つ作り、GAS からあなたの api にデータを送るメニューを付ける。URL は `PropertiesService` に置き、同時に2回動いたときに何が起きるかを確かめる
- 既存のテーブルに列を足す migration を golang-migrate で書き、up と down の両方を試す
- 一覧のデータを Hive に保存し、api が止まっているときは前回の内容を表示する
- 一定間隔で動く処理（SeeFT の `scheduler` のような、決まった間隔で関数を呼ぶループ）を api に足す
- PR ごとに lint とテストを走らせる GitHub Actions の workflow

### 「完了」の意味

次を、調べずに説明できることです。

全員:

- api と mobile の間で JSON のキー名が食い違ったとき、いつ・どうやって気づくか
- SQL を文字列連結で組み立てると、どんな入力で何が起きるか

api 担当:

- SeeFT のテストで、本物の DB の代わりに何を、どこに差し込んでいるか
- テーブルに列を1つ足したとき、Go のコードで直す必要があるのはどこか。直し忘れるといつ気づくか

mobile 担当:

- `await` の後で `mounted` を確かめないと何が起きるか
- `mobile/env/.env` の値を変えたのに画面が変わらないのはなぜか。そこにトークンを書いてはいけないのはなぜか

GAS を触る人:

- `gas/` のコードを読んで「今はこう動いている」と判断してはいけないのはなぜか

これらに答えられるようになったら、本番の issue に取り組む準備ができています。

---

## 14. ブックマーク

- A Tour of Go — https://go.dev/tour/
- Effective Go — https://go.dev/doc/effective_go
- Echo — https://echo.labstack.com/docs
- `database/sql` のチュートリアル — https://go.dev/doc/tutorial/database-access
- go-sqlmock — https://github.com/DATA-DOG/go-sqlmock
- golang-migrate — https://github.com/golang-migrate/migrate
- PostgreSQL チュートリアル — https://www.postgresql.org/docs/current/tutorial.html
- MDN の CORS の解説 — https://developer.mozilla.org/ja/docs/Web/HTTP/CORS
- Dart の言語ツアー — https://dart.dev/language
- Flutter — https://docs.flutter.dev/
- fvm — https://fvm.app/
- Hive — https://pub.dev/packages/hive_flutter
- Apps Script — https://developers.google.com/apps-script
- clasp — https://github.com/google/clasp
- Docker Compose — https://docs.docker.com/compose/

リポジトリ内：

- [AGENTS.md](../../AGENTS.md)：コードの書き方の規約
- [docs/operations/](../operations/README.md)：技大祭での運用と、過去の事故の記録
- [workflow.md](workflow.md)：issue から本番反映までの仕事の進め方
- [docs/decisions/](../decisions/README.md)：何をなぜ決めたか、何を選ばなかったかの記録（ADR）
- [test-roadmap.md](test-roadmap.md)：テストの進め方
- [gas/README.md](../../gas/README.md)：GAS の取得と反映

---

## 15. 助けの求め方

しばらく自分で試してから、聞いてください。丸一日詰まったままにするのも、2分で聞くのも避けてください。聞くときは次を含めてください。

- 何をしようとしていたか
- 何が起きると思っていたか
- 実際には何が起きたか
- 正確なエラー文（スクリーンショットより、コピーした文字のほうが検索できて助かります）
- すでに試したこと

この形なら早く答えが返ってきます。issue も同じ項目で書きます。担当している issue で詰まったら、Slack で PM に連絡してください。
