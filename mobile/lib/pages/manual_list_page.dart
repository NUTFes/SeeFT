// import 'dart:developer';

import 'package:seeft_mobile/configs/importer.dart';
// import 'package:flutter_local_notifications/flutter_local_notifications.dart';
// import 'package:http/http.dart' as http;
import 'package:url_launcher/url_launcher.dart';
import 'package:seeft_mobile/pages/wait_page.dart';

class ManualListPage extends StatefulWidget {
  const ManualListPage({super.key});

  @override
  State<ManualListPage> createState() => _ManualListPageState();
}

class _ManualListPageState extends State<ManualListPage> {
// notification関連をinitStateに書き出さなきゃいけないので書いてたけどutilとかに書いてもいいかもね

//  FlutterLocalNotificationsPlugin flutterLocalNotificationsPlugin;
//  NotificationDetails platformChannelSpecifics;
  List<dynamic> _allManuals = [];
  String _searchQuery = '';
  final TextEditingController _searchController = TextEditingController();
  final FocusNode _searchFocusNode = FocusNode();
  Timer? _debounce;
  bool _isLoading = true;
  bool _hasError = false;

  @override
  void initState() {
    super.initState();
    // onChangedではなくcontrollerのリスナーを使う。
    // IMEの変換確定は「テキストは同じでcomposing範囲だけがクリアされる」
    // 変化として来ることがあり、onChangedでは取りこぼす場合があるため。
    _searchController.addListener(_onSearchChanged);
    _loadData();
  }

  @override
  void dispose() {
    _debounce?.cancel();
    _searchController.removeListener(_onSearchChanged);
    _searchController.dispose();
    _searchFocusNode.dispose();
    super.dispose();
  }

  // 検索入力のハンドラ（controllerリスナー）。
  // 入力のたびに200msのデバウンスタイマーを張り直し、タイマー発火時点で
  // まだIME変換中（composingが有効）ならフィルタ反映を見送る（仕様5.2）。
  // 変換確定時はcomposingがクリアされて再度このリスナーが呼ばれるため、
  // その後のタイマーで確実に反映される。
  void _onSearchChanged() {
    // クリアボタンの表示/非表示を即時反映する
    setState(() {});
    _debounce?.cancel();
    _debounce = Timer(const Duration(milliseconds: 200), () {
      if (!mounted) return;
      // IME変換中はまだ確定していないので反映しない
      if (_searchController.value.composing.isValid) return;
      setState(() {
        _searchQuery = _searchController.text;
      });
    });
  }

  void _clearSearch() {
    _debounce?.cancel();
    _searchController.clear();
    setState(() {
      _searchQuery = '';
    });
    // 連続検索しやすいようにフォーカスは検索ボックスに残す（仕様5.5）
    _searchFocusNode.requestFocus();
  }

  Future<void> _loadData() async {
    try {
      final res = await api.getAllManual();
      if (mounted) {
        setState(() {
          _allManuals = res as List<dynamic>;
          _hasError = false;
          _isLoading = false;
        });
      }
    } catch (err) {
      logger.e('don`t response. error message: $err');
      if (mounted) {
        // ロード失敗は「0件ヒット」とは別物として扱う。
        // これがないとクエリ未入力時に空リスト＝真っ白な画面になる。
        setState(() {
          _hasError = true;
          _isLoading = false;
        });
      }
    }
  }

  // エラー表示からの再読み込み。
  void _retry() {
    setState(() {
      _isLoading = true;
      _hasError = false;
    });
    _loadData();
  }

  List<dynamic> get _filteredManuals {
    // クエリ・タスク名の双方をトリム＋小文字化して部分一致（仕様5.3/6.3）
    final query = _searchQuery.trim().toLowerCase();
    if (query.isEmpty) return _allManuals;
    return _allManuals
        .where((m) => m["task"].toString().trim().toLowerCase().contains(query))
        .toList();
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return const WaitPage(message: 'マニュアルを読み込んでいます');
    }
    if (_hasError) {
      return Scaffold(
        backgroundColor: AppColors.base,
        body: _buildErrorView(),
      );
    }
    final manuals = _filteredManuals;
    return Scaffold(
      backgroundColor: AppColors.base,
      body: Container(
        padding: const EdgeInsets.all(32.0),
        child: Column(
          children: <Widget>[
            Semantics(
              label: 'マニュアル検索',
              textField: true,
              child: TextField(
                controller: _searchController,
                focusNode: _searchFocusNode,
                decoration: InputDecoration(
                  hintText: '検索',
                  prefixIcon: const Icon(Icons.search),
                  suffixIcon: _searchController.text.isNotEmpty
                      ? IconButton(
                          icon: const Icon(Icons.clear),
                          onPressed: _clearSearch,
                        )
                      : null,
                  enabledBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(AppBorderRadius.normal),
                    borderSide: const BorderSide(color: AppColors.grayLight),
                  ),
                  focusedBorder: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(AppBorderRadius.normal),
                    borderSide: const BorderSide(color: AppColors.main),
                  ),
                ),
              ),
            ),
            const SizedBox(height: 16.0),
            Flexible(
              child: manuals.isEmpty && _searchQuery.trim().isNotEmpty
                  ? const Center(
                      child: Text('該当するマニュアルが見つかりませんでした'),
                    )
                  : ListView.builder(
                      itemCount: manuals.length,
                      itemBuilder: (BuildContext context, int index) {
                        return SizedBox(
                          height: 40,
                          child: _manualItem(manuals, index, context),
                        );
                      },
                    ),
            ),
          ],
        ),
      ),
    );
  }

  // ロード失敗時の表示（再読み込み可能）。
  Widget _buildErrorView() {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32.0),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Text(
              'マニュアルの読み込みに失敗しました',
              textAlign: TextAlign.center,
              style: TextStyle(
                color: AppColors.textBlack,
                fontSize: AppFontSizes.md,
              ),
            ),
            const SizedBox(height: 16.0),
            ElevatedButton(
              style: ElevatedButton.styleFrom(
                backgroundColor: AppColors.main,
                foregroundColor: AppColors.textWhite,
              ),
              onPressed: _retry,
              child: const Text('再読み込み'),
            ),
          ],
        ),
      ),
    );
  }

  Widget _manualItem(var manuals, index, context) {
    // スライド版マニュアルのURL。未設定のタスクではボタンを出さない
    final manualUrl = (manuals[index]["manualUrl"] ?? '').toString();
    return Container(
      decoration: BoxDecoration(
        border: Border(bottom: BorderSide(
          width: 1.0,
          color: AppColors.grayLight,
        )),
      ),
      child: ListTile(
        title: Text(
          manuals[index]["task"].toString(),
          style: TextStyle(
            color: AppColors.textBlack,
            fontSize: AppFontSizes.md,
          ),
        ),
        trailing: manualUrl.isEmpty
            ? null
            : IconButton(
                icon: const Icon(Icons.slideshow),
                onPressed: () async {
                  if (await canLaunchUrl(Uri.parse(manualUrl))) {
                    await launchUrl(Uri.parse(manualUrl));
                  }
                },
              ),
        onTap: () async {
          if (await canLaunchUrl(Uri.parse(manuals[index]["url"].toString()))) {
            await launchUrl(Uri.parse((manuals[index]["url"].toString())));
          }
        },
      ),
    );
  }
}
