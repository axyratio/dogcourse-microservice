package repositories

import (
	"fmt"
	"booking-service/config"
	"booking-service/internal/models"
	"time"

	"gorm.io/gorm"
)

func GetBookings(bookings *[]models.Booking) error {
	result := config.DB.
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("user_id", "name", "last_name", "email", "role_id", "created_at", "updated_at")
		}).
		Preload("Dogs").
		Preload("Dogs.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("user_id", "name", "last_name", "email", "role_id", "created_at", "updated_at")
		}).
		Preload("Course").
		Find(bookings)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}


func CreateBooking(booking *models.Booking) error {
	return config.DB.Create(booking).Error
}


func CreateBookingWithDogs(booking *models.Booking, dogIDs []uint, dogAges []string) (*models.Booking, error) {
	if len(dogIDs) == 0 || len(dogIDs) != len(dogAges) {
		return nil, fmt.Errorf("invalid dog arrays")
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// สร้าง booking
		if err := tx.Create(booking).Error; err != nil {
			return err
		}

		// สร้าง booking_dogs
		for i, dogID := range dogIDs {
			bd := models.BookingDog{
				BookingID: booking.BookingID,
				DogID:     dogID,
				DogAge:    dogAges[i],
			}
			if err := tx.Create(&bd).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// ✅ preload หลังจากสร้างเสร็จ
	var created models.Booking
	err = config.DB.
		Preload("User", func(db *gorm.DB) *gorm.DB {
    return db.Select("user_id", "name", "last_name", "email", "role_id", "created_at", "updated_at")
}).
		// ✅ ใช้ alias 'bd' แล้วบังคับใช้ชื่อคอลัมน์จริง 'booking_id'
		Preload("Dogs", func(db *gorm.DB) *gorm.DB {
			return db.Joins("JOIN booking_dogs AS bd ON bd.dog_id = dogs.dog_id").
				Where("bd.booking_id = ?", booking.BookingID)
		}).
		First(&created, booking.BookingID).Error

	if err != nil {
		return booking, fmt.Errorf("created booking but preload failed: %v", err)
	}

	return &created, nil
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
    booking.Status = "APPROVED"
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
    booking.Status = "REJECTED"
    booking.CancelAt = &now
    booking.CompleteAt = nil

    if err := config.DB.Save(&booking).Error; err != nil {
        return nil, err
    }

    return &booking, nil
}