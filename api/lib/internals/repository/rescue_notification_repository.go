package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/pkg/errors"
)

type rescueNotificationRepository struct {
	client db.Client
}

type RescueNotificationRepository interface {
	Create(ctx context.Context, n entity.RescueNotification) error
	FindUnsent(ctx context.Context) (*sql.Rows, error)
	MarkAsSent(ctx context.Context, ids []int) error
}

func NewRescueNotificationRepository(c db.Client) RescueNotificationRepository {
	return &rescueNotificationRepository{client: c}
}

// Create 未送信の通知を1件積む
func (r *rescueNotificationRepository) Create(ctx context.Context, n entity.RescueNotification) error {
	query := `
		INSERT INTO rescue_notifications (rescue_type, rescue_id, user_id, kind, status, response, is_sent)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	var responsePtr *string
	if n.Response != "" {
		responsePtr = &n.Response
	}

	_, err := r.client.DB().ExecContext(ctx, query, n.RescueType, n.RescueID, n.UserID, n.Kind, n.Status, responsePtr, false)
	if err != nil {
		return errors.Wrapf(err, "failed to create rescue notification")
	}
	return nil
}

// FindUnsent 未送信の通知を古い順に全件取得
func (r *rescueNotificationRepository) FindUnsent(ctx context.Context) (*sql.Rows, error) {
	query := `
		SELECT id, rescue_type, rescue_id, user_id, kind, status, response, is_sent, created_at
		FROM rescue_notifications
		WHERE is_sent = false
		ORDER BY created_at ASC, id ASC`

	rows, err := r.client.DB().QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get unsent rescue notifications")
	}
	return rows, nil
}

// MarkAsSent 指定された通知を送信済みにする
func (r *rescueNotificationRepository) MarkAsSent(ctx context.Context, ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		UPDATE rescue_notifications
		SET is_sent = true
		WHERE id IN (%s)`, strings.Join(placeholders, ","))

	if _, err := r.client.DB().ExecContext(ctx, query, args...); err != nil {
		return errors.Wrapf(err, "failed to mark rescue notifications as sent")
	}
	return nil
}
