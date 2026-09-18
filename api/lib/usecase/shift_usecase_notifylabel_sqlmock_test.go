package usecase

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newNotifyLabelTestUseCase(t *testing.T) (*shiftUseCase, sqlmock.Sqlmock, func()) {
	t.Helper()

	client, mock := newFakeDBClient(t)
	crud := abstract.NewCrud(client)

	uc := &shiftUseCase{
		rep:           repository.NewShiftRepository(client, crud),
		taskRep:       repository.NewTaskRepository(client, crud),
		userRep:       repository.NewUserRepository(client, crud),
		actionLogRepo: repository.NewActionLogRepository(client),
	}
	return uc, mock, client.CloseDB
}

// action_logに書き込まれるdiff_payloadのold/newを照合する。
// Slack通知はこの文字列をそのまま表示するので、ここが利用者の目に入る文言になる
type taskChangePayload struct {
	old, new string
}

func (p taskChangePayload) Match(v driver.Value) bool {
	b, ok := v.([]byte)
	if !ok {
		return false
	}
	var payload struct {
		Changes []map[string]string `json:"changes"`
	}
	if err := json.Unmarshal(b, &payload); err != nil || len(payload.Changes) != 1 {
		return false
	}
	c := payload.Changes[0]
	return c["field"] == "task_name" && c["old"] == p.old && c["new"] == p.new
}

func TestUpdateShiftsFromGAS_通知の文言で空欄と読み取り失敗を区別する(t *testing.T) {
	const (
		existingShiftID = 100
		dateIDOf1日目     = 2
	)

	cases := []struct {
		name        string
		oldTask     []driver.Value // nilなら旧タスクを読み取れない
		oldTaskID   int
		newTaskID   int
		newTaskName string
		want        taskChangePayload
	}{
		{
			// 2026-09-18、スプシで食事のセルを消しただけの変更が「朝 食事 (1日目) → （不明）」と届き、
			// 何も触っていない受け手から問い合わせになった
			name:        "セルを空欄にした",
			oldTask:     taskValues(350, "朝 食事 (1日目)"),
			oldTaskID:   350,
			newTaskID:   1,
			newTaskName: "",
			want:        taskChangePayload{old: "朝 食事 (1日目)", new: "（割り当てなし）"},
		},
		{
			name:        "空欄のセルに割り当てた",
			oldTask:     taskValues(1, ""),
			oldTaskID:   1,
			newTaskID:   7,
			newTaskName: "受付",
			want:        taskChangePayload{old: "（割り当てなし）", new: "受付"},
		},
		{
			name:        "旧タスクを読み取れない",
			oldTask:     nil,
			oldTaskID:   999,
			newTaskID:   7,
			newTaskName: "受付",
			want:        taskChangePayload{old: "（不明）", new: "受付"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			uc, mock, closeDB := newNotifyLabelTestUseCase(t)
			defer closeDB()

			expectTwoUsers(mock)
			mock.ExpectQuery(`(?s)FROM tasks.*ORDER BY id`).WillReturnRows(
				sqlmock.NewRows(taskColumnNames).AddRow(taskValues(tc.newTaskID, tc.newTaskName)...))

			// 既存シフトがあるのでUPDATEの経路に入る
			mock.ExpectQuery(`SELECT \* FROM shifts`).WillReturnRows(sqlmock.NewRows([]string{
				"id", "task_id", "user_id", "year_id", "date_id", "time_id", "weather_id",
				"is_attendance", "created_at", "updated_at",
			}).AddRow(existingShiftID, tc.oldTaskID, 1, 45, dateIDOf1日目, 31, 1, false, fixedTestTime, fixedTestTime))

			oldTaskRows := sqlmock.NewRows(taskColumnNames)
			if tc.oldTask != nil {
				oldTaskRows.AddRow(tc.oldTask...)
			}
			mock.ExpectQuery(`FROM tasks WHERE id = \$1`).
				WithArgs(strconv.Itoa(tc.oldTaskID)).
				WillReturnRows(oldTaskRows)

			mock.ExpectExec(`INSERT INTO action_logs`).
				WithArgs(existingShiftID, 1, dateIDOf1日目, "UPDATE", tc.want, false).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec(`UPDATE\s+shifts`).WillReturnResult(sqlmock.NewResult(0, 1))

			err := uc.UpdateShiftsFromGAS(context.Background(), entity.ShiftChangeRequest{
				Changes: []entity.ShiftChange{
					{YearID: 45, TimeID: 31, Date: "1日目", Weather: "晴れ", UserName: "山田太郎", TaskName: tc.newTaskName},
				},
			})

			require.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTaskNameForNotification(t *testing.T) {
	assert.Equal(t, "（割り当てなし）", taskNameForNotification(""))
	assert.Equal(t, "受付", taskNameForNotification("受付"))
}
