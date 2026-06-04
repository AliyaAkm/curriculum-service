package progress

import (
	"context"
	progressdomain "curriculum-service/internal/domain/progress"

	"github.com/google/uuid"
)

func (u *UseCase) CompleteLesson(ctx context.Context, userID uuid.UUID, lessonID uuid.UUID) (*progressdomain.CourseProgress, error) {
	result, err := u.repo.CompleteLesson(ctx, userID, lessonID)
	if err != nil {
		return nil, err
	}

	if result != nil && result.NewlyCompleted && u.notification != nil {
		if data, err := u.repo.GetLessonNotificationData(ctx, lessonID); err == nil && data != nil {
			_ = u.notification.SendEvent(ctx, userID, "lesson_completed", map[string]any{
				"lessonTitle": data.LessonTitle,
				"lessonId":    data.LessonID.String(),
				"courseId":    data.CourseID.String(),
				"courseTitle": data.CourseTitle,
			})
		}
	}

	return result, nil
}

func (u *UseCase) GetCourseProgress(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (*progressdomain.CourseProgress, error) {
	return u.repo.GetCourseProgress(ctx, userID, courseID)
}

func (u *UseCase) ListCourseProgress(ctx context.Context, userID uuid.UUID) ([]progressdomain.CourseProgress, error) {
	return u.repo.ListCourseProgress(ctx, userID)
}
