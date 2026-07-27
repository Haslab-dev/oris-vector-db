# Oris

**Next-generation Embedded Vector Retrieval Engine**

Oris is a lightweight, high-performance, embedded vector retrieval engine for modern AI applications. Import it, open a collection, insert vectors, search — no server, no cluster, no external dependencies.

## Quick Start

```bash
# Install
go get github.com/hasdev/oris

# Run the demo server
make server

# In another terminal, start the web UI
make web
```

Open http://localhost:5173 in your browser.

## Usage

```go
package main

import (
    "fmt"
    "github.com/hasdev/oris/api"
)

func main() {
    cfg := api.DefaultConfig("my-collection", 128)
    col, _ := api.Open("./data", cfg)
    defer col.Close()

    // Insert points.
    for i := uint64(0); i < 1000; i++ {
        vec := make([]float32, 128)
        vec[i%128] = float32(i) / 1000
        col.Insert(i, vec, nil, nil, []byte(fmt.Sprintf("point-%d", i)))
    }

    // Search.
    query := make([]float32, 128)
    query[0] = 0.5
    results, _ := col.Search(query, 10)
    for _, r := range results {
        fmt.Printf("ID=%d Score=%.4f\n", r.ID, r.Score)
    }
}
```

## Architecture

```
                    Application
                         │
                  Public API (api package)
                         │
                Query Planner Engine
                         │
     ┌─────────────┬─────┴──────┬──────────────┐
     ▼             ▼            ▼              ▼
Dense Engine   Sparse Engine  Metadata     Distance
 (HNSW/Flat)     (BM25)      (Roaring)     (SIMD)
     │             │            │
     └─────────────┼────────────┘
                   ▼
             Segment Manager
                   │
     ┌─────────────┼─────────────┐
     ▼             ▼             ▼
Mutable      Immutable      Immutable
 Segment       Segment       Segment
     │
     ▼
Storage Engine (Pebble / In-Memory)
  - WAL for crash recovery
  - Snapshots (BLAKE3 verification)
```

## Engines

| Engine | Technology | Description |
|--------|-----------|-------------|
| Dense | HNSW + Flat | Approximate nearest neighbor search. Auto-selects Flat for < 1K vectors. |
| Sparse | BM25 | Lexical retrieval with inverted index, TF-IDF scoring. |
| Metadata | Roaring Bitmap | Filter expressions (AND, OR, NOT, IN, Range, Exists). |
| Distance | Cosine / Dot / Euclidean | Auto-detects NEON, AVX2, AVX-512. Pure Go fallback. |
| Storage | Pebble (default) + In-Memory | Pluggable key-value interface. WAL + BLAKE3 snapshots. |
| Segment | Mutable → Immutable → Merge | Lucene-inspired segment lifecycle with background compaction. |
| Planner | Query routing | Dense-only, sparse-only, or hybrid search with alpha weighting. |

## Performance Targets

| Metric | Target |
|--------|--------|
| Insert latency | < 1 ms per point |
| Search latency (100K vectors) | < 10 ms |
| Search latency (1M vectors) | < 30 ms |
| Recall@10 | ≥ 0.95 |
| Startup time | < 100 ms (small collections) |

## Project Structure

```
oris/
├── api/              # Public API: Collection, Open, Insert, Search, HTTP server
├── cmd/
│   ├── oris/         # CLI: init, bench, inspect, compact
│   └── orisd/        # Daemon: serves API + web UI, optional seed data
├── engine/
│   ├── dense/        # HNSW + Flat + factory auto-select
│   ├── sparse/       # BM25 inverted index
│   ├── metadata/     # Filter AST + Roaring Bitmap
│   ├── distance/     # SIMD-aware distance kernels
│   ├── planner/      # Query planner (dense/sparse/hybrid)
│   └── segment/      # Segment lifecycle & compaction
├── storage/          # Storage interface + Pebble + WAL + snapshots
├── internal/         # Binary encoding utilities
├── web/              # React + Vite + Tailwind web UI
├── examples/         # Usage examples
├── bench/            # Benchmarks
├── docs/             # Architecture docs
└── Makefile          # `make server` + `make web`
```

## Web UI

The web interface provides a full-featured UI for managing and querying collections:

- **Dashboard** — real-time stats across all collections
- **Playground** — end-to-end: create collections, insert points, search, delete
- **Collections** — browse collections and their metadata
- **Search** — vector search with result details
- **Performance** — engine metrics
- **Documentation** — full architecture and API reference
- **Settings** — instance configuration

```bash
# Start both server and web UI
make
```

## Development

```bash
# Run tests
go test ./...

# Build the daemon
go build -o bin/orisd ./cmd/orisd

# Start API server with 100K seeded demo data
./bin/orisd -workspace ./data -addr :8080 -seed 100000 -static ./web/dist

# Web UI dev mode (proxies /api to :8080)
cd web && bun run dev
```

## License

MIT
