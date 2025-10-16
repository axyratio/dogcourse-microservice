package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// VerifyUserExists ตรวจสอบว่าผู้ใช้มีอยู่ใน users-service หรือไม่
func VerifyToken(token string) (uint, string, string, error) {
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		userServiceURL = "http://localhost:8081" // ค่า default
	}

	// ✅ ตัด Bearer ออกก่อนถ้ามี (ป้องกันซ้ำ)
	token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer"))

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/verify", userServiceURL), nil)
	if err != nil {
		return 0, "", "", fmt.Errorf("failed to create request: %v", err)
	}

	// ✅ ส่ง header ที่ถูกต้อง
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, "", "", fmt.Errorf("cannot connect to users-service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, "", "", fmt.Errorf("invalid token or unauthorized: %s", resp.Status)
	}

	var data struct {
		UserID uint   `json:"user_id"`
		Role   string `json:"role"`
		Email  string `json:"email"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, "", "", fmt.Errorf("invalid response from users-service: %v", err)
	}

	return data.UserID, data.Role, data.Email, nil
}