package usecase

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/pkg/errors"
)

type rescueUnifiedUseCase struct {
	questionRescueRepository    repository.QuestionRescueRepository
	shorthandedRescueRepository repository.ShorthandedRescueRepository
	troubleRescueRepository     repository.TroubleRescueRepository
	userRepository              repository.UserRepository
	taskRepository              repository.TaskRepository
}

type RescueUnifiedUseCase interface {
	GetAllRescues(context.Context) ([]entity.RescueResponse, error)
	GetRescuesByUserID(context.Context, string) ([]entity.RescueResponse, error)
	SaveRescueToSpreadsheet(map[string]interface{}) error
}

func NewRescueUnifiedUseCase(
	qr repository.QuestionRescueRepository,
	sr repository.ShorthandedRescueRepository,
	tr repository.TroubleRescueRepository,
	ur repository.UserRepository,
	tar repository.TaskRepository,
) RescueUnifiedUseCase {
	return &rescueUnifiedUseCase{qr, sr, tr, ur, tar}
}

// 全件取得
func (ru *rescueUnifiedUseCase) GetAllRescues(c context.Context) ([]entity.RescueResponse, error) {
	var allRescues []entity.RescueResponse

	// Question Rescues取得
	questionRescues, err := ru.getQuestionRescues(c, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get question rescues")
	}
	allRescues = append(allRescues, questionRescues...)

	// Shorthanded Rescues取得
	shorthandedRescues, err := ru.getShorthandedRescues(c, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get shorthanded rescues")
	}
	allRescues = append(allRescues, shorthandedRescues...)

	// Trouble Rescues取得
	troubleRescues, err := ru.getTroubleRescues(c, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get trouble rescues")
	}
	allRescues = append(allRescues, troubleRescues...)

	// 時刻でソート（降順）
	sort.Slice(allRescues, func(i, j int) bool {
		return allRescues[i].Time > allRescues[j].Time
	})

	return allRescues, nil
}

// ユーザーID別取得
func (ru *rescueUnifiedUseCase) GetRescuesByUserID(c context.Context, userID string) ([]entity.RescueResponse, error) {
	var allRescues []entity.RescueResponse

	// Question Rescues取得
	questionRescues, err := ru.getQuestionRescues(c, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get question rescues")
	}
	allRescues = append(allRescues, questionRescues...)

	// Shorthanded Rescues取得
	shorthandedRescues, err := ru.getShorthandedRescues(c, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get shorthanded rescues")
	}
	allRescues = append(allRescues, shorthandedRescues...)

	// Trouble Rescues取得
	troubleRescues, err := ru.getTroubleRescues(c, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get trouble rescues")
	}
	allRescues = append(allRescues, troubleRescues...)

	// 時刻でソート（降順）
	sort.Slice(allRescues, func(i, j int) bool {
		return allRescues[i].Time > allRescues[j].Time
	})

	return allRescues, nil
}

// Question Rescues取得
func (ru *rescueUnifiedUseCase) getQuestionRescues(c context.Context, userID string) ([]entity.RescueResponse, error) {
	var rows *sql.Rows
	var err error

	if userID != "" {
		rows, err = ru.questionRescueRepository.FindByUserID(c, userID)
	} else {
		rows, err = ru.questionRescueRepository.All(c)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rescues []entity.RescueResponse
	for rows.Next() {
		var questionRescue entity.QuestionRescueForGet
		var response sql.NullString
		err := rows.Scan(&questionRescue.ID, &questionRescue.UserID, &questionRescue.Question, &questionRescue.Status, &response, &questionRescue.Time, &questionRescue.CreatedAt, &questionRescue.UpdatedAt)
		if err != nil {
			return nil, err
		}
		questionRescue.Response = response.String

		// ユーザー名を取得
		userName, err := ru.getUserName(c, strconv.Itoa(questionRescue.UserID))
		if err != nil {
			userName = "不明なユーザー"
		}

		rescue := entity.NewQuestionRescueResponse(&questionRescue, userName)
		rescues = append(rescues, *rescue)
	}
	return rescues, nil
}

// Shorthanded Rescues取得
func (ru *rescueUnifiedUseCase) getShorthandedRescues(c context.Context, userID string) ([]entity.RescueResponse, error) {
	var rows *sql.Rows
	var err error

	if userID != "" {
		rows, err = ru.shorthandedRescueRepository.FindByUserID(c, userID)
	} else {
		rows, err = ru.shorthandedRescueRepository.All(c)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rescues []entity.RescueResponse
	for rows.Next() {
		var shorthandedRescue entity.ShorthandedRescueForGet
		var place, response sql.NullString
		err := rows.Scan(&shorthandedRescue.ID, &shorthandedRescue.UserID, &shorthandedRescue.TaskID, &shorthandedRescue.MissingNumber, &place, &shorthandedRescue.Status, &response, &shorthandedRescue.Time, &shorthandedRescue.CreatedAt, &shorthandedRescue.UpdatedAt)
		if err != nil {
			return nil, err
		}
		shorthandedRescue.Place = place.String
		shorthandedRescue.Response = response.String

		// ユーザー名を取得
		userName, err := ru.getUserName(c, strconv.Itoa(shorthandedRescue.UserID))
		if err != nil {
			userName = "不明なユーザー"
		}

		// タスク名を取得
		taskName, err := ru.getTaskName(c, strconv.Itoa(shorthandedRescue.TaskID))
		if err != nil {
			taskName = "不明なタスク"
		}

		rescue := entity.NewShorthandedRescueResponse(&shorthandedRescue, userName, taskName)
		rescues = append(rescues, *rescue)
	}
	return rescues, nil
}

// Trouble Rescues取得
func (ru *rescueUnifiedUseCase) getTroubleRescues(c context.Context, userID string) ([]entity.RescueResponse, error) {
	var rows *sql.Rows
	var err error

	if userID != "" {
		rows, err = ru.troubleRescueRepository.FindByUserID(c, userID)
	} else {
		rows, err = ru.troubleRescueRepository.All(c)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rescues []entity.RescueResponse
	for rows.Next() {
		var troubleRescue entity.TroubleRescueForGet
		var place, response sql.NullString
		err := rows.Scan(&troubleRescue.ID, &troubleRescue.UserID, &troubleRescue.TaskID, &place, &troubleRescue.Detail, &troubleRescue.Status, &response, &troubleRescue.Time, &troubleRescue.CreatedAt, &troubleRescue.UpdatedAt)
		if err != nil {
			return nil, err
		}
		troubleRescue.Place = place.String
		troubleRescue.Response = response.String

		// ユーザー名を取得
		userName, err := ru.getUserName(c, strconv.Itoa(troubleRescue.UserID))
		if err != nil {
			userName = "不明なユーザー"
		}

		// タスク名を取得
		taskName, err := ru.getTaskName(c, strconv.Itoa(troubleRescue.TaskID))
		if err != nil {
			taskName = "タスク外"
		}

		rescue := entity.NewTroubleRescueResponse(&troubleRescue, userName, taskName)
		rescues = append(rescues, *rescue)
	}
	return rescues, nil
}

// ユーザー名取得
func (ru *rescueUnifiedUseCase) getUserName(c context.Context, userID string) (string, error) {
	row, err := ru.userRepository.Find(c, userID)
	if err != nil {
		return "", err
	}

	var id int
	var name, mail, studentNumber, tel, password string
	var gradeID, departmentID, bureauID, roleID int
	var createdAt, updatedAt string
	var slackUserID sql.NullString
	err = row.Scan(&id, &name, &mail, &gradeID, &departmentID, &bureauID, &roleID, &studentNumber, &tel, &password, &createdAt, &updatedAt, &slackUserID)
	if err != nil {
		return "", err
	}

	return name, nil
}

// タスク名取得
func (ru *rescueUnifiedUseCase) getTaskName(c context.Context, taskID string) (string, error) {
	row, err := ru.taskRepository.Find(c, taskID)
	if err != nil {
		return "", err
	}

	var id int
	var task, url, color, remark string
	var placeID, bureauID, maxMember, yearID int
	var createdAt, updatedAt string
	err = row.Scan(&id, &task, &placeID, &url, &bureauID, &maxMember, &color, &remark, &yearID, &createdAt, &updatedAt)
	if err != nil {
		return "", err
	}

	return task, nil
}

// スプシ保存用関数（Google Sheets APIのラッパーを想定）
func (ru *rescueUnifiedUseCase) SaveRescueToSpreadsheet(data map[string]interface{}) error {
	// Google Sheets APIで保存処理
	// GAS送信
	err := ru.SendRescueToGAS(data)
	if err != nil {
		fmt.Println("GAS送信失敗:", err)
		return err
	}
	return nil
}

// GAS送信関数
func (ru *rescueUnifiedUseCase) SendRescueToGAS(data map[string]interface{}) error {
	// GASのURLを環境変数から取得
	gasURL := os.Getenv("RESCUE_GAS_URL")
	if gasURL == "" {
		return errors.New("GAS URLが設定されていません (RESCUE_GAS_URL)")
	}

	// JSONに変換
	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "レスキューデータのJSON変換失敗")
	}

	req, err := http.NewRequest("POST", gasURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.Wrap(err, "GASリクエスト作成失敗")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "GASへの送信失敗")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("GASが非OKステータスを返しました: %d", resp.StatusCode)
	}
	return nil
}
