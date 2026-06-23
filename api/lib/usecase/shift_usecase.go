package usecase

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/NUTFes/SeeFT/api/lib/entity"
	rep "github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/pkg/errors"
)

type shiftUseCase struct {
	rep           rep.ShiftRepository
	shiftCardRep  rep.ShiftCardRepository
	taskRep       rep.TaskRepository
	userRep       rep.UserRepository
	yearRep       rep.YearRepository
	dateRep       rep.DateRepository
	timeRep       rep.TimeRepository
	weatherRep    rep.WeatherRepository
	placeRep      rep.PlaceRepository
	gradeRep      rep.GradeRepository
	bureauRep     rep.BureauRepository
	actionLogRepo rep.ActionLogRepository
}

type ShiftUseCase interface {
	GetShiftCardsByUserAndDateAndWeather(context.Context, string, string, string) ([]entity.ShiftCard, error)
	GetUsersByShift(context.Context, string, string, string, string, string) (entity.ShiftUsers, error)
	GetShiftsAdmin(context.Context) ([]entity.ShiftAdmin, error)
	GetShiftAdminByID(context.Context, string) (entity.ShiftAdmin, error)
	CreateShiftAdmin(context.Context, string, string, string, string, string, string, string) (entity.ShiftAdmin, error)
	UpdateShiftAdmin(context.Context, string, string, string, string, string, string, string, string) (entity.ShiftAdmin, error)
	DeleteShiftAdmin(context.Context, string) error
	GetShiftsAdminByDateAndWeather(context.Context, string, string) ([]entity.ShiftAdmin, error)
	GetShiftsAdminByDateAndWeatherAndTime(context.Context, string, string, string, string) ([]entity.ShiftAdmin, error)
	GetMaxID(context.Context) (int, error)
	SaveShiftData(context.Context, entity.ShiftRequest) error
	SendToGAS(context.Context, entity.ShiftRequest) error
	UpdateShiftsFromGAS(context.Context, entity.ShiftChangeRequest) error
}

func NewShiftUseCase(
	rep rep.ShiftRepository,
	shiftCardRep rep.ShiftCardRepository,
	taskRep rep.TaskRepository,
	userRep rep.UserRepository,
	yearRep rep.YearRepository,
	dateRep rep.DateRepository,
	timeRep rep.TimeRepository,
	weatherRep rep.WeatherRepository,
	placeRep rep.PlaceRepository,
	gradeRep rep.GradeRepository,
	bureauRep rep.BureauRepository,
	actionLogRepo rep.ActionLogRepository) ShiftUseCase {
	return &shiftUseCase{rep, shiftCardRep, taskRep, userRep, yearRep, dateRep, timeRep, weatherRep, placeRep, gradeRep, bureauRep, actionLogRepo}
}

var TaskID, UserID, YearID, DateID, TimeID, WeatherID, PlaceID string

// 時間でソート
type ByTime []entity.Shift

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].Time.ID < a[j].Time.ID }

func (a *shiftUseCase) GetUsersByShift(c context.Context, task string, year string, date string, time string, weather string) (entity.ShiftUsers, error) {

	users := entity.User{}
	var shiftUsers entity.ShiftUsers

	// クエリー実行
	rows, err := a.rep.Users(c, task, year, date, time, weather)
	if err != nil {
		return shiftUsers, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		// JOINで取得したユーザー情報を直接Scan（個別SQLは不要、passwordは取得しない）
		err := rows.Scan(
			&users.ID,
			&users.Name,
			&users.Mail,
			&users.GradeID,
			&users.DepartmentID,
			&users.BureauID,
			&users.RoleID,
			&users.StudentNumber,
			&users.Tel,
			&users.CreatedAt,
			&users.UpdatedAt,
			&users.SlackUserID,
		)

		if err != nil {
			return shiftUsers, errors.Wrapf(err, "cannot connect SQL")
		}

		shiftUsers.Users = append(shiftUsers.Users, users)
	}

	row, err := a.yearRep.Find(c, year)
	err = row.Scan(
		&shiftUsers.Year.ID,
		&shiftUsers.Year.Year,
		&shiftUsers.Year.CreatedAt,
		&shiftUsers.Year.UpdatedAt,
	)

	row, err = a.dateRep.Find(c, date)
	err = row.Scan(
		&shiftUsers.Date.ID,
		&shiftUsers.Date.YearID,
		&shiftUsers.Date.Name,
		&shiftUsers.Date.Date,
		&shiftUsers.Date.CreatedAt,
		&shiftUsers.Date.UpdatedAt,
	)
	row, err = a.timeRep.Find(c, time)
	err = row.Scan(
		&shiftUsers.Time.ID,
		&shiftUsers.Time.Time,
		&shiftUsers.Time.CreatedAt,
		&shiftUsers.Time.UpdatedAt,
	)
	row, err = a.weatherRep.Find(c, weather)
	err = row.Scan(
		&shiftUsers.Weather.ID,
		&shiftUsers.Weather.Weather,
		&shiftUsers.Weather.CreatedAt,
		&shiftUsers.Weather.UpdatedAt,
	)

	if err != nil {
		return shiftUsers, errors.Wrapf(err, "cannot connect SQL")
	}

	return shiftUsers, nil
}

// Webアプリ用API

func (a *shiftUseCase) GetShiftsAdmin(c context.Context) ([]entity.ShiftAdmin, error) {
	shift := entity.ShiftAdmin{}
	var shifts []entity.ShiftAdmin

	// クエリー実行
	rows, err := a.rep.All(c)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		err := rows.Scan(
			&shift.ID,
			&shift.TaskID,
			&shift.UserID,
			&shift.YearID,
			&shift.DateID,
			&shift.TimeID,
			&shift.WeatherID,
			&shift.IsAttendance,
			&shift.CreatedAt,
			&shift.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot connect SQL")
		}

		shifts = append(shifts, shift)
	}
	return shifts, nil
}

func (a *shiftUseCase) GetShiftAdminByID(c context.Context, id string) (entity.ShiftAdmin, error) {
	var shift entity.ShiftAdmin
	row, err := a.rep.Find(c, id)
	err = row.Scan(
		&shift.ID,
		&shift.TaskID,
		&shift.UserID,
		&shift.YearID,
		&shift.DateID,
		&shift.TimeID,
		&shift.WeatherID,
		&shift.IsAttendance,
		&shift.CreatedAt,
		&shift.UpdatedAt,
	)

	if err != nil {
		return shift, err
	}
	return shift, nil
}

func (u *shiftUseCase) CreateShiftAdmin(c context.Context, taskID string, userID string, yearID string, dateID string, timeID string, weatherID string, isAttendance string) (entity.ShiftAdmin, error) {
	var latastShift entity.ShiftAdmin
	err := u.rep.Create(c, taskID, userID, yearID, dateID, timeID, weatherID, isAttendance)
	row, err := u.rep.FindLatestRecord(c)
	err = row.Scan(
		&latastShift.ID,
		&latastShift.TaskID,
		&latastShift.UserID,
		&latastShift.YearID,
		&latastShift.DateID,
		&latastShift.TimeID,
		&latastShift.WeatherID,
		&latastShift.IsAttendance,
		&latastShift.CreatedAt,
		&latastShift.UpdatedAt,
	)
	if err != nil {
		return latastShift, err
	}
	return latastShift, err
}

func (u *shiftUseCase) UpdateShiftAdmin(c context.Context, id string, taskID string, userID string, yearID string, dateID string, timeID string, weatherID string, isAttendance string) (entity.ShiftAdmin, error) {
	updatedShift := entity.ShiftAdmin{}
	var shift entity.ShiftAdmin

	row, err := u.rep.Find(c, id)
	err = row.Scan(
		&shift.ID,
		&shift.TaskID,
		&shift.UserID,
		&shift.YearID,
		&shift.DateID,
		&shift.TimeID,
		&shift.WeatherID,
		&shift.IsAttendance,
		&shift.CreatedAt,
		&shift.UpdatedAt,
	)
	if err != nil {
		return shift, err
	}

	if err = u.rep.Update(c, id, taskID, userID, yearID, dateID, timeID, weatherID, isAttendance); err != nil {
		return updatedShift, err
	}
	row, err = u.rep.Find(c, id)
	err = row.Scan(
		&updatedShift.ID,
		&updatedShift.TaskID,
		&updatedShift.UserID,
		&updatedShift.YearID,
		&updatedShift.DateID,
		&updatedShift.TimeID,
		&updatedShift.WeatherID,
		&updatedShift.IsAttendance,
		&updatedShift.CreatedAt,
		&updatedShift.UpdatedAt,
	)
	if err != nil {
		return updatedShift, err
	}
	return updatedShift, nil
}

func (u *shiftUseCase) DeleteShiftAdmin(c context.Context, id string) error {
	// 削除前にシフト情報を取得してaction_logに記録
	shift, err := u.GetShiftAdminByID(c, id)
	if err == nil {
		// タスク名を取得
		taskName := "（不明）"
		taskRow, taskErr := u.taskRep.Find(c, strconv.Itoa(shift.TaskID))
		if taskErr == nil {
			var task entity.Task
			if scanErr := taskRow.Scan(&task.ID, &task.Task, &task.PlaceID, &task.Url, &task.BureauID, &task.MaxMember, &task.Color, &task.Remark, &task.YearID, &task.CreatedAt, &task.UpdatedAt); scanErr == nil {
				taskName = task.Task
			}
		}

		diffPayload := map[string]interface{}{
			"deleted_task": taskName,
		}
		if u.actionLogRepo != nil {
			if logErr := u.actionLogRepo.Create(c, shift.ID, shift.UserID, shift.DateID, "DELETE", diffPayload); logErr != nil {
				log.Printf("action_log記録失敗(DELETE): %v", logErr)
			}
		}
	}

	err = u.rep.Destroy(c, id)
	return err
}

func (a *shiftUseCase) GetShiftsAdminByDateAndWeather(c context.Context, date string, weather string) ([]entity.ShiftAdmin, error) {
	shift := entity.ShiftAdmin{}
	var shifts []entity.ShiftAdmin

	// クエリー実行
	rows, err := a.rep.DateAndWeather(c, date, weather)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		err := rows.Scan(
			&shift.ID,
			&shift.TaskID,
			&shift.UserID,
			&shift.YearID,
			&shift.DateID,
			&shift.TimeID,
			&shift.WeatherID,
			&shift.IsAttendance,
			&shift.CreatedAt,
			&shift.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot connect SQL")
		}

		shifts = append(shifts, shift)
	}
	return shifts, nil
}

func (a *shiftUseCase) GetShiftsAdminByDateAndWeatherAndTime(c context.Context, date string, weather string, lower string, upper string) ([]entity.ShiftAdmin, error) {
	shift := entity.ShiftAdmin{}
	var shifts []entity.ShiftAdmin

	// クエリー実行
	rows, err := a.rep.DateAndWeatherAndTime(c, date, weather, lower, upper)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		err := rows.Scan(
			&shift.ID,
			&shift.TaskID,
			&shift.UserID,
			&shift.YearID,
			&shift.DateID,
			&shift.TimeID,
			&shift.WeatherID,
			&shift.IsAttendance,
			&shift.CreatedAt,
			&shift.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot connect SQL")
		}

		shifts = append(shifts, shift)
	}
	return shifts, nil
}

func (a *shiftUseCase) GetShiftCardsByUserAndDateAndWeather(c context.Context, userID string, dateID string, weatherID string) ([]entity.ShiftCard, error) {
	// グローバル変数を使わずローカル変数で処理
	shiftCards := make([]entity.ShiftCard, 0)

	// N+1問題対策: Grade/Bureauのマップを事前取得
	gradeMap, err := a.loadGradeMap(c)
	if err != nil {
		return nil, err
	}
	bureauMap, err := a.loadBureauMap(c)
	if err != nil {
		return nil, err
	}

	// 新規ShiftCardリポジトリを使用してデータ取得
	shiftData, err := a.shiftCardRep.GetOptimizedShiftData(c, userID, dateID, weatherID)
	if err != nil {
		return nil, err
	}

	// ShiftCardDataをentity.Shiftに変換
	shifts := a.convertShiftCardDataToShifts(shiftData)

	// 空タスクとNGタスクをフィルタリング（UseCase層での二重チェック）
	var validShifts []entity.Shift
	for _, shift := range shifts {
		// 二重チェック: リポジトリ層でも除外しているが、念のため再確認
		if shift.Task.Task == "" || shift.Task.Task == "NG" {
			continue // スキップ
		}
		validShifts = append(validShifts, shift)
	}

	// タスクIDごとにグループ化し、連続性を判定
	taskGroups := make(map[int][]entity.Shift)
	for _, shift := range validShifts {
		taskGroups[shift.Task.ID] = append(taskGroups[shift.Task.ID], shift)
	}

	// 各タスクグループを処理
	for _, taskShifts := range taskGroups {
		// 時間順にソート
		sort.Sort(ByTime(taskShifts))

		// 連続するTimeIDでグループ化
		continuousGroups := a.groupContinuousShifts(taskShifts)

		for _, group := range continuousGroups {
			if len(group) == 0 {
				continue
			}

			card := a.createShiftCardFromGroup(c, group, gradeMap, bureauMap)
			shiftCards = append(shiftCards, card)
		}
	}

	// 最終的な時系列ソート
	sort.Slice(shiftCards, func(i, j int) bool {
		return a.compareTimeStrings(shiftCards[i].StartTime, shiftCards[j].StartTime) < 0
	})

	return shiftCards, nil
}

// convertShiftCardDataToShifts はShiftCardDataをentity.Shiftに変換
func (a *shiftUseCase) convertShiftCardDataToShifts(data []entity.ShiftCardData) []entity.Shift {
	var shifts []entity.Shift

	for _, d := range data {
		yearInt, _ := strconv.Atoi(d.YearValue)
		shift := entity.Shift{
			ID: d.ShiftID,
			Task: entity.TaskMobile{
				ID:    d.TaskID,
				Task:  d.TaskName,
				Color: d.TaskColor,
				Place: d.PlaceName,
				Url:   d.TaskURL,
			},
			User: entity.User{
				ID:       d.UserID,
				Name:     d.UserName,
				BureauID: d.UserBureauID,
				GradeID:  d.UserGradeID,
			},
			Year: entity.Year{
				ID:   d.YearID,
				Year: yearInt,
			},
			Date: entity.Date{
				ID:   d.DateID,
				Date: d.DateValue,
			},
			Time: entity.Time{
				ID:   d.TimeID,
				Time: d.TimeValue,
			},
			Weather: entity.Weather{
				ID:      d.WeatherID,
				Weather: d.WeatherValue,
			},
			IsAttendance: d.IsAttendance,
		}
		shifts = append(shifts, shift)
	}

	return shifts
}

// groupContinuousShifts は同じタスクで時間が連続するシフトをグループ化
func (a *shiftUseCase) groupContinuousShifts(taskShifts []entity.Shift) [][]entity.Shift {
	if len(taskShifts) == 0 {
		return [][]entity.Shift{}
	}

	var groups [][]entity.Shift
	currentGroup := []entity.Shift{taskShifts[0]}

	for i := 1; i < len(taskShifts); i++ {
		prev := taskShifts[i-1]
		curr := taskShifts[i]

		// 連続条件: 同じタスクID && TimeIDが+1
		if prev.Task.ID == curr.Task.ID && curr.Time.ID == prev.Time.ID+1 {
			currentGroup = append(currentGroup, curr)
		} else {
			// 新しいグループを開始
			groups = append(groups, currentGroup)
			currentGroup = []entity.Shift{curr}
		}
	}

	// 最後のグループを追加
	if len(currentGroup) > 0 {
		groups = append(groups, currentGroup)
	}

	return groups
}

// compareTimeStrings は時刻文字列を比較
func (a *shiftUseCase) compareTimeStrings(time1, time2 string) int {
	// "8:00" と "10:00" のような時刻文字列を比較
	t1Parts := strings.Split(time1, ":")
	t2Parts := strings.Split(time2, ":")

	if len(t1Parts) != 2 || len(t2Parts) != 2 {
		return 0
	}

	h1, _ := strconv.Atoi(t1Parts[0])
	m1, _ := strconv.Atoi(t1Parts[1])
	h2, _ := strconv.Atoi(t2Parts[0])
	m2, _ := strconv.Atoi(t2Parts[1])

	t1Minutes := h1*60 + m1
	t2Minutes := h2*60 + m2

	if t1Minutes < t2Minutes {
		return -1
	} else if t1Minutes > t2Minutes {
		return 1
	}
	return 0
}

func (a *shiftUseCase) createShiftCardFromGroup(c context.Context, group []entity.Shift, gradeMap map[int]string, bureauMap map[int]string) entity.ShiftCard {
	if len(group) == 0 {
		// 空の場合でも必ず配列を初期化
		return entity.ShiftCard{
			TaskName:     "",
			StartTime:    "",
			EndTime:      "",
			Place:        "",
			Url:          "",
			ShiftMembers: []entity.ShiftMembers{}, // 空配列で初期化
			BeforeMembers: entity.ShiftMembers{
				STime:   "",
				ETime:   "",
				Members: []entity.ShiftMember{}, // 空配列で初期化
			},
			AfterMembers: entity.ShiftMembers{
				STime:   "",
				ETime:   "",
				Members: []entity.ShiftMember{}, // 空配列で初期化
			},
		}
	}

	first := group[0]
	last := group[len(group)-1]

	// 終了時刻を取得
	endTime, err := a.getNextTimeString(c, last.Time.ID)
	if err != nil {
		endTime = last.Time.Time // フォールバック
	}

	// ShiftCardを初期化（全フィールドを明示的に設定）
	shiftCard := entity.ShiftCard{
		TaskName:     first.Task.Task,
		StartTime:    first.Time.Time,
		EndTime:      endTime,
		Place:        first.Task.Place,
		Url:          first.Task.Url,
		ShiftMembers: []entity.ShiftMembers{}, // 空配列で初期化
	}

	// ShiftMembersを生成（15分ごと）
	var shiftMembers []entity.ShiftMembers
	for i, shift := range group {
		members := a.getShiftMembersForTime(c, shift, gradeMap, bureauMap)

		// メンバーがnilの場合は空配列を設定
		if members == nil {
			members = []entity.ShiftMember{}
		}

		// 終了時刻を計算
		var eTime string
		if i < len(group)-1 {
			eTime = group[i+1].Time.Time
		} else {
			nextTime, err := a.getNextTimeString(c, shift.Time.ID)
			if err != nil {
				eTime = shift.Time.Time
			} else {
				eTime = nextTime
			}
		}

		shiftMember := entity.ShiftMembers{
			STime:   shift.Time.Time,
			ETime:   eTime,
			Members: members,
		}

		shiftMembers = append(shiftMembers, shiftMember)
	}
	shiftCard.ShiftMembers = shiftMembers

	// BeforeMembersを取得（nilチェック付き）
	beforeMembers := a.getBeforeMembers(c, first, gradeMap, bureauMap)
	if beforeMembers.Members == nil {
		beforeMembers.Members = []entity.ShiftMember{}
	}
	shiftCard.BeforeMembers = beforeMembers

	// AfterMembersを取得（nilチェック付き）
	afterMembers := a.getAfterMembers(c, last, gradeMap, bureauMap)
	if afterMembers.Members == nil {
		afterMembers.Members = []entity.ShiftMember{}
	}
	shiftCard.AfterMembers = afterMembers

	return shiftCard
}

// getShiftMembersForTime は指定時間のメンバーを取得
func (a *shiftUseCase) getShiftMembersForTime(c context.Context, shift entity.Shift, gradeMap map[int]string, bureauMap map[int]string) []entity.ShiftMember {
	shiftUsers, err := a.GetUsersByShift(c,
		strconv.Itoa(shift.Task.ID),
		strconv.Itoa(shift.Year.ID),
		strconv.Itoa(shift.Date.ID),
		strconv.Itoa(shift.Time.ID),
		strconv.Itoa(shift.Weather.ID))

	if err != nil {
		return []entity.ShiftMember{} // エラー時は空配列を返す
	}

	var members []entity.ShiftMember
	for _, user := range shiftUsers.Users {
		// マップからGrade/Bureau情報を取得（N+1問題を回避）
		grade := gradeMap[user.GradeID]
		bureau := bureauMap[user.BureauID]

		member := entity.ShiftMember{
			Name:   user.Name,
			Grade:  grade,
			Bureau: bureau,
		}
		members = append(members, member)
	}

	if members == nil {
		return []entity.ShiftMember{}
	}

	return members
}

// getBeforeMembers は前の時間のメンバーを取得
func (a *shiftUseCase) getBeforeMembers(c context.Context, firstShift entity.Shift, gradeMap map[int]string, bureauMap map[int]string) entity.ShiftMembers {
	// 初期化（空配列を保証）
	result := entity.ShiftMembers{
		STime:   "",
		ETime:   "",
		Members: []entity.ShiftMember{},
	}

	// 前の時間のID
	prevTimeID := firstShift.Time.ID - 1
	if prevTimeID < 1 {
		return result
	}

	// 前の時間の文字列を取得
	prevTimeString, err := a.getPreviousTimeString(c, firstShift.Time.ID)
	if err != nil {
		return result
	}

	// 前の時間のメンバーを取得
	prevUsers, err := a.GetUsersByShift(c,
		strconv.Itoa(firstShift.Task.ID),
		strconv.Itoa(firstShift.Year.ID),
		strconv.Itoa(firstShift.Date.ID),
		strconv.Itoa(prevTimeID),
		strconv.Itoa(firstShift.Weather.ID))

	if err != nil {
		return result
	}

	var members []entity.ShiftMember
	for _, user := range prevUsers.Users {
		// マップからGrade/Bureau情報を取得（N+1問題を回避）
		grade := gradeMap[user.GradeID]
		bureau := bureauMap[user.BureauID]

		member := entity.ShiftMember{
			Name:   user.Name,
			Grade:  grade,
			Bureau: bureau,
		}
		members = append(members, member)
	}

	if members == nil {
		members = []entity.ShiftMember{}
	}

	result.STime = prevTimeString
	result.ETime = firstShift.Time.Time
	result.Members = members

	return result
}

// getAfterMembers は後の時間のメンバーを取得（同様の実装）
func (a *shiftUseCase) getAfterMembers(c context.Context, lastShift entity.Shift, gradeMap map[int]string, bureauMap map[int]string) entity.ShiftMembers {
	// 次の時間帯情報（空配列は保証）
	result := entity.ShiftMembers{
		STime:   "",
		ETime:   "",
		Members: []entity.ShiftMember{},
	}

	// 次の時間ID
	nextTimeID := lastShift.Time.ID + 1
	if nextTimeID < 1 {
		return result
	}

	// 次の時間帯の開始時刻（lastShift の直後）
	nextStart, err := a.getNextTimeString(c, lastShift.Time.ID)
	if err != nil || nextStart == "" {
		return result
	}

	// 次の次の時間帯の開始時刻（終了時刻として使用）
	nextEnd, err := a.getNextTimeString(c, nextTimeID)
	if err != nil {
		// フォールバック（環境によっては末尾時刻で終端）
		nextEnd = nextStart
	}

	// 次の時間帯に属するメンバーを取得
	nextUsers, err := a.GetUsersByShift(c,
		strconv.Itoa(lastShift.Task.ID),
		strconv.Itoa(lastShift.Year.ID),
		strconv.Itoa(lastShift.Date.ID),
		strconv.Itoa(nextTimeID),
		strconv.Itoa(lastShift.Weather.ID))
	if err != nil {
		return result
	}

	var members []entity.ShiftMember
	for _, user := range nextUsers.Users {
		// マップからGrade/Bureau情報を取得（N+1問題を回避）
		grade := gradeMap[user.GradeID]
		bureau := bureauMap[user.BureauID]

		member := entity.ShiftMember{
			Name:   user.Name,
			Grade:  grade,
			Bureau: bureau,
		}
		members = append(members, member)
	}
	if members == nil {
		members = []entity.ShiftMember{}
	}

	result.STime = nextStart
	result.ETime = nextEnd
	result.Members = members

	return result
}

// -----------------

// 前の時間の文字列を取得するヘルパー関数
func (a *shiftUseCase) getPreviousTimeString(c context.Context, currentTimeID int) (string, error) {
	if currentTimeID <= 1 {
		return "", errors.New("no previous time available")
	}

	row, err := a.timeRep.Find(c, strconv.Itoa(currentTimeID-1))
	if err != nil {
		return "", err
	}

	var time entity.Time
	err = row.Scan(&time.ID, &time.Time, &time.CreatedAt, &time.UpdatedAt)
	if err != nil {
		return "", err
	}

	return time.Time, nil
}

// 次の時間の文字列を取得するヘルパー関数
func (a *shiftUseCase) getNextTimeString(c context.Context, currentTimeID int) (string, error) {
	row, err := a.timeRep.Find(c, strconv.Itoa(currentTimeID+1))
	if err != nil {
		return "", err
	}

	var time entity.Time
	err = row.Scan(&time.ID, &time.Time, &time.CreatedAt, &time.UpdatedAt)
	if err != nil {
		return "", err
	}

	return time.Time, nil
}

func (a *shiftUseCase) GetMaxID(c context.Context) (int, error) {
	maxID := 0

	// クエリー実行
	row, err := a.rep.MaxID(c)
	if err != nil {
		return 0, err
	}

	err = row.Scan(
		&maxID,
	)
	if err != nil {
		return maxID, err
	}

	return maxID, nil
}

func (u *shiftUseCase) SaveShiftData(ctx context.Context, req entity.ShiftRequest) error {
	// DB保存処理（仮実装）
	// 実際にはリポジトリを通じてDBに保存する
	return nil
}

func (u *shiftUseCase) SendToGAS(ctx context.Context, req entity.ShiftRequest) error {
	var gasData entity.GASShiftData
	gasData.Name = req.Name // ユーザー名を設定（必要に応じて変更）

	for _, shift := range req.Shift {
		var gasShift struct {
			Date     int `json:"date"`
			Contents []struct {
				Row    int  `json:"row"`
				Column int  `json:"column"`
				Value  bool `json:"value"`
			} `json:"contents"`
		}

		gasShift.Date = shift.Date
		for _, content := range shift.Contents {
			// TimeID を基に列番号を計算 (TimeID=1 が G列=7)
			column := content.TimeID

			gasShift.Contents = append(gasShift.Contents, struct {
				Row    int  `json:"row"`
				Column int  `json:"column"`
				Value  bool `json:"value"`
			}{
				Row:    0,      // ユーザーIDを行番号として使用
				Column: column, // 計算した列番号
				Value:  content.IsAttend,
			})
		}

		gasData.Shift = append(gasData.Shift, gasShift)
	}

	// JSONに変換して送信
	jsonData, err := json.Marshal(gasData)
	if err != nil {
		log.Printf("failed to marshal GAS data: %v", err)
		return err
	}

	fmt.Printf("Sending data to GAS: %s\n", string(jsonData))

	// GASエンドポイントに送信
	url := "https://script.google.com/macros/s/AKfycbxjhHW10LZPgfvnUD2wzPqECq-k49fFt02ggx5RGivfhdenMvtFnSFKEdLSO37QVGQh/exec" // GASのエンドポイントURLを設定
	reqBody := bytes.NewBuffer(jsonData)
	httpReq, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		log.Printf("failed to create request: %v", err)
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		log.Printf("failed to send request to GAS: %v", err)
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		log.Printf("GAS returned non-OK status: %d", resp.StatusCode)
		return fmt.Errorf("GAS returned non-OK status: %d", resp.StatusCode)
	}

	log.Println("Data successfully sent to GAS")
	return nil
}

// GASからのシフト変更通知を受けてDBを更新
func (u *shiftUseCase) UpdateShiftsFromGAS(ctx context.Context, req entity.ShiftChangeRequest) error {
	// N+1問題対策: 事前に必要なユーザー名とタスク名を収集
	userNameSet := make(map[string]bool)
	taskNameSet := make(map[string]bool)

	for _, change := range req.Changes {
		userNameSet[change.UserName] = true
		taskName := strings.ReplaceAll(change.TaskName, "　", " ")
		taskNameSet[taskName] = true
	}

	// ユーザー名のリストを作成
	userNames := make([]string, 0, len(userNameSet))
	for name := range userNameSet {
		userNames = append(userNames, name)
	}

	// タスク名のリストを作成
	taskNames := make([]string, 0, len(taskNameSet))
	for name := range taskNameSet {
		taskNames = append(taskNames, name)
	}

	// 一括でユーザーを取得してマップ化
	userMap := make(map[string]entity.User)
	if len(userNames) > 0 {
		userRows, err := u.userRep.FindByNames(ctx, userNames)
		if err != nil {
			return errors.Wrap(err, "ユーザー一括取得失敗")
		}
		defer func() { _ = userRows.Close() }()

		for userRows.Next() {
			var user entity.User
			var slackUserID sql.NullString
			if err := userRows.Scan(&user.ID, &user.Name, &user.Mail, &user.GradeID, &user.DepartmentID, &user.BureauID, &user.RoleID, &user.StudentNumber, &user.Tel, &user.Password, &user.CreatedAt, &user.UpdatedAt, &slackUserID); err != nil {
				continue
			}
			if slackUserID.Valid {
				user.SlackUserID = slackUserID.String
			}
			userMap[user.Name] = user
		}
	}

	// 一括でタスクを取得してマップ化
	taskMap := make(map[string]entity.Task)
	if len(taskNames) > 0 {
		taskRows, err := u.taskRep.FindByNames(ctx, taskNames)
		if err != nil {
			return errors.Wrap(err, "タスク一括取得失敗")
		}
		defer func() { _ = taskRows.Close() }()

		for taskRows.Next() {
			var task entity.Task
			if err := taskRows.Scan(&task.ID, &task.Task, &task.PlaceID, &task.Url, &task.BureauID, &task.MaxMember, &task.Color, &task.Remark, &task.YearID, &task.CreatedAt, &task.UpdatedAt); err != nil {
				continue
			}
			taskMap[task.Task] = task
		}
	}

	for _, change := range req.Changes {
		// 年からYearIDを取得
		yearID := strings.ReplaceAll(strconv.Itoa(change.YearID), " ", "")
		yearID = strings.ReplaceAll(yearID, "　", "")

		// 時刻からTimeIDを取得
		timeID := strings.ReplaceAll(strconv.Itoa(change.TimeID), " ", "")
		timeID = strings.ReplaceAll(timeID, "　", "")

		// 日付名からDateIDを取得
		date := strings.ReplaceAll(change.Date, " ", "")
		date = strings.ReplaceAll(date, "　", "")
		var dateID string
		switch date {
		case "準備日":
			dateID = "1"
		case "1日目":
			dateID = "2"
		case "2日目":
			dateID = "3"
		case "片付け日":
			dateID = "4"
		default:
			return errors.New("不正な日付名: " + change.Date)
		}

		// 天気からWeatherIDを取得
		weather := strings.ReplaceAll(change.Weather, " ", "")
		weather = strings.ReplaceAll(weather, "　", "")
		var weatherID string
		switch weather {
		case "晴れ":
			weatherID = "1"
		case "雨":
			weatherID = "2"
		case "none":
			weatherID = "3"
		default:
			return errors.New("不正な天気: " + change.Weather)
		}

		// マップからユーザーを取得（N+1問題を回避）
		user, exists := userMap[change.UserName]
		if !exists {
			return errors.Errorf("ユーザーが存在しません: %v", change.UserName)
		}
		userID := strconv.Itoa(user.ID)

		// マップからタスクを取得（N+1問題を回避）
		taskName := strings.ReplaceAll(change.TaskName, "　", " ")
		task, exists := taskMap[taskName]
		if !exists {
			// タスクが存在しない場合は新規作成
			name := change.TaskName
			placeID := "1"
			url := ""
			bureauID := "1"
			maxMember := "1"
			color := "000000"
			remark := ""
			createErr := u.taskRep.Create(ctx, name, placeID, url, bureauID, maxMember, color, remark, yearID)
			if createErr != nil {
				return errors.Wrapf(createErr, "タスク新規作成失敗: %v", change.TaskName)
			}
			// 新規作成したタスクを再取得してマップに追加
			taskRow, _ := u.taskRep.FindByName(ctx, change.TaskName)
			if err := taskRow.Scan(&task.ID, &task.Task, &task.PlaceID, &task.Url, &task.BureauID, &task.MaxMember, &task.Color, &task.Remark, &task.YearID, &task.CreatedAt, &task.UpdatedAt); err != nil {
				return errors.Wrapf(err, "タスク再取得失敗: %v", change.TaskName)
			}
			taskMap[task.Task] = task
		}
		taskID := strconv.Itoa(task.ID)

		// 4. 既存シフトがあるか確認
		existRow, err := u.rep.FindByUnique(ctx, userID, dateID, timeID, weatherID)
		if err != nil {
			return errors.Wrapf(err, "シフト検索失敗: user=%s, date=%s", user.Name, dateID)
		}
		var existShift entity.ShiftAdmin
		dateIDInt, err := strconv.Atoi(dateID)
		if err != nil {
			return errors.Wrapf(err, "dateIDパース失敗: %v", dateID)
		}
		scanErr := existRow.Scan(&existShift.ID, &existShift.TaskID, &existShift.UserID, &existShift.YearID, &existShift.DateID, &existShift.TimeID, &existShift.WeatherID, &existShift.IsAttendance, &existShift.CreatedAt, &existShift.UpdatedAt)
		if scanErr != nil && !errors.Is(scanErr, sql.ErrNoRows) {
			// 未存在以外のDBエラーを新規作成へフォールバックさせない（重複生成防止）
			return errors.Wrapf(scanErr, "シフト既存確認失敗: user=%s, date=%s", user.Name, dateID)
		}
		if scanErr == nil && existShift.ID != 0 {
			// 既存があれば更新（タスクが変更された場合のみaction_logに記録）
			oldTaskID := existShift.TaskID
			newTaskID, err := strconv.Atoi(taskID)
			if err != nil {
				return errors.Wrapf(err, "taskIDパース失敗: %v", taskID)
			}
			if oldTaskID != newTaskID {
				// タスクが変更された場合
				oldTaskRow, _ := u.taskRep.Find(ctx, strconv.Itoa(oldTaskID))
				var oldTask entity.Task
				oldTaskName := "（不明）"
				if oldTaskRow != nil {
					if err := oldTaskRow.Scan(&oldTask.ID, &oldTask.Task, &oldTask.PlaceID, &oldTask.Url, &oldTask.BureauID, &oldTask.MaxMember, &oldTask.Color, &oldTask.Remark, &oldTask.YearID, &oldTask.CreatedAt, &oldTask.UpdatedAt); err == nil {
						oldTaskName = oldTask.Task
					}
				}

				newTaskName := task.Task
				if newTaskName == "" {
					newTaskName = "（不明）"
				}

				// action_logに記録
				diffPayload := map[string]interface{}{
					"changes": []map[string]string{
						{"field": "task_name", "old": oldTaskName, "new": newTaskName},
					},
				}
				if u.actionLogRepo != nil {
					if logErr := u.actionLogRepo.Create(ctx, existShift.ID, user.ID, dateIDInt, "UPDATE", diffPayload); logErr != nil {
						log.Printf("action_log記録失敗(UPDATE): %v", logErr)
					}
				}
			}
			isAttendance := false
			if err := u.rep.Update(ctx, strconv.Itoa(existShift.ID), taskID, userID, strconv.Itoa(existShift.YearID), dateID, timeID, weatherID, strconv.FormatBool(isAttendance)); err != nil {
				return errors.Wrapf(err, "シフト更新失敗: %v", existShift.ID)
			}
		} else {
			// なければ新規作成（RETURNING idで確実にIDを取得）
			isAttendance := false
			newShiftID, err := u.rep.CreateAndReturnID(ctx, taskID, userID, yearID, dateID, timeID, weatherID, strconv.FormatBool(isAttendance))
			if err != nil {
				return errors.Wrapf(err, "シフト新規作成失敗: user=%s, date=%s", user.Name, dateID)
			}
			// action_logに記録
			newTaskName := task.Task
			if newTaskName == "" {
				newTaskName = "（新規）"
			}
			diffPayload := map[string]interface{}{
				"changes": []map[string]string{
					{"field": "task_name", "old": "なし", "new": newTaskName},
				},
			}
			if u.actionLogRepo != nil {
				if logErr := u.actionLogRepo.Create(ctx, newShiftID, user.ID, dateIDInt, "CREATE", diffPayload); logErr != nil {
					log.Printf("action_log記録失敗(CREATE): %v", logErr)
				}
			}
		}
	}
	return nil
}

// loadGradeMap は全Gradeを取得してマップ化（N+1問題対策）
func (a *shiftUseCase) loadGradeMap(ctx context.Context) (map[int]string, error) {
	gradeMap := make(map[int]string)
	rows, err := a.gradeRep.All(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var grade entity.Grade
		if err := rows.Scan(&grade.ID, &grade.Grade, &grade.CreatedAt, &grade.UpdatedAt); err != nil {
			continue
		}
		gradeMap[grade.ID] = grade.Grade
	}
	return gradeMap, nil
}

// loadBureauMap は全Bureauを取得してマップ化（N+1問題対策）
func (a *shiftUseCase) loadBureauMap(ctx context.Context) (map[int]string, error) {
	bureauMap := make(map[int]string)
	rows, err := a.bureauRep.All(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var bureau entity.Bureau
		if err := rows.Scan(&bureau.ID, &bureau.Bureau, &bureau.Color, &bureau.CreatedAt, &bureau.UpdatedAt); err != nil {
			continue
		}
		bureauMap[bureau.ID] = bureau.Bureau
	}
	return bureauMap, nil
}
