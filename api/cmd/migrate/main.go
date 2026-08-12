package main

import (
	"log"
	"os"
	"time"
)

func main() {
	config, err := loadDBConfig()
	if err != nil {
		log.Fatalf(
			"データベース設定の読み込みに失敗しました: %v",
			err,
		)
	}

	db, err := openDB(config)
	if err != nil {
		log.Fatalf(
			"データベース接続の初期化に失敗しました: %v",
			err,
		)
	}
	defer closeDB(db)

	if err := pingDB(db, 10*time.Second); err != nil {
		log.Fatalf(
			"データベースへの接続に失敗しました: %v",
			err,
		)
	}

	command := "up"
	if len(os.Args) >= 2 {
		command = os.Args[1]
	}

	switch command {
	case "up":
		if err := applySchemas(db); err != nil {
			log.Fatalf(
				"schemaの適用に失敗しました: %v",
				err,
			)
		}

		if err := applyMigrations(config); err != nil {
			log.Fatalf(
				"migrationのupに失敗しました: %v",
				err,
			)
		}

		log.Println("schemaとmigrationのupが完了しました")

	case "down":
		if len(os.Args) < 3 {
			log.Fatal("downするmigration数、またはallを指定してください")
		}

		downArg := os.Args[2]

		if err := downMigrations(config, downArg); err != nil {
			log.Fatalf(
				"migrationのdownに失敗しました: %v",
				err,
			)
		}

		log.Println("migrationのdownが完了しました")

	default:
		log.Fatalf("不明なコマンドです: %s", command)
	}
}
