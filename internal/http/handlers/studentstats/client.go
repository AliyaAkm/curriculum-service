package studentstats

import (
	"context"
	studentstatsdomain "curriculum-service/internal/domain/studentstats"

	"github.com/google/uuid"
)

type client interface {
	GetStatistics(ctx context.Context, userID uuid.UUID) (*studentstatsdomain.Statistics, error)
}
