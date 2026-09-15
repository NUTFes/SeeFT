/**
 * SeeFT で見るマニュアルを、共有ドライブの置き場に配置する。
 *
 * 対象は「マニュアル」シートで、進捗が「完成(ドキュメント)」「完成(AI)」の行。
 * 置き場は 06_SeeFT閲覧用マニュアル置き場/<担当局>/ で、置いたもののURLを「PDF配置」列に書く。
 *
 *   Googleドキュメント・スライド  PDFに書き出して置く
 *   Word                         一時的にGoogleドキュメントへ変換してからPDFに書き出して置く（元のWordは変更しない）
 *   スプレッドシート・PDF        コピーせずショートカットを置き、元ファイルのURLを書く（コピーは内容が止まるため）
 *
 * PDF配置列のURLは、マニュアルURLシートのB列 → tasks.url → アプリのボタンへと流れる。
 * URLが変わると全部に波及するので、出し直しは refreshPdfs で同じファイルの中身だけを差し替える。
 */

const ROOT_FOLDER_ID = '1J8-yJg2RozOwbSOLcijptq7C-st64Tba'; // 06_SeeFT閲覧用マニュアル置き場
const MANUAL_SHEET_GID = 129469480; // 「マニュアル」シート。名前が変わっても追えるよう gid で探す
const TARGET_STATUSES = ['完成(ドキュメント)', '完成(AI)'];
const HEADERS = { name: 'マニュアル名', status: '進捗', bureau: '担当局', dept: '担当部門', placed: 'PDF配置' };
const TIME_LIMIT_MS = 5 * 60 * 1000; // GASの6分制限の手前で止める

const EXPORT_AS_PDF = [MimeType.GOOGLE_DOCS, MimeType.GOOGLE_SLIDES];
// Word はそのままではPDFに書き出せないので、一時的にGoogleドキュメントへ変換してから書き出す
const CONVERT_TO_DOCS = [MimeType.MICROSOFT_WORD, MimeType.MICROSOFT_WORD_LEGACY];
const PLACE_SHORTCUT = [MimeType.GOOGLE_SHEETS, MimeType.PDF];
const KIND_LABEL = {
  [MimeType.GOOGLE_DOCS]: 'ドキュメント',
  [MimeType.GOOGLE_SLIDES]: 'スライド',
  [MimeType.MICROSOFT_WORD]: 'Word',
  [MimeType.MICROSOFT_WORD_LEGACY]: 'Word',
  [MimeType.GOOGLE_SHEETS]: 'スプレッドシート',
  [MimeType.PDF]: 'PDF',
};
const ACTION_LABEL = {
  pdf: 'PDFを作る',
  shortcut: 'ショートカットを作る',
  placed: '配置済み',
  error: '要確認',
};

function onOpen() {
  SpreadsheetApp.getUi()
    .createMenu('SeeFT')
    .addItem('PDF配置：確認だけする', 'dryRunPlacePdfs')
    .addItem('PDF配置：実行する', 'placePdfs')
    .addSeparator()
    .addItem('PDFを出し直す（URLは変わらない）', 'refreshPdfs')
    .addSeparator()
    .addItem('マニュアルURLシートB列：確認だけする', 'dryRunReplaceManualUrls')
    .addItem('マニュアルURLシートB列：差し替える', 'replaceManualUrls')
    .addToUi();
}

// 何も書き込まずに、各行をどう扱うかをログに出す
function dryRunPlacePdfs() {
  runPlacement_(false);
}

function placePdfs() {
  withLock_(() => runPlacement_(true));
}

function runPlacement_(apply) {
  const started = Date.now();
  const root = DriveApp.getFolderById(ROOT_FOLDER_ID);
  console.log(`置き場: ${root.getName()} ／ 既存の局フォルダ: ${folderNames_(root).join('、') || 'なし'}`);
  console.log(apply ? '実行モード' : '確認モード（何も書き込まない）');

  const { sheet, cols, rows } = readTargetRows_();
  const plans = rows.map((r) => planRow_(r, root));
  renameCollisions_(plans);
  flagExistingNames_(plans);

  console.log(['行', '担当局', 'マニュアル名', '元の形式', '処理', '置く名前', '備考'].join('\t'));
  const counts = {};
  for (const p of plans) {
    if (apply && (p.action === 'pdf' || p.action === 'shortcut')) {
      if (Date.now() - started > TIME_LIMIT_MS) {
        console.log('時間切れが近いので中断した。もう一度実行すると、配置済みの行を飛ばして続きから進む');
        break;
      }
      try {
        const url = place_(p, root);
        // 1件ごとに書く。途中で止まっても、置いた分は次回「配置済み」として飛ばされる
        sheet.getRange(p.row, cols.placed + 1).setValue(url);
        p.note = url;
      } catch (e) {
        p.action = 'error';
        p.note = `実行時エラー: ${e.message}`;
      }
    }
    counts[p.action] = (counts[p.action] || 0) + 1;
    console.log([p.row, p.bureau, p.name, p.kind, ACTION_LABEL[p.action], p.fileName, p.note].join('\t'));
  }
  console.log(`対象 ${plans.length} 行 ／ ` + Object.keys(counts).map((k) => `${ACTION_LABEL[k]} ${counts[k]}`).join(' ／ '));
}

function readTargetRows_() {
  const ss = SpreadsheetApp.getActiveSpreadsheet();
  const sheet = ss.getSheets().find((s) => s.getSheetId() === MANUAL_SHEET_GID);
  if (!sheet) throw new Error(`gid=${MANUAL_SHEET_GID} のシートが見つからない`);

  const values = sheet.getDataRange().getValues();
  const trimmed = (row) => row.map((v) => String(v).trim());
  const headerIndex = values.slice(0, 10).findIndex((row) => {
    const t = trimmed(row);
    return t.includes(HEADERS.name) && t.includes(HEADERS.status);
  });
  if (headerIndex < 0) throw new Error('見出し行（マニュアル名・進捗）が見つからない');

  const header = trimmed(values[headerIndex]);
  const cols = {};
  Object.keys(HEADERS).forEach((key) => {
    cols[key] = header.indexOf(HEADERS[key]);
    if (cols[key] < 0) throw new Error(`見出し「${HEADERS[key]}」が見つからない`);
  });

  const firstDataIndex = headerIndex + 1;
  if (values.length <= firstDataIndex) return { sheet, cols, rows: [] };
  const richTexts = sheet
    .getRange(firstDataIndex + 1, cols.name + 1, values.length - firstDataIndex, 1)
    .getRichTextValues();

  const rows = [];
  for (let i = firstDataIndex; i < values.length; i++) {
    const v = values[i];
    if (!TARGET_STATUSES.includes(String(v[cols.status]).trim())) continue;
    rows.push({
      row: i + 1,
      name: String(v[cols.name]).replace(/\s+/g, ' ').trim(),
      bureau: String(v[cols.bureau]).trim(),
      dept: String(v[cols.dept]).trim(),
      placed: String(v[cols.placed]).trim(),
      link: firstLinkOf_(richTexts[i - firstDataIndex][0]),
    });
  }
  fillLinksViaSheetsApi_(sheet, cols.name, rows.filter((r) => !r.link));
  return { sheet, cols, rows };
}

function firstLinkOf_(richText) {
  if (!richText) return '';
  if (richText.getLinkUrl()) return richText.getLinkUrl();
  const run = richText.getRuns().find((r) => r.getLinkUrl());
  return run ? run.getLinkUrl() : '';
}

// スマートチップのリンクは RichTextValue から取れないことがあるので、Sheets API で読み直す
function fillLinksViaSheetsApi_(sheet, nameCol, rows) {
  if (rows.length === 0) return;
  const sheetName = sheet.getName().replace(/'/g, "''");
  const col = columnLetter_(nameCol + 1);
  const ranges = rows.map((r) => `'${sheetName}'!${col}${r.row}`);
  const spreadsheetId = sheet.getParent().getId();

  let res;
  try {
    res = Sheets.Spreadsheets.get(spreadsheetId, { ranges, fields: 'sheets.data.rowData.values(hyperlink,chipRuns)' });
  } catch (e) {
    console.log(`チップ情報を読めなかったので、ハイパーリンクだけで読み直す: ${e.message}`);
    res = Sheets.Spreadsheets.get(spreadsheetId, { ranges, fields: 'sheets.data.rowData.values(hyperlink)' });
  }

  const grids = (res.sheets && res.sheets[0] && res.sheets[0].data) || [];
  rows.forEach((r, i) => {
    const rowData = (grids[i] && grids[i].rowData) || [];
    const cell = (rowData[0] && rowData[0].values && rowData[0].values[0]) || {};
    const chip = (cell.chipRuns || []).find((c) => c.chip && c.chip.richLinkProperties && c.chip.richLinkProperties.uri);
    r.link = cell.hyperlink || (chip ? chip.chip.richLinkProperties.uri : '');
  });
}

function planRow_(r, root) {
  const p = Object.assign({}, r, { action: 'error', note: '', kind: '', fileName: '' });
  if (r.placed) {
    p.action = 'placed';
    return p;
  }
  if (!r.bureau) {
    p.note = '担当局が空';
    return p;
  }
  const sourceId = fileIdOf_(r.link);
  if (!sourceId) {
    p.note = r.link ? `リンクからファイルIDが取れない: ${r.link}` : 'マニュアル名にリンクが無い';
    return p;
  }

  let source;
  try {
    source = DriveApp.getFileById(sourceId);
  } catch (e) {
    p.note = `元ファイルを開けない: ${e.message}`;
    return p;
  }
  if (source.isTrashed()) {
    p.note = '元ファイルがゴミ箱にある';
    return p;
  }

  p.sourceId = sourceId;
  p.mime = source.getMimeType();
  p.kind = KIND_LABEL[p.mime] || p.mime;
  p.folder = findFolder_(root, r.bureau); // 無ければ実行時に作る

  if (EXPORT_AS_PDF.includes(p.mime) || CONVERT_TO_DOCS.includes(p.mime)) {
    p.action = 'pdf';
    p.fileName = `${r.name}.pdf`;
    if (CONVERT_TO_DOCS.includes(p.mime)) p.note = 'Wordを変換してPDFにする（元のWordは変更しない）';
  } else if (PLACE_SHORTCUT.includes(p.mime)) {
    p.action = 'shortcut';
    p.fileName = r.name;
  } else {
    p.note = 'この形式は自動では扱わない';
  }
  return p;
}

// 同じ局フォルダに同じ名前が並ぶと区別できないので、担当部門を付け、それでも重なれば行番号を付ける
function renameCollisions_(plans) {
  const targets = plans.filter((p) => p.action === 'pdf' || p.action === 'shortcut');
  [(p) => p.dept || `${p.row}行目`, (p) => `${p.row}行目`].forEach((suffixOf) => {
    const groups = {};
    targets.forEach((p) => {
      const key = `${p.bureau}/${p.fileName}`;
      (groups[key] = groups[key] || []).push(p);
    });
    Object.keys(groups)
      .filter((key) => groups[key].length > 1)
      .forEach((key) =>
        groups[key].forEach((p) => {
          p.fileName = withSuffix_(p.fileName, suffixOf(p));
          p.note = '同じ名前があるので名前を変えた';
        })
      );
  });
}

function withSuffix_(fileName, suffix) {
  return fileName.endsWith('.pdf') ? `${fileName.slice(0, -4)}（${suffix}）.pdf` : `${fileName}（${suffix}）`;
}

// PDF配置列が空なのに同じ名前が既にあるのは、手で置いたか前回の実行が途中で止まったとき。上書きせず人が見る
function flagExistingNames_(plans) {
  plans
    .filter((p) => (p.action === 'pdf' || p.action === 'shortcut') && p.folder)
    .forEach((p) => {
      if (p.folder.getFilesByName(p.fileName).hasNext()) {
        p.action = 'error';
        p.note = '同じ名前のファイルが局フォルダに既にある';
      }
    });
}

function place_(p, root) {
  const folder = p.folder || findFolder_(root, p.bureau) || root.createFolder(p.bureau);
  if (p.action === 'pdf') {
    const blob = pdfBlobOf_(p.sourceId, p.mime).setName(p.fileName);
    return viewUrl_(folder.createFile(blob).getId());
  }
  folder.createShortcut(p.sourceId).setName(p.fileName);
  return p.mime === MimeType.GOOGLE_SHEETS ? `https://docs.google.com/spreadsheets/d/${p.sourceId}/edit` : viewUrl_(p.sourceId);
}

function pdfBlobOf_(sourceId, mime) {
  if (!CONVERT_TO_DOCS.includes(mime)) return DriveApp.getFileById(sourceId).getAs(MimeType.PDF);

  // 変換用の一時ドキュメントを置き場の直下に作り、PDFを取り出したらゴミ箱に移す
  const tempName = `一時変換_${sourceId}`;
  const temp = Drive.Files.copy(
    { name: tempName, mimeType: MimeType.GOOGLE_DOCS, parents: [ROOT_FOLDER_ID] },
    sourceId,
    { supportsAllDrives: true }
  );
  try {
    const bytes = DriveApp.getFileById(temp.id).getAs(MimeType.PDF).getBytes();
    return Utilities.newBlob(bytes, MimeType.PDF);
  } finally {
    try {
      DriveApp.getFileById(temp.id).setTrashed(true);
    } catch (e) {
      console.log(`一時ファイルをゴミ箱に移せなかったので、手で消す: ${tempName}（${e.message}）`);
    }
  }
}

// 元ドキュメントを直したあとに使う。PDF配置列が指すPDFの中身だけを差し替えるので、URLは変わらない
function refreshPdfs() {
  withLock_(() => {
    const started = Date.now();
    const root = DriveApp.getFolderById(ROOT_FOLDER_ID);
    const { rows } = readTargetRows_();
    let updated = 0;
    for (const r of rows) {
      const pdfId = fileIdOf_(r.placed);
      const sourceId = fileIdOf_(r.link);
      // ショートカットの行は元ファイルのURLを書いているので、差し替えるものが無い
      if (!pdfId || !sourceId || pdfId === sourceId) continue;
      if (Date.now() - started > TIME_LIMIT_MS) {
        console.log('時間切れが近いので中断した。もう一度実行すると、まだ古いPDFだけを出し直す');
        break;
      }
      try {
        const pdf = DriveApp.getFileById(pdfId);
        const source = DriveApp.getFileById(sourceId);
        if (pdf.getMimeType() !== MimeType.PDF || !isInPlacementFolder_(pdf, root)) {
          console.log(`${r.row}\t${r.name}\t対象外（置き場のPDFではない）`);
          continue;
        }
        if (source.getLastUpdated() <= pdf.getLastUpdated()) continue; // 元が更新されていない
        Drive.Files.update({}, pdfId, pdfBlobOf_(sourceId, source.getMimeType()), { supportsAllDrives: true });
        updated++;
        console.log(`${r.row}\t${r.name}\t出し直した`);
      } catch (e) {
        console.log(`${r.row}\t${r.name}\t要確認: ${e.message}`);
      }
    }
    console.log(`出し直し ${updated} 件`);
  });
}

function isInPlacementFolder_(file, root) {
  const parents = file.getParents();
  while (parents.hasNext()) {
    const grandParents = parents.next().getParents();
    while (grandParents.hasNext()) {
      if (grandParents.next().getId() === root.getId()) return true;
    }
  }
  return false;
}

// 二重に実行すると同じPDFが2つできるので、動いている間は後から来た実行を止める
function withLock_(fn) {
  const lock = LockService.getDocumentLock();
  if (!lock.tryLock(1000)) {
    console.log('ほかの実行が動いているので中止した');
    SpreadsheetApp.getActiveSpreadsheet().toast('ほかの実行が動いているので中止しました');
    return;
  }
  try {
    fn();
  } finally {
    lock.releaseLock();
  }
}

function findFolder_(parent, name) {
  const it = parent.getFoldersByName(name);
  return it.hasNext() ? it.next() : null;
}

function folderNames_(parent) {
  const names = [];
  const it = parent.getFolders();
  while (it.hasNext()) names.push(it.next().getName());
  return names;
}

function viewUrl_(fileId) {
  return `https://drive.google.com/file/d/${fileId}/view`;
}

function fileIdOf_(url) {
  const m = String(url || '').match(/\/d\/([A-Za-z0-9_-]{20,})|[?&]id=([A-Za-z0-9_-]{20,})/);
  return m ? m[1] || m[2] : '';
}

function columnLetter_(n) {
  let s = '';
  for (let x = n; x > 0; x = Math.floor((x - 1) / 26)) {
    s = String.fromCharCode(65 + ((x - 1) % 26)) + s;
  }
  return s;
}
