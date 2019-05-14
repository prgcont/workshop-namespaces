package nshandler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/prgcont/workshop-namespaces/pkg/nshandler"
)

func TestCookieAuth(t *testing.T) {
	tt := []struct {
		description         string
		initialCookie       *http.Cookie
		expectedCookieValue string
	}{
		{
			description:         "Cookie is not set",
			initialCookie:       nil,
			expectedCookieValue: "", // unknown, generation not tested
		},
		{
			description: "Cookie is set",
			initialCookie: &http.Cookie{
				Name:    "auth.user",
				Value:   "dummy",
				Path:    "/",
				Expires: time.Now().Add(time.Hour * 2),
			},
			expectedCookieValue: "dummy",
		},
	}

	cookieName := "auth.user"

	for _, table := range tt {
		t.Run(table.description, func(r *testing.T) {
			runAssert := assert.New(r)
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestCookie, err := r.Cookie(cookieName)
				requestCookieValue := requestCookie.Value
				runAssert.NoErrorf(err, "Cookie %s is not set", cookieName)
				if len(requestCookieValue) == 0 {
					runAssert.Failf("Value of cookie %s must be larger than 0", cookieName)
				}

				if table.expectedCookieValue != "" {
					runAssert.Equalf(requestCookieValue, table.expectedCookieValue, "Unexpected value of cookie %s", cookieName)
				}
			})

			req, err := http.NewRequest(
				"POST",
				"/namespace",
				nil,
			)

			runAssert.NoError(err, "Test Request can't be created")
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			if table.initialCookie != nil {
				req.AddCookie(table.initialCookie)
			}

			rr := httptest.NewRecorder()
			h := nshandler.CookieAuth(testHandler)
			h.ServeHTTP(rr, req)
		})
	}
}
