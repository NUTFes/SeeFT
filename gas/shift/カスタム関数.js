// function JOIN_FILTERED_NAMES(targetValue) {
//   const sheet = SpreadsheetApp.getActiveSpreadsheet().getSheetByName("準々備日");
//   const names = sheet.getRange("A3:A" + sheet.getLastRow()).getValues().flat();
//   const targets = sheet.getRange("BO3:BO" + sheet.getLastRow()).getValues().flat();

//   const filteredNames = names.filter((name, i) => targets[i] === targetValue);

//   return filteredNames.join(", ");
// }
