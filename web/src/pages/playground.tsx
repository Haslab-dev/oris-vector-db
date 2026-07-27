import { useState, useEffect } from "react"
import { Button } from "../components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "../components/ui/card"
import { Badge } from "../components/ui/badge"
import { Trash2, Plus, Search, Database, HardDrive } from "lucide-react"
import {
  fetchCollections,
  fetchStats,
  searchCollection,
  insertPoint,
  deletePoint,
  fetchCount,
  type CollectionInfo,
  type StatsResult,
  type SearchResult,
} from "../lib/api"

export default function PlaygroundPage() {
  const [stats, setStats] = useState<StatsResult | null>(null)
  const [collections, setCollections] = useState<CollectionInfo[]>([])
  const [selectedCol, setSelectedCol] = useState("")
  const [count, setCount] = useState(0)

  // Insert
  const [insertId, setInsertId] = useState("1")
  const [insertDim, setInsertDim] = useState("4")
  const [insertPayload, setInsertPayload] = useState("")
  const [insertMsg, setInsertMsg] = useState("")
  const [numToInsert, setNumToInsert] = useState("100")

  // Search
  const [query, setQuery] = useState("0.5, 0.5, 0, 0")
  const [topK, setTopK] = useState("10")
  const [searchResults, setSearchResults] = useState<SearchResult[] | null>(null)
  const [searchMsg, setSearchMsg] = useState("")

  // Delete
  const [deleteId, setDeleteId] = useState("")
  const [deleteMsg, setDeleteMsg] = useState("")

  // Create collection
  const [newColName, setNewColName] = useState("playground")
  const [newColDim, setNewColDim] = useState("4")
  const [createMsg, setCreateMsg] = useState("")

  const [loading, setLoading] = useState(false)

  useEffect(() => {
    refresh()
  }, [])

  async function refresh() {
    const [s, c] = await Promise.all([
      fetchStats().catch(() => null),
      fetchCollections().catch(() => []),
    ])
    setStats(s)
    setCollections(c)
  }

  async function refreshCount(col: string) {
    if (!col) return
    const c = await fetchCount(col).catch(() => 0)
    setCount(c)
  }

  async function handleSelectCol(name: string) {
    setSelectedCol(name)
    setSearchResults(null)
    setSearchMsg("")
    setInsertMsg("")
    setDeleteMsg("")
    if (name) {
      const c = await fetchCount(name).catch(() => 0)
      setCount(c)
    } else {
      setCount(0)
    }
  }

  async function handleCreate() {
    setCreateMsg("")
    setLoading(true)
    try {
      const res = await fetch("/api/collections", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name: newColName, dimension: parseInt(newColDim) || 4 }),
      })
      if (!res.ok) {
        const err = await res.json()
        setCreateMsg(`Error: ${err.error}`)
        return
      }
      setCreateMsg(`Collection "${newColName}" created`)
      await refresh()
      await handleSelectCol(newColName)
    } catch (e: any) {
      setCreateMsg(`Error: ${e.message}`)
    }
    setLoading(false)
  }

  async function handleInsertOne() {
    setInsertMsg("")
    const dim = parseInt(insertDim)
    const id = parseInt(insertId)
    if (isNaN(id) || isNaN(dim)) {
      setInsertMsg("Invalid ID or dimension")
      return
    }
    const vec = new Array(dim).fill(0)
    vec[0] = (id % 100) / 100
    try {
      await insertPoint(selectedCol, id, vec, insertPayload || undefined)
      setInsertMsg(`Inserted ID ${id}`)
      await refreshCount(selectedCol)
    } catch (e: any) {
      setInsertMsg(`Error: ${e.message}`)
    }
  }

  async function handleInsertMany() {
    setInsertMsg("")
    const dim = parseInt(insertDim)
    const n = parseInt(numToInsert)
    if (isNaN(n) || isNaN(dim) || n < 1 || n > 100000) {
      setInsertMsg("Invalid count (1-100000)")
      return
    }
    setLoading(true)
    let inserted = 0
    for (let i = 0; i < n; i++) {
      const vec = new Array(dim).fill(0)
      vec[i % dim] = (i % 100) / 100
      try {
        await insertPoint(selectedCol, i + 1, vec, `point-${i + 1}`)
        inserted++
      } catch {
        break
      }
    }
    setInsertMsg(`Inserted ${inserted} points`)
    await refreshCount(selectedCol)
    setLoading(false)
  }

  async function handleSearch() {
    setSearchMsg("")
    setSearchResults(null)
    const k = parseInt(topK)
    const parts = query.split(",").map((s) => parseFloat(s.trim())).filter((n) => !isNaN(n))
    if (parts.length === 0 || isNaN(k)) {
      setSearchMsg("Invalid query or topK")
      return
    }
    setLoading(true)
    try {
      const results = await searchCollection(selectedCol, parts, k)
      setSearchResults(results)
      if (results.length === 0) setSearchMsg("No results")
    } catch (e: any) {
      setSearchMsg(`Error: ${e.message}`)
    }
    setLoading(false)
  }

  async function handleDelete() {
    setDeleteMsg("")
    const id = parseInt(deleteId)
    if (isNaN(id)) {
      setDeleteMsg("Invalid ID")
      return
    }
    try {
      await deletePoint(selectedCol, id)
      setDeleteMsg(`Deleted ID ${id}`)
      await refreshCount(selectedCol)
    } catch (e: any) {
      setDeleteMsg(`Error: ${e.message}`)
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Playground</h1>
        <Badge variant="green">Live</Badge>
      </div>

      {/* Quick Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-3">
        <Card>
          <CardContent className="p-3 flex items-center gap-3">
            <Database className="h-4 w-4 text-zinc-400" />
            <div>
              <p className="text-xs text-zinc-500">Collections</p>
              <p className="text-sm font-bold">{stats?.count ?? "—"}</p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-3 flex items-center gap-3">
            <Database className="h-4 w-4 text-zinc-400" />
            <div>
              <p className="text-xs text-zinc-500">Total Vectors</p>
              <p className="text-sm font-bold">{stats?.totalVectors.toLocaleString() ?? "—"}</p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-3 flex items-center gap-3">
            <HardDrive className="h-4 w-4 text-zinc-400" />
            <div>
              <p className="text-xs text-zinc-500">Storage</p>
              <p className="text-sm font-bold">{stats ? `${stats.storageMB.toFixed(1)} MB` : "—"}</p>
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="p-3 flex items-center gap-3">
            <Database className="h-4 w-4 text-zinc-400" />
            <div>
              <p className="text-xs text-zinc-500">Selected</p>
              <p className="text-sm font-bold">{selectedCol || "—"}</p>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Collection Selector + Create */}
      <div className="flex flex-wrap gap-2 items-end">
        <div className="flex-1 min-w-[200px]">
          <label className="text-xs font-medium text-zinc-500 block mb-1">Active Collection</label>
          <select
            value={selectedCol}
            onChange={(e) => handleSelectCol(e.target.value)}
            className="w-full h-10 px-3 rounded-md border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-sm"
          >
            <option value="">Select...</option>
            {collections.map((c) => (
              <option key={c.name} value={c.name}>{c.name} ({c.vectors.toLocaleString()} vec)</option>
            ))}
          </select>
        </div>
        <div className="w-24">
          <label className="text-xs font-medium text-zinc-500 block mb-1">New Name</label>
          <input
            value={newColName}
            onChange={(e) => setNewColName(e.target.value)}
            className="w-full h-10 px-3 rounded-md border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-sm"
          />
        </div>
        <div className="w-20">
          <label className="text-xs font-medium text-zinc-500 block mb-1">Dim</label>
          <input
            value={newColDim}
            onChange={(e) => setNewColDim(e.target.value)}
            className="w-full h-10 px-3 rounded-md border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-sm text-center"
          />
        </div>
        <Button onClick={handleCreate} disabled={loading || !newColName}>
          <Plus className="h-4 w-4 mr-1" /> Create
        </Button>
        <Button variant="outline" onClick={refresh}>
          ↻ Refresh
        </Button>
      </div>
      {createMsg && <p className={`text-xs ${createMsg.startsWith("Error") ? "text-red-500" : "text-emerald-600"}`}>{createMsg}</p>}

      {selectedCol && (
        <>
          <p className="text-xs text-zinc-400">Collection: <strong>{selectedCol}</strong> · {count.toLocaleString()} vectors</p>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            {/* Insert */}
            <Card>
              <CardHeader className="pb-2"><CardTitle className="text-sm flex items-center gap-2"><Plus className="h-4 w-4" /> Insert</CardTitle></CardHeader>
              <CardContent className="space-y-2">
                <div className="flex gap-2">
                  <input
                    placeholder="ID"
                    value={insertId}
                    onChange={(e) => setInsertId(e.target.value)}
                    className="w-20 h-9 px-2 rounded-md border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-sm"
                  />
                  <input
                    placeholder="Dim"
                    value={insertDim}
                    onChange={(e) => setInsertDim(e.target.value)}
                    className="w-16 h-9 px-2 rounded-md border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-sm"
                  />
                  <input
                    placeholder="Payload (opt)"
                    value={insertPayload}
                    onChange={(e) => setInsertPayload(e.target.value)}
                    className="flex-1 h-9 px-2 rounded-md border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-sm"
                  />
                  <Button size="sm" onClick={handleInsertOne}>Add</Button>
                </div>
                <div className="flex gap-2 items-center">
                  <span className="text-xs text-zinc-500">Bulk:</span>
                  <input
                    value={numToInsert}
                    onChange={(e) => setNumToInsert(e.target.value)}
                    className="w-20 h-8 px-2 rounded-md border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-sm"
                  />
                  <Button size="sm" variant="outline" onClick={handleInsertMany} disabled={loading}>
                    Insert N
                  </Button>
                </div>
                {insertMsg && <p className={`text-xs ${insertMsg.startsWith("Error") ? "text-red-500" : "text-emerald-600"}`}>{insertMsg}</p>}
              </CardContent>
            </Card>

            {/* Search */}
            <Card>
              <CardHeader className="pb-2"><CardTitle className="text-sm flex items-center gap-2"><Search className="h-4 w-4" /> Search</CardTitle></CardHeader>
              <CardContent className="space-y-2">
                <input
                  value={query}
                  onChange={(e) => setQuery(e.target.value)}
                  placeholder="0.5, 0.5, 0, 0"
                  className="w-full h-9 px-2 rounded-md border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-sm font-mono"
                />
                <div className="flex gap-2">
                  <input
                    value={topK}
                    onChange={(e) => setTopK(e.target.value)}
                    className="w-16 h-9 px-2 rounded-md border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-sm"
                  />
                  <Button size="sm" onClick={handleSearch} disabled={loading}><Search className="h-3 w-3 mr-1" /> Search</Button>
                </div>
                {searchMsg && <p className="text-xs text-zinc-500">{searchMsg}</p>}
                {searchResults && searchResults.length > 0 && (
                  <div className="max-h-48 overflow-y-auto space-y-1">
                    {searchResults.map((r, i) => (
                      <div key={r.id} className="flex justify-between items-center bg-zinc-50 dark:bg-zinc-800/50 px-2 py-1.5 rounded text-xs">
                        <span className="font-mono text-zinc-900 dark:text-zinc-50">#{i + 1} · ID {r.id}</span>
                        <Badge>{(r.score * 100).toFixed(1)}%</Badge>
                      </div>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>

            {/* Delete */}
            <Card>
              <CardHeader className="pb-2"><CardTitle className="text-sm flex items-center gap-2"><Trash2 className="h-4 w-4" /> Delete</CardTitle></CardHeader>
              <CardContent className="space-y-2">
                <div className="flex gap-2">
                  <input
                    value={deleteId}
                    onChange={(e) => setDeleteId(e.target.value)}
                    placeholder="Point ID"
                    className="flex-1 h-9 px-2 rounded-md border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-sm"
                  />
                  <Button size="sm" variant="outline" onClick={handleDelete}><Trash2 className="h-3 w-3 mr-1" /> Delete</Button>
                </div>
                {deleteMsg && <p className={`text-xs ${deleteMsg.startsWith("Error") ? "text-red-500" : "text-emerald-600"}`}>{deleteMsg}</p>}
              </CardContent>
            </Card>
          </div>
        </>
      )}
    </div>
  )
}
