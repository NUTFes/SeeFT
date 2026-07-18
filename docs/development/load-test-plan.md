# SeeFT API 負荷試験 事前調査・試験計画

作成日: 2026-07-16

対象: SeeFT 45th production stack（`docker-compose.prod.yml` 構成）

状態: **調査完了・試験計画ドラフト（未実行）**。試験の実行は issue 化と環境準備を経てから行う。

参考: NUTFes Bingo 本番構成・負荷テスト統合報告書（2026-07-03）。本書は同報告書の「原因の切り分け」「内部直結 vs 公開URL経由の比較」「p95/p99 での pass/fail 基準」の方法論を踏襲する。

---

## 1. ざっくり結論

調査の結果、負荷試験の焦点は次の5点に絞られる。

1. **当日の負荷はほぼ mobile 由来で、本丸は `GET /shift-cards` の朝バーストである。** mobile にはポーリングが存在せず（`Timer.periodic` ゼロ、手動リフレッシュのみ）、負荷は「人の操作」起点。技大祭当日の朝、実行委員 300 人以上が一斉にアプリを開いたとき、1人あたり1〜2リクエストが `GET /shift-cards` に集中する。このエンドポイントは1リクエストで数十クエリを発行する多段クエリ構造を持つ（後述 2.5）。
2. **全 API トラフィックが Cloudflare Tunnel を通る。** mobile は Flutter Web（ブラウザ実行）であり、API は `https://seeft-api.nutfes.net` として単独ホスト名でトンネル公開されている。Bingo で 1000 clients 不合格の主因となった経路が、SeeFT では静的ファイルだけでなく全 API リクエストに乗る。したがって「内部直結（Track A）vs 検証用トンネル経由（Track B）」の切り分けが必須。
3. **本番 DB へ向けた試験は最初から選択肢にない。** 本番 DB は Patroni 管理の共有 HA クラスタで、Postgres（SeeFT）と MySQL（GM・FinanSu）が同じ3台の物理ノード上で同居している。API は DB 接続プールの上限も設定していない（`SetMaxOpenConns` なし = 無制限）。負荷をかければ Postgres 側の接続数が青天井で伸び、同じノードの CPU・メモリ・スワップを食い潰して同居する MySQL 側（GM・FinanSu）にまで影響しうる。試験は隔離 DB（compose 内の `postgres:18`）に対してのみ行う。
4. **外部副作用の遮断が前提条件になる。** `POST /rescues` は同期・タイムアウトなしで GAS（Google スプレッドシート）へ送信し、API プロセス内の通知スケジューラは5分間隔で Slack DM を送る。スタブ URL と無効トークンへの差し替えなしに試験すると、実スプシ書き込みと実 DM が発生する。
5. **試験前に直すべき欠陥が2つある。** 認証3エンドポイントの `Access-Token` ヘッダ欠落時 panic（500 化）と、DB 接続プール無制限。前者は 5xx 判定を汚し、後者は隔離 DB 上でも接続枯渇でエラーモードを歪めるため、いずれも試験の測定品質に直結する。

なお、依頼時の前提のうち2点は現状と異なっていた。`api/go.mod` は 2026-07-09 のコミット `06935be` で Go 1.26.0 に固定済み、自動テストは `api/lib/usecase/shift_usecase_test.go` が存在する（テストロードマップのフェーズ1着手済み）。負荷試験ツールの導入実績がない点は変わらない（リポジトリ内に k6 / vegeta / Locust の痕跡なし）。

---

## 2. 調査結果

### 2.1 エンドポイント棚卸し

`api/lib/router/router.go` の `ProvideRouter` に 73 ルートが定義されている。ルーター外では、サーバー起動処理が `/swagger/*` を直接登録しており、合計 74 エンドポイントが公開される。

```go
// api/lib/externals/server/server.go:42
e.GET("/swagger/*", echoSwagger.WrapHandler)
```

ミドルウェアは Recover・Logger・CORS の3つのみで（`server.go:21-36`）、認証ミドルウェアは存在しない。認証はハンドラ個別実装で、後述の3エンドポイントに限られる。

呼び出し元は次の方法で確定した。mobile は全 HTTP 呼び出しが `mobile/lib/utils/api.dart` に集約されている（AGENTS.md の規約通り）ため同ファイルを起点に追跡、admin は `admin/next-project/seeft-admin/src/` 全体の grep、GAS は `gas/` 配下の `UrlFetchApp.fetch` 呼び出しを追跡した。

以下、リソースごとの一覧。controller は `api/lib/internals/controller/`、usecase は `api/lib/usecase/` 配下のファイルを指す。「呼び出し元」の *admin(凍結)* は、admin のコードに呼び出しは存在するが `AGENTS.md:11` で「使用していません」と明記された凍結画面であることを示す。

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
| `DELETE /mail_auth/web_signout` | `mail_auth_controller.go:66` → `mail_auth_usecase.go:175` | **必要** | admin(凍結) |
| `GET /mail_auth/web_is_signin` | `mail_auth_controller.go:76` → `mail_auth_usecase.go:183` | **必要** | なし（デッド） |

#### ユーザー

| ルート | 実装 | 認証 | 呼び出し元 |
|---|---|---|---|
| `GET /users` | `user_controller.go:30` → `user_usecase.go:35` | 不要 | admin(凍結)。mobile 側参照は未ルーティング画面のみ |
| `GET /users/:id` | `user_controller.go:39` → `user_usecase.go:78` | 不要 | admin(凍結) |
| `POST /users` | `user_controller.go:49` → `user_usecase.go:112` | 不要 | admin(凍結) |
| `PUT /users/:id` | `user_controller.go:67` → `user_usecase.go:152` | 不要 | admin(凍結) |
| `DELETE /users` | `user_controller.go:85` → `user_usecase.go:214` | 不要 | admin(凍結) |
| `GET /current_user` | `user_controller.go:95` → `user_usecase.go:219` | **必要** | admin(凍結) |
| `POST /api/update_users` | `user_controller.go:107` → `user_usecase.go:268` | 不要 | GAS（`gas/user/コード.js:156`） |

#### シフト（mobile 向け）

| ルート | 実装 | 認証 | 呼び出し元 |
|---|---|---|---|
| `GET /shifts/tasks/:task_id/years/:year_id/dates/:date_id/times/:time_id/weathers/:weather_id` | `shift_controller.go:33` → `shift_usecase.go:73` | 不要 | mobile（シフト表セルタップで最大3連続） |
| `GET /shift-cards/users/:user_id/dates/:date_id/weathers/:weather_id` | `shift_controller.go:46` → `shift_usecase.go:365` | 不要 | mobile（ホーム画面。**最重要**） |
| `POST /shift-cards` | `shift_controller.go:57` → `shift_usecase.go:365` | 不要 | なし（GET と同一処理の body 版。呼び出しゼロ） |
| `POST /request_shifts` | `shift_controller.go:154` → `shift_usecase.go:836,842` | 不要 | なし（デッド。後述） |

#### シフト（admin 向け）

| ルート | 実装 | 認証 | 呼び出し元 |
|---|---|---|---|
| `POST /shifts-admin` | `shift_controller.go:82` → `shift_usecase.go:195` | 不要 | admin(凍結)。登録1回で768連続 POST（後述） |
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
| `GET /tasks/shifts/:shift` | `task_controller.go:47` → `task_usecase.go:96` | 不要 | なし（mobile 側の呼び出しはコメントアウト済み） |
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
| `POST /rescues` | `rescue_unified_controller.go:170` | 不要 | mobile（trouble / question / shorthanded 送信） |
| `GET /rescues` | `rescue_unified_controller.go:332` → `rescue_unified_usecase.go:45` | 不要 | mobile（返答タブ「全体」） |
| `GET /rescues/users/:user_id` | `rescue_unified_controller.go:341` → `rescue_unified_usecase.go:78` | 不要 | mobile（返答タブ・個人） |

#### 救援（個別エンドポイント）

| ルート | 実装 | 認証 | 呼び出し元 |
|---|---|---|---|
| `PUT /question-rescues/:id` | `question_rescue_controller.go:87` → `question_rescue_usecase.go:131` | 不要 | GAS（`gas/rescue/onChange.js:86-98`） |
| `PUT /shorthanded-rescues/:id` | `shorthanded_rescue_controller.go:107` → `shorthanded_rescue_usecase.go:168` | 不要 | GAS（同上） |
| `PUT /trouble-rescues/:id` | `trouble_rescue_controller.go:106` → `trouble_rescue_usecase.go:164` | 不要 | GAS（同上） |
| 上記以外の個別系 17 ルート（GET 一覧/1件/user別/task別、POST、DELETE × 3種） | `question_rescue_controller.go:30-110` ほか | 不要 | なし（デッド） |

#### レビュー

| ルート | 実装 | 認証 | 呼び出し元 |
|---|---|---|---|
| `POST /reviews` | `review_controller.go:50` → `review_usecase.go:95` | 不要 | mobile（シフトカードのレビュー送信） |
| `GET /reviews`・`GET /reviews/:id`・`PUT /reviews/:id`・`DELETE /reviews/:id` | `review_controller.go:28-107` | 不要 | なし（デッド） |

#### 棚卸しから言えること

74 エンドポイントの内訳は、**実際にトラフィックが流れるのが mobile 発 9 ルート + GAS 発 6 ルートの計 15 ルート**、admin(凍結) 専用が 30 ルート（`web_signin`/`web_signup`/`web_signout` を含む。本番でほぼ無トラフィック）、監視・開発用（healthcheck / swagger）が 2 ルート、そして残り 27 ルートは呼び出し元が存在しないデッドルートである。負荷試験のシナリオは現役の 15 ルートに絞ってよい。

`POST /request_shifts` は特筆に値する。mobile・admin のどちらからも呼ばれておらず（grep ゼロ）、実装は DB 保存が no-op で、ハードコードされた GAS URL への同期送信だけを行う。

```go
// api/lib/usecase/shift_usecase.go:836-840
func (u *shiftUseCase) SaveShiftData(ctx context.Context, req entity.ShiftRequest) error {
	// DB保存処理（仮実装）
	// 実際にはリポジトリを通じてDBに保存する
	return nil
}
```

テストロードマップの調査でも、この GAS URL（`shift_usecase.go:885`）に対応する `doPost` が `gas/shift/` に存在しないことが指摘済みであり、エンドポイントごと閉塞する判断ができる（7章の issue 提案参照）。

### 2.2 認証の実装

`Access-Token` ヘッダを検証するのは次の3箇所のみで、grep による網羅確認済み。依頼時の想定と一致する。

```go
// api/lib/internals/controller/mail_auth_controller.go:68  (WebSignOut)
// api/lib/internals/controller/mail_auth_controller.go:78  (WebIsSignIn)
// api/lib/internals/controller/user_controller.go:97       (GetCurrentUser)
accessToken := c.Request().Header["Access-Token"][0]
```

3箇所とも同じ書き方で、ヘッダが欠落していると map アクセスが空スライスを返し `[0]` で **index out of range panic** を起こす。Recover ミドルウェアが 500 に変換するため外形上はエラー応答になるが、負荷試験では「本来 401 相当の入力」が 5xx として計上され、pass/fail 判定を汚す。試験前の修正を推奨する（7章）。

トークンのライフサイクルは次の通り。`POST /mail_auth/web_signin` が bcrypt（cost 10）でパスワード照合し、成功時に既存セッションを削除してから 10 文字のランダムトークンを `session` テーブルに保存する（`mail_auth_usecase.go:124-172`）。つまり**同一ユーザーで再サインインすると旧トークンが無効化される**。認証付きエンドポイントを試験する際は、テストユーザーごとに1回だけサインインしてトークンを使い回す設計にする（4.5 節）。

クライアント側の実態は、認証がほぼ使われていないことを示している。mobile はトークンを一切扱わず、ログイン応答の `id` / `roleID` を端末に保存するだけで、以降の全リクエストは認証ヘッダなしで送られる（`mobile/lib/utils/api.dart` 全域）。admin もトークンを送るのは `GET /current_user` と `DELETE /mail_auth/web_signout` の2呼び出しだけで、users / tasks / shifts-admin の全 CRUD は認証ヘッダなしで呼ぶ。サーバー側も検証しないため整合はしているが、書き込み系を含む大半のエンドポイントが無認証で公開されている事実は、負荷試験とは別に認識しておくべき状態である。

### 2.3 rescues 統一/個別エンドポイントの実態

「移行途中の新旧並存」ではなく、**役割分担による恒常的な並存**であることが確定した。

- 作成・閲覧（mobile）: 統一エンドポイントのみ。`POST /rescues`（type フィールドで trouble / question / shorthanded を分岐、`mobile/lib/utils/api.dart:246,281,316`）、`GET /rescues`、`GET /rescues/users/:user_id`（同 `:227,:212`）。
- ステータス更新（GAS）: 個別エンドポイントのみ。部門長がスプレッドシート上で対応状況を変更すると、`onChange` トリガーが `PUT /{type}-rescues/:id` を叩く。

```javascript
// gas/rescue/onChange.js:86
const url = baseUrl + "/" + type + "-rescues/" + changes[index].id;
```

- admin は統一・個別のどちらも一切呼んでいない（src 全体 grep ゼロ）。

したがって、統一系 3 ルートと個別系の PUT 3 ルートは現役、個別系の残り 17 ルート（GET / POST / DELETE）はデッドである。負荷試験では統一系 3 ルートを対象とし、個別 PUT は GAS バッチシナリオ（S3）に低頻度で含める。

### 2.4 本番構成と DB 経路

`docker-compose.prod.yml` で起動するのは4コンテナ。

```text
cloudflare (cloudflared tunnel)  … トンネル資格情報はデプロイ先ホストの ./web/prod（リポジトリ外）
mobile     (Flutter Web, python server.py, :45029)
api        (Go/Echo, :1234, go run main.go)
admin      (Next.js, :5000→:3000, npm run dev)   … 凍結画面だが起動はされる
```

公開ホスト名は CORS 設定（`server.go:34`）と admin の設定から、`https://seeft.nutfes.net`（mobile）、`https://seeft-admin.nutfes.net`（admin）、`https://seeft-api.nutfes.net`（API）の3つ。

```js
// admin/next-project/seeft-admin/next.config.js:15-18
env: {
  SSR_API_URI: isProd ? 'https://seeft-api.nutfes.net' : 'http://nutfes-seeft-api:1234',
  CSR_API_URI: isProd ? 'https://seeft-api.nutfes.net' : 'http://localhost:1234'
}
```

mobile は Flutter Web としてブラウザ上で動くため、API 呼び出しはユーザーのブラウザから `seeft-api.nutfes.net` へ、すなわち **Cloudflare Tunnel を経由して**届く。Bingo の検証で「内部直結 1000 は合格、トンネル経由 1000 は不合格」という結果が出た経路が、SeeFT では初回ページ取得だけでなく全 API リクエストに適用される。これが Track A / Track B を分けて測る理由である。

DB はこの compose に含まれない。API は環境変数のみで接続先を決める。

```go
// api/lib/externals/db/db.go:30-34（関数名 ConnectMySQL は歴史的名残で、実体は PostgreSQL 接続）
dbUser := os.Getenv("NUTMEG_DB_USER")
dbPassword := os.Getenv("NUTMEG_DB_PASSWORD")
dbHost := os.Getenv("NUTMEG_DB_HOST")
dbPort := os.Getenv("NUTMEG_DB_PORT")
dbName := os.Getenv("NUTMEG_DB_NAME")
```

本番値は `api/env/seeft.env`（リポジトリ外）にあり、接続先は Patroni 管理の共有 HA クラスタである。Postgres（SeeFT）と MySQL（GM・FinanSu）は別クラスタ・別接続プールだが、同じ3台の物理ノード上で同居している。実アドレスはアクセス制御された運用資料側にのみ記録し、本書には記載しない。

#### 共有 DB に向けた試験 vs 隔離 DB に向けた試験

Bingo 報告書が「トンネル経由 vs 内部直結」で失敗要因を切り分けたのと同じ発想で、DB についても2つの試験対象を区別する。ただし結論が異なる。経路の切り分けは両方測る価値があるが、**DB は共有クラスタに向けた試験を最初から実施しない**。理由は2つ。

第一に、巻き込み事故のリスクが構造的に高い。API は接続プールの上限を設定しておらず（`db.go:43` の `sql.Open` 後に `SetMaxOpenConns` / `SetMaxIdleConns` / `SetConnMaxLifetime` の呼び出しが存在しない）、`database/sql` の既定値は「接続数無制限」である。GM・FinanSu は同じクラスタの MySQL 側（別接続プール・別ポート）を使うため、SeeFT が Postgres の接続数を食い潰しても両者の接続プールが直接競合するわけではない。ただし MySQL Server・PostgreSQL・Patroni・etcd は同じ3台の物理ノード上で同居しており、Postgres 側の高負荷が CPU・メモリ・スワップを圧迫すれば、同じノードで動く MySQL 側（GM・FinanSu）にもノイジーネイバーとして影響しうる。この種の資源圧迫は既に実例がある。MySQL 側のスワップ枯渇でノード全体が圧迫され、`systemctl restart mysql` でようやく解消した運用記録が残っており、1つの DB エンジンの負荷が同じノード上の他エンジンを巻き込みうることを裏付けている。

第二に、切り分けとしても不要である。DB 単体の性能はクラスタ側の資源とチューニングに支配され、SeeFT の試験で知りたい「API 実装のボトルネック」は隔離 DB でも同じように観測できる。共有クラスタ固有の挙動（フェイルオーバー、他プロジェクトとの資源競合）は負荷試験ではなく運用監視の領域である。

ただし、この2つの理由は「試験しない」ことの正当化であって、「隔離 DB の結果が本番と同じ」という意味ではない。隔離された `postgres:18` はハードウェア・設定・同時実行負荷が本番 Patroni クラスタと異なり、他プロジェクトとの資源競合も再現しない。したがって本書の試験結果は**「API 実装 + 隔離 DB」という基準線**であり、本番クラスタでの実性能を保証するものではない。この前提は結果を報告する際にも明記する。

隔離 DB の土台は既にリポジトリにある。`docker-compose.yml` の `db` サービス（`postgres:18`、`postgresql/db/` を initdb にマウント、healthcheck 付き）をそのまま使い、試験用スタックは本番 compose から DB 接続先だけを隔離 DB に向けた構成で立てる。

シードデータは本番規模に遠い（`postgresql/db/seed.sql`: users 3 / tasks 6 / shifts 0。times は 96 スロットで本番同等）。試験には規模を合わせたデータ生成が必須になる（4.3 節）。

#### 静的配信のボトルネック（API 試験とは別トラック）

mobile コンテナの実体は次の 11 行である。

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

`socketserver.TCPServer` は **シングルスレッド**であり（`ThreadingMixIn` / `ThreadingHTTPServer` 不使用）、Flutter Web バンドル（`main.dart.js` + canvaskit で数 MB 規模）を1接続ずつ順番に配る。朝の一斉アクセスでは、API に到達する前にこの静的配信が詰まる可能性が高い。Bingo 報告書の「初回ページ取得の失敗」に相当する箇所が SeeFT ではここになる。本書の API 試験計画とは独立に、初回ロード試験（静的配信 + トンネル）を別トラックとして持つべきであり、修正候補（`ThreadingHTTPServer` 化ないし nginx 等への置き換え）も 7 章に含めた。

### 2.5 負荷特性上のホットスポット

#### GET /shift-cards — 1リクエストで数十クエリ

`GetShiftCardsByUserAndDateAndWeather`（`shift_usecase.go:365`）は grade / bureau のマップ事前ロードと `GetOptimizedShiftData` による一括取得で部分的に最適化済みだが、カード生成段階に多段クエリが残っている。

- `createShiftCardFromGroup`（`shift_usecase.go:532`）が 15 分スロットごとに `getShiftMembersForTime` → `GetUsersByShift`（`shift_usecase.go:73`）を呼ぶ
- `GetUsersByShift` は1回あたり JOIN クエリ1本 + year / date / time / weather の `Find` 4本 = **5 クエリ**
- グループ前後の時間帯メンバー取得（`getBeforeMembers` / `getAfterMembers`、`shift_usecase.go:659,717`）でさらに `GetUsersByShift` 2回分

1日フルにシフトが入った部員（例: 8 スロット × 2 タスクグループ）で、おおよそ 50〜60 クエリ / リクエストに達する。300 人以上が集合時刻前後にアプリを開く場合、1 人あたり 1〜2 リクエスト（signin の有無で変動）を送るため、S1 の ramp（5 分）で均せば平均は 1〜2 req/s 程度にしかならない。この程度の req/s でも、1 リクエストが数十クエリに増幅されるため、DB には平均で毎秒 100〜200 クエリが届く計算になる。さらに集合時刻の直前 30 秒〜1 分にアクセスが偏った場合は瞬間的に 10〜20 req/s、DB には毎秒 500〜1,000 クエリ規模のバーストがあり得る。ここが試験の主戦場になる。実際のバースト幅は当日の運用（何時にアプリを開くよう案内するか）に依存し確たる根拠がないため、試験では ramp 時間を 5 分/1 分/30 秒で切り替えて感度を確認する（4.4 節の VU 段階とは別に、ramp 時間そのものを変数として扱う）。

#### GET /rescues — N+1

`GetAllRescues`（`rescue_unified_usecase.go:45`）は3種の救援テーブルを全件取得した後、**行ごとに** ユーザー名（`getUserName`、`rescue_unified_usecase.go:236`）とタスク名（`getTaskName`、同 `:256`）を個別クエリで引く。救援 R 件に対し 3 + R×2 クエリ。当日の救援が数百件たまった状態で返答タブが開かれるたびに数百クエリが走る。試験データに救援件数を現実的に積むこと（4.3 節）で、この劣化が観測できる。

#### POST /rescues — 同期 GAS 送信（タイムアウトなし）

DB 保存前に送信者情報の解決で 4 クエリ（user / grade / bureau / task、`rescue_unified_controller.go:81-107`）、保存後にスプレッドシートへ同期送信する。

```go
// api/lib/usecase/rescue_unified_usecase.go:310-311
client := &http.Client{}
resp, err := client.Do(req)
```

`http.Client` に `Timeout` が設定されていないため、GAS 側が遅延・無応答になるとリクエストを処理する goroutine と DB 接続が無期限に滞留する。試験では `RESCUE_GAS_URL` を必ずローカルスタブに向ける（4.6 節）。加えて、スタブに遅延を注入すれば「GAS が遅いときに API 全体がどう劣化するか」も測れる（S3 シナリオの変種）。

#### GAS バッチ — 逐次・非トランザクション

`POST /api/update_shifts`（`shift_usecase.go:912`）はユーザー・タスクの名前解決こそ一括化済みだが、変更行ごとに 既存確認 → 更新/作成 → action_log 記録 を逐次実行し、全体をトランザクションで括っていない。数百行の変更を送ると数百×3 クエリが1リクエスト内で直列に走る。S3 シナリオで読み取りレイテンシへの影響を測る。

#### 通知スケジューラと action_log

API プロセス内で 5 分間隔の Slack DM flush が走る（`di.go:112`、`scheduler.go`）。シフトの作成・更新・削除は action_log への書き込みを伴う（`shift_usecase.go:271-297,1087-1123`）。試験中に `SLACK_BOT_TOKEN` が有効だと実 DM が飛ぶため、無効化が必須（4.6 節）。なお `SLACK_BOT_TOKEN` 未設定時は DI が通知スケジューラ自体を起動しない安全設計になっている（`di.go:102-114`）。

#### 細かいが積もるもの

- `shift_repository.go` の一部メソッド（`:58,69,80,91,102,120`）はクエリ実行のたびに無条件で `fmt.Printf` する。`abstract_repository.go` は `DEBUG_SQL` ガード済みだが、shift 系の直接実装に残っており、`GetUsersByShift` のホットパス上にあるため毎秒数千行の stdout 書き込みになる。Docker のログドライバ経由では無視できないオーバーヘッドで、試験結果を歪める可能性がある。
- Echo の Logger ミドルウェアも全リクエストで stdout に書く。試験時のログ出力先と量は記録しておく。
- repository 層に文字列連結 SQL が残る（既知 issue #363 / #266）。プリペアドステートメント再利用が効かないため、高 QPS ではパースコストが上乗せされる。

---

## 3. 想定負荷パターン

### 3.1 前提: 負荷は「人の操作」起点

mobile / admin ともにポーリング・自動リフレッシュが存在しないことをコードで確認した（mobile: `Timer.periodic` 等ゼロ、更新は画面 initState と手動リフレッシュボタンのみ。admin: `setInterval` / SWR / react-query ゼロ）。したがって Bingo のような「全クライアントが一定間隔で叩き続ける」モデルではなく、**ユーザー行動モデル**でシナリオを組む。

技大祭当日の mobile ユーザー行動は次の3つに整理できる。

1. **朝イチのシフト確認**: アプリを開く。ログイン済みなら `GET /shift-cards` 1本、セッション切れ・初回なら `POST /mail_auth/signin` → `GET /shift-cards` の2本。全員がほぼ同時刻（集合時刻前）に行う。
2. **日中の確認・操作**: シフトの手動リフレッシュ、日付タブ切替（いずれも `GET /shift-cards`）、シフト表セルタップ（`GET /shifts/tasks/...` が現在/前/後で最大3連続）、マニュアル閲覧（`GET /tasks`）、救援送信（`GET /tasks/users/:id` → `POST /rescues`）、返答確認（`GET /rescues` または `GET /rescues/users/:id`）、タスク終了後のレビュー（`POST /reviews`）。
3. **本部・スプシ側の操作（GAS）**: 名簿・タスク・シフトの一括同期（`POST /api/update_*`）、救援対応ステータスの反映（`PUT /{type}-rescues/:id`）。人数は少ないが1リクエストが重い。

read/write 比はおおよそ 97:3 で読み取り支配。書き込みで重いのは GAS バッチと `POST /rescues`（副作用込み）の2つ。

### 3.2 シナリオ定義

**S1: 朝の一斉アクセス（最優先）**

集合時刻前の 5 分間に全員がアプリを開く状況を模す。

- ramp: 0 → 目標 VU を 5 分で昇圧、そのまま 5 分保持
- 各 VU の行動: 30% が `POST /mail_auth/signin` → `GET /shift-cards`、70% が `GET /shift-cards` のみ（自動ログイン組）。その後 30〜60 秒間隔で `GET /shift-cards` を再取得（リフレッシュ・タブ切替の模擬）
- 補足: signin は bcrypt cost 10 の照合で1回あたり数十 ms の CPU を食う。signin 比率を 100% にした変種は CPU 飽和の確認用に別途1回走らせる

**S2: 日中定常**

- 一定 VU で 10 分保持
- 1 VU・60 秒サイクルごとに、次の操作を**独立試行**として判定する（排他的な分岐ではないため合計は100%を超えてよい）:
  - `GET /shift-cards`: 毎サイクル確実に1回
  - `GET /shifts/tasks/...`（現在/前/後の3リクエスト）: 30%の確率で発生
  - `GET /rescues` または `GET /rescues/users/:id`: 20%の確率で発生
  - `GET /tasks`: 10%の確率で発生
  - `GET /tasks/users/:id` → `POST /rescues`: 2%の確率で発生
  - `POST /reviews`: 2%の確率で発生
- 期待値ベースでは、300 VU・1分あたり平均で shift-cards 300回、shifts/tasks 系 90回（3リクエスト相当で270）、rescues参照 60回、tasks参照 30回、rescue送信 6回、review送信 6回になる

**S3: GAS バッチ併走**

- S2 を 300 VU で回している最中に、`POST /api/update_shifts`（200〜700 changes 程度の現実的なペイロード）を1本投入し、投入前後の読み取り p95 / p99 の変化を測る
- 併せて `PUT /question-rescues/:id` 等の個別 PUT を毎分数本混ぜる（スプシ側の対応反映の模擬）

**優先度**（当日のアクセス集中度とコード上の重さの積で判断）

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

### 4.1 ツール選定: k6

k6（Grafana k6）を採用する。理由は次の通り。

- 単体バイナリで導入でき、リポジトリに依存を持ち込まない（Go 1.26 固定や fvm 運用と干渉しない）
- シナリオ（scenarios / stages）で S1 の ramp、S2 の定常、S3 の併走を宣言的に書ける
- p95 / p99 とエラー率を `thresholds` として定義でき、pass/fail が exit code で機械判定できる（Bingo 型の段階試験の自動化に直結）
- タグ付けでエンドポイント別のレイテンシ集計が標準機能

vegeta は固定レートの単発試験には優れるが、ユーザー行動モデル（サイクル内の確率分岐、signin → shift-cards の依存関係）を表現しにくい。Locust はシナリオ表現力では同等だが Python 環境の管理が増える。SeeFT の制約（負荷試験ツールの運用経験ゼロ、保守モード）では、スクリプト1ファイル + バイナリ1個で完結する k6 が最も撤退コストが低い。

### 4.2 試験環境

2トラック構成とし、いずれも**本番 DB・本番 Cloudflare Tunnel には一切向けない**。

**Track A: 内部直結（API 実装の性能を測る）**

```text
負荷生成 (k6) → api:1234 (Docker network 直結) → db:5432 (postgres:18, 隔離)
```

- `docker-compose.yml` をベースに、api + db + GAS スタブの試験用スタックを起動
- 負荷生成は同一ホストの別コンテナまたはホスト上の k6。ホスト資源の同居影響を避けるため、可能なら負荷生成だけ別マシン（LAN 内）から行う
- ここで得た限界値が「API 実装 + DB の素の性能」。Bingo の「LXC 内部直結 1000 合格」に相当する基準線

**Track B: 検証用トンネル経由（本番経路の性能を測る）**

```text
負荷生成 (k6, 外部) → Cloudflare edge → 検証用 tunnel → api:1234 → db:5432 (隔離)
```

- 検証用の一時トンネル（seeft-stg 相当。本番トンネルとは別の tunnel credential を発行）を用意できることは確認済み
- Track A と同一のスタック・同一のデータ・同一のシナリオで実施し、差分をトンネル経路の寄与として切り分ける
- Bingo の教訓（トンネル経由で初回取得・大量同時接続が崩れる / 負荷生成元の回線品質が結果を汚す）を踏まえ、負荷生成元は Cloudflare への経路が安定した外部 VPS を優先する
- 試験終了後は tunnel credential をローテーションする（Bingo 報告書 11 章と同じ後始末）

**やらないこと**

- 本番 DB（共有 Patroni クラスタ）へ向けた試験。理由は 2.4 節
- 本番トンネル・本番ホスト名（seeft-api.nutfes.net 等）へ向けた試験
- mobile 静的配信の初回ロード試験は本計画のスコープ外の別トラック（2.4 節末尾）

### 4.3 試験データ

隔離 DB に 45th 想定規模のデータを投入する生成スクリプト（SQL または Go）を用意する。

| テーブル | 現行 seed | 試験規模（案） | 根拠 |
|---|---|---|---|
| users | 3 | 400 | 実行委員 300 人以上（PM 確認済み）+ 余裕 |
| tasks | 6 | 100〜150 | 45th の実タスク数をスプシから確認して合わせる |
| shifts | 0 | 2万〜5万 | 400人 × 2日 × 実働スロット数。GAS 同期後の実数をスプシから確認 |
| question/shorthanded/trouble_rescues | 0 | 各 100〜200 | GET /rescues の N+1 を現実条件で踏むため必須 |
| session | 0 | テストユーザー分 | 認証3エンドポイント試験用 |

パスワードは全テストユーザーで固定値の bcrypt ハッシュを使い回し、生成を高速化する。**着手前に 45th の実データ規模（名簿人数・シフト行数・タスク数）をスプレッドシートで確認し、表の数値を確定させる**こと。

### 4.4 VU 段階と pass/fail 基準

実行委員 300 人以上という実利用規模に対し、安全係数を掛けた 600 VU までを段階昇圧で確認する。Bingo の段階試験（600 VU 安定 / 650 VU 以上でレイテンシ超過）と同じ刻みを含めることで、構成間の比較も可能にする。

```text
smoke 10 VU（シナリオ動作確認） → 50 → 100 → 200 → 300 → 450 → 600 VU
各段階: 目標 VU 到達後 5 分保持。FAIL した段階で打ち切り、直前の PASS 段階を安定ラインとする
```

pass/fail 基準（段階ごとに全条件を満たして PASS）:

| 項目 | 基準 |
|---|---|
| HTTP 失敗率（接続エラー・タイムアウト含む） | 1% 未満 |
| 5xx 応答 | 0 件（panic・接続枯渇の検出を兼ねる） |
| `GET /shift-cards` p95 / p99 | 500 ms / 1,000 ms 未満 |
| `GET /rescues` p95 | 500 ms 未満 |
| その他読み取り p95 | 300 ms 未満 |
| 書き込み（POST /rescues 含む、スタブ応答込み）p95 | 1,000 ms 未満 |
| 試験直後の `GET /` | 正常応答（Bingo の health/ready 確認に相当） |

計測・記録項目（各段階で保存）:

- k6 のサマリ（エンドポイントタグ別 p50/p95/p99、req/s、失敗内訳）
- `docker stats` の各コンテナ CPU / メモリ推移
- 隔離 DB の `pg_stat_activity` 接続数の最大値（プール無制限問題の実測。修正後の効果測定にも使う）
- API コンテナの stdout ログ量（無条件 Printf の影響確認）

### 4.5 認証付き3エンドポイントのトークン運用

1. 試験データ投入時にテストユーザー N 人を作成（bcrypt ハッシュ固定）
2. 試験セットアップ（k6 の `setup()`）で各テストユーザーにつき **1回だけ** `POST /mail_auth/web_signin` を実行し、返ったトークンを配列で全 VU に共有
3. VU はトークンを読み取り専用で使い回す。**試験中に同一ユーザーで再サインインしない**（`mail_auth_usecase.go:158` の既存セッション削除で旧トークンが無効化されるため）
4. `Access-Token` ヘッダは必ず付与する。付け忘れは 401 ではなく panic 由来の 500 になり、サーバー実装の 5xx と区別できなくなる（修正前に試験する場合の注意点）

認証3エンドポイントは admin(凍結) 専用のため優先度は低いが、panic 修正の回帰確認とトークン検証クエリ（session lookup）の性能確認として、S2 に低頻度（1% 程度）で `GET /current_user` を混ぜる。

### 4.6 外部副作用の遮断

試験用スタックでは次の環境変数差し替えを必須とする。

| 変数 | 試験時の値 | 遮断される副作用 |
|---|---|---|
| `RESCUE_GAS_URL` | ローカルスタブ（`https` 必須。専用の信頼済み CA から発行した証明書 + 200 固定応答の軽量サーバー） | POST /rescues の実スプシ書き込み |
| `SLACK_BOT_TOKEN` | 未設定 | 通知スケジューラごと無効化（`di.go:102-114` で安全にスキップされる） |
| `NUTMEG_DB_*` | 隔離 DB | 共有クラスタへの接続 |

注意点が3つ。第一に、`RESCUE_GAS_URL` は実装が `https` スキームを検証する（`rescue_unified_usecase.go:293-295`）ため、スタブも https で立てる必要がある。第二に、送信処理は `&http.Client{}` の既定 TLS 検証をそのまま使う（`rescue_unified_usecase.go:310`）ため、自己署名証明書をそのまま使うと `x509: certificate signed by unknown authority` で失敗する。スタブ証明書を発行した CA を試験用 API コンテナの信頼ストアに追加する（ローカル CA を発行し `update-ca-certificates` を通す等）ことで解決し、`InsecureSkipVerify` は使わない。第三に、`POST /request_shifts` はハードコード URL（`shift_usecase.go:885`）のため環境変数では遮断できないが、呼び出し元が存在しないため試験対象から除外すれば実害はない。

### 4.7 既知の制約と扱い

- `api/` の自動テストは shift_usecase の純関数テストのみ（テストロードマップのフェーズ1段階）。負荷試験は自動テストの代替にはならず、機能の正しさはスコープ外とする
- Go は 1.26.0 に固定済み（コミット `06935be`）。k6 は API のビルドチェーンと独立なので影響なし
- 負荷試験ツールの運用経験がチームにない。最初の smoke（10 VU）を「ツールの学習」を兼ねた独立ステップとして扱い、いきなり段階試験に入らない

---

## 5. リスクと対応

| リスク | 影響 | 対応 |
|---|---|---|
| DB 接続プール無制限のまま試験すると、高 VU で接続数が数百に達し DB 側 `max_connections` 超過のエラーが噴く | 「API の限界」ではなく「設定不備の限界」を測ってしまう | 試験前に `SetMaxOpenConns` 等の設定を入れる（7章 issue 案）。あえて未設定のまま1回測り、修正効果を対比する選択肢もある |
| Access-Token panic により、ヘッダなしリクエストが 500 になる | 5xx=0 基準が実装バグで汚れる | 試験前に修正（7章）。修正までは試験スクリプト側で必ずヘッダを付与 |
| POST /rescues の同期 GAS 送信（タイムアウトなし） | スタブが遅いと goroutine 滞留で全体が劣化 | スタブは即時応答を既定とし、遅延注入は独立した変種試験として実施 |
| 無条件 `fmt.Printf` + Logger の stdout 出力 | 高 QPS でログ I/O が測定を歪める | ログ量を計測項目に含め、影響が見えたらガード修正後に再測 |
| mobile 静的配信がシングルスレッド | 当日は API より先に初回ロードが詰まる | 本計画のスコープ外だが、別トラックの issue として起票（7章）。当日運用でも「集合時刻より前に開いておく」案内で分散 |
| 検証用トンネルの負荷生成元の回線品質 | Bingo で観測された「負荷生成側の経路不良による偽 FAIL」 | 外部 VPS を負荷生成元に使い、FAIL 時は生成元を変えて再現確認 |
| 8月の引き継ぎ期限 | 試験〜修正〜再試験のループが1巡しかできない可能性 | Track A を先行し、修正が要るものを早期に issue 化。Track B は Track A の合格後に1回で決める |

---

## 6. 依頼時前提の訂正

調査開始時に与えられた前提のうち、以下は現状と異なっていたため訂正して記録する。

- `api/go.mod` が Go 1.16 のまま → 2026-07-09 のコミット `06935be`（issue #385 対応）で `go 1.26.0` に固定済み。`AGENTS.md:7` の「Go 1.16」表記が未更新で残っている
- `api/` 配下に自動テストが 0 件 → `api/lib/usecase/shift_usecase_test.go` が存在（テストロードマップ フェーズ1の成果）
- 負荷試験ツール導入実績なし → 変わらず（確認済み）

---

## 7. 次に立てるべき issue の提案

実行は issue 化してから着手する運用に従い、次の分割を提案する。上から価値順。

**issue 1: 負荷試験の実行環境整備**

k6 の導入手順（バイナリ配置と Makefile ターゲット）、試験用 compose（隔離 DB + GAS スタブ + 環境変数差し替え）、45th 規模の試験データ生成スクリプト。完了条件は「smoke 10 VU が S1/S2 シナリオで green になること」。着手前に 45th 実データ規模（名簿・シフト行数）のスプシ確認を含める。

**issue 2: 試験前修正（測定品質に直結する2件）**

DB 接続プール上限の設定（`db.go` に `SetMaxOpenConns` / `SetMaxIdleConns` / `SetConnMaxLifetime`）と、`Access-Token` ヘッダ欠落時 panic の修正（3箇所、`c.Request().Header.Get("Access-Token")` + 空チェックで 401 応答へ）。どちらも小さく独立した PR にできる。panic 修正はテストロードマップのフェーズ5（controller テスト）の最初の題材と兼ねられる。

**issue 3: Track A（内部直結）試験の実行**

本書 4.4 の段階昇圧を実施し、Bingo 報告書と同形式で結果を記録する。FAIL 時のボトルネック特定（`pg_stat_activity`、pprof の導入検討）まで含む。

**issue 4: Track B（検証用トンネル経由）試験の実行**

検証用 tunnel credential の発行、外部 VPS からの負荷生成、Track A との差分分析、終了後のトークンローテーション。

**issue 5（測ってから判断）: ホットスポット改修**

`GET /shift-cards` の多段クエリ解消と `GET /rescues` の N+1 解消（既存 #264 / #247 の N+1 issue と関連）。**Track A の結果が基準を満たすなら着手しない**（保守モードの方針に従い、測定で必要性が証明されたときだけ直す）。同様に、`shift_repository.go` の無条件 `fmt.Printf` の `DEBUG_SQL` ガード化、`POST /rescues` の GAS 送信タイムアウト設定もここに含める。

**issue 6（別トラック）: mobile 静的配信の初回ロード対策**

`server.py` の `ThreadingHTTPServer` 化（数行の変更）または nginx 等への置き換え検討と、初回ロード（静的ファイル + トンネル）の負荷試験。API 試験とは独立に進められる。

**issue 7（任意・整理）: デッドエンドポイントの閉塞**

`POST /request_shifts`（ハードコード GAS URL・宛先 doPost 不在）、救援個別系の未使用 17 ルート、`GET /reviews` 系ほか。攻撃面の縮小と棚卸しの固定化が目的で、負荷試験の前提条件ではない。
