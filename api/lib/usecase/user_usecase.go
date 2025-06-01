package usecase

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/NUTFes/SeeFT/api/lib/entity"
	rep "github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	userRep    rep.UserRepository
	burearRep  rep.BureauRepository
	taskRep    rep.TaskRepository
	sessionRep rep.SessionRepository
}

type UserUseCase interface {
	GetUsers(context.Context) ([]entity.User, error)
	GetUserByID(context.Context, string) (entity.User, error)
	CreateUser(context.Context, string, string, string, string, string, string, string, string, string) (entity.User, error)
	UpdateUser(context.Context, string, string, string, string, string, string, string, string, string) (entity.User, error)
	DeleteUser(context.Context, string) error
	GetCurrentUser(context.Context, string) (entity.User, error)
}

func NewUserUseCase(userRep rep.UserRepository, sessionRep rep.SessionRepository) UserUseCase {
	return &userUseCase{userRep: userRep, sessionRep: sessionRep}
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
			&user.Mail,
			&user.GradeID,
			&user.DepartmentID,
			&user.BureauID,
			&user.RoleID,
			&user.StudentNumber,
			&user.Tel,
			&user.Password,
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
		&user.Mail,
		&user.GradeID,
		&user.DepartmentID,
		&user.BureauID,
		&user.RoleID,
		&user.StudentNumber,
		&user.Tel,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (u *userUseCase) CreateUser(c context.Context, name string, mail string, gradeID string, departmentID string, bureauID string, roleID string, studentNumber string, tel string, password string) (entity.User, error) {
	latastUser := entity.User{}
	password = strings.ReplaceAll(password, " ", "")
	password = strings.ReplaceAll(password, "　", "")
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	err := u.userRep.Create(c, name, mail, gradeID, departmentID, bureauID, roleID, studentNumber, tel, string(hashedPassword))
	row, err := u.userRep.FindNewRecord(c)
	err = row.Scan(
		&latastUser.ID,
		&latastUser.Name,
		&latastUser.Mail,
		&latastUser.GradeID,
		&latastUser.DepartmentID,
		&latastUser.BureauID,
		&latastUser.RoleID,
		&latastUser.StudentNumber,
		&latastUser.Tel,
		&latastUser.Password,
		&latastUser.CreatedAt,
		&latastUser.UpdatedAt,
	)
	if err != nil {
		return latastUser, err
	}
	return latastUser, err
}

func (u *userUseCase) UpdateUser(c context.Context, id string, name string, mail string, gradeID string, departmentID string, bureauID string, roleID string, studentNumber string, tel string) (entity.User, error) {
	updatedUser := entity.User{}
	var user entity.User

	row, err := u.userRep.Find(c, id)
	err = row.Scan(
		&user.ID,
		&user.Name,
		&user.Mail,
		&user.GradeID,
		&user.DepartmentID,
		&user.BureauID,
		&user.RoleID,
		&user.StudentNumber,
		&user.Tel,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}

	u.userRep.Update(c, id, name, mail, gradeID, departmentID, bureauID, roleID, studentNumber, tel, user.Password)
	row, err = u.userRep.Find(c, id)
	err = row.Scan(
		&updatedUser.ID,
		&updatedUser.Name,
		&updatedUser.Mail,
		&updatedUser.GradeID,
		&updatedUser.DepartmentID,
		&updatedUser.BureauID,
		&updatedUser.RoleID,
		&updatedUser.StudentNumber,
		&updatedUser.Tel,
		&updatedUser.Password,
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

func (u *userUseCase) GetCurrentUser(c context.Context, accessToken string) (entity.User, error) {
	var session = entity.Session{}
	var user = entity.User{}
	var row *sql.Row
	var err error
	// アクセストークンからmail_authを取得
	row = u.sessionRep.FindSessionByAccessToken(c, accessToken)
	err = row.Scan(
		&session.ID,
		&session.UserID,
		&session.AccessToken,
		&session.CreatedAt,
		&session.UpdatedAt,
	)
	if err != nil {
		return user, err
	}

	// userIDの該当するuserを取得
	row, err = u.userRep.Find(c, strconv.Itoa(session.UserID))
	err = row.Scan(
		&user.ID,
		&user.Name,
		&user.Mail,
		&user.GradeID,
		&user.DepartmentID,
		&user.BureauID,
		&user.RoleID,
		&user.StudentNumber,
		&user.Tel,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}
	return user, nil
}

// GASからのユーザー変更通知を受けてDBを更新
func (u *userUseCase) UpdateUsersFromGAS(ctx context.Context, req entity.UserChangeRequest) error {
	for _, change := range req.Changes {
		// 局名からBureauIDを取得
		var bureauID string
		bureau := strings.ReplaceAll(change.Bureau, " ", "")
		bureau = strings.ReplaceAll(bureau, "　", "")
		switch bureau {
		case `執行部`:
			bureauID = "1"
		case `執行部補佐`:
			bureauID = "2"
		case `総務局`:
			bureauID = "3"
		case `企画局`:
			bureauID = "4"
		case `渉外局`:
			bureauID = "5"
		case `財務局`:
			bureauID = "6"
		case `制作局`:
			bureauID = "7"
		case `情報局`:
			bureauID = "8"
		default:
			bureauID = "0"
		}
		// 学年からGradeIDを取得
		var gradeID string
		grade := strings.ReplaceAll(change.Grade, " ", "")
		grade = strings.ReplaceAll(grade, "　", "")
		switch grade {
		case `B1`:
			gradeID = "1"
		case `B2`:
			gradeID = "2"
		case `B3`:
			gradeID = "3"
		case `B4`:
			gradeID = "4"
		case `M1`:
			gradeID = "5"
		case `M2`:
			gradeID = "6"
		case `D1`:
			gradeID = "7"
		case `D2`:
			gradeID = "8"
		case `D3`:
			gradeID = "9"
		case `OB`:
			gradeID = "10"
		default:
			gradeID = "0"
		}
		// 学科名からDepartmentIDを取得
		var departmentID string
		department := strings.ReplaceAll(change.Department, " ", "")
		department = strings.ReplaceAll(department, "　", "")
		switch department {
		case `未所属`:
			departmentID = "1"
		case `機械工学分野`:
			departmentID = "2"
		case `電気電子情報工学分野`:
			departmentID = "3"
		case `情報・経営システム工学分野`:
			departmentID = "4"
		case `物質生物工学分野`:
			departmentID = "5"
		case `環境社会基盤工学分野`:
			departmentID = "6"
		case `量子・原子力統合工学分野`:
			departmentID = "7"
		case `技術科学イノベーション`:
			departmentID = "8"
		default:
			departmentID = "0"
		}

		studentNumber := strconv.Itoa(change.StudentNumber)
		tel := change.Tel

		// 1. ユーザー名からUserID取得
		userRow, _ := u.userRep.FindByName(ctx, change.Name) // Rowはユーザー名が入っている前提
		var user entity.User
		if err := userRow.Scan(&user.ID, &user.Name, &user.Mail, gradeID, departmentID, bureauID, &user.RoleID, studentNumber, tel, &user.Password, &user.CreatedAt, &user.UpdatedAt); err == nil {
			// ユーザーが存在すれば更新
			u.userRep.Update(ctx, strconv.Itoa(user.ID), change.Name, "", gradeID, departmentID, bureauID, strconv.Itoa(user.RoleID), studentNumber, tel, user.Password)
		} else { // ユーザーがいなければ新規作成
			if err.Error() == "sql: no rows in result set" {
				// 必要な情報は仮値でOK（必要に応じて修正）
				name := change.Name
				mail := ""
				roleID := "1"
				password := "password" // 仮のパスワード（必要に応じて変更）
				hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
				createErr := u.userRep.Create(ctx, name, mail, gradeID, departmentID, bureauID, roleID, studentNumber, tel, string(hashed))
				if createErr != nil {
					return errors.Wrapf(createErr, "ユーザー新規作成失敗: %v", change.Name)
				}
				// 再取得
				userRow, _ = u.userRep.FindByName(ctx, change.Name)
				if err := userRow.Scan(&user.ID, &user.Name, &user.Mail, &user.GradeID, &user.DepartmentID, &user.BureauID, &user.RoleID, &user.StudentNumber, &user.Tel, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
					return errors.Wrapf(err, "ユーザー再取得失敗: %v", change.Name)
				}
			} else {
				return errors.Wrapf(err, "ユーザー取得失敗: %v", change.Name)
			}
		}
	}
	return nil
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
