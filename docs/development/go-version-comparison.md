# Goバージョン比較（issue #385）

issue #385「Goのバージョンを固定する」は、3か所でばらばらになっているGoの版の指定を1つの版に揃える作業である。3か所とは、`api/go.mod`の`go 1.16`（下限を宣言しているだけ）、`api/Dockerfile` / `api/prod.Dockerfile`の`golang:latest`（版を固定しない）、CI（`go-lint.yml` / `go-test.yml`）の`go-version: stable`（版を固定しない）である。この文書は、`docs/development/test-roadmap.md`フェーズ0の一環として、採用する版を決めるための比較情報をまとめたもの。issue本文の完了条件（採用版が1つに定まる／golangci-lintの指摘が増えない／本番ビルドが固定タグで通る）を満たすことをゴールに、実際にリポジトリを触って確かめた結果を書く。

調査日は2026-07-09。バージョンの実体（Go公式のリリース履歴、依存ライブラリの要求バージョンなど）は調査時点のもので、後から状況が変わりうる。

## 要約

現行コードが要求する言語機能の下限はGo1.13相当である。ただし、依存ライブラリ（特に`jackc/pgx/v5`が要求する`go 1.19`）が実質的な下限を引き上げている。さらにGo公式のセキュリティサポートの方針（直近2メジャーのみ）を踏まえると、実際に選べる候補は**Go 1.25（1つ前のメジャー）かGo 1.26（最新メジャー）の2択**に絞られる。ローカルで実際に試すと、go.modのgoディレクティブをこの範囲内のどの値にしても、`golangci-lint run ./...`の指摘件数は0件のまま変わらなかった。セキュリティサポートが残る期間の差が大きいため、結論として**Go 1.26を推奨**する（詳細は「推奨」節）。

## 現状の3か所の食い違い

| ファイル | 現在の指定 | 問題点 |
|---|---|---|
| `api/go.mod:3` | `go 1.16` | 下限を宣言しているだけで、実際に使う版と離れている。toolchain行なし |
| `api/Dockerfile:1` / `api/prod.Dockerfile:1` | `FROM golang:latest` | 版を固定していない。ビルドのたびに中身が変わりうる |
| `.github/workflows/go-lint.yml:20` / `go-test.yml:20` | `go-version: stable` | 版を固定していない。CIランナーが引く「最新安定版」に依存 |

## 比較するバージョンの選び方

比較するのは、判断に役立つバージョンだけに絞る。

| 分類 | バージョン | 位置づけ |
|---|---|---|
| 現状（比較の起点） | 1.16 | go.modの現在値 |
| 依存の絶対下限 | 1.19 | `jackc/pgx/v5 v5.3.1`（gormのpostgresドライバ経由の間接依存）が要求する実質的な下限。これより下は選べない |
| 仕様の境界点 | 1.21 | go directiveの意味が「努力目標」から「ハード下限＋toolchain自動切替」に変わった版 |
| 実採用候補 | 1.25（N-1） | 2025-08-12リリース |
| 実採用候補 | 1.26（N、最新） | 2026-02-10リリース。ローカル環境（go1.26.4）もこの系列 |

1.16・1.19・1.21は「なぜ他を選ばないか」を示すために並べた版で、実際の採用候補は1.25と1.26の2つである。Go公式のリリース履歴ページ（go.dev/doc/devel/release）によれば、サポートの方針は次のとおり。

> Each major Go release is supported until there are two newer major releases.

つまり、あるメジャー版は、次の次のメジャー版が出た時点でサポート対象外になる。2026-02-10に1.26が出た時点で、1.24以前は既にサポート対象外（EOL）である。今サポート対象なのは1.25と1.26だけ。

## 比較表

| 観点 | 1.16（現状） | 1.19（依存下限） | 1.21（仕様境界） | 1.25（N-1） | 1.26（N） |
|---|---|---|---|---|---|
| go directiveの意味 | 努力目標 | 努力目標 | ハード下限＋toolchain自動切替 | ハード下限＋toolchain自動切替 | ハード下限＋toolchain自動切替 |
| `go build ./...`（`go mod tidy`後） | 成功（現状） | 成功 | 成功 | 成功 | 成功 |
| testify要求（v1.11.1、go1.17） | 満たさない | 満たす | 満たす | 満たす | 満たす |
| go-sqlmock要求（v1.5.2、go1.15） | 満たす | 満たす | 満たす | 満たす | 満たす |
| golangci-lint v2.12.2の実行結果（実機） | 0 issues | 0 issues | 0 issues | 0 issues | 0 issues |
| セキュリティサポート（2026-07-09時点） | EOL | EOL | EOL | サポート中 | サポート中（最新） |
| `golang:X.Y`タグの存在 | 未確認（採らない版のため確かめていない） | 未確認（同上） | 未確認（同上） | `1.25-bookworm` / `1.25-trixie`があることを確認 | `1.26-bookworm` / `1.26-trixie`があることを確認 |

## golangci-lintの指摘件数への影響（実機で試した結果）

test-roadmap.mdに書かれていた懸念（「新しめのgolangci-lintはgo.modの版を参照して指摘の有無を変えることがある」）を、実際に`api/go.mod`の`go`行を書き換えて確かめた。

作業ブランチ上で`go`行を1.19 / 1.21 / 1.25 / 1.26に書き換え、そのたびに`go mod tidy`→`go build ./...`→`golangci-lint run ./...`を実行し、最後に`git checkout -- go.mod go.sum`で元に戻す。これを版ごとに繰り返した（コミットはしていない）。

結果、**1.16（現状）を含めどの版でも`golangci-lint run ./...`は「0 issues.」で一致**した。`api/`配下のコードにはgoroutine・channel・generics・range-over-int/funcが1つもなく（別の静的な調査で確認した）、Go1.22でループ変数が反復ごとに作られるようになった変更のような、版による挙動の違いが影響する余地がそもそもないためと考えられる。つまりこのリポジトリでは、goディレクティブを上げても新しい指摘は出ない。

副作用として1点確かめられたのは、`go mod tidy`の実行結果が変わることである。Go1.16から1.17以降へgoディレクティブを変えると、Goのモジュールグラフ剪定（module graph pruning、Go1.17で導入）の対象になる。そのため、`go.mod`に間接依存が約33行明示的に追加され、`go.sum`から不要なチェックサムが約27行削られる。これは一度きりの機械的な変更で、実装PRで`go mod tidy`を実行すれば自動的に反映される。

```diff
# go.mod の変化（例: go 1.19 にした場合の go mod tidy 結果）
-github.com/mattn/go-isatty v0.0.17 // indirect  （requireブロック内の位置が変わる）
+require (
+	github.com/KyleBanks/depth v1.2.1 // indirect
+	... 約30行の間接依存が明示化される
+)
```

## Dockerタグの固定方法

`docker manifest inspect golang:latest`と`docker run --rm golang:latest cat /etc/os-release` / `go version`で確かめたところ、現時点の`golang:latest`の実体はgo1.26.5、Debian 13 (trixie)だった。つまり、今`golang:latest`のまま使っているのは、実質的にGo1.26系である。

候補タグの`golang:1.26-bookworm`（Debian 12）、`golang:1.26-trixie`（Debian 13）、`golang:1.25-bookworm`、`golang:1.25-trixie`は、両系列ともあることを確かめた。`Dockerfile`が`apt-get install locales`というDebian系のコマンドに依存しているため、alpine系のタグには変えられないという制約が既にある（bookworm/trixieはどちらもDebian系なので問題ない）。

bookwormとtrixieのどちらでもビルド自体は通ると見られる。ただ、trixieは2025年にDebian 13として安定版になったばかりで、bookworm（Debian 12、2023年から安定版）の方がパッケージリポジトリの実績が長い。安定性を優先するなら`-bookworm`系を推奨する。

## CIでのバージョン指定方法

`actions/setup-go`は`go-version`（値を直接指定）と`go-version-file`（`go.mod`などのファイルから読み取る）の両方に対応しており、両方を指定すると`go-version`が優先される。`go-version-file`は、`go.mod`に`toolchain`行があればそれを、なければ`go`行を読む。

`go-version-file: api/go.mod`を使うと、バージョンを決める場所をgo.modの1か所にできる（バージョンを上げたいときはgo.modを直すだけで、CIもそれに合わせて変わる）。ただし、Go1.21以降の自動toolchain切替（GOTOOLCHAIN）が有効なままだと注意が必要である。go.modの`toolchain`行や依存関係の要求によっては、CIが意図しない新しいバージョンを自動でダウンロードすることがある。`actions/setup-go`のメンテナ自身が、go-version-fileでバージョンを明示指定する場合は`GOTOOLCHAIN=local`を設定し、CIマトリクスのバージョンを厳密に固定することを推奨している。

## 推奨

Go 1.26を推奨する。

根拠は主に、セキュリティサポートが残る期間である。Go公式のリリース間隔は例年2月・8月の半年おきで、直近の実績は1.24（2025-02）→1.25（2025-08-12）→1.26（2026-02-10）である。この間隔のまま進むと、次のメジャー（1.27）は2026年8月ごろになる可能性が高い（確定ではなく、公式のアナウンスも確かめていない）。

もしそうなると、今1.25（N-1）を採用した場合は、1ヶ月ほど先に1.27がリリースされた時点で「2つ新しいメジャーが出た」状態にさしかかり、サポートが残る期間が非常に短い。一方、今1.26（N）を採用した場合は、1.28がリリースされるまでサポートが続くため、残る期間は1年以上になる。

さらに、前節で実際に試した結果、golangci-lintの指摘件数は1.25でも1.26でも変わらなかった。そのため、「枯れているから1.25の方が安全」という理由づけはこのリポジトリには当てはまらない。ローカルの開発環境も既にgo1.26.4で、Docker側も版を固定しない`golang:latest`が既に事実上1.26系を指している。1.26への固定は、「新しい版に上げる」というより「今すでに使っている版を固定する」に近い。

具体的には、次のように変える。

- `api/go.mod`: `go 1.26.0`（`toolchain`行は追加しない。開発者のローカルGoが1.26未満の場合、Go1.21以降の自動toolchain切替でビルド時に自動取得される）
- `api/Dockerfile` / `api/prod.Dockerfile`: `FROM golang:1.26-bookworm`
- `.github/workflows/go-lint.yml` / `go-test.yml`: `go-version-file: api/go.mod`に変更し、`GOTOOLCHAIN=local`を環境変数として設定
- 実装PR内で`cd api && go mod tidy`を実行し、モジュールグラフ剪定によるgo.mod/go.sumの変化をそのままコミットする

## 未確認事項

- Go 1.27の正式リリース日は確かめていない（公式のアナウンスは探しておらず、過去の間隔からの推測にとどまる）。実装PRのタイミングによってはN/N-1の値が変わりうるため、着手時にgo.dev/doc/devel/releaseを見直すこと
- `golang:1.16`系や`1.19`系などの古いDockerタグが今もあるかは確かめていない（採らない候補のため）
- CI（GitHub Actions ubuntu-latestランナー）上で`golangci-lint-action@v2.12`が、実際にどのGoバージョンでビルドしたバイナリを配布するかは、ローカルのHomebrew版（`golangci-lint 2.12.2、go1.26.2でビルド`）と完全に一致する保証がない。実装PRのCIログで`golangci-lint version`の出力を見るのが確実
- golangci-lintがgo.modの`go`ディレクティブを、どこまで厳密に解析の対象の版へ反映するか（ドキュメントで「go.modのgoディレクティブを既定値として使用、GOVERSION環境変数、さらに1.17へフォールバック」という記述を見ただけで、実装の詳細までは調べていない）

## 既存issue・関連文書との対応

| 項目 | 関連 |
|---|---|
| このドキュメントの発端 | #385 |
| 親ロードマップ | `docs/development/test-roadmap.md`フェーズ0 |
