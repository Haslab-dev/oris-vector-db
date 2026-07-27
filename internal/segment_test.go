package internal

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSegmentMetaRoundtrip(t *testing.T) {
	m := MakeSegmentMeta(42, 1000, 1, 5000, 128)

	data, err := MarshalSegmentMeta(m)
	require.NoError(t, err)

	m2, err := UnmarshalSegmentMeta(data)
	require.NoError(t, err)

	assert.Equal(t, m.ID, m2.ID)
	assert.Equal(t, m.DocCount, m2.DocCount)
	assert.Equal(t, m.MinID, m2.MinID)
	assert.Equal(t, m.MaxID, m2.MaxID)
	assert.Equal(t, m.DenseDim, m2.DenseDim)
}

func TestSegmentMetaTampered(t *testing.T) {
	m := MakeSegmentMeta(1, 10, 0, 100, 4)
	data, err := MarshalSegmentMeta(m)
	require.NoError(t, err)

	// Corrupt the checksum.
	data[48] ^= 0xFF
	_, err = UnmarshalSegmentMeta(data)
	assert.ErrorIs(t, err, ErrInvalidEncoding)
}

func TestSegmentMetaFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "segmeta-*")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	path := dir + "/meta.bin"
	m := MakeSegmentMeta(7, 500, 10, 999, 64)
	require.NoError(t, WriteSegmentMetaFile(path, m))

	m2, err := ReadSegmentMetaFile(path)
	require.NoError(t, err)
	assert.Equal(t, m.ID, m2.ID)
	assert.Equal(t, m.DocCount, m2.DocCount)
}

func TestWriteReadSegmentPayload(t *testing.T) {
	var buf bytes.Buffer
	err := WriteSegmentPayload(&buf, 100, []byte("hello world"))
	require.NoError(t, err)

	id, payload, err := ReadSegmentPayload(&buf)
	require.NoError(t, err)
	assert.Equal(t, uint64(100), id)
	assert.Equal(t, []byte("hello world"), payload)
}

func TestWriteReadEmptyPayload(t *testing.T) {
	var buf bytes.Buffer
	err := WriteSegmentPayload(&buf, 1, nil)
	require.NoError(t, err)

	id, payload, err := ReadSegmentPayload(&buf)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), id)
	assert.Empty(t, payload)
}
