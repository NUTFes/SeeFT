// 一時的な調査用。R/S列のVLOOKUPが引けているかを確認する。確認後に削除してよい。
function inspectManualUrlLookup() {
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const lines = [];

  // 対応表シート
  const map = ss.getSheetByName("マニュアルURL");
  if (!map) {
    lines.push("「マニュアルURL」シートが見つかりません。シート一覧: " + ss.getSheets().map(s => s.getName()).join(" / "));
  } else {
    const mv = map.getRange(1, 1, Math.min(map.getLastRow(), 10), 3).getValues();
    lines.push("=== マニュアルURL シート（" + map.getLastRow() + "行） ===");
    mv.forEach((r, i) => lines.push((i + 1) + "行 A=[" + r[0] + "] B=" + String(r[1]).slice(0, 50) + " C=" + String(r[2]).slice(0, 50)));
  }

  // タスク一覧の縁日運営行
  const tasks = ss.getSheetByName("タスク一覧");
  const tv = tasks.getRange(4, 1, tasks.getLastRow() - 3, 19).getValues();
  const tf = tasks.getRange(4, 1, tasks.getLastRow() - 3, 19).getFormulas();
  lines.push("");
  lines.push("=== タスク一覧で M列が空でない行（先頭5件） ===");
  let shown = 0;
  tv.forEach((r, i) => {
    if (shown >= 5 || !String(r[12]).trim()) return;
    shown++;
    lines.push((i + 4) + "行 A=[" + r[0] + "] M=[" + r[12] + "] R=[" + String(r[17]).slice(0, 50) + "] S=[" + String(r[18]).slice(0, 50) + "]");
    lines.push("      R数式: " + (tf[i][17] || "(なし)").slice(0, 70));
  });
  lines.push("");
  lines.push("=== 縁日運営 の行 ===");
  tv.forEach((r, i) => {
    if (String(r[0]).indexOf("縁日運営") === -1) return;
    lines.push((i + 4) + "行 A=[" + r[0] + "] M=[" + r[12] + "] R=[" + r[17] + "] S=[" + r[18] + "]");
  });

  const text = lines.join("\n");
  Logger.log(text);
  SpreadsheetApp.getUi().alert(text.slice(0, 3000));
}

// タスク一覧のR列/S列にVLOOKUPの数式を貼る。手作業の貼り間違いを避けるため。
// 何度実行しても同じ結果になる（既存の数式を上書きするだけ）。
function fillManualUrlFormulas() {
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const tasks = ss.getSheetByName("タスク一覧");
  if (!tasks) return SpreadsheetApp.getUi().alert("「タスク一覧」シートが見つかりません");

  const startRow = 4;
  // IMPORTRANGEで行が増えても追従できるよう、最終行より少し先まで貼る
  const endRow = Math.max(tasks.getLastRow(), startRow) + 50;
  const count = endRow - startRow + 1;

  const rFormulas = [];
  const sFormulas = [];
  for (let r = startRow; r <= endRow; r++) {
    rFormulas.push([`=IFERROR(VLOOKUP($M${r}, 'マニュアルURL'!$A:$C, 2, FALSE), "")`]);
    sFormulas.push([`=IFERROR(VLOOKUP($M${r}, 'マニュアルURL'!$A:$C, 3, FALSE), "")`]);
  }
  tasks.getRange(startRow, 18, count, 1).setFormulas(rFormulas); // R列
  tasks.getRange(startRow, 19, count, 1).setFormulas(sFormulas); // S列

  SpreadsheetApp.getUi().alert(`タスク一覧の R${startRow}:S${endRow} に数式を貼りました（${count}行）\n\n続けて inspectManualUrlLookup で結果を確認してください`);
}

// シフト送信が「タスク一覧に無い名前」として自動作成してしまった値を、
// 値ごと・シートごと・メンバーごとに集計する。セル単位の一覧は件数が多すぎて読めないため。
// エディタ実行を想定して UI は使わない（ログにだけ出す）。
function findJunkTaskCells() {
  const JUNK = ["飛翔", "飛翔？", "一旦飛翔？", "海外", "研究", "体調不良", "実務", "不参加", "あｎ",
                "トラパ運用方付け（講義棟北駐車場側）", "トラパ運用方付け（正面入り口）", "案内所"];
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const byValue = {};   // 値 → { sheet → { name → count } }
  let total = 0;

  ["準備日", "1日目", "2日目"].forEach(function (sheetName) {
    const sheet = ss.getSheetByName(sheetName);
    if (!sheet) return;
    const lastRow = sheet.getLastRow();
    const lastCol = sheet.getLastColumn();
    const names = sheet.getRange(NAME_ROW, GRID_START_COL, 1, lastCol - GRID_START_COL + 1).getValues()[0];
    const grid = sheet.getRange(GRID_START_ROW, GRID_START_COL, lastRow - GRID_START_ROW + 1, lastCol - GRID_START_COL + 1).getValues();

    grid.forEach(function (row) {
      row.forEach(function (cell, c) {
        const v = String(cell || "").replace(/　/g, " ").trim();
        if (JUNK.indexOf(v) === -1) return;
        total++;
        const name = String(names[c] || "?");
        byValue[v] = byValue[v] || {};
        byValue[v][sheetName] = byValue[v][sheetName] || {};
        byValue[v][sheetName][name] = (byValue[v][sheetName][name] || 0) + 1;
      });
    });
  });

  const lines = ["該当セル " + total + " 件", ""];
  Object.keys(byValue).forEach(function (v) {
    lines.push("=== [" + v + "] ===");
    Object.keys(byValue[v]).forEach(function (sheetName) {
      const people = byValue[v][sheetName];
      const names = Object.keys(people);
      lines.push("  " + sheetName + ": " + names.length + "人");
      names.forEach(function (n) { lines.push("    " + n + " (" + people[n] + "コマ)"); });
    });
  });
  Logger.log(lines.join("\n"));
}
