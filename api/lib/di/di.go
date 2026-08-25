package di

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/NUTFes/SeeFT/api/lib/externals/scheduler"
	"github.com/NUTFes/SeeFT/api/lib/externals/server"
	"github.com/NUTFes/SeeFT/api/lib/externals/slack"
	"github.com/NUTFes/SeeFT/api/lib/internals/controller"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
	"github.com/NUTFes/SeeFT/api/lib/router"
	"github.com/NUTFes/SeeFT/api/lib/usecase"
)

func InitializeServer(ctx context.Context) (db.Client, error) {
	// DB接続
	client, err := db.ConnectMySQL()
	if err != nil {
		return nil, err
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

	// マニュアル配信（nutfes限定, issue #444）: 設定値は環境変数から読み取るのみでロジックは持たない
	manualDir := os.Getenv("MANUAL_DIR")
	if manualDir == "" {
		manualDir = "/manuals"
	}
	manualClientID := os.Getenv("MANUAL_OAUTH_CLIENT_ID")
	manualClientSecret := os.Getenv("MANUAL_OAUTH_CLIENT_SECRET")
	manualRedirectURL := os.Getenv("MANUAL_OAUTH_REDIRECT_URL")
	// アップロード(PUT /manuals/:id)用トークン（issue #448）。閲覧機能の有効条件（OAuth3変数）とは独立しており、
	// 未設定でも閲覧機能はそのまま動作する（アップロードのみusecase側でfail-closedに無効化される）
	manualUploadToken := os.Getenv("MANUAL_UPLOAD_TOKEN")

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

	// MANUAL_OAUTH_* が揃っていない場合はマニュアル配信を無効化する。
	// 特にシークレットが空だとHMAC署名が空鍵になりCookie偽造で認証を素通りできるため、
	// 設定不備のまま /manuals ルートを公開してはならない（router側はnilで未登録になる）
	var manualController controller.ManualController
	if manualClientID != "" && manualClientSecret != "" && manualRedirectURL != "" {
		manualUseCase := usecase.NewManualUseCase(usecase.ManualConfig{
			ClientID:     manualClientID,
			ClientSecret: manualClientSecret,
			RedirectURL:  manualRedirectURL,
			ManualDir:    manualDir,
			UploadToken:  manualUploadToken,
		})
		manualController = controller.NewManualController(manualUseCase)
	} else {
		log.Printf("manual_gate: MANUAL_OAUTH_CLIENT_ID / MANUAL_OAUTH_CLIENT_SECRET / MANUAL_OAUTH_REDIRECT_URL が未設定のため、マニュアル配信ルートを無効化します")
	}

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
		manualController,
	)

	// Scheduler: 5分間隔で未送信通知を flush する（goroutine で起動し即 return）
	slackService, err := slack.NewSlackService()
	if err != nil {
		log.Printf("slack init failed, notification scheduler disabled: %v", err)
	}
	if err == nil {
		notificationUseCase := usecase.NewNotificationUseCase(
			actionLogRepository, slackService,
			userRepository, dateRepository, timeRepository,
			taskRepository, shiftRepository, weatherRepository,
		)
		scheduler.New("notification", 5*time.Minute, notificationUseCase.ProcessUnsentNotifications).Start(ctx)

	}

	// Server
	if err := server.RunServer(ctx, router); err != nil {
		return client, err
	}

	return client, nil
}
