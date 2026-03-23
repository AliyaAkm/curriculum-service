package topic

import (
	"context"
	"curriculum-service/internal/domain/topic"
)

type Repository interface {
	GetAllTopics(ctx context.Context) ([]topic.Topic, error)
}
