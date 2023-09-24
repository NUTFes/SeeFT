package router

import (
	"github.com/NUTFes/SeeFT/api/lib/internals/controller"
	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo"
)

type router struct {
	healthcheckController     controller.HealthcheckController
	mailAuthController		  controller.MailAuthController
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
	mailAuthController controller.MailAuthController,
	bureauController controller.BureauController,
	shiftContoller controller.ShiftController,
	taskController controller.TaskController,
	timeController controller.TimeController,
	userController controller.UserController,
) Router {
	return router{
		healthController,
		mailAuthController,
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

	// mail auth
	e.POST("/mail_auth/signin", r.mailAuthController.SignIn)
	// e.GET("/mail_auth/signin/:student_number", r.mailAuthController.SignIn)

	// bureauのRoute
	e.GET("/bureaus", r.bureauController.IndexBureau)
	e.GET("/bureaus/:id", r.bureauController.ShowBureau)

	// shiftのRoute
	e.GET("/shifts", r.shiftController.IndexShift)
	e.GET("/shifts/:id", r.shiftController.ShowShift)
	e.GET("/shifts/users/:user_id", r.shiftController.ShowShiftsByUser)
	e.GET("/shifts/tasks/:task_id/years/:year_id/dates/:date_id/times/:time_id/weathers/:weather_id", r.shiftController.ShowUsersByShift)
	e.GET("/shifts/users/:user_id/dates/:date/weathers/:weather", r.shiftController.ShowShiftsByUserAndDateAndWeather)

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
	e.POST("/users/:name/:mail/:grade_id/:department_id/:bureau_id/:role_id", r.userController.CreateUser)
	e.PUT("/users/:id/:name/:mail/:grade_id/:department_id/:bureau_id/:role_id", r.userController.UpdateUser)
	e.DELETE("/users/:id", r.userController.DeleteUser)
}