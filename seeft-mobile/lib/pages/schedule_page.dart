import 'dart:developer';

import 'package:seeft_mobile/configs/importer.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:http/http.dart' as http;
import 'package:seeft_mobile/pages/schedule_page_first_day_rainy.dart';
import 'package:seeft_mobile/pages/schedule_page_first_day_sunny.dart';
import 'package:seeft_mobile/pages/schedule_page_second_day_rainy.dart';
import 'package:seeft_mobile/pages/schedule_page_second_day_sunny.dart';

class SchedulePage extends StatefulWidget {
  @override
  _SchedulePageState createState() => _SchedulePageState();
}

class TabInfo {
  String label;
  Widget widget;
  TabInfo(this.label, this.widget);
}

class _SchedulePageState extends State<SchedulePage>   with SingleTickerProviderStateMixin {
  final List<TabInfo> _tabs = [
    TabInfo("9/10 晴天時", SchedulePageFirstDaySunny()),
    TabInfo("9/10 雨天時", SchedulePageFirstDayRainy()),
    TabInfo("9/11 晴天時", SchedulePageSecondDaySunny()),
    TabInfo("9/11 雨天時", SchedulePageSecondDayRainy()),
  ];
  late TabController _tabController;

  // notification関連をinitStateに書き出さなきゃいけないので書いてたけどutilとかに書いてもいいかもね

//  FlutterLocalNotificationsPlugin flutterLocalNotificationsPlugin;
//  NotificationDetails platformChannelSpecifics;

  @override
  void initState() {
    _tabController = TabController(length: _tabs.length, vsync: this);
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    final Size size = MediaQuery.of(context).size;
    return Scaffold(
      appBar: AppBar(
        title: const Text('タイムスケジュール'),
        actions: <Widget>[],
        bottom: PreferredSize(
          child: TabBar(
            isScrollable: true,
            tabs: _tabs.map((TabInfo tab) {
              return Tab(text: tab.label);
            }).toList(),
            controller: _tabController,
          ),
          preferredSize: Size.fromHeight(30.0),
        ),
        // debug
      ),
      drawer: drawer.applicationDrawer(context),
      body: TabBarView(
          controller: _tabController,
          children: _tabs.map((tab) => tab.widget).toList()),
    );
  }
}