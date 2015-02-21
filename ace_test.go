package ace

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPMethod(t *testing.T) {
	assert := assert.New(t)

	a := Default()
	a.GET("/", func(c *C) {
		c.String(200, "Test")
	})

	a.POST("/", func(c *C) {
		c.String(200, c.Request.FormValue("test"))
	})

	a.PUT("/", func(c *C) {
		c.String(200, c.Request.FormValue("test"))
	})

	a.PATCH("/", func(c *C) {
		c.String(200, c.Request.FormValue("test"))
	})

	a.DELETE("/", func(c *C) {
		c.String(200, "deleted")
	})

	a.OPTIONS("/", func(c *C) {
		c.String(200, "options")
	})

	a.HEAD("/", func(c *C) {
		c.String(200, "head")
	})

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("Test", w.Body.String())

	r, _ = http.NewRequest("POST", "/", nil)
	r.ParseForm()
	r.Form.Add("test", "hello")
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("hello", w.Body.String())

	r, _ = http.NewRequest("PUT", "/", nil)
	r.ParseForm()
	r.Form.Add("test", "hello")
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("hello", w.Body.String())

	r, _ = http.NewRequest("PATCH", "/", nil)
	r.ParseForm()
	r.Form.Add("test", "hello")
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("hello", w.Body.String())

	r, _ = http.NewRequest("DELETE", "/", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("deleted", w.Body.String())

	r, _ = http.NewRequest("OPTIONS", "/", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("options", w.Body.String())

	r, _ = http.NewRequest("HEAD", "/", nil)
	w = httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("head", w.Body.String())
}
