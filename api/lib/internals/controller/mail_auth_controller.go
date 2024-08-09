package controller

import (
	"net/http"

	"github.com/NUTFes/SeeFT/api/lib/usecase"
	//"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/labstack/echo/v4"
)

type mailAuthController struct {
	u usecase.MailAuthUseCase
}

type MailAuthController interface {
	SignIn(echo.Context) error
	WebSignUp(echo.Context) error
	WebSignIn(echo.Context) error
	WebSignOut(echo.Context) error
	WebIsSignIn(echo.Context) error
}


func NewMailAuthController(u usecase.MailAuthUseCase) MailAuthController {
	return &mailAuthController{u}
}

// sign in
func (auth *mailAuthController) SignIn(c echo.Context) error {
	studentNumber := c.QueryParam("student_number")
	password := c.QueryParam("password")
	user, err := auth.u.SignIn(c.Request().Context(), studentNumber, password)
	if (err != nil) {
		return err
	}
	return c.JSON(http.StatusOK, user)
}

func (auth *mailAuthController) WebSignUp(c echo.Context) error {
	name := c.QueryParam("name")
	mail := c.QueryParam("mail")
	studentNumber := c.QueryParam("student_number")
	gradeID := c.QueryParam("grade_id")
	departmentID := c.QueryParam("department_id")
	bureauID := c.QueryParam("bureau_id")
	roleID := c.QueryParam("role_id")
	tel := c.QueryParam("tel")
	password := c.QueryParam("password")
	token, err := auth.u.WebSignUp(c.Request().Context(), name, mail, gradeID, departmentID, bureauID, roleID, studentNumber, tel, password)
	if (err != nil) {
		return err
	}
	return c.JSON(http.StatusOK, token)
}

func (auth *mailAuthController) WebSignIn(c echo.Context) error {
	studentNumber := c.QueryParam("student_number")
	password := c.QueryParam("password")
	token, err := auth.u.WebSignIn(c.Request().Context(), studentNumber, password)
	if (err != nil) {
		return err
	}
	return c.JSON(http.StatusOK, token)
}

func (auth *mailAuthController)WebSignOut(c echo.Context) error {
	// headerからトークンを取得する
	accessToken := c.Request().Header["Access-Token"][0]
	err := auth.u.WebSignOut(c.Request().Context(), accessToken)
	if (err != nil) {
		return err
	}
	return c.JSON(http.StatusOK, "Success Sign Out")
}

func (auth *mailAuthController) WebIsSignIn(c echo.Context) error {
	// headerからトークンを取得する
	accessToken := c.Request().Header["Access-Token"][0]
	isSignIn, err := auth.u.WebIsSignIn(c.Request().Context(), accessToken)
	if (err != nil) {
		return err
	}
	return c.JSON(http.StatusOK, isSignIn)
}
