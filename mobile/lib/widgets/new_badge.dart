import 'package:seeft_mobile/configs/importer.dart';

class NewBadge extends StatelessWidget {
  final String text;
  final Color backgroundColor;
  final Color textColor;
  final double fontSize;
  final EdgeInsetsGeometry padding;
  final double borderRadius;

  const NewBadge({
    super.key,
    this.text = 'new!!',
    this.backgroundColor = AppColors.main,
    this.textColor = AppColors.textWhite,
    this.fontSize = AppFontSizes.xs,
    this.padding = const EdgeInsets.symmetric(horizontal: 4, vertical: 2),
    this.borderRadius = AppBorderRadius.normal,
  });

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
