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

// 対応表（マニュアルURLシート）とタスク一覧M列の整合を検査する。
// 紐付けの事故はどれも「エラーが出ずに空になる」形で起きるため、送信前にここで洗い出す。
//   不足:   M列にあるのに対応表A列に無い → そのタスクはマニュアルが引けない
//   惜しい: 空白を除けば一致する → 末尾スペース・全角の揺れ（実例: 技大カジノ、ステージ横テント受付）
//   余り:   対応表A列にあるのにM列に無い → どのタスクにも効いていない行
//   重複:   対応表A列に同じ値が複数行 → VLOOKUPは最初の行しか引かず2行目以降は死ぬ
//   先勝ち: 同名タスクの最初の行だけがM列空 → 送信は先勝ちなので空が送られる（実例: 謎解き）
//   B/C列: 使われている行のURL欠け → シフトカードにボタンが出ない
function checkManualUrlMapping() {
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const ui = SpreadsheetApp.getUi();

  const map = ss.getSheetByName("マニュアルURL");
  if (!map) return ui.alert("「マニュアルURL」シートが見つかりません");
  const tasks = ss.getSheetByName(TASK_LIST_SHEET);
  if (!tasks) return ui.alert("「" + TASK_LIST_SHEET + "」シートが見つかりません");

  // 対応表を読む（1行目はヘッダー）
  const entries = new Map(); // A値 → {row, doc, slide, taskNames:Set}
  const dupKeys = [];
  if (map.getLastRow() >= 2) {
    const mv = map.getRange(2, 1, map.getLastRow() - 1, 3).getValues();
    mv.forEach(function (r, i) {
      const key = String(r[0] || "");
      if (!key.trim()) return;
      if (entries.has(key)) { dupKeys.push((i + 2) + "行: " + key); return; }
      entries.set(key, { row: i + 2, doc: String(r[1] || ""), slide: String(r[2] || ""), taskNames: new Set() });
    });
  }

  // タスク一覧を読む。送信ロジック（buildTaskChanges_）と同じ「タスク名の先勝ち」で
  // 1タスク1行に畳む。ここを揃えないと、送信では空になる行を「紐付いている」と誤報する
  const lastRow = tasks.getLastRow();
  const rows = lastRow >= TASK_LIST_START_ROW
    ? tasks.getRange(TASK_LIST_START_ROW, 1, lastRow - TASK_LIST_START_ROW + 1, 13).getValues()
    : [];
  const firstM = new Map();   // タスク名 → 先勝ちで採用されるM値
  const laterM = new Map();   // タスク名 → 2行目以降にだけ現れた非空のM値
  rows.forEach(function (r) {
    const taskName = String(r[0] || "").replace(/　/g, " ").trim();
    if (!taskName) return;
    const m = String(r[12] || "");
    if (!firstM.has(taskName)) {
      firstM.set(taskName, m);
    } else if (m.trim() && !String(firstM.get(taskName)).trim()) {
      laterM.set(taskName, m);
    }
  });

  // 突き合わせ
  const missing = new Map(); // M値 → タスク名[]
  firstM.forEach(function (m, taskName) {
    if (!m.trim()) return;
    if (entries.has(m)) { entries.get(m).taskNames.add(taskName); return; }
    if (!missing.has(m)) missing.set(m, []);
    missing.get(m).push(taskName);
  });

  const lines = [];

  if (missing.size) {
    lines.push("■ 不足: M列にあるのに対応表A列に無い（" + missing.size + "種）");
    missing.forEach(function (taskNames, m) {
      lines.push("・[" + m + "] ← " + taskNames.join(" / "));
      entries.forEach(function (_, key) {
        if (key !== m && key.trim() === m.trim()) {
          lines.push("    ※惜しい: 対応表の [" + key + "] と空白の有無だけが違う");
        }
      });
    });
    lines.push("");
  }

  const unused = [];
  const urlProblems = [];
  entries.forEach(function (e, key) {
    if (e.taskNames.size === 0) { unused.push(e.row + "行: [" + key + "]"); return; }
    // 送信時と同じ判定（httpUrlOrEmpty_）で、実際にボタンが出るかを見る
    if (!httpUrlOrEmpty_(e.doc)) urlProblems.push(e.row + "行 [" + key + "] B列(ドキュメント版)が" + (e.doc.trim() ? "URLでない" : "空"));
    if (!httpUrlOrEmpty_(e.slide)) urlProblems.push(e.row + "行 [" + key + "] C列(HTML版)が" + (e.slide.trim() ? "URLでない" : "空"));
  });

  if (unused.length) {
    lines.push("■ 余り: 対応表にあるのにM列のどのタスクからも参照されていない（" + unused.length + "行）");
    unused.forEach(function (s) { lines.push("・" + s); });
    lines.push("");
  }
  if (dupKeys.length) {
    lines.push("■ 重複: 対応表A列に同じ値が複数ある（2行目以降は引かれない）");
    dupKeys.forEach(function (s) { lines.push("・" + s); });
    lines.push("");
  }
  if (laterM.size) {
    lines.push("■ 先勝ちの罠: 同名タスクの最初の行だけM列が空（このままだと空が送られる）");
    laterM.forEach(function (m, taskName) { lines.push("・" + taskName + "（後の行には [" + m + "] がある）"); });
    lines.push("");
  }
  if (urlProblems.length) {
    lines.push("■ URL欠け: 使われている行のB/C列");
    urlProblems.forEach(function (s) { lines.push("・" + s); });
    lines.push("");
  }

  lines.push("■ 紐付けの要約");
  let okKeys = 0;
  entries.forEach(function (e, key) {
    if (e.taskNames.size === 0) return;
    okKeys++;
    lines.push("・[" + key + "] → " + e.taskNames.size + "タスク");
  });
  if (!okKeys) lines.push("・引けている行がありません");

  const problems = missing.size + unused.length + dupKeys.length + laterM.size + urlProblems.length;
  lines.unshift(problems ? "問題 " + problems + " 件。タスク送信の前に直してください。" : "問題なし。タスク送信できます。", "");

  const text = lines.join("\n");
  Logger.log(text);
  ui.alert(text.slice(0, 3000) + (text.length > 3000 ? "\n…（続きはログを参照）" : ""));
}
