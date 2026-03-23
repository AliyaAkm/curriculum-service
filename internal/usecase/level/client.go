package level

import (
	"context"
	"curriculum-service/internal/domain/level"
)

type Repository interface {
	GetAllLevel(ctx context.Context) ([]level.Level, error)
}
