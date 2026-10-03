# オンボーディング: 技術スタック

SeeFTに新しく入った人向けの資料です。GoもFlutterもApps Scriptも初めて、という人を想定しています。授業でPythonを少し書いたことがあれば読めるように書きました。

この文書が扱うのは、機能ではなく技術です。技術ごとに、何を使っているか、SeeFTでどう使っているか、何を学ぶか、どこを読めばよいかを説明します。学んだことを確かめるために自分で作るものは、13節にまとめました。SeeFTが技大祭で何をしているかには、技術の説明に要る範囲でだけ触れます。機能については、issueに取り組み、レビューを受けるうちに分かってきます。運用の手順は[docs/operations/](../operations/README.md)にまとまっています。

## この資料でいちばん伝えたいこと

SeeFTには、書き間違えてもエラーが出ず、動かしても間違った結果のまま進んでしまう場所がいくつもあります。コンパイラ（コードを実行できる形に変換し、そのときに誤りを検査するプログラム）もテストもCIも、そこでは誤りを見つけてくれません。たとえば次の場所です。

- apiとmobile・GASの間でやり取りするJSONのキー名（0節）
- migrationのファイルの番号と、適用済みのファイルの書き換え（4節）
- SQLの`SELECT *`と、結果を受け取る変数の並び（5節）
- ビルドのときに埋め込まれるmobileの設定値（6節）
- 端末に保存したまま古くなったデータ（7節）
- リポジトリの`gas/`と、実際に動いているGASのずれ（8節）
- 偽のDBを相手にしたテストでは分からない誤り（10節）

これをいちばん伝えたい理由は3つあります。

まず、45thの事故の多くが、この形で起きたからです。事故の記録（[incidents-45th.md](../operations/incidents-45th.md)の「共通する原因」）は、共通点を2つ挙げています。データがスプレッドシートからapiに渡る境目で起きていることと、失敗が画面に出ないことです。

次に、エラーが出る誤りは、エラーの文を読めば直し始められます。エラーが出ない誤りは、どこで起きうるかを先に知っていなければ、探すこともできません。

最後に、言語やフレームワークの使い方は、公式のチュートリアルで学べます。SeeFTのどこが危ないかは、SeeFTの中でしか学べません。

各節では、こうした場所を「間違っていてもエラーが出ないところ」という小見出しで示します。13節の「完了」の問いも、多くはここから出します。技術を1つずつ学ぶのと並行して、「この変更で、どこに間違いがあってもエラーが出ないか」を考えながら読んでください。

## この資料の読み方

初回は上から通して読んでください。0節で全体の形を見て、1節で言葉をそろえ、2節で手元で動かしてから、技術ごとの節（3〜8節）に進みます。前の節に出てくる言葉のうち、詳しい説明が後ろの節にあるものには、節の番号を添えました。分からない言葉は、まず1節の早見表を見てください。

通読したら、自分の担当に合わせて深く学びます。担当はapiかmobileで、PM（プロジェクトの取りまとめ役）がissueを割り振るときに決まります。GAS（8節）は、GASに関わるissueを持ったときに読めば足ります。各節の見出しの後ろに、その節がとくに必要な人を書きました。「全員」の節は担当にかかわらず読んでください。

コードの書き方の規約は[AGENTS.md](../../AGENTS.md)にあります。AGENTS.mdはもともとAIのコーディング支援ツール向けに書いたものですが、人が読んでもそのままSeeFTの規約です。この資料を読み終えてから読むと、なぜその規約があるのかが分かるはずです。

---

## 0. システムの全体像（全員）

1つのリポジトリに、動くものが3種類入っています。

```text
スプレッドシート ◄──► GAS                               （Google のクラウド）
                      │ ▲
              (1)(3b) │ │ (3a)
                      ▼ │
                      api (Go + Echo) ──► PostgreSQL    （サーバー）
                       ▲              └──► Slack（本人への DM）
               (2)(3a) │
                    mobile (Flutter Web)                （参加者のブラウザ）
```

矢印は、通信を始める側から受ける側に向いています。データは、応答として逆向きにも流れます（(2)ではmobileがapiに問い合わせ、apiがシフトを返します）。スプレッドシートとGASの間だけは両向きにしました。GASがシートを読み書きし、シートのメニューや変更がGASを動かすためです。番号は、下のデータの流れの番号です。

- (1)シフト・名簿・タスク。シフトを組むのは人間で、その作業場所はスプレッドシートです。GAS（スプレッドシートに付いたスクリプト。8節）が、その内容をapiに送ります。SeeFTは、組まれたシフトを受け取り、参加者に表示・通知する側です。
- (2)シフトの閲覧。参加者はスマホのブラウザでmobileを開きます。mobileはapiから自分のシフトを取ってきて表示します。
- (3)レスキュー（当日の人手不足やトラブルの連絡）。(3a)レスキューはmobileからapiに送られます。apiはそれをGASに送り、GASがスプレッドシートに書き込みます。(3b)本部の人がスプレッドシート上で対応状況を変えると、GASがそれをapiに送り返します。

シフトが変わったときや、レスキューの対応状況が変わったときは、apiがSlackで本人にDMを送ります。

3つは動く場所も違います。

| | 書く言語 | 動く場所 | リポジトリにあるもの |
| --- | --- | --- | --- |
| api | Go | サーバー上のDockerコンテナ | 本体のコード |
| mobile | Dart | 参加者のブラウザ（サーバーはビルド済みのファイルを配るだけ） | 本体のコード |
| GAS | JavaScript | Googleのクラウド | ある時点の写し（8節） |

#### 間違っていてもエラーが出ないところ

言語が場所ごとに違うので、apiとmobile / GASの間で型を共有する仕組みはありません。たとえばGoのstructに項目を足しても、Dartのクラスは自動では変わりません。両者の間の約束は、どのURLにどのメソッド（GETやPOST。5節）で、どんな形のJSONを送り、何が返ってくるかだけです。JSONの形とは、キー名、値の型（文字列か数値か）、値が`null`になりうるか、配列かどうかです。エラーのときに`{"error": "..."}`を返すことや、空の一覧を`[]`で返すことも、この約束に含まれます（5節）。

この約束は、コンパイラには検査できません。食い違ったときの起き方は、食い違いの種類で2通りあります。

値の型が食い違った場合（数値のところに文字列を送った、など）、多くのエンドポイントは、受け取る形をstructで決めています。その場合、apiは受け取るときにエラーを返します。動かせば気づけます。

キー名が食い違った場合、api（Go）は、知らないキーをエラーにせずに捨てます。受け取る側の値は空のまま（数値なら`0`、文字列なら`""`）で処理が進み、そのままDBに保存されることがあります。動かしても気づけないことがあります。逆向きに、apiの応答のキー名を変えたときは、mobileの側でその値が`null`になります。nullを許す型ならエラーにならずに空になり、許さない型なら実行時に例外になります。

最初に理解してほしいのはこの点です。キー名の食い違いではありませんが、apiが空の値をそのまま受け入れる、という同じ性質で起きた事故があります。45thでは、名簿を送るGASが学籍番号に常に`0`を送っていて、apiはそれをそのまま受け入れました。本番で全員がログインできなくなり、そこで初めて気づきました（[45thの事故の記録](../operations/incidents-45th.md)）。AGENTS.mdが「既存entityのJSONキー命名変更」や「APIレスポンス形式の変更」を事前相談（AGENTS.mdの「Ask First」。実装の前にPMに相談する）の対象にしているのは、このためです。

環境は3つあります。この資料で「ローカル」と書いたら自分のPC、「本番」と書いたら技大祭で参加者が使うサーバーです。その間に、本番に入れる前に通しで試す検証環境があります（[staging-rehearsal.md](../operations/staging-rehearsal.md)）。ローカルのデータは何度消しても作り直せますが、本番のDBやスプレッドシートのデータは、書き換えたり消したりすると元に戻せません。GASにはローカルがありません。GASを試すときは、検証環境のスプレッドシートを使います。

リポジトリの構成は次のとおりです。

```text
api/                 Go で書いた api のサーバー（ローカルでは port 1234）
  main.go            起動するだけ。部品の組み立ては lib/di（5節）
  lib/               本体（5節）
  cmd/migrate        DB のテーブルを作る・変えるコマンド（4節）
  cmd/seed           初期データを入れるコマンド
  cmd/send-notifications  未送信の Slack 通知を手動で送るコマンド
mobile/              Flutter Web のアプリ（ローカルでは port 45029）
  lib/               本体（7節）
  test/              画面の部品（Widget。6節）のテスト
  python/server.py   本番でビルド済みのファイルを配る小さなサーバー
gas/                 Apps Script の写し（8節）
postgresql/db/       テーブル定義（schema/）、既存テーブルへの変更（migrations/）、初期データ（seed.sql）
docs/                この資料を含む文書
admin/               管理画面。今は開発していない（読まなくてよい）。本番の compose には残っている
```

---

## 1. 用語の早見表（全員）

この資料で説明なしに使う言葉です。最初は全部覚えなくてよいので、出てきたときに戻ってきてください。

| 言葉 | 意味 |
| --- | --- |
| PM | プロジェクトの取りまとめ役。issueの割り振りやマージをする |
| issue | GitHub上の作業票。やることと理由が書いてある。`#461`のような番号で呼ぶ（PRも同じ番号の付け方） |
| ブランチ、develop、マージ | ブランチは、作業ごとに分けたコードの流れ。developはSeeFTの本流のブランチで、各自のブランチの変更をdevelopに取り込むことをマージと呼ぶ |
| PR（プルリクエスト） | 自分のブランチの変更をdevelopに取り込んでもらうための申請。ここでレビュー（他の人による確認）を受ける |
| CI | PRを出すたびにGitHubが自動で走らせる検査（11節） |
| lint | コードの書き方の問題を機械的に見つけるツール |
| コンパイル | ソースコードを実行できる形に変換すること。GoとDartでは、このときに型の誤りなどが見つかる。Pythonにはこの段階がない |
| Docker、コンテナ | アプリを、必要なものごと箱に詰めて動かす仕組み。その箱がコンテナ |
| make | よく使う長いコマンドに短い名前を付けて呼ぶ道具。定義はリポジトリ直下の`Makefile` |
| ポート | 1台のPCの中で、どのプログラムに通信を届けるかを区別する番号 |
| HTTP、API | HTTPはブラウザとサーバーがやり取りする決まり。APIはプログラムから呼べる窓口で、SeeFTでは`GET /reviews`のようなURLとメソッドの組（エンドポイント）の集まり |
| JSON | `{"user_id": 1, "comment": "..."}`のような、データをやり取りするための書き方 |
| ライブラリ、フレームワーク | 他人が書いた部品。フレームワークはプログラムの骨組みまで決めるもの |
| ビルド | ソースコードを、動かせる形（実行ファイルや配信用ファイル）に変換すること |
| デプロイ | ビルドしたものを本番のサーバーに置いて動かすこと。GASでは、Webアプリとして公開する版を決める操作もデプロイと呼ぶ（8節） |
| migration | DBのテーブルの作り方や変え方をファイルに書いておき、順に適用するやり方。SeeFTでは、作るSQL（`schema/`）と変えるSQL（`migrations/`）を分けている（4節） |
| seed | 開発用にDBに入れる初期データ |
| 偽物、モック | テストで本物の代わりに差し込む作り物を、この資料では偽物と呼ぶ。モックは偽物の一種で、どう呼ばれるはずかを先に宣言しておき、そのとおりに呼ばれたかを確かめるもの（10節） |
| タスク | 技大祭で委員が受け持つ仕事の単位。シフトは、人・日付・時間帯・天気ごとにタスクを割り当てたもの |
| 天気（weather） | 技大祭の1日目・2日目は、晴れ用と雨用の2通りのシフトを組む。そのためDBとapiでは、シフトを日付と天気で引く。mobileは天気を選ぶ機能を外していて、晴れのシフトだけを表示する |

表記の約束として、小文字の「api」は`api/`ディレクトリとそのサーバーを指します。大文字の「API」は、サーバーが外部に公開している窓口（エンドポイントの集まり）を指します。「mobile」は`mobile/`とそのアプリのことです。

---

## 2. 手元で動かす（全員）

読み進める前に、一度SeeFTを手元で動かしてください。動くものを触りながら読むほうが、ずっと頭に入ります。この節はMacを前提に書いています。Windowsの人はPMに相談してください。

用意するものは次のとおりです。

- [Docker Desktop](https://www.docker.com/products/docker-desktop/)（入れたら起動しておく）
- [fvm](https://fvm.app/)（Flutterのバージョン管理ツール。6節）
- Gitとmake（Macでは、ターミナルで`xcode-select --install`を実行すると両方入ります）
- Google Chrome（mobileの画面と、開発者ツールでの確認に使う）
- Node.js（GASを触るときだけ）

リポジトリを取ってきます。

```bash
git clone https://github.com/NUTFes/SeeFT.git
cd SeeFT
```

以下はリポジトリのルート（`SeeFT/`のいちばん上）で実行します。まずDBとapiを用意します。

```bash
make build     # コンテナの元になるイメージを作る
make migrate   # DB を起動して、テーブルを作る（4節）
make seed      # 開発用の初期データを入れる
make up-api    # api を起動する
```

`make build`は使っていないadminもビルドするので、初回は時間がかかります。`make up-api`はログを出し続けるので、そのターミナルは開いたままにします。ブラウザでhttp://localhost:1234/bureausを開いて、JSONが返ってくればapiは動いています。

次に、別のターミナルでmobileを起動します。初回だけ、設定ファイルの雛形をコピーし、Flutterを入れます。

```bash
cp mobile/env/.env.example mobile/env/.env
cd mobile && fvm install && cd ..
make mobile-up
```

ブラウザでhttp://localhost:45029を開くとログイン画面が出ます。`127.0.0.1:45029`で開いたり、`make mobile-up`のポートを許可リストに無いものに変えたりすると、何も表示されません。apiが通信を許可しているページの場所（オリジン）と一致しないためです（5節のCORS）。ログインに使うテスト用アカウントはPMに聞いてください。seedにはシフトが入っていないので、ログインできてもマイシフトは空です。

ローカルでは、次の機能は動かないのが正常です。

- Slackの通知。`api/env/dev.env`の`SLACK_BOT_TOKEN`はダミーの値です。通知の仕組みは起動しますが、送信は失敗します。seedのユーザーにはSlackのIDが無いので、ふだんは送ろうともしません
- マニュアルの閲覧（`/manuals`）。Googleログインの設定が無いので、apiはこのURLを登録しません
- レスキューのスプレッドシートへの送信。`api/env/dev.env`の送り先`RESCUE_GAS_URL`がダミーの値なので失敗します。レスキュー自体はDBに保存されますが、mobileの画面には送信の失敗が表示されます

apiのコンテナは起動のたびに`go mod tidy`（依存ライブラリの整理）を実行するので、起動にはネット接続が要ります。

片付けるときは`make down`です。`make down`はコンテナを消し、DBの中身も消えます。次に動かすときは`make migrate`と`make seed`からやり直します。DBを残したまま止めたいときは`docker compose stop`を使います。この場合、次は`make up-api`だけで動きます（migrateとseedは要りません）。`make down`を繰り返すと、使われなくなったDBのデータがディスクに残っていきます。ディスクを空けたいときは`docker compose down -v`で、それも消します。

---

## 3. Go（api担当）

#### 何か

Googleが作った、静的型付け（変数の型をコンパイル時に検査する）の言語です。文法が小さく、書き方の選択肢が少ないので、他人のコードが読みやすいのが特徴です。バージョンはGo 1.26で、`api/go.mod`と`api/Dockerfile`で指定しています。

#### SeeFTでの使い方

apiはすべてGoで書かれています。ローカルではDockerコンテナの中で動かすので、手元にGoを入れなくてもapiは起動できます。ただしエディタの補完（VS CodeのGo拡張が使うgopls）のために、手元にも同じバージョンのGoを入れておくことを勧めます。[公式のダウンロードページ](https://go.dev/dl/)から入れられます。

Pythonを書いたことがある人がつまずくのは、次の3つです。

- 失敗は例外ではなく戻り値で知らせます。失敗しうる関数は最後の戻り値に`error`を返し、呼んだ側が毎回`if err != nil { return ..., err }`と書きます（`nil`はPythonの`None`にあたる「値が無い」）。冗長に見えますが、どこで失敗しうるかがコードに全部書いてあるということです。`panic`という仕組みもありますが、これは続行できない異常のためのもので、普段のエラー処理には使いません。
- クラスがありません。`struct`（データの入れ物）にメソッドを定義し、`interface`（メソッドの一覧）で抽象化します。ある型がinterfaceを満たすかどうかは、メソッドが揃っているかだけで決まります。「このinterfaceを実装する」という宣言は書きません。
- 大文字で始まる名前だけが、パッケージの外から使えます。SeeFTでは、interfaceを大文字、それを実装するstructを小文字で名付けています。たとえば`BureauController`（interface）と、それを作る関数`NewBureauController`は外から使えますが、実装のstruct `bureauController`は使えません。外からはinterface越しにしか触れないようにするためです。

#### 学ぶこと

- 変数、関数、複数の戻り値、`struct`、メソッド、ポインタ（`*T`と`&x`）
- `error`の扱い方。エラーに「どこで失敗したか」を足すには、標準では`fmt.Errorf("...: %w", err)`を使います。SeeFTの`api/lib/`では、外部ライブラリの`errors.Wrap`（`github.com/pkg/errors`）のほうを多く使っています。新しく書くときは、周りのコードに合わせてください
- `interface`と、それを満たすとはどういうことか
- スライス（`[]T`）とマップ、`for ... range`
- パッケージとimport、大文字・小文字による公開範囲
- `context.Context`がなぜどの関数にも第1引数で渡されているか（リクエストが打ち切られたことを下の処理に伝えるため）
- goroutine（軽量な並行処理）。自分で書く機会は少ないですが、GoのHTTPサーバー（標準ライブラリの`net/http`。Echoはその上で動く）は、接続ごとにgoroutineを起動して、複数のリクエストを並行に処理します。複数のリクエストから同時に書き換えられる変数を作らない、という前提だけは覚えておいてください（同時に書き換えて、片方の書き込みが消えるなど、値が予想と違ってしまうことを、データ競合と呼びます）

#### やってみる

[A Tour of Go](https://go.dev/tour/)を最後まで。2〜3日あれば終わります。

---

## 4. PostgreSQLとSQL（api担当）

PostgreSQLがデータベースです。ローカルではDockerの`postgres:18`で動きます。apiからはGo標準の`database/sql`でSQLを直接書いてアクセスします。ORM（SQLを書かずにオブジェクトの操作でDBを扱うライブラリ）のGORMも入っていますが、使っているのは`shift_card_repository.go`の1ファイルだけです。SeeFTでDBを触るとは、ほぼSQLを書くことです。

テーブルを作るファイルと、変えるファイルがあります。

1. `postgresql/db/schema/`は、`create1Bureaus.sql`のような、テーブルを作る`CREATE TABLE`のファイルの置き場です。ファイル名の数字が実行順で、同じ数字のファイルは名前の順に実行します。外部キーで参照される側ほど番号が小さくなっています（`create1Bureaus` → `create2Users` → `create4Shifts`）。
2. `postgresql/db/migrations/`は、既存のテーブルへの変更（列の追加など）の置き場です。`2026082201_add_manual_url_to_tasks.up.sql`（適用するときのSQL）と`.down.sql`（取り消すときのSQL）の組で書き、[golang-migrate](https://github.com/golang-migrate/migrate)が適用します。

`make migrate`はschemaを番号順に適用し、続けてmigrationsを適用します。何度実行しても、同じ変更が二重に適用されることはありません。ただし、適用済みかどうかを記録するやり方が、schemaとmigrationsで違います。schemaは、ファイルごとに、適用済みかどうかと中身のハッシュ（中身から計算した値。中身が変わると値も変わる）をDBに記録します。一方、golang-migrateは、適用済みの一番新しい番号だけを記録し、それより新しい番号のファイルだけを適用します。

#### 間違っていてもエラーが出ないところ

この違いのせいで、失敗したときの起き方も違います。

- 適用済みのschemaのファイルを書き換えると、`make migrate`がエラーで止まります。すぐ気づけます。
- 適用済みのmigrationのファイルを書き換えても、その番号はもう適用済みなので、変更はエラーも出ずに無視されます。
- 自分より新しい番号のmigrationが先に適用された環境では、自分のmigrationは「番号が古い」という理由で、エラーも出ずに飛ばされます。マージする前に、自分のファイルの番号（先頭の日時）がdevelopにあるmigrationのどれよりも新しいことを確かめてください。

**適用済みのファイルは編集しないでください**。変更は、migrationsに新しいファイルを足して書きます。過去にschema側へ列の追加を書いたファイル（`create6AddSlackUserIdToUsers.sql`）が1つ残っていますが、これはmigrationの仕組みを入れる前の名残で、整理中です（#461）。DBスキーマの変更は、AGENTS.mdで事前相談の対象です。

#### 学ぶこと

- SQLの基礎（`SELECT`、`JOIN`、`WHERE`、`GROUP BY`、`ORDER BY`、`INSERT` / `UPDATE` / `DELETE`）
- 主キーと外部キー、1対多、なぜ中間テーブルが要るのか
- `SELECT *`と、列を名前で並べる`SELECT id, name, ...`の違い（次の節で必要になります）
- プレースホルダ（`WHERE id = $1`と書き、値は別に渡す）と、値を文字列連結でSQLに埋め込むとなぜSQLインジェクションになるのか（10節）
- インデックスとは何か。`shifts`は主キー以外の検索に使う列（`user_id`など）にインデックスが無く、件数が増えると遅くなる課題があります（#500）
- `psql`でテーブルを眺めること。`docker compose exec db psql -U seeft -d seeft_db`でDBのコンテナに入れます（抜けるときは`\q`）

#### どこを見るか

まず`postgresql/er.png`（テーブルの関係を描いたER図）で全体を見て、`schema/`を番号順に読んでください。`users`、`tasks`、`shifts`の3つが中心で、残りの多くは局・学年・天気・時間帯のような選択肢の表です。

---

## 5. apiの構成：Echoと層（api担当）

#### Echoとは

GoのWebフレームワークです。URLとメソッドの対応づけ（ルーティング）、JSONの読み書き、ミドルウェア（すべてのリクエストの前後に挟む共通処理）を提供します。Echo自体の機能は少なく、ディレクトリの分け方や層の構成はSeeFTの規約で決めています。

#### 層の構成

リクエストは入口のrouterから、controller → usecase → repositoryの3つの層を順に下ってDBに届き、結果は逆順に上がってきます。routerは層ではなく入口で、entityも層ではなく、どの層からも使うデータの型です。

| 名前 | 場所 | 役割 |
| --- | --- | --- |
| router（入口） | `lib/router/router.go` | URLとメソッドをcontrollerの関数に対応づける。全エンドポイントの一覧はここ |
| controller | `lib/internals/controller/` | HTTPの入出力だけ。リクエストを読み、usecaseを呼び、JSONで返す。DBには触らない |
| usecase | `lib/usecase/` | 業務の処理。repositoryから受け取った行をentityに詰める |
| repository | `lib/internals/repository/` | SQLを実行するだけ。結果は行（`*sql.Rows`、1件なら`*sql.Row`）のまま返す |
| entity（共通の型） | `lib/entity/` | データの形。応答のJSONのキー名は、主にここのstructタグ（`json:"user_id"`のように、項目名の後ろに書く注記）で決まる |

リクエストのJSONのキー名は、entityではなくcontrollerの中で定義したstructで決めていることもあります（`review_controller.go`の`CreateReview`など）。

部品の組み立ては`di.go`に手で書いてあります。部品が使う別の部品を、部品の中で作らずに外から渡すことを、依存性注入（DI）と呼びます。NestJSやSpringのようにDIを自動でやってくれるフレームワークもあります。SeeFTはそれを使わず、`lib/di/di.go`でrepository → usecase → controller → routerの順に`NewXxx(...)`を呼び、作ったものを次の`NewXxx(...)`に渡しています。

```go
// api/lib/di/di.go（抜粋）
// client は DB への接続、crud は SQL を実行する共通の部品。di.go の少し上で作っている
bureauRepository := repository.NewBureauRepository(client, crud)  // repository に DB を渡す
bureauUseCase := usecase.NewBureauUseCase(bureauRepository)        // usecase に repository を渡す
bureauController := controller.NewBureauController(bureauUseCase)  // controller に usecase を渡す
```

外から渡すので、テストでは本物の代わりに偽物を渡せます。ただしrepositoryから下で偽物に差し替えているのは、一番下でDBにつないでいる部品（`db.Client`）だけです。repositoryが行を返す作りなので、repositoryの偽物を作るのが難しいためです（10節）。Slackへの送信のように、DB以外の部品を偽物にしているテストはあります。

新しいエンドポイントを足すときは、まずcontroller・usecase・repository（と必要ならentity）にファイルを足します。次に`di.go`で組み立て、`router.go`にルートを書きます。`di.go`には組み立てだけを書き、業務の判断はusecaseに書きます。

repositoryは行を返し、usecaseがそれをstructに詰めます。ここはSeeFTの少し変わったところです。よく知られたrepositoryの作り方では、repositoryが行をstructに詰めて返します。SeeFTでは、その仕事をusecaseがしています。

```go
// api/lib/usecase/bureau_usecase.go（抜粋）
// a はこの usecase 自身、a.rep はその repository、c はリクエストの context
rows, err := a.rep.All(c)          // repository は SQL を実行して行を返すだけ
...
for rows.Next() {                  // 行を1つずつ読む
    // 行の各列を struct の項目に書き込む。& は「この変数の場所」を渡す記号で、Scan はそこに値を書き込む
    err := rows.Scan(&bureau.ID, &bureau.Bureau, &bureau.Color, ...)
```

例外は、GORMを使う`shift_card_repository.go`です。ここだけはrepositoryがstructに詰めて返します。

#### 間違っていてもエラーが出ないところ

repositoryの多くは`SELECT *`で全列を取っていて、`Scan`に渡す変数はテーブルの列の順番と一致していなければなりません。テーブルに列を1つ足すと、そのテーブルを`SELECT *`して`Scan`している箇所をすべて直す必要があります。直し漏れは、コンパイルでは見つかりません。実行時に、列の数が合わなければ`Scan`がエラーを返します。列の数が合っていて順番だけずれた場合は、型が変換できる限りエラーにならず、違う列の値が入ります。

#### 層をまたいで使っている仕組み

- ミドルウェア（`lib/externals/server/server.go`）は、処理中の予期しない異常（panic）からの復帰、アクセスログ、CORSです。
- CORS。オリジンとは、`http://localhost:45029`のようなスキーム・ホスト・ポートの組です。ブラウザは、ページを配ったのとは別のオリジンに通信するとき、相手のサーバーが許可していなければ、その応答をページに渡しません。JSONをPOSTするときは、先に許可を問い合わせ、断られると本体の通信を送りません。mobileのページ（`http://localhost:45029`）からapi（`http://localhost:1234`）への通信は、ポートが違うので別のオリジンへの通信です。そこでapiは、許可するオリジンを`server.go`に列挙しています。`127.0.0.1`で開いたページは許可リストと一致しないので、画面に何も出ません（許可リストには`127.0.0.1:45029`とも書いてありますが、`http://`が抜けているので一致しません）。これはブラウザが守る決まりなので、GASや`curl`からの通信には関係しません。
- Slack通知（`lib/externals/slack/`、`lib/externals/scheduler/`）。シフトの変更は、apiが`action_logs`テーブルに記録しておき、5分ごとに動く処理がまとめて本人にDMします。レスキューの通知はその場で送ります。詳しくはAGENTS.mdの「通知の定期実行」にあります。
- 自動再ビルド。ローカルのコンテナは[air](https://github.com/air-verse/air)で起動しているので、`.go`ファイルを保存すると自動で再ビルドしてサーバーを起動し直します。`di.go`を変えてうまく動かないときは、`make up-api`を動かしているターミナルでCtrl+Cを押して止め、もう一度`make up-api`を実行してください。`make down`でも直りますが、DBの中身も消えます（2節）。

#### 手本にするファイルの選び方

古いファイルには、今は禁止している書き方が残っています。新しいファイルの中にも、一部の関数に古い書き方が残っていることがあります。全部の規約を満たした完全な手本は、まだありません。規約ごとに手本を選んでください。

| 規約 | まだ直っていない例 | 規約どおりの書き方 |
| --- | --- | --- |
| SQLはプレースホルダで書く | `bureau_repository.go`（文字列連結） | `review_repository.go` |
| エラーは`{"error": "..."}`のJSONで返す | `bureau_controller.go`（`return err`でEchoに任せる） | `review_controller.go` |
| 空の一覧は`null`ではなく`[]`を返す | `bureau_usecase.go`、`review_usecase.go`の`GetReviewsGAS`（0件で`nil`） | `var res = []entity.Review{}`のように、空の配列で初期化する |
| 見つからないときは404を返す | `review_controller.go`の`ShowReview`（500を返す） | `sql.ErrNoRows`（「行が無い」を表すエラー）を見て404にする |

`nil`の一覧をJSONにすると`null`になり、配列を期待しているmobileでは、`as List<dynamic>`のように配列として読むところで例外（`TypeError`）が出ます。

#### 学ぶこと

- HTTPそのもの。メソッド（GET / POST / PUT / DELETE）、ステータスコード（200 / 201 / 400 / 404 / 500）、ヘッダー、JSONボディ
- Echo。`e.GET(...)`によるルーティング、`c.Param`、`c.Bind`、`c.JSON`
- interfaceを受け取る設計と、それがテストのしやすさにつながる理由。SeeFTでその利点をどこまで生かせているかも、10節と合わせて考えてください
- CORSと、なぜブラウザだけが検査するのか

#### どこを見るか

`review`の一式（`review_controller.go` → `review_usecase.go` → `review_repository.go` → `entity/review.go`、`entity/reviewGAS.go`）が、一覧・取得・作成・更新・削除の揃った例です。上の表のとおり規約に合っていない箇所もあるので、それを見つけながら読んでください。repositoryにはテスト（`review_repository_sqlmock_test.go`）もあります。

エンドポイントの一覧は`router.go`を読んでください。`/swagger/index.html`にもAPIの説明ページがありますが、2025年7月から更新されておらず、今のエンドポイントとは一致しません。

---

## 6. DartとFlutter Web（mobile担当）

#### 何か

FlutterはGoogleのUIフレームワークで、言語はDartです。1つのコードからAndroid / iOS / Webなどを作れますが、SeeFTはWebとしてだけ動かしています。参加者はアプリをインストールせず、スマホのブラウザでURLを開きます。`mobile/`の下に`android/`・`ios/`・`macos/`・`windows/`も残っていますが、使っていません。

#### Flutterのバージョン

Flutterのバージョンは[fvm](https://fvm.app/)で固定しています（`mobile/.fvmrc`で3.27.3）。`flutter`コマンドは必ず`fvm flutter ...`の形で打ってください。素の`flutter`を打つと、手元に入っている別のバージョンが動き、依存ライブラリの解決やビルド結果が変わります。バージョンを書いている場所は3か所あり、正は`mobile/.fvmrc`です。本番のコンテナのビルド（`mobile/Dockerfile`の`FLUTTER_VERSION`）とCI（`.github/workflows/flutter-lint.yml`の`flutter-version`）はfvmを使わず、それぞれに同じ値を書いています。上げるときは3か所を同時に上げてください。1か所だけ上げると、手元・CI・本番で別のバージョンが動きます。

Flutter自体を上げる作業は、lintの導入や機能の追加などのほかの変更と混ぜず、専用のブランチで行います（#514）。3.27から3.35のように大きく上げると、不具合が出たときに、原因がバージョンなのか一緒に入れた変更なのかを切り分けられなくなるためです。

#### 画面はWidgetの木

画面のすべて（文字、ボタン、余白、並べ方）がWidgetで、Widgetの中にWidgetを入れて画面を組み立てます。`StatelessWidget`は状態を持たないもので、見た目は、渡された引数（と、アプリ全体のテーマのような周りの設定）で決まります。`StatefulWidget`は状態を持つもので、状態は、組になる`State`クラスのオブジェクトが持ちます。Widget自体は描き直すたびに作り直されますが、`State`は画面にある間ずっと残ります。

```dart
// 状態を変えるときは setState の中で。Flutter が build を呼び直して描き直す
setState(() => _isLoading = false);
```

#### 間違っていてもエラーが出ないところ

設定値はビルド時に埋め込まれます。apiのURLなどは`mobile/env/.env`に書き、`--dart-define-from-file=env/.env`で渡しています。コードからは`String.fromEnvironment('API_BASE_URL')`（`lib/configs/constant.dart`）で読みますが、これはビルドした時点で値がコードに埋め込まれる仕組みです。

- 値を変えても、ビルドし直すまで反映されません。
- `--dart-define-from-file`を渡し忘れても、エラーにはなりません。既定の`http://localhost:1234`に向きます。
- 埋め込まれた値はブラウザに配られるJavaScriptにそのまま入り、誰でも読めます。`.env`には公開してよい値だけを書き、**トークンやパスワードは絶対に書かないでください**。

#### 学ぶこと

- Dart。型、`final`と`const`、null安全（`String?`、`!`、`??`）、`async` / `await`と`Future`
- Widgetの木、`build`メソッド、`StatelessWidget`と`StatefulWidget`の違い
- レイアウト。`Column`、`Row`、`Padding`、`Expanded`、`ListView`
- `setState`、`initState`、`dispose`の役割
- `Navigator`による画面遷移
- Chromeの開発者ツールで、mobileが出す通信（Networkタブ）とログ（Consoleタブ）を見ること

#### やってみる

公式の[Write your first Flutter app](https://docs.flutter.dev/get-started/codelab)を、Webを対象にしてやってください。

---

## 7. mobileの構成と状態の持ち方（mobile担当）

`mobile/lib/`の下はこう分かれています。

```text
main.dart  起動の入口
pages/     画面。1ファイル1画面（my_shift_page.dart がマイシフト画面）。rescue/ はレスキューの画面群
widgets/   複数の画面で使う部品（shift_card.dart がシフトの1コマ）
models/    データの形
utils/     api.dart（API の呼び出し）、logger.dart、permanent_store.dart（端末への保存）
configs/   constant.dart（設定値）、importer.dart（よく使う import をまとめたもの）
theme/     tokens.dart（色・文字サイズ）、theme.dart（アプリ全体の見た目）
assets/    アプリのアイコンと、同梱した日本語フォント
```

#### APIの呼び出し

APIは1つのオブジェクト経由で呼びます。`lib/utils/api.dart`の一番上で`api`というオブジェクトを1つだけ作り、アプリ全体で共有しています（このように1つのオブジェクトを共有する作りを、シングルトンと呼ぶこともあります）。そこに、エンドポイントごとの関数（`api.getAllManual()`、`api.postReview(...)`など）が並んでいます。画面側は`import 'package:seeft_mobile/configs/importer.dart';`だけで`api`を使えます。HTTP通信の処理を`api.dart`だけにまとめるため、`package:http`を`api.dart`以外からimportするのは禁止です。

#### 状態の持ち方

状態管理ライブラリは使いません。RiverpodやProviderのような、画面をまたいで状態を共有するためのライブラリは入れず、各画面の`State`の中に状態を持ち、`setState`で描き直します。AGENTS.mdで導入を禁止しています。画面数が少なく、画面をまたいで共有する状態がほとんど無いので、それで足りています。

非同期の処理を待った後は`mounted`を確かめます。apiの応答を待っている間に利用者がその画面を閉じると、その画面の`State`は破棄されます。破棄された`State`で`setState`を呼ぶと例外になります。

```dart
final manuals = await api.getAllManual();
if (!mounted) return;   // 待っている間に画面が閉じられていたら何もしない
setState(() => _allManuals = manuals);
```

`manual_list_page.dart`の`_loadData`も、書き方は少し違いますが（`if (mounted) { setState(...) }`）、同じ考え方で`mounted`を確かめてから`setState`しています。

#### 端末への保存

端末への保存は2種類あります（`lib/utils/permanent_store.dart`）。1つはSharedPreferencesで、ログイン中のユーザーIDのような小さな値に使います。もう1つはHiveで、シフトカードやレスキューの一覧のような、まとまったデータに使います。マイシフト画面は、前回取得したシフトカードを日付ごとにHiveから読み出しています（`my_shift_page.dart`。保存のキーは日付と天気の組ですが、天気は晴れに固定しています）。

Webなので、どちらも実体はブラウザのストレージです。ブラウザがデータを消すこともあるので、消えても取り戻せるもの（apiから取り直せるデータや、ログインし直せば入るID）だけを置きます。

#### 間違っていてもエラーが出ないところ

保存したデータは、消えることより、古いまま残ることのほうが危険です。シフトが変わったのに前回の内容を表示し続けると、参加者は古いシフトのとおりに動いてしまいます。検証中に表示が古いままのときは、ブラウザのサイトデータを消すと直ることがあります。

#### 見た目とログの規約

色と文字サイズは`AppColors.main`や`AppFontSizes.md`（`lib/theme/tokens.dart`）を使い、`Color(0xFF...)`を直接書きません。ログは`print`ではなく`logger.i(...)` / `logger.e(...)`で出します。

#### 学ぶこと

- `StatefulWidget`のライフサイクル（`initState` → `build` → `dispose`）
- `FutureBuilder`と、`initState`でFutureを作っておく理由（`widgets/first_jump_selector.dart`のコメントに書いてあります）
- `mounted`がなぜ必要か
- SharedPreferencesとHiveの使い分け

#### どこを見るか

起動の流れは`main.dart` → `widgets/first_jump_selector.dart`（ログイン済みかを見て、ログイン画面かマイシフト画面に振り分ける）です。画面を1つ読むなら`pages/manual_list_page.dart`が、APIの呼び出し・読み込み中の表示・エラーの表示・検索が1画面に収まっていて手頃です。

---

## 8. Google Apps Script（GAS）とclasp（GASを触る人）

#### 何か

Googleのスプレッドシートに紐づけて動かせるJavaScriptです。SeeFTでは、スプレッドシートに作ったメニューから、名簿・タスク・シフトをapiに送ります。また、apiから届いたレスキューをスプレッドシートに書き込みます。

#### 間違っていてもエラーが出ないところ

`gas/`はクラウド上のコードの写しです。GASのコードは、各スプレッドシートに紐づいたGoogleのクラウド上のプロジェクト（スプレッドシートごとに1つ）にあり、動いているのはそちらです。リポジトリの`gas/`にあるのは、[clasp](https://github.com/google/clasp)（GASのコードをコマンドラインで取得・反映するツール）で取ってきた、ある時点の写しです。誰かがスプレッドシートのエディタで直接編集すると、クラウド上のプロジェクトだけが更新され、`gas/`は古いままになります。`gas/`の中には、45thでは使っていないディレクトリもあります（`task/`など）。

**GASを触る作業は、`gas/`を読むことではなく、`clasp`でクラウド上の最新のコードを取ってくることから始めます**。手順は[gas/README.md](../../gas/README.md)にあり、流れは次のとおりです。

1. `npm --prefix gas install`でclaspを入れ、リポジトリの外の一時ディレクトリに`clasp clone`で取ってくる
2. 取ってきたものと`gas/`を`diff`で比べ、ずれがあるかを先に確かめる
3. 編集と`clasp push`（クラウドへの反映）は、その一時ディレクトリで行う
4. pushしたら、同じ作業の中で`gas/`に写してcommitする

#### 反映と公開は別の操作

`clasp push`すると、スプレッドシートのメニューやトリガーで動くコードは、すぐ新しくなります。一方、Webアプリとして公開している入口（`doPost`）は、デプロイした時点の版で動き続けます。`clasp push`のあとに、今のデプロイのIDを指定して`clasp deploy -i <デプロイID>`で公開する版を更新しないと、apiから届くレスキューの処理は変わりません。IDを指定せずに`clasp deploy`すると、別のURLの新しいデプロイができるだけで、apiは古いほうを呼び続けます。エラーは出ません。

#### コードの誤りに気づく仕組みが無い

GASには、コードの誤りに気づく仕組みがありません。apiとmobileにはCIがありますが、`gas/`のコードを検査するCIはありません。JavaScriptには静的な型の検査も無く、`clasp push`で誤りのあるコードを反映しても、それを検出する仕組みがありません。しかもGASは全ファイルを1つのスコープで読み込みます。別ファイルで同じ名前の`const`を定義しただけでスクリプト全体が読み込めなくなり、スプレッドシートのメニューごと消えます。そのためgas/README.mdの手順では、全ファイルを1つにつなげてから`node --check`で構文を確かめます（1ファイルずつ確かめても、ファイルをまたぐ重複は見つかりません）。ただし、この検査で見つかるのは、`const`・`let`・`class`のどれかが絡む重複だけです。同じ名前の`function`どうし（や`var`どうし）はエラーにならず、後から読み込まれたほうが前のものを上書きします。

#### GAS特有の道具

- `UrlFetchApp.fetch(url, options)`。apiを呼ぶ
- `PropertiesService.getScriptProperties()`。apiのURLやトークンの置き場。コードに直接書かない
- `LockService`。同じスクリプトが同時に動かないように、実行を順番待ちにするロック。取得したら`finally`で必ず解放する。待つ時間には上限があります（SeeFTの多くは`waitLock(30000)`で30秒）。メニューを2回押すと、2回目は1回目が終わるのを待ち、30秒以内に終わらなければエラーで止まります。二重の送信を防ぐには、送る側・受ける側で重複を見分ける処理が別に要ります
- `doPost(e)`。GASをWebアプリとして公開し、外からPOSTを受ける入口。apiから届くレスキューをここで受ける
- `onChange(e)`。スプレッドシートが変更されたときに動かす関数。`onOpen`のように名前だけで自動で動く関数とは違い、トリガーとして登録しないと動きません

#### 学ぶこと

- JavaScriptの基礎（`const` / `let`、関数、配列の`map` / `filter`、オブジェクト、`JSON.stringify`）
- スプレッドシートの操作（`SpreadsheetApp`、`getRange`、`getValues` / `setValues`）。1セルずつ読み書きすると遅いので、範囲でまとめて読む
- GASの実行時間の上限（6分）。シフトの送信は件数が多く、1回では送り切れないので、途中で区切り、次の実行で続きから送る作りになっています（`gas/shift/コード.js`）
- claspの`clone`、`pull`、`push`、`deploy`
- Webアプリの「デプロイ」と「バージョン」

#### どこを見るか

`gas/rescue/コード.js`の`doPost`が、apiから届いたレスキューをスプレッドシートに書き込む処理です。`gas/rescue/onChange.js`が、スプレッドシートで変えた対応状況をapiに送り返す処理です。

---

## 9. データの流れを1本追う（全員）

リポジトリを読むときは、1つの操作を端から端まで追ってください。例として、mobileでタスクのレビューを送る操作を追います。ここでのレビューは、委員がタスクごとに付ける評価（人手が足りていたか、マニュアルは役に立ったか）のことで、PRのレビューとは別のものです。

1. `mobile/lib/widgets/review_bottom_sheet.dart`で、送信ボタンが`api.postReview(...)`を呼ぶ
2. `mobile/lib/utils/api.dart`が、`POST /reviews`にJSONを送る
3. `api/lib/router/router.go`の`e.POST("/reviews", r.reviewController.CreateReview)`
4. `api/lib/internals/controller/review_controller.go`の`CreateReview`が、JSONをstructに読み込み、usecaseを呼ぶ
5. `api/lib/usecase/review_usecase.go`の`CreateReview`が、タスク名からタスクのIDを引き、repositoryに保存させる
6. `api/lib/internals/repository/review_repository.go`が、`INSERT INTO reviews ...`をプレースホルダで実行する
7. `postgresql/db/schema/create4Review.sql`に、`reviews`テーブルの定義がある

手順2でDartが組み立てるJSONのキー名（`user_id`、`task_name`、`comment`など）と、手順4のcontrollerのstructタグが一致していることを確かめてください。0節で、apiとmobileの間の約束はURLとJSONの形だけだと書きました。その具体例がこれです。

mobileが`comment`を`comments`と書き間違えても、コンパイルは通ります。apiは知らないキー`comments`を捨て、コメントを空文字として保存します。エラーはどこにも出ません。一方、mobileが`staffing_rating`を`"5"`と文字列で送ると、`c.Bind`がエラーを返し、apiは400を返します。こちらは動かせば気づけます。

この追跡を一度丁寧にやると、チュートリアルを読むより多くのことが分かります。慣れたら、次の2本も追ってみてください。

- シフトの送信。`gas/shift/コード.js` → `POST /api/update_shifts` → `shift_controller.go`の`UpdateShiftsFromGAS` → `shift_usecase.go`。関わるファイルと処理が多く、SeeFTで最も重要な処理です
- レスキュー。`mobile/lib/pages/rescue/` → `POST /rescues` → `rescue_unified_controller.go`の`CreateRescue`（DBに保存する）→ `rescue_unified_usecase.go`の`SaveRescueToSpreadsheet` → `SendRescueToGAS` → `gas/rescue/コード.js`の`doPost`。apiがGASを呼ぶ、逆向きの通信が入っています

---

## 10. テスト（api担当・mobile担当）

#### api（Go）

標準の`go test`で動くテストに、[testify](https://github.com/stretchr/testify)（`assert` / `require`で結果を確かめるライブラリ）と[go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)（偽のDBを作るライブラリ）を足して書いています。テストファイルは、対象と同じディレクトリに`xxx_test.go`として置きます。手元では`make test`で全部のテストを実行できます。

今あるテストは、repositoryもusecaseも、実際のDBではなくgo-sqlmockの偽のDBを相手にしています。「このSQLがこの引数で実行されるはず」と先に宣言しておき、実行された中身を確かめる形です（1節のモック）。5節のDIのおかげで、`db.Client`を偽のDBに差し替えるだけで、本物のrepositoryをそのまま使ってテストできます。usecaseのテストでも、偽のrepositoryは作らず、偽のDBを載せた本物のrepositoryを渡しています。repositoryは行（`*sql.Rows`）を返す作りで、その偽物を作るのが難しいためです。

```go
// api/lib/internals/repository/review_repository_sqlmock_test.go（抜粋）
// 「この INSERT が、この5つの引数で実行されるはず」と宣言する
mock.ExpectExec(`INSERT INTO reviews \(user_id, task_id, staffing_rating, manual_rating, comment\)`).
    WithArgs("12", "34", "5", "4", testReviewComment).
    WillReturnResult(sqlmock.NewResult(1, 1))  // 実行されたら「1行追加した」と返す
```

引数の`"12"`や`"5"`が文字列なのは、reviewのrepositoryが値を文字列で受け取り、SQLの中で`$1::int`のように数値に変換しているからです。

#### 間違っていてもエラーが出ないところ

go-sqlmockが確かめられるのは「送ったSQLの文字列と引数」だけです。SQLが実際のテーブルで正しく動くか、`Scan`の列の順番が合っているかは確かめられません。テストが通っても、5節の`SELECT *`の問題は見つかりません。[test-roadmap.md](test-roadmap.md)では、repositoryを実際のDBで確かめるテスト（「フェーズ2」）も計画していますが、まだありません。

repositoryのテストは、`sqlmock_helper_test.go`の`newDBMock`を使うと短く書けます（repositoryのパッケージの中でだけ使えます）。どの層からどう書き進めるか、書く前にテストのケースの表を作ってレビューを受ける進め方は、test-roadmap.mdにあります。テストが書かれていない関数のissue（題名が`[test]`で始まるもの）が練習に向いています。

`review_repository_sqlmock_test.go`の冒頭には、アポストロフィ（`'`）を含むコメントのテストがあります。SQLを文字列連結で組み立てると、この`'`がSQLの文字列の終わりと解釈されて、SQLの文として正しくなくなり、DBがエラーを返します。エラーになるだけではありません。`'`の後ろにSQLを書いた入力を送れば、送った人がSQLの意味を書き換えられます。これがSQLインジェクションです。プレースホルダは、SQLの文と値を別々にDBに渡すので、値の中身がSQLとして解釈されることがありません。

#### mobile（Flutter）

`flutter_test`でwidgetテストを書いています。画面の部品を仮想的に描画して、表示や、押したときの動きを確かめます。

```bash
cd mobile && fvm flutter test
```

`test/widget_test.dart`はFlutterの雛形のカウンタを確かめる中身が残っていて、必ず落ちます。気にしないでください（#493で消す予定です）。widgetテストはCIでは走らないので、mobileを変えたら手元で実行してください。

#### 学ぶこと

- Goの`func TestXxx(t *testing.T)`、テーブル駆動テスト（入力と期待値の組を表にして並べる書き方）、`t.Run`、`assert`と`require`の違い
- go-sqlmockの`ExpectQuery` / `ExpectExec`、`WithArgs`、`WillReturnRows`、`ExpectationsWereMet`
- Flutterの`testWidgets`、`pumpWidget`、`find.text`、`tester.tap`
- 何をテストする価値があるか。境界の値や、過去に事故が起きた入力（アポストロフィを含む文字列、空の値）

GASには自動テストがありません。検証環境のスプレッドシートで実際に動かして確かめます（[staging-rehearsal.md](../operations/staging-rehearsal.md)）。

---

## 11. ツール、Docker、CI（全員）

- Docker / Docker Compose。ローカルではDBとapiがコンテナで動きます（`docker-compose.yml`）。mobileはコンテナではなく手元のfvmで動かします。本番もDocker Composeで、api・mobile・admin・外部公開用のトンネル（Cloudflare Tunnel）がコンテナで動いています。本番のDBはcomposeの外にあります（`docker-compose.prod.yml`、手順は[deploy.md](../operations/deploy.md)）。
- Makefile。よく使うdocker composeのコマンドに短い名前を付けたものです。一覧はリポジトリ直下の`Makefile`を開くのが確実です。Mac用のcomposeファイル（`docker-compose.mac.yml`）を使う`mac-`付きのコマンドもあります。
- GitHub Actions（CI）。PRを出すと、変更した場所に応じて次が走ります。
  - `api/`を変えたとき、golangci-lint（`go-lint.yml`。PRで新しく増えた指摘だけを見る）と`go test`（`go-test.yml`）
  - `mobile/`を変えたとき、`flutter analyze --fatal-infos`（`flutter-lint.yml`。最も軽いinfoレベルの指摘でも落ちる）。**テスト（`flutter test`）は走りません**。mobileを変えたら、10節のコマンドで手元で実行してください。45thでは、別々のPRで変わったボタンの文言とテストが食い違ったままdevelopに入り、マージ後に気づきました（#516）
  - `gas/`を変えたとき、GASのコードを検査するものは走らない
  - どのPRでも、ADR（`docs/decisions/`）に書いた前提のファイルや関数がまだあるかの点検（`docs-refcheck.yml`）
- lintは言語ごとに選んでいます。Goはgolangci-lint、Dartは`flutter analyze`（flutter_lintsのルール）です。ESLintはJavaScript専用なので、GoやDartには使えません。GAS（JavaScript）には今はlintがありません（`gas/package.json`に入っているのはclaspだけです）。入れるならESLintですが、`gas/`はある時点の写しなので（[ADR 0007](../decisions/0007-gas-live-is-source.md)）、claspでライブの版を取ってきてから、それに対して実行します。
- CodeRabbit。PRにAIがレビューコメントを付けます。指摘は参考です。すべてに従う必要はなく、スコープ外のものは理由を書いて見送って構いません。
- Git。issueを立て、`種類/名前/issue番号/内容`のブランチで作業し、PRを出します。種類は`feat`（機能）・`fix`（修正）・`docs`（文書）です。コミットメッセージは日本語で、`feat:` / `fix:` / `docs:`を付けます。PRはスカッシュマージ（コミットを1つにまとめる）でdevelopに入り、マージするのはPMです。詳しくは[開発の進め方](workflow.md)にあります。

#### 学ぶこと

- Dockerの基礎。image（コンテナの元になる型）とcontainerの違い、`docker compose up / down / logs / exec`、volume、port
- lintの指摘を読み、無効化のコメントで抑えるのではなく原因を直すこと
- workflowファイルを読める程度のYAML
- Gitのブランチ、merge、コンフリクトの解消、1つのPRに1つの目的

#### mobileのlintの違反を自動で直せるか

mobileのlintの違反を`dart fix`で直せるかは、ルール名からは分かりません。`dart fix`が直せるのは、その指摘に書き換えの手順が用意されているものだけで、Flutter自身の古い書き方の指摘でも、用意されていないことがあります。45thでは「自動で直せる」と分類した違反を新しく入った人に割り振ったところ、`dart fix`が何も直せませんでした（#288）。割り振る前に、`fvm flutter pub get`のあとで`fvm flutter analyze`の件数と`fvm dart fix --dry-run`の件数をルールごとに比べ、どこまで自動で直せるかを確かめてください（#286）。

---

## 12. 現実的な学習順序（全員）

全部を一度に学ぼうとしないでください。通読した後は、担当する場所から深く学べば十分です。下の段階1〜6の目安は、apiとmobileの両方を触る人でも2〜3週間です。13節で作るものの締め切りは、これとは別にPMと決めます。

| 段階 | 重点 |
| --- | --- |
| 1 | 2節で手元で動かす。GitとGitHubの流れ（issue → ブランチ → PR）を[workflow.md](workflow.md)で読む |
| 2 | 担当がapiならGoの基礎（A Tour of Go）、mobileならDartとFlutterの基礎 |
| 3 | apiはSQLの基礎と`schema/`の読解。`psql`でテーブルを眺める |
| 4 | apiは層の構成とDI。`review`の一式を読む。mobileは画面1つ（`manual_list_page.dart`）を読む |
| 5 | 9節のデータの流れを1本追う |
| 6 | テスト（go-sqlmockまたはwidgetテスト）を1本書く。DockerとCIの基礎 |
| 7 | GASとclasp（GASに関わるissueを持ったときでよい） |

---

## 13. 何か作る（全員）

読むだけでは足りません。同じ技術を使った自分の小さなシステムを、新しいリポジトリで作ってください。SeeFTとは関係ないもの（本棚、習慣の記録、サークルの備品の貸し出し表など）で構いません。大事なのは、つなぐ部分（画面とapi、apiとDB）を自分で書くことです。冒頭に挙げた「間違っていてもエラーが出ないところ」の多くも、このつなぐ部分にあります。

担当が片方だけなら、作る範囲を減らして構いません。apiだけの担当なら、画面は作らず、`curl`などでapiを呼んで動きを確かめます。mobileだけの担当なら、apiは作らず、SeeFTのローカルのapiを呼ぶ画面を作ります。題材はSeeFTのapiにあるもの（局の一覧、レビューの送信など）になります。自分のアプリは`http://localhost:45029`で開いてください。SeeFTのapiがCORSで許可しているオリジンの1つです（一覧は`api/lib/externals/server/server.go`）。SeeFTのmobileを止めてから、自分のアプリのディレクトリで`fvm flutter run -d web-server --web-port 45029 --dart-define-from-file=env/.env`を実行します（`Makefile`の`mobile-up`と同じ形です）。

締め切りや報告の仕方は、PMと相談して決めてください。

### 最低限

#### api（Go + Echo）

- Docker Compose上のPostgreSQL
- 関連する2つ以上のテーブル（1対多）を`CREATE TABLE`のSQLで定義し、リポジトリにコミットする
- controller / usecase / repositoryの3層とentityに分け、SeeFTの`di.go`と同じ形で手で組み立てる
- 1つのテーブルについて、一覧・1件取得・作成・更新・削除のエンドポイント
- SQLはすべてプレースホルダで書く
- 見つからないときは404、入力が不正なときは400を、`{"error": "..."}`のJSONで返す。SeeFT自身は、見つからないときの多くを500で返しています（5節の表）。SeeFTをまねるだけでは満たせません
- 空の一覧は`null`ではなく`[]`を返す
- go-sqlmockを使ったrepositoryのテストを1本以上

#### 画面（Flutter Web）

- 一覧画面と、一覧から開く詳細画面
- APIの呼び出しを1つのファイル（SeeFTの`api.dart`のような）にまとめる
- 読み込み中の表示と、失敗したときの表示
- 作成フォームを1つ。送信後に一覧が新しい内容で表示される
- `await`の後、`setState`の前に`mounted`を確かめる
- apiのURLは`--dart-define-from-file`で渡し、コードに直接書かない
- widgetテストを1本以上（`fvm flutter test`で通ること）

### 発展（時間があれば）

- スプレッドシートを1つ作り、GASから自分のapiにデータを送るメニューを付ける。URLは`PropertiesService`に置き、同時に2回動いたときに何が起きるかを確かめる
- 既存のテーブルに列を足すmigrationをgolang-migrateで書き、upとdownの両方を試す
- 一覧のデータをHiveに保存し、apiが止まっているときは前回の内容を表示する。そのとき、表示しているのが前回の内容だと利用者に分かるようにする（7節）
- 一定間隔で動く処理（SeeFTの`api/lib/externals/scheduler/`のような、決まった間隔で関数を呼ぶループ）をapiに足す
- PRごとにlintとテストを走らせるGitHub Actionsのworkflow

### 「完了」の意味

次を、調べずに説明できることです。多くは、この資料の「間違っていてもエラーが出ないところ」から出しています。括弧の中は、答えが書いてある節です。

全員が答えられるようにすること。

- apiとmobileの間でJSONのキー名が食い違ったとき、いつ・どうやって気づくか。値の型が食い違ったときとは何が違うか（0節、9節）
- SQLを文字列連結で組み立てると、どんな入力で何が起きるか（10節）

api担当が答えられるようにすること。

- SeeFTのテストで、本物のDBの代わりに何を、どこに差し込んでいるか。そのテストでは見つからない誤りは何か（10節）
- テーブルに列を1つ足したとき、Goのコードで直す必要があるのはどこか。直し忘れるといつ気づくか（5節）

mobile担当が答えられるようにすること。

- `await`の後で`mounted`を確かめないと何が起きるか（7節）
- `mobile/env/.env`の値を変えたのに画面が変わらないのはなぜか。そこにトークンを書いてはいけないのはなぜか（6節）

GASを触る人が答えられるようにすること。

- `gas/`のコードを読んで「今はこう動いている」と判断してはいけないのはなぜか。`clasp push`しただけでは変わらないのはどこか（8節）

これらに答えられるようになったら、実際のissueに取り組む準備ができています。

---

## 14. ブックマーク（全員）

- [A Tour of Go](https://go.dev/tour/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Echo](https://echo.labstack.com/docs)
- [`database/sql`のチュートリアル](https://go.dev/doc/tutorial/database-access)
- [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [PostgreSQLチュートリアル](https://www.postgresql.org/docs/current/tutorial.html)
- [MDNのCORSの解説](https://developer.mozilla.org/ja/docs/Web/HTTP/CORS)
- [Dartの言語ツアー](https://dart.dev/language)
- [Flutter](https://docs.flutter.dev/)
- [fvm](https://fvm.app/)
- [Hive](https://pub.dev/packages/hive_flutter)
- [Apps Script](https://developers.google.com/apps-script)
- [clasp](https://github.com/google/clasp)
- [Docker Compose](https://docs.docker.com/compose/)

リポジトリの中の文書は次のとおりです。

- [AGENTS.md](../../AGENTS.md)（コードの書き方の規約）
- [docs/operations/](../operations/README.md)（技大祭での運用と、過去の事故の記録）
- [workflow.md](workflow.md)（issueから本番反映までの仕事の進め方）
- [docs/decisions/](../decisions/README.md)（何をなぜ決めたか、何を選ばなかったかの記録。ADR）
- [test-roadmap.md](test-roadmap.md)（テストの進め方）
- [gas/README.md](../../gas/README.md)（GASの取得と反映）

---

## 15. 助けの求め方（全員）

しばらく自分で試してから、聞いてください。丸一日詰まったままにするのも、2分で聞くのも避けてください。聞くときは次を含めてください。

- 何をしようとしていたか
- 何が起きると思っていたか
- 実際には何が起きたか
- 正確なエラー文（スクリーンショットより、コピーした文字のほうが検索できて助かります）
- すでに試したこと

この形なら早く答えが返ってきます。担当しているissueで詰まったら、SlackでPMに連絡してください。issueを立てるときは、テンプレートに沿って書きます（[workflow.md](workflow.md)の「2. issue」）。

ただし、セキュリティの問題（たとえば、文字列連結のSQLに外から値を入れられる箇所を見つけた）は、公開のissueやPRに書かないでください。手口がそのまま公開されます。PMにSlackのDMで知らせてください。
