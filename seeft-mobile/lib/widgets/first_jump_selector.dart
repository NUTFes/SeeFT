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
    final _isUserID = await store.isUserID();
    final _userID = await store.getUserID();
    logger.w(_userID);
    logger.w(_isUserID);
    return _isUserID;
  }

  @override
  Widget build(BuildContext context) {
    logger.i('navigated Splash.');

    return FutureBuilder(
      future: getPrefRead(),
      builder: (BuildContext context, AsyncSnapshot<bool> snapshot) {
        // var hasData = snapshot;
        logger.i('===============================');
        logger.i(snapshot.connectionState);

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

        var app = new MaterialApp(
          title: constant.appName,
          theme: ThemeData(
            textTheme:
                GoogleFonts.mPlus1pTextTheme(Theme.of(context).textTheme),
            primarySwatch: Colors.teal,
            secondaryHeaderColor: Colors.teal,
            focusColor: Colors.teal,
            backgroundColor: Colors.white,
            cardColor: Colors.white,
            dialogBackgroundColor: Colors.white,
            pageTransitionsTheme: const PageTransitionsTheme(
              builders: <TargetPlatform, PageTransitionsBuilder>{
                TargetPlatform.android: CupertinoPageTransitionsBuilder(),
                TargetPlatform.iOS: CupertinoPageTransitionsBuilder()
              },
            ),
          ),
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
          },
        );

        return app;
      },
    );
  }
}
