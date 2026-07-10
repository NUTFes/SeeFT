import 'package:seeft_mobile/configs/importer.dart';
import 'package:seeft_mobile/pages/sign_in_page.dart';
import 'package:seeft_mobile/widgets/layout.dart';

class FirstJumpSelector extends StatelessWidget {
  const FirstJumpSelector({super.key});

  @override
  Widget build(BuildContext context) {
    final materialTheme = MaterialTheme(Typography().black); // Typography() は変更が必要な可能性があります

    // MaterialApp はここで1回だけ作る。ローディング中／未ログイン／ログイン済みの
    // 画面切り替えは home の中身（_StartupScreen）だけで行い、
    // MaterialApp 自体の設定（home / routes の形）は変えない。
    return MaterialApp(
      title: constant.appName,
      theme: materialTheme.light(),
      darkTheme: materialTheme.dark(),
      themeMode: ThemeMode.light,
      home: const _StartupScreen(),
      routes: {
        '/signin': (context) => SignInPage(),
        '/layout': (context) => Layout(),
      },
    );
  }
}

// アプリ起動直後に表示する画面。
// ログイン済みかどうかを確認している間はローディング表示、
// 確認が終わったらログイン画面 or マイシフト画面を表示する。
class _StartupScreen extends StatefulWidget {
  const _StartupScreen();

  @override
  State<_StartupScreen> createState() => _StartupScreenState();
}

class _StartupScreenState extends State<_StartupScreen> {
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

    return FutureBuilder<bool>(
      future: getPrefRead(),
      builder: (BuildContext context, AsyncSnapshot<bool> snapshot) {
        logger.i('===============================');
        logger.i(snapshot.connectionState);

        if (snapshot.hasError) {
          return Scaffold(
            appBar: AppBar(
              title: const Text('Error Message'),
            ),
            body: Column(children: [Text(snapshot.error.toString())]),
          );
        }

        if (!snapshot.hasData) {
          // 読み込み中（waiting）のデフォルト表示
          return Scaffold(
            body: Center(
              child: CircularProgressIndicator(color: AppColors.main),
            ),
          );
        }

        final isUserID = snapshot.data!;
        logger.i(snapshot.connectionState);
        if (isUserID) {
          logger.i('select Layout.');
          return Layout();
        } else {
          logger.i('select MainPage.');
          return SignInPage();
        }
      },
    );
  }
}

