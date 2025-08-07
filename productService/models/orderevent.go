package models

import (
	"github.com/google/uuid"
)

type ItemRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}

type OrderCreateDTO struct {
	UserID uuid.UUID     `json:"user_id"`
	Items  []ItemRequest `json:"items"`
}
type OrderCreatedEvent struct {
	OrderID uuid.UUID
	Items   []ItemRequest
	Total   float64
}

type StockUpdatedEvent struct {
	OrderID uuid.UUID `json:"order_id"`
}

type StockUpdateFailedEvent struct {
	OrderID uuid.UUID `json:"order_id"`
	Reason  string    `json:"reason"`
}
