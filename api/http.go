package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Server is an HTTP server that manages multiple collections.
type Server struct {
	mu          sync.RWMutex
	collections map[string]*Collection
	workspace   string
	cfg         CollectionConfig
	mux         *http.ServeMux
}

// NewServer creates a new HTTP API server managing collections in a workspace.
func NewServer(workspace string, cfg CollectionConfig) *Server {
	s := &Server{
		collections: make(map[string]*Collection),
		workspace:   workspace,
		cfg:         cfg,
		mux:         http.NewServeMux(),
	}
	s.mux.HandleFunc("/collections", s.handleCollections)
	s.mux.HandleFunc("/collections/", s.handleCollection)
	s.mux.HandleFunc("/stats", s.handleStats)
	s.mux.HandleFunc("/insert", s.handleInsert)
	s.mux.HandleFunc("/search", s.handleSearch)
	s.mux.HandleFunc("/delete", s.handleDelete)
	s.mux.HandleFunc("/count", s.handleCount)
	s.mux.HandleFunc("/snapshot", s.handleSnapshot)
	return s
}

// ServeHTTP implements http.Handler with CORS.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	s.mux.ServeHTTP(w, r)
}

// ListenAndServe starts the HTTP server.
func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s)
}

type collectionInfo struct {
	Name      string `json:"name"`
	Vectors   int    `json:"vectors"`
	Dimension int    `json:"dimension"`
	Distance  string `json:"distance"`
	Segments  int    `json:"segments"`
	SizeBytes int64  `json:"sizeBytes"`
}

type statsResponse struct {
	Count        int               `json:"count"`
	TotalVectors int               `json:"totalVectors"`
	StorageMB    float64           `json:"storageMB"`
	Collections  []collectionInfo  `json:"collections"`
}

func (s *Server) handleCollections(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listCollections(w, r)
	case http.MethodPost:
		s.createCollection(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) listCollections(w http.ResponseWriter, r *http.Request) {
	entries, err := os.ReadDir(s.workspace)
	if err != nil {
		jsonError(w, err.Error())
		return
	}

	var infos []collectionInfo
	s.mu.RLock()
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		col, ok := s.collections[name]
		if !ok {
			info := collectionInfo{Name: name}
			infos = append(infos, info)
			continue
		}
		info := collectionInfo{
			Name:      name,
			Vectors:   col.Count(),
			Dimension: col.cfg.Dimension,
			Distance:  col.cfg.Distance,
		}
		// Count segments (directories under the collection's segments folder)
		segDir := filepath.Join(s.workspace, name, "segments")
		if segEntries, err := os.ReadDir(segDir); err == nil {
			info.Segments = len(segEntries)
		}
		infos = append(infos, info)
	}
	s.mu.RUnlock()

	jsonResp(w, infos)
}

func (s *Server) createCollection(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Dimension   int    `json:"dimension"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request: "+err.Error())
		return
	}
	if req.Name == "" {
		jsonError(w, "name is required")
		return
	}
	if req.Dimension <= 0 {
		req.Dimension = s.cfg.Dimension
	}

	colPath := filepath.Join(s.workspace, req.Name)
	if _, err := os.Stat(colPath); err == nil {
		jsonError(w, "collection already exists")
		return
	}

	cfg := s.cfg
	cfg.Dimension = req.Dimension
	col, err := Open(colPath, cfg)
	if err != nil {
		jsonError(w, err.Error())
		return
	}

	s.mu.Lock()
	s.collections[req.Name] = col
	s.mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	jsonResp(w, map[string]string{"name": req.Name})
}

func (s *Server) handleCollection(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/collections/")
	if name == "" {
		s.listCollections(w, r)
		return
	}

	s.mu.RLock()
	col, ok := s.collections[name]
	s.mu.RUnlock()

	if !ok {
		// Try to open it.
		colPath := filepath.Join(s.workspace, name)
		if _, err := os.Stat(colPath); os.IsNotExist(err) {
			jsonError(w, "collection not found")
			return
		}
		var openErr error
		col, openErr = Open(colPath, s.cfg)
		if openErr != nil {
			jsonError(w, openErr.Error())
			return
		}
		s.mu.Lock()
		s.collections[name] = col
		s.mu.Unlock()
	}

	switch r.Method {
	case http.MethodGet:
		jsonResp(w, collectionInfo{
			Name:      name,
			Vectors:   col.Count(),
			Dimension: col.cfg.Dimension,
			Distance:  col.cfg.Distance,
		})
	case http.MethodPost:
		// Could handle insert/delete/search sub-routes here.
		jsonResp(w, map[string]string{"name": name, "status": "ready"})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	entries, _ := os.ReadDir(s.workspace)
	var totalVectors int
	var totalSize int64
	var infos []collectionInfo

	s.mu.RLock()
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		col, ok := s.collections[name]
		if !ok {
			continue
		}
		cnt := col.Count()
		totalVectors += cnt
		info := collectionInfo{
			Name:      name,
			Vectors:   cnt,
			Dimension: col.cfg.Dimension,
			Distance:  col.cfg.Distance,
		}
		if fi, err := os.Stat(filepath.Join(s.workspace, name)); err == nil {
			totalSize += fi.Size()
		}
		infos = append(infos, info)
	}
	s.mu.RUnlock()

	jsonResp(w, statsResponse{
		Count:        len(infos),
		TotalVectors: totalVectors,
		StorageMB:    float64(totalSize) / 1024 / 1024,
		Collections:  infos,
	})
}

func (s *Server) handleInsert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Collection string    `json:"collection"`
		ID         uint64    `json:"id"`
		Dense      []float32 `json:"dense"`
		Payload    string    `json:"payload,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request: "+err.Error())
		return
	}

	col, err := s.getCollection(req.Collection)
	if err != nil {
		jsonError(w, err.Error())
		return
	}

	if err := col.Insert(req.ID, req.Dense, nil, nil, []byte(req.Payload)); err != nil {
		jsonError(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	jsonResp(w, map[string]uint64{"id": req.ID})
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Collection string    `json:"collection"`
		Query      []float32 `json:"query"`
		TopK       int       `json:"topK"`
	}

	if r.Method == http.MethodGet {
		q := r.URL.Query()
		req.Collection = q.Get("collection")
		req.TopK, _ = strconv.Atoi(q.Get("topK"))
		if req.TopK <= 0 {
			req.TopK = 10
		}
		queryStr := q.Get("query")
		if queryStr != "" {
			parts := strings.Split(queryStr, ",")
			req.Query = make([]float32, len(parts))
			for i, p := range parts {
				v, _ := strconv.ParseFloat(strings.TrimSpace(p), 32)
				req.Query[i] = float32(v)
			}
		}
	} else if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			jsonError(w, "invalid request: "+err.Error())
			return
		}
	} else {
		http.Error(w, "GET or POST required", http.StatusMethodNotAllowed)
		return
	}

	if req.TopK <= 0 {
		req.TopK = 10
	}

	col, err := s.getCollection(req.Collection)
	if err != nil {
		jsonError(w, err.Error())
		return
	}

	results, err := col.Search(req.Query, req.TopK)
	if err != nil {
		jsonError(w, err.Error())
		return
	}

	type searchResult struct {
		ID    uint64  `json:"id"`
		Score float32 `json:"score"`
	}
	resp := make([]searchResult, len(results))
	for i, r := range results {
		resp[i] = searchResult{ID: r.ID, Score: r.Score}
	}
	jsonResp(w, resp)
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Collection string `json:"collection"`
		ID         uint64 `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request: "+err.Error())
		return
	}
	col, err := s.getCollection(req.Collection)
	if err != nil {
		jsonError(w, err.Error())
		return
	}
	if err := col.Delete(req.ID); err != nil {
		jsonError(w, err.Error())
		return
	}
	jsonResp(w, map[string]uint64{"deleted": req.ID})
}

func (s *Server) handleCount(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("collection")
	col, err := s.getCollection(name)
	if err != nil {
		jsonError(w, err.Error())
		return
	}
	jsonResp(w, map[string]int{"count": col.Count()})
}

func (s *Server) handleSnapshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}
	name := r.URL.Query().Get("collection")
	col, err := s.getCollection(name)
	if err != nil {
		jsonError(w, err.Error())
		return
	}
	snap, err := col.Snapshot(fmt.Sprintf("snap-%d", time.Now().Unix()))
	if err != nil {
		jsonError(w, err.Error())
		return
	}
	jsonResp(w, map[string]string{"snapshot": snap.Name})
}

func (s *Server) getCollection(name string) (*Collection, error) {
	if name == "" {
		return nil, fmt.Errorf("collection name is required")
	}
	s.mu.RLock()
	col, ok := s.collections[name]
	s.mu.RUnlock()
	if ok {
		return col, nil
	}

	// Try to open.
	colPath := filepath.Join(s.workspace, name)
	if _, err := os.Stat(colPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("collection '%s' not found", name)
	}
	col, err := Open(colPath, s.cfg)
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	s.collections[name] = col
	s.mu.Unlock()
	return col, nil
}

func jsonResp(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// SeedData inserts sample data into a collection for demo purposes.
func SeedData(col *Collection, n int) {
	for i := uint64(0); i < uint64(n); i++ {
		vec := make([]float32, col.cfg.Dimension)
		vec[i%uint64(col.cfg.Dimension)] = float32(i) / float32(n)
		if err := col.Insert(i, vec, nil, nil, []byte(fmt.Sprintf("point-%d", i))); err != nil {
			log.Printf("seed insert error: %v", err)
		}
	}
	log.Printf("Seeded %d points into collection", n)
}
