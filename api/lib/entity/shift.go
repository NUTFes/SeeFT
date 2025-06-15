package entity

import (
	"time"
)

// スマホアプリ用
type Shift struct {
	ID           int        `json:"id"`
	Task         TaskMobile `json:"task"`
	User         User       `json:"user"`
	Year         Year       `json:"year"`
	Date         Date       `json:"date"`
	Time         Time       `json:"time"`
	Weather      Weather    `json:"weather"`
	IsAttendance bool       `json:"isAttendance"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

type ShiftUsers struct {
	Task    Task    `json:"task"`
	Users   []User  `json:"users"`
	Year    Year    `json:"year"`
	Date    Date    `json:"date"`
	Time    Time    `json:"time"`
	Weather Weather `json:"weather"`
}

// Webアプリ用
type ShiftAdmin struct {
	ID           int       `json:"id"`
	TaskID       int       `json:"taskID"`
	UserID       int       `json:"userID"`
	YearID       int       `json:"yearID"`
	DateID       int       `json:"dateID"`
	TimeID       int       `json:"timeID"`
	WeatherID    int       `json:"weatherID"`
	IsAttendance bool      `json:"isAttendance"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// タスクの結合用エンティティ
type ShiftMember struct {
	Name   string `json:"name"`
	Grade  string `json:"grade"`
	Bureau string `json:"bureau"`
}

type ShiftMembers struct {
	STime   string        `json:"s_time"`
	ETime   string        `json:"e_time"`
	Members []ShiftMember `json:"members"`
}

type ShiftCard struct {
	TaskName      string         `json:"task_name"`
	StartTime     string         `json:"start_time"`
	EndTime       string         `json:"end_time"`
	Place         string         `json:"place"`
	Url           string         `json:"url"`
	ShiftMembers  []ShiftMembers `json:"shift_members"`
	BeforeMembers ShiftMembers   `json:"before_members"`
	AfterMembers  ShiftMembers   `json:"after_members"`
}

// シフト希望
type ShiftRequest struct {
	Name  string `json:"name"` // ユーザーID
	Shift []struct {
		Date     int `json:"date"` // 日付
		Contents []struct {
			TimeID   int  `json:"timeID"`   // 時間ID
			IsAttend bool `json:"isAttend"` // 出席フラグ
		} `json:"contents"`
	} `json:"shift"`
}

type GASShiftData struct {
	Name  string `json:"name"` // ユーザー名
	Shift []struct {
		Date     int `json:"date"` // 日付
		Contents []struct {
			Row    int  `json:"row"`    // 行番号
			Column int  `json:"column"` // 列番号
			Value  bool `json:"value"`  // セルの値
		} `json:"contents"`
	} `json:"shift"`
}

// GASから送られてくるシフトの変更内容
type ShiftChange struct {
	YearID   int    `json:"yearID"`   // yearID
	TimeID   int    `json:"timeID"`   // timeID
	Date     string `json:"date"`     // 日付
	Weather  string `json:"weather"`  // 天気
	UserName string `json:"userName"` // ユーザー名
	TaskName string `json:"taskName"` // タスク名
}

type ShiftChangeRequest struct {
	Changes []ShiftChange `json:"changes"`
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
