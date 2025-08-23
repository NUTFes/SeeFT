import 'package:seeft_mobile/configs/importer.dart';

// 独自のスタイルを当てたテキストフィールド
class CustomTextField extends StatelessWidget {
  final TextEditingController? controller;
  final String? hintText;
  final String? labelText;
  final ValueChanged<String>? onChanged;
  final bool obscureText;
  final TextInputType? keyboardType;

  const CustomTextField({
    super.key,
    this.controller,
    this.hintText,
    this.labelText,
    this.onChanged,
    this.obscureText = false,
    this.keyboardType,
  });

  @override
  Widget build(BuildContext context) {
    return TextField(
      controller: controller,
      obscureText: obscureText, // パスワード入力などで使用
      onChanged: onChanged,
      keyboardType: keyboardType,
      // keyboardTypeがTextInputType.numberの場合数字のみの入力に制限する
      inputFormatters: keyboardType == TextInputType.number ? <TextInputFormatter>[
        FilteringTextInputFormatter.digitsOnly,
      ] : null,
      // 共通のスタイルをここで適用
      cursorColor: AppColors.main, 
      style: const TextStyle(
        color: AppColors.textBlack,
        fontSize: AppFontSizes.md,
      ),
      decoration: InputDecoration(
        // 枠線の色を指定
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(8.0),
          borderSide: const BorderSide(
            color: AppColors.grayDark,
            width: 1.0,
          ),
        ),
        // 有効時の枠線の色を指定
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(8.0),
          borderSide: const BorderSide(
            color: AppColors.grayDark,
            width: 1.0,
          ),
        ),
        // フォーカス時の枠線の色を指定
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(8.0),
          borderSide: const BorderSide(
            color: AppColors.main,
            width: 2.0,
          ),
        ),
        // エラー時の枠線の色を指定
        errorBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(8.0),
          borderSide: const BorderSide(
            color: AppColors.error,
            width: 1.0,
          ),
        ),
        // フォーカス時のエラー枠線の色を指定
        focusedErrorBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(8.0),
          borderSide: const BorderSide(
            color: AppColors.error,
            width: 2.0,
          ),
        ),
        // ラベルのテキスト
        labelText: labelText,
        // ラベルのスタイル
        labelStyle: const TextStyle(
          color: AppColors.grayDark,
          fontSize: AppFontSizes.md,
        ),
        // フローティングラベルのスタイル
        floatingLabelStyle: WidgetStateTextStyle.resolveWith(
          (Set<WidgetState> states) {
            if (states.contains(WidgetState.error)) {
              return const TextStyle(
                color: AppColors.error,
                fontSize: AppFontSizes.md,
              );
            }
            if (states.contains(WidgetState.focused)) {
              return const TextStyle(
                color: AppColors.main,
                fontSize: AppFontSizes.md,
              );
            }
            return const TextStyle(
              color: AppColors.grayDark,
              fontSize: AppFontSizes.md,
            );
          },
        ),
        // プレースホルダーのテキスト
        hintText: hintText,
        // プレースホルダーのスタイル
        hintStyle: const TextStyle(
          color: AppColors.grayDark,
          fontSize: AppFontSizes.md,
        ),
      ),
    );
  }
}
