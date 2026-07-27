// Package memory provides an in-memory implementation of the storage interface.
//
// Useful for testing and small-scale use cases where durability is not required.
package memory

import (
	"sort"
	"sync"

	"github.com/hasdev/oris/storage"
)

// Storage implements storage.Storage using an in-memory map.
type Storage struct {
	mu   sync.RWMutex
	data map[string][]byte
}

// New creates a new in-memory storage.
func New() *Storage {
	return &Storage{data: make(map[string][]byte)}
}

func (s *Storage) Get(key []byte) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[string(key)]
	if !ok {
		return nil, storage.ErrNotFound
	}
	return v, nil
}

func (s *Storage) Set(key, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[string(key)] = value
	return nil
}

func (s *Storage) Delete(key []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, string(key))
	return nil
}

func (s *Storage) NewBatch() storage.Batch {
	return &batch{store: s, ops: make(map[string][]byte)}
}

func (s *Storage) Snapshot() (storage.Snapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cp := make(map[string][]byte, len(s.data))
	for k, v := range s.data {
		cp[k] = v
	}
	return &snapshot{data: cp}, nil
}

func (s *Storage) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = nil
	return nil
}

// batch implements storage.Batch.
type batch struct {
	store *Storage
	ops   map[string][]byte
	dels  map[string]struct{}
}

func (b *batch) Set(key, value []byte) {
	b.ops[string(key)] = value
}

func (b *batch) Delete(key []byte) {
	b.ops[string(key)] = nil // nil marks deletion
}

func (b *batch) Commit() error {
	b.store.mu.Lock()
	defer b.store.mu.Unlock()
	for k, v := range b.ops {
		if v == nil {
			delete(b.store.data, k)
		} else {
			b.store.data[k] = v
		}
	}
	return nil
}

func (b *batch) Close() {
	b.ops = nil
}

// snapshot implements storage.Snapshot.
type snapshot struct {
	data map[string][]byte
}

func (s *snapshot) Get(key []byte) ([]byte, error) {
	v, ok := s.data[string(key)]
	if !ok {
		return nil, storage.ErrNotFound
	}
	return v, nil
}

func (s *snapshot) NewIterator(start, end []byte) (storage.Iterator, error) {
	sk, ek := string(start), string(end)
	var keys []string
	for k := range s.data {
		if k >= sk && k < ek {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return &iterator{data: s.data, keys: keys, pos: -1}, nil
}

func (s *snapshot) Close() {}

// iterator implements storage.Iterator.
type iterator struct {
	data map[string][]byte
	keys []string
	pos  int
}

func (it *iterator) Next() bool {
	it.pos++
	return it.pos < len(it.keys)
}

func (it *iterator) KeyValue() (key, value []byte, err error) {
	if it.pos < 0 || it.pos >= len(it.keys) {
		return nil, nil, storage.ErrSnapshotDone
	}
	k := it.keys[it.pos]
	return []byte(k), it.data[k], nil
}

func (it *iterator) Close() {}
