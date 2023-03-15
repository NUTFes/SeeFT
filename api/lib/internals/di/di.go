package di

import (
	"log"

	"github.com/NUTFes/SeeFT/api/lib/drivers/db"
	"github.com/NUTFes/SeeFT/api/lib/drivers/server"
	"github.com/NUTFes/SeeFT/api/lib/externals/controller"
	"github.com/NUTFes/SeeFT/api/lib/externals/repository"
	"github.com/NUTFes/SeeFT/api/lib/externals/repository/external"
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

	// UseCase
	bureauUseCase := usecase.NewBureauUseCase(bureauRepository)

	// Controller
	healthcheckController := controller.NewHealthCheckController()
	bureauController := controller.NewBureauController(bureauUseCase)

	// router
	router := router.NewRouter(
		healthcheckController,
		bureauController,
	)

	// Server
	server.RunServer(router)

	return client
}