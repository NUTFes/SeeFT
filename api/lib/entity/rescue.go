package entity

import (
	"time"
)

// 基本のレスキュー構造体
type RescueResponse struct {
	Type     string      `json:"type"`
	ID       int         `json:"id"`
	UserName string      `json:"user_name"`
	Time     string      `json:"time"`
	Content  interface{} `json:"content"`
	Status   string      `json:"status"`
	Response string      `json:"response"`
}

// トラブルレスキュー
type TroubleRescue struct {
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

// 質問レスキュー
type QuestionRescue struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Question  string    `json:"question"`
	Status    string    `json:"status"`
	Response  string    `json:"response"`
	Time      time.Time `json:"time"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 人手不足レスキュー
type ShorthandedRescue struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	TaskID        int       `json:"task_id"`
	MissingNumber int       `json:"missing_number"`
	Place         string    `json:"place"`
	Status        string    `json:"status"`
	Response      string    `json:"response"`
	Time          time.Time `json:"time"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// レスキュー更新用のリクエスト構造体
type RescueUpdateRequest struct {
	Type     string `json:"type"`
	ID       int    `json:"id"`
	Status   string `json:"status"`
	Response string `json:"response"`
}

// 複数のレスキュー更新用
type RescueUpdatesRequest struct {
	Updates []RescueUpdateRequest `json:"updates"`
}

// レスキュー作成用のリクエスト構造体
type TroubleRescueRequest struct {
	UserID int    `json:"user_id"`
	TaskID int    `json:"task_id"`
	Place  string `json:"place"`
	Detail string `json:"detail"`
}

type QuestionRescueRequest struct {
	UserID   int    `json:"user_id"`
	Question string `json:"question"`
}

type ShorthandedRescueRequest struct {
	UserID        int    `json:"user_id"`
	TaskID        int    `json:"task_id"`
	MissingNumber int    `json:"missing_number"`
	Place         string `json:"place"`
}
