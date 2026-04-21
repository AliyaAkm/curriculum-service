package course

import (
	"curriculum-service/internal/http/dto/durationcategory"
	"curriculum-service/internal/http/dto/level"
	"curriculum-service/internal/http/dto/status"
	"curriculum-service/internal/http/dto/tag"
	"curriculum-service/internal/http/dto/topic"
	"github.com/google/uuid"
	"time"
)

type Courses struct {
	ID               uuid.UUID                         `json:"id"`
	Title            string                            `json:"title"`
	SubTitle         string                            `json:"sub_title"`
	Description      string                            `json:"description"`
	ExpectedHours    int                               `json:"expected_hours"`
	Rating           float64                           `json:"rating"`
	RatingCount      int                               `json:"rating_count"`
	StudentsCount    int                               `json:"students_count"`
	LessonsCount     int                               `json:"lessons_count"`
	HasCertificate   bool                              `json:"has_certificate"`
	CoverImageUrl    string                            `json:"cover_image_url"`
	PublishedAt      time.Time                         `json:"published_at"`
	CreatedAt        time.Time                         `json:"created_at"`
	UpdatedAt        time.Time                         `json:"updated_at"`
	Status           status.Status                     `json:"status"`
	Level            level.Level                       `json:"level"`
	DurationCategory durationcategory.DurationCategory `json:"duration_category"`
	Author           User                              `json:"author"`
	Tags             []tag.Tag                         `json:"tags"`
	Topic            topic.Topic                       `json:"topic"`
	LearningOutcomes []string                          `json:"learning_outcome"`
}

type Subscription struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	CourseID uuid.UUID `json:"course_id"`
}

type SubscriptionRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	CourseID uuid.UUID `json:"course_id"`
}

type CourseRequest struct {
	Title              string      `json:"title"`
	SubTitle           string      `json:"sub_title"`
	Description        string      `json:"description"`
	ExpectedHours      int         `json:"expected_hours"`
	StudentsCount      int         `json:"students_count"`
	LessonsCount       int         `json:"lessons_count"`
	HasCertificate     bool        `json:"has_certificate"`
	CoverImageUrl      string      `json:"cover_image_url"`
	LearningOutcomes   []string    `json:"learning_outcomes"`
	LevelID            uuid.UUID   `json:"level_id"`
	StatusID           uuid.UUID   `json:"status_id"`
	DurationCategoryID uuid.UUID   `json:"duration_category_id"`
	AuthorID           uuid.UUID   `json:"author_id"`
	TagIDs             []uuid.UUID `json:"tag_ids"`
	TopicID            uuid.UUID   `json:"topic_id"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Roles     []Role    `json:"roles"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type Role struct {
	ID           uuid.UUID `json:"id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	IsDefault    bool      `json:"is_default"`
	IsPrivileged bool      `json:"is_privileged"`
	IsSupport    bool      `json:"is_support"`
	CreatedAt    time.Time `json:"created_at"`
}
