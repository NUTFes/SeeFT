import 'package:seeft_mobile/configs/importer.dart';

final ShiftTable table = ShiftTable();

class ShiftTable {
  Widget shiftTable(var shifts, context) {
    var sortedShifts =
        shifts.sort((i, v) => i["time"]["time"].compareTo(v["time"]["time"]));
    // print(shifts);
    return Table(
        border: TableBorder.all(color: Colors.black26),
        columnWidths: const <int, TableColumnWidth>{
          // 0: IntrinsicColumnWidth(),
          0: FlexColumnWidth(1),
          1: FlexColumnWidth(3),
          // 2: FixedColumnWidth(100.0),
        },
        defaultVerticalAlignment: TableCellVerticalAlignment.middle,
        children: [
          TableRow(children: [
            TableCell(
                child: Container(
              child: Text(
                "日時",
                style: TextStyle(
                  color: Colors.white,
                ),
              ),
              alignment: Alignment.center,
              color: Colors.orangeAccent,
            )),
            TableCell(
              child: Container(
                child: Text(
                  "シフト",
                  style: TextStyle(
                    color: Colors.white,
                  ),
                ),
                alignment: Alignment.center,
                color: Colors.orangeAccent,
              ),
            )
          ]),
          for (var index = 0; index < sortedShifts.length; index++)
            TableRow(
                decoration: BoxDecoration(color: Colors.white60),
                children: [
                  TableCell(
                      child: Container(
                    alignment: Alignment.center,
                    child: new Text(sortedShifts[index]["time"]["time"]),
                  )),
                  TableCell(
                      /*
                      child: Container(
                    alignment: Alignment.center,
                    child: new Text(shift["Work"].toString()),
                    // margin: EdgeInsets.only(bottom: 10.0),
                    height: 25,
                    */
                      child: Container(
                    //color: HexColor(shifts[index]["Color"].toString()),
                    width: double.infinity,
                    height: 40,
                    //color: HexColor(shifts[index]["Color"].toString()),
                    padding: EdgeInsets.all(0),
                    //height: 25.0,
                    child: new Material(
                      type: MaterialType.button,
                      color: HexColor(sortedShifts[index]["task"]["color"]),
                      child: InkWell(
                        splashColor: Colors.orangeAccent,
                        onTap: () async {
                          if (sortedShifts[index]["task"]["task"] != "") {
                            logger.i(sortedShifts[index]["task"]["task"]);
                            await openShiftDialog(
                                context,
                                sortedShifts[index]["task"],
                                sortedShifts[index]["user"],
                                sortedShifts[index]["year"],
                                sortedShifts[index]["date"],
                                sortedShifts[index]["time"],
                                sortedShifts[index]["weather"]);
                          }
                        },
                        //child: Center(child: new Text(shifts[index]["Work"].toString())),
                        child: Container(
                          child: Center(
                              child: new Text(sortedShifts[index]["task"]
                                      ["task"]
                                  .toString())),
                        ),
                      ),
                    ),
                  ))
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
