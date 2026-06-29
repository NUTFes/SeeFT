package usecase

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/pkg/errors"
)

type shorthandedRescueUseCase struct {
	shorthandedRescueRepository repository.ShorthandedRescueRepository
}

type ShorthandedRescueUseCase interface {
	GetShorthandedRescues(context.Context) ([]entity.ShorthandedRescueForGet, error)
	GetShorthandedRescueByID(context.Context, string) (*entity.ShorthandedRescueForGet, error)
	GetShorthandedRescuesByUserID(context.Context, string) ([]entity.ShorthandedRescueForGet, error)
	GetShorthandedRescuesByTaskID(context.Context, string) ([]entity.ShorthandedRescueForGet, error)
	CreateShorthandedRescue(context.Context, string, string, string, string, string) (*entity.ShorthandedRescueForGet, error)
	UpdateShorthandedRescue(context.Context, string, string, string) (*entity.ShorthandedRescueForGet, error)
	DeleteShorthandedRescue(context.Context, string) error
}

func NewShorthandedRescueUseCase(sr repository.ShorthandedRescueRepository) ShorthandedRescueUseCase {
	return &shorthandedRescueUseCase{sr}
}

// 全件取得
func (su *shorthandedRescueUseCase) GetShorthandedRescues(c context.Context) ([]entity.ShorthandedRescueForGet, error) {
	rows, err := su.shorthandedRescueRepository.All(c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get shorthanded rescues")
	}
	defer func() { _ = rows.Close() }()

	var shorthandedRescues []entity.ShorthandedRescueForGet
	for rows.Next() {
		var shorthandedRescue entity.ShorthandedRescueForGet
		var place, response sql.NullString
		err := rows.Scan(&shorthandedRescue.ID, &shorthandedRescue.UserID, &shorthandedRescue.TaskID, &shorthandedRescue.MissingNumber, &place, &shorthandedRescue.Status, &response, &shorthandedRescue.Time, &shorthandedRescue.CreatedAt, &shorthandedRescue.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan shorthanded rescue")
		}
		shorthandedRescue.Place = place.String
		shorthandedRescue.Response = response.String
		shorthandedRescues = append(shorthandedRescues, shorthandedRescue)
	}
	return shorthandedRescues, nil
}

// 1件取得
func (su *shorthandedRescueUseCase) GetShorthandedRescueByID(c context.Context, id string) (*entity.ShorthandedRescueForGet, error) {
	row, err := su.shorthandedRescueRepository.Find(c, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get shorthanded rescue")
	}

	var shorthandedRescue entity.ShorthandedRescueForGet
	var place, response sql.NullString
	err = row.Scan(&shorthandedRescue.ID, &shorthandedRescue.UserID, &shorthandedRescue.TaskID, &shorthandedRescue.MissingNumber, &place, &shorthandedRescue.Status, &response, &shorthandedRescue.Time, &shorthandedRescue.CreatedAt, &shorthandedRescue.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to scan shorthanded rescue")
	}
	shorthandedRescue.Place = place.String
	shorthandedRescue.Response = response.String
	return &shorthandedRescue, nil
}

// ユーザIDで検索
func (su *shorthandedRescueUseCase) GetShorthandedRescuesByUserID(c context.Context, userID string) ([]entity.ShorthandedRescueForGet, error) {
	rows, err := su.shorthandedRescueRepository.FindByUserID(c, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get shorthanded rescues by user ID")
	}
	defer func() { _ = rows.Close() }()

	var shorthandedRescues []entity.ShorthandedRescueForGet
	for rows.Next() {
		var shorthandedRescue entity.ShorthandedRescueForGet
		var place, response sql.NullString
		err := rows.Scan(&shorthandedRescue.ID, &shorthandedRescue.UserID, &shorthandedRescue.TaskID, &shorthandedRescue.MissingNumber, &place, &shorthandedRescue.Status, &response, &shorthandedRescue.Time, &shorthandedRescue.CreatedAt, &shorthandedRescue.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan shorthanded rescue")
		}
		shorthandedRescue.Place = place.String
		shorthandedRescue.Response = response.String
		shorthandedRescues = append(shorthandedRescues, shorthandedRescue)
	}
	return shorthandedRescues, nil
}

// タスクIDで検索
func (su *shorthandedRescueUseCase) GetShorthandedRescuesByTaskID(c context.Context, taskID string) ([]entity.ShorthandedRescueForGet, error) {
	rows, err := su.shorthandedRescueRepository.FindByTaskID(c, taskID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get shorthanded rescues by task ID")
	}
	defer func() { _ = rows.Close() }()

	var shorthandedRescues []entity.ShorthandedRescueForGet
	for rows.Next() {
		var shorthandedRescue entity.ShorthandedRescueForGet
		var place, response sql.NullString
		err := rows.Scan(&shorthandedRescue.ID, &shorthandedRescue.UserID, &shorthandedRescue.TaskID, &shorthandedRescue.MissingNumber, &place, &shorthandedRescue.Status, &response, &shorthandedRescue.Time, &shorthandedRescue.CreatedAt, &shorthandedRescue.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan shorthanded rescue")
		}
		shorthandedRescue.Place = place.String
		shorthandedRescue.Response = response.String
		shorthandedRescues = append(shorthandedRescues, shorthandedRescue)
	}
	return shorthandedRescues, nil
}

// 作成
func (su *shorthandedRescueUseCase) CreateShorthandedRescue(c context.Context, userID string, taskID string, missingNumber string, place string, status string) (*entity.ShorthandedRescueForGet, error) {
	// 入力バリデーション
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if taskID == "" {
		return nil, errors.New("task ID is required")
	}
	if missingNumber == "" {
		return nil, errors.New("missing number is required")
	}
	if status == "" {
		status = "todo"
	}

	// 数値チェック
	if _, err := strconv.Atoi(userID); err != nil {
		return nil, errors.New("invalid user ID")
	}
	if _, err := strconv.Atoi(taskID); err != nil {
		return nil, errors.New("invalid task ID")
	}
	if _, err := strconv.Atoi(missingNumber); err != nil {
		return nil, errors.New("invalid missing number")
	}

	err := su.shorthandedRescueRepository.Create(c, userID, taskID, missingNumber, place, status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create shorthanded rescue")
	}

	// 作成したレコードを取得
	row, err := su.shorthandedRescueRepository.FindNewRecord(c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get created shorthanded rescue")
	}

	var shorthandedRescue entity.ShorthandedRescueForGet
	var placeNull, response sql.NullString
	err = row.Scan(&shorthandedRescue.ID, &shorthandedRescue.UserID, &shorthandedRescue.TaskID, &shorthandedRescue.MissingNumber, &placeNull, &shorthandedRescue.Status, &response, &shorthandedRescue.Time, &shorthandedRescue.CreatedAt, &shorthandedRescue.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to scan created shorthanded rescue")
	}
	shorthandedRescue.Place = placeNull.String
	shorthandedRescue.Response = response.String
	return &shorthandedRescue, nil
}

// 更新
func (su *shorthandedRescueUseCase) UpdateShorthandedRescue(c context.Context, id string, status string, response string) (*entity.ShorthandedRescueForGet, error) {
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

	err := su.shorthandedRescueRepository.Update(c, id, status, response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update shorthanded rescue")
	}

	// 更新したレコードを取得
	return su.GetShorthandedRescueByID(c, id)
}

// 削除
func (su *shorthandedRescueUseCase) DeleteShorthandedRescue(c context.Context, id string) error {
	// 入力バリデーション
	if id == "" {
		return errors.New("ID is required")
	}

	// IDの数値チェック
	if _, err := strconv.Atoi(id); err != nil {
		return errors.New("invalid ID")
	}

	err := su.shorthandedRescueRepository.Delete(c, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete shorthanded rescue")
	}
	return nil
}
