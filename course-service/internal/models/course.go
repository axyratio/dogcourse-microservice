package models

import "time"

type Course struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CourseName  string    `json:"course_name"`
	Description string    `json:"description"`
	UserID      uint      `json:"user_id"` // trainer
	Address     string    `json:"address"`
	Province    string    `json:"province"`
	PostalCode  string    `json:"postal_code"`
	Price       float32   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
