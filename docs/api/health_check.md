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
id2(Controller) -- Text --> id1(Router)
```

# 概要
## Router
### 責務
以下のAPIを叩くことを許可する

[GET] `/`

### 内部仕様
`e.GET("/", ...)` で controller に渡す

### 実装箇所
`api/lib/router/router.go`

## Controller
### 責務
レスポンスをテキストで返すことができる

### 内部仕様
`IndexHealthcheck(c echo.Context) error`
Status200で以下のテキストを返す（JSON ではない）

```text
healthcheck: ok
```

### 実装箇所
`api/lib/internals/controller/health_controller.go`

# テスト
仕様通りのResponseが返ってくることをテストする

# メモ

