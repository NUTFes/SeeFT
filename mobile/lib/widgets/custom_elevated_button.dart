import 'package:seeft_mobile/configs/importer.dart';

// 独自のスタイルを当てたElevatedButton
class CustomElevatedButton extends StatelessWidget {
  final VoidCallback onPressed;
  final String label;
  final IconData? icon;
  final bool isDisabled;

  const CustomElevatedButton({
    super.key,
    required this.onPressed,
    required this.label,
    this.icon,
    this.isDisabled = false,
  });

  @override
  Widget build(BuildContext context) {
    return ElevatedButton.icon(
      onPressed: isDisabled? null : onPressed,
      icon: icon != null
        ? Icon(
          icon,
          color: isDisabled ? AppColors.grayDark : AppColors.textWhite,
        )
        : null,
      label: Text(
        label,
        style: TextStyle(
          fontSize: AppFontSizes.sm,
          color: isDisabled ? AppColors.grayDark : AppColors.textWhite,
        ),
      ),
      style: ElevatedButton.styleFrom(
        backgroundColor: isDisabled ? AppColors.grayLight : AppColors.main,
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
