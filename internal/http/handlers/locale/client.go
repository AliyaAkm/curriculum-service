package locale

import (
	"context"
	"curriculum-service/internal/domain/locale"
)

type client interface {
	GetAllLocales(ctx context.Context) ([]locale.Locale, error)
}
