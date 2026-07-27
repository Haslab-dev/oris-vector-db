const API_BASE = "/api"

export interface CollectionInfo {
  name: string
  vectors: number
  dimension: number
  distance: string
  segments: number
  sizeBytes: number
}

export interface StatsResult {
  count: number
  totalVectors: number
  storageMB: number
  collections: CollectionInfo[]
}

export interface SearchResult {
  id: number
  score: number
}

export async function fetchStats(): Promise<StatsResult> {
  const res = await fetch(`${API_BASE}/stats`)
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

export async function fetchCollections(): Promise<CollectionInfo[]> {
  const res = await fetch(`${API_BASE}/collections`)
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

export async function fetchCollection(name: string): Promise<CollectionInfo> {
  const res = await fetch(`${API_BASE}/collections/${encodeURIComponent(name)}`)
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

export async function searchCollection(collection: string, query: number[], topK: number): Promise<SearchResult[]> {
  const res = await fetch(`${API_BASE}/search`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ collection, query, topK }),
  })
  if (!res.ok) throw new Error(await res.text())
  return res.json()
}

export async function insertPoint(collection: string, id: number, dense: number[], payload?: string): Promise<void> {
  const res = await fetch(`${API_BASE}/insert`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ collection, id, dense, payload }),
  })
  if (!res.ok) throw new Error(await res.text())
}

export async function deletePoint(collection: string, id: number): Promise<void> {
  const res = await fetch(`${API_BASE}/delete`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ collection, id }),
  })
  if (!res.ok) throw new Error(await res.text())
}

export async function fetchCount(collection: string): Promise<number> {
  const res = await fetch(`${API_BASE}/count?collection=${encodeURIComponent(collection)}`)
  if (!res.ok) throw new Error(await res.text())
  const data = await res.json()
  return data.count
}
