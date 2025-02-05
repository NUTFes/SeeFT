import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:seeft_mobile/configs/importer.dart';
import 'package:seeft_mobile/pages/sign_in_page.dart';
import 'package:seeft_mobile/pages/my_shift_page.dart';
import 'package:seeft_mobile/pages/all_shift_page.dart';
import 'package:seeft_mobile/pages/manual_list_page.dart';
import 'package:seeft_mobile/pages/schedule_page.dart';
import 'package:seeft_mobile/pages/contact_page.dart';
import 'package:seeft_mobile/pages/wait_page.dart';
import 'package:seeft_mobile/pages/users_page.dart';
import 'package:seeft_mobile/theme/theme.dart';

class FirstJumpSelector extends StatefulWidget {
  @override
  _FirstJumpSelectorState createState() => new _FirstJumpSelectorState();
}

class _FirstJumpSelectorState extends State<FirstJumpSelector> {
  @override
  void initState() {
    super.initState();
    Timer(Duration(seconds: 3), () => logger.w('progress duration:'));
  }

  Future<bool> getPrefRead() async {
    try {
      final _isUserID = await store.isUserID();
      final _userID = await store.getUserID();
      logger.w(_userID);
      logger.w(_isUserID);
      return _isUserID;
    } catch (e) {
      print(e);
      return false;
    }
  }

  @override
  Widget build(BuildContext context) {
    logger.i('navigated Splash.');
    final _materialTheme = MaterialTheme(Typography().black); // Typography() は変更が必要な可能性があります

    return FutureBuilder(
      future: getPrefRead(),
      builder: (BuildContext context, AsyncSnapshot<bool> snapshot) {
        // var hasData = snapshot;
        logger.i('===============================');
        logger.i(snapshot.connectionState);
        var app;
        if (snapshot.hasData) {
          var isUserID = snapshot.data;
          logger.i(snapshot.connectionState);
          var homeWidget;
          if (isUserID!) {
            // logger.i('select SignInPage.');
            // homeWidget = '/my_shift_page';
            logger.i('select WaitPage.');
            homeWidget = '/my_shift_page';
          } else {
            logger.i('select MainPage.');
            homeWidget = '/signin';
          }

          app = new MaterialApp(
            title: constant.appName,
            theme: _materialTheme.light(), // ライトモードのテーマ
              darkTheme: _materialTheme.dark(), // ダークモードのテーマ
              themeMode: ThemeMode.light, // ライトモードを適用
            //home: homeWidget,
            initialRoute: homeWidget,
            routes: {
              '/signin': (context) => SignInPage(),
              '/my_shift_page': (context) => MyShiftPage(),
              '/all_shift_page': (context) => AllShiftPage(),
              '/manual_list_page': (context) => ManualListPage(),
              '/schedule_page': (context) => SchedulePage(),
              '/contact_page': (context) => ContactPage(),
              '/wait_page': (context) => WaitPage(),
              '/users_page': (context) => UsersPage(),
            },
          );
        } else if (snapshot.hasError) {
          app = new MaterialApp(
            title: constant.appName,
            theme: _materialTheme.light(), // ライトモードのテーマ
              darkTheme: _materialTheme.dark(), // ダークモードのテーマ
              themeMode: ThemeMode.light, // ライトモードを適用
            home: Scaffold(
              appBar: AppBar(
                title: const Text('Error Message'),
              ),
              body: Column(children: [Text(snapshot.error.toString())]),
            ),
          );
        }
        return app;
      },
    );
  }
}
