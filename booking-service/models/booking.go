package models

import "time"

type Booking struct {
    BookingID  uint       `json:"booking_id" gorm:"primaryKey"`
    CourseID   uint       `json:"course_id"`
    UserID     uint       `json:"user_id"` // เก็บเฉพาะ ID ของ Auth Service
    StartTime  *time.Time `json:"start_time"`
    EndTime    *time.Time `json:"end_time"`
    SlipUrl    string     `json:"slip" gorm:"column:slip_url"`
    SlipStatus string     `json:"slip_status"`
    Status     string     `json:"status" gorm:"default:PENDING"`
    BookingAt  time.Time  `json:"booking_at" gorm:"autoCreateTime"`
    CancelAt   *time.Time `json:"cancel_at"`
    CompleteAt *time.Time `json:"complete_at"`

    BookingDogs []BookingDog `gorm:"foreignKey:BookingID;constraint:OnDelete:CASCADE;"`
}
