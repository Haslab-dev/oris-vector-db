package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshalUnmarshalPoint(t *testing.T) {
	p := &Point{
		ID:    42,
		Dense: []float32{1.5, 2.5, 3.5},
		Sparse: SparseVector{
			Indices: []uint32{0, 5, 10},
			Values:  []float32{0.1, 0.2, 0.3},
		},
		Payload: []byte("hello"),
	}

	data, err := MarshalPoint(p)
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	p2, err := UnmarshalPoint(data)
	require.NoError(t, err)

	assert.Equal(t, p.ID, p2.ID)
	assert.Equal(t, p.Dense, p2.Dense)
	assert.Equal(t, p.Sparse.Indices, p2.Sparse.Indices)
	assert.Equal(t, p.Sparse.Values, p2.Sparse.Values)
	assert.Equal(t, p.Payload, p2.Payload)
}

func TestMarshalPointEmptyFields(t *testing.T) {
	p := &Point{
		ID:      1,
		Dense:   []float32{},
		Sparse:  SparseVector{},
		Payload: nil,
	}

	data, err := MarshalPoint(p)
	require.NoError(t, err)

	p2, err := UnmarshalPoint(data)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), p2.ID)
	assert.Empty(t, p2.Dense)
	assert.Empty(t, p2.Sparse.Indices)
	assert.Empty(t, p2.Sparse.Values)
	assert.Empty(t, p2.Payload)
}

func TestUnmarshalInvalidData(t *testing.T) {
	_, err := UnmarshalPoint([]byte{0x01})
	assert.ErrorIs(t, err, ErrInvalidEncoding)

	_, err = UnmarshalPoint(nil)
	assert.ErrorIs(t, err, ErrInvalidEncoding)
}
