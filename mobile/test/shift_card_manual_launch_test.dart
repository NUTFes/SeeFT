import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:seeft_mobile/models/shift_card.dart';
import 'package:seeft_mobile/widgets/shift_card.dart';
import 'package:url_launcher_platform_interface/link.dart';
import 'package:url_launcher_platform_interface/url_launcher_platform_interface.dart';

class _FakeUrlLauncherPlatform extends UrlLauncherPlatform {
  _FakeUrlLauncherPlatform(this.launchResult);

  final bool launchResult;

  @override
  LinkDelegate? get linkDelegate => null;

  @override
  Future<bool> canLaunch(String url) async => true;

  @override
  Future<bool> launchUrl(String url, LaunchOptions options) async {
    return launchResult;
  }
}

ShiftCardData _fakeDataWithManual() {
  return ShiftCardData(
    taskName: '受付',
    startTime: '10:00',
    endTime: '12:00',
    place: '正門',
    url: 'https://example.com/manual',
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

void main() {
  testWidgets('外部起動に失敗した場合はエラーメッセージを表示する', (tester) async {
    final originalPlatform = UrlLauncherPlatform.instance;
    UrlLauncherPlatform.instance = _FakeUrlLauncherPlatform(false);
    addTearDown(() => UrlLauncherPlatform.instance = originalPlatform);

    await tester.pumpWidget(_wrap());
    await tester.pumpAndSettle();

    await tester.tap(find.text('マニュアルを開く'));
    await tester.pumpAndSettle();

    await tester.ensureVisible(find.text('別のタブで開く'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('別のタブで開く'));
    await tester.pumpAndSettle();

    expect(find.text('マニュアルを開けませんでした'), findsOneWidget);
  });

  testWidgets('外部起動に成功した場合はエラーメッセージを表示しない', (tester) async {
    final originalPlatform = UrlLauncherPlatform.instance;
    UrlLauncherPlatform.instance = _FakeUrlLauncherPlatform(true);
    addTearDown(() => UrlLauncherPlatform.instance = originalPlatform);

    await tester.pumpWidget(_wrap());
    await tester.pumpAndSettle();

    await tester.tap(find.text('マニュアルを開く'));
    await tester.pumpAndSettle();

    await tester.ensureVisible(find.text('別のタブで開く'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('別のタブで開く'));
    await tester.pumpAndSettle();

    expect(find.text('マニュアルを開けませんでした'), findsNothing);
  });
}
