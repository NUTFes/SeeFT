import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:seeft_mobile/theme/tokens.dart';
import 'package:seeft_mobile/widgets/manual_list_item.dart';
import 'package:url_launcher_platform_interface/link.dart';
import 'package:url_launcher_platform_interface/url_launcher_platform_interface.dart';

class _FakeUrlLauncherPlatform extends UrlLauncherPlatform {
  _FakeUrlLauncherPlatform(this.launchResult, {this.throwOnLaunch = false});

  final bool launchResult;

  // プラットフォーム側が例外を投げる場合。_open は戻り値だけでなく例外も
  // 失敗として扱うので、その経路を通すために使う
  final bool throwOnLaunch;

  // 直近のlaunchUrl呼び出しの引数（URL・起動モード）を検証用に記録する
  String? lastLaunchedUrl;
  LaunchOptions? lastLaunchOptions;

  @override
  LinkDelegate? get linkDelegate => null;

  // canLaunchが常にtrueでも、URLが空の行は押せてはいけない。
  // 事前チェックに頼らず、押せるかどうかをURLの有無で決めていることを確かめる
  @override
  Future<bool> canLaunch(String url) async => true;

  @override
  Future<bool> launchUrl(String url, LaunchOptions options) async {
    lastLaunchedUrl = url;
    lastLaunchOptions = options;
    if (throwOnLaunch) {
      throw Exception('launch failed');
    }
    return launchResult;
  }
}

const _taskName = 'Music Nuts 参加';
const _documentUrl = 'https://docs.google.com/document/d/test/edit';
const _slideUrl = 'https://seeft-api.nutfes.net/manuals/test';

Widget _wrap({
  String url = _documentUrl,
  String manualUrl = _slideUrl,
}) {
  return MaterialApp(
    home: Scaffold(
      body: ManualListItem(
        taskName: _taskName,
        url: url,
        manualUrl: manualUrl,
      ),
    ),
  );
}

_FakeUrlLauncherPlatform _installFakePlatform(
  bool launchResult, {
  bool throwOnLaunch = false,
}) {
  final originalPlatform = UrlLauncherPlatform.instance;
  final fake = _FakeUrlLauncherPlatform(
    launchResult,
    throwOnLaunch: throwOnLaunch,
  );
  UrlLauncherPlatform.instance = fake;
  addTearDown(() => UrlLauncherPlatform.instance = originalPlatform);
  return fake;
}

void main() {
  // issue #524 の症状。45thの本番では490件中235件がこの状態だった
  group('ドキュメント版が無い行', () {
    testWidgets('タップしても何も開かない', (tester) async {
      final fake = _installFakePlatform(true);

      await tester.pumpWidget(_wrap(url: ''));
      await tester.pumpAndSettle();

      await tester.tap(find.text(_taskName));
      await tester.pumpAndSettle();

      expect(fake.lastLaunchedUrl, isNull);
    });

    testWidgets('onTapがnullで、押下エフェクトも出ない', (tester) async {
      _installFakePlatform(true);

      await tester.pumpWidget(_wrap(url: ''));
      await tester.pumpAndSettle();

      expect(tester.widget<ListTile>(find.byType(ListTile)).onTap, isNull);
    });

    testWidgets('タスク名がグレーで表示される', (tester) async {
      _installFakePlatform(true);

      await tester.pumpWidget(_wrap(url: ''));
      await tester.pumpAndSettle();

      expect(
        tester.widget<Text>(find.text(_taskName)).style?.color,
        AppColors.grayDark,
      );
    });
  });

  group('ドキュメント版がある行', () {
    testWidgets('タップすると外部アプリモードで開く', (tester) async {
      final fake = _installFakePlatform(true);

      await tester.pumpWidget(_wrap());
      await tester.pumpAndSettle();

      await tester.tap(find.text(_taskName));
      await tester.pumpAndSettle();

      expect(fake.lastLaunchedUrl, _documentUrl);
      // 認証付き配信は埋め込めない。シフトカードと同じく別タブで開く
      expect(
        fake.lastLaunchOptions?.mode,
        PreferredLaunchMode.externalApplication,
      );
    });

    testWidgets('タスク名が通常の文字色で表示される', (tester) async {
      _installFakePlatform(true);

      await tester.pumpWidget(_wrap());
      await tester.pumpAndSettle();

      expect(
        tester.widget<Text>(find.text(_taskName)).style?.color,
        AppColors.textBlack,
      );
    });

    testWidgets('起動に失敗したらエラーメッセージを出す', (tester) async {
      _installFakePlatform(false);

      await tester.pumpWidget(_wrap());
      await tester.pumpAndSettle();

      await tester.tap(find.text(_taskName));
      await tester.pumpAndSettle();

      expect(find.text('マニュアルを開けませんでした'), findsOneWidget);
    });

    // _open の try/catch は Uri.parse と launchUrl の両方の例外を拾う。
    // こちらは launchUrl 側が投げた場合
    testWidgets('launchUrlが例外を投げたらエラーメッセージを出す', (tester) async {
      _installFakePlatform(true, throwOnLaunch: true);

      await tester.pumpWidget(_wrap());
      await tester.pumpAndSettle();

      await tester.tap(find.text(_taskName));
      await tester.pumpAndSettle();

      expect(find.text('マニュアルを開けませんでした'), findsOneWidget);
    });

    // WebのlaunchUrlはjavascript:以外でfalseを返さないので、
    // 実際に失敗扱いになるのは主にURLの解析エラー
    testWidgets('URLが壊れていたらエラーメッセージを出す', (tester) async {
      _installFakePlatform(true);

      await tester.pumpWidget(_wrap(url: 'https://[broken'));
      await tester.pumpAndSettle();

      await tester.tap(find.text(_taskName));
      await tester.pumpAndSettle();

      expect(find.text('マニュアルを開けませんでした'), findsOneWidget);
    });
  });

  group('スライド版', () {
    testWidgets('無ければアイコンを出さない', (tester) async {
      _installFakePlatform(true);

      await tester.pumpWidget(_wrap(manualUrl: ''));
      await tester.pumpAndSettle();

      expect(find.byIcon(Icons.slideshow), findsNothing);
    });

    testWidgets('あればアイコンから外部アプリモードで開く', (tester) async {
      final fake = _installFakePlatform(true);

      await tester.pumpWidget(_wrap());
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.slideshow));
      await tester.pumpAndSettle();

      expect(fake.lastLaunchedUrl, _slideUrl);
      expect(
        fake.lastLaunchOptions?.mode,
        PreferredLaunchMode.externalApplication,
      );
    });

    testWidgets('起動に失敗したらエラーメッセージを出す', (tester) async {
      _installFakePlatform(false);

      await tester.pumpWidget(_wrap());
      await tester.pumpAndSettle();

      await tester.tap(find.byIcon(Icons.slideshow));
      await tester.pumpAndSettle();

      expect(find.text('マニュアル（スライド版）を開けませんでした'), findsOneWidget);
    });
  });

  // ドキュメント版が無くスライド版だけがある場合、行はグレーで押せないが
  // 右のアイコンからスライド版は開ける
  testWidgets('スライド版だけがある行は、行が押せなくてもアイコンからは開ける', (tester) async {
    final fake = _installFakePlatform(true);

    await tester.pumpWidget(_wrap(url: ''));
    await tester.pumpAndSettle();

    expect(tester.widget<ListTile>(find.byType(ListTile)).onTap, isNull);

    await tester.tap(find.byIcon(Icons.slideshow));
    await tester.pumpAndSettle();

    expect(fake.lastLaunchedUrl, _slideUrl);
  });
}
