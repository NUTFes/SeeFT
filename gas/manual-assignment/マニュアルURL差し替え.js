/**
 * シフトスプレッドシートの「マニュアルURL」シートのB列（ドキュメント版のURL）を、配置したPDFのURLに差し替える。
 *
 * 照合は元ドキュメントのファイルIDで行う。割り当て表とマニュアルURLシートでは名前の書き方が違うので、名前では照合しない。
 * B列は タスク一覧のR列（VLOOKUP）→ タスク送信 → tasks.url とつながっているので、差し替えたあとにタスク送信が要る。
 */

const SHIFT_SPREADSHEET_ID = '1b5FhiuT7M6kcAM_BkFRu1UVLj-Ssbt-VoBjGmEAN3-I'; // 45th_シフト_ver0
const MANUAL_URL_SHEET_GID = 1172690727; // 「マニュアルURL」シート
const MANUAL_URL_HEADERS = { key: 'マニュアル名', doc: 'ドキュメントURL' };

// 何も書き込まずに、どの行をどう差し替えるかをログに出す
function dryRunReplaceManualUrls() {
  runReplaceManualUrls_(false);
}

function replaceManualUrls() {
  withLock_(() => runReplaceManualUrls_(true));
}

function runReplaceManualUrls_(apply) {
  console.log(apply ? '実行モード' : '確認モード（何も書き込まない）');

  const pdfBySourceId = {}; // 元ドキュメントのID → 配置したPDF
  const pdfByPlacedId = {}; // 配置したPDFのID → 同じもの。差し替え済みの行を見分ける
  const linkedIds = new Set(); // ショートカットを置いた行。B列は元ファイルのURLのままでよい
  readTargetRows_().rows.forEach((r) => {
    const sourceId = fileIdOf_(r.link);
    const placedId = fileIdOf_(r.placed);
    if (!sourceId || !placedId) return;
    if (sourceId === placedId) {
      linkedIds.add(sourceId);
      return;
    }
    const pdf = { name: r.name, url: r.placed, used: false };
    pdfBySourceId[sourceId] = pdf;
    pdfByPlacedId[placedId] = pdf;
  });

  const sheet = SpreadsheetApp.openById(SHIFT_SPREADSHEET_ID)
    .getSheets()
    .find((s) => s.getSheetId() === MANUAL_URL_SHEET_GID);
  if (!sheet) throw new Error(`gid=${MANUAL_URL_SHEET_GID} のシートが見つからない`);
  const values = sheet.getDataRange().getValues();
  const header = values[0].map((v) => String(v).trim());
  const keyCol = header.indexOf(MANUAL_URL_HEADERS.key);
  const docCol = header.indexOf(MANUAL_URL_HEADERS.doc);
  if (keyCol < 0 || docCol < 0) throw new Error('マニュアルURLシートの見出し（マニュアル名・ドキュメントURL）が見つからない');

  console.log(['行', 'キー', '処理', '割り当て表のマニュアル名', '差し替え後のURL'].join('\t'));
  const counts = { replace: 0, done: 0, linked: 0, none: 0 };
  for (let i = 1; i < values.length; i++) {
    const key = String(values[i][keyCol]).trim();
    const current = String(values[i][docCol]).trim();
    if (!key && !current) continue;
    const currentId = fileIdOf_(current);

    if (pdfBySourceId[currentId]) {
      const pdf = pdfBySourceId[currentId];
      pdf.used = true;
      counts.replace++;
      if (apply) sheet.getRange(i + 1, docCol + 1).setValue(pdf.url);
      console.log([i + 1, key, apply ? '差し替えた' : '差し替える', pdf.name, pdf.url].join('\t'));
    } else if (pdfByPlacedId[currentId]) {
      pdfByPlacedId[currentId].used = true;
      counts.done++;
    } else if (linkedIds.has(currentId)) {
      counts.linked++;
    } else {
      counts.none++;
      console.log([i + 1, key, '対応するPDFが無いのでそのまま', '', ''].join('\t'));
    }
  }

  const unused = Object.keys(pdfBySourceId)
    .map((id) => pdfBySourceId[id])
    .filter((pdf) => !pdf.used)
    .map((pdf) => pdf.name);
  console.log(`差し替え ${counts.replace} ／ 差し替え済み ${counts.done} ／ 元ファイルのまま ${counts.linked} ／ 対応するPDFなし ${counts.none}`);
  if (unused.length) console.log(`マニュアルURLシートに行が無いPDF ${unused.length}本（タスクに紐づいていない）: ${unused.join('、')}`);
  if (apply && counts.replace) console.log('次は、シフトスプレッドシートで点検コマンド → タスク送信を実行する');
}
