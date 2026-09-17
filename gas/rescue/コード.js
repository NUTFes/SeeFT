function myFunction() {
  
}

function doPost(e) {
  try {
    var data = JSON.parse(e.postData.contents);
    var sheetName = getSheetName(data.rescue_type);
    if (!sheetName) {
      return ContentService.createTextOutput('Invalid rescue_type').setMimeType(ContentService.MimeType.TEXT);
    }

    var ss = SpreadsheetApp.getActiveSpreadsheet();
    var sheet = ss.getSheetByName(sheetName);
    if (!sheet) {
      return ContentService.createTextOutput('Sheet not found: ' + sheetName).setMimeType(ContentService.MimeType.TEXT);
    }

    var nextId = data.rescue_id;

    var row = [];
    switch (data.rescue_type) {
      case "trouble":
        row = [
          nextId,
          data.sender_name || "",
          data.student_number || "",
          data.phone_number || "",
          data.grade || "",
          data.bureau || "",
          data.answered_at || "",
          data.place || "",
          data.task_name || "",
          data.detail || ""
        ];
        break;
      case "question":
        row = [
          nextId,
          data.sender_name || "",
          data.student_number || "",
          data.phone_number || "",
          data.grade || "",
          data.bureau || "",
          data.answered_at || "",
          data.question || ""
        ];
        break;
      case "shorthanded":
        row = [
          nextId,
          data.sender_name || "",
          data.student_number || "",
          data.phone_number || "",
          data.grade || "",
          data.bureau || "",
          data.answered_at || "",
          data.place || "",
          data.task_name || "",
          data.missing_number || ""
        ];
        break;
      default:
        return ContentService.createTextOutput('Invalid rescue_type').setMimeType(ContentService.MimeType.TEXT);
    }

    // 電話番号(D列)は先頭の0が落ちないよう、テキストとして入力する。
    // appendRow は数字だけの文字列を数値として解釈するので、09012345678 が 9012345678 になる。
    // 先頭のアポストロフィは「以降をテキストとして扱う」入力で、セルの値には残らない。
    // 表示形式を直接テキストにする手は、このシートがテーブルで列に型があるため使えない
    // （UI・GASとも「型付きの列でセルの数値形式を設定することはできません」で弾かれる）
    var PHONE_INDEX = 3; // 0始まり。対応番号 / 送信者名 / 学籍番号 / 電話番号 の4列目は3タイプ共通
    if (row.length > PHONE_INDEX && row[PHONE_INDEX] !== "") {
      row[PHONE_INDEX] = "'" + String(row[PHONE_INDEX]);
    }

    sheet.appendRow(row);
    return ContentService.createTextOutput('Success').setMimeType(ContentService.MimeType.TEXT);
  } catch (err) {
    return ContentService.createTextOutput('Error: ' + err).setMimeType(ContentService.MimeType.TEXT);
  }
}

// rescue_typeからシート名を返す
function getSheetName(rescueType) {
  switch (rescueType) {
    case "trouble":
      return "トラブル";
    case "question":
      return "質問";
    case "shorthanded":
      return "人が来ない";
    default:
      return null;
  }
}

// 45th版を作った直後に1回だけ実行する。コピー元(44th)の実データ行を消し、対応番号を連番のまま残す
function clearInheritedRescueData() {
var ss = SpreadsheetApp.getActiveSpreadsheet();
var targets = ["事件事故", "トラブル", "人が来ない", "質問"];
var report = [];
targets.forEach(function (name) {
var sh = ss.getSheetByName(name);
if (!sh) { report.push(name + ": シートなし"); return; }
var lastRow = sh.getLastRow();
var lastCol = sh.getLastColumn();
if (lastRow < 3) { report.push(name + ": 消すデータなし"); return; }
// 3行目以降のB列以降(対応番号のA列は残す)をクリア
sh.getRange(3, 2, lastRow - 2, lastCol - 1).clearContent();
report.push(name + ": " + (lastRow - 2) + "行分をクリア");
});
SpreadsheetApp.flush();
Logger.log(report.join(" / "));
}

// 全タスクシートのA3をver0のタスク一覧参照に切り替える。UI編集だと巻き戻ったためGASから確実に書き込む
function setAllTaskFormula() {
  var sh = SpreadsheetApp.getActiveSpreadsheet().getSheetByName("全タスク");
  var u = "https://docs.google.com/spreadsheets/d/1b5FhiuT7M6kcAM_BkFRu1UVLj-Ssbt-VoBjGmEAN3-I/edit?gid=18729372";
  var f = '=LET(タスク,IMPORTRANGE("' + u + '","タスク一覧!$A$4:Q"),タスク名,IMPORTRANGE("' + u + '","タスク一覧!$A$4:A"),FILTER(タスク,タスク名<>""))';
  sh.getRange("A3").setFormula(f);
  SpreadsheetApp.flush();
  Utilities.sleep(4000);
  Logger.log("formula=" + sh.getRange("A3").getFormula());
  Logger.log("A3=" + sh.getRange("A3").getDisplayValue());
  Logger.log("A4=" + sh.getRange("A4").getDisplayValue());
}