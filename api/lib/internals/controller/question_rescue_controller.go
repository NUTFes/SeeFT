package controller

import (
	"net/http"
	"strconv"

	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/NUTFes/SeeFT/api/lib/usecase"
	"github.com/labstack/echo/v4"
)

type questionRescueController struct {
	u usecase.QuestionRescueUseCase
}

type QuestionRescueController interface {
	IndexQuestionRescue(echo.Context) error
	ShowQuestionRescue(echo.Context) error
	ShowQuestionRescuesByUserID(echo.Context) error
	CreateQuestionRescue(echo.Context) error
	UpdateQuestionRescue(echo.Context) error
	DeleteQuestionRescue(echo.Context) error
}

func NewQuestionRescueController(u usecase.QuestionRescueUseCase) QuestionRescueController {
	return &questionRescueController{u}
}

// 全件取得
func (qr *questionRescueController) IndexQuestionRescue(c echo.Context) error {
	questionRescues, err := qr.u.GetQuestionRescues(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, questionRescues)
}

// 1件取得
func (qr *questionRescueController) ShowQuestionRescue(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}
	
	questionRescue, err := qr.u.GetQuestionRescueByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, questionRescue)
}

// ユーザIDで検索
func (qr *questionRescueController) ShowQuestionRescuesByUserID(c echo.Context) error {
	userID := c.Param("user_id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User ID is required"})
	}
	
	questionRescues, err := qr.u.GetQuestionRescuesByUserID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, questionRescues)
}

// 作成
func (qr *questionRescueController) CreateQuestionRescue(c echo.Context) error {
	var req entity.QuestionRescueCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}
	if req.UserID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User ID is required"})
	}
	if req.Question == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Question is required"})
	}
	userIDStr := strconv.Itoa(req.UserID)
	createdQuestionRescue, err := qr.u.CreateQuestionRescue(c.Request().Context(), userIDStr, req.Question, req.Status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, createdQuestionRescue)
}

// 更新
func (qr *questionRescueController) UpdateQuestionRescue(c echo.Context) error {
	var req entity.QuestionRescueUpdateRequest
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
	updatedQuestionRescue, err := qr.u.UpdateQuestionRescue(c.Request().Context(), id, req.Status, req.Response)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, updatedQuestionRescue)
}

// 削除
func (qr *questionRescueController) DeleteQuestionRescue(c echo.Context) error {
	var req entity.QuestionRescueDeleteRequest
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
	err := qr.u.DeleteQuestionRescue(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Question rescue deleted successfully"})
}
