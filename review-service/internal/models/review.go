package models

type User struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
	// Add other fields as needed
}

type Course struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Title string `json:"title"`
	// Add other fields as needed
}

type Review struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	UserID   uint   `json:"user_id"`
	CourseID uint   `json:"course_id"`
	Rating   int    `json:"rating"`
	Comment  string `json:"comment"`

	User   User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Course Course `gorm:"foreignKey:CourseID;constraint:OnDelete:CASCADE;"`
}
