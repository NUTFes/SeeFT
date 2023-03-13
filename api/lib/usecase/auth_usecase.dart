import '../entity/entity.dart';
import './repository/repository.dart';

abstract class AuthUsecase {
  Future<User> signIn(ctx, User req);
}

class AuthUsecaseImpl implements AuthUsecase {
  UserRepository userRepository;

  AuthUsecaseImpl(this.userRepository);

  @override
  Future<User> signIn(ctx, User req) async {
    final user = await userRepository.getUserByMail(ctx, req);
    return user;
  }
}
