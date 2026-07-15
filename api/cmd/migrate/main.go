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

	// 初期schemaを適用する。
	if err := applySchemas(db); err != nil {
		log.Fatalf(
			"schemaの適用に失敗しました: %v",
			err,
		)
	}

	// schema適用後、未適用のup migrationを適用する。
	if err := applyMigrations(config); err != nil {
		log.Fatalf(
			"migrationの適用に失敗しました: %v",
			err,
		)
	}

	log.Println("schemaとmigrationの適用が完了しました")
}
