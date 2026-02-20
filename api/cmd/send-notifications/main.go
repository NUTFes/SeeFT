// 未送信の action_logs をその場で1回だけ処理し、Slack に通知を送る。
// 5分待たずに通知を確認したいときに実行する。
//
// 実行方法（api ディレクトリで）:
//
//   # ホストから実行する場合（DB が localhost:5432 のとき）
//   NUTMEG_DB_HOST=localhost go run ./cmd/send-notifications
//
//   # 環境変数は api/env/dev.env を読み込む。未読み込みなら事前に export する。
//   # または API コンテナ内で実行（そのまま DB に接続できる）:
//   docker exec -it nutfes-seeft-api sh -c "cd /app && go run ./cmd/send-notifications"
package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/NUTFes/SeeFT/api/lib/externals/slack"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
	"github.com/NUTFes/SeeFT/api/lib/usecase"
)

func main() {
	_ = godotenv.Load("env/dev.env")

	client, err := db.ConnectMySQL()
	if err != nil {
		log.Fatalf("DB connect: %v", err)
	}
	defer client.CloseDB()

	crud := abstract.NewCrud(client)
	actionLogRepo := repository.NewActionLogRepository(client)
	userRepo := repository.NewUserRepository(client, crud)
	dateRepo := repository.NewDateRepository(client, crud)
	timeRepo := repository.NewTimeRepository(client, crud)
	taskRepo := repository.NewTaskRepository(client, crud)
	shiftRepo := repository.NewShiftRepository(client, crud)
	weatherRepo := repository.NewWeatherRepository(client, crud)

	slackService, err := slack.NewSlackService()
	if err != nil {
		log.Fatalf("Slack init: %v", err)
	}

	uc := usecase.NewNotificationUseCase(
		actionLogRepo,
		slackService,
		userRepo,
		dateRepo,
		timeRepo,
		taskRepo,
		shiftRepo,
		weatherRepo,
	)

	ctx := context.Background()
	if err := uc.ProcessUnsentNotifications(ctx); err != nil {
		log.Fatalf("ProcessUnsentNotifications: %v", err)
	}

	log.Println("Done. Check Slack.")
	os.Exit(0)
}
