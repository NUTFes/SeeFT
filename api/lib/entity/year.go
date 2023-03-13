package entity

import "time"

type Year struct {
	id		int 	`gorm:"primay_key"`
	year	int 	`gorm:"size:128"`
	createdAt	time.Time
	updatedAt	time.Time
	deletedAt	time.Time
}

// class Year {
//   int id;
//   int year;
//   DateTime createdAt;
//   DateTime updatedAt;
//   DateTime? deletedAt;

//   Year({
//     this.id = 0,
//     this.year = 0,
//     DateTime? createdAt,
//     DateTime? updatedAt,
//     this.deletedAt,
//   })  : createdAt = createdAt ?? DateTime(0),
//         updatedAt = updatedAt ?? DateTime(0);

//   bool get isDeleted => deletedAt != null;

//   Map<String, dynamic> toJson() => {
//         'id': id,
//         'year': year,
//         'createdAt': createdAt.toString(),
//         'updatedAt': updatedAt.toString(),
//         'isDeleted': isDeleted,
//       };
// }
