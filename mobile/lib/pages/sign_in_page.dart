import 'package:seeft_mobile/configs/importer.dart';
import 'package:seeft_mobile/widgets/app_bar.dart';

class SignInPage extends StatefulWidget {
  const SignInPage({super.key});

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
      var resRoleId = res["roleID"];
      if (resId != null) {
        await store.setUserID(resId);
        await store.setRoleID(resRoleId);
        // userIdをstoreにset出来てるか確認
        var userID = await store.getUserID();
        setState(() {
          infoText = "Your ID : $userID";
        });
        Navigator.pushNamedAndRemoveUntil(
            context,
            '/layout',
            (Route<dynamic> route) => false);
      } else {
        setState(() {
          ScaffoldMessenger.of(context).showSnackBar(const SnackBar(
            content: const Text('学籍番号もしくはパスワードが違います'),
            backgroundColor: AppColors.error,
          ));
        });
      }
    } catch (e) {
      setState(() {
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(
          content: const Text('学籍番号もしくはパスワードが違います'),
          backgroundColor: AppColors.error,
        ));
      });
    }
  }

  /// 画面描写
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: const CustomAppBar(
        title: 'ログイン',
      ),
      backgroundColor: AppColors.base,
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
                            cursorColor: AppColors.main, 
                            style: const TextStyle(
                              color: AppColors.textBlack,
                              fontSize: AppFontSizes.md,
                            ),
                            decoration: InputDecoration(
                              // 枠線の色を指定
                              border: OutlineInputBorder(
                                borderRadius: BorderRadius.circular(8.0),
                                borderSide: const BorderSide(
                                  color: AppColors.grayDark,
                                  width: 1.0,
                                ),
                              ),
                              // 有効時の枠線の色を指定
                              enabledBorder: OutlineInputBorder(
                                borderRadius: BorderRadius.circular(8.0),
                                borderSide: const BorderSide(
                                  color: AppColors.grayDark,
                                  width: 1.0,
                                ),
                              ),
                              // フォーカス時の枠線の色を指定
                              focusedBorder: OutlineInputBorder(
                                borderRadius: BorderRadius.circular(8.0),
                                borderSide: const BorderSide(
                                  color: AppColors.main,
                                  width: 2.0,
                                ),
                              ),
                              // エラー時の枠線の色を指定
                              errorBorder: OutlineInputBorder(
                                borderRadius: BorderRadius.circular(8.0),
                                borderSide: const BorderSide(
                                  color: AppColors.error,
                                  width: 1.0,
                                ),
                              ),
                              // フォーカス時のエラー枠線の色を指定
                              focusedErrorBorder: OutlineInputBorder(
                                borderRadius: BorderRadius.circular(8.0),
                                borderSide: const BorderSide(
                                  color: AppColors.error,
                                  width: 2.0,
                                ),
                              ),
                              // ラベルのテキスト
                              labelText: '学籍番号',
                              // ラベルのスタイル
                              labelStyle: const TextStyle(
                                color: AppColors.grayDark,
                                fontSize: AppFontSizes.md,
                              ),
                              // フローティングラベルのスタイル
                              floatingLabelStyle: WidgetStateTextStyle.resolveWith(
                                (Set<WidgetState> states) {
                                  if (states.contains(WidgetState.error)) {
                                    return const TextStyle(
                                      color: AppColors.error,
                                      fontSize: AppFontSizes.md,
                                    );
                                  }
                                  if (states.contains(WidgetState.focused)) {
                                    return const TextStyle(
                                      color: AppColors.main,
                                      fontSize: AppFontSizes.md,
                                    );
                                  }
                                  return const TextStyle(
                                    color: AppColors.grayDark,
                                    fontSize: AppFontSizes.md,
                                  );
                                },
                              ),
                              // プレースホルダーのテキスト
                              hintText: '例) 12345678',
                              // プレースホルダーのスタイル
                              hintStyle: const TextStyle(
                                color: AppColors.grayDark,
                                fontSize: AppFontSizes.md,
                              ),
                            ),
                            onChanged: (String value) {
                              setState(() {
                                studentNumber = value;
                              });
                            },
                          ),
                          const SizedBox(height: 24.0),
                          TextField(
                            obscureText: _isObscure,
                            cursorColor: AppColors.main,
                            style: const TextStyle(
                              color: AppColors.textBlack,
                              fontSize: AppFontSizes.md,
                            ),
                            decoration: InputDecoration(
                              // 枠線の色を指定
                              border: OutlineInputBorder(
                                borderRadius: BorderRadius.circular(8.0),
                                borderSide: const BorderSide(
                                  color: AppColors.grayDark,
                                  width: 1.0,
                                ),
                              ),
                              // 有効時の枠線の色を指定
                              enabledBorder: OutlineInputBorder(
                                borderRadius: BorderRadius.circular(8.0),
                                borderSide: const BorderSide(
                                  color: AppColors.grayDark,
                                  width: 1.0,
                                ),
                              ),
                              // フォーカス時の枠線の色を指定
                              focusedBorder: OutlineInputBorder(
                                borderRadius: BorderRadius.circular(8.0),
                                borderSide: const BorderSide(
                                  color: AppColors.main,
                                  width: 2.0,
                                ),
                              ),
                              // エラー時の枠線の色を指定
                              errorBorder: OutlineInputBorder(
                                borderRadius: BorderRadius.circular(8.0),
                                borderSide: const BorderSide(
                                  color: AppColors.error,
                                  width: 1.0,
                                ),
                              ),
                              // フォーカス時のエラー枠線の色を指定
                              focusedErrorBorder: OutlineInputBorder(
                                borderRadius: BorderRadius.circular(8.0),
                                borderSide: const BorderSide(
                                  color: AppColors.error,
                                  width: 2.0,
                                ),
                              ),
                              // ラベルのテキスト
                              labelText: 'パスワード',
                              // ラベルのスタイル
                              labelStyle: const TextStyle(
                                color: AppColors.grayDark,
                                fontSize: AppFontSizes.md,
                              ),
                              // フローティングラベルのスタイル
                              floatingLabelStyle: WidgetStateTextStyle.resolveWith(
                                (Set<WidgetState> states) {
                                  if (states.contains(WidgetState.error)) {
                                    return const TextStyle(
                                      color: AppColors.error,
                                      fontSize: AppFontSizes.md,
                                    );
                                  }
                                  if (states.contains(WidgetState.focused)) {
                                    return const TextStyle(
                                      color: AppColors.main,
                                      fontSize: AppFontSizes.md,
                                    );
                                  }
                                  return const TextStyle(
                                    color: AppColors.grayDark,
                                    fontSize: AppFontSizes.md,
                                  );
                                },
                              ),
                              // アイコンの表示
                              suffixIcon: IconButton(
                                // 文字の表示・非表示でアイコンを変える
                                icon: Icon(_isObscure
                                    ? Icons.visibility_off
                                    : Icons.visibility),
                                iconSize: 24.0,
                                color: AppColors.grayDark,
                                // アイコンがタップされたら現在と反対の状態をセットする
                                onPressed: () {
                                  setState(() {
                                    _isObscure = !_isObscure;
                                  });
                                },
                              ),
                            ),
                            onChanged: (String value) {
                              setState(() {
                                password = value;
                              });
                            },
                          ),
                        ],
                      ),
                      const SizedBox(height: 48.0),
                      Container(
                        width: double.infinity,
                        height: 54.0,
                        decoration: BoxDecoration(
                          borderRadius: BorderRadius.circular(8.0),
                        ),
                        child: ElevatedButton(
                          style: ElevatedButton.styleFrom(
                            backgroundColor: AppColors.main,
                            padding: const EdgeInsets.symmetric(
                              vertical: 16.0,
                              horizontal: 24.0,
                            ),
                          ),
                          onPressed: () async {
                            _signIn();
                          },
                          child: const Text(
                            'ログイン',
                            style: TextStyle(
                              fontSize: AppFontSizes.sm,
                              color: AppColors.textWhite,
                            ),
                          ),
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
