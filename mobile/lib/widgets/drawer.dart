import 'package:seeft_mobile/configs/importer.dart';
import 'package:url_launcher/url_launcher.dart';

final ApplicationDrawer drawer = ApplicationDrawer();

class ApplicationDrawer {
  Widget applicationDrawer(context) {
    return Drawer(
        child: ListView(
      children: <Widget>[
        ListTile(
          title: Text("マイシフト"),
          leading: const Icon(Icons.dvr),
          onTap: () => {
            Navigator.pushNamedAndRemoveUntil(
                context, '/my_shift_page', (Route<dynamic> route) => false)
          },
        ),
        // ListTile(
        //   title: Text("全体シフト"),
        //   leading: Icon(Icons.dynamic_feed),
        //   // onTap: () => {
        //   //   Navigator.pushNamedAndRemoveUntil(
        //   //       context, '/all_shift_page', (Route<dynamic> route) => false)
        //   // },
        // ),
        ListTile(
          title: Text("全体シフト"),
          leading: Icon(Icons.dynamic_feed),
          onTap: () async {
            var url =
                "https://docs.google.com/spreadsheets/d/1KVBNRupRBFeL6Ixn9IUXzbM0wr29bpFfuaOsOeaxUQ0/edit?gid=1158066929#gid=1158066929";
            if (await canLaunchUrl(Uri.parse(url))) {
              await launchUrl(Uri.parse(url));
            } else {
              final Error error = ArgumentError('Could not launch $url');
              throw error;
            }
          },
        ),
        ListTile(
          title: Text("マニュアル一覧"),
          leading: Icon(Icons.list_alt),
          onTap: () => {
            Navigator.pushNamedAndRemoveUntil(
                context, '/manual_list_page', (Route<dynamic> route) => false)
          },
        ),
        ListTile(
          title: Text("ユーザ一覧"),
          leading: Icon(Icons.supervised_user_circle),
          onTap: () => {
            Navigator.pushNamedAndRemoveUntil(
                context, '/users_page', (Route<dynamic> route) => false)
          },
        ),
        // ListTile(
        //   title: Text("タイムスケジュール"),
        //   leading: Icon(Icons.schedule),
        //   onTap: () => {
        //     Navigator.pushNamedAndRemoveUntil(
        //         context, '/schedule_page', (Route<dynamic> route) => false)
        //   },
        // ),
        ListTile(
          title: Text("本部連絡先"),
          leading: Icon(Icons.contact_phone),
          onTap: () => {
            Navigator.pushNamedAndRemoveUntil(
                context, '/contact_page', (Route<dynamic> route) => false)
          },
        ),
        ListTile(
          title: Text("再ログイン"),
          leading: Icon(Icons.login),
          onTap: () => {
            Navigator.pushNamedAndRemoveUntil(
                context, '/signin', (Route<dynamic> route) => false)
          },
        ),
        ListTile(
          title: Text("ヘルプ"),
          leading: Icon(Icons.help),
          onTap: () async {
            var url =
            "https://docs.google.com/presentation/d/1ukPkDkkVSXWmEDY_MBOwEHPtkgm3DL64nDdLQjoTQ_0/edit#slide=id.p1";
                //"https://docs.google.com/document/d/1zCiz6rcrQuAXdVNg15MWCun2c2Babz0umPJxfJLD-Wg";
            if (await canLaunchUrl(Uri.parse(url))) {
              await launchUrl(Uri.parse(url));
            } else {
              final Error error = ArgumentError('Could not launch $url');
              throw error;
            }
          },
        ),
        ListTile(
          title: Text("その他"),
          leading: Icon(Icons.supervised_user_circle),
          onTap: () => {
            Navigator.pushNamedAndRemoveUntil(
                context, '/etc_page', (Route<dynamic> route) => false)
          },
        ),
      ],
    ));
  }
}
