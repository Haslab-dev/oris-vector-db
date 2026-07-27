// Package sparse provides sparse vector indexing and retrieval engines.
package sparse

// Engine is the interface for sparse (lexical) retrieval.
type Engine interface {
	// IndexDocument indexes a document with the given ID and text tokens.
	IndexDocument(id uint64, tokens []string) error

	// DeleteDocument removes a document from the index.
	DeleteDocument(id uint64) error

	// Search finds the top-k most relevant documents for the query tokens.
	Search(tokens []string, topK int) ([]Result, error)

	// Len returns the number of indexed documents.
	Len() int
}

// Result is a single search result.
type Result struct {
	ID    uint64
	Score float32
}

// Tokenizer splits text into tokens for indexing and search.
type Tokenizer interface {
	Tokenize(text string) []string
}
