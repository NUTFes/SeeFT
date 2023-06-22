package controller

import (
	"net/http"

	"github.com/NUTFes/SeeFT/api/lib/usecase"
	"github.com/labstack/echo/v4"
)

type taskController struct {
	u usecase.TaskUseCase
}

type TaskController interface {
	IndexTask(echo.Context) error
	ShowTask(echo.Context) error
	ShowTasksByShift(echo.Context) error
}

func NewTaskController(u usecase.TaskUseCase) TaskController {
	return &taskController{u}
}


func (b *taskController) IndexTask(c echo.Context) error {
	tasks, err := b.u.GetTasks(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, tasks)
}

func (b *taskController) ShowTask(c echo.Context) error {
	id := c.Param("id")
	task, err := b.u.GetTaskByID(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, task)
}

func (b *taskController) ShowTasksByShift(c echo.Context) error {
	shift := c.Param("shift")
	tasks, err := b.u.GetTasksByShift(c.Request().Context(), shift)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, tasks)
}

// type TaskController struct{}

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