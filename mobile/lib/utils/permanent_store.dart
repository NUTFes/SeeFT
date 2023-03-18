import 'package:shared_preferences/shared_preferences.dart';
import 'package:seeft_mobile/configs/importer.dart';
import 'package:seeft_mobile/utils/logger.dart';

final store = PermanentStore.getInstance();

class PermanentStore {
  static PermanentStore _instance = PermanentStore();

  static getInstance() {
    return _instance;
  }

  void setUserID(resId) async {
    final SharedPreferences prefs = await SharedPreferences.getInstance();
    await prefs.setInt('userID', resId);
  }

  getUserID() async {
    final SharedPreferences prefs = await SharedPreferences.getInstance();
    final userID = prefs.getInt('userID') ?? 0;
    logger.d('load parmeanent store: $userID');
    return userID;
  }

  isUserID() async {
    final SharedPreferences prefs = await SharedPreferences.getInstance();
    var userID = prefs.getInt('userID');
    var isUserID = true;
    if (userID == null) {
       isUserID = false;
    } else if (userID == 0) {
       isUserID = false;
    }
    logger.d('load isUserID parmeanent store: $isUserID');
    return isUserID;
  }
}
