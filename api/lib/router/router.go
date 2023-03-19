package router

import (
	"github.com/NUTFes/SeeFT/api/lib/externals/controller"
	"github.com/labstack/echo/v4"
)

type router struct {
	healthcheckController     controller.HealthcheckController
	bureauController          controller.BureauController
	shiftController 		  		controller.ShiftController
	taskController			  		controller.TaskController
	timeController			  		controller.TimeController
	userController						controller.UserController
}

type Router interface {
	ProvideRouter(*echo.Echo)
}

func NewRouter(
	healthController controller.HealthcheckController,
	bureauController controller.BureauController,
	shiftContoller controller.ShiftController,
	taskController controller.TaskController,
	timeController controller.TimeController,
	userController controller.UserController,
) Router {
	return router{
		healthController,
		bureauController,
		shiftContoller,
		taskController,
		timeController,
		userController,
	}
}

func (r router) ProvideRouter(e *echo.Echo) {
	// Healthcheck
	e.GET("/", r.healthcheckController.IndexHealthcheck)

	// bureauのRoute
	e.GET("/bureaus", r.bureauController.IndexBureau)
	e.GET("/bureaus/:id", r.bureauController.ShowBureau)

	// shiftのRoute
	e.GET("/shifts", r.shiftController.IndexShift)
	e.GET("/shifts/:id", r.shiftController.ShowShift)
	e.GET("/shifts/users/:user_id", r.shiftController.ShowShiftsByUser)
	e.GET("/shifts/users/:user_id/date/:date/weather/:weather", r.shiftController.ShowShiftsByUserAndDateAndWeather)

	// taskのRoute
	e.GET("/tasks", r.taskController.IndexTask)
	e.GET("/tasks/:id", r.taskController.ShowTask)
	e.GET("/tasks/shifts/:shift", r.taskController.ShowTasksByShift)

	// timeのRoute
	e.GET("/times", r.timeController.IndexTime)
	e.GET("/times/:id", r.timeController.ShowTime)

	// users
	e.GET("/users", r.userController.IndexUser)
	e.GET("/users/:id", r.userController.ShowUser)
	e.POST("/users", r.userController.CreateUser)
	e.PUT("/users/:id", r.userController.UpdateUser)
	e.DELETE("/users/:id", r.userController.DeleteUser)
}