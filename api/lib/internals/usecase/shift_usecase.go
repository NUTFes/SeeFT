package usecase

import (
	"context"

	rep "github.com/NUTFes/SeeFT/api/lib/externals/repository"
	"github.com/NUTFes/SeeFT/api/lib/internals/entity"
	"github.com/pkg/errors"
)

type shiftUseCase struct {
  rep rep.ShiftRepository
}

type ShiftUseCase interface {
  GetShifts(context.Context) ([]entity.Shift, error)
  GetShiftByID(context.Context, string) (entity.Shift, error)
  GetShiftsByUser(context.Context, string) ([]entity.Shift, error)
  GetShiftsByUserAndDateAndWeather(context.Context, string, string, string) ([]entity.Shift, error)
}

func NewShiftUseCase(rep rep.ShiftRepository) ShiftUseCase {
  return &shiftUseCase{rep}
}

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

func (b *shiftUseCase) GetShiftByID(c context.Context, id string) (entity.Shift, error) {
	var shift entity.Shift
	row, err := b.rep.Find(c, id)
	err = row.Scan(
		&shift.ID,
		&shift.UserID,
		&shift.TaskID,
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
			&shift.UserID,
			&shift.TaskID,
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
			&shift.UserID,
			&shift.TaskID,
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
