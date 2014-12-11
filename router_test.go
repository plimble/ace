package ace

import (
	"testing"
)

func BenchmarkCombine(b *testing.B) {
	a := New()
	a.Use(func(c *C) {})
	a.Use(func(c *C) {})
	a.Use(func(c *C) {})

	for i := 0; i < b.N; i++ {
		a.combineHandlers([]HandlerFunc{
			func(c *C) {},
			func(c *C) {},
			func(c *C) {},
		})
	}
}
