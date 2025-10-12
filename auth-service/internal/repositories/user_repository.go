package repositories

import (
	"errors"

	"gorm.io/gorm"
	"auth-service/internal/models"
)

type Repository struct {
	DB *gorm.DB
}


func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := r.DB.Preload("Role").Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *Repository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) IsEmailExists(email string) bool {
	var count int64
	r.DB.Model(&models.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

func (r *Repository) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *Repository) UpdateUser(user *models.User) error {
	return r.DB.Save(user).Error
}

func (r *Repository) DeactivateUserAccount(userID uint) error {
	var user models.User
	if err := r.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("ไม่พบผู้ใช้")
		}
		return err
	}
	return r.DB.Delete(&user).Error
}

func (r *Repository) DeleteUserByID(id uint) error {
	var user models.User
	if err := r.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("ไม่พบผู้ใช้ที่ระบุ")
		}
		return err
	}
	return r.DB.Delete(&user).Error
}

func (r *Repository) GetDeletedUsers() ([]models.User, error) {
	var users []models.User
	err := r.DB.Unscoped().Where("deleted_at IS NOT NULL").Find(&users).Error
	return users, err
}

func (r *Repository) RestoreUserByID(id uint) error {
	return r.DB.Unscoped().
		Model(&models.User{}).
		Where("user_id = ?", id).
		Update("deleted_at", nil).Error
}