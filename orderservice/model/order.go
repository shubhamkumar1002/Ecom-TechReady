package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string
type PaymentStatus string

const (
	Pending       PaymentStatus = "PENDING"
	Paid          PaymentStatus = "PAID"
	Cacelled      PaymentStatus = "CANCELLED"
	OrderPlaced   OrderStatus   = "ORDER PLACED"
	Shipped       OrderStatus   = "SHIPPED"
	Delivered     OrderStatus   = "DELIVERED"
	Failed        OrderStatus   = "FAILED"
	OrderCanceled OrderStatus   = "ORDER CANCELLED"
)

type Order struct {
	ID          uuid.UUID   `gorm:"type:char(255);primaryKey"`
	UserID      uuid.UUID   `gorm:"type:char(36);not null"`
	OrderItems  []OrderItem `gorm:"foreignKey:OrderID"`
	TotalAmount float64     `gorm:"not null"` // price * quantity
	OrderStatus OrderStatus `gorm:"type:varchar(20);default:'ORDER PLACED'"`
	CreatedAt   time.Time   `gorm:"autoCreateTime"`
	UpdatedAt   time.Time   `gorm:"autoUpdateTime"`
}

type OrderItem struct {
	ID        uint      `gorm:"primaryKey"`
	OrderID   uuid.UUID `gorm:"type:char(36);not null"`
	ProductID string    `gorm:"not null"`
	Quantity  int       `gorm:"not null"`
	UnitPrice float64   `gorm:"not null"`
}

/*
	type OrderCreateDTO struct {
		UserID      uuid.UUID `gorm:"type:char(36);not null"`
		ProductID   string    `gorm:"not null"`
		Quantity    int       `gorm:"not null;default:1"`
		TotalAmount float64   `gorm:"not null"`
	}
*/
type OrderUpdateDTO struct {
	OrderStatus   string `json:"order_status"`
	PaymentStatus string `json:"payment_status"`
}
type ItemRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Price     int       `json:"price"`
}

type OrderCreateDTO struct {
	UserID uuid.UUID     `json:"user_id"`
	Items  []ItemRequest `json:"items"`
	Total  int           `json:"total"`
}
type OrderCreatedEvent struct {
	OrderID uuid.UUID
	Items   []ItemRequest
	Total   float64
}
