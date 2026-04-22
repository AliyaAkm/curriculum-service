package order

import (
	"context"
	"curriculum-service/internal/domain/order"
	"curriculum-service/internal/domain/orderstatus"
	"github.com/google/uuid"
)

type Repository interface {
	CreateOrder(ctx context.Context, value *order.Order) error
	GetOrderByID(ctx context.Context, id uuid.UUID) (*order.Order, error)
}
type StatusRepository interface {
	GetOrderStatusByCode(ctx context.Context, code string) (*orderstatus.OrderStatus, error)
}
