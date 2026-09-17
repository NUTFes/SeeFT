package repository

import (
	"context"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 採番はINSERTのRETURNINGから受け取る。作成後に最新行を読み直す実装だと、
// 同時に作成された別のユーザーの行を掴む。呼び出し元はこのidでセッションを張るため、
// 取り違えるとログイン中の利用者が別人として扱われる。
func TestUserRepositoryCreate_ReturnsIDFromInsert(t *testing.T) {
	client, mock := newDBMock(t)
	repo := NewUserRepository(client, abstract.NewCrud(client))

	mock.ExpectQuery(`(?s)INSERT INTO\s+users.*RETURNING id`).
		WithArgs("新入 太郎", "taro@example.com", "1", "1", "4", "1", "12345678", "09012345678", "hashed", "U0NEW").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(31))

	id, err := repo.Create(context.Background(), "新入 太郎", "taro@example.com", "1", "1", "4", "1", "12345678", "09012345678", "hashed", "U0NEW")

	require.NoError(t, err)
	assert.Equal(t, 31, id)
	require.NoError(t, mock.ExpectationsWereMet())
}
