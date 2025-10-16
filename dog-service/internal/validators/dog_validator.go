package validators

// type Dog struct {
//     DogID  int64  `gorm:"primaryKey;autoIncrement" json:"dog_id"`
//     UserID int64  `json:"user_id"`
//     Name   string `json:"name"`
//     Gender string `json:"gender"` // ใช้ enum ที่ layer validation
//     Age    int    `json:"age"`
//     Weight float64 `json:"weight"`
//     Breed  string `json:"breed"`
// }
// this model

type CreateDog struct {
	Name   string  `json:"name" binding:"required"`
	Gender string  `json:"gender" binding:"required"`
	Weight float64 `json:"weight" binding:"required"`
	Breed  string  `json:"breed" binding:"required"`
}

// "name" : "momo",
// "gender" : "female",
// "age" : 4,
// "weight" : 5.5,
// "breed" : "poodle"

type UpdateDog struct {
	Name   string  `json:"name" binding:"omitempty"`
	Gender string  `json:"gender" binding:"omitempty"`
	Weight float64 `json:"weight" binding:"omitempty"`
	Breed  string  `json:"breed" binding:"omitempty"`
}

type ResponseDog struct {
	Message string  `json:"message,omitempty"` // Optional message field
	DogID   int64   `json:"dog_id"`
	Name    string  `json:"name"`
	Gender  string  `json:"gender"`
	Weight  float64 `json:"weight"`
	Breed   string  `json:"breed"`
}

type DogBatchRequest struct {
	DogIDs []uint `json:"dog_ids" binding:"required"`
}