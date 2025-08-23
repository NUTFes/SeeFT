package repository

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
)

type questionRescueRepository struct {
	client db.Client
	crud   abstract.Crud
}

type QuestionRescueRepository interface {
	All(context.Context) (*sql.Rows, error)
	Find(context.Context, string) (*sql.Row, error)
	FindByUserID(context.Context, string) (*sql.Rows, error)
	Create(context.Context, string, string, string) error
	Update(context.Context, string, string, string) error
	Delete(context.Context, string) error
	FindNewRecord(context.Context) (*sql.Row, error)
}

func NewQuestionRescueRepository(c db.Client, ac abstract.Crud) QuestionRescueRepository {
	return &questionRescueRepository{c, ac}
}

// 全件取得
func (qr *questionRescueRepository) All(c context.Context) (*sql.Rows, error) {
	query := "SELECT * FROM question_rescues ORDER BY created_at DESC"
	return qr.crud.Read(c, query)
}

// 1件取得
func (qr *questionRescueRepository) Find(c context.Context, id string) (*sql.Row, error) {
	query := "SELECT * FROM question_rescues WHERE id = $1"
	return qr.client.DB().QueryRowContext(c, query, id), nil
}

// ユーザIDで検索
func (qr *questionRescueRepository) FindByUserID(c context.Context, userID string) (*sql.Rows, error) {
	query := "SELECT * FROM question_rescues WHERE user_id = $1 ORDER BY created_at DESC"
	return qr.client.DB().QueryContext(c, query, userID)
}

// 作成（セキュリティ強化：プレースホルダーを使用）
func (qr *questionRescueRepository) Create(c context.Context, userID string, question string, status string) error {
	query := `
		INSERT INTO question_rescues (user_id, question, status, time, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}
	
	now := time.Now()
	_, err = qr.client.DB().ExecContext(c, query, userIDInt, question, status, now, now, now)
	return err
}

// 更新（レスポンスとステータスを更新）
func (qr *questionRescueRepository) Update(c context.Context, id string, status string, response string) error {
	query := `
		UPDATE question_rescues
		SET status = $1, response = $2, updated_at = $3
		WHERE id = $4`
	
	now := time.Now()
	_, err := qr.client.DB().ExecContext(c, query, status, response, now, id)
	return err
}

// 削除
func (qr *questionRescueRepository) Delete(c context.Context, id string) error {
	query := "DELETE FROM question_rescues WHERE id = $1"
	_, err := qr.client.DB().ExecContext(c, query, id)
	return err
}

// 最新レコード取得
func (qr *questionRescueRepository) FindNewRecord(c context.Context) (*sql.Row, error) {
	query := "SELECT * FROM question_rescues ORDER BY id DESC LIMIT 1"
	return qr.crud.ReadByID(c, query)
}
