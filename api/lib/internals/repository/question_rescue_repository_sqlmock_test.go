package repository

import (
	"context"
	"database/sql"
	"strconv"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// テスト用の最小限の db.Client フェイク。Create は DB() だけを使う。
type fakeQuestionRescueDBClient struct {
	sqlDB *sql.DB
}

func (f *fakeQuestionRescueDBClient) DB() *sql.DB      { return f.sqlDB }
func (f *fakeQuestionRescueDBClient) GormDB() *gorm.DB { return nil }
func (f *fakeQuestionRescueDBClient) CloseDB()         { _ = f.sqlDB.Close() }

func newQuestionRescueRepoWithMock(t *testing.T) (QuestionRescueRepository, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)
	t.Cleanup(func() { _ = sqlDB.Close() })

	client := &fakeQuestionRescueDBClient{sqlDB: sqlDB}
	return NewQuestionRescueRepository(client, abstract.NewCrud(client)), mock
}

// 採番はINSERTのRETURNINGから受け取る。作成後に最新行を読み直す実装だと、
// 同時に送信された別のレスキューのidを拾い、スプシの対応番号と中身がずれる（#536）。
func TestQuestionRescueRepositoryCreate_ReturnsIDFromInsert(t *testing.T) {
	repo, mock := newQuestionRescueRepoWithMock(t)

	mock.ExpectQuery(`(?s)INSERT INTO question_rescues.*RETURNING id`).
		WithArgs(42, "テスト質問", "todo", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(7))

	id, err := repo.Create(context.Background(), "42", "テスト質問", "todo")

	require.NoError(t, err)
	assert.Equal(t, 7, id)
	require.NoError(t, mock.ExpectationsWereMet())
}

// user_id が数値でないときはSQLを投げずに弾く。
// 「エラーが返ること」だけを見ると、バリデーションより先にSQLを投げる実装でも
// sqlmockが未宣言のクエリでエラーを返すため通ってしまう。変換エラーそのものを確認する。
func TestQuestionRescueRepositoryCreate_RejectsNonNumericUserID(t *testing.T) {
	repo, mock := newQuestionRescueRepoWithMock(t)

	id, err := repo.Create(context.Background(), "abc", "テスト質問", "todo")

	require.ErrorIs(t, err, strconv.ErrSyntax)
	assert.Equal(t, 0, id)
	require.NoError(t, mock.ExpectationsWereMet())
}
