const SURVEY_FILE_URLS = [
  "https://docs.google.com/spreadsheets/d/1yG5YynhIlCM_eywbkFcwdTdbc0yHpDfrJf_TxJz9L6s/edit",
  "https://docs.google.com/spreadsheets/d/1Mc70VEwm4qxmkSFb4uGeqJ1VXClG0bpl6XJn7LTv4Gs/edit",
  "https://docs.google.com/spreadsheets/d/1rH3szY0EPj_ZQs3Ci1AfhapltrAp5g9IEVFCezil3rc/edit",
  "https://docs.google.com/spreadsheets/d/1LK18X56uSR7jR4OBRcwCQN15XjChoa4-VaDSMzdb8a0/edit",
  "https://docs.google.com/spreadsheets/d/1r7YqYPouSsI-pIGMYzMWEZkUSNhA-COijoVXMRFGzoU/edit",
  "https://docs.google.com/spreadsheets/d/1lv8gULRfoBZmiOkGSkuT_Be0bea_AxLpb5SyVl_dbyw/edit",
  "https://docs.google.com/spreadsheets/d/1ZKf__kXREzyKW2vUnLuExxe0cVWPNhqqt4tYaTUqPA4/edit",
];

const AVAILABILITY_GROUPS = [
  { real: ["準々備日"], surveyTab: "準々備日(9/17木)" },
  { real: ["準備日"], surveyTab: "準備日(9/18金)" },
  { real: ["1日目_晴れ", "1日目_雨"], surveyTab: "1日目(9/19土)" },
  { real: ["2日目_晴れ", "2日目_雨"], surveyTab: "2日目(9/20日)" },
  { real: ["片付け日"], surveyTab: "片付け日(9/21月)" },
];

const SLOT_COUNT = 64;
const PEOPLE_COUNT = 350;

// 参加不可セルの塗り色（希望調査の黒塗りに合わせる）
const BLOCK_COLOR = "#000000";

function applyAvailabilityNotes() {
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const files = SURVEY_FILE_URLS.map(function (url) {
    return SpreadsheetApp.openByUrl(url);
  });
  let clearedCount = 0;
  AVAILABILITY_GROUPS.forEach(function (group) {
    const availMap = buildAvailMap_(files, group.surveyTab);
    group.real.forEach(function (realName) {
      const sheet = ss.getSheetByName(realName);
      if (!sheet) return;
      const names = sheet.getRange(3, 2, 1, PEOPLE_COUNT).getValues()[0];
      const grid = sheet.getRange(11, 2, SLOT_COUNT, PEOPLE_COUNT);
      const bgs = grid.getBackgrounds();
      const vals = grid.getValues();
      const clearCells = [];
      const notes = [];
      for (let t = 0; t < SLOT_COUNT; t++) {
        notes.push(new Array(PEOPLE_COUNT).fill(""));
      }
      names.forEach(function (n, c) {
        const name = normalizeName_(n);
        if (!name) return;
        const slots = availMap[name];
        for (let t = 0; t < SLOT_COUNT; t++) {
          if (!slots) {
            notes[t][c] = "参加不可: 未回答";
          } else if (slots[t]) {
            notes[t][c] = "参加不可: " + slots[t];
            bgs[t][c] = BLOCK_COLOR;
            if (String(vals[t][c]).trim()) {
              // 参加不可セルに入っている割当を取り消して、黒塗りが見えるようにする
              clearCells.push(sheet.getRange(11 + t, 2 + c).getA1Notation());
            }
          } else if (isBlockedBg_(bgs[t][c])) {
            bgs[t][c] = null;
          }
        }
      });
      grid.setNotes(notes);
      grid.setBackgrounds(bgs);
      if (clearCells.length) {
        sheet.getRangeList(clearCells).clearContent();
        clearedCount += clearCells.length;
      }
    });
  });
  SpreadsheetApp.getUi().alert(
    "参加不可の警告（メモ＋黒塗り）を更新しました。" +
    (clearedCount ? "\n参加不可セルに入っていた割当 " + clearedCount + " 件を取り消しました。" : ""));
}

function normalizeName_(v) {
  return String(v).replace(/[\s　]/g, "");
}

function isBlockedBg_(hex) {
  const s = String(hex || "").trim().toLowerCase();
  if (s.length !== 7 || s.charAt(0) !== "#") return false;
  const r = parseInt(s.substring(1, 3), 16);
  const g = parseInt(s.substring(3, 5), 16);
  const b = parseInt(s.substring(5, 7), 16);
  return 0.299 * r + 0.587 * g + 0.114 * b < 128;
}

function buildAvailMap_(files, tabName) {
  const map = {};
  files.forEach(function (file) {
    let tab = file.getSheetByName(tabName);
    if (!tab) {
      const prefix = tabName.split("(")[0];
      tab = file.getSheets().filter(function (s) {
        return s.getName().indexOf(prefix) === 0;
      })[0];
    }
    if (!tab) return;
    const lastRow = tab.getLastRow();
    if (lastRow < 3) return;
    const range = tab.getRange(3, 1, lastRow - 2, 70);
    const values = range.getValues();
    const bgs = range.getBackgrounds();
    values.forEach(function (row, r) {
      const name = normalizeName_(row[0]);
      if (!name) return;
      if (name in map) return;
      const slots = [];
      let reason = "";
      for (let t = 0; t < SLOT_COUNT; t++) {
        const c = 6 + t;
        const text = String(row[c]).trim();
        if (text) reason = text;
        if (text || isBlockedBg_(bgs[r][c])) {
          slots.push(reason || "理由未記入");
        } else {
          reason = "";
        slots.push("");
        }
      }
      map[name] = slots;
    });
  });
  return map;
}