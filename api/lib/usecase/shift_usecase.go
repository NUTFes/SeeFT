package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/NUTFes/SeeFT/api/lib/entity"
	rep "github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/pkg/errors"
)

type shiftUseCase struct {
	rep        rep.ShiftRepository
	taskRep    rep.TaskRepository
	userRep    rep.UserRepository
	yearRep    rep.YearRepository
	dateRep    rep.DateRepository
	timeRep    rep.TimeRepository
	weatherRep rep.WeatherRepository
	placeRep   rep.PlaceRepository
}

type ShiftUseCase interface {
	GetShifts(context.Context) ([]entity.Shift, error)
	GetShiftByID(context.Context, string) (entity.Shift, error)
	GetShiftsByUser(context.Context, string) ([]entity.Shift, error)
	GetShiftsByUserAndDateAndWeather(context.Context, string, string, string) ([]entity.Shift, error)
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
	taskRep rep.TaskRepository,
	userRep rep.UserRepository,
	yearRep rep.YearRepository,
	dateRep rep.DateRepository,
	timeRep rep.TimeRepository,
	weatherRep rep.WeatherRepository,
	placeRep rep.PlaceRepository) ShiftUseCase {
	return &shiftUseCase{rep, taskRep, userRep, yearRep, dateRep, timeRep, weatherRep, placeRep}
}

var TaskID, UserID, YearID, DateID, TimeID, WeatherID, PlaceID string

// 時間でソート
type ByTime []entity.Shift

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].Time.ID < a[j].Time.ID }

func (a *shiftUseCase) GetShifts(c context.Context) ([]entity.Shift, error) {
	shift := entity.Shift{}
	var shifts []entity.Shift
	place := entity.Place{}

	// クエリー実行
	rows, err := a.rep.All(c)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&shift.ID,
			&TaskID,
			&UserID,
			&YearID,
			&DateID,
			&TimeID,
			&WeatherID,
			&shift.IsAttendance,
			&shift.CreatedAt,
			&shift.UpdatedAt,
		)

		row, err := a.taskRep.Find(c, TaskID)
		err = row.Scan(
			&shift.Task.ID,
			&shift.Task.Task,
			&PlaceID,
			&shift.Task.Url,
			&shift.Task.BureauID,
			&shift.Task.MaxMember,
			&shift.Task.Color,
			&shift.Task.Remark,
			&shift.Task.YearID,
			&shift.Task.CreatedAt,
			&shift.Task.UpdatedAt,
		)

		row, err = a.placeRep.Find(c, PlaceID)
		err = rows.Scan(
			&place.ID,
			&shift.Task.Place,
			&place.Remark,
			&place.CreatedAt,
			&place.UpdatedAt,
		)

		row, err = a.userRep.Find(c, UserID)
		err = row.Scan(
			&shift.User.ID,
			&shift.User.Name,
			&shift.User.Mail,
			&shift.User.GradeID,
			&shift.User.DepartmentID,
			&shift.User.BureauID,
			&shift.User.RoleID,
			&shift.User.StudentNumber,
			&shift.User.Tel,
			&shift.User.Password,
			&shift.User.CreatedAt,
			&shift.User.UpdatedAt,
		)

		row, err = a.yearRep.Find(c, YearID)
		err = row.Scan(
			&shift.Year.ID,
			&shift.Year.Year,
			&shift.Year.CreatedAt,
			&shift.Year.UpdatedAt,
		)

		row, err = a.dateRep.Find(c, DateID)
		err = row.Scan(
			&shift.Date.ID,
			&shift.Date.YearID,
			&shift.Date.Name,
			&shift.Date.Date,
			&shift.Date.CreatedAt,
			&shift.Date.UpdatedAt,
		)

		row, err = a.timeRep.Find(c, TimeID)
		err = row.Scan(
			&shift.Time.ID,
			&shift.Time.Time,
			&shift.Time.CreatedAt,
			&shift.Time.UpdatedAt,
		)

		row, err = a.weatherRep.Find(c, WeatherID)
		err = row.Scan(
			&shift.Weather.ID,
			&shift.Weather.Weather,
			&shift.Weather.CreatedAt,
			&shift.Weather.UpdatedAt,
		)

		if err != nil {
			return nil, errors.Wrapf(err, "cannot connect SQL")
		}

		shifts = append(shifts, shift)
	}
	return shifts, nil
}

func (a *shiftUseCase) GetShiftByID(c context.Context, id string) (entity.Shift, error) {
	var shift entity.Shift
	place := entity.Place{}

	row, err := a.rep.Find(c, id)
	err = row.Scan(
		&shift.ID,
		&TaskID,
		&UserID,
		&YearID,
		&DateID,
		&TimeID,
		&WeatherID,
		&shift.IsAttendance,
		&shift.CreatedAt,
		&shift.UpdatedAt,
	)

	row, err = a.taskRep.Find(c, TaskID)
	err = row.Scan(
		&shift.Task.ID,
		&shift.Task.Task,
		&PlaceID,
		&shift.Task.Url,
		&shift.Task.BureauID,
		&shift.Task.MaxMember,
		&shift.Task.Color,
		&shift.Task.Remark,
		&shift.Task.YearID,
		&shift.Task.CreatedAt,
		&shift.Task.UpdatedAt,
	)

	row, err = a.placeRep.Find(c, PlaceID)
	err = row.Scan(
		&place.ID,
		&shift.Task.Place,
		&place.Remark,
		&place.CreatedAt,
		&place.UpdatedAt,
	)

	row, err = a.userRep.Find(c, UserID)
	err = row.Scan(
		&shift.User.ID,
		&shift.User.Name,
		&shift.User.Mail,
		&shift.User.GradeID,
		&shift.User.DepartmentID,
		&shift.User.BureauID,
		&shift.User.RoleID,
		&shift.User.StudentNumber,
		&shift.User.Tel,
		&shift.User.Password,
		&shift.User.CreatedAt,
		&shift.User.UpdatedAt,
	)

	row, err = a.yearRep.Find(c, YearID)
	err = row.Scan(
		&shift.Year.ID,
		&shift.Year.Year,
		&shift.Year.CreatedAt,
		&shift.Year.UpdatedAt,
	)

	row, err = a.dateRep.Find(c, DateID)
	err = row.Scan(
		&shift.Date.ID,
		&shift.Date.YearID,
		&shift.Date.Name,
		&shift.Date.Date,
		&shift.Date.CreatedAt,
		&shift.Date.UpdatedAt,
	)

	row, err = a.timeRep.Find(c, TimeID)
	err = row.Scan(
		&shift.Time.ID,
		&shift.Time.Time,
		&shift.Time.CreatedAt,
		&shift.Time.UpdatedAt,
	)

	row, err = a.weatherRep.Find(c, WeatherID)
	err = row.Scan(
		&shift.Weather.ID,
		&shift.Weather.Weather,
		&shift.Weather.CreatedAt,
		&shift.Weather.UpdatedAt,
	)

	if err != nil {
		return shift, err
	}
	return shift, nil
}

func (a *shiftUseCase) GetShiftsByUser(c context.Context, id string) ([]entity.Shift, error) {

	shift := entity.Shift{}
	var shifts []entity.Shift
	place := entity.Place{}

	// クエリー実行
	rows, err := a.rep.User(c, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&shift.ID,
			&TaskID,
			&UserID,
			&YearID,
			&DateID,
			&TimeID,
			&WeatherID,
			&shift.IsAttendance,
			&shift.CreatedAt,
			&shift.UpdatedAt,
		)

		row, err := a.taskRep.Find(c, TaskID)
		err = row.Scan(
			&shift.Task.ID,
			&shift.Task.Task,
			&PlaceID,
			&shift.Task.Url,
			&shift.Task.BureauID,
			&shift.Task.MaxMember,
			&shift.Task.Color,
			&shift.Task.Remark,
			&shift.Task.YearID,
			&shift.Task.CreatedAt,
			&shift.Task.UpdatedAt,
		)

		row, err = a.placeRep.Find(c, PlaceID)
		err = row.Scan(
			&place.ID,
			&shift.Task.Place,
			&place.Remark,
			&place.CreatedAt,
			&place.UpdatedAt,
		)

		row, err = a.userRep.Find(c, UserID)
		err = row.Scan(
			&shift.User.ID,
			&shift.User.Name,
			&shift.User.Mail,
			&shift.User.GradeID,
			&shift.User.DepartmentID,
			&shift.User.BureauID,
			&shift.User.RoleID,
			&shift.User.StudentNumber,
			&shift.User.Tel,
			&shift.User.Password,
			&shift.User.CreatedAt,
			&shift.User.UpdatedAt,
		)

		row, err = a.yearRep.Find(c, YearID)
		err = row.Scan(
			&shift.Year.ID,
			&shift.Year.Year,
			&shift.Year.CreatedAt,
			&shift.Year.UpdatedAt,
		)

		row, err = a.dateRep.Find(c, DateID)
		err = row.Scan(
			&shift.Date.ID,
			&shift.Date.YearID,
			&shift.Date.Name,
			&shift.Date.Date,
			&shift.Date.CreatedAt,
			&shift.Date.UpdatedAt,
		)

		row, err = a.timeRep.Find(c, TimeID)
		err = row.Scan(
			&shift.Time.ID,
			&shift.Time.Time,
			&shift.Time.CreatedAt,
			&shift.Time.UpdatedAt,
		)

		row, err = a.weatherRep.Find(c, WeatherID)
		err = row.Scan(
			&shift.Weather.ID,
			&shift.Weather.Weather,
			&shift.Weather.CreatedAt,
			&shift.Weather.UpdatedAt,
		)

		if err != nil {
			return nil, errors.Wrapf(err, "cannot connect SQL")
		}

		shifts = append(shifts, shift)
	}
	sort.Sort(ByTime(shifts))
	return shifts, nil
}

func (a *shiftUseCase) GetShiftsByUserAndDateAndWeather(c context.Context, id string, date string, weather string) ([]entity.Shift, error) {

	shift := entity.Shift{}
	var shifts []entity.Shift
	place := entity.Place{}

	// クエリー実行
	rows, err := a.rep.UserAndDateAndWeather(c, id, date, weather)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&shift.ID,
			&TaskID,
			&UserID,
			&YearID,
			&DateID,
			&TimeID,
			&WeatherID,
			&shift.IsAttendance,
			&shift.CreatedAt,
			&shift.UpdatedAt,
		)

		row, err := a.taskRep.Find(c, TaskID)
		err = row.Scan(
			&shift.Task.ID,
			&shift.Task.Task,
			&PlaceID,
			&shift.Task.Url,
			&shift.Task.BureauID,
			&shift.Task.MaxMember,
			&shift.Task.Color,
			&shift.Task.Remark,
			&shift.Task.YearID,
			&shift.Task.CreatedAt,
			&shift.Task.UpdatedAt,
		)

		row, err = a.placeRep.Find(c, PlaceID)
		err = row.Scan(
			&place.ID,
			&shift.Task.Place,
			&place.Remark,
			&place.CreatedAt,
			&place.UpdatedAt,
		)

		row, err = a.userRep.Find(c, UserID)
		err = row.Scan(
			&shift.User.ID,
			&shift.User.Name,
			&shift.User.Mail,
			&shift.User.GradeID,
			&shift.User.DepartmentID,
			&shift.User.BureauID,
			&shift.User.RoleID,
			&shift.User.StudentNumber,
			&shift.User.Tel,
			&shift.User.Password,
			&shift.User.CreatedAt,
			&shift.User.UpdatedAt,
		)

		row, err = a.yearRep.Find(c, YearID)
		err = row.Scan(
			&shift.Year.ID,
			&shift.Year.Year,
			&shift.Year.CreatedAt,
			&shift.Year.UpdatedAt,
		)

		row, err = a.dateRep.Find(c, DateID)
		err = row.Scan(
			&shift.Date.ID,
			&shift.Date.YearID,
			&shift.Date.Name,
			&shift.Date.Date,
			&shift.Date.CreatedAt,
			&shift.Date.UpdatedAt,
		)

		row, err = a.timeRep.Find(c, TimeID)
		err = row.Scan(
			&shift.Time.ID,
			&shift.Time.Time,
			&shift.Time.CreatedAt,
			&shift.Time.UpdatedAt,
		)

		row, err = a.weatherRep.Find(c, WeatherID)
		err = row.Scan(
			&shift.Weather.ID,
			&shift.Weather.Weather,
			&shift.Weather.CreatedAt,
			&shift.Weather.UpdatedAt,
		)

		if err != nil {
			return nil, errors.Wrapf(err, "cannot connect SQL")
		}

		shifts = append(shifts, shift)
	}
	sort.Sort(ByTime(shifts))
	return shifts, nil
}

func (a *shiftUseCase) GetUsersByShift(c context.Context, task string, year string, date string, time string, weather string) (entity.ShiftUsers, error) {

	users := entity.User{}
	var shiftUsers entity.ShiftUsers

	// クエリー実行
	rows, err := a.rep.Users(c, task, year, date, time, weather)
	if err != nil {
		return shiftUsers, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&UserID,
		)

		row, err := a.userRep.Find(c, UserID)
		err = row.Scan(
			&users.ID,
			&users.Name,
			&users.Mail,
			&users.GradeID,
			&users.DepartmentID,
			&users.BureauID,
			&users.RoleID,
			&users.StudentNumber,
			&users.Tel,
			&users.Password,
			&users.CreatedAt,
			&users.UpdatedAt,
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
	defer rows.Close()

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

	u.rep.Update(c, id, taskID, userID, yearID, dateID, timeID, weatherID, isAttendance)
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
	err := u.rep.Destroy(c, id)
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
	defer rows.Close()

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
	defer rows.Close()

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
				Row:    0, // ユーザーIDを行番号として使用
				Column: column,    // 計算した列番号
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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("GAS returned non-OK status: %d", resp.StatusCode)
		return fmt.Errorf("GAS returned non-OK status: %d", resp.StatusCode)
	}

	log.Println("Data successfully sent to GAS")
	return nil
}

// GASからのシフト変更通知を受けてDBを更新
func (u *shiftUseCase) UpdateShiftsFromGAS(ctx context.Context, req entity.ShiftChangeRequest) error {
	for _, change := range req.Changes {
		// 1. ユーザー名からUserID取得
		userRow, err := u.userRep.FindByName(ctx, change.Row) // Rowはユーザー名が入っている前提
		var user entity.User
		if err := userRow.Scan(&user.ID, &user.Name, &user.Mail, &user.GradeID, &user.DepartmentID, &user.BureauID, &user.RoleID, &user.StudentNumber, &user.Tel, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
			// ユーザーがいなければ新規作成
			if err.Error() == "sql: no rows in result set" {
					// 必要な情報は仮値でOK（必要に応じて修正）
					name := change.Row
					mail := ""
					gradeID := "1"
					departmentID := "1"
					bureauID := "1"
					roleID := "1"
					studentNumber := "0"
					tel := ""
					password := ""
					createErr := u.userRep.Create(ctx, name, mail, gradeID, departmentID, bureauID, roleID, studentNumber, tel, password)
					if createErr != nil {
							return errors.Wrapf(createErr, "ユーザー新規作成失敗: %v", change.Row)
					}
					// 再取得
					userRow, err = u.userRep.FindByName(ctx, change.Row)
					if err := userRow.Scan(&user.ID, &user.Name, &user.Mail, &user.GradeID, &user.DepartmentID, &user.BureauID, &user.RoleID, &user.StudentNumber, &user.Tel, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
							return errors.Wrapf(err, "ユーザー再取得失敗: %v", change.Row)
					}
			} else {
					return errors.Wrapf(err, "ユーザー取得失敗: %v", change.Row)
			}
	}
		userID := strconv.Itoa(user.ID)

		// 2. Task名からTaskID取得
		taskRow, err := u.taskRep.FindByName(ctx, change.Value)
		var task entity.Task
		if err := taskRow.Scan(&task.ID, &task.Task, &task.PlaceID, &task.Url, &task.BureauID, &task.MaxMember, &task.Color, &task.Remark, &task.YearID, &task.CreatedAt, &task.UpdatedAt); err != nil {
			if err.Error() == "sql: no rows in result set" {
					// 必要な情報は仮値でOK（必要に応じて修正）
					name := change.Value
					placeID := "1"
					url := ""
					bureauID := "1"
					maxMember := "0"
					color := ""
					remark := ""
					yearID := "43"
					createErr := u.taskRep.Create(ctx, name, placeID, url, bureauID, maxMember, color, remark, yearID)
					if createErr != nil {
							return errors.Wrapf(createErr, "タスク新規作成失敗: %v", change.Value)
					}
					// 再取得
					taskRow, err = u.taskRep.FindByName(ctx, change.Value)
					if err := taskRow.Scan(&task.ID, &task.Task, &task.PlaceID, &task.Url, &task.BureauID, &task.MaxMember, &task.Color, &task.Remark, &task.YearID, &task.CreatedAt, &task.UpdatedAt); err != nil {
							return errors.Wrapf(err, "タスク再取得失敗: %v", change.Value)
					}
			} else {
					return errors.Wrapf(err, "タスク取得失敗: %v", change.Value)
			}
	}
		taskID := strconv.Itoa(task.ID)

		// 3. DateIDはSheetName、WeatherIDは1で固定
		// dateID := change.SheetName
		dateID := "1"
		weatherID := "1"
		timeID := strconv.Itoa(change.Column)

		// 4. 既存シフトがあるか確認
		existRow, err := u.rep.FindByUnique(ctx, taskID, userID, dateID, timeID, weatherID)
		var existShift entity.ShiftAdmin
		err = existRow.Scan(&existShift.ID, &existShift.TaskID, &existShift.UserID, &existShift.YearID, &existShift.DateID, &existShift.TimeID, &existShift.WeatherID, &existShift.IsAttendance, &existShift.CreatedAt, &existShift.UpdatedAt)
		if err == nil && existShift.ID != 0 {
			// 既存があれば更新
			isAttendance := false
			u.rep.Update(ctx, strconv.Itoa(existShift.ID), taskID, userID, strconv.Itoa(existShift.YearID), dateID, timeID, weatherID, strconv.FormatBool(isAttendance))
		} else {
			// なければ新規作成
			yearID := "43" // 必要に応じて適切なYearIDを設定
			isAttendance := false
			u.rep.Create(ctx, taskID, userID, yearID, dateID, timeID, weatherID, strconv.FormatBool(isAttendance))
		}
	}
	return nil
}