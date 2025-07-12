package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/NUTFes/SeeFT/api/lib/usecase"
	"github.com/labstack/echo/v4"
)

// 共通のリクエスト構造体
type RescueRequest struct {
Type    string      `json:"type"`
UserID  int         `json:"user_id"`
Content interface{} `json:"content"`
}

// 各タイプ別のContent構造体
type TroubleContent struct {
	TaskID int    `json:"task_id"`
	Place  string `json:"place"`
	Detail string `json:"detail"`
}

type QuestionContent struct {
	Question string `json:"question"`
}

type ShorthandedContent struct {
	TaskID        int    `json:"task_id"`
	MissingNumber int    `json:"missing_number"`
	Place         string `json:"place"`
}

type rescueUnifiedController struct {
	questionRescueUseCase    usecase.QuestionRescueUseCase
	shorthandedRescueUseCase usecase.ShorthandedRescueUseCase
	troubleRescueUseCase     usecase.TroubleRescueUseCase
	rescueUnifiedUseCase     usecase.RescueUnifiedUseCase
}

type RescueUnifiedController interface {
	CreateRescue(echo.Context) error
	GetAllRescues(echo.Context) error
	GetRescuesByUserID(echo.Context) error
}

func NewRescueUnifiedController(
	qu usecase.QuestionRescueUseCase,
	su usecase.ShorthandedRescueUseCase,
	tu usecase.TroubleRescueUseCase,
	ru usecase.RescueUnifiedUseCase,
) RescueUnifiedController {
	return &rescueUnifiedController{qu, su, tu, ru}
}

// 統一されたレスキューリクエスト処理
func (r *rescueUnifiedController) CreateRescue(c echo.Context) error {
	var req RescueRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	// 基本的な入力バリデーション
	if req.UserID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}
	userIDStr := strconv.Itoa(req.UserID)

	switch req.Type {
	case "trouble":
		return r.handleTroubleRescue(c, req, userIDStr)
	case "question":
		return r.handleQuestionRescue(c, req, userIDStr)
	case "shorthanded":
		return r.handleShorthandedRescue(c, req, userIDStr)
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid rescue type"})
	}
}

func (r *rescueUnifiedController) handleTroubleRescue(c echo.Context, req RescueRequest, userIDStr string) error {
	// content をTroubleContentに変換
	contentBytes, err := json.Marshal(req.Content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid trouble content"})
	}

	var content TroubleContent
	if err := json.Unmarshal(contentBytes, &content); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid trouble content format"})
	}

	// バリデーション
	if content.TaskID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid task ID"})
	}
	if content.Detail == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Detail is required"})
	}

	taskIDStr := strconv.Itoa(content.TaskID)
	createdTroubleRescue, err := r.troubleRescueUseCase.CreateTroubleRescue(
		c.Request().Context(),
		userIDStr,
		taskIDStr,
		content.Place,
		content.Detail,
		"todo",
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Trouble rescue created successfully",
		"data":    createdTroubleRescue,
	})
}

func (r *rescueUnifiedController) handleQuestionRescue(c echo.Context, req RescueRequest, userIDStr string) error {
	// content をQuestionContentに変換
	contentBytes, err := json.Marshal(req.Content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid question content"})
	}

	var content QuestionContent
	if err := json.Unmarshal(contentBytes, &content); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid question content format"})
	}

	// バリデーション
	if content.Question == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Question is required"})
	}

	createdQuestionRescue, err := r.questionRescueUseCase.CreateQuestionRescue(
		c.Request().Context(),
		userIDStr,
		content.Question,
		"todo",
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Question rescue created successfully",
		"data":    createdQuestionRescue,
	})
}

func (r *rescueUnifiedController) handleShorthandedRescue(c echo.Context, req RescueRequest, userIDStr string) error {
	// content をShorthandedContentに変換
	contentBytes, err := json.Marshal(req.Content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid shorthanded content"})
	}

	var content ShorthandedContent
	if err := json.Unmarshal(contentBytes, &content); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid shorthanded content format"})
	}

	// バリデーション
	if content.TaskID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid task ID"})
	}
	if content.MissingNumber <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing number must be greater than 0"})
	}

	taskIDStr := strconv.Itoa(content.TaskID)
	missingNumberStr := strconv.Itoa(content.MissingNumber)

	createdShorthandedRescue, err := r.shorthandedRescueUseCase.CreateShorthandedRescue(
		c.Request().Context(),
		userIDStr,
		taskIDStr,
		missingNumberStr,
		content.Place,
		"todo",
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Shorthanded rescue created successfully",
		"data":    createdShorthandedRescue,
	})
}

// 全件取得
func (r *rescueUnifiedController) GetAllRescues(c echo.Context) error {
	rescues, err := r.rescueUnifiedUseCase.GetAllRescues(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, rescues)
}

// ユーザーID別取得
func (r *rescueUnifiedController) GetRescuesByUserID(c echo.Context) error {
	userID := c.Param("user_id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User ID is required"})
	}

	// 数値チェック
	if _, err := strconv.Atoi(userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	rescues, err := r.rescueUnifiedUseCase.GetRescuesByUserID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, rescues)
}
