// Package api provides the public API for interacting with Oris collections.
//
// This is the primary entry point for embedding Oris into Go applications.
package api

import "errors"

// CollectionConfig holds configuration for creating a collection.
type CollectionConfig struct {
	Name             string
	Dimension        int
	Distance         string // "cosine", "dot", "euclidean"
	MutableSegments  int    // max mutable segments before flush
	SegmentMaxPoints int    // max points per segment before seal
}

// DefaultConfig returns a default collection configuration.
func DefaultConfig(name string, dimension int) CollectionConfig {
	return CollectionConfig{
		Name:             name,
		Dimension:        dimension,
		Distance:         "cosine",
		MutableSegments:  1,
		SegmentMaxPoints: 10000,
	}
}

// Result represents a single search result.
type Result struct {
	ID      uint64
	Score   float32
	Payload []byte
}

// ErrCollectionNotFound is returned when a collection does not exist.
var ErrCollectionNotFound = errors.New("collection not found")

// ErrCollectionExists is returned when creating a duplicate collection.
var ErrCollectionExists = errors.New("collection already exists")
