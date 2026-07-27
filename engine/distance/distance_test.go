package distance

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCosine(t *testing.T) {
	k := NewKernel()

	// Same vector should have cosine distance of 0.
	a := []float32{1, 0, 0}
	b := []float32{1, 0, 0}
	assert.InDelta(t, 0, k.Cosine(a, b), 1e-6)

	// Orthogonal vectors should have cosine distance of 1.
	c := []float32{0, 1, 0}
	assert.InDelta(t, 1, k.Cosine(a, c), 1e-6)

	// Opposite vectors should have cosine distance of 2.
	d := []float32{-1, 0, 0}
	assert.InDelta(t, 2, k.Cosine(a, d), 1e-6)
}

func TestDot(t *testing.T) {
	k := NewKernel()

	a := []float32{1, 2, 3}
	b := []float32{4, 5, 6}
	assert.InDelta(t, 32, k.Dot(a, b), 1e-6)

	// Dot with zero vector.
	z := []float32{0, 0, 0}
	assert.InDelta(t, 0, k.Dot(a, z), 1e-6)
}

func TestEuclidean(t *testing.T) {
	k := NewKernel()

	a := []float32{0, 0}
	b := []float32{3, 4}
	assert.InDelta(t, 5, k.Euclidean(a, b), 1e-6)

	// Same point.
	assert.InDelta(t, 0, k.Euclidean(a, a), 1e-6)
}

func TestKernelInterface(t *testing.T) {
	k := NewKernel()
	_, ok := k.(*genericKernel)
	assert.True(t, ok)
}

func TestTypeString(t *testing.T) {
	assert.Equal(t, "cosine", Cosine.String())
	assert.Equal(t, "dot", Dot.String())
	assert.Equal(t, "euclidean", Euclidean.String())
	assert.Equal(t, "unknown", Type(99).String())
}
