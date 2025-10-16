package controllers

import (
	"booking-service/internal/models"
	"booking-service/internal/repositories"
	"booking-service/internal/service"
	"booking-service/internal/utils"
	"booking-service/internal/validators"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	// "gorm.io/gorm"
)


func isAdmin(c *gin.Context) bool {
	v, ok := c.Get("role")
	if !ok { return false }
	s, _ := v.(string)
	return s == "admin"
}

func CheckCourseBooked(c *gin.Context) {
	courseID := c.Param("id")

	// ต้องมี middleware แปะ user_id ไว้ก่อนหน้า
	userID := c.GetUint("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	booked, err := repositories.IsCourseBookedByUser(courseID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถตรวจสอบสถานะการจองได้"})
		return
	}

	fmt.Println(userID)
	fmt.Println(courseID)
	fmt.Println(booked)

	// ตอบกลับรูปแบบที่ client ของคุณคาดหวัง: { "booked": bool }
	c.JSON(http.StatusOK, gin.H{"booked": booked})
}


func GetBookings(c *gin.Context) {
	uid, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// แอดมิน: ทั้งหมด หรือกรองด้วย ?user_id=
	if isAdmin(c) {
		// ถ้ามี ?user_id= ให้กรองเฉพาะคนนั้น
		if q := c.Query("user_id"); q != "" {
			u64, err := strconv.ParseUint(q, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
				return
			}
			items, err := repositories.GetBookingsByUserID(uint(u64))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot list bookings"})
				return
			}
			c.JSON(http.StatusOK, items)
			return
		}
		// ไม่มี query → ทั้งหมด
		items, err := repositories.GetBookings()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot list bookings"})
			return
		}
		c.JSON(http.StatusOK, items)
		return
	}

	// ผู้ใช้ทั่วไป: เฉพาะของตัวเอง
	items, err := repositories.GetBookingsByUserID(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot list bookings"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func GetBookingByID(c *gin.Context) {
	// 1️⃣ แปลง id จาก path parameter
	idStr := c.Param("id")
	id64, _ := strconv.ParseUint(idStr, 10, 64)
	id := uint(id64)

	currentUID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 3) เช็กสิทธิ์: เจ้าของ หรือ แอดมิน เท่านั้น
	

	// 2️⃣ ดึงข้อมูล booking หลักจาก DB
	booking, err := repositories.GetBookingByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบข้อมูลการจอง"})
		return
	}

	if booking.UserID != currentUID && !isAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error":"forbidden"})
		return
	}

	// 3️⃣ ดึงข้อมูล booking_dogs จาก DB (ใช้ booking_id)
	bookingDogs, err := repositories.GetBookingDogsByBookingID(booking.BookingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูล booking_dogs ได้"})
		return
	}

	// 4️⃣ แปลง bookingDogs → dogIDs เพื่อยิงไป dog-service
	dogIDs := make([]uint, 0, len(bookingDogs))
	dogAges := make([]float64, 0, len(bookingDogs))
	for _, bd := range bookingDogs {
		dogIDs = append(dogIDs, bd.DogID)
		dogAges = append(dogAges, bd.DogAge)
	}

	// 5️⃣ ดึง token จาก header
	token := c.GetHeader("Authorization")

	// 6️⃣ ดึงข้อมูลจาก microservices อื่น
	user, _ := services.FetchUserByID(token)
	course, _ := services.FetchCourseByID(booking.CourseID, token)
	dogs, _ := services.FetchDogsByIDs(dogIDs, token)

	// 7️⃣ สร้าง ResponseBooking
	response := validators.ResponseBooking{
		BookingID:  booking.BookingID,
		CourseID:   booking.CourseID,
		UserID:     booking.UserID,
		StartTime:  booking.StartTime,
		EndTime:    booking.EndTime,
		SlipUrl:    booking.SlipUrl,
		SlipStatus: booking.SlipStatus,
		BookingAt:  booking.BookingAt,

		Course:  course,
		User:    user,
		Dogs:    dogs,
		DogAges: dogAges,
	}

	// 8️⃣ ส่ง response
	c.JSON(http.StatusOK, response)
}


// func GetBookingByID(c *gin.Context) {
// 	idParam := c.Param("id")

// 	// แปลง id จาก string → uint
// 	var id uint
// 	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "รหัสการจองไม่ถูกต้อง"})
// 		return
// 	}

// 	// ดึงข้อมูล booking พร้อม preload (User, Dogs, Course)
// 	booking, err := repositories.GetBookingByIDWithPreload(id)
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบการจองที่ระบุ"})
// 		} else {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูลการจองได้"})
// 		}
// 		return
// 	}

// 	// ✅ แสดงทั้ง object ที่ preload มาเลย
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "ดึงข้อมูลการจองสำเร็จ",
// 		"data":    booking,
// 	})
// }

func CreateBooking(c *gin.Context) {
    var input validators.CreateBooking

    if err := c.ShouldBindWith(&input, binding.FormMultipart); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error":"ข้อมูลไม่ถูกต้อง", "detail": err.Error()})
        return
    }

    // แปลง dogs
    if err := json.Unmarshal([]byte(input.DogsRaw), &input.Dogs); err != nil || len(input.Dogs) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error":"รูปแบบ dogs ไม่ถูกต้อง หรือว่าง"})
        return
    }

    // user_id
    v, ok := c.Get("user_id")
    if !ok { c.JSON(http.StatusUnauthorized, gin.H{"error":"Unauthorized"}); return }
    userID := v.(uint)

    // ตัด Bearer ซ้ำ
    rawAuth := strings.TrimSpace(c.GetHeader("Authorization"))
    token := strings.TrimSpace(strings.TrimPrefix(rawAuth, "Bearer"))

    // ตรวจ course
    course, err := services.FetchCourseByID(input.CourseID, token)
    if err != nil || course == nil {
        c.JSON(http.StatusBadRequest, gin.H{"error":"ไม่พบคอร์สนี้", "detail": fmt.Sprintf("%v", err)})
        return
    }

    // เตรียม dogIDs/ages แล้วไปเช็ค service
    var dogIDs []uint
    var dogAges []float64
    for _, d := range input.Dogs {
        dogIDs = append(dogIDs, d.DogID)
        dogAges = append(dogAges, d.DogAge)
    }

	fmt.Println(dogIDs)

    dogs, err := services.FetchDogsByIDs(dogIDs, token)
if err != nil {
    // จะเห็นชัดว่าได้ 401/404 หรือ body เป็นอะไร
    c.JSON(http.StatusBadRequest, gin.H{"error":"ดึงข้อมูลสุนัขล้มเหลว", "detail": err.Error()})
    return
}
if len(dogs) != len(dogIDs) {
    c.JSON(http.StatusBadRequest, gin.H{"error":"ไม่พบข้อมูลสุนัขครบถ้วน", "detail": gin.H{"sent_ids": dogIDs, "got": len(dogs)}})
    return
}

	fmt.Println(dogs)


    // จัดการสลิป (ไฟล์หรือ URL)
    slipPath := ""
    if input.SlipFile != nil {
        filename := filepath.Base(input.SlipFile.Filename)
        savePath := filepath.Join("uploads/slips", filename)
        if err := c.SaveUploadedFile(input.SlipFile, savePath); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error":"อัปโหลด slip ไม่สำเร็จ", "detail": err.Error()})
            return
        }
        slipPath = "/uploads/slips/" + filename
    } else if input.SlipURL != "" {
        slipPath = input.SlipURL
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error":"กรุณาแนบ slip_image (ไฟล์) หรือ slip_url"})
        return
    }

    // บันทึก
    booking := models.Booking{
        CourseID:   input.CourseID,
        UserID:     userID,
        Status:     "PENDING",
        BookingAt:  time.Now(),
        StartTime:  &input.StartTime,
        EndTime:    &input.EndTime,
        SlipUrl:    slipPath,
        SlipStatus: "PENDING",
    }

    createdBooking, err := repositories.CreateBookingWithDogs(&booking, dogIDs, dogAges)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error":"ไม่สามารถสร้างการจองได้", "detail": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message":"สร้างการจองสำเร็จ", "data": gin.H{
        "booking_id":  createdBooking.BookingID,
        "course_id":   createdBooking.CourseID,
        "user_id":     createdBooking.UserID,
        "status":      createdBooking.Status,
        "start_time":  createdBooking.StartTime,
        "end_time":    createdBooking.EndTime,
        "booking_at":  createdBooking.BookingAt,
        "slip_url":    createdBooking.SlipUrl,
        "slip_status": createdBooking.SlipStatus,
    }})
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
            "slip_status":     booking.SlipStatus,
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
            "slip_status":     booking.SlipStatus,
        },
    })
}