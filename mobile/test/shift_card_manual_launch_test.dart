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

const _documentUrl = 'https://drive.google.com/file/d/test/view';
const _slideUrl = 'https://seeft-api.nutfes.net/manuals/test';

ShiftCardData _fakeData({
  String url = _documentUrl,
  String manualUrl = _slideUrl,
}) {
  return ShiftCardData(
    taskName: '受付',
    startTime: '10:00',
    endTime: '12:00',
    place: '正門',
    url: url,
    manualUrl: manualUrl,
    shiftMembers: const [],
    beforeMembers: ShiftMembers(sTime: '', eTime: '', members: const []),
    afterMembers: ShiftMembers(sTime: '', eTime: '', members: const []),
  );
}

Widget _wrap(ShiftCardData data) {
  return MaterialApp(
    home: Scaffold(
      body: SingleChildScrollView(
        child: ShiftCard(data: data, userID: 1),
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

Future<void> _tapText(WidgetTester tester, String text) async {
  await tester.ensureVisible(find.text(text));
  await tester.pumpAndSettle();
  await tester.tap(find.text(text));
  await tester.pumpAndSettle();
}

void main() {
  testWidgets('ドキュメント版の起動成功時はエラーを出さず外部アプリモードで開く', (tester) async {
    final fake = _installFakePlatform(true);

    await tester.pumpWidget(_wrap(_fakeData()));
    await tester.pumpAndSettle();

    await _tapText(tester, 'ドキュメント版を別のタブで開く');

    expect(fake.lastLaunchedUrl, _documentUrl);
    // 閲覧を技大祭アカウントに限定しているため埋め込み不可。必ず外部アプリ(別タブ)モードで開くこと
    expect(
      fake.lastLaunchOptions?.mode,
      PreferredLaunchMode.externalApplication,
    );
    expect(find.text('マニュアルを開けませんでした'), findsNothing);
  });

  testWidgets('ドキュメント版の起動失敗時はエラーメッセージを表示する', (tester) async {
    _installFakePlatform(false);

    await tester.pumpWidget(_wrap(_fakeData()));
    await tester.pumpAndSettle();

    await _tapText(tester, 'ドキュメント版を別のタブで開く');

    expect(find.text('マニュアルを開けませんでした'), findsOneWidget);
  });

  testWidgets('ドキュメント版を押してもカードの中に埋め込み表示しない', (tester) async {
    _installFakePlatform(true);

    await tester.pumpWidget(_wrap(_fakeData()));
    await tester.pumpAndSettle();

    await _tapText(tester, 'ドキュメント版を別のタブで開く');

    expect(find.byType(HtmlElementView), findsNothing);
    expect(find.text('ドキュメント版を閉じる'), findsNothing);
  });

  testWidgets('ドキュメント版が無いタスクは押しても何も開かない', (tester) async {
    final fake = _installFakePlatform(true);

    await tester.pumpWidget(_wrap(_fakeData(url: '')));
    await tester.pumpAndSettle();

    await _tapText(tester, 'ドキュメント版なし');

    expect(fake.lastLaunchedUrl, isNull);
  });

  testWidgets('スライド版の起動成功時はエラーを出さず外部アプリモードで開く', (tester) async {
    final fake = _installFakePlatform(true);

    await tester.pumpWidget(_wrap(_fakeData()));
    await tester.pumpAndSettle();

    await _tapText(tester, 'スライド版を別のタブで開く');

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

    await tester.pumpWidget(_wrap(_fakeData()));
    await tester.pumpAndSettle();

    await _tapText(tester, 'スライド版を別のタブで開く');

    expect(find.text('マニュアル（スライド版）を開けませんでした'), findsOneWidget);
  });
}
