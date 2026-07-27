package segment

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMutableInsertAndLen(t *testing.T) {
	s := NewMutable(4, "cosine")
	assert.Equal(t, StateMutable, s.State())
	assert.Equal(t, 0, s.Len())

	err := s.Insert(1, []float32{1, 0, 0, 0}, nil, nil, []byte("hello"))
	require.NoError(t, err)
	assert.Equal(t, 1, s.Len())

	err = s.Insert(2, []float32{0, 1, 0, 0}, nil, nil, []byte("world"))
	require.NoError(t, err)
	assert.Equal(t, 2, s.Len())
}

func TestMutableDelete(t *testing.T) {
	s := NewMutable(4, "cosine")
	s.Insert(1, []float32{1, 0, 0, 0}, nil, nil, nil)
	s.Insert(2, []float32{0, 1, 0, 0}, nil, nil, nil)
	assert.Equal(t, 2, s.Len())

	err := s.Delete(1)
	require.NoError(t, err)
	assert.Equal(t, 1, s.Len())
}

func TestMutableSeal(t *testing.T) {
	dir, err := os.MkdirTemp("", "segment-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	s := NewMutable(4, "cosine")
	s.Insert(1, []float32{1, 0, 0, 0}, nil, nil, []byte("payload1"))
	s.Insert(2, []float32{0, 1, 0, 0}, nil, nil, []byte("payload2"))

	err = s.Seal(dir)
	require.NoError(t, err)
	assert.Equal(t, StateImmutable, s.State())

	// Sealed segment should reject writes.
	err = s.Insert(3, nil, nil, nil, nil)
	assert.ErrorIs(t, err, ErrSegmentSealed)
}

func TestMutableSearch(t *testing.T) {
	s := NewMutable(4, "cosine")
	s.Insert(1, []float32{1, 0, 0, 0}, nil, nil, nil)
	s.Insert(2, []float32{0, 1, 0, 0}, nil, nil, nil)

	results, err := s.Search([]float32{0.9, 0.1, 0, 0}, 5, nil)
	require.NoError(t, err)
	require.Len(t, results, 2)
	assert.Equal(t, uint64(1), results[0].ID)
	assert.Equal(t, uint64(2), results[1].ID)
}

func TestOpenImmutable(t *testing.T) {
	dir, err := os.MkdirTemp("", "segment-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	s := NewMutable(4, "cosine")
	s.Insert(100, []float32{1, 0, 0, 0}, nil, nil, []byte("data"))
	require.NoError(t, s.Seal(dir))

	segDir := fmt.Sprintf("%s/segment_%06d", dir, s.meta.ID)
	imm, err := OpenImmutable(segDir, 4, "cosine")
	require.NoError(t, err)
	assert.Equal(t, 1, imm.Len())
	assert.Equal(t, StateImmutable, imm.State())

	err = imm.Insert(200, nil, nil, nil, nil)
	assert.ErrorIs(t, err, ErrSegmentSealed)
}

func TestSegmentManagerInsertAndSeal(t *testing.T) {
	dir, err := os.MkdirTemp("", "segment-manager-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	sm := NewManager(dir, 4, "cosine", 5, 3)

	for i := uint64(0); i < 12; i++ {
		vec := []float32{float32(i) / 10, 1 - float32(i)/10, 0, 0}
		err := sm.Insert(i, vec, nil, nil, []byte("p"))
		require.NoError(t, err)
	}
	assert.Equal(t, 12, sm.Len())
}

func TestSegmentManagerSearchAcrossSegments(t *testing.T) {
	dir, err := os.MkdirTemp("", "segment-manager-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	sm := NewManager(dir, 4, "cosine", 5, 10)

	for i := uint64(0); i < 12; i++ {
		vec := []float32{float32(i) / 10, 1 - float32(i)/10, 0, 0}
		require.NoError(t, sm.Insert(i, vec, nil, nil, nil))
	}

	results, err := sm.SearchAcross([]float32{1, 0, 0, 0}, 3)
	require.NoError(t, err)
	require.Len(t, results, 3)

	// The closest point to {1,0,0,0} should have score near 0.
	assert.Less(t, results[0].Score, float32(0.1))
}

func TestSegmentManagerCompaction(t *testing.T) {
	dir, err := os.MkdirTemp("", "segment-compact-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	// maxPoints=3, maxSegments=2 -> triggers compaction after 6 inserts.
	sm := NewManager(dir, 4, "cosine", 3, 2)

	for i := uint64(0); i < 7; i++ {
		vec := []float32{float32(i) / 10, 1 - float32(i)/10, 0, 0}
		require.NoError(t, sm.Insert(i, vec, nil, nil, nil))
	}
	assert.Equal(t, 7, sm.Len())

	// Search should still work after compaction.
	results, err := sm.SearchAcross([]float32{0.5, 0.5, 0, 0}, 3)
	require.NoError(t, err)
	require.Len(t, results, 3)
}

func TestSegmentInterface(t *testing.T) {
	var s Segment = NewMutable(4, "cosine")
	_ = s

	dir, _ := os.MkdirTemp("", "seg-interface-*")
	defer os.RemoveAll(dir)
	mut := NewMutable(4, "cosine")
	mut.Insert(1, nil, nil, nil, nil)
	mut.Seal(dir)
	segDir := fmt.Sprintf("%s/segment_%06d", dir, mut.meta.ID)
	imm, _ := OpenImmutable(segDir, 4, "cosine")
	var is Segment = imm
	_ = is
}

func TestStateString(t *testing.T) {
	assert.Equal(t, "mutable", StateMutable.String())
	assert.Equal(t, "sealing", StateSealing.String())
	assert.Equal(t, "immutable", StateImmutable.String())
	assert.Equal(t, "compacting", StateCompacting.String())
	assert.Equal(t, "deleted", StateDeleted.String())
	assert.Equal(t, "unknown", State(99).String())
}
