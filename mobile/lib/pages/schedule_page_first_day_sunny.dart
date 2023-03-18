import 'dart:developer';

import 'package:seeft_mobile/configs/importer.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:http/http.dart' as http;

class SchedulePageFirstDaySunny extends StatefulWidget {
  @override
  _SchedulePageFirstDaySunnyState createState() =>
      _SchedulePageFirstDaySunnyState();
}

class _SchedulePageFirstDaySunnyState extends State<SchedulePageFirstDaySunny> {
// notification関連をinitStateに書き出さなきゃいけないので書いてたけどutilとかに書いてもいいかもね

//  FlutterLocalNotificationsPlugin flutterLocalNotificationsPlugin;
//  NotificationDetails platformChannelSpecifics;

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    final Size size = MediaQuery.of(context).size;
    return Container(
      padding: const EdgeInsets.all(40.0),
      child: Container(
        child: SingleChildScrollView(
          scrollDirection: Axis.vertical,
          child: Column(
            children: <Widget>[
              Container(
                child: Table(
                  border: TableBorder.all(color: Colors.black),
                  columnWidths: const <int, TableColumnWidth>{
                    0: FlexColumnWidth(10),
                    1: FlexColumnWidth(15),
                    2: FlexColumnWidth(15),
                    3: FlexColumnWidth(15),
                    4: FlexColumnWidth(15),
                  },
                  defaultVerticalAlignment: TableCellVerticalAlignment.middle,
                  children: [
                    TableRow(
                        decoration: BoxDecoration(color: Colors.teal),
                        children: [
                          TableCell(
                              child: Container(
                            child: Text("日時"),
                            alignment: Alignment.center,
                          )),
                          TableCell(
                            child: Container(
                              child: Text("メインステージ"),
                              alignment: Alignment.center,
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text("体育館"),
                              alignment: Alignment.center,
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text("A講義室"),
                              alignment: Alignment.center,
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text("D講義室"),
                              alignment: Alignment.center,
                            ),
                          ),
                        ]),

                    // 2行目
                    TableRow(children: [
                      TableCell(
                          child: Container(
                        child: Text("10:00~\n10:30"),
                        alignment: Alignment.center,
                      )),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text("悠久太鼓愛好会\n鶴亀会"),
                          alignment: Alignment.center,
                          color: Colors.lightGreen,
                        ),
                      ),
                      TableCell(
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text("Popular Music Club"),
                          alignment: Alignment.center,
                          color: Colors.lightGreen,
                        ),
                      ),
                    ]),

                    // 3行目
                    TableRow(children: [
                      TableCell(
                          child: Container(
                        child: Text("10:30~\n11:00"),
                        alignment: Alignment.center,
                      )),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text("合唱サークル"),
                          alignment: Alignment.center,
                          color: Colors.lightGreen,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                          color: Colors.lightGreen,
                        ),
                      ),
                    ]),

                    // 4行目
                    TableRow(children: [
                      TableCell(
                        child: Container(
                          child: Text("11:00~\n11:15"),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                          color: Colors.lightGreen,
                        ),
                      ),
                    ]),

                    // 5行目
                    TableRow(children: [
                      TableCell(
                        child: Container(
                          child: Text("11:15~\n11:30"),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text("Mexico"),
                          alignment: Alignment.center,
                          color: Colors.lightGreen,
                        ),
                      ),
                      TableCell(
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                          color: Colors.lightGreen,
                        ),
                      ),
                    ]),

                    // 6行目
                    TableRow(children: [
                      TableCell(
                        child: Container(
                          child: Text("11:30~\n12:00"),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                          color: Colors.lightGreen,
                        ),
                      ),
                    ]),

                    // 7行目
                    TableRow(children: [
                      TableCell(
                        child: Container(
                          child: Text("12:00~\n12:30"),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text("Street Style"),
                          alignment: Alignment.center,
                          color: Colors.lightGreen,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                          color: Colors.lightGreen,
                        ),
                      ),
                    ]),

                    // 8行目
                    TableRow(children: [
                      TableCell(
                        child: Container(
                          child: Text("12:30~\n13:00"),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                          color: Colors.lightGreen,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                          color: Colors.lightGreen,
                        ),
                      ),
                    ]),

                    // 9行目
                    TableRow(children: [
                      TableCell(
                        child: Container(
                          child: Text("13:00~\n13:30"),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                          color: Colors.lightGreen,
                        ),
                      ),
                    ]),

                    // 10行目
                    TableRow(children: [
                      TableCell(
                        child: Container(
                          child: Text("13:30~\n14:00"),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text("Sri Lanka"),
                          alignment: Alignment.center,
                          color: Colors.lightGreen,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                          color: Colors.lightGreen,
                        ),
                      ),
                    ]),

                    // 12行目
                    TableRow(children: [
                      TableCell(
                        child: Container(
                          child: Text("14:00~\n14:30"),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                          color: Colors.lightGreen,
                        ),
                      ),
                    ]),

                    // 11行目
                    TableRow(children: [
                      TableCell(
                        child: Container(
                          child: Text("14:30~\n15:00"),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text("吹奏楽部"),
                          alignment: Alignment.center,
                          color: Colors.lightGreen,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                          color: Colors.lightGreen,
                        ),
                      ),
                    ]),

                    // 13行目
                    TableRow(children: [
                      TableCell(
                        child: Container(
                          child: Text("15:00~\n15:30"),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                          color: Colors.lightGreen,
                        ),
                      ),
                    ]),

                    // 14行目
                    TableRow(children: [
                      TableCell(
                        child: Container(
                          child: Text("15:30~\n16:00"),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                          color: Colors.lightGreen,
                        ),
                      ),
                    ]),

                    // 15行目
                    TableRow(children: [
                      TableCell(
                        child: Container(
                          child: Text("16:00~\n16:30"),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                          color: Colors.lightGreen,
                        ),
                      ),
                    ]),

                    // 16行目
                    TableRow(children: [
                      TableCell(
                        child: Container(
                          child: Text("16:30~\n17:00"),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text("ヒーローショー"),
                          alignment: Alignment.center,
                          color: Colors.lightGreen,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                          color: Colors.lightGreen,
                        ),
                      ),
                    ]),

                    // 15行目
                    TableRow(children: [
                      TableCell(
                        child: Container(
                          child: Text("17:00~\n17:30"),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                          color: Colors.lightGreen,
                        ),
                      ),
                    ]),

                    // 17行目
                    TableRow(children: [
                      TableCell(
                        child: Container(
                          child: Text("17:30~\n18:00"),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                        ),
                      ),
                      TableCell(
                        verticalAlignment: TableCellVerticalAlignment.fill,
                        child: Container(
                          child: Text(""),
                          alignment: Alignment.center,
                          color: Colors.lightGreen,
                        ),
                      ),
                    ]),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
