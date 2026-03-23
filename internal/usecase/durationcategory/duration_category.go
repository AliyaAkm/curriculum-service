package durationcategory

import (
	"context"
	"curriculum-service/internal/domain/durationcategory"
)

func (u UseCase) GetAllDurationCategories(ctx context.Context) ([]durationcategory.DurationCategory, error) {
	return u.repo.GetAllDurationCategories(ctx)
}
