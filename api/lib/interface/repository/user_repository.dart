import '../../usecase/repository/user_repository.dart';
import './external/database.dart';
import '../../entity/entity.dart';

class UserRepositoryImpl implements UserRepository {
  Database database;

  UserRepositoryImpl(this.database);

  @override
  Future<List<User>> getUsers(ctx) async {
    String sql = '''
SELECT * FROM users;
''';

    List<Map<String, dynamic>> data = await database.finds(ctx, sql);
    List<User> list = [];

    for (Map<String, dynamic> d in data) {
      User user = _convertUser(d);

      if (user.isDeleted) {
        continue;
      }

      list.add(user);
    }

    return list;
  }

  @override
  Future<User> getUser(ctx, int id) async {
    String sql = '''
SELECT * FROM users WHERE id=$id;
''';
    Map<String, dynamic> data = await database.find(ctx, sql);

    return _convertUser(data);
  }

  @override
  Future<User> insertUser(ctx, User req) async {
    String sql = '''
INSERT INTO users (name, bureau_id, grade_id)
    VALUES ("${req.name}", "${req.bureauId}", "${req.gradeId}") 
    returning *;
''';

    var data = await database.insert(ctx, sql);
    if (data["id"] == 0) {
      print("error at UserRepository.insertUser");
    }

    return _convertUser(data);
  }

  @override
  Future<User> updateUser(ctx, User user) async {
    String sql = '''
UPDATE users SET name="${user.name}" WHERE id=${user.id};
''';
    String returningSQL = '''
SELECT * FROM users WHERE id=${user.id};
''';

    var data = await database.update(ctx, sql, returningSQL);

    return _convertUser(data);
  }

  @override
  Future<User> deleteUser(ctx, User user) async {
    String sql = '''
UPDATE users SET deleted_at=NOW() WHERE id=${user.id};
''';
    String returningSQL = '''
SELECT * FROM users WHERE id=${user.id};
''';

    var data = await database.update(ctx, sql, returningSQL);

    return _convertUser(data);
  }

  @override
  Future<User> getUserByMail(ctx, User req) async {
    String sql = '''
SELECT id, name, mail, bureau_id, grade_id, created_at, updated_at, deleted_at
FROM users WHERE mail="${req.mail}";
''';

    Map<String, dynamic> data = await database.find(ctx, sql);
    print(data);

    return _convertUser(data);
  }

  User _convertUser(Map<String, dynamic> data) {
    return User(
      id: data['id'],
      name: data['name'],
      mail: data['mail'],
      bureauId: data['bureau_id'],
      gradeId: data['grade_id'],
      createdAt: data['created_at'],
      updatedAt: data['updated_at'],
      deletedAt: data['deleted_at'],
    );
  }
}
