package entity

import "time"

type Weather struct {
	id		int 	`gorm:"primay_key"`
	weather	string 	`gorm:"size:128"`
	createdAt	time.Time
	updatedAt	time.Time
	deletedAt	time.Time
}

// class Weather {
//   int id;
//   String weather;
//   DateTime createdAt;
//   DateTime updatedAt;
//   DateTime? deletedAt;

//   Weather({
//     this.id = 0,
//     this.weather = '',
//     DateTime? createdAt,
//     DateTime? updatedAt,
//     this.deletedAt,
//   })  : createdAt = createdAt ?? DateTime(0),
//         updatedAt = updatedAt ?? DateTime(0);

//   bool get isDeleted => deletedAt != null;

//   Map<String, dynamic> toJson() => {
//         'id': id,
//         'weather': weather,
//         'createdAt': createdAt.toString(),
//         'updatedAt': updatedAt.toString(),
//         'isDeleted': isDeleted,
//       };
// }
