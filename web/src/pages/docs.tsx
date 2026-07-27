import { Card, CardContent, CardHeader, CardTitle } from "../components/ui/card"

const sections = [
  {
    title: "Architecture Overview",
    content: `Oris is organized into five core engines: Dense (HNSW), Sparse (BM25), Metadata (RoaringBitmap),
Query Planner, and Storage (Pebble). Each engine is independently replaceable behind a stable interface,
enabling future backends like DiskANN, SPLADE, or custom storage without architectural rewrites.

The query flow: User Query → Planner → Metadata Filter → Dense Search → Sparse Search → Merge → Rank → Results.`,
  },
  {
    title: "Segment Lifecycle",
    content: `Segments are the unit of storage and indexing. Each segment owns a dense index, sparse index,
metadata index, and payload store.

Lifecycle: Mutable (accepting writes) → Seal → Immutable (read-only, on disk) → Compaction → Merge → Delete.

The SegmentManager coordinates this: writes go to the active mutable segment. When it reaches capacity,
it seals to disk as an immutable segment. When too many immutables accumulate, compaction merges them
into a single segment. This is inspired by Lucene, RocksDB, and Qdrant.`,
  },
  {
    title: "HNSW (Dense Engine)",
    content: `Hierarchical Navigable Small World graphs provide approximate nearest neighbor search.
Oris implements a custom HNSW with:

- Multi-layer graph with random level generation
- Greedy search from top level down to layer 0
- searchLayer with min-heap candidates and max-heap results
- Edge connection and shrinking when exceeding Mmax
- Dynamic delete with edge cleanup and entry point re-selection
- Parameters: M (connections), efConstruction (build quality), efSearch (search quality)

A Flat (brute-force) backend auto-selects for small collections (< 1000 vectors).`,
  },
  {
    title: "BM25 (Sparse Engine)",
    content: `Best Matching 25 for lexical retrieval. Configurable k1 (term saturation) and b (length normalization).

- Inverted index with posting lists (doc ID → term frequency)
- TF-IDF scoring with Robertson IDF formula
- Document length normalization
- Tokenizer: lowercase, unicode-aware, splits on punctuation/whitespace
- Top-K search sorted by BM25 score descending`,
  },
  {
    title: "Metadata Engine",
    content: `Metadata filtering uses Roaring Bitmaps for efficient set operations.

Supported filters: AND (intersection), OR (union), NOT (negation), IN (multi-value),
Range (numeric bounds), Exists (field presence).

Three indexes: string field index (value → bitmap), numeric index (sorted entries for range queries),
exists index (field name → bitmap).`,
  },
  {
    title: "Storage Layout",
    content: `workspace/
  collections/
    docs/
    wal/
    snapshots/
    segments/
      segment_NNN/
        dense.idx
        sparse.idx
        metadata.idx
        payload.dat

Storage interface supports Pebble (default) and in-memory (testing).
WAL provides crash recovery by replaying uncommitted entries on open.
Snapshots use BLAKE3 checksums for integrity verification.`,
  },
  {
    title: "Query Flow",
    content: `1. User submits Query (dense vector, sparse tokens, metadata filter, topK)
2. Planner routes based on mode: DenseOnly, SparseOnly, or Hybrid
3. Dense search: segment manager fans out across mutable + immutable segments
4. Sparse search: BM25 scores each matching document
5. Hybrid: scores are normalized to [0,1], then combined with alpha weighting
6. Metadata filter: applied as post-filter to refine results
7. Final results sorted by score (ascending distance)`,
  },
]

export default function DocsPage() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Documentation</h1>
        <p className="text-sm text-zinc-500 mt-1">Oris architecture and engine internals</p>
      </div>

      <div className="grid grid-cols-1 gap-4">
        {sections.map((section) => (
          <Card key={section.title}>
            <CardHeader>
              <CardTitle className="text-base">{section.title}</CardTitle>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-zinc-600 dark:text-zinc-400 leading-relaxed whitespace-pre-line font-mono text-xs">
                {section.content}
              </p>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}
