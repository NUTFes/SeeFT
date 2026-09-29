# seeft_mobile

SeeFT の参加者向けアプリ。Flutter で書き、Web としてだけ運用している（Android / iOS 向けの配布はしていない）。

構成とコードの書き方は [AGENTS.md](../AGENTS.md)、はじめて触る人向けの説明は [技術オンボーディング](../docs/development/onboarding.md) にある。

## Requirements
- fvm
- Flutter 3.27.3（`.fvmrc` で固定。`flutter` コマンドは必ず `fvm flutter ...` で打つ）

## 起動

初回だけ、設定ファイルの雛形をコピーし、Flutter を入れる。以下は `mobile/` で実行する。

```bash
cp env/.env.example env/.env
fvm install
```

リポジトリのルートで `make mobile-up` を実行し、http://localhost:45029 で開く。API の CORS が許可しているのはこのポートだけなので、変えないこと。

`env/.env` の値はビルド時に埋め込まれる（`String.fromEnvironment`）。値を変えたら起動し直す。

## テスト

```bash
fvm flutter test --platform chrome
```

`--platform chrome` が必須。付けないと Web 向けのコードを読み込めずに落ちる。`test/widget_test.dart` は Flutter の雛形のまま壊れている（#493）。

## 本番

`Dockerfile` で Web 向けにビルドし、`python/server.py` がビルド済みのファイルを配る（`docker-compose.prod.yml`）。手順は [docs/operations/deploy.md](../docs/operations/deploy.md)。
