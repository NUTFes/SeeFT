N+1解消の効果測定として、あなたの方針（Goベンチマーク＝主効果測定、JOIN書き換え＝JSON安全網）に沿って、いまのコードベース上の「当たり筋」を前提にテスト設計を書きます。

まず対象（どこがN+1になっている前提か）
シフトカード（メイン推奨）
api/lib/usecase/shift_usecase.go の

createShiftCardFromGroup → getShiftMembersForTime
getBeforeMembers
getAfterMembers の中で、ループごとに userRep.Find / gradeRep.Find / bureauRep.Find が走るため、まさにN+1が支配的になります。
レスキュー統合（次点）
api/lib/usecase/rescue_unified_usecase.go の

getQuestionRescues
getShorthandedRescues
getTroubleRescues で rows.Next() の各行ごとに getUserName() / getTaskName() が走ります。
ベンチマーク（大量レコード投入→前後で実行時間比較）
1) 呼ぶべき関数（usecase直叩き推奨）
HTTPを挟まず、usecaseの返却生成までで評価するのがブレが少ないです。

シフト：ShiftUseCase.GetShiftCardsByUserAndDateAndWeather(userID, dateID, weatherID)
レスキュー：RescueUnifiedUseCase.GetRescuesByUserID(userID)（必要なら GetAllRescues も）
この戻り値の生成はN+1箇所に直結しています。

2) データ投入設計（Nを増やす場所を間違えない）
エンドポイント引数に「対象ユーザ1人」がいる一方で、N+1は「その時刻枠に紐づく全ユーザ（メンバー）」を取りに行く部分で爆発するので、投入の増やし方はこうするのが効きます。

target user は1人固定（ベンチ対象の userID）
ただし shifts は 同じ task/year/date/time/weather に属するユーザを大量（ユーザ数 N を増やす）
time_id は連番（getBeforeMembers/getAfterMembers と連続判定に効く）
結果として、

getUsersByShift(...) が N 件返す
その N 件それぞれで Find が複数回走る ので、JOIN化できたときの差が時間に出やすいです。
3) ベンチ実装の要点（testing.B）
ベンチでは必ずこの形にしてください。

DBへの重い投入は b.ResetTimer() の前に1回だけ（sync.Once か TestMain 相当）
反復中（b.N のループ）で json.Marshal までやると比較が歪むので、基本は「usecaseが返す構造体生成まで」で評価（ただし整合性テストは別でやる）
b.ReportAllocs() を有効にしてGC/alloc差も見られるようにする
4) 重要な注意（現状のSQLログ出力）
現状リポジトリ層で fmt.Printf によりSQLが大量出力される実装があり、ベンチ計測を壊す可能性が高いです（N+1で呼ぶたびに出る）。
ベンチを「N+1によるクエリ数差」ではなく「標準出力の遅さ」で測ってしまうリスクがあります。

対策としては（どれか1つ）：

SQLログを DEBUG_SQL 等の環境変数で無効化できるようにする（推奨）
もしくはベンチ実行を標準出力潰しで行う（応急処置。根本対策より弱い）
データ整合性テスト（JOIN書き換えでJSONが同一）
「速くなったのに返すデータが壊れてたら意味がない」ので、ベンチとは別に安全網を用意します。

1) 比較の粒度
“完全一致”を求めるなら、JSONの順序（配列順）も含めて一致させる必要があります
ただし現状 shiftRepository.Users が ORDER BY 無しで SELECT user_id ... しているので、メンバー配列順が揺れる可能性があります
実務的にはおすすめは：

「意味的に同一」を担保するため、テスト側で ShiftMembers.Members（名前/grade/bureau等のキー）をソートしてから canonical JSON化→比較
2) 方式（2つ）
ゴールデン方式：修正前のコミットで一度だけ canonical JSON（またはハッシュ）を保存し、修正後はそれと一致チェック
同一実行方式：コード内で旧実装/新実装を切り替え可能にして、同じDB状態・同じ入力で結果を突き合わせ
どちらが楽かは「旧実装が今同時にビルドできるか」で決まります（現状このリポジトリにはテスト基盤が無く、旧コードが同時に存在しないことが多いので、通常はゴールデン方式か git worktree 2回実行になります）。

次に確認したいこと（ここが決まると実装まで一気に進められます）
「修正前後」の旧実装は、同じリポジトリ内で切り替え可能ですか？それとも git の別コミット/別ブランチで2回ベンチしますか？
JSONの一致は「完全一致（順序も含む）」が必須ですか？それともテスト側で配列を正規化して“意味的に同一”を狙う形でOKですか？
ベンチで増やしたい規模感はどれくらいですか？（例：users N=500/2000、time T=8/16 など）
この3点を教えてくれれば、api/lib/usecase 配下に

ベンチ（大量投入＋前後比較）
整合性テスト（canonical JSON or ハッシュ） を具体的な形（ファイル構成/SQL投入方法/実行コマンド）に落として提案します。