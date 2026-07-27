package api

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenAndInsert(t *testing.T) {
	dir, err := os.MkdirTemp("", "api-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	cfg := DefaultConfig("test", 4)
	cfg.SegmentMaxPoints = 100
	col, err := Open(dir, cfg)
	require.NoError(t, err)
	defer col.Close()

	err = col.Insert(1, []float32{1, 0, 0, 0}, nil, nil, []byte("hello"))
	require.NoError(t, err)
	assert.Equal(t, 1, col.Count())
}

func TestInsertMultipleAndCount(t *testing.T) {
	dir, err := os.MkdirTemp("", "api-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	col, err := Open(dir, DefaultConfig("test", 4))
	require.NoError(t, err)
	defer col.Close()

	for i := uint64(0); i < 10; i++ {
		require.NoError(t, col.Insert(i, []float32{float32(i) / 10, 1 - float32(i)/10, 0, 0}, nil, nil, nil))
	}
	assert.Equal(t, 10, col.Count())
}

func TestSearch(t *testing.T) {
	dir, err := os.MkdirTemp("", "api-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	col, err := Open(dir, DefaultConfig("test", 4))
	require.NoError(t, err)
	defer col.Close()

	for i := uint64(0); i < 10; i++ {
		require.NoError(t, col.Insert(i, []float32{float32(i) / 10, 1 - float32(i)/10, 0, 0}, nil, nil, nil))
	}

	results, err := col.Search([]float32{1, 0, 0, 0}, 3)
	require.NoError(t, err)
	require.Len(t, results, 3)
	assert.Less(t, results[0].Score, float32(0.1))
}

func TestDelete(t *testing.T) {
	dir, err := os.MkdirTemp("", "api-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	col, err := Open(dir, DefaultConfig("test", 4))
	require.NoError(t, err)
	defer col.Close()

	col.Insert(1, []float32{1, 0, 0, 0}, nil, nil, nil)
	col.Insert(2, []float32{0, 1, 0, 0}, nil, nil, nil)
	assert.Equal(t, 2, col.Count())

	err = col.Delete(1)
	require.NoError(t, err)
	assert.Equal(t, 1, col.Count())
}

func TestUpdate(t *testing.T) {
	dir, err := os.MkdirTemp("", "api-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	col, err := Open(dir, DefaultConfig("test", 4))
	require.NoError(t, err)
	defer col.Close()

	col.Insert(1, []float32{1, 0, 0, 0}, nil, nil, []byte("old"))
	col.Update(1, []float32{1, 0, 0, 0}, nil, nil, []byte("new"))
	assert.Equal(t, 1, col.Count())
}

func TestCreateDropListCollections(t *testing.T) {
	dir, err := os.MkdirTemp("", "api-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	cfg := DefaultConfig("mycol", 4)
	err = CreateCollection(dir, "mycol", cfg)
	require.NoError(t, err)

	names, err := ListCollections(dir)
	require.NoError(t, err)
	assert.Contains(t, names, "mycol")

	err = CreateCollection(dir, "mycol", cfg)
	assert.ErrorIs(t, err, ErrCollectionExists)

	err = DropCollection(dir, "mycol")
	require.NoError(t, err)

	names, err = ListCollections(dir)
	require.NoError(t, err)
	assert.NotContains(t, names, "mycol")

	err = DropCollection(dir, "nonexistent")
	assert.ErrorIs(t, err, ErrCollectionNotFound)
}

func TestSnapshot(t *testing.T) {
	dir, err := os.MkdirTemp("", "api-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	col, err := Open(dir, DefaultConfig("test", 4))
	require.NoError(t, err)
	defer col.Close()

	col.Insert(1, []float32{1, 0, 0, 0}, nil, nil, nil)
	snap, err := col.Snapshot("test-snap")
	require.NoError(t, err)
	assert.NotEmpty(t, snap.Name)
}

func TestBatch(t *testing.T) {
	dir, err := os.MkdirTemp("", "api-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	col, err := Open(dir, DefaultConfig("test", 4))
	require.NoError(t, err)
	defer col.Close()

	batch := NewBatch()
	for i := uint64(0); i < 5; i++ {
		batch.Add(i, []float32{float32(i) / 10, 1 - float32(i)/10, 0, 0}, nil, nil, nil)
	}
	assert.Equal(t, 5, batch.Len())

	err = batch.Execute(col)
	require.NoError(t, err)
	assert.Equal(t, 5, col.Count())
}
