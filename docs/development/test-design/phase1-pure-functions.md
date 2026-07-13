# フェーズ1 テスト設計書: 依存ゼロの純関数（第1弾）

[テストロードマップ](../test-roadmap.md)の「テスト設計ステージ」のサイクルを実際に回して作成した設計書。対象はフェーズ1 の純関数 8 つ。各ケース表は AI が起案したのち、**全ケースを使い捨てテストとして実際に実行し、期待値を実挙動と突き合わせて裏取り済み**。つまりこの表の期待値は推測ではなく、現時点の develop の実際の出力である。

レビューで見てほしい点は2つ。

1. ケースの過不足: 守りたい入力パターンが漏れていないか
2. 「要判断」の裁定: 現状挙動をこのまま仕様として固定してよいか、バグとして直すか。直す場合はテスト実装より先に修正 issue を切る

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

[]entity.ActionLog を受け取り、各ログの UserID と DateID から "%d_%d" 形式の文字列キーを生成して map[string][]entity.ActionLog にグルーピングする純関数（L138-145）。同一キー内の要素順は入力スライスの出現順が append によりそのまま保持される。レシーバ n のフィールドには一切依存しないため、テストでは &notificationUseCase{}（ゼロ値）で呼び出し可能であることを実行で確認済み。なお呼び出し元 ProcessUnsentNotifications（L105-117）がこのキーを strings.Split + Atoi で再パースしており、キー書式 "%d_%d" は事実上の契約なのでテストで固定する価値が高い。

### ケース表

1. **複数ユーザー・複数日付の混在を正しく分配する**（正常系）
   - 入力: `logs := []entity.ActionLog{{ID: 1, UserID: 1, DateID: 10}, {ID: 2, UserID: 2, DateID: 10}, {ID: 3, UserID: 1, DateID: 10}, {ID: 4, UserID: 1, DateID: 11}}`
   - 期待値: `map[string][]entity.ActionLog{"1_10": {{ID: 1, UserID: 1, DateID: 10}, {ID: 3, UserID: 1, DateID: 10}}, "2_10": {{ID: 2, UserID: 2, DateID: 10}}, "1_11": {{ID: 4, UserID: 1, DateID: 11}}}（reflect.DeepEqual で比較可能。キー数は 3）`
   - 根拠: 関数の中核挙動（UserID×DateID の組ごとの分配）の回帰防止。同一ユーザーでも DateID が違えば別グループ、同一日付でも UserID が違えば別グループになることを同時に検証する
2. **同一グループ内で入力順が保持される**（正常系）
   - 入力: `logs := []entity.ActionLog{{ID: 5, UserID: 7, DateID: 3}, {ID: 2, UserID: 7, DateID: 3}, {ID: 9, UserID: 7, DateID: 3}}`
   - 期待値: `map[string][]entity.ActionLog{"7_3": {{ID: 5, UserID: 7, DateID: 3}, {ID: 2, UserID: 7, DateID: 3}, {ID: 9, UserID: 7, DateID: 3}}}（ID の並びが入力順 5, 2, 9 のまま。ID 昇順にソートされない）`
   - 根拠: append による挿入順保持の回帰防止。後段の sortLogsByTime が時刻順ソートを担う分業になっており、この関数が勝手にソートしない現状を固定する
3. **空スライスなら非nilの空mapを返す**（境界値）
   - 入力: `logs := []entity.ActionLog{}`
   - 期待値: `got != nil かつ len(got) == 0（make で初期化された空 map。nil map ではない）`
   - 根拠: 空入力での挙動固定。呼び出し元は len(logs)==0 でガードしているが、関数単体として非nil空mapを返す契約を固定し、将来の直呼びでの nil map range/代入事故を防ぐ
4. **nilスライスでも非nilの空mapを返す**（境界値）
   - 入力: `var logs []entity.ActionLog（nil スライスのまま渡す）`
   - 期待値: `got != nil かつ len(got) == 0（空スライスと同一の結果。panic しない）`
   - 根拠: Go では nil スライスの range は安全だが、その事実に依存している現状挙動を明示的に固定する。nil と empty で挙動が分かれないことの回帰防止
5. **要素1個は単一キー・単一要素になり全フィールドが素通しされる**（境界値）
   - 入力: `` logs := []entity.ActionLog{{ID: 1, ShiftID: 55, UserID: 42, DateID: 7, ActionType: "UPDATE", DiffPayload: json.RawMessage(`{"changes":[]}`), IsSent: false, CreatedAt: time.Date(2026, 7, 8, 12, 0, 0, 0, time.UTC)}} ``
   - 期待値: `len(got) == 1、got["42_7"] が長さ1のスライスで、その唯一の要素が入力の entity.ActionLog と reflect.DeepEqual で完全一致（ShiftID・ActionType・DiffPayload・CreatedAt 等が改変されずに素通しされる）`
   - 根拠: 最小構成の正常動作と、UserID/DateID 以外のフィールドを一切書き換えないパススルー性の回帰防止。後段の processGroup は DiffPayload や ShiftID をそのまま使うため改変されると通知内容が壊れる
6. **ゼロ値のActionLogはキー "0_0" にグルーピングされる**（境界値）
   - 入力: `logs := []entity.ActionLog{{}}（全フィールドゼロ値。UserID=0, DateID=0, DiffPayload=nil）`
   - 期待値: `len(got) == 1、got["0_0"] == []entity.ActionLog{{}}（panic せず、nil の DiffPayload もそのまま保持）`
   - 根拠: ゼロ値入力での安全性の固定。DB スキャン異常等で UserID/DateID が 0 のまま流れてきても本関数は落ちず "0_0" に集約する、という現状の防波堤なしの挙動を明示する
7. **キー書式 "%d_%d" により (1,23) と (12,3) が衝突しない**（境界値）
   - 入力: `logs := []entity.ActionLog{{ID: 1, UserID: 1, DateID: 23}, {ID: 2, UserID: 12, DateID: 3}}`
   - 期待値: `len(got) == 2、キー "1_23" と "12_3" がそれぞれ存在し、各スライスの長さは 1（1つのキーに混ざらない）`
   - 根拠: キー書式そのものの回帰防止。区切り文字 "_" を外して "%d%d" 等に変えると "123" に衝突し別人の通知が混ざる。さらに呼び出し元（L108-117）が "_" で Split して Atoi する契約に依存しているため、書式変更は即座に検知したい
8. **負のIDでも検証なしでそのままキー化される**（異常系）
   - 入力: `logs := []entity.ActionLog{{ID: 1, UserID: -1, DateID: -2}}`
   - 期待値: `len(got) == 1、got["-1_-2"] が長さ1のスライスとして存在する（エラーにも panic にもならず、負値がそのままキーに埋め込まれる）`
   - 根拠: 不正なドメイン値に対する現状のノーバリデーション挙動の固定。なお生成キー "-1_-2" はアンダースコアが1個のため呼び出し元の Split/Atoi（L108-117）でも正しく -1, -2 に復元でき、下流も壊れないことを確認済み

実行検証: worktree の api/lib/usecase/ に使い捨てテスト design_verify_groupnotificationsbyuseranddate_test.go を作成し、(&notificationUseCase{}) ゼロ値レシーバから GroupNotificationsByUserAndDate を直接呼び出して全8ケースをサブテストとして実装。cd api && go test ./lib/usecase/... -run TestDesignVerifyGroupNotificationsByUserAndDate -v で実行し、8ケース全て初回 PASS（8/8 起案どおり、修正 0件、削除 0件、needs_judgment 0件）。非nil空map（空/nil スライス両方）、挿入順保持、フィールド素通し（DiffPayload/CreatedAt 含む DeepEqual 完全一致）、キー "0_0"・"1_23"/"12_3" 非衝突・"-1_-2" の生成も全て実挙動で裏取り済み。テストファイルは削除し、git status --porcelain が空（clean）であることを確認。commit/push は行っていない。

## sortLogsByTime

`api/lib/usecase/notification_usecase.go:317`

ActionLog のスライスを、各ログの ShiftID をキーに shiftMap から引いた ShiftAdmin.TimeID の昇順で並べ替えて新しいスライスとして返す。shiftMap にキーが無いログは結果から黙って除外されるため、ソートとフィルタを兼ねる。レシーバ n のフィールドには一切依存しない純関数であり、テストではゼロ値レシーバ &notificationUseCase{} で呼び出せる（参照フィールドは ActionLog.ShiftID と ShiftAdmin.TimeID のみ）。全9ケースを実コード実行で裏取り済み。

### 要判断の論点

- **shiftMapにキーが無いログは黙って除外される**: 実行確認済み: 実挙動は起案どおり（黙って除外、ID 列 [3, 1]）。論点は挙動の食い違いではなく設計判断。シフトが削除済み等で shiftMap に載らないログは通知から無言で欠落し、エラーにもログにも出ない。現状挙動を仕様として固定するか、除外時に警告ログを出す／件数を返す等に改善するか判断が必要。テスト自体は現状挙動（黙って除外）を期待値とする
- **同一TimeIDのタイは全件保持される（相対順序は未保証）**: 実行確認済み: 今回の実行（go 環境の現行バージョン、3要素入力）では観測順序は挿入順 [1, 2, 3] だったが、sort.Slice は非安定ソートで同一 TimeID 内の相対順序は仕様上未定義（Go バージョンや入力サイズで変わり得る）。テストは集合一致（順序不問）でパスする。同時刻グループ内の通知行の並びを安定させたいなら sort.SliceStable への変更（＋順序をアサートするテスト）を検討すべき。順序不問テストで現状を許容するか、SliceStable 化して順序を仕様化するか判断が必要

### ケース表

1. **複数ログをTimeID昇順に並べ替える**（正常系）
   - 入力: `logs := []entity.ActionLog{{ID: 1, ShiftID: 10}, {ID: 2, ShiftID: 20}, {ID: 3, ShiftID: 30}} / shiftMap := map[int]entity.ShiftAdmin{10: {ID: 10, TimeID: 3}, 20: {ID: 20, TimeID: 1}, 30: {ID: 30, TimeID: 2}}`
   - 期待値: `返り値の ID 列が [2, 3, 1]（TimeID 1→2→3 の順）。長さ3`
   - 根拠: 関数の主目的である昇順ソートの基本動作を固定する。ソート条件を CreatedAt 等に誤って変えた場合の回帰を検出
2. **要素1個はそのまま返る**（境界値）
   - 入力: `logs := []entity.ActionLog{{ID: 1, ShiftID: 10}} / shiftMap := map[int]entity.ShiftAdmin{10: {ID: 10, TimeID: 5}}`
   - 期待値: `長さ1で、返り値[0].ID == 1（入力と同一内容のログ）`
   - 根拠: ソート対象が単一要素でも欠落・破壊が起きないことを保証する最小境界
3. **空スライスは空の非nilスライスを返す**（境界値）
   - 入力: `logs := []entity.ActionLog{} / shiftMap := map[int]entity.ShiftAdmin{10: {ID: 10, TimeID: 1}}`
   - 期待値: `len(result) == 0 かつ result != nil（make([]entity.ActionLog, 0) 由来の空スライス）`
   - 根拠: 呼び出し元 L284 が len(sortedLogs) > 0 で分岐しており、空入力で panic せず空が返ることが通知スキップの前提
4. **nilスライスでもpanicせず空を返す**（境界値）
   - 入力: `var logs []entity.ActionLog = nil / shiftMap := map[int]entity.ShiftAdmin{10: {ID: 10, TimeID: 1}}`
   - 期待値: `panic せず len(result) == 0 かつ result != nil`
   - 根拠: nil スライスへの range と make(cap=0) が安全である現状挙動を固定。将来 logs[0] 参照などの変更が入った際の回帰検出
5. **nilマップでは全ログが除外され空を返す**（境界値）
   - 入力: `logs := []entity.ActionLog{{ID: 1, ShiftID: 10}, {ID: 2, ShiftID: 20}} / var shiftMap map[int]entity.ShiftAdmin = nil`
   - 期待値: `panic せず len(result) == 0 かつ result != nil（nil マップの lookup は ok=false になり全件 continue）`
   - 根拠: loadShiftMap が空 map や nil を返すケースでの安全性を固定。nil マップ read は Go では合法だが、書き込みを伴う実装に変わると panic するため回帰検出になる
6. **shiftMapにキーが無いログは黙って除外される**（異常系）【要判断】
   - 入力: `logs := []entity.ActionLog{{ID: 1, ShiftID: 10}, {ID: 2, ShiftID: 99}, {ID: 3, ShiftID: 30}} / shiftMap := map[int]entity.ShiftAdmin{10: {ID: 10, TimeID: 2}, 30: {ID: 30, TimeID: 1}}（ShiftID 99 のエントリ無し）`
   - 期待値: `長さ2、ID 列は [3, 1]。ShiftID=99 のログ（ID: 2）は結果に含まれない。エラーも panic も無し`
   - 根拠: 関数名は sort だがフィルタを兼ねる現状の中核挙動。呼び出し元 L285 の shiftMap[sortedLogs[0].ShiftID] がヒットする前提を支えているため、この除外が消えると別の場所でゼロ値参照が起きる
7. **TimeIDゼロ値・負値は正値より前に並ぶ**（境界値）
   - 入力: `logs := []entity.ActionLog{{ID: 1, ShiftID: 10}, {ID: 2, ShiftID: 20}, {ID: 3, ShiftID: 30}} / shiftMap := map[int]entity.ShiftAdmin{10: {ID: 10, TimeID: 2}, 20: {ID: 20}, 30: {ID: 30, TimeID: -1}}（ShiftID 20 は TimeID がゼロ値 0）`
   - 期待値: `ID 列が [3, 2, 1]（TimeID -1 → 0 → 2 の順）`
   - 根拠: TimeID 未設定（ゼロ値）のシフトが先頭に来る現状挙動の固定。Scan 漏れや未設定データが混入した際に通知の並びがどうなるかを明文化する
8. **同一TimeIDのタイは全件保持される（相対順序は未保証）**（境界値）【要判断】
   - 入力: `logs := []entity.ActionLog{{ID: 1, ShiftID: 10}, {ID: 2, ShiftID: 20}, {ID: 3, ShiftID: 30}} / shiftMap := map[int]entity.ShiftAdmin{10: {ID: 10, TimeID: 1}, 20: {ID: 20, TimeID: 1}, 30: {ID: 30, TimeID: 1}}`
   - 期待値: `長さ3で、ID の集合が {1, 2, 3}（順序は不問のアサートにする。例: ID を収集してソート後 [1, 2, 3] と比較）`
   - 根拠: タイでも要素が落ちない・重複しないことを保証する。順序をアサートしないのは sort.Slice（非安定）採用のため
9. **入力スライスを破壊しない**（正常系）
   - 入力: `logs := []entity.ActionLog{{ID: 1, ShiftID: 10}, {ID: 2, ShiftID: 20}} / shiftMap := map[int]entity.ShiftAdmin{10: {ID: 10, TimeID: 2}, 20: {ID: 20, TimeID: 1}}（ソートで順序が入れ替わる入力）`
   - 期待値: `返り値の ID 列は [2, 1] だが、呼び出し後の logs の ID 列は [1, 2] のまま。かつ返り値は logs と別のスライス（&result[0] != &logs[0]）`
   - 根拠: 新規スライスを組み立てる非破壊実装の固定。呼び出し元 L277 以降で logs と sortedLogs が別物として扱われており、in-place ソートに書き換わると呼び出し元の挙動が暗黙に変わるため

実行検証: worktree の api/lib/usecase/design_verify_sortlogsbytime_test.go（package usecase、ゼロ値レシーバ &notificationUseCase{} 使用）に全9ケースをサブテストとして実装し、cd api && go test ./lib/usecase/... -run TestDesignVerifySortLogsByTime -v で実行。9ケース全て起案どおりの期待値で PASS（修正 0 件、削除 0 件、実行不能ケースなし）。補足観測: 空スライス・nil スライス・nil マップの3ケースいずれも返り値は非 nil の空スライスであることをアサートで確認。タイケースは t.Logf で観測順序 [1 2 3]（挿入順）を記録したが、sort.Slice 非安定のためアサートは集合一致のみ。needs_judgment=true の2件は実挙動と期待値の食い違いではなく設計判断の論点（サイレント除外の是非、SliceStable 化の要否）として維持。実行後テストファイルを削除し、git status --porcelain が空（clean）であることを確認。commit/push は行っていない。

## formatTimeRange

`api/lib/usecase/notification_usecase.go:359`

timeMap から開始スロット（startTimeID）の時刻と、終了スロットの「次のスロット」（endTimeID+1）の時刻を引いて "HH:MM 〜 HH:MM" 形式（半角スペース + U+301C 波ダッシュ + 半角スペース区切り）の文字列を返す。次スロットのキーが無い場合のみ終了時刻を "0:00" にフォールバックする一方、開始時刻側には存在チェックが無くゼロ値（空文字）がそのまま出る。レシーバ n のフィールドには一切依存しない純関数であり、非公開メソッドのため package usecase 内テストで &notificationUseCase{} を使って呼び出せる。全ケースを実コード実行で裏取り済み。

### 要判断の論点

- **startTimeID がマップに存在しない**: 実挙動は起案期待値どおり " 〜 10:00" と確認済み。ただし endTime 側（L363-366）は ok チェック＋"0:00" フォールバックがあるのに、startTime 側（L361）は存在チェックが無く空文字がそのまま Slack 通知文面に出る。非対称でバグの疑いあり。空文字表示を仕様として固定するか、endTime と同様のフォールバック（"0:00" や "（不明）" 等）を追加するか、人間の判断が必要。

### ケース表

1. **連続スロット範囲のフォーマット**（正常系）
   - 入力: `timeMap := map[int]entity.Time{1: {Time: "9:00"}, 2: {Time: "10:00"}, 3: {Time: "11:00"}}; startTimeID := 1; endTimeID := 2`
   - 期待値: `"9:00 〜 11:00"（終了時刻は endTimeID=2 自身の "10:00" ではなく、endTimeID+1=3 の "11:00"）【実行確認済み】`
   - 根拠: 終了時刻に endTimeID+1 のスロット開始時刻を使うという中核仕様の回帰防止。リファクタで +1 が外れると "9:00 〜 10:00" に変わり検知できる。
2. **同一スロット指定（実際の呼び出し形態）**（正常系）
   - 入力: `timeMap := map[int]entity.Time{5: {Time: "13:00"}, 6: {Time: "14:00"}}; startTimeID := 5; endTimeID := 5`
   - 期待値: `"13:00 〜 14:00"【実行確認済み】`
   - 根拠: 唯一の呼び出し元 buildChangesWithTime（L385）は start と end に同じ shift.TimeID を渡す。本番で実際に通る唯一のパスの回帰防止として最重要。
3. **最終スロットで次スロットが無い（要素1個のマップ）**（境界値）
   - 入力: `timeMap := map[int]entity.Time{10: {Time: "22:00"}}; startTimeID := 10; endTimeID := 10`
   - 期待値: `"22:00 〜 0:00"（endTimeID+1=11 が不在のため "0:00" フォールバック）【実行確認済み】`
   - 根拠: 1日の最終スロットに対する明示的フォールバック（entity.Time{Time: "0:00"}）の固定化。要素1個のマップでの動作確認も兼ねる。
4. **nil マップとゼロ値 ID**（境界値）
   - 入力: `var timeMap map[int]entity.Time = nil; startTimeID := 0; endTimeID := 0`
   - 期待値: `" 〜 0:00"（先頭は開始時刻が空文字のため半角スペースで始まる。panic しない）【実行確認済み】`
   - 根拠: nil マップの読み取りは Go 仕様でゼロ値を返すため panic しない。呼び出し元で timeMap 構築に失敗した場合の最悪ケースでもクラッシュしないことの固定化。
5. **startTimeID がマップに存在しない**（異常系）【要判断】
   - 入力: `timeMap := map[int]entity.Time{2: {Time: "10:00"}}; startTimeID := 5; endTimeID := 1（endTimeID+1=2 は存在する）`
   - 期待値: `" 〜 10:00"（開始時刻が空文字のまま出力される。現状挙動の固定化）【実行確認済み】`
   - 根拠: shift.TimeID が timeMap に無い（データ不整合）ケース。startTime 側の欠落時挙動を明示的にテストで固定し、無言の仕様変更を防ぐ。
6. **次スロットは存在するが Time フィールドがゼロ値（空文字）**（境界値）
   - 入力: `timeMap := map[int]entity.Time{1: {Time: "9:00"}, 2: {}}; startTimeID := 1; endTimeID := 1`
   - 期待値: `"9:00 〜 "（末尾は半角スペースで終わる。キーは存在するため "0:00" フォールバックは発火しない）【実行確認済み】`
   - 根拠: "0:00" フォールバックの発火条件が「値が空かどうか」ではなく「キーが存在するかどうか」であることの固定化。ゼロ値 entity.Time の扱いを明示する。
7. **負の ID（endTimeID+1 の算術で key 0 を参照）**（異常系）
   - 入力: `timeMap := map[int]entity.Time{0: {Time: "8:00"}}; startTimeID := -1; endTimeID := -1`
   - 期待値: `" 〜 8:00"（startTimeID=-1 は不在で空文字、endTimeID+1=0 は key 0 にヒットし "8:00"）【実行確認済み】`
   - 根拠: ID の妥当性検証が無いこと、および endTimeID+1 の算術が負値にもそのまま適用されることの固定化。不正入力でも panic しないことを保証する。
8. **逆転範囲（startTimeID > endTimeID）**（異常系）
   - 入力: `timeMap := map[int]entity.Time{1: {Time: "9:00"}, 2: {Time: "10:00"}, 3: {Time: "11:00"}, 4: {Time: "12:00"}}; startTimeID := 3; endTimeID := 1`
   - 期待値: `"11:00 〜 10:00"（開始 > 終了のまま検証なしで出力される）【実行確認済み】`
   - 根拠: 範囲の前後関係を検証しないことの固定化。呼び出し元は同一 ID を渡すため実運用では発生しない GIGO ケースだが、将来バリデーションを追加する際に意図的にこのケースを更新する目印になる。

実行検証: worktree の api/lib/usecase/design_verify_formattimerange_test.go に全8ケースをテーブル駆動テスト（レシーバは &notificationUseCase{} のゼロ値）で実装し、cd api && go test ./lib/usecase/... -run TestDesignVerifyFormatTimeRange -v で実行。8ケース全て PASS（1回目の実行で全一致）。起案どおり: 8件、実挙動との食い違いによる expected 修正: 0件、実行不能で削除・修正したケース: 0件。「startTimeID がマップに存在しない」の needs_judgment=true は実挙動の食い違いではなく、start/end のフォールバック非対称（L361 に ok チェック無し）を仕様とするかバグとするかの設計判断が残るため維持。検証後テストファイルを削除し、git status --porcelain が空（clean）であることを確認。commit/push は一切していない。

## buildChangesWithTime

`api/lib/usecase/notification_usecase.go:371`

アクションログの配列から Slack 通知本文を組み立てる純関数。各ログの diff_payload(JSON) をパースし、shiftMap/timeMap から「開始 〜 終了：」の時間プレフィックスを付け、ActionType（CREATE/UPDATE/DELETE/その他）ごとに整形した行を改行で結合して返す。レシーバ n は n.formatTimeRange の呼び出しにのみ使われ、formatTimeRange も含めてレシーバのフィールド（repo 群・slackService）には一切依存しないため、&notificationUseCase{} のゼロ値レシーバでモック無しにテスト可能（実行で確認済み）。

### 要判断の論点

- **timeMap に TimeID 自体が無い → 開始時刻が空文字で出力**: 実行により起案どおりの挙動（" 〜 0:00：受付 → 警備"）を確認済み。終了側欠落には "0:00" フォールバックがあるのに開始側欠落は空文字のまま Slack 通知に載る非対称。開始側にもフォールバック値を入れる／時間プレフィックス自体を省略するのが本来の意図ではないかというバグ疑いあり。現状挙動を仕様として固定してよいか要判断（修正するなら期待値が変わる）。
- **diff_payload がパース不能（不正 JSON / nil）なログは黙ってスキップ**: 実行により起案どおりの挙動を確認済み。パース不能なログが通知から黙って消える（ログ出力すら無い）ため、シフト変更の通知漏れに直結しうる。欠落を許容仕様として固定するか、警告ログ追加や「（不明）」行としての出力に改めるか要判断。テスト期待値は現状のスキップ挙動としている。

### ケース表

1. **UPDATE 正常系（全マップ完備）**（正常系）
   - 入力: `` logs := []entity.ActionLog{{ID: 1, ShiftID: 10, ActionType: "UPDATE", DiffPayload: json.RawMessage(`{"changes":[{"field":"task","old":"受付","new":"警備"}]}`)}} shiftMap := map[int]entity.ShiftAdmin{10: {ID: 10, TaskID: 100, TimeID: 5}} taskMap := map[int]entity.Task{100: {ID: 100, Task: "警備"}} timeMap := map[int]entity.Time{5: {ID: 5, Time: "10:00"}, 6: {ID: 6, Time: "11:00"}} ``
   - 期待値: `"10:00 〜 11:00：受付 → 警備"（〜は全角チルダ、：は全角コロン、→ の前後は半角スペース）【実行確認済み】`
   - 根拠: 最頻パスの回帰防止。時間プレフィックスが timeMap[TimeID] 〜 timeMap[TimeID+1] で作られること、payload の old/new が採用されることを固定する。
2. **CREATE 正常系（payload の new が DB 現在値より優先）**（正常系）
   - 入力: `` logs := []entity.ActionLog{{ID: 2, ShiftID: 10, ActionType: "CREATE", DiffPayload: json.RawMessage(`{"changes":[{"field":"task","old":"","new":"設営"}]}`)}} shiftMap := map[int]entity.ShiftAdmin{10: {TaskID: 100, TimeID: 5}} taskMap := map[int]entity.Task{100: {Task: "現DB名"}} timeMap := map[int]entity.Time{5: {Time: "10:00"}, 6: {Time: "11:00"}} ``
   - 期待値: `"10:00 〜 11:00：設営（新規）"（taskMap の "現DB名" ではなく payload の "設営" が出る）【実行確認済み】`
   - 根拠: CREATE 分岐で payload 側のタスク名が DB 現在値フォールバックより優先される仕様を固定。taskMap の値をわざとずらすことで優先順位の回帰を検出できる。
3. **複数ログの改行結合と入力順保持**（正常系）
   - 入力: `` logs := []entity.ActionLog{ {ShiftID: 10, ActionType: "UPDATE", DiffPayload: json.RawMessage(`{"changes":[{"old":"受付","new":"警備"}]}`)}, {ShiftID: 0, ActionType: "DELETE", DiffPayload: json.RawMessage(`{"deleted_task":"受付"}`)}, } shiftMap := map[int]entity.ShiftAdmin{10: {TaskID: 100, TimeID: 5}} taskMap := map[int]entity.Task{100: {Task: "警備"}} timeMap := map[int]entity.Time{5: {Time: "10:00"}, 6: {Time: "11:00"}} ``
   - 期待値: `"10:00 〜 11:00：受付 → 警備\n受付（削除）"（strings.Join による \n 結合。DELETE 行は ShiftID=0 が shiftMap に無いため時間プレフィックスなし）【実行確認済み】`
   - 根拠: 複数行が入力順のまま \n 結合されることを固定する（この関数はソートしない。並べ替えは呼び出し前の sortLogsByTime の責務、という責務分担の回帰防止）。
4. **logs が nil / 空スライス**（境界値）
   - 入力: `logs := []entity.ActionLog(nil)（[]entity.ActionLog{} でも同一挙動） shiftMap := map[int]entity.ShiftAdmin{} taskMap := map[int]entity.Task{} timeMap := map[int]entity.Time{}`
   - 期待値: `""（changes が nil のまま strings.Join(nil, "\n") == ""）【nil・空スライスの両方を別サブテストで実行し、いずれも "" を確認済み】`
   - 根拠: 空入力で空文字を返す境界を固定。呼び出し元 buildGroupedMessage は len==0 を先に弾くが、この関数単体でも安全であることを保証する。
5. **全マップ nil + DELETE ログ（shift_id NULL 相当の ShiftID=0）**（境界値）
   - 入力: `` logs := []entity.ActionLog{{ShiftID: 0, ActionType: "DELETE", DiffPayload: json.RawMessage(`{"deleted_task":"撤収作業"}`)}} shiftMap := map[int]entity.ShiftAdmin(nil) taskMap := map[int]entity.Task(nil) timeMap := map[int]entity.Time(nil) ``
   - 期待値: `"撤収作業（削除）"（nil マップ読み取りは Go では安全。時間プレフィックスなしで出力される）【実行確認済み: panic せず期待どおり】`
   - 根拠: nil マップで panic しないこと、および DELETE だけは shiftMap/taskMap のルックアップ成功を要求せず diff_payload 単独で出力される特別扱い（DB の shift_id NULL → ShiftID=0 ケース）を固定する。
6. **shiftMap / taskMap にキーが無い非 DELETE ログはスキップ**（境界値）
   - 入力: `` logs := []entity.ActionLog{ {ShiftID: 99, ActionType: "UPDATE", DiffPayload: json.RawMessage(`{"changes":[{"old":"a","new":"b"}]}`)}, {ShiftID: 10, ActionType: "UPDATE", DiffPayload: json.RawMessage(`{"changes":[{"old":"a","new":"b"}]}`)}, } shiftMap := map[int]entity.ShiftAdmin{10: {TaskID: 100, TimeID: 5}} // 99 は無い taskMap := map[int]entity.Task{} // 100 も無い timeMap := map[int]entity.Time{5: {Time: "10:00"}, 6: {Time: "11:00"}} ``
   - 期待値: `""（1件目は shiftMap 欠落、2件目は taskMap 欠落で両方 continue され、結果は空文字）【実行確認済み】`
   - 根拠: 2 段のルックアップ失敗（shiftMap 欠落・taskMap 欠落）がどちらも黙ってスキップになる現状挙動を 1 ケースで固定する。
7. **timeMap に TimeID+1 が無い（最終時間帯）→ 終了時刻 0:00 フォールバック**（境界値）
   - 入力: `` logs := []entity.ActionLog{{ShiftID: 10, ActionType: "UPDATE", DiffPayload: json.RawMessage(`{"changes":[{"old":"受付","new":"警備"}]}`)}} shiftMap := map[int]entity.ShiftAdmin{10: {TaskID: 100, TimeID: 5}} taskMap := map[int]entity.Task{100: {Task: "警備"}} timeMap := map[int]entity.Time{5: {Time: "22:00"}} // キー 6 が無い ``
   - 期待値: `"22:00 〜 0:00：受付 → 警備"【実行確認済み】`
   - 根拠: formatTimeRange の endTimeID+1 欠落時に entity.Time{Time: "0:00"} へフォールバックする明示的な実装（最終スロットは翌 0:00 終了の想定）を固定する。時間 ID が連番であることに依存した設計の目印にもなる。
8. **timeMap に TimeID 自体が無い → 開始時刻が空文字で出力**（異常系）【要判断】
   - 入力: `` logs := []entity.ActionLog{{ShiftID: 10, ActionType: "UPDATE", DiffPayload: json.RawMessage(`{"changes":[{"old":"受付","new":"警備"}]}`)}} shiftMap := map[int]entity.ShiftAdmin{10: {TaskID: 100, TimeID: 5}} taskMap := map[int]entity.Task{100: {Task: "警備"}} timeMap := map[int]entity.Time{} // 空 ``
   - 期待値: `" 〜 0:00：受付 → 警備"（先頭が半角スペース始まり。開始時刻はゼロ値 Time の空文字がそのまま入る）【実行確認済み: 先頭スペース込みで完全一致】`
   - 根拠: formatTimeRange は startTime を ok チェックなしで timeMap[startTimeID] のゼロ値のまま使うため、開始側欠落だと通知文が「 〜 0:00：」と壊れて見える。この非対称フォールバックを現状挙動としてピン留めしつつ論点化する。
9. **diff_payload がパース不能（不正 JSON / nil）なログは黙ってスキップ**（異常系）【要判断】
   - 入力: `` logs := []entity.ActionLog{ {ID: 1, ShiftID: 10, ActionType: "UPDATE", DiffPayload: nil}, {ID: 2, ShiftID: 10, ActionType: "UPDATE", DiffPayload: json.RawMessage(`{invalid`)}, {ID: 3, ShiftID: 10, ActionType: "UPDATE", DiffPayload: json.RawMessage(`{"changes":[{"old":"受付","new":"警備"}]}`)}, } shiftMap := map[int]entity.ShiftAdmin{10: {TaskID: 100, TimeID: 5}} taskMap := map[int]entity.Task{100: {Task: "警備"}} timeMap := map[int]entity.Time{5: {Time: "10:00"}, 6: {Time: "11:00"}} ``
   - 期待値: `"10:00 〜 11:00：受付 → 警備"（ID:1 は nil で unmarshal エラー、ID:2 は不正 JSON でエラー、いずれも警告なしにスキップされ ID:3 の 1 行のみ）【実行確認済み: nil DiffPayload も panic せず同経路でスキップ】`
   - 根拠: json.Unmarshal 失敗時の continue による欠落挙動を固定。DiffPayload が nil でも "unexpected end of JSON input" で同経路に落ちることを含めて確認する。
10. **UPDATE で payload に changes キーが無い → 不明→DB現在値フォールバック**（異常系）
   - 入力: `` logs := []entity.ActionLog{{ShiftID: 10, ActionType: "UPDATE", DiffPayload: json.RawMessage(`{}`)}} shiftMap := map[int]entity.ShiftAdmin{10: {TaskID: 100, TimeID: 5}} taskMap := map[int]entity.Task{100: {Task: "警備"}} timeMap := map[int]entity.Time{5: {Time: "10:00"}, 6: {Time: "11:00"}} // DiffPayload が `{"changes":[]}`（空配列）でも len>0 が偽になり同一挙動 ``
   - 期待値: `` "10:00 〜 11:00：（不明） → 警備"（old は既定の「（不明）」、new は taskMap の DB 現在値 "警備" にフォールバック）【`{}` と `{"changes":[]}` の両方を別サブテストで実行し、同一出力を確認済み】 ``
   - 根拠: 実装コメントで「フォールバック: DB現在値」と明記された意図的な既定値ロジックの回帰防止。changes キー欠落と空配列の両方が同じ分岐に落ちる境界も兼ねる。

実行検証: worktree の api/lib/usecase/design_verify_buildchangeswithtime_test.go に起案10ケースを12サブテスト（「logs nil/空スライス」と「changes キー無し/空配列」の同一挙動主張をそれぞれ2サブテストに分割）として実装し、cd api && go test ./lib/usecase/... -run TestDesignVerifyBuildChangesWithTime -v を実行。12/12 PASS で全ケースが起案期待値と完全一致（修正0件、削除0件、実行不能0件）。ゼロ値レシーバ &notificationUseCase{} でモック無しに呼べることも実証。開始時刻欠落の先頭半角スペース、nil マップ・nil DiffPayload の非 panic、全角チルダ/コロンの文字種まで %q 比較で裏取り済み。needs_judgment の2件（開始時刻フォールバック非対称、パース不能ログの黙殺）は実挙動が起案どおりであることを確認した上で設計判断の論点として維持。実行後テストファイルを削除し、git status --porcelain が空（clean）であることを確認。commit/push は一切していない。

## convertShiftCardDataToShifts

`api/lib/usecase/shift_usecase.go:433`

DB の JOIN 結果フラット構造体 []entity.ShiftCardData を、モバイル API 用のネスト構造体 []entity.Shift へ 1:1 で詰め替える純関数。YearValue のみ strconv.Atoi で文字列→int 変換し、エラーは黙殺して 0 にフォールバックする。レシーバ *shiftUseCase のフィールドには一切依存せず、(&shiftUseCase{}) のゼロ値レシーバからモック無しで呼び出せることを実行で確認済み。マップ引数は持たないため「nil マップ／キー欠落」の境界は適用外であり、代わりに nil スライス・ゼロ値要素・Atoi 失敗を境界ケースとして採用した。全 8 ケースを go test で実行し、起案期待値と実挙動の一致を裏取り済み。

### 要判断の論点

- **YearValue が数値でない文字列**: Atoi のエラーを _ で捨てて 0 にフォールバックするのは意図的か要判断。DB 由来の値なので通常は数値だが、汚損時に「0年」として API レスポンスに載り静かに壊れる。ログ出力やエラー返却に変える選択肢もあるが、シグネチャ変更（error 追加）は呼び出し側に波及するため、保守モードでは現状固定が妥当かの確認を求める（実挙動は起案どおりで食い違い無し。論点は設計意図の確認のみ）
- **未マッピングフィールドは出力に伝播しない**: TaskMobile には Remark・MaxMember・BureauID フィールドが存在するのに ShiftCardData の対応値（TaskRemark/MaxMember/TaskBureauID）を詰めていない。モバイル画面で不要だから意図的に省いたのか、マッピング漏れなのか要判断。モバイル側（Flutter）がこれらを参照していれば実バグ。現状固定でテストを書くが、期待値を「詰める」に変える判断はモバイル側の参照調査後にすべき（実挙動は起案どおりで食い違い無し。論点は設計意図の確認のみ）

### ケース表

1. **全フィールド設定済み1要素の完全マッピング**（正常系）
   - 入力: `[]entity.ShiftCardData{{ShiftID: 1, UserID: 2, TaskID: 3, YearID: 4, DateID: 5, TimeID: 6, WeatherID: 7, IsAttendance: true, TaskName: "受付", TaskColor: "#FF0000", TaskURL: "https://example.com/manual", PlaceName: "体育館", TimeValue: "9:00", UserName: "山田太郎", UserBureauID: 8, UserGradeID: 9, YearValue: "2024", DateValue: "9/13", WeatherValue: "晴れ"}}`
   - 期待値: `[]entity.Shift{{ID: 1, Task: entity.TaskMobile{ID: 3, Task: "受付", Color: "#FF0000", Place: "体育館", Url: "https://example.com/manual"}, User: entity.User{ID: 2, Name: "山田太郎", BureauID: 8, GradeID: 9}, Year: entity.Year{ID: 4, Year: 2024}, Date: entity.Date{ID: 5, Date: "9/13"}, Time: entity.Time{ID: 6, Time: "9:00"}, Weather: entity.Weather{ID: 7, Weather: "晴れ"}, IsAttendance: true}} と reflect.DeepEqual で一致（CreatedAt/UpdatedAt は両辺ゼロ値 time.Time のため一致する）【実行で確認済み】`
   - 根拠: 各フィールドの対応関係（ShiftID→ID、PlaceName→Task.Place、UserBureauID→User.BureauID 等）の取り違え回帰を防ぐ基準ケース。IsAttendance=true の引き継ぎも兼ねる
2. **複数要素で入力順が保存される**（正常系）
   - 入力: `[]entity.ShiftCardData{{ShiftID: 10, TimeID: 3, YearValue: "2024"}, {ShiftID: 20, TimeID: 1, YearValue: "2024"}, {ShiftID: 30, TimeID: 2, YearValue: "2024"}}（TimeID を昇順にしないことでソートされないことを検証）`
   - 期待値: `len == 3 かつ shifts[0].ID == 10, shifts[1].ID == 20, shifts[2].ID == 30（入力順のまま。ソートは呼び出し側 GetShiftCardsByUserAndDateAndWeather の責務であり本関数では行わない）【実行で確認済み】`
   - 根拠: 本関数に将来ソートやフィルタが混入する回帰を防ぐ。変換専任であることの仕様固定
3. **空スライス入力は nil を返す**（境界値）
   - 入力: `[]entity.ShiftCardData{}`
   - 期待値: `戻り値は nil（var shifts []entity.Shift に append されないため）。len == 0 かつ shifts == nil であり、reflect.DeepEqual(nil戻り値, []entity.Shift{}) は false になることまで実行で確認済み。assert では len == 0 に加えて shifts == nil を固定する`
   - 根拠: nil スライス返却という現状挙動の固定。呼び出し側は range しかしないため実害はないが、将来 json 化などで nil/empty の差が問題になったとき検知できる
4. **nil スライス入力は nil を返す**（境界値）
   - 入力: `var data []entity.ShiftCardData = nil を渡す`
   - 期待値: `panic せず nil を返す（range over nil slice は 0 回イテレーションで安全）【実行で確認済み】`
   - 根拠: リポジトリがエラー無しで nil を返した場合の安全性の固定。nil 入力での panic 回帰を防ぐ
5. **ゼロ値要素1個の変換**（境界値）
   - 入力: `[]entity.ShiftCardData{{}}（全フィールドゼロ値。YearValue は ""）`
   - 期待値: `len == 1 かつ shifts[0] が entity.Shift{} と reflect.DeepEqual で一致（strconv.Atoi("") はエラーだが黙殺され Year.Year == 0）【実行で確認済み】`
   - 根拠: 要素1個という最小非空入力と、全ゼロ値での panic 無しを同時に固定する
6. **YearValue が数値でない文字列**（異常系）【要判断】
   - 入力: `[]entity.ShiftCardData{{ShiftID: 1, YearValue: "abc"}, {ShiftID: 2, YearValue: " 2024"}}（先頭空白は Atoi が受理しない）の2要素`
   - 期待値: `いずれも Year.Year == 0 になり、エラーもログも発生せず正常に変換が完了する（現状挙動）【実行で確認済み。" 2024" も Atoi エラーで 0 になることを確認】`
   - 根拠: strconv.Atoi のエラー黙殺という現状挙動の固定。DB の year_value 汚損時に year=0 で静かに返るパスの検知
7. **YearValue が負数文字列**（境界値）
   - 入力: `[]entity.ShiftCardData{{YearValue: "-1"}}`
   - 期待値: `Year.Year == -1（strconv.Atoi は符号付きをそのまま受理するため負値が通過する）【実行で確認済み】`
   - 根拠: Atoi の受理範囲（符号付き整数）をそのまま通す現状挙動の固定。バリデーション追加時に検知できる
8. **未マッピングフィールドは出力に伝播しない**（異常系）【要判断】
   - 入力: `[]entity.ShiftCardData{{TaskID: 3, TaskName: "受付", TaskRemark: "注意事項あり", MaxMember: 10, TaskBureauID: 5, PlaceID: 7, YearValue: "2024"}}`
   - 期待値: `shifts[0].Task.Remark == "", shifts[0].Task.MaxMember == 0, shifts[0].Task.BureauID == 0, shifts[0].Task.YearID == 0（TaskRemark/MaxMember/TaskBureauID/PlaceID は変換で捨てられる。TaskMobile に対応フィールドがあるのに詰め替えていない）【実行で確認済み】`
   - 根拠: 入力にあって出力に無いフィールドの扱いを明文化する。将来マッピングを追加・削除したとき必ずこのテストに引っかかる

実行検証: worktree の api/lib/usecase/design_verify_convertshiftcarddatatoshifts_test.go に全 8 ケースをテスト関数 8 本として実装し、cd api && go test ./lib/usecase/ -run TestDesignVerifyConvertShiftCardDataToShifts -v で実行。8 ケース全て PASS し、起案どおりが 8 件、実挙動との食い違いによる expected 修正は 0 件、実行不能による削除・修正も 0 件。ゼロ値レシーバ (&shiftUseCase{}) からモック無しで呼び出せることも実行で裏付けた。空スライス入力ケースでは戻り値が nil であることに加え、reflect.DeepEqual(nil, []entity.Shift{}) が false になる点（起案の注記）も assert で確認した。needs_judgment=true の 2 件（Atoi エラー黙殺、TaskRemark/MaxMember/TaskBureauID の未マッピング）は実挙動と期待値の食い違いではなく、現状挙動を仕様として固定してよいかの設計意図確認として据え置き。検証後テストファイルは削除し、git status --porcelain が空（clean）であることを確認。commit/push は行っていない。

## groupContinuousShifts

`api/lib/usecase/shift_usecase.go:478`

ソート済みの同一タスクのシフト列を受け取り、隣接要素が「同じ Task.ID かつ Time.ID がちょうど +1」のとき同じグループに連結し、それ以外で新グループを開始して [][]entity.Shift を返す関数（実装 L478-506）。空入力（nil 含む）は非 nil の空スライスを返す（L479-481 の早期リターン）。レシーバ a *shiftUseCase のフィールドには一切依存しない純関数であり、同パッケージ内テストから (&shiftUseCase{}).groupContinuousShifts(...) とゼロ値レシーバで呼べることを実行で確認済み。entity は github.com/NUTFes/SeeFT/api/lib/entity（Shift.Task は entity.TaskMobile、Shift.Time は entity.Time、いずれも ID int）。全10ケースを実際に go test で実行し、起案期待値と実挙動の食い違いはゼロだった。

### 要判断の論点

- **同一 TimeID の重複は別グループに分割される**: 実行により現状挙動が起案どおり「分割」であることは確認済み。残る論点は設計意図: 同時刻のシフトを別カードに分割するのが意図か要判断。呼び出し元は userID でフィルタ済みのため通常は同一 TimeID が重複しないが、もし「同時刻・同タスクは同じグループにまとめる（curr.Time.ID == prev.Time.ID も連続扱い）」が本来の意図なら現状はバグの疑いがある。保守モード方針では現状挙動（分割）を正解としてテストを書き、変更するなら別 issue とするのが妥当と考えるが、期待値の確定は人間の判断が必要

### ケース表

1. **連続3件が1グループになる**（正常系）
   - 入力: `ヘルパー sh(task, time int) entity.Shift { return entity.Shift{Task: entity.TaskMobile{ID: task}, Time: entity.Time{ID: time}} } を定義し、[]entity.Shift{sh(1,1), sh(1,2), sh(1,3)} を渡す（以降のケースも同じ sh を使用）`
   - 期待値: `reflect.DeepEqual で [][]entity.Shift{{sh(1,1), sh(1,2), sh(1,3)}} と一致（グループ数1、要素順は入力順のまま）。実行で確認済み`
   - 根拠: 最も基本的な連結パス（curr.Time.ID == prev.Time.ID+1 の真側）の回帰防止。グループ内の要素順序が保存されることも固定する
2. **TimeID の欠番でグループが分割される**（正常系）
   - 入力: `[]entity.Shift{sh(1,1), sh(1,2), sh(1,4), sh(1,5)}`
   - 期待値: `reflect.DeepEqual で [][]entity.Shift{{sh(1,1), sh(1,2)}, {sh(1,4), sh(1,5)}} と一致（TimeID 2→4 のギャップで分割、グループ数2）。実行で確認済み`
   - 根拠: 分割パス（連続条件の偽側→currentGroup 確定→新グループ開始）と、最終グループの取りこぼし防止（ループ後の append、L501-503）を同時に検証する
3. **空スライスは非 nil の空結果**（境界値）
   - 入力: `[]entity.Shift{}`
   - 期待値: `戻り値 got について got != nil かつ len(got) == 0（reflect.DeepEqual(got, [][]entity.Shift{}) が true）。実行で確認済み`
   - 根拠: L479-481 の早期リターンが nil でなく空スライスを返す仕様の固定。呼び出し元の range が安全に空回りすること、および JSON 化時に null にならない性質の回帰防止
4. **nil スライスも空スライス扱い**（境界値）
   - 入力: `var in []entity.Shift = nil として groupContinuousShifts(in) を呼ぶ`
   - 期待値: `空スライスと同じく got != nil かつ len(got) == 0。実行で確認済み（panic せず非 nil 空スライスが返る）`
   - 根拠: Go では len(nil) == 0 なので空スライスと同経路だが、nil 入力で panic しないこと・出力が nil に化けないことを明示的に固定する
5. **要素1個は単独グループ**（境界値）
   - 入力: `[]entity.Shift{sh(1,5)}`
   - 期待値: `reflect.DeepEqual で [][]entity.Shift{{sh(1,5)}} と一致（グループ数1、要素数1）。実行で確認済み`
   - 根拠: ループ本体（i=1 以降）を一度も通らず、初期 currentGroup がそのまま最終 append される最小パスの回帰防止
6. **同一 TimeID の重複は別グループに分割される**（境界値）【要判断】
   - 入力: `[]entity.Shift{sh(1,3), sh(1,3), sh(1,4)}`
   - 期待値: `現状挙動: reflect.DeepEqual で [][]entity.Shift{{sh(1,3)}, {sh(1,3), sh(1,4)}} と一致（1件目の sh(1,3) は単独グループ、2件目の sh(1,3) から sh(1,4) が連結）。実行で確認済み`
   - 根拠: 連続条件が「+1 ちょうど」であり同値（+0）を連続とみなさないことの固定。同タスク・同時刻に複数シフトが存在するデータが来た場合の挙動を明示する
7. **TaskID が異なる隣接シフトは TimeID が連続でも分割**（異常系）
   - 入力: `[]entity.Shift{sh(1,1), sh(2,2)}（TimeID は 1→2 で +1 連続だが TaskID が 1 と 2 で異なる）`
   - 期待値: `reflect.DeepEqual で [][]entity.Shift{{sh(1,1)}, {sh(2,2)}} と一致（グループ数2）。実行で確認済み`
   - 根拠: 呼び出し元は taskID ごとに分割してから渡すため通常混在しないが、その不変条件が破れた入力でもタスク跨ぎで誤連結しないガード（prev.Task.ID == curr.Task.ID 条件、L491）の回帰防止
8. **未ソート入力（降順）はマージされない**（異常系）
   - 入力: `[]entity.Shift{sh(1,3), sh(1,2), sh(1,1)}`
   - 期待値: `reflect.DeepEqual で [][]entity.Shift{{sh(1,3)}, {sh(1,2)}, {sh(1,1)}} と一致（全要素が単独グループ、入力順保存、内部でソートされない）。実行で確認済み`
   - 根拠: 本関数は隣接比較のみで内部ソートを行わず、呼び出し元の sort.Sort(ByTime(taskShifts)) が事前条件であることをテストとして文書化する。将来内部にソートを足す変更が入った場合に検知できる
9. **ゼロ値要素2個は別グループ**（境界値）
   - 入力: `[]entity.Shift{{}, {}}（両方とも Task.ID=0, Time.ID=0 のゼロ値）`
   - 期待値: `reflect.DeepEqual で [][]entity.Shift{{{}}, {{}}} と一致（0 == 0+1 が偽のため2グループ）。実行で確認済み`
   - 根拠: ゼロ値の Shift が来ても panic せず、同一 TimeID(0) の重複として分割される挙動の固定。変換層の不具合でゼロ値が混入した場合の振る舞いを明示する
10. **ゼロ値と TimeID=1 は連続扱いで1グループ**（境界値）
   - 入力: `[]entity.Shift{{}, sh(0,1)}（1件目はゼロ値: Task.ID=0, Time.ID=0。2件目は Task.ID=0, Time.ID=1）`
   - 期待値: `reflect.DeepEqual で [][]entity.Shift{{{}, sh(0,1)}} と一致（Task.ID 0==0 かつ Time.ID 1 == 0+1 で連結、グループ数1）。実行で確認済み`
   - 根拠: TimeID=0 起点でも +1 判定が機械的に成立することの固定。ID にドメイン的な下限チェックが無い（0 や負値も通る）ことをテストとして文書化する

実行検証: worktree の api/lib/usecase/design_verify_groupcontinuousshifts_test.go（package usecase）に全10ケースを実装し、cd api && go test ./lib/usecase/... -run TestDesignVerifyGroupContinuousShifts -v で実行。10ケース全て一発 PASS（起案どおり10件、実挙動との食い違いによる修正0件、実行不能による削除0件）。unexported 関数はゼロ値レシーバ (&shiftUseCase{}) から問題なく呼べた。entity のパスは api/lib/entity（import: github.com/NUTFes/SeeFT/api/lib/entity）で、entity.Shift{Task: entity.TaskMobile{ID: ...}, Time: entity.Time{ID: ...}} のリテラル構築がそのまま通ることを確認。needs_judgment の1件（同一 TimeID 重複）は実挙動が起案の「現状挙動: 分割」と一致したため期待値は据え置き、設計意図の判断のみ人間に残す。検証後テストファイルは削除し、git status --porcelain が空（clean）であることを確認。commit/push は行っていない。

## compareTimeStrings

`api/lib/usecase/shift_usecase.go:509`

"8:00" のような時刻文字列2つを時分→総分に換算して数値比較し、-1/0/1 を返す比較関数。sort.Slice（L426）で ShiftCard.StartTime の整列に使われ、辞書順比較だと "10:00" < "8:00" になる問題を回避するのが存在意義。レシーバ a のフィールドには一切依存せず（a は本体で未使用）、entity 型も使わない純関数のため、&shiftUseCase{} のゼロ値レシーバでテスト可能。引数が string のみなのでスライス/マップ系境界は該当せず、空文字列・ゼロ値時刻・不正書式に読み替えて設計した。全10ケースを実コード実行で裏取り済み。

### 要判断の論点

- **空文字列は正当な時刻とも等価扱い**: 実行確認済み（実挙動も 0）。片方が不正な時点で、もう片方が正当な時刻でも一律 0（等価）を返す。sort.Slice は非安定ソートのため、不正な StartTime を持つ ShiftCard の並び順が実行ごとに変わり得る。これを仕様（不正入力は順序不問）とみなすか、エラーを返す設計に改めるべきかは人間の判断が必要。フェーズ1では現状の 0 を期待値とする
- **非数値の時分は暗黙に0:00扱い**: 実行確認済み（実挙動も -1）。Atoi のエラー黙殺により "aa:bb" が 0:00（深夜0時）として扱われ、あらゆる正当な時刻より前に整列される。ガード節の「不正は等価(0)」という方針とも不整合（不正の種類で挙動が変わる）で、バグの疑いがある。エラー時に 0 を返すか、パースエラーを呼び出し元へ伝播するかは人間の判断が必要。フェーズ1では現状の -1 を期待値とする

### ケース表

1. **桁数が異なる時刻の数値比較**（正常系）
   - 入力: `time1: "8:00", time2: "10:00"`
   - 期待値: `-1`
   - 根拠: 辞書順比較では "10:00" < "8:00" となり逆転する。この関数の存在理由そのものであり、文字列比較への安易な書き換えによる回帰を防ぐ最重要ケース
2. **逆順で正の値を返す**（正常系）
   - 入力: `time1: "10:00", time2: "8:00"`
   - 期待値: `1`
   - 根拠: 比較関数の対称性（引数入れ替えで符号反転）を固定。sort.Slice の比較子として順序が安定する前提を守る
3. **同一時刻は等価**（正常系）
   - 入力: `time1: "9:15", time2: "9:15"`
   - 期待値: `0`
   - 根拠: 反射律（同値入力で0）の固定。等価判定の分岐（L531）のカバレッジ
4. **同時間帯での分単位の差**（正常系）
   - 入力: `time1: "9:15", time2: "9:30"`
   - 期待値: `-1`
   - 根拠: 時が同じで分だけ異なる場合の比較。h*60+m の分換算式のうち m 項が効いていることの確認（m を無視する退行の防止）
5. **ゼロ値時刻と一日の最大時刻**（境界値）
   - 入力: `time1: "0:00", time2: "23:59"`
   - 期待値: `-1`
   - 根拠: 総分換算の下限（0分）と上限（1439分）の境界。ゼロ値時刻が最小として正しく整列されることの固定
6. **ゼロパディング表記の同値性**（境界値）
   - 入力: `time1: "08:05", time2: "8:05"`
   - 期待値: `0`
   - 根拠: strconv.Atoi が先頭ゼロを許容するため "08" と "8" は同値。DB/入力由来で表記ゆれがあっても比較結果が変わらないことの固定
7. **非正規化の分表記は換算後に等価**（境界値）
   - 入力: `time1: "1:30", time2: "0:90"`
   - 期待値: `0`
   - 根拠: 比較が文字列一致ではなく h*60+m の総分換算で行われるという内部仕様の固定。分が60以上でも換算されて等価になる
8. **空文字列は正当な時刻とも等価扱い**（異常系）【要判断】
   - 入力: `time1: "", time2: "10:00"`
   - 期待値: `0`
   - 根拠: 空文字列は strings.Split で要素数1となりガード節（L514-516）に入る。不正入力時のフォールバック挙動の固定
9. **コロンが2個以上の書式は等価扱い**（異常系）
   - 入力: `time1: "10:00:00", time2: "9:00"`
   - 期待値: `0`
   - 根拠: 秒付き書式（要素数3）もガード節で 0 になる経路のカバレッジ。要素数超過側の分岐を固定する。挙動の是非の論点は空文字列ケースと同一のためそちらに集約
10. **非数値の時分は暗黙に0:00扱い**（異常系）【要判断】
   - 入力: `time1: "aa:bb", time2: "8:00"`
   - 期待値: `-1`
   - 根拠: 書式（コロン1個）は通るが Atoi が失敗する入力の経路。L518-521 でエラーが _ で黙殺され h=0, m=0 となる現状挙動の固定

実行検証: worktree の api/lib/usecase/design_verify_comparetimestrings_test.go に全10ケースをテーブル駆動で実装し、go test ./lib/usecase/ -run TestDesignVerifyCompareTimeStrings -v を実行。10ケース全て PASS し、起案どおりが10件、実挙動との食い違いによる expected 修正は0件、実行不能で削除・修正したケースも0件。needs_judgment=true の2件（空文字列→0、"aa:bb"→-1）は期待値自体は実挙動と一致しており、設計判断の論点（不正入力の扱い）として judgment_note に「実行確認済み」を追記した。検証後にテストファイルを削除し、git status --porcelain が空（clean）であることを確認。commit/push は一切していない。

## BuildMessageBlocks

`api/lib/externals/slack/slack_service.go:88`

Slack通知用のBlock Kitメッセージ（ヘッダー、基本情報セクション、任意の変更内容セクション、区切り線）を組み立てる純粋な構築関数。レシーバ *SlackService のフィールド（client, channelID）には一切依存せず、(&SlackService{}).BuildMessageBlocks(params) で環境変数なしに呼び出せることを実行で確認済み。入力は MessageParams（string 5フィールドのみ）で、全8ケースを reflect.DeepEqual によるブロック構造全体比較で裏取りした。起案期待値と実挙動は全ケース一致。

### 要判断の論点

- **Changesが空白のみ**: 空白のみの Changes で「*変更内容*」見出しだけの実質空セクションが表示されるのは直感に反する可能性がある。判定を strings.TrimSpace(params.Changes) != "" にすべきかは、呼び出し元 buildGroupedMessage が空白のみの文字列を返しうるかに依存するため人間の判断が必要。実行の結果、現状はそのまま追加される挙動であることを確認済み
- **TitleがSlackヘッダー上限150文字超**: 上限超のブロックはSlack API送信時に invalid_blocks エラーになりうる。切り詰め責務をこの関数に持たせるか呼び出し元の入力制約とするかは設計判断。現状は唯一の呼び出し元がTitleを固定文字列「シフト変更通知」で渡すため実害はなく、実行の結果、切り詰めなしでそのまま生成する挙動であることを確認済み
- **UserNameにmrkdwn特殊文字・メンション構文**: Slackのmrkdwn仕様では & < > のエスケープが推奨され、verbatim=false のため <!channel> 等がメンションとして解釈されうる（mrkdwnインジェクション）。DMのみ運用のため影響は本人宛に限られるが、エスケープを入れるべきかは人間の判断が必要。実行の結果、現状は無加工で埋め込む挙動であることを確認済み

### ケース表

1. **全フィールド設定（変更内容あり）**（正常系）
   - 入力: `slack.MessageParams{Title: "シフト変更通知", UserName: "山田太郎", Date: "9月13日(土)", Weather: "晴れ", Changes: "・10:00-12:00 受付 → 会場整理"}`
   - 期待値: `長さ4の []slack.Block。[0] slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", "🔔 シフト変更通知", false, false))、[1] slack.NewSectionBlock(nil, []*slack.TextBlockObject{slack.NewTextBlockObject("mrkdwn", "ユーザー: 山田太郎", false, false), slack.NewTextBlockObject("mrkdwn", "日付: 9月13日(土)", false, false), slack.NewTextBlockObject("mrkdwn", "天気: 晴れ", false, false)}, nil)、[2] slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*変更内容*\n・10:00-12:00 受付 → 会場整理", false, false), nil, nil)、[3] slack.NewDividerBlock()。reflect.DeepEqual で全体比較可能（実行で一致確認済み）`
   - 根拠: 唯一の呼び出し元 notification_usecase.go L300 と同型の代表入力。4ブロック構成・各テキストの整形フォーマット（プレフィックス「🔔 」「ユーザー: 」「*変更内容*\n」等）の回帰を防ぐ本命ケース
2. **Changesが空文字列（変更内容ブロック省略）**（正常系）
   - 入力: `slack.MessageParams{Title: "シフト変更通知", UserName: "山田太郎", Date: "9月13日(土)", Weather: "曇り", Changes: ""}`
   - 期待値: `長さ3の []slack.Block。[0] HeaderBlock("🔔 シフト変更通知")、[1] 3フィールドのSectionBlock（ケース1と同形式）、[2] slack.NewDividerBlock()。「*変更内容*」を含むSectionBlockは存在しない（実行で一致確認済み）`
   - 根拠: 関数内唯一の分岐 if params.Changes != "" の偽側を固定する。変更内容なし通知でブロックが3個になる（区切り線は残る）現状仕様の回帰防止
3. **ゼロ値 MessageParams{}（全フィールド空）**（境界値）
   - 入力: `slack.MessageParams{}`
   - 期待値: `長さ3の []slack.Block。ヘッダーtextは "🔔 "（ベル絵文字＋半角スペースのみ）、フィールドは順に "ユーザー: ", "日付: ", "天気: "（いずれも末尾半角スペース付きでラベルのみ）、末尾は DividerBlock。panicしない（実行で一致確認済み）`
   - 根拠: 構造体ゼロ値でも安全にブロックが構築できることの検証。ラベル文字列とプレースホルダの連結仕様（trimしない・省略しない）を固定する
4. **Changesが空白のみ**（境界値）【要判断】
   - 入力: `slack.MessageParams{Title: "通知", UserName: "u", Date: "d", Weather: "w", Changes: " "}`
   - 期待値: `長さ4の []slack.Block。[2] は slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*変更内容*\n ", false, false), nil, nil)。空白のみでも変更内容ブロックが追加される（実行で一致確認済み）`
   - 根拠: 分岐条件が params.Changes != "" の完全一致判定であることを固定する。TrimSpace化などの変更が入った際に検知できる
5. **Changesが複数行（実運用のグループ化メッセージ形式）**（境界値）
   - 入力: `slack.MessageParams{Title: "シフト変更通知", UserName: "山田太郎", Date: "9月13日(土)", Weather: "晴れ", Changes: "【追加】\n・9:00-10:00 設営\n【削除】\n・13:00-14:00 受付"}`
   - 期待値: `長さ4の []slack.Block。[2] のtextは "*変更内容*\n【追加】\n・9:00-10:00 設営\n【削除】\n・13:00-14:00 受付"。改行はエスケープ・加工されずそのまま保持される（実行で一致確認済み）`
   - 根拠: 実運用では buildGroupedMessage が複数行文字列を渡す。見出し "*変更内容*" と本文が \n 1個で連結され、本文の改行が透過されることの回帰防止
6. **TitleがSlackヘッダー上限150文字超**（境界値）【要判断】
   - 入力: `slack.MessageParams{Title: strings.Repeat("あ", 151), UserName: "u", Date: "d", Weather: "w", Changes: ""}`
   - 期待値: `長さ3の []slack.Block。ヘッダーtextは "🔔 " + strings.Repeat("あ", 151)（切り詰め・バリデーションなしでそのまま格納される）（実行で一致確認済み）`
   - 根拠: Slack APIのheader block plain_textは150文字上限。本関数は長さ検証も切り詰めもしないという現状挙動を明示的に固定し、暗黙のtruncate追加を検知する
7. **UserNameにmrkdwn特殊文字・メンション構文**（異常系）【要判断】
   - 入力: `slack.MessageParams{Title: "通知", UserName: "<@U12345> & *bold* <!channel>", Date: "d", Weather: "w", Changes: ""}`
   - 期待値: `長さ3の []slack.Block。[1] の第1フィールドtextは "ユーザー: <@U12345> & *bold* <!channel>" で、&, <, > はエスケープされずそのまま格納される（verbatim=false）（実行で一致確認済み）`
   - 根拠: UserNameはDB由来のユーザー入力（user.Name）。エスケープ処理が存在しないという現状挙動を固定し、将来エスケープを追加した際にテストで意図的に更新させる
8. **制御文字・NULバイトを含む入力**（異常系）
   - 入力: `slack.MessageParams{Title: "a\x00b", UserName: "tab\tsep", Date: "line1\nline2", Weather: "\r", Changes: "end\x1b[0m"}`
   - 期待値: `panicせず長さ4の []slack.Block を返す。ヘッダーtextは "🔔 a\x00b"、フィールドは "ユーザー: tab\tsep", "日付: line1\nline2", "天気: \r"、[2] のtextは "*変更内容*\nend\x1b[0m"。全て入力文字列がそのまま連結される（実行で一致確認済み）`
   - 根拠: エラー戻り値のない関数の異常系として、どんなバイト列でもpanicせずブロックを返すこと（サニタイズなしの透過）を保証する。fmt.Sprintf連結のみで例外経路がないことの確認

実行検証: worktree の api/lib/externals/slack/design_verify_buildmessageblocks_test.go に全8ケースを実装し、cd api && go test ./lib/externals/slack/ -run TestDesignVerifyBuildMessageBlocks -v で実行。8ケース全PASS（起案どおり8件、実挙動との食い違いによる修正0件、実行不能による削除0件）。比較は各ケースの期待ブロック列を slack-go のコンストラクタ（NewHeaderBlock/NewSectionBlock/NewDividerBlock/NewTextBlockObject）で構築し reflect.DeepEqual による構造全体一致で検証したため、ブロック数だけでなく type/text/verbatim フラグ含め完全一致を確認。ゼロ値レシーバ (&SlackService{}) からの呼び出しも環境変数なしで問題なく動作。needs_judgment=true の3件（空白のみChanges・150文字超Title・mrkdwn特殊文字）は起案どおりの実挙動を確認したうえで設計判断の論点として維持。検証後テストファイルを削除し、git status --porcelain が空（clean）であることを確認。commit/push は行っていない。

