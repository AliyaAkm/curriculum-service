package streak

import (
	"context"
	"curriculum-service/internal/domain/streak"
	"github.com/google/uuid"
	"time"
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
		err = u.repo.CreateStreak(ctx, entity)
		if err != nil {
			return nil, err
		}
		return entity, nil
	}
	days := int(now.Sub(entity.LastLogin).Hours() / 24)
	if days == 0 {
		return entity, nil // сегодня зашел
	}
	if days == 1 {
		entity.Streak++ // зашел вчера
	} else {
		entity.Streak = 1 // пропустил
	}
	entity.LastLogin = now

	err = u.repo.UpdateStreak(ctx, entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}
