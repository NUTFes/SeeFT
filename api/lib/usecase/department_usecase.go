package usecase

import (
	"context"

	rep "github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/pkg/errors"
)

type departmentUseCase struct {
	rep rep.DepartmentRepository
}

type DepartmentUseCase interface {
	GetDepartments(context.Context) ([]entity.Department, error)
	GetDepartmentByID(context.Context, string) (entity.Department, error)
}

func NewDepartmentUseCase(rep rep.DepartmentRepository) DepartmentUseCase {
	return &departmentUseCase{rep}
}

func (u *departmentUseCase) GetDepartments(c context.Context) ([]entity.Department, error){

	department := entity.Department{}
	var departments []entity.Department

	//クエリー実行
	  rows, err := u.rep.All(c)
	  if err != nil {
		  return nil, err
	  }
	  defer rows.Close()

	  for rows.Next() {
		  err := rows.Scan(
			  &department.ID,
			  &department.Department,
			  &department.CreatedAt,
			  &department.UpdatedAt,
		  )

		  if err != nil {
			return nil, errors.Wrapf(err, "cannot connect SQL")
		  }

		  departments = append(departments, department)
	  }
	  return departments, nil
}

func (u *departmentUseCase) GetDepartmentByID(c context.Context, id string)(entity.Department, error){
	var department entity.Department
	row, err := u.rep.Find(c, id)
	if err != nil {
		return department, err
	}
	err = row.Scan(
		&department.ID,
		&department.Department,
		&department.CreatedAt,
		&department.UpdatedAt,
	)
	if err != nil{
		return department, err
	}
	return department, nil
}