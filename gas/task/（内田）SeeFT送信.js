// 年度。年度が変わったらここを直す
const YEAR_ID = 45;

// タスクシートの名前。YEAR_ID から組み立てるので年度と食い違わない
const TASK_SHEET_NAMES = [
  `${YEAR_ID}thTask_準々備日`,
  `${YEAR_ID}thTask_準備日`,
  `${YEAR_ID}thTask_当日1日目`,
  `${YEAR_ID}thTask_当日2日目`,
  `${YEAR_ID}thTask_片付け日`,
];

// 列番号(1列目が0)
const COL_NAME = 0;        // タスク名
const COL_BUREAU = 5;      // 管轄局
const COL_PLACE = 7;       // 集合場所
const COL_MAX_MEMBER = 11; // 最大人数
const COL_URL = 12;        // マニュアル(ドキュメント版)のURL
const COL_MANUAL_URL = 13; // マニュアル(スライド版)のURL

// 一括取得する列数。列を足したときに取得範囲の広げ忘れが起きないよう定数から導出する
const COLUMN_COUNT = COL_MANUAL_URL + 1;

// 見出し行のタスク名。編集範囲に見出しが混ざったときに弾くために使う
const HEADER_TASK_NAME = 'シフト名';

const DEFAULT_PLACE = '本部(電気棟1F)';
const DEFAULT_MAX_MEMBER = 1;
const LOCK_WAIT_MS = 30000;

// 送信対象のタスクシートかどうか
function isTaskSheet(sheetName) {
  return TASK_SHEET_NAMES.indexOf(sheetName) !== -1;
}

// セルの値を文字列として取り出す。未入力や、その列がまだ存在しない場合は空文字にする。
// String(undefined) は "undefined" という truthy な文字列になり、
// 呼び出し側の `|| デフォルト値` をすり抜けてしまうため、ここで吸収する。
function cellText(valueRow, colIndex) {
  const value = valueRow ? valueRow[colIndex] : undefined;
  return (value === null || value === undefined) ? '' : String(value).trim();
}

// セルからURLを取り出す。
// ハイパーリンクが張られていればそのリンク先を、数式やベタ書きで文字列として
// 入っていればその値を使う。manual_url は対応表から VLOOKUP で引く運用のため
// リンクではなく文字列になる。getLinkUrl() だけだと常に空になってしまう。
function extractUrl(richTextRow, valueRow, colIndex) {
  const richTextCell = richTextRow ? richTextRow[colIndex] : null;
  const linkUrl = richTextCell ? richTextCell.getLinkUrl() : null;
  if (linkUrl) return String(linkUrl).trim();

  const text = cellText(valueRow, colIndex);
  return /^https?:\/\//.test(text) ? text : '';
}

// 行データから送信用の変更内容を組み立てる
function buildChanges(values, richTextValues, startIndex) {
  const changes = [];
  for (let i = startIndex; i < values.length; i++) {
    const taskName = cellText(values[i], COL_NAME);
    // タスク名が入力されていないか、見出し行の場合はスキップ
    if (!taskName || taskName === HEADER_TASK_NAME) continue;

    const maxMemberCell = values[i][COL_MAX_MEMBER];
    const change = {
      yearID: YEAR_ID,
      taskName: taskName,
      bureau: cellText(values[i], COL_BUREAU), // 空ならAPI側で執行部になる
      place: cellText(values[i], COL_PLACE) || DEFAULT_PLACE,
      url: extractUrl(richTextValues[i], values[i], COL_URL),
      manualUrl: extractUrl(richTextValues[i], values[i], COL_MANUAL_URL),
      maxMember: (maxMemberCell && typeof maxMemberCell === 'number') ? maxMemberCell : DEFAULT_MAX_MEMBER,
    };

    changes.push(change);
    console.log(change);
  }
  return changes;
}

// 組み立てた変更内容をAPIへ送信する
function postChanges(changes) {
  // プロパティストアからベースURLを取得
  const properties = PropertiesService.getScriptProperties();
  const baseUrl = properties.getProperty("API_BASE_URL");

  const url = baseUrl + "/api/update_tasks_and_places";
  const options = {
    method: "post",
    contentType: "application/json",
    payload: JSON.stringify({ changes: changes }),
    muteHttpExceptions: false
  };

  const response = UrlFetchApp.fetch(url, options);
  Logger.log("Response from API: " + response.getContentText());
}

// エラーをログに残し、可能ならダイアログでも知らせる
function alertError(error) {
  // UIが使えない文脈でも失われないよう、先にログへ出す
  Logger.log("Error: " + error.message);
  try {
    const ui = SpreadsheetApp.getUi();
    ui.alert(`エラーが発生しました\n
      エラーを修正して再実行してください\n
      ------エラーメッセージ-----\n` +
      error.message +
      `\n---------------------------`);
  } catch (uiError) {
    // トリガー実行などUIを開けない文脈では何もしない
    Logger.log("ダイアログを表示できませんでした: " + uiError.message);
  }
}

// スプシ編集時にDBのタスクを更新する関数
function onChange(e) {
  const lock = LockService.getScriptLock(); // スクリプト全体で共通のロック

  try {
    // ロックを試みる。最大待機時間を設定 (例: 30秒)
    // この時間は、API呼び出しにかかる最大時間などを考慮して調整
    lock.waitLock(LOCK_WAIT_MS);

    // --- ここからクリティカルセクション ---
    const sheet = e.source.getActiveSheet();
    const sheetName = sheet.getName();

    // タスクのシートでなければキャンセル
    if (!isTaskSheet(sheetName)) {
      Logger.log("タスクシートの編集ではありませんでした");
      return;
    }

    const editRange = sheet.getActiveRange();
    const startRow = editRange.getRow();
    const lastRow = editRange.getLastRow();

    // 編集された範囲のタスク名からマニュアル(スライド版)の列までを一括取得する
    const range = sheet.getRange(startRow, 1, (lastRow - startRow) + 1, COLUMN_COUNT);
    const values = range.getValues();
    const richTextValues = range.getRichTextValues();

    // 編集範囲に見出し行が混ざりうるので0行目から見る(buildChanges側で弾く)
    postChanges(buildChanges(values, richTextValues, 0));
    // --- ここまでクリティカルセクション ---
  } catch (error) {
    alertError(error);
  } finally {
    lock.releaseLock(); // 必ずロックを解放する
  }
}

// 現在のスプシをDBに反映する関数
function updateTasksAndPlaces() {
  const lock = LockService.getScriptLock(); // スクリプト全体で共通のロック

  try {
    // ロックを試みる。最大待機時間を設定 (例: 30秒)
    // この時間は、API呼び出しにかかる最大時間などを考慮して調整
    lock.waitLock(LOCK_WAIT_MS);

    // uiを取得
    const ui = SpreadsheetApp.getUi();

    // --- ここからクリティカルセクション ---
    const sheet = SpreadsheetApp.getActiveSpreadsheet().getActiveSheet();
    const sheetName = sheet.getName(); // シート名取得（dateIDとweatherIDに対応）

    // ダイアログで送信内容を確認
    const confirm = ui.alert(
      `以下のシートのタスクをSeeFTに送信してよろしいですか？\n【 ` + sheetName + ` 】`,
      ui.ButtonSet.OK_CANCEL);
    if (confirm === ui.Button.CANCEL) {
      Logger.log("送信がキャンセルされました");
      return;
    }

    // タスクのシートでなければキャンセル
    if (!isTaskSheet(sheetName)) {
      Logger.log("タスクシートではないため送信をキャンセルしました");
      return ui.alert(`タスクシートではないため送信をキャンセルしました`);
    }

    const range = sheet.getDataRange();
    const values = range.getValues();
    const richTextValues = range.getRichTextValues();

    // 1行目は見出しなので飛ばす
    postChanges(buildChanges(values, richTextValues, 1));
    // --- ここまでクリティカルセクション ---
  } catch (error) {
    alertError(error);
  } finally {
    lock.releaseLock(); // 必ずロックを解放する
  }
}
