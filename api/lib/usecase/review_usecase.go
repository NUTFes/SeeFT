package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/NUTFes/SeeFT/api/lib/entity"
	rep "github.com/NUTFes/SeeFT/api/lib/internals/repository"
	"github.com/pkg/errors"
)

type reviewUseCase struct {
	reviewRep rep.ReviewRepository
	taskRep   rep.TaskRepository
}

type ReviewUseCase interface {
	GetReviewsGAS(ctx context.Context) ([]entity.ReviewGAS, error)
	GetReviewGASByID(ctx context.Context, id string) (entity.ReviewGAS, error)
	CreateReview(ctx context.Context, userID string, taskName string, staffingRating string, manualRating string, comment string) (entity.Review, error)
	UpdateReview(ctx context.Context, id string, userID string, taskName string, staffingRating string, manualRating string, comment string) (entity.Review, error)
	DeleteReview(ctx context.Context, id string) error
}

// 生成関数
func NewReviewUseCase(r rep.ReviewRepository, t rep.TaskRepository) ReviewUseCase {
	return &reviewUseCase{reviewRep: r, taskRep: t}
}

func (u *reviewUseCase) GetReviewsGAS(ctx context.Context) ([]entity.ReviewGAS, error) {
	rows, err := u.reviewRep.AllWithDetails(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "reviewRep.AllWithDetails")
	}
	defer func() { _ = rows.Close() }()

	loc, locErr := time.LoadLocation("Asia/Tokyo")

	var res []entity.ReviewGAS
	for rows.Next() {
		var r entity.ReviewGAS
		var created, updated time.Time
		if err := rows.Scan(
			&r.ID,
			&r.UserName,
			&r.UserBureau,
			&r.UserGrade,
			&r.UserStudentNo,
			&r.TaskName,
			&r.StaffingRating,
			&r.ManualRating,
			&r.Comment,
			&created,
			&updated,
		); err != nil {
			return nil, errors.Wrap(err, "rows.Scan")
		}
		if locErr == nil {
        r.CreatedAt = created.In(loc).Format("2006/01/02 15:04:05")
        r.UpdatedAt = updated.In(loc).Format("2006/01/02 15:04:05")
    	} else {
		r.CreatedAt = created.Format("2006/01/02 15:04:05")
		r.UpdatedAt = updated.Format("2006/01/02 15:04:05")
		}
		res = append(res, r)		
	}
	return res, nil
}

// 1件取得
func (u *reviewUseCase) GetReviewGASByID(ctx context.Context, id string) (entity.ReviewGAS, error) {
	row, err := u.reviewRep.FindWithDetails(ctx, id)
	if err != nil {
		return entity.ReviewGAS{}, errors.Wrap(err, "reviewRep.FindWithDetails")
	}

	loc, locErr := time.LoadLocation("Asia/Tokyo")

	var r entity.ReviewGAS
	var created, updated time.Time
	if err := row.Scan(
		&r.ID,
		&r.UserName,
		&r.UserBureau,
		&r.UserGrade,
		&r.UserStudentNo,
		&r.TaskName,
		&r.StaffingRating,
		&r.ManualRating,
		&r.Comment,
		&created,
		&updated,
	); err != nil {
		return entity.ReviewGAS{}, err
	}
	if locErr == nil {
    r.CreatedAt = created.In(loc).Format("2006/01/02 15:04:05")
    r.UpdatedAt = updated.In(loc).Format("2006/01/02 15:04:05")
	} else {
    r.CreatedAt = created.Format("2006/01/02 15:04:05")
    r.UpdatedAt = updated.Format("2006/01/02 15:04:05")
	}
	return r, nil
}

// 作成
func (u *reviewUseCase) CreateReview(c context.Context, userID string, taskName string, staffingRating string, manualRating string, comment string) (entity.Review, error) {
	latestReview := entity.Review{}

	row, err := u.taskRep.FindByName(c, taskName)
	if err != nil {
		return latestReview, errors.Wrap(err, "taskRep.FindByName")
	}
	var task entity.Task
	if err := row.Scan(
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
	); err != nil {
		if err == sql.ErrNoRows {
			return latestReview, errors.New("指定されたtask_nameのタスクが存在しません")
		}
		return latestReview, errors.Wrap(err, "row.Scan (task entity)")
	}

	if err := u.reviewRep.Create(c, userID, fmt.Sprintf("%d", task.ID), staffingRating, manualRating, comment); err != nil {
		return latestReview, errors.Wrap(err, "reviewRep.Create")
	}

	// 作成直後の最新レビューを取得して返す
	row, err = u.reviewRep.FindNewRecord(c)
	if err != nil {
		return latestReview, errors.Wrap(err, "reviewRep.FindNewRecord")
	}
	if err := row.Scan(
		&latestReview.ID,
		&latestReview.UserID,
		&latestReview.TaskID,
		&latestReview.StaffingRating,
		&latestReview.ManualRating,
		&latestReview.Comment,
		&latestReview.CreatedAt,
		&latestReview.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return latestReview, errors.New("レビューの作成後、該当データが見つかりません")
		}
		return latestReview, errors.Wrap(err, "row.Scan (review)")
	}
	return latestReview, nil
}

// 編集
func (u *reviewUseCase) UpdateReview(c context.Context, ID string, userID string, taskName string, staffingRating string, manualRating string, comment string) (entity.Review, error) {
	updatedReview := entity.Review{}
	
	row, err := u.taskRep.FindByName(c, taskName)
	if err != nil {
		return entity.Review{}, errors.Wrap(err, "taskRep.FindByName")
	}
	var task entity.Task
	if err := row.Scan(
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
	); err != nil {
		if err == sql.ErrNoRows {
			return updatedReview, errors.New("指定されたtask_nameのタスクが存在しません")
		}
		return updatedReview, errors.Wrap(err, "row.Scan (task entity)")
	}

	var review entity.Review
	row, err = u.reviewRep.Find(c, ID)
	if err != nil {
		return review, errors.Wrap(err, "reviewRep.Find (before update)")
	}
	if err := row.Scan(
		&review.ID,
		&review.UserID,
		&review.TaskID,
		&review.StaffingRating,
		&review.ManualRating,
		&review.Comment,
		&review.CreatedAt,
		&review.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return review, errors.New("指定されたIDのレビューが存在しません")
		}
		return review, errors.Wrap(err, "row.Scan (review before update)")
	}

	if err := u.reviewRep.Update(c, ID, userID, fmt.Sprintf("%d", task.ID), staffingRating, manualRating, comment); err != nil {
		return updatedReview, errors.Wrap(err, "reviewRep.Update")
	}
	row, err = u.reviewRep.Find(c, ID)
	if err != nil {
		return updatedReview, errors.Wrap(err, "reviewRep.Find (after update)")
	}
	if err := row.Scan(
		&updatedReview.ID,
		&updatedReview.UserID,
		&updatedReview.TaskID,
		&updatedReview.StaffingRating,
		&updatedReview.ManualRating,
		&updatedReview.Comment,
		&updatedReview.CreatedAt,
		&updatedReview.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return updatedReview, errors.New("更新後のレビューが見つかりません")
		}
		return updatedReview, errors.Wrap(err, "row.Scan (review after update)")
	}
	return updatedReview, nil
}

// 削除
func (u *reviewUseCase) DeleteReview(ctx context.Context, id string) error {
	if err := u.reviewRep.Delete(ctx, id); err != nil {
		return errors.Wrap(err, "reviewRep.Delete")
	}
	return nil
}
