package entity

import (
	"time"
)

type Permission struct {
	ID				int 		`gorm:"primay_key"`
	AllowShift		bool		`json:"allowShift"`
	AllowTask		bool		`json:"allowTask"`
	AllowUser		bool		`json:"allowUser"`
	AllowProperty 	bool		`json:"allowPreperty"`
	CreatedAt		time.Time	`json:"createdAt"`
	UpdatedAt		time.Time	`json:"updatedAt"`
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
