// 名簿(users)とタスク・集合場所(tasks/places)をSeeFTに送信する。
// シフト送信の前提となるマスタで、投入順序は 名簿 → タスク → シフト。
// シフトはユーザー名・タスク名の完全一致で紐づくため、ここで送る表記と
// 日程シート上の表記が一致している必要がある。

// 名簿の取得元。日程シートの3〜5行目(名前/局/学年)を使う。
// 「名簿」シートではなく日程シートを使うのは、シフト送信と同じセルから
// 名前を読むことで表記のズレを構造的に防ぐため。
const ROSTER_SOURCE_SHEET = "準々備日";
const ROSTER_NAME_ROW = 3;
const ROSTER_BUREAU_ROW = 4;
const ROSTER_GRADE_ROW = 5;
const ROSTER_START_COL = 2;

// APIが受け付ける局名・学年（user_usecase.goの変換表と対応）。
// 一致しないとIDが0になり外部キー違反で失敗するため、送信前に検証する。
const VALID_BUREAUS = ["執行部", "執行部補佐", "総務局", "企画局", "渉外局", "財務局", "制作局", "情報局", "産学局"];
const VALID_GRADES = ["B1", "B2", "B3", "B4", "M1", "M2", "D1", "D2", "D3", "OB"];

// タスク一覧シートの列（1始まり）。A=シフト名 / F=管轄局 / L=最大人数
const TASK_LIST_SHEET = "タスク一覧";
const TASK_LIST_START_ROW = 4;
const TASK_COL_NAME = 1;
const TASK_COL_BUREAU = 6;
const TASK_COL_MAX = 12;

// 名簿をSeeFTに送信する
function updateUsers() {
  ui = ui || SpreadsheetApp.getUi();
  try {
    lock.waitLock(30000);

    const sheet = SpreadsheetApp.getActiveSpreadsheet().getSheetByName(ROSTER_SOURCE_SHEET);
    if (!sheet) {
      return ui.alert(`名簿の取得元シート「${ROSTER_SOURCE_SHEET}」が見つかりません`);
    }

    const built = buildUserChanges_(sheet);
    if (built.errors.length) {
      return ui.alert(
        `送信できない値がありました。修正して再実行してください\n\n` +
        built.errors.slice(0, 20).join("\n") +
        (built.errors.length > 20 ? `\n...他${built.errors.length - 20}件` : ""));
    }
    if (!built.changes.length) {
      return ui.alert(`送信対象のメンバーが見つかりませんでした`);
    }

    const confirm = ui.alert(
      `名簿 ${built.changes.length} 件をSeeFTに送信してよろしいですか？\n` +
      `（既存のユーザーは更新、未登録のユーザーは新規作成されます）`,
      ui.ButtonSet.OK_CANCEL);
    if (confirm === ui.Button.CANCEL) {
      Logger.log("送信がキャンセルされました");
      return;
    }

    postToSeeFT_("/api/update_users", built.changes);
    ui.alert(`名簿を送信しました\n${built.changes.length} 件`);
  } catch (error) {
    ui.alert(`エラーが発生しました\n
      エラーを修正して再実行してください\n
      ------エラーメッセージ-----\n` +
      error.message +
      `\n---------------------------`);
    Logger.log("Error: " + error.message);
  } finally {
    lock.releaseLock();
  }
}

// 局名から主務(SeeFTに登録する1つの局)を取り出す。
// 兼務は「企画局, 財務局」「企画局/財務局」のように区切り文字で併記されるため、
// 先頭を主務として扱う（users.bureau_idが1つしか持てないため）
function primaryBureauOf_(raw) {
  return String(raw || "").split(/[,、\/／]/)[0].replace(/[\s　]/g, "");
}

// 日程シートの3〜5行目から名簿の送信内容を組み立てる。
// 学籍番号と電話番号は日程シートに無いため、名簿作成アンケートの回答から名前で引く。
// 学籍番号はログイン(学籍番号+パスワード)の識別子なので、欠けたまま送ってはいけない。
function buildUserChanges_(sheet) {
  const lastCol = sheet.getLastColumn();
  const count = lastCol - ROSTER_START_COL + 1;
  const changes = [];
  const errors = [];
  if (count < 1) return { changes: changes, errors: errors };

  // 取得元が読めない場合は送信を止める。学籍番号0のまま投入すると全員がログインできず、
  // しかもエラーが出ないため気づけない(2026-08-26に本番で実際に発生)
  let srcMap;
  try {
    srcMap = loadStudentNumberMap_().map;
  } catch (e) {
    errors.push("名簿作成アンケートを読めません: " + e.message);
    return { changes: changes, errors: errors };
  }

  const names = sheet.getRange(ROSTER_NAME_ROW, ROSTER_START_COL, 1, count).getValues()[0];
  const bureaus = sheet.getRange(ROSTER_BUREAU_ROW, ROSTER_START_COL, 1, count).getValues()[0];
  const grades = sheet.getRange(ROSTER_GRADE_ROW, ROSTER_START_COL, 1, count).getValues()[0];

  const seen = {};
  for (let i = 0; i < count; i++) {
    // シフト送信と同じ扱い（trimのみ）にして、DBに入る表記を一致させる
    const name = String(names[i] || "").trim();
    if (!name) continue;

    if (seen[name]) {
      errors.push(`${name}: 同じ名前が複数の列にあります`);
      continue;
    }
    seen[name] = true;

    // 兼務は「企画局, 財務局」のようにカンマ区切りで書かれるが、users.bureau_idは
    // 1つしか持てないため、先頭を主務として送る
    const bureau = primaryBureauOf_(bureaus[i]);
    const grade = String(grades[i] || "").trim();

    if (VALID_BUREAUS.indexOf(bureau) === -1) {
      errors.push(`${name}: 局名「${bureau}」はSeeFTに登録できません`);
      continue;
    }
    if (VALID_GRADES.indexOf(grade) === -1) {
      errors.push(`${name}: 学年「${grade}」はSeeFTに登録できません`);
      continue;
    }

    // 学籍番号・電話番号をアンケート回答から引く。照合は空白を除いた名前で行うが、
    // 送信する名前は日程シートの表記のまま(シフト送信と一致させる必要があるため)
    const src = srcMap[normalizeForMatch_(name)];
    if (!src) {
      errors.push(`${name}: 名簿作成アンケートに回答が見つかりません`);
      continue;
    }
    if (!/^\d{8}$/.test(src.studentNumber)) {
      errors.push(`${name}: 学籍番号「${src.studentNumber}」が8桁ではありません`);
      continue;
    }
    // 電話番号はハイフン付きで入力されることがあるため数字だけにする
    const tel = String(src.tel || "").replace(/[^0-9]/g, "");

    changes.push({
      name: name,
      bureau: bureau,
      grade: grade,
      // 課程はアンケートの表記(電気/物生など)がAPIの変換表(電気電子情報工学分野など)と
      // 噛み合わず、既定の未所属に落ちる。どこからも参照されないため技大祭後に対応する
      department: "未所属",
      studentNumber: Number(src.studentNumber),
      tel: tel,
      // 技大祭メールアドレス。Slack ID紐付け(users.lookupByEmail)で使う
      mail: src.mail
    });
  }
  return { changes: changes, errors: errors };
}

// タスクと集合場所をSeeFTに送信する
function updateTasksAndPlaces() {
  ui = ui || SpreadsheetApp.getUi();
  try {
    lock.waitLock(30000);

    const sheet = SpreadsheetApp.getActiveSpreadsheet().getSheetByName(TASK_LIST_SHEET);
    if (!sheet) {
      return ui.alert(`「${TASK_LIST_SHEET}」シートが見つかりません`);
    }

    const changes = buildTaskChanges_(sheet);
    if (!changes.length) {
      return ui.alert(`送信対象のタスクが見つかりませんでした`);
    }

    const confirm = ui.alert(
      `タスク ${changes.length} 件をSeeFTに送信してよろしいですか？\n` +
      `（既存のタスクは更新、未登録のタスクは新規作成されます）`,
      ui.ButtonSet.OK_CANCEL);
    if (confirm === ui.Button.CANCEL) {
      Logger.log("送信がキャンセルされました");
      return;
    }

    postToSeeFT_("/api/update_tasks_and_places", changes);
    ui.alert(`タスクを送信しました\n${changes.length} 件`);
  } catch (error) {
    ui.alert(`エラーが発生しました\n
      エラーを修正して再実行してください\n
      ------エラーメッセージ-----\n` +
      error.message +
      `\n---------------------------`);
    Logger.log("Error: " + error.message);
  } finally {
    lock.releaseLock();
  }
}

// タスク一覧シートから送信内容を組み立てる
function buildTaskChanges_(sheet) {
  const lastRow = sheet.getLastRow();
  if (lastRow < TASK_LIST_START_ROW) return [];

  const rows = sheet.getRange(TASK_LIST_START_ROW, 1, lastRow - TASK_LIST_START_ROW + 1, TASK_COL_MAX).getValues();
  const changes = [];
  const seen = {};

  rows.forEach(function (row) {
    // シフト送信は全角スペースを半角に寄せて照合するため、こちらも同じ扱いにする
    const taskName = String(row[TASK_COL_NAME - 1] || "").replace(/　/g, " ").trim();
    if (!taskName) return;
    if (seen[taskName]) return; // タスク一覧には同名が複数日程分並ぶため先勝ちで1件にまとめる
    seen[taskName] = true;

    const maxMember = Number(row[TASK_COL_MAX - 1]);

    changes.push({
      yearID: YEAR_ID,
      taskName: taskName,
      bureau: String(row[TASK_COL_BUREAU - 1] || "").trim(),
      place: "",   // 集合場所はタスク一覧に列が無いためAPI側の既定値(ID=1)に任せる
      url: "",     // マニュアルURLは別途アップロードAPIで紐づけるためここでは送らない
      maxMember: isFinite(maxMember) && maxMember > 0 ? maxMember : 1
    });
  });
  return changes;
}

// SeeFTのエンドポイントにchangesを送信する共通処理
function postToSeeFT_(path, changes) {
  const options = {
    method: "post",
    contentType: "application/json",
    payload: JSON.stringify({ changes: changes }),
    muteHttpExceptions: true
  };
  const response = UrlFetchApp.fetch(baseUrl + path, options);
  const code = response.getResponseCode();
  const body = response.getContentText();
  Logger.log("Response from API (" + code + "): " + body);
  if (code < 200 || code >= 300) {
    throw new Error(`APIがエラーを返しました (HTTP ${code})\n${body}`);
  }
  return response;
}

// 名簿の取得元(準々備日の3〜5行目)を準備日と突き合わせ、送信できない値を洗い出す。
//
// 名前行が壊れていた事故を受けて、局・学年の行にも同じ破損が無いかを確認するために使う。
// 併せてVALID_BUREAUS / VALID_GRADESに一致しない値を列挙する。APIは一致しない局名を
// ID=0として保存しようとして外部キー違反で落ちるため、送信前にここで潰す。
function compareRosterRows() {
  ui = ui || SpreadsheetApp.getUi();
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const target = ss.getSheetByName(ROSTER_SOURCE_SHEET);
  const base = ss.getSheetByName("準備日");
  if (!target || !base) return ui.alert("準々備日 または 準備日 が見つかりません");

  const tCount = target.getLastColumn() - ROSTER_START_COL + 1;
  const bCount = base.getLastColumn() - ROSTER_START_COL + 1;
  const lines = [];
  lines.push("準々備日の列数 " + tCount + " / 準備日の列数 " + bCount);

  const rows = [
    { row: ROSTER_NAME_ROW, label: "3行目(名前)" },
    { row: ROSTER_BUREAU_ROW, label: "4行目(局)" },
    { row: ROSTER_GRADE_ROW, label: "5行目(学年)" }
  ];
  const tVals = {}, bVals = {};
  rows.forEach(function (r) {
    tVals[r.row] = target.getRange(r.row, ROSTER_START_COL, 1, tCount).getValues()[0];
    bVals[r.row] = base.getRange(r.row, ROSTER_START_COL, 1, bCount).getValues()[0];
  });

  // 行ごとに準備日との食い違いを数える
  const n = Math.min(tCount, bCount);
  rows.forEach(function (r) {
    const diff = [];
    for (let i = 0; i < n; i++) {
      if (String(tVals[r.row][i] || "").trim() !== String(bVals[r.row][i] || "").trim()) diff.push(i);
    }
    if (!diff.length) { lines.push(r.label + ": 準備日と完全一致"); return; }
    const f = diff[0], l = diff[diff.length - 1];
    lines.push(r.label + ": 食い違い " + diff.length + " 列 (" +
      colLetter_(ROSTER_START_COL + f) + "〜" + colLetter_(ROSTER_START_COL + l) + ", " +
      ((l - f + 1) === diff.length ? "連続" : "飛び地あり") + ")");
    lines.push("    先頭 " + colLetter_(ROSTER_START_COL + f) + ": 準々備日「" +
      String(tVals[r.row][f]).trim() + "」/ 準備日「" + String(bVals[r.row][f]).trim() + "」");
  });

  // 送信時に弾かれる値を洗い出す
  const badBureau = [], badGrade = [];
  const bureauCount = {};
  for (let i = 0; i < tCount; i++) {
    const name = String(tVals[ROSTER_NAME_ROW][i] || "").trim();
    if (!name) continue;
    const col = colLetter_(ROSTER_START_COL + i);
    const bureau = primaryBureauOf_(tVals[ROSTER_BUREAU_ROW][i]);
    const grade = String(tVals[ROSTER_GRADE_ROW][i] || "").trim();
    bureauCount[bureau] = (bureauCount[bureau] || 0) + 1;
    if (VALID_BUREAUS.indexOf(bureau) === -1) badBureau.push(col + " " + name + ": 局「" + bureau + "」");
    if (VALID_GRADES.indexOf(grade) === -1) badGrade.push(col + " " + name + ": 学年「" + grade + "」");
  }

  lines.push("");
  lines.push("局の内訳: " + Object.keys(bureauCount).sort().map(function (k) {
    return (k || "(空)") + " " + bureauCount[k];
  }).join(" / "));
  lines.push("送信できない局 " + badBureau.length + " 件 / 送信できない学年 " + badGrade.length + " 件");
  badBureau.slice(0, 8).forEach(function (m) { lines.push("    " + m); });
  if (badBureau.length > 8) lines.push("    ...他 " + (badBureau.length - 8) + " 件はログ参照");
  badGrade.slice(0, 8).forEach(function (m) { lines.push("    " + m); });
  if (badGrade.length > 8) lines.push("    ...他 " + (badGrade.length - 8) + " 件はログ参照");

  const text = lines.join("\n");
  Logger.log(text);
  ui.alert("名簿行の突き合わせ（準々備日 vs 準備日）\n\n" + text.slice(0, 1500));
}

// タスク一覧シートの管轄局を集計し、SeeFTに登録できない局名を洗い出す。
//
// 名簿送信と違い、タスク送信はAPI側に局名の検証が無い。task_usecase.goの変換表に
// 一致しない局名は default で bureauID=1（執行部）として黙って保存されるため、
// 取り違えに気づけない。送信前後の照合用にシート側の実数をここで数える。
function checkTaskBureaus() {
  ui = ui || SpreadsheetApp.getUi();
  const sheet = SpreadsheetApp.getActiveSpreadsheet().getSheetByName(TASK_LIST_SHEET);
  if (!sheet) return ui.alert(`「${TASK_LIST_SHEET}」シートが見つかりません`);

  const lastRow = sheet.getLastRow();
  if (lastRow < TASK_LIST_START_ROW) return ui.alert("タスクがありません");
  const rows = sheet.getRange(TASK_LIST_START_ROW, 1, lastRow - TASK_LIST_START_ROW + 1, TASK_COL_MAX).getValues();

  const seen = {};
  const count = {};   // 局名 -> 件数
  const bad = [];     // 変換表に無い局名のタスク
  let total = 0;

  rows.forEach(function (row, i) {
    // buildTaskChanges_と同じ前処理（全角スペースを半角に寄せ、同名は先勝ち）
    const taskName = String(row[TASK_COL_NAME - 1] || "").replace(/　/g, " ").trim();
    if (!taskName) return;
    if (seen[taskName]) return;
    seen[taskName] = true;
    total++;

    // task_usecase.goはスペースを全て除去してから変換表と照合する
    const bureau = String(row[TASK_COL_BUREAU - 1] || "").replace(/[\s　]/g, "");
    count[bureau] = (count[bureau] || 0) + 1;
    if (VALID_BUREAUS.indexOf(bureau) === -1) {
      bad.push((TASK_LIST_START_ROW + i) + "行目 " + taskName + ": 局「" + (bureau || "(空)") + "」");
    }
  });

  const lines = [];
  lines.push("送信対象タスク " + total + " 件（同名は1件に集約後）");
  lines.push("");
  lines.push("局の内訳:");
  Object.keys(count).sort().forEach(function (k) {
    const ok = VALID_BUREAUS.indexOf(k) !== -1;
    lines.push("  " + (k || "(空)") + " " + count[k] + (ok ? "" : "  ← 執行部として保存される"));
  });
  lines.push("");
  lines.push(bad.length === 0
    ? "全ての局名が変換表に一致します。"
    : "変換表に無い局名が " + bad.length + " 件あります。これらは執行部(bureau_id=1)になります:");
  bad.slice(0, 10).forEach(function (m) { lines.push("    " + m); });
  if (bad.length > 10) lines.push("    ...他 " + (bad.length - 10) + " 件はログ参照");

  const text = lines.join("\n");
  Logger.log(text);
  ui.alert("タスク一覧の管轄局チェック\n\n" + text.slice(0, 1400));
}

// 名簿シートの列構成を調べる調査用。
// 日程シートには学籍番号・電話番号・課程が無いため、名簿シートから補える列があるかを確認する。
// ログインは学籍番号+パスワードで行うため、学籍番号が取れないと誰もログインできない。
function inspectRosterSheet() {
  ui = ui || SpreadsheetApp.getUi();
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const sheet = ss.getSheetByName("名簿");
  if (!sheet) {
    const names = ss.getSheets().map(function (s) { return s.getName(); });
    return ui.alert("「名簿」シートが見つかりません\n\nシート一覧:\n" + names.join("\n").slice(0, 1200));
  }

  const lastRow = Math.min(sheet.getLastRow(), 12);
  const lastCol = Math.min(sheet.getLastColumn(), 20);
  const vals = sheet.getRange(1, 1, lastRow, lastCol).getValues();

  const lines = [];
  lines.push("名簿シート: " + sheet.getLastRow() + "行 × " + sheet.getLastColumn() + "列");
  lines.push("");
  for (let r = 0; r < lastRow; r++) {
    const cells = [];
    for (let c = 0; c < lastCol; c++) {
      const v = String(vals[r][c] || "").trim();
      // 学籍番号など個人情報が並ぶため、値は先頭8文字までに切って構造だけ見る
      cells.push(colLetter_(c + 1) + ":" + (v.length > 8 ? v.slice(0, 8) + "…" : v));
    }
    lines.push((r + 1) + "行 | " + cells.join(" | "));
  }

  const text = lines.join("\n");
  Logger.log(text);
  ui.alert("名簿シートの構造\n\n" + text.slice(0, 3000));
}

// 学籍番号の取得元スプレッドシート。
// ver0の名簿シートには学籍番号の列が無く、ログイン(学籍番号+パスワード)ができないため
// ここから名前で引いて補う。IMPORTRANGEではなくGASのopenByIdで読むのは、
// ファイルペアごとの接続許可が要らず、コピー作成時に切れる問題も起きないため。
const STUDENT_NUMBER_SOURCE_ID = "10JQK0AGBJKtm-dQMfysQjEvSkh3leurp7e6kNqY5eLw";
const STUDENT_NUMBER_SOURCE_GID = 211416656;

// 取得元シートの構造を調べる調査用。どの列が名前で、どの列が学籍番号かを特定する。
function inspectStudentNumberSource() {
  ui = ui || SpreadsheetApp.getUi();
  let ss;
  try {
    ss = SpreadsheetApp.openById(STUDENT_NUMBER_SOURCE_ID);
  } catch (e) {
    return ui.alert("取得元スプレッドシートを開けません\n\n" + e.message +
      "\n\n技大祭アカウントでアクセス権があるか確認してください。");
  }

  const sheets = ss.getSheets();
  let sheet = null;
  sheets.forEach(function (s) { if (s.getSheetId() === STUDENT_NUMBER_SOURCE_GID) sheet = s; });
  if (!sheet) {
    return ui.alert("gid=" + STUDENT_NUMBER_SOURCE_GID + " のシートが見つかりません\n\nシート一覧:\n" +
      sheets.map(function (s) { return s.getName() + " (gid=" + s.getSheetId() + ")"; }).join("\n").slice(0, 1200));
  }

  const lastRow = Math.min(sheet.getLastRow(), 8);
  const lastCol = Math.min(sheet.getLastColumn(), 20);
  const vals = sheet.getRange(1, 1, lastRow, lastCol).getValues();

  const lines = [];
  lines.push("ファイル: " + ss.getName());
  lines.push("シート: " + sheet.getName() + " (" + sheet.getLastRow() + "行 × " + sheet.getLastColumn() + "列)");
  lines.push("");
  for (let r = 0; r < lastRow; r++) {
    const cells = [];
    for (let c = 0; c < lastCol; c++) {
      const v = String(vals[r][c] || "").trim();
      // 学籍番号などが並ぶため、構造が分かる範囲に切り詰めて表示する
      cells.push(colLetter_(c + 1) + ":" + (v.length > 8 ? v.slice(0, 8) + "…" : v));
    }
    lines.push((r + 1) + "行 | " + cells.join(" | "));
  }

  const text = lines.join("\n");
  Logger.log(text);
  ui.alert("学籍番号の取得元シート\n\n" + text.slice(0, 3000));
}

// 学籍番号取得元(フォームの回答)の列。1始まり。
const SRC_COL_STUDENT_NUMBER = 3;  // C: 学籍番号(8桁)
const SRC_COL_NAME = 4;            // D: 氏名
const SRC_COL_MAIL = 6;            // F: 技大祭メールアドレス
const SRC_COL_DEPARTMENT = 9;      // I: 課程・分野
const SRC_COL_TEL = 12;            // L: 電話番号
const SRC_DATA_START_ROW = 2;      // 1行目はヘッダー

// 名前の表記ゆれを吸収するための正規化。
// 姓名の間の空白は全角・半角・複数が混在するため、照合時のみ全て除去する。
// 保存する名前は日程シートの表記をそのまま使うので、ここでの正規化は照合専用。
function normalizeForMatch_(s) {
  return String(s || "").replace(/[\s　]/g, "").trim();
}

// フォームの回答を名前をキーにしたマップで返す
function loadStudentNumberMap_() {
  const ss = SpreadsheetApp.openById(STUDENT_NUMBER_SOURCE_ID);
  let sheet = null;
  ss.getSheets().forEach(function (s) { if (s.getSheetId() === STUDENT_NUMBER_SOURCE_GID) sheet = s; });
  if (!sheet) throw new Error("gid=" + STUDENT_NUMBER_SOURCE_GID + " のシートが見つかりません");

  const lastRow = sheet.getLastRow();
  if (lastRow < SRC_DATA_START_ROW) return { map: {}, dup: [], rows: 0 };
  const vals = sheet.getRange(SRC_DATA_START_ROW, 1, lastRow - SRC_DATA_START_ROW + 1, SRC_COL_TEL).getValues();

  const map = {}, dup = [];
  let rows = 0;
  vals.forEach(function (row) {
    const name = normalizeForMatch_(row[SRC_COL_NAME - 1]);
    if (!name) return;
    rows++;
    if (map[name]) { dup.push(name); return; } // 重複回答は先勝ち
    map[name] = {
      studentNumber: String(row[SRC_COL_STUDENT_NUMBER - 1] || "").trim(),
      mail: String(row[SRC_COL_MAIL - 1] || "").trim(),
      department: String(row[SRC_COL_DEPARTMENT - 1] || "").trim(),
      tel: String(row[SRC_COL_TEL - 1] || "").trim()
    };
  });
  return { map: map, dup: dup, rows: rows };
}

// 日程シートの353人が、フォームの回答で引けるかを確認する調査用。
// 引けない人がいると学籍番号が0のままになりログインできないため、書き込み前に洗い出す。
function checkStudentNumberMatch() {
  ui = ui || SpreadsheetApp.getUi();
  let src;
  try {
    src = loadStudentNumberMap_();
  } catch (e) {
    return ui.alert("取得元を読めません\n\n" + e.message);
  }

  const sheet = SpreadsheetApp.getActiveSpreadsheet().getSheetByName(ROSTER_SOURCE_SHEET);
  if (!sheet) return ui.alert("「" + ROSTER_SOURCE_SHEET + "」が見つかりません");
  const count = sheet.getLastColumn() - ROSTER_START_COL + 1;
  const names = sheet.getRange(ROSTER_NAME_ROW, ROSTER_START_COL, 1, count).getValues()[0];

  const missing = [], noNumber = [], badNumber = [];
  let matched = 0, filled = 0;
  for (let i = 0; i < count; i++) {
    const raw = String(names[i] || "").trim();
    if (!raw) continue;
    filled++;
    const hit = src.map[normalizeForMatch_(raw)];
    if (!hit) { missing.push(colLetter_(ROSTER_START_COL + i) + " " + raw); continue; }
    matched++;
    if (!hit.studentNumber) noNumber.push(raw);
    else if (!/^\d{8}$/.test(hit.studentNumber)) badNumber.push(raw + ": 「" + hit.studentNumber + "」");
  }

  const lines = [];
  lines.push("取得元の回答 " + src.rows + " 行 / 名前ユニーク " + Object.keys(src.map).length +
    (src.dup.length ? " / 同名の重複回答 " + src.dup.length + "件(先勝ちで採用)" : ""));
  lines.push(ROSTER_SOURCE_SHEET + "の名前 " + filled + " 人 / 引けた " + matched + " 人 / 引けない " + missing.length + " 人");
  lines.push("");
  if (missing.length) {
    lines.push("取得元に見つからない人:");
    missing.forEach(function (m) { lines.push("    " + m); });
    lines.push("");
  }
  if (noNumber.length) { lines.push("学籍番号が空: " + noNumber.length + " 人"); noNumber.slice(0, 10).forEach(function (m) { lines.push("    " + m); }); }
  if (badNumber.length) { lines.push("学籍番号が8桁でない: " + badNumber.length + " 人"); badNumber.slice(0, 10).forEach(function (m) { lines.push("    " + m); }); }

  lines.push("");
  lines.push(missing.length === 0 && noNumber.length === 0 && badNumber.length === 0
    ? "全員分の学籍番号が取得できます。名簿シートへの書き込みに進めます。"
    : "上記を解消してから書き込んでください。");

  const text = lines.join("\n");
  Logger.log(text);
  ui.alert("学籍番号の突き合わせ\n\n" + text.slice(0, 3000));
}

// フォーム回答の学籍番号の入力ミスを、取得元シートに書き戻して直す。
//
// GAS側に補正表を持つ方法もあるが、同じ事実が「フォーム回答」と「コード」の2箇所に
// 存在することになり、片方だけ直したときに食い違う。取得元を正にして直す。
//
// 氏名と学籍番号（ログインの識別子）をコードに書くとリポジトリの履歴に永続的に残るため、
// この表には値を置かない。訂正が必要になったら取得元スプレッドシートを直接直すか、
// 一時的にここへ書いて実行したあと必ず空に戻し、コミットしないこと。
const STUDENT_NUMBER_FIXES = {};

function fixStudentNumbersInSource() {
  ui = ui || SpreadsheetApp.getUi();
  let ss, sheet = null;
  try {
    ss = SpreadsheetApp.openById(STUDENT_NUMBER_SOURCE_ID);
    ss.getSheets().forEach(function (s) { if (s.getSheetId() === STUDENT_NUMBER_SOURCE_GID) sheet = s; });
  } catch (e) {
    return ui.alert("取得元を開けません\n\n" + e.message);
  }
  if (!sheet) return ui.alert("gid=" + STUDENT_NUMBER_SOURCE_GID + " のシートが見つかりません");

  const lastRow = sheet.getLastRow();
  const names = sheet.getRange(SRC_DATA_START_ROW, SRC_COL_NAME, lastRow - SRC_DATA_START_ROW + 1, 1).getValues();
  const nums = sheet.getRange(SRC_DATA_START_ROW, SRC_COL_STUDENT_NUMBER, lastRow - SRC_DATA_START_ROW + 1, 1).getValues();

  // 対象行を特定し、現在値と変更後を突き合わせてから書く
  const plan = [];
  const wanted = {};
  Object.keys(STUDENT_NUMBER_FIXES).forEach(function (k) { wanted[normalizeForMatch_(k)] = STUDENT_NUMBER_FIXES[k]; });

  for (let i = 0; i < names.length; i++) {
    const key = normalizeForMatch_(names[i][0]);
    if (!wanted[key]) continue;
    const before = String(nums[i][0] || "").trim();
    const after = wanted[key];
    if (before === after) continue;   // 既に直っていれば触らない
    plan.push({ row: SRC_DATA_START_ROW + i, name: String(names[i][0]).trim(), before: before, after: after });
  }

  if (!plan.length) return ui.alert("修正対象はありません。既に訂正済みの可能性があります。");

  const confirm = ui.alert(
    "取得元「" + ss.getName() + "」の学籍番号を訂正します\n\n" +
    plan.map(function (p) {
      return p.row + "行 " + p.name + ": 「" + p.before + "」→「" + p.after + "」";
    }).join("\n") +
    "\n\n実行しますか？",
    ui.ButtonSet.OK_CANCEL);
  if (confirm === ui.Button.CANCEL) return;

  try {
    plan.forEach(function (p) {
      // 文字列として書き込む。数値だと表示形式によって桁が丸められる可能性がある
      sheet.getRange(p.row, SRC_COL_STUDENT_NUMBER).setValue(p.after);
    });
    SpreadsheetApp.flush();
    ui.alert("訂正しました\n\n" + plan.length + " 件\n\n" +
      "【調査】学籍番号が全員分引けるか調べる で確認してください。");
  } catch (error) {
    ui.alert("書き込みに失敗しました\n\n" + error.message +
      "\n\n取得元スプレッドシートの編集権限を確認してください。");
    Logger.log("Error: " + error.message);
  }
}

// ver0の名簿シートに学籍番号の列を追加する。
//
// DB送信は buildUserChanges_ がアンケート回答から直接引くため、この列は送信経路には
// 関与しない。人が名簿上で学籍番号を確認できるようにするための表示用。
// 数式(IMPORTRANGE)ではなく値で書くのは、ファイルペアごとの接続許可が要らず、
// コピー作成時に切れる問題も起きないため。学籍番号は年度中に変わらない。
const ROSTER_SHEET_NAME = "名簿";
const ROSTER_SHEET_NAME_COL = 1;      // A列: 名前
const ROSTER_SHEET_HEADER_ROW = 1;
const STUDENT_NUMBER_HEADER = "学籍番号";

function addStudentNumberColumnToRoster() {
  ui = ui || SpreadsheetApp.getUi();
  const sheet = SpreadsheetApp.getActiveSpreadsheet().getSheetByName(ROSTER_SHEET_NAME);
  if (!sheet) return ui.alert("「" + ROSTER_SHEET_NAME + "」シートが見つかりません");

  let srcMap;
  try {
    srcMap = loadStudentNumberMap_().map;
  } catch (e) {
    return ui.alert("取得元を読めません\n\n" + e.message);
  }

  // 既に学籍番号の列があればそこへ、無ければ右端の次の列に作る
  const lastCol = sheet.getLastColumn();
  const headers = sheet.getRange(ROSTER_SHEET_HEADER_ROW, 1, 1, lastCol).getValues()[0];
  let targetCol = -1;
  headers.forEach(function (h, i) {
    if (String(h || "").trim() === STUDENT_NUMBER_HEADER) targetCol = i + 1;
  });
  const isNew = targetCol === -1;
  if (isNew) targetCol = lastCol + 1;

  const lastRow = sheet.getLastRow();
  const names = sheet.getRange(ROSTER_SHEET_HEADER_ROW + 1, ROSTER_SHEET_NAME_COL,
    lastRow - ROSTER_SHEET_HEADER_ROW, 1).getValues();

  const values = [];
  const missing = [];
  let filled = 0;
  for (let i = 0; i < names.length; i++) {
    const raw = String(names[i][0] || "").trim();
    if (!raw) { values.push([""]); continue; }
    const src = srcMap[normalizeForMatch_(raw)];
    if (!src || !/^\d{8}$/.test(src.studentNumber)) {
      values.push([""]);
      missing.push((ROSTER_SHEET_HEADER_ROW + 1 + i) + "行 " + raw);
      continue;
    }
    // 文字列で書く。数値だと表示形式によって指数表記や桁落ちが起きうる
    values.push([src.studentNumber]);
    filled++;
  }

  const confirm = ui.alert(
    (isNew ? "名簿シートの " + colLetter_(targetCol) + " 列に「" + STUDENT_NUMBER_HEADER + "」を新設します"
           : "名簿シートの既存の " + colLetter_(targetCol) + " 列「" + STUDENT_NUMBER_HEADER + "」を更新します") +
    "\n\n対象 " + names.length + " 行 / 学籍番号を書ける " + filled + " 行" +
    "\n引けない " + missing.length + " 行" +
    (missing.length ? "\n    " + missing.slice(0, 8).join("\n    ") +
      (missing.length > 8 ? "\n    ...他 " + (missing.length - 8) : "") : "") +
    "\n\n実行しますか？",
    ui.ButtonSet.OK_CANCEL);
  if (confirm === ui.Button.CANCEL) return;

  try {
    sheet.getRange(ROSTER_SHEET_HEADER_ROW, targetCol).setValue(STUDENT_NUMBER_HEADER);
    sheet.getRange(ROSTER_SHEET_HEADER_ROW + 1, targetCol, values.length, 1).setValues(values);
    SpreadsheetApp.flush();
    ui.alert("書き込みました\n\n" + colLetter_(targetCol) + "列 / " + filled + " 行に学籍番号を設定");
  } catch (error) {
    ui.alert("書き込みに失敗しました\n\n" + error.message +
      "\n\n名簿シートが保護されている可能性があります。");
    Logger.log("Error: " + error.message);
  }
}

// タスクのGoogleドキュメントURLの取得元。
// A列にドキュメントへのリンクが並ぶ。どのタスクに対応するかを示す列を特定するために調べる。
const TASK_DOC_SOURCE_ID = "1a2pvM1M8NWQNLNaYnsqTzpB_oGE129-ibpbN3Be1Z1Q";
const TASK_DOC_SOURCE_GID = 129469480;

// 取得元シートの構造を調べる調査用。
// セルがハイパーリンク（表示文字とURLが別）の場合、getValues()では文字しか取れないため
// getRichTextValues()でリンク先も確認する。
function inspectTaskDocSource() {
  ui = ui || SpreadsheetApp.getUi();
  let ss, sheet = null;
  try {
    ss = SpreadsheetApp.openById(TASK_DOC_SOURCE_ID);
    ss.getSheets().forEach(function (s) { if (s.getSheetId() === TASK_DOC_SOURCE_GID) sheet = s; });
  } catch (e) {
    return ui.alert("取得元を開けません\n\n" + e.message);
  }
  if (!sheet) {
    return ui.alert("gid=" + TASK_DOC_SOURCE_GID + " が見つかりません\n\nシート一覧:\n" +
      ss.getSheets().map(function (s) { return s.getName() + " (gid=" + s.getSheetId() + ")"; })
        .join("\n").slice(0, 1200));
  }

  const lastRow = Math.min(sheet.getLastRow(), 8);
  const lastCol = Math.min(sheet.getLastColumn(), 12);
  const vals = sheet.getRange(1, 1, lastRow, lastCol).getValues();
  const rich = sheet.getRange(1, 1, lastRow, 1).getRichTextValues();   // A列のリンク先
  const formulas = sheet.getRange(1, 1, lastRow, 1).getFormulas();     // HYPERLINK数式の場合

  const lines = [];
  lines.push("ファイル: " + ss.getName());
  lines.push("シート: " + sheet.getName() + " (" + sheet.getLastRow() + "行 × " + sheet.getLastColumn() + "列)");
  lines.push("");
  for (let r = 0; r < lastRow; r++) {
    const cells = [];
    for (let c = 0; c < lastCol; c++) {
      const v = String(vals[r][c] || "").trim();
      cells.push(colLetter_(c + 1) + ":" + (v.length > 14 ? v.slice(0, 14) + "…" : v));
    }
    lines.push((r + 1) + "行 | " + cells.join(" | "));

    // A列のリンク先を別行で出す。表示文字とURLが違う場合にどちらを使うか判断するため
    let link = "";
    const rt = rich[r][0];
    if (rt) {
      link = rt.getLinkUrl() || "";
      if (!link && rt.getRuns) {
        rt.getRuns().forEach(function (run) { if (!link) link = run.getLinkUrl() || ""; });
      }
    }
    const f = String(formulas[r][0] || "");
    if (link) lines.push("      A列のリンク先: " + link.slice(0, 60) + "…");
    else if (f) lines.push("      A列の数式: " + f.slice(0, 60) + "…");
  }

  const text = lines.join("\n");
  Logger.log(text);
  ui.alert("タスクのドキュメント一覧\n\n" + text.slice(0, 3000));
}
