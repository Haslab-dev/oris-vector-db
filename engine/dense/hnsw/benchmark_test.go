package hnsw

import (
	"math/rand"
	"testing"

	"github.com/hasdev/oris/engine/dense"
)

var sinkResult []dense.Result

func BenchmarkHNSWInsert(b *testing.B) {
	h := newBenchHNSW()
	rng := rand.New(rand.NewSource(42))
	dim := 128

	vecs := make([][]float32, b.N)
	for i := 0; i < b.N; i++ {
		vec := make([]float32, dim)
		for j := 0; j < dim; j++ {
			vec[j] = rng.Float32()
		}
		vecs[i] = vec
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Insert(uint64(i), vecs[i])
	}
}

func BenchmarkHNSWSearch(b *testing.B) {
	h := newBenchHNSW()
	rng := rand.New(rand.NewSource(42))
	dim := 128

	for i := 0; i < 10000; i++ {
		vec := make([]float32, dim)
		for j := 0; j < dim; j++ {
			vec[j] = rng.Float32()
		}
		h.Insert(uint64(i), vec)
	}

	query := make([]float32, dim)
	for i := range query {
		query[i] = rng.Float32()
	}

	b.ResetTimer()
	var r []dense.Result
	for i := 0; i < b.N; i++ {
		r, _ = h.Search(query, 10)
	}
	sinkResult = r
}

func newBenchHNSW() *HNSW {
	return New(dense.Config{
		Dimension:      128,
		Distance:       "cosine",
		M:              16,
		EfConstruction: 200,
		EfSearch:       100,
	})
}
