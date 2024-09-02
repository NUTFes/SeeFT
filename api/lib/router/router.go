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
	gradeController						controller.GradeController
	placeController						controller.PlaceController
	departmentController      controller.DepartmentController
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
	gradeController controller.GradeController,
	placeController controller.PlaceController,
	departmentController controller.DepartmentController,
	shiftContoller controller.ShiftController,
	taskController controller.TaskController,
	timeController controller.TimeController,
	userController controller.UserController,
) Router {
	return router{
		healthController,
		mailAuthController,
		bureauController,
		gradeController,
		placeController,
		departmentController,
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
	e.POST("/mail_auth/web_signin", r.mailAuthController.WebSignIn)
	e.POST("/mail_auth/web_signup", r.mailAuthController.WebSignUp)
	e.DELETE("/mail_auth/web_signout", r.mailAuthController.WebSignOut)
	e.GET("/mail_auth/web_is_signin", r.mailAuthController.WebIsSignIn)
	// e.GET("/mail_auth/signin/:student_number", r.mailAuthController.SignIn)

	// bureauのRoute
	e.GET("/bureaus", r.bureauController.IndexBureau)
	e.GET("/bureaus/:id", r.bureauController.ShowBureau)

	// gradeのRoute
	e.GET("/grades", r.gradeController.IndexGrade)
	e.GET("/grades/:id", r.gradeController.ShowGrade)

	// placeのRoute
	e.GET("/places", r.placeController.IndexPlace)
	e.GET("/places/:id", r.placeController.ShowPlace)
	e.POST("/places", r.placeController.CreatePlace)
	e.PUT("/places/:id", r.placeController.UpdatePlace)
	e.DELETE("/places", r.placeController.DeletePlace)

	// departmentのRoute
	e.GET("/departments", r.departmentController.IndexDepartment)
	e.GET("/departments/:id", r.departmentController.ShowDepartment)

	// shift(スマホ)のRoute
	e.GET("/shifts", r.shiftController.IndexShift)
	e.GET("/shifts/:id", r.shiftController.ShowShift)
	e.GET("/shifts/users/:user_id", r.shiftController.ShowShiftsByUser)
	e.GET("/shifts/tasks/:task_id/years/:year_id/dates/:date_id/times/:time_id/weathers/:weather_id", r.shiftController.ShowUsersByShift)
	e.GET("/shifts/users/:user_id/dates/:date/weathers/:weather", r.shiftController.ShowShiftsByUserAndDateAndWeather)

	// shift(Web)のRoute
	e.GET("/shifts-admin", r.shiftController.IndexShiftAdmin)
	e.GET("/shifts-admin/:id", r.shiftController.ShowShiftAdmin)
	e.POST("/shifts-admin", r.shiftController.CreateShiftAdmin)
	e.PUT("/shifts-admin/:id", r.shiftController.UpdateShiftAdmin)
	e.DELETE("/shifts-admin", r.shiftController.DeleteShiftAdmin)
	e.GET("/shifts-admin/dates/:date/weathers/:weather", r.shiftController.ShowShiftAdminByDateAndWeather)

	// taskのRoute
	e.GET("/tasks", r.taskController.IndexTask)
	e.GET("/tasks/:id", r.taskController.ShowTask)
	e.GET("/tasks/shifts/:shift", r.taskController.ShowTasksByShift)
	e.POST("/tasks", r.taskController.CreateTask)
	e.PUT("/tasks/:id", r.taskController.UpdateTask)
	e.DELETE("/tasks", r.taskController.DeleteTask)

	// timeのRoute
	e.GET("/times", r.timeController.IndexTime)
	e.GET("/times/:id", r.timeController.ShowTime)

	// users
	e.GET("/users", r.userController.IndexUser)
	e.GET("/users/:id", r.userController.ShowUser)
	e.POST("/users", r.userController.CreateUser)
	e.PUT("/users/:id", r.userController.UpdateUser)
	e.DELETE("/users", r.userController.DeleteUser)

	// current_user
	e.GET("/current_user", r.userController.GetCurrentUser)
}