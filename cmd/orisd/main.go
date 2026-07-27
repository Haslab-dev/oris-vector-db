package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hasdev/oris/api"
)

func main() {
	workspace := flag.String("workspace", "./data", "path to Oris workspace")
	addr := flag.String("addr", ":8080", "HTTP server address")
	dim := flag.Int("dim", 128, "default vector dimension for new collections")
	seed := flag.Int("seed", 0, "seed N sample points into 'demo' collection")
	static := flag.String("static", "", "path to static web UI files to serve")
	flag.Parse()

	if err := os.MkdirAll(*workspace, 0755); err != nil {
		log.Fatalf("failed to create workspace: %v", err)
	}

	cfg := api.DefaultConfig("default", *dim)
	cfg.SegmentMaxPoints = 1000

	server := api.NewServer(*workspace, cfg)

	// Seed demo data if requested.
	if *seed > 0 {
		demoPath := filepath.Join(*workspace, "demo")
		if _, err := os.Stat(demoPath); os.IsNotExist(err) {
			col, err := api.Open(demoPath, cfg)
			if err != nil {
				log.Fatalf("failed to create demo collection: %v", err)
			}
			api.SeedData(col, *seed)
			col.Close()
			log.Printf("Demo collection seeded with %d points", *seed)
		} else {
			log.Printf("Demo collection already exists, skipping seed")
		}
	}

	mux := http.NewServeMux()

	// Serve API under /api/.
	mux.Handle("/api/", http.StripPrefix("/api", server))

	// Serve static web UI if provided.
	if *static != "" {
		mux.Handle("/", http.FileServer(http.Dir(*static)))
		log.Printf("Serving static UI from %s", *static)
	} else {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"service":"oris","version":"0.1.0"}`))
		})
	}

	log.Printf("Oris server starting on %s", *addr)
	log.Printf("Workspace: %s", *workspace)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
