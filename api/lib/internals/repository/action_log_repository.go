package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/pkg/errors"
)

type actionLogRepository struct {
	client db.Client
}

type ActionLogRepository interface {
	Create(ctx context.Context, shiftID, userID, dateID int, actionType string, diffPayload interface{}) error
	GetUnsentLogs(ctx context.Context) (*sql.Rows, error)
	GetUnsentLogsByUserAndDate(ctx context.Context, userID, dateID int) (*sql.Rows, error)
	MarkAsSent(ctx context.Context, logIDs []int) error
}

func NewActionLogRepository(c db.Client) ActionLogRepository {
	return &actionLogRepository{client: c}
}

// Create action_logにレコードを作成
func (r *actionLogRepository) Create(ctx context.Context, shiftID, userID, dateID int, actionType string, diffPayload interface{}) error {
	payloadBytes, err := json.Marshal(diffPayload)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal diff payload")
	}

	query := `
		INSERT INTO action_logs (shift_id, user_id, date_id, action_type, diff_payload, is_sent)
		VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = r.client.DB().ExecContext(ctx, query, shiftID, userID, dateID, actionType, payloadBytes, false)
	if err != nil {
		return errors.Wrapf(err, "failed to create action log")
	}

	return nil
}

// GetUnsentLogs 未送信のログを全件取得
func (r *actionLogRepository) GetUnsentLogs(ctx context.Context) (*sql.Rows, error) {
	query := `
		SELECT id, shift_id, user_id, date_id, action_type, diff_payload, is_sent, created_at
		FROM action_logs
		WHERE is_sent = false
		ORDER BY created_at ASC`

	rows, err := r.client.DB().QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get unsent logs")
	}

	return rows, nil
}

// GetUnsentLogsByUserAndDate ユーザーIDと日付IDで未送信ログを取得（グルーピング用）
func (r *actionLogRepository) GetUnsentLogsByUserAndDate(ctx context.Context, userID, dateID int) (*sql.Rows, error) {
	query := `
		SELECT id, shift_id, user_id, date_id, action_type, diff_payload, is_sent, created_at
		FROM action_logs
		WHERE is_sent = false
		AND user_id = $1
		AND date_id = $2
		ORDER BY created_at ASC`

	rows, err := r.client.DB().QueryContext(ctx, query, userID, dateID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get unsent logs by user and date")
	}

	return rows, nil
}

// MarkAsSent 指定されたログIDを送信済みとしてマーク
func (r *actionLogRepository) MarkAsSent(ctx context.Context, logIDs []int) error {
	if len(logIDs) == 0 {
		return nil
	}

	// IN句用のプレースホルダーを生成
	placeholders := make([]string, len(logIDs))
	args := make([]interface{}, len(logIDs))
	for i, id := range logIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		UPDATE action_logs
		SET is_sent = true
		WHERE id IN (%s)`, strings.Join(placeholders, ","))

	_, err := r.client.DB().ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrapf(err, "failed to mark logs as sent")
	}

	return nil
}
