package validators

import "time"

type CreateBooking struct {
    CourseID  uint      `json:"course_id" binding:"required"`
    StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
    Slip      string    `json:"slip" binding:"required"`
    SlipStatus string `json:"slip_status" binding:"required"`
    DogID    []uint    `json:"dog_ids" binding:"required"`
    DogAge   []string     `json:"dog_ages" binding:"required"`
}

type BookingDog struct {
    BookingID uint
    DogID     uint
}