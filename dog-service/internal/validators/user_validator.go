package validators

type User struct {
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
}

type BookingDog struct {
	BookingDogID int64   `json:"booking_dog_id"`
	DogID        int64   `json:"dog_id"`
	DogAge       float64 `json:"dog_age"`
}
