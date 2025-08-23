import 'package:flutter/material.dart';
import 'package:flutter_web_plugins/flutter_web_plugins.dart';
import 'package:seeft_mobile/widgets/first_jump_selector.dart';
import 'package:seeft_mobile/utils/permanent_store.dart';

Future<void> main() async {
  setUrlStrategy(PathUrlStrategy());
  await initHive(); // Hiveの初期化を行う
  runApp(FirstJumpSelector());
}
