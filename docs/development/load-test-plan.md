# SeeFT API負荷試験 事前調査・試験計画

作成日: 2026-07-16

対象: SeeFT 45th production stack（`docker-compose.prod.yml`構成）

状態: 調査完了・試験計画ドラフト（未実行）。試験の実行は、issue化と環境準備を経てから行う。

参考: NUTFes Bingo本番構成・負荷テスト統合報告書（2026-07-03）。本書は、同報告書の「原因の切り分け」「内部直結vs公開URL経由の比較」「p95/p99でのpass/fail基準」のやり方を踏襲する。

---

## 1. ざっくり結論

調査の結果、負荷試験の焦点は次の5点に絞られる。

1. 当日の負荷はほぼmobileから来ており、いちばん重いのは朝に`GET /shift-cards`へ集中するアクセスである。mobileにはポーリングがなく（`Timer.periodic`ゼロ、手動リフレッシュのみ）、負荷は「人の操作」から生じる。技大祭当日の朝、実行委員300人以上が一斉にアプリを開くと、1人あたり1〜2リクエストが`GET /shift-cards`に集中する。このエンドポイントは、1リクエストで数十のクエリを何段にも分けて発行する（後述2.5）。
2. 全APIトラフィックがCloudflare Tunnelを通る。mobileはFlutter Web（ブラウザで実行）で、APIは`https://seeft-api.nutfes.net`という単独のホスト名でトンネル経由で公開されている。Bingoで1000 clientsが不合格になった主な原因の経路を、SeeFTでは静的ファイルだけでなく全APIリクエストが通る。したがって、「内部直結（Track A）vs検証用トンネル経由（Track B）」の切り分けが必須である。
3. 本番DBへ向けた試験は、最初から選択肢にない。本番DBはPatroni管理の共有HAクラスタで、Postgres（SeeFT）とMySQL（GM・FinanSu）が同じ3台の物理ノード上で同居している。APIはDB接続プールの上限も設定していない（`SetMaxOpenConns`なし = 無制限）。負荷をかければPostgres側の接続数が青天井で増え、同じノードのCPU・メモリ・スワップを使い切って、同居するMySQL側（GM・FinanSu）にまで影響しうる。試験は隔離DB（compose内の`postgres:18`）に対してだけ行う。
4. 外部への副作用を止めることが前提条件になる。`POST /rescues`は同期・タイムアウトなしでGAS（Googleスプレッドシート）へ送信し、APIプロセス内の通知スケジューラは5分間隔でSlack DMを送る。スタブURLと無効なトークンに差し替えずに試験すると、実際のスプシへの書き込みと実際のDMの送信が起きる。
5. 試験の前に直すべき欠陥が2つある。認証3エンドポイントで`Access-Token`ヘッダがないときに起きるpanic（500になる）と、DB接続プールの上限がないことである。前者は5xxの判定を不正確にし、後者は隔離DB上でも接続が尽きてエラーの出方を歪める。そのため、どちらも試験の測定の質に直接関わる。

なお、依頼時の前提のうち2点は現状と異なっていた。`api/go.mod`は2026-07-09のコミット`06935be`でGo 1.26.0に固定済みで、自動テストは`api/lib/usecase/shift_usecase_test.go`がある（テストロードマップのフェーズ1に着手済み）。負荷試験ツールを導入したことがない点は変わらない（リポジトリ内にk6 / vegeta / Locustの痕跡はない）。

---

## 2. 調査結果

### 2.1 エンドポイント棚卸し

`api/lib/router/router.go`の`ProvideRouter`に73ルートが定義されている。ルーターの外では、サーバー起動処理が`/swagger/*`を直接登録しており、合計74エンドポイントが公開される。

>追記（2026-09-29）：この節の件数と一覧は調査時点（2026-07-17、`2dd901a`）のもの。その後、マニュアル配信の3ルート（`GET /manuals/:id`・`GET /manuals/oauth/callback`・`PUT /manuals/:id`。`MANUAL_OAUTH_*`の設定時だけ登録）が2026-08-20に追加され、`POST /request_shifts`が#565で削除された。2026-09-29時点ではrouterに75ルート、`/swagger/*`を含めて76エンドポイントである。

```go
// api/lib/externals/server/server.go:42
e.GET("/swagger/*", echoSwagger.WrapHandler)
```

ミドルウェアはRecover・Logger・CORSの3つだけで（`server.go:21-36`）、認証ミドルウェアはない。認証は各ハンドラが個別に実装しており、後述の3エンドポイントに限られる。

呼び出し元は、次のように追って確定した。mobileは、全HTTP呼び出しが`mobile/lib/utils/api.dart`にまとめられている（AGENTS.mdの規約どおり）ため、このファイルから追った。adminは`admin/next-project/seeft-admin/src/`全体をgrepし、GASは`gas/`配下の`UrlFetchApp.fetch`呼び出しを追った。

以下は、リソースごとの一覧である。controllerは`api/lib/internals/controller/`、usecaseは`api/lib/usecase/`配下のファイルを指す。「呼び出し元」の *admin(凍結)* は、adminのコードに呼び出しはあるが、`AGENTS.md:11`で「使用していません」と明記された凍結画面であることを表す。

#### ヘルスチェック・Swagger

| ルート | 実装 | 認証 | 呼び出し元 |
|---|---|---|---|
| `GET /` | `health_controller.go:19` | 不要 | 監視用途 |
| `GET /swagger/*` | `server.go:42` | 不要 | 開発者 |

#### 認証（mail_auth）

| ルート | 実装 | 認証 | 呼び出し元 |
|---|---|---|---|
| `POST /mail_auth/signin` | `mail_auth_controller.go:29` → `mail_auth_usecase.go:34` | 不要 | mobile（ログイン画面） |
| `POST /mail_auth/web_signin` | `mail_auth_controller.go:56` → `mail_auth_usecase.go:124` | 不要 | admin(凍結) |
| `POST /mail_auth/web_signup` | `mail_auth_controller.go:39` → `mail_auth_usecase.go:73` | 不要 | admin(凍結) |
| `DELETE /mail_auth/web_signout` | `mail_auth_controller.go:66` → `mail_auth_usecase.go:175` | 必要 | admin(凍結) |
| `GET /mail_auth/web_is_signin` | `mail_auth_controller.go:76` → `mail_auth_usecase.go:183` | 必要 | なし（デッド） |

#### ユーザー

| ルート | 実装 | 認証 | 呼び出し元 |
|---|---|---|---|
| `GET /users` | `user_controller.go:30` → `user_usecase.go:35` | 不要 | admin(凍結)。mobile側参照は未ルーティング画面のみ |
| `GET /users/:id` | `user_controller.go:39` → `user_usecase.go:78` | 不要 | admin(凍結) |
| `POST /users` | `user_controller.go:49` → `user_usecase.go:112` | 不要 | admin(凍結) |
| `PUT /users/:id` | `user_controller.go:67` → `user_usecase.go:152` | 不要 | admin(凍結) |
| `DELETE /users` | `user_controller.go:85` → `user_usecase.go:214` | 不要 | admin(凍結) |
| `GET /current_user` | `user_controller.go:95` → `user_usecase.go:219` | 必要 | admin(凍結) |
| `POST /api/update_users` | `user_controller.go:107` → `user_usecase.go:268` | 不要 | GAS（`gas/user/コード.js:156`） |

#### シフト（mobile向け）

| ルート | 実装 | 認証 | 呼び出し元 |
|---|---|---|---|
| `GET /shifts/tasks/:task_id/years/:year_id/dates/:date_id/times/:time_id/weathers/:weather_id` | `shift_controller.go:33` → `shift_usecase.go:73` | 不要 | mobile（シフト表セルタップで最大3連続） |
| `GET /shift-cards/users/:user_id/dates/:date_id/weathers/:weather_id` | `shift_controller.go:46` → `shift_usecase.go:365` | 不要 | mobile（ホーム画面。最重要） |
| `POST /shift-cards` | `shift_controller.go:57` → `shift_usecase.go:365` | 不要 | なし（GETと同一処理のbody版。呼び出しゼロ） |
| `POST /request_shifts` | `shift_controller.go:154` → `shift_usecase.go:836,842` | 不要 | なし（デッド。後述。#565で削除） |

#### シフト（admin向け）

| ルート | 実装 | 認証 | 呼び出し元 |
|---|---|---|---|
| `POST /shifts-admin` | `shift_controller.go:82` → `shift_usecase.go:195` | 不要 | admin(凍結)。登録1回で768連続POST（後述） |
| `PUT /shifts-admin/:id` | `shift_controller.go:98` → `shift_usecase.go:222` | 不要 | admin(凍結) |
| `DELETE /shifts-admin` | `shift_controller.go:115` → `shift_usecase.go:271` | 不要 | admin(凍結) |
| `GET /shifts-admin/dates/:date/weathers/:weather` | `shift_controller.go:124` → `shift_usecase.go:299` | 不要 | admin(凍結) |
| `GET /shifts-admin/dates/:date/weathers/:weather/lower/:lower/upper/:upper` | `shift_controller.go:134` → `shift_usecase.go:332` | 不要 | admin(凍結) |
| `GET /shifts-admin/max-id` | `shift_controller.go:146` → `shift_usecase.go:817` | 不要 | admin(凍結) |
| `POST /api/update_shifts` | `shift_controller.go:174` → `shift_usecase.go:912` | 不要 | GAS（`gas/shift/コード.js:104,198,304`） |

#### タスク・マスタ系

| ルート | 実装 | 認証 | 呼び出し元 |
|---|---|---|---|
| `GET /tasks` | `task_controller.go:30` → `task_usecase.go:36` | 不要 | mobile（マニュアル一覧）・admin(凍結) |
| `GET /tasks/:id` | `task_controller.go:38` → `task_usecase.go:69` | 不要 | admin(凍結) |
| `GET /tasks/shifts/:shift` | `task_controller.go:47` → `task_usecase.go:96` | 不要 | なし（mobile側の呼び出しはコメントアウト済み） |
| `GET /tasks/users/:user_id` | `task_controller.go:57` → `task_usecase.go:129` | 不要 | mobile（救援フォームのタスク選択） |
| `POST /tasks` | `task_controller.go:67` → `task_usecase.go:162` | 不要 | admin(凍結) |
| `PUT /tasks/:id` | `task_controller.go:84` → `task_usecase.go:190` | 不要 | admin(凍結) |
| `DELETE /tasks` | `task_controller.go:102` → `task_usecase.go:241` | 不要 | admin(凍結) |
| `POST /api/update_tasks_and_places` | `task_controller.go:112` → `task_usecase.go:247` | 不要 | GAS（`gas/task/（内田）SeeFT送信.js:78,183`） |
| `GET /bureaus`・`GET /bureaus/:id` | `bureau_controller.go:23,31` → `bureau_usecase.go:24,54` | 不要 | admin(凍結) |
| `GET /grades`・`GET /grades/:id` | `grade_controller.go:23,31` → `grade_usecase.go:24,53` | 不要 | admin(凍結) |
| `GET /departments`・`GET /departments/:id` | `department_controller.go:23,31` → `department_usecase.go:24,53` | 不要 | admin(凍結) |
| `GET /times`・`GET /times/:id` | `time_controller.go:23,31` → `time_usecase.go:24,51` | 不要 | なし |
| `GET /places`・`GET /places/:id`・`POST /places`・`PUT /places/:id`・`DELETE /places` | `place_controller.go:26-67` → `place_usecase.go:27-137` | 不要 | admin(凍結) |

#### 救援（統一エンドポイント）

| ルート | 実装 | 認証 | 呼び出し元 |
|---|---|---|---|
| `POST /rescues` | `rescue_unified_controller.go:170` | 不要 | mobile（trouble / question / shorthanded送信） |
| `GET /rescues` | `rescue_unified_controller.go:332` → `rescue_unified_usecase.go:45` | 不要 | mobile（返答タブ「全体」） |
| `GET /rescues/users/:user_id` | `rescue_unified_controller.go:341` → `rescue_unified_usecase.go:78` | 不要 | mobile（返答タブ・個人） |

#### 救援（個別エンドポイント）

| ルート | 実装 | 認証 | 呼び出し元 |
|---|---|---|---|
| `PUT /question-rescues/:id` | `question_rescue_controller.go:87` → `question_rescue_usecase.go:131` | 不要 | GAS（`gas/rescue/onChange.js:86-98`） |
| `PUT /shorthanded-rescues/:id` | `shorthanded_rescue_controller.go:107` → `shorthanded_rescue_usecase.go:168` | 不要 | GAS（同上） |
| `PUT /trouble-rescues/:id` | `trouble_rescue_controller.go:106` → `trouble_rescue_usecase.go:164` | 不要 | GAS（同上） |
| 上記以外の個別系17ルート（GET一覧/1件/user別/task別、POST、DELETE × 3種） | `question_rescue_controller.go:30-110`ほか | 不要 | なし（デッド） |

#### レビュー

| ルート | 実装 | 認証 | 呼び出し元 |
|---|---|---|---|
| `POST /reviews` | `review_controller.go:50` → `review_usecase.go:95` | 不要 | mobile（シフトカードのレビュー送信） |
| `GET /reviews`・`GET /reviews/:id`・`PUT /reviews/:id`・`DELETE /reviews/:id` | `review_controller.go:28-107` | 不要 | なし（デッド） |

#### 棚卸しから言えること

74エンドポイントのうち、**実際にトラフィックが流れるのはmobile発9ルートとGAS発6ルートの計15ルート**である。残りは、admin(凍結)専用が30ルート（`web_signin`/`web_signup`/`web_signout`を含む。本番ではほぼトラフィックがない）、監視・開発用（healthcheck / swagger）が2ルート、呼び出し元がないデッドルートが27ルートである。負荷試験のシナリオは、現役の15ルートに絞ってよい。

>追記（2026-09-29）：この内訳も調査時点のもの。#565で`POST /request_shifts`を削除したので、デッドルートは26になった。調査後に追加されたマニュアル配信の3ルートは、この内訳に含めていない。

`POST /request_shifts`は特に目立つ例である。mobile・adminのどちらからも呼ばれておらず（grepゼロ）、実装はDB保存がno-opで、ハードコードされたGAS URLへの同期送信だけを行う。

```go
// api/lib/usecase/shift_usecase.go:836-840
func (u *shiftUseCase) SaveShiftData(ctx context.Context, req entity.ShiftRequest) error {
	// DB保存処理（仮実装）
	// 実際にはリポジトリを通じてDBに保存する
	return nil
}
```

テストロードマップの調査でも、このGAS URL（`shift_usecase.go:885`）に対応する`doPost`が`gas/shift/`にないことが指摘されている。そのため、エンドポイントごと閉じると判断できる（7章のissue提案参照）。

>追記：#565でエンドポイントごと削除した。送信先は`gas/`に写しの無い別プロジェクト（44thのテスト用スプレッドシートに紐づくもの）で、`doPost`はそちらにあった。45th以降は使っていない。

### 2.2 認証の実装

`Access-Token`ヘッダを検証するのは次の3箇所だけで、grepで漏れがないことを確認した。依頼時の想定と一致する。

```go
// api/lib/internals/controller/mail_auth_controller.go:68  (WebSignOut)
// api/lib/internals/controller/mail_auth_controller.go:78  (WebIsSignIn)
// api/lib/internals/controller/user_controller.go:97       (GetCurrentUser)
accessToken := c.Request().Header["Access-Token"][0]
```

3箇所とも同じ書き方で、ヘッダがないとmapアクセスが空スライスを返し、`[0]`でindex out of rangeのpanicを起こす。Recoverミドルウェアが500に変換するため、外から見るとエラー応答になる。しかし負荷試験では、「本来401相当の入力」が5xxとして数えられ、pass/failの判定が不正確になる。試験前の修正を推奨する（7章）。

トークンの発行から無効化までの流れを見ると、`POST /mail_auth/web_signin`がbcrypt（cost 10）でパスワードを照合し、成功したら既存セッションを削除してから、10文字のランダムトークンを`session`テーブルに保存する（`mail_auth_usecase.go:124-172`）。つまり、**同一ユーザーで再サインインすると旧トークンが無効になる**。認証付きエンドポイントを試験するときは、テストユーザーごとに1回だけサインインし、トークンを使い回す設計にする（4.5節）。

クライアント側を見ると、認証はほとんど使われていない。mobileはトークンをまったく扱わず、ログイン応答の`id` / `roleID`を端末に保存するだけで、以降の全リクエストは認証ヘッダなしで送られる（`mobile/lib/utils/api.dart`全域）。adminもトークンを送るのは`GET /current_user`と`DELETE /mail_auth/web_signout`の2呼び出しだけで、users / tasks / shifts-adminの全CRUDは認証ヘッダなしで呼ぶ。サーバー側も検証しないため、クライアントとサーバーで食い違いはない。ただし、書き込み系を含む大半のエンドポイントが認証なしで公開されていることは、負荷試験とは別に認識しておくべきである。

### 2.3 rescues統一/個別エンドポイントの実態

統一と個別のエンドポイントは、「移行途中の新旧並存」ではなく、**役割を分けて恒常的に並存している**ことが確定した。

作成と閲覧（mobile）は、統一エンドポイントだけを使う。`POST /rescues`（typeフィールドでtrouble / question / shorthandedを分岐、`mobile/lib/utils/api.dart:246,281,316`）、`GET /rescues`、`GET /rescues/users/:user_id`（同`:227,:212`）である。ステータスの更新（GAS）は、個別エンドポイントだけを使う。部門長がスプレッドシート上で対応状況を変えると、`onChange`トリガーが`PUT /{type}-rescues/:id`を呼ぶ。

```javascript
// gas/rescue/onChange.js:86
const url = baseUrl + "/" + type + "-rescues/" + changes[index].id;
```

adminは、統一・個別のどちらもまったく呼んでいない（src全体のgrepゼロ）。

したがって、統一系3ルートと個別系のPUT 3ルートは現役で、個別系の残り17ルート（GET / POST / DELETE）はデッドである。負荷試験では統一系3ルートを対象とし、個別のPUTはGASバッチのシナリオ（S3）に低い頻度で含める。

### 2.4 本番構成とDB経路

`docker-compose.prod.yml`で起動するのは、次の4つのコンテナである。

```text
cloudflare (cloudflared tunnel)  … トンネル資格情報はデプロイ先ホストの ./web/prod（リポジトリ外）
mobile     (Flutter Web, python server.py, :45029)
api        (Go/Echo, :1234, go run main.go)
admin      (Next.js, :5000→:3000, npm run dev)   … 凍結画面だが起動はされる
```

CORS設定（`server.go:34`）とadminの設定から、公開ホスト名は`https://seeft.nutfes.net`（mobile）、`https://seeft-admin.nutfes.net`（admin）、`https://seeft-api.nutfes.net`（API）の3つと分かる。

```js
// admin/next-project/seeft-admin/next.config.js:15-18
env: {
  SSR_API_URI: isProd ? 'https://seeft-api.nutfes.net' : 'http://nutfes-seeft-api:1234',
  CSR_API_URI: isProd ? 'https://seeft-api.nutfes.net' : 'http://localhost:1234'
}
```

mobileはFlutter Webとしてブラウザ上で動くため、API呼び出しはユーザーのブラウザから`seeft-api.nutfes.net`へ、すなわちCloudflare Tunnelを経由して届く。Bingoの検証で「内部直結1000は合格、トンネル経由1000は不合格」という結果が出た経路を、SeeFTでは初回のページ取得だけでなく、全APIリクエストが通る。これが、Track AとTrack Bを分けて測る理由である。

DBはこのcomposeに含まれない。APIは、環境変数だけで接続先を決める。

```go
// api/lib/externals/db/db.go:30-34（関数名 ConnectMySQL は歴史的名残で、実体は PostgreSQL 接続）
dbUser := os.Getenv("NUTMEG_DB_USER")
dbPassword := os.Getenv("NUTMEG_DB_PASSWORD")
dbHost := os.Getenv("NUTMEG_DB_HOST")
dbPort := os.Getenv("NUTMEG_DB_PORT")
dbName := os.Getenv("NUTMEG_DB_NAME")
```

本番の値は`api/env/seeft.env`（リポジトリ外）にあり、接続先はPatroni管理の共有HAクラスタである。Postgres（SeeFT）とMySQL（GM・FinanSu）は別クラスタ・別接続プールだが、同じ3台の物理ノード上で同居している。実際のアドレスはアクセス制御された運用資料の側にだけ記録し、本書には書かない。

#### 共有DBに向けた試験vs隔離DBに向けた試験

Bingo報告書が「トンネル経由vs内部直結」で失敗の原因を切り分けたのと同じ考え方で、DBについても2つの試験対象を区別する。ただし、結論は経路のときと異なる。経路は両方を測る価値があるが、**DBは共有クラスタに向けた試験を最初から実施しない**。理由は2つある。

第一に、同居するGM・FinanSu側まで巻き込むおそれが、構成上高い。APIは接続プールの上限を設定しておらず（`db.go:43`の`sql.Open`の後に`SetMaxOpenConns` / `SetMaxIdleConns` / `SetConnMaxLifetime`の呼び出しがない）、`database/sql`の既定値は「接続数無制限」である。GM・FinanSuは同じクラスタのMySQL側（別接続プール・別ポート）を使う。そのため、SeeFTがPostgresの接続数を使い切っても、両者の接続プールが直接取り合うわけではない。ただしMySQL Server・PostgreSQL・Patroni・etcdは同じ3台の物理ノード上で同居している。Postgres側の高負荷がCPU・メモリ・スワップを圧迫すれば、同じノードで動くMySQL側（GM・FinanSu）にもノイジーネイバーとして影響しうる。この種の資源の圧迫には、既に実例がある。MySQL側のスワップが尽きてノード全体が圧迫され、`systemctl restart mysql`でようやく解消した運用記録が残っている。この記録は、1つのDBエンジンの負荷が、同じノード上のほかのエンジンを巻き込みうることを裏付けている。

第二に、切り分けとしても不要である。DB単体の性能はクラスタ側の資源とチューニングで決まり、SeeFTの試験で知りたい「API実装のボトルネック」は、隔離DBでも同じように観測できる。共有クラスタに固有の挙動（フェイルオーバー、ほかのプロジェクトとの資源の取り合い）は、負荷試験ではなく運用監視で扱う領域である。

ただし、この2つの理由は「試験しない」ことの正当化であって、「隔離DBの結果が本番と同じ」という意味ではない。隔離された`postgres:18`は、ハードウェア・設定・同時実行の負荷が本番のPatroniクラスタと異なり、ほかのプロジェクトとの資源の取り合いも再現しない。したがって、本書の試験結果は **「API実装 + 隔離DB」という基準線**であり、本番クラスタでの実際の性能を保証するものではない。この前提は、結果を報告するときにも明記する。

隔離DBに使えるものは、既にリポジトリにある。`docker-compose.yml`の`db`サービス（`postgres:18`、`postgresql/db/`をinitdbにマウント、healthcheck付き）をそのまま使う。試験用スタックは、本番composeからDBの接続先だけを隔離DBに向けた構成で立てる。

シードデータは本番の規模から遠い（`postgresql/db/seed.sql`: users 3 / tasks 6 / shifts 0。timesは96スロットで本番と同じ）。試験には、規模を合わせたデータの生成が必須になる（4.3節）。

#### 静的配信のボトルネック（API試験とは別トラック）

mobileコンテナの実体は、次の11行である。

```python
# mobile/python/server.py
from http.server import SimpleHTTPRequestHandler
import socketserver

PORT = 45029

class MyHandler(SimpleHTTPRequestHandler):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, directory="./build/web", **kwargs)

with socketserver.TCPServer(("", PORT), MyHandler) as httpd:
    print(f"Serving at port {PORT}")
    httpd.serve_forever()
```

`socketserver.TCPServer`はシングルスレッドで（`ThreadingMixIn` / `ThreadingHTTPServer`不使用）、Flutter Webのバンドル（`main.dart.js` + canvaskitで数MB規模）を1接続ずつ順番に配る。朝の一斉アクセスでは、APIに届く前に、この静的配信が詰まる可能性が高い。Bingo報告書の「初回ページ取得の失敗」に当たる箇所は、SeeFTではここになる。本書のAPI試験計画とは別に、初回ロード試験（静的配信 + トンネル）を別トラックとして持つべきである。修正の候補（`ThreadingHTTPServer`化、またはnginx等への置き換え）も7章に含めた。

### 2.5 負荷特性上のホットスポット

#### GET /shift-cardsは1リクエストで数十クエリを発行する

`GetShiftCardsByUserAndDateAndWeather`（`shift_usecase.go:365`）は、grade / bureauのマップの事前ロードと、`GetOptimizedShiftData`による一括取得で部分的に最適化済みである。しかし、カードを生成する段階に多段のクエリが残っている。

カードの生成では、`createShiftCardFromGroup`（`shift_usecase.go:532`）が15分スロットごとに`getShiftMembersForTime` → `GetUsersByShift`（`shift_usecase.go:73`）を呼ぶ。`GetUsersByShift`は1回あたり、JOINクエリ1本とyear / date / time / weatherの`Find` 4本で、計5クエリを発行する。さらに、グループの前後の時間帯のメンバー取得（`getBeforeMembers` / `getAfterMembers`、`shift_usecase.go:659,717`）で、`GetUsersByShift`を2回分呼ぶ。

1日フルにシフトが入った部員（例: 8スロット × 2タスクグループ）で、おおよそ50〜60クエリ/リクエストに達する。300人以上が集合時刻の前後にアプリを開く場合、1人あたり1〜2リクエスト（signinの有無で変わる）を送るため、S1のramp（5分）で均せば、平均は1〜2 req/s程度にしかならない。この程度のreq/sでも、1リクエストが数十クエリに増えるため、DBには平均で毎秒100〜200クエリが届く計算になる。さらに、集合時刻の直前30秒〜1分にアクセスが偏った場合は、瞬間的に10〜20 req/s、DBには毎秒500〜1,000クエリ規模のバーストがあり得る。試験でいちばん重点的に見るのはここである。実際のバーストの幅は当日の運用（何時にアプリを開くよう案内するか）で変わり、確かな根拠はない。そのため試験では、ramp時間を5分/1分/30秒で切り替えて感度を確かめる（4.4節のVU段階とは別に、ramp時間そのものを変数として扱う）。

#### GET /rescuesのN+1

`GetAllRescues`（`rescue_unified_usecase.go:45`）は、3種の救援テーブルを全件取得した後、行ごとにユーザー名（`getUserName`、`rescue_unified_usecase.go:236`）とタスク名（`getTaskName`、同`:256`）を個別のクエリで引く。救援R件に対して3 + R×2クエリになる。当日の救援が数百件たまった状態では、返答タブが開かれるたびに数百のクエリが走る。試験データに現実的な件数の救援を入れておけば（4.3節）、この劣化を観測できる。

#### POST /rescuesの同期GAS送信（タイムアウトなし）

DBに保存する前に、送信者の情報を解決するために4クエリ（user / grade / bureau / task、`rescue_unified_controller.go:81-107`）を発行し、保存した後にスプレッドシートへ同期で送信する。

```go
// api/lib/usecase/rescue_unified_usecase.go:310-311
client := &http.Client{}
resp, err := client.Do(req)
```

`http.Client`に`Timeout`が設定されていないため、GAS側が遅れたり応答しなくなったりすると、リクエストを処理するgoroutineとDB接続が無期限にたまる。試験では、`RESCUE_GAS_URL`を必ずローカルのスタブに向ける（4.6節）。加えて、スタブに遅延を入れれば、「GASが遅いときにAPI全体がどう劣化するか」も測れる（S3シナリオの変種）。

#### GASバッチの逐次・非トランザクション処理

`POST /api/update_shifts`（`shift_usecase.go:912`）は、ユーザー・タスクの名前解決こそ一括化済みだが、変更行ごとに「既存確認 → 更新/作成 → action_log記録」を逐次実行し、全体をトランザクションで括っていない。数百行の変更を送ると、数百×3のクエリが1リクエストの中で直列に走る。S3シナリオで、読み取りのレイテンシへの影響を測る。

#### 通知スケジューラとaction_log

APIプロセス内で、5分間隔のSlack DMのflushが走る（`di.go:112`、`scheduler.go`）。シフトの作成・更新・削除では、action_logへの書き込みも起きる（`shift_usecase.go:271-297,1087-1123`）。試験中に`SLACK_BOT_TOKEN`が有効だと実際にDMが送られるため、無効にすることが必須である（4.6節）。なお、`SLACK_BOT_TOKEN`が未設定のときは、DIが通知スケジューラ自体を起動しない安全な作りになっている（`di.go:102-114`）。

#### 細かいが積もるもの

- `shift_repository.go`の一部のメソッド（`:58,69,80,91,102,120`）は、クエリを実行するたびに無条件で`fmt.Printf`する。`abstract_repository.go`は`DEBUG_SQL`でガード済みだが、shift系の直接実装に残っている。これは`GetUsersByShift`のホットパス上にあるため、毎秒数千行のstdout書き込みになる。Dockerのログドライバ経由では無視できないオーバーヘッドで、試験結果を歪める可能性がある。
- EchoのLoggerミドルウェアも、全リクエストでstdoutに書く。試験時のログの出力先と量は記録しておく。
- repository層に、文字列連結のSQLが残っている（既知issue #363 / #266）。プリペアドステートメントが再利用されないため、高QPSではパースのコストが上乗せされる。

---

## 3. 想定負荷パターン

### 3.1 前提として、負荷は「人の操作」から生じる

mobile / adminのどちらにも、ポーリング・自動リフレッシュがないことをコードで確認した（mobile: `Timer.periodic`等ゼロ、更新は画面のinitStateと手動リフレッシュボタンのみ。admin: `setInterval` / SWR / react-queryゼロ）。したがって、Bingoのような「全クライアントが一定間隔で叩き続ける」モデルではなく、ユーザー行動モデルでシナリオを組む。

技大祭当日のmobileのユーザー行動は、次の3つに整理できる。

1. 朝イチのシフト確認: アプリを開く。ログイン済みなら`GET /shift-cards` 1本、セッション切れ・初回なら`POST /mail_auth/signin` → `GET /shift-cards`の2本。全員がほぼ同じ時刻（集合時刻の前）に行う。
2. 日中の確認・操作: シフトの手動リフレッシュ、日付タブの切替（いずれも`GET /shift-cards`）、シフト表のセルタップ（`GET /shifts/tasks/...`が現在/前/後で最大3連続）、マニュアル閲覧（`GET /tasks`）、救援送信（`GET /tasks/users/:id` → `POST /rescues`）、返答確認（`GET /rescues`または`GET /rescues/users/:id`）、タスク終了後のレビュー（`POST /reviews`）。
3. 本部・スプシ側の操作（GAS）: 名簿・タスク・シフトの一括同期（`POST /api/update_*`）、救援対応ステータスの反映（`PUT /{type}-rescues/:id`）。人数は少ないが、1リクエストが重い。

read/write比はおおよそ97:3で、読み取りが大半を占める。書き込みで重いのは、GASバッチと`POST /rescues`（副作用込み）の2つである。

### 3.2 シナリオ定義

#### S1: 朝の一斉アクセス（最優先）

集合時刻前の5分間に、全員がアプリを開く状況を模す。

- ramp: 0 → 目標VUを5分で昇圧し、そのまま5分保持
- 各VUの行動: 30%が`POST /mail_auth/signin` → `GET /shift-cards`、70%が`GET /shift-cards`のみ（自動ログイン組）。その後、30〜60秒間隔で`GET /shift-cards`を再取得（リフレッシュ・タブ切替の模擬）
- 補足: signinはbcrypt cost 10の照合で、1回あたり数十msのCPUを使う。signin比率を100%にした変種は、CPU飽和の確認用に別に1回走らせる

#### S2: 日中定常

- 一定のVUで10分保持
- 1 VU・60秒のサイクルごとに、次の操作を独立試行として判定する（排他的な分岐ではないため、合計は100%を超えてよい）。
  - `GET /shift-cards`: 毎サイクル確実に1回
  - `GET /shifts/tasks/...`（現在/前/後の3リクエスト）: 30%の確率で発生
  - `GET /rescues`または`GET /rescues/users/:id`: 20%の確率で発生
  - `GET /tasks`: 10%の確率で発生
  - `GET /tasks/users/:id` → `POST /rescues`: 2%の確率で発生
  - `POST /reviews`: 2%の確率で発生
- 期待値で見ると、300 VU・1分あたり平均でshift-cards 300回、shifts/tasks系90回（3リクエスト相当で270）、rescues参照60回、tasks参照30回、rescue送信6回、review送信6回になる

#### S3: GASバッチ併走

S2を300 VUで回している最中に、`POST /api/update_shifts`（200〜700 changes程度の現実的なペイロード）を1本投入し、投入の前後で読み取りのp95 / p99がどう変わるかを測る。併せて、`PUT /question-rescues/:id`等の個別PUTを毎分数本混ぜる（スプシ側の対応反映の模擬）。

#### 優先度（当日のアクセス集中度とコード上の重さの積で判断）

```text
1. GET /shift-cards        （全員が朝に叩く × 1リクエスト数十クエリ）
2. GET /rescues            （本部・部員が随時開く × N+1）
3. POST /rescues           （トラブル集中時にバースト × 同期GAS送信）
4. GET /shifts/tasks/...   （セルタップで3連続）
5. POST /mail_auth/signin  （bcrypt CPU 負荷）
6. POST /api/update_shifts （低頻度 × 1本が重い）
admin 系・デッドルートは対象外（admin 凍結の経緯と 768 連続 POST の存在のみ 2.1 に記録）
```

---

## 4. 試験計画（ドラフト・未実行）

### 4.1 ツール選定（k6）

k6（Grafana k6）を採用する。理由は次のとおり。

- 単体のバイナリで導入でき、リポジトリに依存を持ち込まない（Go 1.26の固定やfvmの運用と干渉しない）
- シナリオ（scenarios / stages）で、S1のramp、S2の定常、S3の併走を宣言的に書ける
- p95 / p99とエラー率を`thresholds`として定義でき、pass/failをexit codeで機械的に判定できる（Bingo型の段階試験の自動化にそのまま使える）
- タグ付けによるエンドポイント別のレイテンシ集計が、標準機能にある

vegetaは固定レートの単発試験には優れるが、ユーザー行動モデル（サイクル内の確率分岐、signin → shift-cardsの依存関係）を表しにくい。Locustはシナリオの表現力では同等だが、Python環境の管理が増える。SeeFTの制約（負荷試験ツールの運用経験ゼロ、保守モード）のもとでは、スクリプト1ファイルとバイナリ1個で完結するk6が、最も撤退コストが低い。

### 4.2 試験環境

2トラック構成とし、どちらも**本番DB・本番Cloudflare Tunnelには一切向けない**。

#### Track A: 内部直結（API実装の性能を測る）

```text
負荷生成 (k6) → api:1234 (Docker network 直結) → db:5432 (postgres:18, 隔離)
```

- `docker-compose.yml`をベースに、api + db + GASスタブの試験用スタックを起動
- 負荷生成は、同じホストの別コンテナ、またはホスト上のk6。同じホストの資源を共有する影響を避けるため、可能なら負荷生成だけ別マシン（LAN内）から行う
- ここで得た限界値が「API実装 + DBの素の性能」。Bingoの「LXC内部直結1000合格」に当たる基準線

#### Track B: 検証用トンネル経由（本番経路の性能を測る）

```text
負荷生成 (k6, 外部) → Cloudflare edge → 検証用 tunnel → api:1234 → db:5432 (隔離)
```

- 検証用の一時トンネル（seeft-stg相当。本番トンネルとは別のtunnel credentialを発行）を用意できることは確認済み
- Track Aと同じスタック・同じデータ・同じシナリオで実施し、差分をトンネル経路の寄与として切り分ける
- Bingoの教訓（トンネル経由だと初回取得・大量同時接続がうまくいかなくなる/負荷生成元の回線品質が結果を不正確にする）を踏まえ、負荷生成元にはCloudflareへの経路が安定した外部VPSを優先して使う
- 試験の終了後は、tunnel credentialをローテーションする（Bingo報告書11章と同じ後始末）

#### やらないこと

- 本番DB（共有Patroniクラスタ）へ向けた試験。理由は2.4節
- 本番トンネル・本番ホスト名（seeft-api.nutfes.net等）へ向けた試験
- mobileの静的配信の初回ロード試験は、本計画のスコープ外の別トラック（2.4節末尾）

### 4.3 試験データ

隔離DBに45th想定規模のデータを投入する生成スクリプト（SQLまたはGo）を用意する。

| テーブル | 現行seed | 試験規模（案） | 根拠 |
|---|---|---|---|
| users | 3 | 400 | 実行委員300人以上（PM確認済み）+ 余裕 |
| tasks | 6 | 100〜150 | 45thの実タスク数をスプシから確認して合わせる |
| shifts | 0 | 2万〜5万 | 400人 × 2日 × 実働スロット数。GAS同期後の実数をスプシから確認 |
| question/shorthanded/trouble_rescues | 0 | 各100〜200 | GET /rescuesのN+1を現実の条件で再現するため必須 |
| session | 0 | テストユーザー分 | 認証3エンドポイント試験用 |

パスワードは全テストユーザーで固定値のbcryptハッシュを使い回し、生成を速くする。**着手前に45thの実データ規模（名簿人数・シフト行数・タスク数）をスプレッドシートで確認し、表の数値を確定させる**こと。

### 4.4 VU段階とpass/fail基準

実行委員300人以上という実利用規模に対し、安全係数を掛けた600 VUまでを段階的に昇圧して確かめる。Bingoの段階試験（600 VU安定/650 VU以上でレイテンシ超過）と同じ刻みを含めて、構成どうしを比べられるようにもする。

```text
smoke 10 VU（シナリオ動作確認） → 50 → 100 → 200 → 300 → 450 → 600 VU
各段階: 目標 VU 到達後 5 分保持。FAIL した段階で打ち切り、直前の PASS 段階を安定ラインとする
```

pass/failの基準は次の表のとおりで、段階ごとに全条件を満たせばPASSとする。

| 項目 | 基準 |
|---|---|
| HTTP失敗率（接続エラー・タイムアウト含む） | 1%未満 |
| 5xx応答 | 0件（panic・接続枯渇の検出を兼ねる） |
| `GET /shift-cards` p95 / p99 | 500 ms / 1,000 ms未満 |
| `GET /rescues` p95 | 500 ms未満 |
| その他読み取りp95 | 300 ms未満 |
| 書き込み（POST /rescues含む、スタブ応答込み）p95 | 1,000 ms未満 |
| 試験直後の`GET /` | 正常応答（Bingoのhealth/ready確認に相当） |

各段階で、次の項目を計測・記録して保存する。

- k6のサマリ（エンドポイントタグ別p50/p95/p99、req/s、失敗の内訳）
- `docker stats`の各コンテナのCPU/メモリの推移
- 隔離DBの`pg_stat_activity`接続数の最大値（プール無制限の問題の実測。修正後の効果測定にも使う）
- APIコンテナのstdoutログ量（無条件Printfの影響の確認）

### 4.5 認証付き3エンドポイントのトークン運用

1. 試験データの投入時に、テストユーザーN人を作成（bcryptハッシュ固定）
2. 試験のセットアップ（k6の`setup()`）で、各テストユーザーにつき1回だけ`POST /mail_auth/web_signin`を実行し、返ったトークンを配列で全VUに共有
3. VUはトークンを読み取り専用で使い回す。試験中に同一ユーザーで再サインインしない（`mail_auth_usecase.go:158`の既存セッション削除で旧トークンが無効になるため）
4. `Access-Token`ヘッダは必ず付ける。付け忘れると401ではなくpanic由来の500になり、サーバー実装の5xxと区別できなくなる（修正前に試験する場合の注意点）

認証3エンドポイントはadmin(凍結)専用のため優先度は低い。ただ、panic修正の回帰確認と、トークン検証クエリ（session lookup）の性能確認として、S2に低い頻度（1%程度）で`GET /current_user`を混ぜる。

### 4.6 外部副作用の遮断

試験用スタックでは、次の環境変数の差し替えを必須とする。

| 変数 | 試験時の値 | 遮断される副作用 |
|---|---|---|
| `RESCUE_GAS_URL` | ローカルスタブ（`https`必須。専用の信頼済みCAから発行した証明書 + 200固定応答の軽量サーバー） | POST /rescuesによる実際のスプシへの書き込み |
| `SLACK_BOT_TOKEN` | 未設定 | 通知スケジューラごと無効化（`di.go:102-114`で安全にスキップされる） |
| `NUTMEG_DB_*` | 隔離DB | 共有クラスタへの接続 |

注意点が3つある。第一に、`RESCUE_GAS_URL`は実装が`https`スキームを検証する（`rescue_unified_usecase.go:293-295`）ため、スタブもhttpsで立てる必要がある。第二に、送信処理は`&http.Client{}`の既定のTLS検証をそのまま使う（`rescue_unified_usecase.go:310`）ため、自己署名証明書をそのまま使うと`x509: certificate signed by unknown authority`で失敗する。スタブの証明書を発行したCAを、試験用APIコンテナの信頼ストアに追加して解決し（ローカルCAを発行し`update-ca-certificates`を通す等）、`InsecureSkipVerify`は使わない。第三に、`POST /request_shifts`はハードコードされたURL（`shift_usecase.go:885`）に送るため、環境変数では止められない。ただし呼び出し元がないため、試験対象から外せば実際に困ることはない（#565でエンドポイントごと削除済み）。

### 4.7 既知の制約と扱い

`api/`の自動テストは、shift_usecaseの純関数テストだけである（テストロードマップのフェーズ1の段階）。負荷試験は自動テストの代わりにはならず、機能の正しさはスコープ外とする。

Goは1.26.0に固定済みである（コミット`06935be`）。k6はAPIのビルドチェーンから独立しているので、影響はない。

チームには負荷試験ツールの運用経験がない。最初のsmoke（10 VU）を「ツールの学習」を兼ねた独立したステップとして扱い、いきなり段階試験に入らない。

---

## 5. リスクと対応

| リスク | 影響 | 対応 |
|---|---|---|
| DB接続プールが無制限のまま試験すると、高VUで接続数が数百に達し、DB側の`max_connections`超過のエラーが大量に出る | 「APIの限界」ではなく「設定不備の限界」を測ってしまう | 試験前に`SetMaxOpenConns`等の設定を入れる（7章issue案）。あえて未設定のまま1回測り、修正の効果と比べる選択肢もある |
| Access-Token panicにより、ヘッダなしのリクエストが500になる | 5xx=0の基準に、実装バグによる500が混ざる | 試験前に修正（7章）。修正までは、試験スクリプト側で必ずヘッダを付ける |
| POST /rescuesの同期GAS送信（タイムアウトなし） | スタブが遅いとgoroutineがたまり、全体が劣化する | スタブは即時応答を既定とし、遅延を入れる試験は独立した変種として実施 |
| 無条件の`fmt.Printf`とLoggerのstdout出力 | 高QPSでログI/Oが測定を歪める | ログ量を計測項目に含め、影響が見えたらガードを修正した後に測り直す |
| mobileの静的配信がシングルスレッド | 当日はAPIより先に初回ロードが詰まる | 本計画のスコープ外だが、別トラックのissueとして起票（7章）。当日の運用でも「集合時刻より前に開いておく」と案内してアクセスを分散 |
| 検証用トンネルでの負荷生成元の回線品質 | Bingoで観測された「負荷生成側の経路不良による偽FAIL」 | 外部VPSを負荷生成元に使い、FAILのときは生成元を変えて再現を確認 |
| 8月の引き継ぎ期限 | 試験〜修正〜再試験のループが1巡しかできない可能性 | Track Aを先に進め、修正が要るものを早めにissue化。Track BはTrack Aの合格後に1回で決める |

---

## 6. 依頼時前提の訂正

調査を始めたときに与えられた前提のうち、以下は現状と異なっていたため、訂正して記録する。

- `api/go.mod`がGo 1.16のまま → 2026-07-09のコミット`06935be`（issue #385対応）で`go 1.26.0`に固定済み。`AGENTS.md:7`の「Go 1.16」の表記が、更新されないまま残っている
- `api/`配下に自動テストが0件 → `api/lib/usecase/shift_usecase_test.go`がある（テストロードマップ フェーズ1の成果）
- 負荷試験ツールの導入実績なし → 変わらず（確認済み）

---

## 7. 次に立てるべきissueの提案

実行はissue化してから着手する運用に従い、次のように分けることを提案する。上から価値の高い順に並べた。

### issue 1: 負荷試験の実行環境整備

k6の導入手順（バイナリの配置とMakefileターゲット）、試験用compose（隔離DB + GASスタブ + 環境変数の差し替え）、45th規模の試験データ生成スクリプトを整備する。完了条件は「smoke 10 VUがS1/S2シナリオでgreenになること」。着手前に、45thの実データ規模（名簿・シフト行数）をスプシで確認することも含める。

### issue 2: 試験前修正（測定品質に直結する2件）

DB接続プールの上限の設定（`db.go`に`SetMaxOpenConns` / `SetMaxIdleConns` / `SetConnMaxLifetime`）と、`Access-Token`ヘッダがないときのpanicの修正（3箇所、`c.Request().Header.Get("Access-Token")` + 空チェックで401応答へ）。どちらも、小さく独立したPRにできる。panicの修正は、テストロードマップのフェーズ5（controllerテスト）の最初の題材を兼ねられる。

### issue 3: Track A（内部直結）試験の実行

本書4.4の段階昇圧を実施し、Bingo報告書と同じ形式で結果を記録する。FAILしたときのボトルネックの特定（`pg_stat_activity`、pprofの導入検討）まで含む。

### issue 4: Track B（検証用トンネル経由）試験の実行

検証用tunnel credentialの発行、外部VPSからの負荷生成、Track Aとの差分の分析、終了後のトークンのローテーション。

### issue 5（測ってから判断）: ホットスポット改修

`GET /shift-cards`の多段クエリの解消と、`GET /rescues`のN+1の解消（既存のN+1 issue #264 / #247と関連）。**Track Aの結果が基準を満たすなら着手しない**（保守モードの方針に従い、測定で必要性が証明されたときだけ直す）。同じく、`shift_repository.go`の無条件`fmt.Printf`の`DEBUG_SQL`ガード化と、`POST /rescues`のGAS送信のタイムアウト設定も、ここに含める。

### issue 6（別トラック）: mobile静的配信の初回ロード対策

`server.py`の`ThreadingHTTPServer`化（数行の変更）またはnginx等への置き換えの検討と、初回ロード（静的ファイル + トンネル）の負荷試験。API試験とは別に進められる。

### issue 7（任意・整理）: デッドエンドポイントの閉塞

`POST /request_shifts`（ハードコードされたGAS URL。#565で削除済み）、救援個別系の未使用17ルート、`GET /reviews`系ほか。攻撃される面を減らし、棚卸しの結果を固定することが目的で、負荷試験の前提条件ではない。
