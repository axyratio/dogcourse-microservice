package validators

import (
	"mime/multipart"
	"time"
)

type GetBookingID struct {
    BookingID uint `json:"booking_id"`
}

type DogInput struct {
	DogID  uint    `json:"dog_id"`
	DogAge float64 `json:"dog_age"`
}

type CreateBooking struct {
	CourseID  uint                  `form:"course_id" binding:"required"`
	StartTime time.Time             `form:"start_time" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	EndTime   time.Time             `form:"end_time"   binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	SlipFile  *multipart.FileHeader `form:"slip_image"` // <- รับไฟล์โดยตรง
	SlipURL   string                `form:"slip_url"`
	DogsRaw   string                `form:"dogs" binding:"required"`
	Dogs      []DogInput            `json:"-"`
}

// type UpdateBookingSlipStatus struct {
//     SlipStatus string `json:"slip_status" binding:"required,oneof=APPROVED REJECTED"`
// }

type ResponseBooking struct {
	BookingID  uint       `json:"booking_id"`
	CourseID   uint       `json:"course_id"` // reference ไปยัง Course Service
	UserID     uint       `json:"user_id"`   // reference ไปยัง User Service
	StartTime  *time.Time `json:"start_time"`
	EndTime    *time.Time `json:"end_time"`
	SlipStatus string     `json:"slip_status"`
	SlipUrl    string     `json:"slip"`
	BookingAt  time.Time  `json:"booking_at"`

	Course  *Course   `json:"course,omitempty"`
	User    *User     `json:"user,omitempty"`
	Dogs    []Dog     `json:"dogs,omitempty"`
	DogAges []float64 `json:"dog_ages,omitempty"`
}

// DTO Struct
// จาก User Service
type User struct {
	UserID uint   `json:"user_id"`
	Name   string `json:"full_name"`
	Email  string `json:"email"`
}

// จาก Course Service
type Course struct {
	CourseID       int64   `json:"id"`
	CourseName     string  `json:"course_name"`
	Description    string  `json:"description"`
	Address        string  `json:"address"`
	Province       string  `json:"province"`
	PostalCode     string  `json:"postal_code"`
	Price          float32 `json:"price"`
	PaymentsName   string  `json:"payments_name"`
	PaymentsMethod string  `json:"payments_method"`
}

// จาก Dog Service
type Dog struct {
	DogID  uint    `json:"dog_id"`
	Name   string  `json:"name"`
	Breed  string  `json:"breed"`
	Weight float64 `json:"weight"`
	Gender string  `json:"gender"`
    // Age    float64 `json:"age"`
}
