function testPing() {
  Logger.log("ping");
}

function testSingleFileOpen() {
  const start = new Date().getTime();
  const ss = SpreadsheetApp.openByUrl("https://docs.google.com/spreadsheets/d/1OTlG5kp-o2p_k4zJzDuwbND-8bSM64m_JxBjpKL5bYY/edit");
  const openedElapsed = (new Date().getTime() - start) / 1000;
  const sheet = ss.getSheets()[0];
  const values = sheet.getDataRange().getValues();
  const totalElapsed = (new Date().getTime() - start) / 1000;
  Logger.log("openSec=" + openedElapsed + " totalSec=" + totalElapsed + " rows=" + values.length);
}

function testTaskSheetRead() {
  const start = new Date().getTime();
  const thresholds = getTaskThresholds_();
  const elapsed = (new Date().getTime() - start) / 1000;
  Logger.log("count=" + Object.keys(thresholds).length + " elapsedSec=" + elapsed);
}

const STAFFING_DAY_SHEETS = ["準々備日", "準備日", "1日目_晴れ", "1日目_雨", "2日目_晴れ", "2日目_雨", "片付け日"];

const TASK_FILE_URLS = [
  "https://docs.google.com/spreadsheets/d/1OTlG5kp-o2p_k4zJzDuwbND-8bSM64m_JxBjpKL5bYY/edit",
  "https://docs.google.com/spreadsheets/d/1bqSEaVXNpSJKwGxixd0WhpiagaVOb-TKkyhxbVpzGek/edit",
  "https://docs.google.com/spreadsheets/d/1mqXprGzGvCt-7KlcoCMR3i4wh7bvtwTTLTuFiSeFSFg/edit",
  "https://docs.google.com/spreadsheets/d/1QgBssh75EnGGpZiffq9B0IjA_aOdYrnbOaIC-DfukRs/edit",
  "https://docs.google.com/spreadsheets/d/1gf9IciVjgpdE7b9arv3pYbxGi30Um3zkBNAb8sP6ovU/edit",
  "https://docs.google.com/spreadsheets/d/1WwqPgzxqr-ifXAtsgkqNn2hmzJ-0VnkUcqbR44eR3v0/edit",
  "https://docs.google.com/spreadsheets/d/1n0PhAnybL-y_EVIkx8k2yWng1BnzC7WWw5aEQx-ogx8/edit",
  "https://docs.google.com/spreadsheets/d/16_gkzTeU8Fj5kreNtUXxsZ9T69KjVqAtbzbJ03J1UIo/edit"
];

function findTaskSheet_(ss) {
  const sheets = ss.getSheets();
  for (let i = 0; i < sheets.length; i++) {
    if (sheets[i].getName().indexOf("45th_task") === 0) return sheets[i];
  }
  return null;
}

function readThresholdsFromSheet_(sheet) {
  const values = sheet.getDataRange().getValues();
  let headerRowIndex = -1;
  for (let r = 0; r < values.length; r++) {
    if (values[r].indexOf("シフト名") !== -1) {
      headerRowIndex = r;
      break;
    }
  }
  const thresholds = {};
  if (headerRowIndex === -1) return thresholds;
  const headerRow = values[headerRowIndex];
  const nameCol = headerRow.indexOf("シフト名");
  const minCol = headerRow.indexOf("最低人数");
  const idealCol = headerRow.indexOf("適性人数");
  const maxCol = headerRow.indexOf("最大人数");

  for (let r = headerRowIndex + 1; r < values.length; r++) {
    const name = values[r][nameCol];
    if (!name) continue;
    thresholds[name] = {
      min: values[r][minCol],
      ideal: values[r][idealCol],
      max: values[r][maxCol]
    };
  }
  return thresholds;
}

function getTaskThresholds_() {
  const thresholds = {};
  TASK_FILE_URLS.forEach(function(url) {
    const ss = SpreadsheetApp.openByUrl(url);
    const sheet = findTaskSheet_(ss);
    if (!sheet) return;
    const local = readThresholdsFromSheet_(sheet);
    Object.keys(local).forEach(function(name) {
      thresholds[name] = local[name];
    });
  });
  return thresholds;
}

function timeToKey_(time) {
  if (time instanceof Date) {
    return Utilities.formatDate(time, Session.getScriptTimeZone(), "HH:mm");
  }
  return String(time);
}

function countStaffingInSheet_(sheet) {
  const lastRow = sheet.getLastRow();
  const lastCol = sheet.getLastColumn();
  const range = sheet.getRange(11, 1, lastRow - 10, lastCol);
  const values = range.getValues();

  const counts = {};
  values.forEach(function(row) {
    const time = row[0];
    if (!time) return;
    const timeKey = timeToKey_(time);
    for (let c = 1; c < row.length; c++) {
      const shiftName = row[c];
      if (!shiftName) continue;
      if (!counts[timeKey]) counts[timeKey] = {};
      counts[timeKey][shiftName] = (counts[timeKey][shiftName] || 0) + 1;
    }
  });
  return counts;
}

function judgeMark_(count, min, max) {
  if (min === "" || max === "" || min === undefined || max === undefined) return "";
  const lo = Number(min);
  const hi = Number(max);
  if (!isFinite(lo) || !isFinite(hi)) return "要確認";
  if (count < lo) return "不足";
  if (count > hi) return "過多";
  return "適正";
}

function checkStaffing() {
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const thresholds = getTaskThresholds_();

  const output = [["日程", "時刻", "シフト名", "人数", "最低", "適性", "最大", "判定"]];

  STAFFING_DAY_SHEETS.forEach(function(dayName) {
    const sheet = ss.getSheetByName(dayName);
    if (!sheet) return;

    const counts = countStaffingInSheet_(sheet);
    Object.keys(counts).forEach(function(timeKey) {
      Object.keys(counts[timeKey]).forEach(function(shiftName) {
        const count = counts[timeKey][shiftName];
        const th = thresholds[shiftName] || {};
        const mark = judgeMark_(count, th.min, th.max);
        output.push([dayName, timeKey, shiftName, count, th.min || "", th.ideal || "", th.max || "", mark]);
      });
    });
  });

  let checkSheet = ss.getSheetByName("人数チェック");
  if (!checkSheet) {
    checkSheet = ss.insertSheet("人数チェック");
  } else {
    checkSheet.clear();
  }
  checkSheet.getRange(1, 1, output.length, output[0].length).setValues(output);
}

const MATRIX_PREFIX = "人数チェック_";
const BUREAU_ORDER = ["総務局", "企画局", "渉外局", "情報局", "制作局", "財務局", "産学局", "執行部"];

// SeeFTメニュー「人数チェック（全日程）を更新」から実行。7日程シートを順番に再生成する。
// 時間主導トリガー(毎朝7時台)からも呼ばれるためUI操作は失敗しても無視する。
function buildAllStaffingMatrices() {
  const errors = [];
  STAFFING_DAY_SHEETS.forEach(function (realName) {
    try {
      buildStaffingMatrix(realName);
    } catch (e) {
      errors.push(realName + ": " + e.message);
    }
  });
  try {
    SpreadsheetApp.getUi().alert(errors.length ? ("一部失敗しました:\n" + errors.join("\n")) : ("全" + STAFFING_DAY_SHEETS.length + "日程の更新が完了しました。"));
  } catch (e) {} // 時間主導トリガーなどUIが無い実行元からの呼び出し時はアラートを出さない
}

// SeeFTメニュー「このチェックシートだけ更新」から実行。開いている「人数チェック_◯◯」シートだけを再生成する。
function buildCurrentStaffingMatrix() {
  const ui = SpreadsheetApp.getUi();
  const name = SpreadsheetApp.getActiveSheet().getName();
  if (name.indexOf(MATRIX_PREFIX) !== 0) {
    ui.alert("「" + MATRIX_PREFIX + "◯◯」のシートを開いた状態で実行してください。");
    return;
  }
  const realName = name.slice(MATRIX_PREFIX.length);
  if (STAFFING_DAY_SHEETS.indexOf(realName) === -1) {
    ui.alert("対応する日程シートが見つかりません: " + realName);
    return;
  }
  buildStaffingMatrix(realName);
  ui.alert("更新しました: " + realName);
}

// 日程シート名(realNameArg。例: "1日目_晴れ")に対応する「人数チェック_◯◯」シートを再生成する。
function buildStaffingMatrix(realNameArg) {
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const realName = String(realNameArg || "");
  const day = realName.replace("_晴れ", "").replace("_雨", "");
  const sheet = ss.getSheetByName(realName);
  if (!sheet) return;
  let out = ss.getSheetByName(MATRIX_PREFIX + realName);
  if (!out) { out = ss.insertSheet(MATRIX_PREFIX + realName); setupMatrixHeader_(out); }

  const metas = getTaskMetasFromTaskList_(); // タスク一覧(ローカル)から読む。外部8ファイルを開かないので高速
  const dd = collectDayData_(sheet);

  const taskSet = {};
  Object.keys(metas).forEach(function (n) { taskSet[n] = true; });
  Object.keys(dd.data).forEach(function (tk) {
    Object.keys(dd.data[tk]).forEach(function (n) { taskSet[n] = true; });
  });
  filterTaskSetByDay_(taskSet, dd, metas, day); // タスク一覧の日付列(metas[name].date)で選択中の日程のタスクだけに絞る
  const tasks = Object.keys(taskSet).sort(function (a, b) {
    const ia = BUREAU_ORDER.indexOf(String((metas[a] || {}).bureau));
    const ib = BUREAU_ORDER.indexOf(String((metas[b] || {}).bureau));
    if (ia !== ib) return (ia < 0 ? 99 : ia) - (ib < 0 ? 99 : ib);
    return a < b ? -1 : 1;
  });

  const header = ["シフト名", "管轄局", "開始", "終了", "レベル", "最低", "適性", "最大"].concat(dd.timeKeys);
  const rows = [header];
  tasks.forEach(function (name) {
    const m = metas[name] || {};
    const row = [name, m.bureau || "", fmtTime_(m.start), fmtTime_(m.end),
                 m.level || "", m.min || "", m.ideal || "", m.max || ""];
    const win = windowOf_(m);
    dd.timeKeys.forEach(function (tk) {
      const e = dd.data[tk] && dd.data[tk][name];
      const count = e ? e.count : 0;
      const inWin = win ? (toMinutes_(tk) >= win[0] && toMinutes_(tk) < win[1]) : false;
      if (!inWin && count === 0) { row.push(""); return; }
      let text = badge_(count, m.min, m.max);
      text += "\n" + ((e && e.senior) ? "✅" : "🟥新入生only");
      if (e) text += "\n" + e.members.join("");
      row.push(text);
    });
    rows.push(row);
  });

  if (out.getMaxColumns() < header.length) {
    out.insertColumnsAfter(out.getMaxColumns(), header.length - out.getMaxColumns());
  }
  if (out.getLastRow() >= 4) {
    out.getRange(4, 1, out.getLastRow() - 3, out.getLastColumn()).clearContent();
  }
  out.getRange(4, 1, rows.length, header.length).setValues(rows);
  // メタ列(B〜H)はタスク一覧への参照式で埋める（値ハードコード回避）
  const XLOOKUP_COLS = ["$F", "$C", "$D", "$G", "$J", "$K", "$L"];
  const metaFormulas = [];
  for (let r = 0; r < rows.length - 1; r++) metaFormulas.push(XLOOKUP_COLS.map(function (c) { return "=XLOOKUP($A" + (5 + r) + ",'タスク一覧'!$A$4:$A,'タスク一覧'!" + c + "$4:" + c + ",\"\")"; }));
  if (metaFormulas.length) out.getRange(5, 2, metaFormulas.length, 7).setFormulas(metaFormulas);
  try { if (metaFormulas.length) out.getRange(5, 3, metaFormulas.length, 2).setNumberFormat("H:mm"); } catch (e) {} // 表(テーブル)の列型と衝突する場合は書式設定を諦める
  try { if (out.getFrozenRows() < 4) out.setFrozenRows(4); if (out.getFrozenColumns() < 1) out.setFrozenColumns(1); } catch (e) {} // 表(テーブル)化されたシートでは固定操作が失敗するため無視
  try { SpreadsheetApp.flush(); } catch (e) {} // 固定行などのflushエラーはここで無視
  out.getRange("F2").setValue("最終更新: " +
    Utilities.formatDate(new Date(), Session.getScriptTimeZone(), "MM/dd HH:mm") + " → " + realName);
}

// 新規作成した「人数チェック_◯◯」シートの初期化。日程はシート自体が固定なのでプルダウンは不要。
function setupMatrixHeader_(sheet) {
  sheet.getRange("A1").setValue("使い方: SeeFTメニューの「人数チェック（全日程）を更新」または「このチェックシートだけ更新」を実行。🟦 不足 🟥 過多 ✅ 適正／2段目: ✅ 在校生あり 🟥 新入生only");
}

function collectDayData_(sheet) {
  const lastRow = sheet.getLastRow();
  const lastCol = sheet.getLastColumn();
  const head = sheet.getRange(3, 2, 5, lastCol - 1).getValues(); // 3〜7行目: 名前/局/学年/部門/新入生
  const names = head[0], bureaus = head[1], grades = head[2], freshmen = head[4];
  const grid = sheet.getRange(11, 1, lastRow - 10, lastCol).getValues();
  const data = {}; const timeKeys = [];
  grid.forEach(function (row) {
    if (!row[0]) return;
    const tk = timeToKey_(row[0]);
    timeKeys.push(tk);
    for (let c = 1; c < row.length; c++) {
      const task = row[c];
      if (!task) continue;
      const p = c - 1;
      if (!data[tk]) data[tk] = {};
      if (!data[tk][task]) data[tk][task] = { count: 0, senior: false, members: [] };
      const e = data[tk][task];
      e.count++;
      if (freshmen[p] !== true) e.senior = true;
      e.members.push("(" + String(bureaus[p] || "").slice(0, 2) + String(grades[p] || "") + ")" + String(names[p] || "").trim());
    }
  });
  return { data: data, timeKeys: timeKeys };
}

// 旧版: タスク8ファイルを直接openByUrlして読むメタ情報取得。低速なので現在はbuildStaffingMatrixから呼ばれていないが、
// checkStaffing()系のデバッグ用途に残置。
function getTaskMetas_() {
  const metas = {};
  TASK_FILE_URLS.forEach(function (url) {
    const ss = SpreadsheetApp.openByUrl(url);
    const sheet = findTaskSheet_(ss);
    if (!sheet) return;
    const values = sheet.getDataRange().getValues();
    let h = -1;
    for (let r = 0; r < values.length; r++) {
      if (values[r].indexOf("シフト名") !== -1) { h = r; break; }
    }
    if (h === -1) return;
    const head = values[h].map(String);
    const col = function (s) { for (let i = 0; i < head.length; i++) if (head[i].indexOf(s) !== -1) return i; return -1; };
    const c = { name: head.indexOf("シフト名"), bureau: col("管轄"), start: col("開始"),
                end: col("終了"), level: col("レベル"), min: col("最低"), ideal: col("適"), max: col("最大") };
    for (let r = h + 1; r < values.length; r++) {
      const name = values[r][c.name];
      if (!name || metas[name]) continue;
      const g = function (i) { return i >= 0 ? values[r][i] : ""; };
      metas[name] = { bureau: g(c.bureau), start: g(c.start), end: g(c.end),
                      level: g(c.level), min: g(c.min), ideal: g(c.ideal), max: g(c.max) };
    }
  });
  return metas;
}

function badge_(count, min, max) {
  const lo = Number(min), hi = Number(max);
  if (min === "" || max === "" || !isFinite(lo) || !isFinite(hi)) return String(count);
  if (count < lo) return "🟦" + count;
  if (count > hi) return "🟥" + count;
  return "✅" + count;
}

function fmtTime_(v) {
  if (v instanceof Date) return Utilities.formatDate(v, Session.getScriptTimeZone(), "H:mm");
  return String(v || "");
}

function toMinutes_(s) {
  const m = String(s).match(/^(\d{1,2}):(\d{2})/);
  return m ? Number(m[1]) * 60 + Number(m[2]) : -1;
}

function windowOf_(m) {
  const a = toMinutes_(fmtTime_(m.start)), b = toMinutes_(fmtTime_(m.end));
  return (a >= 0 && b >= 0) ? [a, b] : null;
}

// タスク一覧の日付情報(metas[name].date)を読み、他の日程専用のタスクをtaskSetから除外する。
// 日付空欄は全日表示、当日シートに割り当てがあるタスクは常に表示。
function filterTaskSetByDay_(taskSet, dd, metas, day) {
  const assigned = {};
  Object.keys(dd.data).forEach(function (tk) {
    Object.keys(dd.data[tk]).forEach(function (n) { assigned[n] = true; });
  });
  Object.keys(taskSet).forEach(function (n) {
    if (assigned[n]) return;
    const d = (metas[n] || {}).date;
    if (!d) return;
    if (d.indexOf(day) === -1) delete taskSet[n];
  });
}

// タスクのメタ情報(管轄局/時間/しきい値/日付)をタスク一覧シートから読む。外部8ファイルをopenByUrlするgetTaskMetas_より
// 大幅に速い(データはIMPORTRANGE経由で同じ)。列17(Q列)まで読み、dateフィールドに日付列を格納する。
function getTaskMetasFromTaskList_() {
  const ts = SpreadsheetApp.getActiveSpreadsheet().getSheetByName("タスク一覧");
  const metas = {};
  if (!ts) return metas;
  const v = ts.getRange(4, 1, Math.max(ts.getLastRow() - 3, 1), 17).getValues();
  v.forEach(function (row) {
    const n = String(row[0] || "");
    if (!n || metas[n]) return;
    metas[n] = { bureau: String(row[5] || ""), start: row[2], end: row[3], level: String(row[6] || ""),
                 min: row[9], ideal: row[10], max: row[11], date: String(row[16] || "").trim() };
  });
  return metas;
}
