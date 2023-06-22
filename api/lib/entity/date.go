package entity

import (
	"time"
)

type Date struct {
	ID			int 		`json:"id"`
	Date		string		`json:"date"`
	CreatedAt	time.Time	`json:"createdAt"`
	UpdatedAt	time.Time	`json:"updatedAt"`
}

// class Date {
//   int id;
//   String date;
//   DateTime createdAt;
//   DateTime updatedAt;
//   DateTime? deletedAt;

//   Date({
//     this.id = 0,
//     this.date = '',
//     DateTime? createdAt,
//     DateTime? updatedAt,
//     this.deletedAt,
//   })  : createdAt = createdAt ?? DateTime(0),
//         updatedAt = updatedAt ?? DateTime(0);

//   bool get isDeleted => deletedAt != null;

//   Map<String, dynamic> toJson() => {
//         'id': id,
//         'date': date,
//         'createdAt': createdAt.toString(),
//         'updatedAt': updatedAt.toString(),
//         'isDeleted': isDeleted,
//       };
// }
