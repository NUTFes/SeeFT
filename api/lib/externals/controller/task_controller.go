// package controller

// import (
// 	"fmt"
// 	"net/http"
// 	"strconv"

// 	"github.com/NUTFes/SeeFT/api/lib/service"
// 	"github.com/gin-gonic/gin"
// )

// type TaskController struct{}

package controller

// func (controller TaskController) TaskDetail(c *gin.Context) {
// 	userID, _ := strconv.Atoi(c.Param("userID"))
// 	day := c.Param("day")
// 	weather := c.Param("weather")
// 	workID, _ := strconv.Atoi(c.Param("workID"))
// 	time := c.Param("time")

// 	var service service.WorkService

// 	p, err := service.TaskWithUser(userID, day, weather, workID, time)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": fmt.Sprintf("%s", err),
// 		})
// 		fmt.Println(err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, p)
// }

// func (controller WorkController) TaskList(c *gin.Context) {
// 	var service service.WorkService

// 	p, err := service.WorkList()
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": fmt.Sprintf("%s", err),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, p)
// }