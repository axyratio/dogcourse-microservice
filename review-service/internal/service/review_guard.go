// services/review_clients.go
package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

var ErrNotFound = errors.New("not found")
func IsNotFound(err error) bool { return errors.Is(err, ErrNotFound) }

func trimSlash(s string) string {
	return strings.TrimRight(s, "/")
}

func getClient() *http.Client {
	return &http.Client{ Timeout: 5 * time.Second }
}

func doGET(url, authHeader string, out interface{}) (int, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil { return 0, err }
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := getClient().Do(req)
	if err != nil { return 0, err }
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNotFound:
		return resp.StatusCode, ErrNotFound
	}
	if resp.StatusCode >= 400 {
		return resp.StatusCode, fmt.Errorf("remote %s returned %d", url, resp.StatusCode)
	}

	if out != nil {
		_ = json.NewDecoder(resp.Body).Decode(out) // ไม่ซีเรียส body ว่าง
	}
	return resp.StatusCode, nil
}

// 1) เช็กคอร์สมีจริง: GET {COURSE_SERVICE_URL}/courses/:id
func AssertCourseExistsRawID(courseID string, authHeader string) error {
	base := os.Getenv("COURSE_SERVICE_URL")
	if base == "" { return fmt.Errorf("COURSE_SERVICE_URL not set") }
	url := fmt.Sprintf("%s/courses/%s", trimSlash(base), courseID)
	fmt.Println(url, "url courses")
	_, err := doGET(url, authHeader, nil)

	return err
}

// 2) เช็กการจอง: GET {BOOKING_SERVICE_URL}/courses/booking/:id  → {"booked": true/false}
func GetBookedRawID(courseID string, authHeader string) (bool, error) {
	base := os.Getenv("BOOKING_SERVICE_URL")
	if base == "" { return false, fmt.Errorf("BOOKING_SERVICE_URL not set") }
	url := fmt.Sprintf("%s/courses/booking/booked/%s", trimSlash(base), courseID)

	fmt.Println(url, "check booked")
	var resp struct {
		Booked bool `json:"booked"`
	}
	
	code, err := doGET(url, authHeader, &resp)
	if err != nil {
		// 404 = ไม่เคยจอง
		if errors.Is(err, ErrNotFound) {
			return false, nil
		}
		// เผื่อปลายทางตอบ 2xx แต่ body ไม่มี booked → ถือว่า true (ผ่าน guard ที่ปลายทางแล้ว)
		if code >= 200 && code < 300 {
			return true, nil
		}
		return false, err
	}
	fmt.Println(code, err, "status booked andd erorr")
	fmt.Println(resp.Booked, "status booked")

	return resp.Booked, nil
}
