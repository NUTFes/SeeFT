package entity

import "time"

type Shift struct {
	ID		   int    		`json:"id"`
	UserID     uint    		`json:"userID"`
	Date       string 		`json:"date"`
	Time       string 		`json:"time"`
	WorkID     uint    		`json:"workID"`
	Weather    string 		`json:"weather"`
	Attendance bool   		`json:"attendance"`
	CreatedAt  time.Time	`json:"createdAt"`
	UpdatedAt  time.Time	`json:"updatedAt"`
}

// class Shift {
//   int id;
//   User user;
//   Task task;
//   Year year;
//   Date date;
//   Time time;
//   Weather weather;
//   bool isAttendance;
//   int createdUserId;
//   int updatedUserId;
//   DateTime createdAt;
//   DateTime updatedAt;
//   DateTime? deletedAt;

//   Shift({
//     this.id = 0,
//     User? user,
//     task,
//     year,
//     date,
//     time,
//     weather,
//     this.isAttendance = false,
//     this.createdUserId = 0,
//     this.updatedUserId = 0,
//     DateTime? createdAt,
//     DateTime? updatedAt,
//     this.deletedAt,
//   })  : user = user ?? User(),
//         task = task ?? Task(),
//         year = year ?? Year(),
//         date = date ?? Date(),
//         time = time ?? Time(),
//         weather = weather ?? Weather(),
//         createdAt = createdAt ?? DateTime(0),
//         updatedAt = updatedAt ?? DateTime(0);

//   bool get isDeleted => deletedAt != null;

//   Map<String, dynamic> toJson() => {
//         'id': id,
//         'user': user,
//         'task': task,
//         'year': year,
//         'date': date,
//         'time': time,
//         'weather': weather,
//         'isAttendance': isAttendance,
//         'createdUserId': createdUserId,
//         'updatedUserId': updatedUserId,
//         'createdAt': createdAt.toString(),
//         'updatedAt': updatedAt.toString(),
//         'isDeleted': isDeleted,
//       };
// }
