package review

import (
	"context"
	"curriculum-service/internal/domain/review"
	"github.com/google/uuid"
)

type client interface {
	GetReviewByID(ctx context.Context, id uuid.UUID) (*review.CourseReview, error)
	CreateReview(ctx context.Context, value *review.CourseReview) (*review.CourseReview, error)
	UpdateReview(ctx context.Context, id uuid.UUID, value *review.CourseReview) (*review.CourseReview, error)
	DeleteReview(ctx context.Context, id uuid.UUID) error
	GetAllReviewsByCourseID(ctx context.Context, courseID uuid.UUID) ([]review.CourseReview, error)
}
