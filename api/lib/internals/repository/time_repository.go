package repository

import (
	"context"
	"database/sql"

	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
)

type timeRepository struct {
	client db.Client
	crud   abstract.Crud
}

type TimeRepository interface {
	All(context.Context) (*sql.Rows, error)
	Find(context.Context, string) (*sql.Row, error)
	FindByTime(context.Context, string) (*sql.Row, error)
}

func NewTimeRepository(c db.Client, ac abstract.Crud) TimeRepository {
	return &timeRepository{c, ac}
}

// 全件取得
func (b *timeRepository) All(c context.Context) (*sql.Rows, error) {
	query := "SELECT * FROM times"
	return b.crud.Read(c, query)
}

// 1件取得
func (b *timeRepository) Find(c context.Context, id string) (*sql.Row, error) {
	query := "SELECT * FROM times WHERE id =" + id
	return b.crud.ReadByID(c, query)
}

// 時刻から検索
func (b *timeRepository) FindByTime(c context.Context, time string) (*sql.Row, error) {
	query := "SELECT * FROM times WHERE time = '" + time + "'"
	return b.crud.ReadByID(c, query)
}
