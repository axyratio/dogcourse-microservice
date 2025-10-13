package models

import "time"

type Booking struct {
    BookingID  uint       `json:"booking_id" gorm:"primaryKey"`
    CourseID   uint       `json:"course_id"`
    UserID     uint       `json:"user_id"`
    StartTime  *time.Time `json:"start_time"`
    EndTime    *time.Time `json:"end_time"`
    SlipUrl    string     `json:"slip" gorm:"column:slip_url"`
    SlipStatus string     `json:"slip_status"`
    Status     string     `json:"status" gorm:"default:PENDING"`
    BookingAt  time.Time  `json:"booking_at" gorm:"autoCreateTime"`
    CancelAt   *time.Time `json:"cancel_at"`
    CompleteAt *time.Time `json:"complete_at"`

    BookingDogs []BookingDog `gorm:"foreignKey:BookingID;constraint:OnDelete:CASCADE;"`

    Course Course `gorm:"foreignKey:CourseID"`
    User   User   `gorm:"foreignKey:UserID"`
    Dogs   []Dog  `gorm:"many2many:booking_dogs;joinForeignKey:BookingID;joinReferences:DogID;constraint:OnDelete:CASCADE;"`
}


type Course struct {
    ID   uint   `json:"id" gorm:"primaryKey"`
    Name string `json:"name"`
}

type User struct {
    ID   uint   `json:"id" gorm:"primaryKey"`
    Name string `json:"name"`
}

type Dog struct {
    ID   uint   `json:"id" gorm:"primaryKey"`
    Name string `json:"name"`
}
