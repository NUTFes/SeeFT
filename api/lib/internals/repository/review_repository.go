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
	AllWithDetails(ctx context.Context) (*sql.Rows, error)
	FindWithDetails(ctx context.Context, id string) (*sql.Row, error)
	Create(context.Context, string, string, string, string, string) error
	Update(context.Context, string, string, string, string, string, string) error
	Delete(context.Context, string) error
	FindNewRecord(context.Context) (*sql.Row, error)
}

func NewReviewRepository(c db.Client, ac abstract.Crud) ReviewRepository {
	return &reviewRepository{c, ac}
}

// 全件取得
func (r *reviewRepository) All(c context.Context) (*sql.Rows, error) {
	query := "SELECT * FROM reviews"
	return r.crud.Read(c, query)
}

// 1件取得 - プリペアドステートメント使用
func (r *reviewRepository) Find(c context.Context, id string) (*sql.Row, error) {
	query := "SELECT * FROM reviews WHERE id = $1"
	return r.client.DB().QueryRowContext(c, query, id), nil
}

// JOIN用定義
const baseJoin = `
SELECT
  r.id,
  COALESCE(u.name, 'Unknown') AS user_name,
  COALESCE(b.bureau, 'Unknown') AS user_bureau,
  COALESCE(g.grade, 'Unknown') AS user_grade,
  COALESCE(u.student_number, '') AS user_studentnumber,
  COALESCE(t.task, 'Unknown') AS task_name,
  r.staffing_rating,
  r.manual_rating,
  r.comment,
  r.created_at,
  r.updated_at
FROM reviews r
LEFT JOIN users   u ON r.user_id  = u.id
LEFT JOIN tasks   t ON r.task_id  = t.id
LEFT JOIN bureaus b ON u.bureau_id = b.id
LEFT JOIN grades  g ON u.grade_id  = g.id
`

// 全件指定時GAS用変換
func (r *reviewRepository) AllWithDetails(c context.Context) (*sql.Rows, error) {
	query := baseJoin + " ORDER BY r.id"
	return r.client.DB().QueryContext(c, query)
}

// 指定時GAS用変換
func (r *reviewRepository) FindWithDetails(c context.Context, id string) (*sql.Row, error) {
	query := baseJoin + " WHERE r.id = $1"
	return r.client.DB().QueryRowContext(c, query, id), nil
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

// FindNewRecord 最新のレビューを取得する
func (r *reviewRepository) FindNewRecord(c context.Context) (*sql.Row, error) {
	query := `
		SELECT
			*
		FROM
			reviews
		ORDER BY
			id DESC
		LIMIT 1`
	return r.crud.ReadByID(c, query)
}
