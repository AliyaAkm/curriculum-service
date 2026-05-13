package streak

import (
	"context"
	"curriculum-service/internal/domain/streak"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *Repo) GetStreak(ctx context.Context, userID uuid.UUID) (*streak.DailyStreak, error) {
	var entity *streak.DailyStreak

	err := r.db.WithContext(ctx).First(&entity, "user_id = ?", userID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *Repo) CreateStreak(ctx context.Context, value *streak.DailyStreak) error {
	err := r.db.WithContext(ctx).Create(value).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *Repo) UpdateStreak(ctx context.Context, value *streak.DailyStreak) error {
	return r.db.WithContext(ctx).Save(value).Error
}
