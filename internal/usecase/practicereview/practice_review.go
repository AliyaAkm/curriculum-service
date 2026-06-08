package practicereview

import (
	"context"
	"curriculum-service/internal/domain"
	practicereviewdomain "curriculum-service/internal/domain/practicereview"
	"strings"

	"github.com/google/uuid"
)

func (u *UseCase) CreateSubmission(ctx context.Context, req practicereviewdomain.CreateSubmissionRequest) (*practicereviewdomain.Submission, error) {
	if req.PracticeID == uuid.Nil || req.StudentID == uuid.Nil || strings.TrimSpace(req.Code) == "" {
		return nil, domain.ErrValidation
	}
	return u.repo.CreateSubmission(ctx, req)
}

func (u *UseCase) ListStudentSubmissions(ctx context.Context, filter practicereviewdomain.StudentListFilter) ([]practicereviewdomain.Submission, error) {
	if filter.StudentID == uuid.Nil {
		return nil, domain.ErrValidation
	}
	if filter.Status != "" && !isSubmissionStatus(filter.Status) {
		return nil, domain.ErrInvalidPracticeReviewStatus
	}
	return u.repo.ListStudentSubmissions(ctx, filter)
}

func (u *UseCase) GetStudentSubmission(ctx context.Context, studentID uuid.UUID, submissionID uuid.UUID) (*practicereviewdomain.Submission, error) {
	if studentID == uuid.Nil || submissionID == uuid.Nil {
		return nil, domain.ErrValidation
	}
	return u.repo.GetStudentSubmission(ctx, studentID, submissionID)
}

func (u *UseCase) ListTeacherSubmissions(ctx context.Context, filter practicereviewdomain.TeacherListFilter) ([]practicereviewdomain.Submission, error) {
	if filter.TeacherID == uuid.Nil {
		return nil, domain.ErrValidation
	}
	if filter.Status != "" && !isSubmissionStatus(filter.Status) {
		return nil, domain.ErrInvalidPracticeReviewStatus
	}
	if filter.Limit <= 0 || filter.Limit > 100 {
		filter.Limit = 50
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	return u.repo.ListTeacherSubmissions(ctx, filter)
}

func (u *UseCase) GetTeacherSubmission(ctx context.Context, teacherID uuid.UUID, isAdmin bool, submissionID uuid.UUID) (*practicereviewdomain.Submission, error) {
	if teacherID == uuid.Nil || submissionID == uuid.Nil {
		return nil, domain.ErrValidation
	}
	return u.repo.GetTeacherSubmission(ctx, teacherID, isAdmin, submissionID)
}

func (u *UseCase) ReviewSubmission(ctx context.Context, req practicereviewdomain.ReviewSubmissionRequest) (*practicereviewdomain.Submission, error) {
	if req.SubmissionID == uuid.Nil || req.TeacherID == uuid.Nil {
		return nil, domain.ErrValidation
	}
	if req.Status != practicereviewdomain.SubmissionStatusApproved && req.Status != practicereviewdomain.SubmissionStatusChangesRequested {
		return nil, domain.ErrInvalidPracticeReviewStatus
	}
	return u.repo.ReviewSubmission(ctx, req)
}

func isSubmissionStatus(status string) bool {
	switch status {
	case practicereviewdomain.SubmissionStatusSubmitted,
		practicereviewdomain.SubmissionStatusInReview,
		practicereviewdomain.SubmissionStatusChangesRequested,
		practicereviewdomain.SubmissionStatusApproved:
		return true
	default:
		return false
	}
}
