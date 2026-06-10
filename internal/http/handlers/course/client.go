package course

import (
	"context"
	"curriculum-service/internal/domain/course"
	dtocourse "curriculum-service/internal/http/dto/course"
	"curriculum-service/internal/http/dto/resubmit"
	"curriculum-service/internal/http/dto/reviewcourse"
	"curriculum-service/internal/http/dto/reviewresult"
	"github.com/google/uuid"
)

type client interface {
	GetAllCourses(ctx context.Context, query dtocourse.GetCoursesQuery) ([]course.Course, error)
	CreateCourse(ctx context.Context, value *course.Course) (*course.Course, error)
	CreateSubscription(ctx context.Context, value *course.Subscription) (*course.Subscription, error)
	GetCourseByID(ctx context.Context, id uuid.UUID) (*course.Course, error)
	DeleteCourse(ctx context.Context, id uuid.UUID) error
	UpdateCourse(ctx context.Context, id uuid.UUID, value *course.Course) (*course.Course, error)
	GetCourseForUser(ctx context.Context, userID uuid.UUID, courseID uuid.UUID, hasFullAccess bool) (*course.CourseForUser, error)
	GetPendingCheckCourses(ctx context.Context) ([]course.Course, error)
	ReviewCourse(ctx context.Context, courseID, adminID uuid.UUID, req reviewcourse.ReviewCourseRequest) (*reviewcourse.ReviewCourseResponse, error)
	ResubmitCourseForReview(ctx context.Context, courseID, userID uuid.UUID) (*resubmit.ResubmitReviewResponse, error)
	GetCourseReview(ctx context.Context, courseID, userID uuid.UUID) (*reviewresult.CourseReviewResponse, error)
}
