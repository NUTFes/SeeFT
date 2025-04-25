import 'dart:developer';
import 'package:seeft_mobile/configs/importer.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:http/http.dart' as http;
import 'package:seeft_mobile/pages/wait_page.dart';
import 'package:seeft_mobile/theme/tokens.dart';

final ShiftRequestTable table = ShiftRequestTable();

class ShiftRequestPage extends StatefulWidget {
  @override
  _ShiftRequestPageState createState() => _ShiftRequestPageState();
}

class _ShiftRequestPageState extends State<ShiftRequestPage> {

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    final Size size = MediaQuery.of(context).size;
    return FutureBuilder(
      future: getData(),
      builder: (ctx, snapshot) {
        if (snapshot.connectionState == AsyncSnapshot.waiting()) {
          logger.w("message");
        }
        // if (!snapshot.hasData) {
        //   // 待機画面を表示
        //   return WaitPage();
        // }
        return Container(
            padding: const EdgeInsets.all(40.0),
            color: AppColors.base,
            child: SingleChildScrollView(
              child: Column(
                children: <Widget>[
                  Text("参加できない時間を選択してください", textAlign: TextAlign.start),
                  Container(
                    child: table.shiftRequestTable(context)),
                ],
              ),
            ));
      },
    );
  }
}

Future getData() async {
  try {
    var userID = await store.getUserID();
    var res = await api.getMyShiftCurrentFirstDayRainy(userID.toString());
    return res;
  } catch (err) {
    logger.e('don`t response. error message: $err');
  }
}


class ShiftRequestTable {
  // Widget shiftRequestTable(var shifts, context) {
  Widget shiftRequestTable(context) {
    // print(shifts);
    return Table(
        border: TableBorder.all(color: Colors.black26),
        columnWidths: const <int, TableColumnWidth>{
          // 0: IntrinsicColumnWidth(),
          // 0: FlexColumnWidth(1),
          // 1: FlexColumnWidth(1),
          // 2: FixedColumnWidth(100.0),
        },
        defaultVerticalAlignment: TableCellVerticalAlignment.middle,
        children: [
          TableRow(
            decoration: BoxDecoration(color: AppColors.base),
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
          for (var hour = 6; hour < 22; hour++)
            for (var minute = 0; minute < 60; minute += 15)
              TableRow(
                decoration: BoxDecoration(color: AppColors.base),
                children: [
                  TableCell(
                    child: Container(
                      alignment: Alignment.center,
                      child: minute == 0
                          ? Text(hour.toString() + ":00 ~ " + hour.toString() + ":15")
                          : minute == 45
                          ? Text(hour.toString() + ":45 ~ " + (hour + 1).toString() + ":00")
                          : Text(hour.toString() + ":" + minute.toString() + " ~ " + hour.toString() + ":" + (minute + 15).toString()),
                  )),
                  TableCell(
                      child: Container(
                    alignment: Alignment.center,
                    child: new Text(""),
                  )),
                  TableCell(
                      child: Container(
                    alignment: Alignment.center,
                    child: new Text(""),
                  )),
                  TableCell(
                      child: Container(
                    alignment: Alignment.center,
                    child: new Text(""),
                  )),
                  TableCell(
                      child: Container(
                    alignment: Alignment.center,
                    child: new Text(""),
                  )),
              ]),
        ]);
  }
}

class HexColor extends Color {
  static int _getColorFromHex(String hexColor) {
    // logger.d(hexColor);
    hexColor = hexColor.toUpperCase().replaceAll('0X', '');
    // logger.d(hexColor);
    if (hexColor.length == 6) {
      hexColor = 'FF' + hexColor;
    }
    // logger.d(hexColor);
    return int.parse(hexColor, radix: 16);
  }

  HexColor(final String hexColor) : super(_getColorFromHex(hexColor));
}
