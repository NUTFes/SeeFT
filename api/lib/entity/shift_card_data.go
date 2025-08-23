package entity


type ShiftCardData struct {
    // Shiftテーブルのカラム
    ShiftID      int       `gorm:"column:shift_id"`
    UserID       int       `gorm:"column:user_id"`
    TaskID       int       `gorm:"column:task_id"`
    YearID       int       `gorm:"column:year_id"`
    DateID       int       `gorm:"column:date_id"`
    TimeID       int       `gorm:"column:time_id"`
    WeatherID    int       `gorm:"column:weather_id"`
    IsAttendance bool      `gorm:"column:is_attendance"`

    // JOINしたデータ
    TaskName      string    `gorm:"column:task_name"`
    TaskColor     string    `gorm:"column:task_color"`
    TaskURL       string    `gorm:"column:task_url"`
    TaskRemark    string    `gorm:"column:task_remark"`
    MaxMember     int       `gorm:"column:max_member"`
    TaskBureauID  int       `gorm:"column:task_bureau_id"`
    PlaceID       int       `gorm:"column:place_id"`
    PlaceName     string    `gorm:"column:place_name"`
    TimeValue     string    `gorm:"column:time_value"`
    UserName      string    `gorm:"column:user_name"`
    UserBureauID  int       `gorm:"column:user_bureau_id"`
    UserGradeID   int       `gorm:"column:user_grade_id"`
    YearValue     string    `gorm:"column:year_value"`
    DateValue     string    `gorm:"column:date_value"`
    WeatherValue  string    `gorm:"column:weather_value"`
}
