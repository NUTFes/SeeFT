import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:seeft_mobile/models/shift_card.dart';
import 'package:seeft_mobile/widgets/shift_card.dart';

ShiftCardData _fakeData(String taskName, String url) {
  return ShiftCardData(
    taskName: taskName,
    startTime: '10:00',
    endTime: '12:00',
    place: '正門',
    url: url,
    shiftMembers: const [],
    beforeMembers: ShiftMembers(sTime: '', eTime: '', members: const []),
    afterMembers: ShiftMembers(sTime: '', eTime: '', members: const []),
  );
}

Widget _wrap(List<ShiftCardData> orderedData) {
  return MaterialApp(
    home: Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: orderedData
              .map(
                (data) => ShiftCard(
                  key: ValueKey(data.taskName),
                  data: data,
                  userID: 1,
                ),
              )
              .toList(),
        ),
      ),
    ),
  );
}

Finder _manualToggleTextWithinCard(String taskName, String text) {
  return find.descendant(
    of: find.widgetWithText(Card, taskName),
    matching: find.text(text),
  );
}

void main() {
  testWidgets('Keyがあれば並び替え後もマニュアルの開閉状態が正しいカードに紐づく', (tester) async {
    final cardA = _fakeData('A', 'https://example.com/a');
    final cardB = _fakeData('B', 'https://example.com/b');

    await tester.pumpWidget(_wrap([cardA, cardB]));
    await tester.pumpAndSettle();

    expect(_manualToggleTextWithinCard('A', 'マニュアルを開く'), findsOneWidget);
    expect(_manualToggleTextWithinCard('B', 'マニュアルを開く'), findsOneWidget);

    await tester.tap(_manualToggleTextWithinCard('A', 'マニュアルを開く'));
    await tester.pumpAndSettle();

    expect(_manualToggleTextWithinCard('A', 'マニュアルを閉じる'), findsOneWidget);
    expect(_manualToggleTextWithinCard('B', 'マニュアルを開く'), findsOneWidget);

    // 並び順を入れ替えて再描画（一覧の再取得・並び替えを想定）
    await tester.pumpWidget(_wrap([cardB, cardA]));
    await tester.pumpAndSettle();

    // 開閉状態は位置ではなくKey（カードの中身）に紐づいたままであること
    expect(_manualToggleTextWithinCard('A', 'マニュアルを閉じる'), findsOneWidget);
    expect(_manualToggleTextWithinCard('B', 'マニュアルを開く'), findsOneWidget);
  });
}
