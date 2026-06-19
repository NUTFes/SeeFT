import 'package:flutter/material.dart';
import 'package:seeft_mobile/theme/tokens.dart';

class CustomAppBar extends StatelessWidget implements PreferredSizeWidget {
  final String title;

  const CustomAppBar({
    super.key,
    required this.title,
  });

  @override
  Widget build(BuildContext context) {
    return AppBar(
      title: Text(title, style: TextStyle(
        color: AppColors.textWhite,
        fontSize: AppFontSizes.lg,
        fontWeight: FontWeight.bold,
      )),
      centerTitle: false,
      toolbarHeight: 63,
      backgroundColor: AppColors.main,
    );
  }

  @override
  Size get preferredSize => Size.fromHeight(63); // AppBarの標準の高さ
}