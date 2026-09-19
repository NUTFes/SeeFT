package usecase

import (
	"context"
	"database/sql"
	"fmt"
	logpkg "log"
	"strconv"

	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/NUTFes/SeeFT/api/lib/externals/slack"
	rep "github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/pkg/errors"
)

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
	entity.RescueNotificationInProgress:     "👀 本部がレスキューを確認しました",
	entity.RescueNotificationDone:           "✅ レスキューの対応が完了しました",
	entity.RescueNotificationResponse:       "💬 本部から返答が届きました",
	entity.RescueNotificationResponseEdited: "✏️ 本部からの返答が編集されました",
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
		// 返答欄が空だったところに書かれたら「届きました」、書いてあった返答の書き換えなら「編集されました」。
		// 一度消してから書き直した場合は空からの書き込みになるので「届きました」になる
		if oldResponse == "" {
			return entity.RescueNotificationResponse, true
		}
		return entity.RescueNotificationResponseEdited, true
	}
	return "", false
}

// RescueMessageSender レスキュー通知の送り先。本番は *slack.SlackService
type RescueMessageSender interface {
	SendRescueMessage(slack.RescueMessageParams, string) error
}

// RescueNotifier レスキューの対応状況・返答が変わったことを送信者本人にSlack DMで知らせる。
// 各Update*Rescueが更新前後の値を渡し、知らせるべき変化ならその場で送る
type RescueNotifier interface {
	TroubleRescueUpdated(ctx context.Context, before, after *entity.TroubleRescueForGet)
	QuestionRescueUpdated(ctx context.Context, before, after *entity.QuestionRescueForGet)
	ShorthandedRescueUpdated(ctx context.Context, before, after *entity.ShorthandedRescueForGet)
}

type rescueNotifier struct {
	sender  RescueMessageSender
	taskRep rep.TaskRepository
	userRep rep.UserRepository
	// 送信の走らせ方。本番はgoroutineで裏に回し、テストではその場で実行して結果を確かめる
	run func(func())
}

func NewRescueNotifier(sender RescueMessageSender, taskRep rep.TaskRepository, userRep rep.UserRepository) RescueNotifier {
	return &rescueNotifier{
		sender:  sender,
		taskRep: taskRep,
		userRep: userRep,
		run:     func(f func()) { go f() },
	}
}

// 送信時刻の整形とタスク名の補完は「本部からの返答」タブ(GET /rescues)と同じものを使う
func (n *rescueNotifier) TroubleRescueUpdated(ctx context.Context, before, after *entity.TroubleRescueForGet) {
	n.notify(ctx, before.Status, before.Response, after.Status, after.Response, after.UserID, func(ctx context.Context) slack.RescueMessageParams {
		taskName, err := findTaskName(ctx, n.taskRep, strconv.Itoa(after.TaskID))
		if err != nil {
			taskName = "タスク外"
		}
		res := entity.NewTroubleRescueResponse(after, "", taskName)
		return slack.RescueMessageParams{
			Number:    fmt.Sprintf("T%d", after.ID),
			TypeLabel: "トラブル",
			Time:      res.Time,
			Details:   nonEmptyLines("発生タスク", taskName, "発生場所", after.Place, "内容", after.Detail),
		}
	})
}

func (n *rescueNotifier) QuestionRescueUpdated(ctx context.Context, before, after *entity.QuestionRescueForGet) {
	n.notify(ctx, before.Status, before.Response, after.Status, after.Response, after.UserID, func(context.Context) slack.RescueMessageParams {
		res := entity.NewQuestionRescueResponse(after, "")
		return slack.RescueMessageParams{
			Number:    fmt.Sprintf("Q%d", after.ID),
			TypeLabel: "質問",
			Time:      res.Time,
			Details:   nonEmptyLines("質問", after.Question),
		}
	})
}

func (n *rescueNotifier) ShorthandedRescueUpdated(ctx context.Context, before, after *entity.ShorthandedRescueForGet) {
	n.notify(ctx, before.Status, before.Response, after.Status, after.Response, after.UserID, func(ctx context.Context) slack.RescueMessageParams {
		taskName, err := findTaskName(ctx, n.taskRep, strconv.Itoa(after.TaskID))
		if err != nil {
			taskName = "不明なタスク"
		}
		res := entity.NewShorthandedRescueResponse(after, "", taskName)
		return slack.RescueMessageParams{
			Number:    fmt.Sprintf("S%d", after.ID),
			TypeLabel: "人が来ない",
			Time:      res.Time,
			Details:   nonEmptyLines("発生タスク", taskName, "送り先の場所", after.Place, "足りない人数", strconv.Itoa(after.MissingNumber)+"人"),
		}
	})
}

// 知らせるべき変化なら、メッセージを組み立てて送信者にDMする。
// 本部の書き込み(GASのPUT)をSlackの応答待ちで遅らせないよう、組み立てと送信は裏で行う。
// 送信に失敗しても送り直さない。返答はアプリの「本部からの返答」タブで見られるため
func (n *rescueNotifier) notify(ctx context.Context, oldStatus, oldResponse, newStatus, newResponse string, userID int, build func(context.Context) slack.RescueMessageParams) {
	kind, ok := rescueNotificationKind(oldStatus, oldResponse, newStatus, newResponse)
	if !ok {
		return
	}
	// リクエストのctxはPUTの応答を返した時点でキャンセルされるので、裏の処理からは切り離す
	ctx = context.WithoutCancel(ctx)
	n.run(func() {
		params := build(ctx)
		params.Title = rescueNotificationTitle[kind]
		params.Status = rescueStatusLabel[newStatus]
		params.Response = newResponse

		slackUserID, err := findSlackUserID(ctx, n.userRep, userID)
		if err != nil {
			logpkg.Printf("レスキュー通知の送信先を取得できませんでした(%s): %v", params.Number, err)
			return
		}
		// Slack IDが無い人にはシフト変更通知と同じく何も送らない
		if slackUserID == "" {
			return
		}
		if err := n.sender.SendRescueMessage(params, slackUserID); err != nil {
			logpkg.Printf("レスキュー通知の送信に失敗しました(%s): %v", params.Number, err)
		}
	})
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
