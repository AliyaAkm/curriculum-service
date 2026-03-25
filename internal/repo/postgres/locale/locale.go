package locale

import (
	"context"
	"curriculum-service/internal/domain/locale"
)

func (r *Repo) GetAllLocales(ctx context.Context) ([]locale.Locale, error) {
	var locales []locale.Locale
	err := r.db.WithContext(ctx).Find(&locales).Error
	if err != nil {
		return nil, err
	}
	return locales, nil
}
