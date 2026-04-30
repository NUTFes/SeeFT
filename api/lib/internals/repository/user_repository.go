package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
	"github.com/lib/pq"
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
	query := "SELECT * FROM users WHERE id = $1"
	return ur.crud.ReadByID(c, query, id)
}

// 学籍番号から取得
func (ur *userRepository) FindByStudentNumber(c context.Context, studentNumber string) *sql.Row {
	query := "SELECT * FROM users WHERE student_number = $1"
	row := ur.client.DB().QueryRowContext(c, query, studentNumber)
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
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	return ur.crud.UpdateDB(c, query, name, mail, gradeID, departmentID, bureauID, roleID, studentNumber, tel, password)
}

// 編集
func (ur *userRepository) Update(c context.Context, id string, name string, mail string, gradeID string, departmentID string, bureauID string, roleID string, studentNumber string, tel string, password string) error {
	query := `
		UPDATE
			users
		SET
			name = $1,
			mail = $2,
			grade_id = $3,
			department_id = $4,
			bureau_id = $5,
			role_id = $6,
			student_number = $7,
			tel = $8,
			password = $9
		WHERE id = $10`
	return ur.crud.UpdateDB(c, query, name, mail, gradeID, departmentID, bureauID, roleID, studentNumber, tel, password, id)
}

// 削除
func (ur *userRepository) Delete(c context.Context, id string) error {
	query := "DELETE FROM users WHERE id = $1"
	return ur.crud.UpdateDB(c, query, id)
}

func (ur *userRepository) FindNewRecord(c context.Context) (*sql.Row, error) {
	query := "SELECT * FROM users ORDER BY id DESC LIMIT 1"
	return ur.crud.ReadByID(c, query)
}

// ユーザ名からユーザを取得する
func (b *userRepository) FindByName(c context.Context, name string) (*sql.Row, error) {
	query := "SELECT * FROM users WHERE name = $1"
	return b.client.DB().QueryRowContext(c, query, name), nil
}

// 複数のユーザー名から一括でユーザーを取得する（N+1問題対策）
func (b *userRepository) FindByNames(c context.Context, names []string) (*sql.Rows, error) {
	if len(names) == 0 {
		// 空の結果を返す
		query := "SELECT * FROM users WHERE 1=0"
		return b.client.DB().QueryContext(c, query)
	}

	query := "SELECT * FROM users WHERE name = ANY($1::text[])"
	return b.client.DB().QueryContext(c, query, pq.Array(names))
}
