package certificate

import (
	"context"
	certificatedomain "curriculum-service/internal/domain/certificate"
	"io"

	"github.com/google/uuid"
)

type Repository interface {
	FindByUserAndCourse(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (*certificatedomain.Certificate, error)
	GetCourseCompletion(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (certificatedomain.CourseCompletion, error)
	Create(ctx context.Context, value *certificatedomain.Certificate) error
}

type Storage interface {
	PutObject(ctx context.Context, objectKey string, body io.ReadSeeker, size int64, contentType string) error
	PresignGetObject(objectKey string) (string, error)
}
