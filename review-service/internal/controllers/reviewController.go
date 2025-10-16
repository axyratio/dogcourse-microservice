package controllers

import (
	"net/http"
	"review-service/internal/models"
	"review-service/internal/repositories"
	// "review-service/internal/service"
	"review-service/internal/validators"
	"strconv"
	// "strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateReview creates a new review.
// func getAuthHeader(c *gin.Context) string {
// 	// ใช้ค่า header ที่ middleware เก็บไว้ก่อน
// 	if v, ok := c.Get("auth_header"); ok {
// 		if s, ok2 := v.(string); ok2 && s != "" {
// 			return s
// 		}
// 	}
// 	h := c.GetHeader("Authorization")
// 	if h == "" {
// 		return ""
// 	}
// 	if !strings.HasPrefix(strings.ToLower(h), "bearer ") {
// 		return "Bearer " + h
// 	}
// 	return h
// }

// func CreateReview(c *gin.Context) {
// 	var review models.Review
// 	if err := c.ShouldBindJSON(&review); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// user จาก middleware
// 	val, ok := c.Get("user_id")
// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "ไม่พบ token หรือไม่ได้เข้าสู่ระบบ"})
// 		return
// 	}
// 	uid, ok := val.(uint)
// 	if !ok || uid == 0 {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id in context"})
// 		return
// 	}
// 	review.UserID = uid

// 	// course_id จาก path
// 	courseIDParam := c.Param("id") // ← เอา paramID ไปเช็กคอร์ส
// 	cidU64, err := strconv.ParseUint(courseIDParam, 10, 64)
// 	if err != nil || cidU64 == 0 {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "course_id ไม่ถูกต้อง"})
// 		return
// 	}
// 	review.CourseID = uint(cidU64)

// 	// แนบ Bearer
// 	authHeader := getAuthHeader(c)
// 	if authHeader == "" {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
// 		return
// 	}

// 	// 1) เช็กว่ามีคอร์สจริงจาก course-service
// 	if err := services.AssertCourseExistsRawID(courseIDParam, authHeader); err != nil {
// 		if services.IsNotFound(err) {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "ไม่พบคอร์สนี้"})
// 			return
// 		}
// 		c.JSON(http.StatusBadGateway, gin.H{"error": "ติดต่อ course-service ไม่สำเร็จ"})
// 		return
// 	}

// 	// 2) เรียก booking-service: getBooked (ตรวจคนยิง ≡ user_id ใน booking)
// 	booked, err := services.GetBookedRawID(courseIDParam, authHeader)
// 	if err != nil {
// 		c.JSON(http.StatusBadGateway, gin.H{"error": "ติดต่อ booking-service ไม่สำเร็จ"})
// 		return
// 	}
// 	if !booked {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "ยังไม่เคยจองคอร์สนี้ ไม่สามารถรีวิวได้"})
// 		return
// 	}

// 	// 3) ไม่ให้รีวิวซ้ำ
// 	if existsAlready, err := repositories.CheckIfReviewed(review.UserID, review.CourseID); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "ตรวจสอบรีวิวไม่สำเร็จ"})
// 		return
// 	} else if existsAlready {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "คุณได้รีวิวคอร์สนี้ไปแล้ว"})
// 		return
// 	}

// 	// 4) validate เนื้อรีวิว
// 	if err := validators.ValidateReview(review); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// 5) บันทึก
// 	if err := repositories.CreateReview(&review); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "สร้างรีวิวไม่สำเร็จ"})
// 		return
// 	}
// 	c.JSON(http.StatusCreated, review)
// }

// CreateReview creates a new review using user_id from JWT.
func CreateReview(c *gin.Context) {
	var review models.Review

	// รับ rating/comment จาก body
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ดึง user_id จาก token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ไม่พบ token หรือไม่ได้เข้าสู่ระบบ"})
		return
	}
	review.UserID = userID.(uint)

	// ดึง course_id จาก URL
	courseIDParam := c.Param("id")
	courseID, err := strconv.ParseUint(courseIDParam, 10, 64)
	if err != nil || courseID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "course_id ไม่ถูกต้อง"})
		return
	}
	review.CourseID = uint(courseID)

	// ตรวจสอบว่าผู้ใช้นี้รีวิวคอร์สนี้ไปแล้วหรือยัง
	existsAlready, err := repositories.CheckIfReviewed(review.UserID, review.CourseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ตรวจสอบรีวิวไม่สำเร็จ"})
		return
	}
	if existsAlready {
		c.JSON(http.StatusBadRequest, gin.H{"error": "คุณได้รีวิวคอร์สนี้ไปแล้ว"})
		return
	}

	// ตรวจสอบความถูกต้อง
	if err := validators.ValidateReview(review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// บันทึกลงฐานข้อมูล
	if err := repositories.CreateReview(&review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "สร้างรีวิวไม่สำเร็จ"})
		return
	}

	c.JSON(http.StatusCreated, review)
}


// GetReviews retrieves all reviews.
func GetReviews(c *gin.Context) {
	reviews, err := repositories.FindAllReviews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ดึงข้อมูลรีวิวไม่สำเร็จ"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// GetReview retrieves a single review by its ID.
func GetReview(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID ไม่ถูกต้อง"})
		return
	}

	review, err := repositories.FindReviewByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบรีวิว"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ดึงข้อมูลรีวิวไม่สำเร็จ"})
		return
	}

	c.JSON(http.StatusOK, review)
}



// UpdateReview updates an existing review.
func UpdateReview(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID ไม่ถูกต้อง"})
		return
	}

	existingReview, err := repositories.FindReviewByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบรีวิว"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ดึงข้อมูลรีวิวไม่สำเร็จ"})
		return
	}

	var updatedReview models.Review
	if err := c.ShouldBindJSON(&updatedReview); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Preserve original IDs
	updatedReview.ID = existingReview.ID
	updatedReview.UserID = existingReview.UserID
	updatedReview.CourseID = existingReview.CourseID

	if err := validators.ValidateReview(updatedReview); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := repositories.UpdateReview(&updatedReview); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "อัปเดตรีวิวไม่สำเร็จ"})
		return
	}

	c.JSON(http.StatusOK, updatedReview)
}

// DeleteReview deletes a review by its ID.
func DeleteReview(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID ไม่ถูกต้อง"})
		return
	}

	review, err := repositories.FindReviewByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบรีวิว"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ดึงข้อมูลรีวิวไม่สำเร็จ"})
		return
	}

	if err := repositories.DeleteReview(review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ลบรีวิวไม่สำเร็จ"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ลบรีวิวสำเร็จ"})
}
