# SeeFT

新しく入った人は、まず [技術オンボーディング](docs/development/onboarding.md) を読んでください。

## Installation
``` fish
make build
make up

make migrate
make seed
```

## APIの起動
``` fish
make up-api
```

## mobileの起動
初回だけ、設定ファイルの雛形をコピーし、`fvm` で Flutter を入れる。
``` fish
cp mobile/env/.env.example mobile/env/.env
cd mobile && fvm install && cd ..
```

``` fish
make mobile-up
```
http://localhost:45029 で開く（45029 は `make mobile-up` の既定ポート）。API は CORS で許可するオリジンを列挙しており（`api/lib/externals/server/server.go`）、`127.0.0.1` で開いたり、列挙に無いポートに変えたりすると、API の応答が画面に届かない。

## データベースの削除
``` fish
docker compose down -v
```

## Note
### 初期テーブルの追加
``` fish
make migrate
```

### 初期データの追加
``` fish
make seed
```

### データベースのみ起動
``` fish
make up-db
```

### マイグレーションのdown
``` fish
make migrate-down N=<整数|all>
```

### diを編集してからうまく動かないとき
一度コンテナをdownさせてからupし直してみてください。

## SchemaSpyでDBスキーマを確認する（PostgreSQL）
- DB初期データは `postgresql/db` ディレクトリにありますが、実際のDBはPostgreSQLです。
- 最終生成物の出力先: `api/docs/er-diagrams/`
- 一時出力先: `api/docs/schemaspy`（処理後に削除）

```bash
# 標準（docker-compose.yml）
make schemaspy

# Mac用composeを使う場合
make mac-schemaspy
```
接続先・認証情報は compose の環境変数（`SCHEMASPY_HOST`, `SCHEMASPY_DB` など）で上書きできます。

## Author
NUTMEG（技大祭実行委員会情報局）
mail: nutfes.info [at] gmail
