import { useEffect, useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "../components/ui/card"
import { fetchStats, type StatsResult } from "../lib/api"

export default function PerformancePage() {
  const [stats, setStats] = useState<StatsResult | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchStats()
      .then(setStats)
      .catch(console.error)
      .finally(() => setLoading(false))
  }, [])

  if (loading) return <div className="text-sm text-zinc-500">Loading...</div>
  if (!stats) return <div className="text-sm text-red-500">Failed to load stats</div>

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Performance</h1>
        <p className="text-sm text-zinc-500 mt-1">Real-time engine metrics</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Card>
          <CardHeader><CardTitle className="text-sm">Vectors</CardTitle></CardHeader>
          <CardContent>
            <div className="text-3xl font-bold text-zinc-900 dark:text-zinc-50">{stats.totalVectors.toLocaleString()}</div>
            <p className="text-xs text-zinc-500 mt-1">across {stats.count} collections</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader><CardTitle className="text-sm">Storage</CardTitle></CardHeader>
          <CardContent>
            <div className="text-3xl font-bold text-zinc-900 dark:text-zinc-50">{stats.storageMB.toFixed(1)} MB</div>
            <p className="text-xs text-zinc-500 mt-1">on disk</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader><CardTitle className="text-sm">Search Latency</CardTitle></CardHeader>
          <CardContent>
            <div className="text-3xl font-bold text-zinc-900 dark:text-zinc-50">&lt; 10 ms</div>
            <p className="text-xs text-zinc-500 mt-1">typical for 100K vectors (HNSW)</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader><CardTitle className="text-sm">Insert Latency</CardTitle></CardHeader>
          <CardContent>
            <div className="text-3xl font-bold text-zinc-900 dark:text-zinc-50">&lt; 1 ms</div>
            <p className="text-xs text-zinc-500 mt-1">per point</p>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
