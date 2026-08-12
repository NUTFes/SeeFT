package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const migrationSourceURL = "file:///db_data/db/migrations"

func applyMigrations(config dbConfig) error {
	migrator, err := newMigrator(config)
	if err != nil {
		return err
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

func downMigrations(config dbConfig, down string) error {
	migrator, err := newMigrator(config)
	if err != nil {
		return err
	}
	defer closeMigrator(migrator)

	if down == "all" {
		if err := migrator.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				log.Println("downするmigrationはありません")
				return nil
			}

			return fmt.Errorf(
				"migrationのdownに失敗しました: %w",
				err,
			)
		}

		return nil
	}

	downNum, err := strconv.Atoi(down)
	if err != nil {
		return fmt.Errorf(
			"downするmigration数が不正です: %s",
			down,
		)
	}
	if downNum <= 0 {
		return fmt.Errorf(
			"downするmigration数は1以上を指定してください: %d",
			downNum,
		)
	}

	appliedCount, err := getAppliedMigrationCount(migrator)
	if err != nil {
		return fmt.Errorf(
			"適用済みmigration数の取得に失敗しました: %w",
			err,
		)
	}
	if downNum > appliedCount {
		return fmt.Errorf(
			"指定されたmigration数が適用済みmigration数を超えています: down=%d applied=%d",
			downNum,
			appliedCount,
		)
	}

	if err := migrator.Steps(-downNum); err != nil {
		return fmt.Errorf(
			"migrationのdownに失敗しました: %w",
			err,
		)
	}

	return nil
}

func newMigrator(config dbConfig) (*migrate.Migrate, error) {
	databaseURL := buildDatabaseURL(config)

	migrator, err := migrate.New(
		migrationSourceURL,
		databaseURL,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"migrationの初期化に失敗しました: source=%s: %w",
			migrationSourceURL,
			err,
		)
	}

	return migrator, nil
}

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

const migrationDirPath = "/db_data/db/migrations"

func getAppliedMigrationCount(migrator *migrate.Migrate) (int, error) {
	currentVersion, _, err := migrator.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			return 0, nil
		}

		return 0, fmt.Errorf(
			"現在のmigration versionの取得に失敗しました: %w",
			err,
		)
	}

	entries, err := os.ReadDir(migrationDirPath)
	if err != nil {
		return 0, fmt.Errorf(
			"migrationディレクトリの読み込みに失敗しました: %w",
			err,
		)
	}

	count := 0

	for _, entry := range entries {
		name := entry.Name()

		if !strings.HasSuffix(name, ".up.sql") {
			continue
		}

		versionStr, _, ok := strings.Cut(name, "_")
		if !ok {
			continue
		}

		version, err := strconv.ParseUint(versionStr, 10, 64)
		if err != nil {
			continue
		}

		if version <= uint64(currentVersion) {
			count++
		}
	}

	return count, nil
}
