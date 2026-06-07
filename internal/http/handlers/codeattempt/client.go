package codeattempt

import (
	"context"
	codeattemptdomain "curriculum-service/internal/domain/codeattempt"
)

type client interface {
	Run(ctx context.Context, req codeattemptdomain.RunRequest) (*codeattemptdomain.RunResult, error)
}
