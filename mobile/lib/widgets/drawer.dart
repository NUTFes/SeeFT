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
        ListTile(
          title: Text("全体シフト"),
          leading: Icon(Icons.dynamic_feed),
          onTap: () => {
            // Navigator.pushNamedAndRemoveUntil(
            //     context, '/all_shift_page', (Route<dynamic> route) => false)
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
          title: Text("タイムスケジュール"),
          leading: Icon(Icons.schedule),
          onTap: () => {
            Navigator.pushNamedAndRemoveUntil(
                 context, '/schedule_page', (Route<dynamic> route) => false)
          },
        ),
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
          onTap: () async{
            var url = "https://docs.google.com/document/d/1siErvaPQnut7R0wklAmuPY5drDDkl19iXZ7KgpH1TMw/edit?usp=sharing";
            if (await canLaunch(url)) {
              await launch(url);
            } else {
              final Error error = ArgumentError('Could not launch $url');
              throw error;
            }
          },
        ),
      ],
    ));
  }
}
