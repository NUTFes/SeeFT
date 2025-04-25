import 'dart:developer';
import 'package:seeft_mobile/configs/importer.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:http/http.dart' as http;
import 'package:seeft_mobile/pages/wait_page.dart';
import 'package:seeft_mobile/theme/tokens.dart';

class ShiftRequestPage extends StatefulWidget {
  @override
  _ShiftRequestPageState createState() => _ShiftRequestPageState();
}

class _ShiftRequestPageState extends State<ShiftRequestPage> {
  bool isLoading = true;
  // シフト希望のデータを格納する配列
  List shiftReqests = <List<bool>>[
    List.filled(64, true),
    List.filled(64, true),
    List.filled(64, true),
    List.filled(64, true),
  ];
    
  @override
  void initState() {
    super.initState();
    fetchData();
  }
  // シフト希望のデータを取得する関数
  Future<void> fetchData() async {
    final result = await getData(); // 非同期で配列取得
    setState(() {
      // shiftReqests = result;
      isLoading = false;
    });
  }
  
  // シフト希望のデータを変更する関数
  void editShiftRequest(int day, int index) {
    setState(() {
      // タップしたセルの参加可/不可を反転
      shiftReqests[day][index] = !shiftReqests[day][index];
    });
  }
  
  // シフト希望のデータを送信する関数
  void sendShiftRequest() {
  }

  @override
  Widget build(BuildContext context) {
    if(isLoading) {
      return WaitPage();  // データ取得中は待機画面を表示
    }
    
    return Container(
      padding: const EdgeInsets.all(32.0),
      decoration: BoxDecoration(
        color: AppColors.base,
      ),
      child: Column(
        spacing: 16,
        children: [
          Text("参加できない時間を選択してください", style: TextStyle(color: AppColors.textBlack, fontSize: AppFontSizes.sm)),
          Expanded(
            child: SingleChildScrollView(
              child: shiftRequestTable(shiftReqests, editShiftRequest, context),  // シフト希望のテーブルを表示
            ),
          ),
          ElevatedButton(
            onPressed: () {
              sendShiftRequest();
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.main,
              padding: const EdgeInsets.symmetric(horizontal: 32, vertical: 16),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(8),
              ),
            ),
            child: Text("シフト希望を送信", style: TextStyle(color: AppColors.textWhite, fontSize: AppFontSizes.sm))
          ),
        ],
      )
    );
  }
}

// データ取得のための関数
Future getData() async {
  // 1秒待機
  await Future.delayed(const Duration(seconds: 1));
  // try {
  //   var userID = await store.getUserID();
  //   var res = await api.getMyShiftCurrentFirstDayRainy(userID.toString());
  //   return res;
  // } catch (err) {
  //   logger.e('don`t response. error message: $err');
  // }
}


// シフト希望のテーブルを表示するWidget
Widget shiftRequestTable(var shiftRequests, Function editShiftRequest, context) {
  return Table(
      border: TableBorder.all(color: AppColors.grayLight),
      columnWidths: const <int, TableColumnWidth>{
        0: IntrinsicColumnWidth(),
        // 1: FlexColumnWidth(1),
      },
      defaultVerticalAlignment: TableCellVerticalAlignment.middle,
      children: [
        // ヘッダー
        TableRow(
          children: [
            TableCell(
              child: Container()
            ),
            TableCell(
              child: Container(
                child: Text(
                  "準備日",
                  style: TextStyle(
                    color: AppColors.textBlack,
                  ),
                ),
              alignment: Alignment.center,
            )),
            TableCell(
              child: Container(
                child: Text(
                  "1日目",
                  style: TextStyle(
                    color: AppColors.textBlack,
                  ),
                ),
              alignment: Alignment.center,
            )),
            TableCell(
              child: Container(
                child: Text(
                  "2日目",
                  style: TextStyle(
                    color: AppColors.textBlack,
                  ),
                ),
              alignment: Alignment.center,
            )),
            TableCell(
              child: Container(
                child: Text(
                  "片付け日",
                  style: TextStyle(
                    color: AppColors.textBlack,
                  ),
                ),
              alignment: Alignment.center,
            )),
          ]
        ),
        for (var index = 0; index < shiftRequests[0].length; index++)
          // シフト希望のテーブル
          TableRow(
            children: [
              TableCell(
                child: Container(
                  alignment: Alignment.center,
                  child: Text(
                    // indexから時間を計算
                    (6 + index ~/ 4).toString() + ":" + (index * 15 % 60).toString() + (index % 4 == 0? "0":"") + "~" + (6 + (index + 1) ~/ 4).toString() + ":" + ((index + 1) * 15 % 60).toString() + ((index + 1) % 4 == 0? "0":""),
                    style: TextStyle(color: AppColors.textBlack,),
                  )
              )),
              GestureDetector(  // 準備日
                onTap: () {
                  editShiftRequest(0, index); // タップしたセルの参加可/不可を反転
                  print(shiftRequests[0][index].toString());
                },
                child: AnimatedContainer(     // アニメーション付きのコンテナ
                  duration: Duration(milliseconds: 200),
                  height: 40,
                  alignment: Alignment.center,
                  decoration: BoxDecoration(
                    color: shiftRequests[0][index]? null: const Color.fromARGB(255, 196, 91, 91),
                  ),
                  child: shiftRequests[0][index]
                    ? null
                    : Icon(Icons.close_outlined, color: const Color.fromARGB(255, 119, 17, 17)),
                ),
              ),
              GestureDetector(  // 1日目
                onTap: () {
                  editShiftRequest(1, index);
                  print(shiftRequests[1][index].toString());
                },
                child: AnimatedContainer(
                  duration: Duration(milliseconds: 200),
                  height: 40,
                  alignment: Alignment.center,
                  decoration: BoxDecoration(
                    color: shiftRequests[1][index]? null: const Color.fromARGB(255, 196, 91, 91),
                  ),
                  child: shiftRequests[1][index]
                    ? null
                    : Icon(Icons.close_outlined, color: const Color.fromARGB(255, 119, 17, 17)),
                ),
              ),
                GestureDetector(  // 2日目
                  onTap: () {
                    editShiftRequest(2, index);
                    print(shiftRequests[2][index].toString());
                  },
                child: AnimatedContainer(
                  duration: Duration(milliseconds: 200),
                  height: 40,
                  alignment: Alignment.center,
                  decoration: BoxDecoration(
                    color: shiftRequests[2][index]? null: const Color.fromARGB(255, 196, 91, 91),
                  ),
                  child: shiftRequests[2][index]
                    ? null
                    : Icon(Icons.close_outlined, color: const Color.fromARGB(255, 119, 17, 17)),
                ),
              ),
              GestureDetector(  // 片付け日
                onTap: () {
                  editShiftRequest(3, index);
                  print(shiftRequests[3][index].toString());
                },
                child: AnimatedContainer(
                  duration: Duration(milliseconds: 200),
                  height: 40,
                  alignment: Alignment.center,
                  decoration: BoxDecoration(
                    color: shiftRequests[3][index]? null: const Color.fromARGB(255, 196, 91, 91),
                  ),
                  child: shiftRequests[3][index]
                    ? null
                    : Icon(Icons.close_outlined, color: const Color.fromARGB(255, 119, 17, 17)),
                ),
              ),
            ]),
      ]);
}
