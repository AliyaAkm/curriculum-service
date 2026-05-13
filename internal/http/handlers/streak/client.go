package streak

import (
	"context"
	"curriculum-service/internal/domain/streak"
	"github.com/google/uuid"
)

type client interface {
	GetStreak(ctx context.Context, userID uuid.UUID) (*streak.DailyStreak, error)
}
