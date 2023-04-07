package repository

import (
	"context"
	"database/sql"

	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
)

type yearRepository struct {
	client db.Client
	crud   abstract.Crud
}

type YearRepository interface {
	All(context.Context) (*sql.Rows, error)
	Find(context.Context, string) (*sql.Row, error)
}

func NewYearRepository(c db.Client, ac abstract.Crud) YearRepository {
	return &yearRepository{c, ac}
}

// 全件取得
func (b *yearRepository) All(c context.Context) (*sql.Rows, error) {
	query := `
		SELECT 
			* 
		FROM 
			years`
	return b.crud.Read(c, query)
}

// 1件取得
func (b *yearRepository) Find(c context.Context, id string) (*sql.Row, error) {
	query := `
		SELECT 
			* 
		FROM 
			years 
		WHERE 
			id =` + id
	return b.crud.ReadByID(c, query)
}