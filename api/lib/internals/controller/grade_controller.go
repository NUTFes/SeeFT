package controller

import (
	"net/http"

	"github.com/NUTFes/SeeFT/api/lib/usecase"
	"github.com/labstack/echo/v4"
)

type gradeController struct {
	u usecase.GradeUseCase
}

type GradeController interface {
	IndexGrade(echo.Context) error
	ShowGrade(echo.Context) error
}

func NewGradeController(u usecase.GradeUseCase) GradeController {
	return &gradeController{u}
}

func (b *gradeController) IndexGrade(c echo.Context) error {
	grades, err := b.u.GetGrades(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, grades)
}

func (b *gradeController) ShowGrade(c echo.Context) error {
	id := c.Param("id")
	grade, err := b.u.GetGradeByID(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, grade)
}

