package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	logpkg "log"
	"sort"
	"strconv"
	"strings"

	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/NUTFes/SeeFT/api/lib/externals/slack"
	rep "github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/pkg/errors"
)

type notificationUseCase struct {
	actionLogRepo rep.ActionLogRepository
	slackService  *slack.SlackService
	userRep       rep.UserRepository
	dateRep       rep.DateRepository
	timeRep       rep.TimeRepository
	taskRep       rep.TaskRepository
	shiftRep      rep.ShiftRepository
	weatherRep    rep.WeatherRepository
}

type NotificationUseCase interface {
	ProcessUnsentNotifications(ctx context.Context) error
}

func NewNotificationUseCase(
	actionLogRepo rep.ActionLogRepository,
	slackService *slack.SlackService,
	userRep rep.UserRepository,
	dateRep rep.DateRepository,
	timeRep rep.TimeRepository,
	taskRep rep.TaskRepository,
	shiftRep rep.ShiftRepository,
	weatherRep rep.WeatherRepository,
) NotificationUseCase {
	return &notificationUseCase{
		actionLogRepo: actionLogRepo,
		slackService:  slackService,
		userRep:       userRep,
		dateRep:       dateRep,
		timeRep:       timeRep,
		taskRep:       taskRep,
		shiftRep:      shiftRep,
		weatherRep:    weatherRep,
	}
}

// ProcessUnsentNotifications 未送信の通知を処理する
func (n *notificationUseCase) ProcessUnsentNotifications(ctx context.Context) error {
	// 1. 未送信ログを取得
	rows, err := n.actionLogRepo.GetUnsentLogs(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to get unsent logs")
	}
	defer rows.Close()

	// 2. ログを読み込んでエンティティに変換
	var logs []entity.ActionLog
	for rows.Next() {
		var actionLog entity.ActionLog
		err := rows.Scan(
			&actionLog.ID,
			&actionLog.ShiftID,
			&actionLog.UserID,
			&actionLog.DateID,
			&actionLog.ActionType,
			&actionLog.DiffPayload,
			&actionLog.IsSent,
			&actionLog.CreatedAt,
		)
		if err != nil {
			logpkg.Printf("Failed to scan action log: %v", err)
			continue
		}
		logs = append(logs, actionLog)
	}

	if len(logs) == 0 {
		return nil // 未送信ログがない場合は何もしない
	}

	// 3. ユーザーIDと日付IDでグルーピング
	grouped := n.GroupNotificationsByUserAndDate(logs)

	// 4. 各グループを処理
	var sentLogIDs []int
	for key, group := range grouped {
		// keyは "userID_dateID" 形式
		parts := strings.Split(key, "_")
		if len(parts) != 2 {
			logpkg.Printf("Invalid group key: %s", key)
			continue
		}
		userID, _ := strconv.Atoi(parts[0])
		dateID, _ := strconv.Atoi(parts[1])

		// メッセージを生成して送信
		if err := n.processGroup(ctx, group, userID, dateID); err != nil {
			logpkg.Printf("Failed to process group %s: %v", key, err)
			continue
		}

		// 送信済みログIDを収集
		for _, log := range group {
			sentLogIDs = append(sentLogIDs, log.ID)
		}
	}

	// 5. 送信済みフラグを更新
	if len(sentLogIDs) > 0 {
		if err := n.actionLogRepo.MarkAsSent(ctx, sentLogIDs); err != nil {
			return errors.Wrapf(err, "failed to mark logs as sent")
		}
	}

	return nil
}

// GroupNotificationsByUserAndDate ユーザーIDと日付IDでグルーピング
func (n *notificationUseCase) GroupNotificationsByUserAndDate(logs []entity.ActionLog) map[string][]entity.ActionLog {
	grouped := make(map[string][]entity.ActionLog)
	for _, log := range logs {
		key := fmt.Sprintf("%d_%d", log.UserID, log.DateID)
		grouped[key] = append(grouped[key], log)
	}
	return grouped
}

// processGroup グループ化されたログを処理してSlackに送信
func (n *notificationUseCase) processGroup(ctx context.Context, logs []entity.ActionLog, userID, dateID int) error {
	// ユーザー情報を取得
	userRow, err := n.userRep.Find(ctx, strconv.Itoa(userID))
	if err != nil {
		return errors.Wrapf(err, "failed to find user")
	}
	var user entity.User
	var slackUserID sql.NullString
	err = userRow.Scan(
		&user.ID, &user.Name, &user.Mail, &user.GradeID, &user.DepartmentID,
		&user.BureauID, &user.RoleID, &user.StudentNumber, &user.Tel,
		&user.Password, &user.CreatedAt, &user.UpdatedAt, &slackUserID,
	)
	if err != nil {
		return errors.Wrapf(err, "failed to scan user")
	}
	if slackUserID.Valid {
		user.SlackUserID = slackUserID.String
	}

	// 日付情報を取得
	dateRow, err := n.dateRep.Find(ctx, strconv.Itoa(dateID))
	if err != nil {
		return errors.Wrapf(err, "failed to find date")
	}
	var date entity.Date
	err = dateRow.Scan(
		&date.ID, &date.YearID, &date.Name, &date.Date,
		&date.CreatedAt, &date.UpdatedAt,
	)
	if err != nil {
		return errors.Wrapf(err, "failed to scan date")
	}

	// 天気情報を取得（最初のログから取得）
	// weatherID := 0
	// if len(logs) > 0 {
	// 	// shift_idから天気情報を取得する必要がある
	// 	// ここでは簡略化のため、diff_payloadから取得を試みる
	// 	// 実際にはshiftテーブルから取得する必要がある
	// }

	// ログを時間順にソート（shift_idからtime_idを取得してソート）
	sortedLogs, err := n.sortLogsByTime(ctx, logs)
	if err != nil {
		return errors.Wrapf(err, "failed to sort logs by time")
	}

	// 連続時間を計算してメッセージを生成
	timeRange, changes := n.buildGroupedMessage(ctx, sortedLogs)

	// 天気情報を取得（最初のシフトから）
	weather := "不明"
	if len(sortedLogs) > 0 {
		shiftRow, err := n.shiftRep.Find(ctx, strconv.Itoa(sortedLogs[0].ShiftID))
		if err == nil {
			var shiftID, taskID, shiftUserID, yearID, shiftDateID, timeID, shiftWeatherID int
			var isAttendance bool
			var createdAt, updatedAt interface{}
			err = shiftRow.Scan(
				&shiftID, &taskID, &shiftUserID, &yearID, &shiftDateID,
				&timeID, &shiftWeatherID, &isAttendance, &createdAt, &updatedAt,
			)
			if err == nil {
				weatherRow, err := n.weatherRep.Find(ctx, strconv.Itoa(shiftWeatherID))
				if err == nil {
					var weatherEntity entity.Weather
					err = weatherRow.Scan(
						&weatherEntity.ID, &weatherEntity.Weather,
						&weatherEntity.CreatedAt, &weatherEntity.UpdatedAt,
					)
					if err == nil {
						weather = weatherEntity.Weather
					}
				}
			}
		}
	}

	// Slackメッセージを構築
	blocks := n.slackService.BuildMessageBlocks(slack.MessageParams{
		Title:     "シフト変更通知",
		UserName:  user.Name,
		Date:      date.Name,
		Weather:   weather,
		TimeRange: timeRange,
		Changes:   changes,
	})

	// Slackに送信
	if err := n.slackService.SendMessage(blocks, user.SlackUserID); err != nil {
		return errors.Wrapf(err, "failed to send slack message")
	}

	return nil
}

// sortLogsByTime ログを時間順にソート
func (n *notificationUseCase) sortLogsByTime(ctx context.Context, logs []entity.ActionLog) ([]entity.ActionLog, error) {
	// 各ログのshift_idからtime_idを取得してソート
	type logWithTime struct {
		log    entity.ActionLog
		timeID int
	}

	logsWithTime := make([]logWithTime, 0, len(logs))
	for _, log := range logs {
		shiftRow, err := n.shiftRep.Find(ctx, strconv.Itoa(log.ShiftID))
		if err != nil {
			continue
		}

		var shiftID, taskID, userID, yearID, dateID, timeID, weatherID int
		var isAttendance bool
		var createdAt, updatedAt interface{}
		err = shiftRow.Scan(
			&shiftID, &taskID, &userID, &yearID, &dateID,
			&timeID, &weatherID, &isAttendance, &createdAt, &updatedAt,
		)
		if err != nil {
			continue
		}

		logsWithTime = append(logsWithTime, logWithTime{
			log:    log,
			timeID: timeID,
		})
	}

	// timeIDでソート
	sort.Slice(logsWithTime, func(i, j int) bool {
		return logsWithTime[i].timeID < logsWithTime[j].timeID
	})

	// ソート済みログを返す
	sortedLogs := make([]entity.ActionLog, len(logsWithTime))
	for i, lwt := range logsWithTime {
		sortedLogs[i] = lwt.log
	}

	return sortedLogs, nil
}

// buildGroupedMessage グルーピング済みメッセージを生成
func (n *notificationUseCase) buildGroupedMessage(ctx context.Context, logs []entity.ActionLog) (string, string) {
	if len(logs) == 0 {
		return "", ""
	}

	// 各ログのtime_idを取得
	timeIDs := make([]int, 0, len(logs))
	for _, log := range logs {
		shiftRow, err := n.shiftRep.Find(ctx, strconv.Itoa(log.ShiftID))
		if err != nil {
			continue
		}

		var shiftID, taskID, userID, yearID, dateID, timeID, weatherID int
		var isAttendance bool
		var createdAt, updatedAt interface{}
		err = shiftRow.Scan(
			&shiftID, &taskID, &userID, &yearID, &dateID,
			&timeID, &weatherID, &isAttendance, &createdAt, &updatedAt,
		)
		if err != nil {
			continue
		}

		timeIDs = append(timeIDs, timeID)
	}

	if len(timeIDs) == 0 {
		return "", ""
	}

	// 連続時間を計算
	timeRanges := n.CalculateTimeRanges(ctx, timeIDs)

	// 変更内容を構築
	changes := n.buildChangesList(ctx, logs)

	return timeRanges, changes
}

// CalculateTimeRanges 連続するTimeIDから時間範囲を計算
func (n *notificationUseCase) CalculateTimeRanges(ctx context.Context, timeIDs []int) string {
	if len(timeIDs) == 0 {
		return ""
	}

	// ソート済みと仮定
	sort.Ints(timeIDs)

	var ranges []string
	start := timeIDs[0]
	end := timeIDs[0]

	for i := 1; i < len(timeIDs); i++ {
		if timeIDs[i] == end+1 {
			// 連続している
			end = timeIDs[i]
		} else {
			// 連続が途切れた
			ranges = append(ranges, n.formatTimeRange(ctx, start, end))
			start = timeIDs[i]
			end = timeIDs[i]
		}
	}
	// 最後の範囲を追加
	ranges = append(ranges, n.formatTimeRange(ctx, start, end))

	return strings.Join(ranges, ", ")
}

// formatTimeRange 時間範囲をフォーマット
func (n *notificationUseCase) formatTimeRange(ctx context.Context, startTimeID, endTimeID int) string {
	// 開始時刻を取得
	startRow, err := n.timeRep.Find(ctx, strconv.Itoa(startTimeID))
	if err != nil {
		return ""
	}
	var startTime entity.Time
	err = startRow.Scan(&startTime.ID, &startTime.Time, &startTime.CreatedAt, &startTime.UpdatedAt)
	if err != nil {
		return ""
	}

	// 終了時刻を取得（endTimeID+1の時刻）
	endRow, err := n.timeRep.Find(ctx, strconv.Itoa(endTimeID+1))
	if err != nil {
		// 次の時刻がない場合は、endTimeIDの時刻を使用
		endRow, err = n.timeRep.Find(ctx, strconv.Itoa(endTimeID))
		if err != nil {
			return startTime.Time
		}
	}
	var endTime entity.Time
	err = endRow.Scan(&endTime.ID, &endTime.Time, &endTime.CreatedAt, &endTime.UpdatedAt)
	if err != nil {
		return startTime.Time
	}

	return fmt.Sprintf("%s 〜 %s", startTime.Time, endTime.Time)
}

// buildChangesList 変更内容のリストを構築
func (n *notificationUseCase) buildChangesList(ctx context.Context, logs []entity.ActionLog) string {
	var changes []string

	for _, log := range logs {
		// diff_payloadをパース
		var payload map[string]interface{}
		if err := json.Unmarshal(log.DiffPayload, &payload); err != nil {
			continue
		}

		// シフト情報を取得
		shiftRow, err := n.shiftRep.Find(ctx, strconv.Itoa(log.ShiftID))
		if err != nil {
			continue
		}

		var shiftID, taskID, userID, yearID, dateID, timeID, weatherID int
		var isAttendance bool
		var createdAt, updatedAt interface{}
		err = shiftRow.Scan(
			&shiftID, &taskID, &userID, &yearID, &dateID,
			&timeID, &weatherID, &isAttendance, &createdAt, &updatedAt,
		)
		if err != nil {
			continue
		}

		// タスク情報を取得
		taskRow, err := n.taskRep.Find(ctx, strconv.Itoa(taskID))
		if err != nil {
			continue
		}
		var task entity.Task
		err = taskRow.Scan(
			&task.ID, &task.Task, &task.PlaceID, &task.Url, &task.BureauID,
			&task.MaxMember, &task.Color, &task.Remark, &task.YearID,
			&task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			continue
		}

		// 時刻を取得
		timeRow, err := n.timeRep.Find(ctx, strconv.Itoa(timeID))
		if err != nil {
			continue
		}
		var time entity.Time
		err = timeRow.Scan(&time.ID, &time.Time, &time.CreatedAt, &time.UpdatedAt)
		if err != nil {
			continue
		}

		// 変更内容を構築
		var changeText string
		switch log.ActionType {
		case "CREATE":
			changeText = fmt.Sprintf("%s（新規）", task.Task)
		case "UPDATE":
			// diff_payloadからold/newを取得
			oldTask := "（不明）"
			newTask := task.Task
			if changes, ok := payload["changes"].([]interface{}); ok && len(changes) > 0 {
				if change, ok := changes[0].(map[string]interface{}); ok {
					if old, ok := change["old"].(string); ok {
						oldTask = old
					}
					if new, ok := change["new"].(string); ok {
						newTask = new
					}
				}
			}
			changeText = fmt.Sprintf("%s → %s", oldTask, newTask)
		case "DELETE":
			oldTask := "（不明）"
			if deleted, ok := payload["deleted_task"].(string); ok {
				oldTask = deleted
			}
			changeText = fmt.Sprintf("%s（削除）", oldTask)
		default:
			changeText = fmt.Sprintf("%s", task.Task)
		}

		changes = append(changes, changeText)
	}

	return strings.Join(changes, "\n")
}
