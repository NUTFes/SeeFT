package usecase

import (
	"context"
	"fmt"
	"strings"
	"strconv"
	"crypto/rand"

	rep "github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/NUTFes/SeeFT/api/lib/entity"
	"golang.org/x/crypto/bcrypt"
	"github.com/pkg/errors"
)

type mailAuthUseCase struct {
	userRep rep.UserRepository
	sessionRep  rep.SessionRepository
}

type MailAuthUseCase interface {
	SignIn(context.Context, string, string) (entity.LoginUser, error)
	WebSignUp(context.Context, string, string, string, string, string, string, string, string, string) (entity.Token, error)
	WebSignIn(context.Context, string, string) (entity.Token, error)
	WebSignOut(context.Context, string) error
	WebIsSignIn(context.Context, string) (entity.IsSignIn, error)
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
	password = strings.ReplaceAll(password, " ", "")
	password = strings.ReplaceAll(password, "　", "")
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	loginUser := entity.LoginUser{ID: user.ID, RoleID: user.RoleID, Mail: user.Mail}

	if err != nil {
		fmt.Println(err)
		return loginUser, err
	}

	return loginUser, nil
}

func (u *mailAuthUseCase) WebSignUp(c context.Context, name string, mail string, gradeID string, departmentID string, bureauID string, roleID string, studentNumber string, tel string, password string) (entity.Token, error) {
	var token entity.Token
	var user entity.User
	// パスワードをハッシュ化
	password = strings.ReplaceAll(password, " ", "")
	password = strings.ReplaceAll(password, "　", "")
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	err := u.userRep.Create(c, name, mail, gradeID, departmentID, bureauID, roleID, studentNumber, tel, string(hashed))
	row, err := u.userRep.FindNewRecord(c)
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
		return token, err
	}

	// トークン発行
	accessToken, err := _makeRandomStr(10)
	if err != nil {
		return token, err
	}
	// ログイン（セッション開始）
	err = u.sessionRep.Create(c, strconv.Itoa(int(user.ID)), accessToken)
	if err != nil {
		return token, err
	}
	token.AccessToken = accessToken
	return token, nil
}

func (u *mailAuthUseCase) WebSignIn(c context.Context, studentNumber string, password string) (entity.Token, error) {
	var token entity.Token
	var user entity.User
	
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
	u.sessionRep.DeleteByUserID(c, strconv.Itoa(int(user.ID)))

	// パスワードがあっているか確認
	password = strings.ReplaceAll(password, " ", "")
	password = strings.ReplaceAll(password, "　", "")
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return token, err
	}
	// トークン発行
	accessToken, err := _makeRandomStr(10)
	// ログイン (セッション開始)
	err = u.sessionRep.Create(c, strconv.Itoa(int(user.ID)), accessToken)
	if err != nil {
		return token, err
	}
	token.AccessToken = accessToken
	return token, nil

}

func (u *mailAuthUseCase) WebSignOut(c context.Context, accessToken string) error {
	err := u.sessionRep.Delete(c, accessToken)
	if err != nil {
		return err
	}
	return nil
}

func (u *mailAuthUseCase) WebIsSignIn(c context.Context, accessToken string) (entity.IsSignIn, error) {
	var session = entity.Session{}
	var isSignIn entity.IsSignIn
	row := u.sessionRep.FindSessionByAccessToken(c, accessToken)
	_ = row.Scan(
		&session.ID,
		&session.UserID,
		&session.AccessToken,
		&session.CreatedAt,
		&session.UpdatedAt,
	)
	if session.ID != 0 {
		isSignIn = entity.IsSignIn{IsSignIn: true}
	} else {
		isSignIn = entity.IsSignIn{IsSignIn: false}
	}
	return isSignIn, nil
}

// アクセストークンを生成
func _makeRandomStr(digit uint32) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 乱数を生成
	b := make([]byte, digit)
	if _, err := rand.Read(b); err != nil {
		return "", errors.New("unexpected error...")
	}

	// letters からランダムに取り出して文字列を生成
	var result string
	for _, v := range b {
		// index が letters の長さに収まるように調整
		result += string(letters[int(v)%len(letters)])
	}
	return result, nil
}


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
