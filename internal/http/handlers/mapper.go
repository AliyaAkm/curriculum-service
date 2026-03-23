package handlers

import (
	"curriculum-service/internal/domain/category"
	"curriculum-service/internal/http/dto/catalog"
)

func toSearchCoursesResponse(result category.CourseSearchResult) catalog.SearchCoursesResponse {
	return catalog.SearchCoursesResponse{
		Items: toCourseResponses(result.Items),
		Pagination: catalog.PaginationResponse{
			Page:       result.Pagination.Page,
			PageSize:   result.Pagination.PageSize,
			TotalItems: result.Pagination.TotalItems,
			TotalPages: result.Pagination.TotalPages,
		},
		AppliedFilters: catalog.AppliedFiltersResponse{
			Query:           result.AppliedFilters.Query,
			Locale:          string(result.AppliedFilters.Locale),
			TopicSlugs:      result.AppliedFilters.TopicSlugs,
			Levels:          levelsToStrings(result.AppliedFilters.Levels),
			MinRating:       result.AppliedFilters.MinRating,
			Durations:       durationsToStrings(result.AppliedFilters.Durations),
			WithCertificate: result.AppliedFilters.WithCertificate,
		},
	}
}

func toCourseResponses(items []category.CourseCard) []catalog.CourseResponse {
	result := make([]catalog.CourseResponse, 0, len(items))
	for _, item := range items {
		result = append(result, catalog.CourseResponse{
			ID:   item.ID.String(),
			Slug: item.Slug,
			//Status:           string(item.Status),
			// Level:            string(item.Level),
			// DurationCategory: string(item.DurationCategory),
			ExpectedHours:    item.ExpectedHours,
			Rating:           item.Rating,
			RatingCount:      item.RatingCount,
			StudentsCount:    item.StudentsCount,
			LessonsCount:     item.LessonsCount,
			HasCertificate:   item.HasCertificate,
			CoverImageURL:    item.CoverImageURL,
			AuthorName:       item.AuthorName,
			Title:            item.Title,
			Subtitle:         item.Subtitle,
			ShortDescription: item.ShortDescription,
			TopicSlugs:       item.TopicSlugs,
			TopicNames:       item.TopicNames,
			Tags:             item.Tags,
			PublishedAt:      item.PublishedAt,
		})
	}
	return result
}

func toFilterOptionsResponse(options category.FilterOptions) catalog.FilterOptionsResponse {
	return catalog.FilterOptionsResponse{
		Topics:               toTopicFilterResponses(options.Topics),
		Levels:               toFilterValueResponses(options.Levels),
		Durations:            toFilterValueResponses(options.Durations),
		Ratings:              toRatingResponses(options.Ratings),
		CertificateAvailable: options.CertificateAvailable,
	}
}

func toTopicFilterResponses(items []category.TopicFilterOption) []catalog.TopicFilterOptionResponse {
	result := make([]catalog.TopicFilterOptionResponse, 0, len(items))
	for _, item := range items {
		result = append(result, catalog.TopicFilterOptionResponse{
			Slug:         item.Slug,
			Name:         item.Name,
			CoursesCount: item.CoursesCount,
		})
	}
	return result
}

func toFilterValueResponses(items []category.FilterValueOption) []catalog.FilterValueOptionResponse {
	result := make([]catalog.FilterValueOptionResponse, 0, len(items))
	for _, item := range items {
		result = append(result, catalog.FilterValueOptionResponse{
			Value:        item.Value,
			Label:        item.Label,
			CoursesCount: item.CoursesCount,
		})
	}
	return result
}

func toRatingResponses(items []category.RatingOption) []catalog.RatingOptionResponse {
	result := make([]catalog.RatingOptionResponse, 0, len(items))
	for _, item := range items {
		result = append(result, catalog.RatingOptionResponse{
			Value: item.Value,
		})
	}
	return result
}

func levelsToStrings(levels []category.CourseLevel) []string {
	result := make([]string, 0, len(levels))
	for _, level := range levels {
		result = append(result, string(level))
	}
	return result
}

func durationsToStrings(durations []category.DurationCategory) []string {
	result := make([]string, 0, len(durations))
	for _, duration := range durations {
		result = append(result, string(duration))
	}
	return result
}
