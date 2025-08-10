package router

import (
	"github.com/NUTFes/SeeFT/api/lib/internals/controller"
	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo"
)

type router struct {
	healthcheckController           controller.HealthcheckController
	mailAuthController              controller.MailAuthController
	bureauController                controller.BureauController
	gradeController                 controller.GradeController
	placeController                 controller.PlaceController
	departmentController            controller.DepartmentController
	shiftController                 controller.ShiftController
	taskController                  controller.TaskController
	timeController                  controller.TimeController
	userController                  controller.UserController
	questionRescueController        controller.QuestionRescueController
	shorthandedRescueController     controller.ShorthandedRescueController
	troubleRescueController         controller.TroubleRescueController
	rescueUnifiedController         controller.RescueUnifiedController
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
	questionRescueController controller.QuestionRescueController,
	shorthandedRescueController controller.ShorthandedRescueController,
	troubleRescueController controller.TroubleRescueController,
	rescueUnifiedController controller.RescueUnifiedController,
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
		questionRescueController,
		shorthandedRescueController,
		troubleRescueController,
		rescueUnifiedController,
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
	e.GET("/shift-cards/users/:user_id/dates/:date_id/weathers/:weather_id", r.shiftController.ShowShiftCardsByUserAndDateAndWeather)
	e.POST("/shift-cards", r.shiftController.PostShiftCards)

	// shift(Web)のRoute
	e.GET("/shifts-admin", r.shiftController.IndexShiftAdmin)
	e.GET("/shifts-admin/:id", r.shiftController.ShowShiftAdmin)
	e.POST("/shifts-admin", r.shiftController.CreateShiftAdmin)
	e.PUT("/shifts-admin/:id", r.shiftController.UpdateShiftAdmin)
	e.DELETE("/shifts-admin", r.shiftController.DeleteShiftAdmin)
	e.GET("/shifts-admin/dates/:date/weathers/:weather", r.shiftController.ShowShiftAdminByDateAndWeather)
	e.GET("/shifts-admin/dates/:date/weathers/:weather/lower/:lower/upper/:upper", r.shiftController.ShowShiftAdminByDateAndWeatherAndTime)
	e.GET("/shifts-admin/max-id", r.shiftController.SerachMaxID)

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

	// rescue（新統一エンドポイント）
	e.POST("/rescues", r.rescueUnifiedController.CreateRescue)
	e.GET("/rescues", r.rescueUnifiedController.GetAllRescues)
	e.GET("/rescues/users/:user_id", r.rescueUnifiedController.GetRescuesByUserID)

	// question rescues
	e.GET("/question-rescues", r.questionRescueController.IndexQuestionRescue)
	e.GET("/question-rescues/:id", r.questionRescueController.ShowQuestionRescue)
	e.GET("/question-rescues/users/:user_id", r.questionRescueController.ShowQuestionRescuesByUserID)
	e.POST("/question-rescues", r.questionRescueController.CreateQuestionRescue)
	e.PUT("/question-rescues/:id", r.questionRescueController.UpdateQuestionRescue)
	e.DELETE("/question-rescues", r.questionRescueController.DeleteQuestionRescue)

	// shorthanded rescues
	e.GET("/shorthanded-rescues", r.shorthandedRescueController.IndexShorthandedRescue)
	e.GET("/shorthanded-rescues/:id", r.shorthandedRescueController.ShowShorthandedRescue)
	e.GET("/shorthanded-rescues/users/:user_id", r.shorthandedRescueController.ShowShorthandedRescuesByUserID)
	e.GET("/shorthanded-rescues/tasks/:task_id", r.shorthandedRescueController.ShowShorthandedRescuesByTaskID)
	e.POST("/shorthanded-rescues", r.shorthandedRescueController.CreateShorthandedRescue)
	e.PUT("/shorthanded-rescues/:id", r.shorthandedRescueController.UpdateShorthandedRescue)
	e.DELETE("/shorthanded-rescues", r.shorthandedRescueController.DeleteShorthandedRescue)

	// trouble rescues
	e.GET("/trouble-rescues", r.troubleRescueController.IndexTroubleRescue)
	e.GET("/trouble-rescues/:id", r.troubleRescueController.ShowTroubleRescue)
	e.GET("/trouble-rescues/users/:user_id", r.troubleRescueController.ShowTroubleRescuesByUserID)
	e.GET("/trouble-rescues/tasks/:task_id", r.troubleRescueController.ShowTroubleRescuesByTaskID)
	e.POST("/trouble-rescues", r.troubleRescueController.CreateTroubleRescue)
	e.PUT("/trouble-rescues/:id", r.troubleRescueController.UpdateTroubleRescue)
	e.DELETE("/trouble-rescues", r.troubleRescueController.DeleteTroubleRescue)

	// shiftの希望日程
	e.POST("/request_shifts", r.shiftController.SubmitShift)

	// GAS用のRoute
	e.POST("/api/update_users", r.userController.UpdateUsersFromGAS)                     // ユーザの更新
	e.POST("/api/update_tasks_and_places", r.taskController.UpdateTasksAndPlacesFromGAS) // タスクと場所の更新
	e.POST("/api/update_shifts", r.shiftController.UpdateShiftsFromGAS)                  // シフトの更新
}
