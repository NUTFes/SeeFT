package repository

import (
	"context"
	"database/sql"

	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
)

type departmentRepository struct {
	client db.Client
	crud   abstruct.Crud
}

type DepartmentRepository interface {
    All(context.Context) (*sql.Rows, error)
	Find(context.Context, string) (*sql.Row, error)
	Create(context.Context, string) error
	Update(context.Context, string, string) error
	Destroy(context.Context, string) error
	FindLatestRecord(context.Context) (*sql.Row, error)
}

func NewDepartmentRepository(c db.Client, ac abstract.Crud) DepartmentRepository {
	return &departmentRepository{c, ac}
}

//全件取得
func (b *departmentRepository) All(c context.Context) (*sql.Rows, error) {
	query := "SELECT * FROM departments"
	return b.crud.Read(c, query)
}

// 1件取得
func (b *departmentRepository) Find(c context.Context, id string) (*sql.Row, error){
	query := "SELECT * FROM departments WHERE id =" + id
	return b.crud.ReadByID(c, query)
}

// 作成
func (b *departmentRepository) Create(c context.Context, name string) error {
	query := "INSERT INTO departments (departments) VALUES ('" + name + "')"
	return b.crud.UpdateDB(c, query)
}

// 編集
func (b *departmentRepository) Update(c context.Context, id string, name string) error {
	query := "UPDATE departments SET department = '" + name + "' WHERE id = " + id
	return b.crud.UpdateDB(c, query)
}

// 削除
func (b *departmentRepository) Destroy(c context.Context, id string) error {
	query := "DELETE FROM departments WHERE id =" + id
	return b.crud.UpdateDB(c, query)
}

// 最新のdepartmentを取得する
func (b *departmentRepository) FindLatestRecord(c context.Context) (*sql.Row, error) {
	query := `
		SELECT
			*
		FROM
			departments
		ORDER BY
			id
		DESC LIMIT 1
	`
	return b.crud.ReadByID(c, query)
}
