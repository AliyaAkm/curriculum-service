package level

import (
	"context"
	"curriculum-service/internal/domain/level"
)

type client interface {
	GetAllLevels(ctx context.Context) ([]level.Level, error)
}
