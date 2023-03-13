package entity

import "time"

type Work struct {
	ID        int    `gorm:"primary_key"`
	Name      string `gorm:"size:128"`
	URL       string `gorm:"size:128"`
	Place     string
	President string
	TEL       string
	Color     string
	CreatedAt time.Time
	UpdatedAt time.Time
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
