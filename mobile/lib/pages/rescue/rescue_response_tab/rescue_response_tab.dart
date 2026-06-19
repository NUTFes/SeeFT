// レスキューの返答を表示するタブ
import 'package:seeft_mobile/configs/importer.dart';
import 'package:seeft_mobile/models/rescue.dart';
import 'package:seeft_mobile/widgets/refresh_button.dart';
import 'package:seeft_mobile/widgets/custom_error_snack_bar.dart';
import 'package:seeft_mobile/widgets/custom_dropdown_button.dart';
import 'package:collection/collection.dart';

// レスキューの一覧を取得する関数
Future<List<dynamic>?> _getRescueResponses(int? userID) async {
  try {
    dynamic res;
    // API呼び出し
    if (userID == null) {
      // userIDがnullの場合は全てのレスキューのレスポンスを取得
      res = await api.getAllRescueResponses();
    } else {
      res = await api.getRescueResponses(userID);
    }
    // final res = await api.getRescueResponses(userID);
    if (res == null || res.isEmpty) {
      logger.i('No rescue responses found.');
      return []; // レスキューのレスポンスがない場合は空のリストを返す
    }
    logger.i('Rescue responses fetched successfully: $res');
    return (res as List<dynamic>);
    // return (res as List).map((item) => RescueResponse.fromJson(item)).toList();
  } catch (e) {
    logger.e('Failed to fetch rescue report: $e');
    return null; // エラーが発生した場合はnullを返す
  }
}

// キャッシュからデータをロードする関数
List<dynamic>? _getCashedRescueResponses(int? userID) {
  // キャッシュからデータを取得
  logger.i('=== キャッシュからデータを取得します ===');
  List<dynamic>? cachedData;
  if (userID == null) {
    // userIDがnullの場合は全てのレスキューのレスポンスを取得
    cachedData = rescueBox.get('all_rescue_responses');
  } else {
    cachedData = rescueBox.get('rescue_responses_by_${userID}');
  }
  // final List<dynamic>? cachedData = rescueBox.get('rescue_responses_by_${userID}');
  logger.i('キャッシュデータの取得に成功しました: ${cachedData != null ? cachedData.length : 'null'} items');
  
  // キャッシュデータがない場合は表示データをnullに設定する
  if (cachedData == null) {
    logger.e('$userID のキャッシュデータがありません。');
    return null;
  }
  
  // キャッシュデータがある場合はそれを使用
  logger.i('$userID のキャッシュデータを発見しました。');
  
  return cachedData; // キャッシュデータを返す
}

class RescueResponseTab extends StatefulWidget {
  const RescueResponseTab({super.key});

  @override
  _RescueResponseTabState createState() => _RescueResponseTabState();
}

class _RescueResponseTabState extends State<RescueResponseTab> {
  String _selectedFilter = 'my'; // 選択されたフィルターの値
  bool _isLoading = true; // 読み込み中のフラグ
  List<RescueResponse>? _rescueResponses; // レスキューのレスポンスを格納する変数
  final deepEq = DeepCollectionEquality.unordered().equals; // キャッシュデータとフェッチデータの比較用の関数
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
  Future<void> _loadRescueResponses(int? userID) async {
    setState(() => _isLoading = true);  // ロード中フラグをtrueに設定
    
    // キャッシュからデータを取得
    final List<dynamic>? cachedData = _getCashedRescueResponses(userID);
    // キャッシュデータをList<RescueResponse>に変換
    final List<RescueResponse>? cachedRescueResponses = cachedData != null
        ? cachedData.map((item) => RescueResponse.fromJson(item)).toList()
        : null;
    
    // 表示データをキャッシュデータで更新
    setState(() {
      _rescueResponses = cachedRescueResponses;
    });
    
    // サーバーからデータを取得
    final List<dynamic>? fetchedData = await _getRescueResponses(userID);
    
    // キャッシュデータとフェッチデータを比較
    if (fetchedData == null || deepEq(fetchedData, cachedData)) {
      // フェッチデータがnullまたはキャッシュデータと同じ場合は、キャッシュデータを使用
      fetchedData == null
        ? {
          logger.w('フェッチデータが空です。キャッシュを使用します。'),
          showCustomErrorSnackBar(context, 'データの取得に失敗しました。'), // スナックバーでエラーメッセージを表示
        }
        : logger.w('フェッチデータが既に最新です。キャッシュを使用します。');
      setState(() => _isLoading = false);   // ロード中フラグをfalseに設定
      return;
    }
    
    // フェッチデータが新しい場合はキャッシュデータと表示データを更新
    if (userID == null) {
      // ユーザIDがnullの場合は全てのレスキューのレスポンスを更新
      rescueBox.put('all_rescue_responses', fetchedData); // hiveのキャッシュデータをフェッチデータで更新
    } else {
      // ユーザIDが指定されている場合は、ユーザIDに紐づくキャッシュデータを更新
      logger.i('Updating cache for user ID: $userID');
    }
    // rescueBox.put('rescue_responses_by_${userID}', fetchedData); // hiveのキャッシュデータをフェッチデータで更新
    
    // フェッチデータをList<RescueResponse>に変換
    final List<RescueResponse> fetchedRescueResponses = fetchedData
        .map((item) => RescueResponse.fromJson(item))
        .toList();
    
    // 表示データをフェッチデータで更新
    setState(() {
      _rescueResponses = fetchedRescueResponses;
      _isLoading = false; // ロード中フラグをfalseに設定
      logger.i('Rescue responses loaded successfully: $_rescueResponses');
    });
    return;
  }
  
  // データを更新する処理
  Future<void> _handleRefresh() async {
    // レスキューのレスポンスを再取得
    if(_selectedFilter == 'my') {
      // 自分が送信したレスキューのみを取得
      await _loadRescueResponses(_userID);
    } else {
      // 全てのレスキューを取得
      await _loadRescueResponses(null);
    }
    // await _loadRescueResponses(_userID);
  }
  
  
  @override
  Widget build(BuildContext context) {
    return Container(
      padding: EdgeInsets.only(top: 16.0, left: 32.0, right: 32.0, bottom: 16.0),
      child: Column(
        spacing: 8.0,  // 子要素間のスペース
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            spacing: 8.0,
            children: [
              Text(
                "表示範囲",
                style: TextStyle(
                  fontSize: AppFontSizes.md,
                  color: AppColors.grayDark,
                ),
              ),
              Expanded(
                child: CustomDropdownButton(
                  value: _selectedFilter,
                  items: [
                    DropdownMenuItem(
                      value: 'my',
                      child: Text("自分が送信したレスキュー"),
                    ),
                    DropdownMenuItem(
                      value: 'all',
                      child: Text("全てのレスキュー"),
                    ),
                  ],
                  onChanged: (value) {
                    // 選択された値に応じて表示データをフィルタリング
                    setState(() {
                      _selectedFilter = value ?? 'my';
                      _handleRefresh(); // フィルター変更時にデータを再取得
                    });
                  },
                  isDense: true, // 高さをコンパクトにする
                ),
              ),
            ],
          ),
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
            onPressed: _handleRefresh, 
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
                  + "発生時刻: " + res.time;
        break;
      // 質問
      case 'question':
        final res = response as QuestionRescueResponse;
        titleRescue = "【質問】" + res.content.question;
        subTitle = "対応番号: Q" + res.id.toString() + "\n"
                  + "送信者: " + res.userName + "\n"
                  + "発生時刻: " + res.time;
        break;
      // 人が来ない
      case 'shorthanded':
        final res = response as ShorthandedRescueResponse;
        titleRescue = "【人が来ない】" + res.content.task + "（" + res.content.missingNumber.toString() + "人）";
        subTitle = "対応番号: S" + res.id.toString() + "\n"
                  + "送信者: " + res.userName + "\n"
                  + "送り先の場所: " + res.content.place + "\n"
                  + "発生時刻: " + res.time;
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