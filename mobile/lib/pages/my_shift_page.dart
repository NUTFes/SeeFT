import 'package:seeft_mobile/configs/importer.dart';
import 'package:seeft_mobile/widgets/shift_card.dart';
import 'package:collection/collection.dart';


Future<List<dynamic>?> _getShiftCardDataList(int userID, int dayID, int weatherID) async {
  logger.i('=== API Call Started ===' 
    '\n'
    'Fetching shift cards for userID: $userID, dayID: $dayID, weatherID: $weatherID');
  
  try {
    var res = await api.getShiftCardsByUserAndDateAndWeather(
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
    
    return res;
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
  final List<dynamic>? cashedData = shiftCardBox.get('${dayID}_${weatherID}');
  logger.i('キャッシュデータの取得に成功しました: ${cashedData != null ? cashedData.length : 'null'} items');
  
  // キャッシュデータがない場合は表示データをnullに設定する
  if (cashedData == null) {
    logger.e('$dayID, $weatherID のキャッシュデータがありません。');
    return null;
  }
  
  // キャッシュデータがある場合はそれを使用
  logger.i('$dayID, $weatherID のキャッシュデータを発見しました。');
  
  return cashedData; // キャッシュデータを返す
}

class MyShiftPage extends StatefulWidget {
  @override
  _MyShiftPageState createState() => _MyShiftPageState();
}

class _MyShiftPageState extends State<MyShiftPage>
    with TickerProviderStateMixin {
  int _selectedWeatherIndex = 1;  // 天気の選択肢のインデックス(1:晴れ, 2:雨)
  int _selectedDayIndex = 1;      // 日付の選択肢のインデックス(1:準備日, 2:1日目, 3:2日目, 4:片付け日)
  ShiftCardDataList? shiftCardDataList; // ShiftCardDataListを格納する変数
  // Map<String, ShiftCardDataList>? allShiftCardDataList; // 全てのシフトカードデータを格納する変数(キーはhiveのキーと同じ)
  Timer? _debounceTimer;
  Timer? _tabSwitchDebounceTimer; // タブ切り替え用のデバウンスタイマー
  bool isLoading = false; // ロード中かどうかのフラグ
  final int debounceTime = 500; // デバウンス時間を設定
  final deepEq = DeepCollectionEquality.unordered().equals; // シフトカードデータの比較用の関数
  
  // 天気ごとのweatherID
  final Map<int, String> _weatherOptions = {
    1: "晴れ",
    2: "雨"
  };
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
    // タブのコントローラーを初期化
    _tabController = TabController(
      length: _dayOptions.length,
      vsync: this
    );
    _tabController.addListener(_handleTabChange); // タブの変更を監視するリスナーを追加
    
    // ユーザIDと初期タブのデータを初期化
    _initialize();
    super.initState();
  }
  
  // ユーザIDと初期タブのデータを初期化する関数(非同期処理は直接initStateで行えないため、別関数に分ける)
  Future<void> _initialize() async {
    // ユーザIDを取得
    _userID = await store.getUserID();
    logger.i('User ID: $_userID');
    // 初期タブのデータを初期化
    await _loadShiftCardDataListForCurrentTab(_userID, _selectedDayIndex, _selectedWeatherIndex);
  }
  
  // ウィジェットが破棄されるときの処理
  @override
  void dispose() {  
    _tabController.removeListener(_handleTabChange);
    _tabController.dispose();
    _debounceTimer?.cancel(); // デバウンスタイマーをキャンセル
    _tabSwitchDebounceTimer?.cancel(); // タブ切り替えデバウンスタイマーをキャンセル
    super.dispose();
  }
  
  // タブ(日付)が切り替わったときの処理
  void _handleTabChange() {
    final int newDayIndex = _tabController.index + 1; // 日付の選択状態を取得
    setState(() => _selectedDayIndex = newDayIndex);  // 選択された日付のインデックスを更新
    // シフトカードのデータを更新
    _loadShiftCardDataListForCurrentTab(_userID, newDayIndex, _selectedWeatherIndex);
  }
  
  // SegmentedButton(天気)の選択状態が変わったときの処理
  void _handleWeatherSelectionChanged(Set<Object> newSelection) {
    final int newWeatherIndex = (newSelection.first as int);  // 天気の選択状態を取得
    setState(() => _selectedWeatherIndex = newWeatherIndex);  // 選択された天気のインデックスを更新
    // シフトカードのデータを更新
    _loadShiftCardDataListForCurrentTab(_userID, _selectedDayIndex, newWeatherIndex);
  }
  
  // 更新ボタンが押されたときの処理
  void _handleRefreshPressed() {
    // シフトカードのデータを再ロード
    _loadShiftCardDataListForCurrentTab(_userID, _selectedDayIndex, _selectedWeatherIndex);
  }
  
  // 指定のユーザID、日付ID、天気IDのシフトカードデータリストをロードする関数
  Future<void> _loadShiftCardDataListForCurrentTab(int userID, int dayID, int weatherID) async {
    _tabSwitchDebounceTimer?.cancel();  // 既存のデバウンスタイマーをキャンセル
    setState(() => isLoading = true);   // ロード中フラグをtrueに設定
    
    // キャッシュからデータを取得
    final List<dynamic>? cashedData = _getCashedShiftCardDataList(dayID, weatherID);
    // キャッシュデータをShiftCardDataListに変換
    final ShiftCardDataList? cashedShiftCardDataList = 
      cashedData != null ? ShiftCardDataList.fromJson(cashedData) : null;
    
    setState(() {
      shiftCardDataList = cashedShiftCardDataList;  // 表示データをキャッシュデータで更新
    });
    
    // デバウンス処理
    _tabSwitchDebounceTimer = Timer(Duration(milliseconds: debounceTime), () async{ // 一定時間日付と天気が切り替わらなかった場合に実行される
      // サーバーからデータを取得
      final List<dynamic>? fetchedData = await _getShiftCardDataList(userID, dayID, weatherID);
      
      // キャッシュデータとフェッチデータを比較
      if (fetchedData == null || deepEq(fetchedData, cashedData)) {
        fetchedData == null
          ? {
            logger.w('フェッチデータが空です。キャッシュを使用します。'),
            _showErrorSnackbar('データの取得に失敗しました。'), // スナックバーでエラーメッセージを表示
          }
          : logger.w('フェッチデータが既に最新です。キャッシュを使用します。');
        setState(() => isLoading = false);   // ロード中フラグをfalseに設定
        return; // フェッチデータがnullまたはキャッシュデータと同じ場合は何もしない
      }
      
      // フェッチデータをShiftCardDataListに変換
      final ShiftCardDataList? fetchedShiftCardDataList = ShiftCardDataList.fromJson(fetchedData);
      
      // フェッチデータが新しい場合はキャッシュデータと表示データを更新
      shiftCardBox.put('${dayID}_${weatherID}', fetchedData); // hiveのキャッシュデータをフェッチデータで更新
      setState(() {
        shiftCardDataList = fetchedShiftCardDataList; // 表示データをフェッチデータで更新
        isLoading = false; // ロード中フラグをfalseに設定
      });
    });
  }
  
  // スナックバーでエラーメッセージを表示する関数
  void _showErrorSnackbar(String message) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Row(
          children: [
            Icon(Icons.warning, color: AppColors.textWhite, size: 16),
            const SizedBox(width: 8),
            Text(
              message,
              style: TextStyle(
                color: AppColors.textWhite,
                fontSize: AppFontSizes.sm,
              )),
          ],
        ),
        backgroundColor: AppColors.error,
        duration: Duration(seconds: 2),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
        margin: const EdgeInsets.only(left: 16, right: 16, bottom: 8),
        behavior: SnackBarBehavior.floating,
      ),
    );
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
        actions: [
          Padding(
            padding: const EdgeInsets.only(right: 20.0),
            child: SizedBox(
              width: 200,
              // height: 30,
              child: _weatherSegmentedButton(), // 天気を選択するセグメントボタン
            ),
          ),
        ],
        bottom: PreferredSize(
          preferredSize: const Size.fromHeight(75.0),
          child: _dayTabBar(),  // 日付を選択するタブバー
        ),
      ),
      body: Column(
        children: [
          // 更新ボタン
          Container(
            padding: const EdgeInsets.symmetric(vertical: 8.0),
            child: _refreshButton(),
          ),
          // TabBarView
          Expanded(
            child: TabBarView(
              controller: _tabController,
              children: List.generate(_dayOptions.length, (index) {
                logger.i('=== タブ $index のデータをロード中 ===');
                logger.i('Has data: ${shiftCardDataList != null}');
                if (shiftCardDataList == null) {
                  logger.i('タブ $index: データがありません');
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
                }else {
                  logger.i('Tab $index: Data available with ${shiftCardDataList!.data.length} items');
                  if (shiftCardDataList!.data.isNotEmpty) {
                    logger.i('Tab $index: First item task name: ${shiftCardDataList!.data[0].taskName}');
                  }
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
            ),
          ),
        ],
      ),
    );
  }
  
  // 天気を選択するセグメントボタン
  Widget _weatherSegmentedButton() {
    return SegmentedButton(
      selected: {_selectedWeatherIndex}, // 選択されている天気のインデックスをセット
      onSelectionChanged: _handleWeatherSelectionChanged,
      style: SegmentedButton.styleFrom(
        backgroundColor: AppColors.main,
        selectedBackgroundColor: AppColors.base,
        side: BorderSide(
          color: AppColors.grayLight,
          width: 1.0,
        ),
      ),
      showSelectedIcon: false,
      segments: [
        ButtonSegment(
          label: Text(
            "晴れ",
            style: TextStyle(
              color: _selectedWeatherIndex == (_weatherOptions["晴れ"] ?? 1) ? AppColors.main : AppColors.grayLight,
              fontSize: AppFontSizes.sm,
            ),
          ),
          icon: Icon(
            Icons.sunny,
            color: _selectedWeatherIndex == (_weatherOptions["晴れ"] ?? 1) ? AppColors.main : AppColors.grayLight,
            size: 18,
          ),
          value: _weatherOptions["晴れ"] ?? 1, // 晴れのインデックス
        ),
        ButtonSegment(
          label: Text(
            "雨",
            style: TextStyle(
              color: _selectedWeatherIndex == (_weatherOptions["雨"] ?? 2) ? AppColors.main : AppColors.grayLight,
              fontSize: AppFontSizes.sm,
            ),
          ),
          icon: Icon(
            Icons.cloudy_snowing,
            color: _selectedWeatherIndex == (_weatherOptions["雨"] ?? 2) ? AppColors.main : AppColors.grayLight,
            size: 18,
          ),
          value: _weatherOptions["雨"] ?? 2, // 雨のインデックス
        ),
      ],
    );
  }
  
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
  
  // 更新ボタン
  Widget _refreshButton() {
    return ElevatedButton.icon(
      onPressed: isLoading ? null : _handleRefreshPressed, // ロード中はボタンを無効化
      icon: Icon(
        Icons.refresh,
        color: isLoading? AppColors.grayDark: AppColors.main,
        size: 16,
      ),
      label: Text(
        isLoading? "読み込み中...": "更新",
        style: TextStyle(
          color: isLoading? AppColors.grayDark: AppColors.main,
          fontSize: AppFontSizes.sm,
          // fontWeight: FontWeight.bold,
        ),
      ),
      style: ElevatedButton.styleFrom(
        backgroundColor: AppColors.base,
        disabledBackgroundColor: AppColors.base,
        elevation: 0,
        side: BorderSide(
          color: isLoading? AppColors.grayDark: AppColors.main,
          width: 1.0,
        ),
        padding: const EdgeInsets.symmetric(vertical: 16.0, horizontal: 24.0),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(100.0),
        ),
      ),
    );
  }
}