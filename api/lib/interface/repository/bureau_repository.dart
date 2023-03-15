import '../../usecase/repository/bureau_repository.dart';
import '../../entity/entity.dart';
import './external/database.dart';

class BureauRepositoryImpl implements BureauRepository {
  Database database;

  BureauRepositoryImpl(this.database);

  @override
  Future<List<Bureau>> getBureaus(ctx) async {
    String sql = '''
SELECT * FROM bureaus;
''';

    List<Map<String, dynamic>> data = await database.finds(ctx, sql);
    List<Bureau> list = [];

    for (var d in data) {
      Bureau bureau = Bureau(
        id: d['id'],
        bureau: d['bureau'],
        color: Color(int.parse(d['color'], radix: 16)),
        createdAt: d['created_at'],
        updatedAt: d['updated_at'],
        deletedAt: d['deleted_at'],
      );

      if (bureau.isDeleted) {
        continue;
      }

      list.add(bureau);
    }

    return list;
  }
}
