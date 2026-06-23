package controller

import (
	"net/http"

	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/NUTFes/SeeFT/api/lib/usecase"
	"github.com/labstack/echo/v4"
)

type shiftController struct {
	u usecase.ShiftUseCase
}

type ShiftController interface {
	ShowUsersByShift(echo.Context) error
	ShowShiftsByUserAndDateAndWeather(echo.Context) error
	ShowShiftCardsByUserAndDateAndWeather(echo.Context) error
	PostShiftCards(echo.Context) error
	IndexShiftAdmin(echo.Context) error
	ShowShiftAdmin(echo.Context) error
	CreateShiftAdmin(echo.Context) error
	UpdateShiftAdmin(echo.Context) error
	DeleteShiftAdmin(echo.Context) error
	ShowShiftAdminByDateAndWeather(echo.Context) error
	ShowShiftAdminByDateAndWeatherAndTime(echo.Context) error
	SerachMaxID(echo.Context) error
	SubmitShift(echo.Context) error
	UpdateShiftsFromGAS(echo.Context) error
}

func NewShiftController(u usecase.ShiftUseCase) ShiftController {
	return &shiftController{u}
}

func (b *shiftController) ShowUsersByShift(c echo.Context) error {
	task := c.Param("task_id")
	year := c.Param("year_id")
	date := c.Param("date_id")
	time := c.Param("time_id")
	weather := c.Param("weather_id")
	shiftUsers, err := b.u.GetUsersByShift(c.Request().Context(), task, year, date, time, weather)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, shiftUsers)
}

func (b *shiftController) ShowShiftsByUserAndDateAndWeather(c echo.Context) error {
	id := c.Param("user_id")
	date := c.Param("date")
	weather := c.Param("weather")
	shifts, err := b.u.GetShiftsByUserAndDateAndWeather(c.Request().Context(), id, date, weather)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, shifts)
}

func (b *shiftController) ShowShiftCardsByUserAndDateAndWeather(c echo.Context) error {
	userID := c.Param("user_id")
	dateID := c.Param("date_id")
	weatherID := c.Param("weather_id")
	shiftCards, err := b.u.GetShiftCardsByUserAndDateAndWeather(c.Request().Context(), userID, dateID, weatherID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, shiftCards)
}

func (b *shiftController) PostShiftCards(c echo.Context) error {
	// Accept JSON body: { "userId": "...", "dateId": "...", "weatherId": "..." }
	var req struct {
		UserID    string `json:"userId"`
		DateID    string `json:"dateId"`
		WeatherID string `json:"weatherId"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid request body"})
	}
	// Required fields validation
	if req.UserID == "" || req.DateID == "" || req.WeatherID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "userId, dateId, weatherId are required"})
	}

	shiftCards, err := b.u.GetShiftCardsByUserAndDateAndWeather(c.Request().Context(), req.UserID, req.DateID, req.WeatherID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, shiftCards)
}

// Webアプリ用API

func (b *shiftController) IndexShiftAdmin(c echo.Context) error {
	shifts, err := b.u.GetShiftsAdmin(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, shifts)
}

func (b *shiftController) ShowShiftAdmin(c echo.Context) error {
	id := c.Param("id")
	shift, err := b.u.GetShiftAdminByID(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, shift)
}

// Create
func (u *shiftController) CreateShiftAdmin(c echo.Context) error {
	taskID := c.QueryParam("task_id")
	userID := c.QueryParam("user_id")
	yearID := c.QueryParam("year_id")
	dateID := c.QueryParam("date_id")
	timeID := c.QueryParam("time_id")
	weatherID := c.QueryParam("weather_id")
	isAttendance := c.QueryParam("is_attendance")
	latastShift, err := u.u.CreateShiftAdmin(c.Request().Context(), taskID, userID, yearID, dateID, timeID, weatherID, isAttendance)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, latastShift)
}

// Update
func (u *shiftController) UpdateShiftAdmin(c echo.Context) error {
	id := c.Param("id")
	taskID := c.QueryParam("task_id")
	userID := c.QueryParam("user_id")
	yearID := c.QueryParam("year_id")
	dateID := c.QueryParam("date_id")
	timeID := c.QueryParam("time_id")
	weatherID := c.QueryParam("weather_id")
	isAttendance := c.QueryParam("is_attendance")
	updatedShift, err := u.u.UpdateShiftAdmin(c.Request().Context(), id, taskID, userID, yearID, dateID, timeID, weatherID, isAttendance)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, updatedShift)
}

// Destroy
func (u *shiftController) DeleteShiftAdmin(c echo.Context) error {
	id := c.QueryParam("id")
	err := u.u.DeleteShiftAdmin(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "Destroy Shift")
}

func (b *shiftController) ShowShiftAdminByDateAndWeather(c echo.Context) error {
	date := c.Param("date")
	weather := c.Param("weather")
	shifts, err := b.u.GetShiftsAdminByDateAndWeather(c.Request().Context(), date, weather)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, shifts)
}

func (b *shiftController) ShowShiftAdminByDateAndWeatherAndTime(c echo.Context) error {
	date := c.Param("date")
	weather := c.Param("weather")
	lower := c.Param("lower")
	upper := c.Param("upper")
	shifts, err := b.u.GetShiftsAdminByDateAndWeatherAndTime(c.Request().Context(), date, weather, lower, upper)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, shifts)
}

func (b *shiftController) SerachMaxID(c echo.Context) error {
	id, err := b.u.GetMaxID(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, id)
}

func (sc *shiftController) SubmitShift(c echo.Context) error {
	var req entity.ShiftRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}

	// DBに保存
	if err := sc.u.SaveShiftData(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to save shift data")
	}

	// GASに送信
	if err := sc.u.SendToGAS(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to send data to GAS")
	}

	return c.JSON(http.StatusOK, "Shift data submitted successfully")
}

// GASからのシフト変更通知を受け取るエンドポイント
func (sc *shiftController) UpdateShiftsFromGAS(c echo.Context) error {
	var req entity.ShiftChangeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}

	// 必要に応じてユースケース層へ処理を委譲
	if err := sc.u.UpdateShiftsFromGAS(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Shifts updated successfully")
}
