import 'package:flutter/material.dart';

// 色
class AppColors {
  static const Color base = Color(0xFFFFFFFF);      // 背景色
  static const Color main = Color(0xFF009688);      // メインカラー
  static const Color sub = Color(0xFFF3AE56);       // サブカラー
  static const Color error = Color(0xFFBA1A1A);     // エラー時の色
  static const Color white = Color(0xFFFFFFFF);     // 白色
  static const Color grayLight = Color(0xFFD9D9D9); // 明るい灰色
  static const Color grayDark = Color(0xFF7B7B7B);  // 暗い灰色
  static const Color black = Color(0xFF000000);     // 黒色
  static const Color textWhite = Color(0xFFFFFFFF); // 明るい文字色
  static const Color textBlack = Color(0xFF000000); // 暗い文字色
  static const Color link = Color(0xFF1264A3);      // 青色のリンク
}

// フォントサイズ
class AppFontSizes {
  static const double xs = 12.0;  // 最小の文字
  static const double sm = 14.0;  // 小さい文字
  static const double md = 16.0;  // 標準の文字
  static const double lg = 20.0;  // 大きい文字
}

// 角丸
class AppBorderRadius {
  static const double normal = 8.0;    // 標準の角丸
  static const double button = 100.0;  // タップできる要素の角丸
}

// 影
class AppShadows {
  // 標準の影
  static const List<BoxShadow> medium = [
    BoxShadow(
      color: Color(0x40000000), // 25%の不透明度の黒
      offset: Offset(0, 2),
      blurRadius: 2,
      spreadRadius: 0,
    ),
    BoxShadow(
      color: Color(0x40000000), // 25%の不透明度の黒
      offset: Offset(2, 2),
      blurRadius: 4,
      spreadRadius: 0,
    ),
  ];
}