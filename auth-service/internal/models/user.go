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
	

	RoleID uint `json:"role_id"`                    // FK ไปยัง Role
	Role   Role `gorm:"foreignKey:RoleID" json:"-"` // ความสัมพันธ์

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`


}

// models/user.go
func (User) TableName() string { return "users" }