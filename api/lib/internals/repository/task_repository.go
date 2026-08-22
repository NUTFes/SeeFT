package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

var taskDebugSQL = os.Getenv("DEBUG_SQL") != "0"

// SELECT * を避けるための明示的なカラムリスト
// manual_urlはNULL許容カラムのため、Scan時のNULL事故防止でCOALESCEしておく（issue #428対応）
const taskColumns = "id, task, place_id, url, COALESCE(manual_url, '') AS manual_url, bureau_id, max_member, color, remark, year_id, created_at, updated_at"

type taskRepository struct {
	client db.Client
	crud   abstract.Crud
}

type TaskRepository interface {
	All(context.Context) (*sql.Rows, error)
	Find(context.Context, string) (*sql.Row, error)
	Shift(context.Context, string) (*sql.Rows, error)
	Create(context.Context, string, string, string, string, string, string, string, string) error
	Update(context.Context, string, string, string, string, string, string, string, string, string) error
	Destroy(context.Context, string) error
	FindNewRecord(context.Context) (*sql.Row, error)
	FindByName(context.Context, string) (*sql.Row, error)
	FindByUserID(context.Context, string) (*sql.Rows, error)
	FindByNames(context.Context, []string) (*sql.Rows, error)
}

func NewTaskRepository(c db.Client, ac abstract.Crud) TaskRepository {
	return &taskRepository{c, ac}
}

// 全件取得
func (b *taskRepository) All(c context.Context) (*sql.Rows, error) {
	query := "SELECT " + taskColumns + " FROM tasks"
	return b.crud.Read(c, query)
}

// 1件取得
func (b *taskRepository) Find(c context.Context, id string) (*sql.Row, error) {
	query := "SELECT " + taskColumns + " FROM tasks WHERE id = $1"
	return b.crud.ReadByID(c, query, id)
}

// 特定のシフト取得
func (b *taskRepository) Shift(c context.Context, name string) (*sql.Rows, error) {
	query := "SELECT " + taskColumns + " FROM tasks WHERE task = $1"
	rows, err := b.client.DB().QueryContext(c, query, name)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot connect SQL")
	}
	if taskDebugSQL {
		fmt.Printf("\x1b[36m%s\n", query)
	}
	return rows, nil
}

// 作成
func (b *taskRepository) Create(c context.Context, name string, placeID string, url string, bureauID string, maxMember string, color string, remark string, yearID string) error {
	query := "INSERT INTO tasks (task, place_id, url, bureau_id, max_member, color, remark, year_id) VALUES ('" + name + "', " + placeID + ", '" + url + "', " + bureauID + ", " + maxMember + ", '" + color + "', '" + remark + "', " + yearID + ")"
	return b.crud.UpdateDB(c, query)
}

// 編集
func (b *taskRepository) Update(c context.Context, id string, name string, placeID string, url string, bureauID string, maxMember string, color string, remark string, yearID string) error {
	query := "UPDATE tasks SET (task, place_id, url, bureau_id, max_member, color, remark, year_id) = ('" + name + "', " + placeID + ", '" + url + "', " + bureauID + ", " + maxMember + ", '" + color + "', '" + remark + "', " + yearID + ") WHERE id = " + id
	return b.crud.UpdateDB(c, query)
}

// 削除
func (b *taskRepository) Destroy(c context.Context, id string) error {
	query := "DELETE FROM tasks WHERE id =" + id
	return b.crud.UpdateDB(c, query)
}

// 最新のtaskを取得する
func (b *taskRepository) FindNewRecord(c context.Context) (*sql.Row, error) {
	query := `
		SELECT
			` + taskColumns + `
		FROM
			tasks
		ORDER BY
			id
		DESC LIMIT 1
	`
	return b.crud.ReadByID(c, query)
}

// タスク名からタスクを取得する
func (b *taskRepository) FindByName(c context.Context, name string) (*sql.Row, error) {
	query := "SELECT " + taskColumns + " FROM tasks WHERE task = $1"
	return b.client.DB().QueryRowContext(c, query, name), nil
}

// 複数のタスク名から一括でタスクを取得する（N+1問題対策）
func (b *taskRepository) FindByNames(c context.Context, names []string) (*sql.Rows, error) {
	if len(names) == 0 {
		// 空の結果を返す
		query := "SELECT " + taskColumns + " FROM tasks WHERE 1=0"
		return b.client.DB().QueryContext(c, query)
	}

	query := "SELECT " + taskColumns + " FROM tasks WHERE task = ANY($1::text[])"
	return b.client.DB().QueryContext(c, query, pq.Array(names))
}

// 指定したuserIDの全てのタスクを取得する
func (b *taskRepository) FindByUserID(c context.Context, userID string) (*sql.Rows, error) {
	query := `
	SELECT DISTINCT t.id, t.task, t.place_id, t.url, COALESCE(t.manual_url, '') AS manual_url, t.bureau_id, t.max_member, t.color, t.remark, t.year_id, t.created_at, t.updated_at
	FROM tasks t
	JOIN shifts s ON t.id = s.task_id
	WHERE s.user_id = $1
	ORDER BY t.task
	`
	if taskDebugSQL {
		fmt.Printf("\x1b[36m%s\n", query)
	}

	return b.client.DB().QueryContext(c, query, userID)
}

// import '../../usecase/repository/task_repository.dart';
// import '../../entity/entity.dart';
// import './external/database.dart';

// class TaskRepositoryImpl implements TaskRepository {
//   Database database;

//   TaskRepositoryImpl(this.database);

//   @override
//   Future<List<Task>> getTasks(ctx) async {
//     String sql = '''
// SELECT * FROM tasks;
// ''';

//     List<Map<String, dynamic>> data = await database.finds(ctx, sql);
//     List<Task> list = [];

//     for (var d in data) {
//       Task task = Task(
//         id: d['id'],
//         task: d['task'],
//         color: Color(int.parse(d['color'], radix: 16)),
//         place: d['place'],
//         url: d['url'],
//         superviser: d['superviser'],
//         notes: d['notes'].toString(),
//         yearId: d['year_id'],
//         createdAt: d['created_at'],
//         updatedAt: d['updated_at'],
//         deletedAt: d['deleted_at'],
//       );

//       if (task.isDeleted) {
//         continue;
//       }

//       list.add(task);
//     }

//     return list;
//   }

//   @override
//   Future<TaskDetail> getTask(ctx, Shift req) async {
//     String sql = '''
// SELECT   tasks.id,
//   tasks.task,
//   tasks.place,
//   tasks.url,
//   tasks.superviser,
//   tasks.notes,
//   tasks.year_id,
//   tasks.created_at,
//   tasks.updated_at,
//   tasks.deleted_at
// FROM tasks
// WHERE tasks.id = ${req.task.id};
// ''';

//     Map<String, dynamic> data = await database.find(ctx, sql);
//     return TaskDetail(
//       id: data['id'],
//       task: data['task'],
//       place: data['place'],
//       url: data['url'],
//       superviser: data['superviser'],
//       yearId: data['year_id'],
//       notes: data['notes'].toString(),
//       users: '',
//       createdAt: data['created_at'],
//       updatedAt: data['updated_at'],
//       deletedAt: data['deleted_at'],
//     );
//   }

//   @override
//   Future<TaskDetail> getTaskByShift(ctx, Shift req) async {
//     String sql = '''
// SELECT
//   tasks.id,
//   tasks.task,
//   tasks.place,
//   tasks.url,
//   tasks.superviser,
//   tasks.notes,
//   tasks.year_id,
//   group_concat(users.name) as users,
//   tasks.created_at,
//   tasks.updated_at,
//   tasks.deleted_at
// FROM shifts
// LEFT JOIN tasks ON shifts.task_id = tasks.id
// LEFT JOIN users ON shifts.user_id = users.id
// WHERE shifts.date_id = ${req.date.id}
//   AND shifts.weather_id = ${req.weather.id}
//   AND shifts.task_id = ${req.task.id}
//   AND shifts.time_id = ${req.time.id}
//   AND shifts.user_id <> ${req.user.id}
// GROUP BY shifts.task_id;
// ''';
//     Map<String, dynamic> data = await database.find(ctx, sql);
//     return TaskDetail(
//       id: data['id'],
//       task: data['task'],
//       place: data['place'],
//       url: data['url'],
//       superviser: data['superviser'],
//       yearId: data['year_id'],
//       notes: data['notes'].toString(),
//       users: data['users'].toString(),
//       createdAt: data['created_at'],
//       updatedAt: data['updated_at'],
//       deletedAt: data['deleted_at'],
//     );
//   }

//   @override
//   Future<int> countUserFromTask(ctx, Shift req) async {
//     String sql = '''
// SELECT COUNT(shifts.user_id) as count
// FROM shifts
// LEFT JOIN tasks ON shifts.task_id = tasks.id
// LEFT JOIN users ON shifts.user_id = users.id
// WHERE shifts.date_id = ${req.date.id}
//   AND shifts.weather_id = ${req.weather.id}
//   AND shifts.task_id = ${req.task.id}
//   AND shifts.time_id = ${req.time.id}
// GROUP BY shifts.task_id;
//     ''';

//     Map<String, dynamic> data = await database.find(ctx, sql);
//     return data['count'];
//   }
// }
