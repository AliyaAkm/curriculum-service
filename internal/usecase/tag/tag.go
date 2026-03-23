package tag

import (
	"context"
	"curriculum-service/internal/domain/tag"
)

func (u UseCase) GetAllTags(ctx context.Context) ([]tag.Tag, error) {
	return u.repo.GetAllTags(ctx)
}
