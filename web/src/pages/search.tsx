import { useState, useEffect } from "react"
import { Button } from "../components/ui/button"
import { Card, CardContent } from "../components/ui/card"
import { Badge } from "../components/ui/badge"
import { Search as SearchIcon } from "lucide-react"
import { searchCollection, fetchCollections, type SearchResult, type CollectionInfo } from "../lib/api"

export default function SearchPage() {
  const [collections, setCollections] = useState<CollectionInfo[]>([])
  const [selectedCol, setSelectedCol] = useState("")
  const [query, setQuery] = useState("")
  const [topK, setTopK] = useState(10)
  const [results, setResults] = useState<SearchResult[] | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchCollections().then(setCollections).catch(console.error)
  }, [])

  async function handleSearch() {
    if (!selectedCol || !query.trim()) return
    setLoading(true)
    setError(null)
    try {
      const parts = query.split(",").map(Number)
      const data = await searchCollection(selectedCol, parts, topK)
      setResults(data)
    } catch (e: any) {
      setError(e.message)
      setResults(null)
    }
    setLoading(false)
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Search</h1>
        <p className="text-sm text-zinc-500 mt-1">Query vectors in your collections</p>
      </div>

      <div className="flex flex-wrap gap-2 items-end">
        <div className="flex-1 min-w-[200px]">
          <label className="text-xs font-medium text-zinc-500 block mb-1">Collection</label>
          <select
            value={selectedCol}
            onChange={(e) => setSelectedCol(e.target.value)}
            className="w-full h-10 px-3 rounded-md border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-sm"
          >
            <option value="">Select collection...</option>
            {collections.map((c) => (
              <option key={c.name} value={c.name}>{c.name}</option>
            ))}
          </select>
        </div>
        <div className="flex-[2] min-w-[300px]">
          <label className="text-xs font-medium text-zinc-500 block mb-1">Query Vector (comma-separated)</label>
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="0.1, 0.2, 0.3, ..."
            className="w-full h-10 px-3 rounded-md border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-sm font-mono focus:outline-none focus:ring-1 focus:ring-zinc-400"
          />
        </div>
        <div className="w-20">
          <label className="text-xs font-medium text-zinc-500 block mb-1">Top K</label>
          <input
            type="number"
            value={topK}
            onChange={(e) => setTopK(Number(e.target.value))}
            min={1}
            max={100}
            className="w-full h-10 px-3 rounded-md border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-sm text-center"
          />
        </div>
        <Button onClick={handleSearch} disabled={loading || !selectedCol || !query.trim()}>
          <SearchIcon className="h-4 w-4 mr-1" />
          {loading ? "Searching..." : "Search"}
        </Button>
      </div>

      {error && <p className="text-sm text-red-500">{error}</p>}

      {results && (
        <div className="space-y-2">
          <p className="text-xs text-zinc-500">{results.length} result{results.length !== 1 ? "s" : ""}</p>
          {results.map((r, i) => (
            <Card key={r.id}>
              <CardContent className="p-4 flex items-center justify-between">
                <div className="flex items-center gap-4">
                  <span className="text-sm text-zinc-400 font-mono">#{i + 1}</span>
                  <span className="text-sm font-mono text-zinc-900 dark:text-zinc-50">ID: {r.id}</span>
                </div>
                <Badge>{(r.score * 100).toFixed(1)}%</Badge>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  )
}
