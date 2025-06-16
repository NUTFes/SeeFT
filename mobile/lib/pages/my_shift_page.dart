import 'package:seeft_mobile/configs/importer.dart';
// import 'package:flutter_local_notifications/flutter_local_notifications.dart';
// import 'package:http/http.dart' as http;
// import 'package:seeft_mobile/pages/my_shift_page_preparation_day.dart';
// import 'package:seeft_mobile/pages/my_shift_page_current_first_day.dart';
// import 'package:seeft_mobile/pages/my_shift_page_current_second_day.dart';
// import 'package:seeft_mobile/pages/pre_preparation_day_shift.dart';
// import 'package:seeft_mobile/pages/cleanup_day_time_schedule.dart';
import 'package:seeft_mobile/pages/wait_page.dart';
import 'package:seeft_mobile/widgets/shift_card.dart';


Future<ShiftCardDataList> _getShiftCardDataList(int userID, int dayID, int weatherID) async {
  try {
    // var userID = await store.getUserID();
    var res = await api.getShiftCardsByUserAndDateAndWeather(
      userID,
      dayID,
      weatherID,
      // DateTime.now().year.toString(),
    );
    logger.i('res: $res');
    // var taskName = res[0]['task_name'] as String? ?? 'No Task Name';
    // logger.i('taskName: $taskName');
    // resをShiftCardDataListに変換
    ShiftCardDataList resList = ShiftCardDataList.fromJson(res);
    logger.i('resList: $resList');
    // return res;
    return resList;
  } catch (err) {
    logger.e('don`t response. error message: $err');
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
  // 天気ごとのweatherID
  final Map<String, int> _weatherOptions = {
    "晴れ": 1,
    "雨": 2,
  };
  // 日付ごとのdateID
  final Map<String, int> _dayOptions = {
    "準備日": 1,
    "1日目": 2,
    "2日目": 3,
    "片付け日": 4,
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
  late List<TabInfo> _tabs = [
    // TabInfo(" 準備日 ", _MyShiftPageTabPage(day: _dayOptions["準備日"]?? 1, weather: _weatherOptions["晴れ"]?? 1)), // 準備日は天気を選択しないので1を固定
    // TabInfo(" １日目 ", _MyShiftPageTabPage(day: _dayOptions["1日目"]?? 2, weather: _selectedWeatherIndex)),
    // TabInfo(" ２日目 ", _MyShiftPageTabPage(day: _dayOptions["2日目"]?? 3, weather: _selectedWeatherIndex)),
    TabInfo(
      " 準備日 ", 
      _dayOptions["準備日"] ?? 1, // 準備日のdateIDを取得
      // _weatherOptions["晴れ"] ?? 1, // 準備日は天気を選択しないので1を固定
      WaitPage()
    ),
    TabInfo(
      " １日目 ", 
      _dayOptions["1日目"] ?? 2, // １日目のdateIDを取得
      // _selectedWeatherIndex, // 選択された天気のインデックスを使用
      WaitPage()
    ),
    TabInfo(
      " ２日目 ", 
      _dayOptions["2日目"] ?? 2, // １日目のdateIDを取得
      // _selectedWeatherIndex, // 選択された天気のインデックスを使用
      WaitPage()
    ),
    TabInfo(
      "片付け日", 
      _dayOptions["片付け日"] ?? 4, // 片付け日のdateIDを取得
      // _weatherOptions["晴れ"] ?? 1, // 片付け日は天気を選択しないので1を固定
      WaitPage()
    ),
  ];
  late TabController _tabController;
  Future<ShiftCardDataList>? _shiftCardDataListFuture; // シフトカードデータのFuture

  @override
  void initState() {
    _tabController = TabController(length: _tabs.length, vsync: this);
    _tabController.addListener(_handleTabChange);
    _loadShiftCardDataList(); // 初期データの読み込み
    super.initState();
  }
  @override
  void dispose() {
    _tabController.removeListener(_handleTabChange);
    _tabController.dispose();
    super.dispose();
  }
  
  // タブが切り替わったときの処理
  void _handleTabChange() {
    if (_tabController.indexIsChanging) {
      // タブが実際に変更された場合のみデータを再読み込み
      _loadShiftCardDataList();
    }
  }
  
  // SegmentedButtonの選択状態が変わったときの処理
  void _handleWeatherSelectionChanged(Set<int> newSelection) {
    setState(() {
      _selectedWeatherIndex = newSelection.first;
      _loadShiftCardDataList();
      // // タブの内容を更新
      // _tabs[1] = TabInfo(
      //   " １日目 ",
      //   _dayOptions["1日目"] ?? 2, // １日目のdateIDを取得
      //   WaitPage(), // 後で実装するページに置き換える
      // );
      // _tabs[2] = TabInfo(
      //   " ２日目 ",
      //   _dayOptions["2日目"] ?? 3, // ２日目のdateIDを取得
      //   WaitPage(), // 後で実装するページに置き換える
      // );
      // _tabController = TabController(
      //   initialIndex: _tabController.index,
      //   length: _tabs.length,
      //   vsync: this,
      // );
      logger.i('Selected weather index: $_selectedWeatherIndex');
    });
  }
  
  // データ読み込みのトリガーとなる関数
  void _loadShiftCardDataList() {
    setState(() {
      _shiftCardDataListFuture = _getShiftCardDataList(
        // _userID,
        1,
        _tabs[_tabController.index].dayID, // タブのインデックスから日付IDを取得
        
        // _dayOptions[_tabs[_tabController.index].label.trim()] ?? 1, // タブのラベルから日付IDを取得
        _selectedWeatherIndex, // 選択された天気のインデックスを使用
      );
      // _weatherDataFuture = fetchWeatherData(
      //   _selectedWeather.first, // SegmentedButtonはSetで選択状態を持つためfirstで取得
      //   _dates[_tabController.index],
      // );
      
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.base,
      appBar: AppBar(
        title: Padding(
          padding: const EdgeInsets.only(top: 20.0),
          child: Text("マイシフト", 
            style: TextStyle(
              color: AppColors.textWhite,
              fontSize: AppFontSizes.lg,
              fontWeight: FontWeight.bold,
            ),
          ),
        ),
        centerTitle: false,
        toolbarHeight: 50,
        backgroundColor: AppColors.main,
        actions: [
          Padding(
            padding: const EdgeInsets.only(top: 20.0, right: 20.0),
            child: SizedBox(
              width: 131,
              height: 30,
              // 天気を選択するセグメントボタン
              child: SegmentedButton(
                selected: {_selectedWeatherIndex},
                onSelectionChanged: _handleWeatherSelectionChanged,
                // onSelectionChanged: (Set<int> newSelection) {
                //   setState(() {
                //     _selectedWeatherIndex = newSelection.first;
                //     // // 選択された天気のインデックスを更新
                //     // _tabs[1] = TabInfo(
                //     //   " １日目 ",
                //     //   _MyShiftPageTabPage(
                //     //     day: _dayOptions["1日目"] ?? 2,
                //     //     weather: _selectedWeatherIndex,
                //     //   ),
                //     // );
                //     // _tabs[2] = TabInfo(
                //     //   " ２日目 ",
                //     //   _MyShiftPageTabPage(
                //     //     day: _dayOptions["2日目"] ?? 3,
                //     //     weather: _selectedWeatherIndex,
                //     //   ),
                //     // );
                //     _tabController = TabController(
                //       initialIndex: _tabController.index,
                //       length: _tabs.length,
                //       vsync: this
                //     );
                //     // _tabController.index = _tabController.index; // 初期タブを準備日に設定
                //     // _tabController.notifyListeners(); // タブの更新を通知
                //     logger.i('Selected weather index: $_selectedWeatherIndex');
                //   });
                // },
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
                    icon: Icon(
                      Icons.sunny,
                      color: _selectedWeatherIndex == (_weatherOptions["晴れ"] ?? 1) ? AppColors.main : AppColors.grayLight,
                      size: 18,
                    ),
                    value: _weatherOptions["晴れ"] ?? 1, // 晴れのインデックス
                  ),
                  ButtonSegment(
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
            tabs: _tabs.map((TabInfo tab) {
              return Tab(text: tab.label, height: 35);
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
      body: FutureBuilder<ShiftCardDataList>(
        future: _shiftCardDataListFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return Center(
              child: CircularProgressIndicator(
                color: AppColors.main, // 待機中のインジケーターの色
              ), // データ取得中は待機画面を表示
            );
          } else if (snapshot.hasError) {
            return Center(child: Text('エラーが発生しました: ${snapshot.error}'));
          } else if (snapshot.hasData) {
            final shiftCardDataList = snapshot.data!;
            // Widgetのリストを作成
            List<Widget> _shiftCards = [];
            for (var data in shiftCardDataList.data) {
              _shiftCards.add(ShiftCard(data: data));
              _shiftCards.add(const SizedBox(height: 16.0));
            }
            return SingleChildScrollView(
              child: Container(
                padding: const EdgeInsets.all(32.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: _shiftCards,
                  // children: _tabs.map((tab) => tab.widget).toList()),
                  // children: shiftCardDataList.data.map((data) => [
                  //   ShiftCard(data: data),
                  //   const SizedBox(height: 16.0), // 各カードの間にスペースを追加
                  // ]).expand((widgets) => widgets).toList(),
                ),
              ),
            );
          } else {
            return WaitPage(); // データがない場合は待機画面を表示
          }
        },
      ),
    );
  }
  
  // // FutureBuilderを使ってデータを取得し、ShiftCardDataListを表示する
  // FutureBuilder<ShiftCardDataList> _buildShiftCardDataList(int userID, int dayID, int weatherID) {
  //   return FutureBuilder<ShiftCardDataList>(
  //     future: _getShiftCardDataList(userID, dayID, weatherID),
  //     builder: (context, snapshot) {
  //       if (snapshot.connectionState == ConnectionState.waiting) {
  //         return WaitPage(); // データ取得中は待機画面を表示
  //       } else if (snapshot.hasError) {
  //         return Center(child: Text('エラーが発生しました: ${snapshot.error}'));
  //       } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
  //         return Center(child: Text('シフトデータがありません'));
  //       } else {
  //         final shiftCardDataList = snapshot.data!;
  //         return SingleChildScrollView(
  //           child: Container(
  //             padding: const EdgeInsets.all(32.0),
  //             child: Column(
  //               children: shiftCardDataList.map((data) => ShiftCard(data: data)).toList(),
  //             ),
  //           ),
  //         );
  //       }
  //     },
  //   );
  // }
}

// class _MyShiftPageTabPage extends StatefulWidget {
//   final int day; // 1:準備日, 2:１日目, 3:２日目, 4:片付け日
//   final int weather; // 0:晴れ, 1:雨
//   const _MyShiftPageTabPage({
//     Key? key,
//     required this.day,
//     required this.weather,
//   }) : super(key: key);
//   @override
//   _MyShiftPageTabPageState createState() => _MyShiftPageTabPageState();
// }

// class _MyShiftPageTabPageState extends State<_MyShiftPageTabPage> {
// // notification関連をinitStateに書き出さなきゃいけないので書いてたけどutilとかに書いてもいいかもね
  
//   late final int day = widget.day; // 1:準備日, 2:１日目, 3:２日目, 4:片付け日
//   late final int weather = widget.weather; // 0:晴れ, 1:雨
//   @override
//   void initState() {
//     super.initState();
//   }
  
//   @override
//   Widget build(BuildContext context) {
//     return FutureBuilder(
//       future: _getData(day, weather),
//       builder: (ctx, snapshot) {
//         if (snapshot.connectionState == AsyncSnapshot.waiting()) {
//           logger.w("message");
//         }
//         if (!snapshot.hasData) {
//           // 待機画面を表示
//           return WaitPage();
//         }
//         return SingleChildScrollView(
//           child: Container(
//             padding: const EdgeInsets.all(32.0),
//             child: Column(
//               children: <Widget>[
//                 for (var shift in snapshot.data)
//                   ShiftCard(data: _parseShiftMembers(shift)),
//               ],
//             ),
//           ),
//         );
//       },
//     );
//   }
  
  
//   }


//   // resからShiftCardDataのリストに変換する処理を作成
//   // ShiftCardDataのリストを返すようにする
//   ShiftCardData _parseShiftMembers(json) {
//     logger.i('parseShiftMembers: $json');
//     var taskName = json['task_name'] as String? ?? 'No Task Name';
//     logger.i('taskName: $taskName');
//     return ShiftCardData(
//       taskName: json['task_name'] as String,
//       startTime: json['start_time'] as String,
//       endTime: json['end_time'] as String,
//       place: json['place'] as String,
//       url: json['url'] as String,
//       shiftMembers: (json['shift_members'] as List<dynamic>)
//           .map((member) => ShiftMembers(
//                 s_time: member['s_time'],
//                 e_time: member['e_time'],
//                 members: (member['members'] as List<dynamic>)
//                     .map((m) => ShiftMember(
//                           name: m['name'],
//                           grade: m['grade'],
//                           bureau: m['bureau'],
//                         ))
//                     .toList(),
//               ))
//           .toList(),
//       beforeMembers: ShiftMembers(
//         s_time: json['before_members']['s_time'],
//         e_time: json['before_members']['e_time'],
//         members: json['before_members']['members'] != null?
//           (json['before_members']['members'] as List<dynamic>)
//             .map((m) => ShiftMember(
//                   name: m['name']?? '',
//                   grade: m['grade']?? '',
//                   bureau: m['bureau']?? '',
//                 ))
//             .toList():
//           [],
//       ),
//       afterMembers: ShiftMembers(
//         s_time: json['after_members']['s_time'],
//         e_time: json['after_members']['e_time'],
        // members: json['after_members']['members'] != null?
        //   (json['after_members']['members'] as List<dynamic>)
        //     .map((m) => ShiftMember(
        //           name: m['name']?? '',
        //           grade: m['grade']?? '',
        //           bureau: m['bureau']?? '',
        //         ))
        //     .toList():
        //   [],
//       ),
//     );
//   }
// }


