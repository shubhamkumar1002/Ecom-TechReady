package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:char(255);primaryKey"`
	Name     string
	Password string
	Email    string   `gorm:"unique"`
	Phone    string   `gorm:"unique"`
	Role     UserRole `gorm:"default:'REGISTERED_USER'"`
}

type UserLogin struct {
	Password string
	Email    string
}

type UserDetails struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Phone string    `json:"phone"`
	Role  UserRole  `json:"role"`
}
type UserRole string

const (
	RegUser UserRole = "REGISTERED_USER"
	Admin   UserRole = "ADMIN"
)
