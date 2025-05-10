import 'package:seeft_mobile/configs/importer.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:http/http.dart' as http;
import 'package:seeft_mobile/pages/my_shift_page_preparation_day.dart';
import 'package:seeft_mobile/pages/my_shift_page_current_first_day.dart';
import 'package:seeft_mobile/pages/my_shift_page_current_second_day.dart';
import 'package:seeft_mobile/pages/pre_preparation_day_shift.dart';
import 'package:seeft_mobile/pages/cleanup_day_time_schedule.dart';
import 'package:seeft_mobile/pages/wait_page.dart';
import 'package:seeft_mobile/theme/tokens.dart';
/*
import 'package:seeft_mobile/pages/my_shift_page_preparation_day_sunny.dart';
import 'package:seeft_mobile/pages/my_shift_page_preparation_day_rainy.dart';
import 'package:seeft_mobile/pages/my_shift_page_current_day_sunny.dart';
import 'package:seeft_mobile/pages/my_shift_page_current_day_rainy.dart';
import 'package:seeft_mobile/pages/my_shift_page_cleanup_day.dart';
*/

class MyShiftPage extends StatefulWidget {
  @override
  _MyShiftPageState createState() => _MyShiftPageState();
}

class TabInfo {
  String label;
  Widget widget;
  TabInfo(this.label, this.widget);
}

class _MyShiftPageState extends State<MyShiftPage>
    with SingleTickerProviderStateMixin {
  final List<TabInfo> _tabs = [
    TabInfo(" 準備日 ", MyShiftPagePreparationDay()),
    TabInfo(" １日目 ", MyShiftPageCurrentFirstDay()),
    TabInfo(" ２日目 ", MyShiftPageCurrentSecondDay()),
    TabInfo("片付け日", WaitPage()),
  ];
  late TabController _tabController;
  int _selectedWeatherIndex = 0;  // 天気の選択肢のインデックス(0:晴れ, 1:雨)

  @override
  void initState() {
    _tabController = TabController(length: _tabs.length, vsync: this);
    super.initState();
  }
  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
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
        actions: [
          Padding(
            padding: const EdgeInsets.only(top: 20.0, right: 20.0),
            child: SizedBox(
              width: 131,
              height: 30,
              // 天気を選択するセグメントボタン
              child: SegmentedButton(
                selected: {_selectedWeatherIndex},
                onSelectionChanged: (Set<int> newSelection) {
                  setState(() {
                    _selectedWeatherIndex = newSelection.first;
                  });
                },
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
                      color: _selectedWeatherIndex == 0 ? AppColors.main : AppColors.grayLight,
                      size: 18,
                    ),
                    value: 0,
                  ),
                  ButtonSegment(
                    icon: Icon(
                      Icons.cloudy_snowing,
                      color: _selectedWeatherIndex == 1 ? AppColors.main : AppColors.grayLight,
                      size: 18,
                    ),
                    value: 1,
                  ),
                ],
              ),
            ),
          ),
        ],
        backgroundColor: AppColors.main,
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
      body: TabBarView(
          controller: _tabController,
          children: _tabs.map((tab) => tab.widget).toList()),
    );
  }
}
