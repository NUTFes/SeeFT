import 'dart:developer';

import 'package:seeft_mobile/configs/importer.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:http/http.dart' as http;
import 'package:url_launcher/url_launcher.dart';

class ManualListPage extends StatefulWidget {
  @override
  _ManualListPageState createState() => _ManualListPageState();
}

class _ManualListPageState extends State<ManualListPage> {
// notification関連をinitStateに書き出さなきゃいけないので書いてたけどutilとかに書いてもいいかもね

//  FlutterLocalNotificationsPlugin flutterLocalNotificationsPlugin;
//  NotificationDetails platformChannelSpecifics;
  int manualLength = 0;

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    final Size size = MediaQuery.of(context).size;
    return Scaffold(
      appBar: AppBar(
        title: const Text('マニュアル一覧'),
        actions: <Widget>[],
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
            const Center(child: CircularProgressIndicator());
          }
          var manualList = snapshot.data;
          return Container(
            padding: const EdgeInsets.all(40.0),
            child: Column(
              children: <Widget>[
                Flexible(
                  child: ListView.builder(
                    itemCount: manualLength,
                    itemBuilder: (BuildContext context, int index) {
                      return Container(
                          height: 40,
                          child: _manualItem(snapshot.data, index, context));
                    },
                  ),
                ),
              ],
            ),
          );
        },
      ),
    );
  }

  Widget _manualItem(var manuals, index, context) {
    return Container(
      decoration: new BoxDecoration(
          border:
              new Border(bottom: BorderSide(width: 1.0, color: Colors.grey))),
      child: ListTile(
        title: Text(
          manuals[index]["task"].toString(),
          style: TextStyle(color: Colors.black, fontSize: 14.0),
        ),
        onTap: () async {
          if (await canLaunch(manuals[index]["url"].toString())) {
            await launch((manuals[index]["url"].toString()));
          }
        },
      ),
    );
  }
  Future getData() async {
    try {
      var res = await api.getAllManual();
      manualLength = res.length;
      return res;
    } catch (err) {
      logger.e('don`t response. error message: $err');
    }
  }
}
