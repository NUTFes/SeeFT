package entity

import (
	"time"
)

type Department struct {
	ID        	int       `json:"id"`
	Department 	string    `json:"department"`
	CreatedAt 	time.Time `json:"createdAt"`
	UpdatedAt 	time.Time `json:"updatedAt"`
}