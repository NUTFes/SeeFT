package entity

import (
	"time"
)

type Grade struct {
	ID			int 		`json:"id"`
	Grade 		string		`json:"grade"`
	CreatedAt	time.Time	`json:"createdAt"`
	UpdatedAt	time.Time	`json:"updatedAt"`
}

// class Grade {
//   int id;
//   String grade;
//   DateTime createdAt;
//   DateTime updatedAt;
//   DateTime? deletedAt;

//   Grade({
//     this.id = 0,
//     this.grade = '',
//     DateTime? createdAt,
//     DateTime? updatedAt,
//     this.deletedAt,
//   })  : createdAt = createdAt ?? DateTime(0),
//         updatedAt = updatedAt ?? DateTime(0);

//   bool get isDeleted => deletedAt != null;

//   Map<String, dynamic> toJson() => {
//         'id': id,
//         'grade': grade,
//         'createdAt': createdAt.toString(),
//         'updatedAt': updatedAt.toString(),
//         'isDeleted': isDeleted,
//       };
// }
