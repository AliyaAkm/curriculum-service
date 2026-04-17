package review

import (
	"context"
	"curriculum-service/internal/domain/review"
	"github.com/google/uuid"
)

func (r *Repo) GetReviewByID(ctx context.Context, id uuid.UUID) (*review.CourseReview, error) {
	var entity review.CourseReview

	err := r.db.WithContext(ctx).Preload("User").First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
func (r *Repo) CreateReview(ctx context.Context, value *review.CourseReview) (uuid.UUID, error) {
	err := r.db.WithContext(ctx).Create(value).Error
	if err != nil {
		return uuid.UUID{}, err
	}
	return value.ID, nil
}
func (r *Repo) UpdateReview(ctx context.Context, id uuid.UUID, value *review.CourseReview) error {
	err := r.db.WithContext(ctx).Model(&review.CourseReview{}).Where("id = ?", id).Updates(&value).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) DeleteReview(ctx context.Context, id uuid.UUID) error {
	var entity review.CourseReview
	err := r.db.WithContext(ctx).Delete(&entity, "id = ?", id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetAllReviewsByCourseID(ctx context.Context, courseID uuid.UUID) ([]review.CourseReview, error) {
	var entity []review.CourseReview
	err := r.db.WithContext(ctx).Where("course_id = ?", courseID).Order("created_at DESC").Find(&entity).Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}
func (r *Repo) ReviewExistsByCourseAndUser(ctx context.Context, courseID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("course_reviews").
		Where("course_id = ? AND user_id = ?", courseID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
