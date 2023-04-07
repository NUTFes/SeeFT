# SeeFT

## Build
``` fish
make build
```

## 初期データの追加
``` fish
make seed
```

## APIの起動
``` fish
make up-db
```

## データベースのみ起動
``` fish
make up-db
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
### Develop
git submodule update --init

### diを編集してからうまく動かないとき
一度コンテナをdownさせてからupし直してみてください。

## Author
NUTMEG（技大祭実行委員会情報局）
mail: nutfes.info [at] gmail
