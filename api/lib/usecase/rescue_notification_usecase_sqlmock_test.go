package usecase

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/NUTFes/SeeFT/api/lib/externals/slack"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRescueNotificationKind_知らせる変更だけを拾う(t *testing.T) {
	cases := []struct {
		name                                           string
		oldStatus, oldResponse, newStatus, newResponse string
		wantKind                                       string
		wantOK                                         bool
	}{
		{"未対応→対応中は確認の通知", "todo", "", "inProgress", "", entity.RescueNotificationInProgress, true},
		{"対応中→対応済みは完了の通知", "inProgress", "", "done", "向かいます", entity.RescueNotificationDone, true},
		{"未対応→対応済みも完了の通知", "todo", "", "done", "", entity.RescueNotificationDone, true},
		{"対応状況そのままで返答が入ったら返答の通知", "inProgress", "", "inProgress", "5分で着きます", entity.RescueNotificationResponse, true},
		{"返答の書き直しも返答の通知", "inProgress", "5分で着きます", "inProgress", "10分で着きます", entity.RescueNotificationResponse, true},
		{"何も変わっていない(関係ない列の編集)は通知しない", "inProgress", "向かいます", "inProgress", "向かいます", "", false},
		{"返答を消しただけは通知しない", "inProgress", "向かいます", "inProgress", "", "", false},
		{"対応済み→未対応(消し間違い)は通知しない", "done", "", "todo", "", "", false},
		{"対応済み→対応中は通知しない", "done", "", "inProgress", "再対応します", "", false},
		{"知らない状態は通知しない", "todo", "", "unknown", "", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			kind, ok := rescueNotificationKind(tc.oldStatus, tc.oldResponse, tc.newStatus, tc.newResponse)
			assert.Equal(t, tc.wantOK, ok)
			assert.Equal(t, tc.wantKind, kind)
		})
	}
}

var troubleRescueColumns = []string{"id", "user_id", "task_id", "place", "detail", "status", "response", "time", "created_at", "updated_at"}

func troubleRescueRow(status string, response interface{}) []driver.Value {
	return []driver.Value{7, 42, 3, "正門", "看板が倒れた", status, response, fixedTestTime, fixedTestTime, fixedTestTime}
}

func TestUpdateTroubleRescue_差分があるときだけ通知を積む(t *testing.T) {
	cases := []struct {
		name        string
		oldStatus   string
		oldResponse interface{}
		newStatus   string
		newResponse interface{}
		wantKind    string
	}{
		{"確認したら通知を積む", "todo", nil, "inProgress", nil, entity.RescueNotificationInProgress},
		{"関係ない列の編集では積まない", "inProgress", "向かいます", "inProgress", "向かいます", ""},
		{"消し間違いで戻っても積まない", "done", nil, "todo", nil, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client, mock := newFakeDBClient(t)
			defer client.CloseDB()
			uc := NewTroubleRescueUseCase(
				repository.NewTroubleRescueRepository(client, abstract.NewCrud(client)),
				repository.NewRescueNotificationRepository(client),
			)

			mock.ExpectQuery(`FROM trouble_rescues WHERE id`).WithArgs("7").
				WillReturnRows(sqlmock.NewRows(troubleRescueColumns).AddRow(troubleRescueRow(tc.oldStatus, tc.oldResponse)...))
			mock.ExpectExec(`UPDATE trouble_rescues`).WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectQuery(`FROM trouble_rescues WHERE id`).WithArgs("7").
				WillReturnRows(sqlmock.NewRows(troubleRescueColumns).AddRow(troubleRescueRow(tc.newStatus, tc.newResponse)...))
			if tc.wantKind != "" {
				mock.ExpectExec(`INSERT INTO rescue_notifications`).
					WithArgs(entity.RescueTypeTrouble, 7, 42, tc.wantKind, tc.newStatus, tc.newResponse, false).
					WillReturnResult(sqlmock.NewResult(1, 1))
			}

			newResponse, _ := tc.newResponse.(string)
			got, err := uc.UpdateTroubleRescue(context.Background(), "7", tc.newStatus, newResponse)
			require.NoError(t, err)
			assert.Equal(t, tc.newStatus, got.Status)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateTroubleRescue_通知の記録に失敗しても更新は成功させる(t *testing.T) {
	client, mock := newFakeDBClient(t)
	defer client.CloseDB()
	uc := NewTroubleRescueUseCase(
		repository.NewTroubleRescueRepository(client, abstract.NewCrud(client)),
		repository.NewRescueNotificationRepository(client),
	)

	mock.ExpectQuery(`FROM trouble_rescues WHERE id`).
		WillReturnRows(sqlmock.NewRows(troubleRescueColumns).AddRow(troubleRescueRow("todo", nil)...))
	mock.ExpectExec(`UPDATE trouble_rescues`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(`FROM trouble_rescues WHERE id`).
		WillReturnRows(sqlmock.NewRows(troubleRescueColumns).AddRow(troubleRescueRow("done", nil)...))
	mock.ExpectExec(`INSERT INTO rescue_notifications`).WillReturnError(sql.ErrConnDone)

	_, err := uc.UpdateTroubleRescue(context.Background(), "7", "done", "")
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateTroubleRescue_通知が無効なら記録しない(t *testing.T) {
	client, mock := newFakeDBClient(t)
	defer client.CloseDB()
	uc := NewTroubleRescueUseCase(repository.NewTroubleRescueRepository(client, abstract.NewCrud(client)), nil)

	mock.ExpectQuery(`FROM trouble_rescues WHERE id`).
		WillReturnRows(sqlmock.NewRows(troubleRescueColumns).AddRow(troubleRescueRow("todo", nil)...))
	mock.ExpectExec(`UPDATE trouble_rescues`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(`FROM trouble_rescues WHERE id`).
		WillReturnRows(sqlmock.NewRows(troubleRescueColumns).AddRow(troubleRescueRow("done", nil)...))

	_, err := uc.UpdateTroubleRescue(context.Background(), "7", "done", "")
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateQuestionRescue_返答が入ったら通知を積む(t *testing.T) {
	client, mock := newFakeDBClient(t)
	defer client.CloseDB()
	uc := NewQuestionRescueUseCase(
		repository.NewQuestionRescueRepository(client, abstract.NewCrud(client)),
		repository.NewRescueNotificationRepository(client),
	)
	cols := []string{"id", "user_id", "question", "status", "response", "time", "created_at", "updated_at"}

	mock.ExpectQuery(`FROM question_rescues WHERE id`).
		WillReturnRows(sqlmock.NewRows(cols).AddRow(4, 42, "雨天時の集合場所は？", "inProgress", nil, fixedTestTime, fixedTestTime, fixedTestTime))
	mock.ExpectExec(`UPDATE question_rescues`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(`FROM question_rescues WHERE id`).
		WillReturnRows(sqlmock.NewRows(cols).AddRow(4, 42, "雨天時の集合場所は？", "inProgress", "体育館です", fixedTestTime, fixedTestTime, fixedTestTime))
	mock.ExpectExec(`INSERT INTO rescue_notifications`).
		WithArgs(entity.RescueTypeQuestion, 4, 42, entity.RescueNotificationResponse, "inProgress", "体育館です", false).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := uc.UpdateQuestionRescue(context.Background(), "4", "inProgress", "体育館です")
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateShorthandedRescue_対応済みで通知を積む(t *testing.T) {
	client, mock := newFakeDBClient(t)
	defer client.CloseDB()
	uc := NewShorthandedRescueUseCase(
		repository.NewShorthandedRescueRepository(client, abstract.NewCrud(client)),
		repository.NewRescueNotificationRepository(client),
	)
	cols := []string{"id", "user_id", "task_id", "missing_number", "place", "status", "response", "time", "created_at", "updated_at"}

	mock.ExpectQuery(`FROM shorthanded_rescues WHERE id`).
		WillReturnRows(sqlmock.NewRows(cols).AddRow(5, 42, 3, 2, "受付", "todo", nil, fixedTestTime, fixedTestTime, fixedTestTime))
	mock.ExpectExec(`UPDATE shorthanded_rescues`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(`FROM shorthanded_rescues WHERE id`).
		WillReturnRows(sqlmock.NewRows(cols).AddRow(5, 42, 3, 2, "受付", "done", nil, fixedTestTime, fixedTestTime, fixedTestTime))
	mock.ExpectExec(`INSERT INTO rescue_notifications`).
		WithArgs(entity.RescueTypeShorthanded, 5, 42, entity.RescueNotificationDone, "done", nil, false).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err := uc.UpdateShorthandedRescue(context.Background(), "5", "done", "")
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

type fakeRescueSender struct {
	sent []sentRescueMessage
}

type sentRescueMessage struct {
	params      slack.RescueMessageParams
	slackUserID string
}

func (f *fakeRescueSender) SendRescueMessage(p slack.RescueMessageParams, slackUserID string) error {
	f.sent = append(f.sent, sentRescueMessage{p, slackUserID})
	return nil
}

var rescueNotificationColumns = []string{"id", "rescue_type", "rescue_id", "user_id", "kind", "status", "response", "is_sent", "created_at"}

var userColumnNames = []string{"id", "name", "mail", "grade_id", "department_id", "bureau_id", "role_id", "student_number", "tel", "password", "created_at", "updated_at", "slack_user_id"}

func newRescueNotificationTestUseCase(t *testing.T, now time.Time) (*rescueNotificationUseCase, *fakeRescueSender, sqlmock.Sqlmock, func()) {
	t.Helper()
	client, mock := newFakeDBClient(t)
	crud := abstract.NewCrud(client)
	sender := &fakeRescueSender{}
	uc := &rescueNotificationUseCase{
		repo:                     repository.NewRescueNotificationRepository(client),
		sender:                   sender,
		questionRescueUseCase:    NewQuestionRescueUseCase(repository.NewQuestionRescueRepository(client, crud), nil),
		shorthandedRescueUseCase: NewShorthandedRescueUseCase(repository.NewShorthandedRescueRepository(client, crud), nil),
		troubleRescueUseCase:     NewTroubleRescueUseCase(repository.NewTroubleRescueRepository(client, crud), nil),
		taskRep:                  repository.NewTaskRepository(client, crud),
		userRep:                  repository.NewUserRepository(client, crud),
		now:                      func() time.Time { return now },
	}
	return uc, sender, mock, client.CloseDB
}

func TestProcessUnsentRescueNotifications_続けた書き込みを1通にまとめる(t *testing.T) {
	now := fixedTestTime.Add(time.Hour)
	uc, sender, mock, closeDB := newRescueNotificationTestUseCase(t, now)
	defer closeDB()

	// 対応中にしてから返答を書いた。見出しは「確認」、対応状況と返答は最後の状態を載せる
	mock.ExpectQuery(`FROM rescue_notifications`).WillReturnRows(sqlmock.NewRows(rescueNotificationColumns).
		AddRow(1, "trouble", 7, 42, "inProgress", "inProgress", nil, false, now.Add(-40*time.Second)).
		AddRow(2, "trouble", 7, 42, "response", "inProgress", "5分で向かいます <担当>", false, now.Add(-30*time.Second)))
	mock.ExpectQuery(`FROM trouble_rescues WHERE id`).WithArgs("7").
		WillReturnRows(sqlmock.NewRows(troubleRescueColumns).AddRow(troubleRescueRow("inProgress", "5分で向かいます <担当>")...))
	mock.ExpectQuery(`FROM tasks WHERE id`).WithArgs("3").
		WillReturnRows(sqlmock.NewRows(taskColumnNames).AddRow(3, "受付", 1, "", "", 1, 3, "", "", 1, fixedTestTime, fixedTestTime))
	mock.ExpectQuery(`FROM users WHERE id`).WithArgs("42").
		WillReturnRows(sqlmock.NewRows(userColumnNames).AddRow(42, "送信者", "a@example.com", 1, 1, 1, 1, "000", "000", "x", fixedTestTime, fixedTestTime, "U0123"))
	mock.ExpectExec(`UPDATE rescue_notifications`).WithArgs(1, 2).WillReturnResult(sqlmock.NewResult(0, 2))

	require.NoError(t, uc.ProcessUnsentRescueNotifications(context.Background()))
	require.NoError(t, mock.ExpectationsWereMet())

	require.Len(t, sender.sent, 1)
	got := sender.sent[0]
	assert.Equal(t, "U0123", got.slackUserID)
	assert.Equal(t, "👀 本部がレスキューを確認しました", got.params.Title)
	assert.Equal(t, "T7", got.params.Number)
	assert.Equal(t, "トラブル", got.params.TypeLabel)
	assert.Equal(t, "対応中", got.params.Status)
	assert.Equal(t, "5分で向かいます <担当>", got.params.Response)
	assert.Equal(t, "2026/07/24 18:00:00", got.params.Time)
	assert.Equal(t, []string{"発生タスク: 受付", "発生場所: 正門", "内容: 看板が倒れた"}, got.params.Details)
}

func TestProcessUnsentRescueNotifications_書き込み直後は次の回に回す(t *testing.T) {
	now := fixedTestTime.Add(time.Hour)
	uc, sender, mock, closeDB := newRescueNotificationTestUseCase(t, now)
	defer closeDB()

	mock.ExpectQuery(`FROM rescue_notifications`).WillReturnRows(sqlmock.NewRows(rescueNotificationColumns).
		AddRow(1, "trouble", 7, 42, "inProgress", "inProgress", nil, false, now.Add(-5*time.Second)))

	require.NoError(t, uc.ProcessUnsentRescueNotifications(context.Background()))
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Empty(t, sender.sent)
}

func TestProcessUnsentRescueNotifications_SlackIDが無い人は送らず送信済みにする(t *testing.T) {
	now := fixedTestTime.Add(time.Hour)
	uc, sender, mock, closeDB := newRescueNotificationTestUseCase(t, now)
	defer closeDB()

	mock.ExpectQuery(`FROM rescue_notifications`).WillReturnRows(sqlmock.NewRows(rescueNotificationColumns).
		AddRow(3, "question", 4, 42, "done", "done", "体育館です", false, now.Add(-time.Minute)))
	mock.ExpectQuery(`FROM question_rescues WHERE id`).WithArgs("4").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "question", "status", "response", "time", "created_at", "updated_at"}).
			AddRow(4, 42, "雨天時の集合場所は？", "done", "体育館です", fixedTestTime, fixedTestTime, fixedTestTime))
	mock.ExpectQuery(`FROM users WHERE id`).WithArgs("42").
		WillReturnRows(sqlmock.NewRows(userColumnNames).AddRow(42, "送信者", "a@example.com", 1, 1, 1, 1, "000", "000", "x", fixedTestTime, fixedTestTime, nil))
	mock.ExpectExec(`UPDATE rescue_notifications`).WithArgs(3).WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, uc.ProcessUnsentRescueNotifications(context.Background()))
	require.NoError(t, mock.ExpectationsWereMet())

	// 実際のSlackServiceは空のIDを受けると何も送らない(シフト変更通知と同じ扱い)
	require.Len(t, sender.sent, 1)
	assert.Equal(t, "", sender.sent[0].slackUserID)
	assert.Equal(t, "Q4", sender.sent[0].params.Number)
	assert.Equal(t, "✅ レスキューの対応が完了しました", sender.sent[0].params.Title)
}

func TestProcessUnsentRescueNotifications_削除されたレスキューの通知は破棄する(t *testing.T) {
	now := fixedTestTime.Add(time.Hour)
	uc, sender, mock, closeDB := newRescueNotificationTestUseCase(t, now)
	defer closeDB()

	mock.ExpectQuery(`FROM rescue_notifications`).WillReturnRows(sqlmock.NewRows(rescueNotificationColumns).
		AddRow(5, "shorthanded", 9, 42, "done", "done", nil, false, now.Add(-time.Minute)))
	mock.ExpectQuery(`FROM shorthanded_rescues WHERE id`).WithArgs("9").
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectExec(`UPDATE rescue_notifications`).WithArgs(5).WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, uc.ProcessUnsentRescueNotifications(context.Background()))
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Empty(t, sender.sent)
}

func TestProcessUnsentRescueNotifications_DBエラーなら未送信のまま残す(t *testing.T) {
	now := fixedTestTime.Add(time.Hour)
	uc, sender, mock, closeDB := newRescueNotificationTestUseCase(t, now)
	defer closeDB()

	mock.ExpectQuery(`FROM rescue_notifications`).WillReturnRows(sqlmock.NewRows(rescueNotificationColumns).
		AddRow(6, "trouble", 7, 42, "done", "done", nil, false, now.Add(-time.Minute)))
	mock.ExpectQuery(`FROM trouble_rescues WHERE id`).WillReturnError(sql.ErrConnDone)

	require.NoError(t, uc.ProcessUnsentRescueNotifications(context.Background()))
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Empty(t, sender.sent)
}
