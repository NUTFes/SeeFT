import 'package:flutter/material.dart';
import 'package:seeft_mobile/configs/importer.dart';

openRescueDialog(BuildContext context,) async {
  showDialog(
    context: context,
    builder: (context) {
      return SimpleDialog(
        contentPadding: EdgeInsets.zero,
        titlePadding: EdgeInsets.zero,
        title: Container(
          height: 550,
          child: Scaffold(
            appBar: AppBar(
              title: Text("レスキュー"), 
              centerTitle: true,
            ),
            body: Container(
              child: TextButton(
                child: Text("送信"),
                onPressed: api.postRescue
              )
            ),
          ),
        ),
      );
    },
  );
}
