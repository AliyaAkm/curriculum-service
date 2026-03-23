package durationcategory

import (
	"context"
	"curriculum-service/internal/domain/durationcategory"
)

type client interface {
	GetAllDurationCategories(ctx context.Context) ([]durationcategory.DurationCategory, error)
}
