package controller

import (
	"net/http"

	"github.com/NUTFes/SeeFT/api/lib/usecase"
	"github.com/labstack/echo/v4"
)

type bureauController struct {
	u usecase.BureauUseCase
}

type BureauController interface {
	IndexBureau(echo.Context) error
	ShowBureau(echo.Context) error
}

func NewBureauController(u usecase.BureauUseCase) BureauController {
	return &bureauController{u}
}

func (b *bureauController) IndexBureau(c echo.Context) error {
	bureaus, err := b.u.GetBureaus(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, bureaus)
}

func (b *bureauController) ShowBureau(c echo.Context) error {
	id := c.Param("id")
	bureau, err := b.u.GetBureauByID(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, bureau)
}


// import 'dart:convert';
// import 'package:shelf/shelf.dart';

// import '../../config/config.dart';
// import '../../entity/entity.dart';
// import '../../usecase/usecase.dart';

// class BureauController {
//   final StatusResponse statusResponse;
//   final BureauUsecase bureauUsecase;

//   BureauController(
//     this.statusResponse,
//     this.bureauUsecase,
//   );

//   getBureaus(Request request) async {
//     try {
//       List<Bureau> bureaus = await bureauUsecase.getBureaus(request.context);
//       return statusResponse.responseOK(jsonEncode(bureaus));
//     } catch (e) {
//       Log.severe('BureauController.getBureaus: ${e.toString()}');
//       var json = jsonEncode({'message': e.toString()});
//       return statusResponse.responseBadRequest(json);
//     }
//   }
// }
