import { useEffect, useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "../components/ui/card"
import { Database, HardDrive, MemoryStick } from "lucide-react"
import { fetchStats, type StatsResult } from "../lib/api"

export default function DashboardPage() {
  const [stats, setStats] = useState<StatsResult | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    fetchStats()
      .then(setStats)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false))
  }, [])

  if (loading) return <div className="text-sm text-zinc-500">Loading...</div>
  if (error) return <div className="text-sm text-red-500">Error: {error}</div>
  if (!stats) return null

  const statCards = [
    { icon: Database, label: "Collections", value: String(stats.count), change: `${stats.count} total` },
    { icon: Database, label: "Vectors", value: stats.totalVectors.toLocaleString(), change: "across all collections" },
    { icon: HardDrive, label: "Storage", value: `${stats.storageMB.toFixed(1)} MB`, change: "on disk" },
    { icon: MemoryStick, label: "Memory", value: "Embedded", change: "in-process" },
  ]

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">Dashboard</h1>
        <p className="text-sm text-zinc-500 mt-1">Real-time overview of your Oris instance</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {statCards.map((s) => (
          <Card key={s.label}>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium text-zinc-500">{s.label}</CardTitle>
              <s.icon className="h-4 w-4 text-zinc-400" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">{s.value}</div>
              <p className="text-xs text-zinc-500 mt-1">{s.change}</p>
            </CardContent>
          </Card>
        ))}
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Collections</CardTitle>
        </CardHeader>
        <CardContent>
          {stats.collections.length === 0 ? (
            <p className="text-sm text-zinc-500">No collections yet. Create one to get started.</p>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-zinc-200 dark:border-zinc-700 text-left">
                    <th className="pb-3 font-medium text-zinc-500">Name</th>
                    <th className="pb-3 font-medium text-zinc-500 text-right">Vectors</th>
                    <th className="pb-3 font-medium text-zinc-500 text-right">Dimension</th>
                    <th className="pb-3 font-medium text-zinc-500">Distance</th>
                  </tr>
                </thead>
                <tbody>
                  {stats.collections.map((col) => (
                    <tr key={col.name} className="border-b border-zinc-100 dark:border-zinc-800">
                      <td className="py-3 font-medium text-zinc-900 dark:text-zinc-50">{col.name}</td>
                      <td className="py-3 text-right text-zinc-600 dark:text-zinc-400">{col.vectors.toLocaleString()}</td>
                      <td className="py-3 text-right text-zinc-600 dark:text-zinc-400">{col.dimension}</td>
                      <td className="py-3 text-zinc-600 dark:text-zinc-400">{col.distance}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
