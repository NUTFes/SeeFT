package usecase

import (
	"context"
	"sort"
	"strconv"

	rep "github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/pkg/errors"
)

type shiftUseCase struct {
  rep rep.ShiftRepository
  taskRep rep.TaskRepository
  userRep rep.UserRepository
  yearRep rep.YearRepository
  dateRep rep.DateRepository
  timeRep rep.TimeRepository
  weatherRep rep.WeatherRepository
	placeRep rep.PlaceRepository
	gradeRep rep.GradeRepository
	bureauRep rep.BureauRepository
}

type ShiftUseCase interface {
  GetShifts(context.Context) ([]entity.Shift, error)
  GetShiftByID(context.Context, string) (entity.Shift, error)
  GetShiftsByUser(context.Context, string) ([]entity.Shift, error)
  GetShiftsByUserAndDateAndWeather(context.Context, string, string, string) ([]entity.Shift, error)
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
}

func NewShiftUseCase(
  rep rep.ShiftRepository, 
  taskRep rep.TaskRepository,
  userRep rep.UserRepository,
  yearRep rep.YearRepository,
  dateRep rep.DateRepository,
  timeRep rep.TimeRepository,
  weatherRep rep.WeatherRepository,
	placeRep rep.PlaceRepository,
	gradeRep rep.GradeRepository,
	bureauRep rep.BureauRepository) ShiftUseCase {
  return &shiftUseCase{rep, taskRep, userRep, yearRep, dateRep, timeRep, weatherRep, placeRep, gradeRep, bureauRep}
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

func (a *shiftUseCase) GetShiftCardsByUserAndDateAndWeather(c context.Context, userID string, dateID string, weatherID string) ([]entity.ShiftCard, error) {
	// 1. 指定したユーザー・日付・天気のシフトを取得
	shifts, err := a.GetShiftsByUserAndDateAndWeather(c, userID, dateID, weatherID)
	if err != nil {
		return nil, err
	}

	// 2. タスクIDと連続するTimeIDでグループ化
	taskGroups := make(map[int][]entity.Shift)
	for _, shift := range shifts {
		taskGroups[shift.Task.ID] = append(taskGroups[shift.Task.ID], shift)
	}

	var shiftCards []entity.ShiftCard

	// 3. 各タスクグループを処理
	for _, taskShifts := range taskGroups {
		// 時間順にソート
		sort.Sort(ByTime(taskShifts))
		
		// 連続するTimeIDでグループ化
		continuousGroups := a.groupContinuousShifts(taskShifts)
		
		for _, group := range continuousGroups {
			if len(group) == 0 {
				continue
			}

			// 最後のグループの終了時刻を取得
			endTime, err := a.getNextTimeString(c, group[len(group)-1].Time.ID)
			if err != nil {
				endTime = group[len(group)-1].Time.Time // フォールバック
			}

			// ShiftCardを作成
			shiftCard := entity.ShiftCard{
				TaskName:  group[0].Task.Task,
				StartTime: group[0].Time.Time,
				EndTime:   endTime,
				Place:     group[0].Task.Place,
				Url:       group[0].Task.Url,
			}

			// ShiftMembersを作成（15分ごと）
			var shiftMembers []entity.ShiftMembers
			for i, shift := range group {
				// その時間のメンバーを取得
				shiftUsers, err := a.GetUsersByShift(c, 
					strconv.Itoa(shift.Task.ID), 
					strconv.Itoa(shift.Year.ID), 
					strconv.Itoa(shift.Date.ID), 
					strconv.Itoa(shift.Time.ID), 
					strconv.Itoa(shift.Weather.ID))
				if err != nil {
					continue
				}

				var members []entity.ShiftMember
				for _, user := range shiftUsers.Users {
					// Grade情報を取得
					var grade entity.Grade
					gradeRow, err := a.gradeRep.Find(c, strconv.Itoa(user.GradeID))
					if err != nil {
						// エラー時はデフォルト値を設定
						grade.Grade = "不明"
					} else {
						err = gradeRow.Scan(&grade.ID, &grade.Grade, &grade.CreatedAt, &grade.UpdatedAt)
						if err != nil {
							// Scanエラー時もデフォルト値を設定
							grade.Grade = "不明"
						}
					}

					// Bureau情報を取得
					var bureau entity.Bureau
					bureauRow, err := a.bureauRep.Find(c, strconv.Itoa(user.BureauID))
					if err != nil {
						// エラー時はデフォルト値を設定
						bureau.Bureau = "不明"
					} else {
						err = bureauRow.Scan(&bureau.ID, &bureau.Bureau, &bureau.Color, &bureau.CreatedAt, &bureau.UpdatedAt)
						if err != nil {
							// Scanエラー時もデフォルト値を設定
							bureau.Bureau = "不明"
						}
					}

					member := entity.ShiftMember{
						Name:   user.Name,
						Grade:  grade.Grade,
						Bureau: bureau.Bureau,
					}
					members = append(members, member)
				}

				// 終了時刻を計算
				var eTime string
				if i < len(group)-1 {
					// 次のシフトの開始時刻を使用
					eTime = group[i+1].Time.Time
				} else {
					// 最後のシフトの場合、次の時間を取得
					nextTime, err := a.getNextTimeString(c, shift.Time.ID)
					if err != nil {
						eTime = shift.Time.Time // フォールバック
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

			// BeforeMembers（前の時間のメンバー）
			if len(group) > 0 {
				firstShift := group[0]
				if firstShift.Time.ID > 1 {
					// 前の時間の情報を取得
					beforeTimeString, err := a.getPreviousTimeString(c, firstShift.Time.ID)
					if err == nil {
						// 前の時間のメンバーを取得
						beforeUsers, err := a.GetUsersByShift(c,
							strconv.Itoa(firstShift.Task.ID),
							strconv.Itoa(firstShift.Year.ID),
							strconv.Itoa(firstShift.Date.ID),
							strconv.Itoa(firstShift.Time.ID-1),
							strconv.Itoa(firstShift.Weather.ID))
						if err == nil {
							var beforeMembers []entity.ShiftMember
							for _, user := range beforeUsers.Users {
								// Grade情報を取得
								var grade entity.Grade
								gradeRow, err := a.gradeRep.Find(c, strconv.Itoa(user.GradeID))
								if err != nil {
									// エラー時はデフォルト値を設定
									grade.Grade = "不明"
								} else {
									err = gradeRow.Scan(&grade.ID, &grade.Grade, &grade.CreatedAt, &grade.UpdatedAt)
									if err != nil {
										// Scanエラー時もデフォルト値を設定
										grade.Grade = "不明"
									}
								}

								// Bureau情報を取得
								var bureau entity.Bureau
								bureauRow, err := a.bureauRep.Find(c, strconv.Itoa(user.BureauID))
								if err != nil {
									// エラー時はデフォルト値を設定
									bureau.Bureau = "不明"
								} else {
									err = bureauRow.Scan(&bureau.ID, &bureau.Bureau, &bureau.Color, &bureau.CreatedAt, &bureau.UpdatedAt)
									if err != nil {
										// Scanエラー時もデフォルト値を設定
										bureau.Bureau = "不明"
									}
								}

								member := entity.ShiftMember{
									Name:   user.Name,
									Grade:  grade.Grade,
									Bureau: bureau.Bureau,
								}
								beforeMembers = append(beforeMembers, member)
							}
							shiftCard.BeforeMembers = entity.ShiftMembers{
								STime:   beforeTimeString,
								ETime:   firstShift.Time.Time,
								Members: beforeMembers,
							}
						}
					}
				}
			}

			// AfterMembers（後の時間のメンバー）
			if len(group) > 0 {
				lastShift := group[len(group)-1]
				// 後の時間の情報を取得
				afterTimeString, err := a.getNextTimeString(c, lastShift.Time.ID)
				if err == nil {
					// 後の時間のメンバーを取得
					afterUsers, err := a.GetUsersByShift(c,
						strconv.Itoa(lastShift.Task.ID),
						strconv.Itoa(lastShift.Year.ID),
						strconv.Itoa(lastShift.Date.ID),
						strconv.Itoa(lastShift.Time.ID+1),
						strconv.Itoa(lastShift.Weather.ID))
					if err == nil {
						var afterMembers []entity.ShiftMember
						for _, user := range afterUsers.Users {
							// Grade情報を取得
							var grade entity.Grade
							gradeRow, err := a.gradeRep.Find(c, strconv.Itoa(user.GradeID))
							if err != nil {
								// エラー時はデフォルト値を設定
								grade.Grade = "不明"
							} else {
								err = gradeRow.Scan(&grade.ID, &grade.Grade, &grade.CreatedAt, &grade.UpdatedAt)
								if err != nil {
									// Scanエラー時もデフォルト値を設定
									grade.Grade = "不明"
								}
							}

							// Bureau情報を取得
							var bureau entity.Bureau
							bureauRow, err := a.bureauRep.Find(c, strconv.Itoa(user.BureauID))
							if err != nil {
								// エラー時はデフォルト値を設定
								bureau.Bureau = "不明"
							} else {
								err = bureauRow.Scan(&bureau.ID, &bureau.Bureau, &bureau.Color, &bureau.CreatedAt, &bureau.UpdatedAt)
								if err != nil {
									// Scanエラー時もデフォルト値を設定
									bureau.Bureau = "不明"
								}
							}

							member := entity.ShiftMember{
								Name:   user.Name,
								Grade:  grade.Grade,
								Bureau: bureau.Bureau,
							}
							afterMembers = append(afterMembers, member)
						}
						
						// 後の時間の終了時刻を計算
						afterEndTime, err := a.getNextTimeString(c, lastShift.Time.ID+1)
						if err != nil {
							afterEndTime = afterTimeString // フォールバック
						}
						
						shiftCard.AfterMembers = entity.ShiftMembers{
							STime:   afterTimeString,
							ETime:   afterEndTime,
							Members: afterMembers,
						}
					}
				}
			}

			shiftCards = append(shiftCards, shiftCard)
		}
	}

	return shiftCards, nil
}

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

// 連続するTimeIDでシフトをグループ化するヘルパー関数
func (a *shiftUseCase) groupContinuousShifts(shifts []entity.Shift) [][]entity.Shift {
	if len(shifts) == 0 {
		return nil
	}

	var groups [][]entity.Shift
	currentGroup := []entity.Shift{shifts[0]}

	for i := 1; i < len(shifts); i++ {
		// 連続するTimeIDかチェック
		if shifts[i].Time.ID == shifts[i-1].Time.ID+1 {
			currentGroup = append(currentGroup, shifts[i])
		} else {
			// 連続していない場合、新しいグループを開始
			groups = append(groups, currentGroup)
			currentGroup = []entity.Shift{shifts[i]}
		}
	}
	
	// 最後のグループを追加
	groups = append(groups, currentGroup)
	
	return groups
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

// import '../entity/entity.dart';
// import './repository/repository.dart';

// abstract class ShiftUsecase {
//   Future<List<Shift>> getShiftsByUser(ctx, User req);
//   Future<List<Shift>> getShiftsByUserAndDateAndWeather(ctx, Shift req);
//   Future<List<Shift>> getShiftsByYearAndDateAndWeather(ctx, Shift req);
// }

// class ShiftUsecaseImpl implements ShiftUsecase {
//   ShiftRepository shiftRepository;

//   ShiftUsecaseImpl(this.shiftRepository);

//   @override
//   Future<List<Shift>> getShiftsByUser(ctx, User req) async {
//     List<Shift> list = await shiftRepository.getShiftsByUser(ctx, req);
//     return list;
//   }

//   @override
//   Future<List<Shift>> getShiftsByUserAndDateAndWeather(ctx, Shift req) async {
//     List<Shift> list = await shiftRepository.getShiftsByUserAndDateAndWeather(ctx, req);
//     return list;
//   }

//   @override
//   Future<List<Shift>> getShiftsByYearAndDateAndWeather(ctx, Shift req) async {
//     List<Shift> list = await shiftRepository.getShiftsByYearAndDateAndWeather(ctx, req);
//     return list;
//   }
// }
