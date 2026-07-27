// Package pebble provides a Pebble-based implementation of the storage interface.
package pebble

import (
	"github.com/cockroachdb/pebble"

	"github.com/hasdev/oris/storage"
)

// Storage implements storage.Storage using Pebble LSM storage.
type Storage struct {
	db *pebble.DB
}

// Options holds configuration for opening a Pebble storage.
type Options struct {
	Path string
}

// Open opens or creates a Pebble database at the given path.
func Open(opts Options) (*Storage, error) {
	db, err := pebble.Open(opts.Path, &pebble.Options{})
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Get(key []byte) ([]byte, error) {
	v, closer, err := s.db.Get(key)
	if err != nil {
		if err == pebble.ErrNotFound {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}
	defer closer.Close()
	buf := make([]byte, len(v))
	copy(buf, v)
	return buf, nil
}

func (s *Storage) Set(key, value []byte) error {
	return s.db.Set(key, value, pebble.Sync)
}

func (s *Storage) Delete(key []byte) error {
	if err := s.db.Delete(key, pebble.Sync); err != nil {
		return err
	}
	return nil
}

func (s *Storage) NewBatch() storage.Batch {
	b := s.db.NewBatch()
	return &batch{batch: b}
}

func (s *Storage) Snapshot() (storage.Snapshot, error) {
	snap := s.db.NewSnapshot()
	return &snapshot{snap: snap}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

// batch implements storage.Batch.
type batch struct {
	batch *pebble.Batch
}

func (b *batch) Set(key, value []byte) {
	b.batch.Set(key, value, pebble.Sync)
}

func (b *batch) Delete(key []byte) {
	b.batch.Delete(key, pebble.Sync)
}

func (b *batch) Commit() error {
	return b.batch.Commit(pebble.Sync)
}

func (b *batch) Close() {
	b.batch.Close()
}

// snapshot implements storage.Snapshot.
type snapshot struct {
	snap *pebble.Snapshot
}

func (s *snapshot) Get(key []byte) ([]byte, error) {
	v, closer, err := s.snap.Get(key)
	if err != nil {
		if err == pebble.ErrNotFound {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}
	defer closer.Close()
	buf := make([]byte, len(v))
	copy(buf, v)
	return buf, nil
}

func (s *snapshot) NewIterator(start, end []byte) (storage.Iterator, error) {
	iter, err := s.snap.NewIter(&pebble.IterOptions{
		LowerBound: start,
		UpperBound: end,
	})
	if err != nil {
		return nil, err
	}
	return &iterator{iter: iter}, nil
}

func (s *snapshot) Close() {
	s.snap.Close()
}

// iterator implements storage.Iterator.
type iterator struct {
	iter *pebble.Iterator
}

func (it *iterator) Next() bool {
	return it.iter.Next()
}

func (it *iterator) KeyValue() (key, value []byte, err error) {
	if !it.iter.Valid() {
		return nil, nil, storage.ErrSnapshotDone
	}
	key = it.iter.Key()
	value, err = it.iter.ValueAndErr()
	if err != nil {
		return nil, nil, err
	}
	return key, value, nil
}

func (it *iterator) Close() {
	it.iter.Close()
}
