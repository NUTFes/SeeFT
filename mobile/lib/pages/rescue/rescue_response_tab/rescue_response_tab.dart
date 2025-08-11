// レスキューの返答を表示するタブ
import 'package:seeft_mobile/configs/importer.dart';
import 'package:seeft_mobile/models/rescue.dart';
import 'package:seeft_mobile/widgets/refresh_button.dart';
import 'package:seeft_mobile/widgets/custom_error_snack_bar.dart';

// レスキューの一覧を取得する関数
Future<List<RescueResponse>?> _getRescueResponses(int userID) async {
  try {
    // API呼び出し
    final res = await api.getRescueResponses(userID);
    if (res == null || res.isEmpty) {
      logger.i('No rescue responses found.');
      return []; // レスキューのレスポンスがない場合は空のリストを返す
    }
    logger.i('Rescue responses fetched successfully: $res');
    return (res as List).map((item) => RescueResponse.fromJson(item)).toList();
  } catch (e) {
    logger.e('Failed to fetch rescue report: $e');
    return null; // エラーが発生した場合はnullを返す
  }
}

class RescueResponseTab extends StatefulWidget {
  const RescueResponseTab({Key? key}) : super(key: key);

  @override
  _RescueResponseTabState createState() => _RescueResponseTabState();
}

class _RescueResponseTabState extends State<RescueResponseTab> {
  bool _isLoading = true; // 読み込み中のフラグ
  List<RescueResponse>? _rescueResponses; // レスキューのレスポンスを格納する変数
  late int _userID; // ユーザIDを格納する変数
  
  
  // ウィジェットの初期化時の処理
  @override
  void initState() {
    // 非同期処理を分離した関数を呼び出す
    _initialize();
    super.initState();
  }
  
  // ウィジェットを初期化する関数(非同期処理は直接initStateで行えないため分離)
  Future<void> _initialize() async {
    // ユーザIDを取得
    _userID = await store.getUserID();
    logger.i('User ID: $_userID');
    // タブのデータを初期化
    await _loadRescueResponses(_userID);
  }
  
  
  // 指定のユーザIDのレスキューのレスポンスを取得する関数
  Future<void> _loadRescueResponses(int userID) async {
    setState(() => _isLoading = true);  // ロード中フラグをtrueに設定
    // API呼び出し
    final fetchedData = await _getRescueResponses(userID);
    
    // if (fetchedData == null || fetchedData.isEmpty) {
    if (fetchedData == null) {
      // レスキューのレスポンスがない場合はnullを設定
      showCustomErrorSnackBar(context, "データの取得に失敗しました。");
      setState(() {
        _rescueResponses = null;
        _isLoading = false; // ロード中フラグをfalseに設定
      });
      // データがない場合は処理を終了
      return;
    }
    
    // レスキューのレスポンスを格納
    setState(() {
      _rescueResponses = fetchedData; // レスキューのレスポンスを格納
      _isLoading = false; // ロード中フラグをfalseに設定
      logger.i('Rescue responses loaded successfully: $_rescueResponses');
    });
    return;
  }
  
  // 更新ボタンが押されたときの処理
  Future<void> _handleRefreshPressed() async {
    // レスキューのレスポンスを再取得
    await _loadRescueResponses(_userID);
  }
  
  
  @override
  Widget build(BuildContext context) {
    return Container(
      padding: EdgeInsets.only(top: 32.0, left: 32.0, right: 32.0, bottom: 16.0),
      child: Column(
        spacing: 8.0,  // 子要素間のスペース
        children: [
          if (_rescueResponses == null || _rescueResponses!.isEmpty) ...[
            Expanded(
              child: Center(
                child: Container(
                  padding: EdgeInsets.symmetric(horizontal: 32.0, vertical: 8.0),
                  child: Text.rich(
                    TextSpan(
                      children: [
                        TextSpan(
                          text: "送信したレスキューはまだありません。\n",
                          style: TextStyle(
                            color: AppColors.textBlack,
                            fontSize: AppFontSizes.md,
                          ),
                        ),
                        TextSpan(
                          text: "レスキューを送信するとここに本部からの返答が表示されます。",
                          style: TextStyle(
                            color: AppColors.grayDark,
                            fontSize: AppFontSizes.md,
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ),
            ),
          ] else ...[
            // レスキューのレスポンスがある場合はリスト表示
            Expanded(
              child: ListView.separated(
                itemBuilder: (context, index) {
                  final response = _rescueResponses![index];
                  // レスキューのレスポンスを表示するウィジェット
                  return _rescueResponseWidget(response);
                },
                separatorBuilder: (context, index) {
                  return Divider(
                    color: AppColors.grayLight,
                    thickness: 1.0,
                  );
                },
                itemCount: _rescueResponses!.length,
              ),
            ),
          ],
          Divider(
            color: AppColors.grayLight,
            thickness: 1.0,
          ),
          Text(
            "最新の情報を取得するには、\n「更新」ボタンを押してください。",
            style: TextStyle(
              fontSize: AppFontSizes.sm,
              color: AppColors.grayDark,
            ),
          ),
          // 更新ボタン
          RefreshButton(
            onPressed: _handleRefreshPressed, 
            isLoading: _isLoading
          ),
        ],
      ),
    );
  }
  
  // レスキューのレスポンスを表示するウィジェット
  Widget _rescueResponseWidget(RescueResponse response) {
    String titleRescue = "";
    String titleResponse = response.response == "" ? "本部からの返答はまだありません。" : "本部からの返答: " + response.response;
    String subTitle = "";
    // レスキューのタイプに応じてタイトルとサブタイトルを設定
    switch(response.type) {
      // トラブル
      case 'trouble':
        final res = response as TroubleRescueResponse;
        titleRescue = "【トラブル】" + res.content.detail;
        subTitle = "対応番号: T" + res.id.toString() + "\n"
                  + "送信者: " + res.userName + "\n"
                  + "発生タスク: " + res.content.task + "\n"
                  + "発生場所: " + res.content.place + "\n"
                  + "発生時刻: " + res.time.toIso8601String();
        break;
      // 質問
      case 'question':
        final res = response as QuestionRescueResponse;
        titleRescue = "【質問】" + res.content.question;
        subTitle = "対応番号: Q" + res.id.toString() + "\n"
                  + "送信者: " + res.userName + "\n"
                  + "発生時刻: " + res.time.toIso8601String();
        break;
      // 人が来ない
      case 'shorthanded':
        final res = response as ShorthandedRescueResponse;
        titleRescue = "【人が来ない】" + res.content.task + "（" + res.content.missingNumber.toString() + "人）";
        subTitle = "対応番号: S" + res.id.toString() + "\n"
                  + "送信者: " + res.userName + "\n"
                  + "送り先の場所: " + res.content.place + "\n"
                  + "発生時刻: " + res.time.toIso8601String();
        break;
    }
    return ListTile(
      title: Text.rich(
        TextSpan(
          children: [
            TextSpan(
              text: titleRescue + "\n",
              style: TextStyle(
                fontSize: AppFontSizes.md,
                color: AppColors.textBlack,
              ),
            ),
            TextSpan(
              text: titleResponse,
              style: TextStyle(
                fontSize: AppFontSizes.md,
                fontWeight: FontWeight.bold,
                color: response.response == "" ? AppColors.grayDark : AppColors.main,
              ),
            ),
          ],
        ),
      ),
      subtitle: Text(
        subTitle,
        style: TextStyle(
          fontSize: AppFontSizes.sm,
          color: AppColors.grayDark,
        ),
      ),
      leading: _statusIcon(response.status),
    );
  }
  
  // レスキューのレスポンスのステータスに応じたアイコンを表示するウィジェット
  Widget _statusIcon(String status) {
    Color _iconColor = AppColors.error;
    String _iconText = "未対応";
    switch (status) {
      case 'done':
        _iconColor = Color(0xFF325c23);
        _iconText = "対応済";
        break;
      case 'inProgress':
        _iconColor = Color(0xFF3A73E6);
        _iconText = "対応中";
        break;
      case 'todo':
        _iconColor = AppColors.error;
        _iconText = "未対応";
        break;
    }
    
    return Container(
      padding: EdgeInsets.all(4.0),
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(8.0),
        border: Border.all(color: _iconColor, width: 2.0),
      ),
      child: Text(
        _iconText,
        style: TextStyle(
          color: _iconColor,
          fontWeight: FontWeight.bold,
          fontSize: AppFontSizes.sm,
        ),
      ),
    );
  }
}