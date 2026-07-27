package bm25

import (
	"testing"

	"github.com/hasdev/oris/engine/sparse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenizer(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"hello world", []string{"hello", "world"}},
		{"Hello, World!", []string{"hello", "world"}},
		{"  spaced  out  ", []string{"spaced", "out"}},
		{"", nil},
		{"one", []string{"one"}},
		{"go123 is #1!", []string{"go123", "is", "1"}},
	}

	for _, tc := range tests {
		got := Tokenize(tc.input)
		assert.Equal(t, tc.want, got, "Tokenize(%q)", tc.input)
	}
}

func TestIndexAndSearch(t *testing.T) {
	b := New(1.2, 0.75)

	err := b.IndexDocument(1, Tokenize("the quick brown fox"))
	require.NoError(t, err)

	err = b.IndexDocument(2, Tokenize("jumps over the lazy dog"))
	require.NoError(t, err)

	err = b.IndexDocument(3, Tokenize("the brown fox runs fast"))
	require.NoError(t, err)

	assert.Equal(t, 3, b.Len())

	// Search for "brown fox" — docs 1 and 3 match.
	results, err := b.Search(Tokenize("brown fox"), 5)
	require.NoError(t, err)
	require.Len(t, results, 2)

	// Doc 1 and 3 both have "brown" and "fox". Doc 3 has 5 tokens vs doc1 has 4,
	// so doc1 should score slightly higher (shorter doc = higher BM25 score).
	assert.Equal(t, uint64(1), results[0].ID)
	assert.Equal(t, uint64(3), results[1].ID)
	assert.Greater(t, results[0].Score, float32(0))
}

func TestSearchReturnsNoResults(t *testing.T) {
	b := New(1.2, 0.75)
	b.IndexDocument(1, Tokenize("hello world"))

	results, err := b.Search(Tokenize("nonexistent"), 5)
	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestSearchEmptyIndex(t *testing.T) {
	b := New(1.2, 0.75)
	results, err := b.Search(Tokenize("anything"), 5)
	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestDeleteDocument(t *testing.T) {
	b := New(1.2, 0.75)

	b.IndexDocument(1, Tokenize("hello world"))
	b.IndexDocument(2, Tokenize("hello foo"))
	assert.Equal(t, 2, b.Len())

	err := b.DeleteDocument(1)
	require.NoError(t, err)
	assert.Equal(t, 1, b.Len())

	// Doc 2 should still be searchable.
	results, err := b.Search(Tokenize("hello"), 5)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, uint64(2), results[0].ID)
}

func TestBM25ScoreOrder(t *testing.T) {
	b := New(1.2, 0.75)

	// Doc 1: matches query exactly but shorter.
	b.IndexDocument(1, Tokenize("cat dog"))
	// Doc 2: matches query but is much longer (lower BM25).
	b.IndexDocument(2, Tokenize("cat dog bird fish mouse elephant bear wolf"))
	// Doc 3: no match.
	b.IndexDocument(3, Tokenize("xyz abc"))

	results, err := b.Search(Tokenize("cat dog"), 5)
	require.NoError(t, err)
	require.Len(t, results, 2)

	// Shorter doc should rank higher.
	assert.Equal(t, uint64(1), results[0].ID)
	assert.Equal(t, uint64(2), results[1].ID)
	assert.Greater(t, results[0].Score, results[1].Score)
}

func TestTermFrequencyEffect(t *testing.T) {
	b := New(1.2, 0.75)

	// Doc 1: mentions "go" once.
	b.IndexDocument(1, Tokenize("go is a language"))
	// Doc 2: mentions "go" three times.
	b.IndexDocument(2, Tokenize("go go go"))
	// Doc 3: no mention.
	b.IndexDocument(3, Tokenize("python rust"))

	results, err := b.Search(Tokenize("go"), 5)
	require.NoError(t, err)
	require.Len(t, results, 2)

	// Doc 2 has higher TF -> higher score.
	assert.Equal(t, uint64(2), results[0].ID)
	assert.Equal(t, uint64(1), results[1].ID)
}

func TestInterface(t *testing.T) {
	b := New(1.2, 0.75)
	var _ sparse.Engine = b
}

func TestReindexAfterDelete(t *testing.T) {
	b := New(1.2, 0.75)

	b.IndexDocument(1, Tokenize("hello world"))
	b.DeleteDocument(1)
	b.IndexDocument(1, Tokenize("hello world updated"))

	assert.Equal(t, 1, b.Len())

	results, err := b.Search(Tokenize("updated"), 5)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, uint64(1), results[0].ID)
}
