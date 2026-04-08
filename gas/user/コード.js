const lock = LockService.getScriptLock(); // スクリプト全体で共通のロック
// const lock = LockService.getDocumentLock(); // 対象ドキュメントに紐づくロック
// const lock = LockService.getUserLock(); // 実行ユーザーに紐づくロック

// uiを取得
const ui = SpreadsheetApp.getUi();

// プロパティストアからベースURLを取得
const properties = PropertiesService.getScriptProperties();
const baseUrl = properties.getProperty("API_BASE_URL");

const bureaus = [
  '執行部',
  '執行部補佐',
  '総務局',
  '企画局',
  '財務局',
  '渉外局',
  '制作局',
  '情報局'
];
const grades = [
  'B1',
  'B2',
  'B3',
  'B4',
  'M1',
  'M2',
  'D1',
  'D2',
  'D3',
  'OB',
];
const departments = [
  '未所属',
  '機械工学分野',
  '電気電子情報工学分野',
  '情報・経営システム工学分野',
  '物質生物工学分野',
  '環境社会基盤工学分野',
  '量子・原子力統合工学分野',
  '技術科学イノベーション'
];
const departmentsInCell = [
  '未所属',
  '機械',
  '電気電子情報',
  '情報経営',
  '物質生物',
  '環境社会',
  '量子・原子力統合工学分野',
  '技術科学イノベーション'
];

// スプシを開いた時に実行される関数
function onOpen() {
  // タブにメニューを追加する
  ui.createMenu("SeeFT")
    .addItem("名簿をSeeFTに反映する", "updateUsers")
    .addToUi();
}

// 現在の名簿をDBに反映する関数
function updateUsers() {
  try {
    // ロックを試みる。最大待機時間を設定 (例: 30秒)
    // この時間は、API呼び出しにかかる最大時間などを考慮して調整
    lock.waitLock(30000);

    // --- ここからクリティカルセクション ---
    const sheet = SpreadsheetApp.getActiveSpreadsheet().getSheetByName("SeeFT送信用");
    const range = sheet.getDataRange();
    const values = range.getValues(); 

    const nameCol = 0;    // 名前の列番号
    const bureauCol = 1;  // 局の列番号
    const gradeCol = 4;   // 学年の列番号
    const departCol = 5;  // 課程の列番号
    const studentNumCol = 14; // 学籍番号の列番号
    const telCol = 13;    // 電話番号の列番号


    // 変更内容構築
    const changes = [];

    for (let i = 2; i < values.length; i++) {
      // 入力チェック
      if(!bureaus.includes(String(values[i][bureauCol]).trim())) {
        return ui.alert(
          `局が不適な局員がいます\n・` + String(values[i][0]).trim() + `\n修正して再実行してください`
        );
      }
      if(!grades.includes(String(values[i][gradeCol]).trim())) {
        return ui.alert(
          `学年が不適な局員がいます\n・` + String(values[i][0]).trim() + `\n修正して再実行してください`
        );
      }
      if(!departmentsInCell.includes(String(values[i][departCol]).trim())) {
        return ui.alert(
          `課程が不適な局員がいます\n・` + String(values[i][0]).trim() + `\n修正して再実行してください`
        );
      }
      if(typeof(values[i][studentNumCol]) != 'number' || String(values[i][studentNumCol]).trim().length != 8){
        return ui.alert(
          `学籍番号が不適な局員がいます\n・` + String(values[i][0]).trim() + `\n修正して再実行してください`
        );
      }
      if(String(values[i][telCol]).trim().length > 15){
        return ui.alert(
          `電話番号が不適な局員がいます\n・` + String(values[i][0]).trim() + `\n修正して再実行してください`
        );
      }

      // 課程名をSeeFTで使用可能な値に変換
      let department;
      switch(String(values[i][departCol]).trim()) {
        case departmentsInCell[0]:
          department = departments[0];
          break;
        case departmentsInCell[1]:
          department = departments[1];
          break;
        case departmentsInCell[2]:
          department = departments[2];
          break;
        case departmentsInCell[3]:
          department = departments[3];
          break;
        case departmentsInCell[4]:
          department = departments[4];
          break;
        case departmentsInCell[5]:
          department = departments[5];
          break;
        case departmentsInCell[6]:
          department = departments[6];
          break;
        case departmentsInCell[7]:
          department = departments[7];
          break;
      }

      // 適切ならば変更内容を格納
      changes.push({
        name: String(values[i][nameCol]).trim(),          // ユーザー名
        bureau: String(values[i][bureauCol]).trim(),      // 局
        grade: String(values[i][gradeCol]).trim(),        // 学年
        department: String(department).trim(),  // 課程
        studentNumber: Number(values[i][studentNumCol]),  // 学籍番号
        tel: String(values[i][telCol]).trim()             // 電話番号
      })
      console.log(changes[i - 2])
    }

    // サーバーに送信
    const url = baseUrl + "/api/update_users";
    const options = {
      method: "post",
      contentType: "application/json",
      payload: JSON.stringify({ changes: changes })
    };
    // console.log(options.payload)

    const response = UrlFetchApp.fetch(url, options);
    Logger.log("Response from API: " + response.getContentText());
    // --- ここまでクリティカルセクション ---
  } catch (error) {
    ui.alert(`エラーが発生しました\n
      エラーを修正して再実行してください\n
      ------エラーメッセージ-----\n` + 
      error.message + 
      `\n---------------------------`);
    Logger.log("Error: " + error.message);
  } finally {
    lock.releaseLock(); // 必ずロックを解放する
  }
}
