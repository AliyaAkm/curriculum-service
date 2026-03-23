package status

import (
	"context"
	"curriculum-service/internal/domain/status"
)

func (u *UseCase) GetAllStatuses(ctx context.Context) ([]status.Status, error) {
	return u.repo.GetAllStatus(ctx)
}
