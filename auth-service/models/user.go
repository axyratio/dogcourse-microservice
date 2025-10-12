package models

import (
	"time"
	"gorm.io/gorm"
)

type Role struct {
	ID   uint   `gorm:"primaryKey" json:"role_id"`
	RoleName string `json:"role_name"`
}


type User struct {
	UserID   uint   `gorm:"primaryKey;autoIncrement" json:"user_id"`
	Name     string `gorm:"type:varchar(255)" json:"name"`
	LastName string `gorm:"type:varchar(255)" json:"lastname"`
	Email    string `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password string `gorm:"type:varchar(255);not null" json:"password"`

	RoleID uint `json:"role_id"`
	Role   Role `gorm:"foreignKey:RoleID" json:"-"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (User) TableName() string {
	return "users"
}
