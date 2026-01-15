import 'package:seeft_mobile/configs/importer.dart';
import 'package:url_launcher/url_launcher.dart';

// シフトカードのウィジェット
class ShiftCard extends StatefulWidget {
  final ShiftCardData data;

  const ShiftCard({super.key, required this.data});

  @override
  State<ShiftCard> createState() => _ShiftCardState();
}

class _ShiftCardState extends State<ShiftCard> {
  // モック用：未読管理フラグ（初期値は未読）
  bool isRead = false;

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
    return Stack(
      children: [
        GestureDetector(
          onTap: () {
            setState(() {
              isRead = true;
            });
          },
          child: Card(
            margin: EdgeInsets.zero,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(8.0),
              side: BorderSide(
                color: AppColors.grayLight,
                width: 1.0,
              ),
            ),
            elevation: 0.5,
            shadowColor: null,
            color: AppColors.base,
            child: Padding(
              padding: const EdgeInsets.all(8.0),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                mainAxisAlignment: MainAxisAlignment.start,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // 時刻とマニュアルを開くボタン
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      // 時刻の表示
                      Row(
                        children: [
                          const Icon(
                            Icons.access_time,
                            color: AppColors.textBlack,
                            size: 16,
                          ),
                          const SizedBox(width: 2.0),
                          Text(
                            widget.data.startTime + "〜" + widget.data.endTime,
                            style: const TextStyle(
                              fontSize: AppFontSizes.sm,
                              color: AppColors.textBlack,
                            ),
                          ),
                        ],
                      ),
                      // マニュアルを開くボタン
                      _manualButton(url: widget.data.url),
                    ],
                  ),
                  // タスク名の表示
                  Text(
                    widget.data.taskName.toString(),
                    style: const TextStyle(
                      fontSize: AppFontSizes.md,
                      color: AppColors.textBlack,
                      fontWeight: FontWeight.bold,
                    ),
                    textAlign: TextAlign.left,
                  ),
                  const SizedBox(height: 6.0),
                  const Divider(
                    height: 1,
                    color: AppColors.grayLight,
                  ),
                  // 集合場所とトグル
                  Theme(
                    data: Theme.of(context).copyWith(dividerColor: Colors.transparent),
                    child: ExpansionTile(
                      tilePadding: EdgeInsets.zero,
                      iconColor: AppColors.textBlack,
                      minTileHeight: 0,
                      // 集合場所の表示
                      title: Row(
                        children: [
                          const Icon(
                            Icons.location_on_outlined,
                            color: AppColors.textBlack,
                            size: 16,
                          ),
                          const SizedBox(width: 2.0),
                          Text(
                            widget.data.place,
                            style: const TextStyle(
                              fontSize: AppFontSizes.sm,
                              color: AppColors.textBlack,
                            ),
                          ),
                        ],
                      ),
                      // トグルが展開されたときの内容
                      children: <Widget>[
                        _buildSection(
                          title: '【集合場所】',
                          content: [widget.data.place],
                        ),
                        _buildManualSection(
                          title: '【マニュアル】',
                          url: widget.data.url != '' ? widget.data.url : null,
                        ),
                        _buildSection(
                          title: '【困った時は】',
                          content: [
                            '以下の順に対応してください。',
                            '1. マニュアルを確認してください。',
                            '2. 近くの人や近くの先輩に相談してください。',
                            '3. 「緊急事対応」ページから本部に連絡してください。',
                          ],
                        ),
                        _buildSection(
                          title: '【担当者の一覧】',
                          content: [
                            for (var member in widget.data.shiftMembers)
                              '${member.s_time}〜${member.e_time}\n' +
                              member.members.map((m) => '(${m.bureau}${m.grade}) ${m.name}').join(', '),
                          ],
                        ),
                        _buildSection(
                          title: '【前の時間の担当者の一覧】',
                          content: [
                            if (widget.data.beforeMembers.members.isNotEmpty)
                              '${widget.data.beforeMembers.s_time}〜${widget.data.beforeMembers.e_time}\n' +
                              widget.data.beforeMembers.members.map((m) => '(${m.bureau}${m.grade}) ${m.name}').join(', ')
                            else
                              '前の時間の担当者はいません',
                          ],
                        ),
                        _buildSection(
                          title: '【次の時間の担当者の一覧】',
                          content: [
                            if (widget.data.afterMembers.members.isNotEmpty)
                              '${widget.data.afterMembers.s_time}〜${widget.data.afterMembers.e_time}\n' +
                              widget.data.afterMembers.members.map((m) => '(${m.bureau}${m.grade}) ${m.name}').join(', ')
                            else
                              '次の時間の担当者はいません',
                          ],
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
        // 「New!」バッジ（未読の場合のみ表示）
        if (!isRead)
          Positioned(
            top: 4.0,
            right: 4.0,
            child: Container(
              padding: const EdgeInsets.symmetric(
                horizontal: 6.0,
                vertical: 2.0,
              ),
              decoration: BoxDecoration(
                color: Colors.red,
                borderRadius: BorderRadius.circular(8.0),
              ),
              child: const Text(
                'New!',
                style: TextStyle(
                  color: Colors.white,
                  fontSize: 10.0,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
          ),
      ],
    );
  }
  // マニュアルボタン
  Widget _manualButton({
    required String url,
  }) {
    final bool isDisabled = url.isEmpty;  // マニュアルがない場合は無効化する
    return ElevatedButton.icon(
      onPressed: isDisabled ? null : () => _launchManualUrl(url), // マニュアルがあればタップで開く
      icon: Icon(
        Icons.quiz_outlined,
        size: 16,
        color: isDisabled ? AppColors.grayDark : AppColors.textWhite,
      ),
      label: Text(
        isDisabled ? 'マニュアルなし' : 'マニュアル',
        style: TextStyle(
          fontSize: AppFontSizes.sm,
          color: isDisabled ? AppColors.grayDark : AppColors.textWhite,
        ),
      ),
      style: ElevatedButton.styleFrom(
        backgroundColor: isDisabled ? AppColors.grayLight : AppColors.link,
        padding: const EdgeInsets.symmetric(
          vertical: 4.0,
          horizontal: 8.0
        ),
        elevation: 0,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(100.0),
        ),
        // 横幅を広げるための設定
        // minimumSize: Size(double.infinity, 40),
      ),
    );
  }
  
  // 各セクションのUIを生成するヘルパーメソッド
  Widget _buildSection({
    required String title, 
    required List<String> content
  }) {
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
  // マニュアルセクションのUIを生成するヘルパーメソッド
  Widget _buildManualSection({
    required String title, 
    String? url
  }) {
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
            GestureDetector(
              onTap: () => url != null ? _launchManualUrl(url) : null,  // マニュアルが存在すればタップでマニュアルを開く
              child: Text(
                url ?? 'マニュアルがありません',
                style: TextStyle(
                  fontSize: AppFontSizes.xs,
                  color: url != null ? AppColors.link : AppColors.textBlack,
                  height: 1.5 // 行間を調整
                )
              ),
            )
          ],
        ),
      ),
    );
  }
}
