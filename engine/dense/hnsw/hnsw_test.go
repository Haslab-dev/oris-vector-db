package hnsw

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hasdev/oris/engine/dense"
)

func newTestHNSW() *HNSW {
	return New(dense.Config{
		Dimension:      8,
		Distance:       "cosine",
		M:              8,
		EfConstruction: 50,
		EfSearch:       50,
	})
}

func TestInsertAndSearch(t *testing.T) {
	h := newTestHNSW()

	for i := 0; i < 20; i++ {
		vec := make([]float32, 8)
		vec[0] = float32(i) / 20.0
		vec[1] = 1.0 - float32(i)/20.0
		err := h.Insert(uint64(i), vec)
		require.NoError(t, err)
	}

	assert.Equal(t, 20, h.Len())

	query := []float32{0.9, 0.1, 0, 0, 0, 0, 0, 0}
	results, err := h.Search(query, 3)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(results), 1)

	assert.Less(t, results[0].Score, float32(0.3))
	assert.GreaterOrEqual(t, results[0].Score, float32(0))
}

func TestSearchReturnsCorrectOrder(t *testing.T) {
	h := newTestHNSW()

	vectors := [][]float32{
		{0.9, 0.1, 0, 0, 0, 0, 0, 0},
		{0, 0.9, 0.1, 0, 0, 0, 0, 0},
		{0, 0, 0.9, 0.1, 0, 0, 0, 0},
		{0, 0, 0, 0.9, 0.1, 0, 0, 0},
		{0, 0, 0, 0, 0.9, 0.1, 0, 0},
	}
	for i, v := range vectors {
		require.NoError(t, h.Insert(uint64(i), v))
	}

	query := []float32{0.9, 0.1, 0, 0, 0, 0, 0, 0}
	results, err := h.Search(query, 5)
	require.NoError(t, err)
	require.Len(t, results, 5)

	assert.Equal(t, uint64(0), results[0].ID)
}

func TestInsertAndSearchRandom(t *testing.T) {
	h := newTestHNSW()

	rng := rand.New(rand.NewSource(42))
	const n = 200
	const dim = 8

	vectors := make([][]float32, n)
	for i := 0; i < n; i++ {
		vec := make([]float32, dim)
		for j := 0; j < dim; j++ {
			vec[j] = rng.Float32()*2 - 1
		}
		var norm float32
		for _, v := range vec {
			norm += v * v
		}
		norm = float32(1.0 / sqrt(float64(norm)))
		for j := range vec {
			vec[j] *= norm
		}
		vectors[i] = vec
		require.NoError(t, h.Insert(uint64(i), vec))
	}

	assert.Equal(t, n, h.Len())

	results, err := h.Search(vectors[0], 10)
	require.NoError(t, err)
	require.Len(t, results, 10)

	assert.Equal(t, uint64(0), results[0].ID)
	assert.Less(t, results[0].Score, float32(0.01))
}

func TestDelete(t *testing.T) {
	h := newTestHNSW()

	for i := 0; i < 20; i++ {
		vec := make([]float32, 8)
		vec[i%8] = 0.9
		vec[(i+1)%8] = 0.1
		require.NoError(t, h.Insert(uint64(i), vec))
	}

	require.NoError(t, h.Delete(uint64(0)))
	assert.Equal(t, 19, h.Len())

	query := []float32{0.9, 0.1, 0, 0, 0, 0, 0, 0}
	results, err := h.Search(query, 10)
	require.NoError(t, err)
	for _, r := range results {
		assert.NotEqual(t, uint64(0), r.ID)
	}
}

func TestEmptySearch(t *testing.T) {
	h := newTestHNSW()
	results, err := h.Search([]float32{1, 0, 0, 0, 0, 0, 0, 0}, 5)
	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestLen(t *testing.T) {
	h := newTestHNSW()
	assert.Equal(t, 0, h.Len())

	h.Insert(1, []float32{1, 0, 0, 0, 0, 0, 0, 0})
	assert.Equal(t, 1, h.Len())

	h.Insert(2, []float32{0, 1, 0, 0, 0, 0, 0, 0})
	assert.Equal(t, 2, h.Len())
}

func TestLargeEFSearch(t *testing.T) {
	cfg := dense.Config{
		Dimension:      4,
		Distance:       "cosine",
		M:              8,
		EfConstruction: 200,
		EfSearch:       200,
	}
	h := New(cfg)

	for i := 0; i < 100; i++ {
		vec := []float32{
			float32(i) / 100,
			float32(100-i) / 100,
			0.5,
			0.5,
		}
		require.NoError(t, h.Insert(uint64(i), vec))
	}

	results, err := h.Search([]float32{0.5, 0.5, 0.5, 0.5}, 5)
	require.NoError(t, err)
	require.Len(t, results, 5)
}

func TestRebuild(t *testing.T) {
	h := newTestHNSW()

	for i := 0; i < 50; i++ {
		vec := make([]float32, 8)
		vec[0] = float32(i) / 50.0
		vec[1] = 1.0 - float32(i)/50.0
		require.NoError(t, h.Insert(uint64(i), vec))
	}
	assert.Equal(t, 50, h.Len())

	err := h.Rebuild()
	require.NoError(t, err)
	assert.Equal(t, 50, h.Len())

	query := []float32{0.9, 0.1, 0, 0, 0, 0, 0, 0}
	results, err := h.Search(query, 5)
	require.NoError(t, err)
	require.Len(t, results, 5)
	assert.Less(t, results[0].Score, float32(0.3))
}

func sqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x / 2
	for i := 0; i < 10; i++ {
		z = (z + x/z) / 2
	}
	return z
}
