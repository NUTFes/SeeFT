package controller

import (
	"net/http"

	"github.com/NUTFes/SeeFT/api/lib/internals/usecase"
	"github.com/labstack/echo/v4"
)

type mailAuthController struct {
	u usecase.MailAuthUseCase
}

type MailAuthController interface {
	SignIn(echo.Context) error
}

func NewMailAuthController(u usecase.MailAuthUseCase) MailAuthController {
	return &mailAuthController{u}
}

// sign in
func (auth *mailAuthController) SignIn(c echo.Context) error {
	email := c.QueryParam("email")
	token, err := auth.u.SignIn(c.Request().Context(), email)
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, token)
	return nil
}