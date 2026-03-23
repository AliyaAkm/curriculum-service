package topic

import (
	"context"
	"curriculum-service/internal/domain/topic"
)

func (r *Repo) GetAllTopics(ctx context.Context) ([]topic.Topic, error) {
	var topics []topic.Topic
	err := r.db.WithContext(ctx).Order("name ASC").Find(&topics).Error
	if err != nil {
		return nil, err
	}
	return topics, nil
}
