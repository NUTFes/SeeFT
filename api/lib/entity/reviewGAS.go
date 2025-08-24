// 変換用レポジトリ
package entity

type ReviewGAS struct {
	ID             int    `json:"id"`
	UserName       string `json:"user_name"`
	UserBureau     string `json:"user_bureau"`
	UserGrade      string `json:"user_grade"`
	UserStudentNo  int    `json:"user_studentnumber"`
	TaskName       string `json:"task_name"`
	StaffingRating int    `json:"staffing_rating"`
	ManualRating   int    `json:"manual_rating"`
	Comment        string `json:"comment"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}
