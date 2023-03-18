package controller

import (
	"net/http"

	"github.com/NUTFes/SeeFT/api/lib/internals/usecase"
	"github.com/labstack/echo/v4"
)

type shiftController struct {
	u usecase.ShiftUseCase
}

type ShiftController interface {
	IndexShift(echo.Context) error
	ShowShift(echo.Context) error
	ShowShiftsByUser(echo.Context) error
	ShowShiftsByUserAndDateAndWeather(echo.Context) error
}

func NewShiftController(u usecase.ShiftUseCase) ShiftController {
	return &shiftController{u}
}

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
