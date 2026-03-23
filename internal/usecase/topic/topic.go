package topic

import (
	"context"
	"curriculum-service/internal/domain/topic"
)

func (u *UseCase) GetAllTopics(ctx context.Context) ([]topic.Topic, error) {
	return u.repo.GetAllTopics(ctx)
}
