package usecase

import (
	"context"
	"sort"

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
}

func NewShiftUseCase(
	rep rep.ShiftRepository, 
	taskRep rep.TaskRepository,
  userRep rep.UserRepository,
  yearRep rep.YearRepository,
  dateRep rep.DateRepository,
  timeRep rep.TimeRepository,
  weatherRep rep.WeatherRepository) ShiftUseCase {
  return &shiftUseCase{rep, taskRep, userRep, yearRep, dateRep, timeRep, weatherRep}
}

var TaskID, UserID, YearID, DateID, TimeID, WeatherID string

// 時間でソート
type ByTime []entity.Shift
func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].Time.ID < a[j].Time.ID }

func (a *shiftUseCase) GetShifts(c context.Context) ([]entity.Shift, error) {
  shift := entity.Shift{}
  var shifts []entity.Shift

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
			&shift.Task.PlaceID,
			&shift.Task.Url,
			&shift.Task.MaxMember,
			&shift.Task.Color,
			&shift.Task.Remark,
			&shift.Task.YearID,
			&shift.Task.CreatedAt,
			&shift.Task.UpdatedAt,
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
		&shift.Task.PlaceID,
		&shift.Task.Url,
		&shift.Task.MaxMember,
		&shift.Task.Color,
		&shift.Task.Remark,
		&shift.Task.YearID,
		&shift.Task.CreatedAt,
		&shift.Task.UpdatedAt,
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
			&shift.Task.PlaceID,
			&shift.Task.Url,
			&shift.Task.MaxMember,
			&shift.Task.Color,
			&shift.Task.Remark,
			&shift.Task.YearID,
			&shift.Task.CreatedAt,
			&shift.Task.UpdatedAt,
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
			&shift.Task.PlaceID,
			&shift.Task.Url,
			&shift.Task.MaxMember,
			&shift.Task.Color,
			&shift.Task.Remark,
			&shift.Task.YearID,
			&shift.Task.CreatedAt,
			&shift.Task.UpdatedAt,
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
