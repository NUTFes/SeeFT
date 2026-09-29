# SeeFT-API

SeeFT の API サーバー。Go + Echo + PostgreSQL。

起動・マイグレーション・シードはリポジトリのルートから `make` で行う（[README](../README.md)）。構成とコードの書き方は [AGENTS.md](../AGENTS.md)、はじめて触る人向けの説明は [技術オンボーディング](../docs/development/onboarding.md) にある。

## 構成

```text
main.go          起動するだけ。中身は lib/di
lib/di           依存性注入。repository → usecase → controller → router を配線する
lib/router       URL と controller の対応づけ（全エンドポイントの一覧）
lib/internals    controller（HTTP の入出力）と repository（SQL の実行）
lib/usecase      業務ロジック
lib/entity       データの形（JSON のキー名はここで決まる）
lib/externals    DB 接続、HTTP サーバー、Slack、通知の定期実行
cmd/migrate      postgresql/db/schema と migrations を流す
cmd/seed         postgresql/db/seed.sql を流す
cmd/send-notifications  未送信の Slack 通知を手動で流す
```

## 開発

ローカルのコンテナは [air](https://github.com/air-verse/air) で起動しており、`.go` ファイルを保存すると自動で再ビルドされる。ポートは 1234。

```bash
make up-api
make test
```

`make up-api` のログは `docker compose logs -f api` でも見られる。

### diを編集してからうまく動かないとき
一度コンテナをdownさせてからupし直してみてください。

## Author
NUTMEG（技大祭実行委員会情報局）
mail: nutfes.info [at] gmail
