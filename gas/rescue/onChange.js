// スプシ編集時にDBを更新する関数
function onChange(e) {
  // スクリプト全体で共通のロック
  const lock = LockService.getScriptLock();
  // const lock = LockService.getDocumentLock(); // 対象ドキュメントに紐づくロック
  // const lock = LockService.getUserLock(); // 実行ユーザーに紐づくロック

  try {
    // ロックを試みる。最大待機時間を設定 (例: 30秒)
    // この時間は、API呼び出しにかかる最大時間などを考慮して調整
    lock.waitLock(30000);

    // プロパティストアからベースURLを取得
    const properties = PropertiesService.getScriptProperties();
    const baseUrl = properties.getProperty("API_BASE_URL");

    // --- ここからクリティカルセクション ---
    const sheet = e.source.getActiveSheet();

    const range = sheet.getActiveRange();
    const values = range.getValues(); 
    const startRow = range.getRow();
    const lastRow = range.getLastRow();

    // ===== 変更内容構築 =====
    const changes = [];
    
    // 問題の種類から送信先のエンドポイントを設定
    let type;
    let statusCol;    // status(対応状況)が記入されている列番号(0スタート)
    let resCol;       // response(返答)が記入されている列番号(0スタート)
    const idCol = 0;  // id(対応番号)が記入されている列番号(0スタート)
    const sheetName = sheet.getName();  // シート名取得

    // シート名から問題の種類を取得
    if(sheetName.includes("トラブル")) {
      type = "trouble";
      statusCol = 10; // K列目
      resCol = 11;    // L列目
    } else if(sheetName.includes("質問")) {
        type = "question";
        statusCol = 8; // I列目
        resCol = 9;    // J列目
    }else if(sheetName.includes("人が来ない")) {
      type = "shorthanded";
      statusCol = 10; // K列目
      resCol = 11;    // L列目
    }else return // シート名がこれ以外の場合は何もしない

    // 編集された範囲の行を一括取得する
    const rows = sheet.getRange(startRow, idCol + 1, (lastRow - startRow) + 1, resCol + 1).getValues();

    // リクエスト
    for (let i = 0; i < values.length; i++) {
      // 対応番号を取得
      const id = String(rows[i][idCol]).trim();
      // 対応状況を取得
      let status;
      switch(String(rows[i][statusCol]).trim()){
        case "未対応":
          status = "todo";
          break;
        case "対応中":
          status = "inProgress";
          break;
        case "対応済み":
          status = "done";
          break;
        default:
          status = "todo";
          break;
          // continue // 上記以外の場合は、スキップする
      }
      // 返答を取得
      const response = String(rows[i][resCol]).trim();

      changes.push({
        id: id,
        status: status,
        response: response
      });
      console.log(changes[i]);
    }

    for(let index in changes){
      // サーバーに送信
      const url = baseUrl + "/" + type + "-rescues/" + changes[index].id;
      const options = {
        method: "put",
        contentType: "application/json",
        payload: JSON.stringify({ 
          status: changes[index].status,
          response: changes[index].response
        })
      };
      console.log(options)

      const response = UrlFetchApp.fetch(url, options);
      Logger.log("Response from API: " + response.getContentText());
    }
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
