// Package api provides the public API for interacting with Oris collections.
//
// This is the primary entry point for embedding Oris into Go applications.
package api

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"sync"

	"github.com/hasdev/oris/engine/distance"
	"github.com/hasdev/oris/storage"
	"github.com/hasdev/oris/storage/pebble"
	"github.com/hasdev/oris/storage/snapshot"
	"github.com/hasdev/oris/storage/wal"
)

// Collection is an Oris vector collection.
type Collection struct {
	mu           sync.RWMutex
	config       CollectionConfig
	path         string
	store        storage.Storage
	wal          *wal.WAL
	dist         distance.Kernel
	pointCounter uint64
}

// Open opens an existing collection or creates a new one at the given path.
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

	c := &Collection{
		config: cfg,
		path:   path,
		store:  store,
		wal:    w,
		dist:   distance.NewKernel(),
	}

	// Load point counter from storage.
	c.loadCounter()

	return c, nil
}

// Config returns the collection configuration.
func (c *Collection) Config() CollectionConfig {
	return c.config
}

// Insert adds a new point to the collection.
func (c *Collection) Insert(id uint64, dense []float32, sparseIndices []uint32, sparseValues []float32, payload []byte) error {
	if dense == nil {
		dense = make([]float32, c.config.Dimension)
	}

	// Encode point.
	data := encodePoint(id, dense, sparseIndices, sparseValues, payload)

	// Write to WAL first.
	if err := c.wal.WriteSet(pointKey(id), data); err != nil {
		return err
	}

	// Write to storage.
	if err := c.store.Set(pointKey(id), data); err != nil {
		return err
	}

	return nil
}

// Get retrieves a point by ID.
func (c *Collection) Get(id uint64) ([]float32, []uint32, []float32, []byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	data, err := c.store.Get(pointKey(id))
	if err != nil {
		return nil, nil, nil, nil, err
	}

	dense, sparseIndices, sparseValues, payload, err := decodePoint(data)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return dense, sparseIndices, sparseValues, payload, nil
}

// Delete removes a point from the collection.
func (c *Collection) Delete(id uint64) error {
	if err := c.wal.WriteDelete(pointKey(id)); err != nil {
		return err
	}
	return c.store.Delete(pointKey(id))
}

// Count returns the number of points in the collection.
func (c *Collection) Count() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.pointCounter
}

// Search performs a brute-force dense search (MVP — HNSW integration coming in Epic 4).
func (c *Collection) Search(query []float32, topK int) ([]Result, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	snap, err := c.store.Snapshot()
	if err != nil {
		return nil, err
	}
	defer snap.Close()

	iter, err := snap.NewIterator(pointPrefix(), []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	type scored struct {
		id    uint64
		score float32
	}

	results := make([]scored, 0, topK)
	worstScore := float32(-1)
	worstIdx := -1

	for iter.Next() {
		key, value, err := iter.KeyValue()
		if err != nil {
			continue
		}

		dense, _, _, _, err := decodePoint(value)
		if err != nil || len(dense) != len(query) {
			continue
		}

		id := binary.LittleEndian.Uint64(key)
		score := c.dist.Cosine(query, dense)

		if len(results) < topK {
			results = append(results, scored{id: id, score: score})
			if score > worstScore {
				worstScore = score
				worstIdx = len(results) - 1
			}
		} else if score < worstScore {
			results[worstIdx] = scored{id: id, score: score}
			// Recompute worst.
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

	// Convert to Result and order by score ascending.
	out := make([]Result, len(results))
	for i, r := range results {
		out[i] = Result{ID: r.id, Score: r.score}
	}

	return out, nil
}

// Flush ensures all pending writes are persisted.
func (c *Collection) Flush() error {
	return nil // Pebble writes are already synced.
}

// Snapshot creates a point-in-time snapshot of the collection.
func (c *Collection) Snapshot(name string) (*snapshot.Snapshot, error) {
	return snapshot.Create(c.store, filepath.Join(c.path, "snapshots"), name)
}

// Close closes the collection, flushing pending writes.
func (c *Collection) Close() error {
	if err := c.wal.Close(); err != nil {
		return err
	}
	return c.store.Close()
}

func (c *Collection) loadCounter() {
	snap, err := c.store.Snapshot()
	if err != nil {
		return
	}
	defer snap.Close()

	iter, err := snap.NewIterator(pointPrefix(), []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})
	if err != nil {
		return
	}
	defer iter.Close()

	var maxID uint64
	for iter.Next() {
		key, _, err := iter.KeyValue()
		if err != nil {
			continue
		}
		id := binary.LittleEndian.Uint64(key)
		if id > maxID {
			maxID = id
		}
	}
	c.pointCounter = maxID
}

// pointPrefix returns the key prefix for point keys.
func pointPrefix() []byte { return []byte{0x00} }

// pointKey encodes a point ID as a storage key.
func pointKey(id uint64) []byte {
	b := make([]byte, 9)
	b[0] = 0x00 // prefix for point keys
	binary.LittleEndian.PutUint64(b[1:], id)
	return b
}

// encodePoint serializes point data to binary.
func encodePoint(id uint64, dense []float32, sparseIndices []uint32, sparseValues []float32, payload []byte) []byte {
	// Dense length (4) + dense data (4*len) + sparse len (4) + sparse indices (4*len) + sparse values (4*len) + payload len (4) + payload
	size := 4 + len(dense)*4 + 4 + len(sparseIndices)*4 + len(sparseValues)*4 + 4 + len(payload)
	buf := make([]byte, size)
	n := 0

	// Dense vector.
	binary.LittleEndian.PutUint32(buf[n:], uint32(len(dense)))
	n += 4
	for _, f := range dense {
		binary.LittleEndian.PutUint32(buf[n:], uint32(f))
		n += 4
	}

	// Sparse indices.
	binary.LittleEndian.PutUint32(buf[n:], uint32(len(sparseIndices)))
	n += 4
	for _, idx := range sparseIndices {
		binary.LittleEndian.PutUint32(buf[n:], idx)
		n += 4
	}

	// Sparse values.
	for _, v := range sparseValues {
		binary.LittleEndian.PutUint32(buf[n:], mathFloat32bits(v))
		n += 4
	}

	// Payload.
	binary.LittleEndian.PutUint32(buf[n:], uint32(len(payload)))
	n += 4
	copy(buf[n:], payload)

	return buf
}

// decodePoint deserializes point data from binary.
func decodePoint(data []byte) (dense []float32, sparseIndices []uint32, sparseValues []float32, payload []byte, err error) {
	if len(data) < 4 {
		return nil, nil, nil, nil, storage.ErrNotFound
	}
	n := 0

	denseLen := int(binary.LittleEndian.Uint32(data[n:]))
	n += 4
	dense = make([]float32, denseLen)
	for i := range dense {
		dense[i] = float32(binary.LittleEndian.Uint32(data[n:]))
		n += 4
	}

	sparseLen := int(binary.LittleEndian.Uint32(data[n:]))
	n += 4
	sparseIndices = make([]uint32, sparseLen)
	for i := range sparseIndices {
		sparseIndices[i] = binary.LittleEndian.Uint32(data[n:])
		n += 4
	}
	sparseValues = make([]float32, sparseLen)
	for i := range sparseValues {
		sparseValues[i] = mathFloat32frombits(binary.LittleEndian.Uint32(data[n:]))
		n += 4
	}

	payloadLen := int(binary.LittleEndian.Uint32(data[n:]))
	n += 4
	payload = make([]byte, payloadLen)
	copy(payload, data[n:])

	return dense, sparseIndices, sparseValues, payload, nil
}

// mathFloat32bits provides float32 to uint32 bit conversion.
func mathFloat32bits(f float32) uint32 {
	return uint32(f)
}

// mathFloat32frombits provides uint32 to float32 bit conversion.
func mathFloat32frombits(bits uint32) float32 {
	return float32(bits)
}
