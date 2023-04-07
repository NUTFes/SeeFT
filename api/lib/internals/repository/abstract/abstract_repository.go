package abstract

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/pkg/errors"
)

type abstractRepository struct {
	client db.Client
}

type Crud interface {
	Read(context.Context, string) (*sql.Rows, error)
	ReadByID(context.Context, string) (*sql.Row, error)
	UpdateDB(context.Context, string) error
}

func NewCrud(client db.Client) Crud {
	return &abstractRepository{client}
}

func (a abstractRepository) Read(ctx context.Context, query string) (*sql.Rows, error) {
	rows, err := a.client.DB().QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot connect SQL")
	}
	fmt.Printf("\x1b[36m%s\n", query)
	return rows, nil
}

func (a abstractRepository) ReadByID(ctx context.Context, query string) (*sql.Row, error) {
	row := a.client.DB().QueryRowContext(ctx, query)
	fmt.Printf("\x1b[36m%s\n", query)
	return row, nil
}

func (a abstractRepository) UpdateDB(ctx context.Context, query string) error {
	_, err := a.client.DB().ExecContext(ctx, query)
	if err != nil {
		return err
	}
	fmt.Printf("\x1b[36m%s\n", query)
	return err
}

// abstract class Database {
//   Future<List<Map<String, dynamic>>> finds(ctx, String sql);
//   Future<Map<String, dynamic>> find(ctx, String sql);
//   Future<Map<String, dynamic>> insert(ctx, String sql);
//   Future<Map<String, dynamic>> update(ctx, String sql, String getSQL);
// }
