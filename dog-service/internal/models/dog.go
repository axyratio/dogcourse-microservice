package models

import "gorm.io/gorm"

type Dog struct {
	DogID  int64   `gorm:"primaryKey;autoIncrement" json:"dog_id"`
	Name   string  `json:"name"`
	Gender string  `json:"gender"`
	Weight float64 `json:"weight"`
	Breed  string  `json:"breed"`

	// ต้อง match type กับ UserID
	UserID uint `gorm:"not null"`
	User   User `gorm:"foreignKey:UserID;references:UserID;constraint:OnDelete:CASCADE;" json:"-"`

	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	BookingDogs []BookingDog   `gorm:"foreignKey:DogID;constraint:OnDelete:CASCADE;"`
}

// Define User struct
type User struct {
	UserID uint   `gorm:"primaryKey"`
	Name   string `json:"name"`
	// Add other fields as needed
}

// Define BookingDog struct
type BookingDog struct {
	BookingDogID int64 `gorm:"primaryKey;autoIncrement"`
	DogID        int64
	// Add other fields as needed
}
