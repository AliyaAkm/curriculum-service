package level

import (
	"context"
	"curriculum-service/internal/domain/level"
)

func (u *UseCase) GetAllLevels(ctx context.Context) ([]level.Level, error) {
	return u.repo.GetAllLevel(ctx)
}
