package entity

import "time"

// レスキューの時刻はJSTで返す。スプレッドシートへ送る時刻（rescue_unified_controller）も
// Asia/Tokyoで整形しているので、表示に使う時刻をここで揃える。
// LoadLocationを使わないのは、ゾーン情報を読めない環境で黙ってUTCに落ち、
// 直したはずの9時間ズレが再発するため。日本標準時は夏時間を持たないので固定オフセットで表せる。
var rescueTimeLocation = time.FixedZone("JST", 9*60*60)

// DBのtimestamptzを表示用の文字列に整形する
func formatRescueTime(t time.Time) string {
	return t.In(rescueTimeLocation).Format("2006/01/02 15:04:05")
}

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
		Time:     formatRescueTime(tr.Time),
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
		Time:     formatRescueTime(qr.Time),
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
		Time:     formatRescueTime(sr.Time),
		Content: ShorthandedResponseContent{
			Task:          taskName,
			MissingNumber: sr.MissingNumber,
			Place:         sr.Place,
		},
		Status:   sr.Status,
		Response: sr.Response,
	}
}
