package storage_test

import (
	"testing"

	"github.com/hasdev/oris/storage"
	"github.com/hasdev/oris/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorageGetSetDelete(t *testing.T) {
	s := memory.New()

	err := s.Set([]byte("key1"), []byte("value1"))
	require.NoError(t, err)

	v, err := s.Get([]byte("key1"))
	require.NoError(t, err)
	assert.Equal(t, []byte("value1"), v)

	_, err = s.Get([]byte("nonexistent"))
	assert.ErrorIs(t, err, storage.ErrNotFound)

	err = s.Delete([]byte("key1"))
	require.NoError(t, err)

	_, err = s.Get([]byte("key1"))
	assert.ErrorIs(t, err, storage.ErrNotFound)
}

func TestMemoryBatch(t *testing.T) {
	s := memory.New()

	b := s.NewBatch()
	b.Set([]byte("a"), []byte("1"))
	b.Set([]byte("b"), []byte("2"))
	b.Close()

	// Batch not committed, should not be visible.
	_, err := s.Get([]byte("a"))
	assert.ErrorIs(t, err, storage.ErrNotFound)

	b2 := s.NewBatch()
	b2.Set([]byte("a"), []byte("1"))
	b2.Set([]byte("b"), []byte("2"))
	require.NoError(t, b2.Commit())
	b2.Close()

	v, err := s.Get([]byte("a"))
	require.NoError(t, err)
	assert.Equal(t, []byte("1"), v)

	v, err = s.Get([]byte("b"))
	require.NoError(t, err)
	assert.Equal(t, []byte("2"), v)
}

func TestMemorySnapshot(t *testing.T) {
	s := memory.New()
	s.Set([]byte("x"), []byte("100"))
	s.Set([]byte("y"), []byte("200"))

	snap, err := s.Snapshot()
	require.NoError(t, err)
	defer snap.Close()

	// Modify after snapshot.
	s.Set([]byte("x"), []byte("changed"))
	s.Delete([]byte("y"))

	// Snapshot should reflect original state.
	v, err := snap.Get([]byte("x"))
	require.NoError(t, err)
	assert.Equal(t, []byte("100"), v)

	v, err = snap.Get([]byte("y"))
	require.NoError(t, err)
	assert.Equal(t, []byte("200"), v)

	// Iterator test.
	iter, err := snap.NewIterator([]byte("x"), []byte("z"))
	require.NoError(t, err)
	defer iter.Close()

	assert.True(t, iter.Next())
	k, v, err := iter.KeyValue()
	require.NoError(t, err)
	assert.Equal(t, []byte("x"), k)
	assert.Equal(t, []byte("100"), v)

	assert.True(t, iter.Next())
	k, v, err = iter.KeyValue()
	require.NoError(t, err)
	assert.Equal(t, []byte("y"), k)
	assert.Equal(t, []byte("200"), v)

	assert.False(t, iter.Next())
}
