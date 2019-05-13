package handler

import (
	"net/http"
	"time"
)

type contextKey string

const (
	// authenticatedUserCookie is users key id in context
	authenticatedUserCookie = "auth.user"
)

// NewCookieAuth creates instance of CookieAuth
func NewCookieAuth(expiry time.Duration) CookieAuth {
	return CookieAuth{
		name:   authenticatedUserCookie,
		expiry: expiry,
		length: 10,
	}
}

// CookieAuth is simplistic unsafe way to identify user based on cookie.
// it implements http.Handler
type CookieAuth struct {
	expiry time.Duration
	name   string
	length int
}

// CookieAuth Simple cookie identifier to set cookie to identify user, no auth at the moment
func (c *CookieAuth) CookieAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := ""
		userCookie, err := r.Cookie(c.name)
		if err == nil {
			user = userCookie.Value
		}
		newUserCookie := c.generateCookie(user)
		http.SetCookie(w, newUserCookie)
		r.AddCookie(newUserCookie)

		next.ServeHTTP(w, r)
	})
}

func (c *CookieAuth) generateCookie(value string) *http.Cookie {
	if value == "" {
		value = randStringRunes(c.length)
	}
	expires := time.Now().Add(c.expiry)
	return &http.Cookie{
		Name:    c.name,
		Value:   value,
		Path:    "/",
		Expires: expires,
	}
}
