package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"sync"

	"github.com/NUTFes/SeeFT/api/lib/entity"
	rep "github.com/NUTFes/SeeFT/api/lib/internals/repository"
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
	SaveShiftData(context.Context, entity.ShiftRequest) (error)
	SendToGAS(context.Context, entity.ShiftRequest) (error)
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
    const maxConcurrentRequests = 20
    var wg sync.WaitGroup
    sem := make(chan struct{}, maxConcurrentRequests)

    for _, shift := range req.Shift {
        wg.Add(1)
        sem <- struct{}{}

        go func(shiftData struct {
            Date     int
            Contents []struct {
                TimeID   int
                IsAttend bool
            }
        }) {
            defer wg.Done()
            defer func() { <-sem }()

						 // GAS用データの変換
			var gasData []struct {
					Row    int  `json:"row"`    // 行番号 (timeID)
					Column int  `json:"column"` // 列番号 (userID)
					Value  bool `json:"value"`  // セルの値 (IsAttend)
			}
			
			for _, content := range shiftData.Contents {
				gasData = append(gasData, struct {
					Row    int  `json:"row"`
					Column int  `json:"column"`
					Value  bool `json:"value"`
				}{
					Row:    content.TimeID,   // timeID を行番号として使用
					Column: req.UserID,       // userID を列番号として使用
					Value:  content.IsAttend, // セルの値として IsAttend を使用
				})
			}
			
			// JSONに変換
			jsonData, err := json.Marshal(gasData)
			if err != nil {
					log.Printf("failed to marshal GAS data: %v", err)
					return
			}

			// GASエンドポイントに送信
			url := "https://script.google.com/macros/s/AKfycbw7rQNDPdB5fjKORCLbM9NU8hbOgCqPoiojkTZ4S2wk4t20DI_tSvbHhjL80kDv6lZ2Og/exec"
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
			if err != nil {
					log.Printf("failed to create request: %v", err)
					return
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
					log.Printf("failed to send request to GAS: %v", err)
					return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
					log.Printf("GAS returned non-OK status: %d", resp.StatusCode)
					return
			}
		}(shift)
	}

	wg.Wait()
	return nil
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
