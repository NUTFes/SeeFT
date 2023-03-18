import 'dart:developer';
import 'dart:ui';

import 'package:seeft_mobile/configs/importer.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:http/http.dart' as http;

class PrePreparationDayShift extends StatefulWidget {
  @override
  _MyShiftPageState createState() => _MyShiftPageState();
}

class _MyShiftPageState extends State<PrePreparationDayShift> {
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
      child: SingleChildScrollView(
        scrollDirection: Axis.horizontal,
        child: Container(
          child: SingleChildScrollView(
            scrollDirection: Axis.vertical,
            child: Column(
              children: <Widget>[
                Container(
                  child: Table(
                    border: TableBorder.all(color: Colors.black),
                    columnWidths: const <int, TableColumnWidth>{
                      // 0: IntrinsicColumnWidth(),
                      0: FixedColumnWidth(80),
                      1: FixedColumnWidth(100),
                      2: FixedColumnWidth(100),
                      3: FixedColumnWidth(100),
                      4: FixedColumnWidth(100),
                      5: FixedColumnWidth(100),
                      6: FixedColumnWidth(100),
                    },
                    defaultVerticalAlignment: TableCellVerticalAlignment.middle,
                    children: [
                      // 1行目
                      TableRow(
                        decoration: BoxDecoration(color: Colors.teal),
                        children: [
                          TableCell(
                            child: Container(
                              child: Text(
                                "",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "Aチーム",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "Bチーム",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "Cチーム",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "Dチーム",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "トラック",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "Eチーム",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                        ],
                      ),

                      // 2行目
                      TableRow(
                        decoration: BoxDecoration(color: Colors.teal),
                        children: [
                          TableCell(
                            child: Container(
                              child: Text(
                                "リーダー",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "政木",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "真下",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "大場",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "宇野",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "本田さん",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "氏家",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                        ],
                      ),

                      // 3行目
                      TableRow(
                        decoration: BoxDecoration(color: Colors.teal),
                        children: [
                          TableCell(
                            child: Container(
                              child: Text(
                                "想定人数",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "10人",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "10人",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "10人",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "12人",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "20人",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                        ],
                      ),

                      // 4行目
                      TableRow(
                        decoration: BoxDecoration(color: Colors.teal),
                        children: [
                          TableCell(
                            child: Container(
                              child: Text(
                                "集合場所",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "講義棟\n掲示板前",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "講義棟\n掲示板前",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "新講義棟1F",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "24",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "体育館",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                        ],
                      ),

                      // 5行目
                      TableRow(
                        decoration: BoxDecoration(color: Colors.teal),
                        children: [
                          TableCell(
                            child: Container(
                              child: Text(
                                "備考",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "制作局",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "企画局+α",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                        ],
                      ),

                      // 6行目 ここからシフト 10:00~10:30
                      TableRow(
                          decoration: BoxDecoration(color: Colors.white60),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "10:00\n~10:30",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "集合, MT",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "集合, MT",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "集合, MT",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "集合, MT",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "集合, MT",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "集合, MT",
                                textAlign: TextAlign.center,
                              ),
                            )),
                          ]),

                      // 7行目 10:30~11:00
                      TableRow(
                          decoration: BoxDecoration(color: Colors.white60),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "10:30\n~11:00",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "シール貼り",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "シール貼り",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "AL2,AL3整理\nD講物品移動",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "入り口看板\nウェルカムボード制作",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "休憩",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "体育館\nマット張り",
                                textAlign: TextAlign.center,
                              ),
                            )),
                          ]),

                      // 8行目 11:00~11:30
                      TableRow(
                          decoration: BoxDecoration(color: Colors.white60),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "11:00\n~11:30",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "3F整理",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "2F\n机・椅子移動",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "AL2,AL3整理\nD講物品移動",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "入り口看板\nウェルカムボード制作",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "休憩",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "体育館\nマット張り",
                                textAlign: TextAlign.center,
                              ),
                            )),
                          ]),

                      // 9行目 11:30~12:00
                      TableRow(
                          decoration: BoxDecoration(color: Colors.white60),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "11:30\n~12:00",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "3F整理",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "2F\n机・椅子移動",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "休憩",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "入り口看板\nウェルカムボード制作",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "休憩",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "体育館\nマット張り",
                                textAlign: TextAlign.center,
                              ),
                            )),
                          ]),

                      // 10行目 12:00~12:30
                      TableRow(
                          decoration: BoxDecoration(color: Colors.white60),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "12:00\n~12:30",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "3F整理",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "休憩",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "休憩",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "休憩",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "休憩",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "体育館\nマット張り",
                                textAlign: TextAlign.center,
                              ),
                            )),
                          ]),

                      // 11行目 12:30~13:00
                      TableRow(
                          decoration: BoxDecoration(color: Colors.white60),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "12:30\n~13:00",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "休憩",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "休憩",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "1F\n机・椅子移動",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "休憩",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "休憩",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "休憩",
                                textAlign: TextAlign.center,
                              ),
                            )),
                          ]),

                      // 12行目 13:00~13:30
                      TableRow(
                          decoration: BoxDecoration(color: Colors.white60),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "13:00\n~13:30",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "休憩",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "2F\n机・椅子移動",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "1F\n机・椅子移動",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "入試課看板・\nテント運搬",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "Dチーム同行",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "休憩",
                                textAlign: TextAlign.center,
                              ),
                            )),
                          ]),

                      // 13行目 13:30~14:00
                      TableRow(
                          decoration: BoxDecoration(color: Colors.white60),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "13:30\n~14:00",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "24下倉庫\nセコム物品運搬",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "2F\n机・椅子移動",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "1F\n机・椅子移動",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "入試課看板作成",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "Aチーム同行",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "体育館\nステージ設営",
                                textAlign: TextAlign.center,
                              ),
                            )),
                          ]),

                      // 14行目 14:00~14:30
                      TableRow(
                          decoration: BoxDecoration(color: Colors.white60),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "14:00\n~14:30",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "24下倉庫\nセコム物品運搬",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "2F\n机・椅子移動",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "1F\n机・椅子移動",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "入試課看板作成",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "Aチーム同行",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "体育館\nステージ設営",
                                textAlign: TextAlign.center,
                              ),
                            )),
                          ]),

                      // 15行目 14:30~15:00
                      TableRow(
                          decoration: BoxDecoration(color: Colors.white60),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "14:30\n~15:00",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "24下倉庫\nセコム物品運搬",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "1F\n机・椅子移動",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "入試課看板作成",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "Aチーム同行",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "体育館\nステージ設営",
                                textAlign: TextAlign.center,
                              ),
                            )),
                          ]),

                      // 16行目 15:00~15:30
                      TableRow(
                          decoration: BoxDecoration(color: Colors.white60),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "15:00\n~15:30",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "体育館\nステージ設営",
                                textAlign: TextAlign.center,
                              ),
                            )),
                          ]),
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
