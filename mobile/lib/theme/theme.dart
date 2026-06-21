import "package:flutter/material.dart";

class MaterialTheme {
  final TextTheme textTheme;

  const MaterialTheme(this.textTheme);

  static ColorScheme lightScheme() {
    return const ColorScheme(
      brightness: Brightness.light,
      primary: Color(0xff006a60),
      surfaceTint: Color(0xff006a60),
      onPrimary: Color(0xffffffff),
      primaryContainer: Color(0xff9ef2e4),
      onPrimaryContainer: Color(0xff005048),
      secondary: Color(0xff4a635f),
      onSecondary: Color(0xffffffff),
      secondaryContainer: Color(0xffcce8e2),
      onSecondaryContainer: Color(0xff334b47),
      tertiary: Color(0xff456179),
      onTertiary: Color(0xffffffff),
      tertiaryContainer: Color(0xffcce5ff),
      onTertiaryContainer: Color(0xff2d4961),
      error: Color(0xffba1a1a),
      onError: Color(0xffffffff),
      errorContainer: Color(0xffffdad6),
      onErrorContainer: Color(0xff93000a),
      surface: Color(0xfff4fbf8),
      onSurface: Color(0xff161d1c),
      onSurfaceVariant: Color(0xff3f4947),
      outline: Color(0xff6f7977),
      outlineVariant: Color(0xffbec9c6),
      shadow: Color(0xff000000),
      scrim: Color(0xff000000),
      inverseSurface: Color(0xff2b3230),
      inversePrimary: Color(0xff82d5c8),
      primaryFixed: Color(0xff9ef2e4),
      onPrimaryFixed: Color(0xff00201c),
      primaryFixedDim: Color(0xff82d5c8),
      onPrimaryFixedVariant: Color(0xff005048),
      secondaryFixed: Color(0xffcce8e2),
      onSecondaryFixed: Color(0xff05201c),
      secondaryFixedDim: Color(0xffb1ccc6),
      onSecondaryFixedVariant: Color(0xff334b47),
      tertiaryFixed: Color(0xffcce5ff),
      onTertiaryFixed: Color(0xff001e31),
      tertiaryFixedDim: Color(0xffadcae6),
      onTertiaryFixedVariant: Color(0xff2d4961),
      surfaceDim: Color(0xffd5dbd9),
      surfaceBright: Color(0xfff4fbf8),
      surfaceContainerLowest: Color(0xffffffff),
      surfaceContainerLow: Color(0xffeff5f2),
      surfaceContainer: Color(0xffe9efed),
      surfaceContainerHigh: Color(0xffe3eae7),
      surfaceContainerHighest: Color(0xffdde4e1),
    );
  }

  ThemeData light() {
    return theme(lightScheme());
  }

  static ColorScheme lightMediumContrastScheme() {
    return const ColorScheme(
      brightness: Brightness.light,
      primary: Color(0xff003e37),
      surfaceTint: Color(0xff006a60),
      onPrimary: Color(0xffffffff),
      primaryContainer: Color(0xff1d7a6f),
      onPrimaryContainer: Color(0xffffffff),
      secondary: Color(0xff223b37),
      onSecondary: Color(0xffffffff),
      secondaryContainer: Color(0xff58726d),
      onSecondaryContainer: Color(0xffffffff),
      tertiary: Color(0xff1b394f),
      onTertiary: Color(0xffffffff),
      tertiaryContainer: Color(0xff547089),
      onTertiaryContainer: Color(0xffffffff),
      error: Color(0xff740006),
      onError: Color(0xffffffff),
      errorContainer: Color(0xffcf2c27),
      onErrorContainer: Color(0xffffffff),
      surface: Color(0xfff4fbf8),
      onSurface: Color(0xff0c1211),
      onSurfaceVariant: Color(0xff2e3836),
      outline: Color(0xff4b5452),
      outlineVariant: Color(0xff656f6d),
      shadow: Color(0xff000000),
      scrim: Color(0xff000000),
      inverseSurface: Color(0xff2b3230),
      inversePrimary: Color(0xff82d5c8),
      primaryFixed: Color(0xff1d7a6f),
      onPrimaryFixed: Color(0xffffffff),
      primaryFixedDim: Color(0xff006056),
      onPrimaryFixedVariant: Color(0xffffffff),
      secondaryFixed: Color(0xff58726d),
      onSecondaryFixed: Color(0xffffffff),
      secondaryFixedDim: Color(0xff405a55),
      onSecondaryFixedVariant: Color(0xffffffff),
      tertiaryFixed: Color(0xff547089),
      onTertiaryFixed: Color(0xffffffff),
      tertiaryFixedDim: Color(0xff3c586f),
      onTertiaryFixedVariant: Color(0xffffffff),
      surfaceDim: Color(0xffc1c8c5),
      surfaceBright: Color(0xfff4fbf8),
      surfaceContainerLowest: Color(0xffffffff),
      surfaceContainerLow: Color(0xffeff5f2),
      surfaceContainer: Color(0xffe3eae7),
      surfaceContainerHigh: Color(0xffd8dedc),
      surfaceContainerHighest: Color(0xffcdd3d0),
    );
  }

  ThemeData lightMediumContrast() {
    return theme(lightMediumContrastScheme());
  }

  static ColorScheme lightHighContrastScheme() {
    return const ColorScheme(
      brightness: Brightness.light,
      primary: Color(0xff00332d),
      surfaceTint: Color(0xff006a60),
      onPrimary: Color(0xffffffff),
      primaryContainer: Color(0xff00534b),
      onPrimaryContainer: Color(0xffffffff),
      secondary: Color(0xff17302d),
      onSecondary: Color(0xffffffff),
      secondaryContainer: Color(0xff354e4a),
      onSecondaryContainer: Color(0xffffffff),
      tertiary: Color(0xff102e44),
      onTertiary: Color(0xffffffff),
      tertiaryContainer: Color(0xff304c63),
      onTertiaryContainer: Color(0xffffffff),
      error: Color(0xff600004),
      onError: Color(0xffffffff),
      errorContainer: Color(0xff98000a),
      onErrorContainer: Color(0xffffffff),
      surface: Color(0xfff4fbf8),
      onSurface: Color(0xff000000),
      onSurfaceVariant: Color(0xff000000),
      outline: Color(0xff252e2c),
      outlineVariant: Color(0xff414b49),
      shadow: Color(0xff000000),
      scrim: Color(0xff000000),
      inverseSurface: Color(0xff2b3230),
      inversePrimary: Color(0xff82d5c8),
      primaryFixed: Color(0xff00534b),
      onPrimaryFixed: Color(0xffffffff),
      primaryFixedDim: Color(0xff003a34),
      onPrimaryFixedVariant: Color(0xffffffff),
      secondaryFixed: Color(0xff354e4a),
      onSecondaryFixed: Color(0xffffffff),
      secondaryFixedDim: Color(0xff1e3733),
      onSecondaryFixedVariant: Color(0xffffffff),
      tertiaryFixed: Color(0xff304c63),
      onTertiaryFixed: Color(0xffffffff),
      tertiaryFixedDim: Color(0xff17354b),
      onTertiaryFixedVariant: Color(0xffffffff),
      surfaceDim: Color(0xffb4bab8),
      surfaceBright: Color(0xfff4fbf8),
      surfaceContainerLowest: Color(0xffffffff),
      surfaceContainerLow: Color(0xffecf2ef),
      surfaceContainer: Color(0xffdde4e1),
      surfaceContainerHigh: Color(0xffcfd6d3),
      surfaceContainerHighest: Color(0xffc1c8c5),
    );
  }

  ThemeData lightHighContrast() {
    return theme(lightHighContrastScheme());
  }

  static ColorScheme darkScheme() {
    return const ColorScheme(
      brightness: Brightness.dark,
      primary: Color(0xff82d5c8),
      surfaceTint: Color(0xff82d5c8),
      onPrimary: Color(0xff003731),
      primaryContainer: Color(0xff005048),
      onPrimaryContainer: Color(0xff9ef2e4),
      secondary: Color(0xffb1ccc6),
      onSecondary: Color(0xff1c3531),
      secondaryContainer: Color(0xff334b47),
      onSecondaryContainer: Color(0xffcce8e2),
      tertiary: Color(0xffadcae6),
      onTertiary: Color(0xff153349),
      tertiaryContainer: Color(0xff2d4961),
      onTertiaryContainer: Color(0xffcce5ff),
      error: Color(0xffffb4ab),
      onError: Color(0xff690005),
      errorContainer: Color(0xff93000a),
      onErrorContainer: Color(0xffffdad6),
      surface: Color(0xff0e1513),
      onSurface: Color(0xffdde4e1),
      onSurfaceVariant: Color(0xffbec9c6),
      outline: Color(0xff899390),
      outlineVariant: Color(0xff3f4947),
      shadow: Color(0xff000000),
      scrim: Color(0xff000000),
      inverseSurface: Color(0xffdde4e1),
      inversePrimary: Color(0xff006a60),
      primaryFixed: Color(0xff9ef2e4),
      onPrimaryFixed: Color(0xff00201c),
      primaryFixedDim: Color(0xff82d5c8),
      onPrimaryFixedVariant: Color(0xff005048),
      secondaryFixed: Color(0xffcce8e2),
      onSecondaryFixed: Color(0xff05201c),
      secondaryFixedDim: Color(0xffb1ccc6),
      onSecondaryFixedVariant: Color(0xff334b47),
      tertiaryFixed: Color(0xffcce5ff),
      onTertiaryFixed: Color(0xff001e31),
      tertiaryFixedDim: Color(0xffadcae6),
      onTertiaryFixedVariant: Color(0xff2d4961),
      surfaceDim: Color(0xff0e1513),
      surfaceBright: Color(0xff343b39),
      surfaceContainerLowest: Color(0xff090f0e),
      surfaceContainerLow: Color(0xff161d1c),
      surfaceContainer: Color(0xff1a2120),
      surfaceContainerHigh: Color(0xff252b2a),
      surfaceContainerHighest: Color(0xff303635),
    );
  }

  ThemeData dark() {
    return theme(darkScheme());
  }

  static ColorScheme darkMediumContrastScheme() {
    return const ColorScheme(
      brightness: Brightness.dark,
      primary: Color(0xff98ecde),
      surfaceTint: Color(0xff82d5c8),
      onPrimary: Color(0xff002b26),
      primaryContainer: Color(0xff4a9e93),
      onPrimaryContainer: Color(0xff000000),
      secondary: Color(0xffc6e2dc),
      onSecondary: Color(0xff112a26),
      secondaryContainer: Color(0xff7c9691),
      onSecondaryContainer: Color(0xff000000),
      tertiary: Color(0xffc3e0fc),
      onTertiary: Color(0xff07283e),
      tertiaryContainer: Color(0xff7894ae),
      onTertiaryContainer: Color(0xff000000),
      error: Color(0xffffd2cc),
      onError: Color(0xff540003),
      errorContainer: Color(0xffff5449),
      onErrorContainer: Color(0xff000000),
      surface: Color(0xff0e1513),
      onSurface: Color(0xffffffff),
      onSurfaceVariant: Color(0xffd4dfdb),
      outline: Color(0xffaab4b1),
      outlineVariant: Color(0xff889290),
      shadow: Color(0xff000000),
      scrim: Color(0xff000000),
      inverseSurface: Color(0xffdde4e1),
      inversePrimary: Color(0xff005249),
      primaryFixed: Color(0xff9ef2e4),
      onPrimaryFixed: Color(0xff001512),
      primaryFixedDim: Color(0xff82d5c8),
      onPrimaryFixedVariant: Color(0xff003e37),
      secondaryFixed: Color(0xffcce8e2),
      onSecondaryFixed: Color(0xff001512),
      secondaryFixedDim: Color(0xffb1ccc6),
      onSecondaryFixedVariant: Color(0xff223b37),
      tertiaryFixed: Color(0xffcce5ff),
      onTertiaryFixed: Color(0xff001321),
      tertiaryFixedDim: Color(0xffadcae6),
      onTertiaryFixedVariant: Color(0xff1b394f),
      surfaceDim: Color(0xff0e1513),
      surfaceBright: Color(0xff3f4644),
      surfaceContainerLowest: Color(0xff040807),
      surfaceContainerLow: Color(0xff181f1e),
      surfaceContainer: Color(0xff232928),
      surfaceContainerHigh: Color(0xff2d3432),
      surfaceContainerHighest: Color(0xff383f3d),
    );
  }

  ThemeData darkMediumContrast() {
    return theme(darkMediumContrastScheme());
  }

  static ColorScheme darkHighContrastScheme() {
    return const ColorScheme(
      brightness: Brightness.dark,
      primary: Color(0xffaffff2),
      surfaceTint: Color(0xff82d5c8),
      onPrimary: Color(0xff000000),
      primaryContainer: Color(0xff7ed1c4),
      onPrimaryContainer: Color(0xff000e0c),
      secondary: Color(0xffdaf6f0),
      onSecondary: Color(0xff000000),
      secondaryContainer: Color(0xffadc8c2),
      onSecondaryContainer: Color(0xff000e0c),
      tertiary: Color(0xffe6f1ff),
      onTertiary: Color(0xff000000),
      tertiaryContainer: Color(0xffa9c6e2),
      onTertiaryContainer: Color(0xff000c18),
      error: Color(0xffffece9),
      onError: Color(0xff000000),
      errorContainer: Color(0xffffaea4),
      onErrorContainer: Color(0xff220001),
      surface: Color(0xff0e1513),
      onSurface: Color(0xffffffff),
      onSurfaceVariant: Color(0xffffffff),
      outline: Color(0xffe8f2ef),
      outlineVariant: Color(0xffbac5c2),
      shadow: Color(0xff000000),
      scrim: Color(0xff000000),
      inverseSurface: Color(0xffdde4e1),
      inversePrimary: Color(0xff005249),
      primaryFixed: Color(0xff9ef2e4),
      onPrimaryFixed: Color(0xff000000),
      primaryFixedDim: Color(0xff82d5c8),
      onPrimaryFixedVariant: Color(0xff001512),
      secondaryFixed: Color(0xffcce8e2),
      onSecondaryFixed: Color(0xff000000),
      secondaryFixedDim: Color(0xffb1ccc6),
      onSecondaryFixedVariant: Color(0xff001512),
      tertiaryFixed: Color(0xffcce5ff),
      onTertiaryFixed: Color(0xff000000),
      tertiaryFixedDim: Color(0xffadcae6),
      onTertiaryFixedVariant: Color(0xff001321),
      surfaceDim: Color(0xff0e1513),
      surfaceBright: Color(0xff4b5150),
      surfaceContainerLowest: Color(0xff000000),
      surfaceContainerLow: Color(0xff1a2120),
      surfaceContainer: Color(0xff2b3230),
      surfaceContainerHigh: Color(0xff363d3b),
      surfaceContainerHighest: Color(0xff414846),
    );
  }

  ThemeData darkHighContrast() {
    return theme(darkHighContrastScheme());
  }


  ThemeData theme(ColorScheme colorScheme) => ThemeData(
    useMaterial3: true,
    brightness: colorScheme.brightness,
    colorScheme: colorScheme,
    textTheme: textTheme.apply(
      bodyColor: colorScheme.onSurface,
      displayColor: colorScheme.onSurface,
    ),
    // textTheme: GoogleFonts.mPlus1pTextTheme(Theme.of(context).textTheme),  
    scaffoldBackgroundColor: colorScheme.surface,
    canvasColor: colorScheme.surface,
    pageTransitionsTheme: const PageTransitionsTheme( // 画面遷移時のアニメーションをiOS風にする
      builders: <TargetPlatform, PageTransitionsBuilder>{
        TargetPlatform.android: CupertinoPageTransitionsBuilder(),
        TargetPlatform.iOS: CupertinoPageTransitionsBuilder()
      },
    ),
  );


  List<ExtendedColor> get extendedColors => [
  ];
}

class ExtendedColor {
  final Color seed, value;
  final ColorFamily light;
  final ColorFamily lightHighContrast;
  final ColorFamily lightMediumContrast;
  final ColorFamily dark;
  final ColorFamily darkHighContrast;
  final ColorFamily darkMediumContrast;

  const ExtendedColor({
    required this.seed,
    required this.value,
    required this.light,
    required this.lightHighContrast,
    required this.lightMediumContrast,
    required this.dark,
    required this.darkHighContrast,
    required this.darkMediumContrast,
  });
}

class ColorFamily {
  const ColorFamily({
    required this.color,
    required this.onColor,
    required this.colorContainer,
    required this.onColorContainer,
  });

  final Color color;
  final Color onColor;
  final Color colorContainer;
  final Color onColorContainer;
}
