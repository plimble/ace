package ace

import (
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSession(t *testing.T) {
	a := New()

	store := sessions.NewCookieStore([]byte("test"))
	a.Session("test", store, nil)

	a.GET("/", func(c *C) {
		c.Session.Set("test", "abc")
		c.String(200, "")
	})

	a.GET("/test", func(c *C) {
		test := c.Session.GetString("test")
		c.String(200, test)
	})

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, req)
	cookie := w.Header().Get("Set-Cookie")

	req, _ = http.NewRequest("GET", "/test", nil)
	req.Header.Set("Cookie", cookie)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, req)

	assert.Equal(t, "abc", w.Body.String())
}
