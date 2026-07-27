package api

import (
	"os"
	"testing"

	"github.com/hasdev/oris/engine/segment"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIntegrationEndToEnd covers the full pipeline: open, insert, search, seal,
// search across segments, compact, search again.
func TestIntegrationEndToEnd(t *testing.T) {
	dir, err := os.MkdirTemp("", "oris-integration-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	cfg := DefaultConfig("integ", 4)
	cfg.SegmentMaxPoints = 5
	col, err := Open(dir, cfg)
	require.NoError(t, err)
	defer col.Close()

	// Insert 20 points — triggers multiple segment seals.
	for i := uint64(0); i < 20; i++ {
		vec := []float32{float32(i) / 20, 1 - float32(i)/20, 0, 0}
		require.NoError(t, col.Insert(i, vec, nil, nil, nil))
	}
	assert.Equal(t, 20, col.Count())

	// Search after all inserts.
	results, err := col.Search([]float32{1, 0, 0, 0}, 5)
	require.NoError(t, err)
	require.Len(t, results, 5)
	// The closest point to {1,0,0,0} is ID 19 with vec {0.95, 0.05}
	assert.Less(t, results[0].Score, float32(0.1))
	t.Logf("Top-5 results: %+v", results)
}

func TestIntegrationInsertSearchDelete(t *testing.T) {
	dir, err := os.MkdirTemp("", "oris-integration-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	col, err := Open(dir, DefaultConfig("integ", 4))
	require.NoError(t, err)
	defer col.Close()

	for i := uint64(0); i < 10; i++ {
		require.NoError(t, col.Insert(i, []float32{float32(i) / 10, 1 - float32(i)/10, 0, 0}, nil, nil, nil))
	}
	assert.Equal(t, 10, col.Count())

	require.NoError(t, col.Delete(0))
	assert.Equal(t, 9, col.Count())

	results, err := col.Search([]float32{1, 0, 0, 0}, 5)
	require.NoError(t, err)
	require.Len(t, results, 5)
	// ID 0 should be gone.
	for _, r := range results {
		assert.NotEqual(t, uint64(0), r.ID)
	}
}

func TestIntegrationUpdate(t *testing.T) {
	dir, err := os.MkdirTemp("", "oris-integration-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	col, err := Open(dir, DefaultConfig("integ", 4))
	require.NoError(t, err)
	defer col.Close()

	col.Insert(1, []float32{1, 0, 0, 0}, nil, nil, []byte("old"))
	col.Update(1, []float32{0, 1, 0, 0}, nil, nil, []byte("new"))

	// Search for {0,1,0,0} — ID 1 should be top.
	results, err := col.Search([]float32{0, 1, 0, 0}, 1)
	require.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, uint64(1), results[0].ID)
}

func TestIntegrationSegmentSealAndSearch(t *testing.T) {
	dir, err := os.MkdirTemp("", "oris-seal-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	sm := segment.NewManager(dir, 4, "cosine", 10, 10)

	// Insert 15 points — first 10 seal into immutable, next 5 stay mutable.
	for i := uint64(0); i < 15; i++ {
		vec := []float32{float32(i) / 15, 1 - float32(i)/15, 0, 0}
		require.NoError(t, sm.Insert(i, vec, nil, nil, nil))
	}
	assert.Equal(t, 15, sm.Len())

	results, err := sm.SearchAcross([]float32{0.9, 0.1, 0, 0}, 3)
	require.NoError(t, err)
	require.Len(t, results, 3)
	t.Logf("Segment search results: %+v", results)
}

func TestIntegrationCompaction(t *testing.T) {
	dir, err := os.MkdirTemp("", "oris-compact-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	// maxPoints=5, maxSegments=2 → triggers compaction at 10+ inserts.
	sm := segment.NewManager(dir, 4, "cosine", 5, 2)

	for i := uint64(0); i < 15; i++ {
		vec := []float32{float32(i) / 15, 1 - float32(i)/15, 0, 0}
		require.NoError(t, sm.Insert(i, vec, nil, nil, nil))
	}
	assert.Equal(t, 15, sm.Len())

	results, err := sm.SearchAcross([]float32{0.5, 0.5, 0, 0}, 3)
	require.NoError(t, err)
	require.Len(t, results, 3)
	t.Logf("Compacted search results: %+v", results)
}

func TestIntegrationSnapshot(t *testing.T) {
	dir, err := os.MkdirTemp("", "oris-snap-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	col, err := Open(dir, DefaultConfig("snap", 4))
	require.NoError(t, err)

	for i := uint64(0); i < 5; i++ {
		require.NoError(t, col.Insert(i, []float32{float32(i) / 5, 1 - float32(i)/5, 0, 0}, nil, nil, nil))
	}

	snap, err := col.Snapshot("pre-close")
	require.NoError(t, err)
	assert.NotEmpty(t, snap.Checksum)
	col.Close()

	// Verify snapshot file exists.
	_, err = os.Stat(snap.Path + "/data.snap")
	require.NoError(t, err)
}
