// Package dense provides dense vector index implementations.
package dense

// Engine is the interface for dense vector indexing and search.
type Engine interface {
	// Insert adds a vector with the given ID.
	Insert(id uint64, vec []float32) error

	// Delete removes the vector with the given ID.
	Delete(id uint64) error

	// Search finds the top-k most similar vectors.
	Search(query []float32, topK int) ([]Result, error)

	// Len returns the number of vectors in the index.
	Len() int

	// Rebuild re-inserts all vectors into a fresh index. Useful after
	// parameter changes or to recover from graph degradation.
	Rebuild() error
}

// SmallThreshold is the number of vectors below which a flat index is used.
// Query it from factory.New to auto-select the right backend.
var SmallThreshold = 1000

// Result is a single search result.
type Result struct {
	ID    uint64
	Score float32
}

// Config holds configuration for a dense engine.
type Config struct {
	Dimension      int
	Distance       string // "cosine", "dot", "euclidean"
	M              int   // HNSW M parameter (default 16)
	EfConstruction int   // HNSW efConstruction (default 200)
	EfSearch       int   // HNSW efSearch (default 100)
}
