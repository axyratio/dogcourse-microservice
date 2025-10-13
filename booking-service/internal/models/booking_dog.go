package models

// BookingDog เป็น join table ระหว่าง Booking และ Dog
type BookingDog struct {
    BookingID uint `gorm:"primaryKey"`
    DogID     uint `gorm:"primaryKey"`
    DogAge    string
}