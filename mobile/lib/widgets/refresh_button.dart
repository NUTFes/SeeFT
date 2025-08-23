import 'package:seeft_mobile/configs/importer.dart';

// 更新ボタン
class RefreshButton extends StatelessWidget {
  final VoidCallback onPressed;
  // final String label;
  // final IconData? icon;
  final bool isLoading;

  const RefreshButton({
    super.key,
    required this.onPressed,
    // required this.label,
    // this.icon,
    this.isLoading = false,
  });

  @override
  Widget build(BuildContext context) {
    return ElevatedButton.icon(
      onPressed: isLoading ? null : onPressed,
      icon: Icon(
        Icons.refresh,
        color: isLoading ? AppColors.grayDark : AppColors.main,
        size: 16,
      ),
      label: Text(
        // label,
        isLoading? "読み込み中...": "更新",
        style: TextStyle(
          fontSize: AppFontSizes.sm,
          color: isLoading ? AppColors.grayDark : AppColors.main,
        ),
      ),
      style: ElevatedButton.styleFrom(
        backgroundColor: AppColors.base,
        disabledBackgroundColor: AppColors.base,
        elevation: 0,
        side: BorderSide(
          color: isLoading? AppColors.grayDark: AppColors.main,
          width: 1.0,
        ),
        padding: const EdgeInsets.symmetric(
          vertical: 10.0,
          horizontal: 24.0
        ),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(100.0),
        ),
        // 横幅を広げるための設定
        // minimumSize: Size(double.infinity, 40),
      ),
    );
  }
}
