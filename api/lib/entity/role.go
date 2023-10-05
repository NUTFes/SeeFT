package entity

import (
	"time"
)

type Role struct {
	ID        int        `json:"id"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}