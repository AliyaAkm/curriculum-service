package review

import (
	"context"
	"curriculum-service/internal/domain"
	"curriculum-service/internal/domain/review"
	"github.com/google/uuid"
)

func (u *UseCase) CreateReview(ctx context.Context, value *review.CourseReview) (*review.CourseReview, error) {
	value.ID = uuid.New()

	exists, err := u.repo.ReviewExistsByCourseAndUser(ctx, value.CourseID, value.UserID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrReviewAlreadyExists
	}
	id, err := u.repo.CreateReview(ctx, value)
	if err != nil {
		return nil, err
	}
	return u.repo.GetReviewByID(ctx, id)
}

func (u *UseCase) GetReviewByID(ctx context.Context, id uuid.UUID) (*review.CourseReview, error) {
	value, err := u.repo.GetReviewByID(ctx, id)
	if err != nil {
		return nil, err
	}
	value.ViewCount++
	err = u.repo.UpdateReview(ctx, id, value)
	if err != nil {
		return nil, err
	}
	return u.repo.GetReviewByID(ctx, id)
}

func (u *UseCase) UpdateReview(ctx context.Context, id uuid.UUID, newValue *review.CourseReview) (*review.CourseReview, error) {
	newValue.ID = id

	oldValue, err := u.repo.GetReviewByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if newValue.Rating != 0 {
		oldValue.Rating = newValue.Rating
	}
	if newValue.Comment != "" {
		oldValue.Comment = newValue.Comment
	}

	err = u.repo.UpdateReview(ctx, id, oldValue)
	if err != nil {
		return nil, err
	}

	return u.repo.GetReviewByID(ctx, id)
}

func (u *UseCase) DeleteReview(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteReview(ctx, id)
}

func (u *UseCase) GetAllReviewsByCourseID(ctx context.Context, courseID uuid.UUID) ([]review.CourseReview, error) {
	value, err := u.repo.GetAllReviewsByCourseID(ctx, courseID)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(value); i++ {
		value[i].ViewCount++
		err = u.repo.UpdateReview(ctx, value[i].ID, &value[i])
		if err != nil {
			return nil, err
		}
	}
	return u.repo.GetAllReviewsByCourseID(ctx, courseID)
}
