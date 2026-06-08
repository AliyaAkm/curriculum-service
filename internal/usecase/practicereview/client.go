package practicereview

import (
	"context"
	practicereviewdomain "curriculum-service/internal/domain/practicereview"

	"github.com/google/uuid"
)

type Repository interface {
	CreateSubmission(ctx context.Context, req practicereviewdomain.CreateSubmissionRequest) (*practicereviewdomain.Submission, error)
	ListStudentSubmissions(ctx context.Context, filter practicereviewdomain.StudentListFilter) ([]practicereviewdomain.Submission, error)
	GetStudentSubmission(ctx context.Context, studentID uuid.UUID, submissionID uuid.UUID) (*practicereviewdomain.Submission, error)
	ListTeacherSubmissions(ctx context.Context, filter practicereviewdomain.TeacherListFilter) ([]practicereviewdomain.Submission, error)
	GetTeacherSubmission(ctx context.Context, teacherID uuid.UUID, isAdmin bool, submissionID uuid.UUID) (*practicereviewdomain.Submission, error)
	ReviewSubmission(ctx context.Context, req practicereviewdomain.ReviewSubmissionRequest) (*practicereviewdomain.Submission, error)
}
