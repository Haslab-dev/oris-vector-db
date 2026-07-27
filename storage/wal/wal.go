// Package wal provides a write-ahead log for durability.
//
// Every write operation is recorded in the WAL before being applied to
// storage, enabling crash recovery by replaying the log on open.
package wal

import (
	"encoding/binary"

	"github.com/cockroachdb/pebble"

	"github.com/hasdev/oris/storage"
)

const (
	opSet    byte = 0x01
	opDelete byte = 0x02
)

// WAL is a write-ahead log backed by Pebble.
type WAL struct {
	db  *pebble.DB
	seq uint64
}

// Open opens or creates a WAL at the given path.
// If a previous WAL exists, it replays outstanding entries.
func Open(path string, store storage.Storage) (*WAL, error) {
	opts := pebble.Options{}
	db, err := pebble.Open(path+"/wal", &opts)
	if err != nil {
		return nil, err
	}

	w := &WAL{db: db}

	// Replay uncommitted entries on open.
	if err := w.replay(store); err != nil {
		db.Close()
		return nil, err
	}

	return w, nil
}

// WriteSet writes a Set operation to the WAL.
func (w *WAL) WriteSet(key, value []byte) error {
	return w.writeEntry(opSet, key, value)
}

// WriteDelete writes a Delete operation to the WAL.
func (w *WAL) WriteDelete(key []byte) error {
	return w.writeEntry(opDelete, key, nil)
}

func (w *WAL) writeEntry(op byte, key, value []byte) error {
	w.seq++

	// Entry format: op(1) + seq(8) + keyLen(4) + key(K) + valLen(4) + val(V) + sentinel(1)
	entry := make([]byte, 1+8+4+len(key)+4+len(value)+1)
	n := 0
	entry[n] = op
	n++
	binary.LittleEndian.PutUint64(entry[n:], w.seq)
	n += 8
	binary.LittleEndian.PutUint32(entry[n:], uint32(len(key)))
	n += 4
	copy(entry[n:], key)
	n += len(key)
	binary.LittleEndian.PutUint32(entry[n:], uint32(len(value)))
	n += 4
	copy(entry[n:], value)
	n += len(value)
	entry[n] = 0xFF // sentinel

	return w.db.Set(seqKey(w.seq), entry, pebble.Sync)
}

// Truncate removes all WAL entries.
func (w *WAL) Truncate() error {
	iter, err := w.db.NewIter(nil)
	if err != nil {
		return err
	}
	defer iter.Close()

	for iter.First(); iter.Valid(); iter.Next() {
		if err := w.db.Delete(iter.Key(), pebble.Sync); err != nil {
			return err
		}
	}
	return nil
}

// Close closes the WAL.
func (w *WAL) Close() error {
	return w.db.Close()
}

func (w *WAL) replay(store storage.Storage) error {
	iter, err := w.db.NewIter(nil)
	if err != nil {
		return err
	}
	defer iter.Close()

	// Minimum entry size: 1(op) + 8(seq) + 4(keyLen) + 0(key) + 4(valLen) + 0(val) + 1(sentinel) = 18
	const minEntryLen = 18

	for iter.First(); iter.Valid(); iter.Next() {
		entry := iter.Value()
		if len(entry) < minEntryLen {
			continue
		}

		op := entry[0]
		keyLen := binary.LittleEndian.Uint32(entry[9:13])

		keyStart := 13
		keyEnd := keyStart + int(keyLen)
		if keyEnd > len(entry)-5 { // need at least valLen(4) + sentinel(1)
			continue
		}
		key := entry[keyStart:keyEnd]

		valLen := binary.LittleEndian.Uint32(entry[keyEnd:])
		valStart := keyEnd + 4
		valEnd := valStart + int(valLen)

		if valEnd > len(entry)-1 {
			continue
		}
		// sentinel check
		if entry[valEnd] != 0xFF {
			continue
		}

		switch op {
		case opSet:
			store.Set(key, entry[valStart:valEnd])
		case opDelete:
			store.Delete(key)
		}
	}

	// Truncate after replay.
	return w.Truncate()
}

func seqKey(seq uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, seq)
	return b
}
