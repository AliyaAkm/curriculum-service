package dto

import "time"

type CourseResponse struct {
	ID               string    `json:"id"`
	Slug             string    `json:"slug"`
	Status           string    `json:"status"`
	Level            string    `json:"level"`
	DurationCategory string    `json:"duration_category"`
	ExpectedHours    int       `json:"expected_hours"`
	Rating           float64   `json:"rating"`
	RatingCount      int       `json:"rating_count"`
	StudentsCount    int       `json:"students_count"`
	LessonsCount     int       `json:"lessons_count"`
	HasCertificate   bool      `json:"has_certificate"`
	CoverImageURL    string    `json:"cover_image_url"`
	AuthorName       string    `json:"author_name"`
	Title            string    `json:"title"`
	Subtitle         string    `json:"subtitle"`
	ShortDescription string    `json:"short_description"`
	TopicSlugs       []string  `json:"topic_slugs"`
	TopicNames       []string  `json:"topic_names"`
	Tags             []string  `json:"tags"`
	PublishedAt      time.Time `json:"published_at"`
}

type PaginationResponse struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type AppliedFiltersResponse struct {
	Query           string   `json:"query"`
	Locale          string   `json:"locale"`
	TopicSlugs      []string `json:"topic_slugs"`
	Levels          []string `json:"levels"`
	MinRating       float64  `json:"min_rating"`
	Durations       []string `json:"durations"`
	WithCertificate bool     `json:"with_certificate"`
}

type SearchCoursesResponse struct {
	Items          []CourseResponse       `json:"items"`
	Pagination     PaginationResponse     `json:"pagination"`
	AppliedFilters AppliedFiltersResponse `json:"applied_filters"`
}

type TopicFilterOptionResponse struct {
	Slug         string `json:"slug"`
	Name         string `json:"name"`
	CoursesCount int    `json:"courses_count"`
}

type FilterValueOptionResponse struct {
	Value        string `json:"value"`
	Label        string `json:"label"`
	CoursesCount int    `json:"courses_count"`
}

type RatingOptionResponse struct {
	Value float64 `json:"value"`
}

type FilterOptionsResponse struct {
	Topics               []TopicFilterOptionResponse `json:"topics"`
	Levels               []FilterValueOptionResponse `json:"levels"`
	Durations            []FilterValueOptionResponse `json:"durations"`
	Ratings              []RatingOptionResponse      `json:"ratings"`
	CertificateAvailable bool                        `json:"certificate_available"`
}
