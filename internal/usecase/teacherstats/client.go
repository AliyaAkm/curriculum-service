package teacherstats

import (
	"context"
	teacherstatsdomain "curriculum-service/internal/domain/teacherstats"
)

type Repository interface {
	GetStatistics(ctx context.Context, filter teacherstatsdomain.Filter) (*teacherstatsdomain.Statistics, error)
}
