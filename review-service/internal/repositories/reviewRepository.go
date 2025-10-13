package repositories

import (
	"review-service/config"
	"review-service/internal/models"
)

// Create adds a new review to the database.
func CreateReview(review *models.Review) error {
	return config.DB.Create(review).Error
}

// FindAll retrieves all reviews from the database.
func FindAllReviews() ([]models.Review, error) {
	var reviews []models.Review
	err := config.DB.Find(&reviews).Error
	return reviews, err
}

// FindReviewByID retrieves a single review by its ID.
func FindReviewByID(id uint) (*models.Review, error) {
	var review models.Review
	err := config.DB.First(&review, id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

// CheckIfReviewed ตรวจสอบว่าผู้ใช้เคยรีวิวคอร์สนี้แล้วหรือไม่
func CheckIfReviewed(userID uint, courseID uint) (bool, error) {
	var count int64
	err := config.DB.Model(&models.Review{}).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}


// Update modifies an existing review in the database.
func UpdateReview(review *models.Review) error {
	return config.DB.Save(review).Error
}

// Delete removes a review from the database.
func DeleteReview(review *models.Review) error {
	return config.DB.Delete(review).Error
}
