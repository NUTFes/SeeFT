package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/NUTFes/SeeFT/api/lib/entity"
	rep "github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/pkg/errors"
)



// structが持つべき関数を定義
type RescueUseCase interface {
	CreateRescue(context.Context, string) error

	// 作成系
	CreateTroubleRescue(context.Context, entity.TroubleRescueRequest) error
	CreateQuestionRescue(context.Context, entity.QuestionRescueRequest) error
	CreateShorthandedRescue(context.Context, entity.ShorthandedRescueRequest) error

	// 取得系
	GetRescuesByUserID(context.Context, int) ([]entity.RescueResponse, error)
	GetAllRescues(context.Context) ([]entity.RescueResponse, error)

	// 更新系（GAS用）
	UpdateRescues(context.Context, entity.RescueUpdatesRequest) error
}

// クラス
type rescueUseCase struct {
	rescueRep rep.RescueRepository
	userRep   rep.UserRepository
	taskRep   rep.TaskRepository
}

// インスタンス化
func NewRescueUseCase(rescueRep rep.RescueRepository, userRep rep.UserRepository, taskRep rep.TaskRepository) RescueUseCase {
	return &rescueUseCase{
		rescueRep: rescueRep,
		userRep:   userRep,
		taskRep:   taskRep,
	}
}

// structの中身を定義
func (u *rescueUseCase) CreateRescue(ctx context.Context, taskID string) error {
	url := "https://script.google.com/macros/s/AKfycbw7rQNDPdB5fjKORCLbM9NU8hbOgCqPoiojkTZ4S2wk4t20DI_tSvbHhjL80kDv6lZ2Og/exec"
	data := map[string]string{"task_id": taskID}
	
	jsonData, err := json.Marshal(data)

	if err != nil {
		return errors.Wrap(err, "failed to marshal data")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to send rescue request")
	}

	return nil
}

//==========
// 各種処理
//==========

// トラブルレスキュー作成
func (r *rescueUseCase) CreateTroubleRescue(c context.Context, req entity.TroubleRescueRequest) error {
	return r.rescueRep.CreateTroubleRescue(c, req.UserID, req.TaskID, req.Place, req.Detail)
}

// 質問レスキュー作成
func (r *rescueUseCase) CreateQuestionRescue(c context.Context, req entity.QuestionRescueRequest) error {
	return r.rescueRep.CreateQuestionRescue(c, req.UserID, req.Question)
}

// 人手不足レスキュー作成
func (r *rescueUseCase) CreateShorthandedRescue(c context.Context, req entity.ShorthandedRescueRequest) error {
	return r.rescueRep.CreateShorthandedRescue(c, req.UserID, req.TaskID, req.MissingNumber, req.Place)
}

// ユーザーのレスキュー一覧取得
func (r *rescueUseCase) GetRescuesByUserID(c context.Context, userID int) ([]entity.RescueResponse, error) {
	var results []entity.RescueResponse

	// トラブルレスキュー取得
	troubleRows, err := r.rescueRep.FindTroubleRescuesByUserID(c, userID)
	if err != nil {
		return nil, err
	}
	defer troubleRows.Close()

	for troubleRows.Next() {
		var trouble entity.TroubleRescue
		err := troubleRows.Scan(&trouble.ID, &trouble.UserID, &trouble.TaskID, &trouble.Place, &trouble.Detail, &trouble.Status, &trouble.Response, &trouble.Time, &trouble.CreatedAt, &trouble.UpdatedAt)
		if err != nil {
			continue
		}

		// ユーザー名とタスク名を取得
		userName, taskName, err := r.getUserAndTaskNames(c, trouble.UserID, trouble.TaskID)
		if err != nil {
			continue
		}

		content := map[string]interface{}{
			"task":   taskName,
			"place":  trouble.Place,
			"detail": trouble.Detail,
		}

		results = append(results, entity.RescueResponse{
			Type:     "trouble",
			ID:       trouble.ID,
			UserName: userName,
			Time:     trouble.Time.Format("2006/01/02 15:04:05"),
			Content:  content,
			Status:   trouble.Status,
			Response: trouble.Response,
		})
	}

	// 質問レスキュー取得
	questionRows, err := r.rescueRep.FindQuestionRescuesByUserID(c, userID)
	if err != nil {
		return nil, err
	}
	defer questionRows.Close()

	for questionRows.Next() {
		var question entity.QuestionRescue
		err := questionRows.Scan(&question.ID, &question.UserID, &question.Question, &question.Status, &question.Response, &question.Time, &question.CreatedAt, &question.UpdatedAt)
		if err != nil {
			continue
		}

		userName, _, err := r.getUserAndTaskNames(c, question.UserID, 0)
		if err != nil {
			continue
		}

		content := map[string]interface{}{
			"question": question.Question,
		}

		results = append(results, entity.RescueResponse{
			Type:     "question",
			ID:       question.ID,
			UserName: userName,
			Time:     question.Time.Format("2006/01/02 15:04:05"),
			Content:  content,
			Status:   question.Status,
			Response: question.Response,
		})
	}

	// 人手不足レスキュー取得
	shorthandedRows, err := r.rescueRep.FindShorthandedRescuesByUserID(c, userID)
	if err != nil {
		return nil, err
	}
	defer shorthandedRows.Close()

	for shorthandedRows.Next() {
		var shorthanded entity.ShorthandedRescue
		err := shorthandedRows.Scan(&shorthanded.ID, &shorthanded.UserID, &shorthanded.TaskID, &shorthanded.MissingNumber, &shorthanded.Place, &shorthanded.Status, &shorthanded.Response, &shorthanded.Time, &shorthanded.CreatedAt, &shorthanded.UpdatedAt)
		if err != nil {
			continue
		}

		userName, taskName, err := r.getUserAndTaskNames(c, shorthanded.UserID, shorthanded.TaskID)
		if err != nil {
			continue
		}

		content := map[string]interface{}{
			"task":           taskName,
			"missing_number": shorthanded.MissingNumber,
			"place":          shorthanded.Place,
		}

		results = append(results, entity.RescueResponse{
			Type:     "shorthanded",
			ID:       shorthanded.ID,
			UserName: userName,
			Time:     shorthanded.Time.Format("2006/01/02 15:04:05"),
			Content:  content,
			Status:   shorthanded.Status,
			Response: shorthanded.Response,
		})
	}

	return results, nil
}

// 全レスキュー取得
func (r *rescueUseCase) GetAllRescues(c context.Context) ([]entity.RescueResponse, error) {
	var results []entity.RescueResponse

	// 全トラブルレスキュー取得
	troubleRows, err := r.rescueRep.AllTroubleRescues(c)
	if err != nil {
		return nil, err
	}
	defer troubleRows.Close()

	for troubleRows.Next() {
		var trouble entity.TroubleRescue
		err := troubleRows.Scan(&trouble.ID, &trouble.UserID, &trouble.TaskID, &trouble.Place, &trouble.Detail, &trouble.Status, &trouble.Response, &trouble.Time, &trouble.CreatedAt, &trouble.UpdatedAt)
		if err != nil {
			continue
		}

		userName, taskName, err := r.getUserAndTaskNames(c, trouble.UserID, trouble.TaskID)
		if err != nil {
			continue
		}

		content := map[string]interface{}{
			"task":   taskName,
			"place":  trouble.Place,
			"detail": trouble.Detail,
		}

		results = append(results, entity.RescueResponse{
			Type:     "trouble",
			ID:       trouble.ID,
			UserName: userName,
			Time:     trouble.Time.Format("2006/01/02 15:04:05"),
			Content:  content,
			Status:   trouble.Status,
			Response: trouble.Response,
		})
	}

	// 全質問レスキュー取得
	questionRows, err := r.rescueRep.AllQuestionRescues(c)
	if err != nil {
		return nil, err
	}
	defer questionRows.Close()

	for questionRows.Next() {
		var question entity.QuestionRescue
		err := questionRows.Scan(&question.ID, &question.UserID, &question.Question, &question.Status, &question.Response, &question.Time, &question.CreatedAt, &question.UpdatedAt)
		if err != nil {
			continue
		}

		userName, _, err := r.getUserAndTaskNames(c, question.UserID, 0)
		if err != nil {
			continue
		}

		content := map[string]interface{}{
			"question": question.Question,
		}

		results = append(results, entity.RescueResponse{
			Type:     "question",
			ID:       question.ID,
			UserName: userName,
			Time:     question.Time.Format("2006/01/02 15:04:05"),
			Content:  content,
			Status:   question.Status,
			Response: question.Response,
		})
	}

	// 全人手不足レスキュー取得
	shorthandedRows, err := r.rescueRep.AllShorthandedRescues(c)
	if err != nil {
		return nil, err
	}
	defer shorthandedRows.Close()

	for shorthandedRows.Next() {
		var shorthanded entity.ShorthandedRescue
		err := shorthandedRows.Scan(&shorthanded.ID, &shorthanded.UserID, &shorthanded.TaskID, &shorthanded.MissingNumber, &shorthanded.Place, &shorthanded.Status, &shorthanded.Response, &shorthanded.Time, &shorthanded.CreatedAt, &shorthanded.UpdatedAt)
		if err != nil {
			continue
		}

		userName, taskName, err := r.getUserAndTaskNames(c, shorthanded.UserID, shorthanded.TaskID)
		if err != nil {
			continue
		}

		content := map[string]interface{}{
			"task":           taskName,
			"missing_number": shorthanded.MissingNumber,
			"place":          shorthanded.Place,
		}

		results = append(results, entity.RescueResponse{
			Type:     "shorthanded",
			ID:       shorthanded.ID,
			UserName: userName,
			Time:     shorthanded.Time.Format("2006/01/02 15:04:05"),
			Content:  content,
			Status:   shorthanded.Status,
			Response: shorthanded.Response,
		})
	}

	return results, nil
}

// レスキュー更新（GAS用）
func (r *rescueUseCase) UpdateRescues(c context.Context, req entity.RescueUpdatesRequest) error {
	for _, update := range req.Updates {
		switch update.Type {
		case "trouble":
			err := r.rescueRep.UpdateTroubleRescue(c, update.ID, update.Status, update.Response)
			if err != nil {
				return errors.Wrapf(err, "failed to update trouble rescue with ID %d", update.ID)
			}
		case "question":
			err := r.rescueRep.UpdateQuestionRescue(c, update.ID, update.Status, update.Response)
			if err != nil {
				return errors.Wrapf(err, "failed to update question rescue with ID %d", update.ID)
			}
		case "shorthanded":
			err := r.rescueRep.UpdateShorthandedRescue(c, update.ID, update.Status, update.Response)
			if err != nil {
				return errors.Wrapf(err, "failed to update shorthanded rescue with ID %d", update.ID)
			}
		default:
			return errors.New("unknown rescue type: " + update.Type)
		}
	}
	return nil
}

// ヘルパー関数：ユーザー名とタスク名を取得
func (r *rescueUseCase) getUserAndTaskNames(c context.Context, userID int, taskID int) (string, string, error) {
	var userName, taskName string

	// ユーザー名取得
	userRow, err := r.userRep.Find(c, strconv.Itoa(userID))
	if err != nil {
		return "", "", errors.Wrapf(err, "failed to find user with ID %d", userID)
	}
	var user entity.User
	err = userRow.Scan(&user.ID, &user.Name, &user.Mail, &user.GradeID, &user.DepartmentID, &user.BureauID, &user.RoleID, &user.StudentNumber, &user.Tel, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return "", "", errors.Wrapf(err, "failed to scan user data for ID %d", userID)
	}
	userName = user.Name

	// タスク名取得（taskIDが0の場合はスキップ）
	if taskID != 0 {
		taskRow, err := r.taskRep.Find(c, strconv.Itoa(taskID))
		if err != nil {
			return userName, "", errors.Wrapf(err, "failed to find task with ID %d", taskID)
		}
		var task entity.Task
		err = taskRow.Scan(&task.ID, &task.Task, &task.PlaceID, &task.Url, &task.BureauID, &task.MaxMember, &task.Color, &task.Remark, &task.YearID, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return userName, "", errors.Wrapf(err, "failed to scan task data for ID %d", taskID)
		}
		taskName = task.Task
	}

	return userName, taskName, nil
}