package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/NUTFes/SeeFT/api/lib/externals/db"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository/abstract"
	"github.com/pkg/errors"
)

type shiftRepository struct {
	client db.Client
	crud   abstract.Crud
}

type ShiftRepository interface {
	All(context.Context) (*sql.Rows, error)
	Find(context.Context, string) (*sql.Row, error)
	User(context.Context, string) (*sql.Rows, error)
	Users(context.Context, string, string, string, string, string) (*sql.Rows, error)
	UserAndDateAndWeather(context.Context, string, string, string) (*sql.Rows, error)
	DateAndWeather(context.Context, string, string) (*sql.Rows, error)
	DateAndWeatherAndTime(context.Context, string, string, string, string) (*sql.Rows, error)
	Create(context.Context, string, string, string, string, string, string, string) error
	Update(context.Context, string, string, string, string, string, string, string, string) error
	Destroy(context.Context, string) error
	FindLatestRecord(context.Context) (*sql.Row, error)
	MaxID(context.Context) (*sql.Row, error)
	FindByUnique(context.Context, string, string, string, string, string) (*sql.Row, error)
	CreateAndReturnID(context.Context, string, string, string, string, string, string, string) (int, error)
}

func NewShiftRepository(c db.Client, ac abstract.Crud) ShiftRepository {
	return &shiftRepository{c, ac}
}

// 全件取得
func (b *shiftRepository) All(c context.Context) (*sql.Rows, error) {
	query := "SELECT * FROM shifts"
	return b.crud.Read(c, query)
}

// 1件取得
func (b *shiftRepository) Find(c context.Context, id string) (*sql.Row, error) {
	query := "SELECT * FROM shifts WHERE id = " + id
	return b.crud.ReadByID(c, query)
}

// 特定のユーザ取得
func (b *shiftRepository) User(c context.Context, id string) (*sql.Rows, error) {
	query := "SELECT * FROM shifts WHERE user_id = " + id
	rows, err := b.client.DB().QueryContext(c, query)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot connect SQL")
	}
	fmt.Printf("\x1b[36m%s\n", query)
	return rows, nil
}

// 特定のタスクのユーザ取得（JOINでユーザー情報も一括取得）
func (b *shiftRepository) Users(c context.Context, task string, year string, date string, time string, weather string) (*sql.Rows, error) {
	query := "SELECT u.id, u.name, u.mail, u.grade_id, u.department_id, u.bureau_id, u.role_id, u.student_number, u.tel, u.created_at, u.updated_at FROM shifts s JOIN users u ON s.user_id = u.id WHERE s.task_id = $1 AND s.year_id = $2 AND s.date_id = $3 AND s.time_id = $4 AND s.weather_id = $5"
	rows, err := b.client.DB().QueryContext(c, query, task, year, date, time, weather)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot connect SQL")
	}
	fmt.Printf("\x1b[36m%s\n", query)
	return rows, nil
}

// 特定のユーザと日時取得
func (b *shiftRepository) UserAndDateAndWeather(c context.Context, id string, date string, weather string) (*sql.Rows, error) {
	query := "SELECT * FROM shifts WHERE user_id =" + id + " AND date_id =" + date + " AND weather_id =" + weather
	rows, err := b.client.DB().QueryContext(c, query)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot connect SQL")
	}
	fmt.Printf("\x1b[36m%s\n", query)
	return rows, nil
}

// 特定と日時取得
func (b *shiftRepository) DateAndWeather(c context.Context, date string, weather string) (*sql.Rows, error) {
	query := "SELECT * FROM shifts WHERE date_id =" + date + " AND weather_id =" + weather
	rows, err := b.client.DB().QueryContext(c, query)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot connect SQL")
	}
	fmt.Printf("\x1b[36m%s\n", query)
	return rows, nil
}

// 特定と日時取得
func (b *shiftRepository) DateAndWeatherAndTime(c context.Context, date string, weather string, lower string, upper string) (*sql.Rows, error) {
	query := "SELECT * FROM shifts WHERE date_id = " + date + " AND weather_id = " + weather + " AND time_id BETWEEN " + lower + " AND " + upper
	rows, err := b.client.DB().QueryContext(c, query)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot connect SQL")
	}
	fmt.Printf("\x1b[36m%s\n", query)
	return rows, nil
}

// 作成
func (b *shiftRepository) Create(c context.Context, taskID string, userID string, yearID string, dateID string, timeID string, weatherID string, isAttendance string) error {
	query := "INSERT INTO shifts (task_id, user_id, year_id, date_id, time_id, weather_id, is_attendance) VALUES (" + taskID + ", " + userID + ", " + yearID + ", " + dateID + ", " + timeID + ", " + weatherID + ", " + isAttendance + ")"
	return b.crud.UpdateDB(c, query)
}

// 作成してIDを返す
func (b *shiftRepository) CreateAndReturnID(c context.Context, taskID string, userID string, yearID string, dateID string, timeID string, weatherID string, isAttendance string) (int, error) {
	query := "INSERT INTO shifts (task_id, user_id, year_id, date_id, time_id, weather_id, is_attendance) VALUES (" + taskID + ", " + userID + ", " + yearID + ", " + dateID + ", " + timeID + ", " + weatherID + ", " + isAttendance + ") RETURNING id"
	var id int
	err := b.client.DB().QueryRowContext(c, query).Scan(&id)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to create shift and return id")
	}
	fmt.Printf("\x1b[36m%s\n", query)
	return id, nil
}

// 編集
func (b *shiftRepository) Update(c context.Context, id string, taskID string, userID string, yearID string, dateID string, timeID string, weatherID string, isAttendance string) error {
	query := `
	UPDATE 
		shifts 
	SET 
		task_id = ` + taskID +
		", user_id = " + userID +
		", year_id = " + yearID +
		", date_id = " + dateID +
		", time_id = " + timeID +
		", weather_id = " + weatherID +
		", is_attendance = " + isAttendance +
		" WHERE id = " + id
	return b.crud.UpdateDB(c, query)
}

// 削除
func (b *shiftRepository) Destroy(c context.Context, id string) error {
	query := "DELETE FROM shifts WHERE id =" + id
	return b.crud.UpdateDB(c, query)
}

// 最新のbureauを取得する
func (b *shiftRepository) FindLatestRecord(c context.Context) (*sql.Row, error) {
	query := `
		SELECT
			*
		FROM
			shifts
		ORDER BY
			id
		DESC LIMIT 1
	`
	return b.crud.ReadByID(c, query)
}

// 最大のIDを取得する
func (b *shiftRepository) MaxID(c context.Context) (*sql.Row, error) {
	query := `
		SELECT
			MAX(id)
		FROM
			shifts
	`
	return b.crud.ReadByID(c, query)
}

func (b *shiftRepository) FindByUnique(c context.Context, taskID, userID, dateID, timeID, weatherID string) (*sql.Row, error) {
	query := "SELECT * FROM shifts WHERE user_id = " + userID + " AND date_id = " + dateID + " AND time_id = " + timeID + " AND weather_id = " + weatherID
	return b.client.DB().QueryRowContext(c, query), nil
}
