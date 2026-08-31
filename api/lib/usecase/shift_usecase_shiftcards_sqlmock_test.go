package usecase

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// fakeDBClient は db.Client を満たす最小限のフェイク。shiftCardRepository は
// GormDB() 経由、他のrepositoryはDB() 経由でクエリを発行するため、両方が
// 同一の sqlmock インスタンスを共有するようにする。
type fakeDBClient struct {
	sqlDB  *sql.DB
	gormDB *gorm.DB
}

func (f *fakeDBClient) DB() *sql.DB      { return f.sqlDB }
func (f *fakeDBClient) GormDB() *gorm.DB { return f.gormDB }
func (f *fakeDBClient) CloseDB()         { _ = f.sqlDB.Close() }

func newFakeDBClient(t *testing.T) (*fakeDBClient, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	require.NoError(t, err)

	return &fakeDBClient{sqlDB: sqlDB, gormDB: gormDB}, mock
}

var fixedTestTime = time.Date(2026, 7, 24, 9, 0, 0, 0, time.UTC)

func userRow(timeID, id int, name string) []driver.Value {
	return []driver.Value{timeID, id, name, name + "@example.com", 1, 1, 1, 1, 90000000 + id, "000", fixedTestTime, fixedTestTime, ""}
}

func TestGetUsersByTimes_SingleBatchedQuery(t *testing.T) {
	client, mock := newFakeDBClient(t)
	defer client.CloseDB()

	crud := abstract.NewCrud(client)
	shiftRep := repository.NewShiftRepository(client, crud)
	uc := &shiftUseCase{rep: shiftRep}

	cols := []string{"time_id", "id", "name", "mail", "grade_id", "department_id", "bureau_id", "role_id", "student_number", "tel", "created_at", "updated_at", "slack_user_id"}
	rows := sqlmock.NewRows(cols).
		AddRow(userRow(41, 1, "Alice")...).
		AddRow(userRow(42, 2, "Bob")...).
		AddRow(42, 3, "Carol", "carol@example.com", 2, 1, 2, 1, 90000003, "000", fixedTestTime, fixedTestTime, "")

	mock.ExpectQuery(`s\.time_id = ANY\(\$4\)`).
		WithArgs("4", "43", "2", sqlmock.AnyArg(), "1").
		WillReturnRows(rows)

	result, err := uc.getUsersByTimes(context.Background(), 4, 43, 2, 1, []int{41, 42})
	require.NoError(t, err)

	require.Len(t, result[41], 1)
	assert.Equal(t, "Alice", result[41][0].Name)

	require.Len(t, result[42], 2)
	assert.Equal(t, "Bob", result[42][0].Name)
	assert.Equal(t, "Carol", result[42][1].Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUsersByTimes_EmptyTimeIDsSkipsQuery(t *testing.T) {
	client, mock := newFakeDBClient(t)
	defer client.CloseDB()

	crud := abstract.NewCrud(client)
	shiftRep := repository.NewShiftRepository(client, crud)
	uc := &shiftUseCase{rep: shiftRep}

	result, err := uc.getUsersByTimes(context.Background(), 4, 43, 2, 1, nil)
	require.NoError(t, err)
	assert.Empty(t, result)

	// time_id が空の場合はリポジトリへ問い合わせすら発生しないこと
	assert.NoError(t, mock.ExpectationsWereMet())
}

var shiftCardCols = []string{
	"shift_id", "user_id", "task_id", "year_id", "date_id", "time_id", "weather_id", "is_attendance",
	"task_name", "task_color", "task_url", "manual_url", "task_remark", "max_member", "task_bureau_id",
	"place_id", "place_name", "time_value",
	"user_name", "user_bureau_id", "user_grade_id",
	"year_value", "date_value", "weather_value",
}

func shiftCardDataRow(shiftID, userID, taskID, yearID, dateID, timeID, weatherID int, taskName, timeValue string) []driver.Value {
	return []driver.Value{
		shiftID, userID, taskID, yearID, dateID, timeID, weatherID, true,
		taskName, "#ffffff", "", "", "", 5, 1,
		1, "1号館", timeValue,
		"シフト太郎", 1, 1,
		"2026", "9/6", "晴れ",
	}
}

// expectNoMoreTimes は timeRep.Find(times WHERE id = ...) への問い合わせを「該当なし」として
// n回分まとめて登録する。境界スロット（前後の枠が存在しない）を模擬するための共通ヘルパー。
func expectNoMoreTimes(mock sqlmock.Sqlmock, n int) {
	emptyCols := []string{"id", "time", "created_at", "updated_at"}
	for range n {
		mock.ExpectQuery(`(?i)FROM times WHERE id =`).WillReturnRows(sqlmock.NewRows(emptyCols))
	}
}

func newFullShiftUseCase(client *fakeDBClient) *shiftUseCase {
	crud := abstract.NewCrud(client)
	return &shiftUseCase{
		rep:          repository.NewShiftRepository(client, crud),
		shiftCardRep: repository.NewShiftCardRepository(client),
		timeRep:      repository.NewTimeRepository(client, crud),
		gradeRep:     repository.NewGradeRepository(client, crud),
		bureauRep:    repository.NewBureauRepository(client, crud),
	}
}

// TestGetShiftCardsByUserAndDateAndWeather_SingleCardBatchesUserFetch は、
// 6スロット連続の単一カードに対してメンバー取得クエリが1回のバッチにまとまることを検証する。
// UsersByTimesへの期待値を1回しか登録しないため、スロットごとの個別呼び出しに戻る回帰があれば
// 「unexpected call」で即座に失敗する。
func TestGetShiftCardsByUserAndDateAndWeather_SingleCardBatchesUserFetch(t *testing.T) {
	client, mock := newFakeDBClient(t)
	defer client.CloseDB()
	// endTime計算・前後判定・最終スロットeTime再計算がUsersByTimes呼び出しの前後に混在するため、
	// 呼び出し順序ではなく発生回数を検証する
	mock.MatchExpectationsInOrder(false)
	uc := newFullShiftUseCase(client)

	mock.ExpectQuery(`(?i)FROM grades`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "grade", "created_at", "updated_at"}).
			AddRow(1, "B4", fixedTestTime, fixedTestTime))
	mock.ExpectQuery(`(?i)FROM bureaus`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "bureau", "color", "created_at", "updated_at"}).
			AddRow(1, "総務局", "#000", fixedTestTime, fixedTestTime))

	shiftRows := sqlmock.NewRows(shiftCardCols)
	for i, timeID := range []int{1, 2, 3, 4, 5, 6} {
		shiftRows.AddRow(shiftCardDataRow(100+i, 10, 4, 43, 2, timeID, 1, "テスト1", "8:"+[]string{"00", "15", "30", "45"}[i%4])...)
	}
	mock.ExpectQuery(`(?i)FROM\s+"shifts"`).WillReturnRows(shiftRows)

	// 6スロット中、time_id=1にAlice、time_id=3にBobがアサインされているケース
	mock.ExpectQuery(`s\.time_id = ANY\(\$4\)`).
		WithArgs("4", "43", "2", sqlmock.AnyArg(), "1").
		WillReturnRows(sqlmock.NewRows([]string{"time_id", "id", "name", "mail", "grade_id", "department_id", "bureau_id", "role_id", "student_number", "tel", "created_at", "updated_at", "slack_user_id"}).
			AddRow(userRow(1, 1, "Alice")...).
			AddRow(userRow(3, 2, "Bob")...))

	// 前後スロット(time_id=0, 7)は存在しない前提。time_id=7のFindはEndTime算出・前後判定・
	// 最終スロットのeTime再計算で共有されるため1回だけ呼ばれ、該当なしを返す
	expectNoMoreTimes(mock, 1)

	cards, err := uc.GetShiftCardsByUserAndDateAndWeather(context.Background(), "10", "2", "1")
	require.NoError(t, err)
	require.Len(t, cards, 1)

	card := cards[0]
	assert.Equal(t, "テスト1", card.TaskName)
	require.Len(t, card.ShiftMembers, 6)
	require.Len(t, card.ShiftMembers[0].Members, 1)
	assert.Equal(t, "Alice", card.ShiftMembers[0].Members[0].Name)
	assert.Equal(t, "B4", card.ShiftMembers[0].Members[0].Grade)
	assert.Empty(t, card.ShiftMembers[1].Members)
	require.Len(t, card.ShiftMembers[2].Members, 1)
	assert.Equal(t, "Bob", card.ShiftMembers[2].Members[0].Name)
	assert.Empty(t, card.BeforeMembers.Members)
	assert.Empty(t, card.AfterMembers.Members)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetShiftCardsByUserAndDateAndWeather_TwoTasksBatchPerCard は、2タスク(=2カード)の
// リクエストでメンバー取得クエリがカードごとに1回(=合計2回)にとどまることを検証する。
// taskGroupsはGoのmapで管理されており反復順が非決定的なため、MatchExpectationsInOrder(false)で
// 順序に依存しない形にする。
func TestGetShiftCardsByUserAndDateAndWeather_TwoTasksBatchPerCard(t *testing.T) {
	client, mock := newFakeDBClient(t)
	defer client.CloseDB()
	mock.MatchExpectationsInOrder(false)
	uc := newFullShiftUseCase(client)

	mock.ExpectQuery(`(?i)FROM grades`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "grade", "created_at", "updated_at"}).
			AddRow(1, "B4", fixedTestTime, fixedTestTime))
	mock.ExpectQuery(`(?i)FROM bureaus`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "bureau", "color", "created_at", "updated_at"}).
			AddRow(1, "総務局", "#000", fixedTestTime, fixedTestTime))

	shiftRows := sqlmock.NewRows(shiftCardCols).
		AddRow(shiftCardDataRow(200, 10, 4, 43, 2, 1, 1, "テスト1", "8:00")...).
		AddRow(shiftCardDataRow(201, 10, 4, 43, 2, 2, 1, "テスト1", "8:15")...).
		AddRow(shiftCardDataRow(210, 10, 5, 43, 2, 1, 1, "テスト2", "8:00")...).
		AddRow(shiftCardDataRow(211, 10, 5, 43, 2, 2, 1, "テスト2", "8:15")...)
	mock.ExpectQuery(`(?i)FROM\s+"shifts"`).WillReturnRows(shiftRows)

	emptyUsersCols := []string{"time_id", "id", "name", "mail", "grade_id", "department_id", "bureau_id", "role_id", "student_number", "tel", "created_at", "updated_at", "slack_user_id"}
	// task_id=4のカード用バッチ取得
	mock.ExpectQuery(`s\.time_id = ANY\(\$4\)`).
		WithArgs("4", "43", "2", sqlmock.AnyArg(), "1").
		WillReturnRows(sqlmock.NewRows(emptyUsersCols).AddRow(userRow(1, 1, "Alice")...))
	// task_id=5のカード用バッチ取得
	mock.ExpectQuery(`s\.time_id = ANY\(\$4\)`).
		WithArgs("5", "43", "2", sqlmock.AnyArg(), "1").
		WillReturnRows(sqlmock.NewRows(emptyUsersCols).AddRow(userRow(1, 2, "Bob")...))

	// 両カードとも time_id=1..2 で連続し、前スロット(0)は存在せず、後ろのFind(3)も該当なし。
	// カードごとに1回(endTime/前後判定/最終スロットeTime再計算で共有)、2カード分で合計2回
	expectNoMoreTimes(mock, 2)

	cards, err := uc.GetShiftCardsByUserAndDateAndWeather(context.Background(), "10", "2", "1")
	require.NoError(t, err)
	require.Len(t, cards, 2)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestToShiftMembers_MapsGradeAndBureau(t *testing.T) {
	uc := &shiftUseCase{}

	users := []entity.User{
		{Name: "Alice", GradeID: 1, BureauID: 1},
		{Name: "Bob", GradeID: 2, BureauID: 2},
	}
	gradeMap := map[int]string{1: "B4", 2: "M1"}
	bureauMap := map[int]string{1: "総務局", 2: "企画局"}

	members := uc.toShiftMembers(users, gradeMap, bureauMap)

	require.Len(t, members, 2)
	assert.Equal(t, "Alice", members[0].Name)
	assert.Equal(t, "B4", members[0].Grade)
	assert.Equal(t, "総務局", members[0].Bureau)
	assert.Equal(t, "Bob", members[1].Name)
	assert.Equal(t, "M1", members[1].Grade)
	assert.Equal(t, "企画局", members[1].Bureau)
}

// TestGetShiftCardsByUserAndDateAndWeather_BreakCardSkipsMemberFetch は、休憩カードが
// 担当者取得クエリを一切発行しないことを検証する。UsersByTimesへの期待値を登録しないため、
// 休憩でメンバーを引く回帰が入れば「unexpected call」で即座に失敗する。
//
// 休憩には全スタッフの大半が同じtask_idでぶら下がるため、ここを素通しにすると
// 15分スロットごとに数百人を引くことになる(#488)。
func TestGetShiftCardsByUserAndDateAndWeather_BreakCardSkipsMemberFetch(t *testing.T) {
	client, mock := newFakeDBClient(t)
	defer client.CloseDB()
	mock.MatchExpectationsInOrder(false)
	uc := newFullShiftUseCase(client)

	mock.ExpectQuery(`(?i)FROM grades`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "grade", "created_at", "updated_at"}).
			AddRow(1, "B4", fixedTestTime, fixedTestTime))
	mock.ExpectQuery(`(?i)FROM bureaus`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "bureau", "color", "created_at", "updated_at"}).
			AddRow(1, "総務局", "#000", fixedTestTime, fixedTestTime))

	shiftRows := sqlmock.NewRows(shiftCardCols).
		AddRow(shiftCardDataRow(300, 10, 9, 43, 2, 1, 1, breakTaskName, "12:00")...).
		AddRow(shiftCardDataRow(301, 10, 9, 43, 2, 2, 1, breakTaskName, "12:15")...)
	mock.ExpectQuery(`(?i)FROM\s+"shifts"`).WillReturnRows(shiftRows)

	// EndTime算出のための次スロット(time_id=3)の問い合わせだけは短絡の前に発生する
	mock.ExpectQuery(`(?i)FROM times WHERE id =`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "time", "created_at", "updated_at"}).
			AddRow(3, "12:30", fixedTestTime, fixedTestTime))

	cards, err := uc.GetShiftCardsByUserAndDateAndWeather(context.Background(), "10", "2", "1")
	require.NoError(t, err)
	require.Len(t, cards, 1)

	card := cards[0]
	assert.Equal(t, breakTaskName, card.TaskName)
	assert.Equal(t, "12:00", card.StartTime)
	assert.Equal(t, "12:30", card.EndTime)
	// 「誰が休憩中か」は見せない運用方針のため、担当者は前後を含めて空で返す
	assert.Empty(t, card.ShiftMembers)
	assert.Empty(t, card.BeforeMembers.Members)
	assert.Empty(t, card.AfterMembers.Members)

	assert.NoError(t, mock.ExpectationsWereMet())
}
