package entity

import "time"

type TaskDetail struct {
	ID			int			`json:"id"`
	Task		string		`json:"task"`
	Place		string		`json:"place"`
	Url 		string		`json:"url"`
	Superviser	string		`json:"superviser"`
	Notes		string		`json:"notes"`
	YearId		uint		`json:"yaerID"`
	Users		string		`json:"users"`
	CreatedAt	time.Time	`json:"createdAt"`
	UpdatedAt	time.Time	`json:"updatedAt"`
}

// class TaskDetail {
//   int id;
//   String task;
//   String place;
//   String url;
//   String superviser;
//   String notes;
//   int yearId;
//   String users;
//   DateTime createdAt;
//   DateTime updatedAt;
//   DateTime? deletedAt;

//   TaskDetail({
//     this.id = 0,
//     this.task = '',
//     this.place = '',
//     this.url = '',
//     this.superviser = '',
//     String? notes,
//     this.yearId = 0,
//     String? users,
//     DateTime? createdAt,
//     DateTime? updatedAt,
//     this.deletedAt,
//   })  : notes = notes ?? '',
//         users = users ?? '',
//         createdAt = createdAt ?? DateTime(0),
//         updatedAt = updatedAt ?? DateTime(0);

//   bool get isDeleted => deletedAt != null;

//   Map<String, dynamic> toJson() => {
//         'id': id,
//         'task': task,
//         'place': place,
//         'url': url,
//         'superviser': superviser,
//         'notes': notes,
//         'users': users,
//         'createdAt': createdAt.toString(),
//         'updatedAt': updatedAt.toString(),
//         'isDeleted': isDeleted,
//       };
// }
