package repository

import (
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// テスト用の最小限の db.Client フェイク。リポジトリのCreate/Update系は
// crud 経由でも DB() しか使わないため GormDB は呼ばれない。
// 既存の fakeTaskDBClient / fakeReviewDBClient / fakeQuestionRescueDBClient と
// 同じものをリポジトリ非依存にした共通版で、新しいテストはこちらを使う。
type fakeDBClient struct {
	sqlDB *sql.DB
}

func (f *fakeDBClient) DB() *sql.DB      { return f.sqlDB }
func (f *fakeDBClient) GormDB() *gorm.DB { return nil }
func (f *fakeDBClient) CloseDB()         { _ = f.sqlDB.Close() }

// クエリの一致は正規表現で見る（プレースホルダや改行の差で落ちないため）
func newDBMock(t *testing.T) (db.Client, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)
	t.Cleanup(func() { _ = sqlDB.Close() })

	return &fakeDBClient{sqlDB: sqlDB}, mock
}
