package controllers

import (
	"auth-service/internal/models"
	"auth-service/internal/utils"
	"auth-service/internal/validators"
	"net/http"

	"auth-service/internal/repositories"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)


type Handler struct {
	Repo *repositories.Repository
}

// NewHandler สร้าง instance ของ Handler พร้อม inject repository
func NewHandler(repo *repositories.Repository) *Handler {
	return &Handler{Repo: repo}
}


// ---------------------- HANDLERS ----------------------

// RegisterHandler สมัครสมาชิกใหม่
func (h *Handler) RegisterHandler(c *gin.Context) {
	var input validators.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.Repo.IsEmailExists(input.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := &models.User{
		Name:     input.Name,
		LastName: input.LastName,
		Email:    input.Email,
		Password: string(hashedPassword),
		RoleID:   1, // default role
	}

	if err := h.Repo.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful",
	})
}

// LoginHandler เข้าสู่ระบบและสร้าง JWT
func (h *Handler) LoginHandler(c *gin.Context) {
	var input validators.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Repo.GetUserByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateJWT(user.UserID, user.Role.RoleName, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	utils.SetAuthCookie(c, token, false, "")

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

// LogoutHandler ออกจากระบบ
func (h *Handler) LogoutHandler(c *gin.Context) {
	utils.ClearAuthCookie(c, false, "")
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// EditProfileHandler แก้ไขโปรไฟล์ผู้ใช้
func (h *Handler) EditProfileHandler(c *gin.Context) {
	var input validators.EditProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDInterface.(uint)

	user, err := h.Repo.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "models.User not found"})
		return
	}

	changed := false
	if input.Name != nil && *input.Name != user.Name {
		user.Name = *input.Name
		changed = true
	}
	if input.LastName != nil && *input.LastName != user.LastName {
		user.LastName = *input.LastName
		changed = true
	}
	if input.Email != nil && *input.Email != user.Email {
		user.Email = *input.Email
		changed = true
	}

	if !changed {
		c.JSON(http.StatusOK, gin.H{"message": "ไม่มีข้อมูลใดถูกเปลี่ยนแปลง"})
		return
	}

	if err := h.Repo.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// DeleteUserHandler สำหรับ admin ใช้ลบผู้ใช้
func (h *Handler) DeleteUserHandler(c *gin.Context) {
	idParam := c.Param("id")
	idUint64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	idUint := uint(idUint64)

	if err := h.Repo.DeleteUserByID(idUint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ลบผู้ใช้สำเร็จ (soft delete)"})
}

// GetDeletedUsersHandler แสดงผู้ใช้ที่ถูกลบแล้ว
func (h *Handler) GetDeletedUsersHandler(c *gin.Context) {
	users, err := h.Repo.GetDeletedUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถดึงข้อมูลผู้ใช้ที่ถูกลบได้"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// RestoreUserHandler กู้คืนผู้ใช้ที่ถูกลบ
func (h *Handler) RestoreUserHandler(c *gin.Context) {
	idParam := c.Param("id")
	idUint64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	idUint := uint(idUint64)

	if err := h.Repo.RestoreUserByID(idUint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถกู้คืนผู้ใช้ได้"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "กู้คืนผู้ใช้สำเร็จ"})
}

// DeactivateAccountHandler ผู้ใช้ปิดบัญชีตัวเอง
func (h *Handler) DeactivateAccountHandler(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ไม่ได้เข้าสู่ระบบ"})
		return
	}
	uid := userID.(uint)

	if err := h.Repo.DeactivateUserAccount(uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถปิดบัญชีได้"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "บัญชีของคุณถูกปิดเรียบร้อยแล้ว (soft delete)"})
}

// ✅ VerifyTokenHandler สำหรับ dog-service ตรวจสอบ token
func (h *Handler) VerifyTokenHandler(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header required"})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token format"})
		return
	}

	userID, role, email, err := utils.ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"role":    role,
		"email":   email,
	})
}