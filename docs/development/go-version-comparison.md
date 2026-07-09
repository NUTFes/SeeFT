# Go バージョン比較（issue #385）

issue #385「Go のバージョンを固定する」は、`api/go.mod` の `go 1.16`（下限宣言のみ）、`api/Dockerfile` / `api/prod.Dockerfile` の `golang:latest`（floating）、CI（`go-lint.yml` / `go-test.yml`）の `go-version: stable`（floating）という三重不整合を1つの版に揃える作業である。この文書は、`docs/development/test-roadmap.md` フェーズ0の一環として、採用版を決めるための比較情報をまとめたもの。issue本文の完了条件（採用版が1つに定まる／golangci-lintの指摘が増えない／本番ビルドが固定タグで通る）を満たすことをゴールに、実際にリポジトリを触って確認した結果を記載する。

調査日は2026-07-09。バージョンの実体（Go公式のリリース履歴、依存ライブラリの要求バージョンなど）は調査時点のもので、後から状況が変わりうる。

## 要約

現行コードが要求する言語機能の下限はGo1.13相当だが、依存ライブラリ（特に`jackc/pgx/v5`が要求する`go 1.19`）が実質的な下限を押し上げており、さらにGo公式のセキュリティサポートポリシー（直近2メジャーのみ）を踏まえると、実際に選べる候補は**Go 1.25（1つ前のメジャー）かGo 1.26（最新メジャー）の2択**に絞られる。ローカルでの実機検証では、go.modのgoディレクティブをこの範囲内のどの値にしても`golangci-lint run ./...`の指摘件数は0件のまま変化しなかった。結論として、セキュリティサポートの残存期間の差が大きいため**Go 1.26を推奨**する（詳細は「推奨」節）。

## 現状の三重不整合

| ファイル | 現在の指定 | 問題点 |
|---|---|---|
| `api/go.mod:3` | `go 1.16` | 下限宣言のみで実態と乖離。toolchain行なし |
| `api/Dockerfile:1` / `api/prod.Dockerfile:1` | `FROM golang:latest` | floating。ビルドのたびに中身が変わりうる |
| `.github/workflows/go-lint.yml:20` / `go-test.yml:20` | `go-version: stable` | floating。CIランナーが引く「最新安定版」に依存 |

## 比較対象バージョンの選定基準

全バージョンを並べるのではなく、判断に効くものだけに絞る。

| 分類 | バージョン | 位置づけ |
|---|---|---|
| 現状（比較の起点） | 1.16 | go.modの現在値 |
| 依存の絶対下限 | 1.19 | `jackc/pgx/v5 v5.3.1`（gormのpostgresドライバ経由の間接依存）が要求する実質フロア。これより下は選択不可 |
| 仕様の境界点 | 1.21 | go directiveの意味が「努力目標」から「ハード下限＋toolchain自動切替」に変わった版 |
| 実採用候補 | 1.25（N-1） | 2025-08-12リリース |
| 実採用候補 | 1.26（N、最新） | 2026-02-10リリース。ローカル環境（go1.26.4）もこの系列 |

1.16・1.19・1.21は「なぜ他を選ばないかを示す参照列」であり、実際の採用候補は1.25と1.26の2つに絞られる。Go公式のリリース履歴ページ（go.dev/doc/devel/release）によれば、サポートポリシーは次の通り。

> Each major Go release is supported until there are two newer major releases.

つまり、あるメジャー版は次の次のメジャー版が出た時点でサポート対象外になる。2026-02-10に1.26が出た時点で、1.24以前は既にサポート対象外（EOL）。現在サポート対象なのは1.25と1.26のみ。

## 比較表

| 観点 | 1.16（現状） | 1.19（依存下限） | 1.21（仕様境界） | 1.25（N-1） | 1.26（N） |
|---|---|---|---|---|---|
| go directiveの意味論 | 努力目標 | 努力目標 | ハード下限＋toolchain自動切替 | ハード下限＋toolchain自動切替 | ハード下限＋toolchain自動切替 |
| `go build ./...`（`go mod tidy`後） | 成功（現状） | 成功 | 成功 | 成功 | 成功 |
| testify要求（v1.11.1、go1.17） | 満たさない | 満たす | 満たす | 満たす | 満たす |
| go-sqlmock要求（v1.5.2、go1.15） | 満たす | 満たす | 満たす | 満たす | 満たす |
| golangci-lint v2.12.2の実行結果（実機） | 0 issues | 0 issues | 0 issues | 0 issues | 0 issues |
| セキュリティサポート状況（2026-07-09時点） | EOL | EOL | EOL | サポート中 | サポート中（最新） |
| `golang:X.Y`タグの存在 | 未確認（不採用のため確認省略） | 未確認（同上） | 未確認（同上） | `1.25-bookworm` / `1.25-trixie` 存在確認済み | `1.26-bookworm` / `1.26-trixie` 存在確認済み |

## golangci-lintの指摘件数への影響（実機検証結果）

test-roadmap.mdに書かれていた懸念（「新しめの golangci-lint は go.mod の版を参照して指摘の有無を変えることがある」）を、実際に`api/go.mod`の`go`行を書き換えて検証した。

手順: 作業ブランチ上で`go`行を1.19 / 1.21 / 1.25 / 1.26に書き換え、都度`go mod tidy`→`go build ./...`→`golangci-lint run ./...`を実行し、最後に`git checkout -- go.mod go.sum`で元に戻す、を繰り返した（コミットはしていない）。

結果、**1.16（現状）を含めどの版でも`golangci-lint run ./...`は「0 issues.」で一致**した。`api/`配下のコードにはgoroutine・channel・generics・range-over-int/funcが一切存在せず（別途の静的調査で確認済み）、Go1.22のループ変数per-iteration化のような版依存の挙動変化が影響する余地がそもそもないためと考えられる。つまりこのリポジトリに関しては、goディレクティブを上げても新規の指摘は発生しない。

一点、副作用として確認できたのは`go mod tidy`の実行結果そのものである。Go1.16→1.17以降のgoディレクティブへの変更は、Goのモジュールグラフ剪定（module graph pruning、Go1.17で導入）の対象になるため、`go.mod`に間接依存が明示的に約33行追加され、`go.sum`から不要なチェックサムが約27行削られる。これは一度きりの機械的な変更であり、実装PRで`go mod tidy`を実行すれば自動的に反映される。

```
# go.mod の変化（例: go 1.19 にした場合の go mod tidy 結果）
-github.com/mattn/go-isatty v0.0.17 // indirect  （requireブロック内の位置が変わる）
+require (
+	github.com/KyleBanks/depth v1.2.1 // indirect
+	... 約30行の間接依存が明示化される
+)
```

## Dockerタグの固定方法

`docker manifest inspect golang:latest` と `docker run --rm golang:latest cat /etc/os-release` / `go version` で確認したところ、現時点の`golang:latest`の実体は **go1.26.5、Debian 13 (trixie)** だった。つまり、今`golang:latest`のまま使っているのは実質Go1.26系である。

候補タグは両系列とも存在を確認済み: `golang:1.26-bookworm`（Debian 12）、`golang:1.26-trixie`（Debian 13）、`golang:1.25-bookworm`、`golang:1.25-trixie`。`Dockerfile`が`apt-get install locales`というDebian系コマンドに依存しているため、alpine系タグには変更できない制約が既にある（bookworm/trixieはどちらもDebian系なので問題ない）。

bookwormとtrixieのどちらでもビルド自体は通ると見られるが、trixieは2025年にDebian 13として安定版になったばかりで、bookworm（Debian 12、2023年から安定版）の方がパッケージリポジトリの実績が長い。安定性を優先するなら`-bookworm`系を推奨する。

## CIでのバージョン指定方法

`actions/setup-go`は`go-version`（固定値指定）と`go-version-file`（`go.mod`等のファイルから読み取り）の両方をサポートしており、両方指定した場合は`go-version`が優先される。`go-version-file`は`go.mod`の`toolchain`行があればそれを、なければ`go`行を読む。

`go-version-file: api/go.mod`を使うと、go.modを唯一の真実源にできる（バージョンを上げたいときgo.modを直すだけでCIも追従する）。ただし1点注意が必要で、Go1.21以降の自動toolchain切替（GOTOOLCHAIN）が有効なままだと、go.modの`toolchain`行や依存関係の要求次第でCIが意図しない新しいバージョンを自動ダウンロードすることがある。`actions/setup-go`のメンテナ自身が、go-version-fileでバージョンを明示指定する場合は`GOTOOLCHAIN=local`を設定してCIマトリクスのバージョンを厳密に固定することを推奨している。

## 推奨

**Go 1.26を推奨する。**

根拠は主にセキュリティサポートの残存期間。Go公式のリリース間隔は例年2月・8月の半年おきで、直近の実績は1.24（2025-02）→1.25（2025-08-12）→1.26（2026-02-10）。この間隔をそのまま延長すると、次のメジャー（1.27）は2026年8月ごろになる可能性が高い（未確定・公式アナウンス未確認）。もしそうなると、

- 今1.25（N-1）を採用した場合: 1.27がリリースされた時点（＝1ヶ月ほど先）で「2つ新しいメジャーが出た」状態に片足がかかり、サポート残存期間が非常に短い
- 今1.26（N）を採用した場合: 1.28がリリースされるまでサポートが続くため、残存期間は1年以上

さらに、実機検証（前節）でgolangci-lintの指摘件数は1.25でも1.26でも変化がなかったため、「枯れているから1.25の方が安全」という理由づけがこのリポジトリには当てはまらない。ローカル開発環境も既にgo1.26.4であり、Docker側もfloatingの`golang:latest`が既に事実上1.26系を指しているため、1.26への固定は「新しい版に上げる」というより「今すでに使っている版を凍結する」に近い。

具体的な変更内容:

- `api/go.mod`: `go 1.26.0`（`toolchain`行は追加しない。開発者のローカルGoが1.26未満の場合、Go1.21以降の自動toolchain切替でビルド時に自動取得される）
- `api/Dockerfile` / `api/prod.Dockerfile`: `FROM golang:1.26-bookworm`
- `.github/workflows/go-lint.yml` / `go-test.yml`: `go-version-file: api/go.mod` に変更し、`GOTOOLCHAIN=local` を環境変数として設定
- 実装PR内で `cd api && go mod tidy` を実行し、モジュールグラフ剪定によるgo.mod/go.sumの変化をそのままコミットする

## 未確認事項

- Go 1.27の正式リリース日は未確認（公式アナウンス未検索、過去の間隔からの推測にとどまる）。実装PRのタイミングによってはN/N-1の値が変わりうるため、着手時に go.dev/doc/devel/release を再確認すること
- `golang:1.16` 系や `1.19` 系などの古いDockerタグが現存するかは確認していない（不採用候補のため）
- CI（GitHub Actions ubuntu-latestランナー）上で`golangci-lint-action@v2.12`が実際にどのGoバージョンでビルドされたバイナリを配布するかは、ローカルのHomebrew版（`golangci-lint 2.12.2、go1.26.2でビルド`）と完全一致する保証がない。実装PRのCIログで`golangci-lint version`の出力を確認するのが確実
- golangci-lintがgo.modの`go`ディレクティブをどこまで厳密に解析ターゲットへ反映するか（ドキュメント上は「go.modのgoディレクティブを既定値として使用、GOVERSION環境変数、さらに1.17へフォールバック」という記述を確認したのみで、実装の詳細までは未調査）

## 既存issue・関連文書との対応

| 項目 | 関連 |
|---|---|
| このドキュメントの発端 | #385 |
| 親ロードマップ | `docs/development/test-roadmap.md` フェーズ0 |
