package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
)

var userDebugSQL = os.Getenv("DEBUG_SQL") != "0"

type userRepository struct {
	client db.Client
	crud   abstract.Crud
}

type UserRepository interface {
	All(context.Context) (*sql.Rows, error)
	Find(context.Context, string) (*sql.Row, error)
	FindByStudentNumber(context.Context, string) *sql.Row
	Create(context.Context, string, string, string, string, string, string, string, string, string) error
	Update(context.Context, string, string, string, string, string, string, string, string, string, string) error
	Delete(context.Context, string) error
	FindNewRecord(context.Context) (*sql.Row, error)
	FindByName(context.Context, string) (*sql.Row, error)
	FindByNames(context.Context, []string) (*sql.Rows, error)
}

func NewUserRepository(c db.Client, ac abstract.Crud) UserRepository {
	return &userRepository{c, ac}
}

// 全件取得
func (ur *userRepository) All(c context.Context) (*sql.Rows, error) {
	query := "SELECT * FROM users"
	return ur.crud.Read(c, query)
}

// 1件取得
func (ur *userRepository) Find(c context.Context, id string) (*sql.Row, error) {
	query := "SELECT * FROM users WHERE id = " + id
	return ur.crud.ReadByID(c, query)
}

// 学籍番号から取得
func (ur *userRepository) FindByStudentNumber(c context.Context, studentNumber string) *sql.Row {
	query := "SELECT * FROM users WHERE student_number = " + studentNumber
	row := ur.client.DB().QueryRowContext(c, query)
	if userDebugSQL {
		fmt.Printf("\x1b[36m%s\n", query)
	}
	return row
}

// 作成
func (ur *userRepository) Create(c context.Context, name string, mail string, gradeID string, departmentID string, bureauID string, roleID string, studentNumber string, tel string, password string) error {
	query := `
		INSERT INTO
			users (name, mail, grade_id, department_id, bureau_id, role_id, student_number, tel, password)
		VALUES ('` + name + "', '" + mail + "', " + gradeID + ", " + departmentID + ", " + bureauID + ", " + roleID + ", " + studentNumber + ", '" + tel + "', '" + password + "')"
	return ur.crud.UpdateDB(c, query)
}

// 編集
func (ur *userRepository) Update(c context.Context, id string, name string, mail string, gradeID string, departmentID string, bureauID string, roleID string, studentNumber string, tel string, password string) error {
	query := `
		UPDATE
			users
		SET
			name = '` + name +
		"', mail = '" + mail +
		"', grade_id = " + gradeID +
		", department_id = " + departmentID +
		", bureau_id = " + bureauID +
		", role_id = " + roleID +
		", student_number = " + studentNumber +
		", tel = '" + tel +
		"', password = '" + password +
		"' WHERE id = " + id
	return ur.crud.UpdateDB(c, query)
}

// 削除
func (ur *userRepository) Delete(c context.Context, id string) error {
	query := "DELETE FROM users WHERE id = " + id
	return ur.crud.UpdateDB(c, query)
}

func (ur *userRepository) FindNewRecord(c context.Context) (*sql.Row, error) {
	query := "SELECT * FROM users ORDER BY id DESC LIMIT 1"
	return ur.crud.ReadByID(c, query)
}

// ユーザ名からユーザを取得する
func (b *userRepository) FindByName(c context.Context, name string) (*sql.Row, error) {
	query := "SELECT * FROM users WHERE name = '" + name + "'"
	return b.client.DB().QueryRowContext(c, query), nil
}

// 複数のユーザー名から一括でユーザーを取得する（N+1問題対策）
func (b *userRepository) FindByNames(c context.Context, names []string) (*sql.Rows, error) {
	if len(names) == 0 {
		// 空の結果を返す
		query := "SELECT * FROM users WHERE 1=0"
		return b.client.DB().QueryContext(c, query)
	}

	// IN句用のプレースホルダーを作成
	placeholders := make([]interface{}, len(names))
	for i, name := range names {
		placeholders[i] = name
	}

	// IN句を動的に構築
	query := "SELECT * FROM users WHERE name IN ("
	for i := range names {
		if i > 0 {
			query += ", "
		}
		query += "?"
	}
	query += ")"

	return b.client.DB().QueryContext(c, query, placeholders...)
}
