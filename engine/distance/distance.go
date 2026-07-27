// Package distance provides similarity/distance computation kernels.
//
// Oris uses these kernels for all similarity computations. The implementation
// auto-selects the fastest available backend based on CPU capabilities.
package distance

import "math"

// Kernel is the interface for computing distances between vectors.
type Kernel interface {
	// Cosine computes the cosine similarity (1 - cosine) distance.
	Cosine(a, b []float32) float32

	// Dot computes the dot product of two vectors.
	Dot(a, b []float32) float32

	// Euclidean computes the Euclidean distance between two vectors.
	Euclidean(a, b []float32) float32
}

// Type represents the distance metric.
type Type int

const (
	Cosine  Type = iota
	Dot          // Dot product
	Euclidean    // Euclidean distance
)

// String returns the string representation of the distance type.
func (t Type) String() string {
	switch t {
	case Cosine:
		return "cosine"
	case Dot:
		return "dot"
	case Euclidean:
		return "euclidean"
	default:
		return "unknown"
	}
}

// NewKernel returns the best available distance kernel for the current CPU.
func NewKernel() Kernel {
	// Start with pure Go; SIMD backends will be auto-selected here later.
	return &genericKernel{}
}

// genericKernel is a pure Go implementation of the Kernel interface.
type genericKernel struct{}

func (k *genericKernel) Cosine(a, b []float32) float32 {
	dot := float32(0)
	normA := float32(0)
	normB := float32(0)
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 1
	}
	return 1 - dot/(sqrt32(normA)*sqrt32(normB))
}

func (k *genericKernel) Dot(a, b []float32) float32 {
	dot := float32(0)
	for i := range a {
		dot += a[i] * b[i]
	}
	return dot
}

func (k *genericKernel) Euclidean(a, b []float32) float32 {
	sum := float32(0)
	for i := range a {
		d := a[i] - b[i]
		sum += d * d
	}
	return sqrt32(sum)
}

func sqrt32(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}

// ensure kernel implements the interface.
var _ Kernel = (*genericKernel)(nil)
