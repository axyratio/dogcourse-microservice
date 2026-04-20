package models

import (
	"time"
)

type Role struct {
	ID   uint   `gorm:"primaryKey" json:"role_id"`
	RoleName string `json:"role_name"`
}


type User struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Name      string
	Role      string    `gorm:"default:user"`
	CreatedAt time.Time
	UpdatedAt time.Time
}


type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (User) TableName() string {
	return "users"
}
