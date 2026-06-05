package progress

import (
	"context"
	achievementdomain "curriculum-service/internal/domain/achievement"
	progressdomain "curriculum-service/internal/domain/progress"

	"github.com/google/uuid"
)

type Repository interface {
	CompleteLesson(ctx context.Context, userID uuid.UUID, lessonID uuid.UUID) (*progressdomain.CourseProgress, error)
	GetLessonNotificationData(ctx context.Context, lessonID uuid.UUID) (*progressdomain.LessonNotificationData, error)
	SyncUserLevel(ctx context.Context, userID uuid.UUID) (*int, error)
	GetCourseProgress(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (*progressdomain.CourseProgress, error)
	ListCourseProgress(ctx context.Context, userID uuid.UUID) ([]progressdomain.CourseProgress, error)
}

type NotificationSender interface {
	SendEvent(ctx context.Context, userID uuid.UUID, event string, data map[string]any) error
}

type AchievementSyncer interface {
	SyncUnlockedAchievements(ctx context.Context, userID uuid.UUID) ([]achievementdomain.Achievement, error)
}
