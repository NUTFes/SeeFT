package controller

import (
	"net/http"

	"github.com/NUTFes/SeeFT/api/lib/usecase"
	"github.com/labstack/echo/v4"
)

type rescueController struct {
	u usecase.RescueUseCase
}

type RescueController interface {
	CreateRescue(echo.Context) error
}

func NewRescueController(u usecase.RescueUseCase) RescueController {
	return &rescueController{u}
}

func (r *rescueController) CreateRescue(c echo.Context) error {
	taskID := c.FormValue("task_id")
	err := r.u.CreateRescue(c.Request().Context(), taskID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Rescue request sent successfully")
}
