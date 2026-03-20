package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Locale string

const (
	LocaleRU Locale = "ru"
	LocaleEN Locale = "en"
	LocaleKZ Locale = "kz"
)

type CourseStatus string

const (
	StatusDraft     CourseStatus = "draft"
	StatusPublished CourseStatus = "published"
	StatusArchived  CourseStatus = "archived"
)

type CourseLevel string

const (
	LevelBeginner     CourseLevel = "beginner"
	LevelIntermediate CourseLevel = "intermediate"
	LevelAdvanced     CourseLevel = "advanced"
)

type DurationCategory string

const (
	DurationQuick   DurationCategory = "quick"
	DurationFocused DurationCategory = "focused"
	DurationDeep    DurationCategory = "deep"
)

type CourseSearchFilter struct {
	Query           string
	Locale          Locale
	TopicSlugs      []string
	Levels          []CourseLevel
	MinRating       float64
	Durations       []DurationCategory
	WithCertificate *bool
	Page            int
	PageSize        int
}

type CourseCard struct {
	ID               uuid.UUID
	Slug             string
	Status           CourseStatus
	Level            CourseLevel
	DurationCategory DurationCategory
	ExpectedHours    int
	Rating           float64
	RatingCount      int
	StudentsCount    int
	LessonsCount     int
	HasCertificate   bool
	CoverImageURL    string
	AuthorName       string
	Title            string
	Subtitle         string
	ShortDescription string
	TopicSlugs       []string
	TopicNames       []string
	Tags             []string
	PublishedAt      time.Time
}

type Pagination struct {
	Page       int
	PageSize   int
	TotalItems int
	TotalPages int
}

type AppliedFilters struct {
	Query           string
	Locale          Locale
	TopicSlugs      []string
	Levels          []CourseLevel
	MinRating       float64
	Durations       []DurationCategory
	WithCertificate bool
}

type CourseSearchResult struct {
	Items          []CourseCard
	Pagination     Pagination
	AppliedFilters AppliedFilters
}

type TopicFilterOption struct {
	Slug         string
	Name         string
	CoursesCount int
}

type FilterValueOption struct {
	Value        string
	Label        string
	CoursesCount int
}

type RatingOption struct {
	Value float64
}

type FilterOptions struct {
	Topics               []TopicFilterOption
	Levels               []FilterValueOption
	Durations            []FilterValueOption
	Ratings              []RatingOption
	CertificateAvailable bool
}

func NormalizeLocale(raw string) Locale {
	switch Locale(strings.ToLower(strings.TrimSpace(raw))) {
	case LocaleEN:
		return LocaleEN
	case LocaleKZ:
		return LocaleKZ
	default:
		return LocaleRU
	}
}

func (f *CourseSearchFilter) Normalize() error {
	f.Query = strings.TrimSpace(f.Query)
	f.Locale = NormalizeLocale(string(f.Locale))
	f.TopicSlugs = normalizeStringList(f.TopicSlugs)

	if f.MinRating < 0 || f.MinRating > 5 {
		return ErrValidation
	}

	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = 12
	}
	if f.PageSize > 50 {
		f.PageSize = 50
	}

	for _, level := range f.Levels {
		switch level {
		case LevelBeginner, LevelIntermediate, LevelAdvanced:
		default:
			return ErrValidation
		}
	}

	for _, duration := range f.Durations {
		switch duration {
		case DurationQuick, DurationFocused, DurationDeep:
		default:
			return ErrValidation
		}
	}

	return nil
}

func ParseCourseLevels(raw []string) ([]CourseLevel, error) {
	items := normalizeStringList(raw)
	if len(items) == 0 {
		return nil, nil
	}

	levels := make([]CourseLevel, 0, len(items))
	for _, item := range items {
		level := CourseLevel(item)
		switch level {
		case LevelBeginner, LevelIntermediate, LevelAdvanced:
			levels = append(levels, level)
		default:
			return nil, ErrValidation
		}
	}

	return levels, nil
}

func ParseDurationCategories(raw []string) ([]DurationCategory, error) {
	items := normalizeStringList(raw)
	if len(items) == 0 {
		return nil, nil
	}

	durations := make([]DurationCategory, 0, len(items))
	for _, item := range items {
		duration := DurationCategory(item)
		switch duration {
		case DurationQuick, DurationFocused, DurationDeep:
			durations = append(durations, duration)
		default:
			return nil, ErrValidation
		}
	}

	return durations, nil
}

func LevelLabel(locale Locale, level CourseLevel) string {
	switch locale {
	case LocaleEN:
		switch level {
		case LevelBeginner:
			return "Beginner"
		case LevelIntermediate:
			return "Intermediate"
		case LevelAdvanced:
			return "Advanced"
		}
	case LocaleKZ:
		switch level {
		case LevelBeginner:
			return "Bastapqy"
		case LevelIntermediate:
			return "Orta"
		case LevelAdvanced:
			return "Zhogary"
		}
	default:
		switch level {
		case LevelBeginner:
			return "Начинающий"
		case LevelIntermediate:
			return "Средний"
		case LevelAdvanced:
			return "Продвинутый"
		}
	}

	return string(level)
}

func DurationLabel(locale Locale, duration DurationCategory) string {
	switch locale {
	case LocaleEN:
		switch duration {
		case DurationQuick:
			return "Quick"
		case DurationFocused:
			return "Focused"
		case DurationDeep:
			return "Deep"
		}
	case LocaleKZ:
		switch duration {
		case DurationQuick:
			return "Zhyldam"
		case DurationFocused:
			return "Naqty"
		case DurationDeep:
			return "Teren"
		}
	default:
		switch duration {
		case DurationQuick:
			return "Быстрый"
		case DurationFocused:
			return "Сфокусированный"
		case DurationDeep:
			return "Глубокий"
		}
	}

	return string(duration)
}

func normalizeStringList(items []string) []string {
	if len(items) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(items))
	result := make([]string, 0, len(items))
	for _, item := range items {
		normalized := strings.ToLower(strings.TrimSpace(item))
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		result = append(result, normalized)
	}

	if len(result) == 0 {
		return nil
	}

	return result
}
