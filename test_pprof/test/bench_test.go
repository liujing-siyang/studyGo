package test

import (
	"test_pprof/example"
	"testing"
)

func BenchmarkFib(b *testing.B) {
	for i := 0; i < b.N; i++ {
		example.Fib(30)
	}
}
