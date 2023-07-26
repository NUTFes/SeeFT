// import 'dart:ffi';

import 'package:seeft_mobile/configs/importer.dart';

class SignInPage extends StatefulWidget {
  @override
  _SignInPageState createState() => _SignInPageState();
}

/// 利用者登録のページ
class _SignInPageState extends State<SignInPage> {
  bool _isObscure = true;
  String studentNumber = '';
  String password = '';
  String infoText = '';

  _signIn() async {
    try {
      var res = await api.signIn(studentNumber, password);
      var resId = res["id"];

      if (resId != null) {
        await store.setUserID(resId);

        // userIdをstoreにset出来てるか確認
        var userID = await store.getUserID();
        setState(() {
          infoText = "Your ID : ${userID}";
        });
        Navigator.pushNamedAndRemoveUntil(
//        context, '/wait_page', (Route<dynamic> route) => false);
            context,
            '/my_shift_page',
            (Route<dynamic> route) => false);
      } else {
        setState(() {
          ScaffoldMessenger.of(context).showSnackBar(const SnackBar(
            content: Text('メールアドレスもしくはパスワードが違います'),
            backgroundColor: Colors.redAccent,
          ));
        });
      }
    } catch (e) {
      print(e);
      setState(() {
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(
          content: Text('メールアドレスもしくはパスワードが違います'),
          backgroundColor: Colors.redAccent,
        ));
      });
    }
  }

  /// 画面描写
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('ログインページ'),
      ),
      body: Center(
        child: LayoutBuilder(
          builder: (context, constraints) {
            return SingleChildScrollView(
              child: ConstrainedBox(
                constraints: BoxConstraints(minHeight: constraints.maxHeight),
                child: Padding(
                  padding: const EdgeInsets.all(32.0),
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Column(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          TextField(
                            decoration: InputDecoration(
                              border: OutlineInputBorder(
                                borderRadius: BorderRadius.circular(8.0),
                              ),
                              labelText: '学籍番号',
                            ),
                            onChanged: (String value) {
                              setState(() {
                                studentNumber = value;
                              });
                            },
                          ),
                          const SizedBox(height: 25.0),
                          TextField(
                            obscureText: _isObscure,
                            decoration: InputDecoration(
                                border: OutlineInputBorder(
                                  borderRadius: BorderRadius.circular(8.0),
                                ),
                                labelText: 'パスワード',
                                suffixIcon: IconButton(
                                  // 文字の表示・非表示でアイコンを変える
                                  icon: Icon(_isObscure
                                      ? Icons.visibility_off
                                      : Icons.visibility),
                                  // アイコンがタップされたら現在と反対の状態をセットする
                                  onPressed: () {
                                    setState(() {
                                      _isObscure = !_isObscure;
                                    });
                                  },
                                )),
                            onChanged: (String value) {
                              setState(() {
                                password = value;
                              });
                            },
                          ),
                        ],
                      ),
                      const SizedBox(height: 50.0),
                      Container(
                        width: double.infinity,
                        height: 54.0,
                        decoration: BoxDecoration(
                          borderRadius: BorderRadius.circular(8.0),
                        ),
                        child: ElevatedButton(
                          child: const Text('ログイン'),
                          style: ElevatedButton.styleFrom(
                            primary: Colors.teal,
                            onPrimary: Colors.white,
                          ),
                          onPressed: () async {
                            _signIn();
                          },
                        ),
                      ),
                      Container(
                        padding: EdgeInsets.all(8),
                        // メッセージ表示
                        child: Text(infoText),
                      ),
                    ],
                  ),
                ),
              ),
            );
          },
        ),
      ),
    );
  }
}
