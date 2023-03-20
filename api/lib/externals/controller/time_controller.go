package controller

import (
	"net/http"

	"github.com/NUTFes/SeeFT/api/lib/internals/usecase"
	"github.com/labstack/echo/v4"
)

type timeController struct {
	u usecase.TimeUseCase
}

type TimeController interface {
	IndexTime(echo.Context) error
	ShowTime(echo.Context) error
}

func NewTimeController(u usecase.TimeUseCase) TimeController {
	return &timeController{u}
}


func (b *timeController) IndexTime(c echo.Context) error {
	times, err := b.u.GetTimes(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, times)
}

func (b *timeController) ShowTime(c echo.Context) error {
	id := c.Param("id")
	time, err := b.u.GetTimeByID(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, time)
}

// import 'dart:convert';
// import 'package:shelf/shelf.dart';

// import '../../config/config.dart';
// import '../../entity/entity.dart';
// import '../../usecase/usecase.dart';

// class TimeController {
//   final StatusResponse statusResponse;
//   final TimeUsecase timeUsecase;

//   TimeController(
//     this.statusResponse,
//     this.timeUsecase,
//   );

//   Future<Response> getTimes(Request request) async {
//     try {
//       List<Time> times = await timeUsecase.getTimes(request.context);
//       return statusResponse.responseOK(jsonEncode(times));
//     } catch (e) {
//       Log.severe('timeController.getTimes: ${e.toString()}');
//       var json = jsonEncode({'message': e.toString()});
//       return statusResponse.responseBadRequest(json);
//     }
//   }
// }
