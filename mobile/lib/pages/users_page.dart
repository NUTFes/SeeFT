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

  @override
  void initState() {
    super.initState();
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
          var userList = snapshot.data;
          return Container(
            padding: const EdgeInsets.all(40.0),
            child: Column(
              children: <Widget>[
                Flexible(
                  child: ListView.builder(
                    itemCount: usersLength,
                    itemBuilder: (BuildContext context, int index) {
                      return Container(
                          height: 40,
                          child: _manualItem(snapshot.data, index, context));
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
        onTap: () async {},
      ),
    );
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
}
