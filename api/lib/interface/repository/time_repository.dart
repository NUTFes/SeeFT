import '../../usecase/repository/time_repository.dart';
import '../../entity/entity.dart';
import './external/database.dart';

class TimeRepositoryImpl implements TimeRepository {
  Database database;

  TimeRepositoryImpl(this.database);

  @override
  Future<List<Time>> getTimes(ctx) async {
    String sql = '''
SELECT * FROM times;
''';

    List<Map<String, dynamic>> data = await database.finds(ctx, sql);
    List<Time> list = [];

    for (var d in data) {
      Time time = Time(
        id: d['id'],
        time: d['time'],
        createdAt: d['created_at'],
        updatedAt: d['updated_at'],
        deletedAt: d['deleted_at'],
      );

      if (time.isDeleted) {
        continue;
      }

      list.add(time);
    }

    return list;
  }
}
