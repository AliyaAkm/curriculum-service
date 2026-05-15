package progress

import (
	"context"
	progressdomain "curriculum-service/internal/domain/progress"

	"github.com/google/uuid"
)

type client interface {
	CompleteLesson(ctx context.Context, userID uuid.UUID, lessonID uuid.UUID) (*progressdomain.CourseProgress, error)
	GetCourseProgress(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (*progressdomain.CourseProgress, error)
	ListCourseProgress(ctx context.Context, userID uuid.UUID) ([]progressdomain.CourseProgress, error)
}
