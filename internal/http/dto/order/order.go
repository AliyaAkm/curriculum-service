package order

import (
	"curriculum-service/internal/domain/orderstatus"
	"github.com/google/uuid"
)

type OrderRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	CourseID uuid.UUID `json:"course_id"`
	Amount   int       `json:"amount"`
	Currency string    `json:"currency"`
}

type OrderResponse struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	CourseID uuid.UUID `json:"course_id"`
	Amount   int       `json:"amount"`
	Currency string    `json:"currency"`

	Status orderstatus.OrderStatus `json:"status"`
}
