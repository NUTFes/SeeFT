package usecase

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/pkg/errors"
)

type questionRescueUseCase struct {
	questionRescueRepository repository.QuestionRescueRepository
}

type QuestionRescueUseCase interface {
	GetQuestionRescues(context.Context) ([]entity.QuestionRescueForGet, error)
	GetQuestionRescueByID(context.Context, string) (*entity.QuestionRescueForGet, error)
	GetQuestionRescuesByUserID(context.Context, string) ([]entity.QuestionRescueForGet, error)
	CreateQuestionRescue(context.Context, string, string, string) (*entity.QuestionRescueForGet, error)
	UpdateQuestionRescue(context.Context, string, string, string) (*entity.QuestionRescueForGet, error)
	DeleteQuestionRescue(context.Context, string) error
}

func NewQuestionRescueUseCase(qr repository.QuestionRescueRepository) QuestionRescueUseCase {
	return &questionRescueUseCase{qr}
}

// 全件取得
func (qu *questionRescueUseCase) GetQuestionRescues(c context.Context) ([]entity.QuestionRescueForGet, error) {
	rows, err := qu.questionRescueRepository.All(c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get question rescues")
	}
	defer rows.Close()

	var questionRescues []entity.QuestionRescueForGet
	for rows.Next() {
		var questionRescue entity.QuestionRescueForGet
		var response sql.NullString
		err := rows.Scan(&questionRescue.ID, &questionRescue.UserID, &questionRescue.Question, &questionRescue.Status, &response, &questionRescue.Time, &questionRescue.CreatedAt, &questionRescue.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan question rescue")
		}
		questionRescue.Response = response.String
		questionRescues = append(questionRescues, questionRescue)
	}
	return questionRescues, nil
}

// 1件取得
func (qu *questionRescueUseCase) GetQuestionRescueByID(c context.Context, id string) (*entity.QuestionRescueForGet, error) {
	row, err := qu.questionRescueRepository.Find(c, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get question rescue")
	}

	var questionRescue entity.QuestionRescueForGet
	var response sql.NullString
	err = row.Scan(&questionRescue.ID, &questionRescue.UserID, &questionRescue.Question, &questionRescue.Status, &response, &questionRescue.Time, &questionRescue.CreatedAt, &questionRescue.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to scan question rescue")
	}
	questionRescue.Response = response.String
	return &questionRescue, nil
}

// ユーザIDで検索
func (qu *questionRescueUseCase) GetQuestionRescuesByUserID(c context.Context, userID string) ([]entity.QuestionRescueForGet, error) {
	rows, err := qu.questionRescueRepository.FindByUserID(c, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get question rescues by user ID")
	}
	defer rows.Close()

	var questionRescues []entity.QuestionRescueForGet
	for rows.Next() {
		var questionRescue entity.QuestionRescueForGet
		var response sql.NullString
		err := rows.Scan(&questionRescue.ID, &questionRescue.UserID, &questionRescue.Question, &questionRescue.Status, &response, &questionRescue.Time, &questionRescue.CreatedAt, &questionRescue.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan question rescue")
		}
		questionRescue.Response = response.String
		questionRescues = append(questionRescues, questionRescue)
	}
	return questionRescues, nil
}

// 作成
func (qu *questionRescueUseCase) CreateQuestionRescue(c context.Context, userID string, question string, status string) (*entity.QuestionRescueForGet, error) {
	// 入力バリデーション
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if question == "" {
		return nil, errors.New("question is required")
	}
	if status == "" {
		status = "todo"
	}

	// ユーザIDの数値チェック
	if _, err := strconv.Atoi(userID); err != nil {
		return nil, errors.New("invalid user ID")
	}

	err := qu.questionRescueRepository.Create(c, userID, question, status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create question rescue")
	}

	// 作成したレコードを取得
	row, err := qu.questionRescueRepository.FindNewRecord(c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get created question rescue")
	}

	var questionRescue entity.QuestionRescueForGet
	var response sql.NullString
	err = row.Scan(&questionRescue.ID, &questionRescue.UserID, &questionRescue.Question, &questionRescue.Status, &response, &questionRescue.Time, &questionRescue.CreatedAt, &questionRescue.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to scan created question rescue")
	}
	questionRescue.Response = response.String
	return &questionRescue, nil
}

// 更新
func (qu *questionRescueUseCase) UpdateQuestionRescue(c context.Context, id string, status string, response string) (*entity.QuestionRescueForGet, error) {
	// 入力バリデーション
	if id == "" {
		return nil, errors.New("ID is required")
	}
	if status == "" {
		return nil, errors.New("status is required")
	}

	// IDの数値チェック
	if _, err := strconv.Atoi(id); err != nil {
		return nil, errors.New("invalid ID")
	}

	err := qu.questionRescueRepository.Update(c, id, status, response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update question rescue")
	}

	// 更新したレコードを取得
	return qu.GetQuestionRescueByID(c, id)
}

// 削除
func (qu *questionRescueUseCase) DeleteQuestionRescue(c context.Context, id string) error {
	// 入力バリデーション
	if id == "" {
		return errors.New("ID is required")
	}

	// IDの数値チェック
	if _, err := strconv.Atoi(id); err != nil {
		return errors.New("invalid ID")
	}

	err := qu.questionRescueRepository.Delete(c, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete question rescue")
	}
	return nil
}
