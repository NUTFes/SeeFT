package di

import (
	"log"

	"github.com/NUTFes/SeeFT/api/lib/drivers/db"
	"github.com/NUTFes/SeeFT/api/lib/drivers/server"
	"github.com/NUTFes/SeeFT/api/lib/externals/controller"
	"github.com/NUTFes/SeeFT/api/lib/externals/repository"
	"github.com/NUTFes/SeeFT/api/lib/externals/repository/abstract"
	"github.com/NUTFes/SeeFT/api/lib/internals/usecase"
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
	bureauRepository := repository.NewBureauRepository(client, crud)
	shiftRepository := repository.NewShiftRepository(client, crud)
	taskRepository := repository.NewTaskRepository(client, crud)
	timeRepository := repository.NewTimeRepository(client, crud)
	userRepository := repository.NewUserRepository(client, crud)

	// UseCase
	bureauUseCase := usecase.NewBureauUseCase(bureauRepository)
	shiftUseCase := usecase.NewShiftUseCase(shiftRepository)
	taskUseCase := usecase.NewTaskUseCase(taskRepository)
	timeUsecase := usecase.NewTimeUseCase(timeRepository)
	userUseCase := usecase.NewUserUseCase(userRepository)

	// Controller
	healthcheckController := controller.NewHealthCheckController()
	bureauController := controller.NewBureauController(bureauUseCase)
	shiftContoller := controller.NewShiftController(shiftUseCase)
	taskController := controller.NewTaskController(taskUseCase)
	timeController := controller.NewTimeController(timeUsecase)
	userController := controller.NewUserController(userUseCase)

	// router
	router := router.NewRouter(
		healthcheckController,
		bureauController,
		shiftContoller,
		taskController,
		timeController,
		userController,
	)

	// Server
	server.RunServer(router)

	return client
}