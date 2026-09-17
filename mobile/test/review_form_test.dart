import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:seeft_mobile/widgets/review_bottom_sheet.dart';

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

// 星は行ごとにキーで特定する。位置で数えるとウィジェットが増えたときに黙ってずれる
Future<void> _tapStar(WidgetTester tester, Key rowKey, int stars) async {
  await tester.tap(
    find
      .descendant(of: find.byKey(rowKey), matching: find.byType(IconButton))
      .at(stars - 1),
  );
  await tester.pump();
}

TextEditingController _commentController(WidgetTester tester) {
  return tester.widget<TextField>(find.byType(TextField)).controller!;
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

    await _tapStar(tester, staffingRatingRowKey, 3);

    expect(find.byIcon(Icons.star), findsNWidgets(3));
    expect(_submitButton(tester).onPressed, isNull);
  });

  testWidgets('2問とも選ぶと送信ボタンが押せるようになる', (tester) async {
    await tester.pumpWidget(_wrap());

    await _tapStar(tester, staffingRatingRowKey, 3);
    await _tapStar(tester, manualRatingRowKey, 4);

    expect(find.byIcon(Icons.star), findsNWidgets(7));
    expect(_submitButton(tester).onPressed, isNotNull);
    expect(find.text('星を選ぶと送信できます。'), findsNothing);
  });

  testWidgets('コメント欄は255文字を超える入力を打ち切る', (tester) async {
    await tester.pumpWidget(_wrap());

    await tester.enterText(find.byType(TextField), 'あ' * 300);
    await tester.pump();

    expect(_commentController(tester).text.characters.length, 255);
  });

  group('clampToCodePoints', () {
    test('上限以下はそのまま返す', () {
      expect(clampToCodePoints('マニュアルが分かりやすかった', 255), 'マニュアルが分かりやすかった');
    });

    // Web では IME 変換中の文字が maxLength の対象外になるため、
    // 送信時にもコードポイント数で切り詰める
    test('上限を超えたらコードポイント数で切り詰める', () {
      expect(clampToCodePoints('あ' * 300, 255).runes.length, 255);
    });

    // 書記素クラスタで数える maxLength を満たしていても、
    // PostgreSQL が数えるコードポイントでは超えることがある
    test('結合文字を含んでもコードポイント数で収まる', () {
      const family = '👨‍👩‍👧‍👦'; // 1書記素 / 11コードポイント
      final input = family * 255;

      expect(input.characters.length, 255);
      expect(input.runes.length, greaterThan(255));
      expect(clampToCodePoints(input, 255).runes.length, 255);
    });
  });
}
