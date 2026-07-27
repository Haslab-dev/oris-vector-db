// Package segment manages the segment lifecycle for Oris collections.
//
// Each segment is a self-contained unit holding dense, sparse, and metadata
// indexes plus payload storage.
package segment

import (
	"fmt"
	"time"
)

// State represents the lifecycle state of a segment.
type State int

const (
	StateMutable   State = iota // accepting writes
	StateSealing               // being sealed (no new writes)
	StateImmutable             // read-only, on disk
	StateCompacting            // being merged into a new segment
	StateDeleted               // removed after compaction
)

func (s State) String() string {
	switch s {
	case StateMutable:
		return "mutable"
	case StateSealing:
		return "sealing"
	case StateImmutable:
		return "immutable"
	case StateCompacting:
		return "compacting"
	case StateDeleted:
		return "deleted"
	default:
		return "unknown"
	}
}

// Metadata holds identifying and sizing information about a segment.
type Metadata struct {
	ID        uint64
	State     State
	DocCount  int
	MinID     uint64
	MaxID     uint64
	CreatedAt time.Time
	SealedAt  time.Time
	SizeBytes int64
	DenseDim  int
	Distance  string
}

// Segment is the interface implemented by both mutable and immutable segments.
type Segment interface {
	ID() uint64
	State() State
	Metadata() Metadata
	Insert(id uint64, dense []float32, sparseIndices []uint32, sparseValues []float32, payload []byte) error
	Delete(id uint64) error
	Len() int
	Seal(dir string) error
	Close() error
}

// ErrSegmentSealed is returned when trying to write to an immutable segment.
var ErrSegmentSealed = fmt.Errorf("segment is sealed")
