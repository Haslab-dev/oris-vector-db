// Package api provides the public API for Oris.
package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// HTTPServer wraps an Oris collection with HTTP endpoints.
type HTTPServer struct {
	col *Collection
	mux *http.ServeMux
}

// NewHTTPServer creates a new HTTP server for the given collection.
func NewHTTPServer(col *Collection) *HTTPServer {
	s := &HTTPServer{
		col: col,
		mux: http.NewServeMux(),
	}
	s.mux.HandleFunc("/insert", s.handleInsert)
	s.mux.HandleFunc("/search", s.handleSearch)
	s.mux.HandleFunc("/delete", s.handleDelete)
	s.mux.HandleFunc("/count", s.handleCount)
	s.mux.HandleFunc("/snapshot", s.handleSnapshot)
	return s
}

// ServeHTTP implements http.Handler.
func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

type insertRequest struct {
	ID      uint64    `json:"id"`
	Dense   []float32 `json:"dense"`
	Payload string    `json:"payload,omitempty"`
}

type searchRequest struct {
	Query []float32 `json:"query"`
	TopK  int       `json:"topK"`
}

type searchResponse struct {
	Results []searchResult `json:"results"`
}

type searchResult struct {
	ID    uint64  `json:"id"`
	Score float32 `json:"score"`
}

func (s *HTTPServer) handleInsert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}
	var req insertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.col.Insert(req.ID, req.Dense, nil, nil, []byte(req.Payload)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"id":%d}`, req.ID)
}

func (s *HTTPServer) handleSearch(w http.ResponseWriter, r *http.Request) {
	var req searchRequest

	if r.Method == http.MethodGet {
		q := r.URL.Query()
		topK, _ := strconv.Atoi(q.Get("topK"))
		if topK <= 0 {
			topK = 10
		}
		queryStr := q.Get("query")
		if queryStr == "" {
			http.Error(w, "query required", http.StatusBadRequest)
			return
		}
		parts := strings.Split(queryStr, ",")
		query := make([]float32, len(parts))
		for i, p := range parts {
			v, _ := strconv.ParseFloat(strings.TrimSpace(p), 32)
			query[i] = float32(v)
		}
		req = searchRequest{Query: query, TopK: topK}
	} else if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "GET or POST required", http.StatusMethodNotAllowed)
		return
	}

	if req.TopK <= 0 {
		req.TopK = 10
	}

	results, err := s.col.Search(req.Query, req.TopK)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := searchResponse{Results: make([]searchResult, len(results))}
	for i, r := range results {
		resp.Results[i] = searchResult{ID: r.ID, Score: r.Score}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *HTTPServer) handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "valid id required", http.StatusBadRequest)
		return
	}
	if err := s.col.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"deleted":%d}`, id)
}

func (s *HTTPServer) handleCount(w http.ResponseWriter, r *http.Request) {
	cnt := s.col.Count()
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"count":%d}`, cnt)
}

func (s *HTTPServer) handleSnapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}
	snap, err := s.col.Snapshot(fmt.Sprintf("snap-%d", time.Now().Unix()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"snapshot":"%s"}`, snap.Name)
}
