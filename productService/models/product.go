package models

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ID          uuid.UUID `gorm:"id"`
	Name        string    `gorm:"name"`
	Description string    `gorm:"description"`
	Quantity    int       `gorm:"quantity"`
	Price       float64   `gorm:"price"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
type DetailsRequest struct {
	ProductIDs []string `json:"product_ids"`
}

type ProductDetailsResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Price    float64   `json:"price"`
	Quantity int       `json:"quantity"`
}
