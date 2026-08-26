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

// ドキュメント版とスライド版はホストが違うので、入れ替わればアサーションで判別できる
const (
	testTaskDocURL   = "https://docs.google.com/document/d/abc/edit"
	testTaskSlideURL = "https://seeft-api.nutfes.net/manuals/enishi"
)

// テスト用の最小限の db.Client フェイク。
// Create / Update は crud 経由で DB() だけを使うため GormDB は呼ばれない。
type fakeTaskDBClient struct {
	sqlDB *sql.DB
}

func (f *fakeTaskDBClient) DB() *sql.DB      { return f.sqlDB }
func (f *fakeTaskDBClient) GormDB() *gorm.DB { return nil }
func (f *fakeTaskDBClient) CloseDB()         { _ = f.sqlDB.Close() }

func newTaskRepoWithMock(t *testing.T) (TaskRepository, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)
	t.Cleanup(func() { _ = sqlDB.Close() })

	client := &fakeTaskDBClient{sqlDB: sqlDB}
	return NewTaskRepository(client, abstract.NewCrud(client)), mock
}

// url と manual_url はどちらも string なので、引数の取り違えをコンパイラは検出できない。
// カラムの並びとプレースホルダに渡す値の並びが一致していることをテストで固定する。
func TestTaskRepositoryCreate_ManualURLFollowsURL(t *testing.T) {
	repo, mock := newTaskRepoWithMock(t)

	mock.ExpectExec(`INSERT INTO tasks \(task, place_id, url, manual_url, bureau_id, max_member, color, remark, year_id\)`).
		WithArgs("縁日運営", "1", testTaskDocURL, testTaskSlideURL, "4", "12", "ffffff", "", "45").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(context.Background(), "縁日運営", "1", testTaskDocURL, testTaskSlideURL, "4", "12", "ffffff", "", "45")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTaskRepositoryUpdateWithManualURL_ManualURLFollowsURL(t *testing.T) {
	repo, mock := newTaskRepoWithMock(t)

	// WHERE 句の id は最後のプレースホルダ（$10）なので、引数も末尾に来る
	mock.ExpectExec(`UPDATE tasks\s+SET task = \$1, place_id = \$2, url = \$3, manual_url = \$4`).
		WithArgs("縁日運営", "1", testTaskDocURL, testTaskSlideURL, "4", "12", "ffffff", "", "45", "7").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UpdateWithManualURL(context.Background(), "7", "縁日運営", "1", testTaskDocURL, testTaskSlideURL, "4", "12", "ffffff", "", "45")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// 管理画面用の Update は manual_url をSET句に含めないこと。
// 読み出した値を書き戻す形だと、GASのタスク送信と競合したときに古い値で上書きする。
func TestTaskRepositoryUpdate_DoesNotTouchManualURL(t *testing.T) {
	repo, mock := newTaskRepoWithMock(t)

	// SET句に manual_url が無い形を固定する（$4 が bureau_id になっている）
	mock.ExpectExec(`UPDATE tasks\s+SET task = \$1, place_id = \$2, url = \$3, bureau_id = \$4,\s+max_member = \$5, color = \$6, remark = \$7, year_id = \$8\s+WHERE id = \$9`).
		WithArgs("縁日運営", "1", testTaskDocURL, "4", "12", "ffffff", "", "45", "7").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(context.Background(), "7", "縁日運営", "1", testTaskDocURL, "4", "12", "ffffff", "", "45")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// スプレッドシート由来の値にシングルクォートが混ざってもクエリが壊れないこと。
// 文字列連結でSQLを組んでいた実装では壊れていた（AGENTS.md「新規SQLはプレースホルダで書く」）。
func TestTaskRepositoryCreate_QuoteInValueIsPassedAsArgument(t *testing.T) {
	repo, mock := newTaskRepoWithMock(t)

	const urlWithQuote = "https://docs.google.com/document/d/a'b/edit"

	mock.ExpectExec(`INSERT INTO tasks`).
		WithArgs("It's a task", "1", urlWithQuote, "", "1", "1", "ffffff", "", "45").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Create(context.Background(), "It's a task", "1", urlWithQuote, "", "1", "1", "ffffff", "", "45")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
