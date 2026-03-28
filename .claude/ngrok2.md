Claude Code にやらせるプロンプト（Flutter Web + トグル内PDFプレビュー + ngrok公開）
このまま貼ってください。

あなたは /Users/eisaki/workspace/SeeFT のリポジトリで作業する。
目的は Flutter Web の1ページを ngrok で一時公開し、そのページで ShiftCard のトグルを開くと docs/manual/01_44th_配線マニュアル.pdf がトグル内でプレビュー表示されること。
すでに mobile/lib/widgets/manual_viewer.dart と mobile/lib/widgets/shift_card.dart があり、ManualViewer は HtmlElementView + iframe で埋め込み表示できる前提。

要件（必須）
HTMLデモは作らない（Flutter Webでやる）
API/DBには依存しない（固定データでOK）
ShiftCard は既存の mobile/lib/widgets/shift_card.dart を必ず使う
トグル内のマニュアルプレビューは既存の ManualViewer を使う
PDFは docs/manual/01_44th_配線マニュアル.pdf を ngrok で公開した PDF直リンクを使う（例：https://<ngrok-host>/01_44th_配線マニュアル.pdf）
実装タスク
mobile/lib/pages/shift_card_manual_demo_page.dart を新規作成し、ShiftCard を1枚表示する
タスク名/時間/集合場所/担当者リストなどは適当な固定値でOK
ShiftCardData.url には String.fromEnvironment('MANUAL_PDF_URL') を入れ、未指定時はプレースホルダでもOK
Webデモ用のエントリポイントを追加する（既存アプリを壊さない）
例：mobile/lib/main_demo_web.dart を追加して、runApp(MaterialApp(home: ShiftCardManualDemoPage()))
Flutter Web をビルドして配信できるように、実行手順（コマンド）を最後に提示する
flutter build web --target lib/main_demo_web.dart --dart-define MANUAL_PDF_URL=<PDF直リンク>
python3 -m http.server で build/web を配信
ngrok http <port> で公開
共有すべきURLを出す
禁止
admin側のプロトタイプページはいじらない
docs配下にHTMLデモを作らない
API/DBを起動してデータを取る方向にしない
完了条件
公開したURLを別端末で開くと、ShiftCardが表示される
トグルを開くと、同一ページ内にPDFが埋め込み表示される
先に1点だけ
PDF直リンク（ngrokのURL）は 毎回変わるので、上のプロンプトでは MANUAL_PDF_URL を --dart-define で渡す形にしています。これで毎回コードを書き換えずに済みます。

必要なら、あなたの手元の今のPDF直リンク（https://fc01-...ngrok-free.app/01_44th_配線マニュアル.pdf）をそのまま MANUAL_PDF_URL に入れて使ってください。