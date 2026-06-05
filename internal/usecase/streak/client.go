package streak

import (
	"context"
	achievementdomain "curriculum-service/internal/domain/achievement"
	"curriculum-service/internal/domain/streak"
	"github.com/google/uuid"
)

type Repository interface {
	GetStreak(ctx context.Context, userID uuid.UUID) (*streak.DailyStreak, error)
	CreateStreak(ctx context.Context, value *streak.DailyStreak) error
	UpdateStreak(ctx context.Context, value *streak.DailyStreak) error
	UpdateUserMaxStreak(ctx context.Context, userID uuid.UUID, value int64) error
}

type NotificationSender interface {
	SendEvent(ctx context.Context, userID uuid.UUID, event string, data map[string]any) error
}

type AchievementSyncer interface {
	SyncUnlockedAchievements(ctx context.Context, userID uuid.UUID) ([]achievementdomain.Achievement, error)
}
