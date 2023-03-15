// package controller

// import (
// 	"fmt"
// 	"net/http"
// 	"strconv"

// 	"github.com/NUTFes/seeft/service"
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