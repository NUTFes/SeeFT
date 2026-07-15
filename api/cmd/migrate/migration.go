package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const migrationSourceURL = "file:///db_data/db/migrations"

// applyMigrationsは、未適用のup migrationをversionの昇順で適用する。
func applyMigrations(config dbConfig) error {
	databaseURL := buildDatabaseURL(config)

	migrator, err := migrate.New(
		migrationSourceURL,
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf(
			"migrationの初期化に失敗しました: source=%s: %w",
			migrationSourceURL,
			err,
		)
	}
	defer closeMigrator(migrator)

	log.Printf(
		"migrationを開始します: source=%s",
		migrationSourceURL,
	)

	if err := migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			version, dirty, versionErr := migrator.Version()
			if versionErr != nil {
				log.Println("未適用のmigrationはありません")
				return nil
			}

			log.Printf(
				"未適用のmigrationはありません: version=%d dirty=%t",
				version,
				dirty,
			)

			return nil
		}

		version, dirty, versionErr := migrator.Version()
		if versionErr == nil {
			return fmt.Errorf(
				"migrationの適用に失敗しました: version=%d dirty=%t: %w",
				version,
				dirty,
				err,
			)
		}

		return fmt.Errorf(
			"migrationの適用に失敗しました: %w",
			err,
		)
	}

	version, dirty, err := migrator.Version()
	if err != nil {
		return fmt.Errorf(
			"migration適用後のversion取得に失敗しました: %w",
			err,
		)
	}

	log.Printf(
		"migrationが完了しました: version=%d dirty=%t",
		version,
		dirty,
	)

	return nil
}

// closeMigratorは、migration sourceとDB接続を終了する。
func closeMigrator(migrator *migrate.Migrate) {
	sourceErr, databaseErr := migrator.Close()

	if sourceErr != nil {
		log.Printf(
			"migration sourceの終了に失敗しました: %v",
			sourceErr,
		)
	}

	if databaseErr != nil {
		log.Printf(
			"migration用DB接続の終了に失敗しました: %v",
			databaseErr,
		)
	}
}
