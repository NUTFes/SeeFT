import 'package:seeft_mobile/configs/importer.dart';

import 'package:flutter/material.dart';

class WaitPage extends StatefulWidget {
  const WaitPage({super.key, this.message = 'シフトが完成するまでお待ちください'});

  /// 待機中に表示する文言。呼び出し元の用途に合わせて差し替える。
  final String message;

  @override
  State<WaitPage> createState() => _WaitPageState();
}

class _WaitPageState extends State<WaitPage> {
    int manualLength = 0;

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(
              widget.message,
              textAlign: TextAlign.center,
              style: const TextStyle(fontSize: 20),
            ),
          ],
        ),
      ),
    );
  }
}
