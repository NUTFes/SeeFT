// タスクごとに背景色を条件付き書式で割り当てる関数
function setConditionalFormatting() {
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const sheet1 = ss.getSheetByName("準々備日")
  const sheet2 =  ss.getSheetByName("準備日")
  const sheet3 =  ss.getSheetByName("1日目_晴れ")
  const sheet4 =  ss.getSheetByName("1日目_雨")
  const sheet5 =  ss.getSheetByName("2日目_晴れ")
  const sheet6 =  ss.getSheetByName("2日目_雨")
  const sheet7 =  ss.getSheetByName("片付け日")
  const taskSheet = ss.getSheetByName("全タスク")

  const sheets= [
    sheet1,
    sheet2,
    sheet3,
    sheet4,
    sheet5,
    sheet6,
    sheet7
  ];

  // シフトのシートごとに、書式を設定する
  for (let i = 0; i < sheets.length; i++) {
    setConditionalFormattingToSheet(sheets[i], taskSheet)
    // setBackgroundColorToSheet(sheets[i], taskSheet)
  }
}

// 現在開いているシートに対して条件付き書式を割り当てる関数
function setConditionalFormattingCurrentSheet(){
  // 現在開いているシートを取得
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const sheet = ss.getActiveSheet();
  const sheetName = sheet.getName();

  // ダイアログで送信内容を確認
  const confirm = ui.alert(
    `以下のシートに条件付き書式を設定してよろしいですか？\n【 ` + sheetName + ` 】`,
    ui.ButtonSet.OK_CANCEL);
  if (confirm === ui.Button.CANCEL) {
    // キャンセルされた場合の処理
    Logger.log("操作がキャンセルされました");
    return 
  }
  // キャンセルされなかった場合の処理
  const taskSheet = ss.getSheetByName("全タスク")    // タスクシートを取得
  setConditionalFormattingToSheet(sheet, taskSheet) // 開いているシートに条件付き書式を割り当てる
  // setBackgroundColorToSheet(sheet, taskSheet) // 開いているシートに背景色を割り当てる
}

// 指定したシートにタスクごとに背景色を割り当てる関数
function setBackgroundColorToSheet(sheet, taskSheet){
  // === 全タスクシートからタスク名とカラーコードのマッピングを作成 ===
  const mappingValues = taskSheet.getRange("A3:Q" + taskSheet.getLastRow()).getValues();
  const colorMap = {};
  mappingValues.forEach(row => {
    const taskName = row[0];  // A列
    const bgColor = row[15];  // P列
    const fontColor = getContrastTextColor(bgColor);
    if (taskName && bgColor) {
      // 例: { "りんご": "#FF0000", "バナナ": "#FFFF00" }
      colorMap[taskName] = {
        bg: bgColor.startsWith("#") ? bgColor : "#" + bgColor,
        font: fontColor
      }
    }
  });

  // === シートの対象範囲を取得 ===
  const range = sheet.getRange("G3:BZ300");
  const values = range.getValues();

  // === 背景色 & 文字色の配列を作成 ===
  const bgColors = [];
  const fontColors = [];

  for (let row = 0; row < values.length; row++) {
    const bgRow = [];
    const fontRow = [];
    for (let cell = 0; cell < values[row].length; cell++) {
      const cellValue = values[row][cell];
      if (cellValue == "") {
        bgRow.push(null);  // セルが空白の場合
        fontRow.push(null);
      } else if (colorMap[cellValue]) {
        bgRow.push(colorMap[cellValue].bg);
        fontRow.push(colorMap[cellValue].font);
      } else {
        bgRow.push("#000000"); // 一致しない場合 → 背景黒
        fontRow.push("#FFFFFF"); // 一致しない場合 → 文字白
      }
    }
    bgColors.push(bgRow);
    fontColors.push(fontRow);
  }

  // === 背景色と文字色を一括適用 ===
  range.setBackgrounds(bgColors);
  range.setFontColors(fontColors);
}

// 指定されたシートにタスクごとに背景色を条件付き書式で割り当てる関数
function setConditionalFormattingToSheet(sheet, taskSheet) {
  const values = taskSheet.getRange("A3:P" + taskSheet.getLastRow()).getValues();

  // 条件付き書式ルールを一旦リセット
  sheet.clearConditionalFormatRules();

  const rules = [];
  const range = sheet.getRange("G3:BZ300");

  for (let i = 0; i < values.length; i++) {
    const keyword = values[i][0]; // A列
    const bgColor = values[i][15]; // P列（背景色にしたいカラーコード）

    if (keyword && /^#[0-9a-fA-F]{6}$/.test(bgColor)) {
      const textColor = getContrastTextColor(bgColor);
      const rule = SpreadsheetApp.newConditionalFormatRule()
        .whenTextEqualTo(keyword)
        .setBackground(bgColor)
        .setFontColor(textColor)
        .setRanges([range])
        .build();
      rules.push(rule);
    }
  }

  // シフト希望のNG理由の書式設定(空白でないかつ、タスクに存在していない値の場合、黒塗りにする)
  const ngformula = `=NOT(OR(COUNTIF(INDIRECT("全タスク!$A$3:A"), G3), G3 = ""))`;
  const ngRule = SpreadsheetApp.newConditionalFormatRule()
    .whenFormulaSatisfied(ngformula)
    .setBackground("#000000")
    .setFontColor("#ffffff")
    .setRanges([range])
    .build();
  rules.push(ngRule);

  sheet.setConditionalFormatRules(rules);
}

// タスクシートにタスクごとに背景色を条件付き書式で割り当てる関数
function setConditionalFormattingToTaskSheet() {
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const taskSheet = ss.getSheetByName("全タスク")
  const lastRow = taskSheet.getLastRow()
  // const values = taskSheet.getRange("A3:P" + lastRow).getValues();
  const values = taskSheet.getRange("A3:P500").getValues();

  for (let i = 0; i < values.length; i++) {
    const bgColor = values[i][15]; // P列（背景色にしたいカラーコード）
    console.log(bgColor)
    if (/^#[0-9a-fA-F]{6}$/.test(bgColor)) {
      const textColor = getContrastTextColor(bgColor);
      taskSheet
        .getRange(i + 3, 16)  // P列のセルの書式を設定
        .setBackground(bgColor)
        .setFontColor(textColor)
    }else {
      taskSheet
        .getRange(i + 3, 16)
        .setBackground(null)
        .setFontColor(null)
    }
  }
}

// 色の明るさを評価して、適切な文字色を返す
function getContrastTextColor(hexColor) {
  const r = parseInt(hexColor.substr(1, 2), 16);
  const g = parseInt(hexColor.substr(3, 2), 16);
  const b = parseInt(hexColor.substr(5, 2), 16);

  // 輝度を計算（簡易式）
  const brightness = 0.299 * r + 0.587 * g + 0.114 * b;
  return brightness < 128 ? '#ffffff' : '#000000';
}