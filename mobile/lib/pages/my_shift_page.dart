import 'package:seeft_mobile/configs/importer.dart';
import 'package:seeft_mobile/pages/wait_page.dart';
import 'package:seeft_mobile/widgets/shift_card.dart';
import 'package:collection/collection.dart';
// import 'package:shared_preferences/shared_preferences.dart'; // キャッシュ用
// import 'dart:convert'; // JSONエンコード/デコード用


Future<List<dynamic>?> _getShiftCardDataList(int userID, int dayID, int weatherID) async {
  logger.i('=== API Call Started ===');
  logger.i('Parameters - userID: $userID, dayID: $dayID, weatherID: $weatherID');
  
  try {
    var res = await api.getShiftCardsByUserAndDateAndWeather(
      userID,
      dayID,
      weatherID,
    );
    
    logger.i('=== API Response Received ===');
    // logger.i('Raw API Response: $res');
    logger.i('Response Type: ${res.runtimeType}');
    
    if (res is List) {
      logger.i('Response is List with ${res.length} items');
      if (res.isNotEmpty) {
        logger.i('First item: ${res[0]}');
      }
    }
    
    return res;
  } catch (err) {
    logger.e('=== API Error ===');
    logger.e('Error message: $err');
    logger.e('Error type: ${err.runtimeType}');
    // エラーが発生した場合は空のリストを返す
    // return ShiftCardDataList([]);
    return null; // エラー時はnullを返す
  }
}

class TabInfo {
  String label;
  int dayID;
  // int weatherID; // 天気IDを指定
  Widget widget;
  TabInfo(
    this.label,
    this.dayID, // 日付IDを指定
    // this.weatherID, // 天気IDを指定
    this.widget, 
    // {
    // this.dayID = 1, // デフォルトは準備日
    // this.weatherID = 1, // デフォルトは晴れ
    // }
  );
}


class MyShiftPage extends StatefulWidget {
  @override
  _MyShiftPageState createState() => _MyShiftPageState();
}

class _MyShiftPageState extends State<MyShiftPage>
    with TickerProviderStateMixin {
  int _selectedWeatherIndex = 1;  // 天気の選択肢のインデックス(1:晴れ, 2:雨)
  int _selectedDayIndex = 1; // 日付の選択肢のインデックス(1:準備日, 2:1日目, 3:2日目, 4:片付け日)
  ShiftCardDataList? shiftCardDataList; // ShiftCardDataListを格納する変数
  // Map<String, ShiftCardDataList>? allShiftCardDataList; // 全てのシフトカードデータを格納する変数(キーはhiveのキーと同じ)
  Timer? _debounceTimer;
  Timer? _tabSwitchDebounceTimer; // タブ切り替え用のデバウンスタイマー
  // int _currentTabIndex = 0; // 現在のタブインデックス
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
  // final int _userID = store.getUserID(); // ユーザIDを取得
  late Future<int> _userID; // ユーザIDを格納する変数
  @override
  void didChangeDependencies() {  // 依存関係が変わったときに呼ばれる
    super.didChangeDependencies();
    // ユーザIDを取得
    _userID = store.getUserID();
    logger.i('User ID: $_userID');
  }
  late TabController _tabController;
  
  // ウィジェットの初期化時の処理
  @override
  void initState() {
    _tabController = TabController(
      // length: _tabs.length,
      length: _dayOptions.length,
      vsync: this
    );
    _tabController.addListener(_handleTabChange); // タブの変更を監視するリスナーを追加
    super.initState();
    
    // 初期タブのデータを取得
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _userID.then((id) {
        logger.i('User ID in initState: $id');
        // 初期タブのデータをロード
        _loadShiftCardDataList(id, _selectedDayIndex, _selectedWeatherIndex);
      });
    });
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
    _tabController.index; // indexを参照しておく
    logger.i('=== Tab Change ===');
    logger.i('New tab index: ${_tabController.index + 1}'); // タブのインデックスを1から始まる日付IDに変換
    // 新しいタブのデータを取得
    setState(() {
      _selectedDayIndex = _tabController.index + 1; // タブのインデックスを1から始まる日付IDに変換
      logger.i('Selected day index: $_selectedDayIndex');
    });
    _userID.then((id) {
      logger.i('User ID in _handleTabChange: $id');
      // タブ切り替え時にデータをロード
      _loadShiftCardDataList(id, _selectedDayIndex, _selectedWeatherIndex);
    });
  }
  
  // SegmentedButton(天気)の選択状態が変わったときの処理
  void _handleWeatherSelectionChanged(Set<Object> newSelection) {
    final oldWeatherIndex = _selectedWeatherIndex;
    final newWeatherIndex = (newSelection.first as int);
    
    logger.i('=== Weather Selection Changed ===');
    logger.i('Old weather index: $oldWeatherIndex');
    logger.i('New weather index: $newWeatherIndex');
    
    // 新しい天気のデータを取得
    setState(() {
      _selectedWeatherIndex = newWeatherIndex; // 選択された天気のインデックスを更新
      logger.i('Selected weather index: $_selectedWeatherIndex');
    });
    _userID.then((id) {
      logger.i('User ID in _handleWeatherSelectionChanged: $id');
      // 新しい天気のデータをロード
      _loadShiftCardDataList(id, _selectedDayIndex, _selectedWeatherIndex);
    });
  }
  
  // キャッシュからデータをロードする関数
  // ShiftCardDataList? _getCashedShiftCardDataList(int dayID, int weatherID) {
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
    // final ShiftCardDataList shiftCardDataList = ShiftCardDataList.fromJson(cashedData); // キャッシュデータを状態に格納
    
    // return shiftCardDataList; // キャッシュデータを返す
    return cashedData; // キャッシュデータを返す
  }
  
  // シフトカードデータリストをロードする関数
  void _loadShiftCardDataList(int userID, int dayID, int weatherID) async {
    logger.i('=== シフトカードのデータを再設定します ===');
    logger.i('Parameters - userID: $userID, dayID: $dayID, weatherID: $weatherID');
    
    
    // 先にキャッシュデータを表示データを設定しておく(キャッシュデータがない場合はnull)
    // final ShiftCardDataList? cashedData = _getCashedShiftCardDataList(dayID, weatherID);
    final List<dynamic>? cashedData = _getCashedShiftCardDataList(dayID, weatherID); // キャッシュからデータを取得
    final ShiftCardDataList? cashedShiftCardDataList = 
      cashedData != null ? ShiftCardDataList.fromJson(cashedData) : null; // キャッシュデータをShiftCardDataListに変換
    // setState(() => shiftCardDataList = cashedData); // キャッシュデータを状態に格納
    setState(() => shiftCardDataList = cashedShiftCardDataList); // キャッシュデータを状態に格納
    logger.i('キャッシュデータを状態に格納しました: ${shiftCardDataList?.data.length ?? 0} items');
    
    // サーバーからデータをフェッチ
    _getShiftCardDataList(userID, dayID, weatherID).then((data) {
      // データフェッチ成功時の処理
      logger.i('=== Shift Card Data List Loaded ===');
      
      final ShiftCardDataList? fetchedData = ShiftCardDataList.fromJson(data!);
      
      // フェッチデータがnullの場合は何もしない
      if (fetchedData == null || deepEq(data, cashedData)) {
      // if (fetchedData == null) {
        fetchedData == null
          ? logger.w('フェッチデータが空です。キャッシュを使用します。')
          : logger.w('フェッチデータが既に最新です。キャッシュを使用します。');
        return;
      } 
      // if (fetchedData == cashedData) {
      //   logger.i('データは既に最新です。キャッシュ更新をスキップします。');
      //   return; // データが既に最新の場合も何もしない
      // } else {
      //   logger.i('新しいデータが取得されました。キャッシュ更新を行います。');
      // }
      logger.i('新しいデータが取得されました。');
      logger.i('fetchedData data : ${data}');
      logger.i('cachedData data : ${cashedData}');
      // print(deepEq(data, cashedData) ? 'データは同じです。' : 'データは異なります。');
      // キャッシュデータとフェッチデータの差分
      // final diff = fetchedData.data.where((item) => 
        // !cashedShiftCardDataList!.data.any((cachedItem) => cachedItem == item)
      // ).toList();
      // logger.i(diff[0]);
      logger.i('キャッシュ更新を行います。');
      shiftCardBox.put('${dayID}_${weatherID}', data); // hiveでキャッシュにデータを保存
      logger.i('dayID: $dayID, weatherID: $weatherID のデータをキャッシュに保存しました。key: ${dayID}_${weatherID}');
      
      logger.i('表示データを更新します。');
      setState(() => shiftCardDataList = fetchedData); // 状態にデータを格納
      logger.i('データを状態に格納しました: ${shiftCardDataList!.data.length} items');
      
      
      // setState(() {
      //   logger.i('=== Updating State with Fetched Data ===');
      //   logger.i('dayID: $dayID, weatherID: $weatherID のデータを状態に格納します。');
      //   // shiftCardDataList = fetchedData; // 状態にデータを格納
      //   shiftCardDataList = ShiftCardDataList.fromJson(fetchedData); // 状態にデータを格納
      //   logger.i('データを状態に格納しました: ${shiftCardDataList!.data.length} items');
        
      //   logger.i('dayID: $dayID, weatherID: $weatherID のデータをキャッシュに保存します。');
      //   shiftCardBox.put('${dayID}_${weatherID}', fetchedData); // hiveでキャッシュにデータを保存
      //   logger.i('データをキャッシュに保存しました: ${dayID}_${weatherID}');
        
      //   // ロード状態を更新
      //   // _loadingStates['${dayID}_${weatherID}'] = false;
      //   // データを状態に格納
      //   // shiftCardDataList = fetchData;
      // });
      // logger.i('=== Shift Card Data List Loaded and Cached ===');
    }).catchError((error) async {
      logger.e('=== Error Loading Shift Card Data List ===');
      // var casheData = await store.getShiftCardDataList('${dayID}_${weatherID}');
      final cashedData = await shiftCardBox.get('${dayID}_${weatherID}');  // キャッシュからデータを取得
      if (cashedData == null) {
        logger.e('$dayID, $weatherID のキャッシュデータがありません。');
        return; // キャッシュデータがない場合は何もしない
      }
      setState(() {
        logger.i('=== Using Cached Data ===');
        // shiftCardDataList = cashedData;
        shiftCardDataList = ShiftCardDataList.fromJson(cashedData); // キャッシュデータを状態に格納
        logger.i('キャッシュデータを状態に格納しました: ${shiftCardDataList!.data.length} items');
        
        // if (casheData != null) {
        //   logger.i('Using cached data for dayID: $dayID, weatherID: $weatherID');
        //   shiftCardDataList = ShiftCardDataList.fromJson(casheData);
        // } else {
        //   logger.e('No cached data available for dayID: $dayID, weatherID: $weatherID');
        // }

        // _loadingStates['${dayID}_${weatherID}'] = false; // ロード状態を更新
      });
    });
  }
  
  // // デバウンス付きデータ読み込み
  // void _loadShiftCardDataListWithDebounce() {
  //   _debounceTimer?.cancel();
  //   _debounceTimer = Timer(Duration(milliseconds: 500), () {
  //     _clearCurrentTabCacheAndReload();
  //   });
  // }
  
  // // 現在のタブのキャッシュのみをクリアして再読み込み
  // void _clearCurrentTabCacheAndReload() {
  //   logger.i('=== Clearing Current Tab Cache and Reloading ===');
    
  //   // 現在のタブに関連するキャッシュのみを削除
  //   final dayID = _tabs[_currentTabIndex].dayID;
  //   final keysToRemove = <String>[];
    
  //   // 現在のタブの全ての天気パターンのキャッシュを削除
  //   for (final weatherID in _weatherOptions.values) {
  //     final cacheKey = '${_currentTabIndex}_${dayID}_${weatherID}';
  //     keysToRemove.add(cacheKey);
  //   }
    
  //   logger.i('Removing cache keys: $keysToRemove');
    
  //   setState(() {
  //     for (final key in keysToRemove) {
  //       _dataCache.remove(key);
  //       _loadingStates.remove(key);
  //       _requestSequences.remove(key);
  //     }
  //   });
    
  //   logger.i('Current tab cache cleared successfully');
  //   // 現在のタブのデータを再取得
  //   _loadDataForCurrentTab();
  // }
  
  // // 全キャッシュをクリアして再読み込み（更新ボタン用）
  // void _clearCacheAndReload() {
  //   logger.i('=== Clearing All Cache and Reloading ===');
  //   logger.i('Cache before clear: ${_dataCache.keys.toList()}');
    
  //   setState(() {
  //     _dataCache.clear(); // キャッシュをクリア
  //     _loadingStates.clear(); // ロード状態もクリア
  //     _requestSequences.clear(); // シーケンスもクリア
  //   });
    
  //   logger.i('All cache cleared successfully');
  //   // 現在のタブのデータを再取得
  //   _loadDataForCurrentTab();
  // }
  
  // // 現在のタブのデータを取得
  // void _loadDataForCurrentTab() {
  //   final dayID = _tabs[_currentTabIndex].dayID;
  //   final weatherID = _selectedWeatherIndex;
  //   final cacheKey = '${_currentTabIndex}_${dayID}_${weatherID}';
    
  //   // シーケンス番号を増加
  //   _globalSequence++;
  //   final currentSequence = _globalSequence;
  //   _requestSequences[cacheKey] = currentSequence;
    
  //   logger.i('=== Loading Data for Current Tab ===');
  //   logger.i('tabIndex: $_currentTabIndex, dayID: $dayID, weatherID: $weatherID');
  //   logger.i('cacheKey: $cacheKey, sequence: $currentSequence');
    
  //   // 既にロード中またはキャッシュに存在する場合はスキップ
  //   if (_loadingStates[cacheKey] == true || _dataCache.containsKey(cacheKey)) {
  //     logger.i('Data already loading or cached for: $cacheKey');
  //     return;
  //   }
    
  //   // ロード状態を設定
  //   setState(() {
  //     _loadingStates[cacheKey] = true;
  //   });
    
  //   // データを取得
  //   _getShiftCardDataList(1, dayID, weatherID).then((data) {
  //     // 最新のリクエストかチェック
  //     if (_requestSequences[cacheKey] == currentSequence) {
  //       setState(() {
  //         _dataCache[cacheKey] = data;
  //         _loadingStates[cacheKey] = false;
  //       });
  //       logger.i('Data loaded and cached for: $cacheKey (sequence: $currentSequence)');
  //     } else {
  //       logger.i('Ignoring outdated response for: $cacheKey (sequence: $currentSequence, current: ${_requestSequences[cacheKey]})');
  //       setState(() {
  //         _loadingStates[cacheKey] = false;
  //       });
  //     }
  //   }).catchError((error) {
  //     // 最新のリクエストかチェック
  //     if (_requestSequences[cacheKey] == currentSequence) {
  //       setState(() {
  //         _loadingStates[cacheKey] = false;
  //       });
  //       logger.e('Error loading data for $cacheKey: $error (sequence: $currentSequence)');
  //     }
  //   });
  // }
  
  // // タブごとのデータを取得する関数（同期版）
  // ShiftCardDataList? _getDataForTab(int tabIndex) {
  //   final dayID = _tabs[tabIndex].dayID;
  //   final weatherID = _selectedWeatherIndex;
  //   final cacheKey = '${tabIndex}_${dayID}_${weatherID}';
    
  //   logger.i('=== _getDataForTab Called (Sync) ===');
  //   logger.i('tabIndex: $tabIndex, dayID: $dayID, weatherID: $weatherID');
  //   logger.i('cacheKey: $cacheKey');
  //   logger.i('Current tab index: $_currentTabIndex');
    
  //   // キャッシュから取得
  //   if (_dataCache.containsKey(cacheKey)) {
  //     logger.i('Cache HIT - returning cached data for: $cacheKey');
  //     return _dataCache[cacheKey];
  //   }
    
  //   // 現在のタブの場合のみデータを取得開始
  //   if (tabIndex == _currentTabIndex) {
  //     logger.i('Current tab - checking if should load data for: $cacheKey');
  //     // ロード中でない場合のみデータ取得を開始
  //     if (_loadingStates[cacheKey] != true) {
  //       logger.i('Not loading - starting data load for: $cacheKey');
  //       _loadDataForCurrentTab();
  //     } else {
  //       logger.i('Already loading - skipping data load for: $cacheKey');
  //     }
  //   } else {
  //     logger.i('Not current tab - skipping data load for: $cacheKey');
  //   }
    
  //   return null; // データがない場合はnullを返す
  // }
  
  // // 手動更新ボタンの処理
  // void _handleRefreshPressed() {
  //   logger.i('=== Manual Refresh Triggered ===');
  //   _clearCacheAndReload(); // 全キャッシュをクリアして再読み込み
  // }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.base,
      appBar: AppBar(
        // title: Padding(
        //   padding: const EdgeInsets.only(top: 20.0),
        //   child: Text("マイシフト", 
        //     style: TextStyle(
        //       color: AppColors.textWhite,
        //       fontSize: AppFontSizes.lg,
        //       fontWeight: FontWeight.bold,
        //     ),
        //   ),
        // ),
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
            // padding: const EdgeInsets.only(top: 20.0, right: 20.0),
            padding: const EdgeInsets.only(right: 20.0),
            child: SizedBox(
              width: 200,
              // height: 30,
              // 天気を選択するセグメントボタン
              child: SegmentedButton(
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
              ),
            ),
          ),
        ],
        bottom: PreferredSize(
          preferredSize: const Size.fromHeight(75.0),
          // 日付のタブバー
          child: TabBar(
            isScrollable: false,
            // tabs: _tabs.map((TabInfo tab) {
            //   return Tab(text: tab.label, height: 35);
            // }).toList(),
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
          ),
        ),
      ),
      body: Column(
        children: [
          // // 更新ボタン
          // Container(
          //   width: double.infinity,
          //   padding: const EdgeInsets.symmetric(horizontal: 20.0, vertical: 10.0),
          //   child: ElevatedButton.icon(
          //     onPressed: _handleRefreshPressed,
          //     icon: Icon(
          //       Icons.refresh,
          //       color: AppColors.main,
          //       size: 20,
          //     ),
          //     label: Text(
          //       "更新",
          //       style: TextStyle(
          //         color: AppColors.main,
          //         fontSize: AppFontSizes.sm,
          //         fontWeight: FontWeight.bold,
          //       ),
          //     ),
          //     style: ElevatedButton.styleFrom(
          //       backgroundColor: AppColors.base,
          //       elevation: 0,
          //       side: BorderSide(
          //         color: AppColors.main,
          //         width: 1.0,
          //       ),
          //       padding: const EdgeInsets.symmetric(vertical: 8.0),
          //       shape: RoundedRectangleBorder(
          //         borderRadius: BorderRadius.circular(8.0),
          //       ),
          //     ),
          //   ),
          // ),
          // TabBarView
          Expanded(
            child: TabBarView(
              controller: _tabController,
              // children: List.generate(_tabs.length, (index) {
              children: List.generate(_dayOptions.length, (index) {
                // final shiftCardDataList = _getDataForTab(index);
                
                // final isLoading = _loadingStates['${index}_${_tabs[index].dayID}_$_selectedWeatherIndex'] == true;
                
                logger.i('=== Building tab $index ===');
                logger.i('Has data: ${shiftCardDataList != null}');
                // logger.i('Is loading: $isLoading');
                
                if (shiftCardDataList == null) {
                  logger.i('Tab $index: No data available, showing loading indicator');
                  // データがない場合はロード中表示
                  return Center(
                    child: CircularProgressIndicator(
                      color: AppColors.main,
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
                      padding: const EdgeInsets.all(32.0),
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
                // if (isLoading) {
                //   logger.i('Tab $index: Showing loading indicator');
                //   return Center(
                //     child: CircularProgressIndicator(
                //       color: AppColors.main,
                //     ),
                //   );
                // } else if (shiftCardDataList != null) {
                //   logger.i('Tab $index: Data available with ${shiftCardDataList.data.length} items');
                  
                //   if (shiftCardDataList.data.isNotEmpty) {
                //     logger.i('Tab $index: First item task name: ${shiftCardDataList.data[0].taskName}');
                //   }
                  
                //   List<Widget> _shiftCards = [];
                //   for (var data in shiftCardDataList.data) {
                //     _shiftCards.add(ShiftCard(data: data));
                //     _shiftCards.add(const SizedBox(height: 16.0));
                //   }
                  
                //   return SingleChildScrollView(
                //     child: Container(
                //       padding: const EdgeInsets.all(32.0),
                //       child: Column(
                //         crossAxisAlignment: CrossAxisAlignment.start,
                //         children: _shiftCards,
                //       ),
                //     ),
                //   );
                // } else {
                //   logger.i('Tab $index: No data available');
                //   if (index == _currentTabIndex) {
                //     // 現在のタブでデータがない場合はロード中表示
                //     return Center(
                //       child: CircularProgressIndicator(
                //         color: AppColors.main,
                //       ),
                //     );
                //   } else {
                //     // 他のタブは空の状態
                //     return Center(child: Text('タブを選択してください'));
                //   }
                // }
              // }),
              })
            ),
          ),
        ],
      ),
    );
  } 
}