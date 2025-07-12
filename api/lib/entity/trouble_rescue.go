package entity

import (
	"database/sql"
	"time"
)

type TroubleRescue struct {
	ID        int            `json:"id" gorm:"primaryKey"`
	UserID    int            `json:"user_id" gorm:"not null"`
	TaskID    int            `json:"task_id" gorm:"not null"`
	Place     sql.NullString `json:"place"`
	Detail    string         `json:"detail" gorm:"not null"`
	Status    string         `json:"status" gorm:"default:todo"`
	Response  sql.NullString `json:"response"`
	Time      time.Time      `json:"time" gorm:"default:CURRENT_TIMESTAMP"`
	CreatedAt time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

type TroubleRescueForGet struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	TaskID    int       `json:"task_id"`
	Place     string    `json:"place"`
	Detail    string    `json:"detail"`
	Status    string    `json:"status"`
	Response  string    `json:"response"`
	Time      time.Time `json:"time"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
