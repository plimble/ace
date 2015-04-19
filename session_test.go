package ace

import (
	"github.com/plimble/sessions/store/cookie"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSession(t *testing.T) {
	a := New()

	store := cookie.NewCookieStore()
	a.Session(store, nil)

	a.GET("/", func(c *C) {
		session1 := c.Sessions("test")
		session1.Set("test1", "123")
		session1.Set("test2", 123)

		session2 := c.Sessions("foo")
		session2.Set("baz1", "123")
		session2.Set("baz2", 123)

		c.String(200, "")
	})

	a.GET("/test", func(c *C) {
		session := c.Sessions("test")
		test1 := session.GetString("test1", "")
		test2 := session.GetInt("test2", 0)

		assert.Equal(t, "123", test1)
		assert.Equal(t, 123, test2)
		c.String(200, "")
	})

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, req)
	cookie := w.Header().Get("Set-Cookie")

	req, _ = http.NewRequest("GET", "/test", nil)
	req.Header.Set("Cookie", cookie)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, req)
}
