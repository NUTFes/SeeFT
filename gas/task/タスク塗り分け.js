function 塗り分けスケジュール() {
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const srcSheet = ss.getSheetByName("44thTask_準備日");
  const destSheet = ss.getSheetByName("準備日_ガントチャート_晴れ");

  if (!srcSheet) throw new Error("シート '44thTask' が見つかりません。");
  if (!destSheet) throw new Error("シート 'ガントチャート' が見つかりません。");

  const srcDataRange = srcSheet.getDataRange();
  const srcData = srcDataRange.getValues();

  // ▼ 局名と色のマップ
  const bureauColorMap = {
    "総務": "#11734b",
    "企画": "#0a53a8",
    "渉外": "#bfe1f6",
    "情報": "#ffc8aa",
    "制作": "#ffe5a0",
    "財務": "#d4edbc",
    "実行委員": "#b10202",
    "執行部": "#e6cff2"
  };

  // ▼ 時刻ヘッダーの取得（F3:BZ3）
  const timeRow = 3;
  const timeStartCol = 6; // F列
  const timeEndCol = destSheet.getLastColumn();
  const timeHeader = destSheet.getRange(timeRow, timeStartCol, 1, timeEndCol - timeStartCol + 1).getValues()[0];

  // ▼ 時刻列マップ作成
  const timeColumnMap = {};
  for (let i = 0; i < timeHeader.length; i++) {
    let val = timeHeader[i];
    let hours, minutes;

    if (val instanceof Date) {
      hours = val.getHours();
      minutes = val.getMinutes();
    } else if (typeof val === "string" && /^\d{1,2}:\d{2}/.test(val)) {
      [hours, minutes] = val.split(":").map(Number);
    } else {
      continue;
    }

    const key = `${hours}:${minutes}`;
    timeColumnMap[key] = timeStartCol + i - 1;
  }

  Logger.log("timeColumnMap keys: " + Object.keys(timeColumnMap).join(", "));

  const timezone = ss.getSpreadsheetTimeZone();

  for (let row = 1; row < srcData.length; row++) {
    const taskName     = srcData[row][0];   // A列
    const startTime    = srcData[row][2];   // C列
    const endTime      = srcData[row][3];   // D列
    const bindTime     = srcData[row][4];   // E列（拘束時間）
    const bureauRaw    = srcData[row][5];   // F列
    const level        = srcData[row][6];   // G列
    const suitable     = srcData[row][10];  // K列
    const properPeople = srcData[row][10];  // K列（人数）

    if (!taskName || !startTime || !endTime || !bureauRaw || properPeople === "") continue;

    const bureau = bureauRaw.toString().replace(/局$/, '').trim();
    const color = bureauColorMap[bureau];
    if (!color) continue;

    const getTimeKey = (d) => {
      if (!(d instanceof Date)) return null;
      return `${d.getHours()}:${d.getMinutes()}`;
    };

    const startKey = getTimeKey(startTime);
    const endKey = getTimeKey(endTime);
    const startCol = timeColumnMap[startKey];
    const endCol = timeColumnMap[endKey];

    if (startCol == null || endCol == null) continue;

    const outputRow = row + 3;

    // ▼ A〜G列に出力
    destSheet.getRange(outputRow, 1).setValue(taskName);              // A: シフト名
    destSheet.getRange(outputRow, 2).setValue(bureauRaw);             // B: 局名
    destSheet.getRange(outputRow, 3).setValue(
      startTime instanceof Date
        ? Utilities.formatDate(startTime, timezone, "HH:mm")
        : startTime
    );                                                                // C: 開始
    destSheet.getRange(outputRow, 4).setValue(
      endTime instanceof Date
        ? Utilities.formatDate(endTime, timezone, "HH:mm")
        : endTime
    );                                                                // D: 終了
    destSheet.getRange(outputRow, 5).setValue(
      bindTime instanceof Date
        ? Utilities.formatDate(bindTime, timezone, "HH:mm")
        : bindTime
    );                                                                // E: 拘束時間（HH:mm）
    destSheet.getRange(outputRow, 6).setValue(level);                 // F: レベル
    destSheet.getRange(outputRow, 7).setValue(suitable);              // G: 適正人数

    // ▼ ガントチャート塗り（右に1列、下に2行）
    for (let col = startCol; col < endCol; col++) {
      const destCell = destSheet.getRange(outputRow, col + 1); // +1列右（+1 for 1-based）
      destCell.setValue(properPeople);
      destCell.setBackground(color);
    }
  }

  SpreadsheetApp.getUi().alert("塗り分け完了！");
}