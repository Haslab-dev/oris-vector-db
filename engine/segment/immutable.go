package segment

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"

	"github.com/hasdev/oris/engine/dense"
	"github.com/hasdev/oris/engine/dense/flat"
)

// ImmutableSegment is a read-only segment stored on disk.
type ImmutableSegment struct {
	mu       sync.RWMutex
	meta     Metadata
	dir      string
	vectors  map[uint64][]float32
	payloads map[uint64][]byte
	flatIdx  *flat.Flat
	dim      int
}

// OpenImmutable opens an immutable segment from disk.
func OpenImmutable(dir string, dim int, dist string) (*ImmutableSegment, error) {
	meta := Metadata{
		State:    StateImmutable,
		DenseDim: dim,
		Distance: dist,
	}
	fmt.Sscanf(filepath.Base(dir), "segment_%d", &meta.ID)

	seg := &ImmutableSegment{
		meta:     meta,
		dir:      dir,
		vectors:  make(map[uint64][]float32),
		payloads: make(map[uint64][]byte),
		dim:      dim,
	}
	seg.flatIdx = flat.New(dense.Config{Dimension: dim, Distance: dist})

	// Load payloads.
	payloadPath := filepath.Join(dir, "payload.dat")
	if payloadData, pErr := os.ReadFile(payloadPath); pErr == nil {
		offset := 0
		for offset+12 <= len(payloadData) {
			id := binary.LittleEndian.Uint64(payloadData[offset:])
			payloadLen := binary.LittleEndian.Uint32(payloadData[offset+8:])
			offset += 12
			end := offset + int(payloadLen)
			if end > len(payloadData) {
				break
			}
			payload := make([]byte, payloadLen)
			copy(payload, payloadData[offset:end])
			seg.payloads[id] = payload
			offset = end
			meta.DocCount++
		}
	}

	// Load vectors and populate flat index.
	vecPath := filepath.Join(dir, "dense.idx")
	if vecData, vErr := os.ReadFile(vecPath); vErr == nil {
		offset := 0
		for offset+12 <= len(vecData) {
			id := binary.LittleEndian.Uint64(vecData[offset:])
			vecLen := int(binary.LittleEndian.Uint32(vecData[offset+8:]))
			offset += 12
			end := offset + vecLen*4
			if end > len(vecData) {
				break
			}
			vec := make([]float32, vecLen)
			for i := 0; i < vecLen; i++ {
				vec[i] = math.Float32frombits(binary.LittleEndian.Uint32(vecData[offset+i*4:]))
			}
			seg.vectors[id] = vec
			seg.flatIdx.Insert(id, vec)
			offset = end
		}
	}

	seg.meta = meta
	return seg, nil
}

func (s *ImmutableSegment) ID() uint64                  { return s.meta.ID }
func (s *ImmutableSegment) State() State                { return s.meta.State }
func (s *ImmutableSegment) Metadata() Metadata          { return s.meta }

func (s *ImmutableSegment) Insert(id uint64, dense []float32, si []uint32, sv []float32, payload []byte) error {
	return ErrSegmentSealed
}

func (s *ImmutableSegment) Delete(id uint64) error {
	return ErrSegmentSealed
}

func (s *ImmutableSegment) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.meta.DocCount
}

func (s *ImmutableSegment) Close() error { return nil }

func (s *ImmutableSegment) Seal(dir string) error {
	return ErrSegmentSealed
}

var _ Segment = (*ImmutableSegment)(nil)
