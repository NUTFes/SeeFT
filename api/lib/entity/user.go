package entity

import "time"

type User struct {
	ID        int    `gorm:"primary_key"`
	Name      string `gorm:"size:128"`
	Mail      string `gorm:"size:128"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// class User {
//   final int id;
//   final String name;
//   final String mail;
//   final int bureauId;
//   final int gradeId;
//   final DateTime createdAt;
//   final DateTime updatedAt;
//   final DateTime? deletedAt;

//   User({
//     this.id = 0,
//     this.name = '',
//     this.mail = '',
//     this.bureauId = 0,
//     this.gradeId = 0,
//     DateTime? createdAt,
//     DateTime? updatedAt,
//     this.deletedAt,
//   })  : createdAt = createdAt ?? DateTime(0),
//         updatedAt = updatedAt ?? DateTime(0);

//   bool get isDeleted => deletedAt != null;

//   Map<String, dynamic> toJson() => {
//         'id': id,
//         'name': name,
//         'mail': mail,
//         'bureauId': bureauId,
//         'gradeId': gradeId,
//         'createdAt': createdAt.toString(),
//         'updatedAt': updatedAt.toString(),
//         'isDeleted': isDeleted
//       };
// }
