package usecase

import (
	"context"

	"curriculum-service/internal/domain"
)

type Catalog struct {
	repo CatalogRepository
}

func NewCatalog(repo CatalogRepository) *Catalog {
	return &Catalog{repo: repo}
}

func (uc *Catalog) SearchCourses(ctx context.Context, filter domain.CourseSearchFilter) (domain.CourseSearchResult, error) {
	if err := filter.Normalize(); err != nil {
		return domain.CourseSearchResult{}, err
	}

	items, total, err := uc.repo.SearchCourses(ctx, filter)
	if err != nil {
		return domain.CourseSearchResult{}, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + filter.PageSize - 1) / filter.PageSize
	}

	result := domain.CourseSearchResult{
		Items: items,
		Pagination: domain.Pagination{
			Page:       filter.Page,
			PageSize:   filter.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
		AppliedFilters: domain.AppliedFilters{
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

func (uc *Catalog) GetFilterOptions(ctx context.Context, locale domain.Locale) (domain.FilterOptions, error) {
	normalizedLocale := domain.NormalizeLocale(string(locale))
	options, err := uc.repo.GetFilterOptions(ctx, normalizedLocale)
	if err != nil {
		return domain.FilterOptions{}, err
	}

	options.Levels = buildLevelOptions(normalizedLocale, options.Levels)
	options.Durations = buildDurationOptions(normalizedLocale, options.Durations)
	options.Ratings = []domain.RatingOption{
		{Value: 3},
		{Value: 4},
		{Value: 4.5},
	}

	return options, nil
}

func buildLevelOptions(locale domain.Locale, existing []domain.FilterValueOption) []domain.FilterValueOption {
	counts := make(map[string]int, len(existing))
	for _, item := range existing {
		counts[item.Value] = item.CoursesCount
	}

	ordered := []domain.CourseLevel{
		domain.LevelBeginner,
		domain.LevelIntermediate,
		domain.LevelAdvanced,
	}

	result := make([]domain.FilterValueOption, 0, len(ordered))
	for _, level := range ordered {
		result = append(result, domain.FilterValueOption{
			Value:        string(level),
			Label:        domain.LevelLabel(locale, level),
			CoursesCount: counts[string(level)],
		})
	}

	return result
}

func buildDurationOptions(locale domain.Locale, existing []domain.FilterValueOption) []domain.FilterValueOption {
	counts := make(map[string]int, len(existing))
	for _, item := range existing {
		counts[item.Value] = item.CoursesCount
	}

	ordered := []domain.DurationCategory{
		domain.DurationQuick,
		domain.DurationFocused,
		domain.DurationDeep,
	}

	result := make([]domain.FilterValueOption, 0, len(ordered))
	for _, duration := range ordered {
		result = append(result, domain.FilterValueOption{
			Value:        string(duration),
			Label:        domain.DurationLabel(locale, duration),
			CoursesCount: counts[string(duration)],
		})
	}

	return result
}
