package validators

import (
	"github.com/go-playground/validator/v10"
	"review-service/internal/models"
)

// ReviewValidator is a struct for validating review creation and updates.
type ReviewValidator struct {
	Rating   int    `validate:"required,min=1,max=5"`
	Comment  string `validate:"required,min=1"`
	UserID   uint   `validate:"required"`
	CourseID uint   `validate:"required"`
}

// ValidateReview validates the given review model.
func ValidateReview(review models.Review) error {
	validate := validator.New()
	
	reviewValidator := ReviewValidator{
		Rating:   review.Rating,
		Comment:  review.Comment,
		UserID:   review.UserID,
		CourseID: review.CourseID,
	}

	return validate.Struct(reviewValidator)
}
