package entity

import (
	"time"
)

type User struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Mail          string    `json:"mail"`
	GradeID       int       `json:"gradeID"`
	DepartmentID  int       `json:"departmentID"`
	BureauID      int       `json:"bureauID"`
	RoleID        int       `json:"roleID"`
	StudentNumber int       `json:"studentNumber"`
	Tel           string    `json:"tel"`
	Password      string    `json:"password"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type LoginUser struct {
	ID     int    `json:"id"`
	RoleID int    `json:"roleID"`
	Mail   string `json:"mail"`
}

// GASから送られてくる名簿の変更内容
type UserChange struct {
	Name          string `json:"name"`
	Bureau        string `json:"bureau"`
	Grade         string `json:"grade"`
	Department    string `json:"department"`
	StudentNumber int    `json:"studentNumber"`
	Tel           string `json:"tel"`
}

type UserChangeRequest struct {
	Changes []UserChange `json:"changes"`
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
