import '../../usecase/repository/task_repository.dart';
import '../../entity/entity.dart';
import './external/database.dart';

class TaskRepositoryImpl implements TaskRepository {
  Database database;

  TaskRepositoryImpl(this.database);

  @override
  Future<List<Task>> getTasks(ctx) async {
    String sql = '''
SELECT * FROM tasks;
''';

    List<Map<String, dynamic>> data = await database.finds(ctx, sql);
    List<Task> list = [];

    for (var d in data) {
      Task task = Task(
        id: d['id'],
        task: d['task'],
        color: Color(int.parse(d['color'], radix: 16)),
        place: d['place'],
        url: d['url'],
        superviser: d['superviser'],
        notes: d['notes'].toString(),
        yearId: d['year_id'],
        createdAt: d['created_at'],
        updatedAt: d['updated_at'],
        deletedAt: d['deleted_at'],
      );

      if (task.isDeleted) {
        continue;
      }

      list.add(task);
    }

    return list;
  }

  @override
  Future<TaskDetail> getTask(ctx, Shift req) async {
    String sql = '''
SELECT   tasks.id,
  tasks.task,
  tasks.place,
  tasks.url,
  tasks.superviser,
  tasks.notes,
  tasks.year_id,
  tasks.created_at,
  tasks.updated_at,
  tasks.deleted_at
FROM tasks
WHERE tasks.id = ${req.task.id};
''';

    Map<String, dynamic> data = await database.find(ctx, sql);
    return TaskDetail(
      id: data['id'],
      task: data['task'],
      place: data['place'],
      url: data['url'],
      superviser: data['superviser'],
      yearId: data['year_id'],
      notes: data['notes'].toString(),
      users: '',
      createdAt: data['created_at'],
      updatedAt: data['updated_at'],
      deletedAt: data['deleted_at'],
    );
  }

  @override
  Future<TaskDetail> getTaskByShift(ctx, Shift req) async {
    String sql = '''
SELECT
  tasks.id,
  tasks.task,
  tasks.place,
  tasks.url,
  tasks.superviser,
  tasks.notes,
  tasks.year_id,
  group_concat(users.name) as users,
  tasks.created_at,
  tasks.updated_at,
  tasks.deleted_at
FROM shifts
LEFT JOIN tasks ON shifts.task_id = tasks.id
LEFT JOIN users ON shifts.user_id = users.id
WHERE shifts.date_id = ${req.date.id}
  AND shifts.weather_id = ${req.weather.id}
  AND shifts.task_id = ${req.task.id}
  AND shifts.time_id = ${req.time.id}
  AND shifts.user_id <> ${req.user.id}
GROUP BY shifts.task_id;
''';
    Map<String, dynamic> data = await database.find(ctx, sql);
    return TaskDetail(
      id: data['id'],
      task: data['task'],
      place: data['place'],
      url: data['url'],
      superviser: data['superviser'],
      yearId: data['year_id'],
      notes: data['notes'].toString(),
      users: data['users'].toString(),
      createdAt: data['created_at'],
      updatedAt: data['updated_at'],
      deletedAt: data['deleted_at'],
    );
  }

  @override
  Future<int> countUserFromTask(ctx, Shift req) async {
    String sql = '''
SELECT COUNT(shifts.user_id) as count
FROM shifts
LEFT JOIN tasks ON shifts.task_id = tasks.id
LEFT JOIN users ON shifts.user_id = users.id
WHERE shifts.date_id = ${req.date.id}
  AND shifts.weather_id = ${req.weather.id}
  AND shifts.task_id = ${req.task.id}
  AND shifts.time_id = ${req.time.id}
GROUP BY shifts.task_id;
    ''';

    Map<String, dynamic> data = await database.find(ctx, sql);
    return data['count'];
  }
}
