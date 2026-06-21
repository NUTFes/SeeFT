import 'package:seeft_mobile/configs/importer.dart';
import 'package:seeft_mobile/pages/sign_in_page.dart';
import 'package:seeft_mobile/widgets/layout.dart';

class FirstJumpSelector extends StatefulWidget {
  const FirstJumpSelector({super.key});

  @override
  State<FirstJumpSelector> createState() => new _FirstJumpSelectorState();
}

class _FirstJumpSelectorState extends State<FirstJumpSelector> {
  @override
  void initState() {
    super.initState();
    Timer(Duration(seconds: 3), () => logger.w('progress duration:'));
  }

  Future<bool> getPrefRead() async {
    try {
      final isUserID = await store.isUserID();
      final userID = await store.getUserID();
      logger.w(userID);
      logger.w(isUserID);
      return isUserID;
    } catch (e) {
      logger.e(e);
      return false;
    }
  }

  @override
  Widget build(BuildContext context) {
    logger.i('navigated Splash.');
    final materialTheme = MaterialTheme(Typography().black); // Typography() は変更が必要な可能性があります

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
            theme: materialTheme.light(), // ライトモードのテーマ
              darkTheme: materialTheme.dark(), // ダークモードのテーマ
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
            theme: materialTheme.light(), // ライトモードのテーマ
              darkTheme: materialTheme.dark(), // ダークモードのテーマ
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

