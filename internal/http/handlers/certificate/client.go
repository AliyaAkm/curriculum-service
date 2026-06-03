package certificate

import (
	"context"
	certificatedomain "curriculum-service/internal/domain/certificate"

	"github.com/google/uuid"
)

type client interface {
	IssueCourseCertificate(ctx context.Context, userID uuid.UUID, courseID uuid.UUID, userName string, isAdmin bool) (*certificatedomain.IssueResult, error)
}
