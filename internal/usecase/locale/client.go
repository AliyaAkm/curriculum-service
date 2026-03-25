package locale

import (
	"context"
	"curriculum-service/internal/domain/locale"
)

type Repository interface {
	GetAllLocales(ctx context.Context) ([]locale.Locale, error)
}
