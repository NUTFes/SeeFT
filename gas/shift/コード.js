const lock = LockService.getScriptLock(); // スクリプト全体で共通のロック
// const lock = LockService.getDocumentLock(); // 対象ドキュメントに紐づくロック
// const lock = LockService.getUserLock(); // 実行ユーザーに紐づくロック

// uiを取得
var ui = SpreadsheetApp.getUi();

// プロパティストアからベースURLを取得
const properties = PropertiesService.getScriptProperties();
const baseUrl = properties.getProperty("API_BASE_URL");

// スプシを開いた時に実行される関数
function onOpen() {
  // タブにメニューを追加する
  const menu = ui.createMenu("SeeFT");
  menu
    // .addItem("サイドバーを開く", "showSidebar")
    // .addSeparator()
    .addSubMenu(
      ui.createMenu("シフトをSeeFTに送信する")
        .addItem("このシート全体", "updateShifts")
        .addItem("このシートの指定した範囲のみ", "updateShiftsRange")
    )
    .addSeparator()
    .addSubMenu(
      ui.createMenu("条件付き書式を設定する")
      .addItem("全てのシートにまとめて設定(時間かかる)", "setConditionalFormatting")
      .addItem("このシートにだけ設定", "setConditionalFormattingCurrentSheet")
      .addItem("タスクシートに色を設定", "setConditionalFormattingToTaskSheet")
    )
    .addToUi();
}

// スプシ編集時にDBのシフトを更新する関数
function onChange(e) {
  try {
    // ロックを試みる。最大待機時間を設定 (例: 30秒)
    // この時間は、API呼び出しにかかる最大時間などを考慮して調整
    lock.waitLock(30000);

    // --- ここからクリティカルセクション ---
    const sheet = e.source.getActiveSheet();

    const range = sheet.getActiveRange();
    const values = range.getValues(); 
    const backGrounds = range.getBackgrounds(); // 背景色を取得

    const shiftStartRow = 2;    // 見出しを除いたシフトデータの開始行
    const shiftStartColumn = 6; // 名簿情報を除いたシフトデータの開始列
    const startRow = range.getRow();
    const lastRow = range.getLastRow();
    const startColumn = range.getColumn();

    // 編集された範囲のユーザ名を一括取得する
    const userRange = sheet.getRange(startRow, 1, (lastRow - startRow) + 1, 1)
    const userValues = userRange.getValues();

    // ===== 変更内容構築 =====
    const changes = [];
    // yearID
    const yearID = 43;
    
    // 日付
    let date;
    const sheetName = sheet.getName();  // シート名取得（dateIDに対応）
    if(sheetName.includes("準備日")) date = "準備日"
    else if(sheetName.includes("1日目")) date = "1日目"
    else if(sheetName.includes("2日目")) date = "2日目"
    else if(sheetName.includes("片付け日")) date = "片付け日"
    else return // シート名がこれ以外の場合は何もしない
    // else return ui.alert(`シート名が不適です\n修正して再実行してください`)
    
    // 天気
    let weather = "晴れ";
    if(sheetName.includes("雨")) weather = "雨"

    // シフト
    for (let i = 0; i < values.length; i++) {
      // A列（1列目）からユーザー名を取得
      // const userName = String(sheet.getRange(startRow + i, 1).getValue()).trim(); // ユーザ名
      const userName = String(userValues[i]).trim();
      for (let j = 0; j < values[i].length; j++) {
        // timeID(暫定で6:00が25になるようにする)
        const timeID = (startColumn - shiftStartColumn) + j + 24;
        // タスク名(セルの背景色が黒色の場合は'NG'にする. 「休憩」は空白として送信する)
        const taskName = backGrounds[i][j] != '#000000'
          ?(values[i][j] != '休憩'
            ? String(values[i][j]).trim()
            : '')
          :'NG';
        changes.push({
          yearID: yearID,     // yearID
          timeID: timeID,     // timeID
          date: date,         // 日付
          weather: weather,   // 天気
          userName: userName, // ユーザー名
          taskName: taskName  // タスク名
        });
        console.log(changes[j - 1]);
      }
    }

    // サーバーに送信
    const url = baseUrl + "/api/update_shifts";
    const options = {
      method: "post",
      contentType: "application/json",
      payload: JSON.stringify({ changes: changes })
    };
    console.log(options)

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

// 現在のスプシをDBに反映する関数
function updateShifts() {
  try {
    // ロックを試みる。最大待機時間を設定 (例: 30秒)
    // この時間は、API呼び出しにかかる最大時間などを考慮して調整
    lock.waitLock(30000);

    // --- ここからクリティカルセクション ---
    const sheet = SpreadsheetApp.getActiveSpreadsheet().getActiveSheet();
    const sheetName = sheet.getName();  // シート名取得（dateIDとweatherIDに対応）

    // ダイアログで送信内容を確認
    const confirm = ui.alert(
      `以下のシートのシフトをSeeFTに送信してよろしいですか？\n【 ` + sheetName + ` 】`,
      ui.ButtonSet.OK_CANCEL);
    if (confirm === ui.Button.CANCEL) {
      Logger.log("送信がキャンセルされました");
      return 
    }

    const range = sheet.getDataRange();
    const values = range.getValues(); 
    const backGrounds = range.getBackgrounds(); // 背景色を取得
    const shiftStartRow = 2;    // 見出しを除いたシフトデータの開始行
    const shiftStartColumn = 6; // 名簿情報を除いたシフトデータの開始列
    const shiftLastColumn = 77; // シフトデータの最終列(最初が0)

    // ===== 変更内容構築 =====
    const changes = [];
    // yearID
    const yearID = 43;

    // 日付
    let date;
    if(sheetName.includes("準備日")) date = "準備日"
    else if(sheetName.includes("1日目")) date = "1日目"
    else if(sheetName.includes("2日目")) date = "2日目"
    else if(sheetName.includes("片付け日")) date = "片付け日"
    else return ui.alert(`シート名が不適です\n修正して再実行してください`)
    
    // 天気
    let weather = "晴れ";
    if(sheetName.includes("雨")) weather = "雨"
    
    // シフト
    for (let i = shiftStartRow; i < values.length; i++) {
      // ユーザ名
      const userName = String(values[i][0]).trim();
      // for (let j = shiftStartColumn; j < shiftLastColumn; j++) {
      for (let j = shiftStartColumn; j <= shiftLastColumn; j++) {
        // timeID(暫定で6:00が25になるようにする)
        const timeID = j - (shiftStartColumn - 1) + 24;
        // タスク名(セルの背景色が黒色の場合は'NG'にする. 「休憩」は空白として送信する)
        const taskName = backGrounds[i][j] != '#000000'
          ?(values[i][j] != '休憩'
            ? String(values[i][j]).trim()
            : '')
          :'NG';
        changes.push({
          yearID: yearID,     // yearID
          timeID: timeID,     // timeID
          date: date,         // 日付
          weather: weather,   // 天気
          userName: userName, // ユーザー名
          taskName: taskName  // タスク名
        });
        Logger.log(changes[changes.length - 1]);
      }
    }

    // サーバーに送信
    const url = baseUrl + "/api/update_shifts";
    const options = {
      method: "post",
      contentType: "application/json",
      payload: JSON.stringify({ changes: changes })
    };

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


// 現在のスプシの指定した開始行~終了行の範囲のシフトをDBに反映する関数
function updateShiftsRange() {
  try {
    // ロックを試みる。最大待機時間を設定 (例: 30秒)
    // この時間は、API呼び出しにかかる最大時間などを考慮して調整
    lock.waitLock(30000);

    // --- ここからクリティカルセクション ---
    const sheet = SpreadsheetApp.getActiveSpreadsheet().getActiveSheet();
    const sheetName = sheet.getName();  // シート名取得（dateIDとweatherIDに対応）

    // ダイアログで送信内容を確認
    const confirm = ui.alert(
      `以下のシートのシフトをSeeFTに送信してよろしいですか？\n【 ` + sheetName + ` 】`,
      ui.ButtonSet.OK_CANCEL);
    if (confirm === ui.Button.CANCEL) {
      Logger.log("送信がキャンセルされました");
      return 
    }
    const dataRange = sheet.getDataRange();
    const firstRow = 3;
    const lastRow = dataRange.getLastRow();

    // 範囲を指定するダイアログを表示
    const selectedRange = getRowRangeFromUser(firstRow, lastRow);
    if(!selectedRange) {
      Logger.log("送信がキャンセルされました");
      return 
    }
    const startRow = selectedRange.startRow;  // 指定された開始行
    const endRow = selectedRange.endRow;      // 指定された終了行
    const shiftStartColumn = 6; // 名簿情報を除いたシフトデータの開始列(最初が0)
    const shiftLastColumn = 77; // シフトデータの最終列(最初が0)

    // 指定された範囲のデータを取得
    const range = sheet.getRange(startRow, 1, (endRow- startRow) + 1, shiftLastColumn + 1)
    const values = range.getValues(); 
    const backGrounds = range.getBackgrounds(); // 背景色を取得
    

    // ===== 変更内容構築 =====
    const changes = [];
    // yearID
    const yearID = 43;

    // 日付
    let date;
    if(sheetName.includes("準備日")) date = "準備日"
    else if(sheetName.includes("1日目")) date = "1日目"
    else if(sheetName.includes("2日目")) date = "2日目"
    else if(sheetName.includes("片付け日")) date = "片付け日"
    else return ui.alert(`シート名が不適です\n修正して再実行してください`)
    
    // 天気
    let weather = "晴れ";
    if(sheetName.includes("雨")) weather = "雨"
    
    // シフト
    for (let i = 0; i < values.length; i++) {
      // ユーザ名
      const userName = String(values[i][0]).trim();
      for (let j = shiftStartColumn; j <= shiftLastColumn; j++) {
        // timeID(暫定で6:00が25になるようにする)
        const timeID = j - (shiftStartColumn - 1) + 24;
        // タスク名(セルの背景色が黒色の場合は'NG'にする. 「休憩」は空白として送信する)
        const taskName = backGrounds[i][j] != '#000000'
          ?(values[i][j] != '休憩'
            ? String(values[i][j]).trim()
            : '')
          :'NG';
        changes.push({
          yearID: yearID,     // yearID
          timeID: timeID,     // timeID
          date: date,         // 日付
          weather: weather,   // 天気
          userName: userName, // ユーザー名
          taskName: taskName  // タスク名
        });
        Logger.log(changes[changes.length - 1]);
      }
    }

    // サーバーに送信
    const url = baseUrl + "/api/update_shifts";
    const options = {
      method: "post",
      contentType: "application/json",
      payload: JSON.stringify({ changes: changes })
    };

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

// 行を入力させるダイアログを表示して、入力された範囲を返す関数(引数で、データ範囲の開始と終了行をもらう)
function getRowRangeFromUser(firstRow, lastRow) {
  // const ui = SpreadsheetApp.getUi();
  
  // ユーザーに入力を促すプロンプトを表示
  const response = ui.prompt(
    `行の範囲を入力してください\n
    開始行,終了行 の形式で入力してください（例: 5,12）\n
    有効な行の範囲は ${firstRow} から ${lastRow} です`,
    ui.ButtonSet.OK_CANCEL
  );
  
  // OKボタンが押されたかを確認
  if (response.getSelectedButton() !== ui.Button.OK) {
    ui.alert('キャンセルされました。');
    return;
  }

  const input = response.getResponseText().trim();
  const parts = input.split(',');

  // 入力チェック
  if (parts.length !== 2) {
    ui.alert('形式が正しくありません。「開始行,終了行」の形式で入力してください。');
    return;
  }

  const startRow = parseInt(parts[0].trim(), 10);
  const endRow = parseInt(parts[1].trim(), 10);

  if (isNaN(startRow) || isNaN(endRow) || startRow < 1 || endRow < 1 || startRow > endRow || startRow < firstRow || endRow > lastRow) {
    ui.alert('有効な数字で「開始行,終了行」を入力してください。');
    return;
  }

  // 結果の出力（ここで処理を行う）
  ui.alert(`開始行: ${startRow}\n終了行: ${endRow}`);

  // 必要であれば返す
  return { startRow, endRow };
}

function resetLock() {
  lock.releaseLock(); // ロックを解放する
}