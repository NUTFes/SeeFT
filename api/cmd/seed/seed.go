package main

import (
	"database/sql"
	"fmt"
	"os"
)

const seedFilePath = "/db_data/db/seed.sql"

func applySeed(db *sql.DB) error {
	content, err := os.ReadFile(seedFilePath)
	if err != nil {
		return fmt.Errorf("seedファイルを読み込めません: path=%s: %w", seedFilePath, err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("トランザクションを開始できません: %w", err)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	if _, err := tx.Exec(string(content)); err != nil {
		return fmt.Errorf("seed SQLの実行に失敗しました: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("seedのコミットに失敗しました: %w", err)
	}

	return nil
}
