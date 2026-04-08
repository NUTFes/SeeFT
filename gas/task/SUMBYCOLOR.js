function SUMBYCOLOR(colorCode, range) {
  if (!colorCode || !range) {
    return '引数が足りません';
  }

  const bgColor = colorCode.toString().toLowerCase();
  const values = range.getValues();
  const backgrounds = range.getBackgrounds();

  let total = 0;
  for (let i = 0; i < values.length; i++) {
    for (let j = 0; j < values[i].length; j++) {
      if (backgrounds[i][j].toLowerCase() === bgColor) {
        const val = values[i][j];
        if (!isNaN(val)) {
          total += Number(val);
        }
      }
    }
  }
  return total;
}