package usecase

import (
	"context"
	// "fmt"
	// "crypto/rand"
	// "strconv"

	rep "github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/NUTFes/SeeFT/api/lib/entity"
	"golang.org/x/crypto/bcrypt"
	// "github.com/pkg/errors"
)

type mailAuthUseCase struct {
	userRep rep.UserRepository
	sessionRep  rep.SessionRepository
}

type MailAuthUseCase interface {
	SignIn(context.Context, string, string) (entity.LoginUser, error)
}

func NewAuthUseCase(userRep rep.UserRepository, sessionRep rep.SessionRepository) MailAuthUseCase {
	return &mailAuthUseCase{userRep: userRep, sessionRep: sessionRep}
}

func (u *mailAuthUseCase) SignIn(c context.Context, studentNumber string, password string) (entity.LoginUser, error) {
	var user = entity.User{}
	
	// メールアドレスの存在確認
	row := u.userRep.FindByStudentNumber(c, studentNumber)
	err := row.Scan(
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
	// パスワードがあっているか確認
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	loginUser := entity.LoginUser{ID: user.ID, RoleID: user.RoleID, Mail: user.Mail}

	if err != nil {
		return loginUser, err
	}

	return loginUser, nil
	
	// // トークン発行
	// accessToken, err := _makeRandomStr(10)

	// // ログイン (セッション開始)
	// err = u.sessionRep.Create(c, strconv.FormatInt(int64(user.ID), 10), strconv.Itoa(int(user.ID)), accessToken)
	// if err != nil {
	// 	return token, err
	// }
	// token.AccessToken = accessToken

}

// アクセストークンを生成
// func _makeRandomStr(digit uint32) (string, error) {
// 	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 	// 乱数を生成
// 	b := make([]byte, digit)
// 	if _, err := rand.Read(b); err != nil {
// 		return "", errors.New("unexpected error...")
// 	}

// 	// letters からランダムに取り出して文字列を生成
// 	var result string
// 	for _, v := range b {
// 		// index が letters の長さに収まるように調整
// 		result += string(letters[int(v)%len(letters)])
// 	}
// 	return result, nil
// }


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
