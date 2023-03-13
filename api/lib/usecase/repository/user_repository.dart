import '../../entity/entity.dart';

abstract class UserRepository {
  Future<List<User>> getUsers(ctx);
  Future<User> getUser(ctx, int id);
  Future<User> insertUser(ctx, User req);
  Future<User> updateUser(ctx, User user);
  Future<User> deleteUser(ctx, User user);
  Future<User> getUserByMail(ctx, User req);
}
