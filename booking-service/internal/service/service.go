package services

import (
	"booking-service/internal/validators"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func FetchUserByID(token string) (*validators.User, error) {
	baseURL := os.Getenv("USER_SERVICE_URL")
	url := fmt.Sprintf("%s/auth/api/me", baseURL)

	fmt.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// เพิ่ม Authorization header ถ้ามี token
	t := strings.TrimSpace(strings.TrimPrefix(token, "Bearer"))
    if t != "" {
        req.Header.Set("Authorization", "Bearer "+t)
    }

	fmt.Println(req.Header)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user validators.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	fmt.Println(user, "user data")

	return &user, nil
}

// ✅ Fetch course from Course Service (with token)
func FetchCourseByID(id uint, token string) (*validators.Course, error) {
	fmt.Println(id, "id course")
	baseURL := os.Getenv("COURSE_SERVICE_URL")
	url := fmt.Sprintf("%s/courses/%d", baseURL, id)
	fmt.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	// เพิ่ม Authorization header ถ้ามี token
	t := strings.TrimSpace(strings.TrimPrefix(token, "Bearer"))
    if t != "" {
        req.Header.Set("Authorization", "Bearer "+t)
    }

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var course validators.Course
	if err := json.NewDecoder(resp.Body).Decode(&course); err != nil {
		return nil, err
	}


	return &course, nil
}

// ✅ Fetch dogs batch (with token)

func FetchDogsByIDs(dogIDs []uint, token string) ([]validators.Dog, error) {
    if len(dogIDs) == 0 { return []validators.Dog{}, nil }

    baseURL := os.Getenv("DOG_SERVICE_URL")
    url := fmt.Sprintf("%s/dogs/batch", baseURL)

    body, _ := json.Marshal(map[string][]uint{"dog_ids": dogIDs})
    req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")

    // กัน “Bearer Bearer …”
    t := strings.TrimSpace(strings.TrimPrefix(token, "Bearer"))
    if t != "" {
        req.Header.Set("Authorization", "Bearer "+t)
    }

    client := &http.Client{Timeout: 5 * time.Second}
    resp, err := client.Do(req)
    if err != nil { return nil, err }
    defer resp.Body.Close()

    b, _ := io.ReadAll(resp.Body)

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("dog service %d: %s", resp.StatusCode, string(b))
    }

    // ลองหลายรูปแบบ (บาง service หุ้ม array)
    var dogs []validators.Dog
    if err := json.Unmarshal(b, &dogs); err != nil {
        var wrap1 struct{ Dogs []validators.Dog `json:"dogs"` }
        var wrap2 struct{ Data []validators.Dog `json:"data"` }
        if err2 := json.Unmarshal(b, &wrap1); err2 == nil && len(wrap1.Dogs) > 0 {
            return wrap1.Dogs, nil
        }
        if err3 := json.Unmarshal(b, &wrap2); err3 == nil && len(wrap2.Data) > 0 {
            return wrap2.Data, nil
        }
        return nil, fmt.Errorf("decode dogs failed: %v; body=%s", err, string(b))
    }
    return dogs, nil
}
