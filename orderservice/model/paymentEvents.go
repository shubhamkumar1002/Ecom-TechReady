package model

import "github.com/google/uuid"

type PaymentEvent struct {
	OrderID       uuid.UUID `json:"order_id"`
	TotalAmount   float64   `json:"total_amount"`
	PaymentStatus string    `json:"payment_status"`
}
