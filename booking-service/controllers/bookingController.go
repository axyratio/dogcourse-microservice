package controllers

import (
	"booking-service/models"
	"booking-service/repositories"
	"booking-service/validators"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type User struct {
	UserID   uint   `json:"user_id"`
	Name     string `json:"name"`
	LastName string `json:"lastname"`
	Email    string `json:"email"`
}

func GetUserInfo(userID uint) (*User, error) {
	resp, err := http.Get(fmt.Sprintf("http://auth_service:8081/api/v1/auth/users/%d", userID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func callDogService(dogIDs []uint) []map[string]interface{} {
	var dogs []map[string]interface{}
	client := &http.Client{Timeout: 5 * time.Second}

	for _, id := range dogIDs {
		req, _ := http.NewRequest("GET", fmt.Sprintf("http://dog_service:8084/api/v1/dogs/%d", id), nil)
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		var dog map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&dog); err != nil {
			continue
		}
		dogs = append(dogs, dog)
	}

	return dogs
}

func GetBookings(c *gin.Context) {
	var bookings []models.Booking

	if err := repositories.GetBookings(&bookings); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบการจองทั้งหมด"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูลการจองทั้งหมดได้"})
		}
		return
	}

	c.JSON(http.StatusOK, bookings)
}

func GetBookingByID(c *gin.Context) {
	idParam := c.Param("id")
	var bookingID uint
	if _, err := fmt.Sscanf(idParam, "%d", &bookingID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking ID ไม่ถูกต้อง"})
		return
	}

	booking, err := repositories.GetBookingByIDWithPreload(bookingID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบการจอง"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูลการจองได้"})
		}
		return
	}

	// ดึง Dog IDs จาก BookingDogs
	var dogIDs []uint
	for _, bd := range booking.BookingDogs {
		dogIDs = append(dogIDs, bd.DogID)
	}

	// เรียก Dog Service เพื่อดึงรายละเอียดสุนัข
	dogs := callDogService(dogIDs)

	// เรียก Auth Service เพื่อดึงข้อมูลผู้จอง
	user, _ := GetUserInfo(booking.UserID)

	c.JSON(http.StatusOK, gin.H{
		"booking": booking,
		"user":    user,
		"dogs":    dogs,
	})
}

func CreateBooking(c *gin.Context) {
	var input validators.CreateBooking
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ดึง user_id จาก context
	var userID uint
	if v, ok := c.Get("user_id"); ok {
		if id, ok2 := v.(uint); ok2 {
			userID = id
		}
	}
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// ตรวจสอบข้อมูลสุนัข
	if len(input.DogID) == 0 || len(input.DogID) != len(input.DogAge) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dog_ids และ dog_ages ไม่ถูกต้อง"})
		return
	}

	booking := models.Booking{
		CourseID:   input.CourseID,
		UserID:     userID,
		Status:     "PENDING",
		BookingAt:  time.Now(),
		StartTime:  &input.StartTime,
		EndTime:    &input.EndTime,
		SlipUrl:    input.Slip,
		SlipStatus: input.SlipStatus,
	}

	createdBooking, err := repositories.CreateBookingWithDogs(&booking, input.DogID, input.DogAge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "ไม่สามารถสร้างการจองได้",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "สร้างการจองสำเร็จ",
		"data": gin.H{
			"booking_id":  createdBooking.BookingID,
			"course_id":   createdBooking.CourseID,
			"user_id":     createdBooking.UserID,
			"status":      createdBooking.Status,
			"start_time":  createdBooking.StartTime,
			"end_time":    createdBooking.EndTime,
			"booking_at":  createdBooking.BookingAt,
			"slip_url":    createdBooking.SlipUrl,
			"slip_status": createdBooking.SlipStatus,
		},
	})
}

func ApproveBooking(c *gin.Context) {
	idParam := c.Param("id")
	var bookingID uint
	if _, err := fmt.Sscanf(idParam, "%d", &bookingID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking ID ไม่ถูกต้อง"})
		return
	}

	booking, err := repositories.ApproveBooking(bookingID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking approved successfully",
		"data": gin.H{
			"booking_id": booking.BookingID,
			"status":     booking.Status,
		},
	})
}

func RejectBooking(c *gin.Context) {
	idParam := c.Param("id")
	var bookingID uint
	if _, err := fmt.Sscanf(idParam, "%d", &bookingID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking ID ไม่ถูกต้อง"})
		return
	}

	booking, err := repositories.RejectBooking(bookingID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking rejected successfully",
		"data": gin.H{
			"booking_id": booking.BookingID,
			"status":     booking.Status,
		},
	})
}
