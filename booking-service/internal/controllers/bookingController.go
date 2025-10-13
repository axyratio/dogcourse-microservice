package controllers

import (
	"fmt"
	"booking-service/internal/models"
	"booking-service/internal/repositories"
	"booking-service/internal/validators"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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

	// แปลง id จาก string → uint
	var id uint
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "รหัสการจองไม่ถูกต้อง"})
		return
	}

	// ดึงข้อมูล booking พร้อม preload (User, Dogs, Course)
	booking, err := repositories.GetBookingByIDWithPreload(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบการจองที่ระบุ"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูลการจองได้"})
		}
		return
	}

	// ✅ แสดงทั้ง object ที่ preload มาเลย
	c.JSON(http.StatusOK, gin.H{
		"message": "ดึงข้อมูลการจองสำเร็จ",
		"data":    booking,
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
	if len(input.DogID) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ต้องส่ง dog_ids อย่างน้อย 1 รายการ"})
		return
	}
	if len(input.DogID) != len(input.DogAge) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "จำนวน dog_ids และ dog_ages ไม่ตรงกัน"})
		return
	}

	// เตรียม booking entity
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

	// เรียก repo ให้บันทึกลง DB
	createdBooking, err := repositories.CreateBookingWithDogs(&booking, input.DogID, input.DogAge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "ไม่สามารถสร้างการจองได้",
			"detail": err.Error(),
		})
		return
	}

	// ✅ ส่งคืนเฉพาะข้อมูลที่จำเป็น
	response := gin.H{
		"booking_id":  createdBooking.BookingID,
		"course_id":   createdBooking.CourseID,
		"user_id":     createdBooking.UserID,
		"status":      createdBooking.Status,
		"start_time":  createdBooking.StartTime,
		"end_time":    createdBooking.EndTime,
		"booking_at":  createdBooking.BookingAt,
		"slip_url":    createdBooking.SlipUrl,
		"slip_status": createdBooking.SlipStatus,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "สร้างการจองสำเร็จ",
		"data":    response,
	})
}



func ApproveBooking(c *gin.Context) {
    id := c.Param("id")

    booking, err := repositories.ApproveBooking(id)
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

// ปฏิเสธการจอง
func RejectBooking(c *gin.Context) {
    id := c.Param("id")

    booking, err := repositories.RejectBooking(id)
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