import 'package:seeft_mobile/configs/importer.dart';
import 'package:url_launcher/url_launcher.dart';

// シフトカードのウィジェット
class ShiftCard extends StatelessWidget {
  final ShiftCardData data;

  const ShiftCard({super.key, required this.data});
  
  // リンクを開くための非同期メソッドを定義
  Future<void> _launchManualUrl(String url) async {
    // 開きたいURLをUriオブジェクトに変換
    final Uri uri = Uri.parse(url);

    // launchUrlを実行
    if (!await launchUrl(uri, mode: LaunchMode.externalApplication)) {  // ここで外部アプリで開くように指定
      // URLが開けなかった場合のエラー処理
      throw Exception('Could not launch $uri');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Card(
      margin:EdgeInsets.zero,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(8.0), // 角丸を設定
        side: BorderSide(
          color: AppColors.grayLight, // 枠線の色
          width: 1.0, // 枠線の太さ
        ),
      ),
      elevation: 0.5,
      shadowColor: null,
      color: AppColors.base, // 背景色を設定
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          // 時刻とマニュアルを開くボタン
          Padding(
            padding: const EdgeInsets.fromLTRB(8.0, 8.0, 8.0, 3.0),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                // 時刻の表示
                Row(
                  children: [
                    const Icon(
                      Icons.access_time,
                      color: AppColors.textBlack,
                      size: 16
                    ),
                    const SizedBox(width: 2.0),
                    Text(
                      data.startTime + "〜" + data.endTime,
                      style: const TextStyle(
                        fontSize: AppFontSizes.sm,
                        color: AppColors.textBlack
                      ),
                    ),
                  ],
                ),
                // マニュアルを開くボタン
                ElevatedButton(
                  onPressed: () {
                    // マニュアルを開くボタンの処理
                    if (data.url.isNotEmpty) {
                      _launchManualUrl(data.url);
                    } else {
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(content: Text('マニュアルのURLが設定されていません')), // マニュアルのURLが空の場合のエラーメッセージ
                      );
                    }
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: AppColors.sub,
                    foregroundColor: AppColors.textWhite,
                    shape: const StadiumBorder(), // 角が丸い形状
                    elevation: 0, // 影を消す
                    padding: const EdgeInsets.symmetric(horizontal: 8.0, vertical: 4.0),
                  ),
                  child: const Text(
                    'マニュアルを開く',
                    style: TextStyle(
                      fontSize: AppFontSizes.xs, 
                      color: AppColors.textWhite
                    ),
                  ),
                ),
              ],
            ),
          ),
          
          Padding(
            padding: const EdgeInsets.fromLTRB(8.0, 3.0, 8.0, 8.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // タスク名の表示
                Text(
                  data.taskName.toString(),
                  style: const TextStyle(
                    fontSize: AppFontSizes.md,
                    color: AppColors.textBlack, 
                    fontWeight: FontWeight.bold
                  ),
                ),
                const SizedBox(height: 6.0),
                const Divider(
                  height: 1,
                  color: AppColors.grayLight, // 区切り線の色
                ),
                // 集合場所とトグル
                Theme(
                  data: Theme.of(context).copyWith(dividerColor: Colors.transparent), // 区切り線の色を透明に
                  child: ExpansionTile(
                    tilePadding: EdgeInsets.zero,
                    iconColor: AppColors.textBlack, // アイコンの色を変更
                    minTileHeight: 0, // タイルの高さを最小に
                    // 集合場所の表示
                    title: Row(
                      children: [
                        const Icon(
                          Icons.location_on_outlined, 
                          color: AppColors.textBlack,
                          size: 16
                        ),
                        const SizedBox(width: 2.0),
                        Text(
                          data.place,
                          style: const TextStyle(
                            fontSize: AppFontSizes.sm,
                            color: AppColors.textBlack
                          ),
                        ),
                      ],
                    ),
                    // トグルが展開されたときの内容
                    children: <Widget>[
                      _buildSection(
                        title: '【集合場所】',
                        content: [data.place]
                      ),
                      _buildSection(
                        title: '【困った時は】',
                        content: [
                          '1. 近くの先輩に聞く',
                          '2. それでも分からなかったら本部に連絡',
                        ],
                      ),
                      _buildSection(
                        title: '【担当者の一覧】',
                        content: [
                          for (var member in data.shiftMembers)
                            '${member.s_time}〜${member.e_time}\n' +
                            member.members.map((m) => '(${m.bureau}${m.grade}) ${m.name}').join(', '),
                        ],
                      ),
                      _buildSection(
                        title: '【前の時間の担当者の一覧】',
                        content: [
                          if (data.beforeMembers.members.isNotEmpty)
                            '${data.beforeMembers.s_time}〜${data.beforeMembers.e_time}\n' +
                            data.beforeMembers.members.map((m) => '(${m.bureau}${m.grade}) ${m.name}').join(', ')
                          else
                            '前の時間の担当者はいません',
                        ],
                      ),
                      _buildSection(
                        title: '【次の時間の担当者の一覧】',
                        content: [
                          if (data.afterMembers.members.isNotEmpty)
                            '${data.afterMembers.s_time}〜${data.afterMembers.e_time}\n' +
                            data.afterMembers.members.map((m) => '(${m.bureau}${m.grade}) ${m.name}').join(', ')
                          else
                            '次の時間の担当者はいません',
                        ],
                      ),
                    ],
                  ),
                )
              ],
            ),
          ),
        ],
      ),
    );
  }
  // 各セクションのUIを生成するヘルパーメソッド
  Widget _buildSection({required String title, required List<String> content}) {
    return Padding(
      padding: const EdgeInsets.only(top: 8.0, bottom: 8.0),
      child: Align(
        alignment: Alignment.centerLeft,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              title,
              style: const TextStyle(
                fontSize: AppFontSizes.xs,
                color: AppColors.textBlack,
                fontWeight: FontWeight.bold
              ),
            ),
            const SizedBox(height: 4.0),
            ...content.map((line) => Text(
              line, 
              style: const TextStyle(
                fontSize: AppFontSizes.xs, 
                color: AppColors.textBlack,
                height: 1.5 // 行間を調整
              )
            )),
          ],
        ),
      ),
    );
  }
}
