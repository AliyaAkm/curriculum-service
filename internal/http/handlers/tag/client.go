package tag

import (
	"context"
	"curriculum-service/internal/domain/tag"
)

type client interface {
	GetAllTags(ctx context.Context) ([]tag.Tag, error)
}
