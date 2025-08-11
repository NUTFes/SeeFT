Constant constant = Constant();

// 環境変数からAPIのベースURLを取得
const _apiBaseUrl = String.fromEnvironment('API_BASE_URL', defaultValue: 'http://localhost:1234');
// 環境変数から委員長の名前と電話番号を取得
const _chairpersonName = String.fromEnvironment('CHAIRPERSON_NAME', defaultValue: '委員長');
const _chairpersonPhoneNumber = String.fromEnvironment('CHAIRPERSON_PHONE_NUMBER', defaultValue: '00011112222');

class Constant {
  final String appName = "SeeFT";
  // final String apiUrl = "https://seeft-api.nutfes.net";
  // final String apiUrl = "http://localhost:1234";
  final String apiUrl = _apiBaseUrl;
  final String chairpersonName = _chairpersonName;
  final String chairpersonPhoneNumber = _chairpersonPhoneNumber;
}
