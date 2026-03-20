package handlers

import (
	"curriculum-service/internal/domain"
	"curriculum-service/internal/http/dto"
)

func toSearchCoursesResponse(result domain.CourseSearchResult) dto.SearchCoursesResponse {
	return dto.SearchCoursesResponse{
		Items: toCourseResponses(result.Items),
		Pagination: dto.PaginationResponse{
			Page:       result.Pagination.Page,
			PageSize:   result.Pagination.PageSize,
			TotalItems: result.Pagination.TotalItems,
			TotalPages: result.Pagination.TotalPages,
		},
		AppliedFilters: dto.AppliedFiltersResponse{
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

func toCourseResponses(items []domain.CourseCard) []dto.CourseResponse {
	result := make([]dto.CourseResponse, 0, len(items))
	for _, item := range items {
		result = append(result, dto.CourseResponse{
			ID:               item.ID.String(),
			Slug:             item.Slug,
			Status:           string(item.Status),
			Level:            string(item.Level),
			DurationCategory: string(item.DurationCategory),
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

func toFilterOptionsResponse(options domain.FilterOptions) dto.FilterOptionsResponse {
	return dto.FilterOptionsResponse{
		Topics:               toTopicFilterResponses(options.Topics),
		Levels:               toFilterValueResponses(options.Levels),
		Durations:            toFilterValueResponses(options.Durations),
		Ratings:              toRatingResponses(options.Ratings),
		CertificateAvailable: options.CertificateAvailable,
	}
}

func toTopicFilterResponses(items []domain.TopicFilterOption) []dto.TopicFilterOptionResponse {
	result := make([]dto.TopicFilterOptionResponse, 0, len(items))
	for _, item := range items {
		result = append(result, dto.TopicFilterOptionResponse{
			Slug:         item.Slug,
			Name:         item.Name,
			CoursesCount: item.CoursesCount,
		})
	}
	return result
}

func toFilterValueResponses(items []domain.FilterValueOption) []dto.FilterValueOptionResponse {
	result := make([]dto.FilterValueOptionResponse, 0, len(items))
	for _, item := range items {
		result = append(result, dto.FilterValueOptionResponse{
			Value:        item.Value,
			Label:        item.Label,
			CoursesCount: item.CoursesCount,
		})
	}
	return result
}

func toRatingResponses(items []domain.RatingOption) []dto.RatingOptionResponse {
	result := make([]dto.RatingOptionResponse, 0, len(items))
	for _, item := range items {
		result = append(result, dto.RatingOptionResponse{
			Value: item.Value,
		})
	}
	return result
}

func levelsToStrings(levels []domain.CourseLevel) []string {
	result := make([]string, 0, len(levels))
	for _, level := range levels {
		result = append(result, string(level))
	}
	return result
}

func durationsToStrings(durations []domain.DurationCategory) []string {
	result := make([]string, 0, len(durations))
	for _, duration := range durations {
		result = append(result, string(duration))
	}
	return result
}
