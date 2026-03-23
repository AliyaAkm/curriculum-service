package usecase

import (
	"context"
	"curriculum-service/internal/domain/category"
)

type Catalog struct {
	repo CatalogRepository
}

func NewCatalog(repo CatalogRepository) *Catalog {
	return &Catalog{repo: repo}
}

func (uc *Catalog) SearchCourses(ctx context.Context, filter category.CourseSearchFilter) (category.CourseSearchResult, error) {

	items, total, err := uc.repo.SearchCourses(ctx, filter)
	if err != nil {
		return category.CourseSearchResult{}, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + filter.PageSize - 1) / filter.PageSize
	}

	result := category.CourseSearchResult{
		Items: items,
		Pagination: category.Pagination{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
		AppliedFilters: category.AppliedFilters{
			Query:      filter.Query,
			Locale:     filter.Locale,
			TopicSlugs: filter.TopicSlugs,
			Levels:     filter.Levels,
			MinRating:  filter.MinRating,
			Durations:  filter.Durations,
		},
	}

	if filter.WithCertificate != nil {
		result.AppliedFilters.WithCertificate = *filter.WithCertificate
	}

	return result, nil
}
