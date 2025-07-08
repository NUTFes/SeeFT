package controller

import (
	"net/http"
	"strconv"

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
	userID := c.FormValue("user_id")
	taskID := c.FormValue("task_id")
	missingNumber := c.FormValue("missing_number")
	place := c.FormValue("place")
	status := c.FormValue("status")
	
	if userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User ID is required"})
	}
	if taskID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Task ID is required"})
	}
	if missingNumber == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing number is required"})
	}
	
	// 数値チェック
	if _, err := strconv.Atoi(userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}
	if _, err := strconv.Atoi(taskID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid task ID"})
	}
	if _, err := strconv.Atoi(missingNumber); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid missing number"})
	}
	
	createdShorthandedRescue, err := sr.u.CreateShorthandedRescue(c.Request().Context(), userID, taskID, missingNumber, place, status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, createdShorthandedRescue)
}

// 更新
func (sr *shorthandedRescueController) UpdateShorthandedRescue(c echo.Context) error {
	id := c.Param("id")
	status := c.FormValue("status")
	response := c.FormValue("response")
	
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}
	if status == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Status is required"})
	}
	
	// 数値チェック
	if _, err := strconv.Atoi(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}
	
	updatedShorthandedRescue, err := sr.u.UpdateShorthandedRescue(c.Request().Context(), id, status, response)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, updatedShorthandedRescue)
}

// 削除
func (sr *shorthandedRescueController) DeleteShorthandedRescue(c echo.Context) error {
	id := c.FormValue("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}
	
	// 数値チェック
	if _, err := strconv.Atoi(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}
	
	err := sr.u.DeleteShorthandedRescue(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Shorthanded rescue deleted successfully"})
}
