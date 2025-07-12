package repository

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
)

type troubleRescueRepository struct {
	client db.Client
	crud   abstract.Crud
}

type TroubleRescueRepository interface {
	All(context.Context) (*sql.Rows, error)
	Find(context.Context, string) (*sql.Row, error)
	FindByUserID(context.Context, string) (*sql.Rows, error)
	FindByTaskID(context.Context, string) (*sql.Rows, error)
	Create(context.Context, string, string, string, string, string) error
	Update(context.Context, string, string, string) error
	Delete(context.Context, string) error
	FindNewRecord(context.Context) (*sql.Row, error)
}

func NewTroubleRescueRepository(c db.Client, ac abstract.Crud) TroubleRescueRepository {
	return &troubleRescueRepository{c, ac}
}

// 全件取得
func (tr *troubleRescueRepository) All(c context.Context) (*sql.Rows, error) {
	query := "SELECT * FROM trouble_rescues ORDER BY created_at DESC"
	return tr.crud.Read(c, query)
}

// 1件取得
func (tr *troubleRescueRepository) Find(c context.Context, id string) (*sql.Row, error) {
	query := "SELECT * FROM trouble_rescues WHERE id = $1"
	return tr.client.DB().QueryRowContext(c, query, id), nil
}

// ユーザIDで検索
func (tr *troubleRescueRepository) FindByUserID(c context.Context, userID string) (*sql.Rows, error) {
	query := "SELECT * FROM trouble_rescues WHERE user_id = $1 ORDER BY created_at DESC"
	return tr.client.DB().QueryContext(c, query, userID)
}

// タスクIDで検索
func (tr *troubleRescueRepository) FindByTaskID(c context.Context, taskID string) (*sql.Rows, error) {
	query := "SELECT * FROM trouble_rescues WHERE task_id = $1 ORDER BY created_at DESC"
	return tr.client.DB().QueryContext(c, query, taskID)
}

// 作成（セキュリティ強化：プレースホルダーを使用）
func (tr *troubleRescueRepository) Create(c context.Context, userID string, taskID string, place string, detail string, status string) error {
	query := `
		INSERT INTO trouble_rescues (user_id, task_id, place, detail, status, time, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}
	
	taskIDInt, err := strconv.Atoi(taskID)
	if err != nil {
		return err
	}
	
	now := time.Now()
	var placePtr *string
	if place != "" {
		placePtr = &place
	}
	
	_, err = tr.client.DB().ExecContext(c, query, userIDInt, taskIDInt, placePtr, detail, status, now, now, now)
	return err
}

// 更新（レスポンスとステータスを更新）
func (tr *troubleRescueRepository) Update(c context.Context, id string, status string, response string) error {
	query := `
		UPDATE trouble_rescues
		SET status = $1, response = $2, updated_at = $3
		WHERE id = $4`
	
	now := time.Now()
	var responsePtr *string
	if response != "" {
		responsePtr = &response
	}
	
	_, err := tr.client.DB().ExecContext(c, query, status, responsePtr, now, id)
	return err
}

// 削除
func (tr *troubleRescueRepository) Delete(c context.Context, id string) error {
	query := "DELETE FROM trouble_rescues WHERE id = $1"
	_, err := tr.client.DB().ExecContext(c, query, id)
	return err
}

// 最新レコード取得
func (tr *troubleRescueRepository) FindNewRecord(c context.Context) (*sql.Row, error) {
	query := "SELECT * FROM trouble_rescues ORDER BY id DESC LIMIT 1"
	return tr.crud.ReadByID(c, query)
}
