package snapshot

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hasdev/oris/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAndRestore(t *testing.T) {
	dir, err := os.MkdirTemp("", "snapshot-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	store := memory.New()
	store.Set([]byte("key1"), []byte("value1"))
	store.Set([]byte("key2"), []byte("value2"))

	snap, err := Create(store, dir, "test-snap")
	require.NoError(t, err)
	assert.NotEmpty(t, snap.Checksum)
	assert.Equal(t, "test-snap", snap.Name)

	// Verify checksum.
	snapPath := filepath.Join(dir, "test-snap", "data.snap")
	ok, err := Verify(snapPath, snap.Checksum)
	require.NoError(t, err)
	assert.True(t, ok)

	// Restore into a new storage.
	store2 := memory.New()
	err = Restore(store2, snapPath)
	require.NoError(t, err)

	v, err := store2.Get([]byte("key1"))
	require.NoError(t, err)
	assert.Equal(t, []byte("value1"), v)

	v, err = store2.Get([]byte("key2"))
	require.NoError(t, err)
	assert.Equal(t, []byte("value2"), v)
}

func TestVerifyTampered(t *testing.T) {
	dir, err := os.MkdirTemp("", "snapshot-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	store := memory.New()
	store.Set([]byte("a"), []byte("b"))

	snap, err := Create(store, dir, "tampered")
	require.NoError(t, err)

	snapPath := filepath.Join(dir, "tampered", "data.snap")

	// Tamper with the data.
	badChecksum := snap.Checksum
	badChecksum[0] ^= 0xFF

	ok, err := Verify(snapPath, badChecksum)
	require.NoError(t, err)
	assert.False(t, ok)
}
