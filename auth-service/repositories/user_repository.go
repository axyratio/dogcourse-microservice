package repositories

import (
	"errors"

	"auth-service/config"  // ✅ ต้องใช้ path ตาม module ของโปรเจกต์ (ดูใน go.mod)
	"auth-service/models"

	"gorm.io/gorm"
)

// 🔹 ดึงข้อมูลผู้ใช้จาก email พร้อม preload role
func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := config.DB.Preload("Role").Where("email = ?", email).First(&user).Error
	return user, err
}

// 🔹 ดึงข้อมูลผู้ใช้จาก ID
func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// 🔹 ตรวจสอบว่า email ซ้ำหรือไม่
func IsEmailExists(email string) bool {
	var count int64
	config.DB.Model(&models.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

// 🔹 สร้างผู้ใช้ใหม่
func CreateUser(user *models.User) error {
	return config.DB.Create(user).Error
}

// 🔹 อัปเดตข้อมูลผู้ใช้
func UpdateUser(user *models.User) error {
	return config.DB.Save(user).Error
}

// 🔹 ปิดบัญชี (soft delete)
func DeactivateUserAccount(userID uint) error {
	var user models.User

	// ตรวจสอบว่าผู้ใช้มีอยู่จริงไหม
	if err := config.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("ไม่พบผู้ใช้")
		}
		return err
	}

	// Soft delete (mark deleted_at)
	return config.DB.Delete(&user).Error
}

// 🔹 ลบผู้ใช้ (soft delete)
func DeleteUserByID(id uint) error {
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("ไม่พบผู้ใช้ที่ระบุ")
		}
		return err
	}
	return config.DB.Delete(&user).Error
}

// 🔹 ดึงรายชื่อผู้ใช้ที่ถูกลบแล้ว (เฉพาะ admin ใช้)
func GetDeletedUsers() ([]models.User, error) {
	var users []models.User
	err := config.DB.Unscoped().Where("deleted_at IS NOT NULL").Find(&users).Error
	return users, err
}

// 🔹 กู้คืนผู้ใช้ที่ถูกลบ
func RestoreUserByID(id uint) error {
	return config.DB.Unscoped().
		Model(&models.User{}).
		Where("id = ?", id). // ✅ เปลี่ยนจาก user_id → id ให้ตรงกับ struct
		Update("deleted_at", nil).Error
}
