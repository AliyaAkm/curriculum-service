package order

import (
	"context"
	"curriculum-service/internal/domain/order"
	"github.com/google/uuid"
)

func (r *Repo) CreateOrder(ctx context.Context, value *order.Order) error {
	err := r.db.WithContext(ctx).Create(value).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetOrderByID(ctx context.Context, id uuid.UUID) (*order.Order, error) {
	var entity *order.Order
	err := r.db.WithContext(ctx).Preload("Status").Where("id = ?", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}
