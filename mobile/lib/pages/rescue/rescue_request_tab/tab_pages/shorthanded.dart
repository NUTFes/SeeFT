// 人が来ないのレスキューを送信するページ
import 'package:seeft_mobile/configs/importer.dart';
import 'package:seeft_mobile/models/rescue.dart';
import 'package:seeft_mobile/widgets/custom_dropdown_button.dart';
import 'package:seeft_mobile/widgets/custom_text_field.dart';
import 'package:seeft_mobile/widgets/custom_elevated_button.dart';
import 'package:seeft_mobile/widgets/custom_elevated_button_outlined.dart';
import 'package:seeft_mobile/widgets/custom_snack_bar.dart';
import 'package:seeft_mobile/widgets/custom_error_snack_bar.dart';

class RescueRequestTabShorthandedPage extends StatelessWidget {
  RescueRequestTabShorthandedPage({
    super.key,
  });
  
  // ユーザIDに紐づくタスクを取得する関数
  Future<List<RescueTaskDropdownMenuItem>?> fetchData() async {
    try {
      final userId = await store.getUserID();
      logger.i('Fetching tasks for user ID: $userId');
      
      final res = await api.getTasksByUserID(userId.toString());
      
      logger.i('=== API Response Received ===');
      logger.i('Response: $res');
      
      if (res is List) {
        logger.i('Response is List with ${res.length} items');
        if (res.isNotEmpty) {
          logger.i('First item: $res');
        }
      } else {
        logger.w('Response is not a List, it is: ${res.runtimeType}');
      }
      
      final resList = res as List<dynamic>;
      // レスポンスをRescueTaskDropdownMenuItemのリストに変換
      final List<RescueTaskDropdownMenuItem> taskList = resList.map((item) => RescueTaskDropdownMenuItem(
        id: item['id'],
        taskName: item['task'],
      )).toList();
      
      // デフォルトの選択肢として追加
      taskList.insert(0, RescueTaskDropdownMenuItem(
        id: 0,
        taskName: 'タスクを選択してください',
      ));
      
      // idが1と2のタスクを除外
      taskList.removeWhere((element) => element.id == 1 || element.id == 2);

      return taskList;
    } catch (err) {
      logger.e('=== API Error ==='
        '\n'
        'Failed to fetch tasks'
        '\n'
        'Error: $err');
      
      return null; // エラー時はnullを返す
    }
  }
  
  // 人が来ないのレスキューを送信する関数
  Future<bool> _sendRescueRequest(
    BuildContext context,
    RescueTaskDropdownMenuItem selectedTask,
    int missingNumber,
    String place,
  ) async {
    if(
      missingNumber == 0 || 
      place == ""
    ) {
      logger.e('Invalid input data');
      showCustomErrorSnackBar(context, "データが入力されていません");
      return false;
    }
    if(
      selectedTask.id == 0
    ) {
      logger.e('taskID is Invalid');
      showCustomErrorSnackBar(context, "タスクを選択してください");
      return false;
    }
    try {
      final userId = await store.getUserID();
      logger.i('Sending shorthanded report:'
        '\nMissing Number: $missingNumber'
        '\nPlace: $place'
        '\nSelected Task: ${selectedTask.taskName}');

      // API呼び出し
      await api.postShorthandedRescue(userId, selectedTask.id, missingNumber, place);

      logger.i('Shorthanded report sent successfully.');
      showCustomSnackBar(context, "レスキューを送信しました");
      return true;
    } catch (e) {
      logger.e('Failed to send shorthanded report: $e');
      showCustomErrorSnackBar(context, "レスキューの送信に失敗しました");
      return false;
    }
  }

  @override
  Widget build(BuildContext context) {
    // 選択されたタスクを格納する変数
    RescueTaskDropdownMenuItem selectedTask = RescueTaskDropdownMenuItem(
      id: 0,
      taskName: 'タスクを選択してください'
    );
    // 不足人数を格納する変数
    int missingNumber = 0;
    // 発生場所を格納する変数
    String place = '';
    // 送信中かどうかを示す変数
    bool isSubmitting = false;

    return FutureBuilder<List<dynamic>?>(
      future: fetchData(),
      builder: (context, snapshot) {
        if (snapshot.connectionState == ConnectionState.waiting) {
          // 読み込み中の画面
          return Center(
            child: Column(
              children: [
                Expanded(
                  child: Container(
                    padding: const EdgeInsets.all(32.0),
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      crossAxisAlignment: CrossAxisAlignment.center,
                      spacing: 24,
                      children: [
                        CircularProgressIndicator(
                          color: AppColors.main,
                        ),
                        Text(
                          "読み込み中です。",
                          style: TextStyle(
                            color: AppColors.grayDark,
                            fontSize: AppFontSizes.md,
                          ),
                          textAlign: TextAlign.start,
                        ),
                      ],
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
                    Navigator.of(context).pop(); // 戻るボタンが押されたときの処理
                    logger.i("戻りました");
                  },
                  label: "戻る",
                ),
                SizedBox(height: 16.0)
              ],
            ),
          );
        } else if (
          snapshot.hasError ||
          !snapshot.hasData ||
          snapshot.data!.isEmpty
        ) {
          logger.i('タスク一覧の取得に失敗しました');
          // データがない場合は「データがありません」のメッセージを表示
          return Center(
            child: Column(
              children: [
                Expanded(
                  child: Container(
                    padding: const EdgeInsets.all(32.0),
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      crossAxisAlignment: CrossAxisAlignment.center,
                      spacing: 24,
                      children: [
                        Icon(Icons.warning, color: AppColors.grayLight, size: 48.0),
                        Text(
                          "データの取得に失敗しました。\nネットワーク接続を確認してやり直してください。",
                          style: TextStyle(
                            color: AppColors.grayDark,
                            fontSize: AppFontSizes.md,
                          ),
                          textAlign: TextAlign.start,
                        ),
                      ],
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
                    Navigator.of(context).pop(); // 戻るボタンが押されたときの処理
                    logger.i("戻りました");
                  },
                  label: "戻る",
                ),
                SizedBox(height: 16.0)
              ],
            ),
          );
        } else {
          final tasks = snapshot.data as List<RescueTaskDropdownMenuItem>;
          // ドロップダウンの選択肢を作成
          final dropdownMenuItems = tasks.map((task) => DropdownMenuItem<String>(
            value: task.id.toString(),
            child: Text(task.taskName),
          )).toList();
          // 人が来ないのレスキューを送信する画面
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
                            "人が来ない",
                            style: TextStyle(
                              fontSize: AppFontSizes.lg,
                              fontWeight: FontWeight.bold,
                              color: AppColors.textBlack,
                            ),
                            textAlign: TextAlign.center,
                          ),
                        ),
                        // 発生したタスクの選択
                        Text(
                          "人が来ないのはどのタスクですか？",
                          style: TextStyle(
                            fontSize: AppFontSizes.md,
                            color: AppColors.textBlack,
                          ),
                          // textAlign: TextAlign.center,
                        ),
                        StatefulBuilder(
                          builder: (context, setState) {
                            return CustomDropdownButton<String>(
                              value: selectedTask.id.toString(),
                              items: dropdownMenuItems,
                              isDense: true, // 高さをコンパクトにする
                              onChanged: (value) {
                                setState(() {
                                  // 選択されたタスクを更新
                                  selectedTask = tasks.firstWhere(
                                    (task) => task.id.toString() == value,
                                    orElse: () => RescueTaskDropdownMenuItem(
                                      id: 0, // タスク外
                                      taskName: 'タスクを選択してください',
                                    ),
                                  );
                                });
                                logger.i("選択されたタスク: ${selectedTask.taskName}");
                              },
                              hintText: 'タスクを選択してください',
                            );
                          },
                        ),
                        SizedBox(height: 16.0), // スペースを追加
                        // 人数
                        Text(
                          "来ていないのは何人ですか？",
                          style: TextStyle(
                            fontSize: AppFontSizes.md,
                            color: AppColors.textBlack,
                          ),
                          textAlign: TextAlign.center,
                        ),
                        StatefulBuilder(
                          builder: (context, setState) {
                            return CustomTextField(
                              keyboardType: TextInputType.number, // 数字のみのキーボードにする
                              labelText: '人数',
                              hintText: '例：2',
                              onChanged: (String value) {
                                // 入力されたときの処理
                                setState(() {
                                  missingNumber = int.tryParse(value) ?? 0;
                                });
                                logger.i("入力された人数: $value");
                              },
                            );
                          },
                        ),
                        SizedBox(height: 16.0), // スペースを追加
                        // 人を送る場所
                        Text(
                          "どこに人を送れば良いですか？",
                          style: TextStyle(
                            fontSize: AppFontSizes.md,
                            color: AppColors.textBlack,
                          ),
                          textAlign: TextAlign.center,
                        ),
                        StatefulBuilder(
                          builder: (context, setState) {
                            return CustomTextField(
                              labelText: '送り先の場所',
                              hintText: '例：案内所 講義棟',
                              onChanged: (String value) {
                                // 入力されたときの処理
                                setState(() {
                                  place = value;
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
                    isDisabled: isSubmitting,
                    onPressed: () async {
                      setState(() {
                        isSubmitting = true;
                      });
                      // レスキューを送信する
                      final isSuccess = await _sendRescueRequest(
                        context,
                        selectedTask,
                        missingNumber,
                        place,
                      );
                      logger.i("レスキューを送信しました");
                      if(isSuccess){
                        // 最初の画面まで戻る
                        Navigator.of(context).popUntil((route) => route.isFirst);
                      }else{
                        setState(() {
                          isSubmitting = false;
                        });
                      }
                    },
                    label: isSubmitting ? "送信中..." : "送信",
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
      },
    );
  }
}