package controller

import (
	"net/http"
	"strconv"

	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/NUTFes/SeeFT/api/lib/usecase"
	"github.com/labstack/echo/v4"
)

type shorthandedRescueController struct {
	u usecase.ShorthandedRescueUseCase
}

type ShorthandedRescueController interface {
	IndexShorthandedRescue(echo.Context) error
	ShowShorthandedRescue(echo.Context) error
	ShowShorthandedRescuesByUserID(echo.Context) error
	ShowShorthandedRescuesByTaskID(echo.Context) error
	CreateShorthandedRescue(echo.Context) error
	UpdateShorthandedRescue(echo.Context) error
	DeleteShorthandedRescue(echo.Context) error
}

func NewShorthandedRescueController(u usecase.ShorthandedRescueUseCase) ShorthandedRescueController {
	return &shorthandedRescueController{u}
}

// 全件取得
func (sr *shorthandedRescueController) IndexShorthandedRescue(c echo.Context) error {
	shorthandedRescues, err := sr.u.GetShorthandedRescues(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, shorthandedRescues)
}

// 1件取得
func (sr *shorthandedRescueController) ShowShorthandedRescue(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}
	
	shorthandedRescue, err := sr.u.GetShorthandedRescueByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, shorthandedRescue)
}

// ユーザIDで検索
func (sr *shorthandedRescueController) ShowShorthandedRescuesByUserID(c echo.Context) error {
	userID := c.Param("user_id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User ID is required"})
	}
	
	shorthandedRescues, err := sr.u.GetShorthandedRescuesByUserID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, shorthandedRescues)
}

// タスクIDで検索
func (sr *shorthandedRescueController) ShowShorthandedRescuesByTaskID(c echo.Context) error {
	taskID := c.Param("task_id")
	if taskID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Task ID is required"})
	}
	
	shorthandedRescues, err := sr.u.GetShorthandedRescuesByTaskID(c.Request().Context(), taskID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, shorthandedRescues)
}

// 作成
func (sr *shorthandedRescueController) CreateShorthandedRescue(c echo.Context) error {
	var req entity.ShorthandedRescueCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}
	if req.UserID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User ID is required"})
	}
	if req.TaskID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Task ID is required"})
	}
	if req.MissingNumber <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing number is required"})
	}
	userIDStr := strconv.Itoa(req.UserID)
	taskIDStr := strconv.Itoa(req.TaskID)
	missingNumberStr := strconv.Itoa(req.MissingNumber)
	createdShorthandedRescue, err := sr.u.CreateShorthandedRescue(c.Request().Context(), userIDStr, taskIDStr, missingNumberStr, req.Place, req.Status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, createdShorthandedRescue)
}

// 更新
func (sr *shorthandedRescueController) UpdateShorthandedRescue(c echo.Context) error {
	var req entity.ShorthandedRescueUpdateRequest
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
	updatedShorthandedRescue, err := sr.u.UpdateShorthandedRescue(c.Request().Context(), id, req.Status, req.Response)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, updatedShorthandedRescue)
}

// 削除
func (sr *shorthandedRescueController) DeleteShorthandedRescue(c echo.Context) error {
	var req entity.ShorthandedRescueDeleteRequest
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
	err := sr.u.DeleteShorthandedRescue(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Shorthanded rescue deleted successfully"})
}
