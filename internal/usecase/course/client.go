package course

import (
	"context"
	"curriculum-service/internal/domain/course"
	"curriculum-service/internal/domain/review"
	dtocourse "curriculum-service/internal/http/dto/course"
	"github.com/google/uuid"
)

type Repository interface {
	GetAllCourses(ctx context.Context, query dtocourse.GetCoursesQuery) ([]course.Course, error)
	CreateCourse(ctx context.Context, value *course.Course) (uuid.UUID, error)
	CreateSubscription(ctx context.Context, value *course.Subscription) error
	GetSubscriptionByID(ctx context.Context, id uuid.UUID) (*course.Subscription, error)
	GetCourseByID(ctx context.Context, id uuid.UUID) (*course.Course, error)
	DeleteCourse(ctx context.Context, id uuid.UUID) error
	DeleteCoursePrice(ctx context.Context, id uuid.UUID) error
	UpdateCourse(ctx context.Context, id uuid.UUID, value *course.Course) error
}
type ReviewRepository interface {
	GetAllReviewsByCourseID(ctx context.Context, courseID uuid.UUID) ([]review.CourseReview, error)
}
