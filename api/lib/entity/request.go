package entity

// QuestionRescue作成リクエスト
type QuestionRescueCreateRequest struct {
	UserID  int    `json:"user_id"`
	Question string `json:"question"`
	Status   string `json:"status"`
}

// QuestionRescue更新リクエスト
type QuestionRescueUpdateRequest struct {
	Status   string `json:"status"`
	Response string `json:"response"`
}

// QuestionRescue削除リクエスト
type QuestionRescueDeleteRequest struct {
	ID string `json:"id"`
}

// ShorthandedRescue作成リクエスト
type ShorthandedRescueCreateRequest struct {
	UserID       int    `json:"user_id"`
	TaskID       int    `json:"task_id"`
	MissingNumber int   `json:"missing_number"`
	Place        string `json:"place"`
	Status       string `json:"status"`
}

// ShorthandedRescue更新リクエスト
type ShorthandedRescueUpdateRequest struct {
	Status   string `json:"status"`
	Response string `json:"response"`
}

// ShorthandedRescue削除リクエスト
type ShorthandedRescueDeleteRequest struct {
	ID string `json:"id"`
}

// TroubleRescue作成リクエスト
type TroubleRescueCreateRequest struct {
	UserID int    `json:"user_id"`
	TaskID int    `json:"task_id"`
	Place  string `json:"place"`
	Detail string `json:"detail"`
	Status string `json:"status"`
}

// TroubleRescue更新リクエスト
type TroubleRescueUpdateRequest struct {
	Status   string `json:"status"`
	Response string `json:"response"`
}

// TroubleRescue削除リクエスト
type TroubleRescueDeleteRequest struct {
	ID string `json:"id"`
}
