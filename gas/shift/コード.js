const lock = LockService.getScriptLock(); // スクリプト全体で共通のロック
// const lock = LockService.getDocumentLock(); // 対象ドキュメントに紐づくロック
// const lock = LockService.getUserLock(); // 実行ユーザーに紐づくロック

// uiを取得
var ui; // onOpen内で初期化(トップレベルでgetUi()するとエディタ実行が落ちるため)

// プロパティストアからベースURLを取得
const properties = PropertiesService.getScriptProperties();
const baseUrl = properties.getProperty("API_BASE_URL");

// 年度ID。スクリプトプロパティYEAR_IDで上書きできる（未設定時は45th）
const YEAR_ID = Number(properties.getProperty("YEAR_ID") || 45);

// 45thの日程シートのレイアウト（転置: 名前=列、時刻=行）
// 3行目B列以降にメンバー名、11行目以降のA列に時刻、その右側がタスク割り当てのグリッド
const NAME_ROW = 3;         // メンバー名の行
const GRID_START_ROW = 11;  // 時刻・割り当てが始まる行
const GRID_START_COL = 2;   // メンバー列の開始列（B列）
// 1日あたりの時間スロット数(64)は参加可否警告.jsのSLOT_COUNTを共有する。
// GASは全ファイルが同一スコープのため、ここで再宣言するとSyntaxErrorになり
// スクリプト全体が読み込めなくなる（メニューも表示されなくなる）

// スプシを開いた時に実行される関数
function onOpen() {
  // タブにメニューを追加する
  ui = SpreadsheetApp.getUi();
  const menu = ui.createMenu("SeeFT");
  menu
    // .addItem("サイドバーを開く", "showSidebar")
    // .addSeparator()
    .addSubMenu(
      ui.createMenu("シフトをSeeFTに送信する")
        .addItem("このシート全体", "updateShifts")
        .addItem("選択した列のメンバーのみ", "updateShiftsRange")
        .addSeparator()
        .addItem("【調査】送信せず件数だけ数える", "countShiftChanges")
        .addItem("【調査】名前の重複を調べる", "checkDuplicateNames")
        .addItem("【調査】このシートのタスク別セル数を数える", "countTasksInSheet")
        .addItem("【調査】タスク名でセル位置を探す", "findTaskCells")
        .addItem("【調査】コピペ残骸を検出する", "findCopyResidue")
        .addItem("【修正】コピペ残骸を削除する", "deleteCopyResidue")
        .addItem("【調査】準々備日の名前行を準備日と比べる", "compareNameRows")
        .addItem("【修正】準々備日の名前行を復元する", "fixJunjunbibiNameRow")
        .addSeparator()
        .addItem("【事前確認】名前がDBに全部あるか調べる", "checkNamesAgainstDB")
        .addItem("【事前確認】タスク名がDBに全部あるか調べる", "checkTaskNamesAgainstDB")
        .addItem("【リセット】最初から送り直す", "resetShiftResumePoints")
    )
    .addSubMenu(
      // シフトはユーザー名・タスク名で紐づくため、先に名簿とタスクを送る必要がある
      ui.createMenu("マスタをSeeFTに送信する（シフトより先に実行）")
        .addItem("1. 名簿を送信", "updateUsers")
        .addItem("2. タスクを送信", "updateTasksAndPlaces")
        .addSeparator()
        .addItem("【調査】名簿の行(名前/局/学年)を点検する", "compareRosterRows")
        .addItem("【調査】名簿シートの列構成を調べる", "inspectRosterSheet")
        .addItem("【調査】学籍番号の取得元を調べる", "inspectStudentNumberSource")
        .addItem("【調査】学籍番号が全員分引けるか調べる", "checkStudentNumberMatch")
        .addItem("【修正】学籍番号の入力ミスを訂正する", "fixStudentNumbersInSource")
        .addItem("【修正】名簿シートに学籍番号の列を追加", "addStudentNumberColumnToRoster")
        .addItem("【調査】タスクの管轄局を点検する", "checkTaskBureaus")
        .addItem("【調査】タスクのドキュメント一覧を調べる", "inspectTaskDocSource")
        .addItem("【事前確認】マニュアルURLの紐付けを点検する", "checkManualUrlMapping")
    )
    .addSeparator()
    .addSubMenu(
      ui.createMenu("条件付き書式を設定する")
      .addItem("全てのシートにまとめて設定(時間かかる)", "setConditionalFormatting")
      .addItem("このシートにだけ設定", "setConditionalFormattingCurrentSheet")
      .addItem("タスクシートに色を設定", "setConditionalFormattingToTaskSheet")
    )
    .addSeparator()
    .addItem("人数チェック（全日程）を更新", "buildAllStaffingMatrices")
    .addItem("このチェックシートだけ更新", "buildCurrentStaffingMatrix")
    .addItem("メンバーを追加", "addMemberColumns")
    .addItem("参加不可の警告を更新", "applyAvailabilityNotes")
    .addToUi();
}

// シート名から日付名を判定する。判定できない場合はnullを返す
// 「準々備日」は「準備日」を部分文字列として含まないが、判定順の誤りを避けるため先に評価する
function dateNameOf_(sheetName) {
  // 「人数チェック_準々備日」等の派生シートを日程シートと誤認しないよう除外する
  if (sheetName.indexOf("人数チェック") !== -1) return null;
  if (sheetName.indexOf("準々備日") !== -1) return "準々備日";
  if (sheetName.indexOf("準備日") !== -1) return "準備日";
  if (sheetName.indexOf("1日目") !== -1) return "1日目";
  if (sheetName.indexOf("2日目") !== -1) return "2日目";
  if (sheetName.indexOf("片付け日") !== -1) return "片付け日";
  return null;
}

// シート名から天気を判定する
function weatherOf_(sheetName) {
  return sheetName.indexOf("雨") !== -1 ? "雨" : "晴れ";
}

// 時刻セルの値からtimeIDを求める。
// timesテーブルは0:00を1として15分刻みで連番（6:00=25）
function timeStrToTimeID_(cell) {
  let minutes;
  if (cell instanceof Date) {
    minutes = cell.getHours() * 60 + cell.getMinutes();
  } else {
    const m = String(cell).match(/^(\d{1,2}):(\d{2})/);
    if (!m) return null;
    minutes = Number(m[1]) * 60 + Number(m[2]);
  }
  if (minutes % 15 !== 0) return null;
  return minutes / 15 + 1;
}

// 日程シートから送信用のchanges配列を作る。
// targetCols を渡すとその列（1始まりのシート上の列番号）のメンバーだけを対象にする。
function buildShiftChanges_(sheet, targetCols) {
  const sheetName = sheet.getName();
  const date = dateNameOf_(sheetName);
  if (!date) return null;
  const weather = weatherOf_(sheetName);

  const lastCol = sheet.getLastColumn();
  const peopleCount = lastCol - GRID_START_COL + 1;
  if (peopleCount < 1) return { date: date, weather: weather, changes: [] };

  const names = sheet.getRange(NAME_ROW, GRID_START_COL, 1, peopleCount).getValues()[0];
  const timeCells = sheet.getRange(GRID_START_ROW, 1, SLOT_COUNT, 1).getValues();
  const grid = sheet.getRange(GRID_START_ROW, GRID_START_COL, SLOT_COUNT, peopleCount);
  const values = grid.getValues();
  const backGrounds = grid.getBackgrounds();

  const changes = [];
  for (let p = 0; p < peopleCount; p++) {
    const userName = String(names[p] || "").trim();
    if (!userName) continue;
    if (targetCols && targetCols.indexOf(GRID_START_COL + p) === -1) continue;

    for (let t = 0; t < SLOT_COUNT; t++) {
      const timeID = timeStrToTimeID_(timeCells[t][0]);
      if (timeID === null) continue; // 時刻として読めない行は送らない

      // タスク名(セルの背景色が黒色の場合は'NG'にする)
      // 44thでは「休憩」を空白に潰していたが、モバイルで休憩カードを出すためそのまま送る。
      // 「誰が休憩中か」を見せない扱いはAPI側(createShiftCardFromGroup)で行う
      const taskName = backGrounds[t][p] != '#000000'
        ? String(values[t][p]).trim()
        : 'NG';

      // 空欄は送っても既存シフトを消せない(APIに削除の概念が無い)ので既定では送らない
      if (SKIP_EMPTY_CELLS && taskName === '') continue;

      changes.push({
        yearID: YEAR_ID,    // yearID
        timeID: timeID,     // timeID
        date: date,         // 日付
        weather: weather,   // 天気
        userName: userName, // ユーザー名
        taskName: taskName  // タスク名
      });
    }
  }
  return { date: date, weather: weather, changes: changes };
}

// 1リクエストで送るchangesの上限。
// APIはchange 1件につきSELECT+INSERT+action_log記録で2〜4クエリ実行するため、
// 実測では1件あたり約22ミリ秒かかっている(10,669件が4分弱で完了)。
// Cloudflareは100秒で接続を切る(エラー524)ので、1,500件(約34秒)にして余裕を66%確保する。
// 小さくしても総処理時間は変わらないため、安全側に倒している。
const SHIFT_CHUNK_SIZE = 1000;

// 割り当ての無い(空欄の)セルを送信対象から除くかどうか。
//
// falseのまま運用する。除外すると送信量は約61%減るが、次の2つの実害がある:
//   1. スプシで割り当てを取り消しても、DBに古いシフトが残り続ける
//      (APIに削除が無いため、空欄を送って「空タスク」で上書きするしか消す手段がない)
//   2. 送信済みのデータに歯抜けができると、表示側のgroupContinuousShiftsが
//      連続したtimeIDでカードをまとめる仕様のため、1枚であるべきカードが分裂して見える
// 送信量の増加はSHIFT_CHUNK_SIZEの分割と中断・再開でカバーする。
const SKIP_EMPTY_CELLS = false;

// GASの実行上限(6分=360秒)に対して、これを超えたら送信を打ち切って再開情報を残す閾値(ミリ秒)。
// 打ち切らずに上限に達すると「起動時間の最大値を超えました」で強制終了し、
// どこまで送れたか分からなくなるため、余裕を持って自主的に止める。
// 内訳: シート読み取り約20秒 + 送信260秒 + 超過判定後に走る最後のchunk最大34秒 = 約314秒。
// GAS上限360秒に対して46秒の余裕を残す。実測ペース(1件22ミリ秒)なら1回で約11,800件送れる。
const SHIFT_TIME_BUDGET_MS = 260 * 1000;

// 中断した位置を覚えておくためのキー(スクリプトプロパティ)
const SHIFT_RESUME_KEY = "SHIFT_RESUME";

// 組み立てたchangesをSeeFTに送信する。
// 件数が多い場合はCloudflareの100秒制限に収まるよう分割して送る。
// startIndex から送信を始め、時間切れになったら中断位置を返す。
// 戻り値: { sent: 今回送った件数, nextIndex: 次に送るべき位置(完了時はnull) }
function postShiftChanges_(changes, startIndex, onProgress) {
  const url = baseUrl + "/api/update_shifts";
  const begin = new Date().getTime();
  let i = startIndex || 0;
  let sent = 0;

  while (i < changes.length) {
    // 次のchunkを送ると時間切れになりそうなら、ここで打ち切って再開できるようにする
    if (new Date().getTime() - begin > SHIFT_TIME_BUDGET_MS) {
      Logger.log("時間切れのため中断: " + i + " / " + changes.length);
      return { sent: sent, nextIndex: i };
    }

    const chunk = changes.slice(i, i + SHIFT_CHUNK_SIZE);
    const options = {
      method: "post",
      contentType: "application/json",
      payload: JSON.stringify({ changes: chunk }),
      muteHttpExceptions: true // エラー時に本文を読んで原因を伝えるため
    };

    const t0 = new Date().getTime();
    const response = UrlFetchApp.fetch(url, options);
    const elapsed = (new Date().getTime() - t0) / 1000;
    const code = response.getResponseCode();
    const body = response.getContentText();
    Logger.log("chunk " + (Math.floor(i / SHIFT_CHUNK_SIZE) + 1) + " (" + chunk.length + "件) → HTTP " +
               code + " " + elapsed.toFixed(1) + "秒");

    if (code < 200 || code >= 300) {
      // 途中で失敗した場合、どこまで送れたかを伝える。
      // APIにトランザクションが無いため、失敗したchunk内も部分的に反映されている可能性がある
      throw new Error(
        "送信に失敗しました (HTTP " + code + ")\n" +
        "この実行で送信: " + sent + " 件 / 全体の " + i + " 件目まで完了\n" +
        "1件あたり" + (elapsed / chunk.length * 1000).toFixed(0) + "ミリ秒かかっています\n\n" +
        body.slice(0, 300));
    }

    i += chunk.length;
    sent += chunk.length;
    if (onProgress) onProgress(i, changes.length, elapsed);
  }
  return { sent: sent, nextIndex: null };
}

// 中断位置をシートごとにスクリプトプロパティへ保存する。
// 総件数も一緒に持ち、シートが編集されて件数が変わった場合は再開位置を無効にする
// （位置がずれたまま再開すると、別のメンバーのシフトを飛ばしてしまうため）
function saveResumePoint_(sheetName, nextIndex, total) {
  const store = JSON.parse(properties.getProperty(SHIFT_RESUME_KEY) || "{}");
  store[sheetName] = { index: nextIndex, total: total };
  properties.setProperty(SHIFT_RESUME_KEY, JSON.stringify(store));
}

function readResumePoint_(sheetName, total) {
  const store = JSON.parse(properties.getProperty(SHIFT_RESUME_KEY) || "{}");
  const saved = store[sheetName];
  if (!saved) return 0;
  if (saved.total !== total) {
    // シートが編集されて件数が変わった。位置がずれるので最初から送り直す
    Logger.log("件数が変わったため再開位置を破棄: " + saved.total + " → " + total);
    clearResumePoint_(sheetName);
    return 0;
  }
  return saved.index;
}

function clearResumePoint_(sheetName) {
  const store = JSON.parse(properties.getProperty(SHIFT_RESUME_KEY) || "{}");
  delete store[sheetName];
  properties.setProperty(SHIFT_RESUME_KEY, JSON.stringify(store));
}

// 中断位置を手動でリセットする（最初から送り直したいとき）
function resetShiftResumePoints() {
  ui = ui || SpreadsheetApp.getUi();
  properties.deleteProperty(SHIFT_RESUME_KEY);
  ui.alert("送信の再開位置をリセットしました。次回は最初から送信します。");
}

// スプシ編集時にDBのシフトを更新する関数
// 現在インストール型トリガーは未設定のため通常は呼ばれない
function onChange(e) {
  try {
    // ロックを試みる。最大待機時間を設定 (例: 30秒)
    // この時間は、API呼び出しにかかる最大時間などを考慮して調整
    lock.waitLock(30000);

    // --- ここからクリティカルセクション ---
    const sheet = e.source.getActiveSheet();
    const range = sheet.getActiveRange();

    // 編集された列に対応するメンバーだけを送信対象にする
    const targetCols = [];
    for (let c = range.getColumn(); c <= range.getLastColumn(); c++) {
      if (c >= GRID_START_COL) targetCols.push(c);
    }
    if (!targetCols.length) return;

    const built = buildShiftChanges_(sheet, targetCols);
    if (!built) return; // 日程シート以外では何もしない
    if (!built.changes.length) return;

    postShiftChanges_(built.changes, 0);
    // --- ここまでクリティカルセクション ---
  } catch (error) {
    // トリガー起動時はUIが無いためalertを出さずログのみに記録する
    Logger.log("Error: " + error.message);
  } finally {
    lock.releaseLock(); // 必ずロックを解放する
  }
}

// 現在のスプシをDBに反映する関数
function updateShifts() {
  ui = ui || SpreadsheetApp.getUi();
  try {
    // ロックを試みる。最大待機時間を設定 (例: 30秒)
    // この時間は、API呼び出しにかかる最大時間などを考慮して調整
    lock.waitLock(30000);

    // --- ここからクリティカルセクション ---
    const sheet = SpreadsheetApp.getActiveSpreadsheet().getActiveSheet();
    const sheetName = sheet.getName();  // シート名取得（dateIDとweatherIDに対応）

    if (!dateNameOf_(sheetName)) {
      return ui.alert(`シート名が不適です\n修正して再実行してください`);
    }

    const built = buildShiftChanges_(sheet);
    if (!built.changes.length) {
      return ui.alert(`送信対象のメンバーが見つかりませんでした`);
    }

    // 前回このシートの送信が時間切れで中断していれば、その位置から再開する
    const resume = readResumePoint_(sheetName, built.changes.length);
    const chunks = Math.ceil((built.changes.length - resume) / SHIFT_CHUNK_SIZE);
    const confirm = ui.alert(
      (resume > 0
        ? `前回の続き（` + resume + ` 件目以降）から送信します\n【 ` + sheetName + ` 】\n\n残り `
        : `以下のシートのシフトをSeeFTに送信してよろしいですか？\n【 ` + sheetName + ` 】\n\n`) +
      (built.changes.length - resume) + ` 件` +
      (chunks > 1 ? `（` + chunks + ` 回に分けて送信）` : ``) +
      `\n\n※GASの6分制限に達した場合は途中で中断し、次回この続きから再開できます`,
      ui.ButtonSet.OK_CANCEL);
    if (confirm === ui.Button.CANCEL) {
      Logger.log("送信がキャンセルされました");
      return
    }

    const result = postShiftChanges_(built.changes, resume, function (done, total, sec) {
      Logger.log("進捗: " + done + " / " + total + " 件 (直近のchunk " + sec.toFixed(1) + "秒)");
    });

    if (result.nextIndex === null) {
      clearResumePoint_(sheetName);
      ui.alert(`送信が完了しました\n【 ` + sheetName + ` 】\n\n` +
               `この実行で ` + result.sent + ` 件 / 全 ` + built.changes.length + ` 件`);
    } else {
      saveResumePoint_(sheetName, result.nextIndex, built.changes.length);
      ui.alert(`時間切れのため中断しました\n【 ` + sheetName + ` 】\n\n` +
               `この実行で ` + result.sent + ` 件送信\n` +
               `全体の ` + result.nextIndex + ` / ` + built.changes.length + ` 件まで完了\n\n` +
               `もう一度同じメニューを実行すると続きから再開します`);
    }
    // --- ここまでクリティカルセクション ---
  } catch (error) {
    ui.alert(`エラーが発生しました\n
      エラーを修正して再実行してください\n
      ------エラーメッセージ-----\n` +
      error.message +
      `\n---------------------------`);
    Logger.log("Error: " + error.message);
  } finally {
    lock.releaseLock(); // 必ずロックを解放する
  }
}

// 選択した列のメンバーのシフトだけをDBに反映する関数
// 45thは名前が列方向に並ぶため、行番号ではなく選択範囲の列を対象にする
function updateShiftsRange() {
  ui = ui || SpreadsheetApp.getUi();
  try {
    // ロックを試みる。最大待機時間を設定 (例: 30秒)
    // この時間は、API呼び出しにかかる最大時間などを考慮して調整
    lock.waitLock(30000);

    // --- ここからクリティカルセクション ---
    const sheet = SpreadsheetApp.getActiveSpreadsheet().getActiveSheet();
    const sheetName = sheet.getName();  // シート名取得（dateIDとweatherIDに対応）

    if (!dateNameOf_(sheetName)) {
      return ui.alert(`シート名が不適です\n修正して再実行してください`);
    }

    const range = sheet.getActiveRange();
    const targetCols = [];
    for (let c = range.getColumn(); c <= range.getLastColumn(); c++) {
      if (c >= GRID_START_COL) targetCols.push(c);
    }
    if (!targetCols.length) {
      return ui.alert(`送信したいメンバーの列を選択してから実行してください`);
    }

    // 選択列のメンバー名を確認用に取得する
    const names = targetCols.map(function (c) {
      return String(sheet.getRange(NAME_ROW, c).getValue() || "").trim();
    }).filter(function (n) { return n; });
    if (!names.length) {
      return ui.alert(`選択した列にメンバーが見つかりませんでした`);
    }

    // ダイアログで送信内容を確認
    const confirm = ui.alert(
      `以下のメンバーのシフトをSeeFTに送信してよろしいですか？\n【 ` + sheetName + ` 】\n` +
      names.join("、"),
      ui.ButtonSet.OK_CANCEL);
    if (confirm === ui.Button.CANCEL) {
      Logger.log("送信がキャンセルされました");
      return
    }

    const built = buildShiftChanges_(sheet, targetCols);
    if (!built.changes.length) {
      return ui.alert(`送信対象のシフトが見つかりませんでした`);
    }

    const rangeResult = postShiftChanges_(built.changes, 0, function (done, total, sec) {
      Logger.log("進捗: " + done + " / " + total + " 件 (直近のchunk " + sec.toFixed(1) + "秒)");
    });
    ui.alert(rangeResult.nextIndex === null
      ? `送信しました\n` + names.join("、") + `\n` + rangeResult.sent + ` 件`
      : `時間切れで中断しました\n` + rangeResult.sent + ` 件送信\n` +
        `選択範囲を狭めて再実行してください`);
    // --- ここまでクリティカルセクション ---
  } catch (error) {
    ui.alert(`エラーが発生しました\n
      エラーを修正して再実行してください\n
      ------エラーメッセージ-----\n` +
      error.message +
      `\n---------------------------`);
    Logger.log("Error: " + error.message);
  } finally {
    lock.releaseLock(); // 必ずロックを解放する
  }
}

function resetLock() {
  lock.releaseLock(); // ロックを解放する
}

// 送信せずに件数だけ数える調査用の関数。
// 空欄セルを除外した場合にどれだけ減るかを実データで測るために使う。
function countShiftChanges() {
  ui = ui || SpreadsheetApp.getUi();
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const lines = [];
  let totalAll = 0;
  let totalAssigned = 0;

  ["準々備日", "準備日", "1日目_晴れ", "1日目_雨", "2日目_晴れ", "2日目_雨", "片付け日"].forEach(function (sn) {
    const sheet = ss.getSheetByName(sn);
    if (!sheet) { lines.push(sn + ": シートなし"); return; }

    const built = buildShiftChanges_(sheet);
    if (!built) { lines.push(sn + ": 日程シートとして認識できず"); return; }

    // SKIP_EMPTY_CELLSが有効なら、この時点で空欄は既に除外されている
    const sending = built.changes.length;
    const ng = built.changes.filter(function (c) { return c.taskName === "NG"; }).length;
    const real = sending - ng;
    const chunks = Math.ceil(sending / SHIFT_CHUNK_SIZE);
    totalAll += sending;
    totalAssigned += real;
    lines.push(sn + ": 送信" + sending + "件 (NG=" + ng + " 実タスク=" + real + ") " +
               chunks + "回に分割 目安" + Math.ceil(sending / 100) + "秒");
  });

  lines.push("");
  lines.push("合計送信件数: " + totalAll + "件 (うち実タスク " + totalAssigned + "件)");
  lines.push("全シート送信時の目安: " + Math.ceil(totalAll / 100) + "秒");
  lines.push("");
  lines.push("設定: 1回あたり" + SHIFT_CHUNK_SIZE + "件 / 空欄除外=" + (SKIP_EMPTY_CELLS ? "有効" : "無効"));
  Logger.log(lines.join("\n"));
  ui.alert("シフト件数の実測\n\n" + lines.join("\n"));
}

// 列番号(1始まり)をA1記法の列文字に変換する
function colLetter_(n) {
  let s = "";
  while (n > 0) {
    const m = (n - 1) % 26;
    s = String.fromCharCode(65 + m) + s;
    n = (n - m - 1) / 26;
  }
  return s;
}

// 日程シートの3行目に同じ名前が複数の列で使われていないか調べる調査用の関数。
//
// APIはユーザー名でシフトの持ち主を引き当て、シフトの一意キーに
// user_id + date_id + time_id + weather_id を使う(列やタスクは含まれない)。
// そのため同じ名前の列が2つあると、後の列が先の列のシフトをINSERTではなくUPDATEで
// 上書きしてしまい、「送信件数」と「DBの行数」がその分だけずれる。
function checkDuplicateNames() {
  ui = ui || SpreadsheetApp.getUi();
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const lines = [];

  ["準々備日", "準備日", "1日目_晴れ", "1日目_雨", "2日目_晴れ", "2日目_雨", "片付け日"].forEach(function (sn) {
    const sheet = ss.getSheetByName(sn);
    if (!sheet) { lines.push(sn + ": シートなし"); return; }

    const peopleCount = sheet.getLastColumn() - GRID_START_COL + 1;
    if (peopleCount < 1) { lines.push(sn + ": メンバー列なし"); return; }

    const names = sheet.getRange(NAME_ROW, GRID_START_COL, 1, peopleCount).getValues()[0];
    const seen = {}; // 名前 -> その名前が現れた列文字の配列
    let filled = 0;
    for (let p = 0; p < peopleCount; p++) {
      const name = String(names[p] || "").trim();
      if (!name) continue;
      filled++;
      if (!seen[name]) seen[name] = [];
      seen[name].push(colLetter_(GRID_START_COL + p));
    }

    const distinct = Object.keys(seen);
    const dupNames = distinct.filter(function (n) { return seen[n].length > 1; });
    const extraCols = dupNames.reduce(function (acc, n) { return acc + seen[n].length - 1; }, 0);

    lines.push(sn + ": 名前のある列 " + filled + " / 実人数 " + distinct.length +
      " / 重複名 " + dupNames.length + "種 (余分な列 " + extraCols + " → 消えるシフト " +
      (extraCols * SLOT_COUNT) + " 件)");
    dupNames.slice(0, 12).forEach(function (n) {
      lines.push("    " + n + " → " + seen[n].join(", ") + " 列");
    });
    if (dupNames.length > 12) lines.push("    ...他 " + (dupNames.length - 12) + " 種はログ参照");
  });

  const text = lines.join("\n");
  Logger.log(text);
  ui.alert("3行目の名前の重複チェック\n\n" + text.slice(0, 1400));
}

// 準々備日の名前行が壊れている件の調査用。
// 準備日(正常なシート)の3行目と列ごとに突き合わせ、どこがどう食い違うかを報告する。
// 併せて、名前が潰れている列に割り当てデータが残っているかも数える。
function compareNameRows() {
  ui = ui || SpreadsheetApp.getUi();
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const base = ss.getSheetByName("準備日");
  const target = ss.getSheetByName("準々備日");
  if (!base || !target) return ui.alert("準備日 または 準々備日 が見つかりません");

  const baseCount = base.getLastColumn() - GRID_START_COL + 1;
  const targetCount = target.getLastColumn() - GRID_START_COL + 1;
  const n = Math.min(baseCount, targetCount);

  const baseNames = base.getRange(NAME_ROW, GRID_START_COL, 1, baseCount).getValues()[0];
  const targetNames = target.getRange(NAME_ROW, GRID_START_COL, 1, targetCount).getValues()[0];
  const targetGrid = target.getRange(GRID_START_ROW, GRID_START_COL, SLOT_COUNT, targetCount).getValues();

  const lines = [];
  lines.push("準備日の列数 " + baseCount + " / 準々備日の列数 " + targetCount);

  let same = 0;
  const diffs = [];
  for (let p = 0; p < n; p++) {
    const b = String(baseNames[p] || "").trim();
    const t = String(targetNames[p] || "").trim();
    if (b === t) { same++; continue; }
    diffs.push({ col: colLetter_(GRID_START_COL + p), base: b, target: t, idx: p });
  }
  lines.push("3行目が一致した列 " + same + " / 食い違った列 " + diffs.length);

  if (diffs.length) {
    const first = diffs[0], last = diffs[diffs.length - 1];
    lines.push("食い違いの範囲: " + first.col + " 列 〜 " + last.col + " 列");
    lines.push("  先頭 " + first.col + ": 準備日「" + first.base + "」/ 準々備日「" + first.target + "」");
    lines.push("  末尾 " + last.col + ": 準備日「" + last.base + "」/ 準々備日「" + last.target + "」");

    // 食い違う範囲が連続した1ブロックかどうか(飛び地があると単純なコピーでは直せない)
    const contiguous = (last.idx - first.idx + 1) === diffs.length;
    lines.push("  連続した1ブロックか: " + (contiguous ? "はい" : "いいえ(飛び地あり)"));

    // 名前が潰れた列に割り当てが残っているか。残っていれば名前行を直すだけで復旧できる
    let colsWithData = 0, cellsWithData = 0;
    diffs.forEach(function (d) {
      let c = 0;
      for (let t = 0; t < SLOT_COUNT; t++) {
        if (String(targetGrid[t][d.idx] || "").trim() !== "") c++;
      }
      if (c > 0) colsWithData++;
      cellsWithData += c;
    });
    lines.push("  食い違う列のうち割り当てが入っている列 " + colsWithData + " / " + diffs.length +
      " (計 " + cellsWithData + " セル)");
  }

  // 準備日側の名前がそもそも一意かどうかも確認しておく
  const seen = {};
  baseNames.forEach(function (v) {
    const s = String(v || "").trim();
    if (s) seen[s] = (seen[s] || 0) + 1;
  });
  const dup = Object.keys(seen).filter(function (k) { return seen[k] > 1; });
  lines.push("準備日の重複名 " + dup.length + "種" + (dup.length ? ": " + dup.slice(0, 5).join("、") : ""));

  const text = lines.join("\n");
  Logger.log(text);
  ui.alert("名前行の突き合わせ（準備日 vs 準々備日）\n\n" + text.slice(0, 1400));
}

// 準々備日の壊れた名前行(C〜BV列が全て「齊藤　翔太」になっている)を準備日から復元する。
//
// 前提として compareNameRows() で次を確認済み:
//   - 食い違いは連続した1ブロックで、その前後の列は準備日と一致している(列順は保たれている)
//   - 準備日の3行目に重複名が無い
// 書き込む前に対象範囲を読み直して同じ条件が成り立つか検証し、成り立たなければ中止する。
function fixJunjunbibiNameRow() {
  ui = ui || SpreadsheetApp.getUi();
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const base = ss.getSheetByName("準備日");
  const target = ss.getSheetByName("準々備日");
  if (!base || !target) return ui.alert("準備日 または 準々備日 が見つかりません");

  const baseCount = base.getLastColumn() - GRID_START_COL + 1;
  const targetCount = target.getLastColumn() - GRID_START_COL + 1;
  if (baseCount !== targetCount) {
    return ui.alert("列数が違います(準備日 " + baseCount + " / 準々備日 " + targetCount +
      ")。列の対応が取れないため中止します。");
  }

  const baseNames = base.getRange(NAME_ROW, GRID_START_COL, 1, baseCount).getValues()[0];
  const targetNames = target.getRange(NAME_ROW, GRID_START_COL, 1, targetCount).getValues()[0];

  // 食い違う列を洗い出す
  const diffIdx = [];
  for (let p = 0; p < baseCount; p++) {
    if (String(baseNames[p] || "").trim() !== String(targetNames[p] || "").trim()) diffIdx.push(p);
  }
  if (!diffIdx.length) return ui.alert("食い違いはありません。修正は不要です。");

  const first = diffIdx[0], last = diffIdx[diffIdx.length - 1];
  if (last - first + 1 !== diffIdx.length) {
    return ui.alert("食い違いが連続していません(飛び地あり)。\n" +
      "列の対応がずれている可能性があるため、自動修正は行いません。");
  }

  // コピー元に重複が無いことを確認する。重複があるとシフトが上書きされるため直す意味が無い
  const seen = {};
  baseNames.forEach(function (v) {
    const s = String(v || "").trim();
    if (s) seen[s] = (seen[s] || 0) + 1;
  });
  const dup = Object.keys(seen).filter(function (k) { return seen[k] > 1; });
  if (dup.length) {
    return ui.alert("準備日の3行目に重複名があります(" + dup.slice(0, 5).join("、") +
      ")。コピー元として使えないため中止します。");
  }

  const width = diffIdx.length;
  const srcCol = GRID_START_COL + first;
  const srcRange = base.getRange(NAME_ROW, srcCol, 1, width);
  const dstRange = target.getRange(NAME_ROW, srcCol, 1, width);

  // 数式が絡んでいないかを確認する(名前行が数式なら値で潰すのは適切でない)
  const srcFormulas = srcRange.getFormulas()[0].filter(function (f) { return f !== ""; }).length;
  const dstFormulas = dstRange.getFormulas()[0].filter(function (f) { return f !== ""; }).length;

  const confirm = ui.alert(
    "準々備日の名前行を準備日から復元します\n\n" +
    "対象: " + colLetter_(srcCol) + "3 〜 " + colLetter_(srcCol + width - 1) + "3 の " + width + " 列\n" +
    "先頭 " + colLetter_(srcCol) + ": 「" + String(targetNames[first]).trim() + "」→「" +
    String(baseNames[first]).trim() + "」\n" +
    "末尾 " + colLetter_(srcCol + width - 1) + ": 「" + String(targetNames[last]).trim() + "」→「" +
    String(baseNames[last]).trim() + "」\n\n" +
    "数式: コピー元 " + srcFormulas + " 個 / 書き込み先 " + dstFormulas + " 個\n" +
    (srcFormulas || dstFormulas ? "※数式が含まれます。値で上書きされます\n" : "") +
    "\n書き込んでよろしいですか？",
    ui.ButtonSet.OK_CANCEL);
  if (confirm === ui.Button.CANCEL) return;

  dstRange.setValues([baseNames.slice(first, first + width)]);
  SpreadsheetApp.flush();

  // 書き込み結果を読み直して検証する
  const after = target.getRange(NAME_ROW, GRID_START_COL, 1, targetCount).getValues()[0];
  let mismatch = 0;
  const afterSeen = {};
  for (let p = 0; p < targetCount; p++) {
    const s = String(after[p] || "").trim();
    if (String(baseNames[p] || "").trim() !== s) mismatch++;
    if (s) afterSeen[s] = (afterSeen[s] || 0) + 1;
  }
  const afterDup = Object.keys(afterSeen).filter(function (k) { return afterSeen[k] > 1; });

  // 名前が変わった以上、前回の送信位置は無意味なので破棄する
  clearResumePoint_("準々備日");

  ui.alert("復元しました\n\n" +
    "書き込んだ列 " + width + "\n" +
    "準備日との食い違い " + mismatch + " 列\n" +
    "準々備日の重複名 " + afterDup.length + " 種 / 実人数 " + Object.keys(afterSeen).length + "\n\n" +
    (mismatch === 0 && afterDup.length === 0
      ? "正常です。準々備日のシフトを送り直してください。\n(再開位置は破棄したので最初から送信されます)"
      : "想定外の状態です。送信前に確認してください。"));
}

// シフト送信前に、シート上の名前が全てDBのusersに存在するかを確かめる。
//
// APIはchangeを1件ずつ処理し、名前が引けなかった時点でエラーを返してその先を捨てる。
// 22,000件を数分かけて送った末に落ちるのを避けるため、送信前にここで突き合わせる。
// DB側の名前はGETした結果をそのまま使う(APIはユーザー名を正規化せず完全一致で引くため、
// 全角半角スペースの違いもここで検出できる)。
function checkNamesAgainstDB() {
  ui = ui || SpreadsheetApp.getUi();
  if (!baseUrl) return ui.alert("スクリプトプロパティ API_BASE_URL が未設定です");

  const response = UrlFetchApp.fetch(baseUrl + "/users", { muteHttpExceptions: true });
  const code = response.getResponseCode();
  if (code < 200 || code >= 300) {
    return ui.alert("usersの取得に失敗しました (HTTP " + code + ")\n" +
      response.getContentText().slice(0, 300));
  }
  const dbUsers = JSON.parse(response.getContentText()) || [];
  const dbNames = {};
  dbUsers.forEach(function (u) {
    const n = String(u.name || "").trim();
    if (n) dbNames[n] = (dbNames[n] || 0) + 1;
  });
  const dbDup = Object.keys(dbNames).filter(function (k) { return dbNames[k] > 1; });

  const lines = [];
  lines.push("DBのusers " + dbUsers.length + " 件 / 名前ユニーク " + Object.keys(dbNames).length +
    (dbDup.length ? " / 重複名 " + dbDup.length + "種: " + dbDup.slice(0, 3).join("、") : ""));

  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const missingAll = {};
  ["準々備日", "準備日", "1日目_晴れ", "1日目_雨", "2日目_晴れ", "2日目_雨", "片付け日"].forEach(function (sn) {
    const sheet = ss.getSheetByName(sn);
    if (!sheet) { lines.push(sn + ": シートなし"); return; }
    const peopleCount = sheet.getLastColumn() - GRID_START_COL + 1;
    if (peopleCount < 1) { lines.push(sn + ": メンバー列なし"); return; }

    const names = sheet.getRange(NAME_ROW, GRID_START_COL, 1, peopleCount).getValues()[0];
    const missing = [];
    let filled = 0;
    for (let p = 0; p < peopleCount; p++) {
      const n = String(names[p] || "").trim();
      if (!n) continue;
      filled++;
      if (!dbNames[n]) {
        missing.push(colLetter_(GRID_START_COL + p) + ":" + n);
        missingAll[n] = true;
      }
    }
    lines.push(sn + ": 名前 " + filled + " 件 / DBに無い " + missing.length + " 件");
    missing.slice(0, 8).forEach(function (m) { lines.push("    " + m); });
    if (missing.length > 8) lines.push("    ...他 " + (missing.length - 8) + " 件はログ参照");
  });

  const total = Object.keys(missingAll).length;
  lines.push("");
  lines.push(total === 0
    ? "全てDBに存在します。シフトを送信できます。"
    : "DBに無い名前が " + total + " 人います。先に「1. 名簿を送信」を実行するか、\n" +
      "スプシ側の表記を名簿と揃えてください。");

  const text = lines.join("\n");
  Logger.log(text);
  ui.alert("シート名簿とDBの突き合わせ\n\n" + text.slice(0, 1400));
}

// 日程シートのタスク別セル数を数える。DBの実測と突き合わせて取り込みの正しさを検証する用途。
//
// 総数だけの一致では「別のタスクが同数入れ替わっている」可能性を排除できないため、
// タスク名ごとの内訳まで比べる。集計はbuildShiftChanges_の結果をそのまま使い、
// 送信ロジックと同じ変換(黒背景→NG)を必ず経由させる。
function countTasksInSheet() {
  ui = ui || SpreadsheetApp.getUi();
  const sheet = SpreadsheetApp.getActiveSpreadsheet().getActiveSheet();
  const built = buildShiftChanges_(sheet);
  if (!built) return ui.alert("日程シートとして認識できません");

  const c = {};
  built.changes.forEach(function (ch) {
    const k = ch.taskName === "" ? "(空欄)" : ch.taskName;
    c[k] = (c[k] || 0) + 1;
  });

  const keys = Object.keys(c).sort(function (a, b) { return c[b] - c[a]; });
  const lines = [];
  lines.push("シート【" + sheet.getName() + "】のタスク別セル数");
  keys.forEach(function (k) {
    lines.push("  " + String(c[k]).padStart(6) + "  " + k);
  });
  lines.push("  合計 " + built.changes.length + " / タスク種類 " + keys.length);

  const text = lines.join("\n");
  Logger.log(text);
  ui.alert("タスク別セル数\n\n" + text.slice(0, 1400));
}

// シフト送信前に、シフト表のタスク名がDBのtasksに存在するかを確かめる。
//
// APIはタスク名が引けない場合、エラーにせず既定値(集合場所1/執行部/最大人数1/色000000)で
// 勝手にタスクを新規作成する。送信は成功扱いになるため、タスク一覧の記入漏れや
// シフト表のタイポが表に出ない。送信前にここで検出する。
function checkTaskNamesAgainstDB() {
  ui = ui || SpreadsheetApp.getUi();
  if (!baseUrl) return ui.alert("スクリプトプロパティ API_BASE_URL が未設定です");

  const response = UrlFetchApp.fetch(baseUrl + "/tasks", { muteHttpExceptions: true });
  const code = response.getResponseCode();
  if (code < 200 || code >= 300) {
    return ui.alert("tasksの取得に失敗しました (HTTP " + code + ")\n" +
      response.getContentText().slice(0, 300));
  }
  const dbTasks = JSON.parse(response.getContentText()) || [];
  const dbNames = {};
  dbTasks.forEach(function (t) {
    // APIは照合前に全角スペースを半角に寄せる。同じ扱いで突き合わせる
    const n = String(t.task || "").replace(/　/g, " ");
    dbNames[n] = true;
  });

  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const lines = [];
  lines.push("DBのtasks " + dbTasks.length + " 件");
  const missingAll = {};   // タスク名 -> 出現したシートと件数

  ["準々備日", "準備日", "1日目_晴れ", "1日目_雨", "2日目_晴れ", "2日目_雨", "片付け日"].forEach(function (sn) {
    const sheet = ss.getSheetByName(sn);
    if (!sheet) { lines.push(sn + ": シートなし"); return; }
    const built = buildShiftChanges_(sheet);
    if (!built) { lines.push(sn + ": 日程シートとして認識できず"); return; }

    const missing = {};
    built.changes.forEach(function (ch) {
      const n = String(ch.taskName || "").replace(/　/g, " ");
      if (n === "") return;                 // 空欄は既定のタスクに載るので対象外
      if (dbNames[n]) return;
      missing[n] = (missing[n] || 0) + 1;
      if (!missingAll[n]) missingAll[n] = [];
      missingAll[n].push(sn);
    });

    const keys = Object.keys(missing);
    lines.push(sn + ": DBに無いタスク名 " + keys.length + " 種 (" +
      keys.reduce(function (a, k) { return a + missing[k]; }, 0) + " セル)");
    // 打ち切らずに全種類を出す。上位数件だけだと欠落の全体像が掴めず、
    // 局タスクファイルを直す作業を何度もやり直すことになる
    keys.sort(function (a, b) { return missing[b] - missing[a]; });
    keys.forEach(function (k) { lines.push("    " + k + " × " + missing[k]); });
  });

  const total = Object.keys(missingAll).length;
  lines.push("");
  lines.push(total === 0
    ? "全てDBに存在します。勝手なタスク生成は起きません。"
    : "DBに無いタスク名が " + total + " 種あります。\n" +
      "このまま送るとAPIが既定値で自動生成します(局=執行部/最大人数1/色=黒)。\n" +
      "タスク一覧に追加するか、シフト表の表記を既存タスクに合わせてください。");

  const text = lines.join("\n");
  Logger.log(text);
  ui.alert("シフト表のタスク名とDBの突き合わせ\n\n" + text.slice(0, 3000));
}

// タスク名を指定して、それが書かれているセルの位置(シート・担当者・時刻)を一覧する。
//
// 件数だけでは「意味のある割り当て」か「コピペの残骸」かを判別できない。
// 誰の列の何時かが分かれば、人が見て一目で判断できる。
function findTaskCells() {
  ui = ui || SpreadsheetApp.getUi();
  const res = ui.prompt("探すタスク名を入力してください（完全一致）", ui.ButtonSet.OK_CANCEL);
  if (res.getSelectedButton() !== ui.Button.OK) return;
  const target = String(res.getResponseText() || "").trim();
  if (!target) return ui.alert("タスク名が空です");

  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const lines = [];
  let total = 0;

  ["準々備日", "準備日", "1日目_晴れ", "1日目_雨", "2日目_晴れ", "2日目_雨", "片付け日"].forEach(function (sn) {
    const sheet = ss.getSheetByName(sn);
    if (!sheet) return;
    const peopleCount = sheet.getLastColumn() - GRID_START_COL + 1;
    if (peopleCount < 1) return;

    const names = sheet.getRange(NAME_ROW, GRID_START_COL, 1, peopleCount).getValues()[0];
    const timeCells = sheet.getRange(GRID_START_ROW, 1, SLOT_COUNT, 1).getValues();
    const values = sheet.getRange(GRID_START_ROW, GRID_START_COL, SLOT_COUNT, peopleCount).getValues();

    for (let c = 0; c < peopleCount; c++) {
      // 連続するコマは1行にまとめる。1コマだけの割り当ては残骸の可能性が高いため目立たせたい
      let runStart = -1;
      for (let r = 0; r <= SLOT_COUNT; r++) {
        const hit = r < SLOT_COUNT && String(values[r][c]).trim() === target;
        if (hit && runStart === -1) runStart = r;
        if (!hit && runStart !== -1) {
          const len = r - runStart;
          total += len;
          const t0 = timeCells[runStart][0], t1 = timeCells[r - 1][0];
          const fmt = function (v) {
            if (v instanceof Date) return Utilities.formatDate(v, "Asia/Tokyo", "HH:mm");
            return String(v);
          };
          lines.push("  " + sn + " / " + colLetter_(GRID_START_COL + c) + "列 " +
            String(names[c]).trim() + " / " + fmt(t0) + "〜" + fmt(t1) +
            " (" + len + "コマ)" + (len === 1 ? "  ← 1コマのみ" : ""));
          runStart = -1;
        }
      }
    }
  });

  const text = "「" + target + "」が書かれたセル: 計 " + total + "\n" +
    (lines.length ? lines.join("\n") : "  見つかりませんでした");
  Logger.log(text);
  ui.alert("タスクのセル位置\n\n" + text.slice(0, 3000));
}

// シート複製時の消し漏れ(コピペ残骸)を洗い出す。
//
// 「タスク名がDBに無い」検査ではDBに存在するタスクの残骸を見つけられない。
// 残骸は「基準シートと同じ座標に同じ値がある」という形で現れるため、
// 座標と値の一致で検出する。353人×64コマの中で座標も値も一致するのは
// 偶然では起こりにくく、複製の痕跡と考えてよい。
function findCopyResidue() {
  ui = ui || SpreadsheetApp.getUi();
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const BASE = "準々備日";
  const base = ss.getSheetByName(BASE);
  if (!base) return ui.alert(BASE + " が見つかりません");

  const baseCount = base.getLastColumn() - GRID_START_COL + 1;
  const baseVals = base.getRange(GRID_START_ROW, GRID_START_COL, SLOT_COUNT, baseCount).getValues();
  const timeCells = base.getRange(GRID_START_ROW, 1, SLOT_COUNT, 1).getValues();
  const fmtTime = function (v) {
    if (v instanceof Date) return Utilities.formatDate(v, "Asia/Tokyo", "HH:mm");
    return String(v);
  };

  const lines = [];
  lines.push("基準シート: " + BASE + "（この座標・この値と一致するセルを残骸候補とみなす）");
  lines.push("");

  ["準備日", "1日目_晴れ", "1日目_雨", "2日目_晴れ", "2日目_雨", "片付け日"].forEach(function (sn) {
    const sheet = ss.getSheetByName(sn);
    if (!sheet) { lines.push(sn + ": シートなし"); return; }
    const n = Math.min(sheet.getLastColumn() - GRID_START_COL + 1, baseCount);
    if (n < 1) { lines.push(sn + ": メンバー列なし"); return; }

    const vals = sheet.getRange(GRID_START_ROW, GRID_START_COL, SLOT_COUNT, n).getValues();
    const byTask = {};
    let hits = 0, minRow = 999, maxRow = -1, minCol = 999, maxCol = -1;

    for (let r = 0; r < SLOT_COUNT; r++) {
      for (let c = 0; c < n; c++) {
        const b = String(baseVals[r][c] || "").trim();
        if (!b) continue;                                  // 基準側が空欄なら比較しない
        if (String(vals[r][c] || "").trim() !== b) continue; // 値が違えば残骸ではない
        hits++;
        byTask[b] = (byTask[b] || 0) + 1;
        if (r < minRow) minRow = r; if (r > maxRow) maxRow = r;
        if (c < minCol) minCol = c; if (c > maxCol) maxCol = c;
      }
    }

    if (!hits) { lines.push(sn + ": 一致なし"); return; }
    lines.push(sn + ": " + hits + " セルが一致  範囲 " +
      colLetter_(GRID_START_COL + minCol) + (GRID_START_ROW + minRow) + " 〜 " +
      colLetter_(GRID_START_COL + maxCol) + (GRID_START_ROW + maxRow) +
      " (" + fmtTime(timeCells[minRow][0]) + "〜" + fmtTime(timeCells[maxRow][0]) + ")");
    Object.keys(byTask).sort(function (a, b) { return byTask[b] - byTask[a]; })
      .forEach(function (k) { lines.push("    " + k + " × " + byTask[k]); });
  });

  const text = lines.join("\n");
  Logger.log(text);
  ui.alert("コピペ残骸の検出\n\n" + text.slice(0, 3000));
}

// 残骸を消す対象から除外するタスク。
// 本部指揮・実働指揮は執行部が毎日同じ時間帯に担当するため、準々備日と座標が一致しても
// 実データである(部分一致にとどまる・晴雨の両方にある、という特徴で見分けた)。
const RESIDUE_KEEP_TASKS = ["本部指揮", "実働指揮"];

// 準々備日からの複製残骸を削除する。
// 対象は「準々備日と同じ座標に同じ値」かつ RESIDUE_KEEP_TASKS 以外のセル。
// 消す前に必ず一覧を出して確認を取る。
function deleteCopyResidue() {
  ui = ui || SpreadsheetApp.getUi();
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const base = ss.getSheetByName("準々備日");
  if (!base) return ui.alert("準々備日 が見つかりません");

  const baseCount = base.getLastColumn() - GRID_START_COL + 1;
  const baseVals = base.getRange(GRID_START_ROW, GRID_START_COL, SLOT_COUNT, baseCount).getValues();

  // 対象は雨シートと片付け日のみ。準備日・晴れシートは指揮以外の一致が無いため触らない
  const TARGETS = ["1日目_雨", "2日目_雨", "片付け日"];
  const plan = [];   // {sheet, a1[], detail{}}

  TARGETS.forEach(function (sn) {
    const sheet = ss.getSheetByName(sn);
    if (!sheet) return;
    const n = Math.min(sheet.getLastColumn() - GRID_START_COL + 1, baseCount);
    if (n < 1) return;
    const vals = sheet.getRange(GRID_START_ROW, GRID_START_COL, SLOT_COUNT, n).getValues();

    const a1 = [], detail = {};
    for (let r = 0; r < SLOT_COUNT; r++) {
      for (let c = 0; c < n; c++) {
        const b = String(baseVals[r][c] || "").trim();
        if (!b) continue;
        if (RESIDUE_KEEP_TASKS.indexOf(b) !== -1) continue;
        if (String(vals[r][c] || "").trim() !== b) continue;
        a1.push(sheet.getRange(GRID_START_ROW + r, GRID_START_COL + c).getA1Notation());
        detail[b] = (detail[b] || 0) + 1;
      }
    }
    if (a1.length) plan.push({ sheet: sheet, name: sn, a1: a1, detail: detail });
  });

  if (!plan.length) return ui.alert("削除対象はありません。既に消えている可能性があります。");

  const total = plan.reduce(function (a, p) { return a + p.a1.length; }, 0);
  const lines = plan.map(function (p) {
    return p.name + ": " + p.a1.length + " セル\n    " +
      Object.keys(p.detail).sort(function (a, b) { return p.detail[b] - p.detail[a]; })
        .map(function (k) { return k + "×" + p.detail[k]; }).join("\n    ");
  });

  const confirm = ui.alert(
    "準々備日からの複製残骸を削除します\n\n" + lines.join("\n") +
    "\n\n合計 " + total + " セル\n" +
    "（本部指揮・実働指揮は実データのため対象外）\n\n実行しますか？",
    ui.ButtonSet.OK_CANCEL);
  if (confirm === ui.Button.CANCEL) return;

  try {
    let done = 0;
    plan.forEach(function (p) {
      p.sheet.getRangeList(p.a1).clearContent();
      done += p.a1.length;
    });
    SpreadsheetApp.flush();
    ui.alert("削除しました\n\n" + done + " セル\n\n" +
      "【調査】コピペ残骸を検出する で、指揮以外の一致が消えたことを確認してください。");
  } catch (error) {
    ui.alert("削除に失敗しました\n\n" + error.message +
      "\n\n保護範囲に当たっている可能性があります。");
    Logger.log("Error: " + error.message);
  }
}
