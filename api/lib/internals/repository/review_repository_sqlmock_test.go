package repository

import (
	"context"
	"database/sql"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// アポストロフィを含むコメント。文字列連結に戻すとクォートが閉じて SQL が壊れる
const testReviewComment = "マニュアルの don't の説明が分かりやすかった"

type fakeReviewDBClient struct {
	sqlDB *sql.DB
}

func (f *fakeReviewDBClient) DB() *sql.DB      { return f.sqlDB }
func (f *fakeReviewDBClient) GormDB() *gorm.DB { return nil }
func (f *fakeReviewDBClient) CloseDB()         { _ = f.sqlDB.Close() }

func newReviewRepoWithMock(t *testing.T) (ReviewRepository, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)
	t.Cleanup(func() { _ = sqlDB.Close() })

	client := &fakeReviewDBClient{sqlDB: sqlDB}
	return NewReviewRepository(client, abstract.NewCrud(client)), mock
}

// コメントがクエリ本文ではなく引数として渡ることを固定する。
// 引数で渡っている限り、アポストロフィを含んでも SQL は壊れない。
func TestReviewRepositoryCreate_PassesCommentAsArgument(t *testing.T) {
	repo, mock := newReviewRepoWithMock(t)

	mock.ExpectExec(`INSERT INTO reviews \(user_id, task_id, staffing_rating, manual_rating, comment\)`).
		WithArgs("12", "34", "5", "4", testReviewComment).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(context.Background(), "12", "34", "5", "4", testReviewComment)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// WHERE 句の id は最後のプレースホルダ（$6）なので、引数も末尾に来る
func TestReviewRepositoryUpdate_PassesCommentAsArgumentAndIDLast(t *testing.T) {
	repo, mock := newReviewRepoWithMock(t)

	mock.ExpectExec(`UPDATE\s+reviews`).
		WithArgs("12", "34", "5", "4", testReviewComment, "7").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(context.Background(), "7", "12", "34", "5", "4", testReviewComment)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestReviewRepositoryDelete_PassesIDAsArgument(t *testing.T) {
	repo, mock := newReviewRepoWithMock(t)

	mock.ExpectExec(`DELETE FROM reviews WHERE id = \$1`).
		WithArgs("7").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Delete(context.Background(), "7")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
