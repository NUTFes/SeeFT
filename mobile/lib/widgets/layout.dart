import 'package:flutter/material.dart';
import 'package:seeft_mobile/theme/tokens.dart';
import 'package:seeft_mobile/widgets/app_bar.dart';
import 'package:seeft_mobile/pages/my_shift_page.dart';
import 'package:seeft_mobile/pages/manual_list_page.dart';
import 'package:seeft_mobile/pages/wait_page.dart';
import 'package:seeft_mobile/pages/rescue/rescue_page.dart';

class Layout extends StatefulWidget {
  const Layout({Key? key}) : super(key: key);

  @override
  State<Layout> createState() => _LayoutState();
}

class _LayoutState extends State<Layout> {
  var _currentPageIndex = 0;  // 現在のページのインデックス

  final _pages = <Widget>[
    MyShiftPage(),    // マイシフト
    ManualListPage(), // マニュアル
    RescuePage(),     // 緊急時対応
    WaitPage(),       // 仮のページ(後で「その他」に差し替える)
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: _currentPageIndex == 0 || _currentPageIndex == 2 ? null : CustomAppBar(  // マイシフト以外のページの場合はアプリバーを表示
        title: _currentPageIndex == 1 ? 'マニュアル' : 'その他',
      ),
      bottomNavigationBar: SizedBox(
        height: 78, // ナビゲーションバーの高さを設定
        child: BottomNavigationBar( // ナビゲーションバー
          onTap: (int index) {
            setState(() {
              _currentPageIndex = index;    // タップしたアイテムのインデックスに変更
            });
          },
          currentIndex: _currentPageIndex,
          elevation: 5,                               // ナビゲーションバーの影
          backgroundColor: AppColors.base,          // ナビゲーションバーの背景色
          type: BottomNavigationBarType.fixed,        // ラベルを常に表示
          selectedIconTheme: const IconThemeData(     // 選択時のアイコンのスタイル
            size: 24, 
            color: AppColors.main,
          ),
          selectedItemColor: AppColors.main,        // 選択時の文字の色
          selectedLabelStyle: const TextStyle(        // 選択時の文字のスタイル
            fontSize: AppFontSizes.xs, 
            fontWeight: FontWeight.bold, 
          ),
          unselectedIconTheme: const IconThemeData(   // 未選択時のアイコンのスタイル
            size: 24, 
            color: AppColors.grayDark,
          ),
          unselectedItemColor: AppColors.grayDark,  // 未選択時の文字の色
          unselectedLabelStyle: const TextStyle(      // 未選択時の文字のスタイル
            fontSize: AppFontSizes.xs, 
            fontWeight: FontWeight.normal, 
          ),
          items: const <BottomNavigationBarItem>[     // ナビゲーションバーのアイテム
            BottomNavigationBarItem(
              icon: Icon(Icons.today_outlined),
              activeIcon: Icon(Icons.today), 
              label: 'マイシフト', 
              tooltip: "マイシフト",
            ),
            BottomNavigationBarItem(
              icon: Icon(Icons.quiz_outlined),
              activeIcon: Icon(Icons.quiz),
              label: 'マニュアル',
              tooltip: "マニュアル"
            ),
            BottomNavigationBarItem(
              icon: Icon(Icons.error_outline),
              activeIcon: Icon(Icons.error),
              label: '緊急時対応',
              tooltip: "緊急時対応",
            ),
            BottomNavigationBarItem(
              icon: Icon(Icons.more_horiz_outlined),
              activeIcon: Icon(Icons.more_horiz),
              label: 'その他', 
              tooltip: "その他"
            ),
          ],
        ),
      ),
      body: _pages[_currentPageIndex],  // 現在のインデックスのページを表示
    );
  }
}