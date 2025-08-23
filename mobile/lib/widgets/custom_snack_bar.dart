import 'package:seeft_mobile/configs/importer.dart';

// スナックバーを表示する関数
void showCustomSnackBar(BuildContext context, String message) {
  ScaffoldMessenger.of(context).showSnackBar(
    SnackBar(
      content: Row(
        children: [
          Icon(Icons.info_outlined, color: AppColors.textWhite, size: 16),
          const SizedBox(width: 8),
          Text(
            message,
            style: TextStyle(
              color: AppColors.textWhite,
              fontSize: AppFontSizes.sm,
            )),
        ],
      ),
      backgroundColor: AppColors.grayDark,
      duration: Duration(seconds: 2),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
      margin: const EdgeInsets.only(left: 16, right: 16, bottom: 8),
      behavior: SnackBarBehavior.floating,
    ),
  );
}