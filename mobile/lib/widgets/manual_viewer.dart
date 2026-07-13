// ignore: avoid_web_libraries_in_flutter
import 'dart:html' as html;
import 'dart:ui_web' as ui;
import 'package:flutter/material.dart';
import 'package:seeft_mobile/theme/tokens.dart';

// Google Docs / Google Slides の URL を埋め込み用に変換する
String _toEmbeddableUrl(String url) {
  final uri = Uri.tryParse(url);
  if (uri == null) return url;

  // Google Docs / Slides / Sheets の /edit や /view を /preview に変換
  if (uri.host == 'docs.google.com') {
    final newPath = uri.path
        .replaceFirst(RegExp(r'/edit$'), '/preview')
        .replaceFirst(RegExp(r'/view$'), '/preview');
    return uri.replace(path: newPath, queryParameters: {}).toString();
  }

  return url;
}

class ManualViewer extends StatefulWidget {
  final String url;

  const ManualViewer({super.key, required this.url});

  @override
  State<ManualViewer> createState() => _ManualViewerState();
}

class _ManualViewerState extends State<ManualViewer> {
  late final String _viewId;

  @override
  void initState() {
    super.initState();
    // viewId はウィジェットごとにユニークにする
    _viewId = 'manual-iframe-${widget.url.hashCode}';

    ui.platformViewRegistry.registerViewFactory(_viewId, (int viewId) {
      final embeddableUrl = _toEmbeddableUrl(widget.url);
      return html.IFrameElement()
        ..src = embeddableUrl
        ..style.border = 'none'
        ..style.width = '100%'
        ..style.height = '100%'
        ..allowFullscreen = true;
    });
  }

  @override
  Widget build(BuildContext context) {
    return ClipRRect(
      borderRadius: BorderRadius.circular(8.0),
      child: Container(
        height: 480,
        decoration: BoxDecoration(
          border: Border.all(color: AppColors.grayLight),
          borderRadius: BorderRadius.circular(8.0),
        ),
        child: HtmlElementView(viewType: _viewId),
      ),
    );
  }
}
