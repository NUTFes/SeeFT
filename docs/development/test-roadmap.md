# SeeFT テストロードマップ

2026-07-06 の MT で「SeeFT 本体（mobile / api / admin）は保守性を高める方向」が決定した。この文書は、その次のアクションとして挙がった「どこからテストを書き始めるべきか」を、実コードの調査結果に基づいて定義する。調査時点で `api/` 配下の `*_test.go` は 0 件、`mobile/test/` はデフォルト残骸 1 件のみ、CI は lint 専用でテスト実行ステップは存在しない。つまりテストはゼロからのスタートである。

## 要約

実行順序は次のとおり。

1. フェーズ0: テスト実行基盤の整備（CI に go test、Go バージョン一本化、壊れた compose マウントの修正）
2. フェーズ1: 依存ゼロの純関数テスト（モックも DB も不要。テストの書き方の規約をここで確立する）
3. フェーズ2: repository 層のゴールデンテスト（実 DB を使う統合テスト。「現状の動作を正解にする」の実装）
4. フェーズ3: repository 層のクエリ修正（フェーズ2 のテストを安全網にしてプレースホルダ化）
5. フェーズ4: usecase 層のテスト（go-sqlmock）
6. フェーズ5: controller 層のテスト（httptest。GAS との「結合テスト」はここの契約テストで代替する）

mobile は api と独立して並行で進められる。GAS と admin は自動テストの対象外とする（理由は後述）。

MT の叩き台（静的解析 → repository 修正 → repository テスト → usecase テスト → handler テスト → GAS テスト → GAS×API 結合テスト）からの変更は 2 点。

- 「静的解析導入」は完了済みのため外した。`.github/workflows/go-lint.yml` で golangci-lint v2.12 が PR ごとに走っており、残っているのは指摘の解消（#314）でありテストとは独立に進められる。
- 「repository 修正 → repository テスト」を逆順にした。テストがない状態でクエリを書き換えると、直したのか壊したのか判定できない。MT で出た「現状の動作を正解にする」アプローチそのものが、修正より先にテストを書く理由になる。

## 前提: コード構造がテスト戦略を規定する

api は素直なレイヤ構成で、全層がコンストラクタ注入とインターフェースで配線されている（`api/lib/di/di.go` に集約）。これはテストの土台として恵まれている。ただし一点だけ大きな制約がある。repository のインターフェースが `*sql.Rows` / `*sql.Row` という `database/sql` の具象型を返し、Scan を usecase 側が行う設計になっている（18 ファイル中 17 ファイル。例外は `shift_card_repository.go` のみ）。

```go
// api/lib/internals/repository/bureau_repository.go
type BureauRepository interface {
	All(context.Context) (*sql.Rows, error)
	Find(context.Context, string) (*sql.Row, error)
	// ...
}
```

`*sql.Row` はテストコードから作れないため、repository のインターフェースをモックしても usecase のテストは書けない。対処は 2 案ある。

- 案A: `db.Client`（`api/lib/externals/db/db.go:19-23`、`DB() *sql.DB` を持つインターフェース）に [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) の偽 `*sql.DB` を差し込む。既存コードの改修は不要で、repository の実装ごとテストを通す形になる。
- 案B: repository の戻り値を entity に変えるリファクタ。テストは書きやすくなるが、141 メソッドに波及する大工事で保守モードの規模を超える。

このロードマップは案A を採用する。以降のフェーズはすべてこの前提で書かれている。

## フェーズ0: テスト実行基盤の整備

テストを 1 本も書く前に、書いたテストが回る環境を作る。すべて小さい独立した PR にできる。

- Go バージョンの一本化（既存 issue #385）。現状は `api/go.mod` が `go 1.16`、`api/Dockerfile` と `api/prod.Dockerfile` が `golang:latest`、`go-lint.yml:20` が `go-version: stable` の三重不整合。go directive が 1.16 のままだと testify や go-sqlmock の新しいバージョンが入らないため、テスト導入の直接の前提になる。採用バージョンの比較・実機検証結果は [go-version-comparison.md](./go-version-comparison.md) にまとめた。
- `docker-compose.yml:8` と `docker-compose.mac.yml` のマウント修正。存在しない `./mysql/db` を initdb にマウントしており（#277 のディレクトリ改名への追随漏れ）、DB の初期化が機能していない。`./postgresql/db` に直す。MySQL 時代の遺物である `my.cnf` のマウントも合わせて整理する。
- CI に go test ジョブを追加。`go-lint.yml` と同型で `working-directory: api`、`go test ./... -count=1`。テストが 0 本でも green になるので、最初に敷いておくと以降のフェーズの PR が全部 CI で検証される。
- `Makefile` に `test` ターゲットを追加し、`AGENTS.md` の Commands 節に記載する。なお `make seed` 系 3 ターゲットは参照先の `/app/seeds/seeds.go` が存在せず全て壊れている（後述のバグ一覧参照）。
- `go.mod` に `stretchr/testify` と `DATA-DOG/go-sqlmock` を追加。あわせて `api/lib/usecase/shift_usecase.go:66` の未使用グローバル変数（`var TaskID, UserID, ... string`）を削除しておく。

## フェーズ1: 依存ゼロの純関数テスト

モックも DB もフェーズ0 の完了すら待たずに書ける関数が既にある。ここでテーブル駆動テストの規約（ファイル配置、命名、テストケースの書き方）を確立し、以降のフェーズの雛形にする。

- `api/lib/usecase/notification_usecase.go` のヘルパー 4 関数: `GroupNotificationsByUserAndDate`（L138）、`sortLogsByTime`（L317）、`formatTimeRange`（L359）、`buildChangesWithTime`（L371）。unexported なレシーバでもフィールドに触れないため、同一パッケージ内テストから `(&notificationUseCase{})` で直接呼べる。`formatTimeRange` には「`endTimeID+1` が timeMap に無いと `0:00` にフォールバックする」という境界条件が既にあり、テストで固定する価値が高い。Slack 通知は壊れると目に見えて困る機能であり、通知文言の回帰防止に直結する。
- `api/lib/usecase/shift_usecase.go` の純関数 3 つ: `groupContinuousShifts`（L478、連続 TimeID のグループ化）、`compareTimeStrings`（L509、ゼロ埋めなし時刻文字列の比較という壊れやすい仕様）、`convertShiftCardDataToShifts`（L433）。
- `api/lib/externals/slack/slack_service.go` の `BuildMessageBlocks`（L88）。Slack API 通信なしで Block 構成の分岐を検証できる。
- `api/lib/externals/scheduler/scheduler.go` の `Start`。カウンタを増やすだけの Job と短い interval で「起動直後の即時実行」「interval ごとの再実行」「ctx キャンセルで停止」「job エラーでもループ継続」を検証できる。このパッケージは PR #322（通知の定期実行）で追加されるため、マージ後に着手する。

## フェーズ2: repository のゴールデンテスト（実 DB 統合）

「現状の動作を正解にする」を実装するフェーズ。ゴールデンテスト（golden master test。レガシーコード文脈では特性化テスト/characterization testとも呼ぶ）とは、正しい仕様を定義するのではなく、いま動いているコードの出力をそのままテストの期待値として固定するもの。これがあると次のフェーズのクエリ書き換えが「テストが緑のまま = 挙動が変わっていない」で機械的に判定できる。

実行環境は GitHub Actions の service container を使う。

- `postgres:18`（`docker-compose.yml` と同じイメージ・資格情報）を service container として起動し、`postgresql/db/*.sql` を番号順（create1 → create6 → seed.sql）に投入する共通セットアップを作る。compose の initdb マウントは壊れているため、テストヘルパー側で SQL を適用する。
- DB が必要なテストには build tag `integration` を付け、通常の `go test ./...` と分離する。
- 接続は環境変数（`NUTMEG_DB_HOST` 等）のみ参照する既存実装のままで済む。CI では `localhost` を指定する。

最初の 3 本は、テストの価値がすぐ実証できるものを選ぶ。

1. `departmentRepository.Create`（`department_repository.go:42-45`）。実カラム名は `department` なのに `departments` へ INSERT しており恒常的に失敗する既存バグがある。1 本目のテストがいきなりバグを検出することで、テスト導入の説得力になる。
2. `reviewRepository.Create / Update`（`review_repository.go:77-96`）。自由記述の comment をクォート連結しているため、`It's fine` のようなアポストロフィ入り文字列で現状は失敗する。このテストはフェーズ3 のプレースホルダ化の完了判定を兼ねる。
3. `shiftRepository` の CRUD 一巡（Create → FindByUnique → Update → Destroy）。文字列連結が 17 行と最多で、シフトというドメインの中核でもある。

## フェーズ3: repository のクエリ修正

MT で「repository のクエリがゴミ」と言われていた部分の実態は、11 ファイル・約 58 行の文字列連結 SQL である。controller は入力検証なしで HTTP パラメータを素通しするため、これは書式の問題ではなく実際に SQL インジェクションが可能な状態を意味する（既存 issue #363 / #266）。

```go
// api/lib/internals/repository/shift_repository.go:47
query := "SELECT * FROM shifts WHERE id = " + id
```

フェーズ2 のテストを安全網として、連結を `$1` プレースホルダに置換していく。同じファイル内に正しい書き方（`shift_repository.go` の `Users` メソッドは `$1〜$5` を使用）が既にあるので、それに揃える。優先順位は、ユーザー自由記述が入る review、次いで連結箇所最多の shift。あわせて次の 2 点も処理する。

- `department_repository.go:42` のカラム名バグ修正。これは挙動の変更（常に失敗 → 成功する）なので、ゴールデンテストの期待値を意識的に更新する。
- 無条件の `fmt.Printf`（クエリのデバッグ出力）を `DEBUG_SQL` 環境変数ガードに統一。テスト出力の汚染防止を兼ねた機械的な整理。

## フェーズ4: usecase のテスト（go-sqlmock）

案A の方式で、`db.Client` を満たす小さなフェイク（`DB()` が sqlmock の `*sql.DB` を返す struct）を 1 つ作り、repository の実装ごと usecase を通す。Scan のカラム順まで再現する必要があるためフェーズ1 より手間はかかるが、フェイクは全 usecase に横展開できる。

- 雛形には最小の `bureau_usecase.go`（71 行）を使う。`GetBureaus`（L24）で複数行の Scan とフィールドマッピング、`GetBureauByID`（L54）で 1 行取得、クエリエラー時のエラー伝播、空リスト時に nil を返す現状挙動（AGENTS.md 規約の `[]Type{}` との差異）を固定する。
- 以降は価値順に進める。`mail_auth_usecase.SignIn`（L33、認証というクリティカルパス。全角スペース除去 → bcrypt 照合のロジック）、`notification_usecase.ProcessUnsentNotifications`（通知の本丸）。
- `notificationUseCase` は具象型 `*slack.SlackService` に直接依存しているため（`notification_usecase.go:21`）、Slack 送信まで含めた end-to-end テストにはインターフェース抽出のリファクタが要る。これは別 issue に切り出し、それまでは sqlmock で届く範囲をテストする。

## フェーズ5: controller のテスト

controller は usecase インターフェースだけに依存する薄い層なので、手書きフェイク + echo の httptest で書ける。薄いぶん網羅は狙わず、価値のある箇所に絞る。

- FromGAS 系 3 ハンドラの契約テスト: `router.go:177-179` の `/api/update_users`、`/api/update_tasks_and_places`、`/api/update_shifts`。GAS が実際に送る JSON をフィクスチャとして与え、パースと usecase への受け渡しを固定する。GAS 側をテストしなくても、api 側の変更が連携を壊したことを検知できる。これが後述する GAS×API「結合テスト」の現実解。
- `user_controller.go:97` の `c.Request().Header["Access-Token"][0]` はヘッダ欠落時に index out of range で panic する疑いがある。テストで最初に検証すべき欠陥候補。
- リクエストの型変換・バリデーション整備（既存 issue #267）に着手する際は、このフェーズのテストを先に書いてから行う。

## 並行トラック: mobile

mobile は api と依存関係がなく、Flutter に慣れたメンバーへ独立して割り振れる。直近 12 ヶ月で 104 コミットと活発なため、テスト投資の回収が見込める。

- 最初の PR: `mobile/test/widget_test.dart`（存在しない `main.dart` の MyApp を参照するデフォルト残骸で、`flutter test` を CI に足すと即 red になる）を削除または書き直し、同じ PR で `flutter-lint.yml` に `fvm flutter test` 相当のステップを追加する。
- `mobile/lib/models/shift_card.dart` の `ShiftCardDataList.fromJson`（L53-105）。依存ゼロの純関数で即着手できる。L78 で `item['before_members']` 自体の null チェックがなくクラッシュする疑いがあるため、正常系と合わせて固定する。最頻変更画面 my_shift_page の入力データを守る一本目として費用対効果が最も高い。
- `mobile/lib/models/rescue.dart` の type ディスパッチ（L33-44）。trouble / question / shorthanded のラウンドトリップと未知 type の挙動。
- New バッジ判定ロジック（`my_shift_page.dart` の `_isCardChanged` L420、`_detectNewOrUpdatedCardKeys` L437 ほか）。壊れても画面上気づきにくいサイレント故障の典型箇所だが、private な State クラス内にあるため、テストするには純関数として `lib/utils/` へ抽出する移動リファクタが先になる。models のテストが定着してから着手する。

## GAS の方針: 自動テストの対象外とする

`gas/` は 4 プロジェクト・約 1,474 行あるが、内容は「スプレッドシートと api の双方向同期グルーコード」で、SpreadsheetApp / UrlFetchApp などの GAS API 依存が全関数に染み込んでいる。user と shift はファイルのトップレベルで `SpreadsheetApp.getUi()` 等を実行するため、Node でロードした瞬間にクラッシュし、テストランナーに載せること自体ができない。ここに自動テストを整備する費用対効果は低いと判断する。

代替と例外は次のとおり。

- 連携の破壊検知はフェーズ5 の FromGAS 契約テスト（api 側）で行う。
- GAS 側を変更する必要が出たときは、変更対象のロジックを引数を取る純関数として別ファイルに抽出してから触る（rescue/onChange.js のステータスマッピングや、shift/コード.js のペイロード構築が候補）。45th 対応で `yearID` 前提の改修が入る可能性が高いのはシフトペイロード構築部分で、ここだけは先行して抽出しておく価値がある。
- 調査で `api/lib/usecase/shift_usecase.go:887` に GAS Web アプリの URL がハードコードされているが、対応する `doPost` が `gas/shift/` に存在しないことが分かった。リポジトリと実デプロイの差分を `clasp pull` で確認し、どちらを正本にするか決める作業が先に必要。

## admin の方針: 対象外

`admin/` は直近 12 ヶ月コミット 0 件の凍結状態で、AGENTS.md にも不使用と明記されている。テスト投資はしない。もし再稼働させる判断をした場合のみ、CI に `next build` のスモークジョブを足すところから始める。

## MT の宿題:「結合テストって最初と最後どっちだっけ」

「結合テスト」が 2 つの意味で使われているので分けて答える。

- api×DB の統合テスト（フェーズ2）は最初側。教科書的な「単体を固めてから結合」はテストピラミッドを新規開発で下から積む場合の話で、テストのないコードを保守する場合は逆になる。単体テストを書くための分解（モック化・リファクタ）自体が挙動を変えるリスクを持つため、先に外側から現状の動作を固定し、その安全網の中で内側を整えるのがレガシーコード改善の定石。
- GAS×API をプロセス横断で通す E2E は最後、というより自動化しない。GAS 側に自動実行環境を組めないため、api 側の契約テスト（フェーズ5）と、GAS 変更時の 1 回の手動確認で代替する。

## 調査で見つかった既存バグ

ロードマップ調査の副産物。issue 化の要否は MT レビュー後に判断する（挙動確認が済んでいない「疑い」を含む）。

- `api/lib/internals/repository/department_repository.go:42-45`: INSERT のカラム名が `departments`（実カラムは `department`）で恒常失敗する。
- `api/lib/internals/repository/review_repository.go:77-96`: comment のクォート連結により、アポストロフィ入り文字列で実行時エラー。SQL インジェクションも可能（#363 の一部）。
- `api/lib/internals/controller/user_controller.go:97`: `Access-Token` ヘッダ欠落時に index out of range で panic する疑い。
- `mobile/lib/models/shift_card.dart:78`: `before_members` が null の JSON でクラッシュする疑い。
- `docker-compose.yml:8`: 存在しない `./mysql/db` を initdb にマウント。DB 初期化が機能していない。`my.cnf`（MySQL 用設定）のマウントも遺物。
- `Makefile` L74 / L79 / L86: seed 系 3 ターゲットが存在しない `/app/seeds/seeds.go` を参照。
- `api/lib/usecase/shift_usecase.go:887`: ハードコードされた GAS URL の宛先 `doPost` がリポジトリに存在しない。

## 既存 issue との対応

| フェーズ | 関連 issue |
|---|---|
| フェーズ0 基盤整備 | #385（Go バージョン固定）、#261（CI/CD パイプライン） |
| フェーズ2〜3 repository | #363（SQL インジェクション）、#266（プレースホルダ化と crud 経由統一） |
| フェーズ3 以降のリファクタ全般 | #264 / #247（N+1 解消。フェーズ2 の安全網ができてから着手推奨） |
| フェーズ5 controller | #267（入力バリデーション強化） |
| 静的解析（完了済み扱い） | #314（golangci-lint 指摘の解消） |
| GAS | #263（CodeRabbit 指摘のバグ精査） |
