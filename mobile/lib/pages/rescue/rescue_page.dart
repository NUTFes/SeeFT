import 'package:seeft_mobile/configs/importer.dart';
import 'package:seeft_mobile/pages/rescue/rescue_request_tab/rescue_request_tab.dart';
import 'package:seeft_mobile/pages/rescue/rescue_response_tab/rescue_response_tab.dart';

class RescuePage extends StatefulWidget {
  const RescuePage({super.key});
  
  @override
  _RescuePageState createState() => _RescuePageState();
}

class _RescuePageState extends State<RescuePage>
    with TickerProviderStateMixin {
  
  final _tabs = <Tab>[
    Tab(text: "レスキュー送信"),
    Tab(text: "本部からの返答"),
  ];
  
  // late int _userID; // ユーザIDを格納する変数
  late TabController _tabController; // タブのコントローラーを格納する変数
  
  // ウィジェットの初期化時の処理
  @override
  void initState() {
    logger.i('RescuePage is being initialized.');
    // タブのコントローラーを初期化
    _tabController = TabController(
      length: 2, // タブの数を設定
      initialIndex: 0, // 初期タブを設定
      vsync: this
    );
    super.initState();
  }
  
  // ウィジェットが破棄されるときの処理
  @override
  void dispose() {  
    // タブコントローラーを破棄
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.base,
      appBar: AppBar(
        title: Text("緊急時対応", 
            style: TextStyle(
              color: AppColors.textWhite,
              fontSize: AppFontSizes.lg,
              fontWeight: FontWeight.bold,
            ),
          ),
        centerTitle: false,
        toolbarHeight: 63,
        backgroundColor: AppColors.main,
        bottom: PreferredSize(
          preferredSize: const Size.fromHeight(75.0),
          child: _tabBar(), // タブバー
        ),
      ),
      body: _tabBarView(), // タブの内容を表示するTabBarView
    );
  }
  
  // タブバー
  Widget _tabBar() {
    return TabBar(
      isScrollable: false,
      tabs: _tabs,
      controller: _tabController,
      padding: const EdgeInsets.only(top: 20.0, bottom: 20.0),
      // 選択中のタブのスタイル
      labelStyle: TextStyle(
        color: AppColors.main,
        fontSize: AppFontSizes.md,
        fontWeight: FontWeight.normal,
      ),
      // 選択されていないタブのスタイル
      unselectedLabelStyle: TextStyle(
        color: AppColors.grayLight,
        fontSize: AppFontSizes.md,
        fontWeight: FontWeight.normal,
      ),
      // タブのインジケーター
      indicator: BoxDecoration(
        color: AppColors.base,
        borderRadius: BorderRadius.circular(100),
      ),
      indicatorSize: TabBarIndicatorSize.label, // 現在使用中のFlutterのver3.27だとTabBarIndicatorSize.tabにするとアニメーションのバグが深刻化するのでlabelにしてます
      splashBorderRadius: BorderRadius.circular(100),
      dividerHeight: 0,
    );
  }
  
  // タブの内容を表示するTabBarView
  Widget _tabBarView() {
    return TabBarView(
      controller: _tabController,
      children: [
        // レスキューを送信するタブ
        RescueRequestTab(),
        // 本部からの返答タブ
        RescueResponseTab(),
      ],
    );
  }
}