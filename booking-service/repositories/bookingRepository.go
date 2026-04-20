package repositories

import (
	"booking-service/config"
	"booking-service/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func GetBookings(bookings *[]models.Booking) error {
	result := config.DB.Find(bookings)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func CreateBookingWithDogs(booking *models.Booking, dogIDs []uint, dogAges []string) (*models.Booking, error) {
	if len(dogIDs) == 0 || len(dogIDs) != len(dogAges) {
		return nil, fmt.Errorf("invalid dog arrays")
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(booking).Error; err != nil {
			return err
		}

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

	return booking, nil
}

func ApproveBooking(id uint) (*models.Booking, error) {
	var booking models.Booking
	if err := config.DB.First(&booking, id).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	booking.Status = "APPROVED"
	booking.CompleteAt = &now
	booking.CancelAt = nil
	if err := config.DB.Save(&booking).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func RejectBooking(id uint) (*models.Booking, error) {
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

func GetBookingByIDWithPreload(bookingID uint) (*models.Booking, error) {
	var booking models.Booking
	if err := config.DB.Preload("BookingDogs").First(&booking, bookingID).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}