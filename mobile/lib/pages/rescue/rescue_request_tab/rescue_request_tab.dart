import 'package:seeft_mobile/configs/importer.dart';
import 'package:seeft_mobile/pages/rescue/rescue_request_tab/tab_pages/home.dart';

// 「レスキューを送信する」タブ
class RescueRequestTab extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    // ネストされたNavigatorを使用して、タブ内で独立したナビゲーションを可能にする
    return Navigator(
      onGenerateRoute: (settings) {
        return MaterialPageRoute(
          builder: (_) => RescueRequestTabHome(), // ホーム画面を表示
        );
      },
    );
  }
}
