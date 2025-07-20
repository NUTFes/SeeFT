package usecase

import (
	"context"
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
	defer rows.Close()

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
		r.CreatedAt = created.Format("2006/01/02 15:04:05")
		r.UpdatedAt = updated.Format("2006/01/02 15:04:05")
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
	r.CreatedAt = created.Format("2006/01/02 15:04:05")
	r.UpdatedAt = updated.Format("2006/01/02 15:04:05")
	return r, nil
}

// 作成
func (u *reviewUseCase) CreateReview(c context.Context, userID string, taskName string, staffingRating string, manualRating string, comment string) (entity.Review, error) {
	latestReview := entity.Review{}

	row, err := u.taskRep.FindByName(c, taskName)
	var taskID string
	row.Scan(&taskID)

	err = u.reviewRep.Create(c, userID, taskID, staffingRating, manualRating, comment)

	row, err = u.reviewRep.Find(c, userID)
	err = row.Scan(
		&latestReview.UserID,
		&latestReview.TaskID,
		&latestReview.StaffingRating,
		&latestReview.ManualRating,
		&latestReview.Comment,
	)
	if err != nil {
		return latestReview, err
	}
	return latestReview, nil
}

// 編集
func (u *reviewUseCase) UpdateReview(c context.Context, ID string, userID string, taskName string, staffingRating string, manualRating string, comment string) (entity.Review, error) {
	row, err := u.taskRep.FindByName(c, taskName)
	var taskID string
	row.Scan(&taskID)

	updatedReview := entity.Review{}
	var review entity.Review

	row, err = u.reviewRep.Find(c, ID)
	err = row.Scan(
		&review.ID,
		&review.UserID,
		&review.TaskID,
		&review.StaffingRating,
		&review.ManualRating,
		&review.Comment,
	)
	if err != nil {
		return review, err
	}

	u.reviewRep.Update(c, ID, userID, taskID, staffingRating, manualRating, comment)
	row, err = u.reviewRep.Find(c, ID)
	err = row.Scan(
		&updatedReview.ID,
		&updatedReview.UserID,
		&updatedReview.TaskID,
		&updatedReview.StaffingRating,
		&updatedReview.ManualRating,
		&updatedReview.Comment,
		&updatedReview.CreatedAt,
		&updatedReview.UpdatedAt,
	)
	if err != nil {
		return updatedReview, err
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
