package entity

type Color struct {
	value	int 	`gorm:"primay_key"`
}

// class Color {
//   final int value;

//   const Color(int value) : value = value & 0xFFFFFF;

//   const Color.fromRGB(int r, int g, int b)
//       : value = (((r & 0xff) << 16) | ((g & 0xff) << 8) | ((b & 0xff) << 0)) & 0xFFFFFF;

//   int get red => (0xFF0000 & value) >> 16;
//   int get green => (0x00FF00 & value) >> 8;
//   int get blue => (0x0000FF & value) >> 0;

//   @override
//   String toString() => "0x${value.toRadixString(16).padLeft(6, '0')}";
// }
