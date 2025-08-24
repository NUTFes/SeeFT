package controller

import (
	"net/http"

	"github.com/NUTFes/SeeFT/api/lib/usecase"
	"github.com/labstack/echo/v4"
)

type timeController struct {
	u usecase.TimeUseCase
}

type TimeController interface {
	IndexTime(echo.Context) error
	ShowTime(echo.Context) error
}

func NewTimeController(u usecase.TimeUseCase) TimeController {
	return &timeController{u}
}

func (b *timeController) IndexTime(c echo.Context) error {
	times, err := b.u.GetTimes(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, times)
}

func (b *timeController) ShowTime(c echo.Context) error {
	id := c.Param("id")
	time, err := b.u.GetTimeByID(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, time)
}
