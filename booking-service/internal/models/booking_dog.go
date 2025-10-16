package models

// BookingDog เป็น join table ระหว่าง Booking และ Dog
type BookingDog struct {
	BookingDogID uint    `json:"booking_dog_id" gorm:"primaryKey;autoIncrement"`
	BookingID    uint    `json:"booking_id"` // reference ไปยัง Booking Service
	DogID        uint    `json:"dog_id"`     // reference ไปยัง Dog Service
	DogAge       float64 `json:"dog_age"`
}