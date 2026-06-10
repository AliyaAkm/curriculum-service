package teacherstats

import (
	"context"
	teacherstatsdomain "curriculum-service/internal/domain/teacherstats"
)

const (
	defaultPeriodDays = 30
	maxPeriodDays     = 365
)

func (u *UseCase) GetStatistics(ctx context.Context, filter teacherstatsdomain.Filter) (*teacherstatsdomain.Statistics, error) {
	if filter.PeriodDays <= 0 {
		filter.PeriodDays = defaultPeriodDays
	}
	if filter.PeriodDays > maxPeriodDays {
		filter.PeriodDays = maxPeriodDays
	}

	return u.repo.GetStatistics(ctx, filter)
}
