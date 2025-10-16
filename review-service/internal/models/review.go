package models


type Review struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	UserID   uint   `json:"user_id"`
	CourseID uint   `json:"course_id"`
	Rating   int    `json:"rating"`
	Comment  string `json:"comment"`

}
