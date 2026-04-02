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