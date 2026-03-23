package durationcategory

import (
	"context"
	"curriculum-service/internal/domain/durationcategory"
)

func (r *Repo) GetAllDurationCategories(ctx context.Context) ([]durationcategory.DurationCategory, error) {
	var durationCategories []durationcategory.DurationCategory
	err := r.db.WithContext(ctx).Order("name ASC").Find(&durationCategories).Error
	if err != nil {
		return nil, err
	}
	return durationCategories, nil
}
