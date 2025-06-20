package repository

import (
	"context"
	"database/sql"

	// "github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
)

type rescueRepository struct {
	client db.Client
	crud   abstract.Crud
}

type RescueRepository interface {
	// トラブルレスキュー
	CreateTroubleRescue(context.Context, int, int, string, string) error
	FindTroubleRescueByID(context.Context, int) (*sql.Row, error)
	FindTroubleRescuesByUserID(context.Context, int) (*sql.Rows, error)
	AllTroubleRescues(context.Context) (*sql.Rows, error)
	UpdateTroubleRescue(context.Context, int, string, string) error

	// 質問レスキュー
	CreateQuestionRescue(context.Context, int, string) error
	FindQuestionRescueByID(context.Context, int) (*sql.Row, error)
	FindQuestionRescuesByUserID(context.Context, int) (*sql.Rows, error)
	AllQuestionRescues(context.Context) (*sql.Rows, error)
	UpdateQuestionRescue(context.Context, int, string, string) error

	// 人手不足レスキュー
	CreateShorthandedRescue(context.Context, int, int, int, string) error
	FindShorthandedRescueByID(context.Context, int) (*sql.Row, error)
	FindShorthandedRescuesByUserID(context.Context, int) (*sql.Rows, error)
	AllShorthandedRescues(context.Context) (*sql.Rows, error)
	UpdateShorthandedRescue(context.Context, int, string, string) error
}

func NewRescueRepository(c db.Client, ac abstract.Crud) RescueRepository {
	return &rescueRepository{c, ac}
}

// トラブルレスキュー作成
func (r *rescueRepository) CreateTroubleRescue(c context.Context, userID int, taskID int, place string, detail string) error {
	query := `INSERT INTO trouble_rescues (user_id, task_id, place, detail)
	         VALUES ($1, $2, $3, $4)`
	_, err := r.client.DB().ExecContext(c, query, userID, taskID, place, detail)
	return err
}

// トラブルレスキュー1件取得
func (r *rescueRepository) FindTroubleRescueByID(c context.Context, id int) (*sql.Row, error) {
	query := "SELECT * FROM trouble_rescues WHERE id = $1"
	return r.client.DB().QueryRowContext(c, query, id), nil
}

// ユーザーのトラブルレスキュー取得
func (r *rescueRepository) FindTroubleRescuesByUserID(c context.Context, userID int) (*sql.Rows, error) {
	query := "SELECT * FROM trouble_rescues WHERE user_id = $1 ORDER BY time DESC"
	return r.client.DB().QueryContext(c, query, userID)
}

// 全トラブルレスキュー取得
func (r *rescueRepository) AllTroubleRescues(c context.Context) (*sql.Rows, error) {
	query := "SELECT * FROM trouble_rescues ORDER BY time DESC"
	return r.crud.Read(c, query)
}

// トラブルレスキュー更新
func (r *rescueRepository) UpdateTroubleRescue(c context.Context, id int, status string, response string) error {
	query := "UPDATE trouble_rescues SET status = $1, response = $2 WHERE id = $3"
	return r.crud.UpdateDB(c, query)
}

// 質問レスキュー作成
func (r *rescueRepository) CreateQuestionRescue(c context.Context, userID int, question string) error {
	query := `INSERT INTO question_rescues (user_id, question)
	         VALUES ($1, $2)`
	_, err := r.client.DB().ExecContext(c, query, userID, question)
	return err
}

// 質問レスキュー1件取得
func (r *rescueRepository) FindQuestionRescueByID(c context.Context, id int) (*sql.Row, error) {
	query := "SELECT * FROM question_rescues WHERE id = $1"
	return r.client.DB().QueryRowContext(c, query, id), nil
}

// ユーザーの質問レスキュー取得
func (r *rescueRepository) FindQuestionRescuesByUserID(c context.Context, userID int) (*sql.Rows, error) {
	query := "SELECT * FROM question_rescues WHERE user_id = $1 ORDER BY time DESC"
	return r.client.DB().QueryContext(c, query, userID)
}

// 全質問レスキュー取得
func (r *rescueRepository) AllQuestionRescues(c context.Context) (*sql.Rows, error) {
	query := "SELECT * FROM question_rescues ORDER BY time DESC"
	return r.crud.Read(c, query)
}

// 質問レスキュー更新
func (r *rescueRepository) UpdateQuestionRescue(c context.Context, id int, status string, response string) error {
	query := "UPDATE question_rescues SET status = $1, response = $2 WHERE id = $3"
	_, err := r.client.DB().ExecContext(c, query, status, response, id)
	return err
}

// 人手不足レスキュー作成
func (r *rescueRepository) CreateShorthandedRescue(c context.Context, userID int, taskID int, missingNumber int, place string) error {
	query := `INSERT INTO shorthanded_rescues (user_id, task_id, missing_number, place)
	         VALUES ($1, $2, $3, $4)`
	_, err := r.client.DB().ExecContext(c, query, userID, taskID, missingNumber, place)
	return err
}

// 人手不足レスキュー1件取得
func (r *rescueRepository) FindShorthandedRescueByID(c context.Context, id int) (*sql.Row, error) {
	query := "SELECT * FROM shorthanded_rescues WHERE id = $1"
	return r.client.DB().QueryRowContext(c, query, id), nil
}

// ユーザーの人手不足レスキュー取得
func (r *rescueRepository) FindShorthandedRescuesByUserID(c context.Context, userID int) (*sql.Rows, error) {
	query := "SELECT * FROM shorthanded_rescues WHERE user_id = $1 ORDER BY time DESC"
	return r.client.DB().QueryContext(c, query, userID)
}

// 全人手不足レスキュー取得
func (r *rescueRepository) AllShorthandedRescues(c context.Context) (*sql.Rows, error) {
	query := "SELECT * FROM shorthanded_rescues ORDER BY time DESC"
	return r.crud.Read(c, query)
}

// 人手不足レスキュー更新
func (r *rescueRepository) UpdateShorthandedRescue(c context.Context, id int, status string, response string) error {
	query := "UPDATE shorthanded_rescues SET status = $1, response = $2 WHERE id = $3"
	_, err := r.client.DB().ExecContext(c, query, status, response, id)
	return err
}
