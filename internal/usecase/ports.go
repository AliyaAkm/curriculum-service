package usecase

import (
	"context"

	"curriculum-service/internal/domain"
)

type CatalogRepository interface {
	SearchCourses(ctx context.Context, filter domain.CourseSearchFilter) ([]domain.CourseCard, int, error)
	GetFilterOptions(ctx context.Context, locale domain.Locale) (domain.FilterOptions, error)
}
