# SeeFTテストロードマップ

2026-07-06のMTで「SeeFT本体（mobile / api / admin）は保守性を高める方向」が決定した。この文書は、その次のアクションとして挙がった「どこからテストを書き始めるべきか」を、実コードの調査結果に基づいて定義する。調査時点で`api/`配下の`*_test.go`は0件、`mobile/test/`はデフォルトのまま残ったファイル1件のみで、CIはlint専用でテスト実行ステップは存在しない。つまり、テストはゼロから始めることになる。

> 2026-09-29時点で、フェーズ0は完了している（Goは1.26に統一、`go-test.yml`で`api/**`を変更するPRごとに`go test`が走る、composeのマウントは`./postgresql`に修正済み）。`api/`配下の`*_test.go`は18件ある。以下の本文は策定時点の記述のまま残している。

## 要約

実行順序は次のとおり。

1. フェーズ0: テスト実行基盤の整備（CIにgo test、Goバージョン一本化、存在しないディレクトリを指すcomposeマウントの修正）
2. フェーズ1: 依存ゼロの純関数テスト（モックもDBも不要。テストの書き方の規約をここで確立する）
3. フェーズ2: repository層のゴールデンテスト（実DBを使う統合テスト。「現状の動作を正解にする」の実装）
4. フェーズ3: repository層のクエリ修正（フェーズ2のテストを安全網にしてプレースホルダ化）
5. フェーズ4: usecase層のテスト（go-sqlmock）
6. フェーズ5: controller層のテスト（httptest。GASとの「結合テスト」はここの契約テストで代替する）

mobileはapiと独立して並行で進められる。GASとadminは自動テストの対象外とする（理由は後述）。

MTの叩き台（静的解析 → repository修正 → repositoryテスト → usecaseテスト → handlerテスト → GASテスト → GAS×API結合テスト）からは、2点を変えた。

1点目として、「静的解析導入」は完了済みのため外した。`.github/workflows/go-lint.yml`でgolangci-lint v2.12がPRごとに走っており、残っているのは指摘の解消（#314）で、これはテストとは独立に進められる。

2点目として、「repository修正 → repositoryテスト」を逆順にした。テストがない状態でクエリを書き換えると、クエリが返す行が意図どおりに変わっただけなのか、意図していない行まで変わったのかを判定できない。MTで出た「現状の動作を正解にする」アプローチそのものが、修正より先にテストを書く理由になる。

## 前提: コード構造がテスト戦略を規定する

apiは素直なレイヤ構成で、repository・usecase・controllerは`api/lib/di/di.go`でコンストラクタに渡して組み立てている。依存はほとんどインターフェースで受けているが、`notificationUseCase`だけは`*slack.SlackService`を具象型のまま受けている。テストを書くうえでは恵まれた構成である。ただし、repositoryのインターフェースが`*sql.Rows` / `*sql.Row`という`database/sql`の具象型を返し、Scanをusecase側が行う設計になっている点だけは、大きな制約である（18ファイル中17ファイル。例外は`shift_card_repository.go`のみ）。

```go
// api/lib/internals/repository/bureau_repository.go
type BureauRepository interface {
	All(context.Context) (*sql.Rows, error)
	Find(context.Context, string) (*sql.Row, error)
	// ...
}
```

`*sql.Row`はテストコードから作れないため、repositoryのインターフェースをモックしてもusecaseのテストは書けない。対処は2案ある。

案Aは、`db.Client`（`api/lib/externals/db/db.go:19-23`、`DB() *sql.DB`を持つインターフェース）に[go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)の偽`*sql.DB`を差し込む方法である。既存コードの改修は不要で、repositoryの実装ごとテストを通す形になる。

案Bは、repositoryの戻り値をentityに変えるリファクタである。テストは書きやすくなるが、141メソッドに及ぶ大がかりな作業で、保守モードの規模を超える。

このロードマップは案Aを採用する。以降のフェーズは、すべて案Aを前提に書いている。

## テスト設計ステージ: 各フェーズの回し方

フェーズ1以降の各フェーズは、いきなりテストコードを書くのではなく、次のサイクルで回す（フェーズ0だけは基盤作業なので設計ステージはない）。

```text
設計起案（AI・強いモデル）→ 設計レビュー（人間: MT または PR）→ 実装量産（AI・Sonnet で可）→ 検証（CI + CodeRabbit）
```

設計起案では、対象関数ごとにAIが「テストケース表」を起案する。表の列は、ケース名、分類（正常系・境界値・異常系）、入力、期待値、根拠、要判断フラグである。期待値は起案のまま使わず、使い捨てテストを実際に実行して実際の動きと比べてから設計書に載せる。推測の期待値は設計書に残さない。

設計レビューで人間がレビューするのは、この表であって、テストコードではない。見る点は、ケースの過不足（この入力パターンが漏れていないか）と、「要判断」ケースの裁定（現状の動きをバグとして直すか、仕様として固定するか）の2つである。要判断ケースの裁定は、このロードマップの方針「現状の動作を正解にする」に対する唯一の例外を決める判断であり、MTの議題にそのまま使える。

実装量産では、承認済みのケース表をプロンプトに貼って実装させる。設計が確定していれば機械的な作業なので、Sonnetクラスのモデルで十分である。

検証は、CIのgo testとCodeRabbitで行う。人間は、diffがケース表と一致しているかを見る。

### Claudeへの依頼方法

依頼単位は1 issue = 1〜2関数とする。生成は安いがレビューは人間がやるので、レビューできる粒度を守る。

設計起案のプロンプトには、対象関数のファイルパスと行番号、このロードマップの該当フェーズの節、上記のケース表の形式、「現状の動作を正解にする。ただし直感に反する挙動には要判断フラグを付ける」という方針を含める。実装のプロンプトには、承認済みのケース表、テストの置き場所（対象と同一パッケージ内の`*_test.go`）、テーブル駆動テストで書くことを含める。

フェーズ1の設計書第1弾は`docs/development/test-design/phase1-pure-functions.md`にある（この節のサイクルを実際に回して作成したもの）。

## フェーズ0: テスト実行基盤の整備

テストを1本も書く前に、書いたテストが回る環境を作る。すべて小さい独立したPRにできる。

- Goバージョンの一本化（既存issue #385）。現状は`api/go.mod`が`go 1.16`、`api/Dockerfile`と`api/prod.Dockerfile`が`golang:latest`、`go-lint.yml:20`が`go-version: stable`と、3か所で食い違っている。go directiveが1.16のままだとtestifyやgo-sqlmockの新しいバージョンが入らないため、テスト導入の直接の前提になる。採用バージョンの比較・実機検証結果は[go-version-comparison.md](./go-version-comparison.md)にまとめた。
- `docker-compose.yml:8`と`docker-compose.mac.yml`のマウント修正。存在しない`./mysql/db`をinitdbにマウントしており（#277のディレクトリ改名への追随漏れ）、DBの初期化が機能していない。`./postgresql/db`に直す。MySQL時代から残っている`my.cnf`のマウントも合わせて整理する。
- CIにgo testジョブを追加。`go-lint.yml`と同型で`working-directory: api`、`go test ./... -count=1`。テストが0本でもgreenになるので、最初に入れておくと以降のフェーズのPRが全部CIで検証される。
- `Makefile`に`test`ターゲットを追加し、`AGENTS.md`のCommands節に記載する。なお`make seed`系3ターゲットは参照先の`/app/seeds/seeds.go`が存在せず、3つとも`go run`がファイルが見つからないエラーで終わる（後述のバグ一覧参照）。
- `go.mod`に`stretchr/testify`と`DATA-DOG/go-sqlmock`を追加。あわせて`api/lib/usecase/shift_usecase.go:66`の未使用グローバル変数（`var TaskID, UserID, ... string`）を削除しておく。

## フェーズ1: 依存ゼロの純関数テスト

モックもDBも使わず、フェーズ0の完了も待たずに書ける関数が既にある。ここでテーブル駆動テストの規約（ファイル配置、命名、テストケースの書き方）を確立し、以降のフェーズの雛形にする。

- `api/lib/usecase/notification_usecase.go`のヘルパー4関数: `GroupNotificationsByUserAndDate`（L138）、`sortLogsByTime`（L317）、`formatTimeRange`（L359）、`buildChangesWithTime`（L371）。unexportedなレシーバでもフィールドに触れないため、同一パッケージ内テストから`(&notificationUseCase{})`で直接呼べる。`formatTimeRange`には「`endTimeID+1`がtimeMapに無いと`0:00`にフォールバックする」という境界条件が既にあり、テストで固定する価値が高い。Slack通知の文言は、これらの関数が組み立てたまま利用者に届く。時刻やタスク名を誤って組み立てると、利用者は誤った時刻や担当を読むことになる。このテストは通知文言の回帰防止に直結する。
- `api/lib/usecase/shift_usecase.go`の純関数3つ: `groupContinuousShifts`（L478、連続TimeIDのグループ化）、`compareTimeStrings`（L509、ゼロ埋めなし時刻文字列の比較。文字列のまま比べると"8:00"が"10:00"より後に並ぶため、時と分を数値に直して比べている）、`convertShiftCardDataToShifts`（L433）。
- `api/lib/externals/slack/slack_service.go`の`BuildMessageBlocks`（L88）。Slack APIと通信せずに、Block構成の分岐を検証できる。
- `api/lib/externals/scheduler/scheduler.go`の`Start`。カウンタを増やすだけのJobと短いintervalで、「起動直後の即時実行」「intervalごとの再実行」「ctxキャンセルで停止」「jobエラーでもループ継続」を検証できる。このパッケージはPR #322（通知の定期実行）で追加されるため、マージ後に着手する。

## フェーズ2: repositoryのゴールデンテスト（実DB統合）

フェーズ2は、「現状の動作を正解にする」を実装するフェーズである。ゴールデンテスト（golden master test。レガシーコードの文脈では特性化テスト/characterization testとも呼ぶ）とは、正しい仕様を定義するのではなく、いま動いているコードの出力をそのままテストの期待値として固定するテストである。ゴールデンテストがあると、次のフェーズでクエリを書き換えたときに「テストが緑のまま = 挙動が変わっていない」と機械的に判定できる。

実行環境にはGitHub Actionsのservice containerを使う。`postgres:18`（`docker-compose.yml`と同じイメージ・資格情報）をservice containerとして起動し、`postgresql/db/*.sql`を番号順（create1 → create6 → seed.sql）に投入する共通セットアップを作る。composeのinitdbマウントは存在しないディレクトリを指していてSQLが適用されないため、テストヘルパー側でSQLを適用する。DBが必要なテストにはbuild tag `integration`を付け、通常の`go test ./...`と分離する。接続は環境変数（`NUTMEG_DB_HOST`等）だけを参照する既存実装のままで済み、CIでは`localhost`を指定する。

最初の3本は、テストの価値がすぐ実証できるものを選ぶ。

1. `departmentRepository.Create`（`department_repository.go:42-45`）。実カラム名は`department`なのに`departments`へINSERTしており、恒常的に失敗する既存バグがある。1本目のテストがいきなりバグを検出するので、テスト導入の説得材料になる。
2. `reviewRepository.Create / Update`（`review_repository.go:77-96`）。自由記述のcommentをクォート連結しているため、`It's fine`のようなアポストロフィ入り文字列で現状は失敗する。このテストはフェーズ3のプレースホルダ化の完了判定を兼ねる。
3. `shiftRepository`のCRUD一巡（Create → FindByUnique → Update → Destroy）。文字列連結が17行と最多で、シフトというドメインの中核でもある。

## フェーズ3: repositoryのクエリ修正

MTで「repositoryのクエリがゴミ」と言われていた部分の実態は、11ファイル・約58行の文字列連結SQLである。controllerは入力検証なしでHTTPパラメータをそのまま渡すため、この文字列連結SQLは書式の問題ではなく、実際にSQLインジェクションができる状態にある（既存issue #363 / #266）。

```go
// api/lib/internals/repository/shift_repository.go:47
query := "SELECT * FROM shifts WHERE id = " + id
```

フェーズ2のテストを安全網として、連結を`$1`プレースホルダに置換していく。同じファイル内に正しい書き方（`shift_repository.go`の`Users`メソッドは`$1〜$5`を使用）が既にあるので、その書き方に揃える。最初に直すのはユーザーの自由記述が入るreviewで、次が連結箇所の最も多いshiftである。あわせて、`department_repository.go:42`のカラム名バグも直す。この修正は挙動の変更（常に失敗 → 成功する）なので、ゴールデンテストの期待値を意識的に更新する。無条件の`fmt.Printf`（クエリのデバッグ出力）も`DEBUG_SQL`環境変数ガードに統一する。テスト出力が汚れるのを防ぐ目的も兼ねた、機械的な整理である。

## フェーズ4: usecaseのテスト（go-sqlmock）

案Aの方式で、`db.Client`を満たす小さなフェイク（`DB()`がsqlmockの`*sql.DB`を返すstruct）を1つ作り、repositoryの実装ごとusecaseを通す。Scanのカラム順まで再現する必要があるためフェーズ1より手間はかかるが、フェイクは全usecaseで使い回せる。

雛形には最小の`bureau_usecase.go`（71行）を使う。`GetBureaus`（L24）で複数行のScanとフィールドマッピング、`GetBureauByID`（L54）で1行取得、クエリエラー時のエラー伝播、空リスト時にnilを返す現状の動き（AGENTS.md規約の`[]Type{}`との差異）を固定する。以降は価値の高い順に、`mail_auth_usecase.SignIn`（L33、認証というクリティカルパス。全角スペース除去 → bcrypt照合のロジック）、`notification_usecase.ProcessUnsentNotifications`（通知の本丸）と進める。

`notificationUseCase`は具象型`*slack.SlackService`に直接依存しているため（`notification_usecase.go:21`）、Slack送信まで含めたend-to-endテストにはインターフェース抽出のリファクタが要る。このリファクタは別issueに切り出し、それまではsqlmockで届く範囲をテストする。

## フェーズ5: controllerのテスト

controllerはusecaseインターフェースだけに依存する薄い層なので、手書きフェイク + echoのhttptestで書ける。薄いぶん網羅は狙わず、価値のある箇所に絞る。

FromGAS系3ハンドラ（`router.go:177-179`の`/api/update_users`、`/api/update_tasks_and_places`、`/api/update_shifts`）には、契約テストを書く。GASが実際に送るJSONをフィクスチャとして与え、パースとusecaseへの受け渡しを固定する。GAS側をテストしなくても、api側の変更で、GASが送るJSONのパースがエラーになったり、usecaseに渡る値が変わったりしたことを検知できる。この契約テストが、後述するGAS×API「結合テスト」の現実的な答えになる。

`user_controller.go:97`の`c.Request().Header["Access-Token"][0]`は、ヘッダ欠落時にindex out of rangeでpanicする疑いがある。テストで最初に検証すべき欠陥候補である。

リクエストの型変換・バリデーション整備（既存issue #267）に着手する際は、このフェーズのテストを先に書いてから行う。

## 並行トラック: mobile

mobileはapiと依存関係がなく、Flutterに慣れたメンバーへ独立して割り振れる。直近12ヶ月で104コミットと活発なため、テスト投資の回収が見込める。

- 最初のPR: `mobile/test/widget_test.dart`（存在しない`main.dart`のMyAppを参照する、デフォルトのまま残ったファイルで、`flutter test`をCIに足すと即redになる）を削除または書き直し、同じPRで`flutter-lint.yml`に`fvm flutter test`相当のステップを追加する。
- `mobile/lib/models/shift_card.dart`の`ShiftCardDataList.fromJson`（L53-105）。依存ゼロの純関数で、すぐ着手できる。L78で`item['before_members']`自体のnullチェックがなくクラッシュする疑いがあるため、正常系と合わせて固定する。最も頻繁に変更される画面であるmy_shift_pageの入力データを守る1本目として、費用対効果が最も高い。
- `mobile/lib/models/rescue.dart`のtypeディスパッチ（L33-44）。trouble / question / shorthandedのラウンドトリップと未知typeの挙動。
- Newバッジ判定ロジック（`my_shift_page.dart`の`_isCardChanged` L420、`_detectNewOrUpdatedCardKeys` L437ほか）。判定を誤っても、変更のあったカードにNewが付かないか、変わっていないカードに付くだけで、エラーは出ない。そのため画面上で気づきにくい箇所の典型だが、privateなStateクラス内にあるため、テストするには純関数として`lib/utils/`へ抽出する移動リファクタが先になる。modelsのテストが定着してから着手する。

## GASの方針: 自動テストの対象外とする

`gas/`は4プロジェクト・約1,474行あるが、内容は「スプレッドシートとapiの双方向同期グルーコード」で、全関数がSpreadsheetApp / UrlFetchAppなどのGAS APIに依存している。userとshiftはファイルのトップレベルで`SpreadsheetApp.getUi()`等を実行するため、Nodeで読み込むとすぐにクラッシュし、テストランナーに載せること自体ができない。ここに自動テストを整備する費用対効果は低いと判断する。

代わりに、GASが送るJSONをapiが受け付けなくなったことは、フェーズ5のFromGAS契約テスト（api側）で検知する。

例外として、GAS側を変更する必要が出たときは、変更対象のロジックを引数を取る純関数として別ファイルに抽出してから触る（rescue/onChange.jsのステータスマッピングや、shift/コード.jsのペイロード構築が候補）。45th対応で`yearID`前提の改修が入る可能性が高いのはシフトペイロード構築部分で、この部分だけは先行して抽出しておく価値がある。

また、調査で、`api/lib/usecase/shift_usecase.go:887`にGAS WebアプリのURLがハードコードされているのに、対応する`doPost`が`gas/shift/`に存在しないことが分かった。リポジトリと実デプロイの差分を`clasp pull`で確認し、どちらを基準にするか決める作業が先に必要である。

## adminの方針: 対象外

`admin/`は直近12ヶ月コミット0件の凍結状態で、AGENTS.mdにも不使用と明記されている。テスト投資はしない。もし再稼働させる判断をした場合のみ、CIに`next build`のスモークジョブを足すところから始める。

## MTの宿題:「結合テストって最初と最後どっちだっけ」

「結合テスト」が2つの意味で使われているので、分けて答える。

api×DBの統合テスト（フェーズ2）は最初の側に来る。教科書的な「単体を固めてから結合」は、テストピラミッドを新規開発で下から積む場合の話で、テストのないコードを保守する場合は逆になる。単体テストを書くための分解（モック化・リファクタ）自体が挙動を変えるリスクを持つため、先に外側から現状の動作を固定し、その安全網の中で内側を整えるのが、レガシーコード改善の定石である。

GAS×APIをプロセス横断で通すE2Eは最後に来る。というより、自動化しない。GAS側に自動実行環境を組めないため、api側の契約テスト（フェーズ5）と、GAS変更時の1回の手動確認で代替する。

## 調査で見つかった既存バグ

ロードマップの調査で副次的に見つかったものである。issue化の要否はMTレビュー後に判断する（挙動確認が済んでいない「疑い」を含む）。

- `api/lib/internals/repository/department_repository.go:42-45`: INSERTのカラム名が`departments`（実カラムは`department`）で恒常的に失敗する。
- `api/lib/internals/repository/review_repository.go:77-96`: commentのクォート連結により、アポストロフィ入り文字列で実行時エラー。SQLインジェクションも可能（#363の一部）。
- `api/lib/internals/controller/user_controller.go:97`: `Access-Token`ヘッダ欠落時にindex out of rangeでpanicする疑い。
- `mobile/lib/models/shift_card.dart:78`: `before_members`がnullのJSONでクラッシュする疑い。
- `docker-compose.yml:8`: 存在しない`./mysql/db`をinitdbにマウント。DB初期化が機能していない。`my.cnf`（MySQL用設定）のマウントも、MySQL時代から残っているもの。
- `Makefile` L74 / L79 / L86: seed系3ターゲットが存在しない`/app/seeds/seeds.go`を参照。
- `api/lib/usecase/shift_usecase.go:887`: ハードコードされたGAS URLの宛先`doPost`がリポジトリに存在しない。

## 既存issueとの対応

| フェーズ | 関連issue |
|---|---|
| フェーズ0基盤整備 | #385（Goバージョン固定）、#261（CI/CDパイプライン） |
| フェーズ2〜3 repository | #363（SQLインジェクション）、#266（プレースホルダ化とcrud経由統一） |
| フェーズ3以降のリファクタ全般 | #264 / #247（N+1解消。フェーズ2の安全網ができてから着手推奨） |
| フェーズ5 controller | #267（入力バリデーション強化） |
| 静的解析（完了済み扱い） | #314（golangci-lint指摘の解消） |
| GAS | #263（CodeRabbit指摘のバグ精査） |
