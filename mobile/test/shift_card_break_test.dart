import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:seeft_mobile/models/shift_card.dart';
import 'package:seeft_mobile/theme/tokens.dart';
import 'package:seeft_mobile/widgets/shift_card.dart';

// 休憩カード。APIは休憩の担当者を返さないが、集合場所だけはタスクの既定値が入って届く。
// mobile側でそれを表示しないことを確かめたいので、あえてplaceを埋めておく
ShiftCardData _breakData() {
  return ShiftCardData(
    taskName: breakTaskName,
    startTime: '12:00',
    endTime: '13:00',
    place: '本部(電気棟1F)',
    url: 'https://example.com/manual',
    manualUrl: 'https://example.com/manual.html',
    shiftMembers: const [],
    beforeMembers: ShiftMembers(sTime: '', eTime: '', members: const []),
    afterMembers: ShiftMembers(sTime: '', eTime: '', members: const []),
  );
}

ShiftCardData _normalData() {
  return ShiftCardData(
    taskName: '受付',
    startTime: '10:00',
    endTime: '12:00',
    place: '正門',
    url: 'https://example.com/doc',
    manualUrl: 'https://example.com/manual.html',
    shiftMembers: [
      ShiftMembers(sTime: '10:00', eTime: '12:00', members: const []),
    ],
    beforeMembers: ShiftMembers(sTime: '', eTime: '', members: const []),
    afterMembers: ShiftMembers(sTime: '', eTime: '', members: const []),
  );
}

Widget _wrap(ShiftCardData data, {bool isNew = false}) {
  return MaterialApp(
    home: Scaffold(
      body: Center(
        child: SizedBox(
          width: 320,
          child: ShiftCard(data: data, userID: 1, isNew: isNew),
        ),
      ),
    ),
  );
}

void main() {
  test('タスク名で休憩を判定する', () {
    expect(_breakData().isBreak, isTrue);
    expect(_normalData().isBreak, isFalse);
  });

  testWidgets('休憩カードは時刻とタスク名だけを表示する', (tester) async {
    await tester.pumpWidget(_wrap(_breakData()));
    await tester.pumpAndSettle();

    expect(find.text('12:00~13:00'), findsOneWidget);
    expect(find.text('休憩'), findsOneWidget);

    // 集合場所とマニュアルは休憩に存在しないため出さない。
    // マニュアルは_ManualToggleが実際に描画する文言で確認する
    // ('マニュアル'という素の文字列はどの状態でも描画されないため)
    expect(find.text('本部(電気棟1F)'), findsNothing);
    expect(find.byIcon(Icons.location_on_outlined), findsNothing);
    expect(find.text('ドキュメント版を開く'), findsNothing);
    expect(find.text('ドキュメント版なし'), findsNothing);
    expect(find.text('スライド版を別のタブで開く'), findsNothing);
  });

  testWidgets('休憩カードには三点メニューを出さない', (tester) async {
    await tester.pumpWidget(_wrap(_breakData()));
    await tester.pumpAndSettle();

    // 担当者一覧もレビューも対象が無いため、入口ごと出さない
    expect(find.byIcon(Icons.more_vert), findsNothing);
  });

  testWidgets('通常のシフトカードは集合場所・マニュアル・三点メニューを出す', (tester) async {
    await tester.pumpWidget(_wrap(_normalData()));
    await tester.pumpAndSettle();

    expect(find.text('正門'), findsOneWidget);
    expect(find.byIcon(Icons.location_on_outlined), findsOneWidget);
    expect(find.byIcon(Icons.more_vert), findsOneWidget);
    expect(find.text('ドキュメント版を開く'), findsOneWidget);
    expect(find.text('スライド版を別のタブで開く'), findsOneWidget);
  });

  testWidgets('休憩カードは背景色で通常シフトと区別される', (tester) async {
    await tester.pumpWidget(_wrap(_breakData()));
    await tester.pumpAndSettle();
    final breakCard = tester.widget<Card>(find.byType(Card));
    expect(breakCard.color, AppColors.grayLight);

    await tester.pumpWidget(_wrap(_normalData()));
    await tester.pumpAndSettle();
    final normalCard = tester.widget<Card>(find.byType(Card));
    expect(normalCard.color, AppColors.base);
  });
}
