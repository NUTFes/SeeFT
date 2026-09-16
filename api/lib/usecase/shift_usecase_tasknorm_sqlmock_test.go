package usecase

import (
	"context"
	"database/sql/driver"
	"strings"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 45thの日程シートに実在した表記。GASは正規化せずそのまま送ってくる
const (
	zenkakuTaskName = "Music Nuts　参加"
	hankakuTaskName = "Music Nuts 参加"
)

func newTaskNormTestUseCase(t *testing.T) (*shiftUseCase, sqlmock.Sqlmock, func()) {
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

// 同じタスク名のシフトを2件送る。1件目で作られたタスクを2件目が再利用できるかを見る
func twoChangesWithSameTask(taskName string) entity.ShiftChangeRequest {
	return entity.ShiftChangeRequest{
		Changes: []entity.ShiftChange{
			{YearID: 45, TimeID: 37, Date: "1日目", Weather: "晴れ", UserName: "山田太郎", TaskName: taskName},
			{YearID: 45, TimeID: 38, Date: "1日目", Weather: "晴れ", UserName: "鈴木花子", TaskName: taskName},
		},
	}
}

// taskColumnNames は task_usecase_manualurl_sqlmock_test.go の定義を共用する

func taskValues(id int, name string) []driver.Value {
	return []driver.Value{id, name, 1, "", "", 1, 1, "000000", "", 45, fixedTestTime, fixedTestTime}
}

func expectTwoUsers(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(`FROM users`).WillReturnRows(sqlmock.NewRows([]string{
		"id", "name", "mail", "grade_id", "department_id", "bureau_id", "role_id",
		"student_number", "tel", "password", "created_at", "updated_at", "slack_user_id",
	}).
		AddRow(1, "山田太郎", "a@example.com", 1, 1, 1, 1, 90000001, "000", "pw", fixedTestTime, fixedTestTime, "").
		AddRow(2, "鈴木花子", "b@example.com", 1, 1, 1, 1, 90000002, "000", "pw", fixedTestTime, fixedTestTime, ""))
}

// 既存シフトは無いものとして新規作成の経路を通す。actionLogRepoはnilなので記録はされない
func expectShiftInsert(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(`SELECT \* FROM shifts`).WillReturnRows(sqlmock.NewRows([]string{
		"id", "task_id", "user_id", "year_id", "date_id", "time_id", "weather_id",
		"is_attendance", "created_at", "updated_at",
	}))
	mock.ExpectQuery(`INSERT INTO shifts`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
}

// 45th本番では、全角スペース入りのタスク名でシフトを送るたびにタスクが1件ずつ作られ、
// 36分間で2504行まで増えた。照合に使うキー(半角に正規化)と、保存・マップ登録に使う
// キー(全角のまま)がずれており、マップが永久にヒットしなかったため。
// 同じ名前が2回出てきても作成は1回で済むことを確認する。
func TestUpdateShiftsFromGAS_全角スペースのタスク名でも作成は1回だけ(t *testing.T) {
	uc, mock, closeDB := newTaskNormTestUseCase(t)
	defer closeDB()

	expectTwoUsers(mock)
	// DBにこのタスクはまだ無い
	mock.ExpectQuery(`(?s)FROM tasks.*ORDER BY id`).WillReturnRows(sqlmock.NewRows(taskColumnNames))

	// 1件目でだけタスクが作られる。保存される名前は半角に正規化されている
	mock.ExpectExec(`INSERT INTO tasks`).
		WithArgs(hankakuTaskName, "1", "", "", "1", "1", "000000", "", "45").
		WillReturnResult(sqlmock.NewResult(10, 1))
	mock.ExpectQuery(`FROM tasks WHERE task = \$1`).
		WithArgs(hankakuTaskName).
		WillReturnRows(sqlmock.NewRows(taskColumnNames).AddRow(taskValues(10, hankakuTaskName)...))
	expectShiftInsert(mock)

	// 2件目はマップに当たるので INSERT INTO tasks は走らない。
	// 走った場合は未定義の呼び出しとして ExpectationsWereMet が落ちる
	expectShiftInsert(mock)

	err := uc.UpdateShiftsFromGAS(context.Background(), twoChangesWithSameTask(zenkakuTaskName))

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// 投入順序が正しく「タスク送信 → シフト送信」の場合。タスク送信側は半角で保存するため、
// 全角の名前が届いても既存行に当たり、新規作成は起きない。
func TestUpdateShiftsFromGAS_既存の半角タスクに全角の名前が届いても作成しない(t *testing.T) {
	uc, mock, closeDB := newTaskNormTestUseCase(t)
	defer closeDB()

	expectTwoUsers(mock)
	// タスク送信が先に走り、半角版が既にある状態
	mock.ExpectQuery(`(?s)FROM tasks.*ORDER BY id`).WillReturnRows(
		sqlmock.NewRows(taskColumnNames).AddRow(taskValues(10, hankakuTaskName)...))

	expectShiftInsert(mock)
	expectShiftInsert(mock)

	err := uc.UpdateShiftsFromGAS(context.Background(), twoChangesWithSameTask(zenkakuTaskName))

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// tasks.task に一意制約が無いため、過去に作られた重複行が残っている可能性がある。
// FindByNames が id 昇順で返し、呼び出し側が先に読んだ行を採るので、
// どのタスクに紐づくかが実行ごとに変わらない。
func TestUpdateShiftsFromGAS_同名タスクが複数あるとき最も古い行に紐づく(t *testing.T) {
	uc, mock, closeDB := newTaskNormTestUseCase(t)
	defer closeDB()

	expectTwoUsers(mock)
	mock.ExpectQuery(`(?s)FROM tasks.*ORDER BY id`).WillReturnRows(
		sqlmock.NewRows(taskColumnNames).
			AddRow(taskValues(10, hankakuTaskName)...).
			AddRow(taskValues(999, hankakuTaskName)...))

	// task_id は後から読んだ999ではなく10になる
	for i := 0; i < 2; i++ {
		mock.ExpectQuery(`SELECT \* FROM shifts`).WillReturnRows(sqlmock.NewRows([]string{
			"id", "task_id", "user_id", "year_id", "date_id", "time_id", "weather_id",
			"is_attendance", "created_at", "updated_at",
		}))
		mock.ExpectQuery(`INSERT INTO shifts .* VALUES \(10,`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	}

	err := uc.UpdateShiftsFromGAS(context.Background(), twoChangesWithSameTask(zenkakuTaskName))

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// FindByNames に渡す名前の配列（pq.Array）に、指定した名前が含まれることを見る。
// taskNameSet は map なので順序が決まらず、配列そのものとは比較できない
type arrayContaining struct {
	want string
}

func (a arrayContaining) Match(v driver.Value) bool {
	s, ok := v.(string)
	return ok && strings.Contains(s, a.want)
}

// adminからの手動登録などでDBに全角スペース入りの行しか無い場合。
// 正規化後の名前だけで引くと既存行を取得できず、半角版を作って並存させてしまう。
// 照合には正規化前の名前も含める必要がある。
func TestUpdateShiftsFromGAS_DBに全角スペースの行しか無くても再利用する(t *testing.T) {
	uc, mock, closeDB := newTaskNormTestUseCase(t)
	defer closeDB()

	expectTwoUsers(mock)
	// 全角のままの名前も検索対象に含まれていなければ、この行は取得できない
	mock.ExpectQuery(`(?s)FROM tasks.*ORDER BY id`).
		WithArgs(arrayContaining{want: zenkakuTaskName}).
		WillReturnRows(sqlmock.NewRows(taskColumnNames).AddRow(taskValues(10, zenkakuTaskName)...))

	// 既存行を再利用するので INSERT INTO tasks は走らず、task_id は10になる
	for i := 0; i < 2; i++ {
		mock.ExpectQuery(`SELECT \* FROM shifts`).WillReturnRows(sqlmock.NewRows([]string{
			"id", "task_id", "user_id", "year_id", "date_id", "time_id", "weather_id",
			"is_attendance", "created_at", "updated_at",
		}))
		mock.ExpectQuery(`INSERT INTO shifts .* VALUES \(10,`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	}

	err := uc.UpdateShiftsFromGAS(context.Background(), twoChangesWithSameTask(zenkakuTaskName))

	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
