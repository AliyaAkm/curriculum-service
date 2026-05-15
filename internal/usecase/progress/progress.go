package progress

import (
	"context"
	progressdomain "curriculum-service/internal/domain/progress"

	"github.com/google/uuid"
)

func (u *UseCase) CompleteLesson(ctx context.Context, userID uuid.UUID, lessonID uuid.UUID) (*progressdomain.CourseProgress, error) {
	return u.repo.CompleteLesson(ctx, userID, lessonID)
}

func (u *UseCase) GetCourseProgress(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (*progressdomain.CourseProgress, error) {
	return u.repo.GetCourseProgress(ctx, userID, courseID)
}

func (u *UseCase) ListCourseProgress(ctx context.Context, userID uuid.UUID) ([]progressdomain.CourseProgress, error) {
	return u.repo.ListCourseProgress(ctx, userID)
}
