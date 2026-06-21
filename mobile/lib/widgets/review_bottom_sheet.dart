import 'package:seeft_mobile/configs/importer.dart';
import 'package:seeft_mobile/widgets/custom_text_field.dart';
import 'package:seeft_mobile/widgets/custom_elevated_button.dart';
import 'package:seeft_mobile/widgets/custom_elevated_button_outlined.dart';
import 'package:seeft_mobile/widgets/custom_error_snack_bar.dart';
import 'package:seeft_mobile/widgets/custom_snack_bar.dart';

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
  int staffingRating = 3;     // シフトの人数評価
  int manualRating = 3;       // マニュアルの評価
  bool _isSubmitting = false; // 送信中フラグ
  bool _isFailed = false;     // 送信失敗フラグ
  final TextEditingController _controller = TextEditingController();

  // 星による5段階評価の行を作成するウィジェット
  Widget buildStarRow(int currentValue, ValueChanged<int> onChanged) {
    return Row(
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
      showCustomSnackBar(context, "レビューを送信しました");
      return true;
    } catch (e) {
      logger.e('レビューの送信に失敗しました: $e');
      showCustomErrorSnackBar(context, "レビューの送信に失敗しました");
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
      _controller.text,
    );
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
          buildStarRow(staffingRating, (value) => staffingRating = value),
          const SizedBox(height: 8),
          const Text("マニュアルは分かりやすかったですか？", style: TextStyle(color: AppColors.textBlack, fontSize: AppFontSizes.md)),
          buildStarRow(manualRating, (value) => manualRating = value),
          const SizedBox(height: 8),
          const Text("他にもあれば教えてください。", style: TextStyle(color: AppColors.textBlack, fontSize: AppFontSizes.md)),
          CustomTextField(
            controller: _controller,
            hintText: "例：マニュアルが分かりやすくて良かった",
          ),
          const SizedBox(height: 8),
          Visibility(
            visible: _isFailed,
            child: Text(
              "送信に失敗しました。もう一度お試しください。",
              style: TextStyle(
                color: AppColors.error,
                fontSize: AppFontSizes.sm,
              )
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
                  isDisabled: _isSubmitting,
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
