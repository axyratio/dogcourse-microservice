package validators


type CreateCourseValidator struct {
	CourseName string `json:"course_name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Address     string `json:"address" binding:"required"`
	Province	string `json:"province" binding:"required"`
	PostalCode  string `json:"postal_code" binding:"required"`
	Price	   float32   `json:"price" binding:"required"`
	PaymentsName string    `json:"payments_name" binding:"required"`
	PaymentsMethod string    `json:"payments_method" binding:"required"`
}


type UpdateCourseValidator struct {
	CourseName string `json:"course_name" binding:"omitempty"`
	Description string `json:"description" binding:"omitempty"`
	Address     string `json:"address" binding:"omitempty"`
	Province	string `json:"province" binding:"omitempty"`
	PostalCode  string `json:"postal_code" binding:"omitempty"`
	Price	   float32   `json:"price" binding:"omitempty"`
	PaymentsName string    `json:"payments_name" binding:"omitempty"`
	PaymentsMethod string    `json:"payments_method" binding:"omitempty"`
}

type ResponseCourse struct {
	Message string `json:"message,omitempty"` // Optional message field
	CourseID   int64  `json:"course_id"`
	CourseName string `json:"course_name"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Province	string `json:"province"`
	PostalCode  string `json:"postal_code"`
	Price	   float32   `json:"price"`
	PaymentsName string    `json:"payments_name"`
	PaymentsMethod string    `json:"payments_method"`
	UserID      int64  `json:"user_id"`
}
