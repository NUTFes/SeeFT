import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:seeft_mobile/widgets/review_bottom_sheet.dart';

// 星は2問×5個の計10個。前半5個がシフトの人数、後半5個がマニュアルの評価
const _staffingFirstStarIndex = 0;
const _manualFirstStarIndex = 5;

Widget _wrap() => const MaterialApp(
  home: Scaffold(
    body: SingleChildScrollView(
      child: ReviewForm(taskName: 'テストタスク', userID: 1),
    ),
  ),
);

// 「送信」ボタンの実体。isDisabled のとき onPressed が null になる
ElevatedButton _submitButton(WidgetTester tester) {
  return tester.widget<ElevatedButton>(
    find.ancestor(
      of: find.text('送信'),
      matching: find.byType(ElevatedButton),
    ),
  );
}

Future<void> _tapStar(WidgetTester tester, int firstIndex, int stars) async {
  await tester.tap(find.byType(IconButton).at(firstIndex + stars - 1));
  await tester.pump();
}

void main() {
  testWidgets('初期状態は未選択で、送信ボタンを押せない', (tester) async {
    await tester.pumpWidget(_wrap());

    // 初期値を入れていないので、点灯している星は1つも無い
    expect(find.byIcon(Icons.star), findsNothing);
    expect(find.byIcon(Icons.star_border), findsNWidgets(10));

    expect(_submitButton(tester).onPressed, isNull);
    expect(find.text('星を選ぶと送信できます。'), findsOneWidget);
  });

  testWidgets('片方だけ選んでも送信ボタンは押せない', (tester) async {
    await tester.pumpWidget(_wrap());

    await _tapStar(tester, _staffingFirstStarIndex, 3);

    expect(find.byIcon(Icons.star), findsNWidgets(3));
    expect(_submitButton(tester).onPressed, isNull);
  });

  testWidgets('2問とも選ぶと送信ボタンが押せるようになる', (tester) async {
    await tester.pumpWidget(_wrap());

    await _tapStar(tester, _staffingFirstStarIndex, 3);
    await _tapStar(tester, _manualFirstStarIndex, 4);

    expect(find.byIcon(Icons.star), findsNWidgets(7));
    expect(_submitButton(tester).onPressed, isNotNull);
    expect(find.text('星を選ぶと送信できます。'), findsNothing);
  });

  testWidgets('コメント欄はDBのVARCHAR(255)に合わせて255文字で打ち切る', (tester) async {
    await tester.pumpWidget(_wrap());

    expect(tester.widget<TextField>(find.byType(TextField)).maxLength, 255);
  });
}
