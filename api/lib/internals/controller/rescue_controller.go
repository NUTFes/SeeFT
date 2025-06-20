package controller

import (
	"net/http"
	"strconv"

	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/NUTFes/SeeFT/api/lib/usecase"
	"github.com/labstack/echo/v4"
)

type rescueController struct {
	u usecase.RescueUseCase
}

type RescueController interface {
	CreateRescue(echo.Context) error

	// 作成系
	CreateTroubleRescue(echo.Context) error
	CreateQuestionRescue(echo.Context) error
	CreateShorthandedRescue(echo.Context) error

	// 取得系
	GetRescuesByUserID(echo.Context) error
	GetAllRescues(echo.Context) error

	// 更新系（GAS用）
	UpdateRescues(echo.Context) error
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

//==========
// 各種処理
//==========

// トラブルレスキュー作成
func (r *rescueController) CreateTroubleRescue(c echo.Context) error {
	var req entity.TroubleRescueRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}

	if err := r.u.CreateTroubleRescue(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, "Trouble rescue created successfully")
}

// 質問レスキュー作成
func (r *rescueController) CreateQuestionRescue(c echo.Context) error {
	var req entity.QuestionRescueRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}

	if err := r.u.CreateQuestionRescue(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, "Question rescue created successfully")
}

// 人手不足レスキュー作成
func (r *rescueController) CreateShorthandedRescue(c echo.Context) error {
	var req entity.ShorthandedRescueRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}

	if err := r.u.CreateShorthandedRescue(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, "Shorthanded rescue created successfully")
}

// ユーザーのレスキュー一覧取得
func (r *rescueController) GetRescuesByUserID(c echo.Context) error {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid user ID")
	}

	rescues, err := r.u.GetRescuesByUserID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, rescues)
}

// 全レスキュー取得
func (r *rescueController) GetAllRescues(c echo.Context) error {
	rescues, err := r.u.GetAllRescues(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, rescues)
}

// レスキュー更新（GAS用）
func (r *rescueController) UpdateRescues(c echo.Context) error {
	var req entity.RescueUpdatesRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}

	if err := r.u.UpdateRescues(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Rescues updated successfully")
}
