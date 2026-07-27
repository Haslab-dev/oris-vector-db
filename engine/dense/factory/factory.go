// Package factory provides auto-selection between flat and HNSW dense engines.
package factory

import (
	"github.com/hasdev/oris/engine/dense"
	"github.com/hasdev/oris/engine/dense/flat"
	"github.com/hasdev/oris/engine/dense/hnsw"
)

// New creates the best dense engine for the given config and expected size.
// Returns a Flat index for small collections (< SmallThreshold), HNSW otherwise.
func New(cfg dense.Config, expectedSize int) dense.Engine {
	if expectedSize < dense.SmallThreshold {
		return flat.New(cfg)
	}
	return hnsw.New(cfg)
}

// NewFlat always returns a flat (brute-force) engine.
func NewFlat(cfg dense.Config) dense.Engine {
	return flat.New(cfg)
}

// NewHNSW always returns an HNSW engine.
func NewHNSW(cfg dense.Config) dense.Engine {
	return hnsw.New(cfg)
}
