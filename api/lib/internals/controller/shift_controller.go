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
	IndexShift(echo.Context) error
	ShowShift(echo.Context) error
	ShowShiftsByUser(echo.Context) error
	ShowUsersByShift(echo.Context) error
	ShowShiftsByUserAndDateAndWeather(echo.Context) error
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

// スマホ用API
func (b *shiftController) IndexShift(c echo.Context) error {
	shifts, err := b.u.GetShifts(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, shifts)
}

func (b *shiftController) ShowShift(c echo.Context) error {
	id := c.Param("id")
	shift, err := b.u.GetShiftByID(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, shift)
}

func (b *shiftController) ShowShiftsByUser(c echo.Context) error {
	id := c.Param("user_id")
	shifts, err := b.u.GetShiftsByUser(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, shifts)
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

// import 'dart:convert';
// import 'package:shelf/shelf.dart';

// import '../../config/config.dart';
// import '../../entity/entity.dart';
// import '../../usecase/usecase.dart';

// class ShiftController {
//   final StatusResponse statusResponse;
//   final ShiftUsecase shiftUsecase;

//   ShiftController(
//     this.statusResponse,
//     this.shiftUsecase,
//   );

//   Future<Response> getShiftsByUser(Request request, String id) async {
//     try {
//       final req = User(id: int.parse(id));
//       final res = await shiftUsecase.getShiftsByUser(request.context, req);

//       return statusResponse.responseOK(jsonEncode(res));
//     } catch (e) {
//       Log.severe('shiftContoller.getShiftsByUser: ${e.toString()}');
//       var json = jsonEncode({'message': e.toString()});
//       return statusResponse.responseBadRequest(json);
//     }
//   }

//   Future<Response> getShiftsByUserAndDateAndWeather(
//       Request request, String userId, String dateId, String weatherId) async {
//     try {
//       final req = Shift(
//           user: User(id: int.parse(userId)),
//           date: Date(id: int.parse(dateId)),
//           weather: Weather(id: int.parse(weatherId)));
//       final res = await shiftUsecase.getShiftsByUserAndDateAndWeather(request.context, req);
//       return statusResponse.responseOK(jsonEncode(res));
//     } catch (e) {
//       Log.severe('shiftController.getShiftsByUserAndDateAndWeather: ${e.toString()}');
//       var json = jsonEncode({'message': e.toString()});
//       return statusResponse.responseBadRequest(json);
//     }
//   }

//   Future<Response> getShiftsByYearAndDateAndWeather(
//       Request request, String yearId, String dateId, String weatherId) async {
//     try {
//       final req = Shift(
//           year: Year(id: int.parse(yearId)),
//           date: Date(id: int.parse(dateId)),
//           weather: Weather(id: int.parse(weatherId)));
//       final res = await shiftUsecase.getShiftsByYearAndDateAndWeather(request.context, req);
//       return statusResponse.responseOK(jsonEncode(res));
//     } catch (e) {
//       Log.severe('shiftController.getShiftsByYearAndDateAndWeather: ${e.toString()}');
//       var json = jsonEncode({'message': e.toString()});
//       return statusResponse.responseBadRequest(json);
//     }
//   }
// }
