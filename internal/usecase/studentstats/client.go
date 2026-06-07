package studentstats

import (
	"context"
	studentstatsdomain "curriculum-service/internal/domain/studentstats"

	"github.com/google/uuid"
)

type Repository interface {
	GetStatistics(ctx context.Context, userID uuid.UUID) (*studentstatsdomain.Statistics, error)
}

type AIAnalyticsClient interface {
	GetUserAnalytics(ctx context.Context, userID uuid.UUID) (*studentstatsdomain.AIStats, []studentstatsdomain.ActivityDay, error)
}
