package usecase

import (
	"context"

	rep "github.com/NUTFes/SeeFT/api/externals/repository"
	"github.com/NUTFes/SeeFT/api/internals/entity"
	"github.com/pkg/errors"
)

type shiftUseCase struct {
  rep rep.ShiftRepository
}

type ShiftUseCase interface {
  getShiftsByUser(context.Context, string) ([]entity.Shift, error)
  getShiftsByUserAndDateAndWeather(context.Context, string, string, string) ([]entity.Shift, error)
}

func NewBureauUseCase(rep rep.BureauRepository) BureauUseCase {
  return &bureauUseCase{rep}
}

func (a *shiftUseCase) getShiftsByUser(c context.Context, id string) ([]entity.Shift, error) {

  shift := entity.Shift{}
  var shifts []entity.Shift

  // クエリー実行
	rows, err := a.rep.Filter(c.UserID == id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&shift.ID,
			&shift.UserID,
      &shift.Date,
      &shift.Time,
      &shift.WorkID,
      &shift.Weather,
      &shift.Attendance,
			&shift.CreatedAt,
			&shift.UpdatedAt,
		)

		if err != nil {
			return nil, errors.Wrapf(err, "cannot connect SQL")
		}

		shifts = append(shifts, bureau)
	}
	return shifts, nil
}

func (a *shiftUseCase) getShiftsByUserAndDateAndWeather(c context.Context, id string, date string, weather string) ([]entity.Shift, error) {

  shift := entity.Shift{}
  var shifts []entity.Shift

  // クエリー実行
	rows, err := a.rep.Filter(c.UserID == id &&  c.date == date && c.weather == weather)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&shift.ID,
			&shift.UserID,
      &shift.Date,
      &shift.Time,
      &shift.WorkID,
      &shift.Weather,
      &shift.Attendance,
			&shift.CreatedAt,
			&shift.UpdatedAt,
		)

		if err != nil {
			return nil, errors.Wrapf(err, "cannot connect SQL")
		}

		shifts = append(shifts, bureau)
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
