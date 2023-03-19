package controller

import (
	"github.com/NUTFes/SeeFT/api/lib/internals/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
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
	gradeID := c.QueryParam("grade_id")
	bureauID := c.QueryParam("bureau_id")
	roleID := c.QueryParam("role_id")
	latastUser, err := u.u.CreateUser(c.Request().Context(), name, gradeID, bureauID, roleID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, latastUser)
}

// Update
func (u *userController) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	name := c.QueryParam("name")
	gradeID := c.QueryParam("grade_id")
	bureauID := c.QueryParam("bureau_id")
	roleID := c.QueryParam("role_id")
	updatedUser, err := u.u.UpdateUser(c.Request().Context(), id, name, gradeID, bureauID, roleID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, updatedUser)
}

// Destroy
func (u *userController) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	err := u.u.DeleteUser(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "Destroy User")
}


// import (
// 	"fmt"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// type UserController struct{}

// func (controller UserController) Index(c *gin.Context) {

// 	var service service.UserService
// 	p, err := service.GetAll()

// 	if err != nil {
// 		c.AbortWithStatus(http.StatusBadRequest)
// 		fmt.Println(err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, p)
// }

// func (pc UserController) FindByID(c *gin.Context) {
// 	id, _ := strconv.Atoi(c.Param("id"))

// 	var s service.UserService
// 	p, err := s.FindByID(id)

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": fmt.Sprintf("%s", err),
// 		})
// 		fmt.Println(err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, p)

// }