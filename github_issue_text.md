# 開発概要
#### 目的
<!-- この機能を実装する目的 または このタスクの目的を記入 -->
- シフト一覧画面のUI改善として、未読のシフトに「New!」バッジを表示するモック機能を実装する
- 未読のシフトを視覚的に識別できるようにし、ユーザーが新しいシフト情報に気づきやすくする
- 将来的なAPI連携に向けたUIプロトタイプ（フロントエンドのみで完結）を作成する

#### 開発期間
- 開始日：
- 締切日：

# 考えられる開発内容
<!-- チェックボックスで書くと「todo: △/○」みたいに表示されるので良き -->
<!-- 編集対象のファイルパスが決まっていれば、なるべく書いておくといいかも -->

## 1. ShiftCard Widgetの変更
**ファイル**: `mobile/lib/widgets/shift_card.dart`

- [x] `StatelessWidget`から`StatefulWidget`に変更
- [x] 未読管理フラグ `bool isRead = false;` を追加（モック用、初期値は未読）
- [x] カード全体を`Stack`で囲み、右上に「New!」バッジを配置
- [x] `GestureDetector`でカード全体をタップ可能にし、タップ時に`setState`で`isRead = true`に変更

### バッジのデザイン仕様
- 背景色: `Colors.red`
- 文字色: `Colors.white`
- 文字サイズ: `10.0`
- 角丸: `BorderRadius.circular(8.0)`
- パディング: `EdgeInsets.symmetric(horizontal: 6.0, vertical: 2.0)`
- 位置: カードの右上（`Positioned(top: 4.0, right: 4.0)`）
- 表示条件: `isRead == false` の時のみ表示

## 2. モックデータの追加
**ファイル**: `mobile/lib/pages/my_shift_page.dart`

- [x] `_getMockShiftCardDataList()` 関数を追加
  - 2つのシフトカード（「受付」「案内」）を含むモックデータを生成
- [x] 開発用フラグ `_useMockDataAlways = true` を追加
  - キャッシュデータを無視して常にモックデータを表示
- [x] `_loadShiftCardDataList()` メソッドを修正
  - キャッシュデータがない、または空の場合にモックデータを自動表示
  - APIから取得したデータが空の場合もモックデータを使用

## 3. 動作確認
- [x] アプリを起動し、`http://localhost:45029/layout` にアクセス
- [x] 「マイシフト」タブを開く
- [x] 各シフトカードの右上に赤い「New!」バッジが表示されることを確認
- [x] シフトカードをタップすると、バッジが消えることを確認（既読状態になる）
- [x] 画面をリロードすると、バッジが再表示されることを確認（モック実装のため状態はリセット）

# 備考
- **モック実装**: 今回はバックエンド（API）やDBへの保存処理は実装していません
- **状態の永続化なし**: 画面をリロードすると未読状態に戻ります
- **既存機能の維持**: `ExpansionTile`の展開/折りたたみ機能やマニュアルボタンのタップイベントは既存のまま動作します
- **開発用フラグ**: `_useMockDataAlways = true` により、常にモックデータが表示されます。本番環境では `false` に変更する必要があります

## 変更ファイル
- `mobile/lib/widgets/shift_card.dart` - シフトカードWidget（StatefulWidget化、バッジUI実装）
- `mobile/lib/pages/my_shift_page.dart` - シフト一覧ページ（モックデータ追加）
- `mobile/lib/models/shift_card.dart` - シフトカードのデータモデル（参照のみ）

## 今後の拡張予定
- API連携による未読状態の取得
- バックエンドへの既読状態の保存
- 未読状態の永続化（ローカルストレージまたはサーバー側）
- 未読数のバッジ表示（複数の未読シフトがある場合）

# 参考

### 開発の流れ
1. PMにIssue（タスク）をもらう
2. 開発をする（↓の「リンク」の『開発のやり方』を見よう！）
3. チェックボックスを押していこう
4. ヤバい状況になったらIssueの右側にあるStatusを「Help」にしてPMにSlackで連絡しよう
5. チェックボックスが全部押せたらプルリクを作ろう
6. レビューを待とう
7. 修正点があれば修正しよう。なければPMがマージします！お疲れ様！

### SeeFTのタスク管理のルール
1. タスクは全てGit-Hub Projectで管理する
2. 全てのタスクに期日を決める
3. 毎週タスクの進捗を確認する（MTに出られない人はSlackで報告）
4. 毎週忙しさ（消化できるタスク量）を共有する
5. Helpは余裕のある人がいれば巻き取る。いなければ期日を変更する

### リンク
- [開発のやり方](https://github.com/NUTFes/SeeFT/wiki/Git%E3%81%A8GitHub%E3%81%AE%E4%BD%BF%E3%81%84%E6%96%B9%EF%BC%88SeeFT%EF%BC%89#%E3%82%88%E3%81%97%E9%96%8B%E7%99%BA%E3%81%97%E3%82%88%E3%81%86%E3%81%A3%E3%81%A6%E3%81%AA%E3%81%A3%E3%81%9F%E3%81%A8%E3%81%8D)
- [SeeFTのタスク管理のルール](https://github.com/NUTFes/SeeFT/wiki/SeeFT%E3%81%AE%E3%82%BF%E3%82%B9%E3%82%AF%E7%AE%A1%E7%90%86%E3%81%AE%E3%83%AB%E3%83%BC%E3%83%AB)
- [環境構築](https://github.com/NUTFes/SeeFT/wiki/%E7%92%B0%E5%A2%83%E6%A7%8B%E7%AF%89)
- [GitとGitHubの使い方（SeeFT）](https://github.com/NUTFes/SeeFT/wiki/Git%E3%81%A8GitHub%E3%81%AE%E4%BD%BF%E3%81%84%E6%96%B9%EF%BC%88SeeFT%EF%BC%89)
- [SeeFT Projects](https://github.com/orgs/NUTFes/projects/1/views/7)
- [SeeFT Wiki](https://github.com/NUTFes/SeeFT/wiki)
- [Figma](https://www.figma.com/files/team/1034705640918925952/project/108792039)
