package models

// BookingDog เป็น join table ระหว่าง Booking และ Dog
type BookingDog struct {
    BookingID uint   `gorm:"column:booking_id;primaryKey"`
    DogID     uint   `gorm:"column:dog_id;primaryKey"`
    DogAge    string `json:"dog_age" gorm:"column:dog_age"`
}
