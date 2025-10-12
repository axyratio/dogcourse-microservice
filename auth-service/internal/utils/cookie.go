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
func SetAuthCookie(c *gin.Context, token string, secure bool, domain string) {
	// แนะนำให้ใช้ Lax ถ้าเป็นเว็บปกติ (ถ้าต้อง cross-site ให้ใช้ None และต้อง secure=true)
	c.SetSameSite(http.SameSiteLaxMode)

	// อายุคุ้กกี้ 24 ชม. (เท่ากับ exp ของ JWT ตัวอย่างของคุณ)
	maxAge := int((24 * time.Hour).Seconds())

	c.SetCookie(
		AuthCookieName, // name
		token,          // value
		maxAge,         // maxAge (seconds)
		"/",            // path
		domain,         // domain
		secure,         // secure (true = ส่งเฉพาะ https)
		true,           // httpOnly
	)
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