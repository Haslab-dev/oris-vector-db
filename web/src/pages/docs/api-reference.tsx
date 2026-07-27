import { Card, CardContent } from "../../components/ui/card"

const sections = [
  {
    title: "Opening a Collection",
    code: `import "github.com/hasdev/oris/api"

// Configure and open a collection.
cfg := api.DefaultConfig("products", 768)
col, err := api.Open("./data/products", cfg)
if err != nil { panic(err) }
defer col.Close()`,
    desc: "Opens an existing collection or creates a new one at the given path. The config specifies the vector dimension and distance metric.",
  },
  {
    title: "Inserting Points",
    code: `// Insert a single point.
err := col.Insert(
    1,                           // unique ID
    []float32{0.1, 0.2, ...},    // dense vector (768 dims)
    []uint32{5, 10},             // sparse indices (optional)
    []float32{0.5, 0.3},         // sparse values (optional)
    []byte("my payload"),        // payload (optional)
)

// Batch insert.
batch := api.NewBatch()
batch.Add(1, vec1, nil, nil, payload1)
batch.Add(2, vec2, nil, nil, payload2)
batch.Add(3, vec3, nil, nil, payload3)
err = batch.Execute(col)`,
    desc: "Points are added to the mutable segment. The segment auto-seals when it reaches capacity, and a new mutable segment is created.",
  },
  {
    title: "Searching",
    code: `// Dense search (semantic).
results, err := col.Search(queryVector, 10)
for _, r := range results {
    fmt.Printf("ID=%d Score=%.4f\\n", r.ID, r.Score)
}

// Search with metadata filter.
filter := &metadata.In{
    Field: "category",
    Values: []string{"electronics", "gadgets"},
}
filteredResults, err := col.SearchWithFilter(queryVector, 10, filter)`,
    desc: "Search returns results sorted by distance (ascending). Lower scores mean more similar. The filter narrows results to matching metadata.",
  },
  {
    title: "Managing Points",
    code: `// Update a point (delete + insert).
col.Update(1, newVector, nil, nil, newPayload)

// Delete a point.
col.Delete(1)

// Get total point count.
count := col.Count()`,
    desc: "Update is implemented as an atomic delete + insert. Count returns the total across all segments (mutable + immutable).",
  },
  {
    title: "Snapshots",
    code: `// Create a snapshot.
snap, err := col.Snapshot("pre-upgrade")
if err != nil { panic(err) }
fmt.Printf("Checksum: %x\\n", snap.Checksum)`,
    desc: "Snapshots create point-in-time copies of all data with BLAKE3 integrity verification.",
  },
  {
    title: "Managing Collections",
    code: `// Create a new collection in a workspace.
err := api.CreateCollection("./workspace", "products", cfg)

// List all collections.
names, err := api.ListCollections("./workspace")

// Drop a collection.
err := api.DropCollection("./workspace", "products")`,
    desc: "Collection management functions operate on the filesystem directly, creating or removing collection directories.",
  },
  {
    title: "Dense Engine (Direct Usage)",
    code: `import "github.com/hasdev/oris/engine/dense"
import "github.com/hasdev/oris/engine/dense/factory"
import "github.com/hasdev/oris/engine/dense/flat"
import "github.com/hasdev/oris/engine/dense/hnsw"

// Auto-select: Flat < 1000 vectors, HNSW otherwise.
cfg := dense.Config{Dimension: 128, Distance: "cosine", M: 16, EfConstruction: 200, EfSearch: 100}
engine := factory.New(cfg, 5000) // HNSW for 5000 expected vectors

// Flat explicitly.
engine = factory.NewFlat(cfg)

// HNSW explicitly.
engine = factory.NewHNSW(cfg)

// Direct HNSW usage.
h := hnsw.New(cfg)
h.Insert(1, vector)
results, _ := h.Search(query, 10)`,
    desc: "The dense engine interface supports HNSW and Flat backends. The factory auto-selects based on expected dataset size.",
  },
  {
    title: "Sparse Engine (Direct Usage)",
    code: `import "github.com/hasdev/oris/engine/sparse/bm25"

// Create BM25 engine (k1=1.2, b=0.75).
engine := bm25.New(1.2, 0.75)

// Index documents.
engine.IndexDocument(1, bm25.Tokenize("the quick brown fox"))
engine.IndexDocument(2, bm25.Tokenize("jumps over the lazy dog"))

// Search.
results, _ := engine.Search(bm25.Tokenize("brown fox"), 10)`,
    desc: "BM25 provides keyword-based retrieval. Text is tokenized (lowercased, split on punctuation) before indexing.",
  },
  {
    title: "Metadata Engine (Direct Usage)",
    code: `import "github.com/hasdev/oris/engine/metadata"

meta := metadata.New()
meta.IndexField(1, "country", "japan")
meta.IndexField(1, "price", float64(100))

// Evaluate country = japan AND price > 50.
min := 50.0
result, _ := meta.Evaluate(&metadata.And{
    Filters: []metadata.Filter{
        &metadata.In{Field: "country", Values: []string{"japan"}},
        &metadata.Range{Field: "price", Min: &min, Max: nil},
    },
})
// result is a RoaringBitmap of matching doc IDs.`,
    desc: "The metadata engine maintains field indexes and evaluates filter expressions to produce Roaring Bitmaps.",
  },
  {
    title: "HTTP API Server",
    code: `import (
    "net/http"
    "github.com/hasdev/oris/api"
)

col, _ := api.Open("./data", cfg)
server := api.NewHTTPServer(col)
http.ListenAndServe(":8080", server)

// POST /insert  {"id":1, "dense":[0.1,...], "payload":"..."}
// POST /search  {"query":[0.1,...], "topK":10}
// GET  /search?query=0.1,0.2,...&topK=10
// POST /delete?id=1
// GET  /count
// POST /snapshot`,
    desc: "The HTTP API wraps the Collection with RESTful endpoints for JSON-based interaction.",
  },
]

export default function ApiReferencePage() {
  return (
    <div className="space-y-8 max-w-3xl">
      <div>
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">API Reference</h1>
        <p className="text-sm text-zinc-500 mt-1">Complete API documentation with code examples.</p>
      </div>

      {sections.map((section) => (
        <section key={section.title} className="space-y-3">
          <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">{section.title}</h2>
          <p className="text-sm text-zinc-600 dark:text-zinc-400">{section.desc}</p>
          <Card>
            <CardContent className="p-4">
              <pre className="text-sm font-mono bg-zinc-900 text-zinc-100 p-4 rounded-md overflow-x-auto leading-relaxed">
                {section.code}
              </pre>
            </CardContent>
          </Card>
        </section>
      ))}
    </div>
  )
}
