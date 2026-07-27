// Package snapshot provides snapshot creation, restore, and verification utilities.
package snapshot

import (
	"encoding/binary"
	"io"
	"os"
	"path/filepath"

	"lukechampine.com/blake3"

	"github.com/hasdev/oris/storage"
)

// Snapshot represents a stored snapshot with metadata.
type Snapshot struct {
	Path     string
	Name     string
	Checksum [32]byte
}

// Create takes a point-in-time snapshot of the storage and writes it to disk.
func Create(store storage.Storage, dir, name string) (*Snapshot, error) {
	snapDir := filepath.Join(dir, name)
	if err := os.MkdirAll(snapDir, 0755); err != nil {
		return nil, err
	}

	snap, err := store.Snapshot()
	if err != nil {
		return nil, err
	}
	defer snap.Close()

	dataPath := filepath.Join(snapDir, "data.snap")
	f, err := os.Create(dataPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	hasher := blake3.New(32, nil)
	mw := io.MultiWriter(f, hasher)

	iter, err := snap.NewIterator([]byte{0x00}, []byte{0xFF})
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	buf := make([]byte, 8)
	for iter.Next() {
		key, value, err := iter.KeyValue()
		if err != nil {
			return nil, err
		}

		// Write key length + key + value length + value.
		binary.LittleEndian.PutUint32(buf[:4], uint32(len(key)))
		if _, err := mw.Write(buf[:4]); err != nil {
			return nil, err
		}
		if _, err := mw.Write(key); err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint32(buf[:4], uint32(len(value)))
		if _, err := mw.Write(buf[:4]); err != nil {
			return nil, err
		}
		if _, err := mw.Write(value); err != nil {
			return nil, err
		}
	}

	var checksum [32]byte
	copy(checksum[:], hasher.Sum(nil))

	return &Snapshot{
		Path:     snapDir,
		Name:     name,
		Checksum: checksum,
	}, nil
}

// Restore reads a snapshot file and writes all key-value pairs into the storage.
func Restore(store storage.Storage, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	batch := store.NewBatch()
	defer batch.Close()

	head := make([]byte, 4)
	for {
		// Read key length.
		_, err := io.ReadFull(f, head)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		keyLen := binary.LittleEndian.Uint32(head)

		key := make([]byte, keyLen)
		if _, err := io.ReadFull(f, key); err != nil {
			return err
		}

		// Read value length.
		if _, err := io.ReadFull(f, head); err != nil {
			return err
		}
		valLen := binary.LittleEndian.Uint32(head)

		value := make([]byte, valLen)
		if _, err := io.ReadFull(f, value); err != nil {
			return err
		}

		batch.Set(key, value)
	}

	return batch.Commit()
}

// Verify checks that the snapshot file matches the stored checksum.
func Verify(path string, checksum [32]byte) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	hasher := blake3.New(32, nil)
	if _, err := io.Copy(hasher, f); err != nil {
		return false, err
	}

	var got [32]byte
	copy(got[:], hasher.Sum(nil))
	return got == checksum, nil
}
