package level

import (
	"context"
	"curriculum-service/internal/domain/level"
)

func (r *Repo) GetAllLevel(ctx context.Context) ([]level.Level, error) {
	var levels []level.Level
	err := r.db.WithContext(ctx).Find(&levels).Error
	if err != nil {
		return nil, err
	}
	return levels, nil
}
