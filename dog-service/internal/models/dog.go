package models

import "gorm.io/gorm"

type Dog struct {
	DogID  int64   `gorm:"primaryKey;autoIncrement" json:"dog_id"`
	Name   string  `json:"name"`
	Gender string  `json:"gender"`
	Weight float64 `json:"weight"`
	Breed  string  `json:"breed"`

	UserID uint `json:"user_id"`
	// User   User `gorm:"foreignKey:UserID;references:UserID;constraint:OnDelete:CASCADE;" json:"user"`

	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	// BookingDogs []BookingDog   `gorm:"foreignKey:DogID;constraint:OnDelete:CASCADE;" json:"-"`
}