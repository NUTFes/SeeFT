package repository

import (
    "context"

    "github.com/NUTFes/SeeFT/api/lib/entity"
    "github.com/NUTFes/SeeFT/api/lib/externals/db"
    "gorm.io/gorm"
)

type shiftCardRepository struct {
    db *gorm.DB
}

type ShiftCardRepository interface {
    GetOptimizedShiftData(context.Context, string, string, string) ([]entity.ShiftCardData, error)
}

func NewShiftCardRepository(client db.Client) ShiftCardRepository {
    return &shiftCardRepository{
        db: client.GormDB(), // GORMのDBインスタンスを使用
    }
}

func (r *shiftCardRepository) GetOptimizedShiftData(ctx context.Context, userID, dateID, weatherID string) ([]entity.ShiftCardData, error) {
    var results []entity.ShiftCardData

    // GORMを使用したJOINクエリ
    err := r.db.WithContext(ctx).
        Table("shifts").
        Select(`
            shifts.id as shift_id,
            shifts.user_id,
            shifts.task_id,
            shifts.year_id,
            shifts.date_id,
            shifts.time_id,
            shifts.weather_id,
            shifts.is_attendance,
            tasks.task as task_name,
            tasks.color as task_color,
            tasks.url as task_url,
            COALESCE(tasks.manual_url, '') AS manual_url,
            tasks.remark as task_remark,
            tasks.max_member,
            tasks.bureau_id as task_bureau_id,
            places.id as place_id,
            places.place as place_name,
            times.time as time_value,
            users.name as user_name,
            users.bureau_id as user_bureau_id,
            users.grade_id as user_grade_id,
            years.year as year_value,
            dates.date as date_value,
            weathers.weather as weather_value
        `).
        Joins("LEFT JOIN tasks ON shifts.task_id = tasks.id").
        Joins("LEFT JOIN places ON tasks.place_id = places.id").
        Joins("LEFT JOIN times ON shifts.time_id = times.id").
        Joins("LEFT JOIN users ON shifts.user_id = users.id").
        Joins("LEFT JOIN years ON shifts.year_id = years.id").
        Joins("LEFT JOIN dates ON shifts.date_id = dates.id").
        Joins("LEFT JOIN weathers ON shifts.weather_id = weathers.id").
        Where("shifts.user_id = ?", userID).
        Where("shifts.date_id = ?", dateID).
        Where("shifts.weather_id = ?", weatherID).
        Where("tasks.task != ?", "").
        Where("tasks.task != ?", "NG").
        Order("shifts.time_id ASC").
        Scan(&results).Error

    if err != nil {
        return nil, err
    }

    return results, nil
}

