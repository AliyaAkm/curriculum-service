package review

import (
	"context"
	"curriculum-service/internal/domain/review"
	"github.com/google/uuid"
)

type Repository interface {
	GetReviewByID(ctx context.Context, id uuid.UUID) (*review.CourseReview, error)
	CreateReview(ctx context.Context, value *review.CourseReview) (uuid.UUID, error)
	ReviewExistsByCourseAndUser(ctx context.Context, courseID, userID uuid.UUID) (bool, error)
	UpdateReview(ctx context.Context, id uuid.UUID, value *review.CourseReview) error
	DeleteReview(ctx context.Context, id uuid.UUID) error
	GetAllReviewsByCourseID(ctx context.Context, courseID uuid.UUID) ([]review.CourseReview, error)
}
