package usecase

import (
	"context"

	"github.com/NUTFes/SeeFT/api/lib/entity"
	rep "github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/pkg/errors"
)

type timeUseCase struct {
	rep rep.TimeRepository
}

type TimeUseCase interface {
	GetTimes(context.Context) ([]entity.Time, error)
	GetTimeByID(context.Context, string) (entity.Time, error)
}

func NewTimeUseCase(rep rep.TimeRepository) TimeUseCase {
	return &timeUseCase{rep}
}

func (b *timeUseCase) GetTimes(c context.Context) ([]entity.Time, error) {
	time := entity.Time{}
	var times []entity.Time

	rows, err := b.rep.All(c)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		err := rows.Scan(
			&time.ID,
			&time.Time,
			&time.CreatedAt,
			&time.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot connect SQL")
		}

		times = append(times, time)
	}

	return times, nil
}

func (b *timeUseCase) GetTimeByID(c context.Context, id string) (entity.Time, error) {
	var time entity.Time
	row, err := b.rep.Find(c, id)
	err = row.Scan(
		&time.ID,
		&time.Time,
		&time.CreatedAt,
		&time.UpdatedAt,
	)
	if err != nil {
		return time, err
	}

	return time, nil
}
