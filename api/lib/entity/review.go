package entity

import (
	"time"
)

type Review struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	TaskID         int       `json:"task_id"`
	StaffingRating int       `json:"staffing_rating"`
	ManualRating   int       `json:"manual_rating"`
	Comment        string    `json:"comment"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
