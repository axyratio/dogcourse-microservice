package repositories

import (
	"booking-service/config"
	"booking-service/internal/models"

	"time"

	"gorm.io/gorm"
)


func ExistsBookingByUserAndCourse(userID, courseID uint) (bool, error) {
	var count int64
	err := config.DB.Model(&models.Booking{}).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Count(&count).Error
	return count > 0, err
}

// mockup data res
func GetBookingByID(id uint) (*models.Booking, error) {
	var booking models.Booking

	// ดึงข้อมูลจาก DB
	if err := config.DB.First(&booking, id).Error; err != nil {
		return nil, err
	}

	return &booking, nil
}


func GetBookings() ([]models.Booking, error) {
	var items []models.Booking
	if err := config.DB.
		Order("booking_id DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func GetBookingsByUserID(userID uint) ([]models.Booking, error) {
	var items []models.Booking
	if err := config.DB.
		Where("user_id = ?", userID).
		Order("booking_id DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}


func GetBookingDogsByBookingID(bookingID uint) ([]models.BookingDog, error) {
	var bookingDogs []models.BookingDog
	if err := config.DB.Where("booking_id = ?", bookingID).Find(&bookingDogs).Error; err != nil {
		return nil, err
	}
	return bookingDogs, nil
}



func CreateBooking(booking *models.Booking) error {
	return config.DB.Create(booking).Error
}



func CreateBookingWithDogs(booking *models.Booking, dogIDs []uint, dogAges []float64) (*models.Booking, error) {
	tx := config.DB.Begin()

	// บันทึก booking หลัก
	if err := tx.Create(booking).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// บันทึก booking_dog
	for i, dogID := range dogIDs {
		bookingDog := models.BookingDog{
			BookingID: booking.BookingID,
			DogID:     dogID,
			DogAge:    dogAges[i],
		}
		if err := tx.Create(&bookingDog).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()
	return booking, nil
}

// GetBookingByIDWithPreload
// ดึง booking เดี่ยวพร้อม preload user และ dogs
func GetBookingByIDWithPreload(id uint) (*models.Booking, error) {
	var booking models.Booking

	err := config.DB.
		Model(&models.Booking{}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
	return db.Select("user_id", "name", "last_name", "email", "role_id", "created_at", "updated_at")
}).
	Preload("Dogs", func(db *gorm.DB) *gorm.DB {
			return db.
				Joins("JOIN booking_dogs bd ON bd.dog_id = dogs.dog_id").
				Where("bd.booking_id = ?", id)
		}).
		// ✅ preload user ของ dog ด้วย join ตรง ๆ
		Preload("Dogs.User", func(db *gorm.DB) *gorm.DB {
	return db.
		Joins("JOIN users ON users.user_id = dogs.user_id").
		Table("dogs").
		Select("users.user_id, users.name, users.last_name, users.email, users.role_id, users.created_at, users.updated_at")
}).
		Preload("Course").
		First(&booking, id).Error

	if err != nil {
		return nil, err
	}

	return &booking, nil
}




func ApproveBooking(id string) (*models.Booking, error) {
    var booking models.Booking

    // ค้นหา booking ตาม ID
    if err := config.DB.First(&booking, id).Error; err != nil {
        return nil, err
    }

    now := time.Now()
    booking.SlipStatus = "APPROVED"
    booking.CompleteAt = &now
    booking.CancelAt = nil

    // บันทึกลงฐานข้อมูล
    if err := config.DB.Save(&booking).Error; err != nil {
        return nil, err
    }

    // คืน booking หลังอัปเดต
    return &booking, nil
}

// ปฏิเสธการจอง
func RejectBooking(id string) (*models.Booking, error) {
    var booking models.Booking

    if err := config.DB.First(&booking, id).Error; err != nil {
        return nil, err
    }

    now := time.Now()
    booking.SlipStatus = "REJECTED"
    booking.CancelAt = &now
    booking.CompleteAt = nil

    if err := config.DB.Save(&booking).Error; err != nil {
        return nil, err
    }

    return &booking, nil
}