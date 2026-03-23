package topic

import (
	"context"
	"curriculum-service/internal/domain/topic"
)

type client interface {
	GetAllTopics(ctx context.Context) ([]topic.Topic, error)
}
