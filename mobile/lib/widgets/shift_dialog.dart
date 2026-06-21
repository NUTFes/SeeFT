import 'package:seeft_mobile/configs/importer.dart';
import 'package:url_launcher/url_launcher.dart';

getShiftDetail(workId, userId, date, weather, time) async {
  // 本当は引数はworkIdだけでいいが、今後ほかの要素を使うことも考えてそのままに
  try {
    logger.w(workId);
    var res = await api.shiftDetail(workId, userId, date, weather, time);
    return res;
  } catch (e) {}
}

_launchURL(url) async {
  if (await canLaunchUrl(Uri.parse(url))) {
    await launchUrl(Uri.parse(url));
  } else {
    throw 'Unable to launch url $url';
  }
}

openShiftDialog(
    // BuildContext context, workId, userId, date, weather, time) async {
    BuildContext context,
    task,
    user,
    year,
    date,
    time,
    weather) async {
  // var res = await getShiftDetail(workId, userId, date, weather, time);
  // logger.i(res);
  // var resName = res["task"];
  // var resURL = res["url"];
  // var resUsers = res["users"];
  // var resPlace = res["place"];
  // var resPresident = res["superviser"];
  // var resPresidentTel = res["TEL"];

  // 現在のシフトのメンバー
  var res = await api.getUsersByShift(
      task["id"].toString(),
      year["id"].toString(),
      date["id"].toString(),
      time["id"].toString(),
      weather["id"].toString());
  var resName = task["task"].toString();
  var resURL = task["url"].toString();
  var resPlace = task["place"].toString();
  // var resPresident = task["superviser"].toString();
  var resUsersNumber = res["users"].length;
  List<String> resUsersList = <String>[];
  for (var index = 0; index < resUsersNumber; index++) {
    resUsersList.add(res["users"][index]["name"].toString());
  }
  var resUsers = resUsersList.join(",");

  // 1つ前の時間のシフトのメンバー
  var resBeforeUsers = '';
  if (time["id"] > 1) {
    try {
      var beforeRes = await api.getUsersByShift(
          task["id"].toString(),
          year["id"].toString(),
          date["id"].toString(),
          (time["id"] - 1).toString(),
          weather["id"].toString());
      List<String> resBeforeUsersList = <String>[];
      var resBeforeUsersNumber = beforeRes["users"].length;
      for (var index = 0; index < resBeforeUsersNumber; index++) {
        resBeforeUsersList.add(beforeRes["users"][index]["name"].toString());
      }
      resBeforeUsers = resBeforeUsersList.join(",");
    } catch (e) {
      resBeforeUsers = 'none';
    }
  } else {
    resBeforeUsers = 'none';
  }

  // 1つ後の時間のシフトのメンバー
  var resAfterUsers = '';
  if (time["id"] != 96) {
    try {
      var afterRes = await api.getUsersByShift(
          task["id"].toString(),
          year["id"].toString(),
          date["id"].toString(),
          (time["id"] + 1).toString(),
          weather["id"].toString());
      List<String> resAfterUsersList = <String>[];
      var resAfterUsersNumber = afterRes["users"].length;
      for (var index = 0; index < resAfterUsersNumber; index++) {
        resAfterUsersList.add(afterRes["users"][index]["name"].toString());
      }
      resAfterUsers = resAfterUsersList.join(",");
    } catch (e) {
      resAfterUsers = 'none';
    }
  } else {
    resAfterUsers = 'none';
  }

  showDialog(
    context: context,
    builder: (context) {
      return SimpleDialog(
        contentPadding: EdgeInsets.zero,
        titlePadding: EdgeInsets.zero,
        title: SizedBox(
          height: 550,
          child: Scaffold(
            appBar: AppBar(
              title: Text(resName), //シフト名
              centerTitle: true,
              actions: <Widget>[
                IconButton(
                  onPressed: () async {
                    _launchURL(resURL);
                  },
                  icon: Icon(Icons.wrap_text),
                  color: Colors.orangeAccent[100],
                ),
              ],
            ),
            body: Container(
              child: ListView(
                children: <Widget>[
                  ListTile(
                    leading: Icon(Icons.place_outlined),
                    title: Text("集合場所"),
                    subtitle: Text(resPlace),
                  ),
                  // 今回はいらないので一時的にコメントアウト
                  // ListTile(
                  //   leading: Icon(Icons.supervisor_account_outlined),
                  //   title: Text("代表者"),
                  //   subtitle: Text(resPresident),
                  // ),
                  // ListTile(
                  //   leading: Icon(Icons.phone),
                  //   title: Text("緊急時連絡先"),
                  //   subtitle: Text(resPresidentTel),
                  // ),
                  ListTile(
                    leading: Icon(Icons.group),
                    title: Text("メンバー"),
                    subtitle: Text(resUsers),
                  ),
                  ListTile(
                    leading: Icon(Icons.group),
                    title: Text("前のメンバー"),
                    subtitle: Text(resBeforeUsers),
                  ),
                  ListTile(
                    leading: Icon(Icons.group),
                    title: Text("次のメンバー"),
                    subtitle: Text(resAfterUsers),
                  ),
                ],
              ),
            ),
            /*
            floatingActionButton: OutlinedButton(
              child: const Text('マニュアルへ'),
              style: OutlinedButton.styleFrom(
                primary: Colors.orangeAccent,
                side: const BorderSide(color: Colors.orangeAccent),
              ),
              onPressed: () async {
                //_launchURL(resURL);
              },
            ),
            */
          ),
        ),
      );
    },
  );
}
