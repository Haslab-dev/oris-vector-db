package planner

import (
	"testing"

	"github.com/hasdev/oris/engine/dense"
	"github.com/hasdev/oris/engine/dense/flat"
	"github.com/hasdev/oris/engine/metadata"
	"github.com/hasdev/oris/engine/sparse/bm25"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestPlanner() *Planner {
	return New(
		flat.New(dense.Config{Dimension: 4, Distance: "cosine"}),
		bm25.New(1.2, 0.75),
		metadata.New(),
		4,
	)
}

func seedDense(p *Planner) {
	for i := uint64(0); i < 5; i++ {
		vec := []float32{float32(i) / 5, 1 - float32(i)/5, 0, 0}
		p.dense.Insert(i, vec)
	}
}

func seedSparse(p *Planner) {
	texts := []string{"hello world", "hello foo", "bar baz", "foo bar", "hello bar"}
	for i, t := range texts {
		tokens := bm25.Tokenize(t)
		p.sparse.IndexDocument(uint64(i), tokens)
	}
}

func seedMeta(p *Planner) {
	meta := p.meta
	meta.IndexField(0, "tag", "a")
	meta.IndexField(1, "tag", "a")
	meta.IndexField(2, "tag", "b")
	meta.IndexField(3, "tag", "b")
	meta.IndexField(4, "tag", "a")
}

func TestDenseOnly(t *testing.T) {
	p := newTestPlanner()
	seedDense(p)

	results, err := p.Execute(Query{
		DenseVector: []float32{1, 0, 0, 0},
		TopK:        3,
		Mode:        DenseOnly,
	})
	require.NoError(t, err)
	require.Len(t, results, 3)
	// Closest to {1,0,0,0} is ID 4 ({0.8,0.2}) then ID 3 ({0.6,0.4})
	assert.Equal(t, uint64(4), results[0].ID)
	assert.Less(t, results[0].FinalScore, float32(0.1))
}

func TestDenseOnlyWithFilter(t *testing.T) {
	p := newTestPlanner()
	seedDense(p)
	seedMeta(p)

	tagA := &metadata.In{Field: "tag", Values: []string{"a"}}
	results, err := p.Execute(Query{
		DenseVector: []float32{1, 0, 0, 0},
		TopK:        10,
		Mode:        DenseOnly,
		Filter:      tagA,
	})
	require.NoError(t, err)
	// Only docs 0, 1, 4 have tag "a"
	for _, r := range results {
		assert.Contains(t, []uint64{0, 1, 4}, r.ID)
	}
}

func TestSparseOnly(t *testing.T) {
	p := newTestPlanner()
	seedSparse(p)

	results, err := p.Execute(Query{
		SparseTokens: bm25.Tokenize("hello"),
		TopK:         3,
		Mode:         SparseOnly,
	})
	require.NoError(t, err)
	require.Len(t, results, 3)
	// Docs 0, 1, 4 all score the same for "hello". Sort is stable enough for this check.
	assert.Equal(t, float32(0.5389965), results[0].SparseScore)
}

func TestSparseOnlyWithFilter(t *testing.T) {
	p := newTestPlanner()
	seedSparse(p)
	seedMeta(p)

	tagB := &metadata.In{Field: "tag", Values: []string{"b"}}
	results, err := p.Execute(Query{
		SparseTokens: bm25.Tokenize("foo"),
		TopK:         10,
		Mode:         SparseOnly,
		Filter:       tagB,
	})
	require.NoError(t, err)
	// Docs 1 (hello foo), 3 (foo bar). Tag "b" is on 2, 3. Only 3 matches both.
	for _, r := range results {
		assert.Equal(t, uint64(3), r.ID)
	}
}

func TestHybrid(t *testing.T) {
	p := newTestPlanner()
	seedDense(p)
	seedSparse(p)

	results, err := p.Execute(Query{
		DenseVector:  []float32{0.8, 0.2, 0, 0},
		SparseTokens: bm25.Tokenize("hello"),
		TopK:         5,
		Mode:         Hybrid,
		Alpha:        0.5,
	})
	require.NoError(t, err)
	require.Len(t, results, 5)
	t.Logf("Hybrid results: %+v", results)
}

func TestHybridWithFilter(t *testing.T) {
	p := newTestPlanner()
	seedDense(p)
	seedSparse(p)
	seedMeta(p)

	tagA := &metadata.In{Field: "tag", Values: []string{"a"}}
	results, err := p.Execute(Query{
		DenseVector:  []float32{0.5, 0.5, 0, 0},
		SparseTokens: bm25.Tokenize("hello"),
		TopK:         10,
		Mode:         Hybrid,
		Alpha:        0.5,
		Filter:       tagA,
	})
	require.NoError(t, err)
	for _, r := range results {
		assert.Contains(t, []uint64{0, 1, 4}, r.ID)
	}
}

func TestEmptyPlanner(t *testing.T) {
	p := New(nil, nil, metadata.New(), 4)

	results, err := p.Execute(Query{
		DenseVector: []float32{1, 0, 0, 0},
		TopK:        5,
		Mode:        DenseOnly,
	})
	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestExecuteDefaultMode(t *testing.T) {
	p := New(flat.New(dense.Config{Dimension: 4, Distance: "cosine"}), nil, metadata.New(), 4)
	p.dense.Insert(1, []float32{1, 0, 0, 0})

	// Empty mode defaults to Hybrid which falls back to DenseOnly.
	results, err := p.Execute(Query{
		DenseVector: []float32{0.9, 0.1, 0, 0},
		TopK:        5,
	})
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, uint64(1), results[0].ID)
}
