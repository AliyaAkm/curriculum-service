package status

import (
	"context"
	"curriculum-service/internal/domain/status"
)

type client interface {
	GetAllStatuses(ctx context.Context) ([]status.Status, error)
}
