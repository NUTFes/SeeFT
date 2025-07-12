package controller

import (
	"net/http"
	"strconv"

	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/NUTFes/SeeFT/api/lib/usecase"
	"github.com/labstack/echo/v4"
)

type troubleRescueController struct {
	u usecase.TroubleRescueUseCase
}

type TroubleRescueController interface {
	IndexTroubleRescue(echo.Context) error
	ShowTroubleRescue(echo.Context) error
	ShowTroubleRescuesByUserID(echo.Context) error
	ShowTroubleRescuesByTaskID(echo.Context) error
	CreateTroubleRescue(echo.Context) error
	UpdateTroubleRescue(echo.Context) error
	DeleteTroubleRescue(echo.Context) error
}

func NewTroubleRescueController(u usecase.TroubleRescueUseCase) TroubleRescueController {
	return &troubleRescueController{u}
}

// 全件取得
func (tr *troubleRescueController) IndexTroubleRescue(c echo.Context) error {
	troubleRescues, err := tr.u.GetTroubleRescues(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, troubleRescues)
}

// 1件取得
func (tr *troubleRescueController) ShowTroubleRescue(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}
	
	troubleRescue, err := tr.u.GetTroubleRescueByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, troubleRescue)
}

// ユーザIDで検索
func (tr *troubleRescueController) ShowTroubleRescuesByUserID(c echo.Context) error {
	userID := c.Param("user_id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User ID is required"})
	}
	
	troubleRescues, err := tr.u.GetTroubleRescuesByUserID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, troubleRescues)
}

// タスクIDで検索
func (tr *troubleRescueController) ShowTroubleRescuesByTaskID(c echo.Context) error {
	taskID := c.Param("task_id")
	if taskID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Task ID is required"})
	}
	
	troubleRescues, err := tr.u.GetTroubleRescuesByTaskID(c.Request().Context(), taskID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, troubleRescues)
}

// 作成
func (tr *troubleRescueController) CreateTroubleRescue(c echo.Context) error {
	var req entity.TroubleRescueCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}
	if req.UserID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User ID is required"})
	}
	if req.TaskID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Task ID is required"})
	}
	if req.Detail == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Detail is required"})
	}
	userIDStr := strconv.Itoa(req.UserID)
	taskIDStr := strconv.Itoa(req.TaskID)
	createdTroubleRescue, err := tr.u.CreateTroubleRescue(c.Request().Context(), userIDStr, taskIDStr, req.Place, req.Detail, req.Status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, createdTroubleRescue)
}

// 更新
func (tr *troubleRescueController) UpdateTroubleRescue(c echo.Context) error {
	var req entity.TroubleRescueUpdateRequest
	id := c.Param("id")
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}
	if req.Status == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Status is required"})
	}
	if _, err := strconv.Atoi(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}
	updatedTroubleRescue, err := tr.u.UpdateTroubleRescue(c.Request().Context(), id, req.Status, req.Response)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, updatedTroubleRescue)
}

// 削除
func (tr *troubleRescueController) DeleteTroubleRescue(c echo.Context) error {
	var req entity.TroubleRescueDeleteRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}
	id := req.ID
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}
	if _, err := strconv.Atoi(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}
	err := tr.u.DeleteTroubleRescue(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Trouble rescue deleted successfully"})
}
