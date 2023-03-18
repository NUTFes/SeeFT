import 'package:seeft_mobile/configs/importer.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:http/http.dart' as http;
import 'package:seeft_mobile/pages/my_shift_page_preparation_day.dart';
import 'package:seeft_mobile/pages/my_shift_page_current_first_day.dart';
import 'package:seeft_mobile/pages/my_shift_page_current_second_day.dart';
import 'package:seeft_mobile/pages/pre_preparation_day_shift.dart';
import 'package:seeft_mobile/pages/cleanup_day_time_schedule.dart';
import 'package:seeft_mobile/pages/wait_page.dart';
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
    TabInfo("準備日", MyShiftPagePreparationDay()),
    TabInfo("当日 1日目", MyShiftPageCurrentFirstDay()),
    TabInfo("当日 2日目", MyShiftPageCurrentSecondDay()),
    TabInfo("片付け日", CleanUpDayTimeSchedule()),
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
        title: const Text('マイシフト'),
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

/*
  @override
  Widget build(BuildContext context) {
    final Size size = MediaQuery.of(context).size;
    return Scaffold(
      appBar: AppBar(
        title: const Text('マイシフト'),
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
      body: FutureBuilder(
        future: getData(),
        builder: (ctx, snapshot) {
          if (snapshot.connectionState == AsyncSnapshot.waiting()) {
            logger.w("message");
          }
          if (!snapshot.hasData) {
            return CircularProgressIndicator();
          }
          return Container(
              padding: const EdgeInsets.all(40.0),
              child: SingleChildScrollView(
                child: Column(
                  children: <Widget>[
                    Container(
                        // height: size.height - 200,
                        // width: size.width - 80,
                        // padding: const EdgeInsets.all(10.0),
                        // decoration: BoxDecoration(
                        //   border: Border.all(color: Colors.black),
                        // ),
                        // child: _contents(size, snapshot.data)),
                        child: _table(snapshot.data, context)),
                  ],
                ),
              ));
        },
      ),
    );
  }
}

Widget _table(var shifts, context) {
  return Table(
      border: TableBorder.all(color: Colors.white),
      columnWidths: const <int, TableColumnWidth>{
        // 0: IntrinsicColumnWidth(),
        0: FlexColumnWidth(1),
        1: FlexColumnWidth(10),
        // 2: FixedColumnWidth(100.0),
      },
      defaultVerticalAlignment: TableCellVerticalAlignment.middle,
      children: [
        TableRow(decoration: BoxDecoration(color: Colors.green[50]), children: [
          TableCell(
              child: Container(
            child: Text("日時"),
            alignment: Alignment.center,
          )),
          TableCell(
            child: Container(
              child: Text("場所"),
              alignment: Alignment.center,
            ),
          )
        ]),
        for (var shift in shifts)
          TableRow(
              decoration: BoxDecoration(color: Colors.grey[200]),
              children: [
                TableCell(
                    child: Container(
                  alignment: Alignment.center,
                  child: new Text(shift["Time"].toString()),
                )),
                TableCell(
                    child: Container(
                  width: double.infinity,
                  height: 54.0,
                  decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(10),
                  ),
                  child: TextButton(
                      child: new Text(shift["Work"].toString()),
                      style: ElevatedButton.styleFrom(
                        onPrimary: Colors.teal,
                      ),
                      onPressed: () async {
                        openShiftDialog(context);
                      }),
                  /*
                  alignment: Alignment.center,
                  child: new Text(shift["Work"].toString()),
                  // margin: EdgeInsets.only(bottom: 10.0),
                  height: 25,
                  */
                ))
              ]),
      ]);
}

Future getData() async {
  try {
    /*
    var userID = await store.getUserID();
    var res = await api.getMyShift(userID.toString());
    */

    var testUserID = 108;
    var res = await api.getMyShift(testUserID.toString());

    return res;
  } catch (err) {
    logger.e('don`t response. error message: $err');
  }
}
=======
      */