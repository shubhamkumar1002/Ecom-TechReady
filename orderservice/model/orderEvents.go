package model

import "github.com/google/uuid"

type OrderEvent struct {
	OrderID       uuid.UUID `json:"order_id"`
	OrderStatus   string    `json:"order_status"`
	PaymentStatus string    `json:"payment_status"`
}
type StockUpdateFailedEvent struct {
	OrderID uuid.UUID `json:"order_id"`
	Reason  string    `json:"reason"`
}
