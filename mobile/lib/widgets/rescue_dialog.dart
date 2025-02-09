import 'package:flutter/material.dart';
import 'package:seeft_mobile/configs/importer.dart';

rescueDialog(BuildContext context,) async {
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
              actions: <Widget>[
                IconButton(
                  onPressed: () async {
                    _launchURL(resURL);
                  },
                  icon: Icon(Icons.wrap_text),
                  color: Colors.orangeAccent[100],
                ),
              ],
            ),
            body: Container(
              child: TextButton(
                child: Text("送信")
              )
            ),
          ),
        ),
      );
    },
  );
}
