package repository

import (
	"context"
	"database/sql"

	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
)

type reviewRepository struct {
	client db.Client
	crud   abstract.Crud
}

type ReviewRepository interface {
	All(context.Context) (*sql.Rows, error)
	Find(context.Context, string) (*sql.Row, error)
	Create(context.Context, string, string, string, string, string) error
	Update(context.Context, string, string, string, string, string, string) error
	Delete(context.Context, string) error
}

func NewReviewRepository(c db.Client, ac abstract.Crud) ReviewRepository {
	return &reviewRepository{c, ac}
}

// 全件取得
func (r *reviewRepository) All(c context.Context) (*sql.Rows, error) {
	query := "SELECT * FROM review"
	return r.crud.Read(c, query)
}

// 1件取得
func (r *reviewRepository) Find(c context.Context, id string) (*sql.Row, error) {
	query := "SELECT * FROM review WHERE id =" + id
	return r.crud.ReadByID(c, query)
}

// 作成
func (r *reviewRepository) Create(c context.Context, userID string, taskID string, staffingRating string, manualRating string, comment string) error {
	query := `
		INSERT INTO	reviews (user_id, task_id, staffing_rating, manual_rating, comment)
			VALUES ('` + userID + "', '" + taskID + "', " + staffingRating + ", " + manualRating + ", " + comment + "')"
	return r.crud.UpdateDB(c, query)
}

// 編集
func (r *reviewRepository) Update(c context.Context, ID string, userID string, taskID string, staffing_rating string, manual_rating string, comment string) error {
	query := `
		UPDATE
			reviews
		SET
			user_id = '` + userID +
		"', task_id = " + taskID +
		", staffing_rating = " + staffing_rating +
		", manual_rating = " + manual_rating +
		", comment = '" + comment + "' WHERE id = " + ID
	return r.crud.UpdateDB(c, query)
}

// 削除
func (r *reviewRepository) Delete(c context.Context, id string) error {
	query := "DELETE FROM reviews WHERE id = " + id
	return r.crud.UpdateDB(c, query)
}
