import 'package:seeft_mobile/configs/importer.dart';
import 'package:seeft_mobile/widgets/custom_error_snack_bar.dart';
import 'package:url_launcher/url_launcher.dart';

/// マニュアル一覧の1行。
///
/// ドキュメント版(url)とスライド版(manualUrl)は別々に存在しうるので、
/// それぞれが有るときだけ押せるようにする。無いほうを押せるままにすると、
/// 押しても何も起きない行ができる(issue #524)。
class ManualListItem extends StatelessWidget {
  const ManualListItem({
    super.key,
    required this.taskName,
    required this.url,
    required this.manualUrl,
  });

  /// 行に表示するタスク名
  final String taskName;

  /// ドキュメント版マニュアルのURL。空なら行をタップできない
  final String url;

  /// スライド版マニュアルのURL。空ならアイコンを出さない
  final String manualUrl;

  bool get _hasDocument => url.isNotEmpty;
  bool get _hasSlide => manualUrl.isNotEmpty;

  /// URLを別のタブで開き、開けなければエラーを出す。
  ///
  /// canLaunchUrlで事前にガードしない。URLが空のときcanLaunchUrlがfalseを返し、
  /// エラーも出ないまま終わっていたのが issue #524 の症状だった。
  /// shift_card.dart の _open と同じく、launchUrlの戻り値と例外で判定する。
  Future<void> _open(
    BuildContext context,
    String target,
    String errorMessage,
  ) async {
    var launched = false;
    try {
      launched = await launchUrl(
        Uri.parse(target),
        // スライド版は認証付き配信のため埋め込み不可。ドキュメント版も別タブで開く
        mode: LaunchMode.externalApplication,
      );
    } catch (_) {
      launched = false;
    }
    if (!launched && context.mounted) {
      showCustomErrorSnackBar(context, errorMessage);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: const BoxDecoration(
        border: Border(
          bottom: BorderSide(
            width: 1.0,
            color: AppColors.grayLight,
          ),
        ),
      ),
      child: ListTile(
        title: Text(
          taskName,
          style: TextStyle(
            // ドキュメント版が無い行は開けない。押せないことを色でも示す
            color: _hasDocument ? AppColors.textBlack : AppColors.grayDark,
            fontSize: AppFontSizes.md,
          ),
        ),
        trailing: _hasSlide
            ? IconButton(
                icon: const Icon(Icons.slideshow),
                onPressed: () => _open(
                  context,
                  manualUrl,
                  'マニュアル（スライド版）を開けませんでした',
                ),
              )
            : null,
        // onTapにnullを渡すとListTileが押下エフェクトを出さなくなる。
        // 「押せそうに見えて何も起きない」状態をなくすため、無効化はここで行う
        onTap:
            _hasDocument ? () => _open(context, url, 'マニュアルを開けませんでした') : null,
      ),
    );
  }
}
