package locale

import (
	"context"
	"curriculum-service/internal/domain/locale"
)

func (u *UseCase) GetAllLocales(ctx context.Context) ([]locale.Locale, error) {
	return u.repo.GetAllLocales(ctx)
}
