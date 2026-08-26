package usecase

import (
	"context"
	"database/sql/driver"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
	"github.com/stretchr/testify/assert"
)

// ドキュメント版とスライド版はホストが違うので、入れ替わればアサーションで判別できる
const (
	testManualDocURL   = "https://docs.google.com/document/d/abc/edit"
	testManualSlideURL = "https://seeft-api.nutfes.net/manuals/enishi"
)

var taskColumnNames = []string{
	"id", "task", "place_id", "url", "manual_url",
	"bureau_id", "max_member", "color", "remark", "year_id",
	"created_at", "updated_at",
}

func taskRow(name string, url string, manualURL string) []driver.Value {
	return []driver.Value{7, name, 1, url, manualURL, 4, 12, "ffffff", "", 45, fixedTestTime, fixedTestTime}
}

func newTaskUseCaseWithMock(t *testing.T) (TaskUseCase, sqlmock.Sqlmock, *fakeDBClient) {
	t.Helper()

	client, mock := newFakeDBClient(t)
	crud := abstract.NewCrud(client)
	u := NewTaskUseCase(
		repository.NewTaskRepository(client, crud),
		repository.NewPlaceRepository(client, crud),
	)
	return u, mock, client
}

// GASのタスク送信がスプレッドシートのスライド版URLをそのまま manual_url へ渡すこと。
// 受け口の entity.TaskAndPlaceChange に ManualUrl が無いと黙って捨てられ、
// エラーにならないまま導線が出ない状態になるため、経路をテストで固定する。
func TestUpdateTasksAndPlacesFromGAS_WritesManualURL(t *testing.T) {
	u, mock, client := newTaskUseCaseWithMock(t)
	defer client.CloseDB()

	// 集合場所は既存扱い（FindByName は文字列連結のクエリなので引数を取らない）
	mock.ExpectQuery(`SELECT \* FROM places WHERE place = '本部\(電気棟1F\)'`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "place", "remark", "created_at", "updated_at"}).
			AddRow(1, "本部(電気棟1F)", "", fixedTestTime, fixedTestTime))
	mock.ExpectExec(`UPDATE places SET`).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// タスクも既存扱い
	mock.ExpectQuery(`SELECT .* FROM tasks WHERE task = \$1`).
		WithArgs("縁日運営").
		WillReturnRows(sqlmock.NewRows(taskColumnNames).
			AddRow(taskRow("縁日運営", "", "")...))

	// url と manual_url がそれぞれ正しい位置に渡ること
	mock.ExpectExec(`UPDATE tasks\s+SET task = \$1, place_id = \$2, url = \$3, manual_url = \$4`).
		WithArgs("縁日運営", "1", testManualDocURL, testManualSlideURL, "4", "12", "ffffff", "", "45", "7").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := u.UpdateTasksAndPlacesFromGAS(context.Background(), entity.TaskAndPlaceChangeRequest{
		Changes: []entity.TaskAndPlaceChange{{
			YearID:    45,
			TaskName:  "縁日運営",
			Bureau:    "企画局",
			Place:     "本部(電気棟1F)",
			Url:       testManualDocURL,
			ManualUrl: testManualSlideURL,
			MaxMember: 12,
		}},
	})

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// 管理画面からの編集は manual_url を送らない。空で上書きすると、タスク名を直しただけで
// シフトカードのスライド版の行が消えてしまう（if (_hasSlide) の条件付き描画）。
// 更新前に読み出した既存値を引き継ぐことを固定する。
func TestUpdateTask_PreservesExistingManualURL(t *testing.T) {
	u, mock, client := newTaskUseCaseWithMock(t)
	defer client.CloseDB()

	// 更新前の読み出し。既存のスライド版URLを持っている
	mock.ExpectQuery(`SELECT .* FROM tasks WHERE id = \$1`).
		WithArgs("7").
		WillReturnRows(sqlmock.NewRows(taskColumnNames).
			AddRow(taskRow("縁日運営", testManualDocURL, testManualSlideURL)...))

	// manual_url には空文字ではなく既存値が渡ること
	mock.ExpectExec(`UPDATE tasks\s+SET task = \$1, place_id = \$2, url = \$3, manual_url = \$4`).
		WithArgs("縁日運営（改）", "1", testManualDocURL, testManualSlideURL, "4", "12", "ffffff", "", "45", "7").
		WillReturnResult(sqlmock.NewResult(0, 1))

	// 更新後の再読み出し
	mock.ExpectQuery(`SELECT .* FROM tasks WHERE id = \$1`).
		WithArgs("7").
		WillReturnRows(sqlmock.NewRows(taskColumnNames).
			AddRow(taskRow("縁日運営（改）", testManualDocURL, testManualSlideURL)...))

	updated, err := u.UpdateTask(context.Background(), "7", "縁日運営（改）", "1", testManualDocURL, "4", "12", "ffffff", "", "45")

	assert.NoError(t, err)
	assert.Equal(t, testManualSlideURL, updated.ManualUrl)
	assert.NoError(t, mock.ExpectationsWereMet())
}
