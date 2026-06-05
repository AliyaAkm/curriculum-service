package achievement

import (
	"context"
	achievementdomain "curriculum-service/internal/domain/achievement"

	"github.com/google/uuid"
)

func (u *UseCase) ListAchievements(ctx context.Context, userID uuid.UUID) ([]achievementdomain.Achievement, error) {
	return u.repo.ListAchievements(ctx, userID)
}

func (u *UseCase) SyncUnlockedAchievements(ctx context.Context, userID uuid.UUID) ([]achievementdomain.Achievement, error) {
	return u.repo.SyncUnlockedAchievements(ctx, userID)
}
