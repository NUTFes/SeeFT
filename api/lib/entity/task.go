package entity

import (
	"time"
)

type Task struct {
	ID        int       `json:"id"`
	Task      string    `json:"task"`
	PlaceID   int       `json:"placeID"`
	Url       string    `json:"url"`
	ManualUrl string    `json:"manualUrl"`
	BureauID  int       `json:"bureauID"`
	MaxMember int       `json:"maxMember"`
	Color     string    `json:"color"`
	Remark    string    `json:"remark"`
	YearID    int       `json:"yearID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatetedAt"`
}

type TaskMobile struct {
	ID        int       `json:"id"`
	Task      string    `json:"task"`
	Place     string    `json:"place"`
	Url       string    `json:"url"`
	ManualUrl string    `json:"manualUrl"`
	BureauID  int       `json:"bureauID"`
	MaxMember int       `json:"maxMember"`
	Color     string    `json:"color"`
	Remark    string    `json:"remark"`
	YearID    int       `json:"yearID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatetedAt"`
}

// GASから送られてくるタスクと集合場所の変更内容
type TaskAndPlaceChange struct {
	YearID    int    `json:"yearID"`    // yearID
	TaskName  string `json:"taskName"`  // タスク名
	Bureau    string `json:"bureau"`    // 管轄局
	Place     string `json:"place"`     // 集合場所
	Url       string `json:"url"`       // マニュアルURL
	MaxMember int    `json:"maxMember"` // 最大人数
}

type TaskAndPlaceChangeRequest struct {
	Changes []TaskAndPlaceChange `json:"changes"` // タスクと集合場所の変更内容のリスト
}

// import './entity.dart';

// class Task {
//   int id;
//   String task;
//   Color color;
//   String place;
//   String url;
//   String superviser;
//   String notes;
//   int yearId;
//   DateTime createdAt;
//   DateTime updatedAt;
//   DateTime? deletedAt;

//   Task({
//     this.id = 0,
//     this.task = '',
//     this.color = const Color(0xFFFAFA),
//     this.place = '',
//     this.url = '',
//     this.superviser = '',
//     String? notes,
//     this.yearId = 0,
//     DateTime? createdAt,
//     DateTime? updatedAt,
//     this.deletedAt,
//   })  : notes = notes ?? '',
//         createdAt = createdAt ?? DateTime(0),
//         updatedAt = updatedAt ?? DateTime(0);

//   bool get isDeleted => deletedAt != null;

//   Map<String, dynamic> toJson() => {
//         'id': id,
//         'task': task,
//         'color': color.toString(),
//         'place': place,
//         'url': url,
//         'superviser': superviser,
//         'notes': notes,
//         'createdAt': createdAt.toString(),
//         'updatedAt': updatedAt.toString(),
//         'isDeleted': isDeleted,
//       };
// }
