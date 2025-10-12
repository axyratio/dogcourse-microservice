package repositories

import (
	"errors"
	"course-service/config"
	"course-service/internal/models"
	"gorm.io/gorm"
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

func UpdateCourse(course *models.Course, userID uint) error {
	var existingCourse models.Course
	result := config.DB.First(&existingCourse, course.ID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("course not found")
		}
		return result.Error
	}

	if existingCourse.UserID != userID {
		return errors.New("unauthorized: you do not own this course")
	}

	return config.DB.Save(course).Error
}


func DeleteCourse(id uint, userID uint) error {
	result := config.DB.Where("user_id = ?", userID).Delete(&models.Course{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("course not found or unauthorized")
	}
	return nil
}
