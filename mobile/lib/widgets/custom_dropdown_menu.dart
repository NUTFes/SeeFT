import 'package:seeft_mobile/configs/importer.dart';

// ドロップダウンメニュー
class CustomDropdownMenu<T> extends StatelessWidget {
  final List<DropdownMenuEntry<T>> dropdownMenuEntries;
  final T? initialSelection;
  final void Function(T?)? onSelected;

  const CustomDropdownMenu({
    super.key,
    required this.dropdownMenuEntries,
    this.initialSelection,
    this.onSelected,
  });

  @override
  Widget build(BuildContext context) {
    return DropdownMenu<T>(
      dropdownMenuEntries: dropdownMenuEntries,
      initialSelection: initialSelection,
      enableFilter: true, // フィルター機能を有効にする
      onSelected: onSelected,
      // 共通のstyleをここで適用
      textStyle: const TextStyle(
        fontSize: AppFontSizes.md,
        color: AppColors.textBlack,
      ),
      menuStyle: MenuStyle(
        backgroundColor: MaterialStatePropertyAll(AppColors.base),
        elevation: MaterialStatePropertyAll(4),
      ),
      width: double.infinity, // 横幅を広げる
      inputDecorationTheme: InputDecorationTheme(
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
        // プレースホルダーのスタイル
        hintStyle: const TextStyle(
          color: AppColors.grayDark,
          fontSize: AppFontSizes.md,
        ),
      ),
    );
  }
}
