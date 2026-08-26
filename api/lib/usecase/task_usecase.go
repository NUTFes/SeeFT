package usecase

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/NUTFes/SeeFT/api/lib/entity"
	rep "github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/pkg/errors"
)

type taskUseCase struct {
	rep      rep.TaskRepository
	placeRep rep.PlaceRepository
}

type TaskUseCase interface {
	GetTasks(context.Context) ([]entity.Task, error)
	GetTaskByID(context.Context, string) (entity.Task, error)
	GetTasksByShift(context.Context, string) ([]entity.Task, error)
	GetTasksByUserID(context.Context, string) ([]entity.Task, error)
	CreateTask(context.Context, string, string, string, string, string, string, string, string) (entity.Task, error)
	UpdateTask(context.Context, string, string, string, string, string, string, string, string, string) (entity.Task, error)
	DeleteTask(context.Context, string) error
	UpdateTasksAndPlacesFromGAS(context.Context, entity.TaskAndPlaceChangeRequest) error
}

func NewTaskUseCase(
	rep rep.TaskRepository,
	placeRep rep.PlaceRepository) TaskUseCase {
	return &taskUseCase{rep, placeRep}
}

func (b *taskUseCase) GetTasks(c context.Context) ([]entity.Task, error) {
	task := entity.Task{}
	var tasks []entity.Task

	rows, err := b.rep.All(c)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		err := rows.Scan(
			&task.ID,
			&task.Task,
			&task.PlaceID,
			&task.Url,
			&task.ManualUrl,
			&task.BureauID,
			&task.MaxMember,
			&task.Color,
			&task.Remark,
			&task.YearID,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot connect SQL")
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (b *taskUseCase) GetTaskByID(c context.Context, id string) (entity.Task, error) {
	var task entity.Task
	row, err := b.rep.Find(c, id)
	if err != nil {
		return task, err
	}

	err = row.Scan(
		&task.ID,
		&task.Task,
		&task.PlaceID,
		&task.Url,
		&task.ManualUrl,
		&task.BureauID,
		&task.MaxMember,
		&task.Color,
		&task.Remark,
		&task.YearID,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return task, err
	}

	return task, nil
}

func (b *taskUseCase) GetTasksByShift(c context.Context, shift string) ([]entity.Task, error) {
	task := entity.Task{}
	var tasks []entity.Task

	rows, err := b.rep.Shift(c, shift)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		err := rows.Scan(
			&task.ID,
			&task.Task,
			&task.PlaceID,
			&task.Url,
			&task.ManualUrl,
			&task.BureauID,
			&task.MaxMember,
			&task.Color,
			&task.Remark,
			&task.YearID,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot connect SQL")
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// 指定したuserIDのタスクを取得
func (b *taskUseCase) GetTasksByUserID(c context.Context, userID string) ([]entity.Task, error) {
	task := entity.Task{}
	var tasks []entity.Task

	rows, err := b.rep.FindByUserID(c, userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		err := rows.Scan(
			&task.ID,
			&task.Task,
			&task.PlaceID,
			&task.Url,
			&task.ManualUrl,
			&task.BureauID,
			&task.MaxMember,
			&task.Color,
			&task.Remark,
			&task.YearID,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot connect SQL")
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (u *taskUseCase) CreateTask(c context.Context, name string, placeID string, url string, bureauID string, maxMember string, color string, remark string, yearID string) (entity.Task, error) {
	latasTask := entity.Task{}
	// 管理画面からの新規作成ではマニュアル（スライド版）を指定できないため空で作る
	const manualURL = ""
	if err := u.rep.Create(c, name, placeID, url, manualURL, bureauID, maxMember, color, remark, yearID); err != nil {
		return latasTask, err
	}
	row, err := u.rep.FindNewRecord(c)
	if err != nil {
		return latasTask, err
	}
	err = row.Scan(
		&latasTask.ID,
		&latasTask.Task,
		&latasTask.PlaceID,
		&latasTask.Url,
		&latasTask.ManualUrl,
		&latasTask.BureauID,
		&latasTask.MaxMember,
		&latasTask.Color,
		&latasTask.Remark,
		&latasTask.YearID,
		&latasTask.CreatedAt,
		&latasTask.UpdatedAt,
	)
	if err != nil {
		return latasTask, err
	}
	return latasTask, err
}

func (u *taskUseCase) UpdateTask(c context.Context, id string, name string, placeID string, url string, bureauID string, maxMember string, color string, remark string, yearID string) (entity.Task, error) {
	updatedTask := entity.Task{}
	var task entity.Task

	row, err := u.rep.Find(c, id)
	if err != nil {
		return updatedTask, err
	}
	err = row.Scan(
		&task.ID,
		&task.Task,
		&task.PlaceID,
		&task.Url,
		&task.ManualUrl,
		&task.BureauID,
		&task.MaxMember,
		&task.Color,
		&task.Remark,
		&task.YearID,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return task, err
	}

	// 管理画面からの編集ではマニュアル（スライド版）を指定できない。空で上書きすると
	// タスク名を直しただけで導線が消えるため、更新前に読み出した既存値を引き継ぐ
	if err = u.rep.Update(c, id, name, placeID, url, task.ManualUrl, bureauID, maxMember, color, remark, yearID); err != nil {
		return task, err
	}
	row, err = u.rep.Find(c, id)
	if err != nil {
		return updatedTask, err
	}
	err = row.Scan(
		&updatedTask.ID,
		&updatedTask.Task,
		&updatedTask.PlaceID,
		&updatedTask.Url,
		&updatedTask.ManualUrl,
		&updatedTask.BureauID,
		&updatedTask.MaxMember,
		&updatedTask.Color,
		&updatedTask.Remark,
		&updatedTask.YearID,
		&updatedTask.CreatedAt,
		&updatedTask.UpdatedAt,
	)
	if err != nil {
		return updatedTask, err
	}
	return updatedTask, nil
}

func (u *taskUseCase) DeleteTask(c context.Context, id string) error {
	err := u.rep.Destroy(c, id)
	return err
}

// GASからのタスクと集合場所変更通知を受けてDBを更新
func (u *taskUseCase) UpdateTasksAndPlacesFromGAS(ctx context.Context, req entity.TaskAndPlaceChangeRequest) error {
	for _, change := range req.Changes {
		// YearIDを取得
		yearID := strings.ReplaceAll(strconv.Itoa(change.YearID), " ", "")
		yearID = strings.ReplaceAll(yearID, "　", "")

		// 局名からBureauIDを取得
		var bureauID string
		bureau := strings.ReplaceAll(change.Bureau, " ", "")
		bureau = strings.ReplaceAll(bureau, "　", "")
		switch bureau {
		case `執行部`:
			bureauID = "1"
		case `執行部補佐`:
			bureauID = "2"
		case `総務局`:
			bureauID = "3"
		case `企画局`:
			bureauID = "4"
		case `渉外局`:
			bureauID = "5"
		case `財務局`:
			bureauID = "6"
		case `制作局`:
			bureauID = "7"
		case `情報局`:
			bureauID = "8"
		case `産学局`:
			bureauID = "9"
		default:
			bureauID = "1" // デフォルトは執行部
		}

		// 集合場所名からPlaceIDを取得
		placeName := strings.ReplaceAll(change.Place, " ", "")
		placeName = strings.ReplaceAll(placeName, "　", "")
		var placeID string
		if placeName == "" {
			// 集合場所が空の場合はデフォルトの集合場所を設定
			placeID = "1" // デフォルトの集合場所ID
		} else {
			placeRow, err := u.placeRep.FindByName(ctx, placeName)
			if err != nil {
				return errors.Wrapf(err, "集合場所検索失敗: %v", placeName)
			}
			var place entity.Place
			if err := placeRow.Scan(&place.ID, &place.Place, &place.Remark, &place.CreatedAt, &place.UpdatedAt); err == nil {
				// 集合場所が存在する場合は更新
				remark := ""
				if err := u.placeRep.Update(ctx, strconv.Itoa(place.ID), placeName, remark); err != nil {
					return errors.Wrapf(err, "集合場所更新失敗: %v", placeName)
				}
			} else if errors.Is(err, sql.ErrNoRows) {
				// 集合場所が存在しない場合は新規作成
				remark := ""
				createErr := u.placeRep.Create(ctx, placeName, remark)
				if createErr != nil {
					return errors.Wrapf(createErr, "集合場所新規作成失敗: %v", change.Place)
				}
				// 再取得
				placeRow, err = u.placeRep.FindByName(ctx, placeName)
				if err != nil {
					return errors.Wrapf(err, "集合場所再検索失敗: %v", placeName)
				}
				if err := placeRow.Scan(&place.ID, &place.Place, &place.Remark, &place.CreatedAt, &place.UpdatedAt); err != nil {
					return errors.Wrapf(err, "集合場所再取得失敗: %v", change.Place)
				}
			} else {
				return errors.Wrapf(err, "集合場所取得失敗: %v", change.Place)
			}
			placeID = strconv.Itoa(place.ID)
		}

		// 3. Task名からTaskID取得
		taskName := strings.ReplaceAll(change.TaskName, "　", " ") // 全角スペースを半角スペースに変換
		// taskName := strings.ReplaceAll(change.TaskName, " ", "")
		// taskName = strings.ReplaceAll(taskName, "　", "")
		taskRow, err := u.rep.FindByName(ctx, taskName)
		if err != nil {
			return errors.Wrapf(err, "タスク検索失敗: %v", taskName)
		}
		var task entity.Task
		if err := taskRow.Scan(&task.ID, &task.Task, &task.PlaceID, &task.Url, &task.ManualUrl, &task.BureauID, &task.MaxMember, &task.Color, &task.Remark, &task.YearID, &task.CreatedAt, &task.UpdatedAt); err == nil {
			// タスクが存在する場合は更新
			color := "ffffff"
			yearID := yearID
			remark := ""
			if err := u.rep.Update(ctx, strconv.Itoa(task.ID), taskName, placeID, change.Url, change.ManualUrl, bureauID, strconv.Itoa(change.MaxMember), color, remark, yearID); err != nil {
				return errors.Wrapf(err, "タスク更新失敗: %v", taskName)
			}
		} else if errors.Is(err, sql.ErrNoRows) {
			// タスクが存在しない場合は新規作成
			color := "ffffff"
			remark := ""
			createErr := u.rep.Create(ctx, taskName, placeID, change.Url, change.ManualUrl, bureauID, strconv.Itoa(change.MaxMember), color, remark, yearID)
			if createErr != nil {
				return errors.Wrapf(createErr, "タスク新規作成失敗: %v", change.TaskName)
			}
			// // 再取得
			// taskRow, _ = u.rep.FindByName(ctx, taskName)
			// if err := taskRow.Scan(&task.ID, &task.Task, &task.PlaceID, &task.Url, &task.BureauID, &task.MaxMember, &task.Color, &task.Remark, &task.YearID, &task.CreatedAt, &task.UpdatedAt); err != nil {
			// 	return errors.Wrapf(err, "タスク再取得失敗: %v", change.TaskName)
			// }
		} else {
			return errors.Wrapf(err, "タスク取得失敗: %v", change.TaskName)
		}
	}
	return nil
}
