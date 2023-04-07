import 'package:flutter/material.dart';
import 'package:flutter_web_plugins/flutter_web_plugins.dart';
import 'package:seeft_mobile/widgets/first_jump_selector.dart';

void main() {
  setUrlStrategy(PathUrlStrategy());
  runApp(FirstJumpSelector());
}
