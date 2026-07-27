// Package api provides the public API for interacting with Oris collections.
//
// This is the primary entry point for embedding Oris into Go applications.
package api

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"sync"

	"github.com/hasdev/oris/engine/dense"
	"github.com/hasdev/oris/engine/metadata"
	"github.com/hasdev/oris/engine/planner"
	"github.com/hasdev/oris/engine/segment"
	"github.com/hasdev/oris/storage"
	"github.com/hasdev/oris/storage/pebble"
	"github.com/hasdev/oris/storage/snapshot"
	"github.com/hasdev/oris/storage/wal"
)

// Collection is an Oris vector collection with full engine integration.
type Collection struct {
	mu      sync.RWMutex
	cfg     CollectionConfig
	path    string
	store   storage.Storage
	wal     *wal.WAL
	segMgr  *segment.SegmentManager
	mgr     *segment.SegmentManager
	planner *planner.Planner
	meta    *metadata.Engine
}

// segmentEngine wraps SegmentManager to implement dense.Engine.
type segmentEngine struct {
	sm *segment.SegmentManager
}

func (s segmentEngine) Insert(id uint64, vec []float32) error {
	return s.sm.Insert(id, vec, nil, nil, nil)
}
func (s segmentEngine) Delete(id uint64) error     { return s.sm.Delete(id) }
func (s segmentEngine) Search(q []float32, k int) ([]dense.Result, error) {
	return s.sm.SearchAcross(q, k)
}
func (s segmentEngine) Len() int                         { return s.sm.Len() }
func (s segmentEngine) Rebuild() error                   { return nil }

// Open opens or creates a collection at the given path.
func Open(path string, cfg CollectionConfig) (*Collection, error) {
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}

	store, err := pebble.Open(pebble.Options{Path: filepath.Join(path, "data")})
	if err != nil {
		return nil, err
	}

	w, err := wal.Open(path, store)
	if err != nil {
		store.Close()
		return nil, err
	}

	segDir := filepath.Join(path, "segments")
	if err := os.MkdirAll(segDir, 0755); err != nil {
		store.Close()
		w.Close()
		return nil, err
	}

	sm := segment.NewManager(segDir, cfg.Dimension, cfg.Distance, cfg.SegmentMaxPoints, 10)
	metaEngine := metadata.New()
	den := segmentEngine{sm: sm}
	pl := planner.New(den, nil, metaEngine, cfg.Dimension)

	c := &Collection{
		cfg:     cfg,
		path:    path,
		store:   store,
		wal:     w,
		segMgr:  sm,
		mgr:     sm,
		planner: pl,
		meta:    metaEngine,
	}

	return c, nil
}

func (c *Collection) Config() CollectionConfig { return c.cfg }

// Insert adds a point to the collection.
func (c *Collection) Insert(id uint64, denseVec []float32, sparseIndices []uint32, sparseValues []float32, payload []byte) error {
	return c.mgr.Insert(id, denseVec, sparseIndices, sparseValues, payload)
}

// Delete removes a point.
func (c *Collection) Delete(id uint64) error {
	return c.mgr.Delete(id)
}

// Update replaces a point (delete + insert).
func (c *Collection) Update(id uint64, dense []float32, sparseIndices []uint32, sparseValues []float32, payload []byte) error {
	if err := c.Delete(id); err != nil {
		return err
	}
	return c.Insert(id, dense, sparseIndices, sparseValues, payload)
}

// Count returns the total number of points.
func (c *Collection) Count() int {
	return c.mgr.Len()
}

// Search performs a dense-only search.
func (c *Collection) Search(query []float32, topK int) ([]Result, error) {
	results, err := c.planner.Execute(planner.Query{
		DenseVector: query,
		TopK:        topK,
		Mode:        planner.DenseOnly,
	})
	if err != nil {
		return nil, err
	}
	out := make([]Result, len(results))
	for i, r := range results {
		out[i] = Result{ID: r.ID, Score: r.FinalScore}
	}
	return out, nil
}

// SearchWithFilter performs dense search with metadata filter.
func (c *Collection) SearchWithFilter(query []float32, topK int, filter metadata.Filter) ([]Result, error) {
	results, err := c.planner.Execute(planner.Query{
		DenseVector: query,
		TopK:        topK,
		Mode:        planner.DenseOnly,
		Filter:      filter,
	})
	if err != nil {
		return nil, err
	}
	out := make([]Result, len(results))
	for i, r := range results {
		out[i] = Result{ID: r.ID, Score: r.FinalScore}
	}
	return out, nil
}

// Snapshot creates a point-in-time snapshot.
func (c *Collection) Snapshot(name string) (*snapshot.Snapshot, error) {
	return snapshot.Create(c.store, filepath.Join(c.path, "snapshots"), name)
}

// Close closes the collection.
func (c *Collection) Close() error {
	c.wal.Close()
	return c.store.Close()
}

// CreateCollection creates a new collection in the given workspace.
func CreateCollection(workspace, name string, cfg CollectionConfig) error {
	colPath := filepath.Join(workspace, name)
	if _, err := os.Stat(colPath); err == nil {
		return ErrCollectionExists
	}
	col, err := Open(colPath, cfg)
	if err != nil {
		return err
	}
	return col.Close()
}

// DropCollection removes a collection.
func DropCollection(workspace, name string) error {
	colPath := filepath.Join(workspace, name)
	if _, err := os.Stat(colPath); os.IsNotExist(err) {
		return ErrCollectionNotFound
	}
	return os.RemoveAll(colPath)
}

// ListCollections returns collection names in the workspace.
func ListCollections(workspace string) ([]string, error) {
	entries, err := os.ReadDir(workspace)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

// Ensure unused imports aren't flagged.
var _ = binary.LittleEndian

// Batch is a batch of insert operations.
type Batch struct {
	mu     sync.Mutex
	points []batchPoint
}

type batchPoint struct {
	id      uint64
	dense   []float32
	si      []uint32
	sv      []float32
	payload []byte
}

// NewBatch creates a new batch.
func NewBatch() *Batch {
	return &Batch{}
}

// Add adds a point to the batch.
func (b *Batch) Add(id uint64, dense []float32, sparseIndices []uint32, sparseValues []float32, payload []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.points = append(b.points, batchPoint{id: id, dense: dense, si: sparseIndices, sv: sparseValues, payload: payload})
}

// Execute inserts all batched points into the collection.
func (b *Batch) Execute(col *Collection) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, p := range b.points {
		if err := col.Insert(p.id, p.dense, p.si, p.sv, p.payload); err != nil {
			return err
		}
	}
	return nil
}

// Len returns the number of points in the batch.
func (b *Batch) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.points)
}
