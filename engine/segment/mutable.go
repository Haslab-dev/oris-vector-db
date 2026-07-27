package segment

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/RoaringBitmap/roaring"

	"github.com/hasdev/oris/engine/dense"
	"github.com/hasdev/oris/engine/distance"
	"github.com/hasdev/oris/engine/metadata"
)

// MutableSegment is a segment accepting writes in memory.
type MutableSegment struct {
	mu       sync.RWMutex
	meta     Metadata
	vectors  map[uint64][]float32
	sparse   map[uint64][]string
	payloads map[uint64][]byte
	metaIdx  *metadata.Engine
	dim      int
	dist     distance.Type
}

// NewMutable creates a new mutable segment.
func NewMutable(dim int, dist string) *MutableSegment {
	var dt distance.Type
	switch dist {
	case "dot":
		dt = distance.Dot
	case "euclidean":
		dt = distance.Euclidean
	default:
		dt = distance.Cosine
	}

	return &MutableSegment{
		meta: Metadata{
			ID:        allocID(),
			State:     StateMutable,
			CreatedAt: time.Now(),
			DenseDim:  dim,
			Distance:  dist,
		},
		vectors:  make(map[uint64][]float32),
		sparse:   make(map[uint64][]string),
		payloads: make(map[uint64][]byte),
		metaIdx:  metadata.New(),
		dim:      dim,
		dist:     dt,
	}
}

func (s *MutableSegment) ID() uint64        { return s.meta.ID }
func (s *MutableSegment) State() State       { return s.meta.State }
func (s *MutableSegment) Metadata() Metadata { return s.meta }

func (s *MutableSegment) Insert(id uint64, denseVec []float32, sparseIndices []uint32, sparseValues []float32, payload []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.meta.State != StateMutable {
		return ErrSegmentSealed
	}
	if denseVec == nil {
		denseVec = make([]float32, s.dim)
	}
	s.vectors[id] = denseVec
	s.payloads[id] = payload
	s.metaIdx.IndexField(id, "_exists", "true")

	if s.meta.DocCount == 0 || id < s.meta.MinID {
		s.meta.MinID = id
	}
	if id > s.meta.MaxID {
		s.meta.MaxID = id
	}
	s.meta.DocCount++
	_ = sparseIndices
	_ = sparseValues
	return nil
}

func (s *MutableSegment) Delete(id uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.vectors, id)
	delete(s.sparse, id)
	delete(s.payloads, id)
	s.metaIdx.RemoveField(id, "_exists", "true")
	s.meta.DocCount--
	return nil
}

func (s *MutableSegment) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.meta.DocCount
}

// Search searches within this mutable segment using brute-force.
func (s *MutableSegment) Search(query []float32, topK int, filter *metadata.Filter) ([]dense.Result, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.meta.DocCount == 0 {
		return nil, nil
	}

	var allowed *roaring.Bitmap
	if filter != nil {
		var err error
		allowed, err = s.metaIdx.Evaluate(*filter)
		if err != nil {
			return nil, err
		}
	}

	kern := distance.NewKernel()
	type scored struct {
		id    uint64
		score float32
	}
	var results []scored
	worstScore := float32(-1)
	worstIdx := -1

	for id, vec := range s.vectors {
		if allowed != nil && !allowed.Contains(uint32(id)) {
			continue
		}
		var score float32
		switch s.dist {
		case distance.Dot:
			score = 1 - kern.Dot(query, vec)
		case distance.Euclidean:
			score = kern.Euclidean(query, vec)
		default:
			score = kern.Cosine(query, vec)
		}
		if len(results) < topK {
			results = append(results, scored{id: id, score: score})
			if score > worstScore {
				worstScore = score
				worstIdx = len(results) - 1
			}
		} else if score < worstScore {
			results[worstIdx] = scored{id: id, score: score}
			worstScore = score
			worstIdx = 0
			for i, r := range results {
				if r.score > worstScore {
					worstScore = r.score
					worstIdx = i
				}
			}
		}
	}

	out := make([]dense.Result, len(results))
	for i, r := range results {
		out[i] = dense.Result{ID: r.id, Score: r.score}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Score < out[j].Score
	})
	return out, nil
}

func (s *MutableSegment) Seal(dir string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.meta.State != StateMutable {
		return ErrSegmentSealed
	}
	s.meta.State = StateSealing
	segDir := filepath.Join(dir, fmt.Sprintf("segment_%06d", s.meta.ID))
	if err := os.MkdirAll(segDir, 0755); err != nil {
		return err
	}

	// Write payloads.
	f, err := os.Create(filepath.Join(segDir, "payload.dat"))
	if err != nil {
		return err
	}
	defer f.Close()

	for id, payload := range s.payloads {
		buf := make([]byte, 8+4+len(payload))
		binary.LittleEndian.PutUint64(buf, id)
		binary.LittleEndian.PutUint32(buf[8:], uint32(len(payload)))
		copy(buf[12:], payload)
		if _, err := f.Write(buf); err != nil {
			return err
		}
	}
	f.Close()

	// Write vectors.
	vf, err := os.Create(filepath.Join(segDir, "dense.idx"))
	if err != nil {
		return err
	}
	defer vf.Close()

	for id, vec := range s.vectors {
		buf := make([]byte, 8+4+len(vec)*4)
		binary.LittleEndian.PutUint64(buf, id)
		binary.LittleEndian.PutUint32(buf[8:], uint32(len(vec)))
		for i, v := range vec {
			binary.LittleEndian.PutUint32(buf[12+i*4:], math.Float32bits(v))
		}
		if _, err := vf.Write(buf); err != nil {
			return err
		}
	}
	vf.Close()
	s.meta.State = StateImmutable
	s.meta.SealedAt = time.Now()
	return nil
}

func (s *MutableSegment) Close() error { return nil }

// ensure interface compliance.
var _ Segment = (*MutableSegment)(nil)
