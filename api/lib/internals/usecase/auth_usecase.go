package usecase

import (
	"context"

	rep "github.com/NUTFes/SeeFT/api/externals/repository"
	"github.com/NUTFes/SeeFT/api/internals/domain"
	"github.com/pkg/errors"
)

// import '../entity/entity.dart';
// import './repository/repository.dart';

// abstract class AuthUsecase {
//   Future<User> signIn(ctx, User req);
// }

// class AuthUsecaseImpl implements AuthUsecase {
//   UserRepository userRepository;

//   AuthUsecaseImpl(this.userRepository);

//   @override
//   Future<User> signIn(ctx, User req) async {
//     final user = await userRepository.getUserByMail(ctx, req);
//     return user;
//   }
// }
