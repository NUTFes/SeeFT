package controller

import (
	"net/http"

	"github.com/NUTFes/SeeFT/api/lib/usecase"
	"github.com/labstack/echo/v4"
)

type placeController struct {
	u usecase.PlaceUseCase
}

type PlaceController interface {
	IndexPlace(echo.Context) error
	ShowPlace(echo.Context) error
	CreatePlace(echo.Context) error
	UpdatePlace(echo.Context) error
	DeletePlace(echo.Context) error
}

func NewPlaceController(u usecase.PlaceUseCase) PlaceController {
	return &placeController{u}
}

func (b *placeController) IndexPlace(c echo.Context) error {
	places, err := b.u.GetPlaces(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, places)
}

func (b *placeController) ShowPlace(c echo.Context) error {
	id := c.Param("id")
	place, err := b.u.GetPlaceByID(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, place)
}

// Create
func (u *placeController) CreatePlace(c echo.Context) error {
	place := c.QueryParam("place")
	remark := c.QueryParam("remark")
	latastPlace, err := u.u.CreatePlace(c.Request().Context(), place, remark)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, latastPlace)
}

// Update
func (u *placeController) UpdatePlace(c echo.Context) error {
	id := c.Param("id")
	place := c.QueryParam("place")
	remark := c.QueryParam("remark")
	updatedPlace, err := u.u.UpdatePlace(c.Request().Context(), id, place, remark)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, updatedPlace)
}

// Destroy
func (u *placeController) DeletePlace(c echo.Context) error {
	id := c.QueryParam("id")
	err := u.u.DeletePlace(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "Destroy Place")
}