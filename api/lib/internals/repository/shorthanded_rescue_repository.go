package repository

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
)

type shorthandedRescueRepository struct {
	client db.Client
	crud   abstract.Crud
}

type ShorthandedRescueRepository interface {
	All(context.Context) (*sql.Rows, error)
	Find(context.Context, string) (*sql.Row, error)
	FindByUserID(context.Context, string) (*sql.Rows, error)
	FindByTaskID(context.Context, string) (*sql.Rows, error)
	Create(context.Context, string, string, string, string, string) error
	Update(context.Context, string, string, string) error
	Delete(context.Context, string) error
	FindNewRecord(context.Context) (*sql.Row, error)
}

func NewShorthandedRescueRepository(c db.Client, ac abstract.Crud) ShorthandedRescueRepository {
	return &shorthandedRescueRepository{c, ac}
}

// 全件取得
func (sr *shorthandedRescueRepository) All(c context.Context) (*sql.Rows, error) {
	query := "SELECT * FROM shorthanded_rescues ORDER BY created_at DESC"
	return sr.crud.Read(c, query)
}

// 1件取得
func (sr *shorthandedRescueRepository) Find(c context.Context, id string) (*sql.Row, error) {
	query := "SELECT * FROM shorthanded_rescues WHERE id = $1"
	return sr.client.DB().QueryRowContext(c, query, id), nil
}

// ユーザIDで検索
func (sr *shorthandedRescueRepository) FindByUserID(c context.Context, userID string) (*sql.Rows, error) {
	query := "SELECT * FROM shorthanded_rescues WHERE user_id = $1 ORDER BY created_at DESC"
	return sr.client.DB().QueryContext(c, query, userID)
}

// タスクIDで検索
func (sr *shorthandedRescueRepository) FindByTaskID(c context.Context, taskID string) (*sql.Rows, error) {
	query := "SELECT * FROM shorthanded_rescues WHERE task_id = $1 ORDER BY created_at DESC"
	return sr.client.DB().QueryContext(c, query, taskID)
}

// 作成（セキュリティ強化：プレースホルダーを使用）
func (sr *shorthandedRescueRepository) Create(c context.Context, userID string, taskID string, missingNumber string, place string, status string) error {
	query := `
		INSERT INTO shorthanded_rescues (user_id, task_id, missing_number, place, status, time, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}
	
	taskIDInt, err := strconv.Atoi(taskID)
	if err != nil {
		return err
	}
	
	missingNumberInt, err := strconv.Atoi(missingNumber)
	if err != nil {
		return err
	}
	
	now := time.Now()
	var placePtr *string
	if place != "" {
		placePtr = &place
	}
	
	_, err = sr.client.DB().ExecContext(c, query, userIDInt, taskIDInt, missingNumberInt, placePtr, status, now, now, now)
	return err
}

// 更新（レスポンスとステータスを更新）
func (sr *shorthandedRescueRepository) Update(c context.Context, id string, status string, response string) error {
	query := `
		UPDATE shorthanded_rescues
		SET status = $1, response = $2, updated_at = $3
		WHERE id = $4`
	
	now := time.Now()
	var responsePtr *string
	if response != "" {
		responsePtr = &response
	}
	
	_, err := sr.client.DB().ExecContext(c, query, status, responsePtr, now, id)
	return err
}

// 削除
func (sr *shorthandedRescueRepository) Delete(c context.Context, id string) error {
	query := "DELETE FROM shorthanded_rescues WHERE id = $1"
	_, err := sr.client.DB().ExecContext(c, query, id)
	return err
}

// 最新レコード取得
func (sr *shorthandedRescueRepository) FindNewRecord(c context.Context) (*sql.Row, error) {
	query := "SELECT * FROM shorthanded_rescues ORDER BY id DESC LIMIT 1"
	return sr.crud.ReadByID(c, query)
}
