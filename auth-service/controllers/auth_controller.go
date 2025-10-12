package controllers

import (
	"encoding/json"
	"net/http"
	"auth-service/models"
	"auth-service/utils"
	"auth-service/validators"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input validators.RegisterInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		user := models.User{
			Name:     input.Name,
			LastName: input.LastName,
			Email:    input.Email,
			Password: string(hashedPassword),
			RoleID:   1, // default role
		}

		if err := db.Create(&user).Error; err != nil {
			http.Error(w, "Email already exists", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "User registered successfully",
		})
	}
}

func Login(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds validators.LoginInput
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		var user models.User
		if err := db.Where("email = ?", creds.Email).First(&user).Error; err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// ตรวจ password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := utils.GenerateJWT(user.UserID)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		utils.SetCookie(w, token)

		json.NewEncoder(w).Encode(map[string]string{
			"token": token,
			"message": "Login successful",
		})
	}
}

func GetUserProfile(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := utils.GetUserIDFromRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		user.Password = "" // อย่าส่ง password กลับ
		json.NewEncoder(w).Encode(user)
	}
}
