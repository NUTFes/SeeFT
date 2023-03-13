package entity

import "time"

type Permission struct {
	id			int 	`gorm:"primay_key"`
	allowShift	bool
	allowTask	bool
	allowUser	bool
	allowProperty bool
	createdAt	time.Time
	updatedAt	time.Time
	deletedAt	time.Time
}

// class Permission {
//   int id;
//   bool allowShift;
//   bool allowTask;
//   bool allowUser;
//   bool allowProperty;
//   DateTime createdAt;
//   DateTime updatedAt;
//   DateTime? deletedAt;

//   Permission({
//     this.id = 0,
//     this.allowShift = false,
//     this.allowTask = false,
//     this.allowUser = false,
//     this.allowProperty = false,
//     DateTime? createdAt,
//     DateTime? updatedAt,
//     this.deletedAt,
//   })  : createdAt = createdAt ?? DateTime(0),
//         updatedAt = updatedAt ?? DateTime(0);

//   bool get isDeleted => deletedAt != null;

//   Map<String, dynamic> toJson() => {
//         'id': id,
//         'allowShift': allowShift,
//         'allowTask': allowTask,
//         'allowUser': allowUser,
//         'allowProperty': allowProperty,
//         'createdAt': createdAt.toString(),
//         'updatedAt': updatedAt.toString(),
//         'isDeleted': isDeleted,
//       };
// }
