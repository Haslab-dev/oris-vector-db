//go:build amd64

package distance

// avx512Kernel implements Kernel using AVX-512 SIMD instructions (future).
// Current implementation is pure Go; assembly will replace these when written.
type avx512Kernel struct{}

func (k *avx512Kernel) Cosine(a, b []float32) float32 {
	return cosineAVX512(a, b)
}

func (k *avx512Kernel) Dot(a, b []float32) float32 {
	return dotAVX512(a, b)
}

func (k *avx512Kernel) Euclidean(a, b []float32) float32 {
	return euclideanAVX512(a, b)
}

// Pure Go implementations — replace with assembly when SIMD kernels are written.
func cosineAVX512(a, b []float32) float32 {
	return genericCosine(a, b)
}

func dotAVX512(a, b []float32) float32 {
	return genericDot(a, b)
}

func euclideanAVX512(a, b []float32) float32 {
	return genericEuclidean(a, b)
}
