package segment

import (
	"fmt"
	"sync"

	"github.com/hasdev/oris/engine/dense"
)

var nextSegID uint64

func allocID() uint64 {
	nextSegID++
	return nextSegID
}

// SegmentManager coordinates all active segments for a collection.
type SegmentManager struct {
	mu          sync.RWMutex
	mutable     *MutableSegment
	immutables  []*ImmutableSegment
	dir         string
	dim         int
	dist        string
	maxPoints   int
	maxSegments int
}

// NewManager creates a new segment manager.
func NewManager(dir string, dim int, dist string, maxPoints, maxSegments int) *SegmentManager {
	return &SegmentManager{
		mutable:     NewMutable(dim, dist),
		dir:         dir,
		dim:         dim,
		dist:        dist,
		maxPoints:   maxPoints,
		maxSegments: maxSegments,
	}
}

// Insert adds a point to the current mutable segment, auto-sealing if full.
func (sm *SegmentManager) Insert(id uint64, dense []float32, si []uint32, sv []float32, payload []byte) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if err := sm.mutable.Insert(id, dense, si, sv, payload); err != nil {
		return err
	}
	if sm.mutable.Len() >= sm.maxPoints {
		return sm.sealCurrent()
	}
	return nil
}

// Delete removes a point from the mutable segment.
func (sm *SegmentManager) Delete(id uint64) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return sm.mutable.Delete(id)
}

// Len returns the total number of points across all segments.
func (sm *SegmentManager) Len() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	total := sm.mutable.Len()
	for _, seg := range sm.immutables {
		total += seg.Len()
	}
	return total
}

// SearchAcross searches all segments and merges results.
func (sm *SegmentManager) SearchAcross(query []float32, topK int) ([]dense.Result, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	all := make(map[uint64]float32)

	// Search mutable segment.
	mutResults, err := sm.mutable.Search(query, topK, nil)
	if err != nil {
		return nil, err
	}
	for _, r := range mutResults {
		all[r.ID] = r.Score
	}

	// Search immutable segments.
	for _, seg := range sm.immutables {
		results, err := seg.flatIdx.Search(query, topK)
		if err != nil {
			return nil, err
		}
		for _, r := range results {
			if existing, ok := all[r.ID]; ok {
				if r.Score < existing {
					all[r.ID] = r.Score
				}
			} else {
				all[r.ID] = r.Score
			}
		}
	}

	type pair struct {
		id    uint64
		score float32
	}
	pairs := make([]pair, 0, len(all))
	for id, score := range all {
		pairs = append(pairs, pair{id: id, score: score})
	}
	for i := 0; i < len(pairs); i++ {
		for j := i + 1; j < len(pairs); j++ {
			if pairs[j].score < pairs[i].score {
				pairs[i], pairs[j] = pairs[j], pairs[i]
			}
		}
	}
	if len(pairs) > topK {
		pairs = pairs[:topK]
	}

	out := make([]dense.Result, len(pairs))
	for i, p := range pairs {
		out[i] = dense.Result{ID: p.id, Score: p.score}
	}
	return out, nil
}

func (sm *SegmentManager) sealCurrent() error {
	if err := sm.mutable.Seal(sm.dir); err != nil {
		return err
	}
	segDir := fmt.Sprintf("%s/segment_%06d", sm.dir, sm.mutable.meta.ID)
	imm, err := OpenImmutable(segDir, sm.dim, sm.dist)
	if err != nil {
		return err
	}
	sm.immutables = append(sm.immutables, imm)
	sm.mutable = NewMutable(sm.dim, sm.dist)

	if len(sm.immutables) >= sm.maxSegments {
		return sm.compact()
	}
	return nil
}

// compact merges all immutable segments into one.
func (sm *SegmentManager) compact() error {
	if len(sm.immutables) < 2 {
		return nil
	}
	merged := NewMutable(sm.dim, sm.dist)
	for _, seg := range sm.immutables {
		for id, payload := range seg.payloads {
			vec, ok := seg.vectors[id]
			if !ok {
				vec = make([]float32, sm.dim)
			}
			if err := merged.Insert(id, vec, nil, nil, payload); err != nil {
				return err
			}
		}
	}
	if err := merged.Seal(sm.dir); err != nil {
		return err
	}
	segDir := fmt.Sprintf("%s/segment_%06d", sm.dir, merged.meta.ID)
	imm, err := OpenImmutable(segDir, sm.dim, sm.dist)
	if err != nil {
		return err
	}
	sm.immutables = []*ImmutableSegment{imm}
	return nil
}
