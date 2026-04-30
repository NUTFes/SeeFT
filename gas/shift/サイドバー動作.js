const ps = PropertiesService.getScriptProperties();

// function onOpen() {
//   var ui = SpreadsheetApp.getUi();
//   ui.createMenu("サイドバー")
//     .addItem("サイドバーを開く", "showSidebar")
//     .addToUi();
// }

function showSidebar() {
  var html = HtmlService.createHtmlOutputFromFile("Sidebar")
    .setTitle("セル情報")
    .setWidth(300);
  SpreadsheetApp.getUi().showSidebar(html);
}

function getSelectedCellValue() {
  var sheet = SpreadsheetApp.getActiveSpreadsheet().getActiveSheet();
  var range = sheet.getActiveCell();
  return range.getA1Notation() + " = " + range.getValue(); // セルの位置と値を取得
}

// 選択中のシフトのその時間帯の割り当て人数を取得する
function getShiftMemberCount(){
  const sheet = SpreadsheetApp.getActiveSpreadsheet().getActiveSheet();
  const range = sheet.getActiveCell();
  const selectedShift = range.getValue(); // 選択中のシフト
  const shifts = sheet.getRange(3, range.getColumn(), sheet.getDataRange().getLastRow() - 2, 1).getValues(); // 選択中の列の全てのシフト
  const count = shifts.filter(shift => shift == selectedShift).length;
  // console.log(sheet.getRange(range.getLastRow() - 2, 1).getA1Notation())
  // console.log(count)
  return "現在の人数: " + count;
}

// 選択中のセル範囲とシート名を取得する関数
function getSelectedRangeAndSheetName(){
  const sheet = SpreadsheetApp.getActiveSpreadsheet().getActiveSheet();
  const range = sheet.getActiveRange();
  return sheet.getName() + range.getA1Notation();
}

// 選択中のシート名を取得する関数
function getSelectedSheetName(){
  const sheet = SpreadsheetApp.getActiveSpreadsheet().getActiveSheet();
  console.log(sheet.getName());
  return sheet.getName();
}

// 選択中のセル範囲をA1形式で取得する関数
function getSelectedRange(){
  const sheet = SpreadsheetApp.getActiveSpreadsheet().getActiveSheet();
  const range = sheet.getActiveRange();
  console.log(range.getA1Notation())
  return range.getA1Notation();
}

// 選択中のセル範囲とシート名をストアに登録する関数
function setSwappingRange(){
  const sheet = SpreadsheetApp.getActiveSpreadsheet().getActiveSheet();
  const range = sheet.getActiveRange();
  const swappingSheetName = sheet.getName();
  const swappingCellRange = range.getA1Notation();
  ps.setProperties({'SwappingSheetName': swappingSheetName, 'SwappingCellRange': swappingCellRange});
  console.log(ps.getProperties());
  return swappingSheetName + swappingCellRange;
}

// 登録されているセル範囲とシート名を取得する関数
function getSwappingRange(){
  // 登録したセル範囲
  const swappingSheetName = ps.getProperty('SwappingSheetName');
  const swappingCellRange = ps.getProperty('SwappingCellRange');

  return swappingSheetName + swappingCellRange;
}

// 登録したセル範囲と、選択中のセル範囲の値を入れ替える関数
function swappingValues(){
  console.log('セル範囲の入れ替えを開始します...');

  // 登録したセル範囲
  const swappingSheetName = ps.getProperty('SwappingSheetName');
  const swappingCellRange = ps.getProperty('SwappingCellRange');

  // セル範囲が登録されていない場合は中断する
  if(swappingSheetName == null || !swappingCellRange == null) {
    console.log('入れ替え先のセル範囲が登録されていません');
    return;
  }

  const swappingSheet = SpreadsheetApp.getActiveSpreadsheet().getSheetByName(swappingSheetName);
  const swappingRange = swappingSheet.getRange(swappingCellRange);
  const swappingValues = swappingRange.getValues(); // 登録したセル範囲の値

  // 選択中のセル範囲
  const sheet = SpreadsheetApp.getActiveSpreadsheet().getActiveSheet();
  const cell = sheet.getActiveCell();
  const range = sheet.getRange(cell.getRow(), cell.getColumn(), swappingRange.getNumRows(), swappingRange.getNumColumns());
  const values = range.getValues();
  // console.log(range.getA1Notation());
  
  // 選択中のセル範囲に登録したセル範囲の値をコピー
  range.setValues(swappingValues);
  // 登録したセル範囲に選択中のセル範囲の値をコピー
  swappingRange.setValues(values);
  console.log('セル範囲の入れ替えが完了しました');
}


// // 入れ替え先のセル範囲を可視化する関数
// トグルにするのが良いかも。
function visualizingSwapArea(){
  // 登録したセル範囲
  const swappingSheetName = ps.getProperty('SwappingSheetName');
  const swappingCellRange = ps.getProperty('SwappingCellRange');

  // 入れ替えるセル範囲が登録されていない場合は中断する
  // if(swappingSheetName == null || swappingCellRange == null) return;

  const swappingSheet = SpreadsheetApp.getActiveSpreadsheet().getSheetByName(swappingSheetName);
  const swappingRange = swappingSheet.getRange(swappingCellRange);
  // 選択中のセル範囲
  const sheet = SpreadsheetApp.getActiveSpreadsheet().getActiveSheet();
  const cell = sheet.getActiveCell();
  const range = sheet.getRange(cell.getRow(), cell.getColumn(), swappingRange.getNumRows(), swappingRange.getNumColumns());
  range.activate();
}

// function getTaskInfo() {

// }