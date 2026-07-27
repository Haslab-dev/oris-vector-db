package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkMarshalPoint(b *testing.B) {
	p := &Point{
		ID:    42,
		Dense: make([]float32, 128),
		Sparse: SparseVector{
			Indices: []uint32{0, 5, 10, 15, 20},
			Values:  []float32{0.1, 0.2, 0.3, 0.4, 0.5},
		},
		Payload: []byte("hello world this is a test payload for benchmarking"),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := MarshalPoint(p)
		require.NoError(b, err)
	}
}

func BenchmarkUnmarshalPoint(b *testing.B) {
	p := &Point{
		ID:    42,
		Dense: make([]float32, 128),
		Sparse: SparseVector{
			Indices: []uint32{0, 5, 10, 15, 20},
			Values:  []float32{0.1, 0.2, 0.3, 0.4, 0.5},
		},
		Payload: []byte("hello world this is a test payload for benchmarking"),
	}
	data, err := MarshalPoint(p)
	require.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := UnmarshalPoint(data)
		require.NoError(b, err)
	}
}

func BenchmarkSegmentMeta(b *testing.B) {
	m := MakeSegmentMeta(42, 10000, 1, 9999, 128)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data, err := MarshalSegmentMeta(m)
		if err != nil {
			b.Fatal(err)
		}
		_, err = UnmarshalSegmentMeta(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}
