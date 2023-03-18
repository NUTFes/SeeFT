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
}

type Router interface {
	ProvideRouter(*echo.Echo)
}

func NewRouter(
	healthController controller.HealthcheckController,
	bureauController controller.BureauController,
	shiftContoller controller.ShiftController,
	taskController controller.TaskController,
) Router {
	return router{
		healthController,
		bureauController,
		shiftContoller,
		taskController,
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
	e.GET("/shifts/user/:user_id", r.shiftController.ShowShiftsByUser)
	e.GET("/shifts/user/:user_id/date/:date/weather/:weather", r.shiftController.ShowShiftsByUserAndDateAndWeather)

	// taskのRoute
	e.GET("/tasks", r.taskController.IndexTask)
	e.GET("/tasks/:id", r.taskController.ShowTask)
	e.GET("/tasks/shift/:shift", r.taskController.ShowTasksByShift)
}