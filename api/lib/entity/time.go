package	entity

import "time"

type Time struct {
	id 		int 	`gorm:"primay_key"`
	time	string	`gorm:"size:128"`
	createdAt	time.Time
	updatedAt	time.Time
	deletedAt	time.Time
}

// class Time {
//   int id;
//   String time;
//   DateTime createdAt;
//   DateTime updatedAt;
//   DateTime? deletedAt;

//   Time({
//     this.id = 0,
//     this.time = '',
//     DateTime? createdAt,
//     DateTime? updatedAt,
//     this.deletedAt,
//   })  : createdAt = createdAt ?? DateTime(0),
//         updatedAt = updatedAt ?? DateTime(0);

//   bool get isDeleted => deletedAt != null;

//   Map<String, dynamic> toJson() => {
//         'id': id,
//         'time': time,
//         'createdAt': createdAt.toString(),
//         'updatedAt': updatedAt.toString(),
//         'isDeleted': isDeleted,
//       };
// }
