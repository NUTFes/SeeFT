import 'dart:developer';

import 'package:seeft_mobile/configs/importer.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:http/http.dart' as http;
import 'package:url_launcher/url_launcher.dart';

class UsersPage extends StatefulWidget {
  @override
  _UsersPageState createState() => _UsersPageState();
}

class _UsersPageState extends State<UsersPage> {
  int usersLength = 0;
  int isSelectedValue = 1;
  List<String> bureauNameList = [];
  List<String> bureauIdList = [];

  @override
  void initState() {
    super.initState();

    var bureauId;
    getBureauData().then((value) async {
      bureauId = await value;
      for (int i = 0; i < bureauId.toString().length; i++) {
        bureauNameList.add(bureauId[i]["bureau"].toString());
        bureauIdList.add(bureauId[i]["id"].toString());
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    final Size size = MediaQuery.of(context).size;
    return Scaffold(
      appBar: AppBar(
        title: const Text('ユーザ一覧'),
        actions: <Widget>[],
        // debug
      ),
      drawer: drawer.applicationDrawer(context),
      body: FutureBuilder(
        future: getData(),
        builder: (ctx, snapshot) {
          if (snapshot.connectionState == AsyncSnapshot.waiting()) {
            logger.w("message");
          }
          if (!snapshot.hasData) {
            const Center(child: CircularProgressIndicator());
          }
          var userList = snapshot.data as List<dynamic>?;
          if (userList == null) {
            return const Center(child: CircularProgressIndicator());
          }
          return Container(
            padding: const EdgeInsets.all(40.0),
            child: Column(
              children: <Widget>[
                DropdownButton(
                  items: bureauNameList.asMap().entries.map((bureau) {
                    int index = bureau.key;
                    return DropdownMenuItem(
                      value: index + 1,
                      child: Text(bureau.value),
                    );
                  }).toList(),
                  value: isSelectedValue,
                  onChanged: (value) {
                    setState(() {
                      isSelectedValue = value as int;
                    });
                  },
                ),
                Flexible(
                  child: ListView.builder(
                    itemCount: usersLength,
                    itemBuilder: (BuildContext context, int index) {
                      if (userList[index]["bureauID"] == isSelectedValue) {
                        return Container(
                            height: 40,
                            child: _manualItem(snapshot.data, index, context));
                      } else {
                        return SizedBox.shrink();
                      }
                    },
                  ),
                ),
              ],
            ),
          );
        },
      ),
    );
  }

  Widget _manualItem(var users, index, context) {
    return Container(
      decoration: new BoxDecoration(
          border:
              new Border(bottom: BorderSide(width: 1.0, color: Colors.grey))),
      child: ListTile(
        title: Text(
          users[index]["name"].toString(),
          style: TextStyle(color: Colors.black, fontSize: 14.0),
        ),
        onTap: () async {
          userInfoModal(users, index);
        },
      ),
    );
  }

  Future getBureauData() async {
    try {
      var res = await api.getBureausId();
      usersLength = res.length;
      return res;
    } catch (err) {
      logger.e('don`t response. error message: $err');
    }
  }

  Future getData() async {
    try {
      var res = await api.getUsers();
      usersLength = res.length;
      return res;
    } catch (err) {
      logger.e('don`t response. error message: $err');
    }
  }

  Future userInfoModal(users, index) async {
    var value = await showDialog(
      context: context,
      builder: (BuildContext context) => new AlertDialog(
        title: new Text('ユーザ情報'),
        content: new Container(
          height: 200,
          width: 300,
          child: ListView(children: <Widget>[
            ListTile(
              leading: Icon(Icons.person),
              title: Text("名前"),
              subtitle: Text(users[index]["name"]),
            ),
            ListTile(
              leading: Icon(Icons.email),
              title: Text("メールアドレス"),
              subtitle: Text(users[index]["mail"]),
            ),
            ListTile(
              leading: Icon(Icons.phone),
              title: Text("緊急連絡先"),
              subtitle: Text(users[index]["tel"]),
            )
          ]),
        ),
        actions: <Widget>[
          new SimpleDialogOption(
            child: new Text('OK'),
            onPressed: () => Navigator.of(context).pop(),
          )
        ],
      ),
    );
    return value;
  }
}
