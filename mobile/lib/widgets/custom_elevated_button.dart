import 'package:seeft_mobile/configs/importer.dart';

// 独自のスタイルを当てたElevatedButton
class CustomElevatedButton extends StatelessWidget {
  final VoidCallback onPressed;
  final String label;
  final IconData? icon;

  const CustomElevatedButton({
    super.key,
    required this.onPressed,
    required this.label,
    this.icon,
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
          color: AppColors.textWhite
        ),
      ),
      style: ElevatedButton.styleFrom(
        backgroundColor: AppColors.main,
        padding: const EdgeInsets.symmetric(
          vertical: 10.0,
          horizontal: 24.0
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
}
