package usecase

import (
	"context"
	"database/sql/driver"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
	"github.com/stretchr/testify/require"
)

var gasUserCols = []string{
	"id", "name", "mail", "grade_id", "department_id", "bureau_id", "role_id",
	"student_number", "tel", "password", "created_at", "updated_at", "slack_user_id",
}

func newGASUserTestUseCase(t *testing.T) (*userUseCase, sqlmock.Sqlmock, func()) {
	t.Helper()

	client, mock := newFakeDBClient(t)
	crud := abstract.NewCrud(client)

	uc := &userUseCase{userRep: repository.NewUserRepository(client, crud)}
	return uc, mock, client.CloseDB
}

// 既存ユーザー1行。slackUserID は nil を渡すと NULL になる
func existingUserRow(id int, name, mail string, slackUserID any) []driver.Value {
	return []driver.Value{id, name, mail, 1, 1, 4, 1, 12345678, "09012345678", "hashed", fixedTestTime, fixedTestTime, slackUserID}
}

func gasUserChange(name, mail, slackUserID string) entity.UserChangeRequest {
	return entity.UserChangeRequest{
		Changes: []entity.UserChange{
			{
				Name:          name,
				Bureau:        "企画局",
				Grade:         "B1",
				Department:    "未所属",
				StudentNumber: 12345678,
				Tel:           "09012345678",
				Mail:          mail,
				SlackUserID:   slackUserID,
			},
		},
	}
}

// 名簿送信が運んできた mail と slackUserID が新規作成時に保存されること
func TestUpdateUsersFromGAS_新規ユーザーはmailとslackUserIDを保存する(t *testing.T) {
	uc, mock, closeDB := newGASUserTestUseCase(t)
	defer closeDB()

	mock.ExpectQuery(`FROM users WHERE name`).WithArgs("新入 太郎").
		WillReturnRows(sqlmock.NewRows(gasUserCols))
	mock.ExpectExec(`INSERT INTO\s+users`).
		WithArgs("新入 太郎", "taro@example.com", "1", "1", "4", "1", "12345678", "09012345678", sqlmock.AnyArg(), "U0NEW").
		WillReturnResult(sqlmock.NewResult(10, 1))
	mock.ExpectQuery(`FROM users WHERE name`).WithArgs("新入 太郎").
		WillReturnRows(sqlmock.NewRows(gasUserCols).AddRow(existingUserRow(10, "新入 太郎", "taro@example.com", "U0NEW")...))

	err := uc.UpdateUsersFromGAS(context.Background(), gasUserChange("新入 太郎", "taro@example.com", "U0NEW"))

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

// 既存ユーザーは名簿送信が運んできた mail と slackUserID で上書きされること
func TestUpdateUsersFromGAS_既存ユーザーはmailとslackUserIDを更新する(t *testing.T) {
	uc, mock, closeDB := newGASUserTestUseCase(t)
	defer closeDB()

	mock.ExpectQuery(`FROM users WHERE name`).WithArgs("既存 花子").
		WillReturnRows(sqlmock.NewRows(gasUserCols).AddRow(existingUserRow(20, "既存 花子", "", nil)...))
	mock.ExpectExec(`UPDATE\s+users`).
		WithArgs("既存 花子", "hanako@example.com", "1", "1", "4", "1", "12345678", "09012345678", "hashed", "U0HANAKO", "20").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := uc.UpdateUsersFromGAS(context.Background(), gasUserChange("既存 花子", "hanako@example.com", "U0HANAKO"))

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

// Slack IDの取得が途中の回や SLACK_BOT_TOKEN 未設定の回に名簿を再送しても、
// 既に入っている mail と slack_user_id が空で上書きされないこと
func TestUpdateUsersFromGAS_空で送られたmailとslackUserIDは既存値を保持する(t *testing.T) {
	uc, mock, closeDB := newGASUserTestUseCase(t)
	defer closeDB()

	mock.ExpectQuery(`FROM users WHERE name`).WithArgs("既存 花子").
		WillReturnRows(sqlmock.NewRows(gasUserCols).AddRow(existingUserRow(20, "既存 花子", "old@example.com", "U0OLD")...))
	mock.ExpectExec(`UPDATE\s+users`).
		WithArgs("既存 花子", "old@example.com", "1", "1", "4", "1", "12345678", "09012345678", "hashed", "U0OLD", "20").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := uc.UpdateUsersFromGAS(context.Background(), gasUserChange("既存 花子", "", ""))

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
