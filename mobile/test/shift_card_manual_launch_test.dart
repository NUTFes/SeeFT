import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:seeft_mobile/models/shift_card.dart';
import 'package:seeft_mobile/widgets/shift_card.dart';
import 'package:url_launcher_platform_interface/link.dart';
import 'package:url_launcher_platform_interface/url_launcher_platform_interface.dart';

class _FakeUrlLauncherPlatform extends UrlLauncherPlatform {
  _FakeUrlLauncherPlatform(this.launchResult);

  final bool launchResult;

  // 直近のlaunchUrl呼び出しの引数（URL・起動モード）を検証用に記録する
  String? lastLaunchedUrl;
  LaunchOptions? lastLaunchOptions;

  @override
  LinkDelegate? get linkDelegate => null;

  @override
  Future<bool> canLaunch(String url) async => true;

  @override
  Future<bool> launchUrl(String url, LaunchOptions options) async {
    lastLaunchedUrl = url;
    lastLaunchOptions = options;
    return launchResult;
  }
}

const _slideUrl = 'https://seeft-api.nutfes.net/manuals/test';

ShiftCardData _fakeDataWithManual() {
  return ShiftCardData(
    taskName: '受付',
    startTime: '10:00',
    endTime: '12:00',
    place: '正門',
    url: 'https://example.com/manual',
    manualUrl: _slideUrl,
    shiftMembers: const [],
    beforeMembers: ShiftMembers(sTime: '', eTime: '', members: const []),
    afterMembers: ShiftMembers(sTime: '', eTime: '', members: const []),
  );
}

Widget _wrap() {
  return MaterialApp(
    home: Scaffold(
      body: SingleChildScrollView(
        child: ShiftCard(data: _fakeDataWithManual(), userID: 1),
      ),
    ),
  );
}

_FakeUrlLauncherPlatform _installFakePlatform(bool launchResult) {
  final originalPlatform = UrlLauncherPlatform.instance;
  final fake = _FakeUrlLauncherPlatform(launchResult);
  UrlLauncherPlatform.instance = fake;
  addTearDown(() => UrlLauncherPlatform.instance = originalPlatform);
  return fake;
}

void main() {
  testWidgets('外部起動に失敗した場合はエラーメッセージを表示する', (tester) async {
    _installFakePlatform(false);

    await tester.pumpWidget(_wrap());
    await tester.pumpAndSettle();

    await tester.tap(find.text('ドキュメント版を開く'));
    await tester.pumpAndSettle();

    await tester.ensureVisible(find.text('ドキュメント版を別のタブで開く'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('ドキュメント版を別のタブで開く'));
    await tester.pumpAndSettle();

    expect(find.text('マニュアルを開けませんでした'), findsOneWidget);
  });

  testWidgets('外部起動に成功した場合はエラーメッセージを表示しない', (tester) async {
    _installFakePlatform(true);

    await tester.pumpWidget(_wrap());
    await tester.pumpAndSettle();

    await tester.tap(find.text('ドキュメント版を開く'));
    await tester.pumpAndSettle();

    await tester.ensureVisible(find.text('ドキュメント版を別のタブで開く'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('ドキュメント版を別のタブで開く'));
    await tester.pumpAndSettle();

    expect(find.text('マニュアルを開けませんでした'), findsNothing);
  });

  testWidgets('スライド版の起動成功時はエラーを出さず外部アプリモードで開く', (tester) async {
    final fake = _installFakePlatform(true);

    await tester.pumpWidget(_wrap());
    await tester.pumpAndSettle();

    await tester.ensureVisible(find.text('スライド版を別のタブで開く'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('スライド版を別のタブで開く'));
    await tester.pumpAndSettle();

    expect(fake.lastLaunchedUrl, _slideUrl);
    // 認証付き配信のため埋め込み不可。必ず外部アプリ(別タブ)モードで開くこと
    expect(
      fake.lastLaunchOptions?.mode,
      PreferredLaunchMode.externalApplication,
    );
    expect(find.text('マニュアル（スライド版）を開けませんでした'), findsNothing);
  });

  testWidgets('スライド版の起動失敗時はエラーメッセージを表示する', (tester) async {
    _installFakePlatform(false);

    await tester.pumpWidget(_wrap());
    await tester.pumpAndSettle();

    await tester.ensureVisible(find.text('スライド版を別のタブで開く'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('スライド版を別のタブで開く'));
    await tester.pumpAndSettle();

    expect(find.text('マニュアル（スライド版）を開けませんでした'), findsOneWidget);
  });
}
