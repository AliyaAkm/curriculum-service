package certificate

import (
	"github.com/google/uuid"
	"time"
)

type Certificate struct {
	ID                uuid.UUID `gorm:"column:id"`
	UserID            uuid.UUID `gorm:"column:user_id"`
	CourseID          uuid.UUID `gorm:"column:course_id"`
	CertificateNumber string    `gorm:"column:certificate_number"`
	PDFObjectKey      string    `gorm:"column:pdf_object_key"`
	IssuedAt          time.Time `gorm:"column:issued_at"`
	CompletedAt       time.Time `gorm:"column:completed_at"`
}

type IssueData struct {
	UserName          string
	UserEmail         string
	CourseTitle       string
	CertificateNumber string
	IssuedAt          time.Time
	CompletedAt       time.Time
}

type CourseCompletion struct {
	CourseID         uuid.UUID
	CourseTitle      string
	UserLogin        string
	HasCertificate   bool
	CompletedAt      *time.Time
	TotalLessons     int
	CompletedLessons int
}

type IssueResult struct {
	Certificate Certificate
	DownloadURL string
}

func (Certificate) TableName() string {
	return "course_certificates"
}
