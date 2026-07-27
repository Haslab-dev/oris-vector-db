// Package bm25 implements a BM25 sparse retrieval engine.
package bm25

import (
	"math"
	"sort"
	"strings"
	"sync"
	"unicode"

	"github.com/hasdev/oris/engine/sparse"
)

// BM25 is a Best Matching 25 sparse retrieval engine.
type BM25 struct {
	mu         sync.RWMutex
	k1         float64
	b          float64
	docs       map[uint64]int       // doc ID -> doc length (in tokens)
	docText    map[uint64][]string  // doc ID -> stored tokens
	invIndex   map[string]posting   // term -> posting list
	docCount   int                  // total indexed documents
	avgDocLen  float64              // average document length
}

// posting is a list of (docID, termFrequency) pairs.
type posting struct {
	docIDs []uint64
	tfs    []int
}

// New creates a BM25 engine with the given parameters.
// Standard defaults: k1=1.2, b=0.75.
func New(k1, b float64) *BM25 {
	if k1 <= 0 {
		k1 = 1.2
	}
	if b <= 0 || b > 1 {
		b = 0.75
	}
	return &BM25{
		k1:       k1,
		b:        b,
		docs:     make(map[uint64]int),
		docText:  make(map[uint64][]string),
		invIndex: make(map[string]posting),
	}
}

// Tokenize splits text into lowercase tokens on whitespace and punctuation boundaries.
func Tokenize(text string) []string {
	var tokens []string
	var current []rune
	for _, r := range strings.ToLower(text) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			current = append(current, r)
		} else {
			if len(current) > 0 {
				tokens = append(tokens, string(current))
				current = nil
			}
		}
	}
	if len(current) > 0 {
		tokens = append(tokens, string(current))
	}
	return tokens
}

func (b *BM25) IndexDocument(id uint64, tokens []string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// If doc already exists, remove it first.
	if _, exists := b.docs[id]; exists {
		b.removeDoc(id)
	}

	b.docs[id] = len(tokens)
	b.docText[id] = tokens
	b.docCount++

	// Count term frequencies in this document.
	tf := make(map[string]int)
	for _, token := range tokens {
		tf[token]++
	}

	// Update inverted index.
	for term, count := range tf {
		p := b.invIndex[term]
		p.docIDs = append(p.docIDs, id)
		p.tfs = append(p.tfs, count)
		b.invIndex[term] = p
	}

	// Update average document length.
	totalLen := 0
	for _, l := range b.docs {
		totalLen += l
	}
	b.avgDocLen = float64(totalLen) / float64(b.docCount)

	return nil
}

func (b *BM25) DeleteDocument(id uint64) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.removeDoc(id)
	return nil
}

func (b *BM25) removeDoc(id uint64) {
	_, exists := b.docs[id]
	if !exists {
		return
	}
	delete(b.docs, id)
	delete(b.docText, id)
	b.docCount--
	for term, p := range b.invIndex {
		for i, docID := range p.docIDs {
			if docID == id {
				p.docIDs = append(p.docIDs[:i], p.docIDs[i+1:]...)
				p.tfs = append(p.tfs[:i], p.tfs[i+1:]...)
				break
			}
		}
		if len(p.docIDs) == 0 {
			delete(b.invIndex, term)
		} else {
			b.invIndex[term] = p
		}
	}

	// Recompute avg doc length.
	totalLen := 0
	for _, l := range b.docs {
		totalLen += l
	}
	if b.docCount > 0 {
		b.avgDocLen = float64(totalLen) / float64(b.docCount)
	} else {
		b.avgDocLen = 0
	}
}

func (b *BM25) Search(tokens []string, topK int) ([]sparse.Result, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.docCount == 0 {
		return nil, nil
	}

	// Score each document that matches at least one query term.
	scores := make(map[uint64]float64)
	for _, token := range tokens {
		p, exists := b.invIndex[token]
		if !exists {
			continue
		}

		// IDF: log(1 + (N - n + 0.5) / (n + 0.5))
		n := len(p.docIDs)
		idf := math.Log(1.0 + (float64(b.docCount)-float64(n)+0.5)/(float64(n)+0.5))

		for i, docID := range p.docIDs {
			tf := float64(p.tfs[i])
			docLen := float64(b.docs[docID])

			// BM25 score for this term.
			num := tf * (b.k1 + 1)
			denom := tf + b.k1*(1-b.b+ b.b*docLen/b.avgDocLen)
			scores[docID] += idf * num / denom
		}
	}

	// Sort by score descending.
	type scored struct {
		id    uint64
		score float64
	}
	results := make([]scored, 0, len(scores))
	for id, score := range scores {
		results = append(results, scored{id: id, score: score})
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	if len(results) > topK {
		results = results[:topK]
	}

	out := make([]sparse.Result, len(results))
	for i, r := range results {
		out[i] = sparse.Result{ID: r.id, Score: float32(r.score)}
	}
	return out, nil
}

func (b *BM25) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.docCount
}
