package ace

import (
	"bytes"
	"log"
	"os"
	"testing"
)

// TestPanicInHandler assert that panic has been recovered.
func TestPanicInHandler(t *testing.T) {
	// SETUP
	log.SetOutput(bytes.NewBuffer(nil)) // Disable panic logs for testing
	r := New()
	r.Use(Recovery())
	r.GET("/recovery", func(_ *C) {
		panic("Oupps, Houston, we have a problem")
	})

	// RUN
	w := PerformRequest(r, "GET", "/recovery")

	// restore logging
	log.SetOutput(os.Stderr)

	if w.Code != 500 {
		t.Errorf("Response code should be Internal Server Error, was: %s", w.Code)
	}
}

// TestPanicWithAbort assert that panic has been recovered even if context.Abort was used.
func TestPanicWithAbort(t *testing.T) {
	// SETUP
	log.SetOutput(bytes.NewBuffer(nil))
	r := New()
	r.Use(Recovery())
	r.GET("/recovery", func(c *C) {
		c.Abort(400)
		panic("Oupps, Houston, we have a problem")
	})

	// RUN
	w := PerformRequest(r, "GET", "/recovery")

	// restore logging
	log.SetOutput(os.Stderr)

	// TEST
	if w.Code != 400 {
		t.Errorf("Response code should be Bad request, was: %v", w.Code)
	}
}
