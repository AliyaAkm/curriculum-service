package streak

import (
	"context"
	"curriculum-service/internal/domain/streak"
	"time"

	"github.com/google/uuid"
)

func (u *UseCase) GetStreak(ctx context.Context, userID uuid.UUID) (*streak.DailyStreak, error) {
	entity, err := u.repo.GetStreak(ctx, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	if entity == nil {
		entity = &streak.DailyStreak{
			ID:        uuid.New(),
			UserID:    userID,
			Streak:    1,
			LastLogin: now,
		}
		if err := u.repo.CreateStreak(ctx, entity); err != nil {
			return nil, err
		}
		if err := u.repo.UpdateUserMaxStreak(ctx, userID, entity.Streak); err != nil {
			return nil, err
		}
		return entity, nil
	}

	days := int(now.Sub(entity.LastLogin).Hours() / 24)
	if days == 0 {
		if err := u.repo.UpdateUserMaxStreak(ctx, userID, entity.Streak); err != nil {
			return nil, err
		}
		return entity, nil
	}
	if days == 1 {
		entity.Streak++
	} else {
		entity.Streak = 1
	}
	entity.LastLogin = now

	if err := u.repo.UpdateStreak(ctx, entity); err != nil {
		return nil, err
	}
	if err := u.repo.UpdateUserMaxStreak(ctx, userID, entity.Streak); err != nil {
		return nil, err
	}
	return entity, nil
}
