package controllers

import (
	"main/models"
	"main/repositories"
	"main/utils"
	"main/validators"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllDogByUserID(c *gin.Context) {

	userIDInterface, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID := userIDInterface.(uint)

	dogs, err := repositories.GetAllDogByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูลสุนัขได้"})
	}
	c.JSON(http.StatusOK, dogs)
}

// func GetDogByID(c *gin.Context) {
// 	idParam := c.Param("id")                            // ได้ string จาก URL
// 	idUint64, err := strconv.ParseUint(idParam, 10, 64) // แปลงเป็น uint64

// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
// 		return
// 	}

// 	idUint := uint(idUint64) // แปลง uint64 → uint

// 	dog, err := repositories.GetDogByID(idUint)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Dog not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, res)
// }

func CreateDog(c *gin.Context) {
	// validate the request body add dog
	var input validators.CreateDog
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDInterface, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "เป็นคนไม่มีสิทธิ์"})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	dog := models.Dog{
		Name:   input.Name,
		Weight: input.Weight,
		Gender: input.Gender,
		Breed:  input.Breed,
		UserID: userID, // เซ็ตตรงนี้เลย
	}

	if err := repositories.AddDog(&dog); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถเพิ่มหมาได้ แก้โค้ดด่วน"})
		return
	}

	c.JSON(http.StatusCreated, validators.ResponseDog{
		DogID:   dog.DogID,
		Name:    dog.Name,
		Weight:  dog.Weight,
		Breed:   dog.Breed,
		Gender:  dog.Gender,
		Message: "เพิ่มหมาสุดที่รักแล้ว",
	})
}

func UpdateDogByID(c *gin.Context) {
	// 1. แปลง id param เป็น uint
	idParam := c.Param("id")
	idUint64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	idUint := uint(idUint64)

	// 2. Bind JSON เข้ากับ struct validator
	var input validators.UpdateDog
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. ดึงข้อมูลสุนัขเดิมจาก DB
	dog, err := repositories.GetDogByID(idUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบสุนัขที่ระบุ"})
		return
	}

	// 4. ดึง userID จาก context
	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 5. เก็บข้อมูลเดิมไว้เปรียบเทียบ
	originalDog := *dog

	// 6. อัปเดตเฉพาะฟิลด์ที่ผู้ใช้ส่งมาและไม่ใช่ค่า default
	if input.Name != "" {
		dog.Name = input.Name
	}
	if input.Weight != 0 {
		dog.Weight = input.Weight
	}
	if input.Breed != "" {
		dog.Breed = input.Breed
	}

	// 7. เรียก UpdateDogAndCheckOwner (มีการเช็กสิทธิ์ใน repo)
	updated, err := repositories.UpdateDogAndCheckOwner(dog, userID, int64(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถอัปเดตข้อมูลสุนัขได้"})
		return
	}
	if !updated {
		c.JSON(http.StatusForbidden, gin.H{"error": "คุณไม่มีสิทธิ์แก้ไขสุนัขตัวนี้"})
		return
	}

	// 8. ตรวจฟิลด์ที่เปลี่ยนแปลง
	changes := make(map[string]interface{})
	if originalDog.Name != dog.Name {
		changes["name"] = dog.Name
	}
	if originalDog.Weight != dog.Weight {
		changes["weight"] = dog.Weight
	}
	if originalDog.Breed != dog.Breed {
		changes["breed"] = dog.Breed
	}

	// ถ้าไม่มีฟิลด์ไหนเปลี่ยนเลย
	if len(changes) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "ไม่มีข้อมูลที่เปลี่ยนแปลง"})
		return
	}

	// 9. สำเร็จ
	c.JSON(http.StatusOK, gin.H{
		"message": "อัปเดตข้อมูลสุนัขสำเร็จ",
		"updated_fields": changes,
	})
}


/*************  ✨ Windsurf Command ⭐  *************/
// DeleteDogByID เป็นฟังก์ชันสำหรับลบสุนัขโดยระบุ id
// ฟังก์ชันนี้จะตรวจสอบว่าผู้ใช้มีสิทธิ์แก้ไขสุนัขตัวนี้หรือไม่
// ถ้าไม่มีสิทธิ์จะ return status 403 Forbidden

func DeleteDogByID(c *gin.Context) {
	idParam := c.Param("id")
	idUint64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	idUint := uint(idUint64)

	userID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	dog, err := repositories.GetDogByID(idUint)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบสุนัขที่ระบุ"})
		return
	}

	deleted, err := repositories.DeleteDogAndCheckOwner(dog, userID, int64(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถลบข้อมูลสุนัขได้"})
		return
	}
	if !deleted {
		c.JSON(http.StatusForbidden, gin.H{"error": "ไม่มีสิทธิ์ลบสุนัขตัวนี้"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ลบข้อมูลสุนัขสำเร็จ (soft delete)"})
}


