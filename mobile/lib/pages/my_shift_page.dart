import 'package:seeft_mobile/configs/importer.dart';
import 'package:seeft_mobile/widgets/shift_card.dart';
import 'package:collection/collection.dart';
import 'package:seeft_mobile/widgets/custom_error_snack_bar.dart';
import 'package:seeft_mobile/widgets/refresh_button.dart';
import 'package:seeft_mobile/widgets/review_bottom_sheet.dart';


Future<List<dynamic>?> _getShiftCardDataList(int userID, int dayID, int weatherID) async {
  logger.i('=== API Call Started ===' 
    '\n'
    'Fetching shift cards for userID: $userID, dayID: $dayID, weatherID: $weatherID');
  
  try {
    final res = await api.getShiftCardsByUserAndDateAndWeather(
      userID,
      dayID,
      weatherID,
    );
    
    logger.i('=== API Response Received ===');
    // logger.i('Response: $res');
    
    if (res is List) {
      logger.i('Response is List with ${res.length} items');
      if (res.isNotEmpty) {
        logger.i('First item: ${res[0]}');
      }
    }
    
    return res as List<dynamic>;
  } catch (err) {
    logger.e('=== API Error ==='
      '\n'
      'Failed to fetch shift cards for userID: $userID, dayID: $dayID, weatherID: $weatherID'
      '\n'
      'Error: $err');
    
    return null; // エラー時はnullを返す
  }
}

// キャッシュからデータをロードする関数
List<dynamic>? _getCashedShiftCardDataList(int dayID, int weatherID) {
  // キャッシュからデータを取得
  logger.i('=== キャッシュからデータを取得します ===');
  final List<dynamic>? cachedData = shiftCardBox.get('${dayID}_${weatherID}');
  logger.i('キャッシュデータの取得に成功しました: ${cachedData != null ? cachedData.length : 'null'} items');
  
  // キャッシュデータがない場合は表示データをnullに設定する
  if (cachedData == null) {
    logger.e('$dayID, $weatherID のキャッシュデータがありません。');
    return null;
  }
  
  // キャッシュデータがある場合はそれを使用
  logger.i('$dayID, $weatherID のキャッシュデータを発見しました。');
  
  return cachedData; // キャッシュデータを返す
}

class MyShiftPage extends StatefulWidget {
  @override
  _MyShiftPageState createState() => _MyShiftPageState();
}

class _MyShiftPageState extends State<MyShiftPage>
    with TickerProviderStateMixin {
  int _selectedDayID = shiftCardBox.get('selectedDayID')?? 1;         // 日付の選択状態(1:準備日, 2:1日目, 3:2日目, 4:片付け日), Hiveから最後に見ていた日付を取得
  // int _selectedWeatherID = shiftCardBox.get('selectedWeatherID')?? 1; // 天気の選択状態(1:晴れ, 2:雨), Hiveから最後に見ていた天気を取得
  ShiftCardDataList? shiftCardDataList; // ShiftCardDataListを格納する変数
  // Map<String, ShiftCardDataList>? allShiftCardDataList; // 全てのシフトカードデータを格納する変数(キーはhiveのキーと同じ)
  bool isLoading = false;       // ロード中かどうかのフラグ
  Timer? _debounceTimer;        // タブとセグメントボタン切り替え時のデバウンスタイマー
  final int debounceTime = 500; // デバウンス時間を設定
  final deepEq = DeepCollectionEquality.unordered().equals; // シフトカードデータの比較用の関数
  
  // // 天気ごとのweatherID
  // final Map<int, String> _weatherOptions = {
  //   1: "晴れ",
  //   2: "雨"
  // };
  // 日付ごとのdateID
  final Map<int, String> _dayOptions = {
    1: "準備日",
    2: "1日目",
    3: "2日目",
    4: "片付け日"
  };
  
  late int _userID; // ユーザIDを格納する変数
  late TabController _tabController; // タブのコントローラーを格納する変数
  
  // ウィジェットの初期化時の処理
  @override
  void initState() {
    // 非同期処理を分離した関数を呼び出す
    _initialize();
    super.initState();
  }
  
  // ウィジェットを初期化する関数(非同期処理は直接initStateで行えないため分離)
  Future<void> _initialize() async {
    logger.i('MyShiftPage is being initialized.');
    // タブのコントローラーを初期化
    _tabController = TabController(
      length: _dayOptions.length,
      initialIndex: _selectedDayID - 1, // 初期タブを設定(最後に見ていた日付)
      vsync: this
    );
    _tabController.addListener(_handleTabChange); // タブの変更を監視するリスナーを追加
    // ユーザIDを取得
    _userID = await store.getUserID();
    logger.i('User ID: $_userID');
    // 初期タブのデータを初期化
    // await _loadShiftCardDataList(_userID, _selectedDayID, _selectedWeatherID);
    await _loadShiftCardDataList(_userID, _selectedDayID, 1); // 天気は初期化時に晴れ(1)で固定

  }
  
  // ウィジェットが破棄されるときの処理
  @override
  void dispose() {  
    // 非同期処理を分離した関数を呼び出す
    _dispose();
    super.dispose();
  }
  
  // ウィジェットが破棄されるときの処理(非同期処理は直接disposeで行えないため分離)
  Future<void> _dispose() async {  
    logger.i('MyShiftPage is being disposed.');
    // 選択された日付と天気のIDをHiveに保存
    // await shiftCardBox.put('selectedDayID', _selectedDayID);
    // await shiftCardBox.put('selectedWeatherID', _selectedWeatherID);
    // 選択された日付IDをHiveに保存（天気は廃止）
    await shiftCardBox.put('selectedDayID', _selectedDayID);
    logger.i('Stored Day ID: $_selectedDayID');
    // タブコントローラーのリスナーを削除し、コントローラーを破棄
    _tabController.removeListener(_handleTabChange);
    _tabController.dispose();
    // デバウンスタイマーをキャンセル    
    _debounceTimer?.cancel();
  }
  
  // タブ(日付)が切り替わったときの処理
  void _handleTabChange() {
    final int newDayID = _tabController.index + 1; // 日付の選択状態を取得
    setState(() {
      _selectedDayID = newDayID;  // 選択された日付のインデックスを更新
    });
    // シフトカードのデータを更新
    // _loadShiftCardDataList(_userID, newDayID, _selectedWeatherID);
    _loadShiftCardDataList(_userID, newDayID, 1); // 天気はタブ切り替え時に晴れ(1)で固定
  }
  
  // SegmentedButton(天気)の選択状態が変わったときの処理
  // void _handleWeatherSelectionChanged(Set<Object> newSelection) {
  //   final int newWeatherID = (newSelection.first as int);  // 天気の選択状態を取得
  //   setState(() {
  //     _selectedWeatherID = newWeatherID;  // 選択された天気のインデックスを更新
  //   });
  //   // シフトカードのデータを更新
  //   _loadShiftCardDataList(_userID, _selectedDayID, newWeatherID);
  // }
  
  // 更新ボタンが押されたときの処理
  void _handleRefreshPressed() {
    // シフトカードのデータを再ロード
    // _loadShiftCardDataList(_userID, _selectedDayID, _selectedWeatherID);
    _loadShiftCardDataList(_userID, _selectedDayID, 1); // 天気は更新時に晴れ(1)で固定
  }
  
  // 指定のユーザID、日付ID、天気IDのシフトカードデータリストをロードする関数
  Future<void> _loadShiftCardDataList(int userID, int dayID, int weatherID) async {
    _debounceTimer?.cancel();           // 既存のデバウンスタイマーをキャンセル
    setState(() => isLoading = true);   // ロード中フラグをtrueに設定
    
    // 日付が準備日と片付け日であれば天気を1(晴れ)に固定
    // if (dayID == 1 || dayID == 4) {
    //   weatherID = 1;
    // }
    
    // キャッシュからデータを取得
    final List<dynamic>? cachedData = _getCashedShiftCardDataList(dayID, weatherID);
    // キャッシュデータをShiftCardDataListに変換
    final ShiftCardDataList? cashedShiftCardDataList = 
      cachedData != null ? ShiftCardDataList.fromJson(cachedData) : null;
    
    // 表示データをキャッシュデータで更新
    setState(() {
      shiftCardDataList = cashedShiftCardDataList;
    });
    
    // デバウンス処理
    _debounceTimer = Timer(Duration(milliseconds: debounceTime), () async{ // 一定時間日付と天気が切り替わらなかった場合に実行される
      // サーバーからデータを取得
      final List<dynamic>? fetchedData = await _getShiftCardDataList(userID, dayID, weatherID);
      
      // キャッシュデータとフェッチデータを比較
      if (fetchedData == null || deepEq(fetchedData, cachedData)) {
        // フェッチデータがnullまたはキャッシュデータと同じ場合は、キャッシュデータを使用
        fetchedData == null
          ? {
            logger.w('フェッチデータが空です。キャッシュを使用します。'),
            showCustomErrorSnackBar(context, 'データの取得に失敗しました。'), // スナックバーでエラーメッセージを表示
          }
          : logger.w('フェッチデータが既に最新です。キャッシュを使用します。');
        setState(() => isLoading = false);   // ロード中フラグをfalseに設定
        
        // ここにレビューの処理を挟む
        _showReviewFormIfNeeded(cashedShiftCardDataList, dayID);

        return;
      }
      // フェッチデータが新しい場合はキャッシュデータと表示データを更新
      shiftCardBox.put('${dayID}_${weatherID}', fetchedData); // hiveのキャッシュデータをフェッチデータで更新
      
      // フェッチデータをShiftCardDataListに変換
      final ShiftCardDataList? fetchedShiftCardDataList = ShiftCardDataList.fromJson(fetchedData);
      
      // ここにレビューの処理を挟む
      _showReviewFormIfNeeded(fetchedShiftCardDataList, dayID);

      // 表示データをフェッチデータで更新
      setState(() {
        shiftCardDataList = fetchedShiftCardDataList;
        isLoading = false; // ロード中フラグをfalseに設定
      });
    });
  }
  
  // レビューを表示する
  void _showReviewFormIfNeeded(ShiftCardDataList? shiftCardDataList, int dayID) {
    if(shiftCardDataList == null) {
      print("シフトカードデータがありません. レビューを表示しません.");
      return;
    }
    
    // 現在時刻を取得
    DateTime now = DateTime.now();
    
    // dayIDから日付を取得
    String targetDate = '2025-09-12';
    switch (dayID) {
      case 1:
        targetDate = constant.nutfesPreparationDay;
        break;
      case 2:
        targetDate = constant.nutfesDay1;
        break;
      case 3:
        targetDate = constant.nutfesDay2;
        break;
      case 4:
        targetDate = constant.nutfesTidyingUpDay;
        break;
    }

    // 各シフトカードに対するレビュー処理
    shiftCardDataList.data.forEach((shiftCard) {
      // シフトカードのタスクが既にレビュー済みかどうかを確認
      final _isReviewed = reviewedTaskNameBox.get(shiftCard.taskName, defaultValue: false) == true;
      if(_isReviewed){
        print("タスク「${shiftCard.taskName}」は既にレビュー済みです. レビューを表示しません。");
        return;
      }
      print("タスク「${shiftCard.taskName}」は未レビューです. レビューを表示します。");
      
      // シフトカードからタスクの終了時刻を取得
      final String endTime = shiftCard.endTime.padLeft(5, '0'); // 1桁時間の場合に備えて0埋め
      DateTime shiftEndTime = DateTime.parse("$targetDate " + endTime);
      print("現在時刻: $now, シフト終了時刻: $shiftEndTime");
      
      // 対象のタスクが終了しているかどうかを判定
      final _isFinished = now.isAfter(shiftEndTime);
      if (_isFinished) {
        // 各シフトカードに対するレビュー処理を実行
        ReviewBottomSheet.show(
          context,
          shiftCard.taskName,
          _userID
        );
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.base,
      appBar: AppBar(
        title: Text("マイシフト", 
            style: TextStyle(
              color: AppColors.textWhite,
              fontSize: AppFontSizes.lg,
              fontWeight: FontWeight.bold,
            ),
          ),
        centerTitle: false,
        toolbarHeight: 63,
        backgroundColor: AppColors.main,
        foregroundColor: AppColors.base,
        // actions: [
        //   Padding(
        //     padding: const EdgeInsets.only(right: 20.0),
        //     child: SizedBox(
        //       // height: 30,
        //       width: 200,
        //       // 天気を選択するセグメントボタン
        //       child: _weatherSegmentedButton(),
        //     ),
        //   ),
        // ],
        // actions: [], // 天気セグメントボタンは廃止
        bottom: PreferredSize(
          preferredSize: const Size.fromHeight(75.0),
          // 日付を選択するタブバー
          child: _dayTabBar(),
        ),
      ),
      body: Column(
        children: [
          Container(
            padding: const EdgeInsets.symmetric(vertical: 8.0),
            // 更新ボタン
            child: RefreshButton(
              onPressed: _handleRefreshPressed, 
              isLoading: isLoading
            ),
          ),
          Expanded(
            // タブの内容を表示するTabBarView
            child: _tabBarView(),
          ),
        ],
      ),
    );
  }
  
  // 天気を選択するセグメントボタン
  // Widget _weatherSegmentedButton() {
  //   return SegmentedButton(
  //     selected: {_selectedWeatherID}, // 選択されている天気のインデックスをセット
  //     onSelectionChanged: _handleWeatherSelectionChanged,
  //     style: SegmentedButton.styleFrom(
  //       backgroundColor: AppColors.main,
  //       selectedBackgroundColor: AppColors.base,
  //       side: BorderSide(
  //         color: AppColors.grayLight,
  //         width: 1.0,
  //       ),
  //     ),
  //     showSelectedIcon: false,
  //     segments: [
  //       ButtonSegment(
  //         label: Text(
  //           "晴れ",
  //           style: TextStyle(
  //             color: _selectedWeatherID == (_weatherOptions["晴れ"] ?? 1) ? AppColors.main : AppColors.grayLight,
  //             fontSize: AppFontSizes.sm,
  //           ),
  //         ),
  //         icon: Icon(
  //           Icons.sunny,
  //           color: _selectedWeatherID == (_weatherOptions["晴れ"] ?? 1) ? AppColors.main : AppColors.grayLight,
  //           size: 18,
  //         ),
  //         value: _weatherOptions["晴れ"] ?? 1, // 晴れのインデックス
  //       ),
  //       ButtonSegment(
  //         label: Text(
  //           "雨",
  //           style: TextStyle(
  //             color: _selectedWeatherID == (_weatherOptions["雨"] ?? 2) ? AppColors.main : AppColors.grayLight,
  //             fontSize: AppFontSizes.sm,
  //           ),
  //         ),
  //         icon: Icon(
  //           Icons.cloudy_snowing,
  //           color: _selectedWeatherID == (_weatherOptions["雨"] ?? 2) ? AppColors.main : AppColors.grayLight,
  //           size: 18,
  //         ),
  //         value: _weatherOptions["雨"] ?? 2, // 雨のインデックス
  //       ),
  //     ],
  //   );
  // }
  // 天気選択UIは削除（常に晴れ=1 を使用）
  
  // 日付を選択するタブバー
  Widget _dayTabBar() {
    return TabBar(
      isScrollable: false,
      tabs: _dayOptions.entries.map((entry) {
        return Tab(
          text: entry.value, // 日付のラベルを使用
          height: 35,
        );
      }).toList(),
      controller: _tabController,
      padding: const EdgeInsets.only(top: 20.0, bottom: 20.0),
      labelStyle: TextStyle(
        color: AppColors.main,
        fontSize: AppFontSizes.md,
        fontWeight: FontWeight.normal,
      ),
      unselectedLabelStyle: TextStyle(
        color: AppColors.grayLight,
        fontSize: AppFontSizes.md,
        fontWeight: FontWeight.normal,
      ),
      indicator: BoxDecoration(
        color: AppColors.base,
        borderRadius: BorderRadius.circular(100),
      ),
      indicatorSize: TabBarIndicatorSize.label, // 現状のver3.27だとTabBarIndicatorSize.ltabにするとアニメーションのバグが深刻化するのでlabelにしてます
      splashBorderRadius: BorderRadius.circular(100),
      dividerHeight: 0,
    ); 
  }
  
  // タブの内容を表示するTabBarView
  Widget _tabBarView() {
    return TabBarView(
      controller: _tabController,
      children: List.generate(_dayOptions.length, (index) {
        logger.i('=== タブ ${index + 1} のデータをロード中 ==='
          '\n'
          'Has data: ${shiftCardDataList != null}'
        );
        if (shiftCardDataList == null || shiftCardDataList!.data.length == 0) {
          logger.i('タブ ${index + 1}: データがありません');
          // データがない場合は「データがありません」のメッセージを表示
          return Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Icon(Icons.warning, color: AppColors.grayLight, size: 48.0),
                const SizedBox(height: 16.0),
                Text(
                  "データがありません",
                  style: TextStyle(
                    color: AppColors.grayDark,
                    fontSize: AppFontSizes.md,
                  ),
                ),
              ],
            ),
          );
        } else {
          logger.i('Tab ${index + 1}: Data available with ${shiftCardDataList!.data.length} items'
            '\n'
            'First item: ${shiftCardDataList!.data.isNotEmpty ? shiftCardDataList!.data[0].taskName : "No items"}'
          );
          // データがある場合はShiftCardを表示
          return SingleChildScrollView(
            child: Container(
              // padding: const EdgeInsets.all(32.0),
              padding: const EdgeInsets.only(right: 32.0, left: 32.0, top: 8.0, bottom: 32.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: shiftCardDataList!.data.map((data) {
                  return Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      ShiftCard(data: data),
                      const SizedBox(height: 16.0),
                    ],
                  );
                }).toList(),
              ),
            ),
          );
        }
      })
    );
  }
}