import '../../usecase/repository/shift_repository.dart';
import '../../entity/entity.dart';
import './external/database.dart';

class ShiftRepositoryImpl implements ShiftRepository {
  Database database;

  ShiftRepositoryImpl(this.database);

  @override
  Future<Shift> getShift(ctx, Shift req) async {
    String sql = '''
SELECT * FROM shifts
WHERE shifts.id = ${req.id};
''';

    Map<String, dynamic> data = await database.find(ctx, sql);
    return Shift(
      id: data['id'],
      user: User(id: data['user_id']),
      task: Task(id: data['task_id']),
      year: Year(id: data['year_id']),
      date: Date(id: data['date_id']),
      time: Time(id: data['time_id']),
      weather: Weather(id: data['weather_id']),
      isAttendance: data['is_attendance'] == 1,
      createdUserId: data['created_user_id'],
      updatedUserId: data['updated_user_id'],
      createdAt: data['created_at'],
      updatedAt: data['updated_at'],
      deletedAt: data['deleted_at'],
    );
  }

  @override
  Future<List<Shift>> getShiftsByUser(ctx, User req) async {
    String sql = '''
SELECT
  shifts.id,
  shifts.user_id,
  users.name as user_name,
  users.bureau_id as user_bureau_id,
  users.grade_id as user_grade_id,
  users.created_at as user_created_at,
  users.updated_at as user_updated_at,
  users.deleted_at as user_deleted_at,
  shifts.task_id,
  tasks.task as task_name,
  tasks.color as task_color,
  tasks.place as task_place,
  tasks.url as task_url,
  tasks.superviser as task_superviser,
  tasks.notes as task_notes,
  tasks.created_at as task_created_at,
  tasks.updated_at as task_updated_at,
  tasks.deleted_at as task_deleted_at,
  shifts.year_id,
  years.year as year_name,
  years.created_at as year_created_at,
  years.updated_at as year_updated_at,
  years.deleted_at as year_deleted_at,
  shifts.date_id,
  dates.date as date_name,
  dates.created_at as date_created_at,
  dates.updated_at as date_updated_at,
  dates.deleted_at as date_deleted_at,
  shifts.time_id,
  times.time as time_name,
  times.created_at as time_created_at,
  times.updated_at as time_updated_at,
  times.deleted_at as time_deleted_at,
  shifts.weather_id,
  weathers.weather as weather_name,
  weathers.created_at as weather_created_at,
  weathers.updated_at as weather_updated_at,
  weathers.deleted_at as weather_deleted_at,
  shifts.is_attendance,
  shifts.created_user_id,
  shifts.updated_user_id,
  shifts.created_at as shift_created_at,
  shifts.updated_at as shift_updated_at,
  shifts.deleted_at as shift_deleted_at
FROM shifts
LEFT JOIN users ON shifts.user_id = users.id
LEFT JOIN tasks ON shifts.task_id = tasks.id
LEFT JOIN years ON shifts.year_id = years.id
LEFT JOIN dates ON shifts.date_id = dates.id
LEFT JOIN times ON shifts.time_id = times.id
LEFT JOIN weathers ON shifts.weather_id = weathers.id
WHERE user_id=${req.id};
''';

    List<Map<String, dynamic>> data = await database.finds(ctx, sql);
    List<Shift> list = [];

    for (var d in data) {
      Shift shift = _convertShift(d);

      if (shift.isDeleted) {
        continue;
      }

      list.add(shift);
    }

    return list;
  }

  @override
  Future<List<Shift>> getShiftsByUserAndDateAndWeather(ctx, Shift req) async {
    String sql = '''
SELECT
  shifts.id,
  shifts.user_id,
  users.name as user_name,
  users.bureau_id as user_bureau_id,
  users.grade_id as user_grade_id,
  users.created_at as user_created_at,
  users.updated_at as user_updated_at,
  users.deleted_at as user_deleted_at,
  shifts.task_id,
  tasks.task as task_name,
  tasks.color as task_color,
  tasks.place as task_place,
  tasks.url as task_url,
  tasks.superviser as task_superviser,
  tasks.notes as task_notes,
  tasks.created_at as task_created_at,
  tasks.updated_at as task_updated_at,
  tasks.deleted_at as task_deleted_at,
  shifts.year_id,
  years.year as year_name,
  years.created_at as year_created_at,
  years.updated_at as year_updated_at,
  years.deleted_at as year_deleted_at,
  shifts.date_id,
  dates.date as date_name,
  dates.created_at as date_created_at,
  dates.updated_at as date_updated_at,
  dates.deleted_at as date_deleted_at,
  shifts.time_id,
  times.time as time_name,
  times.created_at as time_created_at,
  times.updated_at as time_updated_at,
  times.deleted_at as time_deleted_at,
  shifts.weather_id,
  weathers.weather as weather_name,
  weathers.created_at as weather_created_at,
  weathers.updated_at as weather_updated_at,
  weathers.deleted_at as weather_deleted_at,
  shifts.is_attendance,
  shifts.created_user_id,
  shifts.updated_user_id,
  shifts.created_at as shift_created_at,
  shifts.updated_at as shift_updated_at,
  shifts.deleted_at as shift_deleted_at
FROM shifts
LEFT JOIN users ON shifts.user_id = users.id
LEFT JOIN tasks ON shifts.task_id = tasks.id
LEFT JOIN years ON shifts.year_id = years.id
LEFT JOIN dates ON shifts.date_id = dates.id
LEFT JOIN times ON shifts.time_id = times.id
LEFT JOIN weathers ON shifts.weather_id = weathers.id
WHERE user_id=${req.user.id} AND
  weather_id=${req.weather.id} AND
  date_id=${req.date.id};
''';

    List<Map<String, dynamic>> data = await database.finds(ctx, sql);
    List<Shift> list = [];

    for (var d in data) {
      Shift shift = _convertShift(d);
      if (shift.isDeleted) {
        continue;
      }

      list.add(shift);
    }

    return list;
  }

  @override
  Future<List<Shift>> getShiftsByYearAndDateAndWeather(ctx, Shift req) async {
    String sql = '''
SELECT
  shifts.id,
  shifts.user_id,
  users.name as user_name,
  users.bureau_id as user_bureau_id,
  users.grade_id as user_grade_id,
  users.created_at as user_created_at,
  users.updated_at as user_updated_at,
  users.deleted_at as user_deleted_at,
  shifts.task_id,
  tasks.task as task_name,
  tasks.color as task_color,
  tasks.place as task_place,
  tasks.url as task_url,
  tasks.superviser as task_superviser,
  tasks.notes as task_notes,
  tasks.created_at as task_created_at,
  tasks.updated_at as task_updated_at,
  tasks.deleted_at as task_deleted_at,
  shifts.year_id,
  years.year as year_name,
  years.created_at as year_created_at,
  years.updated_at as year_updated_at,
  years.deleted_at as year_deleted_at,
  shifts.date_id,
  dates.date as date_name,
  dates.created_at as date_created_at,
  dates.updated_at as date_updated_at,
  dates.deleted_at as date_deleted_at,
  shifts.time_id,
  times.time as time_name,
  times.created_at as time_created_at,
  times.updated_at as time_updated_at,
  times.deleted_at as time_deleted_at,
  shifts.weather_id,
  weathers.weather as weather_name,
  weathers.created_at as weather_created_at,
  weathers.updated_at as weather_updated_at,
  weathers.deleted_at as weather_deleted_at,
  shifts.is_attendance,
  shifts.created_user_id,
  shifts.updated_user_id,
  shifts.created_at as shift_created_at,
  shifts.updated_at as shift_updated_at,
  shifts.deleted_at as shift_deleted_at
FROM shifts
LEFT JOIN users ON shifts.user_id = users.id
LEFT JOIN tasks ON shifts.task_id = tasks.id
LEFT JOIN years ON shifts.year_id = years.id
LEFT JOIN dates ON shifts.date_id = dates.id
LEFT JOIN times ON shifts.time_id = times.id
LEFT JOIN weathers ON shifts.weather_id = weathers.id
WHERE shifts.year_id=${req.year.id} AND
      date_id=${req.date.id} AND
      weather_id=${req.weather.id};
''';

    List<Map<String, dynamic>> data = await database.finds(ctx, sql);
    List<Shift> list = [];

    for (var d in data) {
      Shift shift = _convertShift(d);
      if (shift.isDeleted) {
        continue;
      }

      list.add(shift);
    }

    return list;
  }
}

Shift _convertShift(d) {
  Shift shift = Shift(
    id: d['id'],
    user: User(
      id: d['user_id'],
      name: d['user_name'],
      bureauId: d['user_bureau_id'],
      gradeId: d['user_grade_id'],
      createdAt: d['user_created_at'],
      updatedAt: d['user_updated_at'],
      deletedAt: d['user_deleted_at'],
    ),
    task: Task(
      id: d['task_id'],
      task: d['task_name'],
      color: Color(int.parse(d['task_color'], radix: 16)),
      place: d['task_place'],
      url: d['task_url'],
      superviser: d['task_superviser'],
      notes: d['task_notes'].toString(),
      createdAt: d['task_created_at'],
      updatedAt: d['task_updated_at'],
      deletedAt: d['task_deleted_at'],
    ),
    year: Year(
      id: d['year_id'],
      year: d['year_name'],
      createdAt: d['year_created_at'],
      updatedAt: d['year_updated_at'],
      deletedAt: d['year_deleted_at'],
    ),
    date: Date(
      id: d['date_id'],
      date: d['date_name'],
      createdAt: d['date_created_at'],
      updatedAt: d['date_updated_at'],
      deletedAt: d['date_deleted_at'],
    ),
    time: Time(
      id: d['time_id'],
      time: d['time_name'],
      createdAt: d['date_created_at'],
      updatedAt: d['date_updated_at'],
      deletedAt: d['date_deleted_at'],
    ),
    weather: Weather(
      id: d['weather_id'],
      weather: d['weather_name'],
      createdAt: d['weather_created_at'],
      updatedAt: d['weather_updated_at'],
      deletedAt: d['weather_deleted_at'],
    ),
    isAttendance: d['is_attendance'] == 1,
    createdUserId: d['created_user_id'],
    updatedUserId: d['updated_user_id'],
    createdAt: d['shift_created_at'],
    updatedAt: d['shift_updated_at'],
    deletedAt: d['shift_deleted_at'],
  );

  return shift;
}
