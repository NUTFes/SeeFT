package controller

import (
	"net/http"
	"strconv"

	"github.com/NUTFes/SeeFT/api/lib/usecase"
	"github.com/labstack/echo/v4"
)

type reviewController struct {
	u usecase.ReviewUseCase
}

type ReviewController interface {
	IndexReview(echo.Context) error
	ShowReview(echo.Context) error
	CreateReview(echo.Context) error
	UpdateReview(echo.Context) error
	DeleteReview(echo.Context) error
}

func NewReviewController(u usecase.ReviewUseCase) ReviewController {
	return &reviewController{u}
}

// 全件取得
func (rc *reviewController) IndexReview(c echo.Context) error {
	reviews, err := rc.u.GetReviewsGAS(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, reviews)
}

// 1件取得
func (rc *reviewController) ShowReview(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}
	review, err := rc.u.GetReviewGASByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, review)
}

// 作成
func (rc *reviewController) CreateReview(c echo.Context) error {
       var req struct {
	       UserID         int    `json:"user_id"`
	       TaskName       string `json:"task_name"`
	       StaffingRating int    `json:"staffing_rating"`
	       ManualRating   int    `json:"manual_rating"`
	       Comment        string `json:"comment"`
       }
       if err := c.Bind(&req); err != nil {
	       return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
       }
       review, err := rc.u.CreateReview(
	       c.Request().Context(),
	       strconv.Itoa(req.UserID),
	       req.TaskName,
	       strconv.Itoa(req.StaffingRating),
	       strconv.Itoa(req.ManualRating),
	       req.Comment,
       )
       if err != nil {
	       return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
       }
       return c.JSON(http.StatusCreated, review)
}

// 編集
func (rc *reviewController) UpdateReview(c echo.Context) error {
       id := c.Param("id")
       if id == "" {
	       return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
       }
       var req struct {
	       UserID         int    `json:"user_id"`
	       TaskName       string `json:"task_name"`
	       StaffingRating int    `json:"staffing_rating"`
	       ManualRating   int    `json:"manual_rating"`
	       Comment        string `json:"comment"`
       }
       if err := c.Bind(&req); err != nil {
	       return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
       }
       review, err := rc.u.UpdateReview(
	       c.Request().Context(),
	       id,
	       strconv.Itoa(req.UserID),
	       req.TaskName,
	       strconv.Itoa(req.StaffingRating),
	       strconv.Itoa(req.ManualRating),
	       req.Comment,
       )
       if err != nil {
	       return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
       }
       return c.JSON(http.StatusOK, review)
}

// 削除
func (rc *reviewController) DeleteReview(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}
	if err := rc.u.DeleteReview(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
