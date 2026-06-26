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
	PaymentsName   string    `json:"payments_name"`
	PaymentsMethod string    `json:"payments_method"`
    CreatedAt *time.Time `gorm:"autoCreateTime" json:"-"`
    UpdatedAt *time.Time `gorm:"autoUpdateTime" json:"-"`
}