package utils

import (
	"net/http"
)

const AuthCookieName = "jwt"

func SetCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     AuthCookieName,
		Value:    token,
		HttpOnly: true,
		Secure:   false, // ✅ ใช้ true ใน production
		Path:     "/",
		MaxAge:   86400,
	})
}

func ClearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   AuthCookieName,
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})
}
