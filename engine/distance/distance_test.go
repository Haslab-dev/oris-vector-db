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

func TestTypeString(t *testing.T) {
	assert.Equal(t, "cosine", Cosine.String())
	assert.Equal(t, "dot", Dot.String())
	assert.Equal(t, "euclidean", Euclidean.String())
	assert.Equal(t, "unknown", Type(99).String())
}

func TestCPUDetection(t *testing.T) {
	// CPUInfo should be populated by init().
	// We can't assert exact features (depends on build machine),
	// but we can verify the fields are accessible.
	t.Logf("NEON: %v, AVX2: %v, AVX512F: %v, AVX512BW: %v, SVE: %v",
		CPUInfo.HasNEON, CPUInfo.HasAVX2, CPUInfo.HasAVX512F, CPUInfo.HasAVX512BW, CPUInfo.HasSVE)
}

func TestNewKernelReturnsKernel(t *testing.T) {
	k := NewKernel()
	assert.NotNil(t, k)
	// Should be able to do basic calculations.
	a := []float32{1, 0, 0}
	b := []float32{0, 1, 0}
	assert.InDelta(t, 1, k.Cosine(a, b), 1e-6)
}

func TestBackendRegistration(t *testing.T) {
	// Register a test backend and verify it's selected.
	RegisterBackend(Backend{
		Name: "test",
		Kernel: &testKernel{},
		Checker: func() bool { return true },
	})

	k := NewKernel()
	tk, ok := k.(*testKernel)
	assert.True(t, ok, "NewKernel should return the most recently registered valid backend")
	_ = tk
}

type testKernel struct{}

func (k *testKernel) Cosine(a, b []float32) float32 { return 0 }
func (k *testKernel) Dot(a, b []float32) float32    { return 0 }
func (k *testKernel) Euclidean(a, b []float32) float32 { return 0 }

