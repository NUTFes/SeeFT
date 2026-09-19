package usecase

import (
	"context"
	"database/sql/driver"
	"fmt"
	"testing"

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
		{"対応状況そのままで返答が初めて入ったら返答の通知", "inProgress", "", "inProgress", "5分で着きます", entity.RescueNotificationResponse, true},
		{"書いてあった返答の書き直しは編集の通知", "inProgress", "5分で着きます", "inProgress", "10分で着きます", entity.RescueNotificationResponseEdited, true},
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

// Update*Rescueから渡された更新前後の対応状況を記録するだけの通知先
type recordingRescueNotifier struct {
	calls []string
}

func (r *recordingRescueNotifier) TroubleRescueUpdated(_ context.Context, before, after *entity.TroubleRescueForGet) {
	r.calls = append(r.calls, fmt.Sprintf("T%d %s→%s", after.ID, before.Status, after.Status))
}

func (r *recordingRescueNotifier) QuestionRescueUpdated(_ context.Context, before, after *entity.QuestionRescueForGet) {
	r.calls = append(r.calls, fmt.Sprintf("Q%d %s→%s", after.ID, before.Status, after.Status))
}

func (r *recordingRescueNotifier) ShorthandedRescueUpdated(_ context.Context, before, after *entity.ShorthandedRescueForGet) {
	r.calls = append(r.calls, fmt.Sprintf("S%d %s→%s", after.ID, before.Status, after.Status))
}

var troubleRescueColumns = []string{"id", "user_id", "task_id", "place", "detail", "status", "response", "time", "created_at", "updated_at"}

var questionRescueColumns = []string{"id", "user_id", "question", "status", "response", "time", "created_at", "updated_at"}

var shorthandedRescueColumns = []string{"id", "user_id", "task_id", "missing_number", "place", "status", "response", "time", "created_at", "updated_at"}

func troubleRescueRow(status string, response interface{}) []driver.Value {
	return []driver.Value{7, 42, 3, "正門", "看板が倒れた", status, response, fixedTestTime, fixedTestTime, fixedTestTime}
}

func TestUpdateTroubleRescue_更新前後を通知先に渡す(t *testing.T) {
	client, mock := newFakeDBClient(t)
	defer client.CloseDB()
	notifier := &recordingRescueNotifier{}
	uc := NewTroubleRescueUseCase(repository.NewTroubleRescueRepository(client, abstract.NewCrud(client)), notifier)

	mock.ExpectQuery(`FROM trouble_rescues WHERE id`).WithArgs("7").
		WillReturnRows(sqlmock.NewRows(troubleRescueColumns).AddRow(troubleRescueRow("todo", nil)...))
	mock.ExpectExec(`UPDATE trouble_rescues`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(`FROM trouble_rescues WHERE id`).WithArgs("7").
		WillReturnRows(sqlmock.NewRows(troubleRescueColumns).AddRow(troubleRescueRow("inProgress", nil)...))

	got, err := uc.UpdateTroubleRescue(context.Background(), "7", "inProgress", "", true)
	require.NoError(t, err)
	assert.Equal(t, "inProgress", got.Status)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, []string{"T7 todo→inProgress"}, notifier.calls)
}

func TestUpdateTroubleRescue_notifyがfalseなら更新前を読まず知らせない(t *testing.T) {
	client, mock := newFakeDBClient(t)
	defer client.CloseDB()
	notifier := &recordingRescueNotifier{}
	uc := NewTroubleRescueUseCase(repository.NewTroubleRescueRepository(client, abstract.NewCrud(client)), notifier)

	// GASが押し直しの重複を「対応済み・まとめました」にする書き込み
	mock.ExpectExec(`UPDATE trouble_rescues`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(`FROM trouble_rescues WHERE id`).WithArgs("7").
		WillReturnRows(sqlmock.NewRows(troubleRescueColumns).AddRow(troubleRescueRow("done", "対応番号6にまとめました")...))

	_, err := uc.UpdateTroubleRescue(context.Background(), "7", "done", "対応番号6にまとめました", false)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Empty(t, notifier.calls)
}

func TestUpdateTroubleRescue_通知が無効なら更新前を読まない(t *testing.T) {
	client, mock := newFakeDBClient(t)
	defer client.CloseDB()
	uc := NewTroubleRescueUseCase(repository.NewTroubleRescueRepository(client, abstract.NewCrud(client)), nil)

	mock.ExpectExec(`UPDATE trouble_rescues`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(`FROM trouble_rescues WHERE id`).
		WillReturnRows(sqlmock.NewRows(troubleRescueColumns).AddRow(troubleRescueRow("done", nil)...))

	_, err := uc.UpdateTroubleRescue(context.Background(), "7", "done", "", true)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateQuestionRescue_更新前後を通知先に渡す(t *testing.T) {
	client, mock := newFakeDBClient(t)
	defer client.CloseDB()
	notifier := &recordingRescueNotifier{}
	uc := NewQuestionRescueUseCase(repository.NewQuestionRescueRepository(client, abstract.NewCrud(client)), notifier)

	mock.ExpectQuery(`FROM question_rescues WHERE id`).
		WillReturnRows(sqlmock.NewRows(questionRescueColumns).AddRow(4, 42, "雨天時の集合場所は？", "inProgress", nil, fixedTestTime, fixedTestTime, fixedTestTime))
	mock.ExpectExec(`UPDATE question_rescues`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(`FROM question_rescues WHERE id`).
		WillReturnRows(sqlmock.NewRows(questionRescueColumns).AddRow(4, 42, "雨天時の集合場所は？", "done", "体育館です", fixedTestTime, fixedTestTime, fixedTestTime))

	_, err := uc.UpdateQuestionRescue(context.Background(), "4", "done", "体育館です", true)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, []string{"Q4 inProgress→done"}, notifier.calls)
}

func TestUpdateShorthandedRescue_更新前後を通知先に渡す(t *testing.T) {
	client, mock := newFakeDBClient(t)
	defer client.CloseDB()
	notifier := &recordingRescueNotifier{}
	uc := NewShorthandedRescueUseCase(repository.NewShorthandedRescueRepository(client, abstract.NewCrud(client)), notifier)

	mock.ExpectQuery(`FROM shorthanded_rescues WHERE id`).
		WillReturnRows(sqlmock.NewRows(shorthandedRescueColumns).AddRow(5, 42, 3, 2, "受付", "todo", nil, fixedTestTime, fixedTestTime, fixedTestTime))
	mock.ExpectExec(`UPDATE shorthanded_rescues`).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(`FROM shorthanded_rescues WHERE id`).
		WillReturnRows(sqlmock.NewRows(shorthandedRescueColumns).AddRow(5, 42, 3, 2, "受付", "done", nil, fixedTestTime, fixedTestTime, fixedTestTime))

	_, err := uc.UpdateShorthandedRescue(context.Background(), "5", "done", "", true)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, []string{"S5 todo→done"}, notifier.calls)
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

var userColumnNames = []string{"id", "name", "mail", "grade_id", "department_id", "bureau_id", "role_id", "student_number", "tel", "password", "created_at", "updated_at", "slack_user_id"}

func rescueSenderUserRow(slackUserID interface{}) []driver.Value {
	return []driver.Value{42, "送信者", "a@example.com", 1, 1, 1, 1, "000", "000", "x", fixedTestTime, fixedTestTime, slackUserID}
}

// 送信をその場で実行する通知先。本番はgoroutineで裏に回すが、テストでは送った結果をすぐ確かめる
func newSyncRescueNotifier(t *testing.T) (*rescueNotifier, *fakeRescueSender, sqlmock.Sqlmock, func()) {
	t.Helper()
	client, mock := newFakeDBClient(t)
	crud := abstract.NewCrud(client)
	sender := &fakeRescueSender{}
	n := &rescueNotifier{
		sender:  sender,
		taskRep: repository.NewTaskRepository(client, crud),
		userRep: repository.NewUserRepository(client, crud),
		run:     func(f func()) { f() },
	}
	return n, sender, mock, client.CloseDB
}

func TestRescueNotifier_トラブルを確認したら送信者にDMする(t *testing.T) {
	n, sender, mock, closeDB := newSyncRescueNotifier(t)
	defer closeDB()

	mock.ExpectQuery(`FROM tasks WHERE id`).WithArgs("3").
		WillReturnRows(sqlmock.NewRows(taskColumnNames).AddRow(3, "受付", 1, "", "", 1, 3, "", "", 1, fixedTestTime, fixedTestTime))
	mock.ExpectQuery(`FROM users WHERE id`).WithArgs("42").
		WillReturnRows(sqlmock.NewRows(userColumnNames).AddRow(rescueSenderUserRow("U0123")...))

	before := &entity.TroubleRescueForGet{ID: 7, UserID: 42, TaskID: 3, Place: "正門", Detail: "看板が倒れた", Status: "todo", Time: fixedTestTime}
	after := *before
	after.Status = "inProgress"
	after.Response = "5分で向かいます <担当>"
	n.TroubleRescueUpdated(context.Background(), before, &after)
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

func TestRescueNotifier_返答の書き直しは編集として送る(t *testing.T) {
	n, sender, mock, closeDB := newSyncRescueNotifier(t)
	defer closeDB()

	mock.ExpectQuery(`FROM users WHERE id`).WithArgs("42").
		WillReturnRows(sqlmock.NewRows(userColumnNames).AddRow(rescueSenderUserRow("U0123")...))

	before := &entity.QuestionRescueForGet{ID: 4, UserID: 42, Question: "雨天時の集合場所は？", Status: "inProgress", Response: "体育館です", Time: fixedTestTime}
	after := *before
	after.Response = "第2体育館です"
	n.QuestionRescueUpdated(context.Background(), before, &after)
	require.NoError(t, mock.ExpectationsWereMet())

	require.Len(t, sender.sent, 1)
	assert.Equal(t, "✏️ 本部からの返答が編集されました", sender.sent[0].params.Title)
	assert.Equal(t, "Q4", sender.sent[0].params.Number)
	assert.Equal(t, "第2体育館です", sender.sent[0].params.Response)
}

func TestRescueNotifier_知らせる変化でなければ何も読まず送らない(t *testing.T) {
	n, sender, mock, closeDB := newSyncRescueNotifier(t)
	defer closeDB()

	// 関係ない列の編集で、対応状況も返答も変わっていないPUT
	before := &entity.ShorthandedRescueForGet{ID: 5, UserID: 42, TaskID: 3, MissingNumber: 2, Status: "inProgress", Response: "向かいます", Time: fixedTestTime}
	after := *before
	n.ShorthandedRescueUpdated(context.Background(), before, &after)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Empty(t, sender.sent)
}

func TestRescueNotifier_SlackIDが無い人には送らない(t *testing.T) {
	n, sender, mock, closeDB := newSyncRescueNotifier(t)
	defer closeDB()

	mock.ExpectQuery(`FROM tasks WHERE id`).WithArgs("3").
		WillReturnRows(sqlmock.NewRows(taskColumnNames).AddRow(3, "受付", 1, "", "", 1, 3, "", "", 1, fixedTestTime, fixedTestTime))
	mock.ExpectQuery(`FROM users WHERE id`).WithArgs("42").
		WillReturnRows(sqlmock.NewRows(userColumnNames).AddRow(rescueSenderUserRow(nil)...))

	before := &entity.ShorthandedRescueForGet{ID: 5, UserID: 42, TaskID: 3, MissingNumber: 2, Place: "受付", Status: "todo", Time: fixedTestTime}
	after := *before
	after.Status = "done"
	n.ShorthandedRescueUpdated(context.Background(), before, &after)
	require.NoError(t, mock.ExpectationsWereMet())
	assert.Empty(t, sender.sent)
}
