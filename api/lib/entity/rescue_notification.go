package entity

// レスキューの通知種別。送信者に伝える内容を表し、Slackの見出しを決める
const (
	RescueNotificationInProgress     = "inProgress"     // 本部が確認した(未対応→対応中)
	RescueNotificationDone           = "done"           // 対応が完了した(→対応済み)
	RescueNotificationResponse       = "response"       // 対応状況はそのままで、返答が初めて書かれた
	RescueNotificationResponseEdited = "responseEdited" // 対応状況はそのままで、書いてあった返答が書き換えられた
)
