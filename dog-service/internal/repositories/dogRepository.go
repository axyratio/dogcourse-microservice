package repositories

import (
	"dog-service/config"
	"dog-service/internal/models"
	"gorm.io/gorm"
	"errors"
)



func GetAllDogByUserID(userID uint) ([]models.Dog, error) {
	var dogs []models.Dog
	err := config.DB.Where("user_id = ?", userID).Find(&dogs).Error
	if err != nil {
		return nil, err
	}

	return dogs, nil
}

func GetDogByID(id uint) (*models.Dog, error) {
	var dog models.Dog
	if err := config.DB.First(&dog, id).Error; err != nil {
		return nil, err
	}
	return &dog, nil
}

func AddDog(dog *models.Dog) error {
	return config.DB.Create(dog).Error
}

func UpdateDogAndCheckOwner(dog *models.Dog, userID uint, dogID int64) (bool, error) {
    // ตรวจสอบว่า dog เป็นของ user หรือไม่
    var checkDog models.Dog
    err := config.DB.Select("user_id").
        Where("dog_id = ? AND user_id = ?", dogID, userID).
        First(&checkDog).Error

    if errors.Is(err, gorm.ErrRecordNotFound) {
        // ไม่พบสุนัขนี้สำหรับ user นี้
        return false, nil
    }
    if err != nil {
        return false, err
    }

    // เช็กเจ้าของ
    isOwner := checkDog.UserID == userID
    if !isOwner {
        return false, nil
    }

    // อัปเดตข้อมูล
    if err := config.DB.Save(dog).Error; err != nil {
        return false, err
    }

    return true, nil
}

func DeleteDogAndCheckOwner(dog *models.Dog, userID uint, dogID int64) (bool, error) {
	// ตรวจว่า dog เป็นของ user หรือไม่
	var checkDog models.Dog
	err := config.DB.Select("user_id").
		Where("dog_id = ? AND user_id = ?", dogID, userID).
		First(&checkDog).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	// ตรวจเจ้าของ
	if checkDog.UserID != userID {
		return false, nil
	}

	// Soft delete (จะไม่ลบจริง แต่ mark ว่าถูกลบ)
	if err := config.DB.Delete(dog).Error; err != nil {
		return false, err
	}

	return true, nil
}

func GetAllDeletedDogs() ([]models.Dog, error) {
	var dogs []models.Dog
	err := config.DB.Unscoped().Where("deleted_at IS NOT NULL").Find(&dogs).Error
	return dogs, err
}