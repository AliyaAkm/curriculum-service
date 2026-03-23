package usecase

import (
	"context"
	"curriculum-service/internal/domain/category"
)

type CatalogRepository interface {
	SearchCourses(ctx context.Context, filter category.CourseSearchFilter) ([]category.CourseCard, int, error)
	GetFilterOptions(ctx context.Context, locale category.Locale) (category.FilterOptions, error)
}
