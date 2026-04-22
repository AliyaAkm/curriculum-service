package order

import (
	"context"
	"curriculum-service/internal/domain/order"
)

type client interface {
	CreateOrder(ctx context.Context, value *order.Order) (*order.Order, error)
}
