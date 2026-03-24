package course

import (
	"curriculum-service/internal/domain/durationcategory"
	"curriculum-service/internal/domain/level"
	"curriculum-service/internal/domain/status"
	"curriculum-service/internal/domain/tag"
	"curriculum-service/internal/domain/topic"
	"github.com/google/uuid"
	"time"
)

type Course struct {
	ID                 uuid.UUID `gorm:"column:id;primary_key"`
	ExpectedHours      int       `gorm:"column:expected_hours"`
	Rating             float64   `gorm:"column:rating"`
	RatingCount        int       `gorm:"column:rating_count"`
	StudentsCount      int       `gorm:"column:students_count"`
	LessonsCount       int       `gorm:"column:lessons_count"`
	HasCertificate     bool      `gorm:"column:has_certificate"`
	CoverImageUrl      string    `gorm:"column:cover_image_url"`
	PublishedAt        time.Time `gorm:"column:published_at"`
	CreatedAt          time.Time `gorm:"column:created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at"`
	StatusID           uuid.UUID `gorm:"column:status_id"`
	LevelID            uuid.UUID `gorm:"column:level_id"`
	DurationCategoryID uuid.UUID `gorm:"column:duration_category_id"`
	AuthorID           uuid.UUID `gorm:"column:author_id"`
	TopicID            uuid.UUID `gorm:"column:topic_id"`

	Status           status.Status
	DurationCategory durationcategory.DurationCategory
	Level            level.Level
	Author           User
	Tags             []tag.Tag `gorm:"many2many:course_course_tags;"`
	Topic            topic.Topic
}

type User struct {
	ID           uuid.UUID `gorm:"column:id;primary_key"`
	Email        string
	PasswordHash string
	Roles        []Role `gorm:"many2many:user_roles;"`
	IsActive     bool
	CreatedAt    time.Time
}

type Role struct {
	ID           uuid.UUID `gorm:"column:id;primary_key"`
	Code         string
	Name         string
	Description  string
	IsDefault    bool
	IsPrivileged bool
	IsSupport    bool
	CreatedAt    time.Time
}
type CourseTag struct {
	CourseID uuid.UUID `gorm:"column:course_id;primaryKey"`
	TagID    uuid.UUID `gorm:"column:tag_id;primaryKey"`
}

func (Course) TableName() string {
	return "courses"
}
func (CourseTag) TableName() string {
	return "course_course_tags"
}
