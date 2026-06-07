package codeattempt

import (
	"context"
	codeattemptdomain "curriculum-service/internal/domain/codeattempt"
	practicedomain "curriculum-service/internal/domain/practice"

	"github.com/google/uuid"
)

type Repository interface {
	CreateAttempt(ctx context.Context, value codeattemptdomain.Attempt) (*codeattemptdomain.Attempt, error)
}

type Runner interface {
	Run(ctx context.Context, language, code string) (RunnerResult, error)
}

type RunnerResult struct {
	Output string
	Error  string
	Passed bool
}

type PracticeProvider interface {
	GetByID(ctx context.Context, id uuid.UUID) (*practicedomain.Task, error)
}
