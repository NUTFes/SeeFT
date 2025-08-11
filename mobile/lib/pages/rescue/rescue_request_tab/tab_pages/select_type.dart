// 問題の種類を選択する画面
import 'package:seeft_mobile/configs/importer.dart';
import 'package:seeft_mobile/widgets/custom_elevated_button_outlined.dart';
import 'package:seeft_mobile/pages/rescue/rescue_request_tab/tab_pages/trouble.dart';
import 'package:seeft_mobile/pages/rescue/rescue_request_tab/tab_pages/shorthanded.dart';
import 'package:seeft_mobile/pages/rescue/rescue_request_tab/tab_pages/question.dart';

class RescueRequestTabSelectTypePage extends StatelessWidget {

  RescueRequestTabSelectTypePage({
    Key? key,
  }) : super(key: key);

  // 問題の種類のリスト
  final List<Map<String, dynamic>> _listItems = [
    {
      'text': 'トラブル', 
      'explanation': '例：物品がない等',
      'icon': Icons.cancel_outlined,
      'page': RescueRequestTabTroublePage(),
    },
    {
      'text': '質問', 
      'explanation': '例：マニュアルで不明な箇所がある等',
      'icon': Icons.help_outline,
      'page': RescueRequestTabQuestionPage(),
    },
    {
      'text': '人が来ない', 
      'explanation': '例：〇〇のタスクで人が来ない',
      'icon': Icons.no_accounts,
      'page': RescueRequestTabShorthandedPage(),
    },
  ];

  @override
  Widget build(BuildContext context) {
    return Column(
      // spacing: 8.0,
      children: [
        Expanded(
          child: SingleChildScrollView(
            child: Container(
              padding: const EdgeInsets.all(32.0),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                crossAxisAlignment: CrossAxisAlignment.center,
                spacing: 24.0,  // 子要素間のスペース
                children: [
                  Text(
                    "発生した問題の種類を選択してください",
                    style: TextStyle(
                      fontSize: AppFontSizes.md,
                      fontWeight: FontWeight.bold,
                      color: AppColors.textBlack,
                    ),
                    // textAlign: TextAlign.center,
                  ),
                  // 問題の種類を選択するリスト
                  ListView.separated(
                    itemCount: _listItems.length,
                    shrinkWrap: true, // リストの高さを自動調整
                    physics: NeverScrollableScrollPhysics(), // スクロールを無効化
                    itemBuilder: (context, index) {
                      return ListTile(
                        leading: Icon(
                          _listItems[index]['icon'],
                          color: AppColors.main,
                        ),
                        title: Text(
                          _listItems[index]['text'],
                          style: TextStyle(
                            fontSize: AppFontSizes.md,
                            fontWeight: FontWeight.bold,
                            color: AppColors.textBlack,
                          ),
                        ),
                        subtitle: Text(
                          _listItems[index]['explanation'],
                          style: TextStyle(
                            fontSize: AppFontSizes.sm,
                            color: AppColors.grayDark,
                          ),
                        ),
                        trailing: Icon(
                          Icons.arrow_forward_ios,
                          color: AppColors.grayDark,
                          size: 16.0,
                        ),
                        onTap: () {
                          // 問題の種類が選択されたときの処理
                          Navigator.of(context).push(
                            // アニメーションなしで画面遷移
                            PageRouteBuilder(
                              pageBuilder: (context, animation, secondaryAnimation) => _listItems[index]['page'],
                              transitionDuration: Duration.zero,        // 遷移アニメーション時間 0
                              reverseTransitionDuration: Duration.zero, // 戻るときもアニメーション 0
                            ),
                          );
                          logger.i("問題の種類が選択されました");
                        },
                      );
                    },
                    separatorBuilder: (context, index) {
                      return Divider(
                        color: AppColors.grayLight,
                        thickness: 1.0,
                      );
                    },
                  ),
                ],
              ),
            ),
          ),
        ),
        Divider(
          color: AppColors.grayLight,
          thickness: 1.0,
        ),
        SizedBox(height: 8.0), // スペースを追加
        // 戻るボタン
        CustomElevatedButtonOutlined(
          onPressed: () {
            // 戻るボタンが押されたときの処理
            Navigator.of(context).pop();
            logger.i("戻るボタンが押されました");
          },
          label: "戻る",
        ),
        SizedBox(height: 16.0), // スペースを追加
      ],
    );
  }
}
