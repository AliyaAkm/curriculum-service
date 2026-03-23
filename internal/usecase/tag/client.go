package tag

import (
	"context"
	"curriculum-service/internal/domain/tag"
)

type Repository interface {
	GetAllTags(ctx context.Context) ([]tag.Tag, error)
}
