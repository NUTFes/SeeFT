package	entity

import "time"

type Time struct {
	ID 			int 		`json:"id"`
	Time		string		`json:"time"`
	CreatedAt	time.Time	`json:"createdAt"`
	UpdatedAt	time.Time	`json:"updatedAt"`
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
