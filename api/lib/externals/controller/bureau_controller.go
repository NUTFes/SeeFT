package controller

import (
	"net/http"

	"github.com/NUTFes/SeeFT/api/lib/internals/usecase"
	"github.com/labstack/echo/v4"
)

type bureauController struct {
	u usecase.BureauUseCase
}

type BureauController interface {
	IndexBureau(echo.Context) error
	ShowBureau(echo.Context) error
	CreateBureau(echo.Context) error
	UpdateBureau(echo.Context) error
	DestroyBureau(echo.Context) error
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

func (b *bureauController) CreateBureau(c echo.Context) error {
	name := c.QueryParam("name")

	latastBureau, err := b.u.CreateBureau(c.Request().Context(), name)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, latastBureau)
}

func (b *bureauController) UpdateBureau(c echo.Context) error {
	id := c.Param("id")
	name := c.QueryParam("name")

	updatedBureau, err := b.u.UpdateBureau(c.Request().Context(), id, name)

	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, updatedBureau)
}

func (b *bureauController) DestroyBureau(c echo.Context) error {
	id := c.Param("id")
	err := b.u.DestroyBureau(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "Destroy Bureau")
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
