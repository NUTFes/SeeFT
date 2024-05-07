package repository

import (
	"context"
	"database/sql"

	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
)

type placeRepository struct {
	client db.Client
	crud   abstract.Crud
}

type PlaceRepository interface {
	All(context.Context) (*sql.Rows, error)
	Find(context.Context, string) (*sql.Row, error)
	Create(context.Context, string, string) error
	Update(context.Context, string, string, string) error
	Destroy(context.Context, string) error
	FindNewRecord(context.Context) (*sql.Row, error)
}

func NewPlaceRepository(c db.Client, ac abstract.Crud) PlaceRepository {
	return &placeRepository{c, ac}
}

// 全件取得
func (b *placeRepository) All(c context.Context) (*sql.Rows, error) {
	query := "SELECT * FROM places"
	return b.crud.Read(c, query)
}

// 1件取得
func (b *placeRepository) Find(c context.Context, id string) (*sql.Row, error) {
	query := "SELECT * FROM places WHERE id =" + id
	return b.crud.ReadByID(c, query)
}

// 作成
func (b *placeRepository) Create(c context.Context, name string, remark string) error {
	query := "INSERT INTO places (place, remark) VALUES ('" + name + "', '" + remark +"')"
	return b.crud.UpdateDB(c, query)
}

// 編集
func (b *placeRepository) Update(c context.Context, id string, name string, remark string) error {
	query := "UPDATE places SET (place, remark) = ('" + name + "', '" + remark +"') WHERE id = " + id
	return b.crud.UpdateDB(c, query)
}

// 削除
func (b *placeRepository) Destroy(c context.Context, id string) error {
	query := "DELETE FROM places WHERE id =" + id
	return b.crud.UpdateDB(c, query)
}

// 最新のgradeを取得する
func (b *placeRepository) FindNewRecord(c context.Context) (*sql.Row, error) {
	query := `
		SELECT
			*
		FROM
			places
		ORDER BY
			id
		DESC LIMIT 1
	`
	return b.crud.ReadByID(c, query)
}

