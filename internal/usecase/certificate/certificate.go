package certificate

import (
	"bytes"
	"context"
	"curriculum-service/internal/domain"
	certificatedomain "curriculum-service/internal/domain/certificate"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const certificateContentType = "application/pdf"

func (u *UseCase) IssueCourseCertificate(ctx context.Context, userID uuid.UUID, courseID uuid.UUID, userName string, isAdmin bool) (*certificatedomain.IssueResult, error) {
	existing, err := u.repo.FindByUserAndCourse(ctx, userID, courseID)
	if err == nil {
		return u.resultWithDownloadURL(existing)
	}
	if !errors.Is(err, domain.ErrCertificateNotFound) {
		return nil, err
	}

	completion, err := u.repo.GetCourseCompletion(ctx, userID, courseID)
	if err != nil {
		return nil, err
	}
	if !completion.HasCertificate {
		return nil, domain.ErrCertificateUnavailable
	}
	if !isAdmin && (completion.TotalLessons == 0 || completion.CompletedLessons < completion.TotalLessons || completion.CompletedAt == nil) {
		return nil, domain.ErrCourseNotCompleted
	}

	now := time.Now().UTC()
	completedAt := now
	if completion.CompletedAt != nil {
		completedAt = completion.CompletedAt.UTC()
	}

	certID := uuid.New()
	certificateNumber := newCertificateNumber(now, certID)
	value := &certificatedomain.Certificate{
		ID:                certID,
		UserID:            userID,
		CourseID:          courseID,
		CertificateNumber: certificateNumber,
		PDFObjectKey:      certificateObjectKey(courseID, userID, certID),
		IssuedAt:          now,
		CompletedAt:       completedAt,
	}

	pdf := renderCertificatePDF(certificatedomain.IssueData{
		UserName:          certificateUserName(userID, userName, completion.UserLogin),
		CourseTitle:       completion.CourseTitle,
		CertificateNumber: certificateNumber,
		IssuedAt:          value.IssuedAt,
		CompletedAt:       value.CompletedAt,
	})

	if u.storage == nil {
		return nil, errors.New("certificate storage is not configured")
	}
	reader := bytes.NewReader(pdf)
	if err = u.storage.PutObject(ctx, value.PDFObjectKey, reader, int64(len(pdf)), certificateContentType); err != nil {
		return nil, err
	}

	if err = u.repo.Create(ctx, value); err != nil {
		if existing, findErr := u.repo.FindByUserAndCourse(ctx, userID, courseID); findErr == nil {
			return u.resultWithDownloadURL(existing)
		}
		return nil, err
	}

	return u.resultWithDownloadURL(value)
}

func (u *UseCase) resultWithDownloadURL(value *certificatedomain.Certificate) (*certificatedomain.IssueResult, error) {
	if value == nil {
		return nil, domain.ErrCertificateNotFound
	}
	if u.storage == nil {
		return nil, errors.New("certificate storage is not configured")
	}

	downloadURL, err := u.storage.PresignGetObject(value.PDFObjectKey)
	if err != nil {
		return nil, err
	}

	return &certificatedomain.IssueResult{
		Certificate: *value,
		DownloadURL: downloadURL,
	}, nil
}

func newCertificateNumber(now time.Time, id uuid.UUID) string {
	return fmt.Sprintf("ZERDE-%d-%s", now.Year(), strings.ToUpper(strings.ReplaceAll(id.String()[:8], "-", "")))
}

func certificateObjectKey(courseID uuid.UUID, userID uuid.UUID, certificateID uuid.UUID) string {
	return "certificates/" + courseID.String() + "/" + userID.String() + "/" + certificateID.String() + ".pdf"
}

func certificateUserName(userID uuid.UUID, values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}

	return "Student " + userID.String()[:8]
}
