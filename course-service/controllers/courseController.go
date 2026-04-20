package controllers

import (
	"net/http"
	"main/config"
	"main/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllCourses(c *gin.Context) {
	var courses []models.Course
	if err := config.DB.Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูลคอร์สได้"})
		return
	}
	c.JSON(http.StatusOK, courses)
}

func GetCourseByID(c *gin.Context) {
	id := c.Param("id")
	var course models.Course

	if err := config.DB.First(&course, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบคอร์สที่ระบุ"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูลคอร์สได้"})
		}
		return
	}
	c.JSON(http.StatusOK, course)
}

func CreateCourse(c *gin.Context) {
	role, exists := c.Get("role")
	println(role)
	if !exists || (role != "trainer" && role != "admin") {
		c.JSON(http.StatusForbidden, gin.H{"error": "คุณไม่มีสิทธิ์ในการสร้างคอร์ส"})
		return
	}

	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถสร้างคอร์สได้"})
		return
	}

	c.JSON(http.StatusCreated, course)
}

func UpdateCourse(c *gin.Context) {
	id := c.Param("id")
	var course models.Course

	shouldBindErr := c.ShouldBindJSON(&course)
	if shouldBindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": shouldBindErr.Error()})
		return
	}
	var existingCourse models.Course
	if err := config.DB.First(&existingCourse, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบคอร์สที่ระบุ"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูลคอร์สได้"})
		return
	}
	existingCourse.CourseName = course.CourseName
	existingCourse.Description = course.Description
	existingCourse.Address = course.Address
	existingCourse.Province = course.Province
	existingCourse.PostalCode = course.PostalCode
	existingCourse.Price = course.Price
	existingCourse.PaymentsName = course.PaymentsName
	existingCourse.PaymmentsMethod = course.PaymmentsMethod

	if err := config.DB.Save(&existingCourse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถอัปเดตคอร์สได้"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "คอร์สถูกอัปเดตสำเร็จ"})
}

func DeleteCourse(c *gin.Context) {
	id := c.Param("id")
	var course models.Course

	if err := config.DB.First(&course, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบคอร์สที่ระบุ"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถลบคอร์สได้"})
		}
		return
	}

	if err := config.DB.Delete(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถลบคอร์สได้"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "คอร์สถูกลบสำเร็จ"})
}

