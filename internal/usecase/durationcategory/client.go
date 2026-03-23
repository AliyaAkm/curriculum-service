package durationcategory

import (
	"context"
	"curriculum-service/internal/domain/durationcategory"
)

type Repository interface {
	GetAllDurationCategories(ctx context.Context) ([]durationcategory.DurationCategory, error)
}
