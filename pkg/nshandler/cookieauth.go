package nshandler

import (
	"net/http"
	"time"
)

type contextKey string

const (
	// authenticatedUserCookie is users key id in context
	authenticatedUserCookie       = "auth.user"
	authenticatedUserCookieLength = 10
	cookieExpiry                  = time.Hour * 5
)

// CookieAuth creates middleware for ensuring user cookie is present
func CookieAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := ""
		userCookie, err := r.Cookie(authenticatedUserCookie)
		if err == nil {
			user = userCookie.Value
		}
		newUserCookie := generateCookie(authenticatedUserCookie, user, cookieExpiry)
		http.SetCookie(w, newUserCookie)
		r.AddCookie(newUserCookie)

		next.ServeHTTP(w, r)
	})
}

func generateCookie(key, value string, expiry time.Duration) *http.Cookie {
	if value == "" {
		value = randStringRunes(authenticatedUserCookieLength)
	}
	expires := time.Now().Add(expiry)
	return &http.Cookie{
		Name:    authenticatedUserCookie,
		Value:   value,
		Path:    "/",
		Expires: expires,
	}
}
