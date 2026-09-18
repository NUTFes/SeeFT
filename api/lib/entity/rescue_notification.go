package entity

import "time"

// レスキューの種類。rescue_notifications.rescue_type と RescueResponse.Type に入る値
const (
	RescueTypeTrouble     = "trouble"
	RescueTypeQuestion    = "question"
	RescueTypeShorthanded = "shorthanded"
)

// レスキューの通知種別。送信者に伝える内容を表し、Slackの見出しを決める
const (
	RescueNotificationInProgress = "inProgress" // 本部が確認した(未対応→対応中)
	RescueNotificationDone       = "done"       // 対応が完了した(→対応済み)
	RescueNotificationResponse   = "response"   // 対応状況はそのままで返答だけ届いた
)

type RescueNotification struct {
	ID         int       `json:"id"`
	RescueType string    `json:"rescue_type"`
	RescueID   int       `json:"rescue_id"`
	UserID     int       `json:"user_id"`
	Kind       string    `json:"kind"`
	Status     string    `json:"status"`
	Response   string    `json:"response"`
	IsSent     bool      `json:"is_sent"`
	CreatedAt  time.Time `json:"created_at"`
}
