package repository

import (
	"context"
	"database/sql"

	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
)

type dateRepository struct {
	client db.Client
	crud   abstract.Crud
}

type DateRepository interface {
	All(context.Context) (*sql.Rows, error)
	Find(context.Context, string) (*sql.Row, error)
}

func NewDateRepository(c db.Client, ac abstract.Crud) DateRepository {
	return &dateRepository{c, ac}
}

// 全件取得
func (b *dateRepository) All(c context.Context) (*sql.Rows, error) {
	query := `
		SELECT 
			* 
		FROM 
			dates`
	return b.crud.Read(c, query)
}

// 1件取得
func (b *dateRepository) Find(c context.Context, id string) (*sql.Row, error) {
	query := `
		SELECT 
			* 
		FROM 
			dates 
		WHERE 
			id =` + id
	return b.crud.ReadByID(c, query)
}