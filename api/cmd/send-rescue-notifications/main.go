// 未送信の rescue_notifications をその場で1回だけ処理し、Slack に通知を送る。
// 30秒待たずに通知を確認したいときに実行する。
// 書き込みから15秒経っていない通知は、まとめ送りのため次の回に回される。
//
// 実行方法（api ディレクトリで）:
//
//	# ホストから実行する場合（DB が localhost:5432 のとき）
//	NUTMEG_DB_HOST=localhost go run ./cmd/send-rescue-notifications
//
//	# または API コンテナ内で実行（そのまま DB に接続できる）:
//	docker exec -it nutfes-seeft-api sh -c "cd /app && go run ./cmd/send-rescue-notifications"
package main

import (
	"context"
	"log"

	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/NUTFes/SeeFT/api/lib/externals/slack"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
	"github.com/NUTFes/SeeFT/api/lib/usecase"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("env/dev.env")

	client, err := db.ConnectMySQL()
	if err != nil {
		log.Fatalf("DB connect: %v", err)
	}
	defer client.CloseDB()

	crud := abstract.NewCrud(client)
	rescueNotificationRepo := repository.NewRescueNotificationRepository(client)

	slackService, err := slack.NewSlackService()
	if err != nil {
		log.Fatalf("Slack init: %v", err)
	}

	// 送信側では更新を記録しないので、各レスキューのusecaseには通知の記録先を渡さない
	uc := usecase.NewRescueNotificationUseCase(
		rescueNotificationRepo,
		slackService,
		usecase.NewQuestionRescueUseCase(repository.NewQuestionRescueRepository(client, crud), nil),
		usecase.NewShorthandedRescueUseCase(repository.NewShorthandedRescueRepository(client, crud), nil),
		usecase.NewTroubleRescueUseCase(repository.NewTroubleRescueRepository(client, crud), nil),
		repository.NewTaskRepository(client, crud),
		repository.NewUserRepository(client, crud),
	)

	if err := uc.ProcessUnsentRescueNotifications(context.Background()); err != nil {
		log.Fatalf("ProcessUnsentRescueNotifications: %v", err)
	}

	log.Println("Done. Check Slack.")
}
