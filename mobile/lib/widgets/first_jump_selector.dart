import 'package:seeft_mobile/configs/importer.dart';
import 'package:seeft_mobile/pages/sign_in_page.dart';
import 'package:seeft_mobile/widgets/layout.dart';

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
        logger.i('===============================');
        logger.i(snapshot.connectionState);
        var app;
        if (snapshot.hasData) {
          var isUserID = snapshot.data;
          logger.i(snapshot.connectionState);
          var homeWidget;
          if (isUserID!) {
            logger.i('select Layout.');
            homeWidget = '/layout';
          } else {
            logger.i('select MainPage.');
            homeWidget = '/signin';
          }

          app = new MaterialApp(
            title: constant.appName,
            theme: _materialTheme.light(), // ライトモードのテーマ
              darkTheme: _materialTheme.dark(), // ダークモードのテーマ
              themeMode: ThemeMode.light, // ライトモードを適用
            initialRoute: homeWidget,
            routes: {
              '/signin': (context) => SignInPage(),
              '/layout': (context) => Layout(),
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

