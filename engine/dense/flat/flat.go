// Package flat implements a brute-force (flat) dense index.
package flat

import (
	"sort"
	"sync"

	"github.com/hasdev/oris/engine/dense"
	"github.com/hasdev/oris/engine/distance"
)

// Flat is a brute-force dense index that scans all vectors.
// Used for small collections where HNSW overhead isn't warranted.
type Flat struct {
	mu       sync.RWMutex
	vectors  map[uint64][]float32
	ids      []uint64
	dim      int
	dist     distance.Kernel
	metric   distance.Type
}

// New creates a new flat index.
func New(cfg dense.Config) *Flat {
	var metric distance.Type
	switch cfg.Distance {
	case "dot":
		metric = distance.Dot
	case "euclidean":
		metric = distance.Euclidean
	default:
		metric = distance.Cosine
	}

	return &Flat{
		vectors: make(map[uint64][]float32),
		ids:     nil,
		dim:     cfg.Dimension,
		dist:    distance.NewKernel(),
		metric:  metric,
	}
}

func (f *Flat) Insert(id uint64, vec []float32) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.vectors[id] = vec
	f.ids = append(f.ids, id)
	return nil
}

func (f *Flat) Delete(id uint64) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.vectors, id)
	for i, v := range f.ids {
		if v == id {
			f.ids = append(f.ids[:i], f.ids[i+1:]...)
			break
		}
	}
	return nil
}

func (f *Flat) Search(query []float32, topK int) ([]dense.Result, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	type scored struct {
		id    uint64
		score float32
	}

	results := make([]scored, 0, topK)
	var worstScore float32 = -1
	var worstIdx int

	for _, id := range f.ids {
		vec := f.vectors[id]
		score := f.distFn(query, vec)

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

	sort.Slice(results, func(i, j int) bool {
		return results[i].score < results[j].score
	})

	out := make([]dense.Result, len(results))
	for i, r := range results {
		out[i] = dense.Result{ID: r.id, Score: r.score}
	}
	return out, nil
}

func (f *Flat) Len() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.vectors)
}

// Rebuild is a no-op for Flat since it has no index structure to rebuild.
func (f *Flat) Rebuild() error {
	return nil
}

func (f *Flat) distFn(a, b []float32) float32 {
	switch f.metric {
	case distance.Dot:
		return 1 - f.dist.Dot(a, b)
	case distance.Euclidean:
		return f.dist.Euclidean(a, b)
	default:
		return f.dist.Cosine(a, b)
	}
}
