// Package metadata provides metadata filtering for Oris collections.
//
// The metadata engine maintains indexes for each metadata field and evaluates
// filter expressions against them, producing Roaring Bitmaps of matching
// document IDs.
package metadata

import (
	"fmt"

	"github.com/RoaringBitmap/roaring"
)

// Filter represents a metadata filter expression.
type Filter interface {
	// String returns a human-readable representation.
	String() string
}

// And is a conjunction of filters.
type And struct {
	Filters []Filter
}

func (a *And) String() string { return "AND" }

// Or is a disjunction of filters.
type Or struct {
	Filters []Filter
}

func (o *Or) String() string { return "OR" }

// Not negates a filter.
type Not struct {
	Filter Filter
}

func (n *Not) String() string { return "NOT" }

// In matches documents where the field's value is in the given set.
type In struct {
	Field  string
	Values []string
}

func (i *In) String() string { return fmt.Sprintf("IN %s", i.Field) }

// Range matches documents where the field's numeric value is in the range.
// If Min or Max is nil, that bound is open.
type Range struct {
	Field string
	Min   *float64
	Max   *float64
}

func (r *Range) String() string { return fmt.Sprintf("RANGE %s", r.Field) }

// Exists matches documents that have the given field.
type Exists struct {
	Field string
}

func (e *Exists) String() string { return fmt.Sprintf("EXISTS %s", e.Field) }

// Engine maintains metadata indexes and evaluates filter expressions.
type Engine struct {
	// fieldIndex maps field name -> value -> bitmap of doc IDs.
	fieldIndex map[string]map[string]*roaring.Bitmap
	// numericIndex maps field name -> sorted entry list for range queries.
	numericIndex map[string][]numericEntry
	// existsIndex maps field name -> bitmap of docs having that field.
	existsIndex map[string]*roaring.Bitmap
	// allDocs is the bitmap of all documents ever indexed.
	allDocs *roaring.Bitmap
}

type numericEntry struct {
	value float64
	docID uint64
}

// New creates a new metadata engine.
func New() *Engine {
	return &Engine{
		fieldIndex:   make(map[string]map[string]*roaring.Bitmap),
		numericIndex: make(map[string][]numericEntry),
		existsIndex:  make(map[string]*roaring.Bitmap),
		allDocs:      roaring.NewBitmap(),
	}
}

// IndexField indexes a metadata field for a document.
// value must be either a string or a float64 for range queries.
func (e *Engine) IndexField(docID uint64, field string, value interface{}) {
	e.allDocs.Add(uint32(docID))

	// exists index
	if e.existsIndex[field] == nil {
		e.existsIndex[field] = roaring.NewBitmap()
	}
	e.existsIndex[field].Add(uint32(docID))

	switch v := value.(type) {
	case string:
		if e.fieldIndex[field] == nil {
			e.fieldIndex[field] = make(map[string]*roaring.Bitmap)
		}
		if e.fieldIndex[field][v] == nil {
			e.fieldIndex[field][v] = roaring.NewBitmap()
		}
		e.fieldIndex[field][v].Add(uint32(docID))
	case float64:
		e.numericIndex[field] = append(e.numericIndex[field], numericEntry{value: v, docID: docID})
	}
}

// RemoveField removes a document's metadata from the indexes.
func (e *Engine) RemoveField(docID uint64, field string, value interface{}) {
	e.allDocs.Remove(uint32(docID))

	if bm, ok := e.existsIndex[field]; ok {
		bm.Remove(uint32(docID))
	}

	switch v := value.(type) {
	case string:
		if idx, ok := e.fieldIndex[field]; ok {
			if bm, ok := idx[v]; ok {
				bm.Remove(uint32(docID))
				if bm.IsEmpty() {
					delete(idx, v)
				}
			}
		}
	}
}

// Evaluate evaluates a filter expression and returns a bitmap of matching document IDs.
func (e *Engine) Evaluate(filter Filter) (*roaring.Bitmap, error) {
	switch f := filter.(type) {
	case *And:
		return e.evaluateAnd(f)
	case *Or:
		return e.evaluateOr(f)
	case *Not:
		return e.evaluateNot(f)
	case *In:
		return e.evaluateIn(f)
	case *Range:
		return e.evaluateRange(f)
	case *Exists:
		return e.evaluateExists(f)
	default:
		return nil, fmt.Errorf("unknown filter type: %T", filter)
	}
}

// AllDocs returns the bitmap of all indexed documents.
func (e *Engine) AllDocs() *roaring.Bitmap {
	return e.allDocs
}

func (e *Engine) evaluateAnd(f *And) (*roaring.Bitmap, error) {
	if len(f.Filters) == 0 {
		return e.allDocs.Clone(), nil
	}
	result, err := e.Evaluate(f.Filters[0])
	if err != nil {
		return nil, err
	}
	result = result.Clone()
	for _, filter := range f.Filters[1:] {
		bm, err := e.Evaluate(filter)
		if err != nil {
			return nil, err
		}
		result.And(bm)
	}
	return result, nil
}

func (e *Engine) evaluateOr(f *Or) (*roaring.Bitmap, error) {
	if len(f.Filters) == 0 {
		return roaring.NewBitmap(), nil
	}
	result, err := e.Evaluate(f.Filters[0])
	if err != nil {
		return nil, err
	}
	result = result.Clone()
	for _, filter := range f.Filters[1:] {
		bm, err := e.Evaluate(filter)
		if err != nil {
			return nil, err
		}
		result.Or(bm)
	}
	return result, nil
}

func (e *Engine) evaluateNot(f *Not) (*roaring.Bitmap, error) {
	inner, err := e.Evaluate(f.Filter)
	if err != nil {
		return nil, err
	}
	result := e.allDocs.Clone()
	result.AndNot(inner)
	return result, nil
}

func (e *Engine) evaluateIn(f *In) (*roaring.Bitmap, error) {
	idx, ok := e.fieldIndex[f.Field]
	if !ok {
		return roaring.NewBitmap(), nil
	}
	result := roaring.NewBitmap()
	for _, v := range f.Values {
		if bm, ok := idx[v]; ok {
			result.Or(bm)
		}
	}
	return result, nil
}

func (e *Engine) evaluateRange(f *Range) (*roaring.Bitmap, error) {
	entries, ok := e.numericIndex[f.Field]
	if !ok {
		return roaring.NewBitmap(), nil
	}
	result := roaring.NewBitmap()
	for _, entry := range entries {
		if f.Min != nil && entry.value < *f.Min {
			continue
		}
		if f.Max != nil && entry.value > *f.Max {
			continue
		}
		result.Add(uint32(entry.docID))
	}
	return result, nil
}

func (e *Engine) evaluateExists(f *Exists) (*roaring.Bitmap, error) {
	bm, ok := e.existsIndex[f.Field]
	if !ok {
		return roaring.NewBitmap(), nil
	}
	return bm.Clone(), nil
}
