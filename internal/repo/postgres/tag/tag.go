package tag

import (
	"context"
	"curriculum-service/internal/domain/tag"
)

func (r Repo) GetAllTags(ctx context.Context) ([]tag.Tag, error) {
	var tags []tag.Tag
	err := r.db.WithContext(ctx).Order("name ASC").Find(&tags).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}
