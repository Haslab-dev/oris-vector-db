// Package internal provides internal encoding utilities for Oris.
//
// Points, segments, and indexes are serialized using binary encoding
// rather than JSON or Protobuf for maximum performance.
package internal

import (
	"encoding/binary"
	"errors"
	"math"
)

var ErrInvalidEncoding = errors.New("invalid binary encoding")

// Point represents a stored point with its vectors and payload.
type Point struct {
	ID      uint64
	Dense   []float32
	Sparse  SparseVector
	Payload []byte
}

// SparseVector represents a sparse vector using parallel index/value arrays.
type SparseVector struct {
	Indices []uint32
	Values  []float32
}

// MarshalPoint serializes a Point to binary.
func MarshalPoint(p *Point) ([]byte, error) {
	buf := make([]byte, 8+4+len(p.Dense)*4+4+4+len(p.Sparse.Indices)*8+4+len(p.Payload))
	n := 0

	binary.LittleEndian.PutUint64(buf[n:], p.ID)
	n += 8

	// Dense vector.
	binary.LittleEndian.PutUint32(buf[n:], uint32(len(p.Dense)))
	n += 4
	for _, f := range p.Dense {
		binary.LittleEndian.PutUint32(buf[n:], math.Float32bits(f))
		n += 4
	}

	// Sparse vector indices.
	binary.LittleEndian.PutUint32(buf[n:], uint32(len(p.Sparse.Indices)))
	n += 4
	for _, idx := range p.Sparse.Indices {
		binary.LittleEndian.PutUint32(buf[n:], idx)
		n += 4
	}
	// Sparse vector values.
	for _, v := range p.Sparse.Values {
		binary.LittleEndian.PutUint32(buf[n:], math.Float32bits(v))
		n += 4
	}

	// Payload.
	binary.LittleEndian.PutUint32(buf[n:], uint32(len(p.Payload)))
	n += 4
	copy(buf[n:], p.Payload)
	n += len(p.Payload)

	return buf[:n], nil
}

// UnmarshalPoint deserializes a Point from binary.
func UnmarshalPoint(data []byte) (*Point, error) {
	if len(data) < 12 {
		return nil, ErrInvalidEncoding
	}

	p := &Point{}
	n := 0

	p.ID = binary.LittleEndian.Uint64(data[n:])
	n += 8

	// Dense vector.
	denseLen := int(binary.LittleEndian.Uint32(data[n:]))
	n += 4
	if len(data) < n+denseLen*4 {
		return nil, ErrInvalidEncoding
	}
	p.Dense = make([]float32, denseLen)
	for i := range p.Dense {
		p.Dense[i] = math.Float32frombits(binary.LittleEndian.Uint32(data[n:]))
		n += 4
	}

	// Sparse vector indices.
	sparseLen := int(binary.LittleEndian.Uint32(data[n:]))
	n += 4
	if len(data) < n+sparseLen*4 {
		return nil, ErrInvalidEncoding
	}
	p.Sparse.Indices = make([]uint32, sparseLen)
	for i := range p.Sparse.Indices {
		p.Sparse.Indices[i] = binary.LittleEndian.Uint32(data[n:])
		n += 4
	}
	// Sparse vector values.
	if len(data) < n+sparseLen*4 {
		return nil, ErrInvalidEncoding
	}
	p.Sparse.Values = make([]float32, sparseLen)
	for i := range p.Sparse.Values {
		p.Sparse.Values[i] = math.Float32frombits(binary.LittleEndian.Uint32(data[n:]))
		n += 4
	}

	// Payload.
	if len(data) < n+4 {
		return nil, ErrInvalidEncoding
	}
	payloadLen := int(binary.LittleEndian.Uint32(data[n:]))
	n += 4
	if len(data) < n+payloadLen {
		return nil, ErrInvalidEncoding
	}
	p.Payload = make([]byte, payloadLen)
	copy(p.Payload, data[n:])

	return p, nil
}
