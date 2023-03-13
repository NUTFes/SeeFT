import '../entity/entity.dart';
import './repository/repository.dart';

abstract class UserUsecase {
  Future<List<User>> getUsers(ctx);
  Future<User> getUser(ctx, int id);
  Future<User> insertUser(ctx, User req);
  Future<User> updateUser(ctx, User req);
  Future<User> deleteUser(ctx, User req);
}

class UserUsecaseImpl implements UserUsecase {
  UserRepository userRepository;

  UserUsecaseImpl(this.userRepository);

  @override
  Future<List<User>> getUsers(ctx) async {
    List<User> users = await userRepository.getUsers(ctx);
    return users;
  }

  @override
  Future<User> getUser(ctx, int id) async {
    User user = await userRepository.getUser(ctx, id);
    return user;
  }

  @override
  Future<User> insertUser(ctx, User req) async {
    User user = await userRepository.insertUser(ctx, req);
    return user;
  }

  @override
  Future<User> updateUser(ctx, User req) async {
    User test = await userRepository.getUser(ctx, req.id);
    if (req.name == test.name) {
      throw Exception('request name is same.');
    }
    User user = await userRepository.updateUser(ctx, req);
    if (test.updatedAt == user.updatedAt) {
      throw Exception('cant updated because request same response');
    }
    return user;
  }

  @override
  Future<User> deleteUser(ctx, User req) async {
    User user = await userRepository.deleteUser(ctx, req);
    return user;
  }
}
