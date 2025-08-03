package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/NUTFes/SeeFT/api/lib/usecase"
	"github.com/labstack/echo/v4"
)

// 共通のリクエスト構造体
type RescueRequest struct {
	Type    string      `json:"type"`
	UserID  int         `json:"user_id"`
	Content interface{} `json:"content"`
}

// 各タイプ別のContent構造体
type TroubleContent struct {
	TaskID int    `json:"task_id"`
	Place  string `json:"place"`
	Detail string `json:"detail"`
}

type QuestionContent struct {
	Question string `json:"question"`
}

type ShorthandedContent struct {
	TaskID        int    `json:"task_id"`
	MissingNumber int    `json:"missing_number"`
	Place         string `json:"place"`
}

// Content共通情報構造体
type RescueCommonInfo struct {
	SenderName    string
	Grade         string
	Bureau        string
	StudentNumber int
	TaskName      string
	AnsweredAt    string
	EditedAt      string
	PhoneNumber   string
}

type rescueUnifiedController struct {
	questionRescueUseCase    usecase.QuestionRescueUseCase
	shorthandedRescueUseCase usecase.ShorthandedRescueUseCase
	troubleRescueUseCase     usecase.TroubleRescueUseCase
	rescueUnifiedUseCase     usecase.RescueUnifiedUseCase
	useruse                  usecase.UserUseCase
	taskuse                  usecase.TaskUseCase
	gradeuse                 usecase.GradeUseCase
	bureauuse                usecase.BureauUseCase
}

type RescueUnifiedController interface {
	CreateRescue(echo.Context) error
	GetAllRescues(echo.Context) error
	GetRescuesByUserID(echo.Context) error
}

func NewRescueUnifiedController(
	qu usecase.QuestionRescueUseCase,
	su usecase.ShorthandedRescueUseCase,
	tu usecase.TroubleRescueUseCase,
	ru usecase.RescueUnifiedUseCase,
	useruse usecase.UserUseCase,
	taskuse usecase.TaskUseCase,
	gradeuse usecase.GradeUseCase,
	bureauuse usecase.BureauUseCase,
) RescueUnifiedController {
	return &rescueUnifiedController{qu, su, tu, ru, useruse, taskuse, gradeuse, bureauuse}
}

// Contentから共通情報を抽出する関数
func (r *rescueUnifiedController) extractRescueCommonInfo(c echo.Context, reqType string, content interface{}, userID int) RescueCommonInfo {
	getUserInfo := func(userID int) (string, string, string, int, string) {
		userIDStr := strconv.Itoa(userID)
		user, err := r.useruse.GetUserByID(c.Request().Context(), userIDStr)
		if err != nil {
			return "不明なユーザー", "不明", "不明", 0, "不明"
		}
		grade, err := r.gradeuse.GetGradeByID(c.Request().Context(), strconv.Itoa(user.GradeID))
		if err != nil {
			return user.Name, "不明", "不明", user.StudentNumber, user.Tel
		}
		bureau, err := r.bureauuse.GetBureauByID(c.Request().Context(), strconv.Itoa(user.BureauID))
		if err != nil {
			return user.Name, grade.Grade, "不明", user.StudentNumber, user.Tel
		}

		return user.Name, grade.Grade, bureau.Bureau, user.StudentNumber, user.Tel
	}

	getTaskName := func(taskID int) string {
		taskIDStr := strconv.Itoa(taskID)
		task, err := r.taskuse.GetTaskByID(c.Request().Context(), taskIDStr)
		if err != nil {
			return ""
		}
		return task.Task
	}

	loc, err := time.LoadLocation("Asia/Tokyo")
	var now string
	if err == nil {
		now = time.Now().In(loc).Format("2006-01-02 15:04:05")
	} else {
		now = time.Now().Format("2006-01-02 15:04:05") // ロケーション取得失敗時はUTC
	}
	senderName, grade, bureau, studentNumber, phoneNumber := getUserInfo(userID)
	var taskName string

	switch reqType {
	case "trouble":
		c := content.(map[string]interface{})
		if tidFloat, ok := c["task_id"].(float64); ok {
			tidInt := int(tidFloat)
			taskName = getTaskName(tidInt)
		}
		
		return RescueCommonInfo{
			SenderName:    senderName,
			Grade:         grade,
			Bureau:        bureau,
			StudentNumber: studentNumber,
			TaskName:      taskName,
			AnsweredAt:    now,
			EditedAt:      now,
			PhoneNumber:   phoneNumber,
		}
	case "question":
		return RescueCommonInfo{
			SenderName:    senderName,
			Grade:         grade,
			Bureau:        bureau,
			StudentNumber: studentNumber,
			TaskName:      "",
			AnsweredAt:    now,
			EditedAt:      now,
			PhoneNumber:   phoneNumber,
		}
	case "shorthanded":
		c := content.(map[string]interface{})
		if tidFloat, ok := c["task_id"].(float64); ok {
			tidInt := int(tidFloat)
			taskName = getTaskName(tidInt)
		}
		return RescueCommonInfo{
			SenderName:    senderName,
			Grade:         grade,
			Bureau:        bureau,
			StudentNumber: studentNumber,
			TaskName:      taskName,
			AnsweredAt:    now,
			EditedAt:      now,
			PhoneNumber:   phoneNumber,
		}
	default:
		return RescueCommonInfo{}
	}
}

// 統一されたレスキューリクエスト処理
func (r *rescueUnifiedController) CreateRescue(c echo.Context) error {
	var req RescueRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	// 基本的な入力バリデーション
	if req.UserID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}
	userIDStr := strconv.Itoa(req.UserID)

	// ここでDB保存用データとスプシ保存用データを作成
	commonInfo := r.extractRescueCommonInfo(c, req.Type, req.Content, req.UserID)
	var rescueData map[string]interface{}
	switch req.Type {
	case "trouble":
		c := req.Content.(map[string]interface{})
		rescueData = map[string]interface{}{
			"rescue_type": req.Type,
			"sender_name": commonInfo.SenderName,
			"student_number": commonInfo.StudentNumber,
			"phone_number": commonInfo.PhoneNumber,
			"grade": commonInfo.Grade,
			"bureau": commonInfo.Bureau,
			"answered_at": commonInfo.AnsweredAt,
			"place": c["place"],
			"task_name": commonInfo.TaskName,
			"detail": c["detail"],
		}
	case "question":
		c := req.Content.(map[string]interface{})
		rescueData = map[string]interface{}{
			"rescue_type": req.Type,
			"sender_name": commonInfo.SenderName,
			"student_number": commonInfo.StudentNumber,
			"phone_number": commonInfo.PhoneNumber,
			"grade": commonInfo.Grade,
			"bureau": commonInfo.Bureau,
			"answered_at": commonInfo.AnsweredAt,
			"question": c["question"],
		}
	case "shorthanded":
		c := req.Content.(map[string]interface{})
		rescueData = map[string]interface{}{
			"rescue_type": req.Type,
			"sender_name": commonInfo.SenderName,
			"student_number": commonInfo.StudentNumber,
			"phone_number": commonInfo.PhoneNumber,
			"grade": commonInfo.Grade,
			"bureau": commonInfo.Bureau,
			"answered_at": commonInfo.AnsweredAt,
			"place": c["place"],
			"task_name": commonInfo.TaskName,
			"missing_number": c["missing_number"],
		}
	default:
		rescueData = map[string]interface{}{}
	}

	var wg sync.WaitGroup
	wg.Add(2)

	var dbErr, sheetErr error
	var createdRescue interface{}

	// DB保存
	go func() {
		defer wg.Done()
		switch req.Type {
		case "trouble":
			dbErr = r.handleTroubleRescue(c, req, userIDStr)
		case "question":
			dbErr = r.handleQuestionRescue(c, req, userIDStr)
		case "shorthanded":
			dbErr = r.handleShorthandedRescue(c, req, userIDStr)
		default:
			dbErr = c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid rescue type"})
		}
	}()

	// スプシ保存
	go func() {
		defer wg.Done()
		sheetErr = r.rescueUnifiedUseCase.SaveRescueToSpreadsheet(rescueData)
	}()

	wg.Wait()

	if dbErr != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": dbErr.Error()})
	}
	if sheetErr != nil {
		// スプシ保存失敗は警告として返す
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Rescue saved to DB, but failed to save to spreadsheet",
			"data":    createdRescue,
			"sheet_error": sheetErr.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Rescue saved to DB and spreadsheet successfully",
		"data":    createdRescue,
	})
}

func (r *rescueUnifiedController) handleTroubleRescue(c echo.Context, req RescueRequest, userIDStr string) error {
	// content をTroubleContentに変換
	contentBytes, err := json.Marshal(req.Content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid trouble content"})
	}

	var content TroubleContent
	if err := json.Unmarshal(contentBytes, &content); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid trouble content format"})
	}

	// バリデーション
	if content.TaskID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid task ID"})
	}
	if content.Detail == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Detail is required"})
	}

	taskIDStr := strconv.Itoa(content.TaskID)
	createdTroubleRescue, err := r.troubleRescueUseCase.CreateTroubleRescue(
		c.Request().Context(),
		userIDStr,
		taskIDStr,
		content.Place,
		content.Detail,
		"todo",
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Trouble rescue created successfully",
		"data":    createdTroubleRescue,
	})
}

func (r *rescueUnifiedController) handleQuestionRescue(c echo.Context, req RescueRequest, userIDStr string) error {
	// content をQuestionContentに変換
	contentBytes, err := json.Marshal(req.Content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid question content"})
	}

	var content QuestionContent
	if err := json.Unmarshal(contentBytes, &content); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid question content format"})
	}

	// バリデーション
	if content.Question == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Question is required"})
	}

	createdQuestionRescue, err := r.questionRescueUseCase.CreateQuestionRescue(
		c.Request().Context(),
		userIDStr,
		content.Question,
		"todo",
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Question rescue created successfully",
		"data":    createdQuestionRescue,
	})
}

func (r *rescueUnifiedController) handleShorthandedRescue(c echo.Context, req RescueRequest, userIDStr string) error {
	// content をShorthandedContentに変換
	contentBytes, err := json.Marshal(req.Content)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid shorthanded content"})
	}

	var content ShorthandedContent
	if err := json.Unmarshal(contentBytes, &content); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid shorthanded content format"})
	}

	// バリデーション
	if content.TaskID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid task ID"})
	}
	if content.MissingNumber <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing number must be greater than 0"})
	}

	taskIDStr := strconv.Itoa(content.TaskID)
	missingNumberStr := strconv.Itoa(content.MissingNumber)

	createdShorthandedRescue, err := r.shorthandedRescueUseCase.CreateShorthandedRescue(
		c.Request().Context(),
		userIDStr,
		taskIDStr,
		missingNumberStr,
		content.Place,
		"todo",
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Shorthanded rescue created successfully",
		"data":    createdShorthandedRescue,
	})
}

// 全件取得
func (r *rescueUnifiedController) GetAllRescues(c echo.Context) error {
	rescues, err := r.rescueUnifiedUseCase.GetAllRescues(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, rescues)
}

// ユーザーID別取得
func (r *rescueUnifiedController) GetRescuesByUserID(c echo.Context) error {
	userID := c.Param("user_id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "User ID is required"})
	}

	// 数値チェック
	if _, err := strconv.Atoi(userID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	rescues, err := r.rescueUnifiedUseCase.GetRescuesByUserID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, rescues)
}
