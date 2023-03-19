package usecase

import (
  "context"

	rep "github.com/NUTFes/SeeFT/api/lib/externals/repository"
	"github.com/NUTFes/SeeFT/api/lib/internals/entity"
	"github.com/pkg/errors"
)

type userUseCase struct {
  userRep rep.UserRepository
}

type UserUseCase interface {
  GetUsers(context.Context) ([]entity.User, error)
  GetUserByID(context.Context, string) (entity.User, error)
  CreateUser(context.Context, string, string, string, string) (entity.User, error)
  UpdateUser(context.Context, string, string, string, string, string) (entity.User, error)
  DeleteUser(context.Context, string) error
}

func NewUserUseCase(rep rep.UserRepository) UserUseCase {
  return &userUseCase{rep}
}

func (u *userUseCase) GetUsers(c context.Context) ([]entity.User, error) {

	user := entity.User{}
	var users []entity.User

	rows, err := u.userRep.All(c)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.GradeID,
			&user.BureauID,
			&user.RoleID,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return nil, errors.Wrapf(err, "cannot connect SQL")
		}

		users = append(users, user)
	}
	return users, nil
}

func (u *userUseCase) GetUserByID(c context.Context, id string) (entity.User, error) {
	var user entity.User

	row, err := u.userRep.Find(c, id)
	err = row.Scan(
		&user.ID,
		&user.Name,
		&user.GradeID,
		&user.BureauID,
		&user.RoleID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (u *userUseCase) CreateUser(c context.Context, name string, gradeID string, bureauID string, roleID string) (entity.User, error) {
	latastUser := entity.User{}
	err := u.userRep.Create(c, name, gradeID, bureauID, roleID)
	row, err := u.userRep.FindNewRecord(c)
	err = row.Scan(
		&latastUser.ID,
		&latastUser.Name,
		&latastUser.GradeID,
		&latastUser.BureauID,
		&latastUser.RoleID,
		&latastUser.CreatedAt,
		&latastUser.UpdatedAt,
	)
	if err != nil {
		return latastUser, err
	}
	return latastUser, err
}

func (u *userUseCase) UpdateUser(c context.Context, id string, name string, gradeID string, bureauID string, roleID string) (entity.User, error) {
	updatedUser := entity.User{}
	u.userRep.Update(c, id, name, gradeID, bureauID, roleID)
	row, err := u.userRep.Find(c, id)
	err = row.Scan(
		&updatedUser.ID,
		&updatedUser.Name,
		&updatedUser.GradeID,
		&updatedUser.BureauID,
		&updatedUser.RoleID,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)
	if err != nil {
		return updatedUser, err
	}
	return updatedUser, nil
}

func (u *userUseCase) DeleteUser(c context.Context, id string) error {
	err := u.userRep.Delete(c, id)
	return err
}

// import '../entity/entity.dart';
// import './repository/repository.dart';

// abstract class UserUsecase {
//   Future<List<User>> getUsers(ctx);
//   Future<User> getUser(ctx, int id);
//   Future<User> insertUser(ctx, User req);
//   Future<User> updateUser(ctx, User req);
//   Future<User> deleteUser(ctx, User req);
// }

// class UserUsecaseImpl implements UserUsecase {
//   UserRepository userRepository;

//   UserUsecaseImpl(this.userRepository);

//   @override
//   Future<List<User>> getUsers(ctx) async {
//     List<User> users = await userRepository.getUsers(ctx);
//     return users;
//   }

//   @override
//   Future<User> getUser(ctx, int id) async {
//     User user = await userRepository.getUser(ctx, id);
//     return user;
//   }

//   @override
//   Future<User> insertUser(ctx, User req) async {
//     User user = await userRepository.insertUser(ctx, req);
//     return user;
//   }

//   @override
//   Future<User> updateUser(ctx, User req) async {
//     User test = await userRepository.getUser(ctx, req.id);
//     if (req.name == test.name) {
//       throw Exception('request name is same.');
//     }
//     User user = await userRepository.updateUser(ctx, req);
//     if (test.updatedAt == user.updatedAt) {
//       throw Exception('cant updated because request same response');
//     }
//     return user;
//   }

//   @override
//   Future<User> deleteUser(ctx, User req) async {
//     User user = await userRepository.deleteUser(ctx, req);
//     return user;
//   }
// }
