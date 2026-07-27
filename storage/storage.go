// Package storage defines the storage interface for Oris.
//
// Every storage engine must implement this interface, enabling Oris to
// support Pebble, in-memory, or custom backends without changing the
// retrieval engine.
package storage

import "errors"

// ErrNotFound is returned when a key does not exist.
var ErrNotFound = errors.New("key not found")

// ErrSnapshotDone is returned when an iterator reaches the end of a snapshot.
var ErrSnapshotDone = errors.New("snapshot iteration complete")

// Storage is the interface for persistent key-value storage.
type Storage interface {
	// Get retrieves the value for the given key.
	// Returns ErrNotFound if the key does not exist.
	Get(key []byte) ([]byte, error)

	// Set stores a key-value pair.
	Set(key, value []byte) error

	// Delete removes the key-value pair for the given key.
	// Returns nil if the key does not exist.
	Delete(key []byte) error

	// NewBatch creates a new atomic batch operation.
	NewBatch() Batch

	// Snapshot creates a point-in-time read-only view of the storage.
	Snapshot() (Snapshot, error)

	// Close closes the storage engine, flushing any pending writes.
	Close() error
}

// Batch represents an atomic batch of write operations.
type Batch interface {
	Set(key, value []byte)
	Delete(key []byte)
	Commit() error
	Close()
}

// Snapshot represents a point-in-time read-only view of storage.
type Snapshot interface {
	Get(key []byte) ([]byte, error)
	NewIterator(start, end []byte) (Iterator, error)
	Close()
}

// Iterator iterates over key-value pairs in a snapshot.
type Iterator interface {
	Next() bool
	KeyValue() (key, value []byte, err error)
	Close()
}
