package usecase

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/NUTFes/SeeFT/api/lib/entity"
	"github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/pkg/errors"
)

type troubleRescueUseCase struct {
	troubleRescueRepository repository.TroubleRescueRepository
}

type TroubleRescueUseCase interface {
	GetTroubleRescues(context.Context) ([]entity.TroubleRescueForGet, error)
	GetTroubleRescueByID(context.Context, string) (*entity.TroubleRescueForGet, error)
	GetTroubleRescuesByUserID(context.Context, string) ([]entity.TroubleRescueForGet, error)
	GetTroubleRescuesByTaskID(context.Context, string) ([]entity.TroubleRescueForGet, error)
	CreateTroubleRescue(context.Context, string, string, string, string, string) (*entity.TroubleRescueForGet, error)
	UpdateTroubleRescue(context.Context, string, string, string) (*entity.TroubleRescueForGet, error)
	DeleteTroubleRescue(context.Context, string) error
}

func NewTroubleRescueUseCase(tr repository.TroubleRescueRepository) TroubleRescueUseCase {
	return &troubleRescueUseCase{tr}
}

// 全件取得
func (tu *troubleRescueUseCase) GetTroubleRescues(c context.Context) ([]entity.TroubleRescueForGet, error) {
	rows, err := tu.troubleRescueRepository.All(c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get trouble rescues")
	}
	defer rows.Close()

	var troubleRescues []entity.TroubleRescueForGet
	for rows.Next() {
		var troubleRescue entity.TroubleRescueForGet
		var place, response sql.NullString
		err := rows.Scan(&troubleRescue.ID, &troubleRescue.UserID, &troubleRescue.TaskID, &place, &troubleRescue.Detail, &troubleRescue.Status, &response, &troubleRescue.Time, &troubleRescue.CreatedAt, &troubleRescue.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan trouble rescue")
		}
		troubleRescue.Place = place.String
		troubleRescue.Response = response.String
		troubleRescues = append(troubleRescues, troubleRescue)
	}
	return troubleRescues, nil
}

// 1件取得
func (tu *troubleRescueUseCase) GetTroubleRescueByID(c context.Context, id string) (*entity.TroubleRescueForGet, error) {
	row, err := tu.troubleRescueRepository.Find(c, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get trouble rescue")
	}

	var troubleRescue entity.TroubleRescueForGet
	var place, response sql.NullString
	err = row.Scan(&troubleRescue.ID, &troubleRescue.UserID, &troubleRescue.TaskID, &place, &troubleRescue.Detail, &troubleRescue.Status, &response, &troubleRescue.Time, &troubleRescue.CreatedAt, &troubleRescue.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to scan trouble rescue")
	}
	troubleRescue.Place = place.String
	troubleRescue.Response = response.String
	return &troubleRescue, nil
}

// ユーザIDで検索
func (tu *troubleRescueUseCase) GetTroubleRescuesByUserID(c context.Context, userID string) ([]entity.TroubleRescueForGet, error) {
	rows, err := tu.troubleRescueRepository.FindByUserID(c, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get trouble rescues by user ID")
	}
	defer rows.Close()

	var troubleRescues []entity.TroubleRescueForGet
	for rows.Next() {
		var troubleRescue entity.TroubleRescueForGet
		var place, response sql.NullString
		err := rows.Scan(&troubleRescue.ID, &troubleRescue.UserID, &troubleRescue.TaskID, &place, &troubleRescue.Detail, &troubleRescue.Status, &response, &troubleRescue.Time, &troubleRescue.CreatedAt, &troubleRescue.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan trouble rescue")
		}
		troubleRescue.Place = place.String
		troubleRescue.Response = response.String
		troubleRescues = append(troubleRescues, troubleRescue)
	}
	return troubleRescues, nil
}

// タスクIDで検索
func (tu *troubleRescueUseCase) GetTroubleRescuesByTaskID(c context.Context, taskID string) ([]entity.TroubleRescueForGet, error) {
	rows, err := tu.troubleRescueRepository.FindByTaskID(c, taskID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get trouble rescues by task ID")
	}
	defer rows.Close()

	var troubleRescues []entity.TroubleRescueForGet
	for rows.Next() {
		var troubleRescue entity.TroubleRescueForGet
		var place, response sql.NullString
		err := rows.Scan(&troubleRescue.ID, &troubleRescue.UserID, &troubleRescue.TaskID, &place, &troubleRescue.Detail, &troubleRescue.Status, &response, &troubleRescue.Time, &troubleRescue.CreatedAt, &troubleRescue.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan trouble rescue")
		}
		troubleRescue.Place = place.String
		troubleRescue.Response = response.String
		troubleRescues = append(troubleRescues, troubleRescue)
	}
	return troubleRescues, nil
}

// 作成
func (tu *troubleRescueUseCase) CreateTroubleRescue(c context.Context, userID string, taskID string, place string, detail string, status string) (*entity.TroubleRescueForGet, error) {
	// 入力バリデーション
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	if taskID == "" {
		return nil, errors.New("task ID is required")
	}
	if detail == "" {
		return nil, errors.New("detail is required")
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

	err := tu.troubleRescueRepository.Create(c, userID, taskID, place, detail, status)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create trouble rescue")
	}

	// 作成したレコードを取得
	row, err := tu.troubleRescueRepository.FindNewRecord(c)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get created trouble rescue")
	}

	var troubleRescue entity.TroubleRescueForGet
	var placeNull, response sql.NullString
	err = row.Scan(&troubleRescue.ID, &troubleRescue.UserID, &troubleRescue.TaskID, &placeNull, &troubleRescue.Detail, &troubleRescue.Status, &response, &troubleRescue.Time, &troubleRescue.CreatedAt, &troubleRescue.UpdatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to scan created trouble rescue")
	}
	troubleRescue.Place = placeNull.String
	troubleRescue.Response = response.String
	return &troubleRescue, nil
}

// 更新
func (tu *troubleRescueUseCase) UpdateTroubleRescue(c context.Context, id string, status string, response string) (*entity.TroubleRescueForGet, error) {
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

	err := tu.troubleRescueRepository.Update(c, id, status, response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update trouble rescue")
	}

	// 更新したレコードを取得
	return tu.GetTroubleRescueByID(c, id)
}

// 削除
func (tu *troubleRescueUseCase) DeleteTroubleRescue(c context.Context, id string) error {
	// 入力バリデーション
	if id == "" {
		return errors.New("ID is required")
	}

	// IDの数値チェック
	if _, err := strconv.Atoi(id); err != nil {
		return errors.New("invalid ID")
	}

	err := tu.troubleRescueRepository.Delete(c, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete trouble rescue")
	}
	return nil
}
