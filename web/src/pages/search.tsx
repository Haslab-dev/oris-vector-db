import { useState } from "react"
import { Button } from "../components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "../components/ui/card"
import { Badge } from "../components/ui/badge"
import { Search as SearchIcon } from "lucide-react"

interface SearchResult {
  id: string
  score: number
  preview: string
  metadata?: Record<string, string>
}

export default function SearchPage() {
  const [query, setQuery] = useState("")
  const [results, setResults] = useState<SearchResult[] | null>(null)
  const [selected, setSelected] = useState<SearchResult | null>(null)
  const [loading, setLoading] = useState(false)

  async function handleSearch() {
    if (!query.trim()) return
    setLoading(true)
    // Simulate API call.
    await new Promise((r) => setTimeout(r, 300))
    setResults([
      { id: "doc_123", score: 0.962, preview: "HNSW is a graph-based approximate nearest neighbor search algorithm...", metadata: { category: "AI", language: "EN" } },
      { id: "doc_456", score: 0.941, preview: "Hierarchical Navigable Small World graphs provide efficient ANN search...", metadata: { category: "AI", language: "EN" } },
      { id: "doc_789", score: 0.887, preview: "The HNSW algorithm builds a multi-layer graph structure for vector search...", metadata: { category: "ML", language: "EN" } },
    ])
    setLoading(false)
    setSelected(null)
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Search</h1>
        <p className="text-sm text-zinc-500 mt-1">Explore your vector collections</p>
      </div>

      <div className="flex gap-2">
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && handleSearch()}
          placeholder="Search your collections..."
          className="flex-1 h-10 px-3 rounded-md border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-sm focus:outline-none focus:ring-1 focus:ring-zinc-400"
        />
        <Button onClick={handleSearch} disabled={loading}>
          <SearchIcon className="h-4 w-4 mr-1" />
          {loading ? "Searching..." : "Search"}
        </Button>
      </div>

      {results && (
        <div className="flex gap-4">
          <div className="flex-1 space-y-2">
            {results.map((r, i) => (
              <Card
                key={r.id}
                className={`cursor-pointer transition-colors ${selected?.id === r.id ? "ring-2 ring-zinc-400 dark:ring-zinc-500" : ""}`}
                onClick={() => setSelected(r)}
              >
                <CardContent className="p-4">
                  <div className="flex items-center justify-between mb-1">
                    <span className="text-sm text-zinc-500">#{i + 1}</span>
                    <Badge>{(r.score * 100).toFixed(0)}%</Badge>
                  </div>
                  <p className="text-sm font-medium text-zinc-900 dark:text-zinc-50">{r.preview}</p>
                </CardContent>
              </Card>
            ))}
          </div>

          {selected && (
            <div className="w-72 shrink-0">
              <Card>
                <CardHeader><CardTitle>Details</CardTitle></CardHeader>
                <CardContent className="space-y-3">
                  <div>
                    <p className="text-xs text-zinc-500">ID</p>
                    <p className="text-sm font-mono text-zinc-900 dark:text-zinc-50">{selected.id}</p>
                  </div>
                  <div>
                    <p className="text-xs text-zinc-500">Score</p>
                    <p className="text-sm font-mono text-zinc-900 dark:text-zinc-50">{selected.score.toFixed(3)}</p>
                  </div>
                  {selected.metadata && (
                    <div>
                      <p className="text-xs text-zinc-500 mb-1">Metadata</p>
                      {Object.entries(selected.metadata).map(([k, v]) => (
                        <div key={k} className="flex justify-between text-sm">
                          <span className="text-zinc-500">{k}</span>
                          <span className="text-zinc-900 dark:text-zinc-50">{v}</span>
                        </div>
                      ))}
                    </div>
                  )}
                </CardContent>
              </Card>
            </div>
          )}
        </div>
      )}
    </div>
  )
}
