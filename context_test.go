package ace

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSONResp(t *testing.T) {
	assert := assert.New(t)

	data := struct {
		s string `json:"s"`
		n int    `json:"n"`
		b bool   `json:"b"`
	}{
		s: "test",
		n: 123,
		b: true,
	}

	a := New()
	a.GET("/", func(c *C) {
		c.JSON(200, data)
	})

	result, _ := json.Marshal(data)

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal(result, w.Body.String())
	assert.Equal("application/json; charset=UTF-8", w.Header().Get("Content-Type"))
}

func TestDownloadResp(t *testing.T) {
	assert := assert.New(t)
	a := New()
	a.GET("/", func(c *C) {
		c.Download(200, []byte("123"))
	})

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("123", w.Body.String())
	assert.Equal("application/octet-stream; charset=UTF-8", w.Header().Get("Content-Type"))
}

func TestCData(t *testing.T) {
	assert := assert.New(t)
	a := New()

	a.Use(func(c *C) {
		c.SetData("test", "123")
		c.Next()
	})

	a.GET("/", func(c *C) {
		c.GetAllData()
		c.String(200, c.GetData("test").(string))
	})

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	a.ServeHTTP(w, r)
	assert.Equal(200, w.Code)
	assert.Equal("123", w.Body.String())
}
