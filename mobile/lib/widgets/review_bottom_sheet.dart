import 'package:seeft_mobile/configs/importer.dart';
import 'package:seeft_mobile/widgets/custom_text_field.dart';
import 'package:seeft_mobile/widgets/custom_elevated_button.dart';
import 'package:seeft_mobile/widgets/custom_elevated_button_outlined.dart';
import 'package:seeft_mobile/widgets/custom_error_snack_bar.dart';
import 'package:seeft_mobile/widgets/custom_snack_bar.dart';

// DB の reviews.comment は VARCHAR(255)。PostgreSQL はコードポイント数で数える。
const int reviewCommentMaxLength = 255;

// コードポイント数で切り詰める。
// TextField の maxLength は Web だと truncateAfterCompositionEnds が既定で、
// IME 変換中の文字は対象外になる。さらに数え方が書記素クラスタなので、絵文字や
// 結合文字を含むと maxLength を満たしていてもコードポイント数は超えうる。
// 送信直前にここで揃えて、INSERT が落ちてコメントごと失われるのを防ぐ。
String clampToCodePoints(String value, int maxCodePoints) {
  final runes = value.runes.toList();
  if (runes.length <= maxCodePoints) {
    return value;
  }
  return String.fromCharCodes(runes.take(maxCodePoints));
}

// 星の行を特定するためのキー
const Key staffingRatingRowKey = ValueKey('staffing_rating_row');
const Key manualRatingRowKey = ValueKey('manual_rating_row');

// レビューを入力するボトムシートのウィジェット
class ReviewBottomSheet {
  // ボトムシートを表示するメソッド
  static void show(
    BuildContext context,
    String taskName,
    int userID
  ) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true, // 高さを自由に調整
      backgroundColor: AppColors.base,
      barrierColor: Colors.black.withValues(alpha: 0.2), // ← デフォルト0.54 → 薄めに調整
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(16)),
      ),
      builder: (context) {
        return Padding(
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(context).viewInsets.bottom, // キーボードの高さを考慮
          ),
          child: ReviewForm(
            taskName: taskName,
            userID: userID,
          ),
        );
      },
    );
  }
}

class ReviewForm extends StatefulWidget {
  final String taskName;
  final int userID;

  const ReviewForm({super.key, required this.taskName, required this.userID});

  @override
  State<ReviewForm> createState() => _ReviewFormState();
}

class _ReviewFormState extends State<ReviewForm> {
  // 0 は未選択。初期値を入れてしまうと、無操作で送信した人と「普通」と答えた人が
  // 同じ値になり集計で区別できなくなるため、選ぶまで送信させない
  int staffingRating = 0;     // シフトの人数評価
  int manualRating = 0;       // マニュアルの評価
  bool _isSubmitting = false; // 送信中フラグ
  bool _isFailed = false;     // 送信失敗フラグ
  final TextEditingController _controller = TextEditingController();

  // 表示と送信可否で条件がずれないよう、判定はここだけに置く
  bool get _isUnrated => staffingRating == 0 || manualRating == 0;

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  // 星による5段階評価の行を作成するウィジェット。
  // key は、どちらの設問の星かをテストから特定するために使う
  Widget buildStarRow(int currentValue, ValueChanged<int> onChanged, {Key? key}) {
    return Row(
      key: key,
      mainAxisAlignment: MainAxisAlignment.center,
      children: List.generate(5, (index) {
        return IconButton(
          icon: Icon(
            index < currentValue ? Icons.star : Icons.star_border,
            color: index < currentValue ? Colors.amber : Colors.grey,
            size: 32,
          ),
          onPressed: () => setState(() => onChanged(index + 1)),
        );
      }),
    );
  }
  
  // レビューを送信する非同期関数
  Future<bool> _sendReview(
    BuildContext context,
    int userID,
    String taskName,
    int staffingRating,
    int manualRating,
    String comment
  ) async {
    try {
      logger.i('レビューを送信中...'
        '\nUser ID: $userID'
        '\nTask Name: $taskName'
        '\nStaffing Rating: $staffingRating'
        '\nManual Rating: $manualRating'
        '\nComments: ${_controller.text}');

      // API呼び出し
      await api.postReview(userID, taskName, staffingRating, manualRating, comment);

      logger.i('レビューの送信に成功しました');
      if (context.mounted) showCustomSnackBar(context, "レビューを送信しました");
      return true;
    } catch (e) {
      logger.e('レビューの送信に失敗しました: $e');
      if (context.mounted) showCustomErrorSnackBar(context, "レビューの送信に失敗しました");
      return false;
    }
  }
  
  // 送信ボタンを押したときの処理
  void _onSubmit() async {
    setState(() {
      _isSubmitting = true; // 送信中フラグを立てる
    });
    // レビューを送信する
    final isSuccess = await _sendReview(
      context,
      widget.userID,
      widget.taskName,
      staffingRating,
      manualRating,
      clampToCodePoints(_controller.text, reviewCommentMaxLength),
    );
    if (!mounted) return;
    if(isSuccess){  // 送信成功時
      // Hiveにタスク名を保存
      reviewedTaskNameBox.put(widget.taskName, true);
      Navigator.of(context).pop(); // ボトムシートを閉じる
    }else{
      setState(() {
        _isSubmitting = false;
        _isFailed = true;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(32.0),
      child: Column(
        mainAxisSize: MainAxisSize.min, // ボトムシートの高さを内容に合わせる
        crossAxisAlignment: CrossAxisAlignment.start,
        spacing: 8,
        children: [
          Text(
            "シフトのレビュー: ${widget.taskName}",
            style: TextStyle(
              color: AppColors.textBlack,
              fontSize: AppFontSizes.md,
              fontWeight: FontWeight.bold,
            ),
          ),
          const Divider(
            thickness: 1, // 区切り線の太さ
            color: AppColors.grayLight,
          ),
          const Text("シフトの人数は適切でしたか？", style: TextStyle(color: AppColors.textBlack, fontSize: AppFontSizes.md)),
          buildStarRow(
            staffingRating,
            (value) => staffingRating = value,
            key: staffingRatingRowKey,
          ),
          const SizedBox(height: 8),
          const Text("マニュアルは分かりやすかったですか？", style: TextStyle(color: AppColors.textBlack, fontSize: AppFontSizes.md)),
          buildStarRow(
            manualRating,
            (value) => manualRating = value,
            key: manualRatingRowKey,
          ),
          const SizedBox(height: 8),
          const Text("他にもあれば教えてください。", style: TextStyle(color: AppColors.textBlack, fontSize: AppFontSizes.md)),
          CustomTextField(
            controller: _controller,
            hintText: "例：マニュアルが分かりやすくて良かった",
            // DB の comment は VARCHAR(255)。超えると INSERT が落ちて入力が失われるため、
            // 入力側で打ち切る
            maxLength: reviewCommentMaxLength,
          ),
          const SizedBox(height: 8),
          // Visibility で隠しても Column の spacing は入るため、
          // 条件付きで子ごと外して余白が残らないようにする
          if (_isUnrated)
            const Text(
              "星を選ぶと送信できます。",
              style: TextStyle(
                color: AppColors.grayDark,
                fontSize: AppFontSizes.sm,
              ),
            ),
          if (_isFailed)
            const Text(
              "送信に失敗しました。もう一度お試しください。",
              style: TextStyle(
                color: AppColors.error,
                fontSize: AppFontSizes.sm,
              ),
            ),
          Row(
            spacing: 8.0,
            children: [
              // スキップするボタン
              Expanded(
                child: CustomElevatedButtonOutlined(
                  onPressed: () {
                    Navigator.of(context).pop(); // ボトムシートを閉じる
                  },
                  label: "スキップ",
                  isExpanded: true,
                ),
              ),
              // 送信ボタン
              Expanded(
                child: CustomElevatedButton(
                  onPressed: _onSubmit,
                  label: "送信",
                  isDisabled: _isSubmitting || _isUnrated,
                  isExpanded: true,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
