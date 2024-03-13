package repository

import (
	"context"
	"database/sql"

	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
)

type gradeRepository struct {
	client db.Client
	crud   abstract.Crud
}

type GradeRepository interface {
	All(context.Context) (*sql.Rows, error)
	Find(context.Context, string) (*sql.Row, error)
	Create(context.Context, string) error
	Update(context.Context, string, string) error
	Destroy(context.Context, string) error
	FindLatestRecord(context.Context) (*sql.Row, error)
}

func NewGradeRepository(c db.Client, ac abstract.Crud) GradeRepository {
	return &gradeRepository{c, ac}
}

// 全件取得
func (b *gradeRepository) All(c context.Context) (*sql.Rows, error) {
	query := "SELECT * FROM grades"
	return b.crud.Read(c, query)
}

// 1件取得
func (b *gradeRepository) Find(c context.Context, id string) (*sql.Row, error) {
	query := "SELECT * FROM gradeus WHERE id =" + id
	return b.crud.ReadByID(c, query)
}

// 作成
func (b *gradeRepository) Create(c context.Context, name string) error {
	query := "INSERT INTO grades (grade) VALUES ('" + name + "')"
	return b.crud.UpdateDB(c, query)
}

// 編集
func (b *gradeRepository) Update(c context.Context, id string, name string) error {
	query := "UPDATE grades SET grade = '" + name + "' WHERE id = " + id
	return b.crud.UpdateDB(c, query)
}

// 削除
func (b *gradeRepository) Destroy(c context.Context, id string) error {
	query := "DELETE FROM grades WHERE id =" + id
	return b.crud.UpdateDB(c, query)
}

// 最新のgradeを取得する
func (b *gradeRepository) FindLatestRecord(c context.Context) (*sql.Row, error) {
	query := `
		SELECT
			*
		FROM
			grades
		ORDER BY
			id
		DESC LIMIT 1
	`
	return b.crud.ReadByID(c, query)
}

