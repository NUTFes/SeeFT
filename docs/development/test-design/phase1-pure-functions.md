# フェーズ1テスト設計書: 依存ゼロの純関数（第1弾）

[テストロードマップ](../test-roadmap.md)の「テスト設計ステージ」のサイクルを実際に回して作った設計書。対象はフェーズ1の純関数8つ。各ケース表はAIが起案したあと、**全ケースを使い捨てテストとして実際に実行し、期待値が実際の動きと合うことを確かめてある**。つまりこの表の期待値は推測ではなく、現時点のdevelopの実際の出力である。

レビューで見てほしい点は2つ。

1. ケースの過不足: 守りたい入力パターンが漏れていないか
2. 「要判断」の裁定: 今の動きをこのまま仕様として固定してよいか、バグとして直すか。直す場合は、テストを実装する前に修正のissueを立てる

## サマリ

| 関数 | 場所 | ケース数 | 要判断 |
|---|---|---|---|
| `GroupNotificationsByUserAndDate` | `api/lib/usecase/notification_usecase.go:138` | 8 | 0 |
| `sortLogsByTime` | `api/lib/usecase/notification_usecase.go:317` | 9 | 2 |
| `formatTimeRange` | `api/lib/usecase/notification_usecase.go:359` | 8 | 1 |
| `buildChangesWithTime` | `api/lib/usecase/notification_usecase.go:371` | 10 | 2 |
| `convertShiftCardDataToShifts` | `api/lib/usecase/shift_usecase.go:433` | 8 | 2 |
| `groupContinuousShifts` | `api/lib/usecase/shift_usecase.go:478` | 10 | 1 |
| `compareTimeStrings` | `api/lib/usecase/shift_usecase.go:509` | 10 | 2 |
| `BuildMessageBlocks` | `api/lib/externals/slack/slack_service.go:88` | 8 | 3 |
| 合計 | | 71 | 13 |

## GroupNotificationsByUserAndDate

`api/lib/usecase/notification_usecase.go:138`

[]entity.ActionLogを受け取り、各ログのUserIDとDateIDから "%d_%d" 形式の文字列キーを作って、map[string][]entity.ActionLogにグループ分けする純関数（L138-145）。同じキーの中の要素は、appendで足していくため、入力スライスに出てきた順のまま並ぶ。レシーバnのフィールドには一切依存しないので、テストでは &notificationUseCase{}（ゼロ値）で呼び出せることを実行で確かめた。なお呼び出し元のProcessUnsentNotifications（L105-117）は、このキーをstrings.SplitとAtoiで分解し直している。そのためキー書式 "%d_%d" は事実上の契約であり、テストで固定する価値が高い。

### ケース表

1. 複数ユーザー・複数日付の混在を正しく分配する（正常系）
   - 入力: `logs := []entity.ActionLog{{ID: 1, UserID: 1, DateID: 10}, {ID: 2, UserID: 2, DateID: 10}, {ID: 3, UserID: 1, DateID: 10}, {ID: 4, UserID: 1, DateID: 11}}`
   - 期待値: `map[string][]entity.ActionLog{"1_10": {{ID: 1, UserID: 1, DateID: 10}, {ID: 3, UserID: 1, DateID: 10}}, "2_10": {{ID: 2, UserID: 2, DateID: 10}}, "1_11": {{ID: 4, UserID: 1, DateID: 11}}}（reflect.DeepEqual で比較可能。キー数は 3）`
   - 根拠: 関数の中心となる動き（UserIDとDateIDの組ごとの分配）の回帰を防ぐ。同じユーザーでもDateIDが違えば別のグループに、同じ日付でもUserIDが違えば別のグループになることを同時に検証する
2. 同じグループの中で入力順が保たれる（正常系）
   - 入力: `logs := []entity.ActionLog{{ID: 5, UserID: 7, DateID: 3}, {ID: 2, UserID: 7, DateID: 3}, {ID: 9, UserID: 7, DateID: 3}}`
   - 期待値: `map[string][]entity.ActionLog{"7_3": {{ID: 5, UserID: 7, DateID: 3}, {ID: 2, UserID: 7, DateID: 3}, {ID: 9, UserID: 7, DateID: 3}}}（ID の並びが入力順 5, 2, 9 のまま。ID 昇順にソートされない）`
   - 根拠: appendによる挿入順の保持の回帰を防ぐ。時刻順のソートは後段のsortLogsByTimeが受け持つ分担になっており、この関数が勝手にソートしない今の動きを固定する
3. 空スライスなら非nilの空mapを返す（境界値）
   - 入力: `logs := []entity.ActionLog{}`
   - 期待値: `got != nil かつ len(got) == 0（make で初期化された空 map。nil map ではない）`
   - 根拠: 空入力での動きを固定する。呼び出し元はlen(logs)==0でガードしているが、関数単体として非nilの空mapを返す契約を固定し、将来ほかの場所から直接呼んだときにnil mapへのrangeや代入で起きる問題を防ぐ
4. nilスライスでも非nilの空mapを返す（境界値）
   - 入力: `var logs []entity.ActionLog（nil スライスのまま渡す）`
   - 期待値: `got != nil かつ len(got) == 0（空スライスと同一の結果。panic しない）`
   - 根拠: Goではnilスライスのrangeは安全だが、その事実に頼っている今の動きを明示的に固定する。nilと空で動きが分かれないことの回帰を防ぐ
5. 要素1個は単一キー・単一要素になり、全フィールドがそのまま渡される（境界値）
   - 入力: `` logs := []entity.ActionLog{{ID: 1, ShiftID: 55, UserID: 42, DateID: 7, ActionType: "UPDATE", DiffPayload: json.RawMessage(`{"changes":[]}`), IsSent: false, CreatedAt: time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)}} ``
   - 期待値: `len(got) == 1、got["42_7"] が長さ1のスライスで、その唯一の要素が入力の entity.ActionLog と reflect.DeepEqual で完全一致（ShiftID・ActionType・DiffPayload・CreatedAt 等が改変されずに素通しされる）`
   - 根拠: 最小構成での正常な動きと、UserID/DateID以外のフィールドを一切書き換えずにそのまま渡すことの回帰を防ぐ。後段のprocessGroupはDiffPayloadやShiftIDをそのまま使うため、書き換えられると、DiffPayloadから読むタスク名や、ShiftIDから引く時刻・タスクが、元のログと違う値のまま通知に載る
6. ゼロ値のActionLogはキー "0_0" にグループ分けされる（境界値）
   - 入力: `logs := []entity.ActionLog{{}}（全フィールドゼロ値。UserID=0, DateID=0, DiffPayload=nil）`
   - 期待値: `len(got) == 1、got["0_0"] == []entity.ActionLog{{}}（panic せず、nil の DiffPayload もそのまま保持）`
   - 根拠: ゼロ値の入力での安全性を固定する。DBスキャンの異常などでUserID/DateIDが0のまま流れてきても本関数は落ちず "0_0" にまとめる、という入力を検査しない今の動きを明示する
7. キー書式 "%d_%d" により(1,23)と(12,3)が衝突しない（境界値）
   - 入力: `logs := []entity.ActionLog{{ID: 1, UserID: 1, DateID: 23}, {ID: 2, UserID: 12, DateID: 3}}`
   - 期待値: `len(got) == 2、キー "1_23" と "12_3" がそれぞれ存在し、各スライスの長さは 1（1つのキーに混ざらない）`
   - 根拠: キー書式そのものの回帰を防ぐ。区切り文字 "_" を外して "%d%d" などに変えると、どちらも "123" になって衝突し、別人の通知が混ざる。さらに呼び出し元（L108-117）は "_" でSplitしてAtoiする契約に依存しているため、書式の変更はすぐに検知したい
8. 負のIDでも検証されずにそのままキーになる（異常系）
   - 入力: `logs := []entity.ActionLog{{ID: 1, UserID: -1, DateID: -2}}`
   - 期待値: `len(got) == 1、got["-1_-2"] が長さ1のスライスとして存在する（エラーにも panic にもならず、負値がそのままキーに埋め込まれる）`
   - 根拠: 不正なドメイン値に対して何も検証しない今の動きを固定する。なお生成されるキー "-1_-2" はアンダースコアが1個なので、呼び出し元のSplit/Atoi（L108-117）でも正しく -1, -2に戻せ、Invalid group keyとして読み飛ばされたり、別のIDとしてprocessGroupに渡ったりしないことを確かめた

実行検証: worktreeのapi/lib/usecase/に使い捨てテストdesign_verify_groupnotificationsbyuseranddate_test.goを作り、ゼロ値レシーバ(&notificationUseCase{})からGroupNotificationsByUserAndDateを直接呼び出して、全8ケースをサブテストとして実装した。cd api && go test ./lib/usecase/... -run TestDesignVerifyGroupNotificationsByUserAndDate -vで実行し、8ケースすべてが初回でPASSした（8/8起案どおり、修正0件、削除0件、needs_judgment 0件）。非nilの空map（空スライスとnilスライスの両方）、挿入順の保持、フィールドをそのまま渡すこと（DiffPayload/CreatedAtを含めてDeepEqualで完全一致）、キー "0_0" の生成、"1_23" と "12_3" が衝突しないこと、"-1_-2" の生成も、すべて実際の動きで確かめた。テストファイルは削除し、git status --porcelainが空であることを確かめた。commit/pushは行っていない。

## sortLogsByTime

`api/lib/usecase/notification_usecase.go:317`

ActionLogのスライスを、各ログのShiftIDをキーにshiftMapから引いたShiftAdmin.TimeIDの昇順で並べ替え、新しいスライスとして返す。shiftMapにキーが無いログは、エラーもログも出さずに結果から除外されるため、ソートとフィルタを兼ねている。レシーバnのフィールドには一切依存しない純関数であり、テストではゼロ値レシーバ &notificationUseCase{} で呼び出せる（参照するフィールドはActionLog.ShiftIDとShiftAdmin.TimeIDだけ）。全9ケースを実際のコードで実行して確かめた。

### 要判断の論点

shiftMapにキーが無いログがエラーもログも出さずに除外される点は、実行で確かめたところ、実際の動きも起案どおりだった（除外され、ID列は[3, 1]）。論点は動きの食い違いではなく設計判断である。シフトが削除済みなどの理由でshiftMapに載らないログは、通知から何の知らせもなく抜け落ち、エラーにもログにも出ない。今の動きを仕様として固定するか、除外したときに警告ログを出す・件数を返すなどの形に改めるか、判断が必要である。テスト自体は今の動き（エラーもログも出さずに除外する）を期待値とする。

同じTimeIDのログは全件残るが、その相対順序は保証されない。今回の実行（go環境の現行バージョン、3要素の入力）で観測した順序は挿入順[1, 2, 3]だった。しかしsort.Sliceは安定ソートではなく、同じTimeIDの中の相対順序は仕様上決まっていない（Goのバージョンや入力サイズで変わり得る）。テストは順序を問わず、集合が一致すればパスする。同じ時刻のグループの中で通知行の並びを安定させたいなら、sort.SliceStableへの変更と、順序を確かめるテストの追加を検討すべきである。順序を問わないテストで今の動きを認めるか、SliceStableにして順序を仕様にするか、判断が必要である。

### ケース表

1. 複数ログをTimeID昇順に並べ替える（正常系）
   - 入力: `logs := []entity.ActionLog{{ID: 1, ShiftID: 10}, {ID: 2, ShiftID: 20}, {ID: 3, ShiftID: 30}} / shiftMap := map[int]entity.ShiftAdmin{10: {ID: 10, TimeID: 3}, 20: {ID: 20, TimeID: 1}, 30: {ID: 30, TimeID: 2}}`
   - 期待値: `返り値の ID 列が [2, 3, 1]（TimeID 1→2→3 の順）。長さ3`
   - 根拠: 関数の主な目的である昇順ソートの基本の動きを固定する。ソート条件を誤ってCreatedAtなどに変えたときの回帰を検出する
2. 要素1個はそのまま返る（境界値）
   - 入力: `logs := []entity.ActionLog{{ID: 1, ShiftID: 10}} / shiftMap := map[int]entity.ShiftAdmin{10: {ID: 10, TimeID: 5}}`
   - 期待値: `長さ1で、返り値[0].ID == 1（入力と同一内容のログ）`
   - 根拠: ソート対象が要素1個でも、抜け落ちたり内容が書き換わったりしないことを保証する最小の境界
3. 空スライスは空の非nilスライスを返す（境界値）
   - 入力: `logs := []entity.ActionLog{} / shiftMap := map[int]entity.ShiftAdmin{10: {ID: 10, TimeID: 1}}`
   - 期待値: `len(result) == 0 かつ result != nil（make([]entity.ActionLog, 0) 由来の空スライス）`
   - 根拠: 呼び出し元のL284はlen(sortedLogs) > 0で分岐しており、空入力でpanicせずに空が返ることが、通知をスキップする前提になっている
4. nilスライスでもpanicせず空を返す（境界値）
   - 入力: `var logs []entity.ActionLog = nil / shiftMap := map[int]entity.ShiftAdmin{10: {ID: 10, TimeID: 1}}`
   - 期待値: `panic せず len(result) == 0 かつ result != nil`
   - 根拠: nilスライスへのrangeとmake(cap=0)が安全である今の動きを固定する。将来logs[0]を参照するような変更が入ったときの回帰を検出する
5. nilマップでは全ログが除外され空を返す（境界値）
   - 入力: `logs := []entity.ActionLog{{ID: 1, ShiftID: 10}, {ID: 2, ShiftID: 20}} / var shiftMap map[int]entity.ShiftAdmin = nil`
   - 期待値: `panic せず len(result) == 0 かつ result != nil（nil マップの lookup は ok=false になり全件 continue）`
   - 根拠: loadShiftMapが空のmapやnilを返すケースでの安全性を固定する。nilマップの読み取りはGoでは合法だが、書き込みを伴う実装に変わるとpanicするため、その回帰を検出できる
6. shiftMapにキーが無いログはエラーもログも出さずに除外される（異常系）【要判断】
   - 入力: `logs := []entity.ActionLog{{ID: 1, ShiftID: 10}, {ID: 2, ShiftID: 99}, {ID: 3, ShiftID: 30}} / shiftMap := map[int]entity.ShiftAdmin{10: {ID: 10, TimeID: 2}, 30: {ID: 30, TimeID: 1}}（ShiftID 99 のエントリ無し）`
   - 期待値: `長さ2、ID 列は [3, 1]。ShiftID=99 のログ（ID: 2）は結果に含まれない。エラーも panic も無し`
   - 根拠: 関数名はsortだが、フィルタを兼ねるのが今の動きの中心である。呼び出し元L285のshiftMap[sortedLogs[0].ShiftID]はキーが見つかる前提で書かれており、その前提をこの除外が保っている。除外がなくなると、別の場所でゼロ値を参照することになる
7. TimeIDのゼロ値・負値は正の値より前に並ぶ（境界値）
   - 入力: `logs := []entity.ActionLog{{ID: 1, ShiftID: 10}, {ID: 2, ShiftID: 20}, {ID: 3, ShiftID: 30}} / shiftMap := map[int]entity.ShiftAdmin{10: {ID: 10, TimeID: 2}, 20: {ID: 20}, 30: {ID: 30, TimeID: -1}}（ShiftID 20 は TimeID がゼロ値 0）`
   - 期待値: `ID 列が [3, 2, 1]（TimeID -1 → 0 → 2 の順）`
   - 根拠: TimeIDが未設定（ゼロ値）のシフトが先頭に来る今の動きを固定する。Scan漏れや未設定のデータが入り込んだとき、通知の並びがどうなるかを明文化する
8. 同一TimeIDのタイは全件保持される（相対順序は未保証）（境界値）【要判断】
   - 入力: `logs := []entity.ActionLog{{ID: 1, ShiftID: 10}, {ID: 2, ShiftID: 20}, {ID: 3, ShiftID: 30}} / shiftMap := map[int]entity.ShiftAdmin{10: {ID: 10, TimeID: 1}, 20: {ID: 20, TimeID: 1}, 30: {ID: 30, TimeID: 1}}`
   - 期待値: `長さ3で、ID の集合が {1, 2, 3}（順序は不問のアサートにする。例: ID を収集してソート後 [1, 2, 3] と比較）`
   - 根拠: タイでも要素が落ちたり重複したりしないことを保証する。順序を確かめないのは、安定ソートでないsort.Sliceを使っているため
9. 入力スライスを書き換えない（正常系）
   - 入力: `logs := []entity.ActionLog{{ID: 1, ShiftID: 10}, {ID: 2, ShiftID: 20}} / shiftMap := map[int]entity.ShiftAdmin{10: {ID: 10, TimeID: 2}, 20: {ID: 20, TimeID: 1}}（ソートで順序が入れ替わる入力）`
   - 期待値: `返り値の ID 列は [2, 1] だが、呼び出し後の logs の ID 列は [1, 2] のまま。かつ返り値は logs と別のスライス（&result[0] != &logs[0]）`
   - 根拠: 新しいスライスを組み立て、元のスライスを書き換えない実装を固定する。呼び出し元のL277以降ではlogsとsortedLogsを別のものとして扱っており、その場で並べ替える実装に書き換わると、呼び出し元の動きが知らないうちに変わるため

実行検証: worktreeのapi/lib/usecase/design_verify_sortlogsbytime_test.go（package usecase、ゼロ値レシーバ &notificationUseCase{} を使う）に全9ケースをサブテストとして実装し、cd api && go test ./lib/usecase/... -run TestDesignVerifySortLogsByTime -vで実行した。9ケースすべてが起案どおりの期待値でPASSした（修正0件、削除0件、実行できないケースなし）。補足として、空スライス・nilスライス・nilマップの3ケースは、いずれも返り値が非nilの空スライスであることをアサートで確かめた。タイのケースはt.Logfで観測した順序[1 2 3]（挿入順）を記録したが、sort.Sliceは安定ソートではないため、アサートは集合の一致だけにした。needs_judgment=trueの2件は、実際の動きと期待値の食い違いではなく、設計判断の論点（エラーもログも出さない除外の是非、SliceStableにするかどうか）として残した。実行後にテストファイルを削除し、git status --porcelainが空であることを確かめた。commit/pushは行っていない。

## formatTimeRange

`api/lib/usecase/notification_usecase.go:359`

timeMapから、開始スロット（startTimeID）の時刻と、終了スロットの「次のスロット」（endTimeID+1）の時刻を引き、"HH:MM 〜 HH:MM" 形式（半角スペース + U+301C波ダッシュ + 半角スペース区切り）の文字列を返す。終了時刻は、次スロットのキーが無い場合に限り "0:00" にフォールバックする。一方で開始時刻の側には存在チェックが無く、ゼロ値（空文字）がそのまま出る。レシーバnのフィールドには一切依存しない純関数であり、非公開メソッドなので、package usecase内のテストで &notificationUseCase{} を使って呼び出せる。全ケースを実際のコードで実行して確かめた。

### 要判断の論点

startTimeIDがマップに存在しないケースは、実際の動きが起案の期待値どおり " 〜 10:00" であることを確かめた。ただしendTime側（L363-366）にはokチェックと "0:00" へのフォールバックがあるのに、startTime側（L361）には存在チェックが無く、空文字がそのままSlack通知の文面に出る。この非対称には、バグの疑いがある。空文字の表示を仕様として固定するか、endTimeと同じようなフォールバック（"0:00" や "（不明）" など）を足すか、人間の判断が必要である。

### ケース表

1. 連続スロット範囲のフォーマット（正常系）
   - 入力: `timeMap := map[int]entity.Time{1: {Time: "9:00"}, 2: {Time: "10:00"}, 3: {Time: "11:00"}}; startTimeID := 1; endTimeID := 2`
   - 期待値: `"9:00 〜 11:00"（終了時刻は endTimeID=2 自身の "10:00" ではなく、endTimeID+1=3 の "11:00"）【実行確認済み】`
   - 根拠: 終了時刻にendTimeID+1のスロットの開始時刻を使う、という中心の仕様の回帰を防ぐ。リファクタで +1が外れると "9:00 〜 10:00" に変わり、検知できる。
2. 同一スロット指定（実際の呼び出し形態）（正常系）
   - 入力: `timeMap := map[int]entity.Time{5: {Time: "13:00"}, 6: {Time: "14:00"}}; startTimeID := 5; endTimeID := 5`
   - 期待値: `"13:00 〜 14:00"【実行確認済み】`
   - 根拠: 唯一の呼び出し元buildChangesWithTime（L385）は、startとendに同じshift.TimeIDを渡す。本番で実際に通る唯一の経路なので、その回帰を防ぐケースとして最も重要である。
3. 最終スロットで次スロットが無い（要素1個のマップ）（境界値）
   - 入力: `timeMap := map[int]entity.Time{10: {Time: "22:00"}}; startTimeID := 10; endTimeID := 10`
   - 期待値: `"22:00 〜 0:00"（endTimeID+1=11 が不在のため "0:00" フォールバック）【実行確認済み】`
   - 根拠: 1日の最終スロットに対する明示的なフォールバック（entity.Time{Time: "0:00"}）を固定する。要素1個のマップでの動作確認も兼ねる。
4. nilマップとゼロ値ID（境界値）
   - 入力: `var timeMap map[int]entity.Time = nil; startTimeID := 0; endTimeID := 0`
   - 期待値: `" 〜 0:00"（先頭は開始時刻が空文字のため半角スペースで始まる。panic しない）【実行確認済み】`
   - 根拠: nilマップの読み取りは、Goの仕様でゼロ値を返すためpanicしない。呼び出し元でtimeMapを作るのに失敗した最悪の場合でも、クラッシュしないことを固定する。
5. startTimeIDがマップに存在しない（異常系）【要判断】
   - 入力: `timeMap := map[int]entity.Time{2: {Time: "10:00"}}; startTimeID := 5; endTimeID := 1（endTimeID+1=2 は存在する）`
   - 期待値: `" 〜 10:00"（開始時刻が空文字のまま出力される。現状挙動の固定化）【実行確認済み】`
   - 根拠: shift.TimeIDがtimeMapに無い（データ不整合）ケース。startTime側が欠けたときの動きをテストで明示的に固定し、知らないうちに仕様が変わるのを防ぐ。
6. 次スロットは存在するがTimeフィールドがゼロ値（空文字）（境界値）
   - 入力: `timeMap := map[int]entity.Time{1: {Time: "9:00"}, 2: {}}; startTimeID := 1; endTimeID := 1`
   - 期待値: `"9:00 〜 "（末尾は半角スペースで終わる。キーは存在するため "0:00" フォールバックは発火しない）【実行確認済み】`
   - 根拠: "0:00" へのフォールバックが起きる条件が、「値が空かどうか」ではなく「キーが存在するかどうか」であることを固定する。ゼロ値のentity.Timeの扱いを明示する。
7. 負のID（endTimeID+1の算術でkey 0を参照）（異常系）
   - 入力: `timeMap := map[int]entity.Time{0: {Time: "8:00"}}; startTimeID := -1; endTimeID := -1`
   - 期待値: `" 〜 8:00"（startTimeID=-1 は不在で空文字、endTimeID+1=0 は key 0 にヒットし "8:00"）【実行確認済み】`
   - 根拠: IDの妥当性を検証しないこと、およびendTimeID+1の算術が負値にもそのまま適用されることを固定する。不正な入力でもpanicしないことを保証する。
8. 逆転範囲（startTimeID > endTimeID）（異常系）
   - 入力: `timeMap := map[int]entity.Time{1: {Time: "9:00"}, 2: {Time: "10:00"}, 3: {Time: "11:00"}, 4: {Time: "12:00"}}; startTimeID := 3; endTimeID := 1`
   - 期待値: `"11:00 〜 10:00"（開始 > 終了のまま検証なしで出力される）【実行確認済み】`
   - 根拠: 範囲の前後関係を検証しないことを固定する。呼び出し元は同じIDを渡すため実運用では起きないGIGOのケースだが、将来バリデーションを足すときに、このケースを意図して更新する目印になる。

実行検証: worktreeのapi/lib/usecase/design_verify_formattimerange_test.goに、全8ケースをテーブル駆動テスト（レシーバは &notificationUseCase{} のゼロ値）として実装し、cd api && go test ./lib/usecase/... -run TestDesignVerifyFormatTimeRange -vで実行した。8ケースすべてがPASSした（1回目の実行で全件一致）。起案どおりが8件、実際の動きとの食い違いによるexpectedの修正が0件、実行できずに削除・修正したケースが0件だった。「startTimeIDがマップに存在しない」のneeds_judgment=trueは、実際の動きとの食い違いによるものではない。start/endのフォールバックの非対称（L361にokチェックが無い）を仕様とするかバグとするか、設計判断が残るため維持した。検証後にテストファイルを削除し、git status --porcelainが空であることを確かめた。commit/pushは一切していない。

## buildChangesWithTime

`api/lib/usecase/notification_usecase.go:371`

アクションログの配列からSlack通知の本文を組み立てる純関数。各ログのdiff_payload(JSON)をパースし、shiftMap/timeMapから「開始 〜 終了：」の時間プレフィックスを付ける。そのうえでActionType（CREATE/UPDATE/DELETE/その他）ごとに整形した行を、改行で結合して返す。レシーバnはn.formatTimeRangeの呼び出しにだけ使われ、formatTimeRangeも含めて、レシーバのフィールド（repo群・slackService）には一切依存しない。そのため &notificationUseCase{} のゼロ値レシーバで、モック無しにテストできる（実行で確かめた）。

### 要判断の論点

timeMapにTimeID自体が無いと、開始時刻が空文字で出力される。実行して、起案どおりの動き（" 〜 0:00：受付 → 警備"）を確かめた。終了側が欠けたときは "0:00" へのフォールバックがあるのに、開始側が欠けたときは空文字のままSlack通知に載る、という非対称がある。開始側にもフォールバックの値を入れるか、時間プレフィックス自体を省くのが本来の意図ではないか、というバグの疑いがある。今の動きを仕様として固定してよいか、判断が必要である（直すなら期待値が変わる）。

diff_payloadがパースできない（不正なJSON / nil）ログは、ログ出力なしでスキップされる。実行して、起案どおりの動きを確かめた。パースできないログは、ログ出力すら無いまま通知から消えるため、シフト変更の通知漏れに直接つながりうる。この欠落を許容する仕様として固定するか、警告ログを足す・「（不明）」の行として出力するなどの形に改めるか、判断が必要である。テストの期待値は、今のスキップする動きにしている。

### ケース表

1. UPDATE正常系（全マップ完備）（正常系）
   - 入力: `` logs := []entity.ActionLog{{ID: 1, ShiftID: 10, ActionType: "UPDATE", DiffPayload: json.RawMessage(`{"changes":[{"field":"task","old":"受付","new":"警備"}]}`)}} shiftMap := map[int]entity.ShiftAdmin{10: {ID: 10, TaskID: 100, TimeID: 5}} taskMap := map[int]entity.Task{100: {ID: 100, Task: "警備"}} timeMap := map[int]entity.Time{5: {ID: 5, Time: "10:00"}, 6: {ID: 6, Time: "11:00"}} ``
   - 期待値: `"10:00 〜 11:00：受付 → 警備"（〜は全角チルダ、：は全角コロン、→ の前後は半角スペース）【実行確認済み】`
   - 根拠: いちばんよく通る経路の回帰を防ぐ。時間プレフィックスがtimeMap[TimeID] 〜 timeMap[TimeID+1]で作られること、payloadのold/newが使われることを固定する。
2. CREATE正常系（payloadのnewがDBの現在値より優先）（正常系）
   - 入力: `` logs := []entity.ActionLog{{ID: 2, ShiftID: 10, ActionType: "CREATE", DiffPayload: json.RawMessage(`{"changes":[{"field":"task","old":"","new":"設営"}]}`)}} shiftMap := map[int]entity.ShiftAdmin{10: {TaskID: 100, TimeID: 5}} taskMap := map[int]entity.Task{100: {Task: "現DB名"}} timeMap := map[int]entity.Time{5: {Time: "10:00"}, 6: {Time: "11:00"}} ``
   - 期待値: `"10:00 〜 11:00：設営（新規）"（taskMap の "現DB名" ではなく payload の "設営" が出る）【実行確認済み】`
   - 根拠: CREATEの分岐で、payload側のタスク名が、DBの現在値へのフォールバックより優先される仕様を固定する。taskMapの値をわざとずらしておくことで、優先順位の回帰を検出できる。
3. 複数ログの改行結合と入力順の保持（正常系）
   - 入力: `` logs := []entity.ActionLog{ {ShiftID: 10, ActionType: "UPDATE", DiffPayload: json.RawMessage(`{"changes":[{"old":"受付","new":"警備"}]}`)}, {ShiftID: 0, ActionType: "DELETE", DiffPayload: json.RawMessage(`{"deleted_task":"受付"}`)}, } shiftMap := map[int]entity.ShiftAdmin{10: {TaskID: 100, TimeID: 5}} taskMap := map[int]entity.Task{100: {Task: "警備"}} timeMap := map[int]entity.Time{5: {Time: "10:00"}, 6: {Time: "11:00"}} ``
   - 期待値: `"10:00 〜 11:00：受付 → 警備\n受付（削除）"（strings.Join による \n 結合。DELETE 行は ShiftID=0 が shiftMap に無いため時間プレフィックスなし）【実行確認済み】`
   - 根拠: 複数の行が入力順のまま \nで結合されることを固定する（この関数はソートしない。並べ替えは呼び出し前のsortLogsByTimeの責務であり、この分担の回帰を防ぐ）。
4. logsがnil/空スライス（境界値）
   - 入力: `logs := []entity.ActionLog(nil)（[]entity.ActionLog{} でも同一挙動） shiftMap := map[int]entity.ShiftAdmin{} taskMap := map[int]entity.Task{} timeMap := map[int]entity.Time{}`
   - 期待値: `""（changes が nil のまま strings.Join(nil, "\n") == ""）【nil・空スライスの両方を別サブテストで実行し、いずれも "" を確認済み】`
   - 根拠: 空入力で空文字を返す境界を固定する。呼び出し元buildGroupedMessageはlen==0を先に弾くが、この関数単体でも安全であることを保証する。
5. 全マップnil + DELETEログ（shift_id NULL相当のShiftID=0）（境界値）
   - 入力: `` logs := []entity.ActionLog{{ShiftID: 0, ActionType: "DELETE", DiffPayload: json.RawMessage(`{"deleted_task":"撤収作業"}`)}} shiftMap := map[int]entity.ShiftAdmin(nil) taskMap := map[int]entity.Task(nil) timeMap := map[int]entity.Time(nil) ``
   - 期待値: `"撤収作業（削除）"（nil マップ読み取りは Go では安全。時間プレフィックスなしで出力される）【実行確認済み: panic せず期待どおり】`
   - 根拠: nilマップでpanicしないことを固定する。あわせて、DELETEだけはshiftMap/taskMapのルックアップの成功を必要とせず、diff_payloadだけで出力されるという特別扱い（DBのshift_id NULL → ShiftID=0のケース）を固定する。
6. shiftMap / taskMapにキーが無い非DELETEログはスキップ（境界値）
   - 入力: `` logs := []entity.ActionLog{ {ShiftID: 99, ActionType: "UPDATE", DiffPayload: json.RawMessage(`{"changes":[{"old":"a","new":"b"}]}`)}, {ShiftID: 10, ActionType: "UPDATE", DiffPayload: json.RawMessage(`{"changes":[{"old":"a","new":"b"}]}`)}, } shiftMap := map[int]entity.ShiftAdmin{10: {TaskID: 100, TimeID: 5}} // 99 は無い taskMap := map[int]entity.Task{} // 100 も無い timeMap := map[int]entity.Time{5: {Time: "10:00"}, 6: {Time: "11:00"}} ``
   - 期待値: `""（1件目は shiftMap 欠落、2件目は taskMap 欠落で両方 continue され、結果は空文字）【実行確認済み】`
   - 根拠: 2段のルックアップの失敗（shiftMapの欠落・taskMapの欠落）が、どちらも警告を出さずにスキップになる今の動きを、1ケースで固定する。
7. timeMapにTimeID+1が無い（最終時間帯）→ 終了時刻0:00フォールバック（境界値）
   - 入力: `` logs := []entity.ActionLog{{ShiftID: 10, ActionType: "UPDATE", DiffPayload: json.RawMessage(`{"changes":[{"old":"受付","new":"警備"}]}`)}} shiftMap := map[int]entity.ShiftAdmin{10: {TaskID: 100, TimeID: 5}} taskMap := map[int]entity.Task{100: {Task: "警備"}} timeMap := map[int]entity.Time{5: {Time: "22:00"}} // キー 6 が無い ``
   - 期待値: `"22:00 〜 0:00：受付 → 警備"【実行確認済み】`
   - 根拠: formatTimeRangeでendTimeID+1が欠けたときにentity.Time{Time: "0:00"} へフォールバックする明示的な実装（最終スロットは翌0:00終了の想定）を固定する。時間IDが連番であることに依存した設計の目印にもなる。
8. timeMapにTimeID自体が無い → 開始時刻が空文字で出力（異常系）【要判断】
   - 入力: `` logs := []entity.ActionLog{{ShiftID: 10, ActionType: "UPDATE", DiffPayload: json.RawMessage(`{"changes":[{"old":"受付","new":"警備"}]}`)}} shiftMap := map[int]entity.ShiftAdmin{10: {TaskID: 100, TimeID: 5}} taskMap := map[int]entity.Task{100: {Task: "警備"}} timeMap := map[int]entity.Time{} // 空 ``
   - 期待値: `" 〜 0:00：受付 → 警備"（先頭が半角スペース始まり。開始時刻はゼロ値 Time の空文字がそのまま入る）【実行確認済み: 先頭スペース込みで完全一致】`
   - 根拠: formatTimeRangeはstartTimeをokチェックなしでtimeMap[startTimeID]のゼロ値のまま使う。そのため開始側が欠けると、通知文が「 〜 0:00：」のように、開始時刻が空で半角スペースから始まる形で届く。この非対称なフォールバックを今の動きとして固定しつつ、論点として挙げる。
9. diff_payloadがパースできない（不正JSON / nil）ログは、ログ出力なしでスキップ（異常系）【要判断】
   - 入力: `` logs := []entity.ActionLog{ {ID: 1, ShiftID: 10, ActionType: "UPDATE", DiffPayload: nil}, {ID: 2, ShiftID: 10, ActionType: "UPDATE", DiffPayload: json.RawMessage(`{invalid`)}, {ID: 3, ShiftID: 10, ActionType: "UPDATE", DiffPayload: json.RawMessage(`{"changes":[{"old":"受付","new":"警備"}]}`)}, } shiftMap := map[int]entity.ShiftAdmin{10: {TaskID: 100, TimeID: 5}} taskMap := map[int]entity.Task{100: {Task: "警備"}} timeMap := map[int]entity.Time{5: {Time: "10:00"}, 6: {Time: "11:00"}} ``
   - 期待値: `"10:00 〜 11:00：受付 → 警備"（ID:1 は nil で unmarshal エラー、ID:2 は不正 JSON でエラー、いずれも警告なしにスキップされ ID:3 の 1 行のみ）【実行確認済み: nil DiffPayload も panic せず同経路でスキップ】`
   - 根拠: json.Unmarshalが失敗したときのcontinueによる欠落を固定する。DiffPayloadがnilでも "unexpected end of JSON input" で同じ経路に入ることも含めて確かめる。
10. UPDATEでpayloadにchangesキーが無い → 不明→DB現在値フォールバック（異常系）
   - 入力: `` logs := []entity.ActionLog{{ShiftID: 10, ActionType: "UPDATE", DiffPayload: json.RawMessage(`{}`)}} shiftMap := map[int]entity.ShiftAdmin{10: {TaskID: 100, TimeID: 5}} taskMap := map[int]entity.Task{100: {Task: "警備"}} timeMap := map[int]entity.Time{5: {Time: "10:00"}, 6: {Time: "11:00"}} // DiffPayload が `{"changes":[]}`（空配列）でも len>0 が偽になり同一挙動 ``
   - 期待値: `` "10:00 〜 11:00：（不明） → 警備"（old は既定の「（不明）」、new は taskMap の DB 現在値 "警備" にフォールバック）【`{}` と `{"changes":[]}` の両方を別サブテストで実行し、同一出力を確認済み】 ``
   - 根拠: 実装のコメントに「フォールバック: DB現在値」と明記された、意図的な既定値のロジックの回帰を防ぐ。changesキーが無い場合と空配列の場合が同じ分岐に入る、という境界も兼ねる。

実行検証: worktreeのapi/lib/usecase/design_verify_buildchangeswithtime_test.goに、起案した10ケースを12サブテストとして実装した。「logs nil/空スライス」と「changesキー無し/空配列」は、同じ動きになるという主張を、それぞれ2サブテストに分けた。cd api && go test ./lib/usecase/... -run TestDesignVerifyBuildChangesWithTime -vを実行し、12/12がPASSして、全ケースが起案の期待値と完全に一致した（修正0件、削除0件、実行不能0件）。ゼロ値レシーバ &notificationUseCase{} でモック無しに呼べることも確かめた。開始時刻が欠けたときの先頭の半角スペース、nilマップ・nil DiffPayloadでpanicしないこと、全角チルダ/コロンの文字の種類まで、%qで比べて確かめた。needs_judgmentの2件（開始時刻のフォールバックの非対称、パースできないログをログ出力なしで捨てること）は、実際の動きが起案どおりであることを確かめたうえで、設計判断の論点として残した。実行後にテストファイルを削除し、git status --porcelainが空であることを確かめた。commit/pushは一切していない。

## convertShiftCardDataToShifts

`api/lib/usecase/shift_usecase.go:433`

DBのJOIN結果のフラットな構造体[]entity.ShiftCardDataを、モバイルAPI用のネストした構造体[]entity.Shiftへ1:1で詰め替える純関数。YearValueだけはstrconv.Atoiで文字列からintに変換し、エラーは捨てて0にフォールバックする。レシーバ *shiftUseCaseのフィールドには一切依存せず、(&shiftUseCase{})のゼロ値レシーバからモック無しで呼び出せることを実行で確かめた。マップの引数は持たないため、「nilマップ／キー欠落」の境界は当てはまらない。代わりに、nilスライス・ゼロ値の要素・Atoiの失敗を境界ケースにした。全8ケースをgo testで実行し、起案の期待値と実際の動きが一致することを確かめた。

### 要判断の論点

YearValueが数値でない文字列のとき、Atoiのエラーを_で捨てて0にフォールバックするのが意図的なものか、判断が必要である。DB由来の値なので通常は数値だが、数値でない文字列が入っていると「0年」としてAPIレスポンスに載る。エラーもログも出ないため、年が0になったことにAPIの側でもmobileの側でも気づけない。ログを出す・エラーを返すように変える選択肢もあるが、シグネチャの変更（errorの追加）は呼び出し側にも影響する。そのため、保守モードでは今の動きのまま固定するのが妥当か、確認を求める（実際の動きは起案どおりで食い違いは無い。論点は設計意図の確認だけである）。

TaskMobileにはRemark・MaxMember・BureauIDのフィールドがあるのに、ShiftCardDataの対応する値（TaskRemark/MaxMember/TaskBureauID）を詰めていないため、これらの値は出力に伝わらない。モバイル画面で要らないから意図して省いたのか、マッピング漏れなのか、判断が必要である。モバイル側（Flutter）がこれらを参照していれば実際のバグになる。テストは今の動きで固定して書くが、期待値を「詰める」に変えるかどうかは、モバイル側の参照を調べてから判断すべきである（実際の動きは起案どおりで食い違いは無い。論点は設計意図の確認だけである）。

### ケース表

1. 全フィールド設定済みの1要素の完全なマッピング（正常系）
   - 入力: `[]entity.ShiftCardData{{ShiftID: 1, UserID: 2, TaskID: 3, YearID: 4, DateID: 5, TimeID: 6, WeatherID: 7, IsAttendance: true, TaskName: "受付", TaskColor: "#FF0000", TaskURL: "https://example.com/manual", PlaceName: "体育館", TimeValue: "9:00", UserName: "山田太郎", UserBureauID: 8, UserGradeID: 9, YearValue: "2024", DateValue: "9/13", WeatherValue: "晴れ"}}`
   - 期待値: `[]entity.Shift{{ID: 1, Task: entity.TaskMobile{ID: 3, Task: "受付", Color: "#FF0000", Place: "体育館", Url: "https://example.com/manual"}, User: entity.User{ID: 2, Name: "山田太郎", BureauID: 8, GradeID: 9}, Year: entity.Year{ID: 4, Year: 2024}, Date: entity.Date{ID: 5, Date: "9/13"}, Time: entity.Time{ID: 6, Time: "9:00"}, Weather: entity.Weather{ID: 7, Weather: "晴れ"}, IsAttendance: true}} と reflect.DeepEqual で一致（CreatedAt/UpdatedAt は両辺ゼロ値 time.Time のため一致する）【実行で確認済み】`
   - 根拠: 各フィールドの対応関係（ShiftID→ID、PlaceName→Task.Place、UserBureauID→User.BureauIDなど）を取り違える回帰を防ぐ、基準のケース。IsAttendance=trueの引き継ぎも兼ねる
2. 複数要素で入力順が保たれる（正常系）
   - 入力: `[]entity.ShiftCardData{{ShiftID: 10, TimeID: 3, YearValue: "2024"}, {ShiftID: 20, TimeID: 1, YearValue: "2024"}, {ShiftID: 30, TimeID: 2, YearValue: "2024"}}（TimeID を昇順にしないことでソートされないことを検証）`
   - 期待値: `len == 3 かつ shifts[0].ID == 10, shifts[1].ID == 20, shifts[2].ID == 30（入力順のまま。ソートは呼び出し側 GetShiftCardsByUserAndDateAndWeather の責務であり本関数では行わない）【実行で確認済み】`
   - 根拠: 本関数に将来ソートやフィルタが入り込む回帰を防ぐ。変換だけを受け持つことを仕様として固定する
3. 空スライスの入力はnilを返す（境界値）
   - 入力: `[]entity.ShiftCardData{}`
   - 期待値: `戻り値は nil（var shifts []entity.Shift に append されないため）。len == 0 かつ shifts == nil であり、reflect.DeepEqual(nil戻り値, []entity.Shift{}) は false になることまで実行で確認済み。assert では len == 0 に加えて shifts == nil を固定する`
   - 根拠: nilスライスを返す今の動きを固定する。呼び出し側はrangeしかしないため実際の影響はないが、将来json化などでnilと空の差が問題になったときに検知できる
4. nilスライスの入力はnilを返す（境界値）
   - 入力: `var data []entity.ShiftCardData = nil を渡す`
   - 期待値: `panic せず nil を返す（range over nil slice は 0 回イテレーションで安全）【実行で確認済み】`
   - 根拠: リポジトリがエラー無しでnilを返した場合の安全性を固定する。nil入力でpanicする回帰を防ぐ
5. ゼロ値の要素1個の変換（境界値）
   - 入力: `[]entity.ShiftCardData{{}}（全フィールドゼロ値。YearValue は ""）`
   - 期待値: `len == 1 かつ shifts[0] が entity.Shift{} と reflect.DeepEqual で一致（strconv.Atoi("") はエラーだが黙殺され Year.Year == 0）【実行で確認済み】`
   - 根拠: 要素1個という最小の非空入力と、全フィールドがゼロ値でもpanicしないことを同時に固定する
6. YearValueが数値でない文字列（異常系）【要判断】
   - 入力: `[]entity.ShiftCardData{{ShiftID: 1, YearValue: "abc"}, {ShiftID: 2, YearValue: " 2024"}}（先頭空白は Atoi が受理しない）の2要素`
   - 期待値: `いずれも Year.Year == 0 になり、エラーもログも発生せず正常に変換が完了する（現状挙動）【実行で確認済み。" 2024" も Atoi エラーで 0 になることを確認】`
   - 根拠: strconv.Atoiのエラーを捨てる今の動きを固定する。DBのyear_valueが数値でない文字列のときに、エラーもログも出さずにyear=0を返す経路を検知する
7. YearValueが負数の文字列（境界値）
   - 入力: `[]entity.ShiftCardData{{YearValue: "-1"}}`
   - 期待値: `Year.Year == -1（strconv.Atoi は符号付きをそのまま受理するため負値が通過する）【実行で確認済み】`
   - 根拠: Atoiが受け付ける範囲（符号付き整数）をそのまま通す今の動きを固定する。バリデーションを足したときに検知できる
8. マッピングされていないフィールドは出力に伝わらない（異常系）【要判断】
   - 入力: `[]entity.ShiftCardData{{TaskID: 3, TaskName: "受付", TaskRemark: "注意事項あり", MaxMember: 10, TaskBureauID: 5, PlaceID: 7, YearValue: "2024"}}`
   - 期待値: `shifts[0].Task.Remark == "", shifts[0].Task.MaxMember == 0, shifts[0].Task.BureauID == 0, shifts[0].Task.YearID == 0（TaskRemark/MaxMember/TaskBureauID/PlaceID は変換で捨てられる。TaskMobile に対応フィールドがあるのに詰め替えていない）【実行で確認済み】`
   - 根拠: 入力にあって出力に無いフィールドの扱いを明文化する。将来マッピングを足したり削ったりしたとき、必ずこのテストに引っかかる

実行検証: worktreeのapi/lib/usecase/design_verify_convertshiftcarddatatoshifts_test.goに全8ケースをテスト関数8本として実装し、cd api && go test ./lib/usecase/ -run TestDesignVerifyConvertShiftCardDataToShifts -vで実行した。8ケースすべてがPASSし、起案どおりが8件、実際の動きとの食い違いによるexpectedの修正は0件、実行できずに削除・修正したものも0件だった。ゼロ値レシーバ(&shiftUseCase{})からモック無しで呼び出せることも、実行で確かめた。空スライスを入力するケースでは、戻り値がnilであることに加え、reflect.DeepEqual(nil, []entity.Shift{})がfalseになる点（起案の注記）もassertで確かめた。needs_judgment=trueの2件（Atoiのエラーを捨てること、TaskRemark/MaxMember/TaskBureauIDがマッピングされていないこと）は、実際の動きと期待値の食い違いではない。今の動きを仕様として固定してよいかという設計意図の確認として、そのまま残した。検証後にテストファイルは削除し、git status --porcelainが空であることを確かめた。commit/pushは行っていない。

## groupContinuousShifts

`api/lib/usecase/shift_usecase.go:478`

ソート済みの同じタスクのシフト列を受け取り、[][]entity.Shiftを返す関数（実装L478-506）。隣り合う要素が「同じTask.IDかつTime.IDがちょうど +1」のときは同じグループにつなげ、それ以外のときは新しいグループを始める。空入力（nilを含む）は非nilの空スライスを返す（L479-481の早期リターン）。レシーバa *shiftUseCaseのフィールドには一切依存しない純関数であり、同じパッケージのテストから(&shiftUseCase{}).groupContinuousShifts(...)とゼロ値レシーバで呼べることを実行で確かめた。entityはgithub.com/NUTFes/SeeFT/api/lib/entity（Shift.Taskはentity.TaskMobile、Shift.Timeはentity.Timeで、どちらもID int）。全10ケースを実際にgo testで実行し、起案の期待値と実際の動きの食い違いはゼロだった。

### 要判断の論点

同じTimeIDが重複すると別のグループに分かれる点は、実行して、今の動きが起案どおり「分割」であることを確かめた。残る論点は設計意図で、同じ時刻のシフトを別のカードに分けるのが意図どおりなのか、判断が必要である。呼び出し元はuserIDで絞り込み済みのため、通常は同じTimeIDが重複しない。しかし、もし「同じ時刻・同じタスクは同じグループにまとめる（curr.Time.ID == prev.Time.IDも連続として扱う）」が本来の意図なら、今の動きはバグの疑いがある。保守モードの方針では、今の動き（分割）を正解としてテストを書き、変えるなら別のissueにするのが妥当と考える。ただし期待値を確定するには、人間の判断が必要である。

### ケース表

1. 連続する3件が1グループになる（正常系）
   - 入力: `ヘルパー sh(task, time int) entity.Shift { return entity.Shift{Task: entity.TaskMobile{ID: task}, Time: entity.Time{ID: time}} } を定義し、[]entity.Shift{sh(1,1), sh(1,2), sh(1,3)} を渡す（以降のケースも同じ sh を使用）`
   - 期待値: `reflect.DeepEqual で [][]entity.Shift{{sh(1,1), sh(1,2), sh(1,3)}} と一致（グループ数1、要素順は入力順のまま）。実行で確認済み`
   - 根拠: いちばん基本の連結の経路（curr.Time.ID == prev.Time.ID+1が真の側）の回帰を防ぐ。グループ内の要素の順序が保たれることも固定する
2. TimeIDの欠番でグループが分かれる（正常系）
   - 入力: `[]entity.Shift{sh(1,1), sh(1,2), sh(1,4), sh(1,5)}`
   - 期待値: `reflect.DeepEqual で [][]entity.Shift{{sh(1,1), sh(1,2)}, {sh(1,4), sh(1,5)}} と一致（TimeID 2→4 のギャップで分割、グループ数2）。実行で確認済み`
   - 根拠: 分割の経路（連続条件が偽の側 → currentGroupを確定 → 新しいグループを開始）と、最後のグループを漏らさないこと（ループ後のappend、L501-503）を同時に検証する
3. 空スライスは非nilの空の結果（境界値）
   - 入力: `[]entity.Shift{}`
   - 期待値: `戻り値 got について got != nil かつ len(got) == 0（reflect.DeepEqual(got, [][]entity.Shift{}) が true）。実行で確認済み`
   - 根拠: L479-481の早期リターンがnilではなく空スライスを返す仕様を固定する。呼び出し元のrangeが何もせずに安全に終わること、およびJSON化したときにnullにならない性質の回帰を防ぐ
4. nilスライスも空スライスとして扱う（境界値）
   - 入力: `var in []entity.Shift = nil として groupContinuousShifts(in) を呼ぶ`
   - 期待値: `空スライスと同じく got != nil かつ len(got) == 0。実行で確認済み（panic せず非 nil 空スライスが返る）`
   - 根拠: Goではlen(nil) == 0なので空スライスと同じ経路を通るが、nil入力でpanicしないこと、出力がnilにならないことを明示的に固定する
5. 要素1個は単独のグループ（境界値）
   - 入力: `[]entity.Shift{sh(1,5)}`
   - 期待値: `reflect.DeepEqual で [][]entity.Shift{{sh(1,5)}} と一致（グループ数1、要素数1）。実行で確認済み`
   - 根拠: ループ本体（i=1以降）を一度も通らず、初めのcurrentGroupがそのまま最後にappendされる最小の経路の回帰を防ぐ
6. 同一TimeIDの重複は別グループに分割される（境界値）【要判断】
   - 入力: `[]entity.Shift{sh(1,3), sh(1,3), sh(1,4)}`
   - 期待値: `現状挙動: reflect.DeepEqual で [][]entity.Shift{{sh(1,3)}, {sh(1,3), sh(1,4)}} と一致（1件目の sh(1,3) は単独グループ、2件目の sh(1,3) から sh(1,4) が連結）。実行で確認済み`
   - 根拠: 連続の条件が「ちょうど +1」であり、同じ値（+0）を連続とみなさないことを固定する。同じタスク・同じ時刻に複数のシフトがあるデータが来た場合の動きを明示する
7. TaskIDが異なる隣り合うシフトは、TimeIDが連続でも分割（異常系）
   - 入力: `[]entity.Shift{sh(1,1), sh(2,2)}（TimeID は 1→2 で +1 連続だが TaskID が 1 と 2 で異なる）`
   - 期待値: `reflect.DeepEqual で [][]entity.Shift{{sh(1,1)}, {sh(2,2)}} と一致（グループ数2）。実行で確認済み`
   - 根拠: 呼び出し元はtaskIDごとに分けてから渡すため、通常は混ざらない。その不変条件が破れた入力でも、タスクをまたいで誤ってつながないガード（prev.Task.ID == curr.Task.IDの条件、L491）の回帰を防ぐ
8. 未ソート入力（降順）はマージされない（異常系）
   - 入力: `[]entity.Shift{sh(1,3), sh(1,2), sh(1,1)}`
   - 期待値: `reflect.DeepEqual で [][]entity.Shift{{sh(1,3)}, {sh(1,2)}, {sh(1,1)}} と一致（全要素が単独グループ、入力順保存、内部でソートされない）。実行で確認済み`
   - 根拠: 本関数は隣り合う要素を比べるだけで、内部ではソートしない。呼び出し元のsort.Sort(ByTime(taskShifts))が事前条件であることを、テストとして文書にする。将来、内部にソートを足す変更が入った場合に検知できる
9. ゼロ値の要素2個は別グループ（境界値）
   - 入力: `[]entity.Shift{{}, {}}（両方とも Task.ID=0, Time.ID=0 のゼロ値）`
   - 期待値: `reflect.DeepEqual で [][]entity.Shift{{{}}, {{}}} と一致（0 == 0+1 が偽のため2グループ）。実行で確認済み`
   - 根拠: ゼロ値のShiftが来てもpanicせず、同じTimeID(0)の重複として分割される動きを固定する。変換層の不具合でゼロ値が入り込んだ場合の振る舞いを明示する
10. ゼロ値とTimeID=1は連続として扱われ1グループ（境界値）
   - 入力: `[]entity.Shift{{}, sh(0,1)}（1件目はゼロ値: Task.ID=0, Time.ID=0。2件目は Task.ID=0, Time.ID=1）`
   - 期待値: `reflect.DeepEqual で [][]entity.Shift{{{}, sh(0,1)}} と一致（Task.ID 0==0 かつ Time.ID 1 == 0+1 で連結、グループ数1）。実行で確認済み`
   - 根拠: TimeID=0を起点にしても +1の判定が機械的に成り立つことを固定する。IDにドメイン上の下限チェックが無い（0や負値も通る）ことを、テストとして文書にする

実行検証: worktreeのapi/lib/usecase/design_verify_groupcontinuousshifts_test.go（package usecase）に全10ケースを実装し、cd api && go test ./lib/usecase/... -run TestDesignVerifyGroupContinuousShifts -vで実行した。10ケースすべてが1回目でPASSした（起案どおり10件、実際の動きとの食い違いによる修正0件、実行できずに削除したもの0件）。unexportedの関数は、ゼロ値レシーバ(&shiftUseCase{})から問題なく呼べた。entityのパスはapi/lib/entity（import: github.com/NUTFes/SeeFT/api/lib/entity）で、entity.Shift{Task: entity.TaskMobile{ID: ...}, Time: entity.Time{ID: ...}} のリテラルがそのまま組み立てられることを確かめた。needs_judgmentの1件（同一TimeIDの重複）は、実際の動きが起案の「現状挙動: 分割」と一致したため期待値はそのままにし、設計意図の判断だけを人間に残す。検証後にテストファイルは削除し、git status --porcelainが空であることを確かめた。commit/pushは行っていない。

## compareTimeStrings

`api/lib/usecase/shift_usecase.go:509`

"8:00" のような時刻文字列2つを、時と分から総分に換算して数値で比べ、-1/0/1を返す比較関数。sort.Slice（L426）でShiftCard.StartTimeを並べるのに使われ、辞書順で比べると "10:00" < "8:00" になる問題を避けるためにある。レシーバaのフィールドには一切依存せず（aは本体で使われていない）、entity型も使わない純関数なので、&shiftUseCase{} のゼロ値レシーバでテストできる。引数がstringだけなのでスライスやマップの境界は当てはまらず、空文字列・ゼロ値の時刻・不正な書式に読み替えて設計した。全10ケースを実際のコードで実行して確かめた。

### 要判断の論点

空文字列が正当な時刻とも等しいと扱われる点は、実行して確かめた（実際の動きも0）。片方が不正であれば、もう片方が正当な時刻でも一律に0（等しい）を返す。sort.Sliceは安定ソートではないため、不正なStartTimeを持つShiftCardの並び順が、実行のたびに変わり得る。これを仕様（不正な入力は順序を問わない）とみなすか、エラーを返す設計に改めるべきかは、人間の判断が必要である。フェーズ1では今の0を期待値とする。

時と分が数値でないと暗黙のうちに0:00として扱われる点も、実行して確かめた（実際の動きも -1）。Atoiのエラーが捨てられるため、"aa:bb" が0:00（深夜0時）として扱われ、どの正当な時刻よりも前に並ぶ。ガード節の「不正なら等価(0)」という方針とも合わず（不正の種類によって動きが変わる）、バグの疑いがある。エラーのときに0を返すか、パースエラーを呼び出し元へ伝えるかは、人間の判断が必要である。フェーズ1では今の -1を期待値とする。

### ケース表

1. 桁数が異なる時刻の数値比較（正常系）
   - 入力: `time1: "8:00", time2: "10:00"`
   - 期待値: `-1`
   - 根拠: 辞書順で比べると "10:00" < "8:00" となり、順序が逆になる。この関数がある理由そのものであり、文字列比較に安易に書き換えたときの回帰を防ぐ、最も重要なケース
2. 逆順で正の値を返す（正常系）
   - 入力: `time1: "10:00", time2: "8:00"`
   - 期待値: `1`
   - 根拠: 比較関数の対称性（引数を入れ替えると符号が反転する）を固定する。sort.Sliceの比較子として順序が安定する前提を守る
3. 同一時刻は等価（正常系）
   - 入力: `time1: "9:15", time2: "9:15"`
   - 期待値: `0`
   - 根拠: 反射律（同じ値を入れると0）を固定する。等価を判定する分岐（L531）をカバーする
4. 同時間帯での分単位の差（正常系）
   - 入力: `time1: "9:15", time2: "9:30"`
   - 期待値: `-1`
   - 根拠: 時が同じで分だけ異なる場合の比較。h*60+mの分換算の式のうち、mの項が結果に反映されていることを確かめる（mを無視するように退行するのを防ぐ）
5. ゼロ値時刻と一日の最大時刻（境界値）
   - 入力: `time1: "0:00", time2: "23:59"`
   - 期待値: `-1`
   - 根拠: 総分に換算したときの下限（0分）と上限（1439分）の境界。ゼロ値の時刻が最小として正しく並ぶことを固定する
6. ゼロパディング表記の同値性（境界値）
   - 入力: `time1: "08:05", time2: "8:05"`
   - 期待値: `0`
   - 根拠: strconv.Atoiが先頭のゼロを受け付けるため、"08" と "8" は同じ値になる。DBや入力に由来する表記ゆれがあっても、比較結果が変わらないことを固定する
7. 非正規化の分表記は換算後に等価（境界値）
   - 入力: `time1: "1:30", time2: "0:90"`
   - 期待値: `0`
   - 根拠: 比較が文字列の一致ではなく、h*60+mの総分換算で行われるという内部の仕様を固定する。分が60以上でも換算されて等価になる
8. 空文字列は正当な時刻とも等価扱い（異常系）【要判断】
   - 入力: `time1: "", time2: "10:00"`
   - 期待値: `0`
   - 根拠: 空文字列はstrings.Splitで要素数1となり、ガード節（L514-516）に入る。不正な入力のときのフォールバックの動きを固定する
9. コロンが2個以上の書式は等価扱い（異常系）
   - 入力: `time1: "10:00:00", time2: "9:00"`
   - 期待値: `0`
   - 根拠: 秒付きの書式（要素数3）もガード節で0になる経路をカバーし、要素数が多すぎる側の分岐を固定する。この動きの是非の論点は空文字列のケースと同じため、そちらにまとめた
10. 非数値の時分は暗黙に0:00扱い（異常系）【要判断】
   - 入力: `time1: "aa:bb", time2: "8:00"`
   - 期待値: `-1`
   - 根拠: 書式（コロン1個）は通るがAtoiが失敗する入力の経路。L518-521でエラーが_で捨てられ、h=0, m=0となる今の動きを固定する

実行検証: worktreeのapi/lib/usecase/design_verify_comparetimestrings_test.goに全10ケースをテーブル駆動で実装し、go test ./lib/usecase/ -run TestDesignVerifyCompareTimeStrings -vを実行した。10ケースすべてがPASSし、起案どおりが10件、実際の動きとの食い違いによるexpectedの修正は0件、実行できずに削除・修正したケースも0件だった。needs_judgment=trueの2件（空文字列→0、"aa:bb"→-1）は、期待値そのものは実際の動きと一致している。設計判断の論点（不正な入力の扱い）として、judgment_noteに「実行確認済み」を書き足した。検証後にテストファイルを削除し、git status --porcelainが空であることを確かめた。commit/pushは一切していない。

## BuildMessageBlocks

`api/lib/externals/slack/slack_service.go:88`

Slack通知用のBlock Kitメッセージ（ヘッダー、基本情報のセクション、任意の変更内容のセクション、区切り線）を組み立てる純粋な構築関数。レシーバ *SlackServiceのフィールド（client, channelID）には一切依存せず、(&SlackService{}).BuildMessageBlocks(params)で環境変数なしに呼び出せることを実行で確かめた。入力はMessageParams（stringの5フィールドだけ）で、全8ケースを、reflect.DeepEqualでブロック構造全体を比べて確かめた。起案の期待値と実際の動きは、全ケースで一致した。

### 要判断の論点

Changesが空白だけでも、「*変更内容*」の見出しだけで中身のないセクションが表示される。これは直感に反する可能性がある。判定をstrings.TrimSpace(params.Changes) != "" にすべきかどうかは、呼び出し元のbuildGroupedMessageが空白だけの文字列を返しうるかによるため、人間の判断が必要である。実行した結果、今はそのまま追加される動きであることを確かめた。

TitleがSlackのヘッダーの上限150文字を超えると、そのブロックはSlack APIへの送信時にinvalid_blocksエラーになりうる。切り詰める責務をこの関数に持たせるか、呼び出し元の入力制約にするかは設計判断である。今は唯一の呼び出し元がTitleを固定の文字列「シフト変更通知」で渡すため、実際の影響はない。実行した結果、切り詰めずにそのまま生成する動きであることを確かめた。

Slackのmrkdwnの仕様では & < >のエスケープが推奨されており、UserNameにこれらの特殊文字やメンション構文が入ると、verbatim=falseのため<!channel>などがメンションとして解釈されうる（mrkdwnインジェクション）。DMだけで運用しているため影響は本人宛てに限られるが、エスケープを入れるべきかは人間の判断が必要である。実行した結果、今は加工せずに埋め込む動きであることを確かめた。

### ケース表

1. 全フィールド設定（変更内容あり）（正常系）
   - 入力: `slack.MessageParams{Title: "シフト変更通知", UserName: "山田太郎", Date: "9月13日(土)", Weather: "晴れ", Changes: "・10:00-12:00 受付 → 会場整理"}`
   - 期待値: `長さ4の []slack.Block。[0] slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", "🔔 シフト変更通知", false, false))、[1] slack.NewSectionBlock(nil, []*slack.TextBlockObject{slack.NewTextBlockObject("mrkdwn", "ユーザー: 山田太郎", false, false), slack.NewTextBlockObject("mrkdwn", "日付: 9月13日(土)", false, false), slack.NewTextBlockObject("mrkdwn", "天気: 晴れ", false, false)}, nil)、[2] slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*変更内容*\n・10:00-12:00 受付 → 会場整理", false, false), nil, nil)、[3] slack.NewDividerBlock()。reflect.DeepEqual で全体比較可能（実行で一致確認済み）`
   - 根拠: 唯一の呼び出し元notification_usecase.go L300と同じ形の代表的な入力。4ブロックの構成と、各テキストの整形フォーマット（プレフィックス「🔔 」「ユーザー: 」「*変更内容*\n」など）の回帰を防ぐ、中心となるケース
2. Changesが空文字列（変更内容ブロック省略）（正常系）
   - 入力: `slack.MessageParams{Title: "シフト変更通知", UserName: "山田太郎", Date: "9月13日(土)", Weather: "曇り", Changes: ""}`
   - 期待値: `長さ3の []slack.Block。[0] HeaderBlock("🔔 シフト変更通知")、[1] 3フィールドのSectionBlock（ケース1と同形式）、[2] slack.NewDividerBlock()。「*変更内容*」を含むSectionBlockは存在しない（実行で一致確認済み）`
   - 根拠: 関数内で唯一の分岐if params.Changes != "" の偽の側を固定する。変更内容がない通知でブロックが3個になる（区切り線は残る）今の仕様の回帰を防ぐ
3. ゼロ値MessageParams{}（全フィールド空）（境界値）
   - 入力: `slack.MessageParams{}`
   - 期待値: `長さ3の []slack.Block。ヘッダーtextは "🔔 "（ベル絵文字＋半角スペースのみ）、フィールドは順に "ユーザー: ", "日付: ", "天気: "（いずれも末尾半角スペース付きでラベルのみ）、末尾は DividerBlock。panicしない（実行で一致確認済み）`
   - 根拠: 構造体がゼロ値でも安全にブロックを組み立てられることの検証。ラベル文字列とプレースホルダの連結の仕様（trimしない・省略しない）を固定する
4. Changesが空白のみ（境界値）【要判断】
   - 入力: `slack.MessageParams{Title: "通知", UserName: "u", Date: "d", Weather: "w", Changes: " "}`
   - 期待値: `長さ4の []slack.Block。[2] は slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*変更内容*\n ", false, false), nil, nil)。空白のみでも変更内容ブロックが追加される（実行で一致確認済み）`
   - 根拠: 分岐の条件がparams.Changes != "" という完全一致の判定であることを固定する。TrimSpaceを使うなどの変更が入ったときに検知できる
5. Changesが複数行（実運用のグループ化メッセージ形式）（境界値）
   - 入力: `slack.MessageParams{Title: "シフト変更通知", UserName: "山田太郎", Date: "9月13日(土)", Weather: "晴れ", Changes: "【追加】\n・9:00-10:00 設営\n【削除】\n・13:00-14:00 受付"}`
   - 期待値: `長さ4の []slack.Block。[2] のtextは "*変更内容*\n【追加】\n・9:00-10:00 設営\n【削除】\n・13:00-14:00 受付"。改行はエスケープ・加工されずそのまま保持される（実行で一致確認済み）`
   - 根拠: 実運用ではbuildGroupedMessageが複数行の文字列を渡す。見出し "*変更内容*" と本文が \n 1個でつながり、本文の改行がそのまま通ることの回帰を防ぐ
6. TitleがSlackヘッダー上限150文字超（境界値）【要判断】
   - 入力: `slack.MessageParams{Title: strings.Repeat("あ", 151), UserName: "u", Date: "d", Weather: "w", Changes: ""}`
   - 期待値: `長さ3の []slack.Block。ヘッダーtextは "🔔 " + strings.Repeat("あ", 151)（切り詰め・バリデーションなしでそのまま格納される）（実行で一致確認済み）`
   - 根拠: Slack APIのheader blockのplain_textは150文字が上限。本関数は長さの検証も切り詰めもしないという今の動きを明示的に固定し、暗黙のtruncateの追加を検知する
7. UserNameにmrkdwn特殊文字・メンション構文（異常系）【要判断】
   - 入力: `slack.MessageParams{Title: "通知", UserName: "<@U12345> & *bold* <!channel>", Date: "d", Weather: "w", Changes: ""}`
   - 期待値: `長さ3の []slack.Block。[1] の第1フィールドtextは "ユーザー: <@U12345> & *bold* <!channel>" で、&, <, > はエスケープされずそのまま格納される（verbatim=false）（実行で一致確認済み）`
   - 根拠: UserNameはDB由来のユーザー入力（user.Name）。エスケープ処理が無いという今の動きを固定し、将来エスケープを足したときに、テストを意図して更新させる
8. 制御文字・NULバイトを含む入力（異常系）
   - 入力: `slack.MessageParams{Title: "a\x00b", UserName: "tab\tsep", Date: "line1\nline2", Weather: "\r", Changes: "end\x1b[0m"}`
   - 期待値: `panicせず長さ4の []slack.Block を返す。ヘッダーtextは "🔔 a\x00b"、フィールドは "ユーザー: tab\tsep", "日付: line1\nline2", "天気: \r"、[2] のtextは "*変更内容*\nend\x1b[0m"。全て入力文字列がそのまま連結される（実行で一致確認済み）`
   - 根拠: エラーの戻り値がない関数の異常系として、どんなバイト列でもpanicせずにブロックを返すこと（サニタイズせずにそのまま通す）を保証する。fmt.Sprintfの連結だけで、例外の経路がないことを確かめる

実行検証: worktreeのapi/lib/externals/slack/design_verify_buildmessageblocks_test.goに全8ケースを実装し、cd api && go test ./lib/externals/slack/ -run TestDesignVerifyBuildMessageBlocks -vで実行した。8ケースすべてがPASSした（起案どおり8件、実際の動きとの食い違いによる修正0件、実行できずに削除したもの0件）。比較では、各ケースで期待するブロック列をslack-goのコンストラクタ（NewHeaderBlock/NewSectionBlock/NewDividerBlock/NewTextBlockObject）で組み立て、reflect.DeepEqualで構造全体が一致するかを見た。そのため、ブロック数だけでなくtype/text/verbatimフラグまで完全に一致することを確かめた。ゼロ値レシーバ(&SlackService{})からの呼び出しも、環境変数なしで問題なく動いた。needs_judgment=trueの3件（空白のみのChanges・150文字を超えるTitle・mrkdwnの特殊文字）は、起案どおりの実際の動きを確かめたうえで、設計判断の論点として残した。検証後にテストファイルを削除し、git status --porcelainが空であることを確かめた。commit/pushは行っていない。

