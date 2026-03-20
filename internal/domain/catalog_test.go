package domain

import (
	"slices"
	"testing"
)

func TestNormalizeLocaleFallsBackToRussian(t *testing.T) {
	if got := NormalizeLocale("de"); got != LocaleRU {
		t.Fatalf("expected %q, got %q", LocaleRU, got)
	}
}

func TestCourseSearchFilterNormalize(t *testing.T) {
	withCertificate := true
	filter := CourseSearchFilter{
		Query:           "  qa automation  ",
		Locale:          "EN",
		TopicSlugs:      []string{"Programming-Languages", "programming-languages", "  ai  "},
		Levels:          []CourseLevel{LevelBeginner, LevelIntermediate},
		Durations:       []DurationCategory{DurationQuick},
		MinRating:       4.5,
		WithCertificate: &withCertificate,
	}

	if err := filter.Normalize(); err != nil {
		t.Fatalf("normalize returned error: %v", err)
	}

	if filter.Query != "qa automation" {
		t.Fatalf("unexpected query: %q", filter.Query)
	}
	if filter.Locale != LocaleEN {
		t.Fatalf("unexpected locale: %q", filter.Locale)
	}
	if filter.Page != 1 {
		t.Fatalf("expected default page 1, got %d", filter.Page)
	}
	if filter.PageSize != 12 {
		t.Fatalf("expected default page size 12, got %d", filter.PageSize)
	}

	expectedTopics := []string{"programming-languages", "ai"}
	if !slices.Equal(filter.TopicSlugs, expectedTopics) {
		t.Fatalf("unexpected topics: %#v", filter.TopicSlugs)
	}
}

func TestParseCourseLevelsRejectsUnknownValues(t *testing.T) {
	if _, err := ParseCourseLevels([]string{"expert"}); err == nil {
		t.Fatal("expected validation error for unsupported level")
	}
}
