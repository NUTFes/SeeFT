package usecase

import (
	"context"

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
			&shift.Task.Place,
			&shift.Task.Url,
			&shift.Task.Superviser,
			&shift.Task.Color,
			&shift.Task.Notes,
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
			&shift.User.Tel,
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
		&shift.Task.Place,
		&shift.Task.Url,
		&shift.Task.Superviser,
		&shift.Task.Color,
		&shift.Task.Notes,
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
		&shift.User.Tel,
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
			&shift.Task.Place,
			&shift.Task.Url,
			&shift.Task.Superviser,
			&shift.Task.Color,
			&shift.Task.Notes,
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
			&shift.User.Tel,
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
			&shift.Task.Place,
			&shift.Task.Url,
			&shift.Task.Superviser,
			&shift.Task.Color,
			&shift.Task.Notes,
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
			&shift.User.Tel,
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
			&users.Tel,
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
