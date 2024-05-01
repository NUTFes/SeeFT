package entity

import (
	"time"
)

type Place struct {
	ID			int 		`json:"id"`
	Place 		string		`json:"place"`
	Remark		string		`json:"remark"`
	CreatedAt	time.Time	`json:"createdAt"`
	UpdatedAt	time.Time	`json:"updatedAt"`
}