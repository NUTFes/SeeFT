import 'package:seeft_mobile/configs/importer.dart';

import 'package:seeft_mobile/configs/importer.dart';
import 'package:flutter/material.dart';
import 'package:url_launcher/url_launcher.dart';

class WaitPage extends StatefulWidget {
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
    final Size size = MediaQuery.of(context).size;
    return Scaffold(
      appBar: AppBar(
        title: const Text(''),
        actions: <Widget>[],
        // debug
      ),
      drawer: drawer.applicationDrawer(context),
      body: Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children:[ Text("「シフトが完成するまでお待ちください」",style: TextStyle(fontSize: 40),),],
        ),
      ),
    );
  }
}
