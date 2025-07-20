package entity

import (
	"time"
)

type Review struct {
	ID             int       `json:"id"`
	User_ID        int       `json:"user_id"`
	task_id        int       `json:"task_ID"`
	staffig_rating int       `json:"staffing_rating"`
	manual_rating  int       `json:"manual_rating"`
	comment        string    `json:"comment"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
