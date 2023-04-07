package repository

import (
	"context"
	"database/sql"

	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
)

type weatherRepository struct {
	client db.Client
	crud   abstract.Crud
}

type WeatherRepository interface {
	All(context.Context) (*sql.Rows, error)
	Find(context.Context, string) (*sql.Row, error)
}

func NewWeatherRepository(c db.Client, ac abstract.Crud) WeatherRepository {
	return &weatherRepository{c, ac}
}

// 全件取得
func (b *weatherRepository) All(c context.Context) (*sql.Rows, error) {
	query := `
		SELECT 
			* 
		FROM 
			weathers`
	return b.crud.Read(c, query)
}

// 1件取得
func (b *weatherRepository) Find(c context.Context, id string) (*sql.Row, error) {
	query := `
		SELECT 
			* 
		FROM 
			weathers 
		WHERE 
			id =` + id
	return b.crud.ReadByID(c, query)
}