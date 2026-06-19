import 'package:seeft_mobile/configs/importer.dart';

import 'package:seeft_mobile/configs/importer.dart';
import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher.dart';

class WaitPage extends StatefulWidget {
  const WaitPage({super.key});

  @override
  _WaitPageState createState() => new _WaitPageState();
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
          children:[ Text("シフトが完成するまでお待ちください",style: TextStyle(fontSize: 20),),],
        ),
      ),
    );
  }
}
