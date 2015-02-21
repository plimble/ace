package ace

import (
	"testing"
)

func BenchmarkCombineHandlers(b *testing.B) {
	testHandler := func(c *C) {}
	r := Router{}

	r.Use(testHandler,
		testHandler,
		testHandler,
		testHandler,
		testHandler,
		testHandler,
	)

	handlers := []HandlerFunc{
		testHandler,
		testHandler,
		testHandler,
		testHandler,
		testHandler,
		testHandler,
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.combineHandlers(handlers)
	}
}
