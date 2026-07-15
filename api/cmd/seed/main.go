package main

import (
	"log"
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

	log.Printf("seedを投入します: path=%s", seedFilePath)

	if err := applySeed(db); err != nil {
		log.Fatalf(
			"seedの投入に失敗しました: %v",
			err,
		)
	}

	log.Println("seedの投入が完了しました")
}
