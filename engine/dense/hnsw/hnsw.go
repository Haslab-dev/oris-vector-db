// Package hnsw implements the Hierarchical Navigable Small World graph.
package hnsw

import (
	"container/heap"
	"math"
	"math/rand"
	"sync"

	"github.com/hasdev/oris/engine/dense"
	"github.com/hasdev/oris/engine/distance"
)

// HNSW is a Hierarchical Navigable Small World graph index.
type HNSW struct {
	mu       sync.RWMutex
	nodes    map[uint64]*node
	vectors  map[uint64][]float32
	entry    uint64             // entry point ID
	maxLevel int                // current max level
	M        int                // number of connections per layer
	Mmax     int                // max connections (2*M for upper layers)
	efSearch int                // ef for search
	efConst  int                // ef for construction
	dim      int                // vector dimension
	dist     distance.Kernel
	metric   distance.Type
	rng      *rand.Rand
	levelMult float64           // 1/ln(M)
}

type node struct {
	id     uint64
	level  int
	edges  map[int][]uint64 // layer -> neighbor IDs
}

// New creates a new HNSW index with the given config.
func New(cfg dense.Config) *HNSW {
	M := cfg.M
	if M <= 0 {
		M = 16
	}
	efConst := cfg.EfConstruction
	if efConst <= 0 {
		efConst = 200
	}
	efSearch := cfg.EfSearch
	if efSearch <= 0 {
		efSearch = 100
	}

	var metric distance.Type
	switch cfg.Distance {
	case "dot":
		metric = distance.Dot
	case "euclidean":
		metric = distance.Euclidean
	default:
		metric = distance.Cosine
	}

	return &HNSW{
		nodes:     make(map[uint64]*node),
		vectors:   make(map[uint64][]float32),
		entry:     0,
		maxLevel:  -1,
		M:         M,
		Mmax:      M, // for layer 0, use M; for upper layers, use M
		efSearch:  efSearch,
		efConst:   efConst,
		dim:       cfg.Dimension,
		dist:      distance.NewKernel(),
		metric:    metric,
		rng:       rand.New(rand.NewSource(42)),
		levelMult: 1.0 / math.Log(float64(M)),
	}
}

func (h *HNSW) Insert(id uint64, vec []float32) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.vectors[id] = vec
	level := h.randomLevel()

	n := &node{id: id, level: level, edges: make(map[int][]uint64)}
	h.nodes[id] = n

	if h.maxLevel < 0 {
		// First node.
		h.entry = id
		h.maxLevel = level
		return nil
	}

	currEntry := h.nodes[h.entry]

	// Phase 1: traverse from top level down to the new node's level.
	// At each level, find the single nearest neighbor.
	ep := currEntry
	for l := h.maxLevel; l > level; l-- {
		ep = h.selectNearestSimple(ep, vec, l)
	}

	// Phase 2: for levels from min(level, maxLevel) down to 0,
	// find efConstruction nearest neighbors and connect.
	for l := min(level, h.maxLevel); l >= 0; l-- {
		visited := h.searchLayer(ep, vec, h.efConst, l)
		neighbors := h.selectNeighbors(visited, vec, h.M, l)

		// Connect the new node to the selected neighbors.
		n.edges[l] = neighbors
		for _, nid := range neighbors {
			neighbor := h.nodes[nid]
			neighbor.edges[l] = append(neighbor.edges[l], id)
			// Shrink if needed.
			if len(neighbor.edges[l]) > h.Mmax {
				neighbor.edges[l] = h.shrinkConnections(neighbor, l)
			}
		}
	}

	// Update entry point if needed.
	if level > h.maxLevel {
		h.entry = id
		h.maxLevel = level
	}

	return nil
}

func (h *HNSW) Delete(id uint64) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.nodes[id]; !ok {
		return nil
	}

	// Remove this node from all remaining nodes' edge lists.
	for _, node := range h.nodes {
		if node.id == id {
			continue
		}
		for l := 0; l <= node.level; l++ {
			node.edges[l] = removeID(node.edges[l], id)
		}
	}

	delete(h.nodes, id)
	delete(h.vectors, id)

	// Pick a new entry point if needed.
	if h.entry == id {
		var newEntry uint64
		for nid := range h.nodes {
			newEntry = nid
			break
		}
		if newEntry == 0 && h.nodes[0] == nil {
			h.entry = 0
			h.maxLevel = -1
			return nil
		}
		h.entry = newEntry
	}

	return nil
}

func (h *HNSW) Search(query []float32, topK int) ([]dense.Result, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.maxLevel < 0 || len(h.nodes) == 0 {
		return nil, nil
	}

	ep := h.nodes[h.entry]

	// Phase 1: greedy search from top level down to level 0.
	for l := h.maxLevel; l > 0; l-- {
		ep = h.selectNearestSimple(ep, query, l)
	}

	// Phase 2: search layer 0 with efSearch.
	visited := h.searchLayer(ep, query, h.efSearch, 0)

	// Sort visited by distance and pick topK.
	type scored struct {
		id    uint64
		dist  float32
	}
	scores := make([]scored, 0, len(visited))
	for _, nid := range visited {
		scores = append(scores, scored{id: nid, dist: h.distFn(h.vectors[nid], query)})
	}
	for i := 0; i < len(scores); i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[j].dist < scores[i].dist {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}
	if len(scores) > topK {
		scores = scores[:topK]
	}

	results := make([]dense.Result, len(scores))
	for i, s := range scores {
		results[i] = dense.Result{ID: s.id, Score: s.dist}
	}
	return results, nil
}

func (h *HNSW) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.nodes)
}

// dist computes the distance between a stored vector and a query.
func (h *HNSW) distFn(a, b []float32) float32 {
	switch h.metric {
	case distance.Dot:
		return 1 - h.dist.Dot(a, b) // convert to distance (lower = closer)
	case distance.Euclidean:
		return h.dist.Euclidean(a, b)
	default:
		return h.dist.Cosine(a, b)
	}
}

// searchLayer finds the ef nearest neighbors to query on the given layer,
// starting from the entry point. Returns the closest items in arbitrary order.
func (h *HNSW) searchLayer(ep *node, query []float32, ef, layer int) []uint64 {
	visited := make(map[uint64]struct{})

	// candidates: min-heap (closest on top).
	candidates := make(minHeap, 0)
	// results: max-heap (farthest on top).
	results := make(maxHeap, 0)

	epDist := h.distFn(h.vectors[ep.id], query)
	heap.Push(&candidates, &distItem{id: ep.id, dist: epDist})
	heap.Push(&results, &distItem{id: ep.id, dist: epDist})
	visited[ep.id] = struct{}{}

	for candidates.Len() > 0 {
		cand := heap.Pop(&candidates).(*distItem)

		// Stop if the closest candidate is farther than the farthest result.
		if cand.dist > results[0].dist {
			break
		}

		for _, nid := range h.nodes[cand.id].edges[layer] {
			if _, seen := visited[nid]; seen {
				continue
			}
			visited[nid] = struct{}{}

			d := h.distFn(h.vectors[nid], query)
			if len(results) < ef || d < results[0].dist {
				heap.Push(&candidates, &distItem{id: nid, dist: d})
				heap.Push(&results, &distItem{id: nid, dist: d})
				if len(results) > ef {
					heap.Pop(&results)
				}
			}
		}
	}

	// Collect results into a slice.
	out := make([]uint64, len(results))
	for i, item := range results {
		out[i] = item.id
	}
	return out
}

// selectNearestSimple finds the single nearest neighbor via greedy search.
func (h *HNSW) selectNearestSimple(ep *node, query []float32, layer int) *node {
	best := ep
	bestDist := h.distFn(h.vectors[ep.id], query)

	visited := make(map[uint64]struct{})
	visited[ep.id] = struct{}{}

	for {
		improved := false
		for _, nid := range best.edges[layer] {
			if _, seen := visited[nid]; seen {
				continue
			}
			visited[nid] = struct{}{}
			d := h.distFn(h.vectors[nid], query)
			if d < bestDist {
				bestDist = d
				best = h.nodes[nid]
				improved = true
			}
		}
		if !improved {
			break
		}
	}

	return best
}

// selectNeighbors picks the best M neighbors from a candidate list.
func (h *HNSW) selectNeighbors(candidates []uint64, query []float32, M, layer int) []uint64 {
	type scored struct {
		id   uint64
		dist float32
	}
	all := make([]scored, len(candidates))
	for i, nid := range candidates {
		all[i] = scored{id: nid, dist: h.distFn(h.vectors[nid], query)}
	}
	for i := 0; i < len(all); i++ {
		for j := i + 1; j < len(all); j++ {
			if all[j].dist < all[i].dist {
				all[i], all[j] = all[j], all[i]
			}
		}
	}
	n := min(M, len(all))
	out := make([]uint64, n)
	for i := 0; i < n; i++ {
		out[i] = all[i].id
	}
	return out
}

// shrinkConnections reduces the number of connections for a node at a layer.
func (h *HNSW) shrinkConnections(n *node, layer int) []uint64 {
	edges := n.edges[layer]
	if len(edges) <= h.Mmax {
		return edges
	}

	// Keep the closest Mmax connections by evaluating distances to the stored vector.
	type edgeDist struct {
		id   uint64
		dist float32
	}

	vec := h.vectors[n.id]
	eds := make([]edgeDist, len(edges))
	for i, eid := range edges {
		eds[i] = edgeDist{id: eid, dist: h.distFn(h.vectors[eid], vec)}
	}

	for i := 0; i < len(eds); i++ {
		for j := i + 1; j < len(eds); j++ {
			if eds[j].dist < eds[i].dist {
				eds[i], eds[j] = eds[j], eds[i]
			}
		}
	}

	result := make([]uint64, h.Mmax)
	for i := 0; i < h.Mmax; i++ {
		result[i] = eds[i].id
	}
	return result
}

// randomLevel generates a random level using the HNSW level distribution.
func (h *HNSW) randomLevel() int {
	// Level distribution: floor(-ln(uniform) * levelMult)
	rv := h.rng.Float64()
	return int(math.Floor(-math.Log(rv) * h.levelMult))
}

// distItem is a heap item with ID and distance.
type distItem struct {
	id   uint64
	dist float32
}

// minHeap is a min-heap (closest on top). Used for candidates.
type minHeap []*distItem

func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].dist < h[j].dist }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *minHeap) Push(x interface{}) {
	*h = append(*h, x.(*distItem))
}

func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

// maxHeap is a max-heap (farthest on top). Used for results.
type maxHeap []*distItem

func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i].dist > h[j].dist }
func (h maxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *maxHeap) Push(x interface{}) {
	*h = append(*h, x.(*distItem))
}

func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

func removeID(s []uint64, id uint64) []uint64 {
	for i, v := range s {
		if v == id {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
