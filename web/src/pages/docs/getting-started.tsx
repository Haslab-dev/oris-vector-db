import { Card, CardContent, CardHeader, CardTitle } from "../../components/ui/card"


export default function GettingStartedPage() {
  return (
    <div className="space-y-8 max-w-3xl">
      <div>
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Getting Started</h1>
        <p className="text-sm text-zinc-500 mt-1">Install, import, and run your first vector search in minutes.</p>
      </div>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Installation</h2>
        <Card>
          <CardContent className="p-4">
            <pre className="text-sm font-mono bg-zinc-900 text-zinc-100 p-4 rounded-md overflow-x-auto">
{`go get github.com/hasdev/oris`}
            </pre>
          </CardContent>
        </Card>
        <p className="text-sm text-zinc-600 dark:text-zinc-400">
          Oris is a pure Go library with zero external runtime dependencies. The only system dependency is Pebble for storage, which is also pure Go.
        </p>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Your First Collection</h2>
        <Card>
          <CardContent className="p-4">
            <pre className="text-sm font-mono bg-zinc-900 text-zinc-100 p-4 rounded-md overflow-x-auto">
{`package main

import (
    "fmt"
    "github.com/hasdev/oris/api"
)

func main() {
    // Create a collection with 128-dimensional cosine vectors.
    cfg := api.DefaultConfig("my-collection", 128)
    col, err := api.Open("./data", cfg)
    if err != nil {
        panic(err)
    }
    defer col.Close()

    // Insert 100 points.
    for i := uint64(0); i < 100; i++ {
        vec := make([]float32, 128)
        vec[i%128] = float32(i) / 100.0
        col.Insert(i, vec, nil, nil, []byte(fmt.Sprintf("point-%d", i)))
    }

    fmt.Printf("Inserted %d points\\n", col.Count())

    // Search with a query vector.
    query := make([]float32, 128)
    query[0] = 0.5
    results, _ := col.Search(query, 5)
    for _, r := range results {
        fmt.Printf("ID=%d Score=%.4f\\n", r.ID, r.Score)
    }
}`}
            </pre>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Core Concepts</h2>
        <div className="grid grid-cols-1 gap-3">
          {[
            { title: "Collection", desc: "A named container of points with a fixed vector dimension and distance metric. All points in a collection share the same configuration." },
            { title: "Point", desc: "A single entry in a collection consisting of an ID, a dense vector, an optional sparse vector, metadata key-value pairs, and a payload blob." },
            { title: "Segment", desc: "The unit of storage and indexing. Segments start mutable (accepting writes), then seal to immutable (read-only) on disk. Background compaction merges small segments into larger ones." },
            { title: "Engine", desc: "Each retrieval subsystem (dense, sparse, metadata) is an independently replaceable engine behind a stable Go interface." },
            { title: "Planner", desc: "The query planner orchestrates dense search, sparse search, metadata filtering, result merging, and ranking." },
          ].map((c) => (
            <Card key={c.title}>
              <CardHeader className="py-3">
                <CardTitle className="text-sm">{c.title}</CardTitle>
              </CardHeader>
              <CardContent className="py-0 pb-3">
                <p className="text-sm text-zinc-600 dark:text-zinc-400">{c.desc}</p>
              </CardContent>
            </Card>
          ))}
        </div>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Performance Targets</h2>
        <Card>
          <CardContent className="p-4">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-zinc-200 dark:border-zinc-700">
                  <th className="text-left pb-2 font-medium text-zinc-500">Metric</th>
                  <th className="text-right pb-2 font-medium text-zinc-500">Target</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-zinc-100 dark:divide-zinc-800">
                {[
                  ["Insert latency", "&lt; 1 ms per point"],
                  ["Search latency (100K vectors)", "&lt; 10 ms"],
                  ["Search latency (1M vectors)", "&lt; 30 ms"],
                  ["Recall@10", "≥ 0.95"],
                  ["Startup time (small collections)", "&lt; 100 ms"],
                ].map(([metric, target]) => (
                  <tr key={metric}>
                    <td className="py-2 text-zinc-900 dark:text-zinc-50">{metric}</td>
                    <td className="py-2 text-right font-mono text-xs text-zinc-600 dark:text-zinc-400">{target}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Next Steps</h2>
        <div className="flex flex-wrap gap-2">
          {[
            { label: "Architecture Overview", path: "/docs/architecture" },
            { label: "HNSW Dense Engine", path: "/docs/dense-engine" },
            { label: "BM25 Sparse Engine", path: "/docs/sparse-engine" },
            { label: "API Reference", path: "/docs/api-reference" },
          ].map((link) => (
            <a
              key={link.label}
              href={`#${link.path}`}
              className="text-sm text-blue-600 dark:text-blue-400 hover:underline"
            >
              {link.label} →
            </a>
          ))}
        </div>
      </section>
    </div>
  )
}
