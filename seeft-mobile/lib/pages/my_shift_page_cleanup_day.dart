import 'dart:developer';

import 'package:seeft_mobile/configs/importer.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:http/http.dart' as http;

class MyShiftPageCleanupDay extends StatefulWidget {
  @override
  _MyShiftPageState createState() => _MyShiftPageState();
}

class _MyShiftPageState extends State<MyShiftPageCleanupDay> {
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
    return FutureBuilder(
      future: getData(),
      builder: (ctx, snapshot) {
        if (snapshot.connectionState == AsyncSnapshot.waiting()) {
          logger.w("message");
        }
        if (!snapshot.hasData) {
          return Center(child: CircularProgressIndicator());
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
                      child: table.shiftTable(snapshot.data, context)),
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
    var res = await api.getMyShiftCleanupDay(userID.toString());
    return res;
  } catch (err) {
    logger.e('don`t response. error message: $err');
  }
}
