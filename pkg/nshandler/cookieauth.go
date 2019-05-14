package nshandler

import (
	"net/http"
	"time"
)

// NewCookieMiddleware returns Middleware for generating cookie to request
func NewCookieMiddleware(key string, generatedLength int, expiry time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := ""
			userCookie, err := r.Cookie(key)
			if err == nil {
				user = userCookie.Value
			}
			newUserCookie := generateCookie(key, user, generatedLength, expiry)
			http.SetCookie(w, newUserCookie)
			r.AddCookie(newUserCookie)

			next.ServeHTTP(w, r)
		})
	}
}

func generateCookie(key, value string, length int, expiry time.Duration) *http.Cookie {
	if value == "" {
		value = randStringRunes(length)
	}
	expires := time.Now().Add(expiry)
	return &http.Cookie{
		Name:    key,
		Value:   value,
		Path:    "/",
		Expires: expires,
	}
}
