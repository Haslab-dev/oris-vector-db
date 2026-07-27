package internal

import (
	"encoding/binary"
	"io"
	"os"
	"time"

	"github.com/cespare/xxhash/v2"
)

// SegmentMeta holds serializable segment metadata.
type SegmentMeta struct {
	ID        uint64
	DocCount  int32
	MinID     uint64
	MaxID     uint64
	CreatedAt int64 // unix nanos
	SealedAt  int64 // unix nanos
	DenseDim  int32
	Checksum  uint64 // xxhash of all preceding fields
}

// MarshalSegmentMeta serializes SegmentMeta to binary.
func MarshalSegmentMeta(m *SegmentMeta) ([]byte, error) {
	// 8+4+8+8+8+8+4 = 48 bytes before checksum
	buf := make([]byte, 56)
	binary.LittleEndian.PutUint64(buf[0:], m.ID)
	binary.LittleEndian.PutUint32(buf[8:], uint32(m.DocCount))
	binary.LittleEndian.PutUint64(buf[12:], m.MinID)
	binary.LittleEndian.PutUint64(buf[20:], m.MaxID)
	binary.LittleEndian.PutUint64(buf[28:], uint64(m.CreatedAt))
	binary.LittleEndian.PutUint64(buf[36:], uint64(m.SealedAt))
	binary.LittleEndian.PutUint32(buf[44:], uint32(m.DenseDim))
	// Checksum of first 48 bytes
	h := xxhash.Sum64(buf[:48])
	binary.LittleEndian.PutUint64(buf[48:], h)
	return buf, nil
}

// UnmarshalSegmentMeta deserializes SegmentMeta and verifies checksum.
func UnmarshalSegmentMeta(data []byte) (*SegmentMeta, error) {
	if len(data) < 56 {
		return nil, ErrInvalidEncoding
	}
	h := xxhash.Sum64(data[:48])
	stored := binary.LittleEndian.Uint64(data[48:])
	if h != stored {
		return nil, ErrInvalidEncoding
	}
	return &SegmentMeta{
		ID:        binary.LittleEndian.Uint64(data[0:]),
		DocCount:  int32(binary.LittleEndian.Uint32(data[8:])),
		MinID:     binary.LittleEndian.Uint64(data[12:]),
		MaxID:     binary.LittleEndian.Uint64(data[20:]),
		CreatedAt: int64(binary.LittleEndian.Uint64(data[28:])),
		SealedAt:  int64(binary.LittleEndian.Uint64(data[36:])),
		DenseDim:  int32(binary.LittleEndian.Uint32(data[44:])),
	}, nil
}

// WriteSegmentMetaFile writes segment metadata to a file.
func WriteSegmentMetaFile(path string, m *SegmentMeta) error {
	data, err := MarshalSegmentMeta(m)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// ReadSegmentMetaFile reads segment metadata from a file.
func ReadSegmentMetaFile(path string) (*SegmentMeta, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return UnmarshalSegmentMeta(data)
}

// WriteSegmentPayload writes a payload entry (8 byte ID + 4 byte length + payload).
func WriteSegmentPayload(w io.Writer, id uint64, payload []byte) error {
	buf := make([]byte, 12+len(payload))
	binary.LittleEndian.PutUint64(buf, id)
	binary.LittleEndian.PutUint32(buf[8:], uint32(len(payload)))
	copy(buf[12:], payload)
	_, err := w.Write(buf)
	return err
}

// ReadSegmentPayload reads one payload entry from a reader.
func ReadSegmentPayload(r io.Reader) (id uint64, payload []byte, err error) {
	head := make([]byte, 12)
	if _, err = io.ReadFull(r, head); err != nil {
		return
	}
	id = binary.LittleEndian.Uint64(head)
	plen := binary.LittleEndian.Uint32(head[8:])
	payload = make([]byte, plen)
	_, err = io.ReadFull(r, payload)
	return
}

// MakeSegmentMeta creates a SegmentMeta from arguments.
func MakeSegmentMeta(id uint64, docCount int, minID, maxID uint64, dim int) *SegmentMeta {
	now := time.Now().UnixNano()
	return &SegmentMeta{
		ID:        id,
		DocCount:  int32(docCount),
		MinID:     minID,
		MaxID:     maxID,
		CreatedAt: now,
		SealedAt:  now,
		DenseDim:  int32(dim),
	}
}
