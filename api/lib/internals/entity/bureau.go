package entity

import (
	"time"
)

type Bureau struct {
	ID			int 		`json:"id"`
	Name		string		`json:"name"`
	CreatedAt	time.Time	`json:"createdAt"`
	UpdatedAt	time.Time	`json:"updatedAt"`
}

// import './color.dart';

// class Bureau {
//   int id;
//   String bureau;
//   Color color;
//   DateTime createdAt;
//   DateTime updatedAt;
//   DateTime? deletedAt;

//   Bureau({
//     this.id = 0,
//     this.bureau = '',
//     this.color = const Color(0xFFFAFA),
//     DateTime? createdAt,
//     DateTime? updatedAt,
//     this.deletedAt,
//   })  : createdAt = createdAt ?? DateTime(0),
//         updatedAt = updatedAt ?? DateTime(0);

//   bool get isDeleted => deletedAt != null;

//   Map<String, dynamic> toJson() => {
//         'id': id,
//         'bureau': bureau,
//         'color': color.toString(),
//         'createdAt': createdAt.toString(),
//         'updatedAt': updatedAt.toString(),
//         'isDeleted': isDeleted,
//       };
// }
