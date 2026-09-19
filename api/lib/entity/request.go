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
	RescueNotifyOption
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
	RescueNotifyOption
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
	RescueNotifyOption
}

// TroubleRescue削除リクエスト
type TroubleRescueDeleteRequest struct {
	ID string `json:"id"`
}

// レスキュー更新リクエストの共通項目
type RescueNotifyOption struct {
	// 送信者にSlack DMで知らせるか。省略時は知らせる。
	// GASが押し直しの重複を「対応済み」にしてまとめるときだけ false を送る。
	// 送らないと、本部が何もしていないのに「対応が完了しました」が届くため
	Notify *bool `json:"notify"`
}

// ShouldNotify 省略(nil)は知らせる扱いにする。本部の書き込み(onChange)は notify を送らないため
func (o RescueNotifyOption) ShouldNotify() bool {
	return o.Notify == nil || *o.Notify
}
