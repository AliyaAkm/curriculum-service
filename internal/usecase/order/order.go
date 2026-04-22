package order

import (
	"context"
	"curriculum-service/internal/domain/order"
	"github.com/google/uuid"
)

const pendingStatus = "pending"

func (u *UseCase) CreateOrder(ctx context.Context, value *order.Order) (*order.Order, error) {
	value.ID = uuid.New()

	status, err := u.statusRepo.GetOrderStatusByCode(ctx, pendingStatus)
	if err != nil {
		return nil, err
	}
	value.StatusID = status.ID

	price, err := u.priceRepo.GetPriceByCourseID(ctx, value.CourseID)
	if err != nil {
		return nil, err
	}
	value.Amount = price.Amount
	value.Currency = price.Currency

	err = u.repo.CreateOrder(ctx, value)
	if err != nil {
		return nil, err
	}
	return u.repo.GetOrderByID(ctx, value.ID)
}
