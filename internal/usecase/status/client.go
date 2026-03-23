package status

import (
	"context"
	"curriculum-service/internal/domain/status"
)

type Repository interface {
	GetAllStatus(ctx context.Context) ([]status.Status, error)
}
