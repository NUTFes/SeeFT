import 'package:flutter/material.dart';
import 'package:seeft_mobile/theme/tokens.dart';

class NewBadge extends StatelessWidget {
  final String text;
  final Color backgroundColor;
  final Color textColor;
  final double fontSize;
  final EdgeInsetsGeometry padding;
  final double borderRadius;

  const NewBadge({
    Key? key,
    this.text = 'new!!',
    this.backgroundColor = AppColors.main,
    this.textColor = AppColors.textWhite,
    this.fontSize = AppFontSizes.xs,
    this.padding = const EdgeInsets.symmetric(horizontal: 4, vertical: 2),
    this.borderRadius = AppBorderRadius.normal,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: padding,
      decoration: BoxDecoration(
        color: backgroundColor,
        borderRadius: BorderRadius.circular(borderRadius),
      ),
      child: Text(
        text,
        style: TextStyle(
          color: textColor,
          fontSize: fontSize,
          fontWeight: FontWeight.bold,
        ),
      ),
    );
  }
}
