import 'package:flutter/material.dart';
import 'package:seeft_mobile/configs/importer.dart';
import 'package:url_launcher/url_launcher.dart';
class EtcPage extends StatefulWidget{
    @override
    _EtcPageState createState()=> _EtcPageState();
}

class _EtcPageState extends State<EtcPage>{
    @override
    Widget build(BuildContext context){

        final List<Map<String, dynamic>>items=[
            {
                'title':'操作説明',
                'icon':Icons.help_outline,
                'content':() async {
                    var url =
                    "https://docs.google.com/presentation/d/1ukPkDkkVSXWmEDY_MBOwEHPtkgm3DL64nDdLQjoTQ_0/edit#slide=id.p1";
                    if (await canLaunch(url)) {
                    await launch(url);
                    } else {
                    final Error error = ArgumentError('Could not launch $url');
                    throw error;
                    }
                },
            },
            {
                'title':'通知',
                'icon':Icons.notifications_outlined,
                'content':null,
            },
            {
                'title':'ログアウト',
                'icon':Icons.logout,
                'content': () {
                    Navigator.pushNamedAndRemoveUntil(
                        context, '/signin', (Route<dynamic> route) => false);
                },
            },

        ];

        return 
        Scaffold(
            body: Container(
                child:Column(
                    children: items.map((item){
                        return Column(
                            children:[
                                etcContent(
                                    item['title'],
                                    item['icon'],
                                    item['content'],
                                ),
                                Divider(),
                            ],
                        );
                    }).toList()
                )
            )
        );
    }
}

Widget etcContent(String title,IconData icon,VoidCallback? content) {
    return Container(
        height:72,
        child:InkWell(
            onTap: content,
            child: Padding(
                padding:EdgeInsets.symmetric(horizontal:12),
                child:Row(
                    children:<Widget>[
                        Icon(icon,size:24,color:Colors.black,),
                        SizedBox(width:12),
                        Text(
                            title,
                            style:TextStyle(
                                fontSize:AppFontSizes.md,
                                color:Colors.black
                            ),
                        ),
                        Spacer(),
                        Icon(Icons.chevron_right),
                    ]
                ),
            ),
        ),
    );
}
