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

type RescueData struct {
	TaskID string `json:"task_id"` // 小文字taskIDを大文字TaskIDに変更、JSONタグも修正
}

func NewRescueController(u usecase.RescueUseCase) RescueController {
	return &rescueController{u}
}

func (r *rescueController) CreateRescue(c echo.Context) error {
	var req RescueData
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "error")
	}
	err = r.u.CreateRescue(c.Request().Context(), req.TaskID) // := を = に修正、repをreqに修正
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Rescue request sent successfully")
}