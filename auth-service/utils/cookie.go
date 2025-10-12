// utils/cookie.go
package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ชื่อคุ้กกี้ที่ใช้เก็บ JWT
const AuthCookieName = "token"

// SetAuthCookie เซ็ต JWT ลงคุ้กกี้แบบปลอดภัย
// - secure: true บังคับส่งแค่ผ่าน HTTPS (ถ้ารัน http://localhost ให้ส่งเป็น false)
// - domain: ใช้ "" ระหว่าง dev หรือตั้งเป็นโดเมนจริงตอนโปรดักชัน เช่น "example.com"

func SetCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   int(time.Hour.Seconds() * 24),
	})
}

// GetAuthCookie ดึงค่า JWT จากคุ้กกี้
func GetAuthCookie(c *gin.Context) (string, error) {
	return c.Cookie(AuthCookieName)
}

// ClearAuthCookie ลบคุ้กกี้ทิ้ง (ตอน logout)
func ClearAuthCookie(c *gin.Context, secure bool, domain string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		AuthCookieName,
		"",
		-1,   // maxAge < 0 = ลบทันที
		"/",
		domain,
		secure,
		true, // httpOnly
	)
}
