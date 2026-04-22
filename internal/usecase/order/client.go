package order

import (
	"context"
	"curriculum-service/internal/domain/order"
	"curriculum-service/internal/domain/orderstatus"
	"curriculum-service/internal/domain/price"
	"github.com/google/uuid"
)

type Repository interface {
	CreateOrder(ctx context.Context, value *order.Order) error
	GetOrderByID(ctx context.Context, id uuid.UUID) (*order.Order, error)
}
type StatusRepository interface {
	GetOrderStatusByCode(ctx context.Context, code string) (*orderstatus.OrderStatus, error)
}
type PriceRepository interface {
	GetPriceByCourseID(ctx context.Context, courseID uuid.UUID) (*price.CoursePrice, error)
}
