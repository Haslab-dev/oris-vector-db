//go:build arm64

package distance

// neonKernel implements Kernel using ARM NEON SIMD instructions (future).
// Current implementation is pure Go; assembly will replace these when written.
type neonKernel struct{}

func (k *neonKernel) Cosine(a, b []float32) float32 {
	return cosineNEON(a, b)
}

func (k *neonKernel) Dot(a, b []float32) float32 {
	return dotNEON(a, b)
}

func (k *neonKernel) Euclidean(a, b []float32) float32 {
	return euclideanNEON(a, b)
}

// Pure Go implementations — replace with assembly when SIMD kernels are written.
func cosineNEON(a, b []float32) float32 {
	return genericCosine(a, b)
}

func dotNEON(a, b []float32) float32 {
	return genericDot(a, b)
}

func euclideanNEON(a, b []float32) float32 {
	return genericEuclidean(a, b)
}
