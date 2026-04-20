package controllers

import (
	"main/models"
	"main/repositories"
	"main/validators"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateReview creates a new review.


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
