// import 'dart:developer';

import 'package:seeft_mobile/configs/importer.dart';
// import 'package:flutter_local_notifications/flutter_local_notifications.dart';
// import 'package:http/http.dart' as http;
import 'package:url_launcher/url_launcher.dart';
import 'package:seeft_mobile/pages/wait_page.dart';

class ManualListPage extends StatefulWidget {
  const ManualListPage({super.key});

  @override
  State<ManualListPage> createState() => _ManualListPageState();
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
    return Scaffold(
      backgroundColor: AppColors.base,
      body: FutureBuilder(
        future: getData(),
        builder: (ctx, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            logger.i("Loading manual list...");
          }
          if (!snapshot.hasData) {
            // 待機画面を表示
            return WaitPage();
          }
          return Container(
            padding: const EdgeInsets.all(32.0),
            child: Column(
              children: <Widget>[
                Flexible(
                  child: ListView.builder(
                    itemCount: manualLength,
                    itemBuilder: (BuildContext context, int index) {
                      return SizedBox(
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
      decoration: BoxDecoration(
        border: Border(bottom: BorderSide(
          width: 1.0,
          color: AppColors.grayLight,
        )),
      ),
      child: ListTile(
        title: Text(
          manuals[index]["task"].toString(),
          style: TextStyle(
            color: AppColors.textBlack, 
            fontSize: AppFontSizes.md
          ),
        ),
        onTap: () async {
          if (await canLaunchUrl(Uri.parse(manuals[index]["url"].toString()))) {
            await launchUrl(Uri.parse((manuals[index]["url"].toString())));
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
