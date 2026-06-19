// その他画面
import 'package:seeft_mobile/configs/importer.dart';
import 'package:url_launcher/url_launcher.dart';

class EtcPage extends StatelessWidget {

  EtcPage({
    super.key,
  });
  
  // 環境変数からSeeFTの操作説明と全体シフトのURLを取得
  final String seeftInstructionsUrl = constant.seeftInstructionsUrl;
  final String wholeShiftUrl = constant.wholeShiftUrl;

  @override
  Widget build(BuildContext context) {
    // 項目リスト
    final List<Map<String, dynamic>> _listItems = [
      {
        'text': '操作説明', 
        'explanation': 'SeeFTの操作説明を開きます',
        'icon': Icons.help_outline,
        'onTap': () async {
          logger.i("操作説明が選択されました");
          // SeeFTの操作説明のスライドを開く
          final url = seeftInstructionsUrl;
          if (await canLaunch(url)) {
            await launch(url);
          } else {
            final Error error = ArgumentError('Could not launch $url');
            throw error;
          }
        },
      },
      {
        'text': '全体シフト', 
        'explanation': '全局員のシフトを確認できます',
        'icon': Icons.today,
        'onTap': () async {
          logger.i("全体シフトが選択されました");
          // 全体シフトのURLを開く
          final url = wholeShiftUrl;
          if (await canLaunch(url)) {
            await launch(url);
          } else {
            final Error error = ArgumentError('Could not launch $url');
            throw error;
          }
        },
      },
      {
        'text': 'ログアウト', 
        'explanation': 'ログアウトします',
        'icon': Icons.logout,
        'onTap': () {
          logger.i("ログアウトが選択されました");
          // ログアウトする
          Navigator.pushNamedAndRemoveUntil(
              context, '/signin', (Route<dynamic> route) => false);
        },
      },
    ];
    
    return Scaffold(
      backgroundColor: AppColors.base,
      body: SingleChildScrollView(
        child: Container(
          padding: const EdgeInsets.all(32.0),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.center,
            // spacing: 24.0,  // 子要素間のスペース
            children: [
              // 問題の種類を選択するリスト
              ListView.separated(
                itemCount: _listItems.length,
                shrinkWrap: true, // リストの高さを自動調整
                physics: NeverScrollableScrollPhysics(), // スクロールを無効化
                itemBuilder: (context, index) {
                  return ListTile(
                    leading: Icon(
                      _listItems[index]['icon'],
                      color: AppColors.main,
                    ),
                    title: Text(
                      _listItems[index]['text'],
                      style: TextStyle(
                        fontSize: AppFontSizes.md,
                        fontWeight: FontWeight.bold,
                        color: AppColors.textBlack,
                      ),
                    ),
                    subtitle: Text(
                      _listItems[index]['explanation'],
                      style: TextStyle(
                        fontSize: AppFontSizes.sm,
                        color: AppColors.grayDark,
                      ),
                    ),
                    trailing: Icon(
                      Icons.arrow_forward_ios,
                      color: AppColors.grayDark,
                      size: 16.0,
                    ),
                    onTap: () {
                      // 問題の種類が選択されたときの処理
                      _listItems[index]['onTap']();
                    },
                  );
                },
                separatorBuilder: (context, index) {
                  return Divider(
                    color: AppColors.grayLight,
                    thickness: 1.0,
                  );
                },
              ),
            ],
          ),
        ),
      ),
    );
  }
}
