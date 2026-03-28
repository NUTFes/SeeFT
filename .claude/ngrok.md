あなたは SeeFT リポジトリの作業をするコーディングエージェントです。目的は「リモート側で、実アプリ（Flutter）のシフトカードのトグル【マニュアル】から docs/manual/01_44th_配線マニュアル.pdf を開ける」ことです。

絶対にやってはいけないこと
シフトカード風のHTMLデモページを新規作成しない（不要）
PDFを埋め込んだWebページを作らない（不要）
目的を「ブラウザでデモを見る」にすり替えない
何を実装/作業するか（必要十分）
docs/manual/01_44th_配線マニュアル.pdf を一時公開できるようにする（ngrok）

ローカルで静的配信（例: python3 -m http.server）して ngrok で外部公開
目的は「PDF直リンク（https〜/01_44th_配線マニュアル.pdf）を得る」こと。デモページは不要
実アプリのトグル【マニュアル】が開くURLを、そのPDF直リンクに切り替える

現状、モバイルは ShiftCard.data.url を url_launcher で開く
バックエンドは tasks.url を shift_card_repository.go で task_url として返し、最終的にモバイルの data.url になる
したがって「DBの tasks.url を ngrok のPDF直リンクに更新する」ことがゴール
変更方法は “コード変更” ではなく “DB更新” を第一候補にする

tasks.url を更新するSQLを提示し、実行方法（psql or docker exec）も書く
もし「どのtaskを更新するか不明」なら、user_id/date_id/weather_id を指定して該当タスクを特定するSQLも提示する
成果物
リモートに共有すべきURLは「PDF直リンク」だけ（ngrokのURL）
実行手順（コマンド列）:
ローカルHTTPサーバ起動
ngrok起動
tasks.url 更新SQL（特定→更新まで）
モバイルでの確認手順（シフトカード→トグル→【マニュアル】タップで開く）

「ngrokのURLを得たら、tasks.url に入れるべきURLは https://<ngrok-host>/01_44th_配線マニュアル.pdf の“PDF直リンク”であること」
「docs/manual/demo/ など既に作ってしまった不要ファイルは削除してよい（ただし今回の最優先は要件達成）」