import 'package:shared_preferences/shared_preferences.dart';
import 'package:seeft_mobile/configs/importer.dart';
import 'package:hive_flutter/hive_flutter.dart';

// ユーザID等の保存のためのshared_preferencesのインスタンス
final PermanentStore store = PermanentStore.getInstance();

// シフトカードのデータを保存するためのHiveのBox
late final Box shiftCardBox;

class PermanentStore {
  static PermanentStore _instance = PermanentStore();

  static getInstance() {
    return _instance;
  }

  Future<void> setUserID(resId) async {
    final SharedPreferences prefs = await SharedPreferences.getInstance();
    await prefs.setInt('userID', resId);
  }

  Future<int> getUserID() async {
    try {
      final SharedPreferences prefs = await SharedPreferences.getInstance();
      final userID = prefs.getInt('userID') ?? 0;
      logger.d('load parmeanent store: $userID');
      return userID;
    } catch (e) {
      print('Error while fetching SharedPreferences(getUserID): $e');
      return 0;
    }
  }

  Future<void> setRoleID(resId) async {
    final SharedPreferences prefs = await SharedPreferences.getInstance();
    await prefs.setInt('roleID', resId);
  }

  Future<int> getRoleID() async {
    final SharedPreferences prefs = await SharedPreferences.getInstance();
    final roleID = prefs.getInt('roleID') ?? 0;
    logger.d('load parmeanent store: $roleID');
    return roleID;
  }

  Future<bool> isUserID() async {
    final SharedPreferences prefs = await SharedPreferences.getInstance();
    var userID = prefs.getInt('userID') ?? 0;
    var isUserID = true;
    if (userID == 0) {
      isUserID = false;
    }
    logger.d('load isUserID parmeanent store: $isUserID');
    return isUserID;
  }
}

// Hiveの初期化とBoxのオープン
Future<void> initHive() async {
  logger.i('Hiveの初期化を開始します。');
  // Hiveの初期化
  await Hive.initFlutter();
  
  // シフトカードのデータを保存するためのBoxを開く
  shiftCardBox = await Hive.openBox('shiftCardData');
  logger.i('Hiveの初期化が完了しました。');
}