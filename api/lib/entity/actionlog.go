package entity

import (
	"encoding/json"
	"time"
)

type ActionLog struct {
	ID          int             `json:"id"`
	ShiftID     int             `json:"shift_id"`
	UserID      int             `json:"user_id"`
	DateID      int             `json:"date_id"`
	ActionType  string          `json:"action_type"` // CREATE, UPDATE, DELETE
	DiffPayload json.RawMessage `json:"diff_payload"`
	IsSent      bool            `json:"is_sent"`
	CreatedAt   time.Time       `json:"created_at"`
}
