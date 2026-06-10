package course

import (
	"context"
	"curriculum-service/internal/domain/course"
	"curriculum-service/internal/domain/module"
	"curriculum-service/internal/domain/review"
	"curriculum-service/internal/domain/reviewlog"
	dtocourse "curriculum-service/internal/http/dto/course"
	"github.com/google/uuid"
)

type Repository interface {
	GetAllCourses(ctx context.Context, query dtocourse.GetCoursesQuery) ([]course.Course, error)
	CreateCourse(ctx context.Context, value *course.Course) (uuid.UUID, error)
	CreateSubscription(ctx context.Context, value *course.Subscription) error
	GetSubscriptionByID(ctx context.Context, id uuid.UUID) (*course.Subscription, error)
	HasSubscription(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (bool, error)
	GetCourseByID(ctx context.Context, id uuid.UUID) (*course.Course, error)
	DeleteCourse(ctx context.Context, id uuid.UUID) error
	DeleteCoursePrice(ctx context.Context, id uuid.UUID) error
	UpdateCourse(ctx context.Context, id uuid.UUID, value *course.Course) error
	GetPendingCheckCourses(ctx context.Context) ([]course.Course, error)
	ReviewCourse(ctx context.Context, log *reviewlog.CourseReviewLog, isApproved bool) error
	SetCourseUnchecked(ctx context.Context, courseID uuid.UUID) error
	GetLatestReviewLog(ctx context.Context, courseID uuid.UUID) (*reviewlog.CourseReviewLog, error)
}
type ReviewRepository interface {
	GetAllReviewsByCourseID(ctx context.Context, courseID uuid.UUID) ([]review.CourseReview, error)
}
type ModuleRepository interface {
	GetModuleByCourseID(ctx context.Context, courseID uuid.UUID) ([]module.Module, error)
	GetLimitedModulesByCourseID(ctx context.Context, courseID uuid.UUID, limit int) ([]module.Module, error)
}

type NotificationSender interface {
	SendEvent(ctx context.Context, userID uuid.UUID, event string, data map[string]any) error
}
