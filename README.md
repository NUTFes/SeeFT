# SeeFT

## Installation
``` fish
make build
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

## Author
NUTMEG（技大祭実行委員会情報局）
mail: nutfes.info [at] gmail
