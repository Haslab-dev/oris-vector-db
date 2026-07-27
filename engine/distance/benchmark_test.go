package distance

import (
	"testing"
)

var sink float32

func BenchmarkCosine(b *testing.B) {
	k := NewKernel()
	a := make([]float32, 128)
	bv := make([]float32, 128)
	for i := range a {
		a[i] = float32(i) / 128
		bv[i] = float32(128-i) / 128
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = k.Cosine(a, bv)
	}
}

func BenchmarkDot(b *testing.B) {
	k := NewKernel()
	a := make([]float32, 128)
	bv := make([]float32, 128)
	for i := range a {
		a[i] = float32(i) / 128
		bv[i] = float32(128-i) / 128
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = k.Dot(a, bv)
	}
}

func BenchmarkEuclidean(b *testing.B) {
	k := NewKernel()
	a := make([]float32, 128)
	bv := make([]float32, 128)
	for i := range a {
		a[i] = float32(i) / 128
		bv[i] = float32(128-i) / 128
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink = k.Euclidean(a, bv)
	}
}
