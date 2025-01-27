import 'dart:developer';
import 'dart:ui';

import 'package:seeft_mobile/configs/importer.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:http/http.dart' as http;

class CleanUpDayTimeSchedule extends StatefulWidget {
  @override
  _MyShiftPageState createState() => _MyShiftPageState();
}

class _MyShiftPageState extends State<CleanUpDayTimeSchedule> {
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
                      7: FixedColumnWidth(100),
                      8: FixedColumnWidth(100),
                      9: FixedColumnWidth(100),
                    },
                    defaultVerticalAlignment: TableCellVerticalAlignment.middle,
                    children: [
                      // 1行目 チーム
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
                                "本部",
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
                                "Eチーム",
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
                                "Fチーム",
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
                                "Gチーム",
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
                                "Hチーム",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                        ],
                      ),

                      // 2行目 リーダー
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
                                "橘",
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
                                "窪坂",
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
                                "山崎(修)",
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
                                "原田",
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
                                "鈴木",
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
                                "黒木",
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
                                "竹内(友)",
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
                                "小武",
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
                                "阿部",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                        ],
                      ),

                      // 3行目 補佐
                      TableRow(
                        decoration: BoxDecoration(color: Colors.teal),
                        children: [
                          TableCell(
                            child: Container(
                              child: Text(
                                "補佐",
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
                                "市原",
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
                                "中井",
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
                                "小林",
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
                                "藤崎",
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
                                "齋藤・善場",
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
                                "野上",
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
                                "小黒",
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
                                "吉川",
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
                                "長島",
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: Colors.white,
                                ),
                              ),
                            ),
                          ),
                        ],
                      ),

                      // 4行目 人数
                      TableRow(
                        decoration: BoxDecoration(color: Colors.teal),
                        children: [
                          TableCell(
                            child: Container(
                              child: Text(
                                "人数",
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
                                "4人",
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
                                "13人",
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
                                "13人",
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
                                "14人",
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
                                "15人",
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
                                "16人",
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
                                "8人",
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
                                "6人",
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
                        ],
                      ),
                      
                      // 集合場所
                      TableRow(
                        decoration: BoxDecoration(color: Colors.teal),                        
                        children: [
                          TableCell(
                            child: Text(
                              "集合場所",
                              textAlign: TextAlign.center,
                            ),
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "本部",
                                textAlign: TextAlign.center,
                              )
                            )
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "講義棟掲示板前",
                                textAlign: TextAlign.center,
                              )
                            )
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "講義棟掲示板前",
                                textAlign: TextAlign.center,
                              )
                            )
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "講義棟掲示板前",
                                textAlign: TextAlign.center,
                              )
                            )
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "D講義室前",
                                textAlign: TextAlign.center,
                              )
                            )
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "体育館",
                                textAlign: TextAlign.center,
                              )
                            )
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "体育館",
                                textAlign: TextAlign.center,
                              )
                            )
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "24",
                                textAlign: TextAlign.center,
                              )
                            )
                          ),
                          TableCell(
                            child: Container(
                              child: Text(
                                "体育館前",
                                textAlign: TextAlign.center,
                              )
                            )
                          ),
                        ],
                      ),
                      // 5行目 ここからシフト 8:00~8:30
                      TableRow(
                          decoration: BoxDecoration(color: Colors.grey[200]),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "8:00",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "朝食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "朝食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "朝食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "朝食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "朝食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "朝食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "朝食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "朝食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "朝食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                          ]),

                      // 6行目 8:30~9:00
                      TableRow(
                          decoration: BoxDecoration(color: Colors.grey[200]),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "8:30",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "MT",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "MT",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "MT",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "MT",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "MT",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "MT",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "MT",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "MT",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "MT",
                                textAlign: TextAlign.center,
                              ),
                            )),
                          ]),

                      // 7行目 9:00~9:30
                      TableRow(
                          decoration: BoxDecoration(color: Colors.white60),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "9:00",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                              child: Container(
                                child: new Text(
                                  "本部",
                                  textAlign: TextAlign.center,
                                ),
                              ),
                            ),
                            TableCell(
                                child: Container(
                                  child: new Text(
                                    "掲示板・パーディション移動\nAL3・外",
                                    textAlign: TextAlign.center,
                                  ),
                                ),
                            ),
                            TableCell(
                                child: Container(
                                  child: new Text(
                                    "掲示板・パーディション移動\nAL3・外",
                                    textAlign: TextAlign.center,
                                  ),
                                ),
                            ),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "掲示板・パーディション移動\nAL3・外",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "物品運搬\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "メインステージ解体\n\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "電ドラはがし\n\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "制作物撤去",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "検温所撤去\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                          ]),

                      // 8行目 9:30~10:00
                      TableRow(
                          decoration: BoxDecoration(color: Colors.white60),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "9:30",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                              child: Container(
                                child: new Text(
                                  "本部",
                                  textAlign: TextAlign.center,
                                ),
                              ),
                            ),
                            TableCell(
                                child: Container(
                                  child: new Text(
                                    "掲示板・パーディション移動\nAL3・外",
                                    textAlign: TextAlign.center,
                                  ),
                                ),
                            ),
                            TableCell(
                                child: Container(
                                  child: new Text(
                                    "掲示板・パーディション移動\nAL3・外",
                                    textAlign: TextAlign.center,
                                  ),
                                ),
                            ),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "掲示板・パーディション移動\nAL3・外",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "物品運搬\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "椅子・音響片付け\n\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "電ドラはがし\n\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "制作物撤去",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "検温所撤去\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                          ]),

                      // 9行目 10:00~10:30
                      TableRow(
                          decoration: BoxDecoration(color: Colors.white60),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "10:00",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "本部",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "机・イス運搬\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "机・イス運搬\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "机・イス運搬\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                           TableCell(
                                child: Container(
                              child: new Text(
                                "物品運搬\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "メインステージ解体\n\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "電ドラはがし\n\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "制作物撤去",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "検温所撤去\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                          ]),

                      // 10行目 10:30~11:00
                      TableRow(
                          decoration: BoxDecoration(color: Colors.white60),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "10:30",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "本部",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "机・イス運搬\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "机・イス運搬\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "机・イス運搬\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                           TableCell(
                                child: Container(
                              child: new Text(
                                "物品運搬\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "メインステージ解体\n\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "電ドラはがし\n\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "制作物撤去",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "検温所撤去\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                         ]),

                      // 11行目 11:00~11:30
                      TableRow(
                          decoration: BoxDecoration(color: Colors.white60),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "11:00",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "本部",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "机・イス運搬\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "机・イス運搬\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "机・イス運搬\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                           TableCell(
                                child: Container(
                              child: new Text(
                                "昼食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "メインステージ解体\n\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "電ドラはがし\n\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "制作物撤去",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "D講・AL3片付け",
                                textAlign: TextAlign.center,
                              ),
                            )),
                         ]),

                      // 12行目 11:30~12:00
                      TableRow(
                          decoration: BoxDecoration(color: Colors.white60),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "11:30",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "本部",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "机・イス運搬\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "机・イス運搬\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "机・イス運搬\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                           TableCell(
                                child: Container(
                              child: new Text(
                                "昼食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "メインステージ解体\n\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "電ドラはがし\n\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "制作物撤去",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "D講・AL3片付け",
                                textAlign: TextAlign.center,
                              ),
                            )),                          ]),

                      // 13行目 12:00~12:30
                      TableRow(
                          decoration: BoxDecoration(color: Colors.white60),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "12:00",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "昼食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "昼食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "昼食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "昼食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                           TableCell(
                                child: Container(
                              child: new Text(
                                "物品運搬\n\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "昼食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "電ドラはがし\n\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "昼食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "昼食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                                                     ]),

                      // 14行目 12:30~13:00
                      TableRow(
                          decoration: BoxDecoration(color: Colors.white60),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "12:30",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "昼食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "昼食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "昼食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "昼食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                           TableCell(
                                child: Container(
                              child: new Text(
                                "物品運搬\n\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "昼食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "昼食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "昼食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "昼食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                          ]),

                      // 15行目 13:00~13:30
                      TableRow(
                          decoration: BoxDecoration(color: Colors.white60),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "13:00",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "本部",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "机・イス運搬, 復帰\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "机・イス運搬, 復帰\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "机・イス運搬, 復帰\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                           TableCell(
                                child: Container(
                              child: new Text(
                                "物品運搬\n\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "体育館ステージ撤去\n\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "昼食",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "24下倉庫・セコム片付け",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "24下倉庫・セコム片付け",
                                textAlign: TextAlign.center,
                              ),
                            )),
                          ]),

                      // 16行目 13:30~14:00
                      TableRow(
                          decoration: BoxDecoration(color: Colors.white60),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "13:30",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "本部",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "机・イス運搬, 復帰\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "机・イス運搬, 復帰\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "机・イス運搬, 復帰\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                           TableCell(
                                child: Container(
                              child: new Text(
                                "物品運搬\n\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "体育館ステージ撤去\n\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "本部解体",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "24下倉庫・セコム片付け",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "24下倉庫・セコム片付け",
                                textAlign: TextAlign.center,
                              ),
                            )), 
                                                     ]),

                      // 17行目 14:00~14:30
                      TableRow(
                          decoration: BoxDecoration(color: Colors.white60),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "14:00",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "本部",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "机・イス運搬, 復帰\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "机・イス運搬, 復帰\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "机・イス運搬, 復帰\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                           TableCell(
                                child: Container(
                              child: new Text(
                                "物品運搬\n\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "体育館ステージ撤去\n\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "本部解体",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "24下倉庫・セコム片付け",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "24下倉庫・セコム片付け",
                                textAlign: TextAlign.center,
                              ),
                            )), 
                          ]),

                      // 18行目 14:30~15:00
                      TableRow(
                          decoration: BoxDecoration(color: Colors.white60),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "14:30",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "本部",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "机・イス運搬, 復帰\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "机・イス運搬, 復帰\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "机・イス運搬, 復帰\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                           TableCell(
                                child: Container(
                              child: new Text(
                                "物品運搬\n\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "体育館ステージ撤去\n\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "本部解体",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "24下倉庫・セコム片付け",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "24下倉庫・セコム片付け",
                                textAlign: TextAlign.center,
                              ),
                            )), 
                           ]),

                      // 19行目 15:00~15:30
                      TableRow(
                          decoration: BoxDecoration(color: Colors.white60),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "15:00",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "本部",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "シールはがし・掃除\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "シールはがし・掃除\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "シールはがし・掃除\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                           TableCell(
                                child: Container(
                              child: new Text(
                                "シールはがし・掃除\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "シートはがし\n\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "本部解体",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "24下倉庫・セコム片付け",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "24下倉庫・セコム片付け",
                                textAlign: TextAlign.center,
                              ),
                            )), 
                           ]),

                      // 20行目 15:30~16:00
                      TableRow(
                          decoration: BoxDecoration(color: Colors.white60),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "15:30",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "本部",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "シールはがし・掃除\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "シールはがし・掃除\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "シールはがし・掃除\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                           TableCell(
                                child: Container(
                              child: new Text(
                                "シールはがし・掃除\n\n講義棟",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "シートはがし\n\n体育館",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "本部解体",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "24下倉庫・セコム片付け",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "24下倉庫・セコム片付け",
                                textAlign: TextAlign.center,
                              ),
                            )), 
                          ]),

                          // 21行目 集合写真
                          TableRow(
                          decoration: BoxDecoration(color: Colors.grey[200]),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "16:00",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "集合写真",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "集合写真",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "集合写真",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "集合写真",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "集合写真",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "集合写真",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "集合写真",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "集合写真",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "集合写真",
                                textAlign: TextAlign.center,
                              ),
                            )),
                          ]),

                          // 22行目 ビンゴ受付
                          TableRow(
                          decoration: BoxDecoration(color: Colors.orange[100]),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "16:30",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ受付",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ受付",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ受付",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ受付",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ受付",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ受付",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ受付",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ受付",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ受付",
                                textAlign: TextAlign.center,
                              ),
                            )),
                          ]),

                          // 23行目 ビンゴ
                          TableRow(
                          decoration: BoxDecoration(color: Colors.orange[100]),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "17:00",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ",
                                textAlign: TextAlign.center,
                              ),
                            )),
                          ]),

                          // 24行目 ビンゴ
                          TableRow(
                          decoration: BoxDecoration(color: Colors.orange[100]),
                          children: [
                            TableCell(
                                child: Container(
                              child: new Text(
                                "17:30",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ",
                                textAlign: TextAlign.center,
                              ),
                            )),
                            TableCell(
                                child: Container(
                              child: new Text(
                                "ビンゴ",
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
