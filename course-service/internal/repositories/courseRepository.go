package repositories

import (
	"errors"
	"course-service/config"
	"course-service/internal/models"
)

func GetAllCourses() ([]models.Course, error) {
	var courses []models.Course
	if err := config.DB.Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

func GetCourseByID(id uint) (*models.Course, error) {
	var course models.Course
	if err := config.DB.First(&course, id).Error; err != nil {
		return nil, err
	}
	return &course, nil
}

func CreateCourse(course *models.Course) error {
	return config.DB.Create(course).Error
}

func UpdateCourse(course *models.Course) error {
	return config.DB.Save(course).Error
}

func DeleteCourse(id uint) error {
	result := config.DB.Delete(&models.Course{}, id)
	if result.RowsAffected == 0 {
		return errors.New("course not found")
	}
	return result.Error
}
