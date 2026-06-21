package usecase

import (
	"context"

	"github.com/NUTFes/SeeFT/api/lib/entity"
	rep "github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/pkg/errors"
)

type placeUseCase struct {
	rep rep.PlaceRepository
}

type PlaceUseCase interface {
	GetPlaces(context.Context) ([]entity.Place, error)
	GetPlaceByID(context.Context, string) (entity.Place, error)
	CreatePlace(context.Context, string, string) (entity.Place, error)
	UpdatePlace(context.Context, string, string, string) (entity.Place, error)
	DeletePlace(context.Context, string) error
}

func NewPlaceUseCase(rep rep.PlaceRepository) PlaceUseCase {
	return &placeUseCase{rep}
}

func (u *placeUseCase) GetPlaces(c context.Context) ([]entity.Place, error) {

	place := entity.Place{}
	var places []entity.Place

	// クエリー実行
	rows, err := u.rep.All(c)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&place.ID,
			&place.Place,
			&place.Remark,
			&place.CreatedAt,
			&place.UpdatedAt,
		)

		if err != nil {
			return nil, errors.Wrapf(err, "cannot connect SQL")
		}

		places = append(places, place)
	}
	return places, nil
}

func (u *placeUseCase) GetPlaceByID(c context.Context, id string) (entity.Place, error) {
	var place entity.Place
	row, err := u.rep.Find(c, id)
	if err != nil {
		return place, err
	}
	err = row.Scan(
		&place.ID,
		&place.Place,
		&place.Remark,
		&place.CreatedAt,
		&place.UpdatedAt,
	)
	if err != nil {
		return place, err
	}
	return place, nil
}

func (u *placeUseCase) CreatePlace(c context.Context, place string, remark string) (entity.Place, error) {
	latastPlace := entity.Place{}
	if err := u.rep.Create(c, place, remark); err != nil {
		return latastPlace, err
	}
	row, err := u.rep.FindNewRecord(c)
	if err != nil {
		return latastPlace, err
	}
	err = row.Scan(
		&latastPlace.ID,
		&latastPlace.Place,
		&latastPlace.Remark,
		&latastPlace.CreatedAt,
		&latastPlace.UpdatedAt,
	)
	if err != nil {
		return latastPlace, err
	}
	return latastPlace, err
}

func (u *placeUseCase) UpdatePlace(c context.Context, id string, placeName string, remark string) (entity.Place, error) {
	updatedPlace := entity.Place{}
	var place entity.Place

	row, err := u.rep.Find(c, id)
	if err != nil {
		return place, err
	}
	err = row.Scan(
		&place.ID,
		&place.Place,
		&place.Remark,
		&place.CreatedAt,
		&place.UpdatedAt,
	)
	if err != nil {
		return place, err
	}

	if err = u.rep.Update(c, id, placeName, remark); err != nil {
		return place, err
	}
	row, err = u.rep.Find(c, id)
	if err != nil {
		return updatedPlace, err
	}
	err = row.Scan(
		&updatedPlace.ID,
		&updatedPlace.Place,
		&updatedPlace.Remark,
		&updatedPlace.CreatedAt,
		&updatedPlace.UpdatedAt,
	)
	if err != nil {
		return updatedPlace, err
	}
	return updatedPlace, nil
}

func (u *placeUseCase) DeletePlace(c context.Context, id string) error {
	err := u.rep.Destroy(c, id)
	return err
}
