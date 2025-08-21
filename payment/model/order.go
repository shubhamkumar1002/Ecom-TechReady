package model

import (
	"github.com/google/uuid"
	"time"
)

type OrderStatus string

type Order struct {
	ID            uuid.UUID     `gorm:"type:char(255);primaryKey"`
	UserID        uint          `gorm:"not null"`
	ProductID     string        `gorm:"not null"`
	Quantity      int           `gorm:"not null;default:1"`
	TotalAmount   float64       `gorm:"not null"`
	OrderStatus   OrderStatus   `gorm:"type:varchar(20);default:'ORDER PLACED'"`
	PaymentStatus PaymentStatus `gorm:"type:varchar(20);default:'PENDING'"`
	CreatedAt     time.Time     `gorm:"autoCreateTime"`
	UpdatedAt     time.Time     `gorm:"autoUpdateTime"`
}
