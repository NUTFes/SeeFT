package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	schemaDirectory        = "/db_data/db/schema"
	schemaExecutionTimeout = 60 * time.Second
)

const createSchemaHistoryTableSQL = `
CREATE TABLE IF NOT EXISTS seeft_schema_history (
	file_name  TEXT PRIMARY KEY,
	checksum   TEXT NOT NULL,
	applied_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
)
`

var schemaFilePattern = regexp.MustCompile(`(?i)^create(\d+).*\.sql$`)

type schemaFile struct {
	name  string
	path  string
	order int
}

// applySchemasは、schemaディレクトリ内のSQLを所定の順序で適用する。
func applySchemas(db *sql.DB) error {
	if err := ensureSchemaHistoryTable(db); err != nil {
		return err
	}

	files, err := findSchemaFiles(schemaDirectory)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf(
			"schemaファイルが見つかりません: directory=%s",
			schemaDirectory,
		)
	}

	for _, file := range files {
		if err := applySchemaFile(db, file); err != nil {
			return err
		}
	}

	return nil
}

// ensureSchemaHistoryTableは、schemaファイルの適用履歴を管理するテーブルを作成する。
func ensureSchemaHistoryTable(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		schemaExecutionTimeout,
	)
	defer cancel()

	if _, err := db.ExecContext(ctx, createSchemaHistoryTableSQL); err != nil {
		return fmt.Errorf(
			"schema適用履歴テーブルの作成に失敗しました: %w",
			err,
		)
	}

	return nil
}

// findSchemaFilesは、create<数字>*.sqlを取得し、数字順に並べる。
func findSchemaFiles(directory string) ([]schemaFile, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf(
			"schemaディレクトリの読み込みに失敗しました: directory=%s: %w",
			directory,
			err,
		)
	}

	files := make([]schemaFile, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		// SQL以外のファイルは対象外とする。
		if !strings.EqualFold(filepath.Ext(name), ".sql") {
			continue
		}

		matches := schemaFilePattern.FindStringSubmatch(name)
		if matches == nil {
			return nil, fmt.Errorf(
				"schemaファイル名が規則に一致しません: file=%s, expected=create<数字>*.sql",
				name,
			)
		}

		order, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, fmt.Errorf(
				"schemaファイルの実行順を取得できません: file=%s: %w",
				name,
				err,
			)
		}

		files = append(files, schemaFile{
			name:  name,
			path:  filepath.Join(directory, name),
			order: order,
		})
	}

	sort.Slice(files, func(i, j int) bool {
		if files[i].order != files[j].order {
			return files[i].order < files[j].order
		}

		return files[i].name < files[j].name
	})

	return files, nil
}

// applySchemaFileは、未適用のschemaファイルをトランザクション内で適用する。
func applySchemaFile(db *sql.DB, file schemaFile) error {
	content, err := os.ReadFile(file.path)
	if err != nil {
		return fmt.Errorf(
			"schemaファイルの読み込みに失敗しました: file=%s: %w",
			file.name,
			err,
		)
	}

	checksum := calculateChecksum(content)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		schemaExecutionTimeout,
	)
	defer cancel()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf(
			"トランザクションの開始に失敗しました: file=%s: %w",
			file.name,
			err,
		)
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	appliedChecksum, err := findAppliedChecksum(ctx, tx, file.name)
	switch {
	case err == nil:
		if appliedChecksum != checksum {
			return fmt.Errorf(
				"適用済みschemaファイルの内容が変更されています: file=%s",
				file.name,
			)
		}

		log.Printf(
			"schemaをスキップします: order=%d file=%s",
			file.order,
			file.name,
		)

		return nil

	case errors.Is(err, sql.ErrNoRows):
		// 未適用なので、この後でSQLを実行する。

	default:
		return fmt.Errorf(
			"schema適用履歴の確認に失敗しました: file=%s: %w",
			file.name,
			err,
		)
	}

	log.Printf(
		"schemaを適用します: order=%d file=%s",
		file.order,
		file.name,
	)

	if _, err := tx.ExecContext(ctx, string(content)); err != nil {
		return fmt.Errorf(
			"schema SQLの実行に失敗しました: file=%s: %w",
			file.name,
			err,
		)
	}

	if err := insertSchemaHistory(
		ctx,
		tx,
		file.name,
		checksum,
	); err != nil {
		return fmt.Errorf(
			"schema適用履歴の登録に失敗しました: file=%s: %w",
			file.name,
			err,
		)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf(
			"schema適用結果のコミットに失敗しました: file=%s: %w",
			file.name,
			err,
		)
	}

	committed = true

	log.Printf(
		"schemaを適用しました: order=%d file=%s",
		file.order,
		file.name,
	)

	return nil
}

// findAppliedChecksumは、適用済みファイルのチェックサムを取得する。
func findAppliedChecksum(
	ctx context.Context,
	tx *sql.Tx,
	fileName string,
) (string, error) {
	const query = `
SELECT checksum
FROM seeft_schema_history
WHERE file_name = $1
`

	var checksum string

	if err := tx.QueryRowContext(ctx, query, fileName).Scan(&checksum); err != nil {
		return "", err
	}

	return checksum, nil
}

// insertSchemaHistoryは、正常に適用したschemaファイルを履歴へ登録する。
func insertSchemaHistory(
	ctx context.Context,
	tx *sql.Tx,
	fileName string,
	checksum string,
) error {
	const query = `
INSERT INTO seeft_schema_history (
	file_name,
	checksum
)
VALUES ($1, $2)
`

	if _, err := tx.ExecContext(
		ctx,
		query,
		fileName,
		checksum,
	); err != nil {
		return err
	}

	return nil
}

// calculateChecksumは、schemaファイルの内容からSHA-256を生成する。
func calculateChecksum(content []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(content))
}
