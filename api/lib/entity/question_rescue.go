package entity

import (
	"database/sql"
	"time"
)

type QuestionRescue struct {
	ID        int            `json:"id" gorm:"primaryKey"`
	UserID    int            `json:"user_id" gorm:"not null"`
	Question  string         `json:"question" gorm:"not null"`
	Status    string         `json:"status" gorm:"default:todo"`
	Response  sql.NullString `json:"response"`
	Time      time.Time      `json:"time" gorm:"default:CURRENT_TIMESTAMP"`
	CreatedAt time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

type QuestionRescueForGet struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Question  string    `json:"question"`
	Status    string    `json:"status"`
	Response  string    `json:"response"`
	Time      time.Time `json:"time"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
