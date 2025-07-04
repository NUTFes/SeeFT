import 'package:seeft_mobile/configs/importer.dart';
import 'package:seeft_mobile/pages/wait_page.dart';
import 'package:seeft_mobile/widgets/shift_card.dart';


Future<ShiftCardDataList> _getShiftCardDataList(int userID, int dayID, int weatherID) async {
  logger.i('=== API Call Started ===');
  logger.i('Parameters - userID: $userID, dayID: $dayID, weatherID: $weatherID');
  
  try {
    var res = await api.getShiftCardsByUserAndDateAndWeather(
      userID,
      dayID,
      weatherID,
    );
    
    logger.i('=== API Response Received ===');
    logger.i('Raw API Response: $res');
    logger.i('Response Type: ${res.runtimeType}');
    
    if (res is List) {
      logger.i('Response is List with ${res.length} items');
      if (res.isNotEmpty) {
        logger.i('First item: ${res[0]}');
      }
    }
    
    // resをShiftCardDataListに変換
    ShiftCardDataList resList = ShiftCardDataList.fromJson(res);
    logger.i('=== Converted to ShiftCardDataList ===');
    logger.i('ShiftCardDataList data count: ${resList.data.length}');
    if (resList.data.isNotEmpty) {
      logger.i('First ShiftCardData: ${resList.data[0].taskName}');
    }
    
    return resList;
  } catch (err) {
    logger.e('=== API Error ===');
    logger.e('Error message: $err');
    logger.e('Error type: ${err.runtimeType}');
    // エラーが発生した場合は空のリストを返す
    return ShiftCardDataList([]);
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
  
  // データフェッチ管理用の変数
  Map<String, ShiftCardDataList> _dataCache = {}; // Futureではなく実際のデータをキャッシュ
  Map<String, bool> _loadingStates = {}; // ロード状態を管理
  Map<String, int> _requestSequences = {}; // リクエストのシーケンス番号を管理
  int _globalSequence = 0; // グローバルシーケンス番号
  Timer? _debounceTimer;
  Timer? _tabSwitchDebounceTimer; // タブ切り替え用のデバウンスタイマー
  // int _currentTabIndex = 0; // 現在のタブインデックス
  
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
  // late List<TabInfo> _tabs = [
  //   // TabInfo(" 準備日 ", _MyShiftPageTabPage(day: _dayOptions["準備日"]?? 1, weather: _weatherOptions["晴れ"]?? 1)), // 準備日は天気を選択しないので1を固定
  //   // TabInfo(" １日目 ", _MyShiftPageTabPage(day: _dayOptions["1日目"]?? 2, weather: _selectedWeatherIndex)),
  //   // TabInfo(" ２日目 ", _MyShiftPageTabPage(day: _dayOptions["2日目"]?? 3, weather: _selectedWeatherIndex)),
  //   TabInfo(
  //     " 準備日 ", 
  //     _dayOptions["準備日"] ?? 1, // 準備日のdateIDを取得
  //     // _weatherOptions["晴れ"] ?? 1, // 準備日は天気を選択しないので1を固定
  //     WaitPage()
  //   ),
  //   TabInfo(
  //     " １日目 ", 
  //     _dayOptions["1日目"] ?? 2, // １日目のdateIDを取得
  //     // _selectedWeatherIndex, // 選択された天気のインデックスを使用
  //     WaitPage()
  //   ),
  //   TabInfo(
  //     " ２日目 ", 
  //     _dayOptions["2日目"] ?? 3, // ２日目のdateIDを取得
  //     // _selectedWeatherIndex, // 選択された天気のインデックスを使用
  //     WaitPage()
  //   ),
  //   TabInfo(
  //     "片付け日", 
  //     _dayOptions["片付け日"] ?? 4, // 片付け日のdateIDを取得
  //     // _weatherOptions["晴れ"] ?? 1, // 片付け日は天気を選択しないので1を固定
  //     WaitPage()
  //   ),
  // ];
  late TabController _tabController;

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
      // _loadDataForCurrentTab();
      _userID.then((id) {
        logger.i('User ID in initState: $id');
        // 初期タブのデータをロード
        _loadShiftCardDataList(id, _selectedDayIndex, _selectedWeatherIndex);
      });
      // _loadShiftCardDataList(1, _selectedDayIndex, _selectedWeatherIndex); // ユーザIDは1、日付IDは1(準備日)、天気IDは選択された天気のインデックスを使用
    });
  }
  @override
  void dispose() {
    _tabController.removeListener(_handleTabChange);
    _tabController.dispose();
    _debounceTimer?.cancel(); // デバウンスタイマーをキャンセル
    _tabSwitchDebounceTimer?.cancel(); // タブ切り替えデバウンスタイマーをキャンセル
    super.dispose();
  }
  
  // タブが切り替わったときの処理
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
    
    // // タブ切り替えのデバウンス処理
    // _tabSwitchDebounceTimer?.cancel();
    // _tabSwitchDebounceTimer = Timer(Duration(milliseconds: 200), () {
    //   setState(() {
    //     // _currentTabIndex = _tabController.index;
    //     _selectedDayIndex = _tabController.index + 1; // タブのインデックスを1から始まる日付IDに変換
    //   });
    //   // 新しいタブのデータを取得
    //   _loadDataForCurrentTab();
    // });
  }
  
  // SegmentedButtonの選択状態が変わったときの処理
  void _handleWeatherSelectionChanged(Set<Object> newSelection) {
    final oldWeatherIndex = _selectedWeatherIndex;
    // final newWeatherIndex = newSelection.first + 1; // セグメントボタンのインデックスは0から始まるので1を足す
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
    
    // setState(() {
    //   _selectedWeatherIndex = newWeatherIndex;
    // });
    
    // if (oldWeatherIndex != newWeatherIndex) {
    //   logger.i('Weather actually changed - triggering debounced reload');
    //   _loadShiftCardDataListWithDebounce(); // デバウンス付きでデータ読み込み
    // } else {
    //   logger.i('Weather not changed - skipping reload');
    // }
  }
  
  // シフトカードデータリストをロードする関数
  void _loadShiftCardDataList(int userID, int dayID, int weatherID) {
    logger.i('=== Loading Shift Card Data List ===');
    logger.i('Parameters - userID: $userID, dayID: $dayID, weatherID: $weatherID');
    
    _getShiftCardDataList(userID, dayID, weatherID).then((data) {
      logger.i('=== Shift Card Data List Loaded ===');
      logger.i('Data count: ${data.data.length}');
      
      setState(() {
        // キャッシュにデータを保存
        _dataCache['${dayID}_${weatherID}'] = data;
        // ロード状態を更新
        _loadingStates['${dayID}_${weatherID}'] = false;
        // データを状態に格納
        shiftCardDataList = data;
        logger.i('Data cached successfully for dayID: $dayID, weatherID: $weatherID');
      });
      
      logger.i('Data cached successfully for dayID: $dayID, weatherID: $weatherID');
    }).catchError((error) {
      logger.e('Error loading shift card data list: $error');
      setState(() {
        _loadingStates['${dayID}_${weatherID}'] = false; // ロード状態を更新
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