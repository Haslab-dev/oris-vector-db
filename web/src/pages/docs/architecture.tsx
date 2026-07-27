import { Card, CardContent } from "../../components/ui/card"

export default function ArchitecturePage() {
  return (
    <div className="space-y-8 max-w-3xl">
      <div>
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Architecture</h1>
        <p className="text-sm text-zinc-500 mt-1">How Oris is organized and how data flows through the system.</p>
      </div>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Design Principles</h2>
        <Card>
          <CardContent className="p-4 space-y-4">
            {[
              { title: "Engine-First Architecture", desc: "Every subsystem (dense, sparse, metadata, storage) is independently replaceable behind a stable Go interface. Want DiskANN instead of HNSW? Implement the dense.Engine interface. Want RocksDB instead of Pebble? Implement the storage.Storage interface." },
              { title: "Immutable Segments", desc: "Data is written into mutable segments that are sealed into immutable segments. This optimizes reads, simplifies concurrency, and enables background compaction without locking." },
              { title: "Storage/Index Separation", desc: "Persistence and ANN indexes evolve independently. The storage layer doesn't know about HNSW graphs, and the index layer doesn't know about WALs or snapshots." },
              { title: "Embedded by Default", desc: "Oris is a library, not a server. Import it, open a collection, insert and search — no daemon, no cluster, no configuration required." },
              { title: "Incremental Evolution", desc: "New retrieval algorithms are added as plugins behind existing interfaces, not architectural rewrites. The flat engine co-exists with HNSW, and both implement the same dense.Engine interface." },
            ].map((p) => (
              <div key={p.title}>
                <h3 className="text-sm font-medium text-zinc-900 dark:text-zinc-50 mb-1">{p.title}</h3>
                <p className="text-sm text-zinc-600 dark:text-zinc-400">{p.desc}</p>
              </div>
            ))}
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">System Architecture</h2>
        <Card>
          <CardContent className="p-4">
            <pre className="text-xs font-mono bg-zinc-900 text-zinc-100 p-4 rounded-md overflow-x-auto leading-relaxed">
{`                     Application
                           │
                    Public API (api package)
                           │
                  Query Planner Engine
                           │
      ┌─────────────┬──────────────┬──────────────┐
      ▼             ▼              ▼              ▼
 Dense Engine   Sparse Engine   Metadata      Distance
  (HNSW/Flat)     (BM25)        (Roaring)     (SIMD)
      │             │              │
      └─────────────┼──────────────┘
                    ▼
              Segment Manager
                    │
      ┌─────────────┼─────────────┐
      ▼             ▼             ▼
 Mutable      Immutable      Immutable
 Segment        Segment        Segment
      │
      ▼
 Storage Engine (Pebble / In-Memory)
   - WAL for crash recovery
   - Snapshots (BLAKE3 verification)
   - Key-value store per collection`}
            </pre>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Query Flow</h2>
        <Card>
          <CardContent className="p-4">
            <pre className="text-xs font-mono bg-zinc-900 text-zinc-100 p-4 rounded-md overflow-x-auto leading-relaxed">
{`User Query
    │
    ▼
Planner
    │
    ├── Mode: DenseOnly   → Dense Search → Merge → Rank
    ├── Mode: SparseOnly  → Sparse Search → Merge → Rank
    └── Mode: Hybrid      → Dense Search ─┐
                           → Sparse Search ─┤
                                            ▼
                                      Normalize scores
                                      Alpha weighting
                                      Merge & Rank
    │
    ▼
Metadata Filter (optional post-filter)
    │
    ▼
Results (sorted by score)`}
            </pre>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Package Map</h2>
        <Card>
          <CardContent className="p-4">
            <pre className="text-xs font-mono bg-zinc-900 text-zinc-100 p-4 rounded-md overflow-x-auto leading-relaxed">
{`oris/
├── api/              # Public API: Collection, Open, Insert, Search
├── cmd/oris/         # CLI: init, bench, inspect, compact
├── engine/
│   ├── dense/        # Dense Engine interface + HNSW + Flat
│   │   ├── hnsw/     # Hierarchical Navigable Small World graph
│   │   ├── flat/     # Brute-force fallback
│   │   └── factory/  # Auto-select by size
│   ├── sparse/       # Sparse Engine interface
│   │   └── bm25/     # BM25 scorer + inverted index
│   ├── metadata/     # Filter AST + Roaring Bitmap index
│   ├── distance/     # Distance kernels (Cosine/Dot/Euclidean)
│   ├── planner/      # Query planner (dense/sparse/hybrid)
│   └── segment/      # Segment lifecycle & manager
├── storage/          # Storage interface
│   ├── pebble/       # Pebble LSM implementation
│   ├── memory/       # In-memory (testing)
│   ├── wal/          # Write-ahead log
│   └── snapshot/     # BLAKE3-verified snapshots
├── internal/         # Binary encoding & utilities
└── web/              # React + Vite web UI`}
            </pre>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Storage Layout on Disk</h2>
        <Card>
          <CardContent className="p-4">
            <pre className="text-xs font-mono bg-zinc-900 text-zinc-100 p-4 rounded-md overflow-x-auto leading-relaxed">
{`workspace/
└── collections/
    ├── data/              # Pebble database (key-value store)
    │   ├── MANIFEST-*
    │   ├── OPTIONS-*
    │   └── *.sst
    ├── wal/               # Write-ahead log (separate Pebble DB)
    ├── snapshots/         # Point-in-time BLAKE3 snapshots
    │   └── snap-001/
    │       └── data.snap
    └── segments/          # Sealed immutable segments
        ├── segment_000001/
        │   ├── dense.idx     # Vector data
        │   └── payload.dat   # Payload storage
        └── segment_000002/`}
            </pre>
          </CardContent>
        </Card>
      </section>
    </div>
  )
}
