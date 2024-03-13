package di

import (
	"log"

	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/NUTFes/SeeFT/api/lib/externals/server"
	"github.com/NUTFes/SeeFT/api/lib/internals/controller"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
	"github.com/NUTFes/SeeFT/api/lib/usecase"
	"github.com/NUTFes/SeeFT/api/lib/router"
)

func InitializeServer() db.Client {
	// DB接続
	client, err := db.ConnectMySQL()
	if err != nil {
		log.Fatal("db error")
	}

	crud := abstract.NewCrud(client)

	// ↓
	// Repository
	sessionRepository := repository.NewSessionRepository(client)
	bureauRepository := repository.NewBureauRepository(client, crud)
	gradeRepository := repository.NewGradeRepository(client, crud)
	shiftRepository := repository.NewShiftRepository(client, crud)
	taskRepository := repository.NewTaskRepository(client, crud)
	timeRepository := repository.NewTimeRepository(client, crud)
	userRepository := repository.NewUserRepository(client, crud)
	yearRepository := repository.NewYearRepository(client, crud)
	dateRepository := repository.NewDateRepository(client, crud)
	weatherRepository := repository.NewWeatherRepository(client, crud)

	// UseCase
	mailAuthUseCase := usecase.NewAuthUseCase(userRepository, sessionRepository)
	bureauUseCase := usecase.NewBureauUseCase(bureauRepository)
	gradeUseCase := usecase.NewGradeUseCase(gradeRepository)
	shiftUseCase := usecase.NewShiftUseCase(shiftRepository, taskRepository, userRepository, yearRepository, dateRepository, timeRepository, weatherRepository)
	taskUseCase := usecase.NewTaskUseCase(taskRepository)
	timeUsecase := usecase.NewTimeUseCase(timeRepository)
	userUseCase := usecase.NewUserUseCase(userRepository)

	// Controller
	healthcheckController := controller.NewHealthCheckController()
	mailAuthController := controller.NewMailAuthController(mailAuthUseCase)
	bureauController := controller.NewBureauController(bureauUseCase)
	gradeController := controller.NewGradeController(gradeUseCase)
	shiftContoller := controller.NewShiftController(shiftUseCase)
	taskController := controller.NewTaskController(taskUseCase)
	timeController := controller.NewTimeController(timeUsecase)
	userController := controller.NewUserController(userUseCase)

	// router
	router := router.NewRouter(
		healthcheckController,
		mailAuthController,
		bureauController,
		gradeController,
		shiftContoller,
		taskController,
		timeController,
		userController,
	)

	// Server
	server.RunServer(router)

	return client
}