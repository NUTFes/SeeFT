import 'dart:developer';

import 'package:seeft_mobile/configs/importer.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:http/http.dart' as http;
import 'package:url_launcher/url_launcher.dart';
import 'package:flutter/material.dart';

class ContactPage extends StatefulWidget {
  @override
  _ContactPageState createState() => _ContactPageState();
}

class _ContactPageState extends State<ContactPage> {
// notification関連をinitStateに書き出さなきゃいけないので書いてたけどutilとかに書いてもいいかもね

//  FlutterLocalNotificationsPlugin flutterLocalNotificationsPlugin;
//  NotificationDetails platformChannelSpecifics;
  @override
  Widget build(BuildContext context) {
    final Size size = MediaQuery.of(context).size;
    return Scaffold(
      appBar: AppBar(
        title: const Text('本部連絡先'),
        actions: <Widget>[],
        // debug
      ),
      drawer: drawer.applicationDrawer(context),
      body:Container(
        alignment: Alignment.center,
        padding: const EdgeInsets.all(40.0),
        child: SingleChildScrollView(  
          child: Column(
                children: <Widget>[
                  Text(
                    "当日の対応などでマニュアルでは\n対応できないことがあったときは\n本部へ連絡して指示を仰いでください。\n",
                    textAlign: TextAlign.center,
                    ),
                  Text(
                    "本部への緊急連絡先は以下になります。\n",
                    textAlign: TextAlign.center,
                    ),
                  Text(
                    "副委員長　小林諒大",
                    textAlign: TextAlign.center,
                  ),
                  TextButton(
                    onPressed: _openPhoneApp,
                    child: Text(
                      "080-9703-4992\n",
                      textAlign: TextAlign.center,
                    ),
                  ),
                  Text(
                    "副委員長が通話中でつながらなかった場合は以下にご連絡ください\n",
                    textAlign: TextAlign.center,
                    ),
                  Text(
                    "委員長　中井温人",
                    textAlign: TextAlign.center,
                  ),
                  TextButton(
                    onPressed: _openPhoneAppNakai,
                    child: Text(
                      "080-8918-4164\n",
                      textAlign: TextAlign.center,
                    ),
                  ),
                  SizedBox(
                    child: Container(
                    height: 100,
                    width: 100,
                    ),
                  ),
                  Text(
                    "SeeFTの使用方法やバグなどの\nお問い合わせは以下へお願い致します。\n",
                    textAlign: TextAlign.center,
                  ),
                  Text(
                    "情報局長　藤崎竜成",
                    textAlign: TextAlign.center,
                  ),
                  TextButton(
                    onPressed: _openMailApp,
                    child: Text(
                      "r.fujisaki.nutfes@gmail.com\n",
                      textAlign: TextAlign.center,
                    ),
                  ),
                ],
          ),
        ),
      )
    );
  }
}

void _openPhoneApp() {
    const tel = '08097034992';
    _launchURL(
      'tel:' + tel,
    );
  }
void _openPhoneAppNakai() {
    const tel = '08089184164';
    _launchURL(
      'tel:' + tel,
    );
  }

  _openMailApp() async {
    final title = Uri.encodeComponent('タイトル');
    final body = Uri.encodeComponent('本文');
    const mailAddress = 'r.fujisaki.nutfes@gmail.com';

    return _launchURL(
      'mailto:$mailAddress?subject=$title&body=$body',
    );
  }

  Future<void> _launchURL(String url) async {
    if (await canLaunch(url)) {
      await launch(url);
    } else {
      final Error error = ArgumentError('Could not launch $url');
      throw error;
    }
  }
