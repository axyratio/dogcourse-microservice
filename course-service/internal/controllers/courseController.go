package controllers

import (
	"course-service/internal/models"
	"course-service/internal/repositories"
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
	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := repositories.CreateCourse(&course); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, course)
}

func UpdateCourse(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	course.ID = uint(id)
	if err := repositories.UpdateCourse(&course); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, course)
}

func DeleteCourse(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := repositories.DeleteCourse(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	

	c.JSON(http.StatusOK, gin.H{"message": "คอร์สถูกลบสำเร็จ"})
}
