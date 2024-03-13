package usecase

import (
	"context"

	rep "github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/pkg/errors"
)

type gradeUseCase struct {
  rep rep.GradeRepository
}

type GradeUseCase interface {
  GetGrades(context.Context) ([]entity.Grade, error)
  GetGradeByID(context.Context, string) (entity.Grade, error) 
}

func NewGradeUseCase(rep rep.GradeRepository) GradeUseCase {
  return &gradeUseCase{rep}
}

func (u *gradeUseCase) GetGrades(c context.Context) ([]entity.Grade, error) {

  grade := entity.Grade{}
  var grades []entity.Grade

  // クエリー実行
	rows, err := u.rep.All(c)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&grade.ID,
			&grade.Grade,
			&grade.CreatedAt,
			&grade.UpdatedAt,
		)

		if err != nil {
			return nil, errors.Wrapf(err, "cannot connect SQL")
		}

		grades = append(grades, grade)
	}
	return grades, nil
}

func (u *gradeUseCase) GetGradeByID(c context.Context, id string) (entity.Grade, error) {
	var grade entity.Grade
	row, err := u.rep.Find(c, id)
	err = row.Scan(
		&grade.ID,
		&grade.Grade,
		&grade.CreatedAt,
		&grade.UpdatedAt,
	)
	if err != nil {
		return grade, err
	}
	return grade, nil
}