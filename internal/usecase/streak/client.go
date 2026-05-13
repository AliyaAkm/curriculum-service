package streak

import (
	"context"
	"curriculum-service/internal/domain/streak"
	"github.com/google/uuid"
)

type Repository interface {
	GetStreak(ctx context.Context, userID uuid.UUID) (*streak.DailyStreak, error)
	CreateStreak(ctx context.Context, value *streak.DailyStreak) error
	UpdateStreak(ctx context.Context, value *streak.DailyStreak) error
}
