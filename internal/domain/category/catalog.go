package category

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

type CourseLevel string

type DurationCategory string


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
