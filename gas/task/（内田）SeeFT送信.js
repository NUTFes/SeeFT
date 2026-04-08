// スプシ編集時にDBのタスクを更新する関数
function onChange(e) {
  const lock = LockService.getScriptLock(); // スクリプト全体で共通のロック
  // const lock = LockService.getDocumentLock(); // 対象ドキュメントに紐づくロック
  // const lock = LockService.getUserLock(); // 実行ユーザーに紐づくロック

  try {
    // ロックを試みる。最大待機時間を設定 (例: 30秒)
    // この時間は、API呼び出しにかかる最大時間などを考慮して調整
    lock.waitLock(30000);

    // uiを取得
    const ui = SpreadsheetApp.getUi();

    // プロパティストアからベースURLを取得
    const properties = PropertiesService.getScriptProperties();
    const baseUrl = properties.getProperty("API_BASE_URL");

    // --- ここからクリティカルセクション ---
    const sheet = e.source.getActiveSheet();
    const sheetName = sheet.getName();

    // タスクのシートでなければキャンセル
    if(
      sheetName != "44thTask_準々備日" && 
      sheetName != "44thTask_準備日" &&
      sheetName != "44thTask_当日1日目" &&
      sheetName != "44thTask_当日2日目" &&
      sheetName != "44thTask_片付け日"
    ){
      Logger.log("タスクシートの編集ではありませんでした");
      return
    }
    const editRange = sheet.getActiveRange();

    const startRow = editRange.getRow();
    const lastRow = editRange.getLastRow();

    console.log(startRow)
    console.log(lastRow)

    // 編集された範囲のタスク名~リンク(12列目)までを一括取得する
    const range = sheet.getRange(startRow, 1, (lastRow - startRow) + 1, 13)
    const values = range.getValues();
    const ritchTextValues = range.getRichTextValues();

    const nameCol = 0;    // タスク名の列番号(1列目が0)
    const bureauCol = 5;  // 管轄局の列番号
    const placeCol = 7;   // 集合場所の列番号
    const urlCol = 12;    // URLの列番号
    const maxMemberCol = 11;  // 最大人数の列番号

    // 変更内容構築
    const changes = [];
    const yearID = 43;
    for (let i = 0; i < values.length; i++) {
      const taskName = String(values[i][nameCol]).trim()
      // タスク名が入力されていないか、列の見出しの場合はスキップ
      if(!taskName || taskName == 'シフト名') continue;
      const bureau = String(values[i][bureauCol]).trim()
      // if(!bureau || bureau == '実行委員') bureau = '執行部'  // goの方でデフォルト値設定していたため不要
      const place = String(values[i][placeCol]).trim()|| '本部(電気棟1F)'; // 集合場所(デフォルトは「本部(電気棟1F)」)
      const url = ritchTextValues[i][urlCol].getLinkUrl() ? String(ritchTextValues[i][urlCol].getLinkUrl()).trim(): '';
      const maxMember = (values[i][maxMemberCol] && typeof(values[i][maxMemberCol]) == 'number')? values[i][maxMemberCol]: 1; // 最大人数(デフォルトは「1」)

      changes.push({
        yearID: yearID,       // yearID
        taskName: taskName,   // タスク名
        bureau: bureau,       // 管轄局(デフォルトは「執行部」)
        place: place,         // 集合場所(デフォルトは「本部(電気棟1F)」)
        url: url,             // マニュアルURL
        maxMember: maxMember  // 最大人数
      });
      console.log(changes[i])
    }

    // サーバーに送信
    const url = baseUrl + "/api/update_tasks_and_places";
    const options = {
      method: "post",
      contentType: "application/json",
      payload: JSON.stringify({ changes: changes }),
      muteHttpExceptions: false
    };
    // console.log(options.payload)

    const response = UrlFetchApp.fetch(url, options);
    Logger.log("Response from API: " + response.getContentText());
    // --- ここまでクリティカルセクション ---
  } catch (error) {
    const ui = SpreadsheetApp.getUi();
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
function updateTasksAndPlaces() {
  const lock = LockService.getScriptLock(); // スクリプト全体で共通のロック
  // const lock = LockService.getDocumentLock(); // 対象ドキュメントに紐づくロック
  // const lock = LockService.getUserLock(); // 実行ユーザーに紐づくロック

  try {
    // ロックを試みる。最大待機時間を設定 (例: 30秒)
    // この時間は、API呼び出しにかかる最大時間などを考慮して調整
    lock.waitLock(30000);
    
    // uiを取得
    const ui = SpreadsheetApp.getUi();

    // プロパティストアからベースURLを取得
    const properties = PropertiesService.getScriptProperties();
    const baseUrl = properties.getProperty("API_BASE_URL");

    // --- ここからクリティカルセクション ---
    // const sheet = SpreadsheetApp.getActiveSpreadsheet().getSheetByName("タスク");
    const sheet = SpreadsheetApp.getActiveSpreadsheet().getActiveSheet();
    const sheetName = sheet.getName();  // シート名取得（dateIDとweatherIDに対応）

    // ダイアログで送信内容を確認
    const confirm = ui.alert(
      `以下のシートのタスクをSeeFTに送信してよろしいですか？\n【 ` + sheetName + ` 】`,
      ui.ButtonSet.OK_CANCEL);
    if (confirm === ui.Button.CANCEL) {
      Logger.log("送信がキャンセルされました");
      return 
    }

    // タスクのシートでなければキャンセル
    if(
      sheetName != "44thTask_準々備日" && 
      sheetName != "44thTask_準備日" &&
      sheetName != "44thTask_当日1日目" &&
      sheetName != "44thTask_当日2日目" &&
      sheetName != "44thTask_片付け日"
    ){
      Logger.log("タスクシートではないため送信をキャンセルしました");
      return ui.alert(`タスクシートではないため送信をキャンセルしました`)
    }

    const range = sheet.getDataRange();
    const values = range.getValues(); 
    const ritchTextValues = range.getRichTextValues();

    const nameCol = 0;    // タスク名の列番号
    const bureauCol = 5;  // 管轄局の列番号
    const placeCol = 7;   // 集合場所の列番号
    const urlCol = 12;    // URLの列番号
    const maxMemberCol = 11;  // 最大人数の列番号

    // 変更内容構築
    const changes = [];
    const yearID = 43;
    for (let i = 1; i < values.length; i++) {
      const taskName = String(values[i][nameCol]).trim();
      // タスク名が入力されていない場合はスキップ
      if(!taskName) continue;
      const bureau = String(values[i][bureauCol]).trim();
      // if(!bureau || bureau == '実行委員') bureau = '執行部'  // goの方でデフォルト値設定していたため不要
      const place = String(values[i][placeCol]).trim()|| '本部(電気棟1F)'; // 集合場所(デフォルトは「本部(電気棟1F)」)
      const url = ritchTextValues[i][urlCol].getLinkUrl() ? String(ritchTextValues[i][urlCol].getLinkUrl()).trim(): '';
      const maxMember = (values[i][maxMemberCol] && typeof(values[i][maxMemberCol]) == 'number')? values[i][maxMemberCol]: 1; // 最大人数(デフォルトは「1」)
      

      changes.push({
        yearID: yearID,       // yearID
        taskName: taskName,   // タスク名
        bureau: bureau,       // 管轄局(デフォルトは「執行部」)
        place: place,         // 集合場所(デフォルトは「本部(電気棟1F)」)
        url: url,             // マニュアルURL
        maxMember: maxMember  // 最大人数
      })
      console.log(changes[i - 1])
    }

    // サーバーに送信
    const url = baseUrl + "/api/update_tasks_and_places";
    const options = {
      method: "post",
      contentType: "application/json",
      payload: JSON.stringify({ changes: changes }),
      muteHttpExceptions: false
    };
    // console.log(options.payload)

    const response = UrlFetchApp.fetch(url, options);
    Logger.log("Response from API: " + response.getContentText());
    // --- ここまでクリティカルセクション ---
  } catch (error) {
    // uiを取得
    const ui = SpreadsheetApp.getUi();
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