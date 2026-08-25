package usecase

import (
	"context"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// UpdateShiftsFromGASは日付判定の前にユーザーとタスクを一括取得する。
// 日付判定そのものを検証したいので、その2クエリだけ空の結果で応答させる。
func expectEmptyUserAndTaskLookup(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(`FROM users`).WillReturnRows(sqlmock.NewRows([]string{
		"id", "name", "mail", "grade_id", "department_id", "bureau_id", "role_id",
		"student_number", "tel", "password", "created_at", "updated_at", "slack_user_id",
	}))
	mock.ExpectQuery(`FROM tasks`).WillReturnRows(sqlmock.NewRows([]string{
		"id", "task", "place_id", "url", "manual_url", "bureau_id", "max_member",
		"color", "remark", "year_id", "created_at", "updated_at",
	}))
}

func newGASDateTestUseCase(t *testing.T) (*shiftUseCase, sqlmock.Sqlmock, func()) {
	t.Helper()

	client, mock := newFakeDBClient(t)
	crud := abstract.NewCrud(client)

	uc := &shiftUseCase{
		rep:     repository.NewShiftRepository(client, crud),
		taskRep: repository.NewTaskRepository(client, crud),
		userRep: repository.NewUserRepository(client, crud),
	}
	return uc, mock, client.CloseDB
}

func gasChange(yearID int, date string) entity.ShiftChangeRequest {
	return entity.ShiftChangeRequest{
		Changes: []entity.ShiftChange{
			{
				YearID:   yearID,
				TimeID:   37,
				Date:     date,
				Weather:  "晴れ",
				UserName: "存在しないユーザー",
				TaskName: "テストタスク",
			},
		},
	}
}

// 準々備日はdates.id=5に固定で割り当てるが、shiftsはyear_idとdate_idを個別に
// 外部キー検証するため、45th以外の年度でも保存できてしまう。
// 年度と日程の不整合を防ぐガードが効いていることを確認する。
func TestUpdateShiftsFromGAS_準々備日は45th以外を拒否する(t *testing.T) {
	uc, mock, closeDB := newGASDateTestUseCase(t)
	defer closeDB()

	expectEmptyUserAndTaskLookup(mock)

	err := uc.UpdateShiftsFromGAS(context.Background(), gasChange(43, "準々備日"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "準々備日は45th")
	assert.Contains(t, err.Error(), "yearID=43")
}

// 45thの準々備日は日付判定を通過する。
// ユーザーが見つからない旨のエラーになれば、日付判定より先へ進めたことが分かる。
func TestUpdateShiftsFromGAS_準々備日は45thなら受け付ける(t *testing.T) {
	uc, mock, closeDB := newGASDateTestUseCase(t)
	defer closeDB()

	expectEmptyUserAndTaskLookup(mock)

	err := uc.UpdateShiftsFromGAS(context.Background(), gasChange(45, "準々備日"))

	require.Error(t, err)
	assert.NotContains(t, err.Error(), "準々備日は45th")
	assert.Contains(t, err.Error(), "ユーザーが存在しません")
}

// 既存の4日程は年度によるガードを持たない（45th以外でも従来どおり受け付ける）。
func TestUpdateShiftsFromGAS_既存日程は年度を問わない(t *testing.T) {
	for _, date := range []string{"準備日", "1日目", "2日目", "片付け日"} {
		t.Run(date, func(t *testing.T) {
			uc, mock, closeDB := newGASDateTestUseCase(t)
			defer closeDB()

			expectEmptyUserAndTaskLookup(mock)

			err := uc.UpdateShiftsFromGAS(context.Background(), gasChange(43, date))

			require.Error(t, err)
			assert.Contains(t, err.Error(), "ユーザーが存在しません")
		})
	}
}

// 未知の日付名は従来どおり弾く。
func TestUpdateShiftsFromGAS_不正な日付名を拒否する(t *testing.T) {
	uc, mock, closeDB := newGASDateTestUseCase(t)
	defer closeDB()

	expectEmptyUserAndTaskLookup(mock)

	err := uc.UpdateShiftsFromGAS(context.Background(), gasChange(45, "存在しない日"))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "不正な日付名")
}
