package controller

import (
	"net/http"

	"github.com/NUTFes/SeeFT/api/lib/entity"
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
	CreateTask(echo.Context) error
	UpdateTask(echo.Context) error
	DeleteTask(echo.Context) error
	UpdateTasksAndPlacesFromGAS(echo.Context) error
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

// Create
func (u *taskController) CreateTask(c echo.Context) error {
	task := c.QueryParam("task")
	placeID := c.QueryParam("place_id")
	url := c.QueryParam("url")
	bureauID := c.QueryParam("bureau_id")
	maxMember := c.QueryParam("max_member")
	color := c.QueryParam("color")
	remark := c.QueryParam("remark")
	yearID := c.QueryParam("year_id")
	latastTask, err := u.u.CreateTask(c.Request().Context(), task, placeID, url, bureauID, maxMember, color, remark, yearID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, latastTask)
}

// Update
func (u *taskController) UpdateTask(c echo.Context) error {
	id := c.Param("id")
	task := c.QueryParam("task")
	placeID := c.QueryParam("place_id")
	url := c.QueryParam("url")
	bureauID := c.QueryParam("bureau_id")
	maxMember := c.QueryParam("max_member")
	color := c.QueryParam("color")
	remark := c.QueryParam("remark")
	yearID := c.QueryParam("year_id")
	updatedTask, err := u.u.UpdateTask(c.Request().Context(), id, task, placeID, url, bureauID, maxMember, color, remark, yearID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, updatedTask)
}

// Destroy
func (u *taskController) DeleteTask(c echo.Context) error {
	id := c.QueryParam("id")
	err := u.u.DeleteTask(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "Destroy Task")
}

// GASからのタスクと集合場所変更通知を受け取るエンドポイント
func (sc *taskController) UpdateTasksAndPlacesFromGAS(c echo.Context) error {
	var req entity.TaskAndPlaceChangeRequest
	if err := c.Bind(&req); err != nil {
		// return c.JSON(http.StatusBadRequest, "Invalid request")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "無効なリクエスト形式"})
	}

	// 必要に応じてユースケース層へ処理を委譲
	if err := sc.u.UpdateTasksAndPlacesFromGAS(c.Request().Context(), req); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Tasks and Places updated successfully")
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
