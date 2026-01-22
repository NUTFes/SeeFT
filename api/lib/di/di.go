package di

import (
	"context"
	"log"
	"time"

	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/NUTFes/SeeFT/api/lib/externals/server"
	"github.com/NUTFes/SeeFT/api/lib/externals/slack"
	"github.com/NUTFes/SeeFT/api/lib/internals/controller"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
	"github.com/NUTFes/SeeFT/api/lib/router"
	"github.com/NUTFes/SeeFT/api/lib/usecase"
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
	placeRepository := repository.NewPlaceRepository(client, crud)
	departmentRepository := repository.NewDepartmentRepository(client, crud)
	shiftRepository := repository.NewShiftRepository(client, crud)
	shiftCardRepository := repository.NewShiftCardRepository(client)
	taskRepository := repository.NewTaskRepository(client, crud)
	timeRepository := repository.NewTimeRepository(client, crud)
	userRepository := repository.NewUserRepository(client, crud)
	yearRepository := repository.NewYearRepository(client, crud)
	dateRepository := repository.NewDateRepository(client, crud)
	weatherRepository := repository.NewWeatherRepository(client, crud)
	questionRescueRepository := repository.NewQuestionRescueRepository(client, crud)
	shorthandedRescueRepository := repository.NewShorthandedRescueRepository(client, crud)
	troubleRescueRepository := repository.NewTroubleRescueRepository(client, crud)
	reviewRepository := repository.NewReviewRepository(client, crud)
	actionLogRepository := repository.NewActionLogRepository(client)

	// SlackService
	slackService, err := slack.NewSlackService()
	if err != nil {
		log.Printf("Warning: Failed to initialize SlackService: %v. Notifications will not be sent.", err)
		slackService = nil
	}

	// UseCase
	mailAuthUseCase := usecase.NewAuthUseCase(userRepository, sessionRepository)
	bureauUseCase := usecase.NewBureauUseCase(bureauRepository)
	gradeUseCase := usecase.NewGradeUseCase(gradeRepository)
	placeUseCase := usecase.NewPlaceUseCase(placeRepository)
	departmentUseCase := usecase.NewDepartmentUseCase(departmentRepository)
	shiftUseCase := usecase.NewShiftUseCase(shiftRepository, shiftCardRepository, taskRepository, userRepository, yearRepository, dateRepository, timeRepository, weatherRepository, placeRepository, gradeRepository, bureauRepository, actionLogRepository)
	taskUseCase := usecase.NewTaskUseCase(taskRepository, placeRepository)
	timeUsecase := usecase.NewTimeUseCase(timeRepository)
	userUseCase := usecase.NewUserUseCase(userRepository, sessionRepository)
	questionRescueUseCase := usecase.NewQuestionRescueUseCase(questionRescueRepository)
	shorthandedRescueUseCase := usecase.NewShorthandedRescueUseCase(shorthandedRescueRepository)
	troubleRescueUseCase := usecase.NewTroubleRescueUseCase(troubleRescueRepository)
	rescueUnifiedUseCase := usecase.NewRescueUnifiedUseCase(questionRescueRepository, shorthandedRescueRepository, troubleRescueRepository, userRepository, taskRepository)
	reviewUseCase := usecase.NewReviewUseCase(reviewRepository, taskRepository)

	// NotificationUseCase (SlackServiceが初期化されている場合のみ)
	var notificationUseCase usecase.NotificationUseCase
	if slackService != nil {
		notificationUseCase = usecase.NewNotificationUseCase(
			actionLogRepository,
			slackService,
			userRepository,
			dateRepository,
			timeRepository,
			taskRepository,
			shiftRepository,
			weatherRepository,
		)

		// 通知ワーカーを起動（バックグラウンドで実行）
		go func() {
			ticker := time.NewTicker(5 * time.Minute)
			defer ticker.Stop()
			ctx := context.Background()
			for range ticker.C {
				if err := notificationUseCase.ProcessUnsentNotifications(ctx); err != nil {
					log.Printf("Failed to process notifications: %v", err)
				}
			}
		}()
	}

	// Controller
	healthcheckController := controller.NewHealthCheckController()
	mailAuthController := controller.NewMailAuthController(mailAuthUseCase)
	bureauController := controller.NewBureauController(bureauUseCase)
	gradeController := controller.NewGradeController(gradeUseCase)
	placeController := controller.NewPlaceController(placeUseCase)
	departmentController := controller.NewDepartmentController(departmentUseCase)
	shiftContoller := controller.NewShiftController(shiftUseCase)
	taskController := controller.NewTaskController(taskUseCase)
	timeController := controller.NewTimeController(timeUsecase)
	userController := controller.NewUserController(userUseCase)
	questionRescueController := controller.NewQuestionRescueController(questionRescueUseCase)
	shorthandedRescueController := controller.NewShorthandedRescueController(shorthandedRescueUseCase)
	troubleRescueController := controller.NewTroubleRescueController(troubleRescueUseCase)
	rescueUnifiedController := controller.NewRescueUnifiedController(questionRescueUseCase, shorthandedRescueUseCase, troubleRescueUseCase, rescueUnifiedUseCase, userUseCase, taskUseCase, gradeUseCase, bureauUseCase)
	reviewController := controller.NewReviewController(reviewUseCase)

	// router
	router := router.NewRouter(
		healthcheckController,
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
		reviewController,
	)

	// Server
	server.RunServer(router)

	return client
}
