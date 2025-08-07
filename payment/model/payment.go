package model

import (
	"github.com/google/uuid"
	"time"
)

type PaymentStatus string

const (
	Pending        PaymentStatus = "PENDING"
	Paid           PaymentStatus = "PAID"
	Cancelled      PaymentStatus = "CANCELLED"
	RefundStarted  PaymentStatus = "REFUND STARTED"
	RefundComplete PaymentStatus = "REFUND COMPLETE"
)

type Payment struct {
	ID            uuid.UUID     `gorm:"type:char(255);primaryKey"`
	OrderID       uuid.UUID     `gorm:"not null"`
	TotalAmount   float64       `gorm:"not null"` // price * quantity
	PaymentStatus PaymentStatus `gorm:"type:varchar(20);default:'PENDING'"`
	CreatedAt     time.Time     `gorm:"autoCreateTime"`
	UpdatedAt     time.Time     `gorm:"autoUpdateTime"`
}

type PaymentCreateDTO struct {
	OrderID       uuid.UUID     `json:"not null"`
	TotalAmount   float64       `json:"not null"`
	PaymentStatus PaymentStatus `json:"type:varchar(20);default:'PENDING'"`
}
