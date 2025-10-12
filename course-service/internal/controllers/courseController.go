package controllers

import (
	"course-service/internal/models"
	"course-service/internal/repositories"
	"course-service/internal/validators"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllCourses(c *gin.Context) {
	courses, err := repositories.GetAllCourses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, courses)
}

func GetCourseByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	course, err := repositories.GetCourseByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}
	c.JSON(http.StatusOK, course)
}

func CreateCourse(c *gin.Context) {
	var validator validators.CreateCourseValidator
	if err := c.ShouldBindJSON(&validator); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(validator.PaymentsName)

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found in context"})
		return
	}
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format in context"})
		return
	}


	fmt.Printf("Type: %T, Value: %v\n", userID, userID)
	course := models.Course{
		CourseName:     validator.CourseName,
		Description:    validator.Description,
		Address:      validator.Address,
		Province:     validator.Province,
		PostalCode:   validator.PostalCode,
		Price:        validator.Price,
		PaymentsName:   validator.PaymentsName,
		PaymentsMethod: validator.PaymentsMethod,
		UserID:         userIDUint,
	}

	if err := repositories.CreateCourse(&course); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := validators.ResponseCourse{
		CourseID:    int64(course.ID),
		CourseName:  course.CourseName,
		Description: course.Description,
		Address:     course.Address,
		Province:    course.Province,
		PostalCode:  course.PostalCode,
		PaymentsMethod: course.PaymentsMethod,
		PaymentsName: course.PaymentsName,
		Price:       course.Price,
		UserID:      int64(course.UserID),
	}

	c.JSON(http.StatusCreated, response)
}

func UpdateCourse(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var validator validators.UpdateCourseValidator

	if err := c.ShouldBindJSON(&validator); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found in context"})
		return
	}
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format in context"})
		return
	}

	course := models.Course{
		ID:           uint(id),
		CourseName:     validator.CourseName,
		Description:    validator.Description,
		Address:      validator.Address,
		Province:     validator.Province,
		PostalCode:   validator.PostalCode,
		Price:        validator.Price,
		PaymentsName:   validator.PaymentsName,
		PaymentsMethod: validator.PaymentsMethod,
		UserID:         userIDUint,
	}

	if err := repositories.UpdateCourse(&course, userIDUint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := validators.ResponseCourse{
		CourseID:    int64(course.ID),
		CourseName:  course.CourseName,
		Description: course.Description,
		Address:     course.Address,
		Province:    course.Province,
		PostalCode:  course.PostalCode,
		UserID:      int64(course.UserID),
	}

	c.JSON(http.StatusOK, response)
}

func DeleteCourse(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id not found in context"})
		return
	}
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format in context"})
		return
	}
	if err := repositories.DeleteCourse(uint(id), userIDUint); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "คอร์สถูกลบสำเร็จ"})
}
