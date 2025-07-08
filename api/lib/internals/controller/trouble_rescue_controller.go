package controller

import (
	"net/http"
	"strconv"

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
	userID := c.FormValue("user_id")
	taskID := c.FormValue("task_id")
	place := c.FormValue("place")
	detail := c.FormValue("detail")
	status := c.FormValue("status")
	
	if userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User ID is required"})
	}
	if taskID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Task ID is required"})
	}
	if detail == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Detail is required"})
	}
	
	// 数値チェック
	if _, err := strconv.Atoi(userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}
	if _, err := strconv.Atoi(taskID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid task ID"})
	}
	
	createdTroubleRescue, err := tr.u.CreateTroubleRescue(c.Request().Context(), userID, taskID, place, detail, status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, createdTroubleRescue)
}

// 更新
func (tr *troubleRescueController) UpdateTroubleRescue(c echo.Context) error {
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
	
	updatedTroubleRescue, err := tr.u.UpdateTroubleRescue(c.Request().Context(), id, status, response)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, updatedTroubleRescue)
}

// 削除
func (tr *troubleRescueController) DeleteTroubleRescue(c echo.Context) error {
	id := c.FormValue("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}
	
	// 数値チェック
	if _, err := strconv.Atoi(id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}
	
	err := tr.u.DeleteTroubleRescue(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Trouble rescue deleted successfully"})
}
