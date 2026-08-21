# SeeFT

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
``` fish
cd ./mobile
docker build ./ -t seeft-mobile
docker run --detach --publish 45029:45029 seeft-mobile
```

## データベースの削除
``` fish
docker compose down -v
```

## Note
### 初期テーブルの追加
```
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

### Develop
git submodule update --init

### diを編集してからうまく動かないとき
一度コンテナをdownさせてからupし直してみてください。

## SchemaSpyでDBスキーマを確認する（PostgreSQL）
- DB初期データは `mysql/db` ディレクトリにありますが、実際のDBはPostgreSQLです。
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
