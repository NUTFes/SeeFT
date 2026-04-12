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

	// 3. 時刻情報を取得（全グループ共通、1回だけ）
	timeMap, err := n.loadTimeMap(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to load time map")
	}

	// 4. ユーザーIDと日付IDでグルーピング
	grouped := n.GroupNotificationsByUserAndDate(logs)

	// 5. 各グループを処理

	for key, group := range grouped {
		// keyは "userID_dateID" 形式
		var logIDs []int
		parts := strings.Split(key, "_")
		if len(parts) != 2 {
			logpkg.Printf("Invalid group key: %s", key)
			continue
		}
		userID, err1 := strconv.Atoi(parts[0])
		dateID, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			continue
		}

		// メッセージを生成して送信
		if err := n.processGroup(ctx, group, userID, dateID, timeMap); err != nil {
			logpkg.Printf("Failed to process group %s: %v", key, err)
			continue
		}

		for _, log := range group {
			logIDs = append(logIDs, log.ID)
		}

		if err := n.actionLogRepo.MarkAsSent(ctx, logIDs); err != nil {
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

// loadShiftMap ログからshift情報をまとめて取得しmapで返す
func (n *notificationUseCase) loadShiftMap(ctx context.Context, logs []entity.ActionLog) (map[int]entity.ShiftAdmin, error) {
	shiftMap := make(map[int]entity.ShiftAdmin)
	for _, log := range logs {
		// 既に取得済みならスキップ（重複クエリを防ぐ）
		if _, ok := shiftMap[log.ShiftID]; ok {
			continue
		}
		// DBから1件取得
		shiftRow, err := n.shiftRep.Find(ctx, strconv.Itoa(log.ShiftID))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to find shift %d", log.ShiftID)
		}
		// DBの行データをShiftAdminに変換
		var shift entity.ShiftAdmin
		err = shiftRow.Scan(
			&shift.ID, &shift.TaskID, &shift.UserID, &shift.YearID, &shift.DateID,
			&shift.TimeID, &shift.WeatherID, &shift.IsAttendance, &shift.CreatedAt, &shift.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to scan shift %d", log.ShiftID)
		}
		// mapに保存
		shiftMap[log.ShiftID] = shift
	}
	return shiftMap, nil
}

// loadTaskMap shiftMapからタスク情報をまとめて取得しmapで返す
func (n *notificationUseCase) loadTaskMap(ctx context.Context, shiftMap map[int]entity.ShiftAdmin) (map[int]entity.Task, error) {
	taskMap := make(map[int]entity.Task)
	for _, shift := range shiftMap {
		// 既に取得済みならスキップ（重複クエリを防ぐ）
		if _, ok := taskMap[shift.TaskID]; ok {
			continue
		}
		// DBから1件取得
		taskRow, err := n.taskRep.Find(ctx, strconv.Itoa(shift.TaskID))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to find task %d", shift.TaskID)
		}
		// DBの行データをTaskに変換
		var task entity.Task
		err = taskRow.Scan(
			&task.ID, &task.Task, &task.PlaceID, &task.Url, &task.BureauID,
			&task.MaxMember, &task.Color, &task.Remark, &task.YearID,
			&task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to scan task %d", shift.TaskID)
		}
		// mapに保存
		taskMap[shift.TaskID] = task
	}
	return taskMap, nil
}

// loadTimeMap 全時刻情報を取得しmapで返す（96件固定）
func (n *notificationUseCase) loadTimeMap(ctx context.Context) (map[int]entity.Time, error) {
	timeMap := make(map[int]entity.Time)
	timeRow, err := n.timeRep.All(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get all times")
	}
	defer timeRow.Close()

	for timeRow.Next() {
		var time entity.Time
		err := timeRow.Scan(&time.ID, &time.Time, &time.CreatedAt, &time.UpdatedAt)

		if err != nil {
			logpkg.Printf("Failed to scan time: %v", err)
			continue
		}
		timeMap[time.ID] = time
	}
	return timeMap, nil
}

// processGroup グループ化されたログを処理してSlackに送信
func (n *notificationUseCase) processGroup(ctx context.Context, logs []entity.ActionLog, userID, dateID int, timeMap map[int]entity.Time) error {
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

	// shift情報をまとめて取得（N+1クエリ対策）
	shiftMap, err := n.loadShiftMap(ctx, logs)
	if err != nil {
		return errors.Wrapf(err, "failed to load shift map")
	}

	taskMap, err := n.loadTaskMap(ctx, shiftMap)
	if err != nil {
		return errors.Wrapf(err, "failed to load task map")
	}

	// ログを時間順にソート（mapから取得）
	sortedLogs := n.sortLogsByTime(logs, shiftMap)

	// 連続時間を計算してメッセージを生成
	timeRange, changes := n.buildGroupedMessage(sortedLogs, shiftMap, taskMap, timeMap)

	// 天気情報を取得（mapから取得）
	weather := "不明"
	if len(sortedLogs) > 0 {
		shift := shiftMap[sortedLogs[0].ShiftID]
		weatherRow, err := n.weatherRep.Find(ctx, strconv.Itoa(shift.WeatherID))
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

// sortLogsByTime ログを時間順にソート（shiftMapから取得）
func (n *notificationUseCase) sortLogsByTime(logs []entity.ActionLog, shiftMap map[int]entity.ShiftAdmin) []entity.ActionLog {
	type logWithTime struct {
		log    entity.ActionLog
		timeID int
	}

	logsWithTime := make([]logWithTime, 0, len(logs))
	for _, log := range logs {
		shift, ok := shiftMap[log.ShiftID]
		if !ok {
			continue
		}
		logsWithTime = append(logsWithTime, logWithTime{
			log:    log,
			timeID: shift.TimeID,
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

	return sortedLogs
}

// buildGroupedMessage グルーピング済みメッセージを生成
func (n *notificationUseCase) buildGroupedMessage(logs []entity.ActionLog, shiftMap map[int]entity.ShiftAdmin, taskMap map[int]entity.Task, timeMap map[int]entity.Time) (string, string) {
	if len(logs) == 0 {
		return "", ""
	}

	// 各ログのtime_idを取得
	timeIDs := make([]int, 0, len(logs))
	for _, log := range logs {
		shift, ok := shiftMap[log.ShiftID]
		if !ok {
			continue
		}

		timeIDs = append(timeIDs, shift.TimeID)
	}

	if len(timeIDs) == 0 {
		return "", ""
	}

	// 連続時間を計算
	timeRanges := n.CalculateTimeRanges(timeMap, timeIDs)

	// 変更内容を構築
	changes := n.buildChangesList(logs, shiftMap, taskMap)

	return timeRanges, changes
}

// CalculateTimeRanges 連続するTimeIDから時間範囲を計算
func (n *notificationUseCase) CalculateTimeRanges(timeMap map[int]entity.Time, timeIDs []int) string {
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
			ranges = append(ranges, n.formatTimeRange(timeMap, start, end))
			start = timeIDs[i]
			end = timeIDs[i]
		}
	}
	// 最後の範囲を追加
	ranges = append(ranges, n.formatTimeRange(timeMap, start, end))

	return strings.Join(ranges, ", ")
}

// formatTimeRange 時間範囲をフォーマット
func (n *notificationUseCase) formatTimeRange(timeMap map[int]entity.Time, startTimeID, endTimeID int) string {

	startTime := timeMap[startTimeID]

	endTime, ok := timeMap[endTimeID+1]
	if !ok {
		endTime = entity.Time{Time: "0:00"}
	}

	return fmt.Sprintf("%s 〜 %s", startTime.Time, endTime.Time)
}

// buildChangesList 変更内容のリストを構築
func (n *notificationUseCase) buildChangesList(logs []entity.ActionLog, shiftMap map[int]entity.ShiftAdmin, taskMap map[int]entity.Task) string {
	var changes []string

	for _, log := range logs {
		// diff_payloadをパース
		var payload map[string]interface{}
		if err := json.Unmarshal(log.DiffPayload, &payload); err != nil {
			continue
		}

		// シフト情報を取得
		shift, ok := shiftMap[log.ShiftID]

		if !ok {
			continue
		}

		task, ok := taskMap[shift.TaskID]
		if !ok {
			continue
		}

		// 変更内容を構築
		var changeText string
		switch log.ActionType {
		case "CREATE":
			// diff_payloadから新規タスク名を取得
			newTask := task.Task // フォールバック: DB現在値
			if items, ok := payload["changes"].([]interface{}); ok && len(items) > 0 {
				if change, ok := items[0].(map[string]interface{}); ok {
					if name, ok := change["new"].(string); ok {
						newTask = name
					}
				}
			}
			changeText = fmt.Sprintf("%s（新規）", newTask)
		case "UPDATE":
			// diff_payloadからold/newを取得
			oldTask := "（不明）"
			newTask := task.Task // フォールバック: DB現在値
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
