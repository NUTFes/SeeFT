package abstract

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/pkg/errors"
)

var debugSQL = os.Getenv("DEBUG_SQL") == "1"

type abstractRepository struct {
	client db.Client
}

type Crud interface {
	Read(context.Context, string, ...interface{}) (*sql.Rows, error)
	ReadByID(context.Context, string, ...interface{}) (*sql.Row, error)
	UpdateDB(context.Context, string, ...interface{}) error
}

func NewCrud(client db.Client) Crud {
	return &abstractRepository{client}
}

func (a abstractRepository) Read(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := a.client.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot connect SQL")
	}
	if debugSQL {
		fmt.Printf("\x1b[36m%s\n", query)
	}
	return rows, nil
}

func (a abstractRepository) ReadByID(ctx context.Context, query string, args ...interface{}) (*sql.Row, error) {
	row := a.client.DB().QueryRowContext(ctx, query, args...)
	if debugSQL {
		fmt.Printf("\x1b[36m%s\n", query)
	}
	return row, nil
}

func (a abstractRepository) UpdateDB(ctx context.Context, query string, args ...interface{}) error {
	_, err := a.client.DB().ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	if debugSQL {
		fmt.Printf("\x1b[36m%s\n", query)
	}
	return err
}

// abstract class Database {
//   Future<List<Map<String, dynamic>>> finds(ctx, String sql);
//   Future<Map<String, dynamic>> find(ctx, String sql);
//   Future<Map<String, dynamic>> insert(ctx, String sql);
//   Future<Map<String, dynamic>> update(ctx, String sql, String getSQL);
// }
