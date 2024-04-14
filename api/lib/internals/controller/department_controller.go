package controller

import (
	"net/http"

	"github.com/NUTFes/SeeFT/api/lib/usecase"
	"github.com/labstack/echo/v4"
)

type departmentController struct {
	u usecase.DepartmentUsecase
}

type DepartmentController interface {
	IndexDepartment(echo.context) error
	ShowDepartment(echo.context) error
}

func NewDepartmentController(u usecase.DepartmentUseCase) DepartmentContoller {
	return &departmentController{u}
}

func (b *departmentController) IndexDepartment(c echo.Context) error {
	departments, err := b.u.GetDepartments(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, departments)
}

func (b *departmentController) ShowDepartment(c echo.Context) error {
	id := c.Param("id")
	department, err := b.u.GetDepartmentByID(c.Request().Context(),id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, department)
}