package entity

// 統一レスキューレスポンス用の構造体
type RescueResponse struct {
	Type     string      `json:"type"`
	ID       int         `json:"id"`
	UserName string      `json:"user_name"`
	Time     string      `json:"time"`
	Content  interface{} `json:"content"`
	Status   string      `json:"status"`
	Response string      `json:"response"`
}

// 各タイプ別のContent構造体（レスポンス用）
type TroubleResponseContent struct {
	Task   string `json:"task"`
	Place  string `json:"place"`
	Detail string `json:"detail"`
}

type QuestionResponseContent struct {
	Question string `json:"question"`
}

type ShorthandedResponseContent struct {
	Task          string `json:"task"`
	MissingNumber int    `json:"missing_number"`
	Place         string `json:"place"`
}

// レスキューレスポンス作成用のヘルパー関数
func NewTroubleRescueResponse(tr *TroubleRescueForGet, userName string, taskName string) *RescueResponse {
	return &RescueResponse{
		Type:     "trouble",
		ID:       tr.ID,
		UserName: userName,
		Time:     tr.Time.Format("2006/01/02 15:04:05"),
		Content: TroubleResponseContent{
			Task:   taskName,
			Place:  tr.Place,
			Detail: tr.Detail,
		},
		Status:   tr.Status,
		Response: tr.Response,
	}
}

func NewQuestionRescueResponse(qr *QuestionRescueForGet, userName string) *RescueResponse {
	return &RescueResponse{
		Type:     "question",
		ID:       qr.ID,
		UserName: userName,
		Time:     qr.Time.Format("2006/01/02 15:04:05"),
		Content: QuestionResponseContent{
			Question: qr.Question,
		},
		Status:   qr.Status,
		Response: qr.Response,
	}
}

func NewShorthandedRescueResponse(sr *ShorthandedRescueForGet, userName string, taskName string) *RescueResponse {
	return &RescueResponse{
		Type:     "shorthanded",
		ID:       sr.ID,
		UserName: userName,
		Time:     sr.Time.Format("2006/01/02 15:04:05"),
		Content: ShorthandedResponseContent{
			Task:          taskName,
			MissingNumber: sr.MissingNumber,
			Place:         sr.Place,
		},
		Status:   sr.Status,
		Response: sr.Response,
	}
}
