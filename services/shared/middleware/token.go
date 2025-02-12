package middleware

import (
	"log"
	"net/http"
	"time"
)

func SetTokenToCookies(w http.ResponseWriter, name string, token string, durr int) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(time.Duration(durr) * time.Hour),
	}
	http.SetCookie(w, cookie)
	log.Printf("Cookie set: %s=%s", name, token) // Log untuk debugging
}
