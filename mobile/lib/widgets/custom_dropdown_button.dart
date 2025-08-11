import 'package:seeft_mobile/configs/importer.dart';

// ドロップダウンボタンのカスタムウィジェット
class CustomDropdownButton<T> extends StatelessWidget {
  final T? value;
  final List<DropdownMenuItem<T>> items;
  final ValueChanged<T?>? onChanged;
  final String? hintText;
  final bool isDense;
  final bool isExpanded;

  const CustomDropdownButton({
    super.key,
    required this.value,
    required this.items,
    this.onChanged,
    this.hintText,
    this.isDense = false,
    this.isExpanded = true,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
      decoration: BoxDecoration(
        border: Border.all(color: AppColors.grayLight, width: 1.0),
        borderRadius: BorderRadius.circular(8.0),
        color: AppColors.base,
      ),
      child: DropdownButton<T>(
        value: value,
        items: items,
        onChanged: onChanged,
        hint: hintText != null
            ? Text(
              hintText!,
              style: TextStyle(
                fontSize: AppFontSizes.md,
                color: AppColors.grayDark,
              ),
            )
            : null,
        padding: EdgeInsets.zero,     // 内側の余白をなくす
        isExpanded: isExpanded,       // 幅いっぱいに広げる
        isDense: isDense,             // 高さをコンパクトにする
        underline: const SizedBox(),  // 下線を消す
        style: const TextStyle(
          fontSize: AppFontSizes.md,
          color: AppColors.textBlack,
        ),
        dropdownColor: AppColors.base, // ドロップダウンの背景色
        iconSize: 24.0,
        icon: const Icon(Icons.arrow_drop_down, color: AppColors.grayDark),
      ),
    );
  }
}
