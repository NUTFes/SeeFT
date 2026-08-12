package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"
)

const (
	seedFilePath         = "/db_data/db/seed.sql"
	seedExecutionTimeout = 60 * time.Second
)

func applySeed(db *sql.DB) error {
	content, err := os.ReadFile(seedFilePath)
	if err != nil {
		return fmt.Errorf(
			"seedファイルを読み込めません: path=%s: %w",
			seedFilePath,
			err,
		)
	}

	if len(content) == 0 {
		return fmt.Errorf(
			"seedファイルが空です: path=%s",
			seedFilePath,
		)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		seedExecutionTimeout,
	)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf(
			"seed用トランザクションを開始できません: %w",
			err,
		)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	if _, err := tx.ExecContext(ctx, string(content)); err != nil {
		return fmt.Errorf(
			"seed SQLの実行に失敗しました: path=%s: %w",
			seedFilePath,
			err,
		)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf(
			"seedのコミットに失敗しました: %w",
			err,
		)
	}

	committed = true

	return nil
}
