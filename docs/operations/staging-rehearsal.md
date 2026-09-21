# 検証環境とリハーサル

本番に入れる前に、スプレッドシート → API → アプリの流れを通しで確かめるための手順である。45th（2026年）は、本番デプロイの前日（8/26）に検証環境で名簿・タスク・シフトの送信とログインを通しで試した。このリハーサルで、**全員の学籍番号が0で送られていて、誰もログインできない**という不具合が見つかった（[事故の記録](incidents-45th.md)）。

リハーサルは、説明会（1日目の約3週間前）の数日前までに1回は行う。

## 構成

45th の検証環境は、Proxmox 上の検証用コンテナ（場所は別紙）に置いた。本番とは次の点が違う。

| | 検証環境 | 本番 |
| --- | --- | --- |
| compose | `docker-compose.yml`（開発用） | `docker-compose.prod.yml` |
| API のコード | ホストのソースをマウントし、`air` が変更を検知して作り直す | イメージに焼き込み、`go run` で起動 |
| DB | compose の中の PostgreSQL（`db` サービス） | 外部の共用 HA クラスタ（接続プール経由、SSL 必須） |
| 外からの入り口 | cloudflared のクイックトンネル（起動のたびに URL が変わる） | Cloudflare の名前付きトンネル（`seeft-api.nutfes.net` など） |
| アプリ | 検証環境には立てず、手元の Mac から起動する | コンテナで配信 |

このため、検証環境では見つからない問題がある。ビルドの漏れ、DB への接続（SSL・接続プール・ポート）、本番の環境変数の誤りは、本番でしか確かめられない。本番のデプロイ後の確認は [本番へのデプロイ](deploy.md) にある。

## 手順

### 1. 検証環境を立てる

検証用のサーバーで develop を取得し、DB と API を起動して、スキーマと初期データを入れる。

```bash
docker compose up -d db api
```

```bash
make migrate
```

```bash
make seed
```

API の環境変数は `api/env/dev.env` から読む。このファイルは git の管理下にあるので、検証のために入れた値（初期パスワード、トークンなど）をコミットしない。

開発用の DB は、弱いパスワードのまま 5432 番を外に公開する設定になっている。共有のネットワークに置きっぱなしにせず、検証が終わったら止める。

### 2. GAS から届く URL を作る

スプレッドシートの GAS は Google のサーバーから API を呼ぶので、検証環境に外から届く URL が要る。45th は cloudflared のクイックトンネルを使った。たとえば次のように起動する。

```bash
cloudflared tunnel --url http://localhost:1234
```

表示された `https://〜.trycloudflare.com` が一時的な URL になる。トンネルを起動し直すと URL が変わる。

### 3. シフトスプシの送信先を切り替える

シフトスプシのスクリプトプロパティ `API_BASE_URL` を、トンネルの URL に書き換える（末尾にスラッシュを付けない）。

**リハーサルが終わったら、必ず本番の URL に戻す。** 戻し忘れると、本番のつもりで送った全員分のデータが検証環境に届く。45th では、本番デプロイ後の最初の作業をこの切り替えにした。

### 4. 名簿 → タスク → シフトの順に送る

SeeFT メニューの送信前の確認（[SeeFT に渡すデータの約束事](seeft-data-contract.md) の「送る前に確認するメニュー」）を先に回してから送る。

送ったら、検証環境の DB で件数を確かめる。

```sql
SELECT count(*), count(DISTINCT student_number) FROM users;
```

学籍番号の数が人数と同じであること（初期データのユーザーも数に入る）。全員が0なら、名簿の送り方に問題がある。

```sql
SELECT date_id, weather_id, count(*) FROM shifts GROUP BY 1, 2 ORDER BY 1, 2;
```

日程ごとの件数が「人数 × 64コマ」になっていること。

```sql
SELECT t.task, count(*) FROM shifts s JOIN tasks t ON t.id = s.task_id WHERE s.date_id = 2 AND t.task NOT IN ('', 'NG') GROUP BY t.task ORDER BY 2 DESC LIMIT 10;
```

タスクごとの件数が、シートの「【調査】このシートのタスク別セル数を数える」と一致すること。総数だけを比べると、参加不可（黒塗り）の数え方の違いでずれるので、タスクごとに比べる。

```sql
SELECT task, year_id, count(*) FROM tasks GROUP BY task, year_id HAVING count(*) > 1;
```

同じ名前のタスクが2行以上ないこと（あれば、全角スペースか送る順番の問題）。

### 5. 手元の Mac からアプリを起動する

アプリは検証環境に立てない。API が通信を許可している相手（CORS）が `http://localhost:45029` なので、手元でこのポートで起動する必要がある。

```bash
cd mobile && fvm flutter run -d chrome --web-port 45029 --dart-define-from-file=env/.env --dart-define=API_BASE_URL=<トンネルのURL>
```

- `--web-port 45029` は省略できない。省略するとポートが毎回変わり、すべての通信が CORS で拒否される。
- ホスト名は `localhost` のままにする。`127.0.0.1` では許可されない。
- `API_BASE_URL` の末尾にスラッシュを付けない。付けると、ログイン画面に「学籍番号もしくはパスワードが違います」と出て、認証の問題に見える。
- コマンドラインの `--dart-define` は、`env/.env` の同じ名前の値を上書きする。日付などは `env/.env` のまま、API だけを差し替えられる。
- `-d chrome` は一時的なプロファイルで起動するので、以前のシフトのキャッシュ（Hive）に惑わされない。
- `mobile` の `flutter` は必ず `fvm` 経由で実行する。

### 6. ログインして確かめる

名簿に入っている人の学籍番号と、初期パスワード（`USER_DEFAULT_PASSWORD`）でログインし、次を確かめる。

- シフトカードが出る。集合場所とタスクが正しい
- マニュアルのボタン（ドキュメント版・解説 HTML）が出て、開ける
- 「操作説明」が今年の資料を開く

レスキューを検証環境で試すときは、転記先（`RESCUE_GAS_URL`）を本番のレスキュースプシに向けない。45th はリハーサルでは試さず、本番へ入れたあとに1件だけ送って往復を確かめ、その行を DB とスプシの両方から消した。

### 7. Slack 通知の順番を試す

`SLACK_BOT_TOKEN` を入れる順番を間違えると、溜まったシフト変更が全員に DM で飛ぶ（[SeeFT に渡すデータの約束事](seeft-data-contract.md) の「Slack 通知を有効にする順番」）。検証環境で先に試し、溜まった記録が「送らずに既読」になることを確かめておく。

```sql
SELECT is_sent, count(*) FROM action_logs GROUP BY 1;
```

45th のリハーサルでは、溜まっていた90,368件がすべて既読になることを確かめた。

## 45th のリハーサルで確かめたこと

| 項目 | 結果 |
| --- | --- |
| 名簿 | 356件。重複0、学籍番号は356種類 |
| シフト | 1日程あたり22,592行（353人 × 64コマ）。同じ人・同じ時刻の重複0。タスクごとの件数がシートと一致 |
| ログイン | 実在の委員のアカウントで成功 |
| アプリ | ログインからシフトカードの表示まで |
| マニュアル | ドキュメント版と解説 HTML のボタンが出る |
| 操作説明 | 今年の資料が開く |

## 罠

- **開発用と本番用の compose が、同じイメージ名 `seeft-api` を使う。** 同じサーバーで本番用を build すると開発用のイメージが上書きされ、開発用の API が `air: command not found` で起動しなくなる。検証環境では開発用の compose だけを使う。
- **API の通信の許可（CORS）は、`curl` で 200 が返っても保証にならない。** ブラウザが送る事前の問い合わせ（プリフライト）が通るかを確かめる。

```bash
curl -s -i -X OPTIONS -H "Origin: http://localhost:45029" -H "Access-Control-Request-Method: POST" -H "Access-Control-Request-Headers: content-type,access-control-allow-origin" "<APIのURL>/mail_auth/signin" | grep -i "^HTTP/\|access-control-"
```

- **Google のファイルが委員に見えるかは、匿名の `curl` で目安が付く。** `404` は存在しない、`401` は存在するが匿名では開けない（技大祭アカウントなら開ける場合がある）、`200` はリンクを知っている全員に公開、である。`401` でも委員が開けないとは限らないので、最後は実際の委員に開いてもらう。
- **`mobile/env/.env` の日付が過去のままだと、レビューを求める画面がシフトカードを覆う。** アプリはシフトの終了時刻を env の日付から計算し、過去なら画面を出す。
