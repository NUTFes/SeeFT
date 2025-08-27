Constant constant = Constant();

// 環境変数からAPIのベースURLを取得
const _apiBaseUrl = String.fromEnvironment('API_BASE_URL', defaultValue: 'http://localhost:1234');
// 環境変数から委員長の名前と電話番号を取得
const _chairpersonName = String.fromEnvironment('CHAIRPERSON_NAME', defaultValue: '委員長');
const _chairpersonPhoneNumber = String.fromEnvironment('CHAIRPERSON_PHONE_NUMBER', defaultValue: '00011112222');
// 環境変数からSeeFTの操作説明のURLを取得
const _seeftInstructionsUrl = String.fromEnvironment('SEEFT_INSTRUCTIONS_URL', defaultValue: 'https://example.com');
// 環境変数から全体シフトのURLを取得
const _wholeShiftUrl = String.fromEnvironment('WHOLE_SHIFT_URL', defaultValue: 'https://example.com');
// 環境変数から、各日程の日付を取得
const _nutfesPreparationDay = String.fromEnvironment('NUTFES_PREPARATION_DAY', defaultValue: '2025-09-12');
const _nutfesDay1 = String.fromEnvironment('NUTFES_DAY1', defaultValue: '2025-09-13');
const _nutfesDay2 = String.fromEnvironment('NUTFES_DAY2', defaultValue: '2025-09-14');
const _nutfesTidyingUpDay = String.fromEnvironment('NUTFES_TIDYING_UP_DAY', defaultValue: '2025-09-15');

class Constant {
  final String appName = "SeeFT";
  final String apiUrl = _apiBaseUrl;
  final String chairpersonName = _chairpersonName;
  final String chairpersonPhoneNumber = _chairpersonPhoneNumber;
  final String seeftInstructionsUrl = _seeftInstructionsUrl;
  final String wholeShiftUrl = _wholeShiftUrl;
  final String nutfesPreparationDay = _nutfesPreparationDay;
  final String nutfesDay1 = _nutfesDay1;
  final String nutfesDay2 = _nutfesDay2;
  final String nutfesTidyingUpDay = _nutfesTidyingUpDay;
}
