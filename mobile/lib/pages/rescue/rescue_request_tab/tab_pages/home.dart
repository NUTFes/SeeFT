import 'package:seeft_mobile/configs/importer.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:seeft_mobile/widgets/custom_elevated_button.dart';
import 'package:seeft_mobile/pages/rescue/rescue_request_tab/tab_pages/select_type.dart';

// レスキューを送信するタブのホーム画面
class RescueRequestTabHome extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Center(
      child: SingleChildScrollView(
        child: Container(
          padding: const EdgeInsets.all(32.0),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.center,
            spacing: 24.0,  // 子要素間のスペース
            children: [
              Text(
                "送信したレスキューの対応状況は\n「本部からの返答」タブで確認できます",
                style: TextStyle(
                  fontSize: AppFontSizes.sm,
                  color: AppColors.grayDark,
                ),
                textAlign: TextAlign.center,
              ),
              // 困った時の説明
              Container(
                padding: const EdgeInsets.all(8.0),
                decoration: BoxDecoration(
                  color: Color(0xFFD4EAE8),
                  borderRadius: BorderRadius.circular(8.0),
                ),
                width: double.infinity,
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  crossAxisAlignment: CrossAxisAlignment.center,
                  spacing: 16.0,  // 子要素間のスペース
                  children: [
                    Text(
                      "【困った時】",
                      style: TextStyle(
                        fontSize: AppFontSizes.lg,
                        fontWeight: FontWeight.bold,
                        color: AppColors.main,
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.only(left: 8.0, right: 8.0),
                      child: Text.rich(
                        TextSpan(
                          text: "現場やマニュアルでは対応できないこと",
                          style: TextStyle(
                            fontSize: AppFontSizes.md,
                            fontWeight: FontWeight.bold,
                            decoration: TextDecoration.underline,
                            color: AppColors.textBlack,
                          ),
                          children: [
                            TextSpan(
                              text: "があった場合は、以下の「レスキューを送信する」ボタンから本部へ連絡して指示を仰いでください。",
                              style: TextStyle(
                                fontSize: AppFontSizes.md,
                                fontWeight: FontWeight.normal,
                                decoration: TextDecoration.none,
                                color: AppColors.textBlack,
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                    CustomElevatedButton(
                      onPressed: () {
                        // 問題の種類の選択画面に遷移
                        Navigator.of(context).push(
                          MaterialPageRoute(builder: (_) => RescueRequestTabSelectTypePage()),
                        );
                        logger.i("レスキュー送信ボタンが押されました");
                      },
                      label: "レスキューを送信する",
                    ),
                  ],
                ),
              ),
              // 緊急時の説明
              Container(
                padding: const EdgeInsets.all(8.0),
                decoration: BoxDecoration(
                  color: Color(0xFFFFDADA),
                  borderRadius: BorderRadius.circular(8.0),
                ),
                width: double.infinity,
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  crossAxisAlignment: CrossAxisAlignment.center,
                  mainAxisSize: MainAxisSize.min,
                  spacing: 16.0,  // 子要素間のスペース
                  children: [
                    Text(
                      "【緊急時】",
                      style: TextStyle(
                        fontSize: AppFontSizes.lg,
                        fontWeight: FontWeight.bold,
                        color: AppColors.error,
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.only(left: 8.0, right: 8.0),
                      child: Row(
                        // mainAxisSize: MainAxisSize.min,
                        // mainAxisAlignment: MainAxisAlignment.center,
                        crossAxisAlignment: CrossAxisAlignment.start,
                        spacing: 8.0,  // 子要素間のスペース
                        children: [
                          Icon(
                            Icons.error,
                            color: AppColors.error,
                            size: 24.0,
                          ),
                          Expanded(
                            child: Text.rich(
                              TextSpan(
                                text: "事件や事故（人が倒れた等）",
                                style: TextStyle(
                                  fontSize: AppFontSizes.md,
                                  fontWeight: FontWeight.bold,
                                  decoration: TextDecoration.underline, 
                                  decorationColor: AppColors.error,
                                  color: AppColors.error,
                                ),
                                children: [
                                  TextSpan(
                                    text: "の場合は、「レスキューを送信する」ボタンではなく、電話またはトランシーバーで本部に直接連絡してください。\n本部の電話番号は以下の通りです。",
                                    style: TextStyle(
                                      fontSize: AppFontSizes.md,
                                      fontWeight: FontWeight.normal,
                                      decoration: TextDecoration.none,
                                      color: AppColors.error,
                                    ),
                                  ),
                                ],
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                    Container(
                      padding: const EdgeInsets.all(8.0),
                      decoration: BoxDecoration(
                        color: Color(0xFFFFDADA),
                        borderRadius: BorderRadius.circular(8.0),
                        border: Border.all(
                          color: AppColors.error,
                          width: 1.0,
                        ),
                      ),
                      width: double.infinity,
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        crossAxisAlignment: CrossAxisAlignment.center,
                        spacing: 8.0,  // 子要素間のスペース
                        children: [
                          Text(
                            //todo: envから読み込むようにする
                            "委員長 太閤良樹\n000-0000-0000",
                            style: TextStyle(
                              fontSize: AppFontSizes.md,
                              color: AppColors.error,
                            ),
                          ),
                          ElevatedButton(
                            onPressed: () {
                              //todo: 電話アプリを開く
                              _openPhoneApp("00000000000");
                              logger.i("電話が掛けられました");
                            },
                            style: ElevatedButton.styleFrom(
                              backgroundColor: AppColors.error,
                              padding: const EdgeInsets.symmetric(
                                vertical: 10.0,
                                horizontal: 24.0
                              ),
                              elevation: 0,
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(100.0),
                              ),
                              // 横幅を広げるための設定
                              // minimumSize: Size(double.infinity, 40),
                            ),
                            child: Text(
                              "電話を掛ける",
                              style: TextStyle(
                                fontSize: AppFontSizes.sm,
                                color: AppColors.textWhite
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}


// 電話アプリを開く関数
void _openPhoneApp(String phoneNumber) async {
  _launchURL(
    'tel:' + phoneNumber,
  );
}

// URLを開く関数
Future<void> _launchURL(String url) async {
  if (await canLaunch(url)) {
    await launch(url);
  } else {
    final Error error = ArgumentError('Could not launch $url');
    throw error;
  }
}
