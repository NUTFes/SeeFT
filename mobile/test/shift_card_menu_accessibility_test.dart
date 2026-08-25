import 'package:flutter/material.dart';
import 'package:flutter/semantics.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:seeft_mobile/models/shift_card.dart';
import 'package:seeft_mobile/widgets/shift_card.dart';

ShiftCardData _fakeData() {
  return ShiftCardData(
    taskName: '受付',
    startTime: '10:00',
    endTime: '12:00',
    place: '正門',
    url: '',
    manualUrl: '',
    shiftMembers: [
      ShiftMembers(sTime: '10:00', eTime: '12:00', members: []),
    ],
    beforeMembers: ShiftMembers(sTime: '', eTime: '', members: []),
    afterMembers: ShiftMembers(sTime: '', eTime: '', members: []),
  );
}

Widget _wrap() {
  return MaterialApp(
    home: Scaffold(
      body: Center(
        child: SizedBox(
          width: 320,
          child: ShiftCard(data: _fakeData(), userID: 1),
        ),
      ),
    ),
  );
}

void main() {
  testWidgets('三点メニューが開き、担当者一覧ダイアログを表示する（タップ操作）', (tester) async {
    await tester.pumpWidget(_wrap());
    await tester.pumpAndSettle();

    await tester.tap(find.byIcon(Icons.more_vert));
    await tester.pumpAndSettle();

    expect(find.text('担当者一覧'), findsOneWidget);
    expect(find.text('レビューを書く'), findsOneWidget);

    await tester.tap(find.text('担当者一覧'));
    await tester.pump(const Duration(milliseconds: 200));
    await tester.pumpAndSettle();

    expect(find.text('【担当者の一覧】'), findsOneWidget);
  });

  testWidgets('三点ボタンとメニュー項目にボタンとしてのSemanticsが付いている', (tester) async {
    final handle = tester.ensureSemantics();
    await tester.pumpWidget(_wrap());
    await tester.pumpAndSettle();

    final moreButtonSemantics =
        tester.getSemantics(find.byIcon(Icons.more_vert));
    expect(moreButtonSemantics.hasFlag(SemanticsFlag.isButton), isTrue);

    await tester.tap(find.byIcon(Icons.more_vert));
    await tester.pumpAndSettle();

    final membersRowSemantics = tester.getSemantics(find.text('担当者一覧'));
    expect(membersRowSemantics.hasFlag(SemanticsFlag.isButton), isTrue);
    expect(membersRowSemantics.label, contains('担当者一覧'));

    handle.dispose();
  });

  testWidgets('三点ボタンにTabでフォーカスが当たりEnterでメニューが開く（キーボード操作）',
      (tester) async {
    await tester.pumpWidget(_wrap());
    await tester.pumpAndSettle();

    await tester.sendKeyEvent(LogicalKeyboardKey.tab);
    await tester.pumpAndSettle();

    expect(FocusManager.instance.primaryFocus, isNotNull);
    expect(find.text('担当者一覧'), findsNothing);

    await tester.sendKeyEvent(LogicalKeyboardKey.enter);
    await tester.pumpAndSettle();

    expect(find.text('担当者一覧'), findsOneWidget);
    expect(find.text('レビューを書く'), findsOneWidget);
  });
}
