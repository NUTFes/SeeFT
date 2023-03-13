# Title
Health Check

# Writer
@mashita1023

# 目的
サーバーが立ち上がってるか確認のためのヘルスチェック

# 背景
サーバーが立ち上がっていてAPIをとりあえず返せるのかをチェックしたい

# アーキテクチャ

``` mermaid
id1(Router) --> id2(Controller)
id2(Controller) -- Json --> id1(Router)
```

# 概要
## Router
### 責務
以下のAPIを叩くことを許可する

[GET] `/healthz`

### 内部仕様
router.get をする

### 実装箇所
`/internal/external/server/router.dart`

## Controller
### 責務
レスポンスをJsonで返すことができる

### 内部仕様
`Response getHealth(Request request)`
Status200で以下のJSONを返す

``` 
{
  status: "healch ok"
}
```

### 実装箇所
`/internal/interface/controller/health_controller.dart`

# テスト
仕様通りのResponseが返ってくることをテストする

# メモ

