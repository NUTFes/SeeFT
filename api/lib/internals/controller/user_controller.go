package controller

import (
	"net/http"

	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/NUTFes/SeeFT/api/lib/usecase"
	"github.com/labstack/echo/v4"
)

type userController struct {
	u usecase.UserUseCase
}

type UserController interface {
	IndexUser(echo.Context) error
	ShowUser(echo.Context) error
	CreateUser(echo.Context) error
	UpdateUser(echo.Context) error
	DeleteUser(echo.Context) error
	GetCurrentUser(echo.Context) error
	UpdateUsersFromGAS(echo.Context) error
}

func NewUserController(u usecase.UserUseCase) UserController {
	return &userController{u}
}

// Index
func (u *userController) IndexUser(c echo.Context) error {
	users, err := u.u.GetUsers(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, users)
}

// Show
func (u *userController) ShowUser(c echo.Context) error {
	id := c.Param("id")
	user, err := u.u.GetUserByID(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)
}

// Create
func (u *userController) CreateUser(c echo.Context) error {
	name := c.QueryParam("name")
	mail := c.QueryParam("mail")
	studentNumber := c.QueryParam("student_number")
	gradeID := c.QueryParam("grade_id")
	departmentID := c.QueryParam("department_id")
	bureauID := c.QueryParam("bureau_id")
	roleID := c.QueryParam("role_id")
	tel := c.QueryParam("tel")
	password := c.QueryParam("password")
	latastUser, err := u.u.CreateUser(c.Request().Context(), name, mail, gradeID, departmentID, bureauID, roleID, studentNumber, tel, password)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, latastUser)
}

// Update
func (u *userController) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	name := c.QueryParam("name")
	mail := c.QueryParam("mail")
	gradeID := c.QueryParam("grade_id")
	departmentID := c.QueryParam("department_id")
	bureauID := c.QueryParam("bureau_id")
	roleID := c.QueryParam("role_id")
	studentNumber := c.QueryParam("student_number")
	tel := c.QueryParam("tel")
	updatedUser, err := u.u.UpdateUser(c.Request().Context(), id, name, mail, gradeID, departmentID, bureauID, roleID, studentNumber, tel)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, updatedUser)
}

// Destroy
func (u *userController) DeleteUser(c echo.Context) error {
	id := c.QueryParam("id")
	err := u.u.DeleteUser(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "Destroy User")
}

// ログインユーザーの取得
func (auth *userController) GetCurrentUser(c echo.Context) error {
	// headerからトークンを取得する
	accessToken := c.Request().Header["Access-Token"][0]
	user, err := auth.u.GetCurrentUser(c.Request().Context(), accessToken)
	if err != nil {
		return c.JSON(http.StatusNotFound, user)
	} else {
		return c.JSON(http.StatusOK, user)
	}
}

// GASからの名簿変更通知を受け取るエンドポイント
func (sc *userController) UpdateUsersFromGAS(c echo.Context) error {
	var req entity.UserChangeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}

	// 必要に応じてユースケース層へ処理を委譲
	if err := sc.u.UpdateUsersFromGAS(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Users updated successfully")
}
