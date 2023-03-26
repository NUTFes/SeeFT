import 'dart:developer';

import 'package:seeft_mobile/configs/importer.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:http/http.dart' as http;
import 'package:seeft_mobile/pages/my_shift_page_preparation_day_sunny.dart';
import 'package:seeft_mobile/pages/my_shift_page_preparation_day_rainy.dart';

class MyShiftPagePreparationDay extends StatefulWidget {
  @override
  _MyShiftPageState createState() => _MyShiftPageState();
}

class TabInfo {
  String label;
  Widget widget;
  TabInfo(this.label, this.widget);
}

class _MyShiftPageState extends State<MyShiftPagePreparationDay>
  with SingleTickerProviderStateMixin {
  final List<TabInfo> _tabs = [
  TabInfo("晴れ", MyShiftPagePreparationDaySunny()),
//  TabInfo("雨", MyShiftPagePreparationDayRainy()),
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
        actions: <Widget>[],
        bottom: PreferredSize(
          child: TabBar(
            isScrollable: true,
            tabs: _tabs.map((TabInfo tab) {
              return Tab(text: tab.label);
            }).toList(),
            controller: _tabController,
          ),
          preferredSize: Size.fromHeight(10.0),
        ),
        // debug
      ),
      body: TabBarView(
          controller: _tabController,
          children: _tabs.map((tab) => tab.widget).toList()),
    );
  }
}