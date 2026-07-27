import { Card, CardContent } from "../../components/ui/card"

export default function SparseEnginePage() {
  return (
    <div className="space-y-8 max-w-3xl">
      <div>
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Sparse Engine (BM25)</h1>
        <p className="text-sm text-zinc-500 mt-1">Best Matching 25 for lexical retrieval.</p>
      </div>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">BM25 Scoring Formula</h2>
        <Card>
          <CardContent className="p-4">
            <div className="bg-zinc-100 dark:bg-zinc-800 p-4 rounded-md text-center">
              <span className="text-sm font-mono text-zinc-900 dark:text-zinc-50">
                Score(D, Q) = Σ IDF(t) · TF(t, D) · (k₁ + 1) / (TF(t, D) + k₁ · (1 - b + b · |D| / avgdl))
              </span>
            </div>
            <div className="mt-4 space-y-2 text-sm text-zinc-600 dark:text-zinc-400">
              <p><strong className="text-zinc-900 dark:text-zinc-50">TF(t, D)</strong> — Term frequency: how many times term t appears in document D</p>
              <p><strong className="text-zinc-900 dark:text-zinc-50">IDF(t)</strong> — Inverse document frequency: log(1 + (N - n + 0.5) / (n + 0.5)) where N is total documents and n is the number containing t</p>
              <p><strong className="text-zinc-900 dark:text-zinc-50">k₁</strong> — Term saturation parameter (default: 1.2). Higher values increase the impact of term frequency.</p>
              <p><strong className="text-zinc-900 dark:text-zinc-50">b</strong> — Length normalization parameter (default: 0.75). b = 0 disables normalization, b = 1 applies full normalization.</p>
              <p><strong className="text-zinc-900 dark:text-zinc-50">|D| / avgdl</strong> — Document length relative to the average. Shorter documents get a boost, longer documents are penalized.</p>
            </div>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Inverted Index</h2>
        <Card>
          <CardContent className="p-4 space-y-3">
            <p className="text-sm text-zinc-600 dark:text-zinc-400">
              BM25 maintains an inverted index mapping each unique term to a posting list of (docID, termFrequency) pairs.
              When a document is indexed, all its tokens are extracted via the tokenizer, and the inverted index is updated.
            </p>
            <pre className="text-xs font-mono bg-zinc-900 text-zinc-100 p-4 rounded-md overflow-x-auto leading-relaxed">
{`Inverted Index (conceptual):
"hello" → [(doc1, 2), (doc3, 1), (doc5, 1)]
"world" → [(doc1, 1), (doc2, 1)]
"oris"  → [(doc4, 3)]

Search for "hello world":
  1. Tokenize query → ["hello", "world"]
  2. Look up posting lists
  3. Score each matching document using BM25 formula
  4. Sort by score descending
  5. Return top-K`}
            </pre>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Tokenizer</h2>
        <Card>
          <CardContent className="p-4">
            <p className="text-sm text-zinc-600 dark:text-zinc-400 mb-3">
              The default tokenizer lowercases text and splits on whitespace and punctuation boundaries.
              Only alphanumeric characters are kept; everything else is treated as a delimiter.
            </p>
            <pre className="text-xs font-mono bg-zinc-900 text-zinc-100 p-4 rounded-md overflow-x-auto">
{`Input:  "Hello, World! Go123 is #1."
Output: ["hello", "world", "go123", "is", "1"]

Input:  "The quick brown fox"
Output: ["the", "quick", "brown", "fox"]`}
            </pre>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">Usage</h2>
        <Card>
          <CardContent className="p-4">
            <pre className="text-sm font-mono bg-zinc-900 text-zinc-100 p-4 rounded-md overflow-x-auto">
{`import "github.com/hasdev/oris/engine/sparse/bm25"

engine := bm25.New(1.2, 0.75)

// Index documents.
engine.IndexDocument(1, bm25.Tokenize("the quick brown fox"))
engine.IndexDocument(2, bm25.Tokenize("jumps over the lazy dog"))

// Search.
results, _ := engine.Search(bm25.Tokenize("brown fox"), 10)
for _, r := range results {
    fmt.Printf("Doc %d score: %.4f\\n", r.ID, r.Score)
}`}
            </pre>
          </CardContent>
        </Card>
      </section>

      <section className="space-y-3">
        <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50">BM25 vs Dense Search</h2>
        <Card>
          <CardContent className="p-4">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-zinc-200 dark:border-zinc-700">
                  <th className="text-left pb-2 font-medium text-zinc-500">Aspect</th>
                  <th className="text-left pb-2 font-medium text-zinc-500">BM25 (Sparse)</th>
                  <th className="text-left pb-2 font-medium text-zinc-500">HNSW (Dense)</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-zinc-100 dark:divide-zinc-800">
                {[
                  ["Match type", "Exact keyword match", "Semantic similarity"],
                  ["Query", "Text tokens", "Dense vector (embedding)"],
                  ["Index", "Inverted index (term → docs)", "Graph (vector → neighbors)"],
                  ["Strength", "Precision, term frequency", "Recall, semantic understanding"],
                  ["Weakness", "Vocabulary mismatch", "Requires good embeddings"],
                ].map(([aspect, sparse, dense]) => (
                  <tr key={aspect}>
                    <td className="py-2 text-xs text-zinc-600 dark:text-zinc-400">{aspect}</td>
                    <td className="py-2 text-xs text-zinc-900 dark:text-zinc-50">{sparse}</td>
                    <td className="py-2 text-xs text-zinc-900 dark:text-zinc-50">{dense}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </CardContent>
        </Card>
      </section>
    </div>
  )
}
