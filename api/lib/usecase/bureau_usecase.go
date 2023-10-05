package usecase

import (
	"context"

	rep "github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/pkg/errors"
)

type bureauUseCase struct {
  rep rep.BureauRepository
}

type BureauUseCase interface {
  GetBureaus(context.Context) ([]entity.Bureau, error)
  GetBureauByID(context.Context, string) (entity.Bureau, error) 
}

func NewBureauUseCase(rep rep.BureauRepository) BureauUseCase {
  return &bureauUseCase{rep}
}

func (a *bureauUseCase) GetBureaus(c context.Context) ([]entity.Bureau, error) {

  bureau := entity.Bureau{}
  var bureaues []entity.Bureau

  // クエリー実行
	rows, err := a.rep.All(c)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&bureau.ID,
			&bureau.Bureau,
			&bureau.Color,
			&bureau.CreatedAt,
			&bureau.UpdatedAt,
		)

		if err != nil {
			return nil, errors.Wrapf(err, "cannot connect SQL")
		}

		bureaues = append(bureaues, bureau)
	}
	return bureaues, nil
}

func (b *bureauUseCase) GetBureauByID(c context.Context, id string) (entity.Bureau, error) {
	var bureau entity.Bureau
	row, err := b.rep.Find(c, id)
	err = row.Scan(
		&bureau.ID,
		&bureau.Bureau,
		&bureau.Color,
		&bureau.CreatedAt,
		&bureau.UpdatedAt,
	)
	if err != nil {
		return bureau, err
	}
	return bureau, nil
}


// import '../entity/entity.dart';
// import './repository/repository.dart';

// abstract class BureauUsecase {
//   Future<List<Bureau>> getBureaus(ctx);
// }

// class BureauUsecaseImpl implements BureauUsecase {
//   BureauRepository bureauRepository;

//   BureauUsecaseImpl(this.bureauRepository);

//   @override
//   Future<List<Bureau>> getBureaus(ctx) async {
//     List<Bureau> list = await bureauRepository.getBureaus(ctx);

//     return list;
//   }
// }
