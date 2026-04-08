// スプシを開いた時に実行される関数
function onOpen() {
  // タブにメニューを追加する
  const ui = SpreadsheetApp.getUi();
  const menu = ui.createMenu("SeeFT");
  menu
    .addItem("タスクに色を設定する", "assignRandomColors")
    .addItem("タスクと集合場所をSeeFTに送信する", "updateTasksAndPlaces")
    .addToUi();
}

// タスクにランダムな色を割り当てる関数(シフトのスプシのセルの背景色として使う)
function assignRandomColors() {
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const sheet1 = ss.getSheetByName("44thTask_準々備日")
  const sheet2 =  ss.getSheetByName("44thTask_準備日")
  const sheet3 =  ss.getSheetByName("44thTask_当日1日目")
  const sheet4 =  ss.getSheetByName("44thTask_当日2日目")
  const sheet5 =  ss.getSheetByName("44thTask_片付け日")

  const sheets= [
    sheet1,
    sheet2,
    sheet3,
    sheet4,
    sheet5
  ];

  for (let i = 0; i < sheets.length; i++) {
    assignRandomColorsToSheet(sheets[i])
  }
}

// 指定されたシートのタスクにランダムな色を割り当てる関数
function assignRandomColorsToSheet(sheet) {
  const dataRange = sheet.getRange("A2:P" + sheet.getLastRow());
  const data = dataRange.getValues();

  for (let i = 0; i < data.length; i++) {
    const name = data[i][0]; // A列（タスク名）
    const colorCode = data[i][15]; // P列（色）

    // タスク名が入力されていて、色が未入力の場合、適当に色を割り当てる
    if (name && !colorCode) {
      // const randomColor = getRandomHexColor();
      const randomColor = getRandomBrightHexColor(100); // 最低の輝度を指定して色を生成
      sheet.getRange(i + 2, 16).setValue(randomColor); // P列に入力（i+2行目）
    }
  }
}

// ランダムな色を生成する関数
function getRandomHexColor() {
  // const hex = Math.floor(Math.random() * 0xffffff).toString(16).padStart(6, "0");
  let hex;
  do {
    hex = Math.floor(Math.random() * 0xffffff).toString(16).padStart(6, "0");
  } while (hex === "000000"); // 真っ黒を避ける
  return `#${hex}`;
}

// 暗い色を避けてランダムな色を生成する関数
function getRandomBrightHexColor(minBrightness = 100) {
  let r, g, b, brightness;

  do {
    r = Math.floor(Math.random() * 256);
    g = Math.floor(Math.random() * 256);
    b = Math.floor(Math.random() * 256);
    brightness = 0.299 * r + 0.587 * g + 0.114 * b;
  } while (brightness < minBrightness); // 暗すぎる色を除外

  return `#${r.toString(16).padStart(2, '0')}${g.toString(16).padStart(2, '0')}${b.toString(16).padStart(2, '0')}`;
}
