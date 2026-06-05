package progress

import (
	"context"
	achievementdomain "curriculum-service/internal/domain/achievement"
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
			if result.CompletedAt != nil && result.TotalLessons > 0 && result.CompletedLessons >= result.TotalLessons {
				_ = u.notification.SendEvent(ctx, userID, "course_completed", map[string]any{
					"courseId":    data.CourseID.String(),
					"courseTitle": data.CourseTitle,
				})
			}
			if levelReached, err := u.repo.SyncUserLevel(ctx, userID); err == nil && levelReached != nil {
				_ = u.notification.SendEvent(ctx, userID, "level_reached", map[string]any{
					"level": *levelReached,
				})
			}
			u.notifyUnlockedAchievements(ctx, userID)
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

func (u *UseCase) notifyUnlockedAchievements(ctx context.Context, userID uuid.UUID) {
	if u.notification == nil || u.achievements == nil {
		return
	}

	items, err := u.achievements.SyncUnlockedAchievements(ctx, userID)
	if err != nil {
		return
	}
	for _, item := range items {
		_ = u.notification.SendEvent(ctx, userID, "achievement_unlocked", map[string]any{
			"achievementId":    item.ID.String(),
			"achievementCode":  item.Code,
			"achievementTitle": achievementTitle(item),
		})
	}
}

func achievementTitle(value achievementdomain.Achievement) string {
	switch {
	case value.Title.RU != "":
		return value.Title.RU
	case value.Title.EN != "":
		return value.Title.EN
	case value.Title.KK != "":
		return value.Title.KK
	default:
		return "новое достижение"
	}
}
