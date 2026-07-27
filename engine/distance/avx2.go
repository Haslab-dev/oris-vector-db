//go:build amd64

package distance

// avx2Kernel implements Kernel using AVX2 SIMD instructions (future).
// Current implementation is pure Go; assembly will replace these when written.
type avx2Kernel struct{}

func (k *avx2Kernel) Cosine(a, b []float32) float32 {
	return cosineAVX2(a, b)
}

func (k *avx2Kernel) Dot(a, b []float32) float32 {
	return dotAVX2(a, b)
}

func (k *avx2Kernel) Euclidean(a, b []float32) float32 {
	return euclideanAVX2(a, b)
}

// Pure Go implementations — replace with assembly when SIMD kernels are written.
func cosineAVX2(a, b []float32) float32 {
	return genericCosine(a, b)
}

func dotAVX2(a, b []float32) float32 {
	return genericDot(a, b)
}

func euclideanAVX2(a, b []float32) float32 {
	return genericEuclidean(a, b)
}
