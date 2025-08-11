// トラブルのレスキューを送信するページ
import 'package:seeft_mobile/configs/importer.dart';
import 'package:seeft_mobile/models/rescue.dart';
import 'package:seeft_mobile/widgets/custom_dropdown_menu.dart';
import 'package:seeft_mobile/widgets/custom_text_field.dart';
import 'package:seeft_mobile/widgets/custom_elevated_button.dart';
import 'package:seeft_mobile/widgets/custom_elevated_button_outlined.dart';
import 'package:seeft_mobile/widgets/custom_snack_bar.dart';
import 'package:seeft_mobile/widgets/custom_error_snack_bar.dart';

class RescueRequestTabTroublePage extends StatelessWidget {
  RescueRequestTabTroublePage({
    Key? key,
  }) : super(key: key);
  
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
          logger.i('First item: ${res}');
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
      
      // タスク外の選択肢を追加
      taskList.insert(0, RescueTaskDropdownMenuItem(
        id: 0,  // (実際の「タスク外」のidは3だが、DropDown内でのidの重複を避けるために0としている)
        taskName: 'タスク外',
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
  
  // トラブルのレスキューを送信する関数
  Future<bool> _sendRescueRequest(
    BuildContext context, 
    RescueTaskDropdownMenuItem selectedTask,
    String detail,
    String place,
  ) async {
    if(
      detail == "" || 
      place == ""
    ) {
      logger.e('Invalid input data');
      showCustomErrorSnackBar(context, "データが入力されていません");
      return false;
    }
    try {
      final userId = await store.getUserID();
      logger.i('Sending trouble report:'
        '\nDetail: $detail'
        '\nplace: $place'
        '\nSelected Task: ${selectedTask.taskName}');
      
      // API呼び出し
      await api.postTroubleRescue(
        userId, 
        selectedTask.id != 0 ? selectedTask.id : 3, // idが0なら「タスク外(実際のidが3)」で送信する
        place, 
        detail
      );

      logger.i('Trouble report sent successfully.');
      showCustomSnackBar(context, "レスキューを送信しました");
      return true;
    } catch (e) {
      logger.e('Failed to send trouble report: $e');
      showCustomErrorSnackBar(context, "レスキューを送信できませんでした");
      return false;
    }
  }

  @override
  Widget build(BuildContext context) {
    // 選択されたタスクを格納する変数
    RescueTaskDropdownMenuItem _selectedTask = RescueTaskDropdownMenuItem(
      id: 0, // 初期値はタスク外(実際の「タスク外」のidは3だが、DropDown内でのidの重複を避けるために0としている)
      taskName: 'タスク外',
    );
    // トラブルの詳細を格納する変数
    String _detail = '';
    // 発生場所を格納する変数
    String _place = '';
    
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
                SizedBox(height: 16.0), // スペースを追加
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
                SizedBox(height: 16.0), // スペースを追加
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
          final dropdownMenuEntries = tasks.map((task) => DropdownMenuEntry<String>(
            value: task.id.toString(),
            label: task.taskName,
          )).toList();
          // トラブルのレスキューを送信する画面
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
                            "トラブル",
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
                          "どのタスクですか？",
                          style: TextStyle(
                            fontSize: AppFontSizes.md,
                            color: AppColors.textBlack,
                          ),
                          // textAlign: TextAlign.center,
                        ),
                        StatefulBuilder(
                          builder: (context, setState) {
                            return CustomDropdownMenu<String>(
                              dropdownMenuEntries: dropdownMenuEntries,
                              initialSelection: _selectedTask.id.toString(),
                              onSelected: (value) {
                                setState(() {
                                  // 選択されたタスクを更新
                                  _selectedTask = tasks.firstWhere(
                                    (task) => task.id.toString() == value,
                                    orElse: () => RescueTaskDropdownMenuItem(
                                      id: 0, // タスク外
                                      taskName: 'タスク外',
                                    ),
                                  );
                                });
                                print('Selected: $value');
                              },
                            );
                          },
                        ),
                        SizedBox(height: 16.0), // スペースを追加
                        // 発生場所
                        Text(
                          "発生場所はどこですか？",
                          style: TextStyle(
                            fontSize: AppFontSizes.md,
                            color: AppColors.textBlack,
                          ),
                          textAlign: TextAlign.center,
                        ),
                        StatefulBuilder(
                          builder: (context, setState) {
                            return CustomTextField(
                              labelText: '発生場所',
                              hintText: '例：案内所 講義棟',
                              onChanged: (String value) {
                                // 入力されたときの処理
                                setState(() {
                                  _place = value;
                                });
                                logger.i("入力された場所: $value");
                              },
                            );
                          },
                        ),
                        SizedBox(height: 16.0), // スペースを追加
                        // トラブルの詳細
                        Text(
                          "トラブルの詳細を記入してください",
                          style: TextStyle(
                            fontSize: AppFontSizes.md,
                            color: AppColors.textBlack,
                          ),
                          textAlign: TextAlign.center,
                        ),
                        StatefulBuilder(
                          builder: (context, setState) {
                            return CustomTextField(
                              labelText: '詳細',
                              hintText: '例：〇〇の物品がない',
                              onChanged: (String value) {
                                // 入力されたときの処理
                                setState(() {
                                  _detail = value;
                                });
                                logger.i("入力された詳細: $value");
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
              CustomElevatedButton(
                onPressed: () async {
                  // レスキューを送信する
                  final isSuccess = await _sendRescueRequest(
                    context,
                    _selectedTask,
                    _place,
                    _detail
                  );
                  logger.i("トラブルを送信しました");
                  if(isSuccess){
                    // 最初の画面まで戻る
                    Navigator.of(context).popUntil((route) => route.isFirst);
                  }
                },
                label: "送信",
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