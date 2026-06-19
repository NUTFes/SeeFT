// 質問のレスキューを送信するページ
import 'package:seeft_mobile/configs/importer.dart';
import 'package:seeft_mobile/widgets/custom_text_field.dart';
import 'package:seeft_mobile/widgets/custom_elevated_button.dart';
import 'package:seeft_mobile/widgets/custom_elevated_button_outlined.dart';
import 'package:seeft_mobile/widgets/custom_snack_bar.dart';
import 'package:seeft_mobile/widgets/custom_error_snack_bar.dart';

class RescueRequestTabQuestionPage extends StatelessWidget {
  RescueRequestTabQuestionPage({
    super.key,
  });
  
  // レスキューを送信する関数
  Future<bool> _sendRescueRequest(
    BuildContext context,
    String question,
  ) async {
    if(
      question == ""
    ) {
      logger.e('Invalid input data');
      showCustomErrorSnackBar(context, "質問内容が入力されていません");
      return false;
    }
    try {
      final userId = await store.getUserID();
      logger.i('Sending question report:'
        '\nPlace: $question');

      // API呼び出し
      await api.postQuestionRescue(userId, question);

      logger.i('Question report sent successfully.');
      showCustomSnackBar(context, "レスキューを送信しました");
      return true;
    } catch (e) {
      logger.e('Failed to send question report: $e');
      showCustomErrorSnackBar(context, "レスキューの送信に失敗しました");
      return false;
    }
  }

  @override
  Widget build(BuildContext context) {
    // 質問内容を格納する変数
    String question = '';
    // 送信中かどうかを示す変数
    bool isSending = false;

    return Column(
      spacing: 8.0,  // 子要素間のスペース
      children: [
        Expanded(
          child: SingleChildScrollView(
            child: Container(
              padding: const EdgeInsets.all(32.0),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                crossAxisAlignment: CrossAxisAlignment.start,
                spacing: 8.0,  // 子要素間のスペース
                children: [
                  Padding(
                    padding: const EdgeInsets.only(bottom: 16.0),
                    child: Text(
                      "質問",
                      style: TextStyle(
                        fontSize: AppFontSizes.lg,
                        fontWeight: FontWeight.bold,
                        color: AppColors.textBlack,
                      ),
                      textAlign: TextAlign.center,
                    ),
                  ),
                  // 質問
                  Text(
                    "以下に質問内容を記入してください。\n「送信」ボタンを押すと本部に送信されます。",
                    style: TextStyle(
                      fontSize: AppFontSizes.md,
                      color: AppColors.textBlack,
                    ),
                    textAlign: TextAlign.center,
                  ),
                  StatefulBuilder(
                    builder: (context, setState) {
                      return CustomTextField(
                        labelText: '質問',
                        hintText: '例：〇〇のタスクの〇〇が分からないです',
                        onChanged: (String value) {
                          // 入力されたときの処理
                          setState(() {
                            question = value;
                          });
                          logger.i("入力された場所: $value");
                        },
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
        // 送信ボタン
        StatefulBuilder(
          builder: (context, setState) {
            return CustomElevatedButton(
              isDisabled: isSending,
              onPressed: () async {
                setState(() {
                  isSending = true;
                });
                // レスキューを送信する
                final isSuccess = await _sendRescueRequest(
                  context,
                  question,
                );
                logger.i("レスキューを送信しました");
                if(isSuccess){
                  // 最初の画面まで戻る
                  Navigator.of(context).popUntil((route) => route.isFirst);
                }else{
                  setState(() {
                    isSending = false;
                  });
                }
              },
              label: isSending ? "送信中..." : "送信",
            );
          }
        ),
        // 戻るボタン
        CustomElevatedButtonOutlined(
          onPressed: () {
            Navigator.of(context).pop(); // 戻るボタンが押されたときの処理
            logger.i("戻りました");
          },
          label: "戻る",
        ),
        SizedBox(height: 8.0), // スペースを追加
      ],
    );
  }
}