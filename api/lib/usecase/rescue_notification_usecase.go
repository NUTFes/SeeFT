package usecase

import (
	"context"
	"database/sql"
	"fmt"
	logpkg "log"
	"strconv"
	"time"

	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/NUTFes/SeeFT/api/lib/externals/slack"
	rep "github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/pkg/errors"
)

// 最後の書き込みからこの時間が経つまで送らずに待つ。
// 本部が「対応状況」と「返答」を続けて書き換えたとき、DMを2通に分けず最後の状態で1通にまとめるため。
// schedulerの間隔(30秒)と合わせても、書き込みから1分以内に届く
const rescueNotificationSettle = 15 * time.Second

// 対応状況の進み具合。戻る変更(対応済み→未対応など)はスプシの消し間違いとみなして通知しない
var rescueStatusRank = map[string]int{
	"todo":       0,
	"inProgress": 1,
	"done":       2,
}

// 表記はアプリの「本部からの返答」タブ(_statusIcon)に合わせる
var rescueStatusLabel = map[string]string{
	"todo":       "未対応",
	"inProgress": "対応中",
	"done":       "対応済",
}

var rescueNotificationTitle = map[string]string{
	entity.RescueNotificationInProgress: "👀 本部がレスキューを確認しました",
	entity.RescueNotificationDone:       "✅ レスキューの対応が完了しました",
	entity.RescueNotificationResponse:   "💬 本部から返答が届きました",
}

// 同じレスキューの通知をまとめたとき、見出しに使う種別の強さ
var rescueNotificationKindRank = map[string]int{
	entity.RescueNotificationResponse:   0,
	entity.RescueNotificationInProgress: 1,
	entity.RescueNotificationDone:       2,
}

// 更新前後の値から、送信者に知らせるべき変更かを判定する。
// GASのonChangeは値が変わっていない行も含めて編集範囲の全行をPUTしてくるので、
// ここで差分を見ないと関係ない列を編集しただけで通知が飛ぶ
func rescueNotificationKind(oldStatus, oldResponse, newStatus, newResponse string) (string, bool) {
	oldRank, okOld := rescueStatusRank[oldStatus]
	newRank, okNew := rescueStatusRank[newStatus]
	if !okOld || !okNew {
		return "", false
	}
	switch {
	case newRank < oldRank:
		return "", false
	case newRank > oldRank:
		if newStatus == "done" {
			return entity.RescueNotificationDone, true
		}
		return entity.RescueNotificationInProgress, true
	case newResponse != oldResponse && newResponse != "":
		return entity.RescueNotificationResponse, true
	}
	return "", false
}

// レスキューの更新時に呼び、通知すべき変更なら未送信キューに積む。
// 通知の記録に失敗してもPUTは失敗させない(GAS側でエラーのアラートが出て本部の作業が止まるため)
func recordRescueNotification(ctx context.Context, repo rep.RescueNotificationRepository, rescueType string, rescueID, userID int, oldStatus, oldResponse, newStatus, newResponse string) {
	if repo == nil {
		return
	}
	kind, ok := rescueNotificationKind(oldStatus, oldResponse, newStatus, newResponse)
	if !ok {
		return
	}
	err := repo.Create(ctx, entity.RescueNotification{
		RescueType: rescueType,
		RescueID:   rescueID,
		UserID:     userID,
		Kind:       kind,
		Status:     newStatus,
		Response:   newResponse,
	})
	if err != nil {
		logpkg.Printf("rescue_notification記録失敗(%s #%d): %v", rescueType, rescueID, err)
	}
}

// RescueMessageSender レスキュー通知の送り先。本番は *slack.SlackService
type RescueMessageSender interface {
	SendRescueMessage(slack.RescueMessageParams, string) error
}

type rescueNotificationUseCase struct {
	repo                     rep.RescueNotificationRepository
	sender                   RescueMessageSender
	questionRescueUseCase    QuestionRescueUseCase
	shorthandedRescueUseCase ShorthandedRescueUseCase
	troubleRescueUseCase     TroubleRescueUseCase
	taskRep                  rep.TaskRepository
	userRep                  rep.UserRepository
	now                      func() time.Time
}

type RescueNotificationUseCase interface {
	ProcessUnsentRescueNotifications(ctx context.Context) error
}

func NewRescueNotificationUseCase(
	repo rep.RescueNotificationRepository,
	sender RescueMessageSender,
	questionRescueUseCase QuestionRescueUseCase,
	shorthandedRescueUseCase ShorthandedRescueUseCase,
	troubleRescueUseCase TroubleRescueUseCase,
	taskRep rep.TaskRepository,
	userRep rep.UserRepository,
) RescueNotificationUseCase {
	return &rescueNotificationUseCase{
		repo:                     repo,
		sender:                   sender,
		questionRescueUseCase:    questionRescueUseCase,
		shorthandedRescueUseCase: shorthandedRescueUseCase,
		troubleRescueUseCase:     troubleRescueUseCase,
		taskRep:                  taskRep,
		userRep:                  userRep,
		now:                      time.Now,
	}
}

// ProcessUnsentRescueNotifications 未送信の通知をレスキューごとに1通にまとめて送信者へDMする
func (u *rescueNotificationUseCase) ProcessUnsentRescueNotifications(ctx context.Context) error {
	notifications, err := u.loadUnsent(ctx)
	if err != nil {
		return err
	}

	now := u.now()
	for _, group := range groupRescueNotifications(notifications) {
		latest := group[len(group)-1]
		// まだ書き換えが続いているかもしれないので次の回に回す
		if now.Sub(latest.CreatedAt) < rescueNotificationSettle {
			continue
		}

		if err := u.sendGroup(ctx, group); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				// 未送信のまま残し、次の回に再送する
				logpkg.Printf("レスキュー通知の送信失敗(%s #%d): %v", latest.RescueType, latest.RescueID, err)
				continue
			}
			// レスキュー自体が削除されていて送れない。残すと毎回失敗し続けるので送信済みにする
			logpkg.Printf("レスキューが見つからないため通知を破棄(%s #%d)", latest.RescueType, latest.RescueID)
		}

		ids := make([]int, len(group))
		for i, n := range group {
			ids[i] = n.ID
		}
		if err := u.repo.MarkAsSent(ctx, ids); err != nil {
			return errors.Wrapf(err, "failed to mark rescue notifications as sent")
		}
	}
	return nil
}

func (u *rescueNotificationUseCase) loadUnsent(ctx context.Context) ([]entity.RescueNotification, error) {
	rows, err := u.repo.FindUnsent(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var notifications []entity.RescueNotification
	for rows.Next() {
		var n entity.RescueNotification
		var userID sql.NullInt64
		var response sql.NullString
		if err := rows.Scan(&n.ID, &n.RescueType, &n.RescueID, &userID, &n.Kind, &n.Status, &response, &n.IsSent, &n.CreatedAt); err != nil {
			return nil, errors.Wrapf(err, "failed to scan rescue notification")
		}
		if userID.Valid {
			n.UserID = int(userID.Int64)
		}
		n.Response = response.String
		notifications = append(notifications, n)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrapf(err, "failed to read rescue notifications")
	}
	return notifications, nil
}

// 同じレスキュー(種類+ID)の通知をまとめる。各グループ内も、グループの並びも古い順を保つ
func groupRescueNotifications(notifications []entity.RescueNotification) [][]entity.RescueNotification {
	var groups [][]entity.RescueNotification
	index := make(map[string]int)
	for _, n := range notifications {
		key := fmt.Sprintf("%s_%d", n.RescueType, n.RescueID)
		i, ok := index[key]
		if !ok {
			i = len(groups)
			index[key] = i
			groups = append(groups, nil)
		}
		groups[i] = append(groups[i], n)
	}
	return groups
}

// まとめた通知から1通のDMを組み立てて送る
func (u *rescueNotificationUseCase) sendGroup(ctx context.Context, group []entity.RescueNotification) error {
	// 見出しはまとめた中で一番強い変化にし、対応状況と返答は最後の状態を載せる
	kind := group[0].Kind
	for _, n := range group[1:] {
		if rescueNotificationKindRank[n.Kind] > rescueNotificationKindRank[kind] {
			kind = n.Kind
		}
	}
	latest := group[len(group)-1]

	params, userID, err := u.buildRescueMessage(ctx, latest.RescueType, latest.RescueID)
	if err != nil {
		return err
	}
	params.Title = rescueNotificationTitle[kind]
	params.Status = rescueStatusLabel[latest.Status]
	params.Response = latest.Response

	slackUserID, err := findSlackUserID(ctx, u.userRep, userID)
	if err != nil {
		return err
	}
	// Slack IDが無い人にはシフト変更通知と同じく何も送らない(SendMessage側で黙ってスキップされる)
	return u.sender.SendRescueMessage(params, slackUserID)
}

// レスキューの送信内容を読み、メッセージの番号・種類・送信時刻・内容を埋める。
// 送信時刻の整形とタスク名の補完は「本部からの返答」タブ(GET /rescues)と同じものを使う
func (u *rescueNotificationUseCase) buildRescueMessage(ctx context.Context, rescueType string, rescueID int) (slack.RescueMessageParams, int, error) {
	id := strconv.Itoa(rescueID)
	switch rescueType {
	case entity.RescueTypeTrouble:
		r, err := u.troubleRescueUseCase.GetTroubleRescueByID(ctx, id)
		if err != nil {
			return slack.RescueMessageParams{}, 0, err
		}
		taskName, err := findTaskName(ctx, u.taskRep, strconv.Itoa(r.TaskID))
		if err != nil {
			taskName = "タスク外"
		}
		res := entity.NewTroubleRescueResponse(r, "", taskName)
		return slack.RescueMessageParams{
			Number:    fmt.Sprintf("T%d", r.ID),
			TypeLabel: "トラブル",
			Time:      res.Time,
			Details:   nonEmptyLines("発生タスク", taskName, "発生場所", r.Place, "内容", r.Detail),
		}, r.UserID, nil
	case entity.RescueTypeQuestion:
		r, err := u.questionRescueUseCase.GetQuestionRescueByID(ctx, id)
		if err != nil {
			return slack.RescueMessageParams{}, 0, err
		}
		res := entity.NewQuestionRescueResponse(r, "")
		return slack.RescueMessageParams{
			Number:    fmt.Sprintf("Q%d", r.ID),
			TypeLabel: "質問",
			Time:      res.Time,
			Details:   nonEmptyLines("質問", r.Question),
		}, r.UserID, nil
	case entity.RescueTypeShorthanded:
		r, err := u.shorthandedRescueUseCase.GetShorthandedRescueByID(ctx, id)
		if err != nil {
			return slack.RescueMessageParams{}, 0, err
		}
		taskName, err := findTaskName(ctx, u.taskRep, strconv.Itoa(r.TaskID))
		if err != nil {
			taskName = "不明なタスク"
		}
		res := entity.NewShorthandedRescueResponse(r, "", taskName)
		return slack.RescueMessageParams{
			Number:    fmt.Sprintf("S%d", r.ID),
			TypeLabel: "人が来ない",
			Time:      res.Time,
			Details:   nonEmptyLines("発生タスク", taskName, "送り先の場所", r.Place, "足りない人数", strconv.Itoa(r.MissingNumber)+"人"),
		}, r.UserID, nil
	}
	return slack.RescueMessageParams{}, 0, errors.Errorf("unknown rescue type: %s", rescueType)
}

// 「見出し, 値」の組から「見出し: 値」の行を作る。値が空の行は載せない
func nonEmptyLines(pairs ...string) []string {
	lines := make([]string, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		if pairs[i+1] == "" {
			continue
		}
		lines = append(lines, pairs[i]+": "+pairs[i+1])
	}
	return lines
}

// ユーザーのSlack IDを取得する。ユーザーが消えている・Slack IDが未登録なら空文字
func findSlackUserID(ctx context.Context, userRep rep.UserRepository, userID int) (string, error) {
	if userID == 0 {
		return "", nil
	}
	row, err := userRep.Find(ctx, strconv.Itoa(userID))
	if err != nil {
		return "", errors.Wrapf(err, "failed to find user")
	}
	var user entity.User
	var slackUserID sql.NullString
	err = row.Scan(
		&user.ID, &user.Name, &user.Mail, &user.GradeID, &user.DepartmentID,
		&user.BureauID, &user.RoleID, &user.StudentNumber, &user.Tel,
		&user.Password, &user.CreatedAt, &user.UpdatedAt, &slackUserID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", errors.Wrapf(err, "failed to scan user")
	}
	return slackUserID.String, nil
}
