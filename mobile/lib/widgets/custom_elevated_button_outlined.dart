import 'package:seeft_mobile/configs/importer.dart';

// 独自のスタイルを当てたElevatedButtonのアウトラインバージョン
class CustomElevatedButtonOutlined extends StatelessWidget {
  final VoidCallback onPressed;
  final String label;
  final IconData? icon;
  final bool isExpanded;  // 横幅を広げるかどうか

  const CustomElevatedButtonOutlined({
    super.key,
    required this.onPressed,
    required this.label,
    this.icon,
    this.isExpanded = false,
  });

  @override
  Widget build(BuildContext context) {
    return ElevatedButton.icon(
      onPressed: onPressed,
      icon: icon != null ? Icon(icon) : null,
      label: Text(
        label,
        style: const TextStyle(
          fontSize: AppFontSizes.sm,
          color: AppColors.grayDark
        ),
      ),
      style: ElevatedButton.styleFrom(
        backgroundColor: AppColors.base,
        padding: const EdgeInsets.symmetric(
          vertical: 10.0,
          horizontal: 24.0
        ),
        elevation: 0,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(100.0),
          side: BorderSide(
            color: AppColors.grayDark,
            width: 1.0,
          ),
        ),
        // 横幅を広げるための設定
        minimumSize: isExpanded ? const Size(double.infinity, 40) : null,
      ),
    );
  }
}
