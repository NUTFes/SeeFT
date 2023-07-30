import 'package:shared_preferences/shared_preferences.dart';
import 'package:seeft_mobile/configs/importer.dart';
import 'package:seeft_mobile/utils/logger.dart';

final store = PermanentStore.getInstance();

class PermanentStore {
  static PermanentStore _instance = PermanentStore();

  static getInstance() {
    return _instance;
  }

  Future<void> setUserID(resId) async {
    try {
      final SharedPreferences prefs = await SharedPreferences.getInstance();
      prefs.setInt('userID', resId);
    } catch (e) {
      print('Error while fetching SharedPreferences(setUserID): $e');
    }
  }

  Future<int> getUserID() async {
    final SharedPreferences prefs = await SharedPreferences.getInstance();
    final userID = prefs.getInt('userID') ?? 0;
    logger.d('load parmeanent store: $userID');
    return userID;
  }

  Future<bool> isUserID() async {
    try {
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
    } catch (e) {
      print('Error while fetching SharedPreferences(isUserID): $e');
      return false;
    }
  }
}
