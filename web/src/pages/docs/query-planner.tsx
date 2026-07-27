import { Card, CardContent } from "../../components/ui/card"

export default function QueryPlannerPage() {
  return (
    <div className="space-y-8 max-w-3xl">
      <div>
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Query Planner</h1>
        <p className="text-sm text-zinc-500 mt-1">Orchestrating dense, sparse, and hybrid retrieval with metadata filtering.</p>
      </div>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Overview</h2>
        <p className="text-sm text-zinc-600 dark:text-zinc-400">
          The Planner is the brain of Oris. It receives a query with optional dense vector, sparse tokens, metadata filter,
          and execution mode, then routes to the appropriate engines and merges results into a unified ranked list.
        </p>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Query Modes</h2>
        <Card>
          <CardContent className="p-0">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-zinc-200 dark:border-zinc-700">
                  <th className="text-left p-3 font-medium text-zinc-500">Mode</th>
                  <th className="text-left p-3 font-medium text-zinc-500">Description</th>
                  <th className="text-left p-3 font-medium text-zinc-500">Use Case</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-zinc-100 dark:divide-zinc-800">
                {[
                  ["DenseOnly", "Semantic search via HNSW/Flat", "Find similar items by vector embedding similarity"],
                  ["SparseOnly", "Lexical search via BM25", "Keyword search with precise term matching"],
                  ["Hybrid", "Combined dense + sparse with alpha weighting", "Best of both worlds — semantic understanding + keyword precision"],
                ].map(([mode, desc, use]) => (
                  <tr key={mode}>
                    <td className="p-3 font-mono text-xs text-zinc-900 dark:text-zinc-50">{mode}</td>
                    <td className="p-3 text-xs text-zinc-600 dark:text-zinc-400">{desc}</td>
                    <td className="p-3 text-xs text-zinc-600 dark:text-zinc-400">{use}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Query Parameters</h2>
        <Card>
          <CardContent className="p-4">
            <pre className="text-sm font-mono bg-zinc-900 text-zinc-100 p-4 rounded-md overflow-x-auto">
{`type Query struct {
    DenseVector   []float32         // Dense query vector
    SparseTokens  []string          // Sparse query tokens
    Filter        metadata.Filter   // Metadata filter (optional)
    TopK          int               // Number of results to return
    Mode          Mode              // DenseOnly, SparseOnly, or Hybrid
    Alpha         float64           // Hybrid weight: 0=sparse, 1=dense
}`}
            </pre>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Hybrid Search Algorithm</h2>
        <Card>
          <CardContent className="p-4">
            <pre className="text-xs font-mono bg-zinc-900 text-zinc-100 p-4 rounded-md leading-relaxed">
{`Hybrid Search (alpha=0.5):
  1. Run dense search: query the Engine with topK*3
  2. Run sparse search: query BM25 with topK*3
  3. Normalize each result set to [0, 1] range:
     score = (score - min) / (max - min)
  4. Merge scores:
     finalScore = alpha * denseScore + (1-alpha) * sparseScore
  5. Sort by finalScore ascending
  6. Optionally apply metadata filter
  7. Return topK results

Score normalization is critical for fair hybrid merging
since dense and sparse scores have different ranges.`}
            </pre>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Metadata Filter Integration</h2>
        <Card>
          <CardContent className="p-4 space-y-3">
            <p className="text-sm text-zinc-600 dark:text-zinc-400">
              The Planner supports metadata filters in two ways:
            </p>
            <ul className="list-disc list-inside text-sm text-zinc-600 dark:text-zinc-400 space-y-1">
              <li><strong className="text-zinc-900 dark:text-zinc-50">Post-filter</strong>: After search results are produced, filter out results not matching the metadata filter. This is the current implementation.</li>
              <li><strong className="text-zinc-900 dark:text-zinc-50">Pre-filter (planned)</strong>: Evaluate the metadata filter first to get candidate IDs, then only search within those candidates. More efficient when the filter is highly selective.</li>
            </ul>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Usage Example</h2>
        <Card>
          <CardContent className="p-4">
            <pre className="text-sm font-mono bg-zinc-900 text-zinc-100 p-4 rounded-md overflow-x-auto">
{`import "github.com/hasdev/oris/engine/planner"

pl := planner.New(denseEngine, sparseEngine, metaEngine, 128)

// Dense-only search.
results, _ := pl.Execute(planner.Query{
    DenseVector: queryVec,
    TopK: 10,
    Mode: planner.DenseOnly,
})

// Hybrid search with metadata filter.
tagA := &metadata.In{Field: "tag", Values: []string{"a"}}
results, _ := pl.Execute(planner.Query{
    DenseVector:  queryVec,
    SparseTokens: bm25.Tokenize("hello world"),
    TopK:         10,
    Mode:         planner.Hybrid,
    Alpha:        0.5,
    Filter:       tagA,
})`}
            </pre>
          </CardContent>
        </Card>
      </section>
    </div>
  )
}
