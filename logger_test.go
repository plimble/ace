package ace

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Logger(t *testing.T) {
	w := httptest.NewRecorder()

	r := New()
	// replace log for testing
	r.Use(Logger())
	r.GET("/", func(c *C) {
		c.AbortWithStatus(404)
	})

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}

	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Status should be not found but got %v", w.Code)
	}
}
