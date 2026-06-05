package achievement

import (
	"context"
	achievementdomain "curriculum-service/internal/domain/achievement"

	"github.com/google/uuid"
)

type client interface {
	ListAchievements(ctx context.Context, userID uuid.UUID) ([]achievementdomain.Achievement, error)
	SyncUnlockedAchievements(ctx context.Context, userID uuid.UUID) ([]achievementdomain.Achievement, error)
}
