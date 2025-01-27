import 'package:flutter/material.dart';
// import 'package:flutter_web_plugins/flutter_web_plugins.dart';
import 'package:seeft_mobile/widgets/first_jump_selector.dart';
import 'package:flutter/foundation.dart' show kIsWeb; // kIsWeb を使う 追加

void main() {
  // if (kIsWeb){
  //   //web環境のみで動作
  //   setUrlStrategy(PathUrlStrategy());
  // }
  runApp(FirstJumpSelector());
}
