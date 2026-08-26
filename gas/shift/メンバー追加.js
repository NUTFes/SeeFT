// 名簿連鎖(固定名簿→希望調査)に登録済みのメンバーを、名簿シートと全日程シートに追加する
const MEMBER_DAY_SHEETS = ["準々備日", "準備日", "1日目_晴れ", "1日目_雨", "2日目_晴れ", "2日目_雨", "片付け日"];

function addMemberColumns() {
  ui = ui || SpreadsheetApp.getUi();
  const ans = ui.prompt("メンバー追加", "名前、局、学年、課程、メール(任意)をカンマ区切りで入力。例: 小林 優人,企画局,B3,機械,24.y.kobayashi.nutfes@gmail.com", ui.ButtonSet.OK_CANCEL);
  if (ans.getSelectedButton() !== ui.Button.OK) return;
  const parts = ans.getResponseText().split(/[,、]/).map(function(s){ return s.trim(); });
  if (parts.length < 3 || !parts[0]) { ui.alert("入力形式が違います。例: 小林 優人,企画局,B3,機械"); return; }
  const name = parts[0], bureau = parts[1], grade = parts[2], course = parts[3] || "", mail = parts[4] || "";
  const norm = function(s){ return String(s || "").replace(/[ 　]/g, ""); };
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const results = [];

  MEMBER_DAY_SHEETS.forEach(function(sn) {
    const sheet = ss.getSheetByName(sn);
    if (!sheet) { results.push(sn + ": シートが見つかりません"); return; }
    const nameRow = sheet.getRange(3, 1, 1, sheet.getLastColumn()).getValues()[0];
    let last = 0;
    for (let c = 0; c < nameRow.length; c++) {
      if (norm(nameRow[c]) === norm(name)) { results.push(sn + ": 既に " + (c + 1) + " 列目に存在"); last = -1; break; }
      if (String(nameRow[c] || "").trim()) last = c + 1;
    }
    if (last === -1) return;
    if (last < 3) { results.push(sn + ": 名前行(3行目)が読めません"); return; }
    sheet.insertColumnBefore(last); // 条件付き書式等の適用範囲を自動拡張させるため最終メンバー列の手前に挿入
    const newCol = last;
    const src = sheet.getRange(1, newCol + 1, sheet.getMaxRows(), 1);
    const dst = sheet.getRange(1, newCol, sheet.getMaxRows(), 1);
    src.copyTo(dst, SpreadsheetApp.CopyPasteType.PASTE_FORMAT, false);
    src.copyTo(dst, SpreadsheetApp.CopyPasteType.PASTE_DATA_VALIDATION, false);
    sheet.getRange(3, newCol).setValue(name);
    sheet.getRange(4, newCol).setValue(bureau);
    sheet.getRange(5, newCol).setValue(grade);
    sheet.getRange(7, newCol).setValue(false);
    sheet.getRange(8, newCol + 1).copyTo(sheet.getRange(8, newCol));
    results.push(sn + ": " + newCol + " 列目に追加");
  });

  const meibo = ss.getSheetByName("名簿");
  if (meibo) {
    const names = meibo.getRange(1, 1, meibo.getLastRow(), 1).getValues();
    let exists = false;
    for (let i = 1; i < names.length; i++) { if (norm(names[i][0]) === norm(name)) { exists = true; break; } }
    if (exists) { results.push("名簿: 既に存在"); }
    else {
      const newRow = meibo.getLastRow() + 1;
      meibo.getRange(newRow, 1, 1, 5).setValues([[name, bureau, grade, course, mail]]);
      results.push("名簿: " + newRow + " 行目に追加");
    }
  } else { results.push("名簿: シートが見つかりません"); }

  try { SpreadsheetApp.flush(); } catch (e) {}
  ui.alert("メンバー追加の結果 → " + results.join(" ／ "));
}
